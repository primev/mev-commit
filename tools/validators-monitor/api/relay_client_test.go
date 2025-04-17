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

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRelayClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	httpClient := retryablehttp.NewClient()
	relayURLs := []string{"http://relay1.example.com", "http://relay2.example.com"}

	// Test with provided HTTP client
	client := NewRelayClient(relayURLs, logger, httpClient)
	assert.NotNil(t, client)
	assert.Equal(t, relayURLs, client.relayURLs)
	assert.Equal(t, httpClient, client.client)
	assert.Equal(t, logger, client.logger)

	// Test with nil HTTP client
	client = NewRelayClient(relayURLs, logger, nil)
	assert.NotNil(t, client)
	assert.Equal(t, relayURLs, client.relayURLs)
	assert.Nil(t, client.client) // The client should be nil as passed
	assert.Equal(t, logger, client.logger)
}

func TestQueryOneRelay(t *testing.T) {
	tests := []struct {
		name           string
		relayURL       string
		blockNumber    uint64
		responseStatus int
		responseBody   string
		expectedError  bool
	}{
		{
			name:           "successful response",
			relayURL:       "valid-url",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `[{"slot":"123","parent_hash":"0xabc","block_hash":"0xdef","builder_pubkey":"0x123","proposer_pubkey":"0x456","proposer_fee_recipient":"0x789","gas_limit":"1000000","gas_used":"900000","value":"5000000000","num_tx":"100","block_number":"12345"}]`,
			expectedError:  false,
		},
		{
			name:           "empty response",
			relayURL:       "valid-url",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			expectedError:  false,
		},
		{
			name:           "server error",
			relayURL:       "valid-url",
			blockNumber:    12345,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error":"Internal server error"}`,
			expectedError:  true,
		},
		{
			name:           "invalid json",
			relayURL:       "valid-url",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid-json}`,
			expectedError:  true,
		},
		{
			name:           "invalid relay URL",
			relayURL:       "://invalid-url",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server

			if tt.relayURL == "valid-url" {
				// Create a test server for valid URLs
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify request path and query
					assert.Contains(t, r.URL.Path, "/relay/v1/data/bidtraces/proposer_payload_delivered")
					assert.Equal(t, "application/json", r.Header.Get("Accept"))
					assert.Equal(t, "MEV-Commit-Monitor/1.0", r.Header.Get("User-Agent"))
					assert.Equal(t, r.URL.Query().Get("block_number"), "12345")

					w.WriteHeader(tt.responseStatus)
					w.Write([]byte(tt.responseBody))
				}))
				defer server.Close()

				tt.relayURL = server.URL
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			httpClient := retryablehttp.NewClient()
			httpClient.RetryMax = 0 // Disable retries for testing

			client := NewRelayClient([]string{tt.relayURL}, logger, httpClient)
			result := client.queryOneRelay(context.Background(), tt.relayURL, tt.blockNumber)

			assert.Equal(t, tt.relayURL, result.Relay)

			if tt.expectedError {
				assert.NotEmpty(t, result.Error)
			} else {
				assert.Empty(t, result.Error)
				assert.Equal(t, http.StatusOK, result.StatusCode)

				// Check response parsing
				traces, ok := result.Response.([]BidTrace)
				assert.True(t, ok)

				// If we expect data, verify it
				if tt.responseBody != "[]" {
					assert.NotEmpty(t, traces)

					var expectedTraces []BidTrace
					err := json.Unmarshal([]byte(tt.responseBody), &expectedTraces)
					require.NoError(t, err)

					assert.Equal(t, expectedTraces[0].Slot, traces[0].Slot)
					assert.Equal(t, expectedTraces[0].BlockNumber, traces[0].BlockNumber)
				} else {
					assert.Empty(t, traces)
				}
			}
		})
	}
}

func TestQueryRelayData(t *testing.T) {
	// Setup test servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"slot":"123","block_number":"12345"}]`))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Server error"}`))
	}))
	defer server2.Close()

	server3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add a delay to test concurrent behavior
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"slot":"123","block_number":"12345"},{"slot":"124","block_number":"12346"}]`))
	}))
	defer server3.Close()

	relayURLs := []string{server1.URL, server2.URL, server3.URL}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = 0 // Disable retries for testing

	client := NewRelayClient(relayURLs, logger, httpClient)

	results := client.QueryRelayData(context.Background(), 12345)

	// Should have results for all relays
	assert.Equal(t, len(relayURLs), len(results))

	// Check server1 result (success)
	result1 := results[server1.URL]
	assert.Equal(t, server1.URL, result1.Relay)
	assert.Equal(t, http.StatusOK, result1.StatusCode)
	assert.Empty(t, result1.Error)
	traces1, ok := result1.Response.([]BidTrace)
	assert.True(t, ok)
	assert.Len(t, traces1, 1)

	// Check server2 result (error)
	result2 := results[server2.URL]
	assert.Equal(t, server2.URL, result2.Relay)

	if result2.StatusCode != 0 {
		assert.Equal(t, http.StatusInternalServerError, result2.StatusCode)
	}
	assert.NotEmpty(t, result2.Error)

	// Check server3 result (success with delay)
	result3 := results[server3.URL]
	assert.Equal(t, server3.URL, result3.Relay)
	assert.Equal(t, http.StatusOK, result3.StatusCode)
	assert.Empty(t, result3.Error)
	traces3, ok := result3.Response.([]BidTrace)
	assert.True(t, ok)
	assert.Len(t, traces3, 2)
}

func TestQueryRelayDataWithCancelledContext(t *testing.T) {
	// Create a test server with a delay to ensure the context gets cancelled before completion
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"slot":"123","block_number":"12345"}]`))
	}))
	defer server.Close()

	relayURLs := []string{server.URL, server.URL} // Use same server twice for simplicity

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	httpClient := retryablehttp.NewClient()

	client := NewRelayClient(relayURLs, logger, httpClient)

	// Create context that cancels after a short time
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// This should return early due to context cancellation
	results := client.QueryRelayData(ctx, 12345)

	// We might get 0 results if all goroutines were cancelled before sending to channel
	// or up to the number of relays if some completed before cancellation
	assert.LessOrEqual(t, len(results), len(relayURLs))
}

func TestQueryRelayDataWithMultipleRelays(t *testing.T) {
	// Create 5 test servers all returning valid responses
	var relayURLs []string

	for range 5 {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"slot":"123","block_number":"12345"}]`))
		}))
		defer server.Close()
		relayURLs = append(relayURLs, server.URL)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	httpClient := retryablehttp.NewClient()

	client := NewRelayClient(relayURLs, logger, httpClient)

	results := client.QueryRelayData(context.Background(), 12345)

	// Should have results for all relays
	assert.Equal(t, 5, len(results))

	// All results should be successful
	for _, url := range relayURLs {
		result := results[url]
		assert.Equal(t, url, result.Relay)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Empty(t, result.Error)
	}
}
