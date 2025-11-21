package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log/slog"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/core/types"
	w3 "github.com/lmittmann/w3"
	eth "github.com/lmittmann/w3/module/eth"
	w3types "github.com/lmittmann/w3/w3types"
)

// ChunkStatus tracks the processing status of a block range
type ChunkStatus struct {
	StartBlock int64     `json:"start_block"`
	EndBlock   int64     `json:"end_block"`
	Status     string    `json:"status"` // "completed" or "failed"
	ErrorMsg   string    `json:"error_msg,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type blockRange struct {
	start int64
	end   int64
}

type fetchResult struct {
	blockRange
	blocks   []*types.Block
	receipts []types.Receipts
	err      error
}

// aggregateMetrics tracks performance across all workers
type aggregateMetrics struct {
	blocksProcessed   atomic.Int64
	txsProcessed      atomic.Int64
	receiptsProcessed atomic.Int64
	logsProcessed     atomic.Int64
	chunksFailed      atomic.Int64
	chunksRetried     atomic.Int64
	startTime         time.Time
	lastResetTime     time.Time
}

func (m *aggregateMetrics) recordSuccess(blocks []*types.Block, receipts []types.Receipts) {
	m.blocksProcessed.Add(int64(len(blocks)))

	txCount := 0
	logCount := 0
	for i := range blocks {
		txCount += len(blocks[i].Transactions())
		for j := range receipts[i] {
			logCount += len(receipts[i][j].Logs)
		}
	}

	m.txsProcessed.Add(int64(txCount))
	m.receiptsProcessed.Add(int64(txCount))
	m.logsProcessed.Add(int64(logCount))
}

func (m *aggregateMetrics) recordFailure() {
	m.chunksFailed.Add(1)
}

func (m *aggregateMetrics) recordRetry() {
	m.chunksRetried.Add(1)
}

func (m *aggregateMetrics) getSnapshot() (blocks, txs, receipts, logs, failed, retried int64, duration time.Duration) {
	return m.blocksProcessed.Load(),
		m.txsProcessed.Load(),
		m.receiptsProcessed.Load(),
		m.logsProcessed.Load(),
		m.chunksFailed.Load(),
		m.chunksRetried.Load(),
		time.Since(m.startTime)
}

// openChunkDB opens the embedded BadgerDB database for tracking chunk status
func openChunkDB(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable BadgerDB's verbose logging
	return badger.Open(opts)
}

// markChunkStatus records the status of a processed chunk
func markChunkStatus(db *badger.DB, start, end int64, status, errMsg string) error {
	cs := ChunkStatus{
		StartBlock: start,
		EndBlock:   end,
		Status:     status,
		ErrorMsg:   errMsg,
		UpdatedAt:  time.Now(),
	}

	data, err := json.Marshal(cs)
	if err != nil {
		return err
	}

	return db.Update(func(txn *badger.Txn) error {
		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, uint64(start))
		return txn.Set(key, data)
	})
}

// findGaps returns all incomplete chunks below the checkpoint
func findGaps(db *badger.DB, checkpoint int64) ([]blockRange, error) {
	var gaps []blockRange

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 100
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var cs ChunkStatus
				if err := json.Unmarshal(v, &cs); err != nil {
					return err
				}

				// Gap: chunk below checkpoint that's not completed
				if cs.EndBlock <= checkpoint && cs.Status != "completed" {
					gaps = append(gaps, blockRange{cs.StartBlock, cs.EndBlock})
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return gaps, err
}

// cleanupOldChunks removes completed chunks far below the checkpoint
func cleanupOldChunks(db *badger.DB, checkpoint int64) error {
	return db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		var keysToDelete [][]byte

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)

			err := item.Value(func(v []byte) error {
				var cs ChunkStatus
				if err := json.Unmarshal(v, &cs); err != nil {
					return err
				}

				// Delete completed chunks >10,000 blocks behind checkpoint
				if cs.Status == "completed" && cs.EndBlock < checkpoint-10000 {
					keysToDelete = append(keysToDelete, key)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		// Delete outside iterator
		for _, key := range keysToDelete {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}

		if len(keysToDelete) > 0 {
			slog.Info("Cleaned up old chunks", "count", len(keysToDelete))
		}

		return nil
	})
}

// parallelForwardIndex runs parallel workers to fetch and index blocks
func parallelForwardIndex(ctx context.Context, client *w3.Client, chunkDB *badger.DB, pollInterval time.Duration, batchSize, numWorkers int, flushInterval time.Duration, kafkaOpts kafkaOptions, startBlock int64, logger *slog.Logger) {
	// Get checkpoint from Kafka
	checkpoint, hasCheckpoint, err := getKafkaCheckpoint(kafkaOpts)
	if err != nil {
		logger.Error("Failed to get Kafka checkpoint", "err", err)
		return
	}

	if !hasCheckpoint && startBlock > 0 {
		checkpoint = startBlock - 1
		logger.Info("No Kafka checkpoint found, using start block", "start_block", startBlock, "checkpoint", checkpoint)
	} else if hasCheckpoint {
		logger.Info("Starting from Kafka checkpoint", "checkpoint", checkpoint, "will_resume_from", checkpoint+1)
	}

	// Find and process any gaps first
	gaps, err := findGaps(chunkDB, checkpoint)
	if err != nil {
		logger.Error("Failed to find gaps", "err", err)
	} else if len(gaps) > 0 {
		logger.Info("Found gaps to process", "count", len(gaps), "checkpoint", checkpoint)
	} else {
		logger.Info("No gaps found, will process new blocks only", "starting_from", checkpoint+1)
	}

	// Initialize aggregate metrics
	metrics := &aggregateMetrics{
		startTime:     time.Now(),
		lastResetTime: time.Now(),
	}

	// Start periodic metrics reporter
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blocks, txs, receipts, logs, failed, retried, duration := metrics.getSnapshot()
				durationSec := duration.Seconds()
				if durationSec > 0 {
					logger.Info("Aggregate throughput",
						"workers", numWorkers,
						"total_blocks", blocks,
						"total_txs", txs,
						"total_receipts", receipts,
						"total_logs", logs,
						"chunks_failed", failed,
						"chunks_retried", retried,
						"duration", duration,
						"blocks_per_sec", float64(blocks)/durationSec,
						"txs_per_sec", float64(txs)/durationSec,
						"receipts_per_sec", float64(receipts)/durationSec,
						"logs_per_sec", float64(logs)/durationSec,
					)
				}
			}
		}
	}()

	// Channels
	fetchQueue := make(chan blockRange, numWorkers*10)
	results := make(chan fetchResult, numWorkers*10)

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		workerID := i
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case r := <-fetchQueue:
					result := fetchChunk(client, r)
					results <- result
					if result.err != nil {
						logger.Warn("Worker fetch failed, will retry",
							"worker", workerID,
							"start", r.start,
							"end", r.end,
							"err", result.err)
					}
				}
			}
		}()
	}

	// Feed work: gaps first, then new blocks
	go func() {
		// Process gaps
		for _, gap := range gaps {
			fetchQueue <- gap
		}

		// Process new blocks
		next := checkpoint + 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var latest *big.Int
			if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
				logger.Error("Failed to get latest block", "err", err)
				time.Sleep(pollInterval)
				continue
			}

			latestBlock := latest.Int64()
			chunksQueued := 0
			startOfBatch := next
			for next <= latestBlock {
				end := next + int64(batchSize) - 1
				if end > latestBlock {
					end = latestBlock
				}
				fetchQueue <- blockRange{next, end}
				chunksQueued++
				next = end + 1
			}

			if chunksQueued > 0 {
				logger.Info("Queued chunks for processing",
					"chunks", chunksQueued,
					"from_block", startOfBatch,
					"to_block", latestBlock,
					"chain_head", latestBlock)
			}

			time.Sleep(pollInterval)
		}
	}()

	// Process results
	for {
		select {
		case <-ctx.Done():
			logger.Info("Parallel indexing stopped")
			return
		case result := <-results:
			if result.err != nil {
				// Mark as failed and retry
				metrics.recordFailure()
				metrics.recordRetry()
				markChunkStatus(chunkDB, result.start, result.end, "failed", result.err.Error())
				fetchQueue <- result.blockRange // Retry
				continue
			}

			// Process batch (write to Kafka)
			if kafkaOpts.enabled {
				err = processBatchKafka(result.blocks, result.receipts, kafkaOpts, nil)
			} else {
				err = processBatchSQL(result.blocks, result.receipts, nil, flushInterval)
			}

			if err != nil {
				logger.Error("Failed to process batch", "start", result.start, "end", result.end, "err", err)
				metrics.recordFailure()
				markChunkStatus(chunkDB, result.start, result.end, "failed", err.Error())
				// For Kafka/SQL write failures, don't retry - fatal error
				return
			}

			// Mark as completed and record metrics
			markChunkStatus(chunkDB, result.start, result.end, "completed", "")
			metrics.recordSuccess(result.blocks, result.receipts)

			// Log completion for visibility
			txCount := 0
			for i := range result.blocks {
				txCount += len(result.blocks[i].Transactions())
			}
			logger.Info("Chunk completed",
				"from_block", result.start,
				"to_block", result.end,
				"num_blocks", len(result.blocks),
				"num_txs", txCount)

			// Periodically cleanup old chunks
			if result.start%1000 == 0 {
				cleanupOldChunks(chunkDB, result.end)
			}
		}
	}
}

// fetchChunk fetches a range of blocks from the RPC endpoint
func fetchChunk(client *w3.Client, r blockRange) fetchResult {
	n := int(r.end - r.start + 1)
	blocks := make([]*types.Block, n)
	receipts := make([]types.Receipts, n)

	var calls []w3types.RPCCaller
	for i := int64(0); i < int64(n); i++ {
		num := big.NewInt(r.start + i)
		calls = append(calls, eth.BlockByNumber(num).Returns(&blocks[i]))
		calls = append(calls, eth.BlockReceipts(num).Returns(&receipts[i]))
	}

	err := client.Call(calls...)
	if err != nil {
		return fetchResult{blockRange: r, err: err}
	}

	return fetchResult{
		blockRange: r,
		blocks:     blocks,
		receipts:   receipts,
		err:        nil,
	}
}
