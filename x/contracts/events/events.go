package events

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
)

// HasRawLog is an optional interface you can implement on your event struct
// so we can store log metadata (e.g., block number) after unpacking.
type HasRawLog interface {
	SetRawLog(types.Log)
}

// EventHandler is a stand-in for the generic event handlers that are used to subscribe
// to events. It is useful to describe the generic event handlers using this interface
// so that they can be referenced in the EventManager.
type EventHandler interface {
	eventName() string
	handle(types.Log) error
	setTopicAndContract(topic common.Hash, contract *abi.ABI)
	topic() common.Hash
}

// eventHandler is a generic implementation of EventHandler for type-safe event handling.
type eventHandler[T any] struct {
	handler  func(*T, uint64)
	name     string
	topicID  common.Hash
	contract *abi.ABI
}

// NewEventHandler creates a new EventHandler for the given event name from the known contracts.
// The handler function is called when an event is received. The event
// handler should be used to subscribe to events using the EventManager.
func NewEventHandler[T any](name string, handler func(*T, uint64)) EventHandler {
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

	// Ensure this log's first topic matches the event's topicID
	if !bytes.Equal(log.Topics[0].Bytes(), h.topicID.Bytes()) {
		return nil
	}

	// Create a new instance of T (your event struct)
	obj := new(T)

	// Unpack non-indexed data
	if len(log.Data) > 0 {
		err := h.contract.UnpackIntoInterface(obj, h.name, log.Data)
		if err != nil {
			return err
		}
	}

	// Unpack indexed fields
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

	// Finally, run the user-provided handler logic
	h.handler(obj, log.BlockNumber)

	return nil
}

func (h *eventHandler[T]) topic() common.Hash {
	return h.topicID
}

// EventManager is an interface for subscribing to contract events. The EventHandler callback
// is called when an event is received. The Subscription returned by the Subscribe
// method can be used to unsubscribe from the event and also to receive any errors
// that occur while parsing the event. The PublishLogEvent method is used to publish
// the log events to the subscribers.
type EventManager interface {
	Subscribe(event ...EventHandler) (Subscription, error)
	PublishLogEvent(ctx context.Context, log types.Log)
}

// Subscription is a reference to the active event subscription. The Unsubscribe method
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
	subscribers map[common.Hash][]*storedEvent
	contracts   []*abi.ABI
	metrics     *metrics
}

func NewListener(
	logger *slog.Logger,
	contracts ...*abi.ABI,
) *Listener {
	return &Listener{
		logger:      logger,
		subscribers: make(map[common.Hash][]*storedEvent),
		contracts:   contracts,
		metrics:     newMetrics(),
	}
}

func (l *Listener) Metrics() []prometheus.Collector {
	return l.metrics.Metrics()
}

type subscription struct {
	unsub func()
	errCh chan error
}

func (s *subscription) Unsubscribe() {
	s.unsub()
}

func (s *subscription) Err() <-chan error {
	return s.errCh
}

type storedEvent struct {
	evt   EventHandler
	errCh chan error
}

func (l *Listener) Subscribe(ev ...EventHandler) (Subscription, error) {
	if len(ev) == 0 {
		return nil, fmt.Errorf("no events provided")
	}

	// Match each event with an ABI in our loaded contracts
	for _, event := range ev {
		found := false
		for _, c := range l.contracts {
			for _, e := range c.Events {
				if e.Name == event.eventName() {
					event.setTopicAndContract(e.ID, c)

					found = true
					break
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("event %s not found", event.eventName())
		}
	}

	l.subMu.Lock()
	defer l.subMu.Unlock()

	errC := make(chan error, len(ev))
	sub := &subscription{
		errCh: errC,
		unsub: func() {
			for _, event := range ev {
				l.unsubscribe(event.topic(), event)
			}
			close(errC)
		},
	}

	// Store the events in our subscriber map
	for _, event := range ev {
		l.subscribers[event.topic()] = append(l.subscribers[event.topic()], &storedEvent{
			evt:   event,
			errCh: sub.errCh,
		})
	}

	return sub, nil
}

func (l *Listener) unsubscribe(topic common.Hash, event EventHandler) {
	l.subMu.Lock()
	defer l.subMu.Unlock()

	events := l.subscribers[topic]
	for i, e := range events {
		if e.evt == event {
			events = append(events[:i], events[i+1:]...)
			break
		}
	}

	l.subscribers[topic] = events
}

func (l *Listener) PublishLogEvent(ctx context.Context, log types.Log) {
	l.subMu.RLock()
	defer l.subMu.RUnlock()

	l.metrics.totalLogs.Inc()

	wg := sync.WaitGroup{}
	events := l.subscribers[log.Topics[0]]
	for _, event := range events {
		ev := event
		l.metrics.totalEvents.Inc()
		l.metrics.eventCounts.WithLabelValues(ev.evt.eventName()).Inc()
		wg.Add(1)
		go func() {
			defer wg.Done()

			defer func(start time.Time) {
				l.metrics.eventHandlerDurations.
					WithLabelValues(ev.evt.eventName()).
					Set(float64(time.Since(start)))
			}(time.Now())

			if err := ev.evt.handle(log); err != nil {
				l.logger.Error("failed to handle log", "error", err)
				select {
				case ev.errCh <- fmt.Errorf("failed to handle event %s: %w", ev.evt.eventName(), err):
				case <-ctx.Done():
				default:
					l.logger.Error("failed to send error to subscriber", "error", err, "event", ev.evt.eventName())
				}
			}
		}()
	}

	wg.Wait()
}
