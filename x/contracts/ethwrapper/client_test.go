package ethwrapper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	_ EthClient = (*ethClientMock)(nil)
)

type ethClientMock struct {
	blockNumberFn    func(context.Context) (uint64, error)
	blockByNumberFn  func(context.Context, *big.Int) (*types.Block, error)
	headerByNumberFn func(context.Context, *big.Int) (*types.Header, error)
	nonceAtFn        func(context.Context, common.Address, *big.Int) (uint64, error)
	filterLogsFn     func(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	chainIDFn        func(context.Context) (*big.Int, error)
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

func (m *ethClientMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return m.nonceAtFn(ctx, account, blockNumber)
}

func (m *ethClientMock) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return m.filterLogsFn(ctx, query)
}

func (m *ethClientMock) ChainID(ctx context.Context) (*big.Int, error) {
	return m.chainIDFn(ctx)
}

func TestClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	t.Run("BlockNumber", func(t *testing.T) {
		t.Parallel()

		var (
			wantNumber = uint64(42)
			wantCalls  = 1
			haveCalls  = 0
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
			blockNumFn: func(c context.Context, cli EthClient) (uint64, error) {
				return cli.BlockNumber(c)
			},
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
			blockNumFn: func(c context.Context, cli EthClient) (uint64, error) {
				return cli.BlockNumber(c)
			},
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
			blockNumFn: func(c context.Context, cli EthClient) (uint64, error) {
				return cli.BlockNumber(c)
			},
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockByNumberFn: func(context.Context, *big.Int) (*types.Block, error) {
						haveCalls++
						body := &types.Body{Transactions: nil, Uncles: nil}
						return types.NewBlock(&types.Header{Number: wantNumber}, body, nil, trie.NewStackTrie(nil)), nil
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
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

	t.Run("NonceAt", func(t *testing.T) {
		t.Parallel()

		var (
			wantAccount = common.Address{}
			wantNumber  = big.NewInt(42)
			wantNonce   = uint64(1)
			wantCalls   = 1
			haveCalls   = 0
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					nonceAtFn: func(context.Context, common.Address, *big.Int) (uint64, error) {
						haveCalls++
						return wantNonce, nil
					},
				},
			}},
			maxRetries: 1,
		}

		haveNonce, err := client.NonceAt(context.Background(), wantAccount, wantNumber)
		if err != nil {
			t.Errorf("NonceAt(...): unexpected error: %v", err)
		}
		if haveNonce != wantNonce {
			t.Errorf("NonceAt(...):\nhave nonce: %d\nwant nonce: %d", haveNonce, wantNonce)
		}
		if haveCalls != wantCalls {
			t.Errorf("NonceAt(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("NonceAt Retry Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("nonce at error")
			wantCalls = 3
			haveCalls = 0
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					nonceAtFn: func(context.Context, common.Address, *big.Int) (uint64, error) {
						haveCalls++
						return 0, mockErr
					},
				},
			}},
			maxRetries: wantCalls,
		}

		_, err := client.NonceAt(context.Background(), common.Address{}, big.NewInt(42))
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("NonceAt(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("NonceAt(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveCalls != wantCalls {
			t.Errorf("NonceAt(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("FilterLogs", func(t *testing.T) {
		t.Parallel()

		var (
			wantQuery = ethereum.FilterQuery{}
			wantLogs  = []types.Log{{}}
			wantCalls = 1
			haveCalls = 0
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					filterLogsFn: func(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
						haveCalls++
						return wantLogs, nil
					},
				},
			}},
			maxRetries: 1,
		}

		haveLogs, err := client.FilterLogs(context.Background(), wantQuery)
		if err != nil {
			t.Errorf("FilterLogs(...): unexpected error: %v", err)
		}
		if !reflect.DeepEqual(haveLogs, wantLogs) {
			t.Errorf("FilterLogs(...):\nhave logs: %v\nwant logs: %v", haveLogs, wantLogs)
		}
		if haveCalls != wantCalls {
			t.Errorf("FilterLogs(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("FilterLogs Retry Error", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr   = errors.New("filter logs error")
			wantCalls = 3
			haveCalls = 0
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					filterLogsFn: func(context.Context, ethereum.FilterQuery) ([]types.Log, error) {
						haveCalls++
						return nil, mockErr
					},
				},
			}},
			maxRetries: wantCalls,
		}

		_, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{})
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("FilterLogs(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("FilterLogs(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveCalls != wantCalls {
			t.Errorf("FilterLogs(...):\nhave calls: %d\nwant calls: %d", haveCalls, wantCalls)
		}
	})

	t.Run("BlockNumber func", func(t *testing.T) {
		t.Parallel()

		var (
			mockErr = errors.New("block number error")
		)
		client := Client{
			logger: logger,
			clients: []struct {
				url string
				cli EthClient
			}{{
				url: "mock",
				cli: &ethClientMock{
					blockNumberFn: func(context.Context) (uint64, error) {
						return uint64(42), nil
					},
				},
			}},
			maxRetries: 1,
			blockNumFn: func(context.Context, EthClient) (uint64, error) {
				return 0, mockErr
			},
		}

		_, err := client.BlockNumber(context.Background())
		if haveErr, wantErr := err, errRetry; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
		if haveErr, wantErr := err, mockErr; !errors.Is(haveErr, wantErr) {
			t.Errorf("BlockNumber(...):\nhave error: %v\nwant error: %v", haveErr, wantErr)
		}
	})
}
