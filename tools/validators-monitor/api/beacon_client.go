package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

// // HTTPDoer defines the interface for HTTP clients
// type HTTPDoer interface {
// 	Do(req *http.Request) (*http.Response, error)
// }

// ProposerDuty represents a validator's proposer duty
type ProposerDuty struct {
	PubKey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}

// ProposerDutiesResponse is the beacon node response structure
type ProposerDutiesResponse struct {
	Data []ProposerDuty `json:"data"`
}

// ProposerDutyInfo contains parsed proposer duty information
type ProposerDutyInfo struct {
	Epoch          uint64
	Slot           uint64
	ValidatorIndex uint64
	PubKey         string
}

// BeaconClient is a client for the Ethereum 2.0 beacon node API
// It uses a retryable HTTP client with exponential backoff.
type BeaconClient struct {
	client  *http.Client
	baseURL *url.URL
	logger  *slog.Logger
}

// NewBeaconClient creates a new beacon node API client.
func NewBeaconClient(baseURL string, logger *slog.Logger, httpClient *http.Client) (*BeaconClient, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	return &BeaconClient{
		client:  httpClient,
		baseURL: parsed,
		logger:  logger,
	}, nil
}

// GetProposerDuties fetches the proposer duties for a given epoch
func (c *BeaconClient) GetProposerDuties(ctx context.Context, epoch uint64) (*ProposerDutiesResponse, error) {
	// build URL
	u := *c.baseURL
	u.Path = path.Join(u.Path, "eth", "v1", "validator", "duties", "proposer", strconv.FormatUint(epoch, 10))
	reqURL := u.String()
	c.logger.Debug(
		"Querying beacon node for proposer duties",
		slog.String("url", reqURL),
		slog.Uint64("epoch", epoch),
	)

	// prepare request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling beacon node: %w", err)
	}
	defer resp.Body.Close()

	logger := c.logger.With(
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("response_time_ms", time.Since(start).Milliseconds()),
	)

	// handle non-200
	if resp.StatusCode != http.StatusOK {
		limit := io.LimitReader(resp.Body, 512)
		msg, _ := io.ReadAll(limit)
		logger.Error(
			"Non-OK status fetching proposer duties",
			slog.String("body_snippet", string(msg)),
		)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// decode JSON
	var dutiesResp ProposerDutiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&dutiesResp); err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}

	logger.Debug(
		"Fetched proposer duties",
		slog.Int("count", len(dutiesResp.Data)),
	)
	return &dutiesResp, nil
}

// ParseProposerDuties converts the API response to ProposerDutyInfo slice
func ParseProposerDuties(epoch uint64, resp *ProposerDutiesResponse) ([]ProposerDutyInfo, error) {
	duties := make([]ProposerDutyInfo, 0, len(resp.Data))
	for _, d := range resp.Data {
		slot, err := strconv.ParseUint(d.Slot, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing slot %q: %w", d.Slot, err)
		}

		idx, err := strconv.ParseUint(d.ValidatorIndex, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing validator index %q: %w", d.ValidatorIndex, err)
		}

		duties = append(duties, ProposerDutyInfo{
			Epoch:          epoch,
			Slot:           slot,
			ValidatorIndex: idx,
			PubKey:         d.PubKey,
		})
	}
	return duties, nil
}

// GetBlockBySlot fetches the block root for a given slot
func (c *BeaconClient) GetBlockBySlot(ctx context.Context, slot uint64) (string, error) {
	// build URL
	u := *c.baseURL
	u.Path = path.Join(u.Path, "eth", "v2", "beacon", "blocks", strconv.FormatUint(slot, 10))
	reqURL := u.String()

	c.logger.Debug(
		"Querying beacon node for block by slot",
		slog.String("url", reqURL),
		slog.Uint64("slot", slot),
	)

	// prepare request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("calling beacon node: %w", err)
	}
	defer resp.Body.Close()

	logger := c.logger.With(
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("response_time_ms", time.Since(start).Milliseconds()),
	)

	// 404 means block not found
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if resp.StatusCode != http.StatusOK {
		limit := io.LimitReader(resp.Body, 512)
		msg, _ := io.ReadAll(limit)
		logger.Error(
			"Non-OK status fetching block",
			slog.String("body_snippet", string(msg)),
		)
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// decode into nested struct
	var wrapper struct {
		Data struct {
			Message struct {
				Body struct {
					ExecutionPayload struct {
						BlockNumber string `json:"block_number"`
					} `json:"execution_payload"`
				} `json:"body"`
			} `json:"message"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return "", fmt.Errorf("decoding JSON response: %w", err)
	}

	blockNumber := wrapper.Data.Message.Body.ExecutionPayload.BlockNumber
	logger.Debug(
		"Fetched block info",
		slog.String("block_number", blockNumber),
	)
	return blockNumber, nil
}
