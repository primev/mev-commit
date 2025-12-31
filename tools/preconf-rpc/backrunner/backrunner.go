package backrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var builders = map[common.Address]string{
	common.HexToAddress("0xE3d71EF44D20917b93AA93e12Bd35b0859824A8F"): "btcs",
	common.HexToAddress("0x2445e5e28890De3e93F39fCA817639c470F4d3b9"): "iobuilder",
	common.HexToAddress("0xB3998135372F1eE16Cb510af70ed212b5155Af62"): "titan",
	common.HexToAddress("0x570e531fB805B5eEbD5F29Eaa2766fBeB4977ddE"): "quasar",
}

type Store interface {
	AddSwapInfo(ctx context.Context, txnHash common.Hash, blockNumber int64, builders []string) error
	GetStartHintForRewards(ctx context.Context) (int64, error)
	UpdateSwapReward(ctx context.Context, reward *big.Int, bundle []string) (bool, error)
	GetSwapRewardee(ctx context.Context, bundle []string) (common.Address, common.Hash, error)
}

type BlockNumberGetter interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type PointsTracker interface {
	AssignPoints(ctx context.Context, userID common.Address, transactionHash common.Hash, mevRevenue *big.Int) error
}

type backrunner struct {
	client  *http.Client
	rpcURL  string
	apiURL  string
	apiKey  string
	store   Store
	points  PointsTracker
	reqChan chan backrunRequest
	metrics *metrics
	logger  *slog.Logger
}

func New(apiKey, apiURL, rpcURL string, store Store, points PointsTracker, logger *slog.Logger) (*backrunner, error) {
	urlParsed, err := url.Parse(rpcURL)
	if err != nil {
		return nil, err
	}
	q := urlParsed.Query()
	q.Add("api_key", apiKey)
	q.Add("mode", "primev")
	urlParsed.RawQuery = q.Encode()

	apiParsed, err := url.Parse(fmt.Sprintf("%s/api/transactions", apiURL))
	if err != nil {
		return nil, err
	}

	return &backrunner{
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
		rpcURL:  urlParsed.String(),
		apiURL:  apiParsed.String(),
		apiKey:  apiKey,
		store:   store,
		points:  points,
		logger:  logger,
		reqChan: make(chan backrunRequest, 100),
		metrics: newMetrics(),
	}, nil
}

func (b *backrunner) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		b.metrics.attempts,
		b.metrics.success,
		b.metrics.fail,
		b.metrics.latency,
		b.metrics.rewards,
		b.metrics.rewardsTotal,
	}
}

type backrunRequest struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
}

func (b *backrunRequest) Builders() []string {
	params, ok := b.Params.([]any)
	if !ok || len(params) == 0 {
		return nil
	}
	paramsMap, ok := params[0].(map[string]any)
	if !ok {
		return nil
	}
	builders, ok := paramsMap["trustedBuilders"].([]string)
	if !ok {
		return nil
	}
	return builders
}

func (b *backrunRequest) String() string {
	buf := bytes.NewBuffer(nil)
	_ = json.NewEncoder(buf).Encode(b)
	return buf.String()
}

func newReq(id int, rawTx string, cmts []*bidderapiv1.Commitment) (backrunRequest, error) {
	if len(cmts) == 0 {
		return backrunRequest{}, errors.New("no commitments provided")
	}

	if id == 0 {
		id = 1
	}

	blkNo := cmts[0].BlockNumber
	var buildersSelected []string
	for _, cmt := range cmts {
		if cmt.BlockNumber == blkNo {
			bldr, found := builders[common.HexToAddress(cmt.ProviderAddress)]
			if found {
				buildersSelected = append(buildersSelected, bldr)
			}
		}
	}

	if len(buildersSelected) == 0 {
		return backrunRequest{}, errors.New("no known builders in commitments")
	}

	return backrunRequest{
		Version: "2.0",
		Method:  "eth_sendBundle",
		Params: []any{
			map[string]any{
				"txs":             []string{rawTx},
				"blockNumber":     fmt.Sprintf("0x%x", blkNo),
				"trustedBuilders": buildersSelected,
			},
		},
		ID: id,
	}, nil
}

func (b *backrunner) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-ticker.C:
				start, err := b.store.GetStartHintForRewards(egCtx)
				if err != nil {
					continue
				}
				if err := b.checkRewards(egCtx, start); err != nil {
					b.logger.Error("checking backrun rewards", "error", err)
					continue
				}
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case req := <-b.reqChan:
				if err := b.doBackrun(egCtx, req); err != nil {
					b.logger.Error("doing backrun", "error", err)
				}
			}
		}
	})

	go func() {
		defer close(done)

		if err := eg.Wait(); err != nil {
			b.logger.Error("backrunner exited with error", "error", err)
		}
	}()

	return done
}

type transactionRecord struct {
	Amount       string   `json:"amount"`
	BundleId     string   `json:"bundleId"`
	BundleHashes []string `json:"bundleHashes"`
}

type transactionRecords struct {
	Records []transactionRecord `json:"records"`
}

type transactionsResponse struct {
	Success bool               `json:"success"`
	Data    transactionRecords `json:"data"`
}

func (b *backrunner) checkRewards(ctx context.Context, start int64) error {
	reqURL, err := url.Parse(b.apiURL)
	if err != nil {
		return fmt.Errorf("parsing backrun API URL: %w", err)
	}
	q := reqURL.Query()
	q.Add("chainId", "1")
	q.Add("revenueType", "Backrun")
	q.Add("start", fmt.Sprintf("%d", start))
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("creating backrun HTTP request: %w", err)
	}

	req.Header.Set("Authorization", b.apiKey)
	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending backrun HTTP request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading backrun HTTP response: %w", err)
	}

	var respStruct transactionsResponse
	if err := json.Unmarshal(respBody, &respStruct); err != nil {
		return fmt.Errorf("unmarshaling backrun HTTP response: %w", err)
	}

	if !respStruct.Success {
		return fmt.Errorf("unsuccessful backrun API response: %s", string(respBody))
	}

	for _, record := range respStruct.Data.Records {
		amount, ok := new(big.Int).SetString(record.Amount, 10)
		if !ok {
			continue
		}
		// 90% to user, 10% to platform
		amount = new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(90)), big.NewInt(100))
		if updated, err := b.store.UpdateSwapReward(ctx, amount, record.BundleHashes); err != nil {
			return fmt.Errorf("updating backrun reward: %w", err)
		} else if updated {
			b.logger.Info("updated backrun reward", "bundle", record.BundleHashes, "amount", amount.String())
			b.metrics.rewards.Inc()
			b.metrics.rewardsTotal.Add(float64(amount.Int64()))
			account, txnHash, err := b.store.GetSwapRewardee(ctx, record.BundleHashes)
			if err != nil {
				b.logger.Error("getting backrun rewardee", "bundle", record.BundleHashes, "error", err)
				continue
			}
			if err := b.points.AssignPoints(ctx, account, txnHash, amount); err != nil {
				b.logger.Error("assigning backrun points", "user", account.Hex(), "tx", txnHash.Hex(), "error", err)
				continue
			}
		}
	}

	return nil
}

func (b *backrunner) Backrun(
	ctx context.Context,
	rawTx string,
	commitments []*bidderapiv1.Commitment,
) error {
	body, err := newReq(1, rawTx, commitments)
	if err != nil {
		return fmt.Errorf("creating backrun request: %w", err)
	}

	txHash := common.HexToHash(commitments[0].TxHashes[0])

	if err := b.store.AddSwapInfo(ctx, txHash, commitments[0].BlockNumber, body.Builders()); err != nil {
		return fmt.Errorf("storing backrun info: %w", err)
	}

	select {
	case b.reqChan <- body:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *backrunner) doBackrun(ctx context.Context, req backrunRequest) error {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return fmt.Errorf("encoding backrun request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, b.rpcURL, buf)
	if err != nil {
		return fmt.Errorf("creating backrun HTTP request: %w", err)
	}

	b.metrics.attempts.Inc()
	start := time.Now()

	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := b.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("sending backrun HTTP request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		b.metrics.fail.Inc()
		b.metrics.latency.Observe(float64(time.Since(start).Milliseconds()))
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status %d: %s", resp.StatusCode, string(respBody))
	}

	b.metrics.success.Inc()
	b.metrics.latency.Observe(float64(time.Since(start).Milliseconds()))
	b.logger.Info("backrun sent", "request", req.String())

	return nil
}
