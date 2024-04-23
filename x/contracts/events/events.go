package events

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventManager is an interface for subscribing to events. This interface is a stand-in for
// the generic event handlers that are used to subscribe to events. Currently the only
// implementation of this interface is in this package which is why the methods are
// unexported.
type EventHandler interface {
	eventName() string
	handle(types.Log) error
	setTopicAndContract(topic common.Hash, contract *abi.ABI)
}

// eventHandler is a generic implementation of EventHandler for type-safe event handling.
type eventHandler[T any] struct {
	handler  func(*T)
	name     string
	topicID  common.Hash
	contract *abi.ABI
}

// NewEventHandler creates a new EventHandler for the given event name from the known contracts.
// The handler function is called when an event is received. The event
// handler should be used to subscribe to events using the EventManager.
func NewEventHandler[T any](name string, handler func(*T)) EventHandler {
	return &eventHandler[T]{
		handler: handler,
		name:    name,
	}
}

func (h *eventHandler[T]) eventName() string {
	return h.name
}

func (h *eventHandler[T]) setTopicAndContract(topic common.Hash, contract *abi.ABI) {
	h.topicID = topic
	h.contract = contract
}

func (h *eventHandler[T]) handle(log types.Log) error {
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

	h.handler(obj)
	return nil
}

// EventManager is an interface for subscribing to events. The EventHandler callback
// is called when an event is received. The Subscription returned by the Subscribe
// method can be used to unsubscribe from the event.
type EventManager interface {
	Subscribe(event EventHandler) (Subscription, error)
	PublishLogEvent(ctx context.Context, log types.Log)
}

// Subscription is an interface for unsubscribing from an event. The Unsubscribe method
// should be called to stop receiving events. The Err method returns a channel that
// will receive an error if there was any error in parsing the event. This would only
// happen if the event handler was created with an incorrect ABI. If the error channel
// is not read from, future errors will be dropped.
type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

type Listener struct {
	logger      *slog.Logger
	subMu       sync.RWMutex
	subscribers map[common.Hash][]*subscription
	contracts   []*abi.ABI
}

func NewListener(
	logger *slog.Logger,
	contracts ...*abi.ABI,
) *Listener {
	return &Listener{
		logger:      logger,
		subscribers: make(map[common.Hash][]*subscription),
		contracts:   contracts,
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
	var topic common.Hash
	for _, c := range l.contracts {
		for _, e := range c.Events {
			if e.Name == event.eventName() {
				event.setTopicAndContract(e.ID, c)
				topic = e.ID
				break
			}
		}
	}

	if topic == (common.Hash{}) {
		return nil, fmt.Errorf("event not found")
	}

	l.subMu.Lock()
	defer l.subMu.Unlock()

	sub := &subscription{
		event: event,
		errCh: make(chan error, 1),
		unsub: func() { l.unsubscribe(topic, event) },
	}

	l.subscribers[topic] = append(l.subscribers[topic], sub)

	return sub, nil
}

func (l *Listener) unsubscribe(topic common.Hash, event EventHandler) {
	l.subMu.Lock()
	defer l.subMu.Unlock()

	events := l.subscribers[topic]
	for i, e := range events {
		if e.event == event {
			events = append(events[:i], events[i+1:]...)
			close(e.errCh)
			break
		}
	}

	l.subscribers[topic] = events
}

func (l *Listener) PublishLogEvent(ctx context.Context, log types.Log) {
	l.subMu.RLock()
	defer l.subMu.RUnlock()

	wg := sync.WaitGroup{}
	events := l.subscribers[log.Topics[0]]
	for _, event := range events {
		ev := event
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := ev.event.handle(log); err != nil {
				l.logger.Error("failed to handle log", "error", err)
				select {
				case ev.errCh <- err:
				case <-ctx.Done():
				default:
					l.logger.Error("failed to send error to subscriber", "error", err, "event", ev.event.eventName())
				}
			}
		}()
	}

	wg.Wait()
}
