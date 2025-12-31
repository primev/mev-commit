package blocktracker

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/sync/errgroup"
)

type EthClient interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
	BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
}

type LatestBlockInfo struct {
	Number      uint64
	Time        int64
	BaseFee     *big.Int
	NextBaseFee *big.Int
}

type blockTracker struct {
	latestBlockInfo atomic.Pointer[LatestBlockInfo]
	blocks          *lru.Cache[uint64, *types.Block]
	client          EthClient
	log             *slog.Logger
	txnToCheckMu    sync.Mutex
	txnsToCheck     map[common.Hash]chan uint64
	newBlockChan    chan uint64
}

func NewBlockTracker(client EthClient, log *slog.Logger) (*blockTracker, error) {
	cache, err := lru.New[uint64, *types.Block](1000)
	if err != nil {
		log.Error("Failed to create LRU cache", "error", err)
		return nil, err
	}
	return &blockTracker{
		blocks:       cache,
		client:       client,
		log:          log,
		txnsToCheck:  make(map[common.Hash]chan uint64),
		newBlockChan: make(chan uint64, 1),
	}, nil
}

func (b *blockTracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		backoff := time.Second
		for {
			if err := egCtx.Err(); err != nil {
				return err
			}

			subCh := make(chan *types.Header, 1)
			sub, err := b.client.SubscribeNewHead(egCtx, subCh)
			if err != nil {
				b.log.Error("Failed to subscribe to new head", "error", err)
				select {
				case <-time.After(backoff):
					if backoff < 5*time.Second {
						backoff *= 2
					}
					continue
				case <-egCtx.Done():
					return egCtx.Err()
				}
			}

			backoff = time.Second
			for {
				select {
				case <-egCtx.Done():
					sub.Unsubscribe()
					return egCtx.Err()
				case err := <-sub.Err():
					if err != nil {
						b.log.Warn("Subscription error, reconnecting", "error", err)
					}
					sub.Unsubscribe()
					goto RESUBSCRIBE
				case header := <-subCh:
					blockNo := header.Number.Uint64()
					block, err := b.client.BlockByNumber(egCtx, big.NewInt(int64(blockNo)))
					if err != nil {
						b.log.Error("Failed to get block by number", "error", err)
						continue
					}
					_ = b.blocks.Add(blockNo, block)
					b.latestBlockInfo.Store(&LatestBlockInfo{
						Number:      block.NumberU64(),
						Time:        int64(block.Time()),
						BaseFee:     copyBigInt(block.BaseFee()),
						NextBaseFee: computeNextBaseFee(block.Header()),
					})
					select {
					case b.newBlockChan <- blockNo:
					case <-egCtx.Done():
						sub.Unsubscribe()
						return egCtx.Err()
					}
					b.log.Debug("New block detected", "number", block.NumberU64(), "hash", block.Hash().Hex())
				}
			}
		RESUBSCRIBE:
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
	if b.latestBlockInfo.Load() == nil {
		return 0
	}
	return b.latestBlockInfo.Load().Number
}

func (b *blockTracker) AccountNonce(
	ctx context.Context,
	account common.Address,
) (uint64, error) {
	return b.client.PendingNonceAt(ctx, account)
}

func (b *blockTracker) LatestBaseFee() *big.Int {
	info := b.latestBlockInfo.Load()
	if info == nil || info.BaseFee == nil {
		return big.NewInt(0)
	}
	return copyBigInt(info.BaseFee)
}

func (b *blockTracker) NextBaseFee() *big.Int {
	info := b.latestBlockInfo.Load()
	if info == nil || info.NextBaseFee == nil {
		return big.NewInt(0)
	}
	return copyBigInt(info.NextBaseFee)
}

func (b *blockTracker) NextBlockNumber() (uint64, time.Duration, error) {
	latestBlockInfo := b.latestBlockInfo.Load()
	if latestBlockInfo == nil {
		return 0, 0, errors.New("no latest block info available")
	}
	blockTime := time.Unix(latestBlockInfo.Time, 0)
	if time.Since(blockTime) >= 12*time.Second {
		return latestBlockInfo.Number + 2, time.Until(blockTime.Add(24 * time.Second)), nil
	}
	return latestBlockInfo.Number + 1, time.Until(blockTime.Add(12 * time.Second)), nil
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

func copyBigInt(v *big.Int) *big.Int {
	if v == nil {
		return nil
	}
	return new(big.Int).Set(v)
}

func computeNextBaseFee(head *types.Header) *big.Int {
	if head == nil || head.BaseFee == nil {
		return big.NewInt(0)
	}

	baseFee := copyBigInt(head.BaseFee)
	targetGas := new(big.Int).Div(new(big.Int).SetUint64(head.GasLimit), big.NewInt(2))
	if targetGas.Sign() == 0 {
		return baseFee
	}
	gasDelta := new(big.Int).Sub(new(big.Int).SetUint64(head.GasUsed), targetGas)

	baseFeeChange := new(big.Int).Mul(baseFee, gasDelta)
	baseFeeChange.Div(baseFeeChange, new(big.Int).Mul(targetGas, big.NewInt(8)))

	nextBaseFee := new(big.Int).Add(baseFee, baseFeeChange)
	if nextBaseFee.Sign() < 0 {
		nextBaseFee.SetInt64(0)
	}
	return nextBaseFee
}
