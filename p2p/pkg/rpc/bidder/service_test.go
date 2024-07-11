package bidderapi_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	bidderapi "github.com/primev/mev-commit/p2p/pkg/rpc/bidder"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	bufferSize      = 101024 * 1024
	blocksPerWindow = 64
)

type bid struct {
	txHex    string
	amount   string
	blockNum int64
}

type testSender struct {
	bids         []bid
	noOfPreconfs int
}

func (s *testSender) SendBid(
	ctx context.Context,
	txHex string,
	amount string,
	blockNum int64,
	decayStartTimestamp int64,
	decayEndTimestamp int64,
	revertedTxns string,
) (chan *preconfpb.PreConfirmation, error) {
	s.bids = append(s.bids, bid{
		txHex:    txHex,
		amount:   amount,
		blockNum: blockNum,
	})

	preconfs := make(chan *preconfpb.PreConfirmation, s.noOfPreconfs)
	for i := 0; i < s.noOfPreconfs; i++ {
		preconfs <- &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				TxHash:              txHex,
				BidAmount:           amount,
				BlockNumber:         blockNum,
				DecayStartTimestamp: decayStartTimestamp,
				DecayEndTimestamp:   decayEndTimestamp,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				RevertingTxHashes:   revertedTxns,
			},
			Digest:          []byte("digest"),
			Signature:       []byte("signature"),
			ProviderAddress: common.HexToAddress(fmt.Sprintf("%x", i)).Bytes(),
		}
	}

	close(preconfs)

	return preconfs, nil
}

type testRegistryContract struct {
	deposit *big.Int
}

func (t *testRegistryContract) DepositForWindow(opts *bind.TransactOpts, _ *big.Int) (*types.Transaction, error) {
	t.deposit = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) DepositForWindows(opts *bind.TransactOpts, _ []*big.Int) (*types.Transaction, error) {
	t.deposit = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) WithdrawBidderAmountFromWindow(
	opts *bind.TransactOpts,
	address common.Address,
	window *big.Int,
) (*types.Transaction, error) {
	return types.NewTransaction(2, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) GetDeposit(_ *bind.CallOpts, _ common.Address, _ *big.Int) (*big.Int, error) {
	return t.deposit, nil
}

func (t *testRegistryContract) ParseBidderRegistered(_ types.Log) (*bidderregistry.BidderregistryBidderRegistered, error) {
	return &bidderregistry.BidderregistryBidderRegistered{
		DepositedAmount: t.deposit,
		WindowNumber:    big.NewInt(1),
	}, nil
}

func (t *testRegistryContract) ParseBidderWithdrawal(_ types.Log) (*bidderregistry.BidderregistryBidderWithdrawal, error) {
	return &bidderregistry.BidderregistryBidderWithdrawal{
		Amount: t.deposit,
		Window: big.NewInt(1),
	}, nil
}

func (t *testRegistryContract) WithdrawFromWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

type testAutoDepositTracker struct {
	mtx       sync.Mutex
	deposits  map[uint64]bool
	isWorking bool
}

func (t *testAutoDepositTracker) Start(ctx context.Context, startWindow, amount *big.Int) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.isWorking = true
	t.deposits[startWindow.Uint64()] = true
	t.deposits[big.NewInt(0).Add(startWindow, big.NewInt(1)).Uint64()] = true
	return nil
}

func (t *testAutoDepositTracker) IsWorking() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.isWorking
}

func (t *testAutoDepositTracker) GetStatus() (map[uint64]bool, bool, *big.Int) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.deposits, t.isWorking, big.NewInt(1)
}

func (t *testAutoDepositTracker) Stop() ([]*big.Int, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.isWorking = false
	var windowNumbers []*big.Int
	for k := range t.deposits {
		windowNumbers = append(windowNumbers, big.NewInt(int64(k)))
		delete(t.deposits, k)
	}
	return windowNumbers, nil
}

type testTxWatcher struct {
	nonce int
}

func (t *testTxWatcher) WaitForReceipt(_ context.Context, tx *types.Transaction) (*types.Receipt, error) {
	t.nonce++
	if tx.Nonce() != uint64(t.nonce) {
		return nil, errors.New("nonce mismatch")
	}
	return &types.Receipt{
		Status: 1,
		Logs: []*types.Log{
			{
				Address: common.Address{},
				Topics:  []common.Hash{},
				Data:    []byte{},
			},
		},
	}, nil
}

type testBlockTrackerContract struct {
	blockNumberToWinner map[uint64]common.Address
	lastBlockNumber     uint64
	blocksPerWindow     uint64
}

func (btc *testBlockTrackerContract) GetCurrentWindow() (*big.Int, error) {
	return big.NewInt(int64(btc.lastBlockNumber / btc.blocksPerWindow)), nil
}

func startServer(t *testing.T) bidderapiv1.BidderClient {
	lis := bufconn.Listen(bufferSize)

	logger := util.NewTestLogger(os.Stdout)
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("error creating validator: %v", err)
	}

	owner := common.HexToAddress("0x00001")
	registryContract := &testRegistryContract{
		deposit: big.NewInt(1000000000000000000),
	}
	sender := &testSender{noOfPreconfs: 2}
	blockTrackerContract := &testBlockTrackerContract{lastBlockNumber: blocksPerWindow + 1, blocksPerWindow: blocksPerWindow, blockNumberToWinner: make(map[uint64]common.Address)}
	testAutoDepositTracker := &testAutoDepositTracker{deposits: make(map[uint64]bool)}
	srvImpl := bidderapi.NewService(
		owner,
		blockTrackerContract.blocksPerWindow,
		sender,
		registryContract,
		blockTrackerContract,
		validator,
		&testTxWatcher{},
		func(ctx context.Context) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{
				From:    owner,
				Context: ctx,
			}, nil
		},
		testAutoDepositTracker,
		logger,
	)

	baseServer := grpc.NewServer()
	bidderapiv1.RegisterBidderServer(baseServer, srvImpl)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			// Ignore "use of closed network connection" error
			if opErr, ok := err.(*net.OpError); !ok || !errors.Is(opErr.Err, net.ErrClosed) {
				t.Logf("server stopped err: %v", err)
			}
		}
	}()

	// nolint:staticcheck
	conn, err := grpc.DialContext(context.TODO(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("error connecting to server: %v", err)
	}

	t.Cleanup(func() {
		err := lis.Close()
		if err != nil {
			t.Errorf("error closing listener: %v", err)
		}
		baseServer.Stop()
	})

	client := bidderapiv1.NewBidderClient(conn)

	return client
}

func TestDepositHandling(t *testing.T) {
	t.Parallel()

	client := startServer(t)

	t.Run("deposit", func(t *testing.T) {
		type testCase struct {
			amount string
			err    string
		}

		for _, tc := range []testCase{
			{
				amount: "",
				err:    "amount must be a valid integer",
			},
			{
				amount: "0000000000000000000",
				err:    "amount must be a valid integer",
			},
			{
				amount: "asdf",
				err:    "amount must be a valid integer",
			},
			{
				amount: "1000000000000000000",
				err:    "",
			},
		} {
			deposit, err := client.Deposit(context.Background(), &bidderapiv1.DepositRequest{Amount: tc.amount})
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error depositing")
				}
			} else {
				if err != nil {
					t.Fatalf("error depositing: %v", err)
				}
				if deposit.Amount != tc.amount {
					t.Fatalf("expected amount to be %v, got %v", tc.amount, deposit.Amount)
				}
			}
		}
	})

	t.Run("get deposit", func(t *testing.T) {
		deposit, err := client.GetDeposit(context.Background(), &bidderapiv1.GetDepositRequest{WindowNumber: wrapperspb.UInt64(1)})
		if err != nil {
			t.Fatalf("error getting deposit: %v", err)
		}
		if deposit.Amount != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", deposit.Amount)
		}
	})

	t.Run("withdraw", func(t *testing.T) {
		resp, err := client.Withdraw(context.Background(), &bidderapiv1.WithdrawRequest{WindowNumber: wrapperspb.UInt64(1)})
		if err != nil {
			t.Fatalf("error withdrawing: %v", err)
		}

		if resp.Amount != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", resp.Amount)
		}

		if resp.WindowNumber.Value != 1 {
			t.Fatalf("expected window number to be 1, got %v", resp.WindowNumber)
		}
	})
}

func TestAutoDepositHandling(t *testing.T) {
	t.Parallel()

	client := startServer(t)

	t.Run("autodeposit", func(t *testing.T) {
		deposit, err := client.AutoDeposit(context.Background(), &bidderapiv1.DepositRequest{
			Amount:       "1000000000000000000",
			WindowNumber: wrapperspb.UInt64(1),
		})
		if err != nil {
			t.Fatalf("error depositing: %v", err)
		}
		if deposit.StartBlockNumber.Value != 1 {
			t.Fatalf("expected start block number to be 1, got %v", deposit.StartBlockNumber)
		}
		if deposit.AmountPerWindow != "1000000000000000000" {
			t.Fatalf("expected amount per window to be 1000000000000000000, got %v", deposit.AmountPerWindow)
		}
	})

	t.Run("get status", func(t *testing.T) {
		status, err := client.AutoDepositStatus(context.Background(), &bidderapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting deposit: %v", err)
		}
		if status.IsWorking != true {
			t.Fatalf("expected is working to be true, got %v", status.IsWorking)
		}
		if len(status.WindowBalances) != 2 {
			t.Fatalf("expected 2 deposits, got %v", len(status.WindowBalances))
		}
		for _, v := range status.WindowBalances {
			if v.WindowNumber.Value != 1 && v.WindowNumber.Value != 2 {
				t.Fatalf("unexpected window number, got %v", v.WindowNumber)
			}
			if v.Amount != "1000000000000000000" {
				t.Fatalf("expected amount to be 1000000000000000000, got %v", v)
			}
			if v.WindowNumber.Value == 1 && v.StartBlockNumber.Value != 1 && v.EndBlockNumber.Value != blocksPerWindow && !v.IsCurrent {
				t.Fatalf("expected correct values for window 1, got %v", v)
			}
			if v.WindowNumber.Value == 2 && v.StartBlockNumber.Value != blocksPerWindow+1 && v.EndBlockNumber.Value != blocksPerWindow*2 && v.IsCurrent {
				t.Fatalf("expected correct values for window 2, got %v", v)
			}
		}
	})

	t.Run("stop autodeposit", func(t *testing.T) {
		resp, err := client.CancelAutoDeposit(context.Background(), &bidderapiv1.CancelAutoDepositRequest{
			Withdraw: true,
		})
		if err != nil {
			t.Fatalf("error stopping autodeposit: %v", err)
		}
		if len(resp.WindowNumbers) != 0 {
			t.Fatalf("expected 0 window numbers, got %v", len(resp.WindowNumbers))
		}
	})

	t.Run("stop no withdraw", func(t *testing.T) {
		_, err := client.AutoDeposit(context.Background(), &bidderapiv1.DepositRequest{
			WindowNumber: wrapperspb.UInt64(5),
			Amount:       "1000000000000000000",
		})
		if err != nil {
			t.Fatalf("error getting deposit: %v", err)
		}

		resp, err := client.CancelAutoDeposit(context.Background(), &bidderapiv1.CancelAutoDepositRequest{})
		if err != nil {
			t.Fatalf("error stopping autodeposit: %v", err)
		}
		if len(resp.WindowNumbers) != 2 {
			t.Fatalf("expected 2 window numbers, got %v", len(resp.WindowNumbers))
		}

		windows := make([]*wrapperspb.UInt64Value, 2)
		copy(windows, resp.WindowNumbers)

		_, err = client.WithdrawFromWindows(context.Background(), &bidderapiv1.WithdrawFromWindowsRequest{
			WindowNumbers: windows,
		})
		if err != nil {
			t.Fatalf("error withdrawing: %v", err)
		}
	})
}

func TestSendBid(t *testing.T) {
	t.Parallel()

	client := startServer(t)

	type testCase struct {
		name                string
		txHexs              []string
		amount              string
		blockNum            int64
		decayStartTimestamp int64
		decayEndTimestamp   int64
		err                 string
	}

	for _, tc := range []testCase{
		{
			name:                "invalid tx hex",
			txHexs:              []string{"asdf"},
			amount:              "1000000000000000000",
			blockNum:            1,
			decayStartTimestamp: 10,
			decayEndTimestamp:   20,
			err:                 "tx_hashes must be a valid array of transaction hashes",
		},
		{
			name:                "no txns",
			txHexs:              nil,
			amount:              "1000000000000000000",
			blockNum:            1,
			decayStartTimestamp: 10,
			decayEndTimestamp:   20,
			err:                 "tx_hashes must be a valid array of transaction hashes",
		},
		{
			name:                "invalid amount",
			txHexs:              []string{common.HexToHash("0x0000ab").Hex()[2:]},
			amount:              "000000000000000000",
			blockNum:            1,
			decayStartTimestamp: 10,
			decayEndTimestamp:   20,
			err:                 "amount must be a valid integer",
		},
		{
			name:                "invalid block number",
			txHexs:              []string{common.HexToHash("0x0000ab").Hex()[2:]},
			amount:              "1000000000000000000",
			blockNum:            0,
			decayStartTimestamp: 10,
			decayEndTimestamp:   20,
			err:                 "block_number must be a valid integer",
		},
		{
			name:                "success",
			txHexs:              []string{common.HexToHash("0x0000ab").Hex()[2:]},
			amount:              "1000000000000000000",
			blockNum:            1,
			decayStartTimestamp: 10,
			decayEndTimestamp:   20,
			err:                 "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rcv, err := client.SendBid(context.Background(), &bidderapiv1.Bid{
				TxHashes:            tc.txHexs,
				Amount:              tc.amount,
				BlockNumber:         tc.blockNum,
				DecayStartTimestamp: tc.decayStartTimestamp,
				DecayEndTimestamp:   tc.decayEndTimestamp,
				RevertingTxHashes:   []string{},
			})
			if err != nil {
				t.Fatalf("error sending bid: %v", err)
			}

			if tc.err != "" {
				_, err := rcv.Recv()
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error sending bid %s, got %v", tc.err, err)
				}
			} else {
				count := 0
				for {
					_, err := rcv.Recv()
					if err != nil {
						if errors.Is(err, io.EOF) {
							break
						}
						t.Fatalf("error receiving preconfs: %v", err)
					}
					count++
				}
				if count != 2 {
					t.Fatalf("expected 2 preconfs, got %v", count)
				}
			}
		})
	}
}
