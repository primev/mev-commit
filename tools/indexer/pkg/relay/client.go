package relay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/primev/mev-commit/indexer/pkg/config"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/primev/mev-commit/indexer/pkg/database"
	httputil "github.com/primev/mev-commit/indexer/pkg/http"

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

		// Use go-ethereum's hexutil for hex parsing
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

// Insert bid rows (relays are only for bids)
func InsertBid(ctx context.Context, db *database.DB, slot int64, relayID int64, bid map[string]any) error {

	if slot <= 0 || relayID <= 0 {
		return fmt.Errorf("invalid slot or relayID")
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
		return nil // skip if no value
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

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Write
	// Change PostgreSQL $1,$2,$3... to MySQL ?,?,?...
	_, err := db.Conn.ExecContext(ctx2, `
    INSERT INTO bids(
        slot, relay_id, builder_pubkey, proposer_pubkey,
        proposer_fee_recipient, value_wei, block_number, timestamp_ms
    )
    VALUES (?,?,?,?,?,?,?,?)`,
		slot, relayID, builder, proposer, feeRec, valStr, blockNum, tsMS,
	)
	return err
}

func UpsertRelaysAndLoad(ctx context.Context, db *database.DB) ([]Row, error) {
	// upsert defaults from code
	if err := db.UpsertRelays(ctx, config.RelaysDefault); err != nil {
		return nil, err
	}
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// load active
	rows, err := db.Conn.QueryContext(ctx2, `SELECT relay_id, base_url FROM relays WHERE is_active = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rws []Row
	for rows.Next() {
		var id int64
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			continue // Skip bad rows
		}
		rws = append(rws, Row{ID: id, URL: url})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rws, nil
}

func FetchBuilderBlocksReceived(httpc *http.Client, relayBase string, slot int64) ([]map[string]any, error) {
	url := fmt.Sprintf("%s/relay/v1/data/bidtraces/builder_blocks_received?slot=%d", strings.TrimRight(relayBase, "/"), slot)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var arr []map[string]any
	if err := httputil.FetchJSONWithRetry(ctx, httpc, url, &arr, 2, 200*time.Millisecond); err != nil {
		return nil, err
	}

	return arr, nil
}
