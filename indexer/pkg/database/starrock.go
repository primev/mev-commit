package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/indexer/pkg/beacon"
	"github.com/primev/mev-commit/indexer/pkg/config"
)

type DB struct {
	Conn *sql.DB
}

func MustConnect(ctx context.Context, dsn string, maxConns, minConns int) (*DB, error) {
	if cfg, err := mysql.ParseDSN(dsn); err == nil {
		if cfg.Params == nil {
			cfg.Params = map[string]string{}
		}
		cfg.Params["interpolateParams"] = "true"
		dsn = cfg.FormatDSN()
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
		conn.Close()
		return nil, fmt.Errorf("StarRocks ping failed: %v", err)
	}

	return &DB{Conn: conn}, nil

}
func (db *DB) Close() {
	db.Conn.Close()
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

	if _, err := db.Conn.ExecContext(ctx2, ddl); err != nil {
		return fmt.Errorf("failed to create state table: %w", err)
	}

	var count int
	err := db.Conn.QueryRowContext(ctx2, `SELECT COUNT(*) FROM ingestor_state WHERE id = 1`).Scan(&count)
	if err != nil || count == 0 {
		_, err = db.Conn.ExecContext(ctx2,
			`INSERT INTO ingestor_state (id, last_block_number) VALUES (1, 0)`)
		if err != nil {
			return fmt.Errorf("failed to insert initial state: %w", err)
		}
	}

	return nil
}

func (db *DB) LoadLastBlockNumber(ctx context.Context) (int64, bool) {
	ctx2, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var bn int64
	err := db.Conn.QueryRowContext(ctx2,
		`SELECT last_block_number FROM ingestor_state WHERE id = 1 LIMIT 1`).Scan(&bn)
	if err != nil {
		return 0, false
	}
	return bn, true
}

func (db *DB) SaveLastBlockNumber(ctx context.Context, bn int64) error {
	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := db.Conn.ExecContext(ctx2, `DELETE FROM ingestor_state WHERE id = 1`)
	if err != nil {
		return fmt.Errorf("failed to delete old state: %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO ingestor_state (id, last_block_number) VALUES (1, %d)`, bn)
	_, err = db.Conn.ExecContext(ctx2, query)
	if err != nil {
		return fmt.Errorf("save last_block_number failed: %w", err)
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
		value := fmt.Sprintf("('%s', '%s', '%s', 1)", r.Name, r.Tag, r.URL)
		values = append(values, value)
	}

	query := fmt.Sprintf(`INSERT INTO relays (name, tag, base_url, is_active) VALUES %s`,
		strings.Join(values, ","))

	_, err := db.Conn.ExecContext(ctx2, query)
	return err
}

func (db *DB) UpsertBlockFromExec(ctx context.Context, ei *beacon.ExecInfo) error {
	if ei == nil || ei.BlockNumber == 0 {
		return fmt.Errorf("upsert block: nil exec info or block_number=0")
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var timestamp, proposerIndex, relayTag, rewardEth string

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

	if ei.RewardEth != nil {
		rewardEth = fmt.Sprintf("%.6f", *ei.RewardEth)
	} else {
		rewardEth = "NULL"
	}

	query := fmt.Sprintf(`
INSERT INTO blocks(
    slot, block_number, timestamp, proposer_index,
    winning_relay, producer_reward_eth
) VALUES (%d, %d, %s, %s, %s, %s)`,
		ei.Slot, ei.BlockNumber, timestamp, proposerIndex, relayTag, rewardEth)

	_, err := db.Conn.ExecContext(ctx2, query)
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
	pubkeyHex := fmt.Sprintf("%x", vpub)
	_, err := db.Conn.ExecContext(ctx2, `
		INSERT INTO blocks (slot, validator_pubkey) VALUES (?, ?)`,
		slot, pubkeyHex)
	if err != nil {
		return fmt.Errorf("update validator slot=%d: %w", slot, err)
	}
	return nil
}
