package points

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const weiPerPoint = 33_333_333_333_333

type pointsTracker struct {
	apiURL string
	apiKey string
	client *http.Client
	logger *slog.Logger
}

func NewPointsTracker(apiURL, apiKey string, logger *slog.Logger) *pointsTracker {
	return &pointsTracker{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				MaxIdleConns:        256,
				MaxIdleConnsPerHost: 256,
				IdleConnTimeout:     90 * time.Second,
				ForceAttemptHTTP2:   true,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
			Timeout: 15 * time.Second,
		},
		logger: logger,
	}
}

func (p *pointsTracker) AssignPoints(
	ctx context.Context,
	userID common.Address,
	transactionHash common.Hash,
	mevRevenue *big.Int,
) error {
	// upper bound of mev_revenue/weiPerPoint
	points := new(big.Int).Div(mevRevenue, big.NewInt(weiPerPoint))
	if points.Cmp(big.NewInt(0)) <= 0 {
		p.logger.Info("no points to assign", "user", userID.Hex(), "tx", transactionHash.Hex(), "mev_revenue", mevRevenue.String())
		return nil
	}
	p.logger.Info("assigning points", "user", userID.Hex(), "tx", transactionHash.Hex(), "points", points.String())

	reqBody := map[string]any{
		"user": map[string]any{
			"identifier_type": "evm_address",
			"identifier":      userID.Hex(),
		},
		"name": "fast_miles",
		"args": map[string]any{
			"value": map[string]any{
				"amount": points.String(),
				"currency": map[string]any{
					"name": "POINT",
				},
			},
			"transaction_hash": transactionHash.Hex(),
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", p.apiURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to assign points, status code: %d", resp.StatusCode)
	}
	p.logger.Info("successfully assigned points", "user", userID.Hex(), "tx", transactionHash.Hex(), "points", points.String())

	return nil
}
