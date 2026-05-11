package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"sync"
	"time"
)

// gasBufferMaxAge bounds how old a sample stays in the buffer. Sized for the
// longest cadence period we use (DAI/USDT at 48h) plus a small margin.
const gasBufferMaxAge = 50 * time.Hour

// gasBufferPersistInterval controls how often the in-memory buffer is
// snapshotted to the meta table for restart resilience.
const gasBufferPersistInterval = 5 * time.Minute

// gasBufferMetaKey is the key under which the serialized buffer is stored in
// fastswap_miles_meta.
const gasBufferMetaKey = "gas_observations_v1"

// gasObservation is one L1 gas-price reading at a point in time.
type gasObservation struct {
	At        time.Time `json:"at"`
	WeiPerGas uint64    `json:"wpg"` // wei per gas; nominally fits comfortably in 64 bits at typical L1 gwei
}

// gasBuffer stores recent L1 gas observations in memory with periodic
// persistence. It is safe for concurrent use.
type gasBuffer struct {
	db     *sql.DB
	logger *slog.Logger

	mu   sync.RWMutex
	data []gasObservation
}

func newGasBuffer(db *sql.DB, logger *slog.Logger) *gasBuffer {
	return &gasBuffer{
		db:     db,
		logger: logger,
		data:   make([]gasObservation, 0, 1024),
	}
}

// Observe records a gas-price reading. Pruning of old samples happens
// opportunistically here.
func (g *gasBuffer) Observe(weiPerGas *big.Int) {
	if weiPerGas == nil || weiPerGas.Sign() < 0 {
		return
	}
	wpg := weiPerGas.Uint64()
	now := time.Now()

	g.mu.Lock()
	defer g.mu.Unlock()
	g.data = append(g.data, gasObservation{At: now, WeiPerGas: wpg})
	g.pruneLocked(now)
}

// Percentile returns the p-th percentile (0-100) of observations within the
// given lookback window, in wei per gas. Returns (0, false) if no data falls
// in the window.
func (g *gasBuffer) Percentile(p int, lookback time.Duration) (uint64, bool) {
	if p < 0 || p > 100 {
		return 0, false
	}
	cutoff := time.Now().Add(-lookback)

	g.mu.RLock()
	values := make([]uint64, 0, len(g.data))
	for _, obs := range g.data {
		if obs.At.After(cutoff) {
			values = append(values, obs.WeiPerGas)
		}
	}
	g.mu.RUnlock()

	if len(values) == 0 {
		return 0, false
	}
	slices.Sort(values)

	// Nearest-rank percentile, 1-indexed, ceiling convention:
	//   rank = ceil(p/100 × N) implemented via integer math.
	rank := min(max((p*len(values)+99)/100, 1), len(values))
	return values[rank-1], true
}

// Size returns the number of samples currently in the buffer.
func (g *gasBuffer) Size() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.data)
}

func (g *gasBuffer) pruneLocked(now time.Time) {
	cutoff := now.Add(-gasBufferMaxAge)
	keep := 0
	for _, obs := range g.data {
		if obs.At.After(cutoff) {
			break
		}
		keep++
	}
	if keep > 0 {
		g.data = append(g.data[:0], g.data[keep:]...)
	}
}

// persist writes the current buffer to the meta table. Best-effort; logs
// failures and returns nil unless the DB call itself errors.
func (g *gasBuffer) persist(ctx context.Context) error {
	g.mu.RLock()
	snapshot := make([]gasObservation, len(g.data))
	copy(snapshot, g.data)
	g.mu.RUnlock()

	blob, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("marshal gas buffer: %w", err)
	}

	_, err = g.db.ExecContext(ctx, `
INSERT INTO mevcommit_57173.fastswap_miles_meta (k, v) VALUES (?, ?)
`, gasBufferMetaKey, string(blob))
	if err != nil {
		return fmt.Errorf("persist gas buffer: %w", err)
	}
	g.logger.Debug("gas buffer persisted", slog.Int("samples", len(snapshot)))
	return nil
}

// load reads the persisted buffer on startup. Missing or invalid data leaves
// the buffer empty (will refill from live observations).
func (g *gasBuffer) load(ctx context.Context) error {
	var raw string
	err := g.db.QueryRowContext(ctx,
		"SELECT v FROM mevcommit_57173.fastswap_miles_meta WHERE k = ?",
		gasBufferMetaKey).Scan(&raw)
	if err == sql.ErrNoRows {
		g.logger.Info("no persisted gas buffer; starting empty")
		return nil
	}
	if err != nil {
		return fmt.Errorf("load gas buffer: %w", err)
	}

	var snapshot []gasObservation
	if err := json.Unmarshal([]byte(raw), &snapshot); err != nil {
		g.logger.Warn("persisted gas buffer is corrupt; starting empty",
			slog.Any("error", err))
		return nil
	}

	now := time.Now()
	cutoff := now.Add(-gasBufferMaxAge)
	g.mu.Lock()
	defer g.mu.Unlock()
	g.data = g.data[:0]
	for _, obs := range snapshot {
		if obs.At.After(cutoff) {
			g.data = append(g.data, obs)
		}
	}
	g.logger.Info("gas buffer loaded", slog.Int("samples", len(g.data)))
	return nil
}

// Run starts the periodic persistence loop. Returns when ctx is cancelled.
// Loads the persisted snapshot before entering the loop.
func (g *gasBuffer) Run(ctx context.Context) {
	if err := g.load(ctx); err != nil {
		g.logger.Warn("gas buffer load failed; continuing with empty buffer",
			slog.Any("error", err))
	}

	ticker := time.NewTicker(gasBufferPersistInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Best-effort final flush.
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = g.persist(flushCtx)
			cancel()
			return
		case <-ticker.C:
			if err := g.persist(ctx); err != nil {
				g.logger.Warn("gas buffer persist failed", slog.Any("error", err))
			}
		}
	}
}
