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
)

type SlotData struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
}

// RecentMissing backfills recent blocks that are missing data
func recentMissing(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, lookback int64, batch int, ch chan<- SlotData) error {
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
		var vpub []byte
		// Schedule async validator pubkey fetch
		if ei.ProposerIdx != nil {
			vctx, vcancel := context.WithTimeout(ctx, 5*time.Second)
			v, verr := beacon.FetchValidatorPubkey(vctx, httpc, cfg.BeaconBase, *ei.ProposerIdx)
			vcancel()

			if verr != nil {
				return fmt.Errorf("validator pubkey fetch failed slot=%d: %w", ei.Slot, verr)
			} else if len(v) > 0 {
				vpub = v
				if err := db.UpdateValidatorPubkey(ctx, ei.Slot, vpub); err != nil {
					return fmt.Errorf("validator pubkey update failed slot=%d: %w", ei.Slot, err)
				}
			}
		}
		select {
		case ch <- SlotData{Slot: ei.Slot, BlockNumber: ei.BlockNumber, ValidatorPubkey: vpub}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	logger.Info("RecentMissing processed", "blocks", len(blocks))
	return nil
}

// RecentBids backfills bid data for ALL recent slots (not just opted-in blocks)
func recentBids(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, slots []int64) error {
	logger := slog.With("component", "backfill")

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

// RunAll executes all backfill operations ensuring complete coverage
func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, relays []relay.Row) error {
	logger := slog.With("component", "backfill")
	logger.Info("Starting comprehensive backfill for ALL blocks (not just opted-in)")

	// Channel to pass slot data from stage 1 to stages 2 & 3
	slotChan := make(chan SlotData, cfg.BackfillBatch)
	errCh := make(chan error, 1)

	// Run recentMissing and collect slot data
	go func() {
		if err := recentMissing(ctx, db, httpc, cfg, cfg.BackfillLookback, cfg.BackfillBatch, slotChan); err != nil {
			errCh <- err
		}
		close(slotChan)
	}()

	// Collect slots and validator data from channel
	var slotsForBids []int64
	var validatorsToCheck []SlotData

	for data := range slotChan {
		slotsForBids = append(slotsForBids, data.Slot)
		if len(data.ValidatorPubkey) > 0 {
			validatorsToCheck = append(validatorsToCheck, data)
		}
	}
	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("RecentMissing failed", "error", err)
			return err
		}
	default:
	}
	for _, v := range validatorsToCheck {
		opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, cfg, v.BlockNumber, v.ValidatorPubkey)
		if err == nil {
			db.UpdateValidatorOptInStatus(ctx, v.Slot, opted)
		}
	}
	if err := recentBids(ctx, db, httpc, relays, slotsForBids); err != nil {
		logger.Error("RecentBids failed", "error", err)
		return err
	}

	logger.Info("Backfill-done")
	return nil

}
