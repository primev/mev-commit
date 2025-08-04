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
	apiKey             string
	log                *slog.Logger
	mu                 sync.RWMutex // Protects currentEstimates
	currentEstimates   map[int64]float64
	currentBlockNumber int64
}

func NewPricer(apiKey string, logger *slog.Logger) (*BidPricer, error) {
	bp := &BidPricer{
		apiKey:           apiKey,
		log:              logger,
		currentEstimates: make(map[int64]float64),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := bp.syncEstimate(ctx); err != nil {
		return nil, err
	}
	return bp, nil
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
				if err := b.syncEstimate(ctx); err != nil {
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
	for confidence, price := range b.currentEstimates {
		estimates[confidence] = price
	}
	return estimates
}

func (b *BidPricer) syncEstimate(ctx context.Context) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return err
	}

	if b.apiKey != "" {
		req.Header.Set("Authorization", b.apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch price estimate: " + resp.Status)
	}

	respBuf, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return err
	}

	bp := new(blockPrices)
	if err := json.Unmarshal(respBuf, bp); err != nil {
		return err
	}

	if len(bp.Prices) == 0 {
		return errors.New("no block prices available")
	}

	if b.currentBlockNumber < bp.CurrentBlockNumber+1 {
		for _, price := range bp.Prices {
			if price.BlockNumber == bp.CurrentBlockNumber+1 {
				b.mu.Lock()
				for _, estimatedPrice := range price.EstimatedPrices {
					switch estimatedPrice.Confidence {
					case 90:
						b.currentEstimates[int64(estimatedPrice.Confidence)] = estimatedPrice.PriorityFeePerGasGwei
					case 95:
						b.currentEstimates[int64(estimatedPrice.Confidence)] = estimatedPrice.PriorityFeePerGasGwei
					case 99:
						b.currentEstimates[int64(estimatedPrice.Confidence)] = estimatedPrice.PriorityFeePerGasGwei
					}
				}
				b.currentBlockNumber = price.BlockNumber
				b.mu.Unlock()
				b.log.Debug("Updated current estimates", "blockNumber", price.BlockNumber, "estimates", b.currentEstimates)
			}
		}
	}

	return nil
}
