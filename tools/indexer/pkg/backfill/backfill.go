package backfill

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"github.com/primev/mev-commit/tools/indexer/pkg/ingest"
)

type SlotData struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
	ProposerIdx     *int64
}


func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, relays []relay.Row) error {
	logger := slog.With("component", "backfill")
	logger.Info("Starting streaming backfill")

	if err := ctx.Err(); err != nil {
		return err
	}

	lastSlotNumber, _ := db.GetMaxSlotNumber(ctx)
	startSlot := lastSlotNumber - cfg.BackfillLookback
	if startSlot < 0 {
		startSlot = 0
	}

	batch := int64(cfg.BackfillBatch)
	totalBatches := (cfg.BackfillLookback + batch - 1) / batch

	logger.Info("Starting backfill",
		"start_slot", startSlot,
		"end_slot_exclusive", lastSlotNumber,
		"lookback_slots", cfg.BackfillLookback,
		"batch_size", cfg.BackfillBatch,
		"total_batches", totalBatches,
	)

	batchSz := int64(cfg.BackfillBatch)
	var processed int64

	for batchIdx := int64(0); batchIdx < totalBatches; batchIdx++ {
		batchStart := startSlot + batchIdx*batch
		batchEnd := batchStart + batchSz
		if batchEnd > lastSlotNumber {
			batchEnd = lastSlotNumber
		}

		logger.Info("Batch begin",
			"batch", batchIdx+1, "of", totalBatches,
			"range", fmt.Sprintf("[%d,%d)", batchStart, batchEnd),
		)

		for slot := batchStart; slot < batchEnd; slot++ {
			tTotal := time.Now()

			convCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			blockNumber, err := ethereum.SlotToExecutionBlockNumber(convCtx, httpc, cfg.BeaconBase, slot)
			cancel()
			if err != nil {
				logger.Error("Failed to convert slot to block number", "slot", slot, "error", err)
				continue
			}
			if blockNumber != 0 {
				if err := ctx.Err(); err != nil {
					return err
				}
				fetchCtx, fetchCancel := context.WithTimeout(ctx, 5*time.Second)
				ei, ferr := beacon.FetchBeaconExecutionBlock(fetchCtx, httpc, cfg.BeaconBase, blockNumber)
				fetchCancel()
				if ferr != nil || ei == nil {
					logger.Error("beacon fetch failed", "block", blockNumber, "error", ferr)
					continue
				}

				if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
					logger.Error("block upsert failed", "slot", ei.Slot, "error", err)
					continue
				}

				// Validator pubkey + opt-in status
				if err := ingest.LaunchValidatorTasks(ctx, cfg, db, httpc, ei, cfg.BeaconBase, logger); err != nil {
					logger.Error("failed to launch async validator tasks", "slot", ei.Slot, "error", err)

				}
				// Relay mode: fetch & insert bids
				if cfg.RelayData {
					if err := ingest.ProcessBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
						logger.Error("failed to process bids", "error", err)
					}
				}

				processed++
				totalMS := time.Since(tTotal).Milliseconds()
				logger.Info("total time taken", "total_ms", totalMS)
			}

			logger.Info("Batch end",
				"batch", batchIdx+1,
				"processed_slots_so_far", processed,
			)
		}
	}
	logger.Info("Backfill completed", "total_slots_processed", processed)
	return nil
}
