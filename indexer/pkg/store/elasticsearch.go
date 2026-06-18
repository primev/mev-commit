package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync/atomic"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	estypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ESClient struct {
	client      *elasticsearch.TypedClient
	bulkIndexer esutil.BulkIndexer
}

func NewESClient(endpoint string, user, pass string) (*ESClient, error) {
	config := elasticsearch.Config{
		Addresses: []string{endpoint},
		Username:  user,
		Password:  pass,
	}
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new elasticsearch client: %w", err)
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     client,
		NumWorkers: 4,
		FlushBytes: 5e+6, // 5MB
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bulk indexer: %w", err)
	}
	typedClient, err := elasticsearch.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new typed elasticsearch client: %w", err)

	}
	return &ESClient{client: typedClient, bulkIndexer: bi}, nil
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

func (c *ESClient) IndexBlock(ctx context.Context, block *IndexBlock) error {
	return c.Bulk(ctx, "blocks", []interface{}{block})
}

func (c *ESClient) IndexTransactions(ctx context.Context, transactions []*IndexTransaction) error {
	docs := make([]interface{}, len(transactions))
	for i, tx := range transactions {
		docs[i] = tx
	}
	return c.Bulk(ctx, "transactions", docs)
}

func pointerInt(v int) *int {
	return &v
}

func pointerString(v string) *string {
	return &v
}

func (c *ESClient) GetAddresses(ctx context.Context) ([]string, error) {
	query := &search.Request{
		Size: pointerInt(0), // We don't need the actual documents, just the aggregations
		Aggregations: map[string]estypes.Aggregations{
			"unique_from_addresses": {
				Terms: &estypes.TermsAggregation{
					Field: pointerString("from.keyword"),
					Size:  pointerInt(10000),
				},
			},
			"unique_to_addresses": {
				Terms: &estypes.TermsAggregation{
					Field: pointerString("to.keyword"),
					Size:  pointerInt(10000),
				},
			},
		},
	}

	// Execute the search request
	res, err := c.client.Search().
		Index("transactions").
		Request(query).
		Size(1000).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search query: %w", err)
	}

	// Extract unique addresses from the aggregations
	addressSet := make(map[string]struct{})

	// Process the "from" addresses aggregation
	fromAgg, ok := res.Aggregations["unique_from_addresses"]
	if ok {
		buckets := fromAgg.(*estypes.StringTermsAggregate).Buckets
		switch b := buckets.(type) {
		case []estypes.StringTermsBucket:
			for _, bucket := range b {
				address := bucket.Key.(string)
				addressSet[address] = struct{}{}
			}
		}
	}

	// Process the "to" addresses aggregation
	toAgg, ok := res.Aggregations["unique_to_addresses"]
	if ok {
		buckets := toAgg.(*estypes.StringTermsAggregate).Buckets
		switch b := buckets.(type) {
		case []estypes.StringTermsBucket:
			for _, bucket := range b {
				address := bucket.Key.(string)
				addressSet[address] = struct{}{}
			}
		}
	}

	// Combine the unique addresses into a slice
	addresses := make([]string, 0, len(addressSet))
	for address := range addressSet {
		if address != "" && address != "0x" {
			addresses = append(addresses, address)
		}
	}

	return addresses, nil

}

func (c *ESClient) IndexAccountBalances(ctx context.Context, accountBalances []AccountBalance) error {
	docs := make([]interface{}, len(accountBalances))
	for i, accBal := range accountBalances {
		docs[i] = accBal
	}
	return c.Bulk(ctx, "accounts", docs)
}

func (c *ESClient) CreateIndices(ctx context.Context) error {
	indices := []string{"blocks", "transactions", "accounts"}
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
						slog.Error("bulk indexing error", "error", err)
					} else {
						slog.Error("bulk indexing error", "type", res.Error.Type, "reason", res.Error.Reason)
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
