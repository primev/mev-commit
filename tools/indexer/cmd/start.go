package main

import (
	"context"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/tools/indexer/pkg/backfill"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"github.com/urfave/cli/v2"
)

func startIndexer(c *cli.Context) error {

	initLogger := slog.With("component", "init")

	dbURL := c.String(optionDatabaseURL.Name)
	infuraRPC := c.String(optionInfuraRPC.Name)
	beaconBase := c.String(optionBeaconBase.Name)

	initLogger.Info("starting blockchain indexer with StarRocks database")
	initLogger.Info("configuration loaded",
		"block_interval", c.Duration("block-interval"),
		"validator_delay", c.Duration("validator-delay"))
	ctx := c.Context
	// Connect to StarRocks database
	db, err := database.Connect(ctx, dbURL, 20, 5)
	if err != nil {
		initLogger.Error("[DB] connection failed", "error", err)
		return err
	}

	defer func() {
		if cerr := db.Close(); cerr != nil {
			initLogger.Error("[DB] close failed", "error", cerr)
		}
	}()
	initLogger.Info("[DB] connected to StarRocks database")

	// Ensure required tables exist
	if err := db.EnsureStateTable(ctx); err != nil {
		initLogger.Error("[DB] failed to ensure state table", "error", err)
		return err
	}
	initLogger.Info("[DB] state table ready")

	// Initialize HTTP client
	httpc := httputil.NewHTTPClient(c.Duration("http-timeout"))
	initLogger.Info("[HTTP] client initialized", "timeout", c.Duration("http-timeout"))

	// Load relay configurations
	relays, err := relay.UpsertRelaysAndLoad(ctx, db)
	if err != nil {
		initLogger.Error("[RELAY] failed to load", "error", err)
	}
	initLogger.Info("[RELAY] loaded active relays", "count", len(relays))
	for _, r := range relays {
		initLogger.Info("[RELAY] relay found", "id", r.ID, "url", r.URL)
	}

	// Initialize starting block number
	lastBN, found := db.LoadLastBlockNumber(ctx)
	if !found || lastBN == 0 {
		initLogger.Info("no previous state found, checking database for latest block")
		lastBN, err = db.GetMaxBlockNumber(ctx)
		if err != nil {
			initLogger.Error("database query failed", "error", err)
		}
	}

	if lastBN == 0 {
		initLogger.Info("getting latest block from Ethereum RPC...")

		latestBlock, err := ethereum.GetLatestBlockNumber(httpc.HTTPClient, infuraRPC)
		if err != nil {
			initLogger.Error("failed to get latest block from RPC", "error", err)
			return err
		}

		lastBN = latestBlock - 10 // Start 10 blocks behind to ensure data availability
		initLogger.Info("starting from block", "block", lastBN, "latest", latestBlock)
	}

	initLogger.Info("starting from block number", "block", lastBN)
	initLogger.Info("indexer configuration", "lookback", c.Int("backfill-lookback"), "batch", c.Int("backfill-batch"))

	if c.Int("backfill-lookback") > 0 {
		initLogger.Info("[BACKFILL] running one-time backfill",
			"lookback", c.Int("backfill-lookback"),
			"batch", c.Int("backfill-batch"))
		if err := backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays); err != nil {
			initLogger.Error("[BACKFILL] failed", "error", err)
		} else {
			initLogger.Info("[BACKFILL] completed startup backfill")
		}
	} else {
		initLogger.Info("[BACKFILL] skipped", "reason", "backfill-lookback=0")
	}

	mainTicker := time.NewTicker(c.Duration("block-interval"))
	defer mainTicker.Stop()

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			initLogger.Info("[SHUTDOWN] graceful shutdown initiated", "reason", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				initLogger.Error("[SHUTDOWN] failed to save last block number", "error", err)
			}
			initLogger.Info("[SHUTDOWN] indexer stopped", "block", lastBN)
			return nil

		case <-mainTicker.C:
			nextBN := lastBN + 1

			// Fetch execution block data
			ei, err := beacon.FetchCombinedBlockData(ctx, httpc, infuraRPC, beaconBase, nextBN)
			if err != nil || ei == nil {
				initLogger.Warn("[BLOCK] not available yet", "block", nextBN, "error", err)
				continue
			}
			fields := []any{
				"block", nextBN,
				"slot", ei.Slot,
			}
			if ei.Timestamp != nil {
				fields = append(fields, "timestamp", ei.Timestamp.Format(time.RFC3339))
			}
			if ei.ProposerIdx != nil {
				fields = append(fields, "proposer_index", *ei.ProposerIdx)
			}
			if ei.RelayTag != nil {
				fields = append(fields, "winning_relay", *ei.RelayTag)
			}
			if ei.BuilderHex != nil {
				pref := *ei.BuilderHex
				if len(pref) > 20 {
					pref = pref[:20]
				}
				fields = append(fields, "builder_pubkey_prefix", pref)
			}
			if ei.RewardEth != nil {
				fields = append(fields, "producer_reward_eth", *ei.RewardEth)
			}
			initLogger.Info("processing block", fields...)

			// Save block to database
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				initLogger.Error("[DB] failed to save block", "block", nextBN, "error", err)
				continue
			}
			initLogger.Info("[DB] block saved successfully", "block", nextBN)

			// Fetch and store bid data from all relays
			totalBids := 0
			successfulRelays := 0
			mainContextCanceled := false
			const batchSize = 500
			for _, rr := range relays {
				// Check if main context is canceled before processing each relay
				if ctx.Err() != nil {
					initLogger.Warn("main context canceled, stopping relay processing")
					mainContextCanceled = true
					break
				}

				bids, err := relay.FetchBuilderBlocksReceived(ctx, httpc, rr.URL, ei.Slot)
				if err != nil {
					initLogger.Error("[RELAY] failed to fetch bids", "relay_id", rr.ID, "url", rr.URL, "error", err)
					continue
				}

				relayBids := 0
				batch := make([]database.BidRow, 0, batchSize)

				for _, bid := range bids {
					// Check if main context is still valid
					if ctx.Err() != nil {
						initLogger.Warn("[BIDS] main context canceled, stopping bid insertion")
						mainContextCanceled = true
						break
					}

					if row, ok := relay.BuildBidInsert(ei.Slot, rr.ID, bid); ok {
						batch = append(batch, row)

						if len(batch) >= batchSize {
							insCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
							if err := db.InsertBidsBatch(insCtx, batch); err != nil {

								initLogger.Error("[DB]batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
							} else {
								relayBids += len(batch)
							}
							cancel()
							batch = batch[:0]
						}
					}
				}

				// final flush
				if len(batch) > 0 {
					if err := db.InsertBidsBatch(ctx, batch); err != nil {
						initLogger.Error("[DB] batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
					} else {
						relayBids += len(batch)
					}
				}

				if mainContextCanceled {
					break
				}

				if relayBids > 0 {
					initLogger.Info("[BIDS] bids collected", "relay_id", rr.ID, "count", relayBids)
					totalBids += relayBids
					successfulRelays++
				}
			}

			initLogger.Info("[BIDS] summary", "block", nextBN, "total_bids", totalBids, "successful_relays", successfulRelays)
			// Async validator pubkey fetch
			if ei.ProposerIdx != nil {
				go func(slot int64, proposerIdx int64) {
					time.Sleep(c.Duration("validator-delay"))

					vpub, err := beacon.FetchValidatorPubkey(ctx, httpc, beaconBase, proposerIdx)
					if err != nil {
						initLogger.Error("[VALIDATOR] failed to fetch pubkey", "proposer", proposerIdx, "error", err)
						return
					}

					if len(vpub) > 0 {
						if err := db.UpdateValidatorPubkey(ctx, slot, vpub); err != nil {
							initLogger.Error("[VALIDATOR] failed to save pubkey", "slot", slot, "error", err)
						} else {
							initLogger.Info("[VALIDATOR] pubkey saved", "proposer", proposerIdx, "slot", slot)
						}
					}
				}(ei.Slot, *ei.ProposerIdx)
			}

			// Async opt-in status check
			if ei.ProposerIdx != nil {
				go func(slot int64, blockNumber int64) {
					time.Sleep(c.Duration("validator-delay") + 500*time.Millisecond)

					// Wait for validator pubkey to be available
					vpk, err := db.GetValidatorPubkeyWithRetry(ctx, slot, 3, time.Second)
					if err != nil {
						initLogger.Error("[VALIDATOR] pubkey not available", "slot", slot, "error", err)
						return
					}

					opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, createOptionsFromCLI(c), blockNumber, vpk)
					if err != nil {
						initLogger.Error("[OPT-IN] failed to check opt-in status", "slot", slot, "error", err)
						return
					}

					err = db.UpdateValidatorOptInStatus(ctx, slot, opted)
					if err != nil {
						initLogger.Error("[OPT-IN] failed to save opt-in status", "slot", slot, "error", err)
					} else {
						initLogger.Info("[OPT-IN] validator opt-in status", "slot", slot, "opted_in", opted)
					}
				}(ei.Slot, ei.BlockNumber)
			}

			lastBN = nextBN
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				initLogger.Error("[PROGRESS] failed to save block number", "block", lastBN, "error", err)
			} else {
				initLogger.Info("[PROGRESS] advanced to block", "block", lastBN)
			}
		}
	}
}
