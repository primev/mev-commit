package publisher

import (
	"context"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EVMClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

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

	mu        sync.RWMutex
	contracts []common.Address
}

func NewHTTPPublisher(
	progressStore ProgressStore,
	logger *slog.Logger,
	evmClient EVMClient,
	subscriber Subscriber,
) *httpPublisher {
	return &httpPublisher{
		progressStore: progressStore,
		logger:        logger,
		evmClient:     evmClient,
		subscriber:    subscriber,
		// contracts can be empty initially
		contracts: make([]common.Address, 0),
	}
}

// AddContract appends a new contract address to the set we listen on.
func (h *httpPublisher) AddContract(addr common.Address) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.contracts {
		if c == addr {
			h.logger.Info("contract address already tracked", "address", addr.Hex())
			return
		}
	}

	h.contracts = append(h.contracts, addr)
	h.logger.Info("added new contract address", "address", addr.Hex())
}

// getContracts safely returns a *copy* of the slice to avoid concurrent modification issues.
func (h *httpPublisher) getContracts() []common.Address {
	h.mu.RLock()
	defer h.mu.RUnlock()

	copied := make([]common.Address, len(h.contracts))
	copy(copied, h.contracts)
	return copied
}

// Start begins the main subscription loop
func (h *httpPublisher) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

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

					// fetch current contract list
					addresses := h.getContracts()
					if len(addresses) == 0 {
						// If we have no addresses, just update lastBlock & continue
						if err := h.progressStore.SetLastBlock(blockNumber); err != nil {
							h.logger.Error("failed to set last block", "error", err)
							return
						}
						lastBlock = blockNumber
						continue
					}

					for startBlock <= blockNumber {
						endBlock := startBlock + 5000 - 1
						if endBlock > blockNumber {
							endBlock = blockNumber
						}

						q := ethereum.FilterQuery{
							FromBlock: big.NewInt(int64(startBlock)),
							ToBlock:   big.NewInt(int64(endBlock)),
							Addresses: addresses,
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

						h.logger.Debug("processed logs",
							"from", startBlock,
							"to", endBlock,
							"count", len(logs),
							"addresses", len(addresses))
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
