package main

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func newTestEstimator() *costEstimator {
	return &costEstimator{
		logger:    slog.Default(),
		estimates: make(map[string]costEstimate),
	}
}

func TestCostEstimator_Get_NoData_FallsBackToLastResort(t *testing.T) {
	// Exercise the real constructor so it's not unused.
	c := newCostEstimator(nil, slog.Default(), common.Address{}, common.Address{})
	got := c.Get("0xdeadbeef")
	if got.Source != "default_no_data" {
		t.Errorf("source = %q, want default_no_data", got.Source)
	}
	if got.PerRowOverhead != costEstimateLastResort {
		t.Errorf("overhead = %v, want %v (last resort)", got.PerRowOverhead, costEstimateLastResort)
	}
	if got.SweepCount != 0 {
		t.Errorf("sweep count = %d, want 0", got.SweepCount)
	}
}

func TestCostEstimator_Get_HighData_UsesP25(t *testing.T) {
	c := newTestEstimator()
	c.estimates["0xtoken"] = costEstimate{
		PerRowOverhead: 5e-5,
		Source:         "p25",
		SweepCount:     50,
		ComputedAt:     time.Now(),
	}
	got := c.Get("0xToken") // case-insensitive lookup
	if got.Source != "p25" {
		t.Errorf("source = %q, want p25", got.Source)
	}
	if got.PerRowOverhead != 5e-5 {
		t.Errorf("overhead = %v, want 5e-5", got.PerRowOverhead)
	}
	if got.SweepCount != 50 {
		t.Errorf("sweep count = %d, want 50", got.SweepCount)
	}
}

func TestCostEstimator_Get_LowData_UsesP75(t *testing.T) {
	c := newTestEstimator()
	c.estimates["0xtoken"] = costEstimate{
		PerRowOverhead: 1.5e-4,
		Source:         "p75_low_data",
		SweepCount:     5,
		ComputedAt:     time.Now(),
	}
	got := c.Get("0xtoken")
	if got.Source != "p75_low_data" {
		t.Errorf("source = %q, want p75_low_data", got.Source)
	}
}

func TestCostEstimator_Get_CaseInsensitive(t *testing.T) {
	c := newTestEstimator()
	c.estimates["0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"] = costEstimate{
		PerRowOverhead: 4e-5,
		Source:         "p25",
		SweepCount:     100,
	}
	upper := c.Get(strings.ToUpper("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"))
	if upper.Source != "p25" {
		t.Errorf("upper-case lookup failed: source = %q", upper.Source)
	}
}

func TestCostEstimateMinSweeps_BoundaryBehavior(t *testing.T) {
	// Sanity-check the constant matches what the doc says (n >= 10 uses p25).
	if costEstimateMinSweeps != 10 {
		t.Errorf("costEstimateMinSweeps = %d, doc says 10", costEstimateMinSweeps)
	}
}

func TestCostEstimateLastResort_Reasonable(t *testing.T) {
	// Sanity: the last-resort value should be conservative but not absurd.
	// Currently 0.001 ETH (~$3 at $3000/ETH) — high enough that miles for a
	// novel-token swap will be modest, low enough that it's not impossible.
	if costEstimateLastResort < 1e-4 || costEstimateLastResort > 1e-2 {
		t.Errorf("costEstimateLastResort = %v, expected within [1e-4, 1e-2] sanity range", costEstimateLastResort)
	}
}

func TestCostEstimator_Constructor_LowercasesAddresses(t *testing.T) {
	// Sweep-bid lookup queries compare against LOWER(user_address) and
	// LOWER(output_token), so the addresses stored on the estimator MUST be
	// lowercased. A mixed-case stored value would silently match zero rows
	// and the sweep bid term would never fire.
	exec := common.HexToAddress("0x959DAD78D5B68986a43cD270134A2704a990aa68")
	weth := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	c := newCostEstimator(nil, slog.Default(), exec, weth)
	if c.executorAddr != strings.ToLower(exec.Hex()) {
		t.Errorf("executorAddr = %q, want %q", c.executorAddr, strings.ToLower(exec.Hex()))
	}
	if c.wethAddr != strings.ToLower(weth.Hex()) {
		t.Errorf("wethAddr = %q, want %q", c.wethAddr, strings.ToLower(weth.Hex()))
	}
}
