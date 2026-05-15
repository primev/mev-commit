package main

import (
	"log/slog"
	"testing"
	"time"
)

func TestReconciliation_Constructor(t *testing.T) {
	r := newReconciliationMonitor(nil, slog.Default(), "0xabc")
	if r.executorAddr != "0xabc" {
		t.Errorf("executorAddr = %q, want 0xabc", r.executorAddr)
	}
}

func TestReconciliation_ConstantsSane(t *testing.T) {
	if reconciliationAlertHigh <= 1.0 {
		t.Errorf("alert high = %v, should be above 1.0 (over-paying threshold)", reconciliationAlertHigh)
	}
	if reconciliationAlertLow >= 1.0 {
		t.Errorf("alert low = %v, should be below 1.0 (under-paying threshold)", reconciliationAlertLow)
	}
	if reconciliationLookbackDays < 1 {
		t.Errorf("lookback = %d, expected at least 1 day", reconciliationLookbackDays)
	}
	if reconciliationInterval < time.Minute {
		t.Errorf("interval = %v, expected at least 1 minute (avoid hammering DB)", reconciliationInterval)
	}
}

func TestReconciliationStats_RatioComputedCorrectly(t *testing.T) {
	// Sanity check the math relationship without hitting a DB.
	// 100 miles awarded × 1e-5 ETH/mile = 1e-3 ETH user share.
	// Realized profit 2e-3 ETH → ratio = 0.5 (under-paying).
	totalMiles := int64(100)
	milesEth := float64(totalMiles) * float64(weiPerPoint) / 1e18
	if milesEth != 1e-3 {
		t.Errorf("miles ETH conversion = %v, want 0.001", milesEth)
	}
	realized := 2e-3
	ratio := milesEth / realized
	if ratio != 0.5 {
		t.Errorf("ratio = %v, want 0.5", ratio)
	}
	if ratio >= reconciliationAlertLow {
		t.Errorf("ratio 0.5 should be below alert low %v", reconciliationAlertLow)
	}
}
