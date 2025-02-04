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

// EVMClient is an interface for interacting with an Ethereum client for event subscription.
type EVMClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

// ProgressStore is an interface for storing the last block number processed by the event listener.
type ProgressStore interface {
	LastBlock() (uint64, error)
	SetLastBlock(block uint64) error
}

type Subscriber interface {
	PublishLogEvent(ctx context.Context, log types.Log)
}

type httpPublisher struct {
	progressStore ProgressStore
	logger        *slog.Logger
	evmClient     EVMClient
	subscriber    Subscriber
}

func NewHTTPPublisher(progressStore ProgressStore, logger *slog.Logger, evmClient EVMClient, subscriber Subscriber) *httpPublisher {
	return &httpPublisher{
		progressStore: progressStore,
		logger:        logger,
		evmClient:     evmClient,
		subscriber:    subscriber,
	}
}

func (h *httpPublisher) Start(ctx context.Context, contracts ...common.Address) <-chan struct{} {
	doneChan := make(chan struct{})

	if len(contracts) == 0 {
		h.logger.Error("no contracts to listen to")
		close(doneChan)
		return doneChan
	}

	go func() {
		defer close(doneChan)

		lastBlock, err := h.progressStore.LastBlock()
		if err != nil {
			h.logger.Error("failed to get last block", "error", err)
			return
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blockNumber, err := h.evmClient.BlockNumber(ctx)
				if err != nil {
					h.logger.Warn("failed to get block number", "error", err)
					continue
				}

				if blockNumber > lastBlock {
					startBlock := lastBlock + 1
					success := true

					for startBlock <= blockNumber {
						endBlock := startBlock + 5000 - 1
						if endBlock > blockNumber {
							endBlock = blockNumber
						}

						q := ethereum.FilterQuery{
							FromBlock: big.NewInt(int64(startBlock)),
							ToBlock:   big.NewInt(int64(endBlock)),
							Addresses: contracts,
						}

						logs, err := h.evmClient.FilterLogs(ctx, q)
						if err != nil {
							h.logger.Warn("failed to filter logs", "error", err)
							success = false
							break
						}

						for _, logMsg := range logs {
							h.subscriber.PublishLogEvent(ctx, logMsg)
						}

						h.logger.Debug("processed logs", "from", startBlock, "to", endBlock, "count", len(logs))
						startBlock = endBlock + 1
					}

					if success {
						if err := h.progressStore.SetLastBlock(blockNumber); err != nil {
							h.logger.Error("failed to set last block", "error", err)
							return
						}
						lastBlock = blockNumber
					}
				}
			}
		}
	}()

	return doneChan
}
