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
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	_ "github.com/go-sql-driver/mysql"
)

type BlockUpdate struct {
	Number    int64 `json:"number"`
	Timestamp int64 `json:"timestamp"`
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
	kafkaBrokers := flag.String("kafka", "mevcommit-message-queue-cluster-kafka-bootstrap.kafka:9092", "Kafka brokers")
	kafkaTopic := flag.String("topic", "block_timestamps", "Kafka topic")
	dsn := flag.String("dsn", "", "Database DSN for reading block range")
	workers := flag.Int("workers", 10, "Number of parallel RPC workers")
	batchSize := flag.Int("batch", 50, "Number of blocks per RPC batch request")
	startBlock := flag.Int64("start-block", 0, "Start block number (optional)")
	endBlock := flag.Int64("end-block", -1, "End block number (optional, -1 for all)")
	maxBlocks := flag.Int64("max-blocks", 1000000, "Maximum blocks to fetch (for testing)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting BATCH RPC timestamp fetch")
	log.Printf("RPC: %s", *rpcURL)
	log.Printf("Workers: %d", *workers)
	log.Printf("RPC Batch Size: %d blocks per HTTP request", *batchSize)
	log.Printf("Max blocks: %d", *maxBlocks)

	// Connect to database to get block range
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

	// Limit to maxBlocks
	if maxBlock-minBlock+1 > *maxBlocks {
		minBlock = maxBlock - *maxBlocks + 1
	}

	totalBlocks := maxBlock - minBlock + 1
	log.Printf("Block range: %d DOWN TO %d (total: %d blocks) - NEWEST FIRST", maxBlock, minBlock, totalBlocks)

	ctx := context.Background()

	// Setup Kafka with async pattern
	brokers := strings.Split(*kafkaBrokers, ",")
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(*kafkaTopic),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.ProducerBatchMaxBytes(1048576), // 1MB
		kgo.ProducerLinger(50*time.Millisecond),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafkaClient.Close()
	log.Printf("Writing to Kafka: %s", *kafkaTopic)

	// Channels - use batch size for work units
	batchChan := make(chan []int64, *workers*2)
	resultChan := make(chan BlockUpdate, *batchSize*(*workers)*2)

	// Stats
	var (
		totalFetched   int64
		totalProduced  int64
		totalErrors    int64
		startTime      = time.Now()
	)

	// Start producer goroutine (async Kafka)
	var producerWg sync.WaitGroup
	producerWg.Add(1)
	go func() {
		defer producerWg.Done()

		for update := range resultChan {
			// Async Kafka produce
			data, _ := json.Marshal(update)
			record := &kgo.Record{
				Key:   []byte(fmt.Sprintf("%d", update.Number)),
				Value: data,
			}

			// Fire and forget with callback
			kafkaClient.Produce(ctx, record, func(r *kgo.Record, err error) {
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
				} else {
					atomic.AddInt64(&totalProduced, 1)
				}
			})
		}

		// Final flush for Kafka
		kafkaClient.Flush(ctx)
	}()

	// Start RPC workers - BATCH RPC REQUESTS
	var fetchWg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		httpClient := &http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}

		fetchWg.Add(1)
		go func(workerID int, client *http.Client) {
			defer fetchWg.Done()

			for blockBatch := range batchChan {
				updates, errors := fetchBlockBatch(client, *rpcURL, blockBatch)

				atomic.AddInt64(&totalFetched, int64(len(updates)))
				atomic.AddInt64(&totalErrors, int64(errors))

				for _, update := range updates {
					resultChan <- update
				}
			}
		}(i, httpClient)
	}

	// Stats reporter
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		lastFetched := int64(0)
		lastTime := time.Now()

		for range ticker.C {
			current := atomic.LoadInt64(&totalFetched)
			produced := atomic.LoadInt64(&totalProduced)
			errors := atomic.LoadInt64(&totalErrors)

			now := time.Now()
			elapsed := now.Sub(lastTime).Seconds()
			rate := float64(current-lastFetched) / elapsed

			progress := float64(current) / float64(totalBlocks) * 100
			remaining := totalBlocks - current
			eta := time.Duration(0)
			if rate > 0 {
				eta = time.Duration(float64(remaining)/rate) * time.Second
			}

			log.Printf("Fetched: %d/%d (%.1f%%) | Rate: %.0f blk/s | Produced: %d | Errors: %d | ETA: %v",
				current, totalBlocks, progress, rate, produced, errors, eta.Round(time.Second))

			lastFetched = current
			lastTime = now
		}
	}()

	// Feed blocks in batches (reverse order)
	go func() {
		batch := make([]int64, 0, *batchSize)
		for blockNum := maxBlock; blockNum >= minBlock; blockNum-- {
			batch = append(batch, blockNum)

			if len(batch) >= *batchSize {
				batchCopy := make([]int64, len(batch))
				copy(batchCopy, batch)
				batchChan <- batchCopy
				batch = batch[:0]
			}
		}

		// Send remaining batch
		if len(batch) > 0 {
			batchChan <- batch
		}
		close(batchChan)
	}()

	// Wait for all fetches
	fetchWg.Wait()
	close(resultChan)

	// Wait for producer
	producerWg.Wait()

	elapsed := time.Since(startTime)
	finalFetched := atomic.LoadInt64(&totalFetched)
	finalProduced := atomic.LoadInt64(&totalProduced)
	finalErrors := atomic.LoadInt64(&totalErrors)

	log.Printf("========================================")
	log.Printf("COMPLETED!")
	log.Printf("Total fetched: %d", finalFetched)
	log.Printf("Total produced: %d", finalProduced)
	log.Printf("Total errors: %d", finalErrors)
	log.Printf("Duration: %v", elapsed)
	if elapsed.Seconds() > 0 {
		log.Printf("Average rate: %.2f blocks/sec", float64(finalFetched)/elapsed.Seconds())
	}
	log.Printf("========================================")
}

// Fetch multiple blocks in a SINGLE JSON-RPC batch request
func fetchBlockBatch(client *http.Client, rpcURL string, blockNums []int64) ([]BlockUpdate, int) {
	if len(blockNums) == 0 {
		return nil, 0
	}

	// Build batch RPC request
	requests := make([]RPCRequest, len(blockNums))
	for i, blockNum := range blockNums {
		blockHex := fmt.Sprintf("0x%x", blockNum)
		requests[i] = RPCRequest{
			JSONRPC: "2.0",
			Method:  "eth_getBlockByNumber",
			Params:  []interface{}{blockHex, false},
			ID:      i,
		}
	}

	reqBytes, err := json.Marshal(requests)
	if err != nil {
		log.Printf("ERROR: Failed to marshal requests: %v", err)
		return nil, len(blockNums)
	}

	req, err := http.NewRequest("POST", rpcURL, bytes.NewReader(reqBytes))
	if err != nil {
		log.Printf("ERROR: Failed to create request: %v", err)
		return nil, len(blockNums)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "curl/7.68.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: HTTP request failed for batch of %d blocks: %v", len(blockNums), err)
		return nil, len(blockNums)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: HTTP status %d for batch of %d blocks", resp.StatusCode, len(blockNums))
		return nil, len(blockNums)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read response body: %v", err)
		return nil, len(blockNums)
	}

	var responses []RPCResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		log.Printf("ERROR: Failed to unmarshal response for batch of %d blocks: %v (body len: %d)", len(blockNums), err, len(body))
		return nil, len(blockNums)
	}

	// Parse results
	updates := make([]BlockUpdate, 0, len(responses))
	errors := 0

	for _, rpcResp := range responses {
		if rpcResp.Error != nil {
			errors++
			continue
		}

		var block BlockResult
		if err := json.Unmarshal(rpcResp.Result, &block); err != nil {
			errors++
			continue
		}

		if block.Timestamp == "" {
			errors++
			continue
		}

		timestampHex := strings.TrimPrefix(block.Timestamp, "0x")
		timestamp, err := strconv.ParseInt(timestampHex, 16, 64)
		if err != nil {
			errors++
			continue
		}

		// Map response ID back to block number
		if rpcResp.ID < len(blockNums) {
			updates = append(updates, BlockUpdate{
				Number:    blockNums[rpcResp.ID],
				Timestamp: timestamp,
			})
		}
	}

	return updates, errors
}
