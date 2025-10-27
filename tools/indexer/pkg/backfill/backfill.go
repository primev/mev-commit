package backfill

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"golang.org/x/time/rate"
)

type SlotData struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
	ProposerIdx     *int64
}

func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, cfg *config.Config, relays []relay.Row) error {
	logger := slog.With("component", "backfill")
	logger.Info("Starting streaming backfill")

	if err := ctx.Err(); err != nil {
		return err
	}
	lastSlotNumber, _ := db.GetMaxSlotNumber(ctx)
	startSlot := lastSlotNumber - cfg.BackfillLookback

	batch := cfg.BackfillBatch
	totalBatches := (cfg.BackfillLookback + int64(batch) - 1) / int64(batch)
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
		batchStart := startSlot + batchIdx*int64(batch)
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
			conversCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			blockNumber, err := ethereum.SlotToExecutionBlockNumber(conversCtx, httpc, cfg.BeaconBase, slot)
			cancel()
			if err != nil {
				logger.Error("Failed to convert slot to block number", "slot", slot, "error", err)
				continue
			}
			if blockNumber != 0 {

				if err := ctx.Err(); err != nil {
					return err
				}

				fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				ei, ferr := beacon.FetchBeaconExecutionBlock(fetchCtx, httpc, beaconLimiter, cfg.BeaconBase, cfg.BeaconchaAPIKey, blockNumber)
				cancel()
				if ferr != nil || ei == nil {
					logger.Error("beacon fetch failed", "block", blockNumber, "error", ferr)
					continue
				}

				if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
					logger.Error("block upsert failed", "slot", ei.Slot, "error", err)
					continue
				}

				var vpub []byte
				if ei.ProposerIdx != nil {
					vctx, vcancel := context.WithTimeout(ctx, 5*time.Second)
					v, verr := beacon.FetchValidatorPubkey(vctx, httpc, beaconLimiter, cfg.BeaconBase, cfg.BeaconchaAPIKey, *ei.ProposerIdx)
					vcancel()
					if verr != nil {
						logger.Error("validator fetch failed", "slot", ei.Slot, "error", verr)
					} else if len(v) > 0 {
						vpub = v

						// Save validator pubkey
						if err := db.UpdateValidatorPubkey(ctx, ei.Slot, vpub); err != nil {
							logger.Error("validator update failed", "slot", ei.Slot, "error", err)
						} else {
							opted, oerr := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, cfg, ei.BlockNumber, vpub)
							if oerr != nil {
								logger.Error("opt-in check failed", "slot", ei.Slot, "error", oerr)
							} else {
								updCtx, updCancel := context.WithTimeout(ctx, 3*time.Second)
								if uerr := db.UpdateValidatorOptInStatus(updCtx, ei.Slot, opted); uerr != nil {
									logger.Error("opt-in update failed", "slot", ei.Slot, "error", uerr)
								}
								updCancel()
							}
						}
					}
				}
				tBids := time.Now()
				var wg sync.WaitGroup
				for _, r := range relays {
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
					}(r)
				}
				wg.Wait()
				bidsMs := time.Since(tBids).Milliseconds()
				logger.Info("Bids fetch and insert", "bids_ms", bidsMs)
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
