package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// PayloadResponse represents the API response for payload requests
type PayloadResponse struct {
	PayloadID        string `json:"payload_id"`
	ExecutionPayload string `json:"execution_payload"`
	BlockHeight      uint64 `json:"block_height"`
	Timestamp        int64  `json:"timestamp"`
}

// PayloadListResponse represents the response for multiple payloads
type PayloadListResponse struct {
	Payloads   []PayloadResponse `json:"payloads"`
	HasMore    bool              `json:"has_more"`
	NextHeight uint64            `json:"next_height,omitempty"`
	TotalCount int               `json:"total_count"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// PayloadClient handles communication with the leader node's API
type PayloadClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewPayloadClient creates a new payload API client
func NewPayloadClient(baseURL string, logger *slog.Logger) *PayloadClient {
	return &PayloadClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.With("component", "PayloadClient"),
	}
}

// GetLatestPayload fetches the latest payload from the leader node
func (pc *PayloadClient) GetLatestPayload(ctx context.Context) (*PayloadResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payload/latest", pc.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errorResp.Message)
	}

	var payloadResp PayloadResponse
	if err := json.Unmarshal(body, &payloadResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	pc.logger.Debug(
		"Retrieved payload from leader",
		"payload_id", payloadResp.PayloadID,
		"block_height", payloadResp.BlockHeight,
	)
	return &payloadResp, nil
}

// GetPayloadsSince fetches payloads with block height >= sinceHeight from the leader node
func (pc *PayloadClient) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) (*PayloadListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payload/since/%d?limit=%d", pc.baseURL, sinceHeight, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errorResp.Message)
	}

	var payloadListResp PayloadListResponse
	if err := json.Unmarshal(body, &payloadListResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	pc.logger.Debug(
		"Retrieved payloads since height from leader",
		"since_height", sinceHeight,
		"count", len(payloadListResp.Payloads),
		"has_more", payloadListResp.HasMore,
		"next_height", payloadListResp.NextHeight,
	)

	return &payloadListResp, nil
}

// GetPayloadByHeight fetches a specific payload by block height from the leader node
func (pc *PayloadClient) GetPayloadByHeight(ctx context.Context, height uint64) (*PayloadResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payload/height/%d", pc.baseURL, height)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errorResp.Message)
	}

	var payloadResp PayloadResponse
	if err := json.Unmarshal(body, &payloadResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	pc.logger.Debug(
		"Retrieved payload by height from leader",
		"height", height,
		"payload_id", payloadResp.PayloadID,
	)
	return &payloadResp, nil
}

// CheckHealth checks if the leader node API is healthy
func (pc *PayloadClient) CheckHealth(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/health", pc.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute health check: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("leader node unhealthy (status %d)", resp.StatusCode)
	}

	return nil
}
