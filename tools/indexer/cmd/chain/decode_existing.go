package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// This file contains functions to decode existing transactions and logs in the database
// that were imported from external ETL and don't have decoded input/logs.
//
// These functions reuse the same ABI loading and decoding logic as the main indexer.
//
// Shared decode functions from indexer.go:
// - decodeTxInput: decodes transaction input data using contract ABI
// - decodeLog: decodes event logs using contract ABI
// - formatValue: formats decoded values to human-readable format
// - getParsedABI: retrieves parsed ABI from cache for a contract address
// - abiCache: global cache of loaded ABIs
// - loadABIs: loads ABIs from manifest file
// - loadABIsFromManifest: loads ABIs from JSON manifest

func decodeExistingTransactions(ctx context.Context, db *sql.DB, batchSize int, dryRun bool, logger *slog.Logger) error {
	// Get list of contract addresses that have ABIs (from abiCache which is already loaded)
	if len(abiCache) == 0 {
		logger.Warn("No ABIs loaded in cache, cannot decode transactions")
		return nil
	}

	// Build list of addresses for IN clause
	addresses := make([]string, 0, len(abiCache))
	for addr := range abiCache {
		addresses = append(addresses, "'"+addr+"'")
	}
	addressList := "(" + strings.Join(addresses, ",") + ")"

	// Count transactions that need decoding for contracts we have ABIs for
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM transactions
		WHERE LOWER(to_address) IN %s
		AND (decoded IS NULL OR decoded = '{}' OR decoded = '')
	`, addressList)

	var totalCount int64
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return fmt.Errorf("failed to count transactions: %w", err)
	}

	if totalCount == 0 {
		logger.Info("No transactions need decoding for contracts with ABIs")
		return nil
	}

	logger.Info("Found transactions to decode",
		"total", totalCount,
		"contracts_with_abis", len(abiCache))

	var processed int64
	var decoded int64
	var failed int64
	startTime := time.Now()
	offset := int64(0)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Transaction decoding interrupted", "processed", processed, "decoded", decoded, "failed", failed)
			return ctx.Err()
		default:
		}

		// Fetch a batch of transactions that need decoding, only for contracts with ABIs
		query := fmt.Sprintf(`
			SELECT hash, to_address, input
			FROM transactions
			WHERE LOWER(to_address) IN %s
			AND (decoded IS NULL OR decoded = '{}' OR decoded = '')
			LIMIT ? OFFSET ?
		`, addressList)

		rows, err := db.QueryContext(ctx, query, batchSize, offset)
		if err != nil {
			return fmt.Errorf("failed to query transactions: %w", err)
		}

		var batchTxs []struct {
			Hash    string
			To      string
			Input   string
			Decoded string
		}

		for rows.Next() {
			var tx struct {
				Hash    string
				To      string
				Input   string
				Decoded string
			}
			if err := rows.Scan(&tx.Hash, &tx.To, &tx.Input); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan transaction: %w", err)
			}
			batchTxs = append(batchTxs, tx)
		}
		rows.Close()

		if len(batchTxs) == 0 {
			break // No more transactions to process
		}

		batchStart := time.Now()
		batchDecoded := 0
		batchFailed := 0

		// Process each transaction in the batch
		for i := range batchTxs {
			tx := &batchTxs[i]

			// Get ABIs for the contract (should always exist since we filtered by address)
			abiObjs, err := getParsedABI(db, tx.To)
			if err != nil {
				// Shouldn't happen but handle gracefully
				logger.Error("ABI not found for contract (shouldn't happen - check filter query)",
					"hash", tx.Hash,
					"to_address", tx.To,
					"err", err)
				batchFailed++
				failed++
				continue
			}

			// Check if input is empty or too short
			if tx.Input == "" || tx.Input == "0x" {
				logger.Warn("Transaction has no input data (plain ETH transfer)",
					"hash", tx.Hash,
					"to_address", tx.To,
					"input", tx.Input)
				batchFailed++
				failed++
				continue
			}

			if len(tx.Input) < 10 {
				logger.Warn("Transaction input too short to decode",
					"hash", tx.Hash,
					"to_address", tx.To,
					"input", tx.Input,
					"input_length", len(tx.Input))
				batchFailed++
				failed++
				continue
			}

			// Get method signature (first 4 bytes)
			methodSig := tx.Input[0:10] // "0x" + 8 hex chars = 10 chars

			// Decode the input (tries all ABIs until one succeeds)
			decodedData := decodeTxInput(tx.Input, abiObjs)
			if decodedData == nil {
				// Could not decode - log details
				logger.Error("Failed to decode transaction input",
					"hash", tx.Hash,
					"to_address", tx.To,
					"method_signature", methodSig,
					"input_preview", tx.Input[0:min(66, len(tx.Input))], // First 32 bytes
					"reason", "method signature not found in ABI or invalid input data")
				batchFailed++
				failed++
				continue
			}

			// Marshal to JSON
			decodedJSON, err := json.Marshal(decodedData)
			if err != nil {
				logger.Error("Failed to marshal decoded transaction to JSON", "hash", tx.Hash, "err", err)
				batchFailed++
				failed++
				continue
			}

			tx.Decoded = string(decodedJSON)
			batchDecoded++
		}

		// Update the database in a batch
		// StarRocks doesn't support UPDATE in explicit transactions, so we use a single UPDATE with CASE
		// This reduces the number of versions created per tablet
		// To avoid exceeding StarRocks' max_scalar_operator_flat_children limit, split into chunks
		if !dryRun && batchDecoded > 0 {
			const maxCaseConditions = 100 // StarRocks limit for CASE WHEN conditions

			// Collect all transactions that need updating
			var txsToUpdate []struct {
				Hash    string
				Decoded string
			}
			for i := range batchTxs {
				tx := &batchTxs[i]
				if tx.Decoded != "" {
					txsToUpdate = append(txsToUpdate, struct {
						Hash    string
						Decoded string
					}{tx.Hash, tx.Decoded})
				}
			}

			// Process in chunks
			for chunkStart := 0; chunkStart < len(txsToUpdate); chunkStart += maxCaseConditions {
				chunkEnd := chunkStart + maxCaseConditions
				if chunkEnd > len(txsToUpdate) {
					chunkEnd = len(txsToUpdate)
				}
				chunk := txsToUpdate[chunkStart:chunkEnd]

				var caseParts []string
				var hashList []string

				for _, tx := range chunk {
					escapedDecoded := strings.ReplaceAll(tx.Decoded, "'", "''")
					caseParts = append(caseParts, fmt.Sprintf("WHEN '%s' THEN parse_json('%s')", tx.Hash, escapedDecoded))
					hashList = append(hashList, "'"+tx.Hash+"'")
				}

				query := fmt.Sprintf(`UPDATE transactions
					SET decoded = CASE hash
						%s
					END
					WHERE hash IN (%s)`,
					strings.Join(caseParts, " "),
					strings.Join(hashList, ","))

				if _, err := db.ExecContext(ctx, query); err != nil {
					logger.Error("Failed to update transaction chunk", "chunk_size", len(chunk), "err", err)
					// Don't fall back to individual updates - it creates too many tablet versions
					// The chunk size of 100 should work for StarRocks limits
				}
			}
		}

		processed += int64(len(batchTxs))
		decoded += int64(batchDecoded)
		offset += int64(len(batchTxs))

		elapsed := time.Since(batchStart)
		overall := time.Since(startTime)
		remaining := totalCount - processed
		var eta time.Duration
		if processed > 0 {
			eta = time.Duration(float64(overall) / float64(processed) * float64(remaining))
		}

		logger.Info("Processed transaction batch",
			"batch_size", len(batchTxs),
			"batch_decoded", batchDecoded,
			"batch_failed", batchFailed,
			"batch_duration", elapsed,
			"total_processed", processed,
			"total_decoded", decoded,
			"total_failed", failed,
			"progress_pct", fmt.Sprintf("%.2f%%", float64(processed)/float64(totalCount)*100),
			"eta", eta.Round(time.Second),
		)
	}

	duration := time.Since(startTime)
	logger.Info("Transaction decoding completed",
		"total_processed", processed,
		"total_decoded", decoded,
		"total_failed", failed,
		"duration", duration,
		"rate", fmt.Sprintf("%.2f tx/sec", float64(processed)/duration.Seconds()),
	)

	return nil
}

func decodeExistingLogs(ctx context.Context, db *sql.DB, batchSize int, dryRun bool, logger *slog.Logger) error {
	// Get list of contract addresses that have ABIs (from abiCache which is already loaded)
	if len(abiCache) == 0 {
		logger.Warn("No ABIs loaded in cache, cannot decode logs")
		return nil
	}

	// Build list of addresses for IN clause
	addresses := make([]string, 0, len(abiCache))
	for addr := range abiCache {
		addresses = append(addresses, "'"+addr+"'")
	}
	addressList := "(" + strings.Join(addresses, ",") + ")"

	// Count logs that need decoding for contracts we have ABIs for
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM logs
		WHERE LOWER(address) IN %s
		AND (decoded IS NULL OR decoded = '{}' OR decoded = '')
	`, addressList)

	var totalCount int64
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return fmt.Errorf("failed to count logs: %w", err)
	}

	if totalCount == 0 {
		logger.Info("No logs need decoding for contracts with ABIs")
		return nil
	}

	logger.Info("Found logs to decode",
		"total", totalCount,
		"contracts_with_abis", len(abiCache))

	var processed int64
	var decoded int64
	var failed int64
	startTime := time.Now()
	offset := int64(0)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Log decoding interrupted", "processed", processed, "decoded", decoded, "failed", failed)
			return ctx.Err()
		default:
		}

		// Fetch a batch of logs that need decoding, only for contracts with ABIs
		query := fmt.Sprintf(`
			SELECT tx_hash, log_index, address, topics, data
			FROM logs
			WHERE LOWER(address) IN %s
			AND (decoded IS NULL OR decoded = '{}' OR decoded = '')
			LIMIT ? OFFSET ?
		`, addressList)

		rows, err := db.QueryContext(ctx, query, batchSize, offset)
		if err != nil {
			return fmt.Errorf("failed to query logs: %w", err)
		}

		var batchLogs []struct {
			TxHash   string
			LogIndex int
			Address  string
			Topics   string
			Data     string
			Decoded  string
		}

		for rows.Next() {
			var log struct {
				TxHash   string
				LogIndex int
				Address  string
				Topics   string
				Data     string
				Decoded  string
			}
			if err := rows.Scan(&log.TxHash, &log.LogIndex, &log.Address, &log.Topics, &log.Data); err != nil {
				rows.Close()
				return fmt.Errorf("failed to scan log: %w", err)
			}
			batchLogs = append(batchLogs, log)
		}
		rows.Close()

		if len(batchLogs) == 0 {
			break // No more logs to process
		}

		batchStart := time.Now()
		batchDecoded := 0
		batchFailed := 0

		// Process each log in the batch
		for i := range batchLogs {
			log := &batchLogs[i]

			// Parse topics from JSON
			var topics []string
			if err := json.Unmarshal([]byte(log.Topics), &topics); err != nil {
				logger.Error("Failed to parse topics JSON",
					"tx_hash", log.TxHash,
					"log_index", log.LogIndex,
					"topics_raw", log.Topics,
					"err", err)
				batchFailed++
				failed++
				continue
			}

			if len(topics) == 0 {
				logger.Warn("Log has no topics (anonymous event or malformed)",
					"tx_hash", log.TxHash,
					"log_index", log.LogIndex,
					"address", log.Address)
				batchFailed++
				failed++
				continue
			}

			// Get ABIs for the contract (should always exist since we filtered by address)
			abiObjs, err := getParsedABI(db, log.Address)
			if err != nil {
				logger.Error("ABI not found for log contract (shouldn't happen - check filter query)",
					"tx_hash", log.TxHash,
					"log_index", log.LogIndex,
					"address", log.Address,
					"err", err)
				batchFailed++
				failed++
				continue
			}

			// Get event signature (first topic)
			eventSig := topics[0]

			// Decode the log (tries all ABIs until one succeeds)
			decodedData := decodeLog(topics, log.Data, abiObjs)
			if decodedData == nil {
				logger.Error("Failed to decode log",
					"tx_hash", log.TxHash,
					"log_index", log.LogIndex,
					"address", log.Address,
					"event_signature", eventSig,
					"topics_count", len(topics),
					"data_preview", log.Data[0:min(66, len(log.Data))], // First 32 bytes
					"reason", "event signature not found in ABI or invalid log data")
				batchFailed++
				failed++
				continue
			}

			// Marshal to JSON
			decodedJSON, err := json.Marshal(decodedData)
			if err != nil {
				logger.Error("Failed to marshal decoded log to JSON",
					"tx_hash", log.TxHash,
					"log_index", log.LogIndex,
					"err", err)
				batchFailed++
				failed++
				continue
			}

			log.Decoded = string(decodedJSON)
			batchDecoded++
		}

		// Update the database in a batch
		// StarRocks doesn't support UPDATE in explicit transactions, so we use a single UPDATE with CASE
		// This reduces the number of versions created per tablet
		// To avoid exceeding StarRocks' max_scalar_operator_flat_children limit, split into chunks
		if !dryRun && batchDecoded > 0 {
			const maxCaseConditions = 100 // StarRocks limit for CASE WHEN conditions

			// Collect all logs that need updating
			var logsToUpdate []struct {
				TxHash   string
				LogIndex int
				Decoded  string
			}
			for i := range batchLogs {
				log := &batchLogs[i]
				if log.Decoded != "" {
					logsToUpdate = append(logsToUpdate, struct {
						TxHash   string
						LogIndex int
						Decoded  string
					}{log.TxHash, log.LogIndex, log.Decoded})
				}
			}

			// Process in chunks
			for chunkStart := 0; chunkStart < len(logsToUpdate); chunkStart += maxCaseConditions {
				chunkEnd := chunkStart + maxCaseConditions
				if chunkEnd > len(logsToUpdate) {
					chunkEnd = len(logsToUpdate)
				}
				chunk := logsToUpdate[chunkStart:chunkEnd]

				var caseParts []string
				var whereConditions []string

				for _, log := range chunk {
					escapedDecoded := strings.ReplaceAll(log.Decoded, "'", "''")
					caseParts = append(caseParts, fmt.Sprintf("WHEN tx_hash = '%s' AND log_index = %d THEN parse_json('%s')",
						log.TxHash, log.LogIndex, escapedDecoded))
					whereConditions = append(whereConditions, fmt.Sprintf("(tx_hash = '%s' AND log_index = %d)",
						log.TxHash, log.LogIndex))
				}

				query := fmt.Sprintf(`UPDATE logs
					SET decoded = CASE
						%s
					END
					WHERE %s`,
					strings.Join(caseParts, " "),
					strings.Join(whereConditions, " OR "))

				if _, err := db.ExecContext(ctx, query); err != nil {
					logger.Error("Failed to update log chunk", "chunk_size", len(chunk), "err", err)
					// Don't fall back to individual updates - it creates too many tablet versions
					// The chunk size of 100 should work for StarRocks limits
				}
			}
		}

		processed += int64(len(batchLogs))
		decoded += int64(batchDecoded)
		offset += int64(len(batchLogs))

		elapsed := time.Since(batchStart)
		overall := time.Since(startTime)
		remaining := totalCount - processed
		var eta time.Duration
		if processed > 0 {
			eta = time.Duration(float64(overall) / float64(processed) * float64(remaining))
		}

		logger.Info("Processed log batch",
			"batch_size", len(batchLogs),
			"batch_decoded", batchDecoded,
			"batch_failed", batchFailed,
			"batch_duration", elapsed,
			"total_processed", processed,
			"total_decoded", decoded,
			"total_failed", failed,
			"progress_pct", fmt.Sprintf("%.2f%%", float64(processed)/float64(totalCount)*100),
			"eta", eta.Round(time.Second),
		)
	}

	duration := time.Since(startTime)
	logger.Info("Log decoding completed",
		"total_processed", processed,
		"total_decoded", decoded,
		"total_failed", failed,
		"duration", duration,
		"rate", fmt.Sprintf("%.2f logs/sec", float64(processed)/duration.Seconds()),
	)

	return nil
}
