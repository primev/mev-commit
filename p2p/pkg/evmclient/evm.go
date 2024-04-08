package evmclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

// EVM is an interface for interacting with the Ethereum Virtual Machine. It includes
// subset of the methods from the go-ethereum ethclient.Client interface.
type EVM interface {
	// Client returns the underlying RPC client.
	Batcher() Batcher
	// NetworkID returns the network ID associated with this client.
	NetworkID(ctx context.Context) (*big.Int, error)
	// BlockNumber returns the most recent block number
	BlockNumber(ctx context.Context) (uint64, error)
	// PendingNonceAt retrieves the current pending nonce associated with an account.
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	// NonceAt retrieves the current nonce associated with an account.
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
	// execution of a transaction.
	// Note after eip 1559 this returns suggested priority fee per gas + suggested base fee per gas.
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	// SuggestGasTipCap retrieves the currently suggested 1559 priority fee to allow
	// a timely execution of a transaction.
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	// EstimateGas tries to estimate the gas needed to execute a specific
	// transaction based on the current pending state of the backend blockchain.
	// There is no guarantee that this is the true gas limit requirement as other
	// transactions may be added or removed by miners, but it should provide a basis
	// for setting a reasonable default.
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	// SendTransaction injects the transaction into the pending pool for execution.
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	// CallContract executes an Ethereum contract call with the specified data as the
	// input.
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// TransactionByHash checks the pool of pending transactions in addition to the
	// blockchain. The isPending return value indicates whether the transaction has been
	// mined yet. Note that the transaction may not be part of the canonical chain even if
	// it's not pending.
	TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

type Batcher interface {
	// BatchCallContext executes multiple Ethereum contract calls concurrently with the
	// specified data as the input.
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type evm struct {
	*ethclient.Client
}

func (e *evm) Batcher() Batcher {
	return e.Client.Client()
}

// TraceTransaction implements Debugger.TraceTransaction interface.
func (e *evm) TraceTransaction(ctx context.Context, txHash common.Hash) (*TransactionTrace, error) {
	result := new(TransactionTrace)
	traceOptions := map[string]interface{}{} // Empty map for default options.
	if err := e.Client.Client().CallContext(
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

func WrapEthClient(client *ethclient.Client) EVM {
	return &evm{
		Client: client,
	}
}
