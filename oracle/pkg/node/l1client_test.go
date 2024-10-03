package node

import (
	"bytes"
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
	_ l1Listener.EthClient = (*ethClientMock)(nil)
)

type ethClientMock struct {
	blockNumberFn    func(context.Context) (uint64, error)
	blockByNumberFn  func(context.Context, *big.Int) (*types.Block, error)
	headerByNumberFn func(context.Context, *big.Int) (*types.Header, error)
}

func (m *ethClientMock) BlockNumber(ctx context.Context) (uint64, error) {
	return m.blockNumberFn(ctx)
}

func (m *ethClientMock) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return m.blockByNumberFn(ctx, number)
}

func (m *ethClientMock) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return m.headerByNumberFn(ctx, number)
}

func TestL1Client(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	t.Run("BlockNumber", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = uint64(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockNumberFn: func(context.Context) (uint64, error) {
						haveCalls++
						return wantNumber, nil
					},
				},
			}},
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

	t.Run("BlockNumber Drift", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = uint64(42)
			wantDrift  = uint64(10)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockNumberFn: func(context.Context) (uint64, error) {
						haveCalls++
						return wantNumber, nil
					},
				},
			}},
			blockNumberDrift: int(wantDrift),
			maxRetries:       1,
		}

		haveNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			t.Errorf("BlockNumber(...): unexpected error: %v", err)
		}
		if haveNumber != wantNumber-wantDrift {
			t.Errorf("BlockNumber(...):\nhave number: %d\nwant number: %d", haveNumber, wantNumber+wantDrift)
		}
		if haveCalls != wantCalls {
			t.Errorf("BlockNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("BlockNumber Retry Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("block number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockNumberFn: func(context.Context) (uint64, error) {
						haveCalls++
						return 0, mockErr
					},
				},
			}},
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

	t.Run("BlockByNumber", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = big.NewInt(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockByNumberFn: func(context.Context, *big.Int) (*types.Block, error) {
						haveCalls++
						return types.NewBlock(&types.Header{Number: wantNumber}, nil, nil, nil, nil), nil
					},
				},
			}},
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

	t.Run("BlockByNumber Retry Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("block by number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockByNumberFn: func(context.Context, *big.Int) (*types.Block, error) {
						haveCalls++
						return nil, mockErr
					},
				},
			}},
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

	t.Run("HeaderByNumber", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = big.NewInt(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					headerByNumberFn: func(context.Context, *big.Int) (*types.Header, error) {
						haveCalls++
						return &types.Header{Number: wantNumber}, nil
					},
				},
			}},
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

	t.Run("HeaderByNumber Winner Override", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = big.NewInt(42)
			wantWinner = "winner"
			wantCalls  = 1
			haveCalls  = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					headerByNumberFn: func(context.Context, *big.Int) (*types.Header, error) {
						haveCalls++
						return &types.Header{Number: wantNumber}, nil
					},
				},
			}},
			winnersOverride: []string{wantWinner},
			maxRetries:      1,
		}

		haveHeader, err := client.HeaderByNumber(context.Background(), wantNumber)
		if err != nil {
			t.Errorf("HeaderByNumber(...): unexpected error: %v", err)
		}
		if haveHeader.Number.Cmp(wantNumber) != 0 {
			t.Errorf("HeaderByNumber(...):\nhave number: %v\nwant number: %v", haveHeader.Number, wantNumber)
		}
		if !bytes.Equal(haveHeader.Extra, []byte(wantWinner)) {
			t.Errorf("HeaderByNumber(...):\nhave winner: %s\nwant winner: %s", haveHeader.Extra, wantWinner)
		}
		if haveCalls != wantCalls {
			t.Errorf("HeaderByNumber(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("HeaderByNumber Retry Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("header by number error")
			wantCalls = 3
			haveCalls = 0
		)
		client := l1Client{
			logger: logger,
			clients: []struct {
				url string
				cli l1Listener.EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					headerByNumberFn: func(context.Context, *big.Int) (*types.Header, error) {
						haveCalls++
						return nil, mockErr
					},
				},
			}},
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
