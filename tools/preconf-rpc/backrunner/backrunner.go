package backrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	AddSwapInfo(ctx context.Context, bundleHash common.Hash, txnHash common.Hash) error
	RewardsToCheck(ctx context.Context) (map[common.Hash]common.Hash, error)
	UpdateSwapReward(ctx context.Context, bundleHash common.Hash, reward *big.Int) error
}

type backrunner struct {
	client *http.Client
	apiURL string
	store  Store
}

func New(client *http.Client, apiKey, apiURL string) (*backrunner, error) {
	urlParsed, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	q := urlParsed.Query()
	q.Add("apiKey", apiKey)
	q.Add("mode", "primev")
	urlParsed.RawQuery = q.Encode()

	return &backrunner{
		client: client,
		apiURL: urlParsed.String(),
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
				bundles, err := b.store.RewardsToCheck(ctx)
				if err != nil {
					continue
				}
			}
		}
	}()

	return done
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.apiURL, buf)
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
	return b.store.AddSwapInfo(ctx, bundleHash, txHash)
}
