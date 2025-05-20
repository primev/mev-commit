package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/lib/pq"
)

// RelayRecord represents a record in the relay_data table
type RelayRecord struct {
	ID                 int64
	Slot               uint64
	BlockNumber        uint64
	ValidatorIndex     uint64
	ValidatorPubkey    string
	MEVReward          *big.Int
	MEVRewardRecipient string
	RelaysWithData     []string
	Winner             string
	TotalCommitments   int
	TotalRewards       int
	TotalSlashes       int
	TotalAmount        string
	CreatedAt          time.Time
}

type CommitmentRecord struct {
	ID                  int64
	BlockNumber         uint64
	CommitmentIndex     []byte
	Bidder              string
	Committer           string
	BidAmount           *big.Int
	SlashAmount         *big.Int
	DecayStartTimestamp uint64
	DecayEndTimestamp   uint64
	TxnHash             string
	RevertingTxHashes   string
	CommitmentDigest    []byte
	DispatchTimestamp   uint64
	CreatedAt           time.Time
}

// PostgresDB handles database connections and operations
type PostgresDB struct {
	db     *sql.DB
	logger *slog.Logger
}

// Config contains database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new database connection
func NewPostgresDB(
	cfg Config,
	logger *slog.Logger,
) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		//nolint:errcheck
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info(
		"Connected to PostgreSQL database",
		"host", cfg.Host,
		"database", cfg.DBName,
	)

	return &PostgresDB{
		db:     db,
		logger: logger,
	}, nil
}

// InitSchema initializes the database schema
func (p *PostgresDB) InitSchema(ctx context.Context) error {
	schema := `
	CREATE TABLE IF NOT EXISTS relay_data (
		id SERIAL PRIMARY KEY,
		slot BIGINT NOT NULL,
		block_number BIGINT NOT NULL,
		validator_index BIGINT NOT NULL,
		validator_pubkey TEXT NOT NULL,
		mev_reward NUMERIC NOT NULL,
		mev_reward_recipient TEXT,
		relays_with_data TEXT[] NOT NULL,
		winner TEXT,
		total_commitments INTEGER,
		total_rewards INTEGER,
		total_slashes INTEGER,
		total_amount TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_relay_data_slot ON relay_data(slot);
	CREATE INDEX IF NOT EXISTS idx_relay_data_block_number ON relay_data(block_number);
	CREATE INDEX IF NOT EXISTS idx_relay_data_validator_pubkey ON relay_data(validator_pubkey);

	CREATE TABLE IF NOT EXISTS block_commitments (
		id SERIAL PRIMARY KEY,
		block_number BIGINT NOT NULL,
		commitment_index BYTEA NOT NULL,
		bidder TEXT NOT NULL,
		committer TEXT NOT NULL,
		bid_amount NUMERIC NOT NULL,
		slash_amount NUMERIC NOT NULL,
		decay_start_timestamp BIGINT NOT NULL,
		decay_end_timestamp BIGINT NOT NULL,
		txn_hash TEXT,
		reverting_tx_hashes TEXT,
		commitment_digest BYTEA NOT NULL,
		dispatch_timestamp BIGINT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	
	CREATE INDEX IF NOT EXISTS idx_block_commitments_block ON block_commitments(block_number);
	CREATE INDEX IF NOT EXISTS idx_block_commitments_bidder ON block_commitments(bidder);
	CREATE INDEX IF NOT EXISTS idx_block_commitments_committer ON block_commitments(committer);
	`

	_, err := p.db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	p.logger.Info("Database schema initialized")
	return nil
}

// SaveRelayData saves relay data to the database
func (p *PostgresDB) SaveRelayData(
	ctx context.Context,
	data *RelayRecord,
) error {
	query := `
	INSERT INTO relay_data (
		slot, block_number, validator_index, validator_pubkey, mev_reward,
		mev_reward_recipient, relays_with_data, winner, total_commitments, total_rewards, 
		total_slashes, total_amount
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, created_at
	`

	// Convert big.Int to string for database storage
	mevRewardStr := data.MEVReward.String()

	// Convert the string array to a pq.StringArray for proper PostgreSQL compatibility
	relaysArray := pq.StringArray(data.RelaysWithData)

	row := p.db.QueryRowContext(
		ctx, query,
		data.Slot, data.BlockNumber, data.ValidatorIndex, data.ValidatorPubkey,
		mevRewardStr, data.MEVRewardRecipient,
		relaysArray, data.Winner, data.TotalCommitments, data.TotalRewards,
		data.TotalSlashes, data.TotalAmount,
	)

	err := row.Scan(&data.ID, &data.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save relay data: %w", err)
	}

	p.logger.Debug(
		"Saved relay data to database",
		slog.Int64("id", data.ID),
		slog.Uint64("block_number", data.BlockNumber),
		slog.Uint64("slot", data.Slot),
		slog.String("validator", data.ValidatorPubkey),
		slog.Time("created_at", data.CreatedAt),
	)

	return nil
}

// GetRelayDataByBlock retrieves relay data for a specific block
func (p *PostgresDB) GetRelayDataByBlock(
	ctx context.Context,
	blockNumber uint64,
) ([]*RelayRecord, error) {
	query := `
	SELECT id, slot, block_number, validator_index, validator_pubkey, 
		   mev_reward, mev_reward_recipient, relays_with_data, winner, total_commitments, 
		   total_rewards, total_slashes, total_amount, created_at
	FROM relay_data
	WHERE block_number = $1
	ORDER BY created_at DESC
	`

	rows, err := p.db.QueryContext(ctx, query, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query relay data: %w", err)
	}
	//nolint:errcheck
	defer rows.Close()

	var records []*RelayRecord

	for rows.Next() {
		var r RelayRecord
		var mevRewardStr string

		err := rows.Scan(
			&r.ID, &r.Slot, &r.BlockNumber, &r.ValidatorIndex, &r.ValidatorPubkey,
			&mevRewardStr, &r.MEVRewardRecipient,
			&r.RelaysWithData, &r.Winner, &r.TotalCommitments,
			&r.TotalRewards, &r.TotalSlashes, &r.TotalAmount, &r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		r.MEVReward = new(big.Int)
		r.MEVReward.SetString(mevRewardStr, 10)

		records = append(records, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// SaveBlockCommitments saves block commitments to the database
func (p *PostgresDB) SaveBlockCommitments(
	ctx context.Context,
	commitments []*CommitmentRecord,
) error {
	// Begin transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Prepare statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO block_commitments (
			block_number, commitment_index, bidder, committer, bid_amount,
			slash_amount, decay_start_timestamp, decay_end_timestamp, txn_hash,
			reverting_tx_hashes, commitment_digest, dispatch_timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT DO NOTHING
		RETURNING id, created_at
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	//nolint:errcheck
	defer stmt.Close()

	for _, commitment := range commitments {
		var id int64
		var createdAt time.Time

		err := stmt.QueryRowContext(
			ctx,
			commitment.BlockNumber,
			commitment.CommitmentIndex,
			commitment.Bidder,
			commitment.Committer,
			commitment.BidAmount.String(),
			commitment.SlashAmount.String(),
			commitment.DecayStartTimestamp,
			commitment.DecayEndTimestamp,
			commitment.TxnHash,
			commitment.RevertingTxHashes,
			commitment.CommitmentDigest,
			commitment.DispatchTimestamp,
		).Scan(&id, &createdAt)

		if err == sql.ErrNoRows {
			// This means it was a duplicate - continue processing
			continue
		}

		if err != nil {
			return fmt.Errorf("failed to save commitment: %w", err)
		}

		commitment.ID = id
		commitment.CreatedAt = createdAt
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.db.Close()
}
