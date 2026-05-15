package main

import (
	"log/slog"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// helper: a buffer pre-loaded with N samples at uniform values so percentile
// math is predictable.
func bufWithUniform(n int, weiPerGas uint64) *gasBuffer {
	g := &gasBuffer{
		logger: slog.Default(),
		data:   make([]gasObservation, 0, n),
	}
	now := time.Now()
	for i := range n {
		g.data = append(g.data, gasObservation{At: now.Add(-time.Duration(i) * time.Second), WeiPerGas: weiPerGas})
	}
	return g
}

func TestDecideSweep_SkipBeforeCadence(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr) // 24h cadence
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-1 * time.Hour), // recent sweep
		CurrentGasWei: 1_000_000,
		Buf:           bufWithUniform(100, 1_000_000_000),
	})
	if out.Decision != SweepSkip {
		t.Errorf("Decision = %v, want SweepSkip", out.Decision)
	}
	if out.Reason != "waiting_for_cadence" {
		t.Errorf("Reason = %q, want waiting_for_cadence", out.Reason)
	}
}

func TestDecideSweep_AttemptAfterCadence_GasOk(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr) // 24h cadence
	// Buffer sample = 5_000_000_000 (5 gwei), current = 1_000_000_000 (1 gwei) → below p25.
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-25 * time.Hour),
		CurrentGasWei: 1_000_000_000,
		Buf:           bufWithUniform(100, 5_000_000_000),
	})
	if out.Decision != SweepAttempt {
		t.Errorf("Decision = %v, want SweepAttempt", out.Decision)
	}
	if !out.HasGasData {
		t.Errorf("HasGasData = false, want true")
	}
}

func TestDecideSweep_SkipAboveGasCap(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr) // 24h cadence, p25 cap
	// All samples at 5 gwei → p25 = 5 gwei. Current at 10 gwei should skip.
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-25 * time.Hour),
		CurrentGasWei: 10_000_000_000,
		Buf:           bufWithUniform(100, 5_000_000_000),
	})
	if out.Decision != SweepSkip {
		t.Errorf("Decision = %v, want SweepSkip", out.Decision)
	}
	if out.Reason != "gas_above_cap" {
		t.Errorf("Reason = %q, want gas_above_cap", out.Reason)
	}
}

func TestDecideSweep_ForceAfterCadenceX1_5(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr) // 24h cadence → 36h force
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-37 * time.Hour),
		CurrentGasWei: 999_999_999_999, // very high gas — would normally skip
		Buf:           bufWithUniform(100, 1_000_000_000),
	})
	if out.Decision != SweepForce {
		t.Errorf("Decision = %v, want SweepForce", out.Decision)
	}
}

func TestDecideSweep_VolatileNoCadenceFloor(t *testing.T) {
	cfg := lookupTokenConfig(pepeAddr) // 0 cadence
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-1 * time.Minute), // very recent
		CurrentGasWei: 1,
		Buf:           bufWithUniform(100, 1_000_000_000),
	})
	if out.Decision != SweepAttempt {
		t.Errorf("Decision = %v, want SweepAttempt (volatile has no cadence floor)", out.Decision)
	}
}

func TestDecideSweep_VolatileForceAt6h(t *testing.T) {
	cfg := lookupTokenConfig(pepeAddr) // 6h force
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-7 * time.Hour),
		CurrentGasWei: 999_999_999_999,
		Buf:           bufWithUniform(100, 1_000_000_000),
	})
	if out.Decision != SweepForce {
		t.Errorf("Decision = %v, want SweepForce", out.Decision)
	}
}

func TestDecideSweep_NoGasData_AllowsAttempt(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr)
	emptyBuf := &gasBuffer{logger: slog.Default(), data: nil}
	out := decideSweep(sweepDecisionInput{
		Now:           time.Now(),
		Cfg:           cfg,
		LastSweepAt:   time.Now().Add(-25 * time.Hour),
		CurrentGasWei: 100,
		Buf:           emptyBuf,
	})
	if out.Decision != SweepAttempt {
		t.Errorf("Decision = %v, want SweepAttempt", out.Decision)
	}
	if out.Reason != "no_gas_data" {
		t.Errorf("Reason = %q, want no_gas_data", out.Reason)
	}
	if out.HasGasData {
		t.Errorf("HasGasData = true, want false")
	}
}

func TestSweepClock_NeverSweptIsZero(t *testing.T) {
	c := newSweepClock()
	if got := c.LastSwept(usdcAddr); !got.IsZero() {
		t.Errorf("LastSwept on never-swept = %v, want zero time", got)
	}
}

func TestSweepClock_MarkAndGet(t *testing.T) {
	c := newSweepClock()
	now := time.Now()
	c.MarkSwept(usdcAddr, now)
	if got := c.LastSwept(usdcAddr); !got.Equal(now) {
		t.Errorf("LastSwept = %v, want %v", got, now)
	}
}

func TestSweepClock_SnapshotAndRestore(t *testing.T) {
	c := newSweepClock()
	now := time.Now()
	c.MarkSwept(usdcAddr, now)
	c.MarkSwept(usdtAddr, now.Add(-1*time.Hour))

	snap := c.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("snapshot size = %d, want 2", len(snap))
	}

	c2 := newSweepClock()
	c2.Restore(snap)
	if got := c2.LastSwept(usdcAddr); !got.Equal(now) {
		t.Errorf("restored USDC = %v, want %v", got, now)
	}
}

func TestSweepClock_RestoreReplacesState(t *testing.T) {
	c := newSweepClock()
	c.MarkSwept(usdcAddr, time.Now())
	// Restore an unrelated snapshot — should clear USDC.
	c.Restore(map[common.Address]time.Time{
		usdtAddr: time.Now().Add(-1 * time.Hour),
	})
	if got := c.LastSwept(usdcAddr); !got.IsZero() {
		t.Errorf("USDC after restore = %v, want zero (not in snapshot)", got)
	}
	// Quiet the unused-import warning if big.Int isn't used elsewhere.
	_ = big.NewInt(0)
}

func TestSweepDecisionString(t *testing.T) {
	cases := []struct {
		d    sweepDecision
		want string
	}{
		{SweepSkip, "skip"},
		{SweepAttempt, "attempt"},
		{SweepForce, "force"},
		{sweepDecision(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.d.String(); got != c.want {
			t.Errorf("decision %d = %q, want %q", c.d, got, c.want)
		}
	}
}
