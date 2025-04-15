package notifier

import (
	"context"
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
	notifySecondsAhead   time.Duration
	lastNotifiedBlockNum uint64
	mu                   sync.Mutex
}

func NewFullNotifier(
	logger *slog.Logger,
	targetBlockNumChan chan uint64,
	l1Client L1Client,
	notifySecondsAhead time.Duration,
) *FullNotifier {
	return &FullNotifier{
		logger:             logger,
		targetBlockNumChan: targetBlockNumChan,
		l1Client:           l1Client,
		notifySecondsAhead: notifySecondsAhead,
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
	currentBlockNum := header.Number.Uint64()
	nextBlockNum := currentBlockNum + 1

	now := time.Now()
	nextBlockTime := now.Add(BlockDuration)
	notificationTime := nextBlockTime.Add(-b.notifySecondsAhead)

	if notificationTime.Before(now) {
		b.sendTargetBlockNotification(nextBlockNum)
		return nil
	}

	delay := notificationTime.Sub(now)

	b.logger.Debug("scheduling notification",
		"target_block_number", nextBlockNum,
		"delay_seconds", delay.Seconds(),
		"expected_block_time", nextBlockTime)

	b.scheduleNotification(ctx, nextBlockNum, delay)
	return nil
}

func (b *FullNotifier) scheduleNotification(ctx context.Context, blockNum uint64, delay time.Duration) {
	go func() {
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
			b.sendTargetBlockNotification(blockNum)
		}
	}()
}

func (b *FullNotifier) sendTargetBlockNotification(targetBlockNum uint64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if targetBlockNum <= b.lastNotifiedBlockNum {
		return
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
}
