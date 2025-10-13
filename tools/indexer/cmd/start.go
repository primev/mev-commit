package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/backfill"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"log/slog"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

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
func safe(p interface{}) interface{} {
	v := reflect.ValueOf(p)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
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
		"proposer_index", safe(ei.ProposerIdx),
		"winning_relay", safe(ei.RelayTag),
		"builder_pubkey_prefix", safe(ei.BuilderPublicKey),
		"mev_reward_eth", safe(ei.MevRewardEth),
	)

	if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
		logger.Error("[DB] failed to save block", "block", nextBN, "error", err)
		return lastBN
	}
	logger.Info("[DB] block saved successfully", "block", nextBN)
	cfg := createOptionsFromCLI(c)

	if cfg.RelayMode {
		if err := processBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
			logger.Error("failed to process bids", "error", err)
			return lastBN
		}
	}
	if err := launchValidatorTasks(ctx, c, db, httpc, ei, beaconBase, logger); err != nil {
		logger.Error("[VALIDATOR] failed to launch async tasks", "slot", ei.Slot, "error", err)
		return lastBN
	}

	saveBlockProgress(db, nextBN, logger)
	return nextBN
}

func processBidsForBlock(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, ei *beacon.ExecInfo, logger *slog.Logger) error {
	logger.Info("[BIDS] processing bids for block", "block", ei.BlockNumber, "slot", ei.Slot)
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
			logger.Info("[BIDS] ok",
				"slot", ei.Slot, "relay_id", rel.ID,
				"bids_in", len(bids), "rows_out", len(rows),
			)
		}(r)
	}
	wg.Wait()

	logger.Info("[BIDS] summary", "block", ei.BlockNumber, "total_bids", totalBids, "successful_relays", successfulRelays)
	return nil
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

func launchValidatorTasks(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, ei *beacon.ExecInfo, beaconBase string, logger *slog.Logger) error { // Async validator pubkey fetch
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
			logger.Error("[VALIDATOR] failed to save pubkey", "slot", ei.Slot, "error", err)
		} else {
			logger.Info("[VALIDATOR] pubkey saved", "proposer", *ei.ProposerIdx, "slot", ei.Slot)
		}
	}

	// Wait for validator pubkey to be available
	getCtx, getCancel := context.WithTimeout(context.Background(), 5*time.Second)
	vpk, err := db.GetValidatorPubkeyWithRetry(getCtx, ei.Slot, 3, time.Second)
	getCancel()

	if err != nil {
		logger.Error("[VALIDATOR] pubkey not available", "slot", ei.Slot, "error", err)
		return fmt.Errorf("save validator pubkey: %w", err)
	}

	opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, createOptionsFromCLI(c), ei.BlockNumber, vpk)
	if err != nil {
		return fmt.Errorf("check opt-in status: %w", err)
	}

	updCtx, updCancel := context.WithTimeout(context.Background(), 3*time.Second)
	err = db.UpdateValidatorOptInStatus(updCtx, ei.Slot, opted)
	updCancel()
	if err != nil {
		return fmt.Errorf("save opt-in status: %w", err)
	} else {
		logger.Info("[OPT-IN] validator opt-in status", "slot", ei.Slot, "opted_in", opted)
	}
	return nil

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
	var relays []relay.Row
	if c.Bool(optionRelayFlag.Name) {
		cfgRelays, err := config.ResolveRelays(c)
		if err != nil {
			return err
		}
		relays = make([]relay.Row, 0, len(cfgRelays))
		for _, r := range cfgRelays {
			relays = append(relays, relay.Row{ID: r.Relay_id, URL: r.URL})
		}
		initLogger.Info("[RELAY] enabled", "count", len(relays))
	} else {
		initLogger.Info("[RELAY] disabled")
		relays = make([]relay.Row, 0, len(config.RelaysDefault))
		for _, r := range config.RelaysDefault {
			relays = append(relays, relay.Row{ID: r.Relay_id, URL: r.URL})
		}
	}
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
