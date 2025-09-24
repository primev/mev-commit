// pkg/backfill/backfill.go
package backfill

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/primev/mev-commit/indexer/pkg/beacon"
	"github.com/primev/mev-commit/indexer/pkg/config"
	"github.com/primev/mev-commit/indexer/pkg/database"
	"github.com/primev/mev-commit/indexer/pkg/ethereum"
	"github.com/primev/mev-commit/indexer/pkg/relay"
)

// RecentMissing backfills recent blocks that are missing data
func RecentMissing(ctx context.Context, db *database.DB, httpc *http.Client, cfg *config.Config, lookback int64, batch int) error {

	rows, err := db.Conn.QueryContext(ctx, `
WITH recent AS (SELECT COALESCE(MAX(slot),0) AS s FROM blocks)
SELECT slot, block_number
FROM blocks, recent
WHERE slot > recent.s - ?
  AND block_number IS NOT NULL
  AND (winning_relay IS NULL
       OR winning_builder_pubkey IS NULL
       OR fee_recipient IS NULL
       OR producer_reward_eth IS NULL
       OR timestamp IS NULL
       OR proposer_index IS NULL)
ORDER BY slot DESC
LIMIT ?`, lookback, batch)
	if err != nil {
		log.Printf("[BACKFILL] RecentMissing query failed: %v", err)
		return err
	}
	defer rows.Close()

	processed := 0
	for rows.Next() {
		var slot int64
		var bn int64
		if err := rows.Scan(&slot, &bn); err != nil {
			log.Printf("[BACKFILL] RecentMissing scan failed: %v", err)
			continue
		}

		// Fetch beacon execution block data
		if ei, err := beacon.FetchBeaconExecutionBlock(httpc, cfg.BeaconBase, bn); err == nil && ei != nil {
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				log.Printf("[BACKFILL] RecentMissing upsert failed for slot %d: %v", slot, err)
				continue
			}

			// Schedule async validator pubkey fetch
			if ei.ProposerIdx != nil {
				go func(slot int64, idx int64) {
					time.Sleep(cfg.ValidatorWait)
					if vpub, err := beacon.FetchValidatorPubkey(httpc, cfg.BeaconBase, idx); err == nil && len(vpub) > 0 {
						if err := db.UpdateValidatorPubkey(context.Background(), slot, vpub); err != nil {
							log.Printf("[BACKFILL] RecentMissing validator pubkey update failed for slot %d: %v", slot, err)
						}
					}
				}(ei.Slot, *ei.ProposerIdx)
			}
			processed++
		} else {
			log.Printf("[BACKFILL] RecentMissing beacon fetch failed for block %d: %v", bn, err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("[BACKFILL] RecentMissing rows iteration failed: %v", err)
		return err
	}

	log.Printf("[BACKFILL] RecentMissing processed %d blocks", processed)
	return nil
}

// RecentBids backfills bid data for ALL recent slots (not just opted-in blocks)
func RecentBids(ctx context.Context, db *database.DB, httpc *http.Client, relays []relay.Row, lookback int64, batch int) error {

	rows, err := db.Conn.QueryContext(ctx, `
WITH recent AS (SELECT COALESCE(MAX(slot),0) AS s FROM blocks)
SELECT DISTINCT slot
FROM blocks, recent
WHERE slot > recent.s - ?
  AND block_number IS NOT NULL
ORDER BY slot DESC
LIMIT ?`, lookback, batch)
	if err != nil {
		log.Printf("[BACKFILL] RecentBids query failed: %v", err)
		return err
	}
	defer rows.Close()

	processed := 0
	totalBids := 0
	for rows.Next() {
		var slot int64
		if err := rows.Scan(&slot); err != nil {
			log.Printf("[BACKFILL] RecentBids scan failed: %v", err)
			continue
		}

		// Fetch bids from ALL relays for this slot
		slotBids := 0
		for _, rr := range relays {
			if bids, err := relay.FetchBuilderBlocksReceived(httpc, rr.URL, slot); err == nil {
				for _, b := range bids {
					if err := relay.InsertBid(ctx, db, slot, rr.ID, b); err != nil {
						log.Printf("[BACKFILL] RecentBids insert failed for slot %d, relay %d: %v", slot, rr.ID, err)
					} else {
						slotBids++
					}
				}
			} else {
				log.Printf("[BACKFILL] RecentBids fetch failed for slot %d, relay %d (%s): %v", slot, rr.ID, rr.URL, err)
			}
		}

		if slotBids > 0 {
			totalBids += slotBids
			processed++
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("[BACKFILL] RecentBids rows iteration failed: %v", err)
		return err
	}

	log.Printf("[BACKFILL] RecentBids processed %d slots with %d total bids from ALL blocks", processed, totalBids)
	return nil
}

// ValidatorOptIn backfills validator opt-in status (this is opt-in specific data)
func ValidatorOptIn(ctx context.Context, db *database.DB, httpc *http.Client, cfg *config.Config, lookback int64, batch int) error {

	rows, err := db.Conn.QueryContext(ctx, `
WITH recent AS (SELECT COALESCE(MAX(slot),0) AS s FROM blocks)
SELECT slot, block_number, validator_pubkey
FROM blocks, recent
WHERE slot > recent.s - ?
  AND block_number IS NOT NULL
  AND validator_pubkey IS NOT NULL
  AND validator_opted_in IS NULL
ORDER BY slot DESC
LIMIT ?`, lookback, batch)
	if err != nil {
		log.Printf("[BACKFILL] ValidatorOptIn query failed: %v", err)
		return err
	}
	defer rows.Close()

	processed := 0
	for rows.Next() {
		var slot, bn int64
		var vpk []byte
		if err := rows.Scan(&slot, &bn, &vpk); err != nil {
			log.Printf("[BACKFILL] ValidatorOptIn scan failed: %v", err)
			continue
		}

		// Check opt-in status
		opted, err := ethereum.CallAreOptedInAtBlock(httpc, cfg, bn, vpk)
		if err == nil {
			updateCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			// Fixed query parameter style for StarRocks
			if _, err := db.Conn.ExecContext(updateCtx, `UPDATE blocks SET validator_opted_in=? WHERE slot=? AND validator_opted_in IS NULL`, opted, slot); err != nil {
				log.Printf("[BACKFILL] ValidatorOptIn update failed for slot %d: %v", slot, err)
			} else {
				processed++
			}
			cancel()
		} else {
			log.Printf("[BACKFILL] ValidatorOptIn check failed for slot %d: %v", slot, err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("[BACKFILL] ValidatorOptIn rows iteration failed: %v", err)
		return err
	}

	log.Printf("[BACKFILL] ValidatorOptIn processed %d validators", processed)
	return nil
}

// AllBlocksBids ensures bid data is collected for ALL blocks, regardless of opt-in status
func AllBlocksBids(ctx context.Context, db *database.DB, httpc *http.Client, relays []relay.Row, startSlot, endSlot int64) error {
	log.Printf("[BACKFILL] AllBlocksBids ensuring bid coverage from slot %d to %d", startSlot, endSlot)

	totalProcessed := 0
	totalBids := 0

	for slot := startSlot; slot <= endSlot; slot++ {
		slotBids := 0

		// Fetch bids from ALL relays for every single slot
		for _, rr := range relays {
			if bids, err := relay.FetchBuilderBlocksReceived(httpc, rr.URL, slot); err == nil {
				for _, b := range bids {
					if err := relay.InsertBid(ctx, db, slot, rr.ID, b); err == nil {
						slotBids++
					}
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
			log.Printf("[BACKFILL] AllBlocksBids cancelled at slot %d", slot)
			return ctx.Err()
		default:
		}
	}

	log.Printf("[BACKFILL] AllBlocksBids completed: %d slots processed, %d total bids collected", totalProcessed, totalBids)
	return nil
}

// RunAll executes all backfill operations ensuring complete coverage
func RunAll(ctx context.Context, db *database.DB, httpc *http.Client, cfg *config.Config, relays []relay.Row) {
	log.Printf("[BACKFILL] Starting comprehensive backfill for ALL blocks (not just opted-in)")

	// Run backfill operations
	if err := RecentMissing(ctx, db, httpc, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		log.Printf("[BACKFILL] RecentMissing failed: %v", err)
	}

	if err := ValidatorOptIn(ctx, db, httpc, cfg, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		log.Printf("[BACKFILL] ValidatorOptIn failed: %v", err)
	}

	// This ensures bid data for ALL blocks, not just mev-commit opted-in blocks
	if err := RecentBids(ctx, db, httpc, relays, cfg.BackfillLookback, cfg.BackfillBatch); err != nil {
		log.Printf("[BACKFILL] RecentBids failed: %v", err)
	}

	log.Printf("[BACKFILL] All operations completed - relay data covers ALL blocks")
}
