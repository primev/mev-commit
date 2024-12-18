package member

import (
	"context"
	"encoding/hex"
	"io"
	"log/slog"
	"math/big"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/cl/pb/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type mockEngineClient struct{}

func (m *mockEngineClient) NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash, beaconRoot *common.Hash) (engine.PayloadStatusV1, error) {
	return engine.PayloadStatusV1{}, nil
}

func (m *mockEngineClient) ForkchoiceUpdatedV3(context.Context, engine.ForkchoiceStateV1, *engine.PayloadAttributes) (engine.ForkChoiceResponse, error) {
	return engine.ForkChoiceResponse{}, nil
}

func (m *mockEngineClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return &types.Header{}, nil
}

type mockBlockBuilder struct {
	mu            sync.Mutex
	finalizeCalls []finalizeCall
}

type finalizeCall struct {
	payloadID        string
	executionPayload string
	messageID        string
}

func (m *mockBlockBuilder) FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.finalizeCalls = append(m.finalizeCalls, finalizeCall{
		payloadID:        payloadIDStr,
		executionPayload: executionPayloadStr,
		messageID:        msgID,
	})
	return nil
}

func (f *mockBlockBuilder) Calls() []finalizeCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]finalizeCall(nil), f.finalizeCalls...)
}

// fakePayloadStreamerServer simulates the PayloadStreamer gRPC service for testing.
type fakePayloadStreamerServer struct {
	pb.UnimplementedPayloadStreamerServer

	mu            sync.Mutex
	subscribed    bool
	sentPayload   bool
	clientID      string
	serverStopped bool
}

func (s *fakePayloadStreamerServer) Subscribe(stream pb.PayloadStreamer_SubscribeServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF || s.serverStopped {
			return nil
		}
		if err != nil {
			return err
		}

		if req := msg.GetSubscribeRequest(); req != nil {
			// Acknowledge subscription
			s.mu.Lock()
			s.subscribed = true
			s.clientID = req.GetClientId()
			s.mu.Unlock()

			// After subscribing, send a single payload message, then close the stream.
			resp := &pb.PayloadMessage{
				PayloadId:        "test-payload-id",
				ExecutionPayload: "test-exec-payload",
				SenderInstanceId: "sender-123",
				MessageId:        "test-msg-id",
			}
			if err := stream.SendMsg(resp); err != nil {
				return err
			}
			s.mu.Lock()
			s.sentPayload = true
			s.mu.Unlock()

			// Wait a moment and then return EOF to stop the stream
			time.Sleep(200 * time.Millisecond)
			return nil
		} else if ack := msg.GetAckPayload(); ack != nil {
			continue
		}
	}
}

// TestMemberClientRun tests the Run method end-to-end with a fake server and fake dependencies.
func TestMemberClientRun(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0") // ephemeral port
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	defer s.Stop()

	streamerServer := &fakePayloadStreamerServer{}
	pb.RegisterPayloadStreamerServer(s, streamerServer)

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.Serve(lis)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	case <-time.After(time.Millisecond * 100):
		// Server started successfully
	}

	clientID := "test-client-id"
	streamerAddr := lis.Addr().String()
	logger := slog.Default()

	engineClient := &mockEngineClient{}
	blockBuilder := &mockBlockBuilder{}

	conn, err := grpc.NewClient(streamerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial test server: %v", err)
	}
	streamerClient := pb.NewPayloadStreamerClient(conn)

	mc := &MemberClient{
		clientID:     clientID,
		streamerAddr: streamerAddr,
		conn:         conn,
		client:       streamerClient,
		logger:       logger,
		engineCl:     engineClient,
		bb:           blockBuilder,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = mc.Run(ctx)
	if err != nil {
		t.Errorf("MemberClient.Run returned an error: %v", err)
	}

	streamerServer.mu.Lock()
	subscribed := streamerServer.subscribed
	sentPayload := streamerServer.sentPayload
	streamerServer.mu.Unlock()

	if !subscribed {
		t.Errorf("Server did not receive subscription from client")
	}
	if !sentPayload {
		t.Errorf("Server did not send a payload message")
	}

	calls := blockBuilder.Calls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 FinalizeBlock call, got %d", len(calls))
	}
	call := calls[0]
	if call.payloadID != "test-payload-id" {
		t.Errorf("Expected payloadID 'test-payload-id', got '%s'", call.payloadID)
	}
	if call.executionPayload != "test-exec-payload" {
		t.Errorf("Expected executionPayload 'test-exec-payload', got '%s'", call.executionPayload)
	}
	if call.messageID != "test-msg-id" {
		t.Errorf("Expected messageID 'test-msg-id', got '%s'", call.messageID)
	}
}

func TestJWTSecretDecodingNoMocks(t *testing.T) {
	validSecret := "deadbeef"
	invalidSecret := "zzzz"

	_, err := hex.DecodeString(validSecret)
	if err != nil {
		t.Errorf("Failed to decode valid secret: %v", err)
	}

	_, err = hex.DecodeString(invalidSecret)
	if err == nil {
		t.Error("Expected error decoding invalid secret, got none")
	}
}