package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"fmt"
)

func setupTestClient(t *testing.T, handler http.Handler) *BeaconClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client := retryablehttp.NewClient()
	client.RetryMax = 1

	beaconClient, err := NewBeaconClient(server.URL, logger, client)
	require.NoError(t, err)

	return beaconClient
}

func TestNewBeaconClient(t *testing.T) {
	// Test with valid URL
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := retryablehttp.NewClient()

	beaconClient, err := NewBeaconClient("http://localhost:5052", logger, client)
	assert.NoError(t, err)
	assert.NotNil(t, beaconClient)
	assert.Equal(t, "http://localhost:5052", beaconClient.baseURL.String())

	// Test with invalid URL
	beaconClient, err = NewBeaconClient("://invalid-url", logger, client)
	assert.Error(t, err)
	assert.Nil(t, beaconClient)
}

func TestGetProposerDuties(t *testing.T) {
	tests := []struct {
		name           string
		epoch          uint64
		responseStatus int
		responseBody   string
		expectedError  bool
		expectedCount  int
	}{
		{
			name:           "successful response",
			epoch:          123,
			responseStatus: http.StatusOK,
			responseBody:   `{"data":[{"pubkey":"0x123456","validator_index":"42","slot":"1234"},{"pubkey":"0x789012","validator_index":"43","slot":"1235"}]}`,
			expectedError:  false,
			expectedCount:  2,
		},
		{
			name:           "empty response",
			epoch:          123,
			responseStatus: http.StatusOK,
			responseBody:   `{"data":[]}`,
			expectedError:  false,
			expectedCount:  0,
		},
		{
			name:           "server error",
			epoch:          123,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"message":"Internal server error"}`,
			expectedError:  true,
			expectedCount:  0,
		},
		{
			name:           "invalid json response",
			epoch:          123,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectedError:  true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request path and method
				assert.Equal(t, "/eth/v1/validator/duties/proposer/"+fmt.Sprint(tt.epoch), r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			})

			client := setupTestClient(t, handler)
			resp, err := client.GetProposerDuties(context.Background(), tt.epoch)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedCount, len(resp.Data))
			}
		})
	}
}

func TestParseProposerDuties(t *testing.T) {
	tests := []struct {
		name          string
		epoch         uint64
		input         *ProposerDutiesResponse
		expectedCount int
		expectedError bool
	}{
		{
			name:  "valid duties",
			epoch: 123,
			input: &ProposerDutiesResponse{
				Data: []ProposerDuty{
					{PubKey: "0xabc123", ValidatorIndex: "42", Slot: "1234"},
					{PubKey: "0xdef456", ValidatorIndex: "43", Slot: "1235"},
				},
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "empty duties",
			epoch:         123,
			input:         &ProposerDutiesResponse{Data: []ProposerDuty{}},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:  "invalid slot",
			epoch: 123,
			input: &ProposerDutiesResponse{
				Data: []ProposerDuty{
					{PubKey: "0xabc123", ValidatorIndex: "42", Slot: "invalid"},
				},
			},
			expectedCount: 0,
			expectedError: true,
		},
		{
			name:  "invalid validator index",
			epoch: 123,
			input: &ProposerDutiesResponse{
				Data: []ProposerDuty{
					{PubKey: "0xabc123", ValidatorIndex: "invalid", Slot: "1234"},
				},
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duties, err := ParseProposerDuties(tt.epoch, tt.input)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, duties)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(duties))

				if tt.expectedCount > 0 {
					// Check first duty values
					assert.Equal(t, tt.epoch, duties[0].Epoch)
					assert.Equal(t, tt.input.Data[0].PubKey, duties[0].PubKey)

					// Check parsed values
					slotValue, _ := json.Number(tt.input.Data[0].Slot).Int64()
					assert.Equal(t, uint64(slotValue), duties[0].Slot)

					validatorValue, _ := json.Number(tt.input.Data[0].ValidatorIndex).Int64()
					assert.Equal(t, uint64(validatorValue), duties[0].ValidatorIndex)
				}
			}
		})
	}
}

func TestGetBlockBySlot(t *testing.T) {
	tests := []struct {
		name           string
		slot           uint64
		responseStatus int
		responseBody   string
		expectedResult string
		expectedError  bool
	}{
		{
			name:           "successful response",
			slot:           1234,
			responseStatus: http.StatusOK,
			responseBody: `{
				"data": {
					"message": {
						"body": {
							"execution_payload": {
								"block_number": "12345678"
							}
						}
					}
				}
			}`,
			expectedResult: "12345678",
			expectedError:  false,
		},
		{
			name:           "block not found",
			slot:           1234,
			responseStatus: http.StatusNotFound,
			responseBody:   `{"message":"Block not found"}`,
			expectedResult: "",
			expectedError:  false,
		},
		{
			name:           "server error",
			slot:           1234,
			responseStatus: http.StatusInternalServerError,
			responseBody:   `{"message":"Internal server error"}`,
			expectedResult: "",
			expectedError:  true,
		},
		{
			name:           "invalid json response",
			slot:           1234,
			responseStatus: http.StatusOK,
			responseBody:   `{invalid json}`,
			expectedResult: "",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request path and method
				assert.Equal(t, "/eth/v2/beacon/blocks/"+fmt.Sprint(tt.slot), r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Accept"))

				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			})

			client := setupTestClient(t, handler)
			blockNumber, err := client.GetBlockBySlot(context.Background(), tt.slot)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, blockNumber)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, blockNumber)
			}
		})
	}
}
