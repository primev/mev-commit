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
	rpcBatchSize := flag.Int("rpc-batch", 50, "Number of blocks to fetch per RPC batch")
	kafkaBatchSize := flag.Int("kafka-batch", 1000, "Number of updates per Kafka batch")
	workers := flag.Int("workers", 10, "Number of concurrent RPC workers")
	startBlock := flag.Int64("start-block", 0, "Start block number (optional)")
	endBlock := flag.Int64("end-block", -1, "End block number (optional, -1 for all)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting Kafka-based timestamp fix")
	log.Printf("RPC: %s", *rpcURL)
	log.Printf("Kafka: %s", *kafkaBrokers)
	log.Printf("Topic: %s", *kafkaTopic)
	log.Printf("RPC batch size: %d", *rpcBatchSize)
	log.Printf("Kafka batch size: %d", *kafkaBatchSize)
	log.Printf("Workers: %d", *workers)

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

	log.Printf("Block range: %d DOWN TO %d (total: %d blocks) - NEWEST FIRST", maxBlock, minBlock, maxBlock-minBlock+1)

	// Create Kafka client
	brokers := strings.Split(*kafkaBrokers, ",")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(*kafkaTopic),
		kgo.ProducerBatchMaxBytes(1048576), // 1MB batches
		kgo.ProducerLinger(100*time.Millisecond),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Channels
	blockChan := make(chan int64, *workers*2)
	resultChan := make(chan BlockUpdate, *kafkaBatchSize*2)

	var wg sync.WaitGroup

	// Stats
	var (
		totalFetched  int64
		totalProduced int64
		totalErrors   int64
		mu            sync.Mutex
	)

	// Start RPC workers with batching
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := &http.Client{Timeout: 30 * time.Second}

			// Collect blocks into batches for RPC
			batch := make([]int64, 0, *rpcBatchSize)
			for blockNum := range blockChan {
				batch = append(batch, blockNum)

				if len(batch) >= *rpcBatchSize {
					results, errors := fetchBlockTimestampBatch(client, *rpcURL, batch)

					for _, update := range results {
						resultChan <- update
					}

					mu.Lock()
					totalFetched += int64(len(results))
					totalErrors += int64(errors)
					mu.Unlock()

					batch = batch[:0]
				}
			}

			// Process remaining batch
			if len(batch) > 0 {
				results, errors := fetchBlockTimestampBatch(client, *rpcURL, batch)

				for _, update := range results {
					resultChan <- update
				}

				mu.Lock()
				totalFetched += int64(len(results))
				totalErrors += int64(errors)
				mu.Unlock()
			}
		}(i)
	}

	// Start Kafka batch producer
	var producerWg sync.WaitGroup
	producerWg.Add(1)
	go func() {
		defer producerWg.Done()

		batch := make([]BlockUpdate, 0, *kafkaBatchSize)
		lastFlush := time.Now()

		for update := range resultChan {
			batch = append(batch, update)

			// Flush when batch is full or every 5 seconds
			if len(batch) >= *kafkaBatchSize || time.Since(lastFlush) > 5*time.Second {
				if len(batch) > 0 {
					produced := produceKafkaBatch(ctx, client, batch)

					mu.Lock()
					totalProduced += int64(produced)
					mu.Unlock()

					log.Printf("Progress: %d fetched, %d produced to Kafka, %d errors - Batch: %d records (blocks %d to %d)",
						totalFetched, totalProduced, totalErrors, len(batch), batch[len(batch)-1].Number, batch[0].Number)

					batch = batch[:0]
					lastFlush = time.Now()
				}
			}
		}

		// Final flush
		if len(batch) > 0 {
			produced := produceKafkaBatch(ctx, client, batch)
			mu.Lock()
			totalProduced += int64(produced)
			mu.Unlock()
			log.Printf("Final batch: %d records produced to Kafka", produced)
		}

		// Flush any remaining messages
		if err := client.Flush(ctx); err != nil {
			log.Printf("Error flushing Kafka: %v", err)
		}
	}()

	// Feed blocks to workers IN REVERSE ORDER
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

	// Wait for Kafka producer
	producerWg.Wait()

	elapsed := time.Since(startTime)

	log.Printf("========================================")
	log.Printf("Timestamp fix completed!")
	log.Printf("Total fetched: %d", totalFetched)
	log.Printf("Total produced to Kafka: %d", totalProduced)
	log.Printf("Total errors: %d", totalErrors)
	log.Printf("Duration: %v", elapsed)
	if elapsed.Seconds() > 0 {
		log.Printf("Rate: %.2f blocks/sec", float64(totalProduced)/elapsed.Seconds())
	}
	log.Printf("========================================")
}

// Batch RPC fetch - sends multiple requests in parallel
func fetchBlockTimestampBatch(client *http.Client, rpcURL string, blockNums []int64) ([]BlockUpdate, int) {
	results := make([]BlockUpdate, 0, len(blockNums))
	errors := 0

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process in mini-batches of 10 for parallel RPC calls
	for i := 0; i < len(blockNums); i += 10 {
		end := i + 10
		if end > len(blockNums) {
			end = len(blockNums)
		}
		miniBatch := blockNums[i:end]

		for _, blockNum := range miniBatch {
			wg.Add(1)
			go func(bn int64) {
				defer wg.Done()

				timestamp, _, err := fetchBlockTimestamp(client, rpcURL, bn)
				if err != nil {
					mu.Lock()
					errors++
					mu.Unlock()
					return
				}

				mu.Lock()
				results = append(results, BlockUpdate{
					Number:    bn,
					Timestamp: timestamp,
				})
				mu.Unlock()
			}(blockNum)
		}
		wg.Wait()
	}

	return results, errors
}

// Batch Kafka produce
func produceKafkaBatch(ctx context.Context, client *kgo.Client, updates []BlockUpdate) int {
	records := make([]*kgo.Record, 0, len(updates))

	for _, update := range updates {
		data, err := json.Marshal(update)
		if err != nil {
			log.Printf("Error marshaling update: %v", err)
			continue
		}

		records = append(records, &kgo.Record{
			Key:   []byte(fmt.Sprintf("%d", update.Number)),
			Value: data,
		})
	}

	if len(records) == 0 {
		return 0
	}

	// Produce batch
	results := client.ProduceSync(ctx, records...)

	successCount := 0
	for _, result := range results {
		if result.Err != nil {
			log.Printf("Error producing to Kafka: %v", result.Err)
		} else {
			successCount++
		}
	}

	return successCount
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
