package txmonitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Debugger defines an interface for EVM debugging tools.
type Debugger interface {
	// RevertReason retrieves the revert reason for a transaction.
	RevertReason(ctx context.Context, receipt *types.Receipt) (string, error)
}

// Result represents the result of a transaction receipt retrieval operation.
type Result struct {
	// Receipt is the transaction receipt.
	Receipt *types.Receipt
	// Err is the error encountered during the operation.
	Err error
}

// BatchReceiptGetter is an interface for retrieving multiple receipts
type BatchReceiptGetter interface {
	// BatchReceipts retrieves multiple receipts for a list of transaction hashes.
	BatchReceipts(ctx context.Context, txHashes []common.Hash) ([]Result, error)
}

type evmHelper struct {
	client *ethclient.Client
	logger *slog.Logger
}

func NewEVMHelperWithLogger(client *ethclient.Client, logger *slog.Logger) *evmHelper {
	return &evmHelper{client, logger}
}

// BatchReceipts retrieves multiple receipts for a list of transaction hashes.
func (e *evmHelper) BatchReceipts(ctx context.Context, txHashes []common.Hash) ([]Result, error) {
	e.logger.Debug("Starting BatchReceipts", "txHashes", txHashes)
	batch := make([]rpc.BatchElem, len(txHashes))

	for i, hash := range txHashes {
		e.logger.Debug("Preparing batch element", "index", i, "hash", hash.Hex())
		batch[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{hash},
			Result: new(types.Receipt),
		}
	}

	var receipts []Result
	var err error
	for attempts := 0; attempts < 50; attempts++ {
		e.logger.Debug("Attempting batch call", "attempt", attempts+1)
		// Execute the batch request
		err = e.client.Client().BatchCallContext(context.Background(), batch)
		if err != nil {
			e.logger.Error("Batch call attempt failed", "attempt", attempts+1, "error", err)
			time.Sleep(1 * time.Second)
		} else {
			e.logger.Debug("Batch call attempt succeeded", "attempt", attempts+1)
			break
		}
	}

	if err != nil {
		e.logger.Error("All batch call attempts failed", "error", err)
		return nil, err
	}

	receipts = make([]Result, len(batch))
	for i, elem := range batch {
		e.logger.Debug("Processing batch element", "index", i, "result", elem.Result, "error", elem.Error)
		receipts[i].Receipt = elem.Result.(*types.Receipt)
		if elem.Error != nil {
			receipts[i].Err = elem.Error
		}
	}

	// Retry individual failed transactions
	for i, receipt := range receipts {
		if receipt.Err != nil {
			e.logger.Info("Retrying failed transaction", "index", i, "hash", txHashes[i].Hex())
			for attempts := 0; attempts < 50; attempts++ {
				e.logger.Info("Attempting individual call", "attempt", attempts+1, "hash", txHashes[i].Hex())
				err = e.client.Client().CallContext(context.Background(), receipt.Receipt, "eth_getTransactionReceipt", txHashes[i])
				if err == nil {
					e.logger.Info("Individual call succeeded", "attempt", attempts+1, "hash", txHashes[i].Hex())
					receipts[i].Err = nil
					break
				}
				e.logger.Error("Individual call attempt failed", "attempt", attempts+1, "hash", txHashes[i].Hex(), "error", err)
				time.Sleep(1 * time.Second)
			}
		}
	}

	e.logger.Debug("BatchReceipts completed successfully", "receipts", receipts)
	return receipts, nil
}

func (e *evmHelper) RevertReason(ctx context.Context, receipt *types.Receipt) (string, error) {
	tx, _, err := e.client.TransactionByHash(ctx, receipt.TxHash)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   tx.To(),
		Data: tx.Data(),
	}

	reason, err := e.client.CallContract(ctx, msg, receipt.BlockNumber)
	if err != nil {
		return "", err
	}

	return string(reason), nil
}
