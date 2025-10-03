package backfill

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
)

// RecentMissing backfills recent blocks that are missing data
func RecentMissing(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")

	blocks, err := db.GetRecentMissingBlocks(ctx, lookback, batch)
	if err != nil {
		logger.Error("RecentMissing query failed", "error", err)
		return err
	}

	processed := 0
	for _, block := range blocks {

		// Fetch beacon execution block data
		if ei, err := beacon.FetchBeaconExecutionBlock(httpc, cfg.BeaconBase, block.BlockNumber); err == nil && ei != nil {
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				logger.Error("RecentMissing upsert failed", "slot", block.Slot, "error", err)

				continue
			}

			// Schedule async validator pubkey fetch
			if ei.ProposerIdx != nil {
				go func(slot int64, idx int64) {
					time.Sleep(cfg.ValidatorWait)
					if vpub, err := beacon.FetchValidatorPubkey(httpc, cfg.BeaconBase, idx); err == nil && len(vpub) > 0 {
						if err := db.UpdateValidatorPubkey(context.Background(), slot, vpub); err != nil {
							logger.Error("RecentMissing validator pubkey update failed", "slot", slot, "error", err)
						}
					}
				}(ei.Slot, *ei.ProposerIdx)
			}
			processed++
		} else {
			logger.Error("RecentMissing beacon fetch failed", "block_number", block.BlockNumber, "error", err)

		}
	}

	logger.Info("RecentMissing processed", "blocks", processed)
	return nil
}

// RecentBids backfills bid data for ALL recent slots (not just opted-in blocks)
func RecentBids(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")
	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	slots, err := db.GetRecentSlotsWithBlocks(opCtx, lookback, batch)
	if err != nil {
		logger.Error("RecentBids query failed", "error", err)
		return err
	}

	logger.Info("RecentBids fetched slots", "count", len(slots))
	processed := 0
	totalBids := 0
	const batchSize = 500

	for _, slot := range slots {
		if ctx.Err() != nil {
			break
		}

		slotBids := 0

		for _, rr := range relays {
			if ctx.Err() != nil {
				break
			}

			fetchCtx, fcancel := context.WithTimeout(ctx, 5*time.Second)
			bids, err := relay.FetchBuilderBlocksReceived(fetchCtx, httpc, rr.URL, slot)
			fcancel()
			if err != nil {
				logger.Error("RecentBids fetch failed", "slot", slot, "relay_id", rr.ID, "relay_url", rr.URL, "error", err)
				continue
			}

			rows := make([]database.BidRow, 0, len(bids))
			for _, b := range bids {
				if row, ok := relay.BuildBidInsert(slot, rr.ID, b); ok {
					rows = append(rows, row)
				}
			}

			for i := 0; i < len(rows); i += batchSize {
				end := i + batchSize
				if end > len(rows) {
					end = len(rows)
				}
				insCtx, icancel := context.WithTimeout(ctx, 5*time.Second)
				if err := db.InsertBidsBatch(insCtx, rows[i:end]); err != nil {
					logger.Error("RecentBids batch insert failed", "slot", slot, "relay_id", rr.ID, "count", end-i, "error", err)
				} else {
					slotBids += end - i
				}
				icancel()
				if ctx.Err() != nil {
					break
				}
			}
		}

		if slotBids > 0 {
			totalBids += slotBids
			processed++
		}
	}

	logger.Info("RecentBids processed", "slots", processed, "total_bids", totalBids)
	return nil
}

// ValidatorOptIn backfills validator opt-in status (this is opt-in specific data)
func ValidatorOptIn(ctx context.Context, db *database.DB, httpc *http.Client, cfg *config.Config, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")
	validators, err := db.GetValidatorsNeedingOptInCheck(ctx, lookback, batch)
	if err != nil {
		logger.Error("ValidatorOptIn query failed", "error", err)
		return err
	}

	processed := 0

	for _, v := range validators {
		opted, err := ethereum.CallAreOptedInAtBlock(httpc, cfg, v.BlockNumber, v.ValidatorPubkey)
		if err == nil {
			if err := db.UpdateValidatorOptInStatus(ctx, v.Slot, opted); err != nil {
				logger.Error("ValidatorOptIn update failed", "slot", v.Slot, "error", err)
			} else {
				processed++
			}
		} else {
			logger.Error("ValidatorOptIn check failed", "slot", v.Slot, "error", err)
		}
	}

	logger.Info("ValidatorOptIn processed", "validators", processed)
	return nil
}

// AllBlocksBids ensures bid data is collected for ALL blocks, regardless of opt-in status
func AllBlocksBids(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, startSlot, endSlot int64) error {
	logger := slog.With("component", "backfill")
	logger.Info("AllBlocksBids starting", "start_slot", startSlot, "end_slot", endSlot)
	const batchSize = 500
	totalProcessed := 0
	totalBids := 0

	for slot := startSlot; slot <= endSlot; slot++ {
		if ctx.Err() != nil {
			break
		}
		slotBids := 0

		for _, rr := range relays {
			fetchCtx, fcancel := context.WithTimeout(ctx, 5*time.Second)
			bids, err := relay.FetchBuilderBlocksReceived(fetchCtx, httpc, rr.URL, slot)
			fcancel()
			if err != nil {
				logger.Error("AllBlocksBids fetch failed", "slot", slot, "relay_id", rr.ID, "url", rr.URL, "error", err)
				continue
			}

			rows := make([]database.BidRow, 0, len(bids))
			for _, b := range bids {
				if row, ok := relay.BuildBidInsert(slot, rr.ID, b); ok {
					rows = append(rows, row)
				}
			}

			for i := 0; i < len(rows); i += batchSize {
				end := i + batchSize
				if end > len(rows) {
					end = len(rows)
				}

				insCtx, icancel := context.WithTimeout(ctx, 5*time.Second)
				if err := db.InsertBidsBatch(insCtx, rows[i:end]); err != nil {
					logger.Error("AllBlocksBids batch insert failed", "slot", slot, "relay_id", rr.ID, "count", end-i, "error", err)
				} else {
					slotBids += end - i
				}
				icancel()

				if ctx.Err() != nil {
					break
				}
			}
		}

		if slotBids > 0 {
			totalBids += slotBids
			totalProcessed++
		}

		// Respect context cancellation
		select {
		case <-ctx.Done():
			logger.Warn("AllBlocksBids cancelled", "current_slot", slot)
			return ctx.Err()
		default:
		}
	}

	logger.Info("AllBlocksBids completed", "slots", totalProcessed, "total_bids", totalBids)
	return nil
}

// RunAll executes all backfill operations ensuring complete coverage
func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, relays []relay.Row) {
	logger := slog.With("component", "backfill")
	logger.Info("Starting comprehensive backfill for ALL blocks (not just opted-in)")

	// Run backfill operations
	if err := RecentMissing(ctx, db, httpc, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("RecentMissing failed", "error", err)
	}

	if err := ValidatorOptIn(ctx, db, httpc.HTTPClient, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("ValidatorOptIn failed", "error", err)
	}

	// This ensures bid data for ALL blocks, not just mev-commit opted-in blocks
	if err := RecentBids(ctx, db, httpc, relays, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("RecentBids failed", "error", err)
	}

	logger.Info("All operations completed - relay data covers ALL blocks")
}
