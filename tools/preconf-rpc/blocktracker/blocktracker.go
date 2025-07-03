package blocktracker

import (
	"context"
	"log/slog"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error)
}

type blockTracker struct {
	latestBlockNo atomic.Uint64
	blocks        *lru.Cache[uint64, *types.Block]
	client        EthClient
	log           *slog.Logger
	checkTrigger  chan struct{}
}

func NewBlockTracker(client EthClient, log *slog.Logger) (*blockTracker, error) {
	cache, err := lru.New[uint64, *types.Block](1000)
	if err != nil {
		log.Error("Failed to create LRU cache", "error", err)
		return nil, err
	}
	return &blockTracker{
		latestBlockNo: atomic.Uint64{},
		blocks:        cache,
		client:        client,
		log:           log,
		checkTrigger:  make(chan struct{}, 1),
	}, nil
}

func (b *blockTracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blockNo, err := b.client.BlockNumber(ctx)
				if err != nil {
					b.log.Error("Failed to get block number", "error", err)
					continue
				}
				if blockNo > b.latestBlockNo.Load() {
					block, err := b.client.BlockByNumber(ctx, big.NewInt(int64(blockNo)))
					if err != nil {
						b.log.Error("Failed to get block by number", "error", err)
						continue
					}
					_ = b.blocks.Add(blockNo, block)
					b.latestBlockNo.Store(block.NumberU64())
					b.log.Info("New block detected", "number", block.NumberU64(), "hash", block.Hash().Hex())
					b.triggerCheck()
				}
			}
		}
	}()
	return done
}

func (b *blockTracker) triggerCheck() {
	select {
	case b.checkTrigger <- struct{}{}:
	default:
		// Non-blocking send, if channel is full, we skip
	}
}

func (b *blockTracker) LatestBlockNumber() uint64 {
	return b.latestBlockNo.Load()
}

func (b *blockTracker) CheckTxnInclusion(
	ctx context.Context,
	txHash common.Hash,
	blockNumber uint64,
) (bool, error) {
WaitForBlock:
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-b.checkTrigger:
			if blockNumber <= b.latestBlockNo.Load() {
				break WaitForBlock
			}
		}
	}

	block, ok := b.blocks.Get(blockNumber)
	if !ok {
		block, err := b.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			b.log.Error("Failed to get block by number", "error", err, "blockNumber", blockNumber)
			return false, err
		}
		_ = b.blocks.Add(blockNumber, block)
	}

	for _, tx := range block.Transactions() {
		if tx.Hash().Cmp(txHash) == 0 {
			return true, nil
		}
	}
	return false, nil
}
