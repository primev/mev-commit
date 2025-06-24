package pricer_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
)

func TestEstimatePrice(t *testing.T) {
	t.Parallel()

	bp := pricer.BidPricer{}

	ctx := context.Background()
	txn := types.NewTransaction(
		0,
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		big.NewInt(1000000000), // 1 Gwei
		21000,                  // gas limit
		big.NewInt(1000000000), // gas price
		nil,                    // no data
	)

	price, err := bp.EstimatePrice(ctx, txn)
	if err != nil {
		t.Fatalf("failed to estimate price: %v", err)
	}

	if price.BlockNumber == 0 {
		t.Error("expected non-zero block number in estimated price")
	}

	if price.BidAmount.Cmp(big.NewInt(0)) <= 0 {
		t.Error("expected estimated price to be greater than zero")
	}

	t.Logf("Estimated price: %s at block %d", price.BidAmount.String(), price.BlockNumber)
}
