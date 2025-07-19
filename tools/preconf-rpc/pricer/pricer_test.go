package pricer_test

import (
	"context"
	"testing"

	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
)

func TestEstimatePrice(t *testing.T) {
	t.Parallel()

	bp := pricer.NewPricer("")

	ctx := context.Background()

	prices, err := bp.EstimatePrice(ctx)
	if err != nil {
		t.Fatalf("failed to estimate price: %v", err)
	}

	if prices.CurrentBlockNumber == 0 {
		t.Error("expected non-zero current block number")
	}

	if len(prices.Prices) == 0 {
		t.Error("expected at least one block price")
	}

	price := prices.Prices[0]

	if price.BlockNumber == 0 {
		t.Error("expected non-zero block number in price")
	}

	if len(price.EstimatedPrices) == 0 {
		t.Error("expected at least one estimated price")
	}

	for _, estPrice := range price.EstimatedPrices {
		if estPrice.PriorityFeePerGasGwei <= 0 {
			t.Errorf("expected positive priority fee per gas, got %f", estPrice.PriorityFeePerGasGwei)
		}
	}

	t.Logf("Estimated prices: %v", prices)
}
