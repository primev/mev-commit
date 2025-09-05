package publisher

import (
	"context"
	"log/slog"
	"math/big"
	"slices"
	"sync"
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

	mu        sync.RWMutex
	contracts []common.Address
	updateCh  chan struct{}
}

func NewWSPublisher(progressStore ProgressStore, logger *slog.Logger, evmClient WSEVMClient, subscriber Subscriber) *wsPublisher {
	return &wsPublisher{
		progressStore: progressStore,
		logger:        logger,
		evmClient:     evmClient,
		subscriber:    subscriber,
		contracts:     make([]common.Address, 0),
		updateCh:      make(chan struct{}, 1),
	}
}

func (w *wsPublisher) AddContracts(addr ...common.Address) {
	added := w.addContracts(addr...)
	if added {
		select {
		case w.updateCh <- struct{}{}:
		default:
		}
	}
}

func (w *wsPublisher) addContracts(addr ...common.Address) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	added := false
	for _, a := range addr {
		if !slices.Contains(w.contracts, a) {
			w.contracts = append(w.contracts, a)
			added = true
			w.logger.Info("ws: added contract address", "address", a.Hex())
		}
	}
	return added
}

func (w *wsPublisher) getContracts() []common.Address {
	w.mu.RLock()
	defer w.mu.RUnlock()
	cp := make([]common.Address, len(w.contracts))
	copy(cp, w.contracts)
	return cp
}

func (w *wsPublisher) Start(ctx context.Context, contractAddr ...common.Address) <-chan struct{} {
	doneChan := make(chan struct{})
	added := w.addContracts(contractAddr...)
	if !added {
		w.logger.Warn("contracts were added before starting the publisher")
	}

	if len(contractAddr) == 0 {
		w.logger.Info("ws: starting with no contracts; waiting for addresses to be added")
	}

	go func() {
		defer close(doneChan)

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

			addresses := w.getContracts()
			if len(addresses) == 0 {
				select {
				case <-ctx.Done():
					return
				case <-w.updateCh:
					continue
				case <-time.After(500 * time.Millisecond):
					continue
				}
			}

			q := ethereum.FilterQuery{
				FromBlock: big.NewInt(int64(lastBlock + 1)),
				ToBlock:   nil,
				Addresses: addresses,
			}

			logChan := make(chan types.Log)

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
				case <-w.updateCh:
					w.logger.Info("ws: contract set updated; resubscribing")
					inactivityStart = time.Now()
					sub.Unsubscribe()
					break PROCESSING
				case err := <-sub.Err():
					// retry after 5 seconds
					w.logger.Warn("subscription error", "error", err)
					inactivityStart = time.Now()
					time.Sleep(5 * time.Second)
					break PROCESSING
				case logMsg := <-logChan:
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
