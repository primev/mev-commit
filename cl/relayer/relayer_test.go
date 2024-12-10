package relayer

import (
	"context"
	"errors"
	"net"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"log/slog"

	"github.com/go-redis/redismock/v9"
	"github.com/primev/mev-commit/cl/pb/pb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestCreateConsumerGroup(t *testing.T) {
	logger := slog.Default()
	db, mock := redismock.NewClientMock()

	r := &Relayer{
		redisClient: db,
		logger:      logger,
		server:      grpc.NewServer(),
	}

	groupName := "member_group:testClient"
	mock.ExpectXGroupCreateMkStream(blockStreamName, groupName, "0").SetVal("OK")
	err := r.createConsumerGroup(context.Background(), groupName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mock.ClearExpect()
	mock.ExpectXGroupCreateMkStream(blockStreamName, groupName, "0").SetErr(errors.New("BUSYGROUP Consumer Group name already exists"))
	err = r.createConsumerGroup(context.Background(), groupName)
	if err != nil {
		t.Fatalf("expected no error on BUSYGROUP, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestAckMessage(t *testing.T) {
	logger := slog.Default()
	db, mock := redismock.NewClientMock()

	r := &Relayer{
		redisClient: db,
		logger:      logger,
	}

	groupName := "member_group:testClient"
	messageID := "123-1"

	mock.ExpectXAck(blockStreamName, groupName, messageID).SetVal(1)

	err := r.ackMessage(context.Background(), messageID, groupName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestReadMessages(t *testing.T) {
	logger := slog.Default()
	db, mock := redismock.NewClientMock()

	r := &Relayer{
		redisClient: db,
		logger:      logger,
	}

	groupName := "member_group:testClient"
	consumerName := "member_consumer:testClient"

	mock.ExpectXReadGroup(&redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{blockStreamName, string(RedisMsgTypePending)},
		Count:    1,
		Block:    time.Second,
	}).SetErr(redis.Nil) // simulating no pending messages

	mock.ExpectXReadGroup(&redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{blockStreamName, string(RedisMsgTypeNew)},
		Count:    1,
		Block:    time.Second,
	}).SetVal([]redis.XStream{
		{
			Stream: blockStreamName,
			Messages: []redis.XMessage{
				{
					ID: "123-1",
					Values: map[string]interface{}{
						"payload_id":         "payload_123",
						"execution_payload":  "some_encoded_payload",
						"sender_instance_id": "instance_abc",
					},
				},
			},
		},
	})

	messages, err := r.readMessages(context.Background(), groupName, consumerName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(messages) != 1 || len(messages[0].Messages) != 1 {
		t.Fatalf("expected 1 message, got %v", messages)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	logger := slog.Default()
	db, mock := redismock.NewClientMock()

	r := &Relayer{
		redisClient: db,
		logger:      logger,
		server:      grpc.NewServer(),
	}

	pb.RegisterRelayerServer(r.server, r)

	lis := bufconn.Listen(1024 * 1024)

	serverDone := make(chan struct{})
	var serverErr atomic.Value

	mock.ExpectXGroupCreateMkStream(blockStreamName, "member_group:testClient", "0").SetVal("OK")

	mock.ExpectXReadGroup(&redis.XReadGroupArgs{
		Group:    "member_group:testClient",
		Consumer: "member_consumer:testClient",
		Streams:  []string{blockStreamName, "0"},
		Count:    1,
		Block:    time.Second,
	}).SetErr(redis.Nil)

	mock.ExpectXReadGroup(&redis.XReadGroupArgs{
		Group:    "member_group:testClient",
		Consumer: "member_consumer:testClient",
		Streams:  []string{blockStreamName, ">"},
		Count:    1,
		Block:    time.Second,
	}).SetVal([]redis.XStream{
		{
			Stream: blockStreamName,
			Messages: []redis.XMessage{
				{
					ID: "123-1",
					Values: map[string]interface{}{
						"payload_id":         "payload_123",
						"execution_payload":  "some_encoded_payload",
						"sender_instance_id": "instance_abc",
					},
				},
			},
		},
	})

	ackCalled := make(chan struct{})

	customMatch := func(expected, actual []interface{}) error {
		if len(actual) >= 1 {
			cmdName, ok := actual[0].(string)
			if ok && strings.ToUpper(cmdName) == "XACK" {
				select {
				case <-ackCalled:
				default:
					close(ackCalled)
				}
			}
		}
		return nil
	}

	mock.CustomMatch(customMatch).ExpectXAck(blockStreamName, "member_group:testClient", "123-1").SetVal(int64(1))

	go func() {
		err := r.server.Serve(lis)
		if err != nil && err != grpc.ErrServerStopped {
			serverErr.Store(err)
		}
		close(serverDone)
	}()

	defer func() {
		r.server.GracefulStop()
		<-serverDone
		if err, ok := serverErr.Load().(error); ok {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Create a gRPC client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		"passthrough:///",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("failed to dial bufconn: %v", err)
	}
	defer conn.Close()

	client := pb.NewRelayerClient(conn)

	// Call Subscribe
	stream, err := client.Subscribe(ctx)
	if err != nil {
		t.Fatalf("failed to call Subscribe: %v", err)
	}

	err = stream.Send(&pb.ClientMessage{
		Message: &pb.ClientMessage_SubscribeRequest{
			SubscribeRequest: &pb.SubscribeRequest{
				ClientId: "testClient",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to send subscribe request: %v", err)
	}

	recvMsg, err := stream.Recv()
	if err != nil {
		t.Fatalf("failed to receive message from server: %v", err)
	}
	if recvMsg.GetPayloadId() != "payload_123" {
		t.Errorf("expected payload_123, got %s", recvMsg.GetPayloadId())
	}

	err = stream.Send(&pb.ClientMessage{
		Message: &pb.ClientMessage_AckPayload{
			AckPayload: &pb.AckPayloadRequest{
				ClientId:  "testClient",
				PayloadId: "payload_123",
				MessageId: "123-1",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to send ack: %v", err)
	}

	select {
	case <-ackCalled:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for XAck to be called")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet redis expectations: %v", err)
	}
}
