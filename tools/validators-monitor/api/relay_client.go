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

// RelayClient queries multiple builder relays with retry/backoff
// It uses a retryable HTTP client under the hood.
type RelayClient struct {
	client    *http.Client
	relayURLs []string
	logger    *slog.Logger
}

// NewRelayClient constructs a RelayClient. If httpClient is nil, a default retryablehttp.Client is used.
// HTTP timeouts and retry policies are configurable on the client.
func NewRelayClient(relayURLs []string, logger *slog.Logger, httpClient *http.Client) *RelayClient {
	return &RelayClient{client: httpClient, relayURLs: relayURLs, logger: logger}
}

// QueryRelayData concurrently queries all relays for bid traces at a block
func (c *RelayClient) QueryRelayData(ctx context.Context, blockNumber uint64) map[string]RelayResult {
	c.logger.Debug("querying relays for block",
		slog.Uint64("block_number", blockNumber),
	)

	results := make(map[string]RelayResult, len(c.relayURLs))
	resultCh := make(chan RelayResult, len(c.relayURLs))

	// spawn one goroutine per relay
	for _, relay := range c.relayURLs {
		r := relay
		go func() {
			res := c.queryOneRelay(ctx, r, blockNumber)
			select {
			case resultCh <- res:
			case <-ctx.Done():
			}
		}()
	}

	// collect
	for range c.relayURLs {
		select {
		case res := <-resultCh:
			results[res.Relay] = res
		case <-ctx.Done():
			return results
		}
	}

	return results
}

// queryOneRelay performs a single relay request with retries and backoff
func (c *RelayClient) queryOneRelay(ctx context.Context, relayURL string, blockNumber uint64) RelayResult {
	result := RelayResult{Relay: relayURL}

	// build URL
	u, err := url.Parse(relayURL)
	if err != nil {
		result.Error = fmt.Sprintf("invalid relay URL: %v", err)
		return result
	}
	u.Path = path.Join(u.Path, "relay", "v1", "data", "bidtraces", "proposer_payload_delivered")
	q := u.Query()
	q.Set("block_number", strconv.FormatUint(blockNumber, 10))
	u.RawQuery = q.Encode()
	reqURL := u.String()

	c.logger.Debug("querying relay",
		slog.String("relay", relayURL),
		slog.String("url", reqURL),
	)

	// prepare request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("building request: %v", err)
		return result
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "MEV-Commit-Monitor/1.0")

	// execute
	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	// attach status and timing
	result.StatusCode = resp.StatusCode
	c.logger.Debug("relay response",
		slog.String("relay", relayURL),
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("latency_ms", time.Since(start).Milliseconds()),
	)

	// handle non-200
	if resp.StatusCode != http.StatusOK {
		limit := io.LimitReader(resp.Body, 512)
		msg, _ := io.ReadAll(limit)
		result.Error = fmt.Sprintf("status %d: %s", resp.StatusCode, string(msg))
		return result
	}

	// decode JSON
	var traces []BidTrace
	if err := json.NewDecoder(resp.Body).Decode(&traces); err != nil {
		result.Error = fmt.Sprintf("parsing JSON: %v", err)
		return result
	}

	result.Response = traces
	c.logger.Debug("parsed bid traces",
		slog.String("relay", relayURL),
		slog.Int("count", len(traces)),
	)

	return result
}
