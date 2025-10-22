package ingest

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
