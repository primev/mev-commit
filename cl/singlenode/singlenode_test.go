package singlenode

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/primev/mev-commit/cl/blockbuilder"
	localstate "github.com/primev/mev-commit/cl/singlenode/state"
	"github.com/primev/mev-commit/cl/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockBlockBuilder implements the BlockBuilder interface for testing
type MockBlockBuilder struct {
	mock.Mock
}

func (m *MockBlockBuilder) GetPayload(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBlockBuilder) FinalizeBlock(ctx context.Context, payloadID string, executionPayload string, extraData string) error {
	args := m.Called(ctx, payloadID, executionPayload, extraData)
	return args.Error(0)
}

// MockConnectionRefused provides a safe implementation for testing
type MockConnectionRefused struct{}

func (m *MockConnectionRefused) IsConnectionRefused(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "connection refused")
}

// setupTestLogger creates a logger for testing
func setupTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// TestNewSingleNodeApp tests the creation of a new SingleNodeApp
func TestNewSingleNodeApp(t *testing.T) {
	ctx := context.Background()
	logger := setupTestLogger()

	validCfg := Config{
		InstanceID:               "test-instance",
		EthClientURL:             "http://localhost:8545",
		JWTSecret:                "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		EVMBuildDelay:            time.Second,
		EVMBuildDelayEmptyBlocks: time.Second * 2,
		PriorityFeeReceipt:       "0x1234567890abcdef1234567890abcdef12345678",
		HealthAddr:               ":8080",
	}

	app, err := NewSingleNodeApp(ctx, validCfg, logger)

	if err == nil && app != nil {
		app.Stop()
	}

	invalidJWTCfg := validCfg
	invalidJWTCfg.JWTSecret = "invalid-jwt"

	_, err = NewSingleNodeApp(ctx, invalidJWTCfg, logger)
	require.Error(t, err, "Expected error with invalid JWT secret")
}

// TestHealthHandler tests the health endpoint
func TestHealthHandler(t *testing.T) {
	logger := setupTestLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := &SingleNodeApp{
		logger:            logger,
		appCtx:            ctx,
		cancel:            cancel,
		connectionRefused: false,
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	app.healthHandler(w, req)

	resp := w.Result()
	//nolint:errcheck
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK for healthy app")

	app.connectionRefused = true
	w = httptest.NewRecorder()
	app.healthHandler(w, req)

	resp = w.Result()
	//nolint:errcheck
	defer resp.Body.Close()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Expected 503 when connection refused")

	app.connectionRefused = false
	app.cancel()
	w = httptest.NewRecorder()
	app.healthHandler(w, req)

	resp = w.Result()
	//nolint:errcheck
	defer resp.Body.Close()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "Expected 503 when context canceled")
}

// TestSetConnectionStatus tests the connection status management
func TestSetConnectionStatus(t *testing.T) {
	logger := setupTestLogger()
	ctx := context.Background()

	app := &SingleNodeApp{
		logger:            logger,
		appCtx:            ctx,
		connectionRefused: false,
	}

	app.setConnectionStatus(nil)
	assert.False(t, app.connectionRefused, "Connection refused should be false after nil error")

	err := fmt.Errorf("connection refused")
	app.setConnectionStatus(err)
	assert.True(t, app.connectionRefused, "Connection refused should be true after connection refused error")

	app.connectionRefused = false
	app.setConnectionStatus(fmt.Errorf("some other error"))
	assert.False(t, app.connectionRefused, "Connection refused should remain false after other error")
}

// TestProduceBlock tests the block production cycle
func TestProduceBlock(t *testing.T) {
	logger := setupTestLogger()
	ctx := context.Background()

	mockBuilder := new(MockBlockBuilder)
	stateMgr := localstate.NewLocalStateManager(logger)

	app := &SingleNodeApp{
		logger:       logger,
		appCtx:       ctx,
		blockBuilder: mockBuilder,
		stateManager: stateMgr,
	}

	err := stateMgr.SaveBlockStateAndPublishToStream(ctx, &types.BlockBuildState{
		CurrentStep:      types.StepBuildBlock,
		PayloadID:        "test-payload-id",
		ExecutionPayload: "test-execution-payload",
	})
	require.NoError(t, err)

	mockBuilder.On("GetPayload", mock.Anything).Return(nil)
	mockBuilder.On("FinalizeBlock", mock.Anything, "test-payload-id", "test-execution-payload", "").Return(nil)

	err = app.produceBlock()
	require.NoError(t, err, "Expected no error from produceBlock")

	mockBuilder.AssertExpectations(t)

	mockBuilder = new(MockBlockBuilder)
	app.blockBuilder = mockBuilder

	mockBuilder.On("GetPayload", mock.Anything).Return(assert.AnError)

	err = app.produceBlock()
	require.Error(t, err, "Expected error from produceBlock when GetPayload fails")
	assert.Contains(t, err.Error(), "failed to get payload", "Expected specific error message")

	mockBuilder.AssertExpectations(t)

	mockBuilder = new(MockBlockBuilder)
	app.blockBuilder = mockBuilder

	err = stateMgr.SaveBlockStateAndPublishToStream(ctx, &types.BlockBuildState{
		CurrentStep: types.StepBuildBlock,
		PayloadID:   "", // Empty payload ID
	})
	require.NoError(t, err)

	mockBuilder.On("GetPayload", mock.Anything).Return(nil)

	err = app.produceBlock()
	require.Error(t, err, "Expected error with empty payload ID")
	assert.Contains(t, err.Error(), "payload ID is empty", "Expected specific error message")

	mockBuilder = new(MockBlockBuilder)
	app.blockBuilder = mockBuilder

	err = stateMgr.SaveBlockStateAndPublishToStream(ctx, &types.BlockBuildState{
		CurrentStep:      types.StepBuildBlock,
		PayloadID:        "test-payload-id",
		ExecutionPayload: "test-execution-payload",
	})
	require.NoError(t, err)

	mockBuilder.On("GetPayload", mock.Anything).Return(nil)
	mockBuilder.On("FinalizeBlock", mock.Anything, "test-payload-id", "test-execution-payload", "").Return(assert.AnError)

	err = app.produceBlock()
	require.Error(t, err, "Expected error from produceBlock when FinalizeBlock fails")
	assert.Contains(t, err.Error(), "failed to finalize block", "Expected specific error message")

	mockBuilder.AssertExpectations(t)

	mockBuilder = new(MockBlockBuilder)
	app.blockBuilder = mockBuilder

	mockBuilder.On("GetPayload", mock.Anything).Return(blockbuilder.ErrEmptyBlock)

	err = app.produceBlock()
	assert.Contains(t, err.Error(), blockbuilder.ErrEmptyBlock.Error(),
		"Expected error to contain ErrEmptyBlock message")

	mockBuilder.AssertExpectations(t)
}

// TestRunLoop tests parts of the run loop that can be isolated
func TestRunLoop(t *testing.T) {
	logger := setupTestLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockBuilder := new(MockBlockBuilder)
	stateMgr := localstate.NewLocalStateManager(logger)

	app := &SingleNodeApp{
		logger:       logger,
		appCtx:       ctx,
		cancel:       cancel,
		blockBuilder: mockBuilder,
		stateManager: stateMgr,
	}

	result := app.resetBlockProduction("Test reset")
	assert.False(t, result, "Expected resetBlockProduction to return false on success")

	state := stateMgr.GetBlockBuildState(ctx)
	assert.Equal(t, types.StepBuildBlock, state.CurrentStep, "Expected state to be reset")

	testCancel := context.CancelFunc(func() {
		testCtx, testCancelFunc := context.WithCancel(context.Background())
		testCancelFunc()
		app.appCtx = testCtx
	})

	originalCancel := app.cancel

	app.cancel = testCancel

	app.shutdownWithError(assert.AnError, "Test shutdown error")

	select {
	case <-app.appCtx.Done():
		// Context was canceled as expected
	default:
		t.Error("Context was not canceled by shutdownWithError")
	}

	app.cancel = originalCancel
}

// TestStartStop tests the Start and Stop methods
func TestStartStop(t *testing.T) {
	logger := setupTestLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stateMgr := localstate.NewLocalStateManager(logger)

	mockBuilder := new(MockBlockBuilder)

	// Setup mock to immediately cancel the context when GetPayload is called
	// This ensures the run loop exits right away
	mockBuilder.On("GetPayload", mock.Anything).Run(func(args mock.Arguments) {
		cancel() // Cancel context to exit the run loop immediately
	}).Return(context.Canceled).Once()

	// Create app with minimal configuration for testing
	app := &SingleNodeApp{
		logger:       logger,
		cfg:          Config{HealthAddr: ":0"},
		appCtx:       ctx,
		cancel:       cancel,
		blockBuilder: mockBuilder,
		stateManager: stateMgr,
	}

	app.Start()

	time.Sleep(100 * time.Millisecond)

	app.Stop()

	mockBuilder.AssertExpectations(t)
}

// TestRunLoopEmptyBlockHandling tests how runLoop handles ErrEmptyBlock
func TestRunLoopEmptyBlockHandling(t *testing.T) {
	logger := setupTestLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockBuilder := new(MockBlockBuilder)
	stateMgr := localstate.NewLocalStateManager(logger)

	app := &SingleNodeApp{
		logger:       logger,
		appCtx:       ctx,
		cancel:       cancel,
		blockBuilder: mockBuilder,
		stateManager: stateMgr,
	}

	mockBuilder.On("GetPayload", mock.Anything).Return(blockbuilder.ErrEmptyBlock).Once()

	mockBuilder.On("GetPayload", mock.Anything).Run(func(args mock.Arguments) {
		cancel()
	}).Return(context.Canceled).Once()

	done := make(chan struct{})
	go func() {
		app.runLoop()
		close(done)
	}()

	// Wait for the run loop to exit
	select {
	case <-done:
		// Run loop exited as expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Timeout waiting for run loop to exit")
	}

	mockBuilder.AssertExpectations(t)
}

// TestIsConnectionRefused tests the connection refused detection logic
func TestIsConnectionRefused(t *testing.T) {
	mock := &MockConnectionRefused{}

	err := fmt.Errorf("connection refused")
	assert.True(t, mock.IsConnectionRefused(err),
		"Should detect 'connection refused' error")

	err = fmt.Errorf("Something with connection refused inside")
	assert.True(t, mock.IsConnectionRefused(err),
		"Should detect error with 'connection refused' substring")

	err = fmt.Errorf("some other error")
	assert.False(t, mock.IsConnectionRefused(err),
		"Should not detect generic error as connection refused")

	assert.False(t, mock.IsConnectionRefused(nil),
		"Should not detect nil as connection refused")
}
