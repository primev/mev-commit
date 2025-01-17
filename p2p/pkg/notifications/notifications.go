package notifications

import "github.com/cskr/pubsub/v2"

const (
	TopicPeerConnected    = "peer_connected"
	TopicPeerDisconnected = "peer_disconnected"
)

type Notification struct {
	Topic string
	Value map[string]any
}

type Notifier interface {
	Notify(*Notification)
}

type Notifiee interface {
	Subscribe(topics ...string) chan *Notification
	Unsubscribe(chan *Notification)
}

type Notifications struct {
	ps *pubsub.PubSub[string, *Notification]
}

func (n *Notifications) Notify(notification *Notification) {
	n.ps.Pub(notification, notification.Topic)
}

func (n *Notifications) Subscribe(topics ...string) chan *Notification {
	return n.ps.Sub(topics...)
}

func (n *Notifications) Unsubscribe(ch chan *Notification) {
	n.ps.Unsub(ch)
}
