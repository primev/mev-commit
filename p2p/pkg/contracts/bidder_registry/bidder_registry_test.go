package bidderregistrycontract_test

import (
	"bytes"
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidder_registrycontract "github.com/primev/mev-commit/p2p/pkg/contracts/bidder_registry"
	"github.com/primev/mev-commit/p2p/pkg/evmclient"
	mockevmclient "github.com/primev/mev-commit/p2p/pkg/evmclient/mock"
	"github.com/primev/mev-commit/x/util"
)

func TestBidderRegistryContract(t *testing.T) {
	t.Parallel()

	owner := common.HexToAddress("abcd")

	t.Run("Deposit", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		txHash := common.HexToHash("abcdef")
		amount := big.NewInt(1000000000000000000)
		window := big.NewInt(1)

		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("depositForSpecificWindow", window)
		if err != nil {
			t.Fatal(err)
		}

		mockClient := mockevmclient.New(
			mockevmclient.WithSendFunc(
				func(ctx context.Context, req *evmclient.TxRequest) (common.Hash, error) {
					if req.To.Cmp(registryContractAddr) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							registryContractAddr.Hex(), req.To.Hex(),
						)
					}
					if !bytes.Equal(req.CallData, expCallData) {
						t.Fatalf("expected call data to be %x, got %x", expCallData, req.CallData)
					}
					if req.Value.Cmp(amount) != 0 {
						t.Fatalf(
							"expected amount to be %s, got %s",
							amount.String(), req.Value.String(),
						)
					}

					return txHash, nil
				},
			),
			mockevmclient.WithWaitForReceiptFunc(
				func(ctx context.Context, txnHash common.Hash) (*types.Receipt, error) {
					if txnHash != txHash {
						t.Fatalf(
							"expected txn hash to be %s, got %s",
							txHash.Hex(), txnHash.Hex(),
						)
					}
					return &types.Receipt{
						Status: 1,
					}, nil
				},
			),
		)

		registryContract := bidder_registrycontract.New(
			owner,
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)
		err = registryContract.DepositForSpecificWindow(context.Background(), amount, big.NewInt(1))
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetDeposit", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		amount := big.NewInt(1000000000000000000)
		address := common.HexToAddress("abcdef")
		window := big.NewInt(1)
		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("getDeposit", address, window)
		if err != nil {
			t.Fatal(err)
		}

		mockClient := mockevmclient.New(
			mockevmclient.WithCallFunc(
				func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error) {
					if req.To.Cmp(registryContractAddr) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							registryContractAddr.Hex(), req.To.Hex(),
						)
					}
					if !bytes.Equal(req.CallData, expCallData) {
						t.Fatalf("expected call data to be %x, got %x", expCallData, req.CallData)
					}

					return amount.FillBytes(make([]byte, 32)), nil
				},
			),
		)

		registryContract := bidder_registrycontract.New(
			owner,
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)
		stakeAmt, err := registryContract.GetDeposit(context.Background(), address, window)
		if err != nil {
			t.Fatal(err)
		}

		if stakeAmt.Cmp(amount) != 0 {
			t.Fatalf("expected stake amount to be %s, got %s", amount.String(), stakeAmt.String())
		}
	})

	t.Run("GetMinimalStake", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		amount := big.NewInt(1000000000000000000)

		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("minDeposit")
		if err != nil {
			t.Fatal(err)
		}

		mockClient := mockevmclient.New(
			mockevmclient.WithCallFunc(
				func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error) {
					if req.To.Cmp(registryContractAddr) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							registryContractAddr.Hex(), req.To.Hex(),
						)
					}
					if !bytes.Equal(req.CallData, expCallData) {
						t.Fatalf("expected call data to be %x, got %x", expCallData, req.CallData)
					}

					return amount.FillBytes(make([]byte, 32)), nil
				},
			),
		)

		registryContract := bidder_registrycontract.New(
			owner,
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		stakeAmt, err := registryContract.GetMinDeposit(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if stakeAmt.Cmp(amount) != 0 {
			t.Fatalf("expected stake amount to be %s, got %s", amount.String(), stakeAmt.String())
		}
	})

	t.Run("CheckBidderDeposit", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		blocksPerWindow := big.NewInt(64)
		amount := new(big.Int).Mul(big.NewInt(1000000000000000000), blocksPerWindow)
		address := common.HexToAddress("abcdef")

		callCount := 0

		mockClient := mockevmclient.New(
			mockevmclient.WithCallFunc(
				func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error) {
					callCount++
					if req.To.Cmp(registryContractAddr) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							registryContractAddr.Hex(), req.To.Hex(),
						)
					}

					if callCount == 1 {
						return new(big.Int).Div(amount, blocksPerWindow).FillBytes(make([]byte, 32)), nil
					}

					return amount.FillBytes(make([]byte, 32)), nil
				},
			),
		)

		registryContract := bidder_registrycontract.New(
			owner,
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		window := big.NewInt(1)
		isRegistered := registryContract.CheckBidderDeposit(context.Background(), address, window, blocksPerWindow)
		if !isRegistered {
			t.Fatal("expected bidder to be registered")
		}
	})
}
