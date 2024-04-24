package publisher

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

		lastBlock, err := w.progressStore.LastBlock()
		if err != nil {
			w.logger.Error("failed to get last block", "error", err)
			return
		}

		q := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(lastBlock + 1)),
			ToBlock:   nil,
			Addresses: contracts,
		}

		logChan := make(chan types.Log)

		sub, err := w.evmClient.SubscribeFilterLogs(ctx, q, logChan)
		if err != nil {
			w.logger.Error("failed to subscribe to logs", "error", err)
			return
		}

		defer sub.Unsubscribe()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-sub.Err():
				w.logger.Error("subscription error", "error", err)
				return
			case logMsg := <-logChan:
				// process log
				w.subscriber.PublishLogEvent(ctx, logMsg)

				if logMsg.BlockNumber > lastBlock {
					if err := w.progressStore.SetLastBlock(logMsg.BlockNumber); err != nil {
						w.logger.Error("failed to set last block", "error", err)
						return
					}
				}
			}
		}
	}()

	return doneChan
}
