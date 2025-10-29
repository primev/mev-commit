package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

// fillGaps detects and fills missing blocks in the database
// This is a one-time operation at startup to fix gaps from previous issues
func fillGaps(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, rpcURL, beaconBase, beaconchaAPIKey string, logger *slog.Logger) error {
	logger.Info("[FILL_GAPS] detecting missing blocks...")

	// Query for gaps
	gaps, err := db.QueryGaps(ctx)
	if err != nil {
		return fmt.Errorf("failed to query gaps: %w", err)
	}

	// Convert gaps to list of missing blocks
	var missingBlocks []int64
	for _, gap := range gaps {
		gapStart, gapEnd := gap[0], gap[1]
		for bn := gapStart; bn <= gapEnd; bn++ {
			missingBlocks = append(missingBlocks, bn)
		}
	}

	if len(missingBlocks) == 0 {
		logger.Info("[FILL_GAPS] no missing blocks found")
		return nil
	}

	logger.Info("[FILL_GAPS] found missing blocks", "count", len(missingBlocks), "blocks", missingBlocks)

	// Fetch missing blocks
	blocks := make([]*beacon.ExecInfo, 0, len(missingBlocks))
	for _, bn := range missingBlocks {
		logger.Info("[FILL_GAPS] fetching block", "block", bn)
		ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, bn)
		if err != nil {
			logger.Error("[FILL_GAPS] failed to fetch block", "block", bn, "error", err)
			return fmt.Errorf("failed to fetch block %d: %w", bn, err)
		}
		blocks = append(blocks, ei)
	}

	// Insert blocks
	logger.Info("[FILL_GAPS] inserting missing blocks", "count", len(blocks))
	if err := db.InsertBlocksBatch(ctx, blocks); err != nil {
		return fmt.Errorf("failed to insert blocks: %w", err)
	}

	// Process validators
	logger.Info("[FILL_GAPS] processing validators for missing blocks")
	if err := processValidatorsBatch(ctx, c, db, httpc, beaconLimiter, blocks, beaconBase, beaconchaAPIKey, logger); err != nil {
		return fmt.Errorf("failed to process validators: %w", err)
	}

	logger.Info("[FILL_GAPS] successfully filled all gaps", "blocks_filled", len(missingBlocks))
	return nil
}

// fetchBlockWithRetry attempts to fetch a block with exponential backoff retries
// Only used for forward indexer where blocks might not be available yet on beaconcha.in
func fetchBlockWithRetry(ctx context.Context, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, rpcURL, beaconBase, beaconchaAPIKey string, blockNum int64, logger *slog.Logger) (*beacon.ExecInfo, error) {
	retryDelays := []time.Duration{
		2 * time.Second,
		4 * time.Second,
		12 * time.Second,
		20 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	var lastErr error
	for attempt := 0; attempt <= len(retryDelays); attempt++ {
		// First attempt (attempt=0) has no delay
		if attempt > 0 {
			delay := retryDelays[attempt-1]
			logger.Warn("retrying block fetch",
				"block", blockNum,
				"attempt", attempt,
				"delay_s", delay.Seconds(),
				"max_attempts", len(retryDelays)+1)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		ei, err := beacon.FetchCombinedBlockData(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, blockNum)
		if err == nil && ei != nil {
			if attempt > 0 {
				logger.Info("block fetch succeeded after retry",
					"block", blockNum,
					"attempt", attempt+1)
			}
			return ei, nil
		}

		lastErr = err
		logger.Debug("block fetch attempt failed",
			"block", blockNum,
			"attempt", attempt+1,
			"error", err)
	}

	// All retries exhausted
	return nil, fmt.Errorf("block fetch failed after %d attempts: %w", len(retryDelays)+1, lastErr)
}

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

		// PHASE 1: Fetch all blocks in batch with parallel workers
		// Workers share the rate limiter (30 RPS), so more workers = faster fetching
		fetchStartTime := time.Now()

		// Parallel fetch using worker pool
		type fetchResult struct {
			blockNum int64
			block    *beacon.ExecInfo
			duration time.Duration
			err      error
		}

		numWorkers := cfg.FetchWorkers
		if numWorkers < 1 {
			numWorkers = 1
		}

		blockNumsChan := make(chan int64, batchSize)
		resultsChan := make(chan fetchResult, batchSize)

		// Launch workers
		var fetchWg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			fetchWg.Add(1)
			go func() {
				defer fetchWg.Done()
				for bn := range blockNumsChan {
					if err := ctx.Err(); err != nil {
						resultsChan <- fetchResult{blockNum: bn, err: err}
						continue
					}

					blockFetchStart := time.Now()
					ei, err := fetchBlockWithRetry(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, bn, logger)
					blockFetchDuration := time.Since(blockFetchStart)

					resultsChan <- fetchResult{
						blockNum: bn,
						block:    ei,
						duration: blockFetchDuration,
						err:      err,
					}
				}
			}()
		}

		// Send block numbers to workers
		go func() {
			for bn := batchEnd; bn >= batchStart; bn-- {
				blockNumsChan <- bn
			}
			close(blockNumsChan)
		}()

		// Wait for all workers to finish
		go func() {
			fetchWg.Wait()
			close(resultsChan)
		}()

		// Collect results and order them
		resultsMap := make(map[int64]fetchResult)
		var totalFetchTime time.Duration
		fetchCount := 0

		for result := range resultsChan {
			if result.err != nil {
				logger.Error("FATAL: block fetch failed after all retries", "block", result.blockNum, "error", result.err)
				return fmt.Errorf("block fetch failed at %d: %w", result.blockNum, result.err)
			}
			if result.block == nil {
				logger.Error("FATAL: block data is nil", "block", result.blockNum)
				return fmt.Errorf("block data nil at %d", result.blockNum)
			}

			resultsMap[result.blockNum] = result
			totalFetchTime += result.duration
			fetchCount++

			// Log progress every 20 blocks
			if fetchCount%20 == 0 {
				avgPerBlock := totalFetchTime / time.Duration(fetchCount)
				logger.Info("fetch progress",
					"blocks_fetched", fetchCount,
					"workers", numWorkers,
					"avg_per_block_ms", avgPerBlock.Milliseconds(),
					"total_elapsed_s", fmt.Sprintf("%.2f", time.Since(fetchStartTime).Seconds()))
			}
		}

		// Build blocks slice in descending order
		blocks := make([]*beacon.ExecInfo, 0, batchSize)
		for bn := batchEnd; bn >= batchStart; bn-- {
			if result, found := resultsMap[bn]; found {
				blocks = append(blocks, result.block)
			} else {
				// This should never happen - all blocks should have been fetched
				logger.Error("FATAL: block missing from results map", "block", bn)
				return fmt.Errorf("block %d missing from results map", bn)
			}
		}

		// Sanity check: verify all blocks were fetched
		if len(blocks) != fetchCount {
			logger.Error("FATAL: block count mismatch", "expected", fetchCount, "got", len(blocks))
			return fmt.Errorf("block count mismatch: expected %d, got %d", fetchCount, len(blocks))
		}

		fetchDuration := time.Since(fetchStartTime)
		avgFetchPerBlock := time.Duration(0)
		if fetchCount > 0 {
			avgFetchPerBlock = totalFetchTime / time.Duration(fetchCount)
		}

		logger.Info("fetch phase complete",
			"blocks", fetchCount,
			"workers", numWorkers,
			"total_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"avg_per_block_ms", avgFetchPerBlock.Milliseconds(),
			"effective_rps", fmt.Sprintf("%.2f", float64(fetchCount)/fetchDuration.Seconds()))

		// PHASE 2: Batch insert all blocks
		insertStartTime := time.Now()
		if err := db.InsertBlocksBatch(ctx, blocks); err != nil {
			logger.Error("FATAL: batch block insert failed", "error", err)
			return fmt.Errorf("batch block insert failed: %w", err)
		}
		insertDuration := time.Since(insertStartTime)

		// PHASE 3: Process validator tasks in batch
		validatorStartTime := time.Now()
		if err := processValidatorsBatch(ctx, c, db, httpc, beaconLimiter, blocks, beaconBase, beaconchaAPIKey, logger); err != nil {
			logger.Error("FATAL: validator batch processing failed", "error", err)
			return fmt.Errorf("validator batch processing failed: %w", err)
		}
		validatorDuration := time.Since(validatorStartTime)

		// Update current block number
		currentBN = blocks[len(blocks)-1].BlockNumber - 1

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
		fetchRPS := float64(blocksProcessed) / fetchDuration.Seconds()
		validatorRPS := float64(blocksProcessed) / validatorDuration.Seconds()

		logger.Info("batch completed",
			"blocks", blocksProcessed,
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"fetch_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"fetch_rps", fmt.Sprintf("%.2f", fetchRPS),
			"insert_s", fmt.Sprintf("%.3f", insertDuration.Seconds()),
			"validator_s", fmt.Sprintf("%.2f", validatorDuration.Seconds()),
			"validator_rps", fmt.Sprintf("%.2f", validatorRPS),
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

		// PHASE 1: Fetch all blocks in batch with parallel workers
		// Workers share the rate limiter (30 RPS), so more workers = faster fetching
		fetchStartTime := time.Now()

		// Parallel fetch using worker pool
		type fetchResult struct {
			blockNum int64
			block    *beacon.ExecInfo
			duration time.Duration
			err      error
		}

		numWorkers := cfg.FetchWorkers
		if numWorkers < 1 {
			numWorkers = 1
		}

		blockNumsChan := make(chan int64, batchSize)
		resultsChan := make(chan fetchResult, batchSize)

		// Launch workers
		var fetchWg sync.WaitGroup
		for w := 0; w < numWorkers; w++ {
			fetchWg.Add(1)
			go func() {
				defer fetchWg.Done()
				for bn := range blockNumsChan {
					if err := ctx.Err(); err != nil {
						resultsChan <- fetchResult{blockNum: bn, err: err}
						continue
					}

					blockFetchStart := time.Now()
					ei, err := fetchBlockWithRetry(ctx, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, bn, logger)
					blockFetchDuration := time.Since(blockFetchStart)

					resultsChan <- fetchResult{
						blockNum: bn,
						block:    ei,
						duration: blockFetchDuration,
						err:      err,
					}
				}
			}()
		}

		// Send block numbers to workers (ascending for forward)
		go func() {
			for bn := batchStart; bn <= batchEnd; bn++ {
				blockNumsChan <- bn
			}
			close(blockNumsChan)
		}()

		// Wait for all workers to finish
		go func() {
			fetchWg.Wait()
			close(resultsChan)
		}()

		// Collect results and order them
		resultsMap := make(map[int64]fetchResult)
		var totalFetchTime time.Duration
		fetchCount := 0

		for result := range resultsChan {
			if result.err != nil {
				logger.Error("FATAL: block fetch failed after all retries", "block", result.blockNum, "error", result.err)
				return fmt.Errorf("block fetch failed at %d: %w", result.blockNum, result.err)
			}
			if result.block == nil {
				logger.Error("FATAL: block data is nil", "block", result.blockNum)
				return fmt.Errorf("block data nil at %d", result.blockNum)
			}

			resultsMap[result.blockNum] = result
			totalFetchTime += result.duration
			fetchCount++

			// Log progress every 20 blocks
			if fetchCount%20 == 0 {
				avgPerBlock := totalFetchTime / time.Duration(fetchCount)
				logger.Info("fetch progress",
					"blocks_fetched", fetchCount,
					"workers", numWorkers,
					"avg_per_block_ms", avgPerBlock.Milliseconds(),
					"total_elapsed_s", fmt.Sprintf("%.2f", time.Since(fetchStartTime).Seconds()))
			}
		}

		// Build blocks slice in ascending order
		blocks := make([]*beacon.ExecInfo, 0, batchSize)
		for bn := batchStart; bn <= batchEnd; bn++ {
			if result, found := resultsMap[bn]; found {
				blocks = append(blocks, result.block)
			} else {
				// This should never happen - all blocks should have been fetched
				logger.Error("FATAL: block missing from results map", "block", bn)
				return fmt.Errorf("block %d missing from results map", bn)
			}
		}

		// Sanity check: verify all blocks were fetched
		if len(blocks) != fetchCount {
			logger.Error("FATAL: block count mismatch", "expected", fetchCount, "got", len(blocks))
			return fmt.Errorf("block count mismatch: expected %d, got %d", fetchCount, len(blocks))
		}

		fetchDuration := time.Since(fetchStartTime)
		avgFetchPerBlock := time.Duration(0)
		if fetchCount > 0 {
			avgFetchPerBlock = totalFetchTime / time.Duration(fetchCount)
		}

		logger.Info("fetch phase complete",
			"blocks", fetchCount,
			"workers", numWorkers,
			"total_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"avg_per_block_ms", avgFetchPerBlock.Milliseconds(),
			"effective_rps", fmt.Sprintf("%.2f", float64(fetchCount)/fetchDuration.Seconds()))

		// PHASE 2: Batch insert all blocks
		insertStartTime := time.Now()
		if err := db.InsertBlocksBatch(ctx, blocks); err != nil {
			logger.Error("FATAL: batch block insert failed", "error", err)
			return fmt.Errorf("batch block insert failed: %w", err)
		}
		insertDuration := time.Since(insertStartTime)

		// PHASE 3: Process validator tasks in batch
		validatorStartTime := time.Now()
		if err := processValidatorsBatch(ctx, c, db, httpc, beaconLimiter, blocks, beaconBase, beaconchaAPIKey, logger); err != nil {
			logger.Error("FATAL: validator batch processing failed", "error", err)
			return fmt.Errorf("validator batch processing failed: %w", err)
		}
		validatorDuration := time.Since(validatorStartTime)

		// Update current block number
		currentBN = blocks[len(blocks)-1].BlockNumber

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
		fetchRPS := float64(blocksProcessed) / fetchDuration.Seconds()
		validatorRPS := float64(blocksProcessed) / validatorDuration.Seconds()

		logger.Info("batch completed",
			"blocks", blocksProcessed,
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"fetch_s", fmt.Sprintf("%.2f", fetchDuration.Seconds()),
			"fetch_rps", fmt.Sprintf("%.2f", fetchRPS),
			"insert_s", fmt.Sprintf("%.3f", insertDuration.Seconds()),
			"validator_s", fmt.Sprintf("%.2f", validatorDuration.Seconds()),
			"validator_rps", fmt.Sprintf("%.2f", validatorRPS),
			"checkpoint_s", fmt.Sprintf("%.3f", checkpointDuration.Seconds()),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSecond))
	}
}

// processValidatorsBatch processes all validators in a batch efficiently
func processValidatorsBatch(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, blocks []*beacon.ExecInfo, beaconBase, beaconchaAPIKey string, logger *slog.Logger) error {
	t0 := time.Now()

	// Step 1: Collect all proposer indices from blocks
	var proposerIndices []int64
	slotToProposerIdx := make(map[int64]int64) // slot -> proposerIdx
	slotToBlock := make(map[int64]*beacon.ExecInfo) // slot -> ExecInfo for opt-in checks
	for _, ei := range blocks {
		if ei.ProposerIdx != nil {
			proposerIndices = append(proposerIndices, *ei.ProposerIdx)
			slotToProposerIdx[ei.Slot] = *ei.ProposerIdx
			slotToBlock[ei.Slot] = ei
		}
	}

	if len(proposerIndices) == 0 {
		logger.Info("[VALIDATOR] no proposers to process in batch")
		return nil
	}

	// Step 2: Batch fetch all validator pubkeys
	t1 := time.Now()
	pubkeyMap, err := beacon.FetchValidatorPubkeysBatch(ctx, httpc, beaconLimiter, beaconBase, beaconchaAPIKey, proposerIndices)
	if err != nil {
		return fmt.Errorf("batch fetch validator pubkeys: %w", err)
	}
	fetchDuration := time.Since(t1)

	// Step 3: Batch save all pubkeys to database
	t2 := time.Now()
	slotToPubkey := make(map[int64][]byte)
	for slot, proposerIdx := range slotToProposerIdx {
		if pubkey, found := pubkeyMap[proposerIdx]; found && len(pubkey) > 0 {
			slotToPubkey[slot] = pubkey
		}
	}
	if err := db.UpdateValidatorPubkeysBatch(ctx, slotToPubkey); err != nil {
		logger.Error("[VALIDATOR] failed to batch save pubkeys", "error", err)
		return fmt.Errorf("batch save validator pubkeys: %w", err)
	}
	saveDuration := time.Since(t2)

	// Step 4: Batch check opt-in status for all validators
	t3 := time.Now()

	// Collect all pubkeys and map them to slots
	var pubkeysForOptIn [][]byte
	pubkeyHexToSlot := make(map[string]int64)
	slotToPubkeyHex := make(map[int64]string)

	for slot, pubkey := range slotToPubkey {
		pubkeysForOptIn = append(pubkeysForOptIn, pubkey)
		pkHex := common.Bytes2Hex(pubkey)
		pubkeyHexToSlot[pkHex] = slot
		slotToPubkeyHex[slot] = pkHex
	}

	// Get block number for opt-in check (use first block's number as they should be close)
	var blockNum int64
	if len(blocks) > 0 {
		blockNum = blocks[0].BlockNumber
	}

	// Batch check opt-in status
	optInMap, err := ethereum.CallAreOptedInAtBlockBatch(httpc.HTTPClient, createOptionsFromCLI(c), blockNum, pubkeysForOptIn)
	if err != nil {
		logger.Error("[VALIDATOR] batch opt-in check failed", "error", err)
		return fmt.Errorf("batch opt-in check: %w", err)
	}

	// Step 5: Batch save opt-in statuses
	slotToOptIn := make(map[int64]bool)
	for slot, pkHex := range slotToPubkeyHex {
		if opted, found := optInMap[pkHex]; found {
			slotToOptIn[slot] = opted
		}
	}

	if err := db.UpdateValidatorOptInStatusBatch(ctx, slotToOptIn); err != nil {
		logger.Error("[VALIDATOR] failed to batch save opt-in status", "error", err)
		return fmt.Errorf("batch save opt-in status: %w", err)
	}
	optInDuration := time.Since(t3)

	totalDuration := time.Since(t0)

	logger.Info("[VALIDATOR] batch processing complete",
		"total_validators", len(proposerIndices),
		"pubkeys_saved", len(slotToPubkey),
		"optin_checked", len(slotToOptIn),
		"total_ms", totalDuration.Milliseconds(),
		"fetch_ms", fetchDuration.Milliseconds(),
		"save_ms", saveDuration.Milliseconds(),
		"optin_ms", optInDuration.Milliseconds())

	return nil
}

func launchValidatorTasks(ctx context.Context, c *cli.Context, db *database.DB, httpc *retryablehttp.Client, beaconLimiter *rate.Limiter, ei *beacon.ExecInfo, beaconBase, beaconchaAPIKey string, logger *slog.Logger) error { // Async validator pubkey fetch
	if ei.ProposerIdx == nil {
		return nil
	}

	t0 := time.Now()

	// Step 1: Fetch validator pubkey from beaconcha.in
	vctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	t1 := time.Now()
	vpub, err := beacon.FetchValidatorPubkey(vctx, httpc, beaconLimiter, beaconBase, beaconchaAPIKey, *ei.ProposerIdx)
	fetchDuration := time.Since(t1)
	if err != nil {
		return fmt.Errorf("fetch validator pubkey: %w", err)
	}

	// Step 2: Save pubkey to database
	t2 := time.Now()
	if len(vpub) > 0 {
		if err := db.UpdateValidatorPubkey(vctx, ei.Slot, vpub); err != nil {
			logger.Error("[VALIDATOR] failed to save pubkey", "slot", ei.Slot, "error", err)
		}
	}
	saveDuration := time.Since(t2)

	// Step 3: Retry read from database
	getCtx, getCancel := context.WithTimeout(context.Background(), 5*time.Second)
	t3 := time.Now()
	vpk, err := db.GetValidatorPubkeyWithRetry(getCtx, ei.Slot, 3, time.Second)
	getCancel()
	retryReadDuration := time.Since(t3)

	if err != nil {
		logger.Error("[VALIDATOR] pubkey not available", "slot", ei.Slot, "error", err)
		return fmt.Errorf("save validator pubkey: %w", err)
	}

	// Step 4: Check opt-in status via Ethereum RPC
	t4 := time.Now()
	opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, createOptionsFromCLI(c), ei.BlockNumber, vpk)
	optInCheckDuration := time.Since(t4)
	if err != nil {
		return fmt.Errorf("check opt-in status: %w", err)
	}

	// Step 5: Save opt-in status
	updCtx, updCancel := context.WithTimeout(context.Background(), 3*time.Second)
	t5 := time.Now()
	err = db.UpdateValidatorOptInStatus(updCtx, ei.Slot, opted)
	updCancel()
	saveOptInDuration := time.Since(t5)

	totalDuration := time.Since(t0)

	if err != nil {
		return fmt.Errorf("save opt-in status: %w", err)
	}

	logger.Info("[VALIDATOR] processing complete",
		"slot", ei.Slot,
		"proposer", *ei.ProposerIdx,
		"opted_in", opted,
		"total_ms", totalDuration.Milliseconds(),
		"fetch_ms", fetchDuration.Milliseconds(),
		"save_ms", saveDuration.Milliseconds(),
		"retry_read_ms", retryReadDuration.Milliseconds(),
		"optin_check_ms", optInCheckDuration.Milliseconds(),
		"save_optin_ms", saveOptInDuration.Milliseconds())

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
	batchSize := int64(10)

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

		// Accumulate ALL bids from ALL blocks in this batch
		// Flush ONCE at the end - no threshold checking
		allAccumulatedBids := make([]database.BidRow, 0, 100000) // ~10 blocks × 10k bids
		totalBlocksWithBids := 0
		totalBidsAccumulated := 0

		// Track tablet distribution in this batch
		tabletVersions := make(map[int64]int) // tablet_id -> version count

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

			// Retry loop for blocks with 0 bids
			retryCount := 0
			const retryDelay = 5 * time.Second
			const maxRetries = 0 // 0 means retry forever until bids appear

		retryBlock:
			// PHASE 2: Fetch bids from all relays in parallel (collect first, insert once)
			blockStartTime := time.Now()
			type relayResult struct {
				relayID       int64
				bidRows       []database.BidRow
				fetchDuration time.Duration
				err           error
			}
			resultsChan := make(chan relayResult, len(relays))
			var wg sync.WaitGroup

			// Calculate which tablet this slot hashes to (for logging)
			tabletID := block.Slot % 10 // BUCKETS 10 in table definition

			for _, rr := range relays {
				wg.Add(1)
				go func(r relay.Row) {
					defer wg.Done()

					relayCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()

					fetchStart := time.Now()
					bids, err := relay.FetchBuilderBlocksReceived(relayCtx, httpc, r.URL, block.Slot)
					fetchDur := time.Since(fetchStart)

					if err != nil {
						resultsChan <- relayResult{relayID: r.ID, fetchDuration: fetchDur, err: err}
						return
					}

					// Convert bids to rows (NO INSERT YET - collect all relays first)
					bidRows := make([]database.BidRow, 0, len(bids))
					for _, bid := range bids {
						if row, ok := relay.BuildBidInsert(block.Slot, r.ID, bid); ok {
							bidRows = append(bidRows, row)
						}
					}

					logger.Debug("relay bids fetched",
						"relay_id", r.ID,
						"slot", block.Slot,
						"tablet", tabletID,
						"bids", len(bidRows),
						"fetch_ms", fetchDur.Milliseconds())

					resultsChan <- relayResult{
						relayID:       r.ID,
						bidRows:       bidRows,
						fetchDuration: fetchDur,
					}
				}(rr)
			}

			// Wait for all relays
			go func() {
				wg.Wait()
				close(resultsChan)
			}()

			// Collect ALL bids from ALL relays first
			allBidRows := make([]database.BidRow, 0, 10000) // Pre-allocate for typical block
			successfulRelays := 0
			relaysWithBids := 0
			var errors []string
			var totalFetchDuration time.Duration

			for result := range resultsChan {
				if result.err != nil {
					errors = append(errors, fmt.Sprintf("relay_%d: %v", result.relayID, result.err))
				} else {
					successfulRelays++ // Relay responded successfully (even if 0 bids)
					if len(result.bidRows) > 0 {
						allBidRows = append(allBidRows, result.bidRows...)
						relaysWithBids++
					}
					totalFetchDuration += result.fetchDuration
				}
			}

			totalBids := len(allBidRows)

			// If all relays failed with errors, exit
			if len(errors) == len(relays) {
				logger.Error("FATAL: all relays failed with errors for block",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"tablet", tabletID,
					"errors", errors)
				return fmt.Errorf("all relays failed for block %d", block.BlockNumber)
			}

			// If some relays had errors and no bids found, retry with delay
			// If all relays succeeded but returned 0 bids, that's valid - move on
			if relaysWithBids == 0 && totalBids == 0 && len(errors) > 0 {
				retryCount++
				logger.Warn("no bids found for block and some relays had errors, retrying after delay",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"tablet", tabletID,
					"retry_attempt", retryCount,
					"retry_delay_s", retryDelay.Seconds(),
					"successful_relays", successfulRelays,
					"relays_with_errors", len(errors))

				// Check if we should stop retrying (if maxRetries > 0)
				if maxRetries > 0 && retryCount >= maxRetries {
					logger.Error("max retries reached for block with no bids, skipping",
						"block", block.BlockNumber,
						"slot", block.Slot,
						"retry_attempts", retryCount)
					// Don't update checkpoint - block will be picked up in next batch
					continue
				}

				// Wait before retry
				select {
				case <-time.After(retryDelay):
					// Retry the same block
					goto retryBlock
				case <-ctx.Done():
					logger.Info("shutdown during retry wait", "block", block.BlockNumber)
					return nil
				}
			}

			// Log if all relays were checked successfully but returned no bids (valid case)
			if relaysWithBids == 0 && totalBids == 0 && len(errors) == 0 {
				logger.Info("no bids found for block (all relays checked successfully)",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"tablet", tabletID,
					"successful_relays", successfulRelays)
			}

			// ACCUMULATE bids - insert ONCE at end of batch
			if totalBids > 0 {
				allAccumulatedBids = append(allAccumulatedBids, allBidRows...)
				totalBidsAccumulated += totalBids
				totalBlocksWithBids++

				// Track tablet usage
				tabletVersions[tabletID]++

				logger.Debug("bids accumulated for block",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"bids_from_block", totalBids,
					"total_accumulated", len(allAccumulatedBids),
					"blocks_with_bids", totalBlocksWithBids)
			}

			blockDuration := time.Since(blockStartTime)

			var avgFetchMs int64
			if successfulRelays > 0 {
				avgFetchMs = totalFetchDuration.Milliseconds() / int64(successfulRelays)
			}

			// Log with appropriate level based on success
			if len(errors) > 0 && relaysWithBids > 0 {
				// Partial success - some relays failed but we got bids from others
				logger.Warn("block bids fetched with some relay failures",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"tablet", tabletID,
					"total_bids", totalBids,
					"relays_with_bids", relaysWithBids,
					"relays_with_errors", len(errors),
					"total_relays", len(relays),
					"avg_fetch_ms", avgFetchMs,
					"fetch_duration_ms", blockDuration.Milliseconds(),
					"errors", errors)
			} else {
				// Full success
				logger.Debug("block bids fetched",
					"block", block.BlockNumber,
					"slot", block.Slot,
					"tablet", tabletID,
					"total_bids", totalBids,
					"relays_with_bids", relaysWithBids,
					"total_relays", len(relays),
					"avg_fetch_ms", avgFetchMs,
					"fetch_duration_ms", blockDuration.Milliseconds())
			}

			// Track current block for checkpoint (will save after batch insert)
			currentBlock = block.BlockNumber
		}

		// SINGLE BATCH INSERT: Insert all accumulated bids at once
		// InsertBidsBatch will internally chunk at 9000 rows per INSERT
		if len(allAccumulatedBids) > 0 {
			insertStartTime := time.Now()
			if err := db.InsertBidsBatch(ctx, allAccumulatedBids); err != nil {
				logger.Error("FATAL: failed to insert accumulated bids",
					"total_bids", len(allAccumulatedBids),
					"blocks", totalBlocksWithBids,
					"error", err)
				return fmt.Errorf("insert accumulated bids: %w", err)
			}
			insertDuration := time.Since(insertStartTime)

			logger.Info("batch bids inserted",
				"total_bids", len(allAccumulatedBids),
				"blocks_with_bids", totalBlocksWithBids,
				"insert_ms", insertDuration.Milliseconds())
		}

		// CHECKPOINT: Save checkpoint AFTER successful batch insert
		// This ensures we only mark blocks as processed after bids are safely in database
		checkpointStartTime := time.Now()
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

		batchDuration := time.Since(batchStartTime)
		blocksPerSecond := float64(len(blocks)) / batchDuration.Seconds()

		// Calculate statistics
		blocksWithBids := 0
		for _, count := range tabletVersions {
			blocksWithBids += count
		}

		logger.Info("bid batch completed",
			"blocks_processed", len(blocks),
			"blocks_with_bids", blocksWithBids,
			"total_bids_accumulated", totalBidsAccumulated,
			"tablets_written", len(tabletVersions),
			"total_s", fmt.Sprintf("%.2f", batchDuration.Seconds()),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSecond))

		// Log detailed tablet distribution
		if len(tabletVersions) > 0 {
			logger.Debug("tablet block distribution", "tablet_blocks", tabletVersions)
		}
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

	// Fill gaps if requested (one-time operation at startup)
	if c.Bool("fill-gaps") {
		initLogger.Info("[FILL_GAPS] gap filling enabled, detecting and filling missing blocks...")
		if err := fillGaps(ctx, c, db, httpc, beaconLimiter, rpcURL, beaconBase, beaconchaAPIKey, initLogger); err != nil {
			initLogger.Error("[FILL_GAPS] failed to fill gaps", "error", err)
			return fmt.Errorf("fill gaps failed: %w", err)
		}
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
