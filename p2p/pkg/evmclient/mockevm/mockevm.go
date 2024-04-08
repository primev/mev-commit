package mockevm

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
)

type mockEvm struct {
	networkID              *big.Int
	batcherFunc            func() evmclient.Batcher
	blockNumFunc           func(ctx context.Context) (uint64, error)
	pendingNonceAtFunc     func(ctx context.Context, account common.Address) (uint64, error)
	nonceAtFunc            func(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	suggestGasPriceFunc    func(ctx context.Context) (*big.Int, error)
	suggestGasTipCapFunc   func(ctx context.Context) (*big.Int, error)
	estimateGasFunc        func(ctx context.Context, call ethereum.CallMsg) (uint64, error)
	sendTransactionFunc    func(ctx context.Context, tx *types.Transaction) error
	callContractFunc       func(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	transactionReceiptFunc func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	transactionByHasFunc   func(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)
}

type Option func(*mockEvm)

type batchFunc func(context.Context, []rpc.BatchElem) error

func (b batchFunc) BatchCallContext(ctx context.Context, batch []rpc.BatchElem) error {
	return b(ctx, batch)
}

func WithBatcherFunc(batcherFunc func(context.Context, []rpc.BatchElem) error) Option {
	return func(m *mockEvm) {
		m.batcherFunc = func() evmclient.Batcher {
			return batchFunc(batcherFunc)
		}
	}
}

func WithBlockNumFunc(blockNumFunc func(ctx context.Context) (uint64, error)) Option {
	return func(m *mockEvm) {
		m.blockNumFunc = blockNumFunc
	}
}

func WithPendingNonceAtFunc(pendingNonceAtFunc func(ctx context.Context, account common.Address) (uint64, error)) Option {
	return func(m *mockEvm) {
		m.pendingNonceAtFunc = pendingNonceAtFunc
	}
}

func WithNonceAtFunc(nonceAtFunc func(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)) Option {
	return func(m *mockEvm) {
		m.nonceAtFunc = nonceAtFunc
	}
}

func WithSuggestGasPriceFunc(suggestGasPriceFunc func(ctx context.Context) (*big.Int, error)) Option {
	return func(m *mockEvm) {
		m.suggestGasPriceFunc = suggestGasPriceFunc
	}
}

func WithSuggestGasTipCapFunc(suggestGasTipCapFunc func(ctx context.Context) (*big.Int, error)) Option {
	return func(m *mockEvm) {
		m.suggestGasTipCapFunc = suggestGasTipCapFunc
	}
}

func WithEstimateGasFunc(estimateGasFunc func(ctx context.Context, call ethereum.CallMsg) (uint64, error)) Option {
	return func(m *mockEvm) {
		m.estimateGasFunc = estimateGasFunc
	}
}

func WithSendTransactionFunc(sendTransactionFunc func(ctx context.Context, tx *types.Transaction) error) Option {
	return func(m *mockEvm) {
		m.sendTransactionFunc = sendTransactionFunc
	}
}

func WithCallContractFunc(callContractFunc func(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)) Option {
	return func(m *mockEvm) {
		m.callContractFunc = callContractFunc
	}
}

func WithTransactionReceiptFunc(transactionReceiptFunc func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)) Option {
	return func(m *mockEvm) {
		m.transactionReceiptFunc = transactionReceiptFunc
	}
}

func WithTransactionByHashFunc(transactionByHashFunc func(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error)) Option {
	return func(m *mockEvm) {
		m.transactionByHasFunc = transactionByHashFunc
	}
}

func NewMockEvm(networkID uint64, opts ...Option) *mockEvm {
	m := &mockEvm{}
	for _, opt := range opts {
		opt(m)
	}
	m.networkID = new(big.Int).SetUint64(networkID)
	return m
}

var ErrNotImplemented = errors.New("not implemented")

func (m *mockEvm) Batcher() evmclient.Batcher {
	if m.batcherFunc != nil {
		return m.batcherFunc()
	}
	return nil
}

func (m *mockEvm) NetworkID(ctx context.Context) (*big.Int, error) {
	return m.networkID, nil
}

func (m *mockEvm) BlockNumber(ctx context.Context) (uint64, error) {
	if m.blockNumFunc != nil {
		return m.blockNumFunc(ctx)
	}
	return 0, ErrNotImplemented
}

func (m *mockEvm) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if m.pendingNonceAtFunc != nil {
		return m.pendingNonceAtFunc(ctx, account)
	}
	return 0, ErrNotImplemented
}

func (m *mockEvm) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	if m.nonceAtFunc != nil {
		return m.nonceAtFunc(ctx, account, blockNumber)
	}
	return 0, ErrNotImplemented
}

func (m *mockEvm) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if m.suggestGasPriceFunc != nil {
		return m.suggestGasPriceFunc(ctx)
	}
	return nil, ErrNotImplemented
}

func (m *mockEvm) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	if m.suggestGasTipCapFunc != nil {
		return m.suggestGasTipCapFunc(ctx)
	}
	return nil, ErrNotImplemented
}

func (m *mockEvm) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	if m.estimateGasFunc != nil {
		return m.estimateGasFunc(ctx, call)
	}
	return 0, ErrNotImplemented
}

func (m *mockEvm) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if m.sendTransactionFunc != nil {
		return m.sendTransactionFunc(ctx, tx)
	}
	return ErrNotImplemented
}

func (m *mockEvm) CallContract(
	ctx context.Context,
	call ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	if m.callContractFunc != nil {
		return m.callContractFunc(ctx, call, blockNumber)
	}
	return nil, ErrNotImplemented
}

func (m *mockEvm) TransactionReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*types.Receipt, error) {
	if m.transactionReceiptFunc != nil {
		return m.transactionReceiptFunc(ctx, txHash)
	}
	return nil, ErrNotImplemented
}

func (m *mockEvm) TransactionByHash(
	ctx context.Context,
	txHash common.Hash,
) (*types.Transaction, bool, error) {
	if m.transactionByHasFunc != nil {
		return m.transactionByHasFunc(ctx, txHash)
	}
	return nil, false, ErrNotImplemented
}
