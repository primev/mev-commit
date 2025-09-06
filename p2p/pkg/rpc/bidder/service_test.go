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
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	preconfstore "github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	bidderapi "github.com/primev/mev-commit/p2p/pkg/rpc/bidder"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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
	b *preconfpb.Bid,
) (chan *preconfpb.PreConfirmation, error) {
	s.bids = append(s.bids, bid{
		txHex:    b.TxHash,
		amount:   b.BidAmount,
		blockNum: b.BlockNumber,
	})

	preconfs := make(chan *preconfpb.PreConfirmation, s.noOfPreconfs)
	for i := 0; i < s.noOfPreconfs; i++ {
		preconfs <- &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				TxHash:              b.TxHash,
				BidAmount:           b.BidAmount,
				SlashAmount:         b.SlashAmount,
				BlockNumber:         b.BlockNumber,
				DecayStartTimestamp: b.DecayStartTimestamp,
				DecayEndTimestamp:   b.DecayEndTimestamp,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				RevertingTxHashes:   b.RevertingTxHashes,
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
	deposit   *big.Int
	perAmount *big.Int
	watcher   *testTxWatcher
}

func (t *testRegistryContract) DepositAsBidder(opts *bind.TransactOpts, _ common.Address) (*types.Transaction, error) {
	if opts != nil && opts.Value != nil {
		t.perAmount = new(big.Int).Set(opts.Value)
	}
	if t.watcher != nil {
		t.watcher.logs = 1
	}
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) DepositEvenlyAsBidder(opts *bind.TransactOpts, providers []common.Address) (*types.Transaction, error) {
	if opts != nil && opts.Value != nil && len(providers) > 0 {
		t.perAmount = new(big.Int).Div(new(big.Int).Set(opts.Value), big.NewInt(int64(len(providers))))
	}
	if t.watcher != nil {
		t.watcher.logs = len(providers)
	}
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) RequestWithdrawalsAsBidder(_ *bind.TransactOpts, providers []common.Address) (*types.Transaction, error) {
	if t.watcher != nil {
		t.watcher.logs = len(providers)
	}
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) WithdrawAsBidder(_ *bind.TransactOpts, providers []common.Address) (*types.Transaction, error) {
	if t.watcher != nil {
		t.watcher.logs = len(providers)
	}
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) GetDeposit(_ *bind.CallOpts, _ common.Address, _ common.Address) (*big.Int, error) {
	return t.deposit, nil
}

func (t *testRegistryContract) ParseBidderDeposited(_ types.Log) (*bidderregistry.BidderregistryBidderDeposited, error) {
	amt := t.perAmount
	if amt == nil {
		amt = t.deposit
	}
	return &bidderregistry.BidderregistryBidderDeposited{
		DepositedAmount: amt,
		Bidder:          common.Address{},
		Provider:        common.Address{},
	}, nil
}

func (t *testRegistryContract) ParseWithdrawalRequested(_ types.Log) (*bidderregistry.BidderregistryWithdrawalRequested, error) {
	amt := t.perAmount
	if amt == nil {
		amt = t.deposit
	}
	return &bidderregistry.BidderregistryWithdrawalRequested{
		AvailableAmount: amt,
		Provider:        common.Address{},
	}, nil
}

func (t *testRegistryContract) ParseBidderWithdrawal(_ types.Log) (*bidderregistry.BidderregistryBidderWithdrawal, error) {
	amt := t.perAmount
	if amt == nil {
		amt = t.deposit
	}
	return &bidderregistry.BidderregistryBidderWithdrawal{
		AmountWithdrawn: amt,
		Provider:        common.Address{},
	}, nil
}

func (t *testRegistryContract) FilterBidderDeposited(_ *bind.FilterOpts, _ []common.Address, _ []common.Address, _ []*big.Int) (*bidderregistry.BidderregistryBidderDepositedIterator, error) {
	return nil, nil
}

type testTxWatcher struct {
	logs int
}

func (t *testTxWatcher) WaitForReceipt(_ context.Context, tx *types.Transaction) (*types.Receipt, error) {
	n := t.logs
	if n <= 0 {
		n = 1
	}
	logs := make([]*types.Log, n)
	for i := 0; i < n; i++ {
		logs[i] = &types.Log{
			Address: common.Address{},
			Topics:  []common.Hash{},
			Data:    []byte{},
		}
	}
	return &types.Receipt{
		Status: 1,
		Logs:   logs,
	}, nil
}

type testStore struct {
	commitments []*preconfstore.Commitment
}

func (t *testStore) ListCommitments(_ *preconfstore.ListOpts) ([]*preconfstore.Commitment, error) {
	return t.commitments, nil
}

type testProviderRegistry struct {
	claim *big.Int
}

func (t *testProviderRegistry) WithdrawSlashedAmount(_ *bind.TransactOpts) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testProviderRegistry) BidderSlashedAmount(_ *bind.CallOpts, _ common.Address) (*big.Int, error) {
	return t.claim, nil
}

func (t *testProviderRegistry) ParseBidderWithdrawSlashedAmount(_log types.Log) (*providerregistry.ProviderregistryBidderWithdrawSlashedAmount, error) {
	return &providerregistry.ProviderregistryBidderWithdrawSlashedAmount{
		Amount: t.claim,
	}, nil
}

func (t *testProviderRegistry) FilterProviderRegistered(_ *bind.FilterOpts, _ []common.Address) (*providerregistry.ProviderregistryProviderRegisteredIterator, error) {
	return nil, nil
}

func (t *testProviderRegistry) AreProvidersValid(_ *bind.CallOpts, _ []common.Address) ([]bool, error) {
	return []bool{true}, nil
}

func (t *testStore) GetCommitments(blockNum int64) ([]*preconfstore.Commitment, error) {
	cmts := make([]*preconfstore.Commitment, 0)
	for _, c := range t.commitments {
		if c.Bid.BlockNumber == blockNum {
			cmts = append(cmts, c)
		}
	}
	return cmts, nil
}

type testSetCodeHelper struct {
}

func (t *testSetCodeHelper) SetCode(ctx context.Context, opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func startServer(t *testing.T) bidderapiv1.BidderClient {
	cs := &testStore{}
	return startServerWithStore(t, cs)
}

func startServerWithStore(t *testing.T, cs *testStore) bidderapiv1.BidderClient {
	lis := bufconn.Listen(bufferSize)

	logger := util.NewTestLogger(os.Stdout)
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("error creating validator: %v", err)
	}

	owner := common.HexToAddress("0x00001")
	watcher := &testTxWatcher{}
	registryContract := &testRegistryContract{
		deposit: big.NewInt(1000000000000000000),
		watcher: watcher,
	}
	providerRegistry := &testProviderRegistry{
		claim: big.NewInt(1000000000000000000),
	}
	sender := &testSender{noOfPreconfs: 2}
	setCodeHelper := &testSetCodeHelper{}
	srvImpl := bidderapi.NewService(
		owner,
		sender,
		registryContract,
		providerRegistry,
		validator,
		watcher,
		func(ctx context.Context) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{
				From:    owner,
				Context: ctx,
			}, nil
		},
		cs,
		15*time.Second,
		logger,
		setCodeHelper,
		nil,
		nil,
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
	)

	baseServer := grpc.NewServer()
	bidderapiv1.RegisterBidderServer(baseServer, srvImpl)
	srvStopped := make(chan struct{})
	go func() {
		defer close(srvStopped)

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

		<-srvStopped
	})

	client := bidderapiv1.NewBidderClient(conn)

	return client
}

func TestDepositHandling(t *testing.T) {
	t.Parallel()

	client := startServer(t)

	t.Run("deposit", func(t *testing.T) {
		type testCase struct {
			amount   string
			provider string
			err      string
		}

		for _, tc := range []testCase{
			{
				amount:   "",
				provider: "0x9323e8B917ED544c26039A4D0Cccd8E3083e5647",
				err:      "amount must be a valid integer",
			},
			{
				amount:   "0000000000000000000",
				provider: "0x9323e8B917ED544c26039A4D0Cccd8E3083e5647",
				err:      "amount must be a valid integer",
			},
			{
				amount:   "asdf",
				provider: "0x9323e8B917ED544c26039A4D0Cccd8E3083e5647",
				err:      "amount must be a valid integer",
			},
			{
				amount:   "1000000000000000000",
				provider: "0x9323e8B917ED544c26039A4D0Cccd8E3083e5647",
				err:      "",
			},
			{
				amount:   "1000000000000000000",
				provider: "string",
				err:      "provider must be a valid Ethereum address",
			},
		} {
			deposit, err := client.Deposit(context.Background(), &bidderapiv1.DepositRequest{
				Amount:   tc.amount,
				Provider: tc.provider,
			})
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
		deposit, err := client.GetDeposit(context.Background(), &bidderapiv1.GetDepositRequest{
			Provider: "0x353ff5537cc2fe78C8632c41Bb848735C2982c45",
		})
		if err != nil {
			t.Fatalf("error getting deposit: %v", err)
		}
		if deposit.Amount != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", deposit.Amount)
		}
	})

	t.Run("deposit evenly", func(t *testing.T) {
		deposit, err := client.DepositEvenly(context.Background(), &bidderapiv1.DepositEvenlyRequest{
			TotalAmount: "1000000000000000000",
			Providers:   []string{"0x8F644015Aa59984BF849259375d9135A89a940e7", "0x353ff5537cc2fe78C8632c41Bb848735C2982c45"},
		})
		if err != nil {
			t.Fatalf("error depositing evenly: %v", err)
		}
		if len(deposit.Providers) != 2 {
			t.Fatalf("expected 2 providers, got %v", len(deposit.Providers))
		}
		if deposit.Amounts[0] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", deposit.Amounts[0])
		}
		if deposit.Amounts[1] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", deposit.Amounts[1])
		}
	})

	t.Run("request withdrawals", func(t *testing.T) {
		resp, err := client.RequestWithdrawals(context.Background(), &bidderapiv1.RequestWithdrawalsRequest{
			Providers: []string{"0x8F644015Aa59984BF849259375d9135A89a940e7", "0x353ff5537cc2fe78C8632c41Bb848735C2982c45"},
		})
		if err != nil {
			t.Fatalf("error requesting withdrawals: %v", err)
		}
		if len(resp.Providers) != 2 {
			t.Fatalf("expected 2 providers, got %v", len(resp.Providers))
		}
		if len(resp.Amounts) != 2 {
			t.Fatalf("expected 2 amounts, got %v", len(resp.Amounts))
		}
		if resp.Amounts[0] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", resp.Amounts[0])
		}
		if resp.Amounts[1] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", resp.Amounts[1])
		}
	})

	t.Run("withdraw", func(t *testing.T) {
		resp, err := client.Withdraw(context.Background(), &bidderapiv1.WithdrawRequest{
			Providers: []string{"0x8F644015Aa59984BF849259375d9135A89a940e7", "0x353ff5537cc2fe78C8632c41Bb848735C2982c45"},
		})
		if err != nil {
			t.Fatalf("error withdrawing: %v", err)
		}

		if len(resp.Providers) != 2 {
			t.Fatalf("expected 2 providers, got %v", len(resp.Providers))
		}
		if len(resp.Amounts) != 2 {
			t.Fatalf("expected 2 amounts, got %v", len(resp.Amounts))
		}
		if resp.Amounts[0] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", resp.Amounts[0])
		}
		if resp.Amounts[1] != "500000000000000000" {
			t.Fatalf("expected amount to be 500000000000000000, got %v", resp.Amounts[1])
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
			err:                 "empty bid",
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
		{
			name:                "success",
			txHexs:              []string{common.HexToHash("0x0000ab").Hex()},
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

func TestClaimSlashedFunds(t *testing.T) {
	t.Parallel()

	client := startServer(t)

	t.Run("claim slashed funds", func(t *testing.T) {
		resp, err := client.ClaimSlashedFunds(context.Background(), &bidderapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error claiming slashed funds: %v", err)
		}
		if resp.Value != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", resp.Value)
		}
	})
}

func TestGetBidInfo(t *testing.T) {
	t.Parallel()

	t.Run("get bid info", func(t *testing.T) {
		testCommitment := &preconfstore.Commitment{
			PreConfirmation: &preconfpb.PreConfirmation{
				Bid: &preconfpb.Bid{
					TxHash:              "0x1234567890abcdef,0x1234567890abcdef",
					BlockNumber:         1,
					BidAmount:           "1000000000000000000",
					DecayStartTimestamp: 123456789,
					DecayEndTimestamp:   123457896,
					Digest:              []byte("digest"),
				},
				DispatchTimestamp: 123456889,
				ProviderAddress:   common.HexToAddress("0x1234").Bytes(),
			},
			Status:  preconfstore.CommitmentStatusOpened,
			Details: "test details",
			Payment: "900000000000000000",
			Refund:  "100000000000000000",
		}

		store := &testStore{
			commitments: []*preconfstore.Commitment{testCommitment},
		}

		client := startServerWithStore(t, store)

		resp, err := client.GetBidInfo(context.Background(), &bidderapiv1.GetBidInfoRequest{
			BlockNumber: 1,
		})
		if err != nil {
			t.Fatalf("error getting bid info: %v", err)
		}
		if len(resp.BlockBidInfo) != 1 {
			t.Fatalf("expected 1 bid info, got %v", len(resp.BlockBidInfo))
		}
		bidInfo := resp.BlockBidInfo[0]
		if bidInfo.BlockNumber != 1 {
			t.Fatalf("expected block number to be 1, got %v", bidInfo.BlockNumber)
		}
		if len(bidInfo.Bids) != 1 {
			t.Fatalf("expected 1 bid, got %v", len(bidInfo.Bids))
		}
		bid := bidInfo.Bids[0]
		if len(bid.TxnHashes) != 2 {
			t.Fatalf("expected 2 transaction hashes, got %v", len(bid.TxnHashes))
		}
		if bid.TxnHashes[0] != "0x1234567890abcdef" || bid.TxnHashes[1] != "0x1234567890abcdef" {
			t.Fatalf("expected transaction hashes to be 0x1234567890abcdef,0x1234567890abcdef, got %v", bid.TxnHashes)
		}
		if bid.BidAmount != "1000000000000000000" {
			t.Fatalf("expected bid amount to be 1000000000000000000, got %v", bid.BidAmount)
		}
		if bid.DecayStartTimestamp != 123456789 {
			t.Fatalf("expected decay start timestamp to be 123456789, got %v", bid.DecayStartTimestamp)
		}
		if bid.DecayEndTimestamp != 123457896 {
			t.Fatalf("expected decay end timestamp to be 123457896, got %v", bid.DecayEndTimestamp)
		}
		if bid.BidDigest != common.Bytes2Hex([]byte("digest")) {
			t.Fatalf("expected bid digest to be 'digest', got %v", bid.BidDigest)
		}
		if len(bid.Commitments) != 1 {
			t.Fatalf("expected 1 commitment, got %v", len(bid.Commitments))
		}
		commitment := bid.Commitments[0]
		if commitment.Details != "test details" {
			t.Fatalf("expected commitment details to be 'test details', got %v", commitment.Details)
		}
		if commitment.Status != string(preconfstore.CommitmentStatusOpened) {
			t.Fatalf("expected commitment status to be opened, got %v", commitment.Status)
		}
		if commitment.Payment != "900000000000000000" {
			t.Fatalf("expected commitment payment to be 900000000000000000, got %v", commitment.Payment)
		}
		if commitment.Refund != "100000000000000000" {
			t.Fatalf("expected commitment refund to be 100000000000000000, got %v", commitment.Refund)
		}
		if commitment.ProviderAddress != strings.TrimPrefix(common.HexToAddress("0x1234").Hex(), "0x") {
			t.Fatalf("expected provider address to be 0x1234, got %v", commitment.ProviderAddress)
		}
	})

	t.Run("get bid info multiple commitments", func(t *testing.T) {
		testCommitments := []*preconfstore.Commitment{
			{
				PreConfirmation: &preconfpb.PreConfirmation{
					Bid: &preconfpb.Bid{
						TxHash:              "0x1234567890abcdef,0x1234567890abcdef",
						BlockNumber:         1,
						BidAmount:           "1000000000000000000",
						DecayStartTimestamp: 123456789,
						DecayEndTimestamp:   123457896,
						Digest:              []byte("digest1"),
					},
					DispatchTimestamp: 123456889,
					ProviderAddress:   common.HexToAddress("0x1234").Bytes(),
				},
				Status:  preconfstore.CommitmentStatusOpened,
				Details: "test details",
				Payment: "900000000000000000",
				Refund:  "100000000000000000",
			},
			{
				PreConfirmation: &preconfpb.PreConfirmation{
					Bid: &preconfpb.Bid{
						TxHash:              "0xabcdef1234567890,0xabcdef1234567890",
						BlockNumber:         2,
						BidAmount:           "2000000000000000000",
						DecayStartTimestamp: 123456789,
						DecayEndTimestamp:   123457896,
						Digest:              []byte("digest2"),
					},
					DispatchTimestamp: 123456889,
					ProviderAddress:   common.HexToAddress("0x5678").Bytes(),
				},
				Status:  preconfstore.CommitmentStatusSettled,
				Details: "another test details",
				Payment: "1800000000000000000",
				Refund:  "200000000000000000",
			},
			{
				PreConfirmation: &preconfpb.PreConfirmation{
					Bid: &preconfpb.Bid{
						TxHash:              "0xabcdef1234567890,0xabcdef1234567890",
						BlockNumber:         2,
						BidAmount:           "2000000000000000000",
						DecayStartTimestamp: 123456789,
						DecayEndTimestamp:   123457896,
						Digest:              []byte("digest2"),
					},
					DispatchTimestamp: 123456889,
					ProviderAddress:   common.HexToAddress("0x1278").Bytes(),
				},
				Status:  preconfstore.CommitmentStatusFailed,
				Details: "yet another test details",
				Payment: "2700000000000000000",
				Refund:  "300000000000000000",
			},
		}

		store := &testStore{
			commitments: testCommitments,
		}

		client := startServerWithStore(t, store)
		resp, err := client.GetBidInfo(context.Background(), &bidderapiv1.GetBidInfoRequest{})
		if err != nil {
			t.Fatalf("error getting bid info: %v", err)
		}

		if len(resp.BlockBidInfo) != 2 {
			t.Fatalf("expected 2 bid info, got %v", len(resp.BlockBidInfo))
		}

		for _, bidInfo := range resp.BlockBidInfo {
			switch bidInfo.BlockNumber {
			case 1:
				if len(bidInfo.Bids) != 1 {
					t.Fatalf("expected 1 bid for block 1, got %v", len(bidInfo.Bids))
				}
				bid := bidInfo.Bids[0]
				if len(bid.TxnHashes) != 2 || bid.TxnHashes[0] != "0x1234567890abcdef" || bid.TxnHashes[1] != "0x1234567890abcdef" {
					t.Fatalf("unexpected transaction hashes for block 1, got %v", bid.TxnHashes)
				}
				if len(bid.Commitments) != 1 {
					t.Fatalf("expected 1 commitment for block 1, got %v", len(bid.Commitments))
				}
				commitment := bid.Commitments[0]
				if commitment.Details != "test details" {
					t.Fatalf("expected commitment details to be 'test details', got %v", commitment.Details)
				}
				if commitment.Status != string(preconfstore.CommitmentStatusOpened) {
					t.Fatalf("expected commitment status to be opened, got %v", commitment.Status)
				}
				if commitment.Payment != "900000000000000000" {
					t.Fatalf("expected commitment payment to be 900000000000000000, got %v", commitment.Payment)
				}
				if commitment.Refund != "100000000000000000" {
					t.Fatalf("expected commitment refund to be 100000000000000000, got %v", commitment.Refund)
				}
				if commitment.ProviderAddress != strings.TrimPrefix(common.HexToAddress("0x1234").Hex(), "0x") {
					t.Fatalf("expected provider address to be 0x1234, got %v", commitment.ProviderAddress)
				}
			case 2:
				if len(bidInfo.Bids) != 1 {
					t.Fatalf("expected 1 bid for block 2, got %v", len(bidInfo.Bids))
				}
				bid := bidInfo.Bids[0]
				if len(bid.Commitments) != 2 {
					t.Fatalf("expected 2 commitments for block 2, got %v", len(bid.Commitments))
				}
				for _, commitment := range bid.Commitments {
					if commitment.Details != "another test details" && commitment.Details != "yet another test details" {
						t.Fatalf("unexpected commitment details for block 2, got %v", commitment.Details)
					}
					if commitment.Status != string(preconfstore.CommitmentStatusSettled) && commitment.Status != string(preconfstore.CommitmentStatusFailed) {
						t.Fatalf("unexpected commitment status for block 2, got %v", commitment.Status)
					}
					if commitment.Payment != "1800000000000000000" && commitment.Payment != "2700000000000000000" {
						t.Fatalf("unexpected commitment payment for block 2, got %v", commitment.Payment)
					}
					if commitment.Refund != "200000000000000000" && commitment.Refund != "300000000000000000" {
						t.Fatalf("unexpected commitment refund for block 2, got %v", commitment.Refund)
					}
					if commitment.ProviderAddress != strings.TrimPrefix(common.HexToAddress("0x5678").Hex(), "0x") &&
						commitment.ProviderAddress != strings.TrimPrefix(common.HexToAddress("0x1278").Hex(), "0x") {
						fmt.Println(strings.TrimPrefix(common.HexToAddress("0x1278").Hex(), "0x"))
						t.Fatalf("unexpected provider address for block 2, got %v", commitment.ProviderAddress)
					}
				}
			default:
				t.Fatalf("unexpected block number %v", bidInfo.BlockNumber)
			}
		}
	})
}
