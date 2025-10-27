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
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
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
	// Parse DSN to extract database name
	dbName, dsnWithoutDB := parseDSN(dsn)

	// If database name is found, try to create it if it doesn't exist
	if dbName != "" {
		if err := ensureDatabase(ctx, dsnWithoutDB, dbName); err != nil {
			return nil, fmt.Errorf("failed to ensure database exists: %w", err)
		}
	}

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

// parseDSN extracts the database name from DSN and returns DSN without the database
// DSN format: username:password@tcp(host:port)/database?params
func parseDSN(dsn string) (string, string) {
	// Find the database name between "/" and "?" or end of string
	slashIdx := strings.LastIndex(dsn, "/")
	if slashIdx == -1 {
		return "", dsn
	}

	afterSlash := dsn[slashIdx+1:]
	questionIdx := strings.Index(afterSlash, "?")

	var dbName string
	if questionIdx == -1 {
		dbName = afterSlash
	} else {
		dbName = afterSlash[:questionIdx]
	}

	// Return DSN without database name
	dsnWithoutDB := dsn[:slashIdx+1]
	if questionIdx != -1 {
		dsnWithoutDB += afterSlash[questionIdx:]
	}

	return dbName, dsnWithoutDB
}

// ensureDatabase creates the database if it doesn't exist
func ensureDatabase(ctx context.Context, dsnWithoutDB, dbName string) error {
	conn, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return fmt.Errorf("failed to open connection without database: %w", err)
	}
	defer conn.Close()

	createCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Try to create database if it doesn't exist
	_, err = conn.ExecContext(createCtx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	return nil
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
func (db *DB) EnsureStateTable(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ddl := `
	CREATE TABLE IF NOT EXISTS ingestor_state (
		id TINYINT,
		last_forward_block BIGINT,
		last_backward_block BIGINT
	) ENGINE=OLAP
	PRIMARY KEY(id)
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
			`INSERT INTO ingestor_state (id, last_forward_block, last_backward_block) VALUES (1, 0, 0)`)
		if err != nil {
			return fmt.Errorf("failed to insert initial state: %w", err)
		}
	}

	return nil
}

func (db *DB) EnsureRelaysTable(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ddl := `
	CREATE TABLE IF NOT EXISTS relays (
		relay_id INT,
		name VARCHAR(255),
		tag VARCHAR(255),
		base_url VARCHAR(255),
		is_active TINYINT
	) ENGINE=OLAP
	PRIMARY KEY(relay_id)
	DISTRIBUTED BY HASH(relay_id) BUCKETS 1
	PROPERTIES (
		"replication_num" = "1"
	)`

	if _, err := db.conn.ExecContext(ctx2, ddl); err != nil {
		return fmt.Errorf("failed to create relays table: %w", err)
	}

	return nil
}

func (db *DB) EnsureBlocksTable(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ddl := `
	CREATE TABLE IF NOT EXISTS blocks (
		slot BIGINT,
		block_number BIGINT,
		timestamp DATETIME,
		proposer_index BIGINT,
		winning_relay VARCHAR(255),
		producer_reward_eth DOUBLE,
		winning_builder_pubkey VARCHAR(255),
		fee_recipient VARCHAR(255),
		validator_pubkey VARCHAR(255),
		validator_opted_in TINYINT
	) ENGINE=OLAP
	PRIMARY KEY(slot)
	DISTRIBUTED BY HASH(slot) BUCKETS 10
	PROPERTIES (
		"replication_num" = "1"
	)`

	if _, err := db.conn.ExecContext(ctx2, ddl); err != nil {
		return fmt.Errorf("failed to create blocks table: %w", err)
	}

	return nil
}

func (db *DB) EnsureBidsTable(ctx context.Context) error {
	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ddl := `
	CREATE TABLE IF NOT EXISTS bids (
		slot BIGINT,
		relay_id BIGINT,
		builder_pubkey VARCHAR(255),
		proposer_pubkey VARCHAR(255),
		proposer_fee_recipient VARCHAR(255),
		value_wei VARCHAR(255),
		block_number BIGINT,
		timestamp_ms BIGINT
	) ENGINE=OLAP
	PRIMARY KEY(slot, relay_id, builder_pubkey)
	DISTRIBUTED BY HASH(slot) BUCKETS 10
	PROPERTIES (
		"replication_num" = "1"
	)`

	if _, err := db.conn.ExecContext(ctx2, ddl); err != nil {
		return fmt.Errorf("failed to create bids table: %w", err)
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
func (db *DB) GetValidatorPubkey(ctx context.Context, slot int64) ([]byte, error) {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var vpk []byte
	err := db.conn.QueryRowContext(ctx2, `SELECT validator_pubkey FROM blocks WHERE slot=?`, slot).Scan(&vpk)
	return vpk, err
}

// LoadForwardCheckpoint returns the last forward indexed block number
func (db *DB) LoadForwardCheckpoint(ctx context.Context) (int64, bool) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bn int64
	err := db.conn.QueryRowContext(ctx2,
		`SELECT last_forward_block FROM ingestor_state WHERE id = 1 LIMIT 1`).Scan(&bn)
	if err != nil {
		return 0, false
	}
	return bn, true
}

// LoadBackwardCheckpoint returns the last backward indexed block number
func (db *DB) LoadBackwardCheckpoint(ctx context.Context) (int64, bool) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bn int64
	err := db.conn.QueryRowContext(ctx2,
		`SELECT last_backward_block FROM ingestor_state WHERE id = 1 LIMIT 1`).Scan(&bn)
	if err != nil {
		return 0, false
	}
	return bn, true
}

// SaveForwardCheckpoint saves the forward indexer progress
func (db *DB) SaveForwardCheckpoint(ctx context.Context, bn int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	q := fmt.Sprintf(`INSERT INTO ingestor_state (id, last_forward_block) VALUES (1, %d)`, bn)
	if _, err := db.conn.ExecContext(ctx2, q); err != nil {
		return fmt.Errorf("save forward checkpoint failed: %w", err)
	}
	return nil
}

// SaveBackwardCheckpoint saves the backward indexer progress
func (db *DB) SaveBackwardCheckpoint(ctx context.Context, bn int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	q := fmt.Sprintf(`INSERT INTO ingestor_state (id, last_backward_block) VALUES (1, %d)`, bn)
	if _, err := db.conn.ExecContext(ctx2, q); err != nil {
		return fmt.Errorf("save backward checkpoint failed: %w", err)
	}
	return nil
}

func (db *DB) UpsertRelays(ctx context.Context, relays []config.Relay) error {
	if len(relays) == 0 {
		return nil
	}

	ctx2, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// StarRocks batch insert approach
	var values []string
	for _, r := range relays {
		value := fmt.Sprintf("(%d, '%s', '%s', '%s', 1)", r.Relay_id, r.Name, r.Tag, r.URL)
		values = append(values, value)
	}

	query := fmt.Sprintf(`INSERT INTO relays (relay_id, name, tag, base_url, is_active) VALUES %s`,
		strings.Join(values, ","))

	_, err := db.conn.ExecContext(ctx2, query)
	return err
}

func (db *DB) UpsertBlockFromExec(ctx context.Context, ei *beacon.ExecInfo) error {
	if ei == nil || ei.BlockNumber == 0 {
		return fmt.Errorf("upsert block: nil exec info or block_number=0")
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var timestamp, proposerIndex, relayTag, rewardEth, builderPubkeyPrefix, feeRecHex string

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
		relayTag = "NULL"
	}
	if ei.BuilderHex != nil {
		builderPubkeyPrefix = fmt.Sprintf("'%s'", (*ei.BuilderHex))
	} else {
		builderPubkeyPrefix = "NULL"
	}
	if ei.FeeRecHex != nil {
		feeRecHex = fmt.Sprintf("'%s'", (*ei.FeeRecHex))
	} else {
		feeRecHex = "NULL"
	}
	if ei.RewardEth != nil {
		rewardEth = fmt.Sprintf("%.6f", *ei.RewardEth)
	} else {
		rewardEth = "NULL"
	}

	query := fmt.Sprintf(`
INSERT INTO blocks(
    slot, block_number, timestamp, proposer_index,
    winning_relay, producer_reward_eth, winning_builder_pubkey, fee_recipient
) VALUES (%d, %d, %s, %s, %s, %s, %s, %s)`,
		ei.Slot, ei.BlockNumber, timestamp, proposerIndex, relayTag, rewardEth, builderPubkeyPrefix, feeRecHex)

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

func (db *DB) GetActiveRelays(ctx context.Context) ([]struct {
	ID  int64
	URL string
}, error) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx2, `SELECT relay_id, base_url FROM relays WHERE is_active = 1`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []struct {
		ID  int64
		URL string
	}
	for rows.Next() {
		var id int64
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			continue // Skip bad rows
		}
		results = append(results, struct {
			ID  int64
			URL string
		}{ID: id, URL: url})
	}
	return results, rows.Err()
}

func (db *DB) GetRecentMissingBlocks(ctx context.Context, lookback int64, batch int) ([]struct {
	Slot        int64
	BlockNumber int64
}, error) {
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
}, error) {
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
	} // TINYINT(1) in StarRocks
	q := fmt.Sprintf("INSERT INTO blocks (slot, validator_opted_in) VALUES (%d, %d)", slot, v)
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
