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
	bufferSize = 101024 * 1024
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
	deposit    *big.Int
	minDeposit *big.Int
}

func (t *testRegistryContract) DepositForSpecificWindow(opts *bind.TransactOpts, _ *big.Int) (*types.Transaction, error) {
	t.deposit = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) DepositForNWindows(opts *bind.TransactOpts, _ *big.Int, _ uint16) (*types.Transaction, error) {
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

func (t *testRegistryContract) MinDeposit(_ *bind.CallOpts) (*big.Int, error) {
	return t.minDeposit, nil
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

type testAutoDepositTracker struct {
	deposits map[uint64]bool
}

func (t *testAutoDepositTracker) DoAutoMoveToAnotherWindow(ads []*bidderapiv1.AutoDeposit) <-chan struct{} {
	for _, ad := range ads {
		t.deposits[ad.WindowNumber.Value] = true
	}

	doneChan := make(chan struct{})
	close(doneChan)

	return doneChan
}

func (t *testAutoDepositTracker) Stop() {
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
	registryContract := &testRegistryContract{minDeposit: big.NewInt(100000000000000000)}
	sender := &testSender{noOfPreconfs: 2}
	blockTrackerContract := &testBlockTrackerContract{blocksPerWindow: 64, blockNumberToWinner: make(map[uint64]common.Address)}
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

	t.Run("get min deposit", func(t *testing.T) {
		deposit, err := client.GetMinDeposit(context.Background(), &bidderapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting min deposit: %v", err)
		}
		if deposit.Amount != "100000000000000000" {
			t.Fatalf("expected amount to be 100000000000000000, got %v", deposit.Amount)
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
