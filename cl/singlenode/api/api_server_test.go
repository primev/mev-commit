package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"log/slog"

	"github.com/primev/mev-commit/cl/types"
)

// Mock implementations for testing

type mockStateManager struct {
	state types.BlockBuildState
}

func (m *mockStateManager) GetBlockBuildState(ctx context.Context) types.BlockBuildState {
	return m.state
}

type mockPayloadRepository struct {
	latestPayload    *types.PayloadInfo
	latestError      error
	payloadsByHeight map[uint64]*types.PayloadInfo
	payloadsSince    []types.PayloadInfo
	payloadsSinceErr error
	getByHeightErr   error
}

func (m *mockPayloadRepository) Close() error {
	return nil
}

func (m *mockPayloadRepository) GetLatestPayload(ctx context.Context) (*types.PayloadInfo, error) {
	return m.latestPayload, m.latestError
}

func (m *mockPayloadRepository) GetPayloadsSince(ctx context.Context, height uint64, limit int) ([]types.PayloadInfo, error) {
	if m.payloadsSinceErr != nil {
		return nil, m.payloadsSinceErr
	}
	return m.payloadsSince, nil
}

func (m *mockPayloadRepository) GetPayloadByHeight(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
	if m.getByHeightErr != nil {
		return nil, m.getByHeightErr
	}
	if payload, exists := m.payloadsByHeight[height]; exists {
		return payload, nil
	}
	return nil, sql.ErrNoRows
}
func (m *mockPayloadRepository) SavePayload(ctx context.Context, info *types.PayloadInfo) error {
	if m.latestPayload != nil && m.latestPayload.PayloadID == info.PayloadID {
		return fmt.Errorf("payload already exists")
	}

	// Simulate saving the payload
	m.latestPayload = info

	if info.BlockHeight == 0 {
		return fmt.Errorf("invalid block height")
	}

	m.payloadsByHeight[info.BlockHeight] = info
	return nil
}

// Helper function to create a test logger
func createTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError, // Only log errors during tests
	}))
}

// Helper function to create test payload info
func createTestPayloadInfo(id string, height uint64) *types.PayloadInfo {
	return &types.PayloadInfo{
		PayloadID:        id,
		ExecutionPayload: fmt.Sprintf("payload-data-%s", id),
		BlockHeight:      height,
		InsertedAt:       time.Now(),
	}
}

func TestNewPayloadServer(t *testing.T) {
	logger := createTestLogger()
	stateManager := &mockStateManager{}
	payloadRepo := &mockPayloadRepository{}

	server := NewPayloadServer("localhost:8080", stateManager, payloadRepo, logger)

	if server == nil {
		t.Fatal("Expected non-nil PayloadServer")
	}

	if server.server.Addr != "localhost:8080" {
		t.Errorf("Expected addr localhost:8080, got %s", server.server.Addr)
	}

	if server.server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout 30s, got %v", server.server.ReadTimeout)
	}

	if server.server.WriteTimeout != 30*time.Second {
		t.Errorf("Expected WriteTimeout 30s, got %v", server.server.WriteTimeout)
	}

	if server.server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", server.server.IdleTimeout)
	}
}

func TestHandleGetLatestPayload_FromRepository(t *testing.T) {
	testPayload := createTestPayloadInfo("test-payload-1", 100)

	stateManager := &mockStateManager{}
	payloadRepo := &mockPayloadRepository{
		latestPayload: testPayload,
	}

	server := NewPayloadServer("localhost:8080", stateManager, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/latest", nil)
	w := httptest.NewRecorder()

	server.handleGetLatestPayload(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PayloadResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.PayloadID != testPayload.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", testPayload.PayloadID, response.PayloadID)
	}

	if response.ExecutionPayload != testPayload.ExecutionPayload {
		t.Errorf("Expected ExecutionPayload %s, got %s", testPayload.ExecutionPayload, response.ExecutionPayload)
	}

	if response.BlockHeight != testPayload.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", testPayload.BlockHeight, response.BlockHeight)
	}
}

func TestHandleGetLatestPayload_FromStateManager(t *testing.T) {
	state := types.BlockBuildState{
		PayloadID:        "state-payload-1",
		ExecutionPayload: "state-execution-data",
	}

	stateManager := &mockStateManager{state: state}
	payloadRepo := &mockPayloadRepository{
		latestPayload: nil,
		latestError:   fmt.Errorf("repository error"),
	}

	server := NewPayloadServer("localhost:8080", stateManager, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/latest", nil)
	w := httptest.NewRecorder()

	server.handleGetLatestPayload(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PayloadResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.PayloadID != state.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", state.PayloadID, response.PayloadID)
	}

	if response.ExecutionPayload != state.ExecutionPayload {
		t.Errorf("Expected ExecutionPayload %s, got %s", state.ExecutionPayload, response.ExecutionPayload)
	}

	if response.BlockHeight != 0 {
		t.Errorf("Expected BlockHeight 0, got %d", response.BlockHeight)
	}
}

func TestHandleGetLatestPayload_NoPayloadAvailable(t *testing.T) {
	stateManager := &mockStateManager{
		state: types.BlockBuildState{},
	}
	payloadRepo := &mockPayloadRepository{
		latestPayload: nil,
		latestError:   fmt.Errorf("no payload"),
	}

	server := NewPayloadServer("localhost:8080", stateManager, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/latest", nil)
	w := httptest.NewRecorder()

	server.handleGetLatestPayload(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Error != "No payload available" {
		t.Errorf("Expected error 'No payload available', got %s", response.Error)
	}
}

func TestHandleGetLatestPayload_MethodNotAllowed(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, &mockPayloadRepository{}, createTestLogger())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payload/latest", nil)
	w := httptest.NewRecorder()

	server.handleGetLatestPayload(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestHandleGetPayloadsSince_Success(t *testing.T) {
	payloads := []types.PayloadInfo{
		*createTestPayloadInfo("payload-1", 101),
		*createTestPayloadInfo("payload-2", 102),
		*createTestPayloadInfo("payload-3", 103),
	}

	payloadRepo := &mockPayloadRepository{
		payloadsSince: payloads,
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/since/100", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadsSince(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PayloadListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Payloads) != 3 {
		t.Errorf("Expected 3 payloads, got %d", len(response.Payloads))
	}

	if response.HasMore != false {
		t.Errorf("Expected HasMore false, got %v", response.HasMore)
	}

	if response.NextHeight != 104 {
		t.Errorf("Expected NextHeight 104, got %d", response.NextHeight)
	}

	if response.TotalCount != 3 {
		t.Errorf("Expected TotalCount 3, got %d", response.TotalCount)
	}
}

func TestHandleGetPayloadsSince_WithLimit(t *testing.T) {
	// Create 6 payloads but set limit to 2
	payloads := make([]types.PayloadInfo, 3) // Return 3 to test hasMore logic
	for i := 0; i < 3; i++ {
		payloads[i] = *createTestPayloadInfo(fmt.Sprintf("payload-%d", i+1), uint64(101+i))
	}

	payloadRepo := &mockPayloadRepository{
		payloadsSince: payloads,
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/since/100?limit=2", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadsSince(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PayloadListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response.Payloads) != 2 {
		t.Errorf("Expected 2 payloads, got %d", len(response.Payloads))
	}

	if response.HasMore != true {
		t.Errorf("Expected HasMore true, got %v", response.HasMore)
	}
}

func TestHandleGetPayloadsSince_NoRepository(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, nil, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/since/100", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadsSince(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}
}

func TestHandleGetPayloadsSince_InvalidHeight(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, &mockPayloadRepository{}, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/since/invalid", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadsSince(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleGetPayloadsSince_RepositoryError(t *testing.T) {
	payloadRepo := &mockPayloadRepository{
		payloadsSinceErr: fmt.Errorf("database error"),
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/since/100", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadsSince(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandleGetPayloadByHeight_Success(t *testing.T) {
	testPayload := createTestPayloadInfo("height-payload", 150)

	payloadRepo := &mockPayloadRepository{
		payloadsByHeight: map[uint64]*types.PayloadInfo{
			150: testPayload,
		},
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/height/150", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadByHeight(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PayloadResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.PayloadID != testPayload.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", testPayload.PayloadID, response.PayloadID)
	}

	if response.BlockHeight != testPayload.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", testPayload.BlockHeight, response.BlockHeight)
	}
}

func TestHandleGetPayloadByHeight_NotFound(t *testing.T) {
	payloadRepo := &mockPayloadRepository{
		payloadsByHeight: map[uint64]*types.PayloadInfo{},
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/height/999", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadByHeight(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHandleGetPayloadByHeight_InvalidHeight(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, &mockPayloadRepository{}, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/height/invalid", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadByHeight(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleGetPayloadByHeight_NoRepository(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, nil, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/height/100", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadByHeight(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}
}

func TestHandleGetPayloadByHeight_RepositoryError(t *testing.T) {
	payloadRepo := &mockPayloadRepository{
		getByHeightErr: fmt.Errorf("database connection error"),
	}

	server := NewPayloadServer("localhost:8080", &mockStateManager{}, payloadRepo, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payload/height/100", nil)
	w := httptest.NewRecorder()

	server.handleGetPayloadByHeight(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandleHealth_Success(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, &mockPayloadRepository{}, createTestLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := strings.TrimSpace(w.Body.String())
	if body != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", body)
	}
}

func TestHandleHealth_MethodNotAllowed(t *testing.T) {
	server := NewPayloadServer("localhost:8080", &mockStateManager{}, &mockPayloadRepository{}, createTestLogger())

	req := httptest.NewRequest(http.MethodPost, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestConvertToResponse(t *testing.T) {
	timestamp := time.Now()
	payload := &types.PayloadInfo{
		PayloadID:        "test-id",
		ExecutionPayload: "test-execution-data",
		BlockHeight:      42,
		InsertedAt:       timestamp,
	}

	response := convertToResponse(payload)

	if response.PayloadID != payload.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", payload.PayloadID, response.PayloadID)
	}

	if response.ExecutionPayload != payload.ExecutionPayload {
		t.Errorf("Expected ExecutionPayload %s, got %s", payload.ExecutionPayload, response.ExecutionPayload)
	}

	if response.BlockHeight != payload.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", payload.BlockHeight, response.BlockHeight)
	}

	if response.Timestamp != timestamp.Unix() {
		t.Errorf("Expected Timestamp %d, got %d", timestamp.Unix(), response.Timestamp)
	}
}
