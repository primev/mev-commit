package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
)

type DB struct {
	conn *sql.DB
}
type BidInsert struct {
	Slot        int64
	RelayID     int64
	BuilderHex  string
	ProposerHex string
	FeeRecHex   string
	ValStr      string
	BlockNum    *int64
	TsMS        *int64
}

func Connect(ctx context.Context, dsn string, maxConns, minConns int) (*DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open StarRocks connection: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(maxConns)
	conn.SetMaxIdleConns(minConns)
	conn.SetConnMaxLifetime(time.Hour)
	conn.SetConnMaxIdleTime(30 * time.Minute)
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := conn.PingContext(pingCtx); err != nil {
		_ = conn.Close()

		return nil, fmt.Errorf("StarRocks ping failed: %v", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetMaxSlotNumber(ctx context.Context) (int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var slot int64
	err := db.conn.QueryRowContext(ctx2, `SELECT COALESCE(MAX(slot),0) FROM blocks`).Scan(&slot)
	return slot, err
}

func (db *DB) GetMinSlotNumber(ctx context.Context) (int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var slot int64
	err := db.conn.QueryRowContext(ctx2, `SELECT COALESCE(MIN(slot), 0) FROM blocks`).Scan(&slot)
	return slot, err
}

func (db *DB) EnsureStateTable(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ddl := `
	CREATE TABLE IF NOT EXISTS ingestor_state (
		id TINYINT,
		last_block_number BIGINT
	) ENGINE=OLAP
	DUPLICATE KEY(id)
	DISTRIBUTED BY HASH(id) BUCKETS 1
	PROPERTIES (
		"replication_num" = "1"
	)`

	if _, err := db.conn.ExecContext(ctx2, ddl); err != nil {
		return fmt.Errorf("failed to create state table: %w", err)
	}

	var count int
	err := db.conn.QueryRowContext(ctx2, `SELECT COUNT(*) FROM ingestor_state WHERE id = 1`).Scan(&count)
	if err != nil || count == 0 {
		_, err = db.conn.ExecContext(ctx2,
			`INSERT INTO ingestor_state (id, last_block_number) VALUES (1, 0)`)
		if err != nil {
			return fmt.Errorf("failed to insert initial state: %w", err)
		}
	}

	return nil
}

func (db *DB) GetMaxBlockNumber(ctx context.Context) (int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bn int64
	err := db.conn.QueryRowContext(ctx2, `SELECT COALESCE(MAX(block_number),0) FROM blocks`).Scan(&bn)
	return bn, err
}

func (db *DB) LoadLastBlockNumber(ctx context.Context) (int64, bool) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bn int64
	err := db.conn.QueryRowContext(ctx2,
		`SELECT last_block_number FROM ingestor_state WHERE id = 1 LIMIT 1`).Scan(&bn)
	if err != nil {
		return 0, false
	}
	return bn, true
}

func (db *DB) SaveLastBlockNumber(ctx context.Context, bn int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	q2 := fmt.Sprintf(`INSERT INTO ingestor_state (id, last_block_number) VALUES (1, %d)`, bn)
	if _, err := db.conn.ExecContext(ctx2, q2); err != nil {
		return fmt.Errorf("save last_block_number failed (insert): %w", err)
	}

	return nil
}

func (db *DB) UpsertBlockFromExec(ctx context.Context, ei *beacon.ExecInfo) error {
	if ei == nil || ei.BlockNumber == 0 {
		return fmt.Errorf("upsert block: nil exec info or block_number=0")
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var timestamp, proposerIndex, relayTag, builderPubkeyPrefix, proposerFeeRecHex, mevRewardEth, feeRecHex, proposerRewardEth string

	if ei.Timestamp != nil {
		timestamp = fmt.Sprintf("'%s'", ei.Timestamp.Format("2006-01-02 15:04:05"))
	} else {
		timestamp = "NULL"
	}
	if ei.ProposerIdx != nil {
		proposerIndex = fmt.Sprintf("%d", *ei.ProposerIdx)
	} else {
		proposerIndex = "NULL"
	}
	if ei.RelayTag != nil {
		relayTag = fmt.Sprintf("'%s'", *ei.RelayTag)
	} else {
		relayTag = "''"
	}
	if ei.BuilderPublicKey != nil {
		builderPubkeyPrefix = fmt.Sprintf("'%s'", (*ei.BuilderPublicKey))
	} else {
		builderPubkeyPrefix = "''"
	}
	if ei.ProposerFeeRecHex != nil {
		proposerFeeRecHex = fmt.Sprintf("'%s'", (*ei.ProposerFeeRecHex))
	} else {
		proposerFeeRecHex = "''"
	}
	if ei.MevRewardEth != nil {
		mevRewardEth = fmt.Sprintf("%.6f", *ei.MevRewardEth)
	} else {
		mevRewardEth = "NULL"
	}
	if ei.FeeRecipient != nil {
		feeRecHex = fmt.Sprintf("'%s'", (*ei.FeeRecipient))
	} else {
		feeRecHex = "''"
	}
	if ei.ProposerRewardEth != nil {
		proposerRewardEth = fmt.Sprintf("%.6f", *ei.ProposerRewardEth)
	} else {
		proposerRewardEth = "''"
	}

	query := fmt.Sprintf(`
INSERT INTO blocks(
    slot, block_number, timestamp, proposer_index,
    winning_relay, winning_builder_pubkey, proposer_fee_recipient, mev_reward, proposer_reward_eth, fee_recipient
) VALUES (%d, %d, %s, %s, %s, %s, %s, %s, %s, %s)`,
		ei.Slot, ei.BlockNumber, timestamp, proposerIndex,
		relayTag, builderPubkeyPrefix, proposerFeeRecHex,
		mevRewardEth, proposerRewardEth, feeRecHex,
	)
	_, err := db.conn.ExecContext(ctx2, query)
	if err != nil {
		return fmt.Errorf("upsert block slot=%d: %w", ei.Slot, err)
	}
	return nil
}

func (db *DB) UpdateValidatorPubkey(ctx context.Context, slot int64, vpub []byte) error {
	if slot == 0 {
		return fmt.Errorf("update validator: slot=0")
	}
	if len(vpub) == 0 {
		return nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	vhex := hexutil.Encode(vpub)

	q := fmt.Sprintf("INSERT INTO blocks (slot, validator_pubkey) VALUES (%d, '%s')", slot, vhex)

	if _, err := db.conn.ExecContext(ctx2, q); err != nil {
		return fmt.Errorf("update validator slot=%d: %w", slot, err)
	}

	return nil
}

// Minimal batching: builds one multi-VALUES INSERT.

type BidRow struct {
	Slot, RelayID             int64
	Builder, Proposer, FeeRec string
	ValStr                    string
	BlockNum, TsMS            *int64
}

func (db *DB) InsertBidsBatch(ctx context.Context, rows []BidRow) error {
	if len(rows) == 0 {
		return nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var sb strings.Builder
	sb.WriteString(`
        INSERT INTO bids(
            slot, relay_id, builder_pubkey, proposer_pubkey,
            proposer_fee_recipient, value_wei, block_number, timestamp_ms
        ) VALUES `)

	for i, r := range rows {
		if i > 0 {
			sb.WriteString(",")
		}
		blockNumSQL := "NULL"
		if r.BlockNum != nil {
			blockNumSQL = fmt.Sprintf("%d", *r.BlockNum)
		}
		tsMSSQL := "NULL"
		if r.TsMS != nil {
			tsMSSQL = fmt.Sprintf("%d", *r.TsMS)
		}

		fmt.Fprintf(&sb, "(%d,%d,'%s','%s','%s','%s',%s,%s)",
			r.Slot, r.RelayID, r.Builder, r.Proposer, r.FeeRec, r.ValStr, blockNumSQL, tsMSSQL)
	}

	_, err := db.conn.ExecContext(ctx2, sb.String())
	return err
}

func (db *DB) GetRecentMissingBlocks(ctx context.Context, lookback int64, batch int) ([]struct {
	Slot        int64
	BlockNumber int64
}, error,
) {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if lookback < 0 || batch < 0 || batch > 10000 {
		return nil, fmt.Errorf("invalid parameters: lookback=%d, batch=%d", lookback, batch)
	}

	query := fmt.Sprintf(`
        WITH recent AS (
            SELECT COALESCE(MAX(slot), 0) AS s FROM blocks
        )
        SELECT slot, block_number
        FROM blocks, recent
        WHERE slot > recent.s - %d
          AND block_number IS NOT NULL
          AND (winning_relay IS NULL 
               OR winning_builder_pubkey IS NULL 
               OR fee_recipient IS NULL 
               OR producer_reward_eth IS NULL 
               OR timestamp IS NULL 
               OR proposer_index IS NULL)
        ORDER BY slot DESC
        LIMIT %d`, lookback, batch)

	rows, err := db.conn.QueryContext(ctx2, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []struct {
		Slot        int64
		BlockNumber int64
	}
	for rows.Next() {
		var slot, bn int64
		if err := rows.Scan(&slot, &bn); err != nil {
			continue
		}
		results = append(results, struct {
			Slot        int64
			BlockNumber int64
		}{Slot: slot, BlockNumber: bn})
	}
	return results, rows.Err()
}

func (db *DB) GetRecentSlotsWithBlocks(ctx context.Context, lookback int64, batch int) ([]int64, error) {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	q := fmt.Sprintf(`
WITH recent AS (SELECT COALESCE(MAX(slot),0) AS s FROM blocks)
SELECT DISTINCT slot
FROM blocks, recent
WHERE slot > recent.s - ?
  AND block_number IS NOT NULL
ORDER BY slot DESC
LIMIT %d`, batch)
	rows, err := db.conn.QueryContext(ctx2, q, lookback)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var slots []int64
	for rows.Next() {
		var slot int64
		if err := rows.Scan(&slot); err != nil {
			continue
		}
		slots = append(slots, slot)
	}
	return slots, rows.Err()
}

func (db *DB) GetValidatorsNeedingOptInCheck(ctx context.Context, lookback int64, batch int) ([]struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
}, error,
) {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	q := fmt.Sprintf(`
WITH recent AS (SELECT COALESCE(MAX(slot),0) AS s FROM blocks)
SELECT slot, block_number, validator_pubkey
FROM blocks, recent
WHERE slot > recent.s - ?
  AND block_number IS NOT NULL
  AND validator_pubkey IS NOT NULL
  AND validator_opted_in IS NULL
ORDER BY slot DESC
LIMIT %d`, batch)
	rows, err := db.conn.QueryContext(ctx2, q, lookback)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []struct {
		Slot            int64
		BlockNumber     int64
		ValidatorPubkey []byte
	}
	for rows.Next() {
		var slot, bn int64
		var vpk []byte
		if err := rows.Scan(&slot, &bn, &vpk); err != nil {
			continue
		}
		results = append(results, struct {
			Slot            int64
			BlockNumber     int64
			ValidatorPubkey []byte
		}{
			Slot: slot, BlockNumber: bn, ValidatorPubkey: vpk,
		})
	}
	return results, rows.Err()
}

func (db *DB) UpdateValidatorOptInStatus(ctx context.Context, slot int64, opted bool) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	v := 0
	if opted {
		v = 1
	}
	q := fmt.Sprintf(
		"UPDATE blocks SET validator_opted_in=%d WHERE slot=%d",
		v, slot,
	)
	_, err := db.conn.ExecContext(ctx2, q)
	return err
}

func (db *DB) GetValidatorPubkeyWithRetry(ctx context.Context, slot int64, retries int, retryDelay time.Duration) ([]byte, error) {
	var vpk []byte
	for i := 0; i < retries; i++ {
		ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
		err := db.conn.QueryRowContext(ctx2, `SELECT validator_pubkey FROM blocks WHERE slot=?`, slot).Scan(&vpk)
		cancel()

		if err == nil && len(vpk) > 0 {
			return vpk, nil
		}

		if i < retries-1 {
			time.Sleep(retryDelay)
		}
	}
	return nil, fmt.Errorf("validator pubkey not available after %d retries", retries)
}

func (db *DB) BatchUpsertBlocksFromExec(ctx context.Context, execInfos []*beacon.ExecInfo) error {
	if len(execInfos) == 0 {
		return nil
	}

	const maxRowsPerInsert = 500

	type row struct {
		sql string
	}

	rows := make([]row, 0, len(execInfos))

	for _, ei := range execInfos {
		if ei == nil || ei.BlockNumber == 0 {
			continue
		}

		var timestamp, proposerIndex, relayTag, builderPubkeyPrefix,
			proposerFeeRecHex, mevRewardEth, feeRecHex, proposerRewardEth string

		if ei.Timestamp != nil {
			timestamp = fmt.Sprintf("'%s'", ei.Timestamp.Format("2006-01-02 15:04:05"))
		} else {
			timestamp = "NULL"
		}
		if ei.ProposerIdx != nil {
			proposerIndex = fmt.Sprintf("%d", *ei.ProposerIdx)
		} else {
			proposerIndex = "NULL"
		}
		if ei.RelayTag != nil {
			relayTag = fmt.Sprintf("'%s'", *ei.RelayTag)
		} else {
			relayTag = "''"
		}
		if ei.BuilderPublicKey != nil {
			builderPubkeyPrefix = fmt.Sprintf("'%s'", *ei.BuilderPublicKey)
		} else {
			builderPubkeyPrefix = "''"
		}
		if ei.ProposerFeeRecHex != nil {
			proposerFeeRecHex = fmt.Sprintf("'%s'", *ei.ProposerFeeRecHex)
		} else {
			proposerFeeRecHex = "''"
		}
		if ei.MevRewardEth != nil {
			mevRewardEth = fmt.Sprintf("%.6f", *ei.MevRewardEth)
		} else {
			mevRewardEth = "''"
		}
		if ei.FeeRecipient != nil {
			feeRecHex = fmt.Sprintf("'%s'", *ei.FeeRecipient)
		} else {
			feeRecHex = "''"
		}
		if ei.ProposerRewardEth != nil {
			proposerRewardEth = fmt.Sprintf("%.6f", *ei.ProposerRewardEth)
		} else {
			proposerRewardEth = "''"
		}

		value := fmt.Sprintf("(%d, %d, %s, %s, %s, %s, %s, %s, %s, %s)",
			ei.Slot,
			ei.BlockNumber,
			timestamp,
			proposerIndex,
			relayTag,
			builderPubkeyPrefix,
			proposerFeeRecHex,
			mevRewardEth,
			proposerRewardEth,
			feeRecHex,
		)

		rows = append(rows, row{sql: value})
	}

	if len(rows) == 0 {
		return nil
	}

	// Insert in chunks
	for i := 0; i < len(rows); i += maxRowsPerInsert {
		j := i + maxRowsPerInsert
		if j > len(rows) {
			j = len(rows)
		}

		query := `
INSERT INTO blocks(
    slot, block_number, timestamp, proposer_index,
    winning_relay, winning_builder_pubkey, proposer_fee_recipient,
    mev_reward, proposer_reward_eth, fee_recipient
) VALUES ` + strings.Join(func(ss []row) []string {
			out := make([]string, len(ss))
			for k := range ss {
				out[k] = ss[k].sql
			}
			return out
		}(rows[i:j]), ",")

		// short timeout per batch
		ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := db.conn.ExecContext(ctx2, query)
		cancel()
		if err != nil {
			return fmt.Errorf("batch upsert blocks [%d:%d]: %w", i, j, err)
		}
	}

	return nil
}

// INSERT into StarRocks PK table acts as UPSERT on slot
func (db *DB) UpsertBlockPubkeysDirect(ctx context.Context, pairs []struct {
	Slot   int64
	Pubkey string
},
) error {
	if len(pairs) == 0 {
		return nil
	}
	const maxRows = 1000
	vals := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if p.Slot == 0 || p.Pubkey == "" {
			continue
		}
		vals = append(vals, fmt.Sprintf("(%d, '%s')", p.Slot, p.Pubkey))
	}
	if len(vals) == 0 {
		return nil
	}

	for i := 0; i < len(vals); i += maxRows {
		j := i + maxRows
		if j > len(vals) {
			j = len(vals)
		}
		q := "INSERT INTO blocks (slot, validator_pubkey) VALUES " + strings.Join(vals[i:j], ",")
		ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
		_, err := db.conn.ExecContext(ctx2, q)
		cancel()
		if err != nil {
			return fmt.Errorf("upsert block pubkeys [%d:%d]: %w", i, j, err)
		}
	}
	return nil
}
