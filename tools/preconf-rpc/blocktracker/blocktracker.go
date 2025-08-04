package blocktracker

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"sync"
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
	checkCond     *sync.Cond
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
		checkCond:     sync.NewCond(&sync.Mutex{}),
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
					b.triggerCheck()
					b.log.Debug("New block detected", "number", block.NumberU64(), "hash", block.Hash().Hex())
				}
			}
		}
	}()
	return done
}

func (b *blockTracker) triggerCheck() {
	b.checkCond.L.Lock()
	b.checkCond.Broadcast()
	b.checkCond.L.Unlock()
}

func (b *blockTracker) LatestBlockNumber() uint64 {
	return b.latestBlockNo.Load()
}

func (b *blockTracker) NextBlockNumber() (uint64, time.Duration, error) {
	block, found := b.blocks.Get(b.latestBlockNo.Load())
	if !found {
		return 0, 0, errors.New("latest block not found in cache")
	}
	blockTime := time.Unix(int64(block.Time()), 0)
	return b.latestBlockNo.Load() + 1, time.Until(blockTime.Add(12 * time.Second)), nil
}

func (b *blockTracker) CheckTxnInclusion(
	ctx context.Context,
	txHash common.Hash,
	blockNumber uint64,
) (bool, error) {
	if blockNumber <= b.latestBlockNo.Load() {
		return b.checkTxnInclusion(ctx, txHash, blockNumber)
	}

	waitCh := make(chan struct{})
	go func() {
		b.checkCond.L.Lock()
		defer b.checkCond.L.Unlock()
		for blockNumber > b.latestBlockNo.Load() {
			b.checkCond.Wait()
		}
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case <-waitCh:
		return b.checkTxnInclusion(ctx, txHash, blockNumber)
	}
}

func (b *blockTracker) checkTxnInclusion(ctx context.Context, txHash common.Hash, blockNumber uint64) (bool, error) {
	var err error
	block, ok := b.blocks.Get(blockNumber)
	if !ok {
		block, err = b.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			b.log.Error("Failed to get block by number", "error", err, "blockNumber", blockNumber)
			return false, err
		}
		_ = b.blocks.Add(blockNumber, block)
	}

	if txn := block.Transaction(txHash); txn != nil {
		return true, nil
	}

	return false, nil
}
