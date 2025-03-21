package ethwrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Methods of the Ethereum client that are overridden.
type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	ChainID(ctx context.Context) (*big.Int, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

// errRetry is returned when retry maxRetries is exhausted.
var errRetry = errors.New("retry attempts exhausted")

// EthClientOptions is a functional option for Client.
type EthClientOptions func(*Client)

// EthClientWithBlockNumberDrift sets the block number drift
// for the BlockNumber method.
func EthClientWithBlockNumberDrift(drift int) EthClientOptions {
	return func(c *Client) { c.blockNumberDrift = drift }
}

// EthClientWithWinnersOverride randomly sets the winner in
// the block header extra data for the HeaderByNumber method.
func EthClientWithWinnersOverride(winners []string) EthClientOptions {
	return func(c *Client) { c.winnersOverride = winners }
}

// EthClientWithMaxRetries sets the maximum number
// of retries on error for all rpc connections.
func EthClientWithMaxRetries(retries int) EthClientOptions {
	return func(c *Client) { c.maxRetries = retries }
}

// EthClientWithBlockNumOverride overrides the block number function.
func EthClientWithBlockNumOverride(fn func(context.Context, EthClient) (uint64, error)) EthClientOptions {
	return func(c *Client) {
		c.blockNumFn = fn
	}
}

// NewClient creates a new ethclient with the given RPC URLs.
func NewClient(logger *slog.Logger, rpcUrls []string, opts ...EthClientOptions) (*Client, error) {
	c := &Client{logger: logger}

	var errs error
	for _, url := range rpcUrls {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cli, err := ethclient.DialContext(ctx, url)
		cancel()
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to dial client RPC URL %s: %w", url, err))
			continue
		}
		c.clients = append(
			c.clients,
			struct {
				url string
				cli EthClient
			}{
				url,
				cli,
			},
		)
	}
	if errs != nil {
		return nil, errs
	}

	for _, opt := range opts {
		opt(c)
	}
	if c.blockNumFn == nil {
		c.blockNumFn = func(ctx context.Context, cli EthClient) (uint64, error) {
			return cli.BlockNumber(ctx)
		}
	}
	return c, nil
}

// Client is an Ethereum client that can connect to multiple RPCs.
// If an operation fails, it will automatically try the next client node in the list.
// When all client nodes are exhausted, it will retry the operation up to maxRetries times.
type Client struct {
	logger  *slog.Logger
	clients []struct {
		url string
		cli EthClient
	}

	// Options.
	blockNumberDrift int
	winnersOverride  []string
	maxRetries       int
	blockNumFn       func(context.Context, EthClient) (uint64, error)
}

// RawClient returns the first raw ethclient.
func (c *Client) RawClient() *ethclient.Client {
	client, ok := c.clients[0].cli.(*ethclient.Client)
	if !ok {
		return nil
	}
	return client
}

// BlockNumber returns the latest block number.
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	var errs error
	for i := range c.maxRetries {
		for _, client := range c.clients {
			switch bn, err := c.blockNumFn(ctx, client.cli); {
			case err == nil:
				c.logger.Debug("get block number succeeded", "attempt", i, "block", bn)
				return bn - uint64(c.blockNumberDrift), nil
			case errors.Is(err, ethereum.NotFound):
				return 0, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get block number from %s: %w", client.url, err))
				c.logger.Warn("get block number failed", "url", client.url, "attempt", i, "error", err)
			}
		}
		if i < c.maxRetries-1 {
			d := time.Duration(1+rand.Int64N(6)) * time.Second
			c.logger.Info("get block number retry", "in", d)
			time.Sleep(d)
		}
	}
	return 0, fmt.Errorf("get block number: %w: %w", errRetry, errs)
}

// BlockByNumber returns the block with the given number.
func (c *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var errs error
	for i := range c.maxRetries {
		for _, client := range c.clients {
			switch b, err := client.cli.BlockByNumber(ctx, number); {
			case err == nil:
				c.logger.Debug("get block by number succeeded", "attempt", i, "block", b)
				return b, nil
			case errors.Is(err, ethereum.NotFound):
				return nil, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get block by number from %s: %w", client.url, err))
				c.logger.Warn("get block by number failed", "url", client.url, "attempt", i, "error", err)
			}
		}
		if i < c.maxRetries-1 {
			d := time.Duration(1+rand.Int64N(6)) * time.Second
			c.logger.Info("get block by number retry", "in", d)
			time.Sleep(d)
		}
	}
	return nil, fmt.Errorf("get block by number: %w: %w", errRetry, errs)
}

// HeaderByNumber returns the block header of the block with the given number.
func (c *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var errs error
	for i := range c.maxRetries {
		for _, client := range c.clients {
			switch h, err := client.cli.HeaderByNumber(ctx, number); {
			case err == nil:
				c.logger.Debug("get header by number succeeded", "attempt", i, "header", h)
				if len(c.winnersOverride) > 0 {
					h.Extra = []byte(c.winnersOverride[rand.IntN(len(c.winnersOverride))])
				}
				return h, nil
			case errors.Is(err, ethereum.NotFound):
				return nil, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get header by number from %s: %w", client.url, err))
				c.logger.Warn("get header by number failed", "url", client.url, "attempt", i, "error", err)
			}
		}
		if i < c.maxRetries-1 {
			d := time.Duration(1+rand.Int64N(6)) * time.Second
			c.logger.Info("get header by number retry", "in", d)
			time.Sleep(d)
		}
	}
	return nil, fmt.Errorf("get header by number: %w: %w", errRetry, errs)
}

// NonceAt returns the nonce of the account at the given block number.
func (c *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	var errs error
	for i := range c.maxRetries {
		for _, client := range c.clients {
			switch nonce, err := client.cli.NonceAt(ctx, account, blockNumber); {
			case err == nil:
				c.logger.Debug("get nonce succeeded", "attempt", i, "nonce", nonce)
				return nonce, nil
			case errors.Is(err, ethereum.NotFound):
				return 0, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get nonce from %s: %w", client.url, err))
				c.logger.Warn("get nonce failed", "url", client.url, "attempt", i, "error", err)
			}
		}
		if i < c.maxRetries-1 {
			d := time.Duration(1+rand.Int64N(6)) * time.Second
			c.logger.Info("get nonce retry", "in", d)
			time.Sleep(d)
		}
	}
	return 0, fmt.Errorf("get nonce: %w: %w", errRetry, errs)
}

// FilterLogs returns the logs that satisfy the given filter query.
func (c *Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var errs error
	for i := range c.maxRetries {
		for _, client := range c.clients {
			switch logs, err := client.cli.FilterLogs(ctx, query); {
			case err == nil:
				c.logger.Debug("filter logs succeeded", "attempt", i, "logs", logs)
				return logs, nil
			case errors.Is(err, ethereum.NotFound):
				return nil, err
			default:
				errs = errors.Join(errs, fmt.Errorf("filter logs from %s: %w", client.url, err))
				c.logger.Warn("filter logs failed", "url", client.url, "attempt", i, "error", err)
			}
		}
		if i < c.maxRetries-1 {
			d := time.Duration(1+rand.Int64N(6)) * time.Second
			c.logger.Info("filter logs retry", "in", d)
			time.Sleep(d)
		}
	}
	return nil, fmt.Errorf("filter logs: %w: %w", errRetry, errs)
}

func (c *Client) ChainID(ctx context.Context) (*big.Int, error) {
	rawClient := c.RawClient()
	if rawClient == nil {
		return nil, fmt.Errorf("no raw client")
	}
	return rawClient.ChainID(ctx)
}

func (c *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	rawClient := c.RawClient()
	if rawClient == nil {
		return 0, fmt.Errorf("no raw client")
	}
	return rawClient.PendingNonceAt(ctx, account)
}
