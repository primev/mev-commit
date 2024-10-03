package node

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/oracle/pkg/l1Listener"
)

var _ l1Listener.EthClient = (*l1Client)(nil)

// errRetry is returned when retry maxRetries is exhausted.
var errRetry = errors.New("retry attempts exhausted")

// l1ClientOptions is a functional option for l1Client.
type l1ClientOptions func(*l1Client)

// l1ClientWithBlockNumberDrift sets the block number drift
// for the BlockNumber method.
func l1ClientWithBlockNumberDrift(drift int) l1ClientOptions {
	return func(c *l1Client) { c.blockNumberDrift = drift }
}

// l1ClientWithWinnersOverride randomly sets the winner in
// the block header extra data for the HeaderByNumber method.
func l1ClientWithWinnersOverride(winners []string) l1ClientOptions {
	return func(c *l1Client) { c.winnersOverride = winners }
}

// l1ClientWithMaxRetries sets the maximum number
// of retries on error for all L1 connections.
func l1ClientWithMaxRetries(retries int) l1ClientOptions {
	return func(c *l1Client) { c.maxRetries = retries }
}

// newL1Client creates a new L1 client with the given RPC URLs.
func newL1Client(logger *slog.Logger, l1RpcUrls []string, opts ...l1ClientOptions) (*l1Client, error) {
	c := &l1Client{logger: logger}

	var errs error
	for _, url := range l1RpcUrls {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cli, err := ethclient.DialContext(ctx, url)
		cancel()
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("failed to dial L1 RPC URL %s: %w", url, err))
			continue
		}
		c.clients = append(
			c.clients,
			struct {
				url string
				cli l1Listener.EthClient
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
	return c, nil
}

// l1Client is an Ethereum client that can connect to multiple L1 nodes.
// If an operation fails, it will automatically try the next L1 node in the list.
// When all L1 nodes are exhausted, it will retry the operation up to maxRetries times.
type l1Client struct {
	logger  *slog.Logger
	clients []struct {
		url string
		cli l1Listener.EthClient
	}

	// Options.
	blockNumberDrift int
	winnersOverride  []string
	maxRetries       int
}

// BlockNumber returns the latest block number.
func (c *l1Client) BlockNumber(ctx context.Context) (uint64, error) {
	var errs error
	for i := range c.maxRetries {
		for _, l1 := range c.clients {
			switch bn, err := l1.cli.BlockNumber(ctx); {
			case err == nil:
				c.logger.Debug("get block number succeeded", "attempt", i, "block", bn)
				return bn - uint64(c.blockNumberDrift), nil
			case errors.Is(err, ethereum.NotFound):
				return 0, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get block number from %s: %w", l1.url, err))
				c.logger.Warn("get block number failed", "url", l1.url, "attempt", i, "error", err)
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
func (c *l1Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var errs error
	for i := range c.maxRetries {
		for _, l1 := range c.clients {
			switch b, err := l1.cli.BlockByNumber(ctx, number); {
			case err == nil:
				c.logger.Debug("get block by number succeeded", "attempt", i, "block", b)
				return b, nil
			case errors.Is(err, ethereum.NotFound):
				return nil, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get block by number from %s: %w", l1.url, err))
				c.logger.Warn("get block by number failed", "url", l1.url, "attempt", i, "error", err)
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
func (c *l1Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var errs error
	for i := range c.maxRetries {
		for _, l1 := range c.clients {
			switch h, err := l1.cli.HeaderByNumber(ctx, number); {
			case err == nil:
				c.logger.Debug("get header by number succeeded", "attempt", i, "header", h)
				if len(c.winnersOverride) > 0 {
					h.Extra = []byte(c.winnersOverride[rand.IntN(len(c.winnersOverride))])
				}
				return h, nil
			case errors.Is(err, ethereum.NotFound):
				return nil, err
			default:
				errs = errors.Join(errs, fmt.Errorf("get header by number from %s: %w", l1.url, err))
				c.logger.Warn("get header by number failed", "url", l1.url, "attempt", i, "error", err)
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
