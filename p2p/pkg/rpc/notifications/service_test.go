package notificationsapi_test

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	notificationsapi "github.com/primev/mev-commit/p2p/pkg/rpc/notifications"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type testNotifiee struct {
	notifications chan *notifications.Notification
}

func (n *testNotifiee) Subscribe(_ ...notifications.Topic) chan *notifications.Notification {
	return n.notifications
}

func (n *testNotifiee) Unsubscribe(ch chan *notifications.Notification) <-chan struct{} {
	close(ch)
	done := make(chan struct{})
	close(done)
	return done
}

func startServer(t *testing.T, n *testNotifiee) notificationsapiv1.NotificationsClient {
	bufferSize := 1024 * 1024
	lis := bufconn.Listen(bufferSize)

	logger := util.NewTestLogger(os.Stdout)
	srvImpl := notificationsapi.NewService(n, logger)

	baseServer := grpc.NewServer()
	notificationsapiv1.RegisterNotificationsServer(baseServer, srvImpl)
	srvStopped := make(chan struct{})
	go func() {
		defer close(srvStopped)

		if err := baseServer.Serve(lis); err != nil {
			// Ignore "use of closed network connection" error
			if opErr, ok := err.(*net.OpError); !ok || !errors.Is(opErr.Err, net.ErrClosed) {
				t.Logf("server stopped err: %v", err)
			}
		}
	}()

	// nolint:staticcheck
	conn, err := grpc.DialContext(context.TODO(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("error connecting to server: %v", err)
	}

	t.Cleanup(func() {
		err := lis.Close()
		if err != nil {
			t.Errorf("error closing listener: %v", err)
		}
		baseServer.Stop()

		<-srvStopped
	})

	client := notificationsapiv1.NewNotificationsClient(conn)

	return client
}

func TestSubscribe(t *testing.T) {
	t.Parallel()

	n := &testNotifiee{
		notifications: make(chan *notifications.Notification),
	}
	client := startServer(t, n)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.Subscribe(ctx, &notificationsapiv1.SubscribeRequest{
		Topics: []string{
			string(notifications.TopicPeerConnected),
			string(notifications.TopicPeerDisconnected),
		},
	})
	if err != nil {
		t.Fatalf("error subscribing: %v", err)
	}

	n1 := notifications.NewNotification(
		notifications.TopicPeerConnected,
		map[string]interface{}{
			"key": "value",
		},
	)

	n.notifications <- n1

	resp, err := stream.Recv()
	if err != nil {
		t.Fatalf("error receiving notification: %v", err)
	}

	if resp.Topic != string(n1.Topic()) {
		t.Errorf("expected topic %q, got %q", n1.Topic(), resp.Topic)
	}

	if resp.Value.Fields["key"].GetStringValue() != n1.Value()["key"] {
		t.Errorf(
			"expected value %q, got %q",
			n1.Value()["key"], resp.Value.Fields["key"].GetStringValue(),
		)
	}

	cancel()

	_, err = stream.Recv()
	if err == nil {
		t.Error("expected error receiving notification")
	}

	if err.Error() != "rpc error: code = Canceled desc = context canceled" {
		t.Errorf("expected context canceled error, got %v", err)
	}
}
