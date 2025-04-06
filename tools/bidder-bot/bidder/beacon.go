package bidder

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type beaconClient struct {
	apiURL string
	logger *slog.Logger
	client *http.Client
}

func newBeaconClient(apiURL string, logger *slog.Logger) *beaconClient {
	return &beaconClient{
		apiURL: apiURL,
		logger: logger,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type beaconHeaderResponse struct {
	Data struct {
		Header struct {
			Message struct {
				Slot string `json:"slot"`
			} `json:"message"`
		} `json:"header"`
	} `json:"data"`
}

type beaconBlockResponse struct {
	Data struct {
		Message struct {
			Body struct {
				ExecutionPayload struct {
					BlockNumber string `json:"block_number"`
				} `json:"execution_payload"`
			} `json:"body"`
		} `json:"message"`
	} `json:"data"`
}

func (bc *beaconClient) getLatestSlot(ctx context.Context) (uint64, error) {
	url := fmt.Sprintf("%s/eth/v1/beacon/headers/head", bc.apiURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		bc.logger.Error("creating request for current slot", "error", err)
		return 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := bc.client.Do(req)
	if err != nil {
		bc.logger.Error("failed to execute request for current slot", "error", err)
		return 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bc.logger.Error("unexpected status code for current slot", "status", resp.StatusCode)
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var headerResp beaconHeaderResponse
	if err := json.NewDecoder(resp.Body).Decode(&headerResp); err != nil {
		bc.logger.Error("failed to decode response for current slot", "error", err)
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	slot, err := strconv.ParseUint(headerResp.Data.Header.Message.Slot, 10, 64)
	if err != nil {
		bc.logger.Error("failed to parse slot", "error", err)
		return 0, fmt.Errorf("parsing slot: %w", err)
	}

	bc.logger.Debug("retrieved current beacon slot", "slot", slot)
	return slot, nil
}

func (bc *beaconClient) getBlockNumForSlot(ctx context.Context, slot uint64) (uint64, error) {
	url := fmt.Sprintf("%s/eth/v2/beacon/blocks/%d", bc.apiURL, slot)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		bc.logger.Error("creating request for block number", "error", err)
		return 0, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := bc.client.Do(req)
	if err != nil {
		bc.logger.Error("failed to execute request for block number", "error", err)
		return 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bc.logger.Error("unexpected status code for block number", "status", resp.StatusCode, "slot", slot)
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var blockResp beaconBlockResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
		bc.logger.Error("failed to decode response for block number", "error", err)
		return 0, fmt.Errorf("decoding response: %w", err)
	}

	blockNumber, err := strconv.ParseUint(blockResp.Data.Message.Body.ExecutionPayload.BlockNumber, 10, 64)
	if err != nil {
		bc.logger.Error("failed to parse block number", "error", err)
		return 0, fmt.Errorf("parsing block number: %w", err)
	}

	bc.logger.Debug("retrieved block number for slot", "slot", slot, "blockNumber", blockNumber)
	return blockNumber, nil
}
