package ethereum

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
)

// SlotToExecutionBlockNumber converts a beacon chain slot to an execution layer block number
// using the beaconcha.in API with retryable HTTP client and context support.
func SlotToExecutionBlockNumber(ctx context.Context, httpc *retryablehttp.Client, beaconBase string, slot int64) (int64, error) {
	url := fmt.Sprintf("%s/slot/%d", beaconBase, slot)

	if _, has := ctx.Deadline(); !has {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	var wrap struct {
		Status string `json:"status"`
		Data   struct {
			ExecutionBlockNumber int64  `json:"exec_block_number"`
			Status               string `json:"status"`
		} `json:"data"`
	}

	if err := httputil.FetchJSON(ctx, httpc, url, &wrap); err != nil {
		return 0, fmt.Errorf("failed to fetch slot %d: %w", slot, err)
	}

	if wrap.Status != "OK" {
		return 0, fmt.Errorf("API returned status: %s for slot %d", wrap.Status, slot)
	}

	// Check if slot was missed
	if wrap.Data.Status != "1" {
		return 0, nil // Missed slot, not an error
	}

	if wrap.Data.ExecutionBlockNumber == 0 {
		return 0, fmt.Errorf("no execution block number for slot %d", slot)
	}

	return wrap.Data.ExecutionBlockNumber, nil
}
