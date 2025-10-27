package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	w3 "github.com/lmittmann/w3"
	eth "github.com/lmittmann/w3/module/eth"
	w3types "github.com/lmittmann/w3/w3types"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Decoded struct {
	Name string      `json:"name"`
	Sig  string      `json:"sig"`
	Args interface{} `json:"args"`
}

type IndexBlock struct {
	Number           int64   `json:"number"`
	CoinbaseAddress  string  `json:"coinbase_address"`
	Hash             string  `json:"hash"`
	ParentHash       string  `json:"parent_hash"`
	Nonce            string  `json:"nonce"`
	Sha3Uncles       string  `json:"sha3_uncles"`
	LogsBloom        string  `json:"logs_bloom"`
	TransactionsRoot string  `json:"transactions_root"`
	StateRoot        string  `json:"state_root"`
	ReceiptsRoot     string  `json:"receipts_root"`
	Miner            string  `json:"miner"`
	Difficulty       string  `json:"difficulty"`
	ExtraData        string  `json:"extra_data"`
	Size             int64   `json:"size"`
	GasLimit         int64   `json:"gas_limit"`
	GasUsed          int64   `json:"gas_used"`
	Timestamp        int64   `json:"timestamp"`
	BaseFeePerGas    *string `json:"base_fee_per_gas,omitempty"`
	BlobGasUsed      *int64  `json:"blob_gas_used,omitempty"`
	ExcessBlobGas    *int64  `json:"excess_blob_gas,omitempty"`
	WithdrawalsRoot  *string `json:"withdrawals_root,omitempty"`
	RequestsHash     *string `json:"requests_hash,omitempty"`
	MixHash          *string `json:"mix_hash,omitempty"`
	TxCount          int     `json:"tx_count"`
}

type IndexTx struct {
	Hash                 string  `json:"hash"`
	Nonce                uint64  `json:"nonce"`
	BlockNumber          int64   `json:"block_number"`
	BlockHash            string  `json:"block_hash"`
	TxIndex              int     `json:"tx_index"`
	From                 string  `json:"from_address"`
	To                   *string `json:"to_address,omitempty"`
	Value                string  `json:"value"`
	Gas                  int64   `json:"gas"`
	GasPrice             *string `json:"gas_price,omitempty"`
	MaxPriorityFeePerGas *string `json:"max_priority_fee_per_gas,omitempty"`
	MaxFeePerGas         *string `json:"max_fee_per_gas,omitempty"`
	EffectiveGasPrice    *string `json:"effective_gas_price,omitempty"`
	Input                string  `json:"input"`
	Type                 uint8   `json:"type"`
	ChainID              *int64  `json:"chain_id,omitempty"`
	AccessListJSON       *string `json:"access_list_json,omitempty"`
	BlobGas              *int64  `json:"blob_gas,omitempty"`
	BlobGasFeeCap        *string `json:"blob_gas_fee_cap,omitempty"`
	BlobHashesJSON       *string `json:"blob_hashes_json,omitempty"`
	V                    *string `json:"v,omitempty"`
	R                    *string `json:"r,omitempty"`
	S                    *string `json:"s,omitempty"`
	DecodedJSON          string  `json:"decoded,omitempty"`
}

type IndexReceipt struct {
	TxHash            string  `json:"tx_hash"`
	Status            uint64  `json:"status"`
	CumulativeGasUsed int64   `json:"cumulative_gas_used"`
	GasUsed           int64   `json:"gas_used"`
	ContractAddress   *string `json:"contract_address,omitempty"`
	LogsBloom         string  `json:"logs_bloom"`
	Type              uint8   `json:"type"`
	BlobGasUsed       uint64  `json:"blob_gas_used"`
	BlobGasPrice      *string `json:"blob_gas_price,omitempty"`
}

type IndexLog struct {
	TxHash         string  `json:"tx_hash"`
	LogIndex       int     `json:"log_index"`
	Address        string  `json:"address"`
	BlockNumber    *int64  `json:"block_number,omitempty"`
	BlockHash      *string `json:"block_hash,omitempty"`
	TxIndex        int     `json:"tx_index"`
	BlockTimestamp *int64  `json:"block_timestamp,omitempty"`
	TopicsJSON     string  `json:"topics"`
	Data           string  `json:"data"`
	Removed        bool    `json:"removed"`
	DecodedJSON    string  `json:"decoded,omitempty"`
}

var abiCache = make(map[string]abi.ABI)

const maxRowsPerInsert = 10_000

type kafkaOptions struct {
	enabled   bool
	chainName string
	client    *kgo.Client
}

func isDynamicType(t *abi.Type) bool {
	switch t.T {
	case abi.TupleTy:
		for _, elem := range t.TupleElems {
			if isDynamicType(elem) {
				return true
			}
		}
		return false
	case abi.SliceTy:
		return true
	case abi.ArrayTy:
		return t.Size == 0 || isDynamicType(t.Elem)
	case abi.StringTy, abi.BytesTy:
		return true
	default:
		return false
	}
}

func ptrString(s string) *string {
	return &s
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	rpcURL := flag.String("rpc", "http://localhost:8545", "Ethereum RPC URL")
	dsn := flag.String("dsn", "root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true", "StarRocks DSN")
	mode := flag.String("mode", "dual", "Mode: dual, forward, or backfill")
	fromBlock := flag.Int64("from", 0, "Starting block for backfill (lowest)")
	toBlock := flag.Int64("to", 0, "Ending block for backfill (highest)")
	pollInterval := flag.Duration("poll", 100*time.Millisecond, "Poll interval for forward mode")
	batchSize := flag.Int("batch-size", 50, "Batch size for RPC calls")
	flushInterval := flag.Duration("flush-interval", 2*time.Second, "Maximum time to buffer rows before flushing to the database")
	startBlock := flag.Int64("start-block", 0, "Starting block for forward mode when no history exists")
	abiDir := flag.String("abi-dir", "./contracts-abi/abi", "Directory containing ABI files")
	abiConfig := flag.String("abi-config", "", "Optional path to ABI manifest JSON")
	enableKafka := flag.Bool("enable-kafka", true, "Enable Kafka output")
	kafkaBrokers := flag.String("kafka-brokers", "mevcommit-message-queue-cluster-kafka-bootstrap.kafka.svc:9092", "Comma-separated list of Kafka brokers")
	chainName := flag.String("chain-name", "eth-mainnet", "Chain name to use as topic prefix (e.g., 'mevcommit', 'eth-mainnet')")
	flag.Parse()

	// Validate inputs
	if *mode != "forward" && *mode != "backfill" && *mode != "dual" {
		logger.Error("Invalid mode", "mode", *mode)
		os.Exit(1)
	}
	if *mode == "backfill" && (*fromBlock >= *toBlock || *toBlock == 0) {
		logger.Error("Invalid backfill range", "from", *fromBlock, "to", *toBlock)
		os.Exit(1)
	}
	if *mode == "dual" && *startBlock == 0 {
		logger.Error("Dual mode requires --start-block to be set (target for backward indexer)")
		os.Exit(1)
	}
	if *enableKafka && strings.TrimSpace(*chainName) == "" {
		logger.Error("Chain name cannot be empty when Kafka is enabled")
		os.Exit(1)
	}
	if *enableKafka && strings.TrimSpace(*kafkaBrokers) == "" {
		logger.Error("Kafka brokers cannot be empty when Kafka is enabled")
		os.Exit(1)
	}

	client, err := w3.Dial(*rpcURL)
	if err != nil {
		logger.Error("Failed to connect to Ethereum RPC", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	// Parse DSN to extract database name and create database if not exists
	cfg, err := mysql.ParseDSN(*dsn)
	if err != nil {
		logger.Error("Failed to parse DSN", "err", err)
		os.Exit(1)
	}
	dbName := cfg.DBName
	if dbName != "" {
		cfg.DBName = ""
		baseDSN := cfg.FormatDSN()
		dbTemp, err := sql.Open("mysql", baseDSN)
		if err != nil {
			logger.Error("Failed to connect to server without database", "err", err)
			os.Exit(1)
		}
		defer dbTemp.Close()

		_, err = dbTemp.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
		if err != nil {
			logger.Error("Failed to create database", "dbName", dbName, "err", err)
			os.Exit(1)
		}
	}

	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		logger.Error("Failed to connect to StarRocks", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create tables if they do not exist (always done regardless of output mode)
	if err := createTables(db); err != nil {
		logger.Error("Failed to create tables", "err", err)
		os.Exit(1)
	}

	// Initialize Kafka client if enabled
	var kafkaOpts kafkaOptions
	if *enableKafka {
		brokers := strings.Split(*kafkaBrokers, ",")
		// Trim whitespace from broker addresses
		for i := range brokers {
			brokers[i] = strings.TrimSpace(brokers[i])
		}

		client, err := kgo.NewClient(
			kgo.SeedBrokers(brokers...),
			kgo.RequiredAcks(kgo.AllISRAcks()),           // Wait for all replicas (durability)
			kgo.ProducerBatchMaxBytes(10_485_760),        // 10MB batch size (matches topic max)
			kgo.ProducerBatchCompression(kgo.ZstdCompression()), // Zstd compression for best ratio
			kgo.ProducerLinger(100*time.Millisecond),     // Batch records for throughput
			kgo.RequestRetries(3),                        // Retry on transient failures
			kgo.RetryBackoffFn(func(attempts int) time.Duration {
				return time.Duration(attempts) * 100 * time.Millisecond
			}),
		)
		if err != nil {
			logger.Error("Failed to create Kafka client", "err", err)
			os.Exit(1)
		}

		// Create topics if they don't exist
		topics := []string{
			fmt.Sprintf("%s-blocks", *chainName),
			fmt.Sprintf("%s-transactions", *chainName),
			fmt.Sprintf("%s-receipts", *chainName),
			fmt.Sprintf("%s-logs", *chainName),
		}
		adminClient := kadm.NewClient(client)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		// Topic configuration: allow large messages (10MB for blob transactions)
		// IMPORTANT: Kafka broker must also have message.max.bytes >= 10MB
		// To update broker config:
		//   kubectl edit configmap kafka-config -n kafka
		//   Add: message.max.bytes=10485760
		//   Restart Kafka pods
		topicConfigs := map[string]*string{
			"max.message.bytes": ptrString("10485760"), // 10MB (handles large calldata/blob txs)
			"retention.ms":      ptrString("604800000"), // 7 days
		}

		for _, topic := range topics {
			_, err := adminClient.CreateTopic(ctx, 1, -1, topicConfigs, topic)
			if err != nil {
				logger.Warn("Failed to create topic (may already exist)", "topic", topic, "err", err)
			}
		}

		logger.Info("Kafka configured for large messages",
			"max_message_bytes", "10MB",
			"compression", "zstd",
			"note", "Ensure Kafka broker message.max.bytes >= 10MB")
		cancel()

		kafkaOpts = kafkaOptions{
			enabled:   true,
			chainName: *chainName,
			client:    client,
		}

		// Graceful shutdown handler
		defer func() {
			logger.Info("Shutting down Kafka client...")
			flushCtx, flushCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer flushCancel()
			if err := client.Flush(flushCtx); err != nil {
				logger.Error("Failed to flush Kafka client", "err", err)
			}
			client.Close()
			logger.Info("Kafka client closed")
		}()
	}

	// Load ABIs
	if err := loadABIs(db, *abiDir, *abiConfig); err != nil {
		logger.Error("Failed to load ABIs", "err", err)
		// Continue or exit based on preference; here continue
	}

	// Log startup configuration
	outputMode := "sql"
	if *enableKafka {
		outputMode = "kafka"
	}
	logger.Info("Starting indexer",
		"mode", *mode,
		"output", outputMode,
		"rpc", *rpcURL,
		"database", dbName,
		"chain_name", *chainName,
		"batch_size", *batchSize,
		"flush_interval", *flushInterval,
	)
	if *enableKafka {
		logger.Info("Kafka configuration", "brokers", *kafkaBrokers, "topics_prefix", *chainName)
	}

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	if *mode == "dual" {
		// Run both forward and backward indexers concurrently
		logger.Info("Starting dual-direction indexer",
			"mode", "dual",
			"forward", "tracks latest blocks",
			"backward", fmt.Sprintf("backfills to block %d", *startBlock))

		var wg sync.WaitGroup
		wg.Add(2)

		// Forward indexer: keeps up with chain head (priority)
		go func() {
			defer wg.Done()
			forwardIndexFromLatest(ctx, client, db, *pollInterval, *batchSize, *flushInterval, kafkaOpts, logger)
		}()

		// Backward indexer: backfills historical data
		go func() {
			defer wg.Done()
			backwardIndexToStart(ctx, client, db, *startBlock, *batchSize, *flushInterval, kafkaOpts, logger)
		}()

		wg.Wait()
		logger.Info("Dual-direction indexer stopped")
	} else if *mode == "forward" {
		forwardIndex(ctx, client, db, *pollInterval, *batchSize, *flushInterval, kafkaOpts, *startBlock, logger)
	} else {
		backfillIndex(ctx, client, db, *fromBlock, *toBlock, *batchSize, *flushInterval, kafkaOpts, logger)
	}
}

func loadABIs(db *sql.DB, abiDir, abiConfig string) error {
	if abiConfig == "" {
		slog.Warn("No ABI manifest provided; contract ABI decoding will be skipped")
		return nil
	}

	return loadABIsFromManifest(db, abiDir, abiConfig)
}

func loadABIsFromManifest(db *sql.DB, abiDir, abiConfig string) error {
	configBytes, err := os.ReadFile(abiConfig)
	if err != nil {
		return err
	}

	var manifest struct {
		Contracts []struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			ABIPath string `json:"abi_path"`
		} `json:"contracts"`
	}
	if err := json.Unmarshal(configBytes, &manifest); err != nil {
		return err
	}

	for _, entry := range manifest.Contracts {
		if entry.Address == "" || entry.ABIPath == "" {
			slog.Error("Skipping ABI manifest entry, missing address or abi_path", "entry", entry)
			continue
		}

		path := entry.ABIPath
		if !filepath.IsAbs(path) {
			path = filepath.Join(abiDir, path)
		}

		abiBytes, err := os.ReadFile(path)
		if err != nil {
			slog.Error("Failed to read ABI file from manifest", "path", path, "err", err)
			continue
		}
		abiStr := string(abiBytes)

		name := entry.Name
		if name == "" {
			base := filepath.Base(entry.ABIPath)
			name = strings.TrimSuffix(base, filepath.Ext(base))
		}

		address := strings.ToLower(entry.Address)
		if _, err := db.Exec("INSERT INTO contract_abis (address, name, abi) VALUES (?, ?, parse_json(?))", address, name, abiStr); err != nil {
			slog.Error("Failed to insert ABI from manifest", "address", address, "err", err)
			continue
		}

		abiObj, err := abi.JSON(strings.NewReader(abiStr))
		if err != nil {
			slog.Error("Failed to parse ABI from manifest", "address", address, "err", err)
			continue
		}
		abiCache[address] = abiObj
	}

	return nil
}

func createTables(db *sql.DB) error {
	createStatements := []string{
		`CREATE TABLE IF NOT EXISTS blocks (
			number BIGINT,
			coinbase_address VARCHAR(255),
			hash VARCHAR(255),
			parent_hash VARCHAR(255),
			nonce VARCHAR(255),
			sha3_uncles VARCHAR(255),
			logs_bloom VARCHAR(65535),
			transactions_root VARCHAR(255),
			state_root VARCHAR(255),
			receipts_root VARCHAR(255),
			miner VARCHAR(255),
			difficulty VARCHAR(255),
			extra_data VARCHAR(65535),
			size BIGINT,
			gas_limit BIGINT,
			gas_used BIGINT,
			timestamp BIGINT,
			base_fee_per_gas VARCHAR(255),
			blob_gas_used BIGINT,
			excess_blob_gas BIGINT,
			withdrawals_root VARCHAR(255),
			requests_hash VARCHAR(255),
			mix_hash VARCHAR(255),
			tx_count INT
		) 
		ENGINE=olap 
		PRIMARY KEY(number) 
		DISTRIBUTED BY HASH(number) BUCKETS 1 
		PROPERTIES("replication_num"="1")`,

		`CREATE TABLE IF NOT EXISTS transactions (
			hash VARCHAR(255),
    		nonce BIGINT,
			block_number BIGINT,
			block_hash VARCHAR(255),
			tx_index INT,
			from_address VARCHAR(255),
			to_address VARCHAR(255),
			value VARCHAR(255),
			gas BIGINT,
			gas_price VARCHAR(255),
			max_priority_fee_per_gas VARCHAR(255),
			max_fee_per_gas VARCHAR(255),
			effective_gas_price VARCHAR(255),
			input VARBINARY(1048576),
			type TINYINT,
			chain_id BIGINT,
			access_list_json JSON,
			blob_gas BIGINT,
			blob_gas_fee_cap VARCHAR(255),
			blob_hashes_json JSON,
			v VARCHAR(255),
			r VARCHAR(255),
			s VARCHAR(255),
			decoded JSON
		) 
		ENGINE=olap 
		PRIMARY KEY(hash) 
		DISTRIBUTED BY HASH(hash) BUCKETS 1 
		PROPERTIES("replication_num"="1")`,

		`CREATE TABLE IF NOT EXISTS receipts (
			tx_hash VARCHAR(255),
			status BIGINT,
			cumulative_gas_used BIGINT,
			gas_used BIGINT,
			contract_address VARCHAR(255),
			logs_bloom VARCHAR(65535),
			type TINYINT,
			blob_gas_used BIGINT,
			blob_gas_price VARCHAR(255)
		) 
		ENGINE=olap 
		PRIMARY KEY(tx_hash) 
		DISTRIBUTED BY HASH(tx_hash) BUCKETS 1 
		PROPERTIES("replication_num"="1")`,

		`CREATE TABLE IF NOT EXISTS logs (
			tx_hash VARCHAR(255),
			log_index INT,
			address VARCHAR(255),
			block_number BIGINT,
			block_hash VARCHAR(255),
			tx_index INT,
			block_timestamp BIGINT,
			topics JSON,
			data VARBINARY(1048576),
			removed BOOLEAN,
			decoded JSON
		) 
		ENGINE=olap 
		PRIMARY KEY(tx_hash, log_index) 
		DISTRIBUTED BY HASH(tx_hash) BUCKETS 1 
		PROPERTIES("replication_num"="1")`,

		`CREATE TABLE IF NOT EXISTS contract_abis (
			address VARCHAR(255),
			name VARCHAR(255),
			abi JSON
		) 
		ENGINE=olap 
		PRIMARY KEY(address) 
		DISTRIBUTED BY HASH(address) BUCKETS 1 
		PROPERTIES("replication_num"="1")`,
	}

	for _, stmt := range createStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func forwardIndex(ctx context.Context, client *w3.Client, db *sql.DB, pollInterval time.Duration, batchSize int, flushInterval time.Duration, kafkaOpts kafkaOptions, startBlock int64, logger *slog.Logger) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Forward indexing stopped due to context cancellation")
			return
		default:
		}
		var latest *big.Int
		if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
			logger.Error("Failed to get latest block", "err", err)
			time.Sleep(pollInterval)
			continue
		}

		// Get checkpoint from Kafka if enabled, otherwise use SQL
		var maxIndexed int64
		var hasIndexed bool
		var err error
		if kafkaOpts.enabled {
			maxIndexed, hasIndexed, err = getKafkaCheckpoint(kafkaOpts)
			if err != nil {
				logger.Error("Failed to get Kafka checkpoint", "err", err)
				time.Sleep(pollInterval)
				continue
			}
		} else {
			maxIndexed, hasIndexed, err = getMaxIndexedBlock(db)
			if err != nil {
				logger.Error("Failed to get max indexed block from database", "err", err)
				time.Sleep(pollInterval)
				continue
			}
		}

		start := startBlock
		if hasIndexed {
			start = maxIndexed + 1
			if start < startBlock {
				start = startBlock
			}
		}

		for start <= latest.Int64() {
			end := start + int64(batchSize) - 1
			if end > latest.Int64() {
				end = latest.Int64()
			}
			blockNums := make([]int64, 0, end-start+1)
			for i := start; i <= end; i++ {
				blockNums = append(blockNums, i)
			}
			if err := processBatch(client, db, blockNums, flushInterval, kafkaOpts); err != nil {
				logger.Error("Failed to process batch", "from", start, "to", end, "err", err)
				logLoadTrackingDetails(logger, db, err)
				break
			}
			start = end + 1
		}

		time.Sleep(pollInterval)
	}
}

// forwardIndexFromLatest continuously indexes new blocks from the latest checkpoint forward
// This runs as a goroutine and keeps up with the chain head
func forwardIndexFromLatest(ctx context.Context, client *w3.Client, db *sql.DB, pollInterval time.Duration, batchSize int, flushInterval time.Duration, kafkaOpts kafkaOptions, logger *slog.Logger) {
	// Get checkpoint range from Kafka
	_, maxBlock, hasData, err := getKafkaCheckpointRange(kafkaOpts)
	if err != nil {
		logger.Error("Forward indexer: Failed to get Kafka checkpoint", "direction", "forward", "err", err)
		return
	}

	var startBlock int64
	if !hasData {
		// No data in Kafka, start from current chain head
		var latest *big.Int
		if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
			logger.Error("Forward indexer: Failed to get latest block", "direction", "forward", "err", err)
			return
		}
		startBlock = latest.Int64()
		logger.Info("Forward indexer: Starting from chain head (no checkpoint)", "direction", "forward", "start_block", startBlock)
	} else {
		// Resume from max block + 1
		startBlock = maxBlock + 1
		logger.Info("Forward indexer: Resuming from checkpoint", "direction", "forward", "checkpoint", maxBlock, "start_block", startBlock)
	}

	next := startBlock
	var pending []int64

	for {
		select {
		case <-ctx.Done():
			logger.Info("Forward indexer stopped", "direction", "forward")
			return
		default:
		}

		var latest *big.Int
		if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
			logger.Error("Forward indexer: Failed to get latest block", "direction", "forward", "err", err)
			time.Sleep(pollInterval)
			continue
		}

		// Collect blocks up to latest
		for next <= latest.Int64() {
			pending = append(pending, next)
			next++

			if len(pending) >= batchSize {
				if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
					logger.Error("Forward indexer: Failed to process batch", "direction", "forward", "from", pending[0], "to", pending[len(pending)-1], "err", err)
					logLoadTrackingDetails(logger, db, err)
					// Exit on error to avoid data loss
					return
				}
				logger.Info("Forward indexer: Processed batch", "direction", "forward", "from", pending[0], "to", pending[len(pending)-1], "blocks", len(pending))
				pending = pending[:0]
			}
		}

		// Flush remaining blocks
		if len(pending) > 0 {
			if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
				logger.Error("Forward indexer: Failed to process batch", "direction", "forward", "from", pending[0], "to", pending[len(pending)-1], "err", err)
				logLoadTrackingDetails(logger, db, err)
				return
			}
			logger.Info("Forward indexer: Processed batch", "direction", "forward", "from", pending[0], "to", pending[len(pending)-1], "blocks", len(pending))
			pending = pending[:0]
		}

		time.Sleep(pollInterval)
	}
}

// backwardIndexToStart indexes historical blocks from latest checkpoint backward to start block
// This runs as a goroutine and backfills historical data
func backwardIndexToStart(ctx context.Context, client *w3.Client, db *sql.DB, startBlock int64, batchSize int, flushInterval time.Duration, kafkaOpts kafkaOptions, logger *slog.Logger) {
	// Get checkpoint range from Kafka
	minBlock, _, hasData, err := getKafkaCheckpointRange(kafkaOpts)
	if err != nil {
		logger.Error("Backward indexer: Failed to get Kafka checkpoint", "direction", "backward", "err", err)
		return
	}

	var backfillFrom int64
	if !hasData {
		// No data in Kafka, start from current chain head - 1
		var latest *big.Int
		if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
			logger.Error("Backward indexer: Failed to get latest block", "direction", "backward", "err", err)
			return
		}
		backfillFrom = latest.Int64() - 1
		logger.Info("Backward indexer: Starting from chain head (no checkpoint)", "direction", "backward", "start_from", backfillFrom, "target", startBlock)
	} else {
		// Resume from min block - 1
		backfillFrom = minBlock - 1
		logger.Info("Backward indexer: Resuming from checkpoint", "direction", "backward", "checkpoint", minBlock, "start_from", backfillFrom, "target", startBlock)
	}

	// Check if already reached target
	if backfillFrom < startBlock {
		logger.Info("Backward indexer: Already reached start block, nothing to backfill", "direction", "backward", "current", backfillFrom, "target", startBlock)
		return
	}

	var pending []int64
	for i := backfillFrom; i >= startBlock; i-- {
		select {
		case <-ctx.Done():
			logger.Info("Backward indexer stopped", "direction", "backward", "reached_block", i)
			return
		default:
		}

		pending = append(pending, i)

		if len(pending) >= batchSize {
			// Sort in ascending order for processing
			sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })

			if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
				logger.Error("Backward indexer: Failed to process batch", "direction", "backward", "from", pending[0], "to", pending[len(pending)-1], "err", err)
				logLoadTrackingDetails(logger, db, err)
				// Exit on error to avoid data loss
				return
			}
			logger.Info("Backward indexer: Processed batch", "direction", "backward", "from", pending[0], "to", pending[len(pending)-1], "blocks", len(pending), "remaining", i-startBlock)
			pending = pending[:0]
		}
	}

	// Flush remaining blocks
	if len(pending) > 0 {
		sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })

		if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
			logger.Error("Backward indexer: Failed to process batch", "direction", "backward", "from", pending[0], "to", pending[len(pending)-1], "err", err)
			logLoadTrackingDetails(logger, db, err)
			return
		}
		logger.Info("Backward indexer: Processed batch", "direction", "backward", "from", pending[0], "to", pending[len(pending)-1], "blocks", len(pending))
	}

	logger.Info("Backward indexer: Completed backfill to start block", "direction", "backward", "target", startBlock)
}

func backfillIndex(ctx context.Context, client *w3.Client, db *sql.DB, from, to int64, batchSize int, flushInterval time.Duration, kafkaOpts kafkaOptions, logger *slog.Logger) {
	var pending []int64
	for i := to; i >= from; i-- {
		select {
		case <-ctx.Done():
			logger.Info("Backfill indexing stopped due to context cancellation")
			return
		default:
		}

		// In Kafka mode, skip existence check (Kafka will deduplicate by key)
		// In SQL mode, check if block exists to avoid reprocessing
		shouldProcess := kafkaOpts.enabled
		if !kafkaOpts.enabled {
			exists, err := blockExists(db, i)
			if err != nil {
				logger.Error("Failed to check if block exists", "block", i, "err", err)
				continue
			}
			shouldProcess = !exists
		}

		if shouldProcess {
			pending = append(pending, i)
			if len(pending) == batchSize {
				sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })
				if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
					logger.Error("Failed to process batch", "err", err)
					logLoadTrackingDetails(logger, db, err)
				}
				pending = pending[:0]
			}
		}
	}
	if len(pending) > 0 {
		sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })
		if err := processBatch(client, db, pending, flushInterval, kafkaOpts); err != nil {
			logger.Error("Failed to process batch", "err", err)
			logLoadTrackingDetails(logger, db, err)
		}
	}
}

func processBatch(client *w3.Client, db *sql.DB, blockNums []int64, flushInterval time.Duration, kafkaOpts kafkaOptions) error {
	startBatch := time.Now()
	n := len(blockNums)
	blocks := make([]*types.Block, n)
	receipts := make([]types.Receipts, n)
	var calls []w3types.RPCCaller
	for i := range blockNums {
		num := big.NewInt(blockNums[i])
		calls = append(calls, eth.BlockByNumber(num).Returns(&blocks[i]))
		calls = append(calls, eth.BlockReceipts(num).Returns(&receipts[i]))
	}
	if err := client.Call(calls...); err != nil {
		return err
	}
	fetchDuration := time.Since(startBatch)
	slog.Info("Fetched batch", "blocks", len(blockNums), "duration", fetchDuration)

	if kafkaOpts.enabled {
		// Kafka mode: process and send immediately (no chunking)
		return processBatchKafka(blocks, receipts, kafkaOpts, db)
	} else {
		// SQL mode: use chunking for efficient batch inserts
		return processBatchSQL(blocks, receipts, db, flushInterval)
	}
}

func batchInsertBlocks(tx *sql.Tx, blocks []IndexBlock) error {
	if len(blocks) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(blocks))
	valueArgs := make([]interface{}, 0, len(blocks)*24)
	for _, b := range blocks {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, b.Number, b.CoinbaseAddress, b.Hash, b.ParentHash, b.Nonce, b.Sha3Uncles, b.LogsBloom, b.TransactionsRoot, b.StateRoot, b.ReceiptsRoot, b.Miner, b.Difficulty, b.ExtraData, b.Size, b.GasLimit, b.GasUsed, b.Timestamp, b.BaseFeePerGas, b.BlobGasUsed, b.ExcessBlobGas, b.WithdrawalsRoot, b.RequestsHash, b.MixHash, b.TxCount)
	}
	stmt := fmt.Sprintf("INSERT INTO blocks (number, coinbase_address, hash, parent_hash, nonce, sha3_uncles, logs_bloom, transactions_root, state_root, receipts_root, miner, difficulty, extra_data, size, gas_limit, gas_used, timestamp, base_fee_per_gas, blob_gas_used, excess_blob_gas, withdrawals_root, requests_hash, mix_hash, tx_count) VALUES %s", strings.Join(valueStrings, ", "))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}

func batchInsertTxs(tx *sql.Tx, txs []IndexTx) error {
	if len(txs) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(txs))
	valueArgs := make([]interface{}, 0, len(txs)*23)
	for _, t := range txs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, t.Hash, t.Nonce, t.BlockNumber, t.BlockHash, t.TxIndex, t.From, t.To, t.Value, t.Gas, t.GasPrice, t.MaxPriorityFeePerGas, t.MaxFeePerGas, t.EffectiveGasPrice, t.Input, t.Type, t.ChainID, t.AccessListJSON, t.BlobGas, t.BlobGasFeeCap, t.BlobHashesJSON, t.V, t.R, t.S, t.DecodedJSON)
	}
	stmt := fmt.Sprintf("INSERT INTO transactions (hash, nonce, block_number, block_hash, tx_index, from_address, to_address, value, gas, gas_price, max_priority_fee_per_gas, max_fee_per_gas, effective_gas_price, input, type, chain_id, access_list_json, blob_gas, blob_gas_fee_cap, blob_hashes_json, v, r, s, decoded) VALUES %s", strings.Join(valueStrings, ", "))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}

func batchInsertReceipts(tx *sql.Tx, receipts []IndexReceipt) error {
	if len(receipts) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(receipts))
	valueArgs := make([]interface{}, 0, len(receipts)*9)
	for _, r := range receipts {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, r.TxHash, r.Status, r.CumulativeGasUsed, r.GasUsed, r.ContractAddress, r.LogsBloom, r.Type, r.BlobGasUsed, r.BlobGasPrice)
	}
	stmt := fmt.Sprintf("INSERT INTO receipts (tx_hash, status, cumulative_gas_used, gas_used, contract_address, logs_bloom, type, blob_gas_used, blob_gas_price) VALUES %s", strings.Join(valueStrings, ", "))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}

func batchInsertLogs(tx *sql.Tx, logs []IndexLog) error {
	if len(logs) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(logs))
	valueArgs := make([]interface{}, 0, len(logs)*11)
	for _, l := range logs {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, l.TxHash, l.LogIndex, l.Address, l.BlockNumber, l.BlockHash, l.TxIndex, l.BlockTimestamp, l.TopicsJSON, l.Data, l.Removed, l.DecodedJSON)
	}
	stmt := fmt.Sprintf("INSERT INTO logs (tx_hash, log_index, address, block_number, block_hash, tx_index, block_timestamp, topics, data, removed, decoded) VALUES %s", strings.Join(valueStrings, ", "))
	_, err := tx.Exec(stmt, valueArgs...)
	return err
}

func produceToKafka(opts kafkaOptions, blocks []IndexBlock, txs []IndexTx, receipts []IndexReceipt, logs []IndexLog) error {
	// Timing: JSON marshaling phase
	marshalStart := time.Now()

	// Collect all records to produce
	var records []*kgo.Record

	// Add blocks
	for _, block := range blocks {
		data, err := json.Marshal(block)
		if err != nil {
			return fmt.Errorf("failed to marshal block %d: %w", block.Number, err)
		}
		records = append(records, &kgo.Record{
			Key:   []byte(strconv.FormatInt(block.Number, 10)),
			Topic: fmt.Sprintf("%s-blocks", opts.chainName),
			Value: data,
		})
	}

	// Add transactions
	for _, tx := range txs {
		data, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction %s: %w", tx.Hash, err)
		}

		// Log if transaction is unusually large
		if len(data) > 1_000_000 {
			slog.Info("Large transaction detected",
				"tx_hash", tx.Hash,
				"size_bytes", len(data),
				"block", tx.BlockNumber)
		}

		records = append(records, &kgo.Record{
			Key:   []byte(tx.Hash),
			Topic: fmt.Sprintf("%s-transactions", opts.chainName),
			Value: data,
		})
	}

	// Add receipts
	for _, receipt := range receipts {
		data, err := json.Marshal(receipt)
		if err != nil {
			return fmt.Errorf("failed to marshal receipt %s: %w", receipt.TxHash, err)
		}
		records = append(records, &kgo.Record{
			Key:   []byte(receipt.TxHash),
			Topic: fmt.Sprintf("%s-receipts", opts.chainName),
			Value: data,
		})
	}

	// Add logs
	for _, log := range logs {
		data, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("failed to marshal log %s:%d: %w", log.TxHash, log.LogIndex, err)
		}
		records = append(records, &kgo.Record{
			Key:   []byte(fmt.Sprintf("%s:%d", log.TxHash, log.LogIndex)),
			Topic: fmt.Sprintf("%s-logs", opts.chainName),
			Value: data,
		})
	}

	if len(records) == 0 {
		return nil
	}

	marshalDuration := time.Since(marshalStart)

	// Retry logic: up to 3 attempts
	const maxRetries = 3
	var totalEnqueueDuration time.Duration
	var totalWaitDuration time.Duration
	var totalBlocksSuccess atomic.Int32
	var totalTxsSuccess atomic.Int32
	var totalReceiptsSuccess atomic.Int32
	var totalLogsSuccess atomic.Int32

	recordsToSend := records
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			slog.Warn("Retrying failed Kafka records",
				"attempt", attempt,
				"max_attempts", maxRetries,
				"records_to_retry", len(recordsToSend))
		}

		// Create new context for each attempt
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		// Timing: Enqueue phase (calling .Produce on all records)
		enqueueStart := time.Now()

		// Produce with callbacks (async pattern from the example)
		var (
			wg                   sync.WaitGroup
			mu                   sync.Mutex
			failedRecords        []*kgo.Record
			successCount         atomic.Int32
			blocksSuccessCount   atomic.Int32
			txsSuccessCount      atomic.Int32
			receiptsSuccessCount atomic.Int32
			logsSuccessCount     atomic.Int32
		)

		wg.Add(len(recordsToSend))

		for _, record := range recordsToSend {
			rec := record

			opts.client.Produce(ctx, rec, func(r *kgo.Record, err error) {
				defer wg.Done()

				if err != nil {
					// Log detailed information about the failed record
					slog.Error("Failed to produce Kafka record",
						"attempt", attempt,
						"topic", r.Topic,
						"key", string(r.Key),
						"value_size_bytes", len(r.Value),
						"error", err.Error(),
					)

					// Store for retry
					mu.Lock()
					failedRecords = append(failedRecords, r)
					mu.Unlock()
				} else {
					successCount.Add(1)

					// Track per-entity-type metrics based on topic suffix
					if strings.HasSuffix(r.Topic, "-blocks") {
						blocksSuccessCount.Add(1)
						totalBlocksSuccess.Add(1)
					} else if strings.HasSuffix(r.Topic, "-transactions") {
						txsSuccessCount.Add(1)
						totalTxsSuccess.Add(1)
					} else if strings.HasSuffix(r.Topic, "-receipts") {
						receiptsSuccessCount.Add(1)
						totalReceiptsSuccess.Add(1)
					} else if strings.HasSuffix(r.Topic, "-logs") {
						logsSuccessCount.Add(1)
						totalLogsSuccess.Add(1)
					}
				}
			})
		}

		enqueueDuration := time.Since(enqueueStart)
		totalEnqueueDuration += enqueueDuration

		// Timing: Wait phase (waiting for all ACKs from Kafka)
		waitStart := time.Now()

		// Wait for all callbacks to complete
		wg.Wait()

		waitDuration := time.Since(waitStart)
		totalWaitDuration += waitDuration

		cancel() // Clean up context

		success := successCount.Load()
		failed := int32(len(failedRecords))
		total := int32(len(recordsToSend))

		slog.Info("Kafka produce attempt result",
			"attempt", attempt,
			"total_records", total,
			"successful", success,
			"failed", failed,
			"blocks_acked", blocksSuccessCount.Load(),
			"txs_acked", txsSuccessCount.Load(),
			"receipts_acked", receiptsSuccessCount.Load(),
			"logs_acked", logsSuccessCount.Load(),
		)

		// If all succeeded, break out
		if failed == 0 {
			break
		}

		// If this was the last attempt, return error
		if attempt == maxRetries {
			// Build detailed error message
			errMsg := fmt.Sprintf("failed to produce %d/%d records to Kafka after %d attempts", failed, total, maxRetries)
			for i, rec := range failedRecords {
				if i < 10 { // Show first 10 failures
					errMsg += fmt.Sprintf("\n  [%d] topic=%s, key=%s, size=%d bytes", i+1, rec.Topic, string(rec.Key), len(rec.Value))
				}
			}
			if len(failedRecords) > 10 {
				errMsg += fmt.Sprintf("\n  ... and %d more failures", len(failedRecords)-10)
			}
			return fmt.Errorf("%s", errMsg)
		}

		// Prepare for retry
		recordsToSend = failedRecords

		// Wait a bit before retrying (exponential backoff)
		backoff := time.Duration(attempt) * 2 * time.Second
		slog.Info("Waiting before retry", "backoff", backoff)
		time.Sleep(backoff)
	}

	// Detailed timing breakdown for produceToKafka
	totalProduceDuration := marshalDuration + totalEnqueueDuration + totalWaitDuration
	slog.Info("produceToKafka timing breakdown",
		"total_duration", totalProduceDuration,
		"marshal_duration", marshalDuration,
		"marshal_pct", fmt.Sprintf("%.1f%%", 100*marshalDuration.Seconds()/totalProduceDuration.Seconds()),
		"enqueue_duration", totalEnqueueDuration,
		"enqueue_pct", fmt.Sprintf("%.1f%%", 100*totalEnqueueDuration.Seconds()/totalProduceDuration.Seconds()),
		"wait_duration", totalWaitDuration,
		"wait_pct", fmt.Sprintf("%.1f%%", 100*totalWaitDuration.Seconds()/totalProduceDuration.Seconds()),
		"total_records", len(records),
	)

	// Calculate and log throughput metrics (only count successful ACKs)
	durationSec := totalWaitDuration.Seconds()
	if durationSec > 0 {
		blocksPerSec := float64(totalBlocksSuccess.Load()) / durationSec
		txsPerSec := float64(totalTxsSuccess.Load()) / durationSec
		receiptsPerSec := float64(totalReceiptsSuccess.Load()) / durationSec
		logsPerSec := float64(totalLogsSuccess.Load()) / durationSec

		slog.Info("Kafka ACK throughput (broker speed)",
			"wait_duration", totalWaitDuration,
			"blocks_acked", totalBlocksSuccess.Load(),
			"txs_acked", totalTxsSuccess.Load(),
			"receipts_acked", totalReceiptsSuccess.Load(),
			"logs_acked", totalLogsSuccess.Load(),
			"blocks_per_sec", fmt.Sprintf("%.2f", blocksPerSec),
			"txs_per_sec", fmt.Sprintf("%.2f", txsPerSec),
			"receipts_per_sec", fmt.Sprintf("%.2f", receiptsPerSec),
			"logs_per_sec", fmt.Sprintf("%.2f", logsPerSec),
		)
	}

	return nil
}

func shouldFlushChunk(blocks []IndexBlock, txs []IndexTx, receipts []IndexReceipt, logs []IndexLog, chunkStart time.Time, flushInterval time.Duration) bool {
	if len(blocks) >= maxRowsPerInsert {
		return true
	}
	if len(txs) >= maxRowsPerInsert {
		return true
	}
	if len(receipts) >= maxRowsPerInsert {
		return true
	}
	if len(logs) >= maxRowsPerInsert {
		return true
	}
	if flushInterval <= 0 {
		flushInterval = 2 * time.Second
	}
	if !chunkStart.IsZero() && time.Since(chunkStart) >= flushInterval {
		return true
	}
	return false
}

func getMaxIndexedBlock(db *sql.DB) (int64, bool, error) {
	var max sql.NullInt64
	err := db.QueryRow("SELECT MAX(number) FROM blocks").Scan(&max)
	if err != nil {
		return 0, false, err
	}
	if max.Valid {
		return max.Int64, true, nil
	}
	return 0, false, nil
}

var loadTrackingJobIDRegexp = regexp.MustCompile(`job_id=([0-9]+)`)

func logLoadTrackingDetails(logger *slog.Logger, db *sql.DB, err error) {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return
	}
	if mysqlErr.Number != 1064 {
		return
	}
	matches := loadTrackingJobIDRegexp.FindStringSubmatch(mysqlErr.Message)
	if len(matches) != 2 {
		return
	}
	jobID := matches[1]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, queryErr := db.QueryContext(ctx, "SELECT tracking_log FROM information_schema.load_tracking_logs WHERE job_id = ?", jobID)
	if queryErr != nil {
		logger.Error("Failed to fetch load tracking log", "job_id", jobID, "err", queryErr)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var logEntry string
		if scanErr := rows.Scan(&logEntry); scanErr != nil {
			logger.Error("Failed to scan load tracking log", "job_id", jobID, "err", scanErr)
			return
		}
		logger.Error("StarRocks load detail", "job_id", jobID, "log", logEntry)
	}
	if err := rows.Err(); err != nil {
		logger.Error("Load tracking log iteration failed", "job_id", jobID, "err", err)
	}
}

func blockExists(db *sql.DB, number int64) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM blocks WHERE number = ?", number).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func getParsedABI(db *sql.DB, address string) (abi.ABI, error) {
	lowAddr := strings.ToLower(address)
	if obj, ok := abiCache[lowAddr]; ok {
		return obj, nil
	}

	// If not in cache, we don't have the ABI - skip decoding
	// All ABIs are loaded into abiCache at startup, so no need to query database
	return abi.ABI{}, sql.ErrNoRows
}

// formatValue converts values to human-readable format
func formatValue(value interface{}) interface{} {
	switch v := value.(type) {
	case common.Address:
		return v.Hex()
	case common.Hash:
		return v.Hex()
	case *big.Int:
		return v.String()
	case []byte:
		return "0x" + hex.EncodeToString(v)
	case [32]byte:
		return "0x" + hex.EncodeToString(v[:])
	default:
		return fmt.Sprintf("%v", v)
	}
}

func decodeTxInput(input string, abiObj abi.ABI) *Decoded {
	if input == "" || !strings.HasPrefix(input, "0x") || len(input) < 10 {
		return nil
	}
	inputBytes := common.FromHex(input)
	if len(inputBytes) < 4 {
		return nil
	}
	sig := inputBytes[:4]
	method, err := abiObj.MethodById(sig)
	if err != nil {
		return nil
	}
	argValues, err := method.Inputs.Unpack(inputBytes[4:])
	if err != nil {
		return nil
	}
	argsMap := make(map[string]interface{})
	for i, input := range method.Inputs {
		argsMap[input.Name] = formatValue(argValues[i])
	}
	return &Decoded{
		Name: method.Name,
		Sig:  method.Sig,
		Args: argsMap,
	}
}

func decodeLog(topics []string, data string, abiObj abi.ABI) *Decoded {
	if len(topics) == 0 {
		return nil
	}
	eventID := common.HexToHash(topics[0])
	event, err := abiObj.EventByID(eventID)
	if err != nil {
		return nil
	}

	// Unpack the data portion using the event name
	dataBytes := common.FromHex(data)
	dataList, err := abiObj.Unpack(event.Name, dataBytes)
	if err != nil {
		return nil
	}

	argsMap := make(map[string]interface{})

	// Process indexed and non-indexed arguments
	topicIndex := 1 // Start at 1 because topics[0] is the event signature
	dataIndex := 0

	for _, input := range event.Inputs {
		if input.Indexed {
			// Indexed arguments are in topics
			if topicIndex >= len(topics) {
				return nil
			}
			topic := common.HexToHash(topics[topicIndex])

			// For dynamic types (string, bytes, arrays), only the hash is stored in topics
			// But we still show it in a readable format with a note
			if isDynamicType(&input.Type) {
				argsMap[input.Name] = map[string]interface{}{
					"indexed_hash": topic.Hex(),
					"note":         "indexed dynamic type - only hash available on-chain",
				}
			} else {
				// For static types, decode the actual value from the topic
				args := abi.Arguments{{Type: input.Type}}
				unpacked, err := args.Unpack(topic[:])
				if err != nil || len(unpacked) == 0 {
					argsMap[input.Name] = topic.Hex()
				} else {
					argsMap[input.Name] = formatValue(unpacked[0])
				}
			}
			topicIndex++
		} else {
			// Non-indexed arguments are in data - these have full values
			if dataIndex >= len(dataList) {
				return nil
			}
			argsMap[input.Name] = formatValue(dataList[dataIndex])
			dataIndex++
		}
	}

	return &Decoded{
		Name: event.Name,
		Sig:  event.Sig,
		Args: argsMap,
	}
}
