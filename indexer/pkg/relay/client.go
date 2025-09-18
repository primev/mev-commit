package relay

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/primev/mev-commit/indexer/pkg/config"

	"strconv"

	"github.com/primev/mev-commit/indexer/pkg/database"
	httputil "github.com/primev/mev-commit/indexer/pkg/http"
	"github.com/primev/mev-commit/indexer/pkg/utils"
)

type Row struct {
	ID  int64
	URL string
}

// Insert bid rows (relays are only for bids)
func InsertBid(ctx context.Context, db *database.DB, slot int64, relayID int64, bid map[string]any) error {
	const batchSize = 200
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
	builder := utils.HexToBytes(fmt.Sprint(get("builder_pubkey", "builderPubkey", "builder")))
	proposer := utils.HexToBytes(fmt.Sprint(get("proposer_pubkey", "proposerPubkey")))
	feeRec := utils.HexToBytes(fmt.Sprint(get("proposer_fee_recipient", "proposerFeeRecipient", "feeRecipient")))

	valStr, ok := utils.ParseBigString(get("value", "value_wei", "valueWei"))
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

	hb := fmt.Sprintf("%x", builder) // hex string, no "0x" prefix
	hp := fmt.Sprintf("%x", proposer)
	hf := fmt.Sprintf("%x", feeRec)

	_, err := db.Conn.ExecContext(ctx2, `
    INSERT INTO bids(
        slot, relay_id, builder_pubkey, proposer_pubkey,
        proposer_fee_recipient, value_wei, block_number, timestamp_ms
    )
    VALUES (?, ?, UNHEX(?), UNHEX(?), UNHEX(?), ?, ?, ?)`,
		slot, relayID, hb, hp, hf, valStr, blockNum, tsMS,
	)

	return err
}

// BatchInsertBids inserts bids using multi-row VALUES with UNHEX(?) for binary fields.
// Falls back to per-row InsertBid on batch error.
func BatchInsertBids(ctx context.Context, db *database.DB, slot, relayID int64, bids []map[string]any) (ok, fail int, lastErr error) {
	if slot <= 0 || relayID <= 0 || len(bids) == 0 {
		return 0, 0, nil
	}
	const batchSize = 200

	// local parser reuses InsertBid’s logic but returns the prepared args
	parse := func(bid map[string]any) (args []any, skip bool) {
		get := func(keys ...string) any {
			for _, k := range keys {
				if v, ok := bid[k]; ok {
					return v
				}
			}
			return nil
		}
		builder := utils.HexToBytes(fmt.Sprint(get("builder_pubkey", "builderPubkey", "builder")))
		proposer := utils.HexToBytes(fmt.Sprint(get("proposer_pubkey", "proposerPubkey")))
		feeRec := utils.HexToBytes(fmt.Sprint(get("proposer_fee_recipient", "proposerFeeRecipient", "feeRecipient")))

		valStr, ok := utils.ParseBigString(get("value", "value_wei", "valueWei"))
		if !ok || valStr == "" {
			return nil, true
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
		hb := fmt.Sprintf("%x", builder) // hex (no 0x); UNHEX(?) will restore bytes
		hp := fmt.Sprintf("%x", proposer)
		hf := fmt.Sprintf("%x", feeRec)
		return []any{slot, relayID, hb, hp, hf, valStr, blockNum, tsMS}, false
	}

	makeSQL := func(n int) string {
		ph := "(?, ?, UNHEX(?), UNHEX(?), UNHEX(?), ?, ?, ?)"
		parts := make([]string, n)
		for i := range parts {
			parts[i] = ph
		}
		return "INSERT INTO bids(" +
			"slot, relay_id, builder_pubkey, proposer_pubkey, " +
			"proposer_fee_recipient, value_wei, block_number, timestamp_ms" +
			") VALUES " + strings.Join(parts, ",")
	}

	for i := 0; i < len(bids); i += batchSize {
		j := i + batchSize
		if j > len(bids) {
			j = len(bids)
		}

		args := make([]any, 0, (j-i)*8)
		rows := 0
		for _, b := range bids[i:j] {
			a, skip := parse(b)
			if skip {
				continue
			}
			args = append(args, a...)
			rows++
		}
		if rows == 0 {
			continue
		}

		batchCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
		_, err := db.Conn.ExecContext(batchCtx, makeSQL(rows), args...)
		cancel()
		if err != nil {
			lastErr = err
			// fallback: per-row using your existing InsertBid
			for _, b := range bids[i:j] {
				oneCtx, c := context.WithTimeout(ctx, 2*time.Second)
				if e := InsertBid(oneCtx, db, slot, relayID, b); e != nil {
					fail++
					lastErr = e
				} else {
					ok++
				}
				c()
			}
		} else {
			ok += rows
		}
	}
	return ok, fail, lastErr
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
