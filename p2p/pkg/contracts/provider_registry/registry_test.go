package registrycontract_test

import (
	"bytes"
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	registrycontract "github.com/primev/mev-commit/p2p/pkg/contracts/provider_registry"
	"github.com/primev/mev-commit/p2p/pkg/evmclient"
	mockevmclient "github.com/primev/mev-commit/p2p/pkg/evmclient/mock"
	"github.com/primev/mev-commit/x/util"
)

func TestRegistryContract(t *testing.T) {
	t.Parallel()

	t.Run("RegisterProvider", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		txHash := common.HexToHash("abcdef")
		amount := big.NewInt(1000000000000000000)

		expCallData, err := registrycontract.RegistryABI().Pack("registerAndStake")
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

		registryContract := registrycontract.New(
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		err = registryContract.RegisterProvider(context.Background(), amount)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetStake", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		amount := big.NewInt(1000000000000000000)
		address := common.HexToAddress("abcdef")

		expCallData, err := registrycontract.RegistryABI().Pack("checkStake", address)
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

		registryContract := registrycontract.New(
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		stakeAmt, err := registryContract.GetStake(context.Background(), address)
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

		expCallData, err := registrycontract.RegistryABI().Pack("minStake")
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

		registryContract := registrycontract.New(
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		stakeAmt, err := registryContract.GetMinStake(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if stakeAmt.Cmp(amount) != 0 {
			t.Fatalf("expected stake amount to be %s, got %s", amount.String(), stakeAmt.String())
		}
	})

	t.Run("CheckProviderRegistered", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		amount := big.NewInt(1000000000000000000)
		address := common.HexToAddress("abcdef")

		mockClient := mockevmclient.New(
			mockevmclient.WithCallFunc(
				func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error) {
					if req.To.Cmp(registryContractAddr) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							registryContractAddr.Hex(), req.To.Hex(),
						)
					}

					return amount.FillBytes(make([]byte, 32)), nil
				},
			),
		)

		registryContract := registrycontract.New(
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		isRegistered := registryContract.CheckProviderRegistered(context.Background(), address)
		if !isRegistered {
			t.Fatalf("expected bidder to be registered")
		}
	})
}
