package pricer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
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

type blockPrices struct {
	CurrentBlockNumber int64        `json:"currentBlockNumber"`
	Prices             []BlockPrice `json:"blockPrices"`
}

type BidPricer struct {
	apiKey           string
	log              *slog.Logger
	mu               sync.RWMutex // Protects currentEstimates
	currentEstimates map[int64]float64
}

func NewPricer(apiKey string, logger *slog.Logger) *BidPricer {
	return &BidPricer{
		apiKey: apiKey,
		log:    logger,
	}
}

func (b *BidPricer) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(2 * time.Second) // Adjust the ticker interval as needed
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := b.syncEstimate(ctx); err != nil {
					b.log.Error("Failed to estimate price", "error", err)
				}
			}
		}
	}()
	return done
}

func (b *BidPricer) EstimatePrice(ctx context.Context) map[int64]float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	estimates := make(map[int64]float64)
	for blockNumber, price := range b.currentEstimates {
		estimates[blockNumber] = price
	}
	return estimates
}

func (b *BidPricer) SyncEstimate(ctx context.Context) error {
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

	bp := new(blockPrices)
	if err := json.Unmarshal(respBuf, bp); err != nil {
		return nil, err
	}

	if len(bp.Prices) == 0 {
		return nil, errors.New("no block prices available")
	}

	return bp, nil
}
