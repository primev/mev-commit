package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func processBatchSQL(blocks []*types.Block, receipts []types.Receipts, db *sql.DB, flushInterval time.Duration) error {
	var chunkBlocks []IndexBlock
	var chunkTxs []IndexTx
	var chunkReceipts []IndexReceipt
	var chunkLogs []IndexLog

	var chunkStart time.Time
	resetChunkTimer := func() { chunkStart = time.Now() }
	resetChunkTimer()

	flushChunk := func() error {
		if len(chunkBlocks) == 0 && len(chunkTxs) == 0 && len(chunkReceipts) == 0 && len(chunkLogs) == 0 {
			return nil
		}

		// Write to SQL
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if err := batchInsertBlocks(tx, chunkBlocks); err != nil {
			tx.Rollback()
			return err
		}
		if err := batchInsertTxs(tx, chunkTxs); err != nil {
			tx.Rollback()
			return err
		}
		if err := batchInsertReceipts(tx, chunkReceipts); err != nil {
			tx.Rollback()
			return err
		}
		if err := batchInsertLogs(tx, chunkLogs); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		slog.Info("Flushed chunk",
			"blocks", len(chunkBlocks),
			"txs", len(chunkTxs),
			"receipts", len(chunkReceipts),
			"logs", len(chunkLogs),
			"duration", time.Since(chunkStart),
			"mode", "sql",
		)

		chunkBlocks = chunkBlocks[:0]
		chunkTxs = chunkTxs[:0]
		chunkReceipts = chunkReceipts[:0]
		chunkLogs = chunkLogs[:0]
		resetChunkTimer()
		return nil
	}

	for k := range blocks {
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
		chunkBlocks = append(chunkBlocks, indexBlock)

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
				abiObjs, err := getParsedABI(db, *indexTx.To)
				if err == nil {
					decoded := decodeTxInput(indexTx.Input, abiObjs)
					if decoded != nil {
						decodedJSON, _ := json.Marshal(decoded)
						indexTx.DecodedJSON = string(decodedJSON)
					}
				}
			}

			chunkTxs = append(chunkTxs, indexTx)

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

			chunkReceipts = append(chunkReceipts, indexReceipt)

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
				abiObjs, err := getParsedABI(db, indexLog.Address)
				if err == nil && len(topics) > 0 {
					decoded := decodeLog(topics, indexLog.Data, abiObjs)
					if decoded != nil {
						decodedJSON, _ := json.Marshal(decoded)
						indexLog.DecodedJSON = string(decodedJSON)
					}
				}

				topicsJSON, _ := json.Marshal(topics)
				indexLog.TopicsJSON = string(topicsJSON)

				chunkLogs = append(chunkLogs, indexLog)
			}
		}

		if shouldFlushChunk(chunkBlocks, chunkTxs, chunkReceipts, chunkLogs, chunkStart, flushInterval) {
			slog.Debug("Chunk threshold reached",
				"blocks", len(chunkBlocks),
				"txs", len(chunkTxs),
				"receipts", len(chunkReceipts),
				"logs", len(chunkLogs),
				"elapsed", time.Since(chunkStart),
			)
			if err := flushChunk(); err != nil {
				return err
			}
		}
	}

	return flushChunk()
}
