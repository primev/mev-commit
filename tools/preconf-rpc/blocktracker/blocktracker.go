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
	"golang.org/x/sync/errgroup"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

type blockTracker struct {
	latestBlockNo atomic.Uint64
	blocks        *lru.Cache[uint64, *types.Block]
	client        EthClient
	log           *slog.Logger
	txnToCheckMu  sync.Mutex
	txnsToCheck   map[common.Hash]chan uint64
	newBlockChan  chan uint64
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
		txnsToCheck:   make(map[common.Hash]chan uint64),
		newBlockChan:  make(chan uint64, 1),
	}, nil
}

func (b *blockTracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-ticker.C:
				blockNo, err := b.client.BlockNumber(egCtx)
				if err != nil {
					b.log.Error("Failed to get block number", "error", err)
					continue
				}
				if blockNo > b.latestBlockNo.Load() {
					block, err := b.client.BlockByNumber(egCtx, big.NewInt(int64(blockNo)))
					if err != nil {
						b.log.Error("Failed to get block by number", "error", err)
						continue
					}
					_ = b.blocks.Add(blockNo, block)
					b.latestBlockNo.Store(block.NumberU64())
					select {
					case b.newBlockChan <- blockNo:
					case <-egCtx.Done():
						return egCtx.Err()
					}
					b.log.Debug("New block detected", "number", block.NumberU64(), "hash", block.Hash().Hex())
				}
			}
		}
	})
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case bNo := <-b.newBlockChan:
				block, ok := b.blocks.Get(bNo)
				if !ok {
					b.log.Error("Block not found in cache", "blockNumber", bNo)
					continue
				}
				txnsToClear := make([]common.Hash, 0)
				b.txnToCheckMu.Lock()
				for txHash, resultCh := range b.txnsToCheck {
					if txn := block.Transaction(txHash); txn != nil {
						resultCh <- bNo
						close(resultCh)
						txnsToClear = append(txnsToClear, txHash)
					}
				}
				for _, txHash := range txnsToClear {
					delete(b.txnsToCheck, txHash)
				}
				b.txnToCheckMu.Unlock()
			}
		}
	})

	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			b.log.Error("Block tracker exited with error", "error", err)
		}
	}()

	return done
}

func (b *blockTracker) LatestBlockNumber() uint64 {
	return b.latestBlockNo.Load()
}

func (b *blockTracker) AccountNonce(
	ctx context.Context,
	account common.Address,
) (uint64, error) {
	return b.client.PendingNonceAt(ctx, account)
}

func (b *blockTracker) NextBlockNumber() (uint64, time.Duration, error) {
	latestBlockNo := b.latestBlockNo.Load()
	block, found := b.blocks.Get(latestBlockNo)
	if !found {
		return 0, 0, errors.New("latest block not found in cache")
	}
	blockTime := time.Unix(int64(block.Time()), 0)
	if time.Since(blockTime) >= 12*time.Second {
		return latestBlockNo + 2, time.Until(blockTime.Add(24 * time.Second)), nil
	}
	return latestBlockNo + 1, time.Until(blockTime.Add(12 * time.Second)), nil
}

func (b *blockTracker) WaitForTxnInclusion(
	txHash common.Hash,
) chan uint64 {
	resultCh := make(chan uint64, 1)
	b.txnToCheckMu.Lock()
	b.txnsToCheck[txHash] = resultCh
	b.txnToCheckMu.Unlock()
	return resultCh
}
