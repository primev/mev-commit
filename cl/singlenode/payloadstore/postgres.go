package payloadstore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/primev/mev-commit/cl/types" // Import shared types
)

// PostgresRepository implements the types.PayloadRepository interface using PostgreSQL.
type PostgresRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewPostgresRepository creates a new PostgresRepository.
// It also attempts to create the necessary table if it doesn't exist.
func NewPostgresRepository(ctx context.Context, dsn string, logger *slog.Logger) (*PostgresRepository, error) {
	l := logger.With("component", "PostgresRepository")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		err := db.Close()
		if err != nil {
			l.Error("Failed to close database connection after error", "error", err)
		}
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	// Create table with enhanced schema for sequential access
	schemaCreationQuery := `
		CREATE TABLE IF NOT EXISTS execution_payloads (
			id SERIAL PRIMARY KEY,
			payload_id VARCHAR(66) NOT NULL, -- e.g., 0x... (32 bytes hex + 0x prefix)
			raw_execution_payload TEXT NOT NULL,
			block_height BIGINT NOT NULL,
			inserted_at TIMESTAMPTZ DEFAULT NOW(),
			
			-- Indexes for efficient querying
			UNIQUE(block_height)
		);
		
		-- Create indexes if they don't exist
		CREATE INDEX IF NOT EXISTS idx_block_height ON execution_payloads(block_height);
		CREATE INDEX IF NOT EXISTS idx_inserted_at ON execution_payloads(inserted_at);
	`
	execCtx, execCancel := context.WithTimeout(ctx, 10*time.Second)
	defer execCancel()
	if _, err := db.ExecContext(execCtx, schemaCreationQuery); err != nil {
		err := db.Close()
		if err != nil {
			l.Error("Failed to close database connection after error", "error", err)
		}
		return nil, fmt.Errorf("failed to create execution_payloads table: %w", err)
	}
	l.Info("Successfully connected to PostgreSQL and ensured table exists.")
	return &PostgresRepository{db: db, logger: l}, nil
}

// SavePayload saves the payload information to the database.
func (r *PostgresRepository) SavePayload(ctx context.Context, info *types.PayloadInfo) error {
	query := `
		INSERT INTO execution_payloads (payload_id, raw_execution_payload, block_height)
		VALUES ($1, $2, $3)
		ON CONFLICT (block_height) DO UPDATE
		SET payload_id = EXCLUDED.payload_id,
		    raw_execution_payload = EXCLUDED.raw_execution_payload,
		    inserted_at = NOW();
	`

	insertCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(insertCtx, query, info.PayloadID, info.ExecutionPayload, info.BlockHeight)
	if err != nil {
		r.logger.Error(
			"Failed to insert payload into postgres",
			"payload_id", info.PayloadID,
			"block_height", info.BlockHeight,
			"error", err,
		)
		return fmt.Errorf("failed to insert payload into postgres: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected > 0 {
		r.logger.Debug(
			"Payload saved to database",
			"payload_id", info.PayloadID,
			"block_height", info.BlockHeight,
		)
	} else if err == nil && rowsAffected == 0 {
		r.logger.Debug(
			"Payload already exists in database or no rows affected",
			"payload_id", info.PayloadID,
			"block_height", info.BlockHeight,
		)
	}

	return nil
}

// GetPayloadsSince retrieves payloads with block height >= sinceHeight, ordered by block height
func (r *PostgresRepository) GetPayloadsSince(ctx context.Context, sinceHeight uint64, limit int) ([]types.PayloadInfo, error) {
	query := `
		SELECT payload_id, raw_execution_payload, block_height, inserted_at
		FROM execution_payloads
		WHERE block_height >= $1
		ORDER BY block_height ASC
		LIMIT $2;
	`

	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(queryCtx, query, sinceHeight, limit)
	if err != nil {
		r.logger.Error(
			"Failed to query payloads since height",
			"since_height", sinceHeight,
			"limit", limit,
			"error", err,
		)
		return nil, fmt.Errorf("failed to query payloads since height %d: %w", sinceHeight, err)
	}
	//nolint:errcheck
	defer rows.Close()

	var payloads []types.PayloadInfo
	for rows.Next() {
		var payload types.PayloadInfo
		err := rows.Scan(
			&payload.PayloadID,
			&payload.ExecutionPayload,
			&payload.BlockHeight,
			&payload.InsertedAt,
		)
		if err != nil {
			r.logger.Error(
				"Failed to scan payload row",
				"error", err,
			)
			return nil, fmt.Errorf("failed to scan payload row: %w", err)
		}
		payloads = append(payloads, payload)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(
			"Error iterating payload rows",
			"error", err,
		)
		return nil, fmt.Errorf("error iterating payload rows: %w", err)
	}

	r.logger.Debug(
		"Retrieved payloads since height",
		"since_height", sinceHeight,
		"count", len(payloads),
		"limit", limit,
	)

	return payloads, nil
}

// GetPayloadByHeight retrieves a specific payload by block height
func (r *PostgresRepository) GetPayloadByHeight(ctx context.Context, height uint64) (*types.PayloadInfo, error) {
	query := `
		SELECT payload_id, raw_execution_payload, block_height, inserted_at
		FROM execution_payloads
		WHERE block_height = $1;
	`

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var payload types.PayloadInfo
	err := r.db.QueryRowContext(queryCtx, query, height).Scan(
		&payload.PayloadID,
		&payload.ExecutionPayload,
		&payload.BlockHeight,
		&payload.InsertedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("Payload not found for height", "height", height)
			return nil, sql.ErrNoRows
		}
		r.logger.Error(
			"Failed to query payload by height",
			"height", height,
			"error", err,
		)
		return nil, fmt.Errorf("failed to query payload by height %d: %w", height, err)
	}

	r.logger.Debug(
		"Retrieved payload by height",
		"height", height,
		"payload_id", payload.PayloadID,
	)

	return &payload, nil
}

// GetLatestPayload retrieves the most recent payload
func (r *PostgresRepository) GetLatestPayload(ctx context.Context) (*types.PayloadInfo, error) {
	query := `
		SELECT payload_id, raw_execution_payload, block_height, inserted_at
		FROM execution_payloads
		ORDER BY block_height DESC
		LIMIT 1;
	`

	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var payload types.PayloadInfo
	err := r.db.QueryRowContext(queryCtx, query).Scan(
		&payload.PayloadID,
		&payload.ExecutionPayload,
		&payload.BlockHeight,
		&payload.InsertedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debug("No payloads found in database")
			return nil, sql.ErrNoRows
		}
		r.logger.Error(
			"Failed to query latest payload",
			"error", err,
		)
		return nil, fmt.Errorf("failed to query latest payload: %w", err)
	}

	r.logger.Debug(
		"Retrieved latest payload",
		"payload_id", payload.PayloadID,
		"block_height", payload.BlockHeight,
	)

	return &payload, nil
}

// Close closes the database connection.
func (r *PostgresRepository) Close() error {
	if r.db != nil {
		r.logger.Info("Closing PostgreSQL connection")
		return r.db.Close()
	}
	return nil
}
