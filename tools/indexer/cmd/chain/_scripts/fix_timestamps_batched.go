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

func main() {
	rpcURL := flag.String("rpc", "https://chainrpc-v1.mev-commit.xyz", "RPC endpoint URL")
	dsn := flag.String("dsn", "", "Database DSN (required)")
	batchSize := flag.Int("batch-size", 500, "Number of blocks to collect before updating database")
	workers := flag.Int("workers", 10, "Number of concurrent RPC workers")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (don't update database)")
	startBlock := flag.Int64("start-block", 0, "Start block number (optional)")
	endBlock := flag.Int64("end-block", -1, "End block number (optional, -1 for all)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting timestamp fix (batched mode)")
	log.Printf("RPC: %s", *rpcURL)
	log.Printf("Batch size: %d blocks per DB update", *batchSize)
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

	log.Printf("Block range: %d to %d (total: %d blocks)", minBlock, maxBlock, maxBlock-minBlock+1)

	ctx := context.Background()

	// Channels
	blockChan := make(chan int64, *workers*2)
	resultChan := make(chan BlockUpdate, *batchSize*2)

	var wg sync.WaitGroup

	// Stats
	var (
		totalFetched int64
		totalUpdated int64
		totalErrors  int64
		mu           sync.Mutex
	)

	// Start RPC workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}

			for blockNum := range blockChan {
				correctTimestamp, blockHash, err := fetchBlockTimestamp(client, *rpcURL, blockNum)
				if err != nil {
					log.Printf("[Worker %d] Error fetching block %d: %v", workerID, blockNum, err)
					mu.Lock()
					totalErrors++
					mu.Unlock()
					continue
				}

				// Validate timestamp
				if correctTimestamp < 1577836800000 || correctTimestamp > 1893456000000 {
					log.Printf("[Worker %d] WARNING: Block %d has suspicious timestamp: %d",
						workerID, blockNum, correctTimestamp)
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

			// Flush batch when full or every 10 seconds
			if len(batch) >= *batchSize || time.Since(lastUpdate) > 10*time.Second {
				if !*dryRun && len(batch) > 0 {
					updated := batchUpdateBlocks(ctx, db, batch)
					mu.Lock()
					totalUpdated += int64(updated)
					mu.Unlock()

					log.Printf("Progress: %d fetched, %d updated, %d errors - Latest batch: %d blocks",
						totalFetched, totalUpdated, totalErrors, len(batch))
				}
				batch = batch[:0]
				lastUpdate = time.Now()
			}
		}

		// Final flush
		if !*dryRun && len(batch) > 0 {
			updated := batchUpdateBlocks(ctx, db, batch)
			mu.Lock()
			totalUpdated += int64(updated)
			mu.Unlock()
			log.Printf("Final batch: %d blocks updated", updated)
		}
	}()

	// Feed blocks to workers
	startTime := time.Now()
	go func() {
		for blockNum := minBlock; blockNum <= maxBlock; blockNum++ {
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
	log.Printf("Total errors: %d", totalErrors)
	log.Printf("Duration: %v", elapsed)
	log.Printf("Rate: %.2f blocks/sec", float64(totalFetched)/elapsed.Seconds())
	log.Printf("========================================")

	if *dryRun {
		log.Printf("DRY RUN MODE - No changes were made to the database")
	}
}

func batchUpdateBlocks(ctx context.Context, db *sql.DB, updates []BlockUpdate) int {
	const maxCaseConditions = 100 // StarRocks limit

	totalUpdated := 0

	// Process in chunks to avoid exceeding StarRocks limits
	for chunkStart := 0; chunkStart < len(updates); chunkStart += maxCaseConditions {
		chunkEnd := chunkStart + maxCaseConditions
		if chunkEnd > len(updates) {
			chunkEnd = len(updates)
		}
		chunk := updates[chunkStart:chunkEnd]

		var caseParts []string
		var numberList []string

		for _, update := range chunk {
			caseParts = append(caseParts, fmt.Sprintf("WHEN %d THEN %d", update.Number, update.Timestamp))
			numberList = append(numberList, fmt.Sprintf("%d", update.Number))
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
			log.Printf("Error updating chunk of %d blocks: %v", len(chunk), err)
			continue
		}

		rows, _ := result.RowsAffected()
		totalUpdated += int(rows)
	}

	return totalUpdated
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
