package events

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

// EventManager is an interface for subscribing to events. This interface is a stand-in for
// the generic event handlers that are used to subscribe to events.
type EventHandler interface {
	EventName() string
	Handle(types.Log) error
	SetTopicAndContract(topic common.Hash, contract *abi.ABI)
	Topic() common.Hash
}

// eventHandler is a generic implementation of EventHandler for type-safe event handling.
type eventHandler[T any] struct {
	handler  func(*T) error
	name     string
	topicID  common.Hash
	contract *abi.ABI
}

// NewEventHandler creates a new EventHandler for the given event name from the known contracts.
// The handler function is called when an event is received. The handler function should
// return an error if the event is a fatal error, otherwise it should return nil. The event
// handler should be used to subscribe to events using the EventManager interface.
func NewEventHandler[T any](name string, handler func(*T) error) EventHandler {
	return &eventHandler[T]{
		handler: handler,
		name:    name,
	}
}

func (h *eventHandler[T]) EventName() string {
	return h.name
}

func (h *eventHandler[T]) SetTopicAndContract(topic common.Hash, contract *abi.ABI) {
	h.topicID = topic
	h.contract = contract
}

func (h *eventHandler[T]) Handle(log types.Log) error {
	if h.contract == nil {
		return fmt.Errorf("contract not set")
	}

	if !bytes.Equal(log.Topics[0].Bytes(), h.topicID.Bytes()) {
		return nil
	}

	obj := new(T)

	if len(log.Data) > 0 {
		err := h.contract.UnpackIntoInterface(obj, h.name, log.Data)
		if err != nil {
			return err
		}
	}

	var indexed abi.Arguments
	for _, arg := range h.contract.Events[h.name].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}

	if len(indexed) > 0 {
		err := abi.ParseTopics(obj, indexed, log.Topics[1:])
		if err != nil {
			return err
		}
	}

	return h.handler(obj)
}

func (h *eventHandler[T]) Topic() common.Hash {
	return h.topicID
}

type EventManager interface {
	Subscribe(event EventHandler) (Subscription, error)
}

type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

type Listener struct {
	logger        *slog.Logger
	evmClient     EVMClient
	progressStore ProgressStore
	subMu         sync.RWMutex
	subscribers   map[common.Hash][]*subscription
	contracts     map[common.Address]*abi.ABI
}

func NewListener(
	logger *slog.Logger,
	evmClient EVMClient,
	progressStore ProgressStore,
	contracts map[common.Address]*abi.ABI,
) *Listener {
	return &Listener{
		logger:        logger,
		evmClient:     evmClient,
		progressStore: progressStore,
		subscribers:   make(map[common.Hash][]*subscription),
		contracts:     contracts,
	}
}

type subscription struct {
	event EventHandler
	unsub func()
	errCh chan error
}

func (s *subscription) Unsubscribe() {
	s.unsub()
}

func (s *subscription) Err() <-chan error {
	return s.errCh
}

func (l *Listener) Subscribe(event EventHandler) (Subscription, error) {
	found := false
	for _, c := range l.contracts {
		for _, e := range c.Events {
			if e.Name == event.EventName() {
				event.SetTopicAndContract(e.ID, c)
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("event not found")
	}

	l.subMu.Lock()
	defer l.subMu.Unlock()

	sub := &subscription{
		event: event,
		errCh: make(chan error),
		unsub: func() { l.unsubscribe(event) },
	}

	l.subscribers[event.Topic()] = append(l.subscribers[event.Topic()], sub)

	return sub, nil
}

func (l *Listener) unsubscribe(event EventHandler) {
	l.subMu.Lock()
	defer l.subMu.Unlock()

	events := l.subscribers[event.Topic()]
	for i, e := range events {
		if e.event == event {
			events = append(events[:i], events[i+1:]...)
			break
		}
	}

	l.subscribers[event.Topic()] = events
}

func (l *Listener) publishLogEvent(ctx context.Context, log types.Log) {
	l.subMu.RLock()
	defer l.subMu.RUnlock()

	wg := sync.WaitGroup{}
	events := l.subscribers[log.Topics[0]]
	for _, event := range events {
		ev := event
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := ev.event.Handle(log); err != nil {
				l.logger.Error("failed to handle log", "error", err)
				select {
				case ev.errCh <- err:
				case <-ctx.Done():
				}
			}
		}()
	}

	wg.Wait()
}

func (l *Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	if len(l.contracts) == 0 {
		close(doneChan)
		return doneChan
	}

	go func() {
		defer close(doneChan)

		lastBlock, err := l.progressStore.LastBlock()
		if err != nil {
			l.logger.Error("failed to get last block", "error", err)
			return
		}

		contracts := make([]common.Address, 0, len(l.contracts))
		for addr := range l.contracts {
			contracts = append(contracts, addr)
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blockNumber, err := l.evmClient.BlockNumber(ctx)
				if err != nil {
					l.logger.Error("failed to get block number", "error", err)
					return
				}

				if blockNumber > lastBlock {
					q := ethereum.FilterQuery{
						FromBlock: big.NewInt(int64(lastBlock + 1)),
						ToBlock:   big.NewInt(int64(blockNumber)),
						Addresses: contracts,
					}

					logs, err := l.evmClient.FilterLogs(ctx, q)
					if err != nil {
						l.logger.Error("failed to filter logs", "error", err)
						return
					}

					for _, logMsg := range logs {
						// process log
						l.publishLogEvent(ctx, logMsg)
					}

					if err := l.progressStore.SetLastBlock(blockNumber); err != nil {
						l.logger.Error("failed to set last block", "error", err)
						return
					}
					l.logger.Debug("processed logs", "from", lastBlock+1, "to", blockNumber, "count", len(logs))
					lastBlock = blockNumber
				}
			}
		}
	}()

	return doneChan
}
