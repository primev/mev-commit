package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	w3 "github.com/lmittmann/w3"
	eth "github.com/lmittmann/w3/module/eth"
	w3types "github.com/lmittmann/w3/w3types"
)

type Decoded struct {
	Name string      `json:"name"`
	Sig  string      `json:"sig"`
	Args interface{} `json:"args"`
}

type IndexBlock struct {
	Number           int64
	CoinbaseAddress  string
	Hash             string
	ParentHash       string
	Nonce            string
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	ReceiptsRoot     string
	Miner            string
	Difficulty       string
	ExtraData        string
	Size             int64
	GasLimit         int64
	GasUsed          int64
	Timestamp        int64
	BaseFeePerGas    *string
	BlobGasUsed      *int64
	ExcessBlobGas    *int64
	WithdrawalsRoot  *string
	RequestsHash     *string
	MixHash          *string
	TxCount          int
}

type IndexTx struct {
	Hash                 string
	BlockNumber          int64
	BlockHash            string
	TxIndex              int
	From                 string
	To                   *string
	Value                string
	Gas                  int64
	GasPrice             *string
	MaxPriorityFeePerGas *string
	MaxFeePerGas         *string
	EffectiveGasPrice    *string
	Input                string
	Type                 uint8
	ChainID              *int64
	AccessListJSON       *string
	BlobGas              *int64
	BlobGasFeeCap        *string
	BlobHashesJSON       *string
	V                    *string
	R                    *string
	S                    *string
	DecodedJSON          string
}

type IndexReceipt struct {
	TxHash            string
	Status            uint64
	CumulativeGasUsed int64
	GasUsed           int64
	ContractAddress   *string
	LogsBloom         string
	Type              uint8
	BlobGasUsed       uint64
	BlobGasPrice      *string
}

type IndexLog struct {
	TxHash         string
	LogIndex       int
	Address        string
	BlockNumber    *int64
	BlockHash      *string
	TxIndex        int
	BlockTimestamp *int64
	TopicsJSON     string
	Data           string
	Removed        bool
	DecodedJSON    string
}

var abiCache = make(map[string]abi.ABI)

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

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	rpcURL := flag.String("rpc", "http://localhost:8545", "Ethereum RPC URL")
	dsn := flag.String("dsn", "root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true", "StarRocks DSN")
	mode := flag.String("mode", "forward", "Mode: forward or backfill")
	fromBlock := flag.Int64("from", 0, "Starting block for backfill (lowest)")
	toBlock := flag.Int64("to", 0, "Ending block for backfill (highest)")
	pollInterval := flag.Duration("poll", 100*time.Millisecond, "Poll interval for forward mode")
	batchSize := flag.Int("batch-size", 25, "Batch size for RPC calls")
	abiDir := flag.String("abi-dir", "./contracts-abi/abi", "Directory containing ABI files")
	contractMapJSON := flag.String("contract-map", `{"BidderRegistry": "0x145a9f4cbae2ec281f417195ea3464dbd04289a2", "BlockTracker": "0x5d64b933739558101f9359e2750acc228f0cb64f", "L1Gateway": "0x5d64b933739558101f9359e2750acc228f0cb64f", "Oracle": "0x37a037d2423221f403cfa146f5fb962e19582d90", "PreconfManager": "0x2ee9e88f57a7db801e114a4df7a99eb7257871e2", "ProviderRegistry": "0xeb6d22309062a86fa194520344530874221ef48c", "SettlementGateway": "0x21f5f1142200a515248a2eef5b0654581c7f2b46"}`, "JSON mapping contract names to addresses")
	flag.Parse()

	if *mode != "forward" && *mode != "backfill" {
		logger.Error("Invalid mode", "mode", *mode)
		os.Exit(1)
	}
	if *mode == "backfill" && (*fromBlock >= *toBlock || *toBlock == 0) {
		logger.Error("Invalid backfill range", "from", *fromBlock, "to", *toBlock)
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

	// Create tables if they do not exist
	if err := createTables(db); err != nil {
		logger.Error("Failed to create tables", "err", err)
		os.Exit(1)
	}

	// Load ABIs
	if err := loadABIs(db, *abiDir, *contractMapJSON); err != nil {
		logger.Error("Failed to load ABIs", "err", err)
		// Continue or exit based on preference; here continue
	}

	if *mode == "forward" {
		forwardIndex(client, db, *pollInterval, *batchSize, logger)
	} else {
		backfillIndex(client, db, *fromBlock, *toBlock, *batchSize, logger)
	}
}

func loadABIs(db *sql.DB, abiDir, contractMapJSON string) error {
	contractMap := make(map[string]string)
	if err := json.Unmarshal([]byte(contractMapJSON), &contractMap); err != nil {
		return err
	}

	files, err := os.ReadDir(abiDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".abi") {
			name := strings.TrimSuffix(file.Name(), ".abi")
			address, ok := contractMap[name]
			if !ok {
				continue
			}
			path := filepath.Join(abiDir, file.Name())
			abiBytes, err := os.ReadFile(path)
			if err != nil {
				slog.Error("Failed to read ABI file", "path", path, "err", err)
				continue
			}
			abiStr := string(abiBytes)

			_, err = db.Exec("INSERT INTO contract_abis (address, name, abi) VALUES (?, ?, parse_json(?))", strings.ToLower(address), name, abiStr)
			if err != nil {
				slog.Error("Failed to insert ABI", "name", name, "err", err)
				continue
			}

			abiObj, err := abi.JSON(strings.NewReader(abiStr))
			if err != nil {
				slog.Error("Failed to parse ABI", "name", name, "err", err)
				continue
			}
			abiCache[strings.ToLower(address)] = abiObj
		}
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
			input VARCHAR(65535),
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
			data VARCHAR(65535),
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

func forwardIndex(client *w3.Client, db *sql.DB, pollInterval time.Duration, batchSize int, logger *slog.Logger) {
	for {
		var latest *big.Int
		if err := client.Call(eth.BlockNumber().Returns(&latest)); err != nil {
			logger.Error("Failed to get latest block", "err", err)
			time.Sleep(pollInterval)
			continue
		}

		maxIndexed, err := getMaxIndexedBlock(db)
		if err != nil {
			logger.Error("Failed to get max indexed block", "err", err)
			time.Sleep(pollInterval)
			continue
		}
		start := int64(0)
		if maxIndexed > 0 {
			start = maxIndexed + 1
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
			if err := processBatch(client, db, blockNums); err != nil {
				logger.Error("Failed to process batch", "from", start, "to", end, "err", err)
				break
			}
			start = end + 1
		}

		time.Sleep(pollInterval)
	}
}

func backfillIndex(client *w3.Client, db *sql.DB, from, to int64, batchSize int, logger *slog.Logger) {
	var pending []int64
	for i := to; i >= from; i-- {
		exists, err := blockExists(db, i)
		if err != nil {
			logger.Error("Failed to check if block exists", "block", i, "err", err)
			continue
		}
		if !exists {
			pending = append(pending, i)
			if len(pending) == batchSize {
				sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })
				if err := processBatch(client, db, pending); err != nil {
					logger.Error("Failed to process batch", "err", err)
				}
				pending = pending[:0]
			}
		}
	}
	if len(pending) > 0 {
		sort.Slice(pending, func(a, b int) bool { return pending[a] < pending[b] })
		if err := processBatch(client, db, pending); err != nil {
			logger.Error("Failed to process batch", "err", err)
		}
	}
}

func processBatch(client *w3.Client, db *sql.DB, blockNums []int64) error {
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

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var indexBlocks []IndexBlock
	var indexTxs []IndexTx
	var indexReceipts []IndexReceipt
	var indexLogs []IndexLog

	for k := range blockNums {
		header := blocks[k].Header()
		var baseFeePerGas *string
		if header.BaseFee != nil {
			s := header.BaseFee.String()
			baseFeePerGas = &s
		}
		var blobGasUsed *int64
		if header.BlobGasUsed != nil {
			u := int64(*header.BlobGasUsed)
			blobGasUsed = &u
		}
		var excessBlobGas *int64
		if header.ExcessBlobGas != nil {
			u := int64(*header.ExcessBlobGas)
			excessBlobGas = &u
		}
		var withdrawalsRoot *string
		if header.WithdrawalsHash != nil {
			s := header.WithdrawalsHash.Hex()
			withdrawalsRoot = &s
		}
		var requestsHash *string // Set to nil if not applicable
		var mixHash *string
		if header.MixDigest != (common.Hash{}) {
			s := header.MixDigest.Hex()
			mixHash = &s
		}

		indexBlock := IndexBlock{
			Number:           blocks[k].Number().Int64(),
			CoinbaseAddress:  header.Coinbase.Hex(),
			Hash:             blocks[k].Hash().Hex(),
			ParentHash:       header.ParentHash.Hex(),
			Nonce:            "0x" + hex.EncodeToString(header.Nonce[:]),
			Sha3Uncles:       header.UncleHash.Hex(),
			LogsBloom:        "0x" + hex.EncodeToString(header.Bloom[:]),
			TransactionsRoot: header.TxHash.Hex(),
			StateRoot:        header.Root.Hex(),
			ReceiptsRoot:     header.ReceiptHash.Hex(),
			Miner:            header.Coinbase.Hex(),
			Difficulty:       header.Difficulty.String(),
			ExtraData:        "0x" + hex.EncodeToString(header.Extra),
			Size:             int64(blocks[k].Size()),
			GasLimit:         int64(header.GasLimit),
			GasUsed:          int64(header.GasUsed),
			Timestamp:        int64(header.Time),
			BaseFeePerGas:    baseFeePerGas,
			BlobGasUsed:      blobGasUsed,
			ExcessBlobGas:    excessBlobGas,
			WithdrawalsRoot:  withdrawalsRoot,
			RequestsHash:     requestsHash,
			MixHash:          mixHash,
			TxCount:          len(blocks[k].Transactions()),
		}
		indexBlocks = append(indexBlocks, indexBlock)

		txs := blocks[k].Transactions()
		for i, txn := range txs {
			chainID := txn.ChainId().Int64()
			chainIDPtr := &chainID
			var to *string
			if txn.To() != nil {
				s := txn.To().Hex()
				to = &s
			}
			var gasPrice *string
			gp := txn.GasPrice()
			if gp != nil {
				s := gp.String()
				gasPrice = &s
			}
			var maxPriorityFeePerGas *string
			var maxFeePerGas *string
			if txn.Type() == types.DynamicFeeTxType || txn.Type() == types.BlobTxType {
				mp := txn.GasTipCap()
				if mp != nil {
					s := mp.String()
					maxPriorityFeePerGas = &s
				}
				mf := txn.GasFeeCap()
				if mf != nil {
					s := mf.String()
					maxFeePerGas = &s
				}
			}
			var accessListJSON *string
			al := txn.AccessList()
			if len(al) > 0 {
				alBytes, _ := json.Marshal(al)
				s := string(alBytes)
				accessListJSON = &s
			}
			var blobGas *int64
			bg := txn.BlobGas()
			if bg != 0 {
				bi := int64(bg)
				blobGas = &bi
			}
			var blobGasFeeCap *string
			if txn.Type() == types.BlobTxType {
				bgfc := txn.BlobGasFeeCap()
				if bgfc != nil {
					s := bgfc.String()
					blobGasFeeCap = &s
				}
			}
			var blobHashesJSON *string
			bh := txn.BlobHashes()
			if len(bh) > 0 {
				bhBytes, _ := json.Marshal(bh)
				s := string(bhBytes)
				blobHashesJSON = &s
			}
			v, r, s := txn.RawSignatureValues()
			var vStr, rStr, sStr *string
			if v != nil {
				vs := v.String()
				vStr = &vs
			}
			if r != nil {
				rs := r.String()
				rStr = &rs
			}
			if s != nil {
				ss := s.String()
				sStr = &ss
			}

			signer := types.LatestSignerForChainID(txn.ChainId())
			from, err := types.Sender(signer, txn)
			if err != nil {
				return fmt.Errorf("failed to recover sender for tx %s: %w", txn.Hash().Hex(), err)
			}

			receipt := receipts[k][i]
			var effectiveGasPrice *string
			if receipt.EffectiveGasPrice != nil {
				s := receipt.EffectiveGasPrice.String()
				effectiveGasPrice = &s
			}

			indexTx := IndexTx{
				Hash:                 txn.Hash().Hex(),
				BlockNumber:          blocks[k].Number().Int64(),
				BlockHash:            blocks[k].Hash().Hex(),
				TxIndex:              i,
				From:                 from.Hex(),
				To:                   to,
				Value:                txn.Value().String(),
				Gas:                  int64(txn.Gas()),
				GasPrice:             gasPrice,
				MaxPriorityFeePerGas: maxPriorityFeePerGas,
				MaxFeePerGas:         maxFeePerGas,
				EffectiveGasPrice:    effectiveGasPrice,
				Input:                "0x" + hex.EncodeToString(txn.Data()),
				Type:                 txn.Type(),
				ChainID:              chainIDPtr,
				AccessListJSON:       accessListJSON,
				BlobGas:              blobGas,
				BlobGasFeeCap:        blobGasFeeCap,
				BlobHashesJSON:       blobHashesJSON,
				V:                    vStr,
				R:                    rStr,
				S:                    sStr,
				DecodedJSON:          "",
			}

			// Attempt to decode tx input if To is set
			if indexTx.To != nil {
				abiObj, err := getParsedABI(db, *indexTx.To)
				if err == nil {
					decoded := decodeTxInput(indexTx.Input, abiObj)
					if decoded != nil {
						decodedJSON, _ := json.Marshal(decoded)
						indexTx.DecodedJSON = string(decodedJSON)
					}
				}
			}

			indexTxs = append(indexTxs, indexTx)

			var contractAddress *string
			if receipt.ContractAddress != (common.Address{}) {
				s := receipt.ContractAddress.Hex()
				contractAddress = &s
			}
			var blobGasPrice *string
			if receipt.BlobGasPrice != nil {
				s := receipt.BlobGasPrice.String()
				blobGasPrice = &s
			}

			indexReceipt := IndexReceipt{
				TxHash:            receipt.TxHash.Hex(),
				Status:            receipt.Status,
				CumulativeGasUsed: int64(receipt.CumulativeGasUsed),
				GasUsed:           int64(receipt.GasUsed),
				ContractAddress:   contractAddress,
				LogsBloom:         "0x" + hex.EncodeToString(receipt.Bloom[:]),
				Type:              receipt.Type,
				BlobGasUsed:       receipt.BlobGasUsed,
				BlobGasPrice:      blobGasPrice,
			}

			indexReceipts = append(indexReceipts, indexReceipt)

			for _, l := range receipt.Logs {
				var blockNumber *int64
				bn := int64(l.BlockNumber)
				blockNumber = &bn
				var blockHash *string
				bh := l.BlockHash.Hex()
				blockHash = &bh
				var blockTimestamp *int64
				ts := int64(header.Time)
				blockTimestamp = &ts

				topics := make([]string, len(l.Topics))
				for j, topic := range l.Topics {
					topics[j] = topic.Hex()
				}

				indexLog := IndexLog{
					TxHash:         l.TxHash.Hex(),
					LogIndex:       int(l.Index),
					Address:        l.Address.Hex(),
					BlockNumber:    blockNumber,
					BlockHash:      blockHash,
					TxIndex:        int(l.TxIndex),
					BlockTimestamp: blockTimestamp,
					TopicsJSON:     "",
					Data:           "0x" + hex.EncodeToString(l.Data),
					Removed:        l.Removed,
					DecodedJSON:    "",
				}

				// Attempt to decode log
				abiObj, err := getParsedABI(db, indexLog.Address)
				if err == nil && len(topics) > 0 {
					decoded := decodeLog(topics, indexLog.Data, abiObj)
					if decoded != nil {
						decodedJSON, _ := json.Marshal(decoded)
						indexLog.DecodedJSON = string(decodedJSON)
					}
				}

				topicsJSON, _ := json.Marshal(topics)
				indexLog.TopicsJSON = string(topicsJSON)

				indexLogs = append(indexLogs, indexLog)
			}
		}
	}

	if err := batchInsertBlocks(tx, indexBlocks); err != nil {
		return err
	}
	if err := batchInsertTxs(tx, indexTxs); err != nil {
		return err
	}
	if err := batchInsertReceipts(tx, indexReceipts); err != nil {
		return err
	}
	if err := batchInsertLogs(tx, indexLogs); err != nil {
		return err
	}

	return tx.Commit()
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
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, t.Hash, t.BlockNumber, t.BlockHash, t.TxIndex, t.From, t.To, t.Value, t.Gas, t.GasPrice, t.MaxPriorityFeePerGas, t.MaxFeePerGas, t.EffectiveGasPrice, t.Input, t.Type, t.ChainID, t.AccessListJSON, t.BlobGas, t.BlobGasFeeCap, t.BlobHashesJSON, t.V, t.R, t.S, t.DecodedJSON)
	}
	stmt := fmt.Sprintf("INSERT INTO transactions (hash, block_number, block_hash, tx_index, from_address, to_address, value, gas, gas_price, max_priority_fee_per_gas, max_fee_per_gas, effective_gas_price, input, type, chain_id, access_list_json, blob_gas, blob_gas_fee_cap, blob_hashes_json, v, r, s, decoded) VALUES %s", strings.Join(valueStrings, ", "))
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

func getMaxIndexedBlock(db *sql.DB) (int64, error) {
	var max sql.NullInt64
	err := db.QueryRow("SELECT MAX(number) FROM blocks").Scan(&max)
	if err != nil {
		return 0, err
	}
	if max.Valid {
		return max.Int64, nil
	}
	return 0, nil
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

	var abiStr string
	err := db.QueryRow("SELECT abi FROM contract_abis WHERE address = ?", lowAddr).Scan(&abiStr)
	if err != nil {
		return abi.ABI{}, err
	}

	obj, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return abi.ABI{}, err
	}

	abiCache[lowAddr] = obj
	return obj, nil
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
