package pricer_test

import (
	"context"
	"io"
	"testing"

	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/x/util"
)

func TestEstimatePrice(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)
	bp, err := pricer.NewPricer("", logger)
	if err != nil {
		t.Fatalf("failed to create pricer: %v", err)
	}

	ctx := context.Background()

	prices := bp.EstimatePrice(ctx)

	if len(prices) != 3 {
		t.Fatalf("expected 3 confidence levels, got %d", len(prices))
	}

	for confidence, price := range prices {
		if confidence <= 0 {
			t.Errorf("expected positive confidence level, got %d", confidence)
		}
		if price <= 0 {
			t.Errorf("expected positive price, got %f", price)
		}
	}

	t.Logf("Estimated prices: %v", prices)
}
