package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// sweepClockPersistInterval bounds how often the in-memory clock is
// snapshotted to the meta table for restart resilience.
const sweepClockPersistInterval = 5 * time.Minute

// sweepClockMetaKey is the row in fastswap_miles_meta that stores the
// JSON-serialized per-token last-swept timestamps.
const sweepClockMetaKey = "sweep_clock_v1"

// sweepDecision is the scheduler's verdict for a single token at a single
// evaluation tick. Profitability is intentionally NOT a scheduler concern —
// the scheduler decides "are we even allowed to attempt"; the caller then
// runs the (expensive) Barter quote and applies the profitability guard.
type sweepDecision int

const (
	// SweepSkip — do not attempt this cycle. Either cadence hasn't elapsed,
	// or gas is above the tier cap.
	SweepSkip sweepDecision = iota
	// SweepAttempt — eligible to sweep if profitability passes. Gas cap
	// passed (or there's no gas data yet).
	SweepAttempt
	// SweepForce — force-sweep window reached. Gas cap is waived; the caller
	// should still apply the profitability guard (we never sweep at a loss).
	SweepForce
)

func (d sweepDecision) String() string {
	switch d {
	case SweepSkip:
		return "skip"
	case SweepAttempt:
		return "attempt"
	case SweepForce:
		return "force"
	default:
		return "unknown"
	}
}

// sweepDecisionInput packages the runtime state needed to make a decision.
type sweepDecisionInput struct {
	Now           time.Time
	Cfg           tokenConfig
	LastSweepAt   time.Time // zero value means "never swept"
	CurrentGasWei uint64    // wei per gas
	Buf           *gasBuffer
}

// sweepDecisionOutput is the scheduler's verdict plus enough context for
// downstream logging and reconciliation.
type sweepDecisionOutput struct {
	Decision   sweepDecision
	Reason     string // short tag for logs
	GasCapWei  uint64 // wei per gas; the cap that was applied (0 if no cap)
	HasGasData bool   // false when the gas buffer had no samples in the window
}

// decideSweep applies cadence, gas-cap, and force-sweep logic. Profitability
// is checked separately by the caller after a successful Barter quote.
func decideSweep(in sweepDecisionInput) sweepDecisionOutput {
	elapsed := in.Now.Sub(in.LastSweepAt)

	// Cadence floor (only applies when SweepCadence > 0; volatile has zero).
	if in.Cfg.SweepCadence > 0 && elapsed < in.Cfg.SweepCadence {
		return sweepDecisionOutput{
			Decision: SweepSkip,
			Reason:   "waiting_for_cadence",
		}
	}

	// Force-sweep window: bypass the gas cap, but profitability still applies
	// downstream so we never sweep at a loss.
	if elapsed >= in.Cfg.forceSweepInterval() {
		return sweepDecisionOutput{
			Decision: SweepForce,
			Reason:   "force_sweep_window",
		}
	}

	// Gas cap check using percentile of recent observations.
	gasCap, ok := in.Buf.Percentile(in.Cfg.Tier.GasCapPercentile(), in.Cfg.gasCapLookback())
	if !ok {
		// No gas data yet (cold start). Allow the attempt; profitability
		// guard is still the absolute floor.
		return sweepDecisionOutput{
			Decision:   SweepAttempt,
			Reason:     "no_gas_data",
			HasGasData: false,
		}
	}

	if in.CurrentGasWei > gasCap {
		return sweepDecisionOutput{
			Decision:   SweepSkip,
			Reason:     "gas_above_cap",
			GasCapWei:  gasCap,
			HasGasData: true,
		}
	}

	return sweepDecisionOutput{
		Decision:   SweepAttempt,
		Reason:     "gas_within_cap",
		GasCapWei:  gasCap,
		HasGasData: true,
	}
}

// sweepClock tracks per-token last-successful-sweep timestamps. Used by the
// scheduler to compute the elapsed-since-last-sweep input.
//
// Persisted across restarts so cadence enforcement isn't reset by pod
// recycles. (Persistence wired through main.go alongside the gas buffer.)
type sweepClock struct {
	mu        sync.RWMutex
	lastSweep map[common.Address]time.Time // key: token address
}

func newSweepClock() *sweepClock {
	return &sweepClock{lastSweep: make(map[common.Address]time.Time)}
}

// MarkSwept records a successful sweep for the token. Call only after the
// sweep transaction has been submitted (or simulated in dry-run).
func (c *sweepClock) MarkSwept(token common.Address, at time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSweep[token] = at
}

// LastSwept returns the time of the most recent successful sweep for the
// token. Returns the zero time if the token has never been swept (which the
// scheduler treats as "elapsed = a very long time" — i.e., it'll force-sweep
// on the first opportunity, subject to profitability).
func (c *sweepClock) LastSwept(token common.Address) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastSweep[token]
}

// Snapshot returns a copy of the per-token timestamps. Used for persistence.
func (c *sweepClock) Snapshot() map[common.Address]time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[common.Address]time.Time, len(c.lastSweep))
	maps.Copy(out, c.lastSweep)
	return out
}

// Restore replaces the in-memory state with the provided snapshot. Used
// during startup-from-persistence.
func (c *sweepClock) Restore(snapshot map[common.Address]time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSweep = make(map[common.Address]time.Time, len(snapshot))
	maps.Copy(c.lastSweep, snapshot)
}

// persist writes the current snapshot to the meta table. Best-effort.
func (c *sweepClock) persist(ctx context.Context, db *sql.DB) error {
	// Map keys serialize to checksum-cased hex addresses; reload normalizes
	// back via common.HexToAddress.
	snapshot := c.Snapshot()
	encoded := make(map[string]time.Time, len(snapshot))
	for k, v := range snapshot {
		encoded[k.Hex()] = v
	}
	blob, err := json.Marshal(encoded)
	if err != nil {
		return fmt.Errorf("marshal sweep clock: %w", err)
	}
	_, err = db.ExecContext(ctx, `
INSERT INTO mevcommit_57173.fastswap_miles_meta (k, v) VALUES (?, ?)
`, sweepClockMetaKey, string(blob))
	if err != nil {
		return fmt.Errorf("persist sweep clock: %w", err)
	}
	return nil
}

// load reads the persisted snapshot on startup. Missing or invalid data
// leaves the clock empty (every token's lastSweep is zero, scheduler
// treats that as "elapsed = a very long time" and is willing to sweep on
// the first opportunity subject to profitability).
func (c *sweepClock) load(ctx context.Context, db *sql.DB, logger *slog.Logger) error {
	var raw string
	err := db.QueryRowContext(ctx,
		"SELECT v FROM mevcommit_57173.fastswap_miles_meta WHERE k = ?",
		sweepClockMetaKey).Scan(&raw)
	if err == sql.ErrNoRows {
		logger.Info("no persisted sweep clock; starting empty")
		return nil
	}
	if err != nil {
		return fmt.Errorf("load sweep clock: %w", err)
	}
	var encoded map[string]time.Time
	if err := json.Unmarshal([]byte(raw), &encoded); err != nil {
		logger.Warn("persisted sweep clock is corrupt; starting empty",
			slog.Any("error", err))
		return nil
	}
	snapshot := make(map[common.Address]time.Time, len(encoded))
	for k, v := range encoded {
		snapshot[common.HexToAddress(k)] = v
	}
	c.Restore(snapshot)
	logger.Info("sweep clock loaded", slog.Int("tokens", len(snapshot)))
	return nil
}

// Run starts the periodic persistence loop. Returns when ctx is cancelled.
// Loads the persisted snapshot before entering the loop.
func (c *sweepClock) Run(ctx context.Context, db *sql.DB, logger *slog.Logger) {
	if err := c.load(ctx, db, logger); err != nil {
		logger.Warn("sweep clock load failed; continuing with empty clock",
			slog.Any("error", err))
	}

	ticker := time.NewTicker(sweepClockPersistInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_ = c.persist(flushCtx, db)
			cancel()
			return
		case <-ticker.C:
			if err := c.persist(ctx, db); err != nil {
				logger.Warn("sweep clock persist failed", slog.Any("error", err))
			}
		}
	}
}
