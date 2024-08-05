package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"runtime"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

const (
	BlockIndexName       = "blocks"
	AccountIndexName     = "accounts"
	TransactionIndexName = "transactions"
)

type (
	BlockItem struct {
		Number       uint64 `json:"number"`
		Hash         string `json:"hash"`
		ParentHash   string `json:"parentHash"`
		Root         string `json:"root"`
		Nonce        uint64 `json:"nonce"`
		Timestamp    string `json:"timestamp"`
		Transactions int    `json:"transactions"`
		BaseFee      uint64 `json:"baseFee"`
		GasLimit     uint64 `json:"gasLimit"`
		GasUsed      uint64 `json:"gasUsed"`
		Difficulty   uint64 `json:"difficulty"`
		ExtraData    string `json:"extraData"`
	}

	TransactionItem struct {
		Hash               string `json:"hash"`
		From               string `json:"from"`
		To                 string `json:"to"`
		Gas                uint64 `json:"gas"`
		GasPrice           uint64 `json:"gasPrice"`
		GasTipCap          uint64 `json:"gasTipCap"`
		GasFeeCap          uint64 `json:"gasFeeCap"`
		Value              string `json:"value"`
		Nonce              uint64 `json:"nonce"`
		BlockHash          string `json:"blockHash"`
		BlockNumber        uint64 `json:"blockNumber"`
		ChainId            string `json:"chainId"`
		V                  string `json:"v"`
		R                  string `json:"r"`
		S                  string `json:"s"`
		Input              string `json:"input"`
		Timestamp          string `json:"timestamp"`
		Status             uint64 `json:"status"`
		GasUsed            uint64 `json:"gasUsed"`
		CumulativeGasUsed  uint64 `json:"cumulativeGasUsed"`
		ContractAddress    string `json:"contractAddress"`
		TransactionIndex   uint   `json:"transactionIndex"`
		ReceiptBlockHash   string `json:"receiptBlockHash"`
		ReceiptBlockNumber uint64 `json:"receiptBlockNumber"`
	}

	AccountBalanceItem struct {
		Address     string `json:"address"`
		Balance     string `json:"balance"`
		Timestamp   string `json:"timestamp"`
		BlockNumber uint64 `json:"blockNumber"`
	}
)

// SortOrder represents the sort order for search results.
type SortOrder int

// String returns the string representation of the sort order.
// It implements the fmt.Stringer interface.
func (s SortOrder) String() string {
	switch s {
	case SortOrderAsc:
		return "asc"
	case SortOrderDesc:
		return "desc"
	default:
		return "unknown"
	}
}

const (
	SortOrderUnknown SortOrder = iota
	SortOrderAsc
	SortOrderDesc
)

// Store defines an interface for interacting with the backend storage.
type Store interface {
	// CreateIndexes creates all required indexes if they do not exist.
	CreateIndexes(context.Context) error

	// IndexBlock indexes a given block.
	IndexBlock(context.Context, *BlockItem) error

	// IndexTransactions indexes a list of transactions.
	IndexTransactions(ctx context.Context, transactions []*TransactionItem) error

	// LastIndexedBlock returns the number of the last indexed block.
	LastIndexedBlock(context.Context, SortOrder) (*big.Int, error)
}

var (
	_ Store = (*ElasticsearchStore)(nil)
)

// ElasticsearchStore implements the Store interface and
// represents a connection to an Elasticsearch cluster.
type ElasticsearchStore struct {
	client      *elasticsearch.Client
	typedClient *elasticsearch.TypedClient
}

// NewElasticsearchStore creates a new ElasticsearchStore instance connected to the specified URL.
func NewElasticsearchStore(url string, user, pass string) (*ElasticsearchStore, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
		Username:  user,
		Password:  pass,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	typedClient, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticsearchStore{
		client:      client,
		typedClient: typedClient,
	}, nil
}

// CreateIndexes implements the Store interface.
func (e *ElasticsearchStore) CreateIndexes(ctx context.Context) error {
	indices := []string{
		BlockIndexName,
		AccountIndexName,
		TransactionIndexName,
	}

	for _, index := range indices {
		exists, err := e.typedClient.Indices.Exists(index).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to check if index %q exists: %w", index, err)
		}
		if exists {
			continue
		}

		res, err := e.typedClient.Indices.Create(index).
			Settings(&types.IndexSettings{
				NumberOfShards:   "1",
				NumberOfReplicas: "0",
			}).
			Mappings(&types.TypeMapping{
				Properties: map[string]types.Property{
					"timestamp": map[string]interface{}{
						"type":   "date",
						"format": "strict_date_optional_time||epoch_millis",
					},
				},
			}).
			Do(ctx)
		if err != nil {
			return fmt.Errorf("create index %q: %w", index, err)
		}
		if !res.Acknowledged {
			return fmt.Errorf("index %q creation not acknowledged", index)
		}
	}
	return nil
}

// IndexBlock implements the Store interface.
func (e *ElasticsearchStore) IndexBlock(ctx context.Context, block *BlockItem) error {
	return e.bulk(ctx, BlockIndexName, []interface{}{block})
}

// IndexTransactions implements the Store interface.
func (e *ElasticsearchStore) IndexTransactions(ctx context.Context, transactions []*TransactionItem) error {
	docs := make([]interface{}, len(transactions))
	for i, tx := range transactions {
		docs[i] = tx
	}
	return e.bulk(ctx, TransactionIndexName, docs)
}

// LastIndexedBlock implements the Store interface.
func (e *ElasticsearchStore) LastIndexedBlock(ctx context.Context, order SortOrder) (*big.Int, error) {
	const index = BlockIndexName

	exists, err := e.typedClient.Indices.Exists(index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("check if index %q exists: %w", index, err)
	}
	if !exists {
		return big.NewInt(0), nil
	}

	cntRes, err := e.typedClient.Count().Index(index).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("count documents in %q index: %w", index, err)
	}
	if cntRes.Count == 0 {
		return big.NewInt(0), nil
	}

	schRes, err := e.typedClient.Search().
		Index(index).
		Sort(map[string]interface{}{
			"number": map[string]interface{}{
				"order": order.String(),
			},
		}).
		Size(1).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("search for last indexed block: %w", err)
	}
	if schRes.Hits.Total.Value == 0 {
		return big.NewInt(0), nil
	}

	var block struct {
		Number uint64 `json:"number"`
	}
	if err := json.Unmarshal(schRes.Hits.Hits[0].Source_, &block); err != nil {
		return nil, fmt.Errorf("unmarshal search result: %w", err)
	}
	return new(big.Int).SetUint64(block.Number), nil
}

func (e *ElasticsearchStore) bulk(ctx context.Context, index string, docs []interface{}) (err error) {
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     e.client,
		NumWorkers: runtime.GOMAXPROCS(0),
		FlushBytes: 5e+6,
	})
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, indexer.Close(ctx))
	}()

	var errs error
	for _, doc := range docs {
		data, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("marshal document: %w", err)
		}
		err = indexer.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action: "index",
				Index:  index,
				Body:   bytes.NewReader(data),
				OnFailure: func(
					_ context.Context,
					_ esutil.BulkIndexerItem,
					_ esutil.BulkIndexerResponseItem, err error,
				) {
					errs = errors.Join(errs, fmt.Errorf("index document %s: %w", string(data), err))
				},
			},
		)
		if err != nil {
			return errors.Join(errs, fmt.Errorf("add document to bulk indexer: %w", err))
		}
	}
	return errs
}
