package bidderregistrycontract

import (
	"context"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primevprotocol/mev-commit/contracts-abi/clients/BidderRegistry"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
)

var bidderRegistryABI = func() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryMetaData.ABI))
	if err != nil {
		panic(err)
	}
	return abi
}

type Interface interface {
	// PrepayAllowance registers a bidder with the bidder_registry contract.
	PrepayAllowance(ctx context.Context, amount *big.Int) error
	// GetAllowance returns the stake of a bidder.
	GetAllowance(ctx context.Context, address common.Address) (*big.Int, error)
	// GetMinAllowance returns the minimum stake required to register as a bidder.
	GetMinAllowance(ctx context.Context) (*big.Int, error)
	// CheckBidderRegistred returns true if bidder is registered
	CheckBidderAllowance(ctx context.Context, address common.Address) bool
}

type bidderRegistryContract struct {
	bidderRegistryABI          abi.ABI
	bidderRegistryContractAddr common.Address
	client                     evmclient.Interface
	logger                     *slog.Logger
}

func New(
	bidderRegistryContractAddr common.Address,
	client evmclient.Interface,
	logger *slog.Logger,
) Interface {
	return &bidderRegistryContract{
		bidderRegistryABI:          bidderRegistryABI(),
		bidderRegistryContractAddr: bidderRegistryContractAddr,
		client:                     client,
		logger:                     logger,
	}
}

func (r *bidderRegistryContract) PrepayAllowance(ctx context.Context, amount *big.Int) error {
	callData, err := r.bidderRegistryABI.Pack("prepay")
	if err != nil {
		r.logger.Error("error packing call data", "error", err)
		return err
	}

	txnHash, err := r.client.Send(ctx, &evmclient.TxRequest{
		To:       &r.bidderRegistryContractAddr,
		CallData: callData,
		Value:    amount,
	})
	if err != nil {
		return err
	}

	receipt, err := r.client.WaitForReceipt(ctx, txnHash)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		r.logger.Error(
			"prepay failed for bidder registry",
			"txnHash", txnHash,
			"receipt", receipt,
		)
		return err
	}

	r.logger.Info("prepay successful for bidder registry", "txnHash", txnHash)

	return nil
}

func (r *bidderRegistryContract) GetAllowance(
	ctx context.Context,
	address common.Address,
) (*big.Int, error) {
	callData, err := r.bidderRegistryABI.Pack("getAllowance", address)
	if err != nil {
		r.logger.Error("error packing call data", "error", err)
		return nil, err
	}

	result, err := r.client.Call(ctx, &evmclient.TxRequest{
		To:       &r.bidderRegistryContractAddr,
		CallData: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := r.bidderRegistryABI.Unpack("getAllowance", result)
	if err != nil {
		r.logger.Error("error unpacking result", "error", err)
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (r *bidderRegistryContract) GetMinAllowance(ctx context.Context) (*big.Int, error) {
	callData, err := r.bidderRegistryABI.Pack("minAllowance")
	if err != nil {
		r.logger.Error("error packing call data", "error", err)
		return nil, err
	}

	result, err := r.client.Call(ctx, &evmclient.TxRequest{
		To:       &r.bidderRegistryContractAddr,
		CallData: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := r.bidderRegistryABI.Unpack("minAllowance", result)
	if err != nil {
		r.logger.Error("error unpacking result", "error", err)
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (r *bidderRegistryContract) CheckBidderAllowance(
	ctx context.Context,
	address common.Address,
) bool {

	minStake, err := r.GetMinAllowance(ctx)
	if err != nil {
		r.logger.Error("error getting min stake", "error", err)
		return false
	}

	stake, err := r.GetAllowance(ctx, address)
	if err != nil {
		r.logger.Error("error getting stake", "error", err)
		return false
	}

	return stake.Cmp(minStake) >= 0
}
