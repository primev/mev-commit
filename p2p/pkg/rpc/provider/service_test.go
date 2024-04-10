package providerapi_test

import (
	"context"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	providerapi "github.com/primevprotocol/mev-commit/p2p/pkg/rpc/provider"
	"github.com/primevprotocol/mev-commit/p2p/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testRegistryContract struct {
	stake    *big.Int
	minStake *big.Int
}

func (t *testRegistryContract) RegisterProvider(ctx context.Context, amount *big.Int) error {
	t.stake = amount
	return nil
}

func (t *testRegistryContract) GetStake(ctx context.Context, address common.Address) (*big.Int, error) {
	return t.stake, nil
}

func (t *testRegistryContract) GetMinStake(ctx context.Context) (*big.Int, error) {
	return t.minStake, nil
}

func (t *testRegistryContract) CheckProviderRegistered(ctx context.Context, address common.Address) bool {
	return t.stake.Cmp(t.minStake) > 0
}

type testEVMClient struct {
	pendingTxns     []evmclient.TxnInfo
	cancelledHashes []common.Hash
}

func (t *testEVMClient) PendingTxns() []evmclient.TxnInfo {
	return t.pendingTxns
}

func (t *testEVMClient) CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error) {
	t.cancelledHashes = append(t.cancelledHashes, txHash)
	return txHash, nil
}

func startServer(t *testing.T, evm *testEVMClient) (providerapiv1.ProviderClient, *providerapi.Service) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	logger := util.NewTestLogger(os.Stdout)
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("error creating validator: %v", err)
	}

	owner := common.HexToAddress("0x00001")
	registryContract := &testRegistryContract{minStake: big.NewInt(100000000000000000)}
	if evm == nil {
		evm = &testEVMClient{}
	}

	srvImpl := providerapi.NewService(
		logger,
		registryContract,
		owner,
		evm,
		validator,
	)

	baseServer := grpc.NewServer()
	providerapiv1.RegisterProviderServer(baseServer, srvImpl)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			t.Errorf("error serving server: %v", err)
		}
	}()

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
	})

	client := providerapiv1.NewProviderClient(conn)

	return client, srvImpl
}

func TestStakeHandling(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t, nil)

	t.Run("register stake", func(t *testing.T) {
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
			stake, err := client.RegisterStake(context.Background(), &providerapiv1.StakeRequest{Amount: tc.amount})
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
			client, svc := startServer(t, nil)

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

func TestCancelTx(t *testing.T) {
	t.Parallel()

	evmClient := &testEVMClient{
		pendingTxns: []evmclient.TxnInfo{
			{
				Hash:    common.HexToHash("0x00001").Hex(),
				Nonce:   1,
				Created: time.Now().String(),
			},
			{
				Hash:    common.HexToHash("0x00002").Hex(),
				Nonce:   2,
				Created: time.Now().String(),
			},
		},
	}
	client, _ := startServer(t, evmClient)

	t.Run("get pending txns", func(t *testing.T) {
		pendingTxns, err := client.GetPendingTxns(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting pending txns: %v", err)
		}

		if len(pendingTxns.PendingTxns) != len(evmClient.pendingTxns) {
			t.Fatalf("expected %v pending txns, got %v", len(evmClient.pendingTxns), len(pendingTxns.PendingTxns))
		}

		for i, pendingTxn := range pendingTxns.PendingTxns {
			if pendingTxn.TxHash != evmClient.pendingTxns[i].Hash {
				t.Fatalf("expected tx hash to be %v, got %v", evmClient.pendingTxns[i].Hash, pendingTxn.TxHash)
			}
			if uint64(pendingTxn.Nonce) != evmClient.pendingTxns[i].Nonce {
				t.Fatalf("expected nonce to be %v, got %v", evmClient.pendingTxns[i].Nonce, pendingTxn.Nonce)
			}
			if pendingTxn.Created != evmClient.pendingTxns[i].Created {
				t.Fatalf("expected created to be %v, got %v", evmClient.pendingTxns[i].Created, pendingTxn.Created)
			}
		}
	})

	t.Run("cancel tx", func(t *testing.T) {
		txHash := common.HexToHash("0x00001")
		cancelTxHash, err := client.CancelTransaction(context.Background(), &providerapiv1.CancelReq{TxHash: txHash.Hex()})
		if err != nil {
			t.Fatalf("error cancelling tx: %v", err)
		}
		if cancelTxHash.TxHash != txHash.Hex() {
			t.Fatalf("expected cancel tx hash to be %v, got %v", txHash.Hex(), cancelTxHash.TxHash)
		}
		if len(evmClient.cancelledHashes) != 1 {
			t.Fatalf("expected 1 cancelled tx, got %v", len(evmClient.cancelledHashes))
		}
		if evmClient.cancelledHashes[0] != txHash {
			t.Fatalf("expected cancelled tx hash to be %v, got %v", txHash, evmClient.cancelledHashes[0])
		}
	})
}
