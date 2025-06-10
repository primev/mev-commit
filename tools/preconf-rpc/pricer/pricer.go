package pricer

import (
	"context"
	"errors"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

var apiURL = "https://api.blocknative.com/gasprices/blockprices?chainid=1"

type bidPricer struct {
	apiURL string
}

func (b *bidPricer) EstimatePrice(ctx context.Context, txn *types.Transaction) (*big.Int, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", b.apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch price estimate: " + resp.Status)
	}

	respBuf, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, err
	}

	return big.NewInt(1000000000), nil // Return a dummy value for now.
}
