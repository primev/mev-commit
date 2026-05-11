package main

import (
	"log/slog"
	"math/big"
	"testing"
	"time"
)

func newTestGasBuffer() *gasBuffer {
	return &gasBuffer{
		logger: slog.Default(),
		data:   make([]gasObservation, 0, 32),
	}
}

func TestGasBuffer_Observe_AppendsSample(t *testing.T) {
	// Exercise the real constructor so it isn't reported as unused.
	g := newGasBuffer(nil, slog.Default())
	g.Observe(big.NewInt(1_000_000_000)) // 1 gwei
	if g.Size() != 1 {
		t.Errorf("Size = %d, want 1", g.Size())
	}
}

func TestGasBuffer_Observe_RejectsNegativeAndNil(t *testing.T) {
	g := newTestGasBuffer()
	g.Observe(nil)
	g.Observe(big.NewInt(-1))
	if g.Size() != 0 {
		t.Errorf("Size = %d, want 0 (negative/nil rejected)", g.Size())
	}
}

func TestGasBuffer_Percentile_EmptyReturnsFalse(t *testing.T) {
	g := newTestGasBuffer()
	_, ok := g.Percentile(50, 1*time.Hour)
	if ok {
		t.Errorf("expected ok=false for empty buffer")
	}
}

func TestGasBuffer_Percentile_SingleSample(t *testing.T) {
	g := newTestGasBuffer()
	g.Observe(big.NewInt(5_000_000_000))
	val, ok := g.Percentile(50, 1*time.Hour)
	if !ok {
		t.Fatalf("expected ok=true with single sample")
	}
	if val != 5_000_000_000 {
		t.Errorf("p50 of single sample = %d, want 5_000_000_000", val)
	}
}

func TestGasBuffer_Percentile_NearestRank(t *testing.T) {
	g := newTestGasBuffer()
	for _, v := range []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		g.Observe(big.NewInt(v))
	}
	cases := []struct {
		p    int
		want uint64
	}{
		{25, 3}, // rank = (25*10)/100 = 2 → values[2-1+1]? nearest-rank: rank=ceil(0.25*10)=3 → values[2]=3
		{50, 5}, // rank = 5 → values[4] = 5
		{75, 8}, // rank = 7.5 → 8 → values[7] = 8
		{90, 9}, // rank = 9 → values[8] = 9
		{100, 10},
	}
	for _, c := range cases {
		got, ok := g.Percentile(c.p, 1*time.Hour)
		if !ok {
			t.Errorf("p%d: ok=false unexpectedly", c.p)
			continue
		}
		if got != c.want {
			t.Errorf("p%d = %d, want %d", c.p, got, c.want)
		}
	}
}

func TestGasBuffer_Percentile_LookbackWindow(t *testing.T) {
	g := newTestGasBuffer()

	// Inject "old" samples directly (before lookback) and "recent" samples.
	now := time.Now()
	g.data = append(g.data,
		gasObservation{At: now.Add(-2 * time.Hour), WeiPerGas: 100},
		gasObservation{At: now.Add(-2 * time.Hour), WeiPerGas: 200},
	)
	g.Observe(big.NewInt(50)) // recent

	val, ok := g.Percentile(50, 30*time.Minute)
	if !ok {
		t.Fatalf("expected ok=true within lookback window")
	}
	if val != 50 {
		t.Errorf("p50 within 30m lookback = %d, want 50 (older samples should be excluded)", val)
	}
}

func TestGasBuffer_Percentile_OutOfRangeReturnsFalse(t *testing.T) {
	g := newTestGasBuffer()
	g.Observe(big.NewInt(1))

	if _, ok := g.Percentile(-1, time.Hour); ok {
		t.Errorf("expected ok=false for p=-1")
	}
	if _, ok := g.Percentile(101, time.Hour); ok {
		t.Errorf("expected ok=false for p=101")
	}
}

func TestGasBuffer_Prune_OldSamplesDropped(t *testing.T) {
	g := newTestGasBuffer()

	// Inject samples older than gasBufferMaxAge directly.
	old := time.Now().Add(-gasBufferMaxAge - time.Hour)
	g.data = append(g.data, gasObservation{At: old, WeiPerGas: 999})

	// Trigger prune via a fresh Observe (which calls pruneLocked).
	g.Observe(big.NewInt(1))

	if g.Size() != 1 {
		t.Errorf("Size = %d, want 1 (old sample should be pruned)", g.Size())
	}
}
