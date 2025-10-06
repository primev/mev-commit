package backfill

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"log/slog"
	"net/http"
	"time"
)

// RecentMissing backfills recent blocks that are missing data
func recentMissing(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")

	blocks, err := db.GetRecentMissingBlocks(ctx, lookback, batch)
	if err != nil {
		logger.Error("RecentMissing query failed", "error", err)
		return err
	}

	for _, block := range blocks {
		fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		// Fetch beacon execution block data
		ei, ferr := beacon.FetchBeaconExecutionBlock(fetchCtx, httpc, cfg.BeaconBase, block.BlockNumber)
		cancel()
		if ferr != nil || ei == nil {
			return fmt.Errorf("beacon fetch failed for block=%d: %w", block.BlockNumber, ferr)
		}
		if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
			logger.Error("RecentMissing upsert failed", "slot", block.Slot, "error", err)

			continue
		}

		// Schedule async validator pubkey fetch
		if ei.ProposerIdx != nil {
			vctx, vcancel := context.WithTimeout(ctx, 5*time.Second)
			vpub, verr := beacon.FetchValidatorPubkey(vctx, httpc, cfg.BeaconBase, *ei.ProposerIdx)
			vcancel()

			if verr != nil {
				return fmt.Errorf("validator pubkey fetch failed slot=%d: %w", ei.Slot, verr)
			}
			if len(vpub) > 0 {
				if err := db.UpdateValidatorPubkey(ctx, ei.Slot, vpub); err != nil {
					return fmt.Errorf("validator pubkey update failed slot=%d: %w", ei.Slot, err)
				}
			}
		}
	}

	logger.Info("RecentMissing processed", "blocks", len(blocks))
	return nil
}

// RecentBids backfills bid data for ALL recent slots (not just opted-in blocks)
func recentBids(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")
	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	slots, err := db.GetRecentSlotsWithBlocks(opCtx, lookback, batch)
	if err != nil {
		logger.Error("RecentBids query failed", "error", err)
		return err
	}

	logger.Info("RecentBids fetched slots", "count", len(slots))

	for _, slot := range slots {
		if ctx.Err() != nil {
			break
		}

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

			if len(rows) > 0 {
				insCtx, icancel := context.WithTimeout(ctx, 5*time.Second)
				if err := db.InsertBidsBatch(insCtx, rows); err != nil {
					icancel()
					return fmt.Errorf("bids insert failed slot=%d relay_id=%d: %w", slot, rr.ID, err)
				}
				icancel()
			}
		}

	}

	logger.Info("RecentBids processed", "slots", len(slots))
	return nil
}

// ValidatorOptIn backfills validator opt-in status (this is opt-in specific data)
func validatorOptIn(ctx context.Context, db *database.DB, httpc *http.Client, cfg *config.Config, lookback int64, batch int) error {
	logger := slog.With("component", "backfill")
	validators, err := db.GetValidatorsNeedingOptInCheck(ctx, lookback, batch)
	if err != nil {
		logger.Error("ValidatorOptIn query failed", "error", err)
		return err
	}

	for _, v := range validators {
		opted, err := ethereum.CallAreOptedInAtBlock(httpc, cfg, v.BlockNumber, v.ValidatorPubkey)
		if err == nil {
			if err := db.UpdateValidatorOptInStatus(ctx, v.Slot, opted); err != nil {
				logger.Error("ValidatorOptIn update failed", "slot", v.Slot, "error", err)
			}

		} else {
			logger.Error("ValidatorOptIn check failed", "slot", v.Slot, "error", err)
		}
	}

	logger.Info("ValidatorOptIn processed", "validators", len(validators))
	return nil
}

// RunAll executes all backfill operations ensuring complete coverage
func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, relays []relay.Row) error {
	logger := slog.With("component", "backfill")
	logger.Info("Starting comprehensive backfill for ALL blocks (not just opted-in)")

	if err := recentMissing(ctx, db, httpc, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("RecentMissing failed", "error", err)
		return err
	}
	if err := validatorOptIn(ctx, db, httpc.HTTPClient, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("ValidatorOptIn failed", "error", err)
		return err
	}
	if err := recentBids(ctx, db, httpc, relays, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		logger.Error("RecentBids failed", "error", err)
		return err
	}

	logger.Info("Backfill-done")
	return nil
}
