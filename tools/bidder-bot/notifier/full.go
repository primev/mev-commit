package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	BlockDuration = 12 * time.Second
)

type FullNotifier struct {
	logger               *slog.Logger
	targetBlockNumChan   chan uint64
	l1Client             L1Client
	lastNotifiedBlockNum uint64
	mu                   sync.Mutex
}

func NewFullNotifier(
	logger *slog.Logger,
	l1Client L1Client,
	targetBlockNumChan chan uint64,
) *FullNotifier {
	return &FullNotifier{
		logger:             logger,
		l1Client:           l1Client,
		targetBlockNumChan: targetBlockNumChan,
	}
}

type L1Client interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
}

func (b *FullNotifier) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		headers := make(chan *types.Header)
		sub, err := b.l1Client.SubscribeNewHead(ctx, headers)
		if err != nil {
			b.logger.Error("failed to subscribe to new heads", "error", err)
			return
		}

		b.logger.Info("subscribed to new block headers")

		for {
			select {
			case <-ctx.Done():
				b.logger.Info("context done")
				return

			case err := <-sub.Err():
				b.logger.Error("subscription error", "error", err)
				sub, err = b.l1Client.SubscribeNewHead(ctx, headers)
				if err != nil {
					b.logger.Error("failed to resubscribe to new heads", "error", err)
					return
				}

			case header := <-headers:
				if err := b.handleHeader(header); err != nil {
					b.logger.Error("error handling header", "error", err)
				}
			}
		}
	}()
	return done
}

func (b *FullNotifier) handleHeader(header *types.Header) error {
	targetBlockNum := header.Number.Uint64() + 1
	b.logger.Debug("scheduling notification from header", "target_block_number", targetBlockNum)

	b.mu.Lock()
	defer b.mu.Unlock()

	if targetBlockNum <= b.lastNotifiedBlockNum {
		return fmt.Errorf("skipping notification for duplicate target block number %d", targetBlockNum)
	}

	select {
	case b.targetBlockNumChan <- targetBlockNum:
		b.logger.Debug("sent target block number", "target_block_number", targetBlockNum)
	default:
		select {
		case drainedTargetBlockNum := <-b.targetBlockNumChan:
			b.logger.Warn("drained buffered target block number", "drained_target_block_number", drainedTargetBlockNum)
		default:
		}
		b.targetBlockNumChan <- targetBlockNum
		b.logger.Warn("sent target block number after draining buffer", "target_block_number", targetBlockNum)
	}

	b.lastNotifiedBlockNum = targetBlockNum
	return nil
}
