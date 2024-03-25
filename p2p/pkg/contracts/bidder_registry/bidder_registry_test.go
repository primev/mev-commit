package bidderregistrycontract_test

import (
	"bytes"
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidder_registrycontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/bidder_registry"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	mockevmclient "github.com/primevprotocol/mev-commit/p2p/pkg/evmclient/mock"
	"github.com/primevprotocol/mev-commit/p2p/pkg/util"
)

func TestBidderRegistryContract(t *testing.T) {
	t.Parallel()

	t.Run("PrepayAllowance", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		txHash := common.HexToHash("abcdef")
		amount := big.NewInt(1000000000000000000)

		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("prepay")
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
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		err = registryContract.PrepayAllowance(context.Background(), amount)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("GetAllowance", func(t *testing.T) {
		registryContractAddr := common.HexToAddress("abcd")
		amount := big.NewInt(1000000000000000000)
		address := common.HexToAddress("abcdef")

		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("getAllowance", address)
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
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		stakeAmt, err := registryContract.GetAllowance(context.Background(), address)
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

		expCallData, err := bidder_registrycontract.BidderRegistryABI().Pack("minAllowance")
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
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		stakeAmt, err := registryContract.GetMinAllowance(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if stakeAmt.Cmp(amount) != 0 {
			t.Fatalf("expected stake amount to be %s, got %s", amount.String(), stakeAmt.String())
		}
	})

	t.Run("CheckBidderAllowance", func(t *testing.T) {
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

		registryContract := bidder_registrycontract.New(
			registryContractAddr,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		isRegistered := registryContract.CheckBidderAllowance(context.Background(), address)
		if !isRegistered {
			t.Fatal("expected bidder to be registered")
		}
	})
}
