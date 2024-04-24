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
	"github.com/ethereum/go-ethereum/common"
	bidderapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/bidderapi/v1"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	bidderapi "github.com/primevprotocol/mev-commit/p2p/pkg/rpc/bidder"
	"github.com/primevprotocol/mev-commit/x/util"
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

func (t *testRegistryContract) DepositForSpecificWindow(ctx context.Context, amount, window *big.Int) error {
	t.deposit = amount
	return nil
}

func (t *testRegistryContract) GetDeposit(ctx context.Context, address common.Address, window *big.Int) (*big.Int, error) {
	return t.deposit, nil
}

func (t *testRegistryContract) GetMinDeposit(ctx context.Context) (*big.Int, error) {
	return t.minDeposit, nil
}

func (t *testRegistryContract) CheckBidderDeposit(ctx context.Context, address common.Address, window, numberOfRounds *big.Int) bool {
	return t.deposit.Cmp(t.minDeposit) > 0
}

func (t *testRegistryContract) WithdrawDeposit(ctx context.Context, window *big.Int) error {
	return nil
}

type testBlockTrackerContract struct {
	blockNumberToWinner map[uint64]common.Address
	lastBlockNumber     uint64
	blocksPerWindow     uint64
}

// GetCurrentWindow returns the current window number.
func (btc *testBlockTrackerContract) GetCurrentWindow(ctx context.Context) (uint64, error) {
	return btc.lastBlockNumber / btc.blocksPerWindow, nil
}

// GetBlocksPerWindow returns the number of blocks per window.
func (btc *testBlockTrackerContract) GetBlocksPerWindow(ctx context.Context) (uint64, error) {
	return btc.blocksPerWindow, nil
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
	srvImpl := bidderapi.NewService(
		sender,
		owner,
		registryContract,
		blockTrackerContract,
		validator,
		logger,
	)

	baseServer := grpc.NewServer()
	bidderapiv1.RegisterBidderServer(baseServer, srvImpl)
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
			// {
			// 	amount: "",
			// 	err:    "amount must be a valid integer",
			// },
			// {
			// 	amount: "0000000000000000000",
			// 	err:    "amount must be a valid integer",
			// },
			// {
			// 	amount: "asdf",
			// 	err:    "amount must be a valid integer",
			// },
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
