package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func processBatchKafka(blocks []*types.Block, receipts []types.Receipts, kafkaOpts kafkaOptions, db *sql.DB) error {
	batchStart := time.Now()
	structBuildStart := time.Now()

	var batchBlocks []IndexBlock
	var batchTxs []IndexTx
	var batchReceipts []IndexReceipt
	var batchLogs []IndexLog

	for k := range blocks {
		var blockData []IndexBlock
		var txsData []IndexTx
		var receiptsData []IndexReceipt
		var logsData []IndexLog

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
		var requestsHash *string
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
		blockData = append(blockData, indexBlock)

		txs := blocks[k].Transactions()
		for i, txn := range txs {
			chainIDBig := txn.ChainId()
			var chainIDPtr *int64
			var signer types.Signer
			if chainIDBig != nil && chainIDBig.Sign() > 0 {
				chainID := chainIDBig.Int64()
				chainIDPtr = &chainID
				signer = types.LatestSignerForChainID(chainIDBig)
			} else {
				signer = types.HomesteadSigner{}
			}
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
				Nonce:                txn.Nonce(),
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

			txsData = append(txsData, indexTx)

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

			receiptsData = append(receiptsData, indexReceipt)

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

				logsData = append(logsData, indexLog)
			}
		}
		
		batchBlocks = append(batchBlocks, blockData...)
		batchTxs = append(batchTxs, txsData...)
		batchReceipts = append(batchReceipts, receiptsData...)
		batchLogs = append(batchLogs, logsData...)
	}
	structBuildDuration := time.Since(structBuildStart)

	// Produce once for the whole fetch batch, then wait once inside produceToKafka
	produceStart := time.Now()
	if err := produceToKafka(kafkaOpts, batchBlocks, batchTxs, batchReceipts, batchLogs); err != nil {
		return fmt.Errorf("failed to produce batch to Kafka: %w", err)
	}
	produceDuration := time.Since(produceStart)
	totalDuration := time.Since(batchStart)

	// Detailed timing breakdown
	slog.Info("Kafka processing timing breakdown",
		"total_duration", totalDuration,
		"struct_build_duration", structBuildDuration,
		"struct_build_pct", fmt.Sprintf("%.1f%%", 100*structBuildDuration.Seconds()/totalDuration.Seconds()),
		"produce_duration", produceDuration,
		"produce_pct", fmt.Sprintf("%.1f%%", 100*produceDuration.Seconds()/totalDuration.Seconds()),
		"blocks", len(batchBlocks),
		"txs", len(batchTxs),
		"receipts", len(batchReceipts),
		"logs", len(batchLogs),
	)

	// Throughput metrics
	dur := totalDuration.Seconds()
	if dur > 0 {
		slog.Info("Kafka batch throughput",
			"blocks_per_sec", fmt.Sprintf("%.2f", float64(len(batchBlocks))/dur),
			"txs_per_sec", fmt.Sprintf("%.2f", float64(len(batchTxs))/dur),
			"receipts_per_sec", fmt.Sprintf("%.2f", float64(len(batchReceipts))/dur),
			"logs_per_sec", fmt.Sprintf("%.2f", float64(len(batchLogs))/dur),
		)
	}

	return nil
}

// getKafkaCheckpointRange reads the min and max block numbers from the Kafka blocks topic
// Returns (minBlock, maxBlock, hasData, error)
func getKafkaCheckpointRange(kafkaOpts kafkaOptions) (int64, int64, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	topic := fmt.Sprintf("%s-blocks", kafkaOpts.chainName)

	// Use kadm admin client to get start and end offsets
	adminClient := kadm.NewClient(kafkaOpts.client)
	offsets, err := adminClient.ListEndOffsets(ctx, topic)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to list end offsets for topic %s: %w", topic, err)
	}

	// Get partition 0 (we only use 1 partition)
	offset, ok := offsets.Lookup(topic, 0)
	if !ok {
		return 0, 0, false, fmt.Errorf("partition 0 not found for topic %s", topic)
	}

	highWatermark := offset.Offset
	if highWatermark == 0 {
		// Topic is empty, no checkpoint exists
		slog.Info("Kafka checkpoint: topic is empty", "topic", topic)
		return 0, 0, false, nil
	}

	// Fetch first and last messages
	// First message at offset 0
	kafkaOpts.client.AddConsumePartitions(map[string]map[int32]kgo.Offset{
		topic: {0: kgo.NewOffset().At(0)},
	})

	fetches := kafkaOpts.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		kafkaOpts.client.RemoveConsumePartitions(map[string][]int32{topic: {0}})
		return 0, 0, false, fmt.Errorf("fetch errors for first message: %v", errs)
	}

	var firstBlock IndexBlock
	fetches.EachRecord(func(r *kgo.Record) {
		if err := json.Unmarshal(r.Value, &firstBlock); err != nil {
			slog.Error("Failed to unmarshal first block from Kafka", "err", err)
		}
	})

	kafkaOpts.client.RemoveConsumePartitions(map[string][]int32{topic: {0}})

	// Last message at offset = highWatermark - 1
	lastOffset := highWatermark - 1
	kafkaOpts.client.AddConsumePartitions(map[string]map[int32]kgo.Offset{
		topic: {0: kgo.NewOffset().At(lastOffset)},
	})

	fetches = kafkaOpts.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		kafkaOpts.client.RemoveConsumePartitions(map[string][]int32{topic: {0}})
		return 0, 0, false, fmt.Errorf("fetch errors for last message: %v", errs)
	}

	var lastBlock IndexBlock
	fetches.EachRecord(func(r *kgo.Record) {
		if err := json.Unmarshal(r.Value, &lastBlock); err != nil {
			slog.Error("Failed to unmarshal last block from Kafka", "err", err)
		}
	})

	kafkaOpts.client.RemoveConsumePartitions(map[string][]int32{topic: {0}})

	if firstBlock.Number == 0 || lastBlock.Number == 0 {
		return 0, 0, false, fmt.Errorf("failed to decode blocks from Kafka topic %s", topic)
	}

	slog.Info("Kafka checkpoint range found",
		"topic", topic,
		"min_block", firstBlock.Number,
		"max_block", lastBlock.Number,
		"total_blocks", highWatermark,
	)

	return firstBlock.Number, lastBlock.Number, true, nil
}

// getKafkaCheckpoint reads the last block number from the Kafka blocks topic
// Returns the last indexed block number, whether any blocks exist, and any error
func getKafkaCheckpoint(kafkaOpts kafkaOptions) (int64, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	topic := fmt.Sprintf("%s-blocks", kafkaOpts.chainName)

	// Use kadm admin client to get end offsets
	adminClient := kadm.NewClient(kafkaOpts.client)
	offsets, err := adminClient.ListEndOffsets(ctx, topic)
	if err != nil {
		return 0, false, fmt.Errorf("failed to list end offsets for topic %s: %w", topic, err)
	}

	// Get partition 0 (we only use 1 partition)
	offset, ok := offsets.Lookup(topic, 0)
	if !ok {
		return 0, false, fmt.Errorf("partition 0 not found for topic %s", topic)
	}

	highWatermark := offset.Offset
	if highWatermark == 0 {
		// Topic is empty, no checkpoint exists
		slog.Info("Kafka checkpoint: topic is empty", "topic", topic)
		return 0, false, nil
	}

	// Fetch the last message (offset = highWatermark - 1)
	lastOffset := highWatermark - 1

	// Use the same client with direct consumption to avoid recreating connection
	kafkaOpts.client.AddConsumePartitions(map[string]map[int32]kgo.Offset{
		topic: {0: kgo.NewOffset().At(lastOffset)},
	})
	defer kafkaOpts.client.RemoveConsumePartitions(map[string][]int32{
		topic: {0},
	})

	// Fetch the last record
	fetches := kafkaOpts.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		return 0, false, fmt.Errorf("fetch errors: %v", errs)
	}

	var lastBlock IndexBlock
	fetches.EachRecord(func(r *kgo.Record) {
		if err := json.Unmarshal(r.Value, &lastBlock); err != nil {
			slog.Error("Failed to unmarshal last block from Kafka", "err", err)
		}
	})

	if lastBlock.Number == 0 {
		return 0, false, fmt.Errorf("failed to decode last block from Kafka topic %s", topic)
	}

	slog.Info("Kafka checkpoint found",
		"topic", topic,
		"last_block", lastBlock.Number,
		"last_block_hash", lastBlock.Hash,
		"offset", lastOffset,
	)

	return lastBlock.Number, true, nil
}
