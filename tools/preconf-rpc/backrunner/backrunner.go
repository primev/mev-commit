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
	"net/http"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
)

var builders = map[common.Address]string{
	common.HexToAddress("0xE3d71EF44D20917b93AA93e12Bd35b0859824A8F"): "btcs",
	common.HexToAddress("0x2445e5e28890De3e93F39fCA817639c470F4d3b9"): "iobuilder",
	common.HexToAddress("0xB3998135372F1eE16Cb510af70ed212b5155Af62"): "titan",
	common.HexToAddress("0x570e531fB805B5eEbD5F29Eaa2766fBeB4977ddE"): "quasar",
}

type Store interface {
	AddBackrunInfo(ctx context.Context, txnHash common.Hash, blockNumber int64) error
	RewardsToCheck(ctx context.Context, uptoBlockNumber int64) ([]common.Hash, uint64, error)
	UpdateSwapReward(ctx context.Context, bundleHash common.Hash, reward *big.Int) error
}

type BlockNumberGetter interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type backrunner struct {
	client    *http.Client
	rpcURL    string
	apiURL    string
	apiKey    string
	store     Store
	bNoGetter BlockNumberGetter
	logger    *slog.Logger
}

func New(client *http.Client, apiKey, apiURL, rpcURL string) (*backrunner, error) {
	urlParsed, err := url.Parse(rpcURL)
	if err != nil {
		return nil, err
	}
	q := urlParsed.Query()
	q.Add("apiKey", apiKey)
	q.Add("mode", "primev")
	urlParsed.RawQuery = q.Encode()

	apiParsed, err := url.Parse(fmt.Sprintf("%s/api/transactions", apiURL))
	if err != nil {
		return nil, err
	}

	return &backrunner{
		client: client,
		rpcURL: urlParsed.String(),
		apiURL: apiParsed.String(),
		apiKey: apiKey,
	}, nil
}

type backrunRequest struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int    `json:"id"`
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
		Params: map[string]any{
			"txs":             []string{rawTx},
			"blockNumber":     fmt.Sprintf("0x%x", blkNo),
			"trustedBuilders": buildersSelected,
		},
		ID: id,
	}, nil
}

func (b *backrunner) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		ticker := time.NewTicker(time.Second * 15)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				currentBlock, err := b.bNoGetter.BlockNumber(ctx)
				if err != nil {
					b.logger.Error("getting current block number", "error", err)
					continue
				}
				txns, start, err := b.store.RewardsToCheck(ctx, int64(currentBlock))
				if err != nil {
					continue
				}
				if len(txns) == 0 {
					continue
				}
				if err := b.checkRewards(ctx, txns, start); err != nil {
					b.logger.Error("checking backrun rewards", "error", err)
					continue
				}
			}
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

func (b *backrunner) checkRewards(ctx context.Context, txns []common.Hash, start uint64) error {
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

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return fmt.Errorf("encoding backrun request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.rpcURL, buf)
	if err != nil {
		return fmt.Errorf("creating backrun HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
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

	respStruct := map[string]any{}
	if err := json.Unmarshal(respBody, &respStruct); err != nil {
		return fmt.Errorf("unmarshaling backrun HTTP response: %w", err)
	}

	result, found := respStruct["result"]
	if !found {
		return fmt.Errorf("no result in backrun response: %s", string(respBody))
	}

	resultMap, ok := result.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid result format in backrun response: %s", string(respBody))
	}

	bundleHashStr, ok := resultMap["bundle_hash"].(string)
	if !ok {
		return fmt.Errorf("invalid bundle_hash format in backrun response: %s", string(respBody))
	}

	bundleHash := common.HexToHash(bundleHashStr)
	b.logger.Info("backrun submitted", "txHash", txHash.Hex(), "bundleHash", bundleHash.Hex())

	return b.store.AddBackrunInfo(ctx, txHash, commitments[0].BlockNumber)
}
