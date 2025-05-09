package providerapi_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/cloudflare/circl/sign/bls"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	stake        *big.Int
	topup        *big.Int
	minStake     *big.Int
	blsKey       []byte
	blsSignature []byte
}

func (t *testRegistryContract) ProviderRegistered(opts *bind.CallOpts, address common.Address) (bool, error) {
	if t.stake.Sign() == 0 {
		return false, nil
	}
	return true, nil
}

func (t *testRegistryContract) RegisterAndStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	t.stake = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) AddVerifiedBLSKey(opts *bind.TransactOpts, blsPublicKey []byte, blsSignature []byte) (*types.Transaction, error) {
	t.blsKey = blsPublicKey
	t.blsSignature = blsSignature
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	t.topup = opts.Value
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) GetProviderStake(_ *bind.CallOpts, address common.Address) (*big.Int, error) {
	return big.NewInt(0).Add(t.stake, t.topup), nil
}

func (t *testRegistryContract) GetBLSKeys(_ *bind.CallOpts, address common.Address) ([][]byte, error) {
	return [][]byte{t.blsKey}, nil
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

func (t *testRegistryContract) ParseBLSKeyAdded(log types.Log) (*providerregistry.ProviderregistryBLSKeyAdded, error) {
	return &providerregistry.ProviderregistryBLSKeyAdded{
		Provider:     common.Address{},
		BlsPublicKey: t.blsKey,
	}, nil
}

func (t *testRegistryContract) Unstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testRegistryContract) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
}

type testBidderRegistryContract struct {
	providerReward *big.Int
}

func (t *testBidderRegistryContract) GetProviderAmount(_ *bind.CallOpts, address common.Address) (*big.Int, error) {
	return t.providerReward, nil
}

func (t *testBidderRegistryContract) WithdrawProviderAmount(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	// Save the current reward amount before setting it to 0
	amount := new(big.Int).Set(t.providerReward)
	// Set to zero as the contract would do
	t.providerReward = big.NewInt(0)
	// Return a transaction with the data field containing the amount (for testing purposes)
	return types.NewTransaction(1, common.Address{}, nil, 0, nil, amount.Bytes()), nil
}

type testWatcher struct{}

func (t *testWatcher) WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	return &types.Receipt{
		Status: 1,
		Logs: []*types.Log{
			{
				Address: common.Address{},
				Topics:  []common.Hash{},
				Data:    tx.Data(),
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
	registryContract := &testRegistryContract{
		stake:    big.NewInt(0),
		topup:    big.NewInt(0),
		minStake: big.NewInt(100000000000000000),
	}

	bidderRegistryContract := &testBidderRegistryContract{
		providerReward: big.NewInt(500000000000000000), // 0.5 ETH initial reward
	}

	srvImpl := providerapi.NewService(
		logger,
		registryContract,
		bidderRegistryContract,
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

	validBLSKey := "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456"
	validSignature := "bbbbbbbbb1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2"
	t.Run("register stake", func(t *testing.T) {
		type testCase struct {
			amount        string
			blsPublicKeys []string
			blsSignatures []string
			err           string
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
				amount:        "1000000000000000000",
				blsPublicKeys: []string{"0x"},
				err:           "bls_public_key must be a valid 48-byte hex string, with optional 0x prefix.",
			},
			{
				amount:        "1000000000000000000",
				blsPublicKeys: []string{"0x12345"},
				err:           "bls_public_key must be a valid 48-byte hex string, with optional 0x prefix.",
			},
			{
				amount:        "1000000000000000000",
				blsPublicKeys: []string{validBLSKey},
				err:           "missing BLS signatures",
			},
			{
				amount:        "1000000000000000000",
				blsPublicKeys: []string{validBLSKey},
				blsSignatures: []string{validSignature},
				err:           "",
			},
			{
				amount: "1000000000000000000",
				err:    "",
			},
		} {
			stake, err := client.Stake(
				context.Background(),
				&providerapiv1.StakeRequest{
					Amount:        tc.amount,
					BlsPublicKeys: tc.blsPublicKeys,
					BlsSignatures: tc.blsSignatures,
				},
			)
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
				if len(tc.blsPublicKeys) > 0 && stake.BlsPublicKeys[0] != tc.blsPublicKeys[0] {
					t.Fatalf("expected bls_public_key to be %v, got %v", tc.blsPublicKeys[0], stake.BlsPublicKeys[0])
				}
			}
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

func TestProviderReward(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t)

	t.Run("get provider reward", func(t *testing.T) {
		reward, err := client.GetProviderReward(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting provider reward: %v", err)
		}
		if reward.Amount != "500000000000000000" {
			t.Fatalf("expected reward amount to be 500000000000000000, got %v", reward.Amount)
		}
	})

	t.Run("withdraw provider reward", func(t *testing.T) {
		withdrawal, err := client.WithdrawProviderReward(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error withdrawing provider reward: %v", err)
		}
		if withdrawal.Amount != "500000000000000000" {
			t.Fatalf("expected withdrawal amount to be 500000000000000000, got %v", withdrawal.Amount)
		}

		// Check that getting the reward now returns 0
		reward, err := client.GetProviderReward(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error getting provider reward after withdrawal: %v", err)
		}
		if reward.Amount != "0" {
			t.Fatalf("expected reward amount to be 0 after withdrawal, got %v", reward.Amount)
		}

		// Try to withdraw again, should still succeed but with 0 amount
		withdrawal, err = client.WithdrawProviderReward(context.Background(), &providerapiv1.EmptyMessage{})
		if err != nil {
			t.Fatalf("error on second provider reward withdrawal: %v", err)
		}
		if withdrawal.Amount != "0" {
			t.Fatalf("expected second withdrawal amount to be 0, got %v", withdrawal.Amount)
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
				SlashAmount:         "0",
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
				SlashAmount:         "0",
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
				SlashAmount:         "0",
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
				SlashAmount:         "0",
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
				SlashAmount:         "0",
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
				SlashAmount:         "0",
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
					if bid.SlashAmount != tc.bid.SlashAmount {
						t.Errorf("expected slash amount to be %v, got %v", tc.bid.SlashAmount, bid.SlashAmount)
					}
					if bid.BlockNumber != tc.bid.BlockNumber {
						t.Errorf("expected block number to be %v, got %v", tc.bid.BlockNumber, bid.BlockNumber)
					}
					bidCh <- bid
				}
			}()

			activeReceiverTimeout := time.Now().Add(2 * time.Second)
			for svc.ActiveReceivers() <= 0 {
				// Check for timeout on each iteration
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

func TestBLSKeys(t *testing.T) {
	// Generate a BLS signature to verify
	message := []byte("adb4257612d45f12570533308b20ac77dbfeb85a")
	hashedMessage := crypto.Keccak256(message)
	ikm := make([]byte, 32)
	privateKey, err := bls.KeyGen[bls.G1](ikm, nil, nil)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	publicKey := privateKey.PublicKey()
	signature := bls.Sign(privateKey, hashedMessage)

	// Verify the signature
	if !bls.Verify(publicKey, hashedMessage, signature) {
		t.Fatalf("Failed to verify generated BLS signature")
	}

	pubkeyb, _ := publicKey.MarshalBinary()
	fmt.Printf("Public Key: %s\n", common.Bytes2Hex(pubkeyb))
	fmt.Printf("Message: %s\n", common.Bytes2Hex(hashedMessage))
	fmt.Printf("Signature: %s\n", common.Bytes2Hex(signature))

}

func TestWithdrawStakedAmount(t *testing.T) {
	t.Parallel()

	client, _ := startServer(t)

	t.Run("withdraw stake", func(t *testing.T) {
		iv := make([]byte, 32)
		_, _ = rand.Read(iv)
		blsPrivKey, _ := bls.KeyGen[bls.G1](iv, []byte{}, []byte{})
		pubKey := blsPrivKey.PublicKey()
		pubKeyBytes, _ := pubKey.MarshalBinary()
		value := common.Hex2Bytes("0x53c61cfb8128ad59244e8c1d26109252ace23d14")
		hash := crypto.Keccak256Hash(value)
		signature := bls.Sign(blsPrivKey, hash.Bytes())

		t.Logf("Public Key: %s", hex.EncodeToString(pubKeyBytes))

		_, err := client.Stake(context.Background(), &providerapiv1.StakeRequest{
			Amount:        "1000000000000000000",
			BlsPublicKeys: []string{hex.EncodeToString(pubKeyBytes)},
			BlsSignatures: []string{hex.EncodeToString(signature)},
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
