package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/indexer/pkg/ethclient"
	"github.com/primev/mev-commit/indexer/pkg/store"
)

const TimeLayOut = "2006-01-02T15:04:05.000Z"

type Config struct {
	EthClient                        ethclient.EthereumClient
	Storage                          store.Storage
	IndexInterval                    time.Duration
	AccountAddresses                 []string
	MinBlocksToFetchAccountAddresses uint
	TimeoutToFetchAccountAddresses   time.Duration
}

type BlockchainIndexer struct {
	ethClient                        ethclient.EthereumClient
	storage                          store.Storage
	forwardBlockChan                 chan *types.Block
	backwardBlockChan                chan *types.Block
	txChan                           chan *types.Transaction
	indexInterval                    time.Duration
	lastForwardIndexedBlock          *big.Int
	lastBackwardIndexedBlock         *big.Int
	logger                           *slog.Logger
	accountAddresses                 []string
	blockCounter                     uint
	minBlocksToFetchAccountAddresses uint
	timeoutToFetchAccountAddresses   time.Duration
}

func NewBlockchainIndexer(config Config) *BlockchainIndexer {
	return &BlockchainIndexer{
		ethClient:                        config.EthClient,
		storage:                          config.Storage,
		forwardBlockChan:                 make(chan *types.Block, 100),
		backwardBlockChan:                make(chan *types.Block, 100),
		txChan:                           make(chan *types.Transaction, 100),
		indexInterval:                    config.IndexInterval,
		logger:                           slog.Default(),
		accountAddresses:                 config.AccountAddresses,
		blockCounter:                     0,
		minBlocksToFetchAccountAddresses: config.MinBlocksToFetchAccountAddresses,
		timeoutToFetchAccountAddresses:   config.TimeoutToFetchAccountAddresses,
	}
}

func (bi *BlockchainIndexer) Start(ctx context.Context) error {
	if err := bi.storage.CreateIndices(ctx); err != nil {
		return fmt.Errorf("failed to create indices: %w", err)
	}

	latestBlockNumber, err := bi.ethClient.BlockNumber(ctx)
	bi.logger.Info("latest block number", "block number", latestBlockNumber)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	if err = bi.initializeForwardIndex(ctx, latestBlockNumber.Uint64()); err != nil {
		return err
	}

	if err = bi.initializeBackwardIndex(ctx, latestBlockNumber.Uint64()); err != nil {
		return err
	}

	go bi.fetchForwardBlocks(ctx)
	go bi.processForwardBlocks(ctx)
	go bi.fetchBackwardBlocks(ctx)
	go bi.processBackwardBlocks(ctx)
	go bi.IndexAccountBalances(ctx)

	<-ctx.Done()
	return ctx.Err()
}

func (bi *BlockchainIndexer) initializeForwardIndex(ctx context.Context, latestBlockNumber uint64) error {
	lastForwardIndexedBlock, err := bi.storage.GetLastIndexedBlock(ctx, "forward")
	if err != nil {
		return fmt.Errorf("failed to get last forward indexed block: %w", err)
	}

	bi.logger.Info("last indexed block", "blockNumber", lastForwardIndexedBlock, "direction", "forward")

	if lastForwardIndexedBlock == nil || lastForwardIndexedBlock.Sign() == 0 {
		bi.lastForwardIndexedBlock = new(big.Int).SetUint64(latestBlockNumber - 1)
	} else {
		bi.lastForwardIndexedBlock = lastForwardIndexedBlock
	}

	return nil
}

func (bi *BlockchainIndexer) initializeBackwardIndex(ctx context.Context, latestBlockNumber uint64) error {
	lastBackwardIndexedBlock, err := bi.storage.GetLastIndexedBlock(ctx, "backward")
	if err != nil {
		return fmt.Errorf("failed to get last backward indexed block: %w", err)
	}

	bi.logger.Info("last indexed block", "blockNumber", lastBackwardIndexedBlock, "direction", "backward")

	if lastBackwardIndexedBlock == nil || lastBackwardIndexedBlock.Sign() == 0 {
		bi.lastBackwardIndexedBlock = new(big.Int).SetUint64(latestBlockNumber)
	} else {
		bi.lastBackwardIndexedBlock = lastBackwardIndexedBlock
	}

	return nil
}

func (bi *BlockchainIndexer) fetchForwardBlocks(ctx context.Context) {
	ticker := time.NewTicker(bi.indexInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			latestBlockNumber, err := bi.ethClient.BlockNumber(ctx)
			if err != nil {
				bi.logger.Error("failed to get latest block number", "error", err)
				continue
			}

			for blockNum := new(big.Int).Add(bi.lastForwardIndexedBlock, big.NewInt(1)); blockNum.Cmp(latestBlockNumber) <= 0; blockNum.Add(blockNum, big.NewInt(5)) {
				endBlockNum := new(big.Int).Add(blockNum, big.NewInt(4))
				if endBlockNum.Cmp(latestBlockNumber) > 0 {
					endBlockNum.Set(latestBlockNumber)
				}

				var blockNums []*big.Int
				for bn := new(big.Int).Set(blockNum); bn.Cmp(endBlockNum) <= 0; bn.Add(bn, big.NewInt(1)) {
					blockNums = append(blockNums, new(big.Int).Set(bn))
				}

				blocks, err := bi.fetchBlocks(ctx, blockNums)
				if err != nil {
					bi.logger.Error("failed to fetch blocks", "start", blockNum, "end", endBlockNum, "error", err)
					continue
				}

				for _, block := range blocks {
					bi.forwardBlockChan <- block
					bi.lastForwardIndexedBlock.Set(block.Number())
					bi.blockCounter++
				}
			}
		}
	}
}

func (bi *BlockchainIndexer) fetchBackwardBlocks(ctx context.Context) {
	ticker := time.NewTicker(bi.indexInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if bi.lastBackwardIndexedBlock.Sign() <= 0 {
				return
			}
			zeroBigNum := big.NewInt(0)
			blockNum := new(big.Int).Sub(bi.lastBackwardIndexedBlock, big.NewInt(1))

			for i := 0; blockNum.Cmp(zeroBigNum) >= 0; i++ {
				endBlockNum := new(big.Int).Sub(blockNum, big.NewInt(4))
				if endBlockNum.Cmp(zeroBigNum) < 0 {
					endBlockNum.Set(zeroBigNum)
				}

				var blockNums []*big.Int
				for bn := new(big.Int).Set(blockNum); bn.Cmp(endBlockNum) >= 0; bn.Sub(bn, big.NewInt(1)) {
					blockNums = append(blockNums, new(big.Int).Set(bn))
				}

				blocks, err := bi.fetchBlocks(ctx, blockNums)
				if err != nil {
					bi.logger.Error("failed to fetch blocks", "start", blockNum, "end", endBlockNum, "error", err)
					break
				}

				for _, block := range blocks {
					bi.backwardBlockChan <- block
					bi.lastBackwardIndexedBlock.Set(block.Number())
					if block.Number().Cmp(zeroBigNum) == 0 {
						bi.logger.Info("done fetching backward blocks...")
						return
					}
				}
				blockNum.Sub(endBlockNum, big.NewInt(1))
			}
		}
	}
}

func (bi *BlockchainIndexer) processForwardBlocks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-bi.forwardBlockChan:
			if err := bi.indexBlock(ctx, block); err != nil {
				bi.logger.Error("failed to index block", "error", err)
			}
			if err := bi.indexTransactions(ctx, block); err != nil {
				bi.logger.Error("failed to index transactions", "error", err)
			}
		}
	}
}

func (bi *BlockchainIndexer) processBackwardBlocks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-bi.backwardBlockChan:
			if err := bi.indexBlock(ctx, block); err != nil {
				bi.logger.Error("failed to index block", "error", err)
			}
			if err := bi.indexTransactions(ctx, block); err != nil {
				bi.logger.Error("failed to index transactions", "error", err)
			}
			if block.Number().Cmp(big.NewInt(0)) == 0 {
				bi.logger.Info("done processing backward blocks...")
				return
			}
		}
	}
}

func (bi *BlockchainIndexer) IndexAccountBalances(ctx context.Context) {
	timer := time.NewTimer(bi.timeoutToFetchAccountAddresses)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := bi.indexBalances(ctx, 0); err != nil {
				return
			}
			bi.blockCounter = 0
			timer.Reset(bi.timeoutToFetchAccountAddresses)
		default:
			if bi.blockCounter >= bi.minBlocksToFetchAccountAddresses {
				if err := bi.indexBalances(ctx, bi.lastForwardIndexedBlock.Uint64()); err != nil {
					return
				}
				bi.blockCounter = 0
				timer.Reset(bi.timeoutToFetchAccountAddresses)
			}
		}
	}
}

func (bi *BlockchainIndexer) indexBalances(ctx context.Context, blockNumber uint64) error {
	addresses, err := bi.storage.GetAddresses(ctx)
	if err != nil {
		return err
	}

	addresses = append(addresses, bi.accountAddresses...)

	addrs := make([]common.Address, len(addresses))
	for i, address := range addresses {
		addrs[i] = common.HexToAddress(address)
	}

	accBalances, err := bi.ethClient.AccountBalances(ctx, addrs, blockNumber)
	if err != nil {
		return err
	}

	return bi.storage.IndexAccountBalances(ctx, accBalances)
}

func (bi *BlockchainIndexer) indexBlock(ctx context.Context, block *types.Block) error {
	timestamp := time.UnixMilli(int64(block.Time())).UTC().Format(TimeLayOut)
	indexBlock := &store.IndexBlock{
		Number:       block.NumberU64(),
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Root:         block.Root().Hex(),
		Nonce:        block.Nonce(),
		Timestamp:    timestamp,
		Transactions: len(block.Transactions()),
		BaseFee:      block.BaseFee().Uint64(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Difficulty:   block.Difficulty().Uint64(),
		ExtraData:    hex.EncodeToString(block.Extra()),
	}

	return bi.storage.IndexBlock(ctx, indexBlock)
}

func (bi *BlockchainIndexer) indexTransactions(ctx context.Context, block *types.Block) error {
	var transactions []*store.IndexTransaction
	var txHashes []string

	for _, tx := range block.Transactions() {
		from, err := types.Sender(types.NewCancunSigner(tx.ChainId()), tx)
		if err != nil {
			return fmt.Errorf("failed to derive sender: %w", err)
		}

		v, r, s := tx.RawSignatureValues()
		timestamp := tx.Time().UTC().Format(TimeLayOut)
		transaction := &store.IndexTransaction{
			Hash:        tx.Hash().Hex(),
			From:        from.Hex(),
			Gas:         tx.Gas(),
			Nonce:       tx.Nonce(),
			BlockHash:   block.Hash().Hex(),
			BlockNumber: block.NumberU64(),
			ChainId:     tx.ChainId().String(),
			V:           v.String(),
			R:           r.String(),
			S:           s.String(),
			Input:       hex.EncodeToString(tx.Data()),
			Timestamp:   timestamp,
		}

		if tx.To() != nil {
			transaction.To = tx.To().Hex()
		}
		if tx.GasPrice() != nil {
			transaction.GasPrice = tx.GasPrice().Uint64()
		}
		if tx.GasTipCap() != nil {
			transaction.GasTipCap = tx.GasTipCap().Uint64()
		}
		if tx.GasFeeCap() != nil {
			transaction.GasFeeCap = tx.GasFeeCap().Uint64()
		}
		if tx.Value() != nil {
			transaction.Value = tx.Value().String()
		}

		transactions = append(transactions, transaction)
		txHashes = append(txHashes, tx.Hash().Hex())
	}

	receipts, err := bi.fetchReceipts(ctx, txHashes)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction receipts: %w", err)
	}

	for _, tx := range transactions {
		if receipt, ok := receipts[tx.Hash]; ok {
			tx.Status = receipt.Status
			tx.GasUsed = receipt.GasUsed
			tx.CumulativeGasUsed = receipt.CumulativeGasUsed
			tx.ContractAddress = receipt.ContractAddress.Hex()
			tx.TransactionIndex = receipt.TransactionIndex
			tx.ReceiptBlockHash = receipt.BlockHash.Hex()
			tx.ReceiptBlockNumber = receipt.BlockNumber.Uint64()
		}
	}

	return bi.storage.IndexTransactions(ctx, transactions)
}

func (bi *BlockchainIndexer) fetchReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error) {
	return bi.ethClient.TxReceipts(ctx, txHashes)
}

func (bi *BlockchainIndexer) fetchBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error) {
	return bi.ethClient.GetBlocks(ctx, blockNums)
}
