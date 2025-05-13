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

		// Convert mevReward from string back to big.Int
		r.MEVReward = new(big.Int)
		r.MEVReward.SetString(mevRewardStr, 10)

		records = append(records, &r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return records, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.db.Close()
}
