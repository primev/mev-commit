package main

import (
	"context"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/backfill"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"

	"github.com/urfave/cli/v2"
)

func initializeDatabase(ctx context.Context, dbURL string, logger *slog.Logger) (*database.DB, error) {
	db, err := database.Connect(ctx, dbURL, 20, 5)
	if err != nil {
		logger.Error("[DB] connection failed", "error", err)
		return nil, err
	}
	logger.Info("[DB] connected to StarRocks database")

	if err := db.EnsureStateTable(ctx); err != nil {
		logger.Error("[DB] failed to ensure state table", "error", err)
		return nil, err
	}
	logger.Info("[DB] state table ready")

	return db, nil
}

func loadRelays(ctx context.Context, db *database.DB, logger *slog.Logger) ([]relay.Row, error) {
	relays, err := relay.UpsertRelaysAndLoad(ctx, db)
	if err != nil {
		logger.Error("[RELAY] failed to load", "error", err)
		return nil, err
	}

	logger.Info("[RELAY] loaded active relays", "count", len(relays))
	for _, r := range relays {
		logger.Info("[RELAY] relay found", "id", r.ID, "url", r.URL)
	}

	return relays, nil
}

func getStartingBlockNumber(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, infuraRPC string, logger *slog.Logger) (int64, error) {
	lastBN, found := db.LoadLastBlockNumber(ctx)

	if !found || lastBN == 0 {
		logger.Info("no previous state found, checking database for latest block")
		var err error
		lastBN, err = db.GetMaxBlockNumber(ctx)
		if err != nil {
			logger.Error("database query failed", "error", err)
		}
	}

	if lastBN == 0 {
		logger.Info("getting latest block from Ethereum RPC...")

		latestBlock, err := ethereum.GetLatestBlockNumber(httpc.HTTPClient, infuraRPC)
		if err != nil {
			logger.Error("failed to get latest block from RPC", "error", err)
			return 0, err
		}

		lastBN = latestBlock - 10 // Start 10 blocks behind to ensure data availability
		logger.Info("starting from block", "block", lastBN, "latest", latestBlock)
	}
	return lastBN, nil
}

func runBackfillIfConfigured(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, logger *slog.Logger) {
	logger.Info("indexer configuration", "lookback", c.Int("backfill-lookback"), "batch", c.Int("backfill-batch"))

	if c.Int("backfill-lookback") > 0 {
		logger.Info("[BACKFILL] running one-time backfill",
			"lookback", c.Int("backfill-lookback"),
			"batch", c.Int("backfill-batch"))
		if err := backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays); err != nil {
			logger.Error("[BACKFILL] failed", "error", err)
		} else {
			logger.Info("[BACKFILL] completed startup backfill")
		}
	} else {
		logger.Info("[BACKFILL] skipped", "reason", "backfill-lookback=0")
	}
}

func runMainLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, infuraRPC, beaconBase string, startBN int64, logger *slog.Logger) error {
	mainTicker := time.NewTicker(c.Duration("block-interval"))
	defer mainTicker.Stop()

	lastBN := startBN

	for {
		select {
		case <-ctx.Done():
			logger.Info("[SHUTDOWN] graceful shutdown initiated", "reason", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				logger.Error("[SHUTDOWN] failed to save last block number", "error", err)
			}
			logger.Info("[SHUTDOWN] indexer stopped", "block", lastBN)
			return nil

		case <-mainTicker.C:
			lastBN = processNextBlock(ctx, c, db, httpc, relays, infuraRPC, beaconBase, lastBN, logger)
		}
	}
}

func processNextBlock(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, infuraRPC, beaconBase string, lastBN int64, logger *slog.Logger) int64 {
	nextBN := lastBN + 1

	ei, err := beacon.FetchCombinedBlockData(ctx, httpc, infuraRPC, beaconBase, nextBN)
	if err != nil || ei == nil {
		logger.Warn("[BLOCK] not available yet", "block", nextBN, "error", err)
		return lastBN
	}

	logger.Info("processing block",
		"block", nextBN,
		"slot", ei.Slot,
		"timestamp", ei.Timestamp,
		"proposer_index", ei.ProposerIdx,
		"winning_relay", ei.RelayTag,
		"builder_pubkey_prefix", ei.BuilderHex,
		"producer_reward_eth", ei.RewardEth,
	)

	if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
		logger.Error("[DB] failed to save block", "block", nextBN, "error", err)
		return lastBN
	}
	logger.Info("[DB] block saved successfully", "block", nextBN)

	processBidsForBlock(ctx, db, httpc, relays, ei, logger)
	launchAsyncValidatorTasks(ctx, c, db, httpc, ei, beaconBase, logger)

	saveBlockProgress(db, nextBN, logger)
	return nextBN
}

func processBidsForBlock(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, ei *beacon.ExecInfo, logger *slog.Logger) {

	// Fetch and store bid data from all relays
	totalBids := 0
	successfulRelays := 0
	mainContextCanceled := false
	const batchSize = 500
	for _, rr := range relays {
		if ctx.Err() != nil {
			logger.Warn("main context canceled, stopping relay processing")
			break
		}

		bids, err := relay.FetchBuilderBlocksReceived(ctx, httpc, rr.URL, ei.Slot)
		if err != nil {
			logger.Error("[RELAY] failed to fetch bids", "relay_id", rr.ID, "url", rr.URL, "error", err)
			continue
		}

		relayBids := 0
		batch := make([]database.BidRow, 0, batchSize)

		for _, bid := range bids {
			// Check if main context is still valid
			if ctx.Err() != nil {
				logger.Warn("[BIDS] main context canceled, stopping bid insertion")
				mainContextCanceled = true
				break
			}

			if row, ok := relay.BuildBidInsert(ei.Slot, rr.ID, bid); ok {
				batch = append(batch, row)

				if len(batch) >= batchSize {
					insCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					if err := db.InsertBidsBatch(insCtx, batch); err != nil {

						logger.Error("[DB]batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
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
			flushCtx, flushCancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := db.InsertBidsBatch(flushCtx, batch); err != nil {
				logger.Error("[DB] batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
			} else {
				relayBids += len(batch)
			}
			flushCancel()
		}

		if mainContextCanceled {
			break
		}

		if relayBids > 0 {
			logger.Info("[BIDS] bids collected", "relay_id", rr.ID, "count", relayBids)
			totalBids += relayBids
			successfulRelays++
		}
	}
	logger.Info("[BIDS] summary", "block", ei.BlockNumber, "total_bids", totalBids, "successful_relays", successfulRelays)

}

func saveBlockProgress(db *database.DB, blockNum int64, logger *slog.Logger) {
	gctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.SaveLastBlockNumber(gctx, blockNum); err != nil {
		logger.Error("[PROGRESS] failed to save block number", "block", blockNum, "error", err)
	} else {
		logger.Info("[PROGRESS] advanced to block", "block", blockNum)
	}

}

func launchAsyncValidatorTasks(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, ei *beacon.ExecInfo, beaconBase string, logger *slog.Logger) { // Async validator pubkey fetch
	if ei.ProposerIdx != nil {
		go func(slot int64, proposerIdx int64) {
			time.Sleep(c.Duration("validator-delay"))

			vctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			vpub, err := beacon.FetchValidatorPubkey(vctx, httpc, beaconBase, proposerIdx)
			if err != nil {
				logger.Error("[VALIDATOR] failed to fetch pubkey", "proposer", proposerIdx, "error", err)
				return
			}

			if len(vpub) > 0 {
				if err := db.UpdateValidatorPubkey(vctx, slot, vpub); err != nil {
					logger.Error("[VALIDATOR] failed to save pubkey", "slot", slot, "error", err)
				} else {
					logger.Info("[VALIDATOR] pubkey saved", "proposer", proposerIdx, "slot", slot)
				}
			}
		}(ei.Slot, *ei.ProposerIdx)
	}

	// Async opt-in status check
	if ei.ProposerIdx != nil {
		go func(slot int64, blockNumber int64) {
			time.Sleep(c.Duration("validator-delay") + 500*time.Millisecond)

			// Wait for validator pubkey to be available
			getCtx, getCancel := context.WithTimeout(context.Background(), 5*time.Second)
			vpk, err := db.GetValidatorPubkeyWithRetry(getCtx, slot, 3, time.Second)

			if err != nil {
				logger.Error("[VALIDATOR] pubkey not available", "slot", slot, "error", err)
				return
			}

			opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, createOptionsFromCLI(c), blockNumber, vpk)
			if err != nil {
				logger.Error("[OPT-IN] failed to check opt-in status", "slot", slot, "error", err)
				return
			}

			err = db.UpdateValidatorOptInStatus(getCtx, slot, opted)
			getCancel()
			if err != nil {
				logger.Error("[OPT-IN] failed to save opt-in status", "slot", slot, "error", err)
			} else {
				logger.Info("[OPT-IN] validator opt-in status", "slot", slot, "opted_in", opted)
			}
		}(ei.Slot, ei.BlockNumber)
	}
}

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

	db, err := initializeDatabase(ctx, dbURL, initLogger)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			initLogger.Error("[DB] close failed", "error", cerr)
		}
	}()

	// Initialize HTTP client
	httpc := httputil.NewHTTPClient(c.Duration("http-timeout"))
	initLogger.Info("[HTTP] client initialized", "timeout", c.Duration("http-timeout"))

	// Load relay configurations
	relays, err := loadRelays(ctx, db, initLogger)
	if err != nil {
		return err
	}

	// Get starting block number
	lastBN, err := getStartingBlockNumber(ctx, db, httpc, infuraRPC, initLogger)
	if err != nil {
		return err
	}

	initLogger.Info("starting from block number", "block", lastBN)
	initLogger.Info("indexer configuration", "lookback", c.Int("backfill-lookback"), "batch", c.Int("backfill-batch"))

	// Run backfill if configured
	go runBackfillIfConfigured(ctx, c, db, httpc, relays, initLogger)
	return runMainLoop(ctx, c, db, httpc, relays, infuraRPC, beaconBase, lastBN, initLogger)
}
