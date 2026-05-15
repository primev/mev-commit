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

// sweepBidGlobalPercentile is the percentile of recent fastswap bid costs
// used as a proxy for a sweep tx's bid. The sweep tx is just another
// fastswap, so the per-user-tx bid distribution is the right reference
// population. p75 (under-promise) — same direction the cron uses for the
// dashboard's bid-cost estimate.
const sweepBidGlobalPercentile = 0.75

// sweepBidFallbackEth is the per-sweep bid used when the global percentile
// query returns nothing (cold start, freshly-deployed pod with empty
// fastswap_miles). Mirrors the cron's FALLBACK_BID_COST_ETH in
// fastprotocolapp/.../miles-estimate-gas/route.ts — post-fix p75 of
// realized bid distribution.
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
// over the configured lookback window. This is the only method that touches
// the database; intended to be called periodically by a background goroutine.
//
// PerRowOverhead is the sum of two terms:
//
//  1. Per-row sweep gas overhead — pro-rata share of the L1 gas the executor
//     paid to convert this row's surplus tokens to ETH. Computed as the
//     per-token p25 (or p75 on low data) of `surplus_eth − net_profit_eth −
//     bid_cost/1e18` across processed ETH-input ERC20-output rows.
//
//  2. Per-row sweep bid contribution — the preconf bid the executor pays to
//     submit the sweep tx itself, amortized across the user rows in the
//     matching batch. The sweep tx IS a fastswap, so we proxy its bid with
//     the p75 of recent realized bid_cost across all processed fastswap
//     rows (tight distribution, far more samples than executor sweeps
//     alone), then multiply by per-token sweep count and divide by
//     per-token user-row count to get the per-row contribution.
//
// Both terms scale together with batch size — low-volume tokens have a small
// number of user rows per sweep so the per-row contribution is high, and
// high-volume tokens dilute both. The miles formula in
// `awardUpfrontERC20Miles` subtracts the combined value, so the sweep tx's
// own bid no longer falls on the protocol's books silently.
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

// computePerTokenSweepBidEth returns the per-row sweep bid contribution (in
// ETH) for each output token, over the same lookback window used for
// PerRowOverhead.
//
// For each output_token T:
//
//	per_row_sweep_bid_eth(T) = (n_sweeps(T) × global_bid_p75_eth) / n_user_rows(T)
//
// Why a global percentile and not the realized per-sweep bid: the sweep tx
// IS a fastswap, so the distribution of all user-row bids is the right
// reference population. It's tight (memory: stddev/mean ≈ 24%) and has
// orders of magnitude more samples than the executor-sweep subset alone,
// so a percentile across the broader set is at least as accurate as
// looking up each sweep's realized bid via tx_view — and avoids the JOIN
// entirely. p75 mirrors the cron's same-purpose proxy on the frontend.
//
// "Executor sweep" = a row in fastswap_miles whose user_address is the
// executor and whose output_token is WETH; its input_token is the ERC20
// being swept (which equals the user's output_token).
func (c *costEstimator) computePerTokenSweepBidEth(ctx context.Context) (map[string]float64, error) {
	bidEth, err := c.queryGlobalBidPercentileEth(ctx)
	if err != nil {
		c.logger.Warn("global bid percentile query failed; using fallback",
			slog.Any("error", err), slog.Float64("fallback_eth", sweepBidFallbackEth))
		bidEth = sweepBidFallbackEth
	}

	sweepCountByToken, err := c.queryExecutorSweepCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("query executor sweep counts: %w", err)
	}
	if len(sweepCountByToken) == 0 {
		return map[string]float64{}, nil
	}

	userRowsByToken, err := c.queryUserRowCounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("query user row counts: %w", err)
	}

	perRowByToken := make(map[string]float64, len(sweepCountByToken))
	for token, nSweeps := range sweepCountByToken {
		nRows := userRowsByToken[token]
		if nSweeps <= 0 || nRows <= 0 {
			continue
		}
		perRowByToken[token] = float64(nSweeps) * bidEth / float64(nRows)
	}
	return perRowByToken, nil
}

// queryGlobalBidPercentileEth returns the chosen percentile of realized bid
// costs (in ETH) across processed fastswap rows in the lookback window.
// Returns 0 when no rows are available so the caller can substitute a
// fallback.
func (c *costEstimator) queryGlobalBidPercentileEth(ctx context.Context) (float64, error) {
	var p float64
	row := c.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT percentile_approx(CAST(bid_cost AS DOUBLE)/1e18, %f) AS p
FROM mevcommit_57173.fastswap_miles
WHERE processed = 1
  AND bid_cost IS NOT NULL
  AND CAST(bid_cost AS DOUBLE) > 0
  AND block_timestamp >= NOW() - INTERVAL %d DAY
`, sweepBidGlobalPercentile, costEstimateLookbackDays))
	if err := row.Scan(&p); err != nil {
		return 0, err
	}
	if !(p > 0) {
		return 0, fmt.Errorf("percentile returned non-positive value: %v", p)
	}
	return p, nil
}

// queryExecutorSweepCounts returns the number of sweeps the executor has
// performed per swept-token (= user's output_token) over the lookback window.
func (c *costEstimator) queryExecutorSweepCounts(ctx context.Context) (map[string]int, error) {
	rows, err := c.db.QueryContext(ctx, fmt.Sprintf(`
SELECT LOWER(input_token) AS swept_token, COUNT(*) AS n_sweeps
FROM mevcommit_57173.fastswap_miles
WHERE LOWER(user_address) = ?
  AND swap_type = 'eth_weth'
  AND LOWER(output_token) = ?
  AND block_timestamp >= NOW() - INTERVAL %d DAY
GROUP BY input_token
`, costEstimateLookbackDays), c.executorAddr, c.wethAddr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make(map[string]int)
	for rows.Next() {
		var token string
		var n int
		if err := rows.Scan(&token, &n); err != nil {
			c.logger.Warn("scan executor sweep count failed", slog.Any("error", err))
			continue
		}
		out[token] = n
	}
	return out, rows.Err()
}

// queryUserRowCounts returns the number of non-executor user rows per
// output_token over the lookback window — the divisor in the per-row sweep
// bid calculation.
func (c *costEstimator) queryUserRowCounts(ctx context.Context) (map[string]int, error) {
	rows, err := c.db.QueryContext(ctx, fmt.Sprintf(`
SELECT LOWER(output_token) AS output_token, COUNT(*) AS n_rows
FROM mevcommit_57173.fastswap_miles
WHERE swap_type = 'erc20'
  AND LOWER(user_address) != ?
  AND block_timestamp >= NOW() - INTERVAL %d DAY
GROUP BY output_token
`, costEstimateLookbackDays), c.executorAddr)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make(map[string]int)
	for rows.Next() {
		var token string
		var n int
		if err := rows.Scan(&token, &n); err != nil {
			c.logger.Warn("scan user row count failed", slog.Any("error", err))
			continue
		}
		out[token] = n
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
