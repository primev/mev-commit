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
	ProposerIdx     *int64
}

func RunAll(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, cfg *config.Config, relays []relay.Row) error {
	logger := slog.With("component", "backfill")
	logger.Info("Starting streaming backfill")

	if err := ctx.Err(); err != nil {
		return err
	}

	blocks, err := db.GetRecentMissingBlocks(ctx, cfg.BackfillLookback, cfg.BackfillBatch)
	if err != nil {
		return fmt.Errorf("get missing blocks: %w", err)
	}

	for _, b := range blocks {
		if err := ctx.Err(); err != nil {
			return err
		}

		fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		ei, ferr := beacon.FetchBeaconExecutionBlock(fetchCtx, httpc, cfg.BeaconBase, b.BlockNumber)
		cancel()
		if ferr != nil || ei == nil {
			logger.Error("beacon fetch failed", "block", b.BlockNumber, "error", ferr)
			continue
		}

		if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
			logger.Error("block upsert failed", "slot", ei.Slot, "error", err)
			continue
		}

		var vpub []byte
		if ei.ProposerIdx != nil {
			vctx, vcancel := context.WithTimeout(ctx, 5*time.Second)
			v, verr := beacon.FetchValidatorPubkey(vctx, httpc, cfg.BeaconBase, *ei.ProposerIdx)
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

		for _, r := range relays {
			if err := ctx.Err(); err != nil {
				return err
			}

			bctx, bcancel := context.WithTimeout(ctx, 5*time.Second)
			bids, berr := relay.FetchBuilderBlocksReceived(bctx, httpc, r.URL, ei.Slot)
			bcancel()
			if berr != nil {
				logger.Debug("bid fetch failed", "slot", ei.Slot, "relay", r.ID, "error", berr)
				continue
			}

			if len(bids) == 0 {
				continue
			}

			rows := make([]database.BidRow, 0, len(bids))
			for _, bid := range bids {
				if row, ok := relay.BuildBidInsert(ei.Slot, r.ID, bid); ok {
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
		}
		logger.Debug("slot processed", "slot", ei.Slot)
	}

	logger.Info("Backfill completed", "blocks_processed", len(blocks))
	return nil
}
