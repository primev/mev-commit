package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"

	"strconv"
)

type Row struct {
	ID  int64
	URL string
}

func parseBigString(v any) (string, bool) {
	switch t := v.(type) {
	case nil:
		return "", false
	case string:
		z := strings.ReplaceAll(strings.TrimSpace(t), ",", "")
		if z == "" {
			return "", false
		}

		if strings.HasPrefix(z, "0x") || strings.HasPrefix(z, "0X") {
			bi, err := hexutil.DecodeBig(z)
			if err != nil {
				return "", false
			}
			return bi.String(), true
		}

		// For decimal strings
		if _, ok := new(big.Int).SetString(z, 10); ok {
			return z, true
		}
		return "", false
	case float64:
		return strconv.FormatFloat(t, 'f', 0, 64), true
	case json.Number:
		return t.String(), true
	default:
		return fmt.Sprintf("%v", t), true
	}
}

func BuildBidInsert(slot int64, relayID int64, bid map[string]any) (database.BidRow, bool) {

	if slot <= 0 || relayID <= 0 {
		return database.BidRow{}, false
	}

	// helper to read alternative keys from different relay schemas
	get := func(keys ...string) any {
		for _, k := range keys {
			if v, ok := bid[k]; ok {
				return v
			}
		}
		return nil
	}

	// Parse fields
	builder := common.FromHex(fmt.Sprint(get("builder_pubkey", "builderPubkey", "builder")))
	proposer := common.FromHex(fmt.Sprint(get("proposer_pubkey", "proposerPubkey")))
	feeRec := common.FromHex(fmt.Sprint(get("proposer_fee_recipient", "proposerFeeRecipient", "feeRecipient")))

	valStr, ok := parseBigString(get("value", "value_wei", "valueWei"))
	if !ok || valStr == "" {
		return database.BidRow{}, false // skip if no value
	}

	var blockNum *int64
	if v := get("block_number", "blockNumber"); v != nil {
		switch t := v.(type) {
		case float64:
			x := int64(t)
			blockNum = &x
		case string:
			if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
				if bi, ok := new(big.Int).SetString(t[2:], 16); ok {
					x := bi.Int64()
					blockNum = &x
				}
			} else if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				blockNum = &n
			}
		}
	}

	var tsMS *int64
	if v := get("timestamp_ms", "timestampMs", "time_ms", "time"); v != nil {
		switch t := v.(type) {
		case float64:
			x := int64(t)
			tsMS = &x
		case string:
			if n, err := strconv.ParseInt(t, 10, 64); err == nil {
				tsMS = &n
			}
		}
	}

	return database.BidRow{
		Slot:     slot,
		RelayID:  relayID,
		Builder:  hexutil.Encode(builder),
		Proposer: hexutil.Encode(proposer),
		FeeRec:   hexutil.Encode(feeRec),
		ValStr:   valStr,
		BlockNum: blockNum,
		TsMS:     tsMS,
	}, true
}

func FetchBuilderBlocksReceived(ctx context.Context, httpc *retryablehttp.Client, relayBase string, slot int64) ([]map[string]any, error) {
	url := fmt.Sprintf("%s/relay/v1/data/bidtraces/builder_blocks_received?slot=%d", strings.TrimRight(relayBase, "/"), slot)

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var arr []map[string]any
	if err := httputil.FetchJSON(ctx2, httpc, url, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}
