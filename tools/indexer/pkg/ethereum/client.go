package ethereum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"

	"github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
)

func CallAreOptedInAtBlock(httpc *http.Client, cfg *config.Config, blockNum int64, pubkey []byte) (bool, error) {
	if len(pubkey) == 0 {
		return false, fmt.Errorf("empty pubkey")
	}
	client, err := ethclient.Dial(cfg.AlchemyRPC)
	if err != nil {
		return false, err
	}
	contract, err := validatoroptinrouter.NewValidatoroptinrouter(common.HexToAddress(cfg.OptInContract), client)
	if err != nil {
		return false, err
	}

	result, err := contract.AreValidatorsOptedIn(&bind.CallOpts{BlockNumber: big.NewInt(blockNum)}, [][]byte{pubkey})
	if err != nil {
		return false, err
	}

	if len(result) == 0 {
		return false, nil
	}
	o := result[0]
	return o.IsVanillaOptedIn || o.IsAvsOptedIn || o.IsMiddlewareOptedIn, nil
}

// GetLatestBlockNumber gets the latest block number from Ethereum RPC
func GetLatestBlockNumber(httpc *http.Client, rpcURL string) (int64, error) {
	payload := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_blockNumber",
		"params":  []any{},
	}

	buf, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", rpcURL, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpc.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Convert hex to int64
	blockNum, err := strconv.ParseInt(result.Result[2:], 16, 64)
	return blockNum, err
}
