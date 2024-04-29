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
	// DepositForSpecificWindow registers a bidder with the bidder_registry contract for a specific window.
	DepositForSpecificWindow(ctx context.Context, amount, window *big.Int) error
	// GetDeposit returns the stake of a bidder.
	GetDeposit(ctx context.Context, address common.Address, window *big.Int) (*big.Int, error)
	// GetMinDeposit returns the minimum stake required to register as a bidder.
	GetMinDeposit(ctx context.Context) (*big.Int, error)
	// CheckBidderRegistred returns true if bidder is registered
	CheckBidderDeposit(ctx context.Context, address common.Address, window, blocksPerWindow *big.Int) bool
	// WithdrawDeposit withdraws the stake of a bidder.
	WithdrawDeposit(ctx context.Context, window *big.Int) (*big.Int, error)
}

type bidderRegistryContract struct {
	owner                      common.Address
	bidderRegistryABI          abi.ABI
	bidderRegistryContractAddr common.Address
	client                     evmclient.Interface
	logger                     *slog.Logger
}

func New(
	owner common.Address,
	bidderRegistryContractAddr common.Address,
	client evmclient.Interface,
	logger *slog.Logger,
) Interface {
	return &bidderRegistryContract{
		owner:                      owner,
		bidderRegistryABI:          bidderRegistryABI(),
		bidderRegistryContractAddr: bidderRegistryContractAddr,
		client:                     client,
		logger:                     logger,
	}
}

func (r *bidderRegistryContract) DepositForSpecificWindow(ctx context.Context, amount, window *big.Int) error {
	callData, err := r.bidderRegistryABI.Pack("depositForSpecificWindow", window)
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
			"deposit failed for bidder registry",
			"txnHash", txnHash,
			"receipt", receipt,
		)
		return err
	}

	return nil
}

func (r *bidderRegistryContract) GetDeposit(
	ctx context.Context,
	address common.Address,
	window *big.Int,
) (*big.Int, error) {
	callData, err := r.bidderRegistryABI.Pack("getDeposit", address, window)
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

	results, err := r.bidderRegistryABI.Unpack("getDeposit", result)
	if err != nil {
		r.logger.Error("error unpacking result", "error", err)
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (r *bidderRegistryContract) GetMinDeposit(ctx context.Context) (*big.Int, error) {
	callData, err := r.bidderRegistryABI.Pack("minDeposit")
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

	results, err := r.bidderRegistryABI.Unpack("minDeposit", result)
	if err != nil {
		r.logger.Error("error unpacking result", "error", err)
		return nil, err
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (r *bidderRegistryContract) WithdrawDeposit(ctx context.Context, window *big.Int) (*big.Int, error) {
	callData, err := r.bidderRegistryABI.Pack("withdrawBidderAmountFromWindow", r.owner, window)
	if err != nil {
		r.logger.Error("error packing call data", "error", err)
		return nil, err
	}

	txnHash, err := r.client.Send(ctx, &evmclient.TxRequest{
		To:       &r.bidderRegistryContractAddr,
		CallData: callData,
	})
	if err != nil {
		return nil, err
	}

	receipt, err := r.client.WaitForReceipt(ctx, txnHash)
	if err != nil {
		return nil, err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		r.logger.Error(
			"withdraw failed for bidder registry",
			"txnHash", txnHash,
			"receipt", receipt,
		)
		return nil, err
	}

	var bidderWithdrawn struct {
		Bidder common.Address
		Amount *big.Int
		Window *big.Int
	}

	for _, log := range receipt.Logs {
		if len(log.Topics) > 1 {
			bidderWithdrawn.Bidder = common.HexToAddress(log.Topics[1].Hex())
		}

		err := r.bidderRegistryABI.UnpackIntoInterface(&bidderWithdrawn, "BidderWithdrawn", log.Data)
		if err != nil {
			r.logger.Debug("Failed to unpack event", "err", err)
			continue
		}
		r.logger.Info("bidder withdrawn", "address", bidderWithdrawn.Bidder, "withdrawn", bidderWithdrawn.Amount.Uint64(), "windowNumber", bidderWithdrawn.Window.Int64())
	}

	r.logger.Info("withdraw successful for bidder registry", "txnHash", txnHash, "bidder", bidderWithdrawn.Bidder)

	return bidderWithdrawn.Amount, nil
}

func (r *bidderRegistryContract) CheckBidderDeposit(
	ctx context.Context,
	address common.Address,
	window *big.Int,
	blocksPerWindow *big.Int,
) bool {
	minStake, err := r.GetMinDeposit(ctx)
	if err != nil {
		r.logger.Error("error getting min stake", "error", err)
		return false
	}

	stake, err := r.GetDeposit(ctx, address, window)
	if err != nil {
		r.logger.Error("error getting stake", "error", err)
		return false
	}
	r.logger.Info("checking bidder deposit",
		"stake", stake.Uint64(),
		"blocksPerWindow", blocksPerWindow.Uint64(),
		"minStake", minStake.Uint64(),
		"window", window.Uint64(),
		"address", address.Hex(),
	)
	return (stake.Div(stake, blocksPerWindow)).Cmp(minStake) >= 0
}
