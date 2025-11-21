package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	workers := flag.Int("workers", 100, "Number of parallel RPC workers")
	startBlock := flag.Int64("start-block", 0, "Start block number (optional)")
	endBlock := flag.Int64("end-block", -1, "End block number (optional, -1 for all)")
	outputFile := flag.String("output", "", "Optional: write to local file instead of Kafka")
	maxBlocks := flag.Int64("max-blocks", 1000000, "Maximum blocks to fetch (for testing)")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting FAST timestamp fetch")
	log.Printf("RPC: %s", *rpcURL)
	log.Printf("Workers: %d", *workers)
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

	// Setup output (Kafka or file)
	var kafkaClient *kgo.Client
	var fileWriter *bufio.Writer
	var outputMu sync.Mutex

	if *outputFile != "" {
		// File output
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer f.Close()
		fileWriter = bufio.NewWriterSize(f, 1024*1024) // 1MB buffer
		defer fileWriter.Flush()
		log.Printf("Writing to file: %s", *outputFile)
	} else {
		// Kafka output with async pattern
		brokers := strings.Split(*kafkaBrokers, ",")
		kafkaClient, err = kgo.NewClient(
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
	}

	// Channels
	blockChan := make(chan int64, *workers*10)
	resultChan := make(chan BlockUpdate, *workers*10)

	// Stats
	var (
		totalFetched   int64
		totalProduced  int64
		totalErrors    int64
		startTime      = time.Now()
	)

	// Start producer goroutine (async Kafka or file writer)
	var producerWg sync.WaitGroup
	producerWg.Add(1)
	go func() {
		defer producerWg.Done()

		for update := range resultChan {
			if fileWriter != nil {
				// Write to file
				data, _ := json.Marshal(update)
				outputMu.Lock()
				fileWriter.Write(data)
				fileWriter.WriteByte('\n')
				outputMu.Unlock()
				atomic.AddInt64(&totalProduced, 1)
			} else {
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
		}

		// Final flush for Kafka
		if kafkaClient != nil {
			kafkaClient.Flush(ctx)
		}
	}()

	// Start RPC workers - MAXIMIZE PARALLELISM
	var fetchWg sync.WaitGroup
	httpClients := make([]*http.Client, *workers)
	for i := 0; i < *workers; i++ {
		httpClients[i] = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		}

		fetchWg.Add(1)
		go func(workerID int, client *http.Client) {
			defer fetchWg.Done()

			for blockNum := range blockChan {
				timestamp, _, err := fetchBlockTimestamp(client, *rpcURL, blockNum)
				if err != nil {
					atomic.AddInt64(&totalErrors, 1)
					continue
				}

				atomic.AddInt64(&totalFetched, 1)
				resultChan <- BlockUpdate{
					Number:    blockNum,
					Timestamp: timestamp,
				}
			}
		}(i, httpClients[i])
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

	// Feed blocks in reverse order
	go func() {
		for blockNum := maxBlock; blockNum >= minBlock; blockNum-- {
			blockChan <- blockNum
		}
		close(blockChan)
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

	resp, err := client.Post(rpcURL, "application/json", strings.NewReader(string(reqBytes)))
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
