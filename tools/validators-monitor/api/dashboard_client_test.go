package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDashboardClient(t *testing.T, handler http.Handler) *DashboardClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	httpClient := &http.Client{Timeout: 10 * time.Second}

	dashboardClient, err := NewDashboardClient(server.URL, logger, httpClient)
	require.NoError(t, err)

	return dashboardClient
}

func TestNewDashboardClient(t *testing.T) {
	// Test with valid URL
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	httpClient := &http.Client{Timeout: 10 * time.Second}

	dashboardClient, err := NewDashboardClient("http://dashboard.example.com", logger, httpClient)
	assert.NoError(t, err)
	assert.NotNil(t, dashboardClient)
	assert.Equal(t, "http://dashboard.example.com", dashboardClient.baseURL.String())

	// Test with invalid URL
	dashboardClient, err = NewDashboardClient("://invalid-url", logger, httpClient)
	assert.Error(t, err)
	assert.Nil(t, dashboardClient)
}

func TestGetBlockInfo(t *testing.T) {
	tests := []struct {
		name           string
		blockNumber    uint64
		responseStatus int
		responseBody   string
		expectedError  bool
		expectedResult *DashboardResponse
	}{
		{
			name:           "successful response",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody: `{
				"number": 12345,
				"winner": "0xabcdef1234567890",
				"window": 42,
				"total_opened_commitments": 10,
				"total_rewards": 5,
				"total_slashes": 2,
				"total_amount": "123.45"
			}`,
			expectedError: false,
			expectedResult: &DashboardResponse{
				Number:                 12345,
				Winner:                 "0xabcdef1234567890",
				Window:                 42,
				TotalOpenedCommitments: 10,
				TotalRewards:           5,
				TotalSlashes:           2,
				TotalAmount:            "123.45",
			},
		},
		{
			name:           "server error",
			blockNumber:    12345,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "Internal server error"}`,
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:           "not found error",
			blockNumber:    12345,
			responseStatus: http.StatusNotFound,
			responseBody:   `{"error": "Block not found"}`,
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:           "invalid json response",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:           "empty response",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `{}`,
			expectedError:  false,
			expectedResult: &DashboardResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request path
				expectedPath := "/block/" + strconv.FormatUint(tt.blockNumber, 10)
				assert.Equal(t, expectedPath, r.URL.Path)

				// Check request method
				assert.Equal(t, http.MethodGet, r.Method)

				// Check headers
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				assert.Equal(t, "MEV-Commit-Monitor/1.0", r.Header.Get("User-Agent"))

				// Send response
				w.WriteHeader(tt.responseStatus)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			})

			client := setupTestDashboardClient(t, handler)
			result, err := client.GetBlockInfo(context.Background(), tt.blockNumber)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Number, result.Number)
				assert.Equal(t, tt.expectedResult.Winner, result.Winner)
				assert.Equal(t, tt.expectedResult.Window, result.Window)
				assert.Equal(t, tt.expectedResult.TotalOpenedCommitments, result.TotalOpenedCommitments)
				assert.Equal(t, tt.expectedResult.TotalRewards, result.TotalRewards)
				assert.Equal(t, tt.expectedResult.TotalSlashes, result.TotalSlashes)
				assert.Equal(t, tt.expectedResult.TotalAmount, result.TotalAmount)
			}
		})
	}
}

func TestGetCommitmentsByBlock(t *testing.T) {
	tests := []struct {
		name           string
		blockNumber    uint64
		responseStatus int
		responseBody   string
		expectedError  bool
		expectedLength int
		checkFirstItem bool // whether to verify the first commitment's details
	}{
		{
			name:           "successful response with commitments",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody: `[
				{
					"commitment_index": [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32],
					"bidder": "0xbidder1",
					"committer": "0xcommitter1",
					"bid_amt": "1000000000000000000",
					"slash_amt": "500000000000000000",
					"block_number": 12345,
					"decay_start_time_stamp": 1000,
					"decay_end_time_stamp": 2000,
					"txn_hash": "0xtxhash1",
					"reverting_tx_hashes": "hash1,hash2",
					"commitment_digest": [33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64],
					"dispatch_timestamp": 3000
				},
				{
					"commitment_index": [65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96],
					"bidder": "0xbidder2",
					"committer": "0xcommitter2",
					"bid_amt": "2000000000000000000",
					"slash_amt": "1000000000000000000",
					"block_number": 12345,
					"decay_start_time_stamp": 1500,
					"decay_end_time_stamp": 2500,
					"txn_hash": "0xtxhash2",
					"reverting_tx_hashes": "hash3",
					"commitment_digest": [97,98,99,100,101,102,103,104,105,106,107,108,109,110,111,112,113,114,115,116,117,118,119,120,121,122,123,124,125,126,127,128],
					"dispatch_timestamp": 3500
				}
			]`,
			expectedError:  false,
			expectedLength: 2,
			checkFirstItem: true,
		},
		{
			name:           "successful response with empty array",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `[]`,
			expectedError:  false,
			expectedLength: 0,
			checkFirstItem: false,
		},
		{
			name:           "not found response",
			blockNumber:    12345,
			responseStatus: http.StatusNotFound,
			responseBody:   `{"error": "No commitments found for this block"}`,
			expectedError:  false, // We treat not found as an empty result
			expectedLength: 0,
			checkFirstItem: false,
		},
		{
			name:           "server error",
			blockNumber:    12345,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"error": "Internal server error"}`,
			expectedError:  true,
			expectedLength: 0,
			checkFirstItem: false,
		},
		{
			name:           "invalid json response",
			blockNumber:    12345,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectedError:  true,
			expectedLength: 0,
			checkFirstItem: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request path
				expectedPath := "/block/" + strconv.FormatUint(tt.blockNumber, 10) + "/commitments"
				assert.Equal(t, expectedPath, r.URL.Path)

				// Check request method
				assert.Equal(t, http.MethodGet, r.Method)

				// Check headers
				assert.Equal(t, "application/json", r.Header.Get("Accept"))
				assert.Equal(t, "MEV-Commit-Monitor/1.0", r.Header.Get("User-Agent"))

				// Send response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)
				_, err := w.Write([]byte(tt.responseBody))
				require.NoError(t, err)
			})

			client := setupTestDashboardClient(t, handler)
			result, err := client.GetCommitmentsByBlock(context.Background(), tt.blockNumber)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedLength, len(result))

				if tt.checkFirstItem && len(result) > 0 {
					// Check the first item's details
					assert.Equal(t, "0xbidder1", result[0].Bidder)
					assert.Equal(t, "0xcommitter1", result[0].Committer)
					assert.Equal(t, "1000000000000000000", result[0].BidAmt)
					assert.Equal(t, "500000000000000000", result[0].SlashAmt)
					assert.Equal(t, uint64(12345), result[0].BlockNumber)
					assert.Equal(t, uint64(1000), result[0].DecayStartTimeStamp)
					assert.Equal(t, uint64(2000), result[0].DecayEndTimeStamp)
					assert.Equal(t, "0xtxhash1", result[0].TxnHash)
					assert.Equal(t, "hash1,hash2", result[0].RevertingTxHashes)
					assert.Equal(t, uint64(3000), result[0].DispatchTimestamp)

					// Check the byte arrays
					expectedIndex := [32]byte{}
					for i := range 32 {
						expectedIndex[i] = byte(i + 1)
					}
					assert.Equal(t, expectedIndex, result[0].CommitmentIndex)

					expectedDigest := [32]byte{}
					for i := range 32 {
						expectedDigest[i] = byte(i + 33)
					}
					assert.Equal(t, expectedDigest, result[0].CommitmentDigest)
				}
			}
		})
	}
}
func TestGetCommitmentsByBlock_ContextCancellation(t *testing.T) {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Just return success - the context cancellation is handled by the HTTP client
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`[{"commitment_index": "0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"}]`))
		require.NoError(t, err)
	})

	client := setupTestDashboardClient(t, handler)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel it immediately to simulate a timeout/cancellation
	cancel()

	// The request should fail due to canceled context
	result, err := client.GetCommitmentsByBlock(ctx, 12345)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetCommitmentsByBlock_NetworkErrors(t *testing.T) {
	// Create a client with an unreachable URL
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	httpClient := &http.Client{Timeout: 10 * time.Second}

	dashboardClient, err := NewDashboardClient("http://localhost:12345", logger, httpClient)
	require.NoError(t, err)

	// The request should fail due to connection refused
	result, err := dashboardClient.GetCommitmentsByBlock(context.Background(), 12345)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetCommitmentsByBlock_HeaderHandling(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that Accept and User-Agent headers are set
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "MEV-Commit-Monitor/1.0", r.Header.Get("User-Agent"))

		// Return a valid response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`[]`))
		require.NoError(t, err)
	})

	client := setupTestDashboardClient(t, handler)
	result, err := client.GetCommitmentsByBlock(context.Background(), 12345)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestContextCancellation(t *testing.T) {
	// Create a simple handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Just return success - the context cancellation is handled by the HTTP client
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"number": 12345}`))
		require.NoError(t, err)
	})

	client := setupTestDashboardClient(t, handler)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel it immediately to simulate a timeout/cancellation
	cancel()

	// The request should fail due to canceled context
	result, err := client.GetBlockInfo(ctx, 12345)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNetworkErrors(t *testing.T) {
	// Create a client with an unreachable URL
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	httpClient := &http.Client{Timeout: 10 * time.Second}

	dashboardClient, err := NewDashboardClient("http://localhost:12345", logger, httpClient)
	require.NoError(t, err)

	// The request should fail due to connection refused
	result, err := dashboardClient.GetBlockInfo(context.Background(), 12345)
	assert.Error(t, err)
	assert.Nil(t, result)
}
