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
)

// DashboardResponse represents the response from the dashboard service
type DashboardResponse struct {
	Number                 uint64 `json:"number"`
	Winner                 string `json:"winner"`
	Window                 uint64 `json:"window"`
	TotalOpenedCommitments int    `json:"total_opened_commitments"`
	TotalRewards           int    `json:"total_rewards"`
	TotalSlashes           int    `json:"total_slashes"`
	TotalAmount            string `json:"total_amount"`
}

// CommitmentData represents a commitment from the dashboard service
type CommitmentData struct {
	CommitmentIndex     [32]byte `json:"commitment_index"`
	Bidder              string   `json:"bidder"`
	Committer           string   `json:"committer"`
	BidAmt              string   `json:"bid_amt"`
	SlashAmt            string   `json:"slash_amt"`
	BlockNumber         uint64   `json:"block_number"`
	DecayStartTimeStamp uint64   `json:"decay_start_time_stamp"`
	DecayEndTimeStamp   uint64   `json:"decay_end_time_stamp"`
	TxnHash             string   `json:"txn_hash"`
	RevertingTxHashes   string   `json:"reverting_tx_hashes"`
	CommitmentDigest    [32]byte `json:"commitment_digest"`
	DispatchTimestamp   uint64   `json:"dispatch_timestamp"`
}

// DashboardClient is a client for interacting with the dashboard service API
// It uses a retryable HTTP client under the hood.
type DashboardClient struct {
	client  *http.Client
	baseURL *url.URL
	logger  *slog.Logger
}

// NewDashboardClient creates a new dashboard client.
func NewDashboardClient(
	baseURL string,
	logger *slog.Logger,
	httpClient *http.Client,
) (*DashboardClient, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	return &DashboardClient{
		client:  httpClient,
		baseURL: parsed,
		logger:  logger,
	}, nil
}

// GetBlockInfo queries the dashboard service for block information
func (c *DashboardClient) GetBlockInfo(
	ctx context.Context,
	blockNumber uint64,
) (*DashboardResponse, error) {
	// build request URL
	u := *c.baseURL
	u.Path = path.Join(u.Path, "block", strconv.FormatUint(blockNumber, 10))
	reqURL := u.String()

	c.logger.Debug(
		"Querying dashboard service for block",
		slog.Uint64("block_number", blockNumber),
		slog.String("url", reqURL),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MEV-Commit-Monitor/1.0")

	// execute with retries
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling dashboard service: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	// check HTTP status
	if resp.StatusCode != http.StatusOK {
		// limit error-body read
		limit := io.LimitReader(resp.Body, 512)
		msg, _ := io.ReadAll(limit)
		return nil, fmt.Errorf("dashboard service %d: %s", resp.StatusCode, string(msg))
	}

	// decode JSON directly
	var dr DashboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}

	c.logger.Debug(
		"Dashboard service response received",
		slog.Uint64("block_number", dr.Number),
		slog.String("winner", dr.Winner),
		slog.Int("total_opened_commitments", dr.TotalOpenedCommitments),
		slog.Int("total_rewards", dr.TotalRewards),
		slog.Int("total_slashes", dr.TotalSlashes),
		slog.String("total_amount", dr.TotalAmount),
	)

	return &dr, nil
}

// GetCommitmentsByBlock queries the dashboard service for commitments by block number
func (c *DashboardClient) GetCommitmentsByBlock(
	ctx context.Context,
	blockNumber uint64,
) ([]CommitmentData, error) {
	// build request URL
	u := *c.baseURL
	u.Path = path.Join(u.Path, "block", strconv.FormatUint(blockNumber, 10), "commitments")
	reqURL := u.String()

	c.logger.Debug(
		"Querying dashboard service for block commitments",
		slog.Uint64("block_number", blockNumber),
		slog.String("url", reqURL),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MEV-Commit-Monitor/1.0")

	// execute with retries
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling dashboard service: %w", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	// check HTTP status
	if resp.StatusCode == http.StatusNotFound {
		// No commitments found for this block
		c.logger.Debug(
			"No commitments found for block",
			slog.Uint64("block_number", blockNumber),
		)
		return []CommitmentData{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		// limit error-body read
		limit := io.LimitReader(resp.Body, 512)
		msg, _ := io.ReadAll(limit)
		return nil, fmt.Errorf("dashboard service %d: %s", resp.StatusCode, string(msg))
	}

	// decode JSON directly
	var commitments []CommitmentData
	if err := json.NewDecoder(resp.Body).Decode(&commitments); err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}

	c.logger.Debug(
		"Dashboard service commitments received",
		slog.Uint64("block_number", blockNumber),
		slog.Int("commitment_count", len(commitments)),
	)

	return commitments, nil
}
