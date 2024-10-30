package providerapi_test

import (
	"context"
	"errors"
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
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	providerapi "github.com/primev/mev-commit/p2p/pkg/rpc/provider"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testRegistryContract struct {
	stake    *big.Int
	topup    *big.Int
	minStake *big.Int
	blsKey   []byte
}

func (t *testRegistryContract) RegisterAndStake(opts *bind.TransactOpts, blsPublicKey []byte) (*types.Transaction, error) {
	t.stake = opts.Value
	t.blsKey = blsPublicKey
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	t.topup = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) GetProviderStake(_ *bind.CallOpts, address common.Address) (*big.Int, error) {
	return big.NewInt(0).Add(t.stake, t.topup), nil
}

func (t *testRegistryContract) GetBLSKey(_ *bind.CallOpts, address common.Address) ([]byte, error) {
	return t.blsKey, nil
}

func (t *testRegistryContract) MinStake(_ *bind.CallOpts) (*big.Int, error) {
	return t.minStake, nil
}

func (t *testRegistryContract) ParseProviderRegistered(log types.Log) (*providerregistry.ProviderregistryProviderRegistered, error) {
	return &providerregistry.ProviderregistryProviderRegistered{
		Provider:     common.Address{},
		StakedAmount: t.stake,
	}, nil
}

func (t *testRegistryContract) ParseFundsDeposited(log types.Log) (*providerregistry.ProviderregistryFundsDeposited, error) {
	return &providerregistry.ProviderregistryFundsDeposited{
		Provider: common.Address{},
		Amount:   t.topup,
	}, nil
}

func (t *testRegistryContract) ParseWithdraw(log types.Log) (*providerregistry.ProviderregistryWithdraw, error) {
	return &providerregistry.ProviderregistryWithdraw{
		Provider: common.Address{},
		Amount:   t.stake,
	}, nil
}

func (t *testRegistryContract) ParseUnstake(log types.Log) (*providerregistry.ProviderregistryUnstake, error) {
	return &providerregistry.ProviderregistryUnstake{
		Provider:  common.Address{},
		Timestamp: new(big.Int).SetInt64(time.Now().Unix()),
	}, nil
}

func (t *testRegistryContract) Unstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

type testWatcher struct{}

func (t *testWatcher) WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
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

func startServer(t *testing.T) (providerapiv1.ProviderClient, *providerapi.Service) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	logger := util.NewTestLogger(os.Stdout)
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("error creating validator: %v", err)
	}

	owner := common.HexToAddress("0x00001")
	registryContract := &testRegistryContract{minStake: big.NewInt(100000000000000000)}

	srvImpl := providerapi.NewService(
		logger,
		registryContract,
		owner,
		&testWatcher{},
		func(context.Context) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{
				From:    owner,
				Context: context.Background(),
			}, nil
		},
		validator,
	)

	baseServer := grpc.NewServer()
	providerapiv1.RegisterProviderServer(baseServer, srvImpl)
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
		t.Fatalf("error connecting to server: %v", err)
	}

	t.Cleanup(func() {
		err := lis.Close()
		if err != nil {
			t.Errorf("error closing listener: %v", err)
		}
		baseServer.Stop()

		<-srvStopped
	})

	client := providerapiv1.NewProviderClient(conn)

	return client, srvImpl
}

func TestStakeHandling(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t)

	t.Run("register stake", func(t *testing.T) {
		type testCase struct {
			amount       string
			blsPublicKey string
			err          string
		}

		for _, tc := range []testCase{
			{
				amount:       "",
				blsPublicKey: "",
				err:          "amount must be a valid integer",
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
				amount:       "1000000000000000000",
				blsPublicKey: "0x",
				err:          "bls_public_key must be a valid 48-byte hex string, with optional 0x prefix.",
			},
			{
				amount:       "1000000000000000000",
				blsPublicKey: "0x12345",
				err:          "bls_public_key must be a valid 48-byte hex string, with optional 0x prefix.",
			},
			{
				amount:       "1000000000000000000",
				blsPublicKey: "0x123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
				err:          "",
			},
			{
				amount:       "1000000000000000000",
				blsPublicKey: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
				err:          "",
			},
		} {
			stake, err := client.RegisterStake(context.Background(),
				&providerapiv1.StakeRequest{Amount: tc.amount, BlsPublicKey: tc.blsPublicKey})
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("expected error staking: %s got %v", tc.err, err)
				}
			} else {
				if err != nil {
					t.Fatalf("error staking: %v", err)
				}
				if stake.Amount != tc.amount {
					t.Fatalf("expected amount to be %v, got %v", tc.amount, stake.Amount)
				}
				tc.blsPublicKey = strings.TrimPrefix(tc.blsPublicKey, "0x")
				if stake.BlsPublicKey != tc.blsPublicKey {
					t.Fatalf("expected bls_public_key to be %v, got %v", tc.blsPublicKey, stake.BlsPublicKey)
				}
			}
		}
	})

	t.Run("get stake", func(t *testing.T) {
		stake, err := client.GetStake(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting stake: %v", err)
		}
		if stake.Amount != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", stake.Amount)
		}
		if stake.BlsPublicKey != "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456" {
			t.Fatalf("expected bls public key to be 123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456, got %v", stake.BlsPublicKey)
		}
	})

	t.Run("get min stake", func(t *testing.T) {
		stake, err := client.GetMinStake(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting min stake: %v", err)
		}
		if stake.Amount != "100000000000000000" {
			t.Fatalf("expected amount to be 100000000000000000, got %v", stake.Amount)
		}
	})
}

func TestBidHandling(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                   string
		bid                    *preconfpb.Bid
		status                 providerapiv1.BidResponse_Status
		noStatus               bool
		processErr             string
		decayDispatchTimestamp int64
	}

	for _, tc := range []testCase{
		{
			name: "accepted bid",
			bid: &preconfpb.Bid{
				TxHash: strings.Join(
					[]string{
						common.HexToHash("0x00001").Hex()[2:], // remove 0x
						common.HexToHash("0x00002").Hex()[2:], // remove 0x
					}, ",", // join with comma
				),
				BidAmount:           "1000000000000000000",
				BlockNumber:         1,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			status:                 providerapiv1.BidResponse_STATUS_ACCEPTED,
			decayDispatchTimestamp: 10,
		},
		{
			name: "rejected bid",
			bid: &preconfpb.Bid{
				TxHash:              common.HexToHash("0x00003").Hex()[2:], // remove 0x
				BidAmount:           "1000000000000000000",
				BlockNumber:         1,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			status:                 providerapiv1.BidResponse_STATUS_REJECTED,
			decayDispatchTimestamp: 10,
		},
		{
			name: "invalid bid status",
			bid: &preconfpb.Bid{
				TxHash:              common.HexToHash("0x00003").Hex()[2:], // remove 0x
				BidAmount:           "1000000000000000000",
				BlockNumber:         1,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			status:                 providerapiv1.BidResponse_STATUS_UNSPECIFIED,
			noStatus:               true,
			decayDispatchTimestamp: 10,
		},
		{
			name: "invalid bid txHash",
			bid: &preconfpb.Bid{
				TxHash:              "asdf",
				BidAmount:           "1000000000000000000",
				BlockNumber:         1,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			processErr:             "tx_hashes: tx_hashes must be a valid array of transaction hashes",
			decayDispatchTimestamp: 10,
		},
		{
			name: "invalid bid amount",
			bid: &preconfpb.Bid{
				TxHash:              common.HexToHash("0x00004").Hex()[2:], // remove 0x
				BidAmount:           "0000000000000000000",
				BlockNumber:         1,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			processErr:             "bid_amount: bid_amount must be a valid integer",
			decayDispatchTimestamp: 10,
		},
		{
			name: "invalid bid block number",
			bid: &preconfpb.Bid{
				TxHash:              common.HexToHash("0x00004").Hex()[2:], // remove 0x
				BidAmount:           "1000000000000000000",
				BlockNumber:         0,
				Digest:              []byte("digest"),
				Signature:           []byte("signature"),
				DecayStartTimestamp: 199,
				DecayEndTimestamp:   299,
			},
			processErr:             "block_number: value must be greater than 0",
			decayDispatchTimestamp: 10,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, svc := startServer(t)

			bidCh := make(chan *providerapiv1.Bid)

			rcvr, err := client.ReceiveBids(context.Background(), &providerapiv1.EmptyMessage{})
			if err != nil {
				t.Fatalf("error receiving bids: %v", err)
			}
			go func() {
				defer func() { _ = rcvr.CloseSend() }()
				for {
					bid, err := rcvr.Recv()
					if err != nil {
						break
					}
					sentBidHashes := strings.Split(tc.bid.TxHash, ",")
					if len(bid.TxHashes) != len(sentBidHashes) {
						t.Errorf("expected %v tx hashes, got %v", len(sentBidHashes), len(bid.TxHashes))
					}
					for i, sentBidHash := range sentBidHashes {
						if sentBidHash != bid.TxHashes[i] {
							t.Errorf("expected tx hash %v to be %v, got %v", i, sentBidHash, bid.TxHashes[i])
						}
					}
					if bid.BidAmount != tc.bid.BidAmount {
						t.Errorf("expected bid amount to be %v, got %v", tc.bid.BidAmount, bid.BidAmount)
					}
					if bid.BlockNumber != tc.bid.BlockNumber {
						t.Errorf("expected block number to be %v, got %v", tc.bid.BlockNumber, bid.BlockNumber)
					}
					bidCh <- bid
				}
			}()

			activeReceiverTimeout := time.Now().Add(2 * time.Second)
			for {
				if svc.ActiveReceivers() > 0 {
					break
				}
				if time.Now().After(activeReceiverTimeout) {
					t.Fatalf("timed out waiting for active receivers")
				}
				time.Sleep(10 * time.Millisecond)
			}

			sndr, err := client.SendProcessedBids(context.Background())
			if err != nil {
				t.Fatalf("error sending processed bids: %v", err)
			}
			go func() {
				defer func() { _ = sndr.CloseSend() }()
				for {
					bid, more := <-bidCh
					if !more {
						break
					}
					err := sndr.Send(&providerapiv1.BidResponse{
						BidDigest:         bid.BidDigest,
						Status:            tc.status,
						DispatchTimestamp: tc.decayDispatchTimestamp,
					})
					if err != nil {
						break
					}
				}
			}()

			respC, err := svc.ProcessBid(context.Background(), tc.bid)
			if err != nil {
				if tc.processErr != "" {
					if !strings.Contains(err.Error(), tc.processErr) {
						t.Fatalf("expected error processing bid: %v, got %v", tc.processErr, err)
					}
					return
				}
				t.Fatalf("error processing bid: %v", err)
			}

			select {
			case resp := <-respC:
				if resp.Status != tc.status {
					t.Fatalf("expected status to be %v, got %v", tc.status, resp.Status)
				}
				if tc.noStatus {
					t.Fatalf("expected no status, got %v", resp)
				}
			case <-time.After(2 * time.Second):
				if !tc.noStatus {
					t.Fatalf("expected status to be %v, got timeout", tc.status)
				}
			}
		})
	}
}

func TestWithdrawStakedAmount(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t)

	t.Run("withdraw stake", func(t *testing.T) {
		_, err := client.RegisterStake(context.Background(), &providerapiv1.StakeRequest{
			Amount:       "1000000000000000000",
			BlsPublicKey: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
		})
		if err != nil {
			t.Fatalf("error depositing stake: %v", err)
		}
		withdrawalResp, err := client.WithdrawStake(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error withdrawing stake: %v", err)
		}
		if withdrawalResp.Amount != "1000000000000000000" {
			t.Fatalf("expected amount to be 1000000000000000000, got %v", withdrawalResp.Amount)
		}
	})
}

func TestRequestWithdrawal(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t)

	t.Run("request withdrawal", func(t *testing.T) {
		_, err := client.Unstake(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error requesting withdrawal: %v", err)
		}
	})
}
