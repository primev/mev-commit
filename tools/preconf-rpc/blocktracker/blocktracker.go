package blocktracker

import (
	"context"
	"log/slog"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error)
}

type blockTracker struct {
	latestBlockNo atomic.Uint64
	blocks        map[uint64]*types.Block
	client        EthClient
	log           *slog.Logger
	checkTrigger  chan struct{}
}

func NewBlockTracker(client EthClient, log *slog.Logger) *blockTracker {
	return &blockTracker{
		latestBlockNo: atomic.Uint64{},
		blocks:        make(map[uint64]*types.Block),
		client:        client,
		log:           log,
		checkTrigger:  make(chan struct{}, 1),
	}
}

func (b *blockTracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer close(done)
		// Simulate block tracking logic
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
					b.blocks[blockNo] = block
					b.latestBlockNo.Store(block.NumberU64())
					b.log.Info("New block detected", "number", block.NumberU64(), "hash", block.Hash().Hex())
					b.checkTrigger <- struct{}{}
				}
			}
		}
	}()
	return done
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

	block, ok := b.blocks[blockNumber]
	if !ok {
		block, err := b.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
		if err != nil {
			b.log.Error("Failed to get block by number", "error", err, "blockNumber", blockNumber)
			return false, err
		}
		b.blocks[blockNumber] = block
	}

	for _, tx := range block.Transactions() {
		if tx.Hash().Cmp(txHash) == 0 {
			return true, nil
		}
	}
	return false, nil
}
