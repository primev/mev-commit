package pricer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

var apiURL = "https://api.blocknative.com/gasprices/blockprices?chainid=1"

type EstimatedPrice struct {
	Confidence            int     `json:"confidence"`
	PriorityFeePerGasGwei float64 `json:"maxPriorityFeePerGas"`
}

type BlockPrice struct {
	BlockNumber     int64            `json:"blockNumber"`
	EstimatedPrices []EstimatedPrice `json:"estimatedPrices"`
}

type BlockPrices struct {
	MsSinceLastBlock   int64        `json:"msSinceLastBlock"`
	CurrentBlockNumber int64        `json:"currentBlockNumber"`
	Prices             []BlockPrice `json:"blockPrices"`
}

type BidPricer struct {
	apiKey string
}

func NewPricer(apiKey string) *BidPricer {
	return &BidPricer{
		apiKey: apiKey,
	}
}

func (b *BidPricer) EstimatePrice(ctx context.Context) (*BlockPrices, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	if b.apiKey != "" {
		req.Header.Set("Authorization", b.apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch price estimate: " + resp.Status)
	}

	respBuf, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, err
	}

	bp := new(BlockPrices)
	if err := json.Unmarshal(respBuf, bp); err != nil {
		return nil, err
	}

	if len(bp.Prices) == 0 {
		return nil, errors.New("no block prices available")
	}

	return bp, nil
}
