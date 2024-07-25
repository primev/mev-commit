package evmclients

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type infiniteRetryL1Client struct {
	EthClient
	logger *slog.Logger
}

func NewInfiniteRetryL1Client(ethClient EthClient, logger *slog.Logger) EthClient {
	return &infiniteRetryL1Client{
		EthClient: ethClient,
		logger:    logger,
	}
}

func (i *infiniteRetryL1Client) BlockNumber(ctx context.Context) (uint64, error) {
	var blkNum uint64
	var err error
	for retries := 50; retries > 0; retries-- {
		blkNum, err = i.EthClient.BlockNumber(ctx)
		if err == nil {
			break
		}
		i.logger.Error("failed to get block number, retrying...", "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return 0, err
	}
	return blkNum, nil
}

func (i *infiniteRetryL1Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var hdr *types.Header
	var err error
	for retries := 50; retries > 0; retries-- {
		hdr, err = i.EthClient.HeaderByNumber(ctx, number)
		if err == nil {
			break
		}
		i.logger.Error("failed to get header by number, retrying...", "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return hdr, nil
}

func (i *infiniteRetryL1Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var blk *types.Block
	var err error
	for retries := 50; retries > 0; retries-- {
		blk, err = i.EthClient.BlockByNumber(ctx, number)
		if err == nil {
			break
		}
		i.logger.Error("failed to get block by number, retrying...", "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}
	return blk, nil
}
