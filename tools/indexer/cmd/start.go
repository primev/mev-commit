package main

import (
	"context"
	"log/slog"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/backfill"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"

	"github.com/urfave/cli/v2"
)

func initializeDatabase(ctx context.Context, dbURL string, logger *slog.Logger) (*database.DB, error) {
	db, err := database.Connect(ctx, dbURL, 20, 5)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		return nil, err
	}
	logger.Info("database connected to StarRocks database")

	if err := db.EnsureStateTable(ctx); err != nil {
		logger.Error("database failed to ensure state table", "error", err)
		return nil, err
	}
	logger.Info("database state table ready")

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
		logger.Info("running one-time backfill",
			"lookback", c.Int("backfill-lookback"),
			"batch", c.Int("backfill-batch"))
		if err := backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays); err != nil {
			logger.Error("failed to backfill", "error", err)
		} else {
			logger.Info("completed startup backfill")
		}
	} else {
		logger.Info("backfill skipped", "reason", "backfill-lookback=0")
	}
}

func runMainLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, infuraRPC, beaconBase string, startBN int64, logger *slog.Logger) error {
	mainTicker := time.NewTicker(c.Duration("block-interval"))
	defer mainTicker.Stop()

	lastBN := startBN

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutdown graceful shutdown initiated", "reason", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				logger.Error("shutdown failed to save last block number", "error", err)
			}
			logger.Info("shutdown indexer stopped", "block", lastBN)
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
		logger.Warn("block not available yet", "block", nextBN, "error", err)
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
		logger.Error("failed to save block", "block", nextBN, "error", err)
		return lastBN
	}
	logger.Info("block saved successfully", "block", nextBN)
	cfg := createOptionsFromCLI(c)

	if cfg.RelayData {
		if err := backfill.ProcessBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
			logger.Error("failed to process bids", "error", err)
			return lastBN
		}
	}
	if err := backfill.LaunchValidatorTasks(ctx, cfg, db, httpc, ei, beaconBase, logger); err != nil {
		logger.Error("failed to launch async validator tasks", "slot", ei.Slot, "error", err)
		return lastBN
	}

	saveBlockProgress(db, nextBN, logger)
	return nextBN
}

func saveBlockProgress(db *database.DB, blockNum int64, logger *slog.Logger) {
	gctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.SaveLastBlockNumber(gctx, blockNum); err != nil {
		logger.Error("progress failed to save block number", "block", blockNum, "error", err)
	} else {
		logger.Info("progress advanced to block", "block", blockNum)
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
		initLogger.Info("relay enabled", "count", len(relays))
	} else {
		initLogger.Info("relay disabled")
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
			initLogger.Error("database close failed", "error", cerr)
		}
	}()

	// Initialize HTTP client
	httpc := httputil.NewHTTPClient(c.Duration("http-timeout"))
	initLogger.Info("http client initialized", "timeout", c.Duration("http-timeout"))

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
