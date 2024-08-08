package node

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/oracle/pkg/l1Listener"
)

var (
	_ l1Listener.EthClient = (*retryL1Client)(nil)
	_ l1Listener.EthClient = (*MockEthClient)(nil)
)

type MockEthClient struct {
	blockNumberFn    func(context.Context) (uint64, error)
	blockByNumberFn  func(context.Context, *big.Int) (*types.Block, error)
	headerByNumberFn func(context.Context, *big.Int) (*types.Header, error)
}

func (m *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	return m.blockNumberFn(ctx)
}

func (m *MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return m.blockByNumberFn(ctx, number)
}

func (m *MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return m.headerByNumberFn(ctx, number)
}

func TestRetryL1Client(t *testing.T) {
	discard := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("BlockNumber Success", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = uint64(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				blockNumberFn: func(context.Context) (uint64, error) {
					haveCalls++
					return wantNumber, nil
				},
			},
			logger:     discard,
			maxRetries: 1,
		}

		haveNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			t.Errorf("BlockNumber(...): unexpected error: %v", err)
		}
		if haveNumber != wantNumber {
			t.Errorf("BlockNumber(...):\nhave number: %d\nwant number: %d", haveNumber, wantNumber)
		}
		if haveCalls != wantCalls {
			t.Errorf("BlockNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("BlockNumber Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("block number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				blockNumberFn: func(context.Context) (uint64, error) {
					haveCalls++
					return 0, mockErr
				},
			},
			logger:     discard,
			maxRetries: wantCalls,
		}

		_, err := client.BlockNumber(context.Background())
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveCalls != wantCalls {
			t.Errorf("BlockNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("BlockByNumber Success", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = big.NewInt(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				blockByNumberFn: func(context.Context, *big.Int) (*types.Block, error) {
					haveCalls++
					return types.NewBlock(&types.Header{Number: wantNumber}, nil, nil, nil, nil), nil
				},
			},
			logger:     discard,
			maxRetries: 1,
		}

		haveBlock, err := client.BlockByNumber(context.Background(), wantNumber)
		if err != nil {
			t.Errorf("BlockByNumber(...): unexpected error: %v", err)
		}
		if haveBlock.Number().Cmp(wantNumber) != 0 {
			t.Errorf("BlockByNumber(...):\nhave number: %v\nwant number: %v", haveBlock.Number(), wantNumber)
		}
		if haveCalls != wantCalls {
			t.Errorf("BlockByNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("BlockByNumber Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("block by number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				blockByNumberFn: func(context.Context, *big.Int) (*types.Block, error) {
					haveCalls++
					return nil, mockErr
				},
			},
			logger:     discard,
			maxRetries: wantCalls,
		}

		_, err := client.BlockByNumber(context.Background(), big.NewInt(42))
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockByNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockByNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveCalls != wantCalls {
			t.Errorf("BlockByNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("HeaderByNumber Success", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = big.NewInt(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				headerByNumberFn: func(context.Context, *big.Int) (*types.Header, error) {
					haveCalls++
					return &types.Header{Number: wantNumber}, nil
				},
			},
			logger:     discard,
			maxRetries: 1,
		}

		haveHeader, err := client.HeaderByNumber(context.Background(), wantNumber)
		if err != nil {
			t.Errorf("HeaderByNumber(...): unexpected error: %v", err)
		}
		if haveHeader.Number.Cmp(wantNumber) != 0 {
			t.Errorf("HeaderByNumber(...):\nhave number: %v\nwant number: %v", haveHeader.Number, wantNumber)
		}
		if haveCalls != wantCalls {
			t.Errorf("HeaderByNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("HeaderByNumber Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("header by number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := &retryL1Client{
			EthClient: &MockEthClient{
				headerByNumberFn: func(context.Context, *big.Int) (*types.Header, error) {
					haveCalls++
					return nil, mockErr
				},
			},
			logger:     discard,
			maxRetries: wantCalls,
		}

		_, err := client.HeaderByNumber(context.Background(), big.NewInt(42))
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("HeaderByNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("HeaderByNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveCalls != wantCalls {
			t.Errorf("HeaderByNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})
}
