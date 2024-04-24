package blocktrackercontract

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
)

var blockTrackerABI = func() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerMetaData.ABI))
	if err != nil {
		panic(err)
	}
	return abi
}()

type Interface interface {
	// GetBlocksPerWindow returns the number of blocks per window.
	GetBlocksPerWindow(ctx context.Context) (uint64, error)
	// GetCurrentWindow returns the current window number.
	GetCurrentWindow(ctx context.Context) (uint64, error)
}

type blockTrackerContract struct {
	blockTrackerABI          abi.ABI
	blockTrackerContractAddr common.Address
	client                   evmclient.Interface
	wsClient                 evmclient.Interface
	logger                   *slog.Logger
}

type NewL1BlockEvent struct {
	BlockNumber *big.Int
	Winner      common.Address
	Window      *big.Int
}

func New(
	blockTrackerContractAddr common.Address,
	client evmclient.Interface,
	wsClient evmclient.Interface,
	logger *slog.Logger,
) Interface {
	return &blockTrackerContract{
		blockTrackerABI:          blockTrackerABI,
		blockTrackerContractAddr: blockTrackerContractAddr,
		client:                   client,
		wsClient:                 wsClient,
		logger:                   logger,
	}
}

// GetBlocksPerWindow returns the number of blocks per window.
func (btc *blockTrackerContract) GetBlocksPerWindow(ctx context.Context) (uint64, error) {
	callData, err := btc.blockTrackerABI.Pack("getBlocksPerWindow")
	if err != nil {
		btc.logger.Error("error packing call data for getBlocksPerWindow", "error", err)
		return 0, err
	}

	result, err := btc.client.Call(ctx, &evmclient.TxRequest{
		To:       &btc.blockTrackerContractAddr,
		CallData: callData,
	})
	if err != nil {
		return 0, err
	}

	results, err := btc.blockTrackerABI.Unpack("getBlocksPerWindow", result)
	if err != nil {
		btc.logger.Error("error unpacking result for getBlocksPerWindow", "error", err)
		return 0, err
	}

	blocksPerWindow, ok := results[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("invalid result type")
	}

	return blocksPerWindow.Uint64(), nil
}

// GetCurrentWindow returns the current window number.
func (btc *blockTrackerContract) GetCurrentWindow(ctx context.Context) (uint64, error) {
	callData, err := btc.blockTrackerABI.Pack("getCurrentWindow")
	if err != nil {
		btc.logger.Error("error packing call data for getCurrentWindow", "error", err)
		return 0, err
	}

	result, err := btc.client.Call(ctx, &evmclient.TxRequest{
		To:       &btc.blockTrackerContractAddr,
		CallData: callData,
	})
	if err != nil {
		return 0, err
	}

	results, err := btc.blockTrackerABI.Unpack("getCurrentWindow", result)
	if err != nil {
		btc.logger.Error("error unpacking result for getCurrentWindow", "error", err)
		return 0, err
	}

	currentWindow, ok := results[0].(*big.Int)
	if !ok {
		return 0, fmt.Errorf("invalid result type")
	}

	return currentWindow.Uint64(), nil
}
