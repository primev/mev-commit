package explorersubmitter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Endpoint string
	ApiKey   string
	AppCode  string
}

type TxData struct {
	ChainID  string `json:"df_chainid"`
	AppCode  string `json:"df_appcode"`
	TxHash   string `json:"df_txhash"`
	ExpireTs int64  `json:"df_expire_ts"`
	TxInfo   TxInfo `json:"df_txinfo"`
}

type TxInfo struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func Submit(ctx context.Context, config Config, chainID string, txHash string, from string, to string) error {
	if config.Endpoint == "" {
		return nil
	}

	expireTs := time.Now().Add(15 * time.Minute).Unix()

	txData := TxData{
		ChainID:  chainID,
		AppCode:  config.AppCode,
		TxHash:   txHash,
		ExpireTs: expireTs,
		TxInfo: TxInfo{
			From: from,
			To:   to,
		},
	}

	jsonData, err := json.Marshal(txData)
	if err != nil {
		return fmt.Errorf("failed to marshal json data: %w", err)
	}

	params := url.Values{}
	params.Add("apikey", config.ApiKey)
	params.Add("action", "submitTxPending")
	params.Add("JsonData", string(jsonData))

	reqURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %w", err)
	}

	reqURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-ok status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
