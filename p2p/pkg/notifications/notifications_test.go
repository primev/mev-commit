package notifications_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
)

func TestNotifications(t *testing.T) {
	t.Parallel()

	// Create a new Notifications object with a buffer capacity of 10.
	n := notifications.New(10)

	// Create a new channel to subscribe to the "peer_connected" topic.
	ch := n.Subscribe(notifications.TopicPeerConnected)

	// Create a new notification with the "peer_connected" topic and a value of
	// map[string]any{"peer_id": "1234"}.
	notification := notifications.NewNotification(
		notifications.TopicPeerConnected,
		map[string]any{"peer_id": "1234"},
	)

	// Notify the subscribers about the notification.
	n.Notify(notification)

	// Receive the notification from the channel.
	receivedNotification := <-ch
	if diff := cmp.Diff(
		notification,
		receivedNotification,
		cmp.AllowUnexported(notifications.Notification{}),
	); diff != "" {
		t.Errorf("unexpected notification (-want +got):\n%s", diff)
	}

	// Create a new notification with the "peer_disconnected" topic and a value of
	// map[string]any{"peer_id": "1234"}.
	notification = notifications.NewNotification(
		notifications.TopicPeerDisconnected,
		map[string]any{"peer_id": "1234"},
	)

	// Notify the subscribers about the notification.
	n.Notify(notification)

	// Should not receive the notification on the channel.
	select {
	case <-ch:
		t.Error("unexpected notification")
	case <-time.After(500 * time.Millisecond):
	}

	// Unsubscribe the channel.
	<-n.Unsubscribe(ch)

	// Create a new notification with the "peer_connected" topic and a value of
	// map[string]any{"peer_id": "1234"}.
	notification = notifications.NewNotification(
		notifications.TopicPeerConnected,
		map[string]any{"peer_id": "1234"},
	)

	// Notify the subscribers about the notification.
	n.Notify(notification)

	_, more := <-ch
	if more {
		t.Error("channel should be closed")
	}

	// Create a new channel to subscribe to the "peer_disconnected" topic.
	ch = n.Subscribe(notifications.TopicPeerDisconnected)

	// Create a new notification with the "peer_disconnected" topic and a value of
	// map[string]any{"peer_id": "1234"}.
	notification = notifications.NewNotification(
		notifications.TopicPeerDisconnected,
		map[string]any{"peer_id": "1234"},
	)

	// Notify the subscribers about the notification.
	n.Notify(notification)

	// Receive the notification from the channel.
	receivedNotification = <-ch
	if diff := cmp.Diff(
		notification,
		receivedNotification,
		cmp.AllowUnexported(notifications.Notification{}),
	); diff != "" {
		t.Errorf("unexpected notification (-want +got):\n%s", diff)
	}

	n.Shutdown()
	if _, ok := <-ch; ok {
		t.Error("channel should be closed")
	}
}
