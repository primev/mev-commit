package backfill

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
)

type SlotData struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
	ProposerIdx     *int64
}

func LaunchValidatorTasks(ctx context.Context, cfg *config.Config, db *database.DB, httpc *retryablehttp.Client, ei *beacon.ExecInfo, beaconBase string, logger *slog.Logger) error { // Async validator pubkey fetch
	if ei.ProposerIdx == nil {
		return nil
	}

	vctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	vpub, err := beacon.FetchValidatorPubkey(vctx, httpc, beaconBase, *ei.ProposerIdx)
	if err != nil {
		return fmt.Errorf("fetch validator pubkey: %w", err)
	}

	if len(vpub) > 0 {
		if err := db.UpdateValidatorPubkey(vctx, ei.Slot, vpub); err != nil {
			logger.Error("validator failed to save pubkey", "slot", ei.Slot, "error", err)
		} else {
			logger.Info("validator pubkey saved", "proposer", *ei.ProposerIdx, "slot", ei.Slot)
		}
	}

	// Wait for validator pubkey to be available
	getCtx, getCancel := context.WithTimeout(context.Background(), 5*time.Second)
	vpk, err := db.GetValidatorPubkeyWithRetry(getCtx, ei.Slot, 3, time.Second)
	getCancel()

	if err != nil {
		logger.Error("validator pubkey not available", "slot", ei.Slot, "error", err)
		return fmt.Errorf("save validator pubkey: %w", err)
	}

	opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, cfg, ei.BlockNumber, vpk)
	if err != nil {
		return fmt.Errorf("check opt-in status: %w", err)
	}

	updCtx, updCancel := context.WithTimeout(context.Background(), 3*time.Second)
	err = db.UpdateValidatorOptInStatus(updCtx, ei.Slot, opted)
	updCancel()
	if err != nil {
		return fmt.Errorf("save opt-in status: %w", err)
	} else {
		logger.Info("validator opt-in status", "slot", ei.Slot, "opted_in", opted)
	}
	return nil

}

func ProcessBidsForBlock(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, ei *beacon.ExecInfo, logger *slog.Logger) error {
	logger.Info("processing bids for block", "block", ei.BlockNumber, "slot", ei.Slot)
	var wg sync.WaitGroup
	var totalBids int64
	var successfulRelays int64
	for _, r := range relays {
		r := r
		wg.Add(1)
		go func(rel relay.Row) {
			defer wg.Done()
			if err := ctx.Err(); err != nil {
				return
			}
			bctx, bcancel := context.WithTimeout(ctx, 5*time.Second)
			bids, berr := relay.FetchBuilderBlocksReceived(bctx, httpc, rel.URL, ei.Slot)
			bcancel()
			if berr != nil {
				logger.Debug("bid fetch failed", "slot", ei.Slot, "relay", r.ID, "error", berr)
				return
			}
			if len(bids) == 0 {
				return
			}
			rows := make([]database.BidRow, 0, len(bids))
			for _, bid := range bids {
				if row, ok := relay.BuildBidInsert(ei.Slot, rel.ID, bid); ok {
					rows = append(rows, row)
				}
			}
			if len(rows) > 0 {
				insCtx, insCancel := context.WithTimeout(ctx, 5*time.Second)
				if ierr := db.InsertBidsBatch(insCtx, rows); ierr != nil {
					logger.Error("bid insert failed", "slot", ei.Slot, "relay", r.ID, "error", ierr)
				}
				insCancel()
			}
			atomic.AddInt64(&totalBids, int64(len(rows)))
			atomic.AddInt64(&successfulRelays, 1)
			logger.Info("bid insert ok",
				"slot", ei.Slot, "relay_id", rel.ID,
				"bids_in", len(bids), "rows_out", len(rows),
			)
		}(r)
	}
	wg.Wait()

	logger.Info("summary", "block", ei.BlockNumber, "total_bids", totalBids, "successful_relays", successfulRelays)
	return nil
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
				if err := LaunchValidatorTasks(ctx, cfg, db, httpc, ei, cfg.BeaconBase, logger); err != nil {
					logger.Error("failed to launch async validator tasks", "slot", ei.Slot, "error", err)

				}
				// Relay mode: fetch & insert bids
				if cfg.RelayData {
					if err := ProcessBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
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
