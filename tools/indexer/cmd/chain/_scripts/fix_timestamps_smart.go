package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type BlockUpdate struct {
	Number    int64
	Timestamp int64
	Hash      string
}

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockResult struct {
	Number    string `json:"number"`
	Timestamp string `json:"timestamp"`
	Hash      string `json:"hash"`
}

var (
	versionErrorCount int64
	versionErrorMu    sync.Mutex
	lastCompactTime   time.Time
)

func main() {
	rpcURL := flag.String("rpc", "https://chainrpc-v1.mev-commit.xyz", "RPC endpoint URL")
	dsn := flag.String("dsn", "", "Database DSN (required)")
	batchSize := flag.Int("batch-size", 2000, "Number of blocks per batch")
	workers := flag.Int("workers", 30, "Number of concurrent RPC workers")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (don't update database)")
	startBlock := flag.Int64("start-block", 0, "Start block number (optional)")
	endBlock := flag.Int64("end-block", -1, "End block number (optional, -1 for all)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting SMART timestamp fix (REVERSE ORDER with auto-pause)")
	log.Printf("RPC: %s", *rpcURL)
	log.Printf("Batch size: %d blocks per batch", *batchSize)
	log.Printf("Workers: %d", *workers)
	log.Printf("Dry run: %v", *dryRun)

	// Connect to database
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get block range
	var minBlock, maxBlock int64
	query := "SELECT MIN(number), MAX(number) FROM blocks"
	if err := db.QueryRow(query).Scan(&minBlock, &maxBlock); err != nil {
		log.Fatalf("Failed to get block range: %v", err)
	}

	if *startBlock > 0 {
		minBlock = *startBlock
	}
	if *endBlock >= 0 && *endBlock < maxBlock {
		maxBlock = *endBlock
	}

	log.Printf("Block range: %d DOWN TO %d (total: %d blocks) - NEWEST FIRST", maxBlock, minBlock, maxBlock-minBlock+1)

	ctx := context.Background()
	lastCompactTime = time.Now()

	// Channels
	blockChan := make(chan int64, *workers*2)
	resultChan := make(chan BlockUpdate, *batchSize*2)

	var wg sync.WaitGroup

	// Stats
	var (
		totalFetched int64
		totalUpdated int64
		totalSkipped int64
		totalErrors  int64
		mu           sync.Mutex
	)

	// Track successfully updated blocks
	updatedBlocks := make(map[int64]bool)
	var updatedMu sync.Mutex

	// Start RPC workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}

			for blockNum := range blockChan {
				// Check if already updated
				updatedMu.Lock()
				if updatedBlocks[blockNum] {
					updatedMu.Unlock()
					mu.Lock()
					totalSkipped++
					mu.Unlock()
					continue
				}
				updatedMu.Unlock()

				correctTimestamp, blockHash, err := fetchBlockTimestamp(client, *rpcURL, blockNum)
				if err != nil {
					log.Printf("[Worker %d] Error fetching block %d: %v", workerID, blockNum, err)
					mu.Lock()
					totalErrors++
					mu.Unlock()
					continue
				}

				resultChan <- BlockUpdate{
					Number:    blockNum,
					Timestamp: correctTimestamp,
					Hash:      blockHash,
				}

				mu.Lock()
				totalFetched++
				mu.Unlock()
			}
		}(i)
	}

	// Start batch updater
	var updateWg sync.WaitGroup
	updateWg.Add(1)
	go func() {
		defer updateWg.Done()

		batch := make([]BlockUpdate, 0, *batchSize)
		lastUpdate := time.Now()

		for update := range resultChan {
			batch = append(batch, update)

			// Flush batch when full or every 5 seconds
			if len(batch) >= *batchSize || time.Since(lastUpdate) > 5*time.Second {
				if !*dryRun && len(batch) > 0 {
					updated, successBlocks := batchUpdateBlocks(ctx, db, batch)

					// Mark successfully updated blocks
					updatedMu.Lock()
					for _, blockNum := range successBlocks {
						updatedBlocks[blockNum] = true
					}
					updatedMu.Unlock()

					mu.Lock()
					totalUpdated += int64(updated)
					mu.Unlock()

					log.Printf("Progress: %d fetched, %d updated, %d skipped, %d errors - Batch: %d blocks (range: %d to %d)",
						totalFetched, totalUpdated, totalSkipped, totalErrors, len(batch), batch[len(batch)-1].Number, batch[0].Number)
				}
				batch = batch[:0]
				lastUpdate = time.Now()
			}
		}

		// Final flush
		if !*dryRun && len(batch) > 0 {
			updated, _ := batchUpdateBlocks(ctx, db, batch)
			mu.Lock()
			totalUpdated += int64(updated)
			mu.Unlock()
			log.Printf("Final batch: %d blocks updated", updated)
		}
	}()

	// Feed blocks to workers IN REVERSE ORDER (max to min)
	startTime := time.Now()
	go func() {
		for blockNum := maxBlock; blockNum >= minBlock; blockNum-- {
			blockChan <- blockNum
		}
		close(blockChan)
	}()

	// Wait for RPC workers
	wg.Wait()
	close(resultChan)

	// Wait for updater
	updateWg.Wait()

	elapsed := time.Since(startTime)

	log.Printf("========================================")
	log.Printf("Timestamp fix completed!")
	log.Printf("Total fetched: %d", totalFetched)
	log.Printf("Total updated: %d", totalUpdated)
	log.Printf("Total skipped: %d", totalSkipped)
	log.Printf("Total errors: %d", totalErrors)
	log.Printf("Duration: %v", elapsed)
	if elapsed.Seconds() > 0 {
		log.Printf("Rate: %.2f blocks/sec", float64(totalUpdated)/elapsed.Seconds())
	}
	log.Printf("========================================")

	if *dryRun {
		log.Printf("DRY RUN MODE - No changes were made to the database")
	}
}

func batchUpdateBlocks(ctx context.Context, db *sql.DB, updates []BlockUpdate) (int, []int64) {
	if len(updates) == 0 {
		return 0, nil
	}

	const maxCaseConditions = 25 // Smaller chunks to reduce version pressure
	totalUpdated := 0
	var successBlocks []int64

	// Process in chunks
	for chunkStart := 0; chunkStart < len(updates); chunkStart += maxCaseConditions {
		chunkEnd := chunkStart + maxCaseConditions
		if chunkEnd > len(updates) {
			chunkEnd = len(updates)
		}
		chunk := updates[chunkStart:chunkEnd]

		// Check if we need to pause for compaction
		versionErrorMu.Lock()
		if versionErrorCount > 0 && time.Since(lastCompactTime) < 3*time.Minute {
			versionErrorMu.Unlock()
			// Pause to let StarRocks compact
			sleepTime := 3*time.Minute - time.Since(lastCompactTime)
			log.Printf("Too many version errors - pausing %v for StarRocks compaction...", sleepTime)
			time.Sleep(sleepTime)
			versionErrorMu.Lock()
			versionErrorCount = 0
			lastCompactTime = time.Now()
			versionErrorMu.Unlock()
		} else {
			versionErrorMu.Unlock()
		}

		var caseParts []string
		var numberList []string
		var chunkBlockNums []int64

		for _, update := range chunk {
			caseParts = append(caseParts, fmt.Sprintf("WHEN %d THEN %d", update.Number, update.Timestamp))
			numberList = append(numberList, fmt.Sprintf("%d", update.Number))
			chunkBlockNums = append(chunkBlockNums, update.Number)
		}

		query := fmt.Sprintf(`UPDATE blocks
			SET timestamp = CASE number
				%s
			END
			WHERE number IN (%s)`,
			strings.Join(caseParts, " "),
			strings.Join(numberList, ","))

		result, err := db.ExecContext(ctx, query)
		if err != nil {
			if strings.Contains(err.Error(), "too many versions") {
				log.Printf("Version limit hit - will pause after this batch")
				versionErrorMu.Lock()
				versionErrorCount++
				versionErrorMu.Unlock()
			} else {
				log.Printf("Error updating chunk of %d blocks: %v", len(chunk), err)
			}
			// Don't mark as success if error
			time.Sleep(2 * time.Second) // Longer pause on error
			continue
		}

		rows, _ := result.RowsAffected()
		totalUpdated += int(rows)
		successBlocks = append(successBlocks, chunkBlockNums...)

		// Minimal delay to maximize throughput
		time.Sleep(50 * time.Millisecond)
	}

	return totalUpdated, successBlocks
}

func fetchBlockTimestamp(client *http.Client, rpcURL string, blockNum int64) (int64, string, error) {
	blockHex := fmt.Sprintf("0x%x", blockNum)

	reqBody := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockHex, false},
		ID:      1,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return 0, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := client.Post(rpcURL, "application/json", bytes.NewReader(reqBytes))
	if err != nil {
		return 0, "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return 0, "", fmt.Errorf("failed to parse response: %w", err)
	}

	if rpcResp.Error != nil {
		return 0, "", fmt.Errorf("RPC error: %s", rpcResp.Error.Message)
	}

	var block BlockResult
	if err := json.Unmarshal(rpcResp.Result, &block); err != nil {
		return 0, "", fmt.Errorf("failed to parse block: %w", err)
	}

	if block.Timestamp == "" {
		return 0, "", fmt.Errorf("block %d does not exist", blockNum)
	}

	timestampHex := strings.TrimPrefix(block.Timestamp, "0x")
	timestamp, err := strconv.ParseInt(timestampHex, 16, 64)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return timestamp, block.Hash, nil
}
