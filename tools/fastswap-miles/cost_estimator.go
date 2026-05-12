package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
)

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
	// PerRowOverhead is the estimated sweep overhead per user row, in ETH.
	// This is the value subtracted in the miles formula in lieu of realized
	// pro-rata sweep gas.
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

	mu        sync.RWMutex
	estimates map[string]costEstimate // key: lowercased token hex
	lastFresh time.Time
}

func newCostEstimator(db *sql.DB, logger *slog.Logger) *costEstimator {
	return &costEstimator{
		db:        db,
		logger:    logger,
		estimates: make(map[string]costEstimate),
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

	c.mu.Lock()
	c.estimates = fresh
	c.lastFresh = time.Now()
	c.mu.Unlock()

	c.logger.Info("cost estimates refreshed",
		slog.Int("tokens", len(fresh)),
		slog.Duration("window", costEstimateLookbackDays*24*time.Hour))
	return nil
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
