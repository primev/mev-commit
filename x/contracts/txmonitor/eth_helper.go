package txmonitor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// Debugger defines an interface for EVM debugging tools.
type Debugger interface {
	// RevertReason retrieves the revert reason for a transaction.
	RevertReason(ctx context.Context, receipt *types.Receipt, from common.Address) (string, error)
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

type EVMHelperImpl struct {
	client       *ethclient.Client
	logger       *slog.Logger
	contractABIs map[common.Address]*abi.ABI
}

func NewEVMHelperWithLogger(client *ethclient.Client, logger *slog.Logger, contracts map[common.Address]*abi.ABI) *EVMHelperImpl {
	return &EVMHelperImpl{client, logger, contracts}
}

// BatchReceipts retrieves multiple receipts for a list of transaction hashes.
func (e *EVMHelperImpl) BatchReceipts(ctx context.Context, txHashes []common.Hash) ([]Result, error) {
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

func (e *EVMHelperImpl) RevertReason(
	ctx context.Context,
	receipt *types.Receipt,
	from common.Address,
) (string, error) {
	tx, _, err := e.client.TransactionByHash(ctx, receipt.TxHash)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	_, err = e.client.CallContract(ctx, msg, receipt.BlockNumber)
	if err == nil {
		return "no revert reason", nil
	}
	var dErr rpc.DataError
	if !errors.As(err, &dErr) {
		return err.Error(), nil
	}

	errStr, ok := dErr.ErrorData().(string)
	if !ok {
		return err.Error(), nil
	}

	buf := common.FromHex(errStr)
	if reason, rErr := abi.UnpackRevert(buf); rErr == nil {
		return reason, nil
	}

	if contractABI, ok := e.contractABIs[*tx.To()]; ok {
		reason := unwrapABIError(err, contractABI)
		return reason, nil
	}

	return err.Error(), nil
}

func unwrapABIError(err error, contractABI *abi.ABI) string {
	var dErr rpc.DataError
	if !errors.As(err, &dErr) {
		return err.Error()
	}

	errStr, ok := dErr.ErrorData().(string)
	if !ok {
		return err.Error()
	}

	buf := common.FromHex(errStr)
	if reason, rErr := abi.UnpackRevert(buf); rErr == nil {
		return reason
	}

	if contractABI == nil {
		return err.Error()
	}

	for _, cErr := range contractABI.Errors {
		if !bytes.Equal(cErr.ID[:4], buf[:4]) {
			continue
		}
		errData, uErr := cErr.Unpack(buf)
		if uErr != nil {
			return err.Error()
		}

		values, ok := errData.([]interface{})
		if !ok {
			values := make([]string, len(cErr.Inputs))
			for i := range values {
				values[i] = "?" // unknown value
			}
		}

		params := make([]string, len(cErr.Inputs))
		for i, input := range cErr.Inputs {
			if input.Name == "" {
				input.Name = fmt.Sprintf("arg%d", i)
			}
			params[i] = fmt.Sprintf("%s=%v", input.Name, values[i])
		}
		return fmt.Sprintf("%s(%s)", cErr.Name, strings.Join(params, ", "))
	}

	return fmt.Sprintf("unknown ABI error: %s", errStr)
}
