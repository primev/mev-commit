package leaderfollower

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/primev/mev-commit/cl/mocks"
	"github.com/primev/mev-commit/cl/redisapp/types"
	"github.com/redis/go-redis/v9"
)

type mockLeaderProc struct {
	startCalled bool
	stopCalled  bool
	stopErr     error
	isLeading   bool
}

func (m *mockLeaderProc) Start() {
	m.startCalled = true
}

func (m *mockLeaderProc) Stop() error {
	m.stopCalled = true
	return m.stopErr
}

func (m *mockLeaderProc) IsLeading() bool {
	return m.isLeading
}

func TestNewLeaderFollowerManager(t *testing.T) {
	// Setup
	instanceID := "test-instance"
	logger := slog.Default()
	redisClient := redis.NewClient(&redis.Options{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateManager := mocks.NewMockCoordinator(ctrl)
	blockBuilder := mocks.NewMockBlockBuilder(ctrl)

	// Execute
	lfm, err := NewLeaderFollowerManager(instanceID, logger, redisClient, stateManager, blockBuilder)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if lfm.instanceID != instanceID {
		t.Errorf("Expected instanceID to be %s, got %s", instanceID, lfm.instanceID)
	}
}

func TestHaveMessagesToProcess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSM := mocks.NewMockCoordinator(ctrl)

	// Prepare mock state manager to return some messages
	messages := []redis.XStream{
		{
			Stream: "test-stream",
			Messages: []redis.XMessage{
				{
					ID: "1-0",
					Values: map[string]interface{}{
						"payload_id":         "test-payload-id",
						"execution_payload":  "test-execution-payload",
						"sender_instance_id": "test-instance-id",
					},
				},
			},
		},
	}

	// Set up expectations
	gomock.InOrder(
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypePending).Return([]redis.XStream{}, nil),
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypeNew).Return(messages, nil),
	)

	lfm := &LeaderFollowerManager{
		stateManager: mockSM,
		logger:       slog.Default(),
	}

	hasMessages := lfm.HaveMessagesToProcess(ctx)

	if !hasMessages {
		t.Errorf("Expected to have messages to process, but got false")
	}
}

func TestHaveMessagesToProcess_NoMessages(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSM := mocks.NewMockCoordinator(ctrl)

	// Set up expectations
	gomock.InOrder(
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypePending).Return([]redis.XStream{}, nil),
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypeNew).Return([]redis.XStream{}, nil),
	)

	lfm := &LeaderFollowerManager{
		stateManager: mockSM,
		logger:       slog.Default(),
	}

	hasMessages := lfm.HaveMessagesToProcess(ctx)

	if hasMessages {
		t.Errorf("Expected to have no messages to process, but got true")
	}
}

func TestLeaderWork_StepBuildBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSM := mocks.NewMockCoordinator(ctrl)
	mockBB := mocks.NewMockBlockBuilder(ctrl)

	lfm := &LeaderFollowerManager{
		stateManager: mockSM,
		blockBuilder: mockBB,
		logger:       slog.Default(),
	}

	mockSM.EXPECT().
		ReadMessagesFromStream(ctx, types.RedisMsgTypePending).
		Return([]redis.XStream{}, nil).
		AnyTimes()

	mockSM.EXPECT().
		ReadMessagesFromStream(ctx, types.RedisMsgTypeNew).
		Return([]redis.XStream{}, nil).
		AnyTimes()

	mockSM.EXPECT().GetBlockBuildState(ctx).Return(types.BlockBuildState{CurrentStep: types.StepBuildBlock}).AnyTimes()
	mockBB.EXPECT().GetPayload(ctx).Return(nil).AnyTimes()
	mockSM.EXPECT().ResetBlockState(ctx).Return(nil).AnyTimes()

	// Run leaderWork in a separate goroutine to allow context cancellation
	done := make(chan error)
	go func() {
		err := lfm.leaderWork(ctx)
		done <- err
	}()

	// Let it run for a short time, then cancel the context
	time.Sleep(100 * time.Millisecond)
	cancel()

	err := <-done
	if err != nil {
		t.Errorf("leaderWork returned error: %v", err)
	}
}

func TestFollowerWork_NoMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSM := mocks.NewMockCoordinator(ctrl)
	mockBB := mocks.NewMockBlockBuilder(ctrl)

	lfm := &LeaderFollowerManager{
		stateManager: mockSM,
		blockBuilder: mockBB,
		logger:       slog.Default(),
		instanceID:   "test-instance",
	}

	// Set up expectations
	gomock.InOrder(
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypePending).Return([]redis.XStream{}, nil),
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypeNew).Return([]redis.XStream{}, nil),
	)

	// Run followerWork in a separate goroutine to allow context cancellation
	done := make(chan error)
	go func() {
		err := lfm.followerWork(ctx)
		done <- err
	}()

	// Wait for the function to return
	err := <-done
	if err != nil {
		t.Errorf("followerWork returned error: %v", err)
	}
}

func TestFollowerWork_WithMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSM := mocks.NewMockCoordinator(ctrl)
	mockBB := mocks.NewMockBlockBuilder(ctrl)

	lfm := &LeaderFollowerManager{
		stateManager: mockSM,
		blockBuilder: mockBB,
		logger:       slog.Default(),
		instanceID:   "test-instance",
	}

	messages := []redis.XStream{
		{
			Stream: "test-stream",
			Messages: []redis.XMessage{
				{
					ID: "1-0",
					Values: map[string]interface{}{
						"payload_id":         "test-payload-id",
						"execution_payload":  "test-execution-payload",
						"sender_instance_id": "other-instance",
					},
				},
			},
		},
	}

	gomock.InOrder(
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypePending).Return([]redis.XStream{}, nil),
		mockSM.EXPECT().ReadMessagesFromStream(ctx, types.RedisMsgTypeNew).Return(messages, nil),
		mockBB.EXPECT().FinalizeBlock(ctx, "test-payload-id", "test-execution-payload", "1-0").Return(nil),
		mockSM.EXPECT().AckMessage(ctx, "1-0").Return(nil).Do(func(ctx context.Context, msgID string) {
			cancel()
		}),
	)

	mockSM.EXPECT().ReadMessagesFromStream(ctx, gomock.Any()).AnyTimes().Return([]redis.XStream{}, nil)

	done := make(chan error)
	go func() {
		err := lfm.followerWork(ctx)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("followerWork returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("followerWork timed out")
	}
}

func TestWaitForGoroutinesToStop(t *testing.T) {
	lfm := &LeaderFollowerManager{
		wg:     sync.WaitGroup{},
		logger: slog.Default(),
	}

	// Simulate running goroutines
	lfm.wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		lfm.wg.Done()
	}()

	err := lfm.WaitForGoroutinesToStop()
	if err != nil {
		t.Errorf("Expected goroutines to stop without error, got: %v", err)
	}
}

func TestHandleLeaderElection_Promoted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLeaderProc := &mockLeaderProc{}
	promotedCh := make(chan time.Time, 1)
	demotedCh := make(chan time.Time)
	erroredCh := make(chan error)

	lfm := &LeaderFollowerManager{
		leaderProc: mockLeaderProc,
		logger:     slog.Default(),
		promotedCh: promotedCh,
		demotedCh:  demotedCh,
		erroredCh:  erroredCh,
	}

	// Simulate promotion
	promotedCh <- time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	lfm.handleLeaderElection(ctx)

	if !lfm.isLeader.Load() {
		t.Errorf("Expected isLeader to be true after promotion")
	}
}

func TestHandleLeaderElection_Demoted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLeaderProc := &mockLeaderProc{}
	promotedCh := make(chan time.Time)
	demotedCh := make(chan time.Time, 1)
	erroredCh := make(chan error)

	lfm := &LeaderFollowerManager{
		leaderProc: mockLeaderProc,
		logger:     slog.Default(),
		promotedCh: promotedCh,
		demotedCh:  demotedCh,
		erroredCh:  erroredCh,
	}

	// Simulate demotion
	demotedCh <- time.Now()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	lfm.handleLeaderElection(ctx)

	if lfm.isLeader.Load() {
		t.Errorf("Expected isLeader to be false after demotion")
	}
}

func TestHandleLeaderElection_Errored(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockLeaderProc := mocks.NewMockLeader(ctrl)
	mockLeaderProc := &mockLeaderProc{}
	promotedCh := make(chan time.Time)
	demotedCh := make(chan time.Time)
	erroredCh := make(chan error, 1)

	lfm := &LeaderFollowerManager{
		leaderProc: mockLeaderProc,
		logger:     slog.Default(),
		promotedCh: promotedCh,
		demotedCh:  demotedCh,
		erroredCh:  erroredCh,
	}

	// Simulate error
	erroredCh <- errors.New("leader election error")
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	lfm.handleLeaderElection(ctx)
	// Check that isLeader remains unchanged (default is false)
	if lfm.isLeader.Load() {
		t.Errorf("Expected isLeader to remain false after error")
	}
}
