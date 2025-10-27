package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"golang.org/x/time/rate"

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

	if err := db.EnsureRelaysTable(ctx); err != nil {
		logger.Error("[DB] failed to ensure relays table", "error", err)
		return nil, err
	}
	logger.Info("[DB] relays table ready")

	if err := db.EnsureBlocksTable(ctx); err != nil {
		logger.Error("[DB] failed to ensure blocks table", "error", err)
		return nil, err
	}
	logger.Info("[DB] blocks table ready")

	if err := db.EnsureBidsTable(ctx); err != nil {
		logger.Error("[DB] failed to ensure bids table", "error", err)
		return nil, err
	}
	logger.Info("[DB] bids table ready")

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

// getStartingPoints returns the forward and backward starting block numbers
func getStartingPoints(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, rpcURL string, logger *slog.Logger) (forwardStart, backwardStart int64, err error) {
	// Try to load existing checkpoints
	forwardBN, forwardFound := db.LoadForwardCheckpoint(ctx)
	backwardBN, backwardFound := db.LoadBackwardCheckpoint(ctx)

	// If both checkpoints exist, resume from them
	if forwardFound && backwardFound && forwardBN > 0 && backwardBN > 0 {
		logger.Info("[CHECKPOINT] resuming from saved checkpoints",
			"forward_block", forwardBN,
			"backward_block", backwardBN)
		return forwardBN, backwardBN, nil
	}

	// Otherwise, get the current latest block from Ethereum
	logger.Info("[INIT] no valid checkpoints found, getting latest block from Ethereum RPC...")
	latestBlock, err := ethereum.GetLatestBlockNumber(httpc.HTTPClient, rpcURL)
	if err != nil {
		logger.Error("[INIT] failed to get latest block from RPC", "error", err)
		return 0, 0, err
	}

	// Start both indexers from the current latest block
	// Forward will go: latest → latest+1 → latest+2...
	// Backward will go: latest-1 → latest-2 → latest-3...
	startBlock := latestBlock
	logger.Info("[INIT] initializing both indexers from current block",
		"start_block", startBlock,
		"latest_block", latestBlock)

	// Save initial checkpoints
	if err := db.SaveForwardCheckpoint(ctx, startBlock); err != nil {
		logger.Warn("[CHECKPOINT] failed to save initial forward checkpoint", "error", err)
	}
	if err := db.SaveBackwardCheckpoint(ctx, startBlock-1); err != nil {
		logger.Warn("[CHECKPOINT] failed to save initial backward checkpoint", "error", err)
	}

	return startBlock, startBlock - 1, nil
}

// runBackwardLoop indexes blocks backwards from the starting point
func runBackwardLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, relays []relay.Row, rpcURL, beaconBase, beaconchaAPIKey string, startBN, stopBlock int64, logger *slog.Logger) error {
	logger = logger.With("indexer", "backward")
	blockInterval := c.Duration("block-interval")
	ticker := time.NewTicker(blockInterval)
	defer ticker.Stop()

	currentBN := startBN
	logger.Info("[BACKWARD] starting backward indexer", "start_block", currentBN, "stop_block", stopBlock)

	for {
		select {
		case <-ctx.Done():
			logger.Info("[BACKWARD] shutdown initiated", "last_block", currentBN)
			if err := db.SaveBackwardCheckpoint(ctx, currentBN); err != nil {
				logger.Error("[BACKWARD] failed to save checkpoint on shutdown", "error", err)
			}
			return nil

		case <-ticker.C:
			if currentBN <= stopBlock {
				logger.Info("[BACKWARD] reached stop block, stopping", "stop_block", stopBlock)
				return nil
			}

			logger.Info("[BACKWARD] processing block", "block", currentBN)

			ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, currentBN)
			if err != nil || ei == nil {
				logger.Warn("[BACKWARD] block not available", "block", currentBN, "error", err)
				continue
			}

			// Process block data
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				logger.Error("[BACKWARD] failed to upsert block", "block", currentBN, "error", err)
				continue
			}

			// Process bids for this block
			if err := processBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
				logger.Error("[BACKWARD] failed to process bids", "block", currentBN, "error", err)
			}

			// Process validator tasks
			if err := launchValidatorTasks(ctx, c, db, httpc, beaconLimiter, ei, beaconBase, beaconchaAPIKey, logger); err != nil {
				logger.Error("[BACKWARD] failed to launch validator tasks", "block", currentBN, "error", err)
			}

			// Save checkpoint
			if err := db.SaveBackwardCheckpoint(ctx, currentBN); err != nil {
				logger.Error("[BACKWARD] failed to save checkpoint", "block", currentBN, "error", err)
			}

			// Move to previous block
			currentBN--
		}
	}
}

// runForwardLoop indexes blocks forward from the starting point
func runForwardLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, relays []relay.Row, rpcURL, beaconBase, beaconchaAPIKey string, startBN int64, logger *slog.Logger) error {
	logger = logger.With("indexer", "forward")
	blockInterval := c.Duration("block-interval")
	ticker := time.NewTicker(blockInterval)
	defer ticker.Stop()

	currentBN := startBN
	logger.Info("[FORWARD] starting forward indexer", "start_block", currentBN)

	for {
		select {
		case <-ctx.Done():
			logger.Info("[FORWARD] shutdown initiated", "last_block", currentBN)
			if err := db.SaveForwardCheckpoint(ctx, currentBN); err != nil {
				logger.Error("[FORWARD] failed to save checkpoint on shutdown", "error", err)
			}
			return nil

		case <-ticker.C:
			nextBN := currentBN + 1
			logger.Info("[FORWARD] processing block", "block", nextBN)

			ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, nextBN)
			if err != nil || ei == nil {
				logger.Warn("[FORWARD] block not available yet", "block", nextBN, "error", err)
				continue
			}

			// Process block data
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				logger.Error("[FORWARD] failed to upsert block", "block", nextBN, "error", err)
				continue
			}

			// Process bids for this block
			if err := processBidsForBlock(ctx, db, httpc, relays, ei, logger); err != nil {
				logger.Error("[FORWARD] failed to process bids", "block", nextBN, "error", err)
			}

			// Process validator tasks
			if err := launchValidatorTasks(ctx, c, db, httpc, beaconLimiter, ei, beaconBase, beaconchaAPIKey, logger); err != nil {
				logger.Error("[FORWARD] failed to launch validator tasks", "block", nextBN, "error", err)
			}

			// Save checkpoint and advance
			if err := db.SaveForwardCheckpoint(ctx, nextBN); err != nil {
				logger.Error("[FORWARD] failed to save checkpoint", "block", nextBN, "error", err)
			}
			currentBN = nextBN
		}
	}
}
func processBidsForBlock(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, ei *beacon.ExecInfo, logger *slog.Logger) error {

	// Fetch and store bid data from all relays
	totalBids := 0
	successfulRelays := 0
	const batchSize = 500
	for _, rr := range relays {
		if err := ctx.Err(); err != nil {
			logger.Warn("main context canceled, stopping relay processing")
			return err
		}

		bids, err := relay.FetchBuilderBlocksReceived(ctx, httpc, rr.URL, ei.Slot)
		if err != nil {
			// logger.Error("[RELAY] failed to fetch bids", "relay_id", rr.ID, "url", rr.URL, "error", err)
			return fmt.Errorf("fetch bids: relay_id=%d url=%s slot=%d: %w", rr.ID, rr.URL, ei.Slot, err)

		}

		relayBids := 0
		batch := make([]database.BidRow, 0, batchSize)

		for _, bid := range bids {

			if err := ctx.Err(); err != nil {
				logger.Warn("[BIDS] main context canceled, stopping bid insertion")
				return err
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

		if relayBids > 0 {
			logger.Info("[BIDS] bids collected", "relay_id", rr.ID, "count", relayBids)
			totalBids += relayBids
			successfulRelays++
		}
	}
	logger.Info("[BIDS] summary", "block", ei.BlockNumber, "total_bids", totalBids, "successful_relays", successfulRelays)
	return nil
}

func launchValidatorTasks(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, ei *beacon.ExecInfo, beaconBase, beaconchaAPIKey string, logger *slog.Logger) error { // Async validator pubkey fetch
	if ei.ProposerIdx == nil {
		return nil
	}

	vctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	vpub, err := beacon.FetchValidatorPubkey(vctx, httpc, beaconLimiter, beaconBase, beaconchaAPIKey, *ei.ProposerIdx)
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
	rpcURL := c.String(optionRPCURL.Name)
	beaconBase := c.String(optionBeaconBase.Name)
	beaconchaAPIKey := c.String(optionBeaconchaAPIKey.Name)
	beaconchaRPS := c.Int(optionBeaconchaRPS.Name)
	backwardStopBlock := c.Int64(optionBackwardStopBlock.Name)

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

	// Initialize rate limiter for beaconcha API
	var beaconLimiter *rate.Limiter
	if beaconchaRPS > 0 {
		beaconLimiter = rate.NewLimiter(rate.Limit(beaconchaRPS), beaconchaRPS)
		initLogger.Info("[RATE_LIMITER] beaconcha rate limiter initialized", "rps", beaconchaRPS)
	} else {
		initLogger.Warn("[RATE_LIMITER] beaconcha rate limiting disabled (rps=0)")
	}

	// Load relay configurations
	relays, err := loadRelays(ctx, db, initLogger)
	if err != nil {
		return err
	}

	// Get starting points for forward and backward indexers
	forwardStart, backwardStart, err := getStartingPoints(ctx, db, httpc, rpcURL, initLogger)
	if err != nil {
		return err
	}

	initLogger.Info("[INIT] dual indexer starting",
		"forward_start", forwardStart,
		"backward_start", backwardStart,
		"block_interval", c.Duration("block-interval"))

	// Create error channel to capture indexer errors
	errChan := make(chan error, 2)

	// Launch forward indexer
	go func() {
		initLogger.Info("[FORWARD] launching forward indexer goroutine")
		err := runForwardLoop(ctx, c, db, httpc, beaconLimiter, relays, rpcURL, beaconBase, beaconchaAPIKey, forwardStart, initLogger)
		errChan <- err
	}()

	// Launch backward indexer
	go func() {
		initLogger.Info("[BACKWARD] launching backward indexer goroutine")
		err := runBackwardLoop(ctx, c, db, httpc, beaconLimiter, relays, rpcURL, beaconBase, beaconchaAPIKey, backwardStart, backwardStopBlock, initLogger)
		errChan <- err
	}()

	// Wait for either indexer to complete or error
	err1 := <-errChan
	err2 := <-errChan

	if err1 != nil {
		initLogger.Error("[INDEXER] indexer 1 stopped with error", "error", err1)
	}
	if err2 != nil {
		initLogger.Error("[INDEXER] indexer 2 stopped with error", "error", err2)
	}

	initLogger.Info("[INDEXER] both indexers stopped")
	return nil
}
