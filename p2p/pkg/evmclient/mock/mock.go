package mockevmclient

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
)

type Option func(*mockEvmClient)

func New(opts ...Option) *mockEvmClient {
	m := &mockEvmClient{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithSendFunc(
	f func(ctx context.Context, req *evmclient.TxRequest) (common.Hash, error),
) Option {
	return func(m *mockEvmClient) {
		m.SendFunc = f
	}
}

func WithWaitForReceiptFunc(
	f func(ctx context.Context, txnHash common.Hash) (*types.Receipt, error),
) Option {
	return func(m *mockEvmClient) {
		m.WaitForReceiptFunc = f
	}
}

func WithCallFunc(
	f func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error),
) Option {
	return func(m *mockEvmClient) {
		m.CallFunc = f
	}
}

func WithCancelFunc(
	f func(ctx context.Context, txHash common.Hash) (common.Hash, error),
) Option {
	return func(m *mockEvmClient) {
		m.CancelFunc = f
	}
}

func WithSubscribeFilterLogs(
	f func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error),
) Option {
	return func(m *mockEvmClient) {
		m.SubscribeFilterLogsFunc = f
	}
}

func WithBlockByNumber(
	f func(ctx context.Context, number *big.Int) (*types.Block, error),
) Option {
	return func(m *mockEvmClient) {
		m.BlockByNumberFunc = f
	}
}

func WithBlockNumber(
	f func(ctx context.Context) (uint64, error),
) Option {
	return func(m *mockEvmClient) {
		m.BlockNumberFunc = f
	}
}

func WithFilterLogs(
	f func(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error),
) Option {
	return func(m *mockEvmClient) {
		m.FilterLogsFunc = f
	}
}

type mockEvmClient struct {
	SendFunc                func(ctx context.Context, req *evmclient.TxRequest) (common.Hash, error)
	WaitForReceiptFunc      func(ctx context.Context, txnHash common.Hash) (*types.Receipt, error)
	CallFunc                func(ctx context.Context, req *evmclient.TxRequest) ([]byte, error)
	CancelFunc              func(ctx context.Context, txHash common.Hash) (common.Hash, error)
	SubscribeFilterLogsFunc func(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	BlockByNumberFunc       func(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumberFunc         func(ctx context.Context) (uint64, error)
	FilterLogsFunc          func(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
}

func (m *mockEvmClient) BlockNumber(ctx context.Context) (uint64, error) {
	if m.BlockNumberFunc == nil {
		return 0, errors.New("not implemented")
	}
	return m.BlockNumberFunc(ctx)
}

func (m *mockEvmClient) Send(
	ctx context.Context,
	req *evmclient.TxRequest,
) (common.Hash, error) {
	if m.SendFunc == nil {
		return common.Hash{}, errors.New("not implemented")
	}
	return m.SendFunc(ctx, req)
}

func (m *mockEvmClient) WaitForReceipt(
	ctx context.Context,
	txnHash common.Hash,
) (*types.Receipt, error) {
	if m.WaitForReceiptFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.WaitForReceiptFunc(ctx, txnHash)
}

func (m *mockEvmClient) Call(ctx context.Context, req *evmclient.TxRequest) ([]byte, error) {
	if m.CallFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.CallFunc(ctx, req)
}

func (m *mockEvmClient) CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error) {
	if m.CancelFunc == nil {
		return common.Hash{}, errors.New("not implemented")
	}
	return m.CancelFunc(ctx, txHash)
}

func (m *mockEvmClient) SubscribeFilterLogs(
	ctx context.Context,
	query ethereum.FilterQuery,
	ch chan<- types.Log,
) (ethereum.Subscription, error) {
	if m.SubscribeFilterLogsFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.SubscribeFilterLogsFunc(ctx, query, ch)
}

func (m *mockEvmClient) BlockByNumber(
	ctx context.Context,
	number *big.Int,
) (*types.Block, error) {
	if m.BlockByNumberFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.BlockByNumber(ctx, number)
}

func (m *mockEvmClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if m.FilterLogsFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.FilterLogsFunc(ctx, query)
}
