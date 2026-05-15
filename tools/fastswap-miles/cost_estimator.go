package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Percentile of recent fastswap bid_cost used as the per-sweep bid proxy
// (the sweep tx is itself a fastswap). p75 = under-promise.
const sweepBidGlobalPercentile = 0.75

// Fallback bid used when the percentile query returns NULL (cold start /
// no processed rows yet). Post-fix realized p75 ≈ 4e-5 ETH.
const sweepBidFallbackEth = 4e-5

// costEstimateLookbackDays is the rolling window over which per-token sweep
// overhead percentiles are computed.
const costEstimateLookbackDays = 14

// costEstimateRefreshInterval is the period between background refreshes of
// the per-token estimate cache.
const costEstimateRefreshInterval = 30 * time.Minute

// costEstimateMinSweeps is the minimum number of recent sweeps required for
// the "primary" percentile (p25 by default). Below this threshold the
// estimator falls back to a more conservative p75 to compensate for noisy
// data on low-volume tokens.
const costEstimateMinSweeps = 10

// costEstimateLastResort is the hardcoded fallback overhead used when a token
// has no historical sweep data at all (brand-new token). Conservative —
// chosen so users still get *some* miles on the first swap of a novel token,
// while leaving room for protocol margin once real data arrives.
const costEstimateLastResort = 0.001 // ETH

// costEstimate holds the per-token sweep overhead estimate (in ETH) used for
// upfront miles awarding. Estimates are refreshed periodically from
// fastswap_miles realized sweep data.
type costEstimate struct {
	// PerRowOverhead is the estimated sweep overhead per user row, in ETH —
	// the sum of pro-rata sweep gas and pro-rata sweep bid. This is the value
	// subtracted in the miles formula in lieu of realized values for the
	// (still-pending) sweep that will eventually convert this row's surplus
	// tokens to ETH. See Refresh for the two components.
	PerRowOverhead float64

	// Source describes how this estimate was computed (for observability).
	// One of: "p25", "p75_low_data", "default_no_data".
	Source string

	// SweepCount is the number of sweeps the estimate is based on.
	// Zero indicates the default-no-data fallback.
	SweepCount int

	// ComputedAt is when this estimate was last computed.
	ComputedAt time.Time
}

// costEstimator maintains an in-memory, periodically refreshed cache of
// per-token cost estimates. It is safe for concurrent use.
type costEstimator struct {
	db     *sql.DB
	logger *slog.Logger

	// Lowercased hex of the executor address and WETH address. Used to
	// identify executor sweep rows (user_address = executor, output = WETH)
	// when computing the per-token sweep bid contribution.
	executorAddr string
	wethAddr     string

	mu        sync.RWMutex
	estimates map[string]costEstimate // key: lowercased token hex
	lastFresh time.Time
}

func newCostEstimator(db *sql.DB, logger *slog.Logger, executorAddr, wethAddr common.Address) *costEstimator {
	return &costEstimator{
		db:           db,
		logger:       logger,
		executorAddr: strings.ToLower(executorAddr.Hex()),
		wethAddr:     strings.ToLower(wethAddr.Hex()),
		estimates:    make(map[string]costEstimate),
	}
}

// Get returns the current estimate for a token, computing it from cache.
// If no estimate exists for this token, returns the default-no-data fallback.
func (c *costEstimator) Get(token string) costEstimate {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if est, ok := c.estimates[strings.ToLower(token)]; ok {
		return est
	}
	return costEstimate{
		PerRowOverhead: costEstimateLastResort,
		Source:         "default_no_data",
		SweepCount:     0,
		ComputedAt:     time.Now(),
	}
}

// Refresh recomputes per-token estimates from realized fastswap_miles data
// over the configured lookback window. PerRowOverhead is the sum of two
// terms: per-token p25/p75 of pro-rata sweep gas (existing) plus per-token
// (n_sweeps × global_bid_p75 / n_user_rows) for the sweep tx's own bid.
// Both scale together with batch size.
func (c *costEstimator) Refresh(ctx context.Context) error {
	// Filter to ETH-input rows so per_row_oh isolates pure sweep_overhead.
	// ERC20-input rows have user_gas baked into (surplus_eth - net_profit_eth),
	// which would inflate the estimate and cause the miles formula (which
	// deducts user_gas separately) to over-deduct.
	rows, err := c.db.QueryContext(ctx, fmt.Sprintf(`
SELECT output_token,
       COUNT(*) as n,
       percentile_approx(per_row_oh, 0.25) as p25,
       percentile_approx(per_row_oh, 0.75) as p75
FROM (
  SELECT output_token,
         surplus_eth - net_profit_eth - CAST(bid_cost AS DOUBLE)/1e18 as per_row_oh
  FROM mevcommit_57173.fastswap_miles
  WHERE processed = 1
    AND swap_type = 'erc20'
    AND miles > 0
    AND surplus_eth > 0
    AND surplus_eth IS NOT NULL
    AND surplus_eth < 1.0
    AND input_token = '0x0000000000000000000000000000000000000000'
    AND block_timestamp >= NOW() - INTERVAL %d DAY
) t
GROUP BY output_token`, costEstimateLookbackDays))
	if err != nil {
		return fmt.Errorf("query cost estimates: %w", err)
	}
	defer func() { _ = rows.Close() }()

	fresh := make(map[string]costEstimate)
	for rows.Next() {
		var token string
		var n int
		var p25, p75 float64
		if err := rows.Scan(&token, &n, &p25, &p75); err != nil {
			c.logger.Warn("cost estimate scan failed", slog.Any("error", err))
			continue
		}

		var overhead float64
		var source string
		if n >= costEstimateMinSweeps {
			overhead = p25
			source = "p25"
		} else {
			overhead = p75
			source = "p75_low_data"
		}

		fresh[strings.ToLower(token)] = costEstimate{
			PerRowOverhead: overhead,
			Source:         source,
			SweepCount:     n,
			ComputedAt:     time.Now(),
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate cost estimates: %w", err)
	}

	// Fold per-token sweep bid contribution into PerRowOverhead. A failure
	// here is logged but non-fatal — the gas-only overhead is still a usable
	// estimate and the alternative would be skipping the refresh entirely.
	bidByToken, err := c.computePerTokenSweepBidEth(ctx)
	if err != nil {
		c.logger.Warn("sweep bid contribution refresh failed; falling back to gas-only overhead",
			slog.Any("error", err))
	} else {
		for token, est := range fresh {
			if bid, ok := bidByToken[token]; ok && bid > 0 {
				est.PerRowOverhead += bid
				fresh[token] = est
			}
		}
	}

	c.mu.Lock()
	c.estimates = fresh
	c.lastFresh = time.Now()
	c.mu.Unlock()

	c.logger.Info("cost estimates refreshed",
		slog.Int("tokens", len(fresh)),
		slog.Int("tokens_with_sweep_bid", len(bidByToken)),
		slog.Duration("window", costEstimateLookbackDays*24*time.Hour))
	return nil
}

// computePerTokenSweepBidEth returns per-row sweep bid contribution (ETH)
// keyed by lowercased output_token. Single round trip: joins per-token
// executor sweep counts × per-token user-row counts × global bid p75 (with
// fallback when no processed rows exist).
func (c *costEstimator) computePerTokenSweepBidEth(ctx context.Context) (map[string]float64, error) {
	query := fmt.Sprintf(`
SELECT s.token, s.n_sweeps, u.n_users, COALESCE(b.p, %f) AS bid_p75
FROM (
  SELECT LOWER(input_token) AS token, COUNT(*) AS n_sweeps
  FROM mevcommit_57173.fastswap_miles
  WHERE LOWER(user_address) = ?
    AND swap_type = 'eth_weth'
    AND LOWER(output_token) = ?
    AND block_timestamp >= NOW() - INTERVAL %d DAY
  GROUP BY input_token
) s
JOIN (
  SELECT LOWER(output_token) AS token, COUNT(*) AS n_users
  FROM mevcommit_57173.fastswap_miles
  WHERE swap_type = 'erc20'
    AND LOWER(user_address) != ?
    AND block_timestamp >= NOW() - INTERVAL %d DAY
  GROUP BY output_token
) u ON u.token = s.token
CROSS JOIN (
  SELECT percentile_approx(CAST(bid_cost AS DOUBLE)/1e18, %f) AS p
  FROM mevcommit_57173.fastswap_miles
  WHERE processed = 1
    AND bid_cost IS NOT NULL
    AND CAST(bid_cost AS DOUBLE) > 0
    AND block_timestamp >= NOW() - INTERVAL %d DAY
) b
`, sweepBidFallbackEth, costEstimateLookbackDays, costEstimateLookbackDays,
		sweepBidGlobalPercentile, costEstimateLookbackDays)

	rows, err := c.db.QueryContext(ctx, query, c.executorAddr, c.wethAddr, c.executorAddr)
	if err != nil {
		return nil, fmt.Errorf("query per-token sweep bid: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make(map[string]float64)
	for rows.Next() {
		var token string
		var nSweeps, nUsers int
		var bidEth float64
		if err := rows.Scan(&token, &nSweeps, &nUsers, &bidEth); err != nil {
			c.logger.Warn("scan per-token sweep bid failed", slog.Any("error", err))
			continue
		}
		if nSweeps <= 0 || nUsers <= 0 || bidEth <= 0 {
			continue
		}
		out[token] = float64(nSweeps) * bidEth / float64(nUsers)
	}
	return out, rows.Err()
}

// Run starts a background loop that refreshes estimates on the configured
// interval. Returns when the context is cancelled. Performs an immediate
// initial refresh on startup so estimates are warm before the first miles
// computation.
func (c *costEstimator) Run(ctx context.Context) {
	if err := c.Refresh(ctx); err != nil {
		c.logger.Warn("initial cost estimate refresh failed; using defaults", slog.Any("error", err))
	}

	ticker := time.NewTicker(costEstimateRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Refresh(ctx); err != nil {
				c.logger.Warn("cost estimate refresh failed", slog.Any("error", err))
			}
		}
	}
}
