package notifications

import (
	"github.com/cskr/pubsub/v2"
)

type Topic string

const (
	TopicPeerConnected    Topic = "peer_connected"
	TopicPeerDisconnected Topic = "peer_disconnected"
	TopicValidatorOptedIn Topic = "validator_opted_in"
)

func IsTopicValid(topic Topic) bool {
	switch topic {
	case TopicPeerConnected, TopicPeerDisconnected:
		return true
	default:
		return false
	}
}

// Notification is a struct that represents a notification. It has a Topic field
// that represents the topic of the notification and a Value field that represents
// the value of the notification. The Value field is a map[string]any, which means
// that it can be any type.
type Notification struct {
	topic Topic
	value map[string]any
}

func (n *Notification) Topic() Topic {
	return n.topic
}

func (n *Notification) Value() map[string]any {
	return n.value
}

func NewNotification(topic Topic, value map[string]any) *Notification {
	return &Notification{
		topic: topic,
		value: value,
	}
}

// Notifier is an interface that is used to notify about a notification. It will be
// used by the publisher to notify the subscribers about the notification.
type Notifier interface {
	Notify(*Notification)
}

// Notifiee is an interface that is used to subscribe to a notification. It will be
// used by the subscriber to subscribe to the notification.
type Notifiee interface {
	Subscribe(topics ...Topic) chan *Notification
	Unsubscribe(chan *Notification) <-chan struct{}
}

// Notifications is the implementation of the Notifier and Notifiee interfaces. It
// uses the pubsub package to implement the Notifier and Notifiee interfaces.
type Notifications struct {
	ps *pubsub.PubSub[Topic, *Notification]
}

func New(bufferCapacity int) *Notifications {
	return &Notifications{
		ps: pubsub.New[Topic, *Notification](bufferCapacity),
	}
}

func (n *Notifications) Notify(notification *Notification) {
	n.ps.Pub(notification, notification.Topic())
}

func (n *Notifications) Subscribe(topics ...Topic) chan *Notification {
	return n.ps.Sub(topics...)
}

func (n *Notifications) Unsubscribe(ch chan *Notification) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		n.ps.Unsub(ch)
	}()

	return done
}

func (n *Notifications) Shutdown() {
	n.ps.Shutdown()
}
