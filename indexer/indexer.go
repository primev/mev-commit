package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// Config holds indexer configuration.
type Config struct {
	logger *slog.Logger
	chain  Chain
	store  Store
}

// ForwardIndexer is an indexer that indexes from the
// last indexed block and retrieves the blocks going forward.
type ForwardIndexer struct {
	*Config

	lastIndexedBlock *big.Int
}

// Run starts the indexer process.
func (fi *ForwardIndexer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			blocks, err := fi.fetchBlocks(ctx)
			if err != nil {
				fi.logger.Error("fetch blocks", "error", err)
				// TODO: try again after some backoff time before continue...
				continue
			}
			err = fi.storeBlocks(ctx, blocks)
			if err != nil {
				fi.logger.Error("store blocks", "error", err)
			}
		}
	}
}

func (fi *ForwardIndexer) fetchBlocks(ctx context.Context) ([]*types.Block, error) {
	const window = 5

	var (
		res []*types.Block
	)

	last := fi.lastIndexedBlock
	curr, err := fi.chain.BlockNumber()
	if err != nil {
		return nil, fmt.Errorf("block number: %w", err)
	}

	var (
		start = new(big.Int).Add(last, big.NewInt(1))
		end   = new(big.Int).Add(start, big.NewInt(window))
	)
	for start.Cmp(curr) <= 0 {
		if end.Cmp(curr) > 0 {
			end.Set(curr)
		}

		blocks, err := fi.chain.Blocks(start, end)
		if err != nil {
			fi.logger.Error("fetch blocks", "start", start, "end", end, "error", err)
			continue
		}
		res = append(res, blocks...)

		start = new(big.Int).Add(end, big.NewInt(1))
		end = new(big.Int).Add(start, big.NewInt(window))
	}

	return res, nil
}

func (fi *ForwardIndexer) storeBlocks(ctx context.Context, blocks []*types.Block) error {
	for _, block := range blocks {
		txs, txh, err := parseBlockTransactions(block)
		if err != nil {
			return fmt.Errorf("parse block #%v transactions: %w", block.Number(), err)
		}

		rcs, err := fi.chain.TxReceipts(txh)
		if err != nil {
			return fmt.Errorf("fetch transaction receipts for block #%v: %w", block.Number(), err)
		}

		for _, tx := range txs {
			if rc, ok := rcs[tx.Hash]; ok {
				tx.Status = rc.Status
				tx.GasUsed = rc.GasUsed
				tx.CumulativeGasUsed = rc.CumulativeGasUsed
				tx.ContractAddress = rc.ContractAddress.Hex()
				tx.TransactionIndex = rc.TransactionIndex
				tx.ReceiptBlockHash = rc.BlockHash.Hex()
				tx.ReceiptBlockNumber = rc.BlockNumber.Uint64()
			}
		}

		bi := &BlockItem{
			Number:       block.NumberU64(),
			Hash:         block.Hash().Hex(),
			ParentHash:   block.ParentHash().Hex(),
			Root:         block.Root().Hex(),
			Nonce:        block.Nonce(),
			Timestamp:    time.UnixMilli(int64(block.Time())).UTC().Format(timeMilliZ),
			Transactions: len(block.Transactions()),
			BaseFee:      block.BaseFee().Uint64(),
			GasLimit:     block.GasLimit(),
			GasUsed:      block.GasUsed(),
			Difficulty:   block.Difficulty().Uint64(),
			ExtraData:    hex.EncodeToString(block.Extra()),
		}

		if err := fi.store.IndexBlock(ctx, bi); err != nil {
			return fmt.Errorf("index block #%v: %w", block.Number(), err)
		}

		if err := fi.store.IndexTransactions(ctx, txs); err != nil {
			return fmt.Errorf("index transactions for block #%v: %w", block.Number(), err)
		}

		fi.lastIndexedBlock = block.Number()
	}

	return nil
}

// BackwardIndexer is an indexer that indexes from the
// last indexed block and retrieves the blocks going backward.
type BackwardIndexer struct {
	*Config

	lastIndexedBlock *big.Int
}

// Run starts the indexer process.
func (bi *BackwardIndexer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// TODO: implement...
		}
	}
}

// BalanceIndexer is an indexer that indexes the
// balance of all accounts in the blockchain.
type BalanceIndexer struct {
	*Config

	lastIndexedBlock *big.Int
}

// Run starts the indexer process.
func (bi *BalanceIndexer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// TODO: implement...
		}
	}
}
