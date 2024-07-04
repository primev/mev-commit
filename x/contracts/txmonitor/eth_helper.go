package txmonitor

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// StructLog represents a single operation (op code) executed during
// a transaction. It is part of the execution trace and provides detailed
// information about each step of the transaction execution.
type StructLog struct {
	// Pc indicates the program counter at the point of execution.
	Pc uint64 `json:"pc,omitempty"`
	// Op is the name of the opcode that was executed.
	Op string `json:"op,omitempty"`
	// Gas represents the gas available at this step of execution.
	Gas uint64 `json:"gas,omitempty"`
	// GasCost is the amount of gas used by this operation.
	GasCost uint64 `json:"gasCost,omitempty"`
	// Depth indicates the call depth of the current operation.
	Depth int `json:"depth,omitempty"`
	// Error reports any exceptions or errors encountered during the
	// execution of the corresponding operation (op code) in the EVM.
	Error string `json:"error,omitempty"`
	// Stack is the state of the EVM stack at this point of execution.
	Stack []string `json:"stack,omitempty"`
	// Memory represents the contents of the EVM memory at this point.
	Memory []string `json:"memory,omitempty"`
	// Storage is a map representing the contract storage state after this
	// operation. The keys are storage slots and the values are the data
	// stored at those slots.
	Storage map[string]string `json:"storage,omitempty"`
}

// TransactionTrace encapsulates the result of tracing an Ethereum transaction,
// providing a high-level overview of the transaction's execution and its
// effects on the EVM state.
type TransactionTrace struct {
	// Failed indicates whether the transaction was successful.
	Failed bool `json:"failed,omitempty"`
	// Gas indicates the total gas used by the transaction.
	Gas uint64 `json:"gas,omitempty"`
	// ReturnValue is the data returned by the transaction if it was a call
	// to a contract.
	ReturnValue string `json:"returnValue,omitempty"`
	// StructLogs is a slice of StructLog entries, each representing a step
	// in the transaction's execution. These logs provide a detailed trace
	// of the EVM operations, allowing for deep inspection of transaction
	// behavior.
	StructLogs []StructLog `json:"structLogs,omitempty"`
}

// Debugger defines an interface for EVM debugging tools.
type Debugger interface {
	// TraceTransaction takes a transaction hash and returns
	// a detailed trace of the transaction's execution.
	TraceTransaction(ctx context.Context, txHash common.Hash) (*TransactionTrace, error)
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
	client *rpc.Client
}

// NewEVMHelper creates a new EVMHelper instance.
func NewEVMHelper(client *rpc.Client) *evmHelper {
	return &evmHelper{client}
}

// TraceTransaction implements Debugger.TraceTransaction interface.
func (e *evmHelper) TraceTransaction(ctx context.Context, txHash common.Hash) (*TransactionTrace, error) {
	result := new(TransactionTrace)
	traceOptions := map[string]interface{}{} // Empty map for default options.
	if err := e.client.CallContext(
		ctx,
		result,
		"debug_traceTransaction",
		txHash,
		traceOptions,
	); err != nil {
		return nil, err
	}
	return result, nil
}

// BatchReceipts retrieves multiple receipts for a list of transaction hashes.
func (e *evmHelper) BatchReceipts(ctx context.Context, txHashes []common.Hash) ([]Result, error) {
	batch := make([]rpc.BatchElem, len(txHashes))

	for i, hash := range txHashes {
		batch[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{hash},
			Result: new(types.Receipt),
		}
	}

	var receipts []Result
	var err error
	for attempts := 0; attempts < 50; attempts++ {
		// Execute the batch request
		err = e.client.BatchCallContext(context.Background(), batch)
		if err != nil {
			log.Printf("Batch call attempt %d failed: %v", attempts+1, err)
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}
	receipts = make([]Result, len(batch))
	for i, elem := range batch {
		receipts[i].Receipt = elem.Result.(*types.Receipt)
		if elem.Error != nil {
			receipts[i].Err = elem.Error
		}
	}

	return receipts, nil
}
