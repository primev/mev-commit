package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
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
func runBackwardLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, rpcURL, beaconBase, beaconchaAPIKey string, startBN, stopBlock int64, logger *slog.Logger) error {
	cfg := createOptionsFromCLI(c)
	currentBN := startBN

	logger.Info("starting backward indexer", "start_block", currentBN, "stop_block", stopBlock, "batch_size", cfg.BatchSize)

	for currentBN > stopBlock {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			logger.Info("shutdown initiated, saving checkpoint", "last_block", currentBN)
			if err := db.SaveBackwardCheckpoint(context.Background(), currentBN); err != nil {
				logger.Error("CRITICAL: checkpoint save failed on shutdown", "error", err)
			}
			return nil
		}

		// Calculate batch range
		batchEnd := currentBN
		batchStart := currentBN - int64(cfg.BatchSize) + 1
		if batchStart <= stopBlock {
			batchStart = stopBlock + 1
		}
		batchSize := batchEnd - batchStart + 1

		logger.Info("processing batch", "start", batchStart, "end", batchEnd, "size", batchSize)
		batchStartTime := time.Now()

		// PHASE 1: Fetch all blocks in batch sequentially (respects 30 RPS rate limit)
		fetchStartTime := time.Now()
		blocks := make([]*beacon.ExecInfo, 0, batchSize)
		for bn := batchEnd; bn >= batchStart; bn-- {
			if err := ctx.Err(); err != nil {
				logger.Info("shutdown during fetch", "last_processed", currentBN)
				db.SaveBackwardCheckpoint(context.Background(), currentBN)
				return nil
			}

			// Fetch block data
			ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, bn)
			if err != nil {
				logger.Error("FATAL: block fetch failed", "block", bn, "error", err)
				return fmt.Errorf("block fetch failed at %d: %w", bn, err)
			}
			if ei == nil {
				logger.Error("FATAL: block data is nil", "block", bn)
				return fmt.Errorf("block data nil at %d", bn)
			}

			blocks = append(blocks, ei)
		}
		fetchDuration := time.Since(fetchStartTime)

		// PHASE 2: Batch insert all blocks
		insertStartTime := time.Now()
		if err := db.InsertBlocksBatch(ctx, blocks); err != nil {
			logger.Error("FATAL: batch block insert failed", "error", err)
			return fmt.Errorf("batch block insert failed: %w", err)
		}
		insertDuration := time.Since(insertStartTime)

		// PHASE 3: Process validator tasks for each block (synchronous - must complete)
		validatorStartTime := time.Now()
		for _, ei := range blocks {
			if err := ctx.Err(); err != nil {
				logger.Info("shutdown during validator tasks", "last_processed", currentBN)
				db.SaveBackwardCheckpoint(context.Background(), currentBN)
				return nil
			}

			if err := launchValidatorTasks(ctx, c, db, httpc, beaconLimiter, ei, beaconBase, beaconchaAPIKey, logger); err != nil {
				logger.Error("FATAL: validator tasks failed", "block", ei.BlockNumber, "error", err)
				return fmt.Errorf("validator tasks failed at %d: %w", ei.BlockNumber, err)
			}

			currentBN = ei.BlockNumber - 1
		}
		validatorDuration := time.Since(validatorStartTime)

		blocksProcessed := len(blocks)

		// PHASE 4: Save checkpoint
		checkpointStartTime := time.Now()
		if err := db.SaveBackwardCheckpoint(ctx, currentBN); err != nil {
			logger.Error("FATAL: checkpoint save failed", "block", currentBN, "error", err)
			return fmt.Errorf("checkpoint save failed at %d: %w", currentBN, err)
		}
		checkpointDuration := time.Since(checkpointStartTime)

		batchDuration := time.Since(batchStartTime)
		blocksPerSecond := float64(blocksProcessed) / batchDuration.Seconds()

		logger.Info("batch completed",
			"blocks", blocksProcessed,
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"fetch_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"insert_s", fmt.Sprintf("%.3f", insertDuration.Seconds()),
			"validator_s", fmt.Sprintf("%.2f", validatorDuration.Seconds()),
			"checkpoint_s", fmt.Sprintf("%.3f", checkpointDuration.Seconds()),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSecond))
	}

	logger.Info("backward indexer completed", "stop_block", stopBlock)
	return nil
}

// runForwardLoop indexes blocks forward from the starting point
func runForwardLoop(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, rpcURL, beaconBase, beaconchaAPIKey string, startBN int64, logger *slog.Logger) error {
	cfg := createOptionsFromCLI(c)
	currentBN := startBN

	logger.Info("starting forward indexer", "start_block", currentBN, "batch_size", cfg.BatchSize)

	for {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			logger.Info("shutdown initiated, saving checkpoint", "last_block", currentBN)
			if err := db.SaveForwardCheckpoint(context.Background(), currentBN); err != nil {
				logger.Error("CRITICAL: checkpoint save failed on shutdown", "error", err)
			}
			return nil
		}

		// Get latest block from chain
		latestBlock, err := ethereum.GetLatestBlockNumber(httpc.HTTPClient, rpcURL)
		if err != nil {
			logger.Error("FATAL: failed to get latest block", "error", err)
			return fmt.Errorf("get latest block failed: %w", err)
		}

		// Calculate batch range
		batchStart := currentBN + 1
		batchEnd := batchStart + int64(cfg.BatchSize) - 1
		if batchEnd > latestBlock {
			batchEnd = latestBlock
		}

		// Check if caught up
		if batchStart > latestBlock {
			logger.Debug("caught up, waiting for new blocks", "current", currentBN, "latest", latestBlock)
			time.Sleep(12 * time.Second)
			continue
		}

		batchSize := batchEnd - batchStart + 1
		behindTip := latestBlock - batchEnd

		logger.Info("processing batch", "start", batchStart, "end", batchEnd, "size", batchSize, "behind_tip", behindTip)
		batchStartTime := time.Now()

		// PHASE 1: Fetch all blocks in batch sequentially (respects 30 RPS rate limit)
		fetchStartTime := time.Now()
		blocks := make([]*beacon.ExecInfo, 0, batchSize)
		for bn := batchStart; bn <= batchEnd; bn++ {
			if err := ctx.Err(); err != nil {
				logger.Info("shutdown during fetch", "last_processed", currentBN)
				db.SaveForwardCheckpoint(context.Background(), currentBN)
				return nil
			}

			// Fetch block data
			ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, bn)
			if err != nil {
				logger.Error("FATAL: block fetch failed", "block", bn, "error", err)
				return fmt.Errorf("block fetch failed at %d: %w", bn, err)
			}
			if ei == nil {
				logger.Error("FATAL: block data is nil", "block", bn)
				return fmt.Errorf("block data nil at %d", bn)
			}

			blocks = append(blocks, ei)
		}
		fetchDuration := time.Since(fetchStartTime)

		// PHASE 2: Batch insert all blocks
		insertStartTime := time.Now()
		if err := db.InsertBlocksBatch(ctx, blocks); err != nil {
			logger.Error("FATAL: batch block insert failed", "error", err)
			return fmt.Errorf("batch block insert failed: %w", err)
		}
		insertDuration := time.Since(insertStartTime)

		// PHASE 3: Process validator tasks for each block (synchronous - must complete)
		validatorStartTime := time.Now()
		for _, ei := range blocks {
			if err := ctx.Err(); err != nil {
				logger.Info("shutdown during validator tasks", "last_processed", currentBN)
				db.SaveForwardCheckpoint(context.Background(), currentBN)
				return nil
			}

			if err := launchValidatorTasks(ctx, c, db, httpc, beaconLimiter, ei, beaconBase, beaconchaAPIKey, logger); err != nil {
				logger.Error("FATAL: validator tasks failed", "block", ei.BlockNumber, "error", err)
				return fmt.Errorf("validator tasks failed at %d: %w", ei.BlockNumber, err)
			}

			currentBN = ei.BlockNumber
		}
		validatorDuration := time.Since(validatorStartTime)

		blocksProcessed := len(blocks)

		// PHASE 4: Save checkpoint
		checkpointStartTime := time.Now()
		if err := db.SaveForwardCheckpoint(ctx, currentBN); err != nil {
			logger.Error("FATAL: checkpoint save failed", "block", currentBN, "error", err)
			return fmt.Errorf("checkpoint save failed at %d: %w", currentBN, err)
		}
		checkpointDuration := time.Since(checkpointStartTime)

		batchDuration := time.Since(batchStartTime)
		blocksPerSecond := float64(blocksProcessed) / batchDuration.Seconds()

		logger.Info("batch completed",
			"blocks", blocksProcessed,
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"fetch_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"insert_s", fmt.Sprintf("%.3f", insertDuration.Seconds()),
			"validator_s", fmt.Sprintf("%.2f", validatorDuration.Seconds()),
			"checkpoint_s", fmt.Sprintf("%.3f", checkpointDuration.Seconds()),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSecond))
	}
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

// runBidWorker continuously queries for blocks without bids and processes them
// direction: "forward" (ascending) or "backward" (descending)
func runBidWorker(ctx context.Context, db *database.DB, httpc *retryablehttp.Client, relays []relay.Row, startBlock int64, direction string, logger *slog.Logger) error {
	// Load checkpoint
	currentBlock := startBlock
	if direction == "forward" {
		if checkpoint, found := db.LoadForwardBidCheckpoint(ctx); found && checkpoint > 0 {
			currentBlock = checkpoint
			logger.Info("resuming from checkpoint", "block", currentBlock)
		}
	} else {
		if checkpoint, found := db.LoadBackwardBidCheckpoint(ctx); found && checkpoint > 0 {
			currentBlock = checkpoint
			logger.Info("resuming from checkpoint", "block", currentBlock)
		}
	}

	logger.Info("bid worker started", "start_block", currentBlock, "direction", direction, "relays", len(relays))
	batchSize := int64(100)

	for {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			logger.Info("shutdown initiated, saving checkpoint", "last_block", currentBlock)
			if direction == "forward" {
				db.SaveForwardBidCheckpoint(context.Background(), currentBlock)
			} else {
				db.SaveBackwardBidCheckpoint(context.Background(), currentBlock)
			}
			return nil
		}

		// PHASE 1: Get blocks without bids
		queryStartTime := time.Now()
		blocks, err := db.GetBlocksWithoutBids(ctx, currentBlock, batchSize, direction)
		queryDuration := time.Since(queryStartTime)

		if err != nil {
			logger.Error("FATAL: failed to query blocks without bids", "error", err)
			return fmt.Errorf("query blocks without bids failed: %w", err)
		}

		if len(blocks) == 0 {
			logger.Debug("no blocks need bid processing, waiting")
			time.Sleep(10 * time.Second)
			continue
		}

		logger.Info("processing blocks for bids", "count", len(blocks), "query_ms", queryDuration.Milliseconds())
		batchStartTime := time.Now()

		// Process each block
		for _, block := range blocks {
			if err := ctx.Err(); err != nil {
				logger.Info("shutdown during processing", "last_block", currentBlock)
				if direction == "forward" {
					db.SaveForwardBidCheckpoint(context.Background(), currentBlock)
				} else {
					db.SaveBackwardBidCheckpoint(context.Background(), currentBlock)
				}
				return nil
			}

			// PHASE 2: Fetch bids from all relays in parallel
			blockStartTime := time.Now()
			type relayResult struct {
				relayID int64
				bids    int
				err     error
			}
			resultsChan := make(chan relayResult, len(relays))
			var wg sync.WaitGroup

			for _, rr := range relays {
				wg.Add(1)
				go func(r relay.Row) {
					defer wg.Done()

					relayCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()

					bids, err := relay.FetchBuilderBlocksReceived(relayCtx, httpc, r.URL, block.Slot)
					if err != nil {
						resultsChan <- relayResult{relayID: r.ID, err: err}
						return
					}

					// Insert bids
					const batchInsertSize = 500
					batch := make([]database.BidRow, 0, batchInsertSize)
					totalInserted := 0

					for _, bid := range bids {
						if row, ok := relay.BuildBidInsert(block.Slot, r.ID, bid); ok {
							batch = append(batch, row)
							if len(batch) >= batchInsertSize {
								if err := db.InsertBidsBatch(ctx, batch); err != nil {
									resultsChan <- relayResult{relayID: r.ID, err: fmt.Errorf("batch insert: %w", err)}
									return
								}
								totalInserted += len(batch)
								batch = batch[:0]
							}
						}
					}

					// Insert remaining
					if len(batch) > 0 {
						if err := db.InsertBidsBatch(ctx, batch); err != nil {
							resultsChan <- relayResult{relayID: r.ID, err: fmt.Errorf("final batch insert: %w", err)}
							return
						}
						totalInserted += len(batch)
					}

					resultsChan <- relayResult{relayID: r.ID, bids: totalInserted}
				}(rr)
			}

			// Wait for all relays
			go func() {
				wg.Wait()
				close(resultsChan)
			}()

			// Collect results
			totalBids := 0
			successfulRelays := 0
			var errors []string

			for result := range resultsChan {
				if result.err != nil {
					errors = append(errors, fmt.Sprintf("relay_%d: %v", result.relayID, result.err))
				} else if result.bids > 0 {
					totalBids += result.bids
					successfulRelays++
				}
			}

			// If all relays failed, exit
			if len(errors) == len(relays) {
				logger.Error("FATAL: all relays failed for block", "block", block.BlockNumber, "errors", errors)
				return fmt.Errorf("all relays failed for block %d", block.BlockNumber)
			}

			blockDuration := time.Since(blockStartTime)

			logger.Info("block bids processed",
				"block", block.BlockNumber,
				"slot", block.Slot,
				"total_bids", totalBids,
				"successful_relays", successfulRelays,
				"total_relays", len(relays),
				"duration_ms", blockDuration.Milliseconds())

			// PHASE 3: Update checkpoint
			checkpointStartTime := time.Now()
			currentBlock = block.BlockNumber
			var checkpointErr error
			if direction == "forward" {
				checkpointErr = db.SaveForwardBidCheckpoint(ctx, currentBlock)
			} else {
				checkpointErr = db.SaveBackwardBidCheckpoint(ctx, currentBlock)
			}
			checkpointDuration := time.Since(checkpointStartTime)

			if checkpointErr != nil {
				logger.Error("FATAL: checkpoint save failed", "block", currentBlock, "error", checkpointErr)
				return fmt.Errorf("checkpoint save failed at %d: %w", currentBlock, checkpointErr)
			}

			logger.Debug("checkpoint saved", "block", currentBlock, "duration_ms", checkpointDuration.Milliseconds())
		}

		batchDuration := time.Since(batchStartTime)
		blocksPerSecond := float64(len(blocks)) / batchDuration.Seconds()
		logger.Info("bid batch completed",
			"blocks_processed", len(blocks),
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSecond))
	}
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

	// Create cancellable context from CLI context
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	// Load relay configurations
	cfgRelays, err := config.ResolveRelays(c)
	if err != nil {
		return err
	}
	relays := make([]relay.Row, 0, len(cfgRelays))
	for _, r := range cfgRelays {
		relays = append(relays, relay.Row{ID: r.Relay_id, URL: r.URL})
	}
	initLogger.Info("relays loaded", "count", len(relays))
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

	// Get starting points for forward and backward indexers
	forwardStart, backwardStart, err := getStartingPoints(ctx, db, httpc, rpcURL, initLogger)
	if err != nil {
		return err
	}

	// Check if bid workers should be enabled (based on whether relays are configured)
	enableBidWorkers := len(relays) > 0
	totalWorkers := 2 // Always run 2 block indexers
	if enableBidWorkers {
		totalWorkers += 2 // Add 2 bid workers if relays configured
	}

	if enableBidWorkers {
		initLogger.Info("[INIT] bid workers enabled - relays configured", "relays", len(relays))
	} else {
		initLogger.Info("[INIT] bid workers disabled - no relays configured")
	}

	// Get max block for backward bid worker (only if needed)
	var maxBlock int64
	if enableBidWorkers {
		maxBlock, err = db.GetMaxBlockNumber(ctx)
		if err != nil {
			initLogger.Warn("[BID_WORKER] failed to get max block, using forward start", "error", err)
			maxBlock = forwardStart
		}
	}

	initLogger.Info("[INIT] starting workers",
		"total_workers", totalWorkers,
		"forward_block_start", forwardStart,
		"backward_block_start", backwardStart,
		"block_interval", c.Duration("block-interval"))

	// Create error channel to capture all goroutine errors
	errChan := make(chan error, totalWorkers)

	// Launch forward block indexer
	go func() {
		logger := initLogger.With("worker", "block-forward")
		logger.Info("launching forward block indexer")
		err := runForwardLoop(ctx, c, db, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, forwardStart, logger)
		errChan <- err
	}()

	// Launch backward block indexer
	go func() {
		logger := initLogger.With("worker", "block-backward")
		logger.Info("launching backward block indexer")
		err := runBackwardLoop(ctx, c, db, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, backwardStart, backwardStopBlock, logger)
		errChan <- err
	}()

	// Conditionally launch bid workers
	if enableBidWorkers {
		// Launch forward bid worker
		go func() {
			logger := initLogger.With("worker", "bid-forward")
			logger.Info("launching forward bid worker")
			err := runBidWorker(ctx, db, httpc, relays, forwardStart, "forward", logger)
			errChan <- err
		}()

		// Launch backward bid worker
		go func() {
			logger := initLogger.With("worker", "bid-backward")
			logger.Info("launching backward bid worker")
			err := runBidWorker(ctx, db, httpc, relays, maxBlock, "backward", logger)
			errChan <- err
		}()
	}

	// Wait for first worker to exit (error or completion)
	// If any worker exits with error, the whole application should exit
	err = <-errChan
	if err != nil {
		initLogger.Error("FATAL: worker failed, shutting down all workers", "error", err)
		cancel() // Cancel context to stop other workers
		return err
	}

	// If one worker completed successfully, wait for others
	initLogger.Info("one worker completed, waiting for others")
	for i := 1; i < totalWorkers; i++ {
		workerErr := <-errChan
		if workerErr != nil && err == nil {
			err = workerErr
		}
	}

	if err != nil {
		initLogger.Error("FATAL: worker(s) failed", "error", err)
		return err
	}

	initLogger.Info("all workers completed successfully")
	return nil
}
