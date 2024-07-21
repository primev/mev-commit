package indexer

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lmittmann/tint"
	"github.com/lmittmann/w3/w3types"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
)

type BlockchainIndexer struct {
	ethClient                EthereumClient
	esClient                 ElasticsearchClient
	blockChan                chan *types.Block
	txChan                   chan *types.Transaction
	indexInterval            time.Duration
	lastForwardIndexedBlock  *big.Int
	lastBackwardIndexedBlock *big.Int
	logger                   *slog.Logger
}

type EthereumClient interface {
	GetBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	TxReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error)
}

type ElasticsearchClient interface {
	Index(ctx context.Context, index string, document interface{}) error
	Search(ctx context.Context, index string, query *estypes.Query) (*search.Response, error)
	GetLastIndexedBlock(ctx context.Context, direction string) (*big.Int, error)
	CreateIndices(ctx context.Context) error
	Bulk(ctx context.Context, indexName string, docs []interface{}) error
}

func NewBlockchainIndexer(ethClient EthereumClient, esClient ElasticsearchClient, indexInterval time.Duration) *BlockchainIndexer {
	return &BlockchainIndexer{
		ethClient:     ethClient,
		esClient:      esClient,
		blockChan:     make(chan *types.Block, 100),
		txChan:        make(chan *types.Transaction, 100),
		indexInterval: indexInterval,
		logger:        slog.Default(),
	}
}

func (bi *BlockchainIndexer) Start(ctx context.Context) error {
	if err := bi.esClient.CreateIndices(ctx); err != nil {
		return fmt.Errorf("failed to create indices: %w", err)
	}

	latestBlockNumber, err := bi.ethClient.BlockNumber(ctx)
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
	go bi.fetchBackwardBlocks(ctx)
	go bi.processBlocks(ctx)

	// Block the main function indefinitely
	select {}
}

func (bi *BlockchainIndexer) initializeForwardIndex(ctx context.Context, latestBlockNumber uint64) error {
	lastForwardIndexedBlock, err := bi.esClient.GetLastIndexedBlock(ctx, "forward")
	if err != nil {
		return fmt.Errorf("failed to get last forward indexed block: %w", err)
	}
	bi.logger.Info("last indexed block", "blockNumber", lastForwardIndexedBlock, "direction", "forward")

	if lastForwardIndexedBlock == nil || lastForwardIndexedBlock.Cmp(big.NewInt(0)) == 0 {
		bi.lastForwardIndexedBlock = new(big.Int).SetUint64(latestBlockNumber - 1)
	} else {
		bi.lastForwardIndexedBlock = lastForwardIndexedBlock
	}

	return nil
}

func (bi *BlockchainIndexer) initializeBackwardIndex(ctx context.Context, latestBlockNumber uint64) error {
	lastBackwardIndexedBlock, err := bi.esClient.GetLastIndexedBlock(ctx, "backward")
	if err != nil {
		return fmt.Errorf("failed to get last backward indexed block: %w", err)
	}
	bi.logger.Info("last indexed block", "blockNumber", lastBackwardIndexedBlock, "direction", "backward")

	if lastBackwardIndexedBlock == nil || lastBackwardIndexedBlock.Cmp(big.NewInt(0)) == 0 {
		bi.lastBackwardIndexedBlock = new(big.Int).SetUint64(latestBlockNumber - 1)
	} else {
		bi.lastBackwardIndexedBlock = new(big.Int).Sub(lastBackwardIndexedBlock, big.NewInt(1))
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
				bi.logger.Error("Failed to get latest block number", "error", err)
				continue
			}

			for blockNum := new(big.Int).Add(bi.lastForwardIndexedBlock, big.NewInt(1)); blockNum.Cmp(latestBlockNumber) <= 0; blockNum.Add(blockNum, big.NewInt(5)) {
				endBlockNum := new(big.Int).Add(blockNum, big.NewInt(4))
				if endBlockNum.Cmp(latestBlockNumber) > 0 {
					endBlockNum.Set(latestBlockNumber)
				}

				blockNums := []*big.Int{}
				for bn := new(big.Int).Set(blockNum); bn.Cmp(endBlockNum) <= 0; bn.Add(bn, big.NewInt(1)) {
					blockNums = append(blockNums, new(big.Int).Set(bn))
				}

				blocks, err := bi.fetchBlocks(ctx, blockNums)
				if err != nil {
					bi.logger.Error("Failed to fetch blocks", "start", blockNum, "end", endBlockNum, "error", err)
					continue
				}

				for _, block := range blocks {
					bi.blockChan <- block
					bi.lastForwardIndexedBlock.Set(block.Number())
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
			if bi.lastBackwardIndexedBlock.Sign() < 0 {
				return
			}
			block, err := bi.ethClient.BlockByNumber(ctx, bi.lastBackwardIndexedBlock)
			if err != nil {
				bi.logger.Error("Failed to get block", "number", bi.lastBackwardIndexedBlock, "error", err)
				continue
			}
			bi.blockChan <- block
			bi.lastBackwardIndexedBlock.Set(new(big.Int).Sub(bi.lastBackwardIndexedBlock, big.NewInt(1)))
		}
	}
}

func (bi *BlockchainIndexer) processBlocks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case block := <-bi.blockChan:
			if err := bi.indexBlock(ctx, block); err != nil {
				bi.logger.Error("Failed to index block", "error", err)
			}
			if err := bi.indexTransactions(ctx, block); err != nil {
				bi.logger.Error("Failed to index transactions", "error", err)
			}
		}
	}
}

func (bi *BlockchainIndexer) indexBlock(ctx context.Context, block *types.Block) error {
	timestamp := time.UnixMilli(int64(block.Time())).UTC().Format("2006-01-02T15:04:05.000Z")
	blockDoc := map[string]interface{}{
		"number":       block.NumberU64(),
		"hash":         block.Hash().Hex(),
		"parentHash":   block.ParentHash().Hex(),
		"root":         block.Root().Hex(),
		"nonce":        block.Nonce(),
		"timestamp":    timestamp,
		"transactions": len(block.Transactions()),
		"baseFee":      block.BaseFee().Uint64(),
		"gasLimit":     block.GasLimit(),
		"gasUsed":      block.GasUsed(),
		"difficulty":   block.Difficulty().Uint64(),
		"extraData":    hex.EncodeToString(block.Extra()),
	}

	if err := bi.esClient.Index(ctx, "blocks", blockDoc); err != nil {
		return fmt.Errorf("failed to index block: %w", err)
	}
	return nil
}

func (bi *BlockchainIndexer) indexTransactions(ctx context.Context, block *types.Block) error {
	txDocs := make([]interface{}, 0, len(block.Transactions()))

	var txHashes []string
	for _, tx := range block.Transactions() {
		from, err := types.Sender(types.NewCancunSigner(tx.ChainId()), tx)
		if err != nil {
			return fmt.Errorf("failed to derive sender: %w", err)
		}

		v, r, s := tx.RawSignatureValues()
		timestamp := tx.Time().UTC().Format("2006-01-02T15:04:05.000Z")
		txDoc := map[string]interface{}{
			"hash":        tx.Hash().Hex(),
			"from":        from.Hex(),
			"gas":         tx.Gas(),
			"nonce":       tx.Nonce(),
			"blockHash":   block.Hash().Hex(),
			"blockNumber": block.NumberU64(),
			"chainId":     tx.ChainId().String(),
			"v":           v.String(),
			"r":           r.String(),
			"s":           s.String(),
			"input":       hex.EncodeToString(tx.Data()),
			"timestamp":   timestamp,
		}

		if tx.To() != nil {
			txDoc["to"] = tx.To().Hex()
		} else {
			txDoc["to"] = ""
		}

		if tx.GasPrice() != nil {
			txDoc["gasPrice"] = tx.GasPrice().Uint64()
		}

		if tx.GasTipCap() != nil {
			txDoc["gasTipCap"] = tx.GasTipCap().Uint64()
		}

		if tx.GasFeeCap() != nil {
			txDoc["gasFeeCap"] = tx.GasFeeCap().Uint64()
		}

		if tx.Value() != nil {
			txDoc["value"] = tx.Value().Uint64()
		}
		txDocs = append(txDocs, txDoc)
		txHashes = append(txHashes, tx.Hash().Hex())
	}

	receipts, err := bi.fetchReceipts(ctx, txHashes)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction receipts: %w", err)
	}
	// Add receipt information to transaction documents
	for _, txDoc := range txDocs {
		txD := txDoc.(map[string]interface{})
		txHash := txD["hash"].(string)
		if receipt, ok := receipts[txHash]; ok {
			txD["status"] = receipt.Status
			txD["gasUsed"] = receipt.GasUsed
			txD["cumulativeGasUsed"] = receipt.CumulativeGasUsed
			txD["receiptContractAddress"] = receipt.ContractAddress.Hex()
			txD["transactionIndex"] = receipt.TransactionIndex
			txD["receiptBlockHash"] = receipt.BlockHash
			txD["receiptBlockNumber"] = receipt.BlockNumber.Uint64()
			txD["logs"] = receipt.Logs
		}
	}

	if err := bi.esClient.Bulk(ctx, "transactions", txDocs); err != nil {
		return fmt.Errorf("bulk indexing of transactions failed: %w", err)
	}
	return nil
}

func (bi *BlockchainIndexer) fetchReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error) {
	return bi.ethClient.TxReceipts(ctx, txHashes)
}

func (bi *BlockchainIndexer) fetchBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error) {
	return bi.ethClient.GetBlocks(ctx, blockNums)
}

type W3EvmClient struct {
	client *w3.Client
}

func NewW3EthereumClient(endpoint string) (*W3EvmClient, error) {
	client, err := w3.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	return &W3EvmClient{client: client}, nil
}

func (c *W3EvmClient) GetBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error) {
	batchBlocksCaller := make([]w3types.RPCCaller, len(blockNums))
	blocks := make([]types.Block, len(blockNums))
	for i, blockNum := range blockNums {
		batchBlocksCaller[i] = eth.BlockByNumber(blockNum).Returns(&blocks[i])
	}
	err := c.client.Call(batchBlocksCaller...)
	if err != nil {
		return nil, err
	}
	var b []*types.Block
	for _, block := range blocks {
		b = append(b, &block)
	}
	return b, nil
}

func (c *W3EvmClient) TxReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error) {
	batchTxReceiptCaller := make([]w3types.RPCCaller, len(txHashes))
	txReceipts := make([]types.Receipt, len(txHashes))
	for i, txHash := range txHashes {
		batchTxReceiptCaller[i] = eth.TxReceipt(w3.H(txHash)).Returns(&txReceipts[i])
	}
	err := c.client.Call(batchTxReceiptCaller...)
	if err != nil {
		return map[string]*types.Receipt{}, nil
	}
	txHashToReceipt := make(map[string]*types.Receipt)
	for _, txReceipt := range txReceipts {
		txHashToReceipt[txReceipt.TxHash.Hex()] = &txReceipt
	}
	return txHashToReceipt, nil
}

func (c *W3EvmClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var block types.Block
	if err := c.client.Call(eth.BlockByNumber(number).Returns(&block)); err != nil {
		return nil, err
	}
	return &block, nil
}

func (c *W3EvmClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber big.Int
	if err := c.client.Call(eth.BlockNumber().Returns(&blockNumber)); err != nil {
		return nil, err
	}
	return &blockNumber, nil
}

type ESClient struct {
	client      *elasticsearch.TypedClient
	bulkIndexer esutil.BulkIndexer
}

func NewESClient(endpoint string) (*ESClient, error) {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{endpoint},
		Username:  "elastic",
		Password:  "mev-commit",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: &elasticsearch.Client{
			BaseClient: client.BaseClient,
			API:        esapi.New(client.Transport),
		},
		NumWorkers: 4,
		FlushBytes: 5e+6, // 5MB
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bulk indexer: %w", err)
	}

	return &ESClient{client: client, bulkIndexer: bi}, nil
}

func (c *ESClient) Index(ctx context.Context, index string, document interface{}) error {
	_, err := c.client.Index(index).Document(document).Do(ctx)
	return err
}

func (c *ESClient) Search(ctx context.Context, index string, query *estypes.Query) (*search.Response, error) {
	return c.client.Search().Index(index).Query(query).Do(ctx)
}

func (c *ESClient) GetLastIndexedBlock(ctx context.Context, direction string) (*big.Int, error) {
	// Check if the index exists
	exists, err := c.client.Indices.Exists("blocks").Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if index exists: %w", err)
	}

	if !exists {
		return big.NewInt(0), nil
	}

	// Check if the index contains any documents
	countRes, err := c.client.Count().Index("blocks").Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count documents in index: %w", err)
	}

	if countRes.Count == 0 {
		return big.NewInt(0), nil
	}

	var sortOrder string
	if direction == "forward" {
		sortOrder = "desc"
	} else if direction == "backward" {
		sortOrder = "asc"
	} else {
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}

	// Perform the search query
	res, err := c.client.Search().
		Index("blocks").
		Sort(map[string]interface{}{
			"number": map[string]interface{}{
				"order": sortOrder,
			},
		}).
		Size(1).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	// Check if there are no hits (index exists but no documents)
	if res.Hits.Total.Value == 0 {
		return big.NewInt(0), nil
	}

	var block struct {
		Number uint64 `json:"number"`
	}

	if err := json.Unmarshal(res.Hits.Hits[0].Source_, &block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %w", err)
	}
	blockNumber := new(big.Int).SetUint64(block.Number)
	return blockNumber, nil
}

func (c *ESClient) CreateIndices(ctx context.Context) error {
	indices := []string{"blocks", "transactions"}
	for _, index := range indices {
		res, err := c.client.Indices.Exists(index).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to check if index %s exists: %w", index, err)
		}

		if !res {
			indexSettings := esapi.IndicesCreateRequest{
				Index: index,
				Body: strings.NewReader(`{
					"settings": {
						"number_of_shards": 1,
						"number_of_replicas": 0
					},
					"mappings": {
						"properties": {
							"timestamp": {
								"type": "date",
								"format": "strict_date_optional_time||epoch_millis"
							}
						}
					}
				}`),
			}

			createRes, err := indexSettings.Do(ctx, c.client)
			if err != nil {
				return fmt.Errorf("failed to create index %s: %w", index, err)
			}
			defer createRes.Body.Close()

			if createRes.IsError() {
				return fmt.Errorf("error creating index %s: %s", index, createRes.String())
			}
		}
	}
	return nil
}

func (c *ESClient) Bulk(ctx context.Context, indexName string, docs []interface{}) error {
	var (
		countSuccessful uint64
		countFailed     uint64
	)

	for _, doc := range docs {
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal document: %w", err)
		}

		err = c.bulkIndexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action: "index",
				Index:  indexName,
				Body:   bytes.NewReader(data),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					atomic.AddUint64(&countFailed, 1)
					if err != nil {
						slog.Error("Bulk indexing error", "error", err)
					} else {
						slog.Error("Bulk indexing error", "type", res.Error.Type, "reason", res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			return fmt.Errorf("failed to add document to bulk indexer: %w", err)
		}
	}
	return nil
}

func (c *ESClient) Close(ctx context.Context) error {
	return c.bulkIndexer.Close(ctx)
}

func SetLogLevel(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	opts := &tint.Options{
		Level:      logLevel,
		TimeFormat: time.Kitchen, // Optional: Customize the time format
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
