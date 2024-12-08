package relayer

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"log/slog"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/primev/mev-commit/cl/pb/pb"
)

const blockStreamName = "mevcommit_block_stream"

type RedisMsgType string

const (
	RedisMsgTypePending RedisMsgType = "0"
	RedisMsgTypeNew     RedisMsgType = ">"
)

type Relayer struct {
	pb.UnimplementedRelayerServer
	redisClient *redis.Client
	logger      *slog.Logger
	server      *grpc.Server
}

func NewRelayer(redisAddr string, logger *slog.Logger) (*Relayer, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	err := redisClient.ConfigSet(context.Background(), "min-replicas-to-write", "1").Err()
	if err != nil {
		logger.Error("Error setting min-replicas-to-write", "error", err)
		return nil, err
	}

	return &Relayer{
		redisClient: redisClient,
		logger:      logger,
		server:      grpc.NewServer(),
	}, nil
}

func (r *Relayer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	pb.RegisterRelayerServer(r.server, r)
	reflection.Register(r.server)

	r.logger.Info("Relayer is listening", "address", address)
	return r.server.Serve(lis)
}

func (r *Relayer) Stop() {
	r.server.GracefulStop()
	if err := r.redisClient.Close(); err != nil {
		r.logger.Error("Error closing Redis client in Relayer", "error", err)
	}
}

func (r *Relayer) Subscribe(stream pb.Relayer_SubscribeServer) error {
	ctx := stream.Context()

	var clientID string
	firstMessage, err := stream.Recv()
	if err != nil {
		r.logger.Error("Failed to receive initial message", "error", err)
		return err
	}
	if req := firstMessage.GetSubscribeRequest(); req != nil {
		clientID = req.ClientId
	} else {
		return fmt.Errorf("expected SubscribeRequest, got %v", firstMessage)
	}

	groupName := "member_group:" + clientID
	consumerName := "member_consumer:" + clientID

	err = r.createConsumerGroup(ctx, groupName)
	if err != nil {
		r.logger.Error("Failed to create consumer group", "clientID", clientID, "error", err)
		return err
	}

	r.logger.Info("Subscriber connected", "clientID", clientID)
	return r.handleBidirectionalStream(stream, clientID, groupName, consumerName)
}

func (r *Relayer) createConsumerGroup(ctx context.Context, groupName string) error {
	err := r.redisClient.XGroupCreateMkStream(ctx, blockStreamName, groupName, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func (r *Relayer) handleBidirectionalStream(stream pb.Relayer_SubscribeServer, clientID, groupName, consumerName string) error {
	ctx := stream.Context()
	var pendingMessageID string

	for {
		if pendingMessageID == "" {
			// No pending message, read the next message from Redis
			messages, err := r.readMessages(ctx, groupName, consumerName)
			if err != nil {
				r.logger.Error("Error reading messages", "clientID", clientID, "error", err)
				return err
			}
			if len(messages) == 0 {
				continue
			}

			msg := messages[0]
			field := msg.Messages[0]
			pendingMessageID = field.ID

			payloadIDStr, ok := field.Values["payload_id"].(string)
			executionPayloadStr, okPayload := field.Values["execution_payload"].(string)
			senderInstanceID, okSenderID := field.Values["sender_instance_id"].(string)
			if !ok || !okPayload || !okSenderID {
				r.logger.Error("Invalid message format", "clientID", clientID)
				// Acknowledge malformed messages to prevent reprocessing
				r.ackMessage(ctx, field.ID, groupName)
				pendingMessageID = ""
				continue
			}

			err = stream.Send(&pb.PayloadMessage{
				PayloadId:        payloadIDStr,
				ExecutionPayload: executionPayloadStr,
				SenderInstanceId: senderInstanceID,
				MessageId:        field.ID,
			})
			if err != nil {
				r.logger.Error("Failed to send message to client", "clientID", clientID, "error", err)
				return err
			}
		}

		clientMsg, err := stream.Recv()
		if err != nil {
			r.logger.Error("Failed to receive acknowledgment", "clientID", clientID, "error", err)
			return err
		}

		if ack := clientMsg.GetAckPayload(); ack != nil {
			if ack.MessageId == pendingMessageID {
				err := r.ackMessage(ctx, pendingMessageID, groupName)
				if err != nil {
					r.logger.Error("Failed to acknowledge message", "clientID", clientID, "error", err)
					return err
				}
				r.logger.Info("Message acknowledged", "clientID", clientID, "messageID", pendingMessageID)
				pendingMessageID = ""
			} else {
				r.logger.Error("Received acknowledgment for unknown message ID", "clientID", clientID, "messageID", ack.MessageId)
			}
		} else {
			r.logger.Error("Expected AckPayloadRequest, got something else", "clientID", clientID)
		}
	}
}

func (r *Relayer) readMessages(ctx context.Context, groupName, consumerName string) ([]redis.XStream, error) {
	messages, err := r.readMessagesFromStream(ctx, RedisMsgTypePending, groupName, consumerName)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 || len(messages[0].Messages) == 0 {
		messages, err = r.readMessagesFromStream(ctx, RedisMsgTypeNew, groupName, consumerName)
		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (s *Relayer) readMessagesFromStream(ctx context.Context, msgType RedisMsgType, groupName, consumerName string) ([]redis.XStream, error) {
	args := &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: consumerName,
		Streams:  []string{blockStreamName, string(msgType)},
		Count:    1,
		Block:    time.Second,
	}

	messages, err := s.redisClient.XReadGroup(ctx, args).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("error reading messages: %w", err)
	}

	return messages, nil
}

func (r *Relayer) ackMessage(ctx context.Context, messageID, groupName string) error {
	return r.redisClient.XAck(ctx, blockStreamName, groupName, messageID).Err()
}
