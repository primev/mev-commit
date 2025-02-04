package publisher

import (
	"context"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	inactivityTimeout = 30 * time.Minute
)

type WSEVMClient interface {
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

type wsPublisher struct {
	progressStore ProgressStore
	logger        *slog.Logger
	evmClient     WSEVMClient
	subscriber    Subscriber
}

func NewWSPublisher(progressStore ProgressStore, logger *slog.Logger, evmClient WSEVMClient, subscriber Subscriber) *wsPublisher {
	return &wsPublisher{
		progressStore: progressStore,
		logger:        logger,
		evmClient:     evmClient,
		subscriber:    subscriber,
	}
}

func (w *wsPublisher) Start(ctx context.Context, contracts ...common.Address) <-chan struct{} {
	doneChan := make(chan struct{})

	if len(contracts) == 0 {
		w.logger.Error("no contracts to listen to")
		close(doneChan)
		return doneChan
	}

	go func() {
		defer close(doneChan)

		w.logger.Info("starting to listen to events")

		lastBlock, err := w.progressStore.LastBlock()
		if err != nil {
			w.logger.Error("failed to get last block", "error", err)
			return
		}

		inactivityStart := time.Now()
		for {
			if time.Since(inactivityStart) > inactivityTimeout {
				w.logger.Error("no activity for 30 minutes, exiting")
				return
			}

			q := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(lastBlock + 1)),
				ToBlock:   nil,
				Addresses: contracts,
			}

			logChan := make(chan types.Log)

			w.logger.Info("subscribing to logs", "query", q)
			sub, err := w.evmClient.SubscribeFilterLogs(ctx, q, logChan)
			if err != nil {
				// retry after 5 seconds
				w.logger.Warn("failed to subscribe to logs", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

		PROCESSING:
			for {
				select {
				case <-ctx.Done():
					sub.Unsubscribe()
					return
				case err := <-sub.Err():
					// retry after 5 seconds
					w.logger.Warn("subscription error", "error", err)
					inactivityStart = time.Now()
					time.Sleep(5 * time.Second)
					break PROCESSING
				case logMsg := <-logChan:
					w.logger.Info("received log", "log", logMsg)
					// process log
					w.subscriber.PublishLogEvent(ctx, logMsg)

					if logMsg.BlockNumber > lastBlock {
						if err := w.progressStore.SetLastBlock(logMsg.BlockNumber); err != nil {
							w.logger.Error("failed to set last block", "error", err)
							return
						}
						lastBlock = logMsg.BlockNumber
					}
				}
			}
		}
	}()

	return doneChan
}
