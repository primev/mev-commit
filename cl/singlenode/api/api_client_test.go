package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewPayloadClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	baseURL := "http://localhost:8080"

	client := NewPayloadClient(baseURL, logger)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30s, got %v", client.httpClient.Timeout)
	}

	if client.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestPayloadClient_GetLatestPayload_Success(t *testing.T) {
	// Create test payload response
	expectedPayload := PayloadResponse{
		PayloadID:        "payload_123",
		ExecutionPayload: "0x1234567890abcdef",
		BlockHeight:      100,
		Timestamp:        time.Now().Unix(),
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/payload/latest" {
			t.Errorf("Expected path /api/v1/payload/latest, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedPayload)
		if err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Create client
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	// Test the method
	ctx := context.Background()
	result, err := client.GetLatestPayload(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PayloadID != expectedPayload.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", expectedPayload.PayloadID, result.PayloadID)
	}

	if result.ExecutionPayload != expectedPayload.ExecutionPayload {
		t.Errorf("Expected ExecutionPayload %s, got %s", expectedPayload.ExecutionPayload, result.ExecutionPayload)
	}

	if result.BlockHeight != expectedPayload.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", expectedPayload.BlockHeight, result.BlockHeight)
	}
}

func TestPayloadClient_GetLatestPayload_ErrorResponse(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := ErrorResponse{
			Error:   "internal_error",
			Code:    500,
			Message: "Internal server error occurred",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			t.Fatalf("Failed to encode error response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetLatestPayload(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedError := "API error: Internal server error occurred"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestPayloadClient_GetLatestPayload_InvalidJSON(t *testing.T) {
	// Create test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("invalid json"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetLatestPayload(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestPayloadClient_GetLatestPayload_ContextCanceled(t *testing.T) {
	// Create test server with delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	// Create context that will be canceled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := client.GetLatestPayload(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestPayloadClient_GetPayloadsSince_Success(t *testing.T) {
	// Create test payload list response
	expectedResponse := PayloadListResponse{
		Payloads: []PayloadResponse{
			{
				PayloadID:        "payload_100",
				ExecutionPayload: "0x100",
				BlockHeight:      100,
				Timestamp:        time.Now().Unix(),
			},
			{
				PayloadID:        "payload_101",
				ExecutionPayload: "0x101",
				BlockHeight:      101,
				Timestamp:        time.Now().Unix(),
			},
		},
		HasMore:    true,
		NextHeight: 102,
		TotalCount: 50,
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/payload/since/100"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		expectedQuery := "limit=10"
		if r.URL.RawQuery != expectedQuery {
			t.Errorf("Expected query %s, got %s", expectedQuery, r.URL.RawQuery)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedResponse)
		if err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetPayloadsSince(ctx, 100, 10)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Payloads) != 2 {
		t.Errorf("Expected 2 payloads, got %d", len(result.Payloads))
	}

	if result.HasMore != true {
		t.Error("Expected HasMore to be true")
	}

	if result.NextHeight != 102 {
		t.Errorf("Expected NextHeight 102, got %d", result.NextHeight)
	}

	if result.TotalCount != 50 {
		t.Errorf("Expected TotalCount 50, got %d", result.TotalCount)
	}
}

func TestPayloadClient_GetPayloadsSince_ErrorResponse(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := ErrorResponse{
			Error:   "not_found",
			Code:    404,
			Message: "No payloads found for the given height",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			t.Fatalf("Failed to encode error response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetPayloadsSince(ctx, 999, 10)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedError := "API error: No payloads found for the given height"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestPayloadClient_GetPayloadByHeight_Success(t *testing.T) {
	// Create test payload response
	expectedPayload := PayloadResponse{
		PayloadID:        "payload_150",
		ExecutionPayload: "0x150",
		BlockHeight:      150,
		Timestamp:        time.Now().Unix(),
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/payload/height/150"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedPayload)
		if err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetPayloadByHeight(ctx, 150)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PayloadID != expectedPayload.PayloadID {
		t.Errorf("Expected PayloadID %s, got %s", expectedPayload.PayloadID, result.PayloadID)
	}

	if result.BlockHeight != expectedPayload.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", expectedPayload.BlockHeight, result.BlockHeight)
	}
}

func TestPayloadClient_GetPayloadByHeight_NotFound(t *testing.T) {
	// Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorResp := ErrorResponse{
			Error:   "not_found",
			Code:    404,
			Message: "Payload not found for height 999",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(errorResp)
		if err != nil {
			t.Fatalf("Failed to encode error response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetPayloadByHeight(ctx, 999)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedError := "API error: Payload not found for height 999"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestPayloadClient_CheckHealth_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/health" {
			t.Errorf("Expected path /api/v1/health, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	err := client.CheckHealth(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestPayloadClient_CheckHealth_Unhealthy(t *testing.T) {
	// Create test server that returns unhealthy status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	err := client.CheckHealth(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "leader node unhealthy (status 500)"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestPayloadClient_CheckHealth_NetworkError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient("http://nonexistent.example.com", logger)

	ctx := context.Background()
	err := client.CheckHealth(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should contain "failed to execute health check"
	if !contains(err.Error(), "failed to execute health check") {
		t.Errorf("Expected error to contain 'failed to execute health check', got %s", err.Error())
	}
}

func TestPayloadClient_ErrorResponse_InvalidJSON(t *testing.T) {
	// Create test server that returns non-JSON error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Bad Request - Not JSON"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)

	ctx := context.Background()
	result, err := client.GetLatestPayload(ctx)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	expectedError := "API error (status 400): Bad Request - Not JSON"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

// Benchmark tests
func BenchmarkPayloadClient_GetLatestPayload(b *testing.B) {
	payload := PayloadResponse{
		PayloadID:        "payload_bench",
		ExecutionPayload: "0xbenchmark",
		BlockHeight:      1000,
		Timestamp:        time.Now().Unix(),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			b.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewPayloadClient(server.URL, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetLatestPayload(ctx)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Table-driven test for multiple scenarios
func TestPayloadClient_GetLatestPayload_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse func(w http.ResponseWriter, r *http.Request)
		expectedError  bool
		errorContains  string
	}{
		{
			name: "successful response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				payload := PayloadResponse{
					PayloadID:   "test_payload",
					BlockHeight: 42,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(payload)
				if err != nil {
					t.Fatalf("Failed to encode response: %v", err)
				}
			},
			expectedError: false,
		},
		{
			name: "server error",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte("Internal Server Error"))
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			},
			expectedError: true,
			errorContains: "API error (status 500)",
		},
		{
			name: "delayed response",
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(50 * time.Millisecond) // Short delay to test but not cause timeout
				payload := PayloadResponse{
					PayloadID:   "delayed_payload",
					BlockHeight: 123,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(payload)
				if err != nil {
					t.Fatalf("Failed to encode response: %v", err)
				}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer server.Close()

			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
			client := NewPayloadClient(server.URL, logger)

			ctx := context.Background()
			result, err := client.GetLatestPayload(ctx)

			if tt.expectedError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %s, got %s", tt.errorContains, err.Error())
				}
				if result != nil {
					t.Error("Expected nil result on error")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if result == nil {
					t.Error("Expected non-nil result on success")
				}
			}
		})
	}
}
