package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Persistence touches the DB, so a full round-trip test requires an in-memory
// SQL stub which the rest of this package doesn't yet have. The
// JSON-encoded representation IS the persisted blob, though, so testing the
// encoding/decoding step covers the failure mode that actually matters in
// practice (data shape changes silently breaking restore).

func TestSweepClock_PersistEncodingRoundtrip(t *testing.T) {
	original := newSweepClock()
	t1 := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 5, 9, 14, 30, 0, 0, time.UTC)
	original.MarkSwept(usdcAddr, t1)
	original.MarkSwept(wbtcAddr, t2)

	// Encode the same way persist() does: snapshot → checksum-hex keys →
	// JSON blob.
	snapshot := original.Snapshot()
	encoded := make(map[string]time.Time, len(snapshot))
	for k, v := range snapshot {
		encoded[k.Hex()] = v
	}
	blob, err := json.Marshal(encoded)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Decode the same way load() does: JSON → checksum-hex keys → addresses.
	var decoded map[string]time.Time
	if err := json.Unmarshal(blob, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rehydrated := make(map[common.Address]time.Time, len(decoded))
	for k, v := range decoded {
		rehydrated[common.HexToAddress(k)] = v
	}

	restored := newSweepClock()
	restored.Restore(rehydrated)

	if !restored.LastSwept(usdcAddr).Equal(t1) {
		t.Errorf("USDC restored = %v, want %v", restored.LastSwept(usdcAddr), t1)
	}
	if !restored.LastSwept(wbtcAddr).Equal(t2) {
		t.Errorf("WBTC restored = %v, want %v", restored.LastSwept(wbtcAddr), t2)
	}
	// A token never marked must come back as zero.
	if !restored.LastSwept(daiAddr).IsZero() {
		t.Errorf("DAI should be zero, got %v", restored.LastSwept(daiAddr))
	}
}

func TestSweepClock_PersistEncodingHandlesEmptyClock(t *testing.T) {
	c := newSweepClock()
	encoded := make(map[string]time.Time)
	for k, v := range c.Snapshot() {
		encoded[k.Hex()] = v
	}
	blob, err := json.Marshal(encoded)
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}
	// Empty map serializes as `{}`; ensure we don't accidentally produce
	// `null` (which would unmarshal back as nil and look like missing).
	if !strings.HasPrefix(string(blob), "{") {
		t.Errorf("empty clock JSON = %q, want object literal", string(blob))
	}
}

func TestSweepLoopConstantsSane(t *testing.T) {
	// 1.05× sweep gas — leaves margin for between-quote gas drift.
	if sweepProfitabilityNumerator != 105 {
		t.Errorf("sweepProfitabilityNumerator = %d, want 105 (5%% margin)", sweepProfitabilityNumerator)
	}
	// Dust threshold should be tiny — most tokens have 6+ decimals so 1000
	// raw units is essentially nothing. If raised significantly, real
	// sweeps could be erroneously skipped.
	if sweepDustRawUnits != 1000 {
		t.Errorf("sweepDustRawUnits = %d, want 1000", sweepDustRawUnits)
	}
	// Barter min-return fraction matches the deferred path's value (0.98)
	// for consistency.
	if sweepBarterMinReturnFraction != 0.98 {
		t.Errorf("sweepBarterMinReturnFraction = %f, want 0.98", sweepBarterMinReturnFraction)
	}
	// Override margin: 20% buffer above raw breakeven before bypassing
	// cadence. Below 1.0 would allow overriding into losing trades.
	if cadenceOverrideMargin < 1.0 {
		t.Errorf("cadenceOverrideMargin = %v; must be >= 1.0 to avoid bypassing cadence into losses", cadenceOverrideMargin)
	}
}

func TestCadenceOverrideMet_AboveThreshold(t *testing.T) {
	// 100 rows × 0.00005 ETH = 0.005 ETH budget.
	// Sweep gas = 0.003 ETH × 1.2 = 0.0036 threshold.
	// 0.005 >= 0.0036 → override fires.
	if !cadenceOverrideMet(100, 0.00005, 0.003) {
		t.Errorf("100 rows at $0.00005 per-row vs $0.003 sweep gas should override")
	}
}

func TestCadenceOverrideMet_BelowThreshold(t *testing.T) {
	// 10 rows × 0.00005 = 0.0005 budget.
	// 0.003 × 1.2 = 0.0036 threshold. Far below.
	if cadenceOverrideMet(10, 0.00005, 0.003) {
		t.Errorf("10 rows should not be enough to override at $0.003 sweep gas")
	}
}

func TestCadenceOverrideMet_ExactlyAtThreshold(t *testing.T) {
	// At exactly the threshold (budget == 1.2 × sweep_gas), allow override.
	// 1000 rows × 0.0001 = 0.1 budget. Sweep gas = 0.1 / 1.2 ≈ 0.0833.
	if !cadenceOverrideMet(1000, 0.0001, 0.1/1.2) {
		t.Errorf("budget exactly at threshold should still trigger override")
	}
}

func TestCadenceOverrideMet_NoVolumeNoOverride(t *testing.T) {
	if cadenceOverrideMet(0, 0.0001, 0.001) {
		t.Errorf("0 rows must never override")
	}
}

func TestCadenceOverrideMet_NoCostEstimateNoOverride(t *testing.T) {
	// Zero or negative per-row overhead means we don't have an estimate;
	// override should never fire because we have no budget to compare to.
	if cadenceOverrideMet(1000, 0, 0.001) {
		t.Errorf("zero per-row estimate must never override")
	}
	if cadenceOverrideMet(1000, -1, 0.001) {
		t.Errorf("negative per-row estimate must never override")
	}
}

func TestIsBarterAmountTooLow(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{errStrFunc("barter API error 400: Amount too low for token"), true},
		{errStrFunc("amount too low"), true},
		{errStrFunc("connection refused"), false},
		{errStrFunc("barter API error 500: internal server error"), false},
	}
	for _, c := range cases {
		if got := isBarterAmountTooLow(c.err); got != c.want {
			t.Errorf("isBarterAmountTooLow(%v) = %v, want %v", c.err, got, c.want)
		}
	}
}

// errStrFunc returns a one-shot error with the given message. Local helper to
// avoid importing the errors package just for this test file.
type strErr string

func (e strErr) Error() string { return string(e) }

func errStrFunc(s string) error { return strErr(s) }
