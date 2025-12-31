package explorersubmitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

type explorerSubmitter struct {
	endpoint string
	apiKey   string
	appCode  string
	client   *http.Client
	reqChan  chan submitRequest
	logger   *slog.Logger
}

type submitRequest struct {
	tx   *types.Transaction
	from common.Address
}

func New(endpoint, apiKey, appCode string, logger *slog.Logger) *explorerSubmitter {
	return &explorerSubmitter{
		endpoint: endpoint,
		apiKey:   apiKey,
		appCode:  appCode,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		reqChan: make(chan submitRequest, 100),
		logger:  logger,
	}
}

func (e *explorerSubmitter) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case req := <-e.reqChan:
				if err := e.doSubmit(egCtx, req); err != nil {
					e.logger.Error("failed to submit to explorer", "error", err)
				}
			}
		}
	})

	go func() {
		defer close(done)
		if err := eg.Wait(); err != nil {
			if errors.Is(err, context.Canceled) {
				e.logger.Info("Explorer submitter stopped")
			} else {
				e.logger.Error("Explorer submitter exited with error", "error", err)
			}
		}
	}()

	return done
}

func (e *explorerSubmitter) Submit(ctx context.Context, tx *types.Transaction, from common.Address) error {
	select {
	case e.reqChan <- submitRequest{tx: tx, from: from}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Drop request if channel is full to prevent blocking
		e.logger.Warn("request channel full, dropping request", "txHash", tx.Hash().Hex())
		return nil
	}
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

func (e *explorerSubmitter) doSubmit(ctx context.Context, req submitRequest) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in doSubmit: %v", r)
		}
	}()

	if e.endpoint == "" {
		return nil
	}

	chainID := "1"
	expireTs := time.Now().Add(15 * time.Minute).Unix()

	txHash := req.tx.Hash().Hex()
	from := req.from.Hex()
	to := ""
	if req.tx.To() != nil {
		to = req.tx.To().Hex()
	}

	txData := TxData{
		ChainID:  chainID,
		AppCode:  e.appCode,
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
	params.Add("apikey", e.apiKey)
	params.Add("action", "submitTxPending")
	params.Add("JsonData", string(jsonData))

	reqURL, err := url.Parse(e.endpoint)
	if err != nil {
		return fmt.Errorf("failed to parse endpoint url: %w", err)
	}

	reqURL.RawQuery = params.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := e.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request to explorer endpoint: %v", err.(*url.Error).Err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("received non-ok status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	e.logger.Debug("Successfully submitted tx to explorer", "hash", txHash)
	return nil
}
