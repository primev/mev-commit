package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
)

type FullNotifier struct {
	logger               *slog.Logger
	targetBlockChan      chan bidder.TargetBlock
	l1Client             L1Client
	lastNotifiedBlockNum uint64
	mu                   sync.Mutex
}

func NewFullNotifier(
	logger *slog.Logger,
	l1Client L1Client,
	targetBlockChan chan bidder.TargetBlock,
) *FullNotifier {
	return &FullNotifier{
		logger:          logger,
		l1Client:        l1Client,
		targetBlockChan: targetBlockChan,
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
				if err := b.handleHeader(ctx, header); err != nil {
					b.logger.Error("error handling header", "error", err)
				}
			}
		}
	}()
	return done
}

func (b *FullNotifier) handleHeader(ctx context.Context, header *types.Header) error {
	targetBlock := bidder.TargetBlock{
		Num:  header.Number.Uint64() + 1,
		Time: time.Unix(int64(header.Time), 0).Add(slotDuration),
	}
	b.logger.Debug("handling header",
		"target_block_number", targetBlock.Num,
		"target_block_time", targetBlock.Time,
	)

	b.mu.Lock()
	defer b.mu.Unlock()

	if targetBlock.Num <= b.lastNotifiedBlockNum {
		return fmt.Errorf("skipping notification for duplicate target block number %d", targetBlock.Num)
	}

	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	select {
	case b.targetBlockChan <- targetBlock:
		b.logger.Debug("sent target block",
			"target_block_number", targetBlock.Num,
			"target_block_time", targetBlock.Time,
		)
	case <-sendCtx.Done():
		return fmt.Errorf("failed to send target block %d", targetBlock.Num)
	}

	b.lastNotifiedBlockNum = targetBlock.Num
	return nil
}
