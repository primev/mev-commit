package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

// reconciliationLookbackDays is the rolling window used by the reconciliation
// metric. Sized to be much larger than the longest sweep cadence so the
// timing mismatch between "miles awarded for a swap" and "sweep settling"
// averages out within the window.
const reconciliationLookbackDays = 7

// reconciliationInterval is how often the metric is recomputed and logged.
const reconciliationInterval = 1 * time.Hour

// reconciliationAlertHigh / Low are the ratio thresholds at which we surface
// log-level alerts. Above 1.0 means we're over-paying users (estimates are
// too aggressive); well below 1.0 means we're under-paying (estimates are
// too conservative).
const (
	reconciliationAlertHigh = 1.05
	reconciliationAlertLow  = 0.75
)

// reconciliationStats is the snapshot of one metric run.
type reconciliationStats struct {
	MilesPaidETH      float64
	RealizedProfitETH float64
	Ratio             float64 // miles_paid / realized_profit
	NRowsAwarded      int
	NSweepsRealized   int
	WindowDays        int
	ComputedAt        time.Time
}

// reconciliationMonitor periodically computes the
// (miles_paid_eth / realized_sweep_profit_eth) ratio over the configured
// lookback window. Used as a tuning signal: if the ratio drifts above 1.0
// for sustained periods, the per-token estimate percentile should be raised.
type reconciliationMonitor struct {
	db           *sql.DB
	logger       *slog.Logger
	executorAddr string // lowercased hex
}

func newReconciliationMonitor(db *sql.DB, logger *slog.Logger, executorAddr string) *reconciliationMonitor {
	return &reconciliationMonitor{
		db:           db,
		logger:       logger,
		executorAddr: executorAddr,
	}
}

// Compute runs one reconciliation pass and returns the stats.
func (r *reconciliationMonitor) Compute(ctx context.Context) (reconciliationStats, error) {
	stats := reconciliationStats{
		WindowDays: reconciliationLookbackDays,
		ComputedAt: time.Now(),
	}

	// Miles awarded to real users (excludes executor) over the lookback.
	// miles_paid_eth = sum(miles) * weiPerPoint / 1e18. This is the 90%
	// user-share basis. Total economic-value basis is miles_paid_eth / 0.9.
	var totalMiles int64
	err := r.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT COALESCE(SUM(miles), 0) AS total_miles, COUNT(*) AS n
FROM mevcommit_57173.fastswap_miles
WHERE processed = 1
  AND miles > 0
  AND LOWER(user_address) != ?
  AND block_timestamp >= NOW() - INTERVAL %d DAY
`, reconciliationLookbackDays), r.executorAddr).Scan(&totalMiles, &stats.NRowsAwarded)
	if err != nil {
		return stats, fmt.Errorf("query miles paid: %w", err)
	}
	stats.MilesPaidETH = float64(totalMiles) * float64(weiPerPoint) / 1e18

	// Realized sweep profit: executor's sweep rows. ETH received from sweep
	// minus L1 gas paid. user_amt_out + surplus is the gross ETH return
	// (since output token is ETH for these rows); gas_cost is the L1 gas cost.
	err = r.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT
  COALESCE(SUM(
    (CAST(user_amt_out AS DOUBLE) + CAST(surplus AS DOUBLE) - CAST(gas_cost AS DOUBLE)) / 1e18
  ), 0) AS realized_eth,
  COUNT(*) AS n_sweeps
FROM mevcommit_57173.fastswap_miles
WHERE LOWER(user_address) = ?
  AND swap_type = 'eth_weth'
  AND block_timestamp >= NOW() - INTERVAL %d DAY
`, reconciliationLookbackDays), r.executorAddr).Scan(&stats.RealizedProfitETH, &stats.NSweepsRealized)
	if err != nil {
		return stats, fmt.Errorf("query realized profit: %w", err)
	}

	if stats.RealizedProfitETH > 0 {
		stats.Ratio = stats.MilesPaidETH / stats.RealizedProfitETH
	}
	return stats, nil
}

// perTokenBreakdown holds a per-output-token slice of the reconciliation
// metric. Useful for spotting tokens that drift individually while the
// aggregate stays balanced.
type perTokenBreakdown struct {
	Token        string
	MilesPaidETH float64
	NRows        int
}

// ComputePerToken returns the per-output-token miles-paid breakdown over
// the same lookback window. Sweep-side realized profit isn't attributed to a
// single output token in fastswap_miles (sweep tx rows have swap_type=eth_weth
// and don't carry the swept token), so this method only reports the
// obligation side. Pair with sweep_executed logs to compare against realized.
func (r *reconciliationMonitor) ComputePerToken(ctx context.Context) ([]perTokenBreakdown, error) {
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
SELECT LOWER(output_token) AS token,
       SUM(miles) AS miles_sum,
       COUNT(*) AS n
FROM mevcommit_57173.fastswap_miles
WHERE processed = 1
  AND miles > 0
  AND LOWER(user_address) != ?
  AND swap_type = 'erc20'
  AND block_timestamp >= NOW() - INTERVAL %d DAY
GROUP BY LOWER(output_token)
ORDER BY miles_sum DESC
`, reconciliationLookbackDays), r.executorAddr)
	if err != nil {
		return nil, fmt.Errorf("per-token reconciliation query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []perTokenBreakdown
	for rows.Next() {
		var token string
		var milesSum int64
		var n int
		if err := rows.Scan(&token, &milesSum, &n); err != nil {
			return out, err
		}
		out = append(out, perTokenBreakdown{
			Token:        token,
			MilesPaidETH: float64(milesSum) * float64(weiPerPoint) / 1e18,
			NRows:        n,
		})
	}
	return out, rows.Err()
}

// Run starts the periodic reconciliation loop. Returns when ctx is
// cancelled. Each tick computes and logs the metric; threshold breaches log
// at warning level.
func (r *reconciliationMonitor) Run(ctx context.Context) {
	ticker := time.NewTicker(reconciliationInterval)
	defer ticker.Stop()

	tick := func() {
		stats, err := r.Compute(ctx)
		if err != nil {
			r.logger.Warn("reconciliation compute failed", slog.Any("error", err))
			return
		}
		level := slog.LevelInfo
		switch {
		case stats.Ratio > reconciliationAlertHigh:
			level = slog.LevelWarn
		case stats.Ratio > 0 && stats.Ratio < reconciliationAlertLow:
			level = slog.LevelWarn
		}
		r.logger.Log(ctx, level, "reconciliation_metric",
			slog.Float64("ratio", stats.Ratio),
			slog.Float64("miles_paid_eth", stats.MilesPaidETH),
			slog.Float64("realized_profit_eth", stats.RealizedProfitETH),
			slog.Int("rows_awarded", stats.NRowsAwarded),
			slog.Int("sweeps_realized", stats.NSweepsRealized),
			slog.Int("window_days", stats.WindowDays))

		// Per-token breakdown: helps spot a single token drifting while
		// the aggregate looks healthy. Best-effort — a query failure here
		// shouldn't suppress the aggregate metric.
		breakdown, err := r.ComputePerToken(ctx)
		if err != nil {
			r.logger.Warn("per-token reconciliation failed", slog.Any("error", err))
			return
		}
		for _, b := range breakdown {
			r.logger.Info("reconciliation_per_token",
				slog.String("token", b.Token),
				slog.Float64("miles_paid_eth", b.MilesPaidETH),
				slog.Int("rows", b.NRows),
				slog.Int("window_days", stats.WindowDays))
		}
	}

	tick()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick()
		}
	}
}
