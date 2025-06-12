package pricer

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

var apiURL = "https://api.blocknative.com/gasprices/blockprices?chainid=1"

type blockPrice struct {
	CurrentBlockNumber int64 `json:"currentBlockNumber"`
	BlockPrices        []struct {
		BlockNumber     int64 `json:"blockNumber"`
		EstimatedPrices []struct {
			Confidence        int     `json:"confidence"`
			PriorityFeePerGas float64 `json:"maxPriorityFeePerGas"`
		}
	}
}

type BlockPrice struct {
	BlockNumber int64
	BidAmount   *big.Int
}

type BidPricer struct{}

func (b *BidPricer) EstimatePrice(ctx context.Context, txn *types.Transaction) (*BlockPrice, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
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

	var bp blockPrice
	if err := json.Unmarshal(respBuf, &bp); err != nil {
		return nil, err
	}

	if len(bp.BlockPrices) == 0 {
		return nil, errors.New("no block prices available")
	}

	for _, price := range bp.BlockPrices {
		if price.BlockNumber == bp.CurrentBlockNumber+1 {
			for _, p := range price.EstimatedPrices {
				if p.Confidence == 99 { // Assuming we want the 99% confidence price
					// Convert the priority fee from Gwei to Wei
					// 1 Gwei = 1e9 Wei
					priorityFee := p.PriorityFeePerGas * 1e9
					bidAmount := big.NewInt(0).Mul(big.NewInt(int64(priorityFee)), big.NewInt(int64(txn.Gas())))
					return &BlockPrice{BlockNumber: price.BlockNumber, BidAmount: bidAmount}, nil
				}
			}
		}
	}

	// If we reach here, it means we didn't find a suitable price.
	// This could happen if the API response format changes or if no 99% confidence price is available.
	return nil, errors.New("no suitable price found for the next block")
}
