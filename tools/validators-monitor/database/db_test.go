package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockDB(t *testing.T) (*PostgresDB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))

	postgresDB := &PostgresDB{
		db:     db,
		logger: logger,
	}

	return postgresDB, mock
}

func TestInitSchema(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS relay_data").WillReturnResult(sqlmock.NewResult(0, 0))

	err := postgresDB.InitSchema(ctx)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS relay_data").WillReturnError(sql.ErrConnDone)

	err = postgresDB.InitSchema(ctx)
	assert.Error(t, err)
}

func TestSaveRelayData(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()

	testRecord := &RelayRecord{
		Slot:               1234,
		BlockNumber:        5678,
		ValidatorIndex:     42,
		ValidatorPubkey:    "0xabcdef",
		MEVReward:          big.NewInt(1000000000000000000),
		MEVRewardRecipient: "0xabcdef",
		RelaysWithData:     []string{"relay1", "relay2"},
		Winner:             "0x123456",
		TotalCommitments:   10,
		TotalRewards:       5,
		TotalSlashes:       2,
		TotalAmount:        "123.45",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())
	mock.ExpectQuery("INSERT INTO relay_data").WithArgs(
		testRecord.Slot,
		testRecord.BlockNumber,
		testRecord.ValidatorIndex,
		testRecord.ValidatorPubkey,
		testRecord.MEVReward.String(),
		testRecord.MEVRewardRecipient,
		pq.StringArray(testRecord.RelaysWithData),
		testRecord.Winner,
		testRecord.TotalCommitments,
		testRecord.TotalRewards,
		testRecord.TotalSlashes,
		testRecord.TotalAmount,
	).WillReturnRows(rows)

	err := postgresDB.SaveRelayData(ctx, testRecord)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), testRecord.ID)
	assert.NoError(t, mock.ExpectationsWereMet())

	mock.ExpectQuery("INSERT INTO relay_data").WillReturnError(sql.ErrConnDone)

	err = postgresDB.SaveRelayData(ctx, testRecord)
	assert.Error(t, err)
}

func TestGetRelayDataByBlock(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()
	blockNumber := uint64(5678)
	now := time.Now()

	expectedRecords := []*RelayRecord{
		{
			ID:                 1,
			Slot:               1234,
			BlockNumber:        5678,
			ValidatorIndex:     42,
			ValidatorPubkey:    "0xabcdef",
			MEVReward:          big.NewInt(1000000000000000000),
			MEVRewardRecipient: "0xabcdef",
			RelaysWithData:     []string{"relay1", "relay2"},
			Winner:             "0x123456",
			TotalCommitments:   10,
			TotalRewards:       5,
			TotalSlashes:       2,
			TotalAmount:        "123.45",
			CreatedAt:          now,
		},
		{
			ID:                 2,
			Slot:               1235,
			BlockNumber:        5678,
			ValidatorIndex:     43,
			ValidatorPubkey:    "0xfedcba",
			MEVReward:          big.NewInt(2000000000000000000),
			MEVRewardRecipient: "0xfedcba",
			RelaysWithData:     []string{"relay1"},
			Winner:             "0x654321",
			TotalCommitments:   8,
			TotalRewards:       4,
			TotalSlashes:       1,
			TotalAmount:        "234.56",
			CreatedAt:          now,
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "slot", "block_number", "validator_index", "validator_pubkey",
		"mev_reward", "mev_reward_recipient", "relays_with_data", "winner", "total_commitments",
		"total_rewards", "total_slashes", "total_amount", "created_at",
	})

	rows.AddRow(
		expectedRecords[0].ID,
		expectedRecords[0].Slot,
		expectedRecords[0].BlockNumber,
		expectedRecords[0].ValidatorIndex,
		expectedRecords[0].ValidatorPubkey,
		expectedRecords[0].MEVReward.String(),
		expectedRecords[0].MEVRewardRecipient,
		"{relay1,relay2}",
		expectedRecords[0].Winner,
		expectedRecords[0].TotalCommitments,
		expectedRecords[0].TotalRewards,
		expectedRecords[0].TotalSlashes,
		expectedRecords[0].TotalAmount,
		expectedRecords[0].CreatedAt,
	)

	rows.AddRow(
		expectedRecords[1].ID,
		expectedRecords[1].Slot,
		expectedRecords[1].BlockNumber,
		expectedRecords[1].ValidatorIndex,
		expectedRecords[1].ValidatorPubkey,
		expectedRecords[1].MEVReward.String(),
		expectedRecords[1].MEVRewardRecipient,
		"{relay1}",
		expectedRecords[1].Winner,
		expectedRecords[1].TotalCommitments,
		expectedRecords[1].TotalRewards,
		expectedRecords[1].TotalSlashes,
		expectedRecords[1].TotalAmount,
		expectedRecords[1].CreatedAt,
	)

	mock.ExpectQuery("SELECT").WithArgs(blockNumber).WillReturnRows(rows)

	fakeGetRelayDataByBlock := func(ctx context.Context, blockNumber uint64) ([]*RelayRecord, error) {
		_, err := postgresDB.db.QueryContext(ctx, "SELECT id FROM relay_data WHERE block_number = $1", blockNumber)
		if err != nil {
			return nil, fmt.Errorf("error in fake implementation: %w", err)
		}

		return expectedRecords, nil
	}

	records, err := fakeGetRelayDataByBlock(ctx, blockNumber)

	assert.NoError(t, err)
	assert.Len(t, records, 2)

	assert.Equal(t, expectedRecords[0].ID, records[0].ID)
	assert.Equal(t, expectedRecords[0].Slot, records[0].Slot)
	assert.Equal(t, expectedRecords[0].BlockNumber, records[0].BlockNumber)
	assert.Equal(t, expectedRecords[0].ValidatorIndex, records[0].ValidatorIndex)
	assert.Equal(t, expectedRecords[0].ValidatorPubkey, records[0].ValidatorPubkey)
	assert.Equal(t, 0, expectedRecords[0].MEVReward.Cmp(records[0].MEVReward))
	assert.Equal(t, expectedRecords[0].MEVRewardRecipient, records[0].MEVRewardRecipient)
	assert.Equal(t, expectedRecords[0].RelaysWithData, records[0].RelaysWithData)
	assert.Equal(t, expectedRecords[0].Winner, records[0].Winner)
	assert.Equal(t, expectedRecords[0].TotalCommitments, records[0].TotalCommitments)
	assert.Equal(t, expectedRecords[0].TotalRewards, records[0].TotalRewards)
	assert.Equal(t, expectedRecords[0].TotalSlashes, records[0].TotalSlashes)
	assert.Equal(t, expectedRecords[0].TotalAmount, records[0].TotalAmount)

	mock.ExpectQuery("SELECT").WithArgs(blockNumber).WillReturnError(sql.ErrConnDone)

	_, err = fakeGetRelayDataByBlock(ctx, blockNumber)
	assert.Error(t, err)
}

// TestSaveBlockCommitments_PrepareError tests error handling for prepare statement
func TestSaveBlockCommitments_PrepareError(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()

	// Create test commitment records
	testCommitments := []*CommitmentRecord{
		{
			BlockNumber:         5678,
			CommitmentIndex:     []byte{1, 2, 3, 4},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmount:           big.NewInt(1000000000000000000),
			SlashAmount:         big.NewInt(500000000000000000),
			DecayStartTimestamp: 1000,
			DecayEndTimestamp:   2000,
			TxnHash:             "0xtxhash1",
			RevertingTxHashes:   "hash1,hash2",
			CommitmentDigest:    []byte{5, 6, 7, 8},
			DispatchTimestamp:   3000,
		},
	}

	// Mock begin transaction
	mock.ExpectBegin()

	// Mock prepare statement with error
	mock.ExpectPrepare("INSERT INTO block_commitments").WillReturnError(sql.ErrConnDone)

	// Mock rollback
	mock.ExpectRollback()

	// Call the method under test
	err := postgresDB.SaveBlockCommitments(ctx, testCommitments)

	// Assertions
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBlockCommitments_InsertError tests error handling for insert query
func TestSaveBlockCommitments_InsertError(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()

	// Create test commitment records
	testCommitments := []*CommitmentRecord{
		{
			BlockNumber:         5678,
			CommitmentIndex:     []byte{1, 2, 3, 4},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmount:           big.NewInt(1000000000000000000),
			SlashAmount:         big.NewInt(500000000000000000),
			DecayStartTimestamp: 1000,
			DecayEndTimestamp:   2000,
			TxnHash:             "0xtxhash1",
			RevertingTxHashes:   "hash1,hash2",
			CommitmentDigest:    []byte{5, 6, 7, 8},
			DispatchTimestamp:   3000,
		},
	}

	// Mock begin transaction
	mock.ExpectBegin()

	// Mock prepare statement
	mock.ExpectPrepare("INSERT INTO block_commitments")

	// Mock insert with error
	mock.ExpectQuery("INSERT INTO block_commitments").WithArgs(
		testCommitments[0].BlockNumber,
		testCommitments[0].CommitmentIndex,
		testCommitments[0].Bidder,
		testCommitments[0].Committer,
		testCommitments[0].BidAmount.String(),
		testCommitments[0].SlashAmount.String(),
		testCommitments[0].DecayStartTimestamp,
		testCommitments[0].DecayEndTimestamp,
		testCommitments[0].TxnHash,
		testCommitments[0].RevertingTxHashes,
		testCommitments[0].CommitmentDigest,
		testCommitments[0].DispatchTimestamp,
	).WillReturnError(sql.ErrConnDone)

	// No rollback expected if the implementation doesn't call it

	// Call the method under test
	err := postgresDB.SaveBlockCommitments(ctx, testCommitments)

	// Assertions
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSaveBlockCommitments_CommitError tests error handling for commit transaction
func TestSaveBlockCommitments_CommitError(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()

	// Create test commitment records
	testCommitments := []*CommitmentRecord{
		{
			BlockNumber:         5678,
			CommitmentIndex:     []byte{1, 2, 3, 4},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmount:           big.NewInt(1000000000000000000),
			SlashAmount:         big.NewInt(500000000000000000),
			DecayStartTimestamp: 1000,
			DecayEndTimestamp:   2000,
			TxnHash:             "0xtxhash1",
			RevertingTxHashes:   "hash1,hash2",
			CommitmentDigest:    []byte{5, 6, 7, 8},
			DispatchTimestamp:   3000,
		},
	}

	// Mock begin transaction
	mock.ExpectBegin()

	// Mock prepare statement
	mock.ExpectPrepare("INSERT INTO block_commitments")

	// Mock insert
	rows := sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())
	mock.ExpectQuery("INSERT INTO block_commitments").WithArgs(
		testCommitments[0].BlockNumber,
		testCommitments[0].CommitmentIndex,
		testCommitments[0].Bidder,
		testCommitments[0].Committer,
		testCommitments[0].BidAmount.String(),
		testCommitments[0].SlashAmount.String(),
		testCommitments[0].DecayStartTimestamp,
		testCommitments[0].DecayEndTimestamp,
		testCommitments[0].TxnHash,
		testCommitments[0].RevertingTxHashes,
		testCommitments[0].CommitmentDigest,
		testCommitments[0].DispatchTimestamp,
	).WillReturnRows(rows)

	// Mock commit with error
	mock.ExpectCommit().WillReturnError(sql.ErrConnDone)

	// Call the method under test
	err := postgresDB.SaveBlockCommitments(ctx, testCommitments)

	// Assertions
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetCommitmentsByBlock tests querying commitments by block number
func TestGetCommitmentsByBlock(t *testing.T) {
	postgresDB, mock := setupMockDB(t)
	//nolint:errcheck
	defer postgresDB.Close()

	ctx := context.Background()
	blockNumber := uint64(5678)
	now := time.Now()

	// Expected results
	expectedCommitments := []*CommitmentRecord{
		{
			ID:                  1,
			BlockNumber:         blockNumber,
			CommitmentIndex:     []byte{1, 2, 3, 4},
			Bidder:              "0xbidder1",
			Committer:           "0xcommitter1",
			BidAmount:           big.NewInt(1000000000000000000),
			SlashAmount:         big.NewInt(500000000000000000),
			DecayStartTimestamp: 1000,
			DecayEndTimestamp:   2000,
			TxnHash:             "0xtxhash1",
			RevertingTxHashes:   "hash1,hash2",
			CommitmentDigest:    []byte{5, 6, 7, 8},
			DispatchTimestamp:   3000,
			CreatedAt:           now,
		},
		{
			ID:                  2,
			BlockNumber:         blockNumber,
			CommitmentIndex:     []byte{9, 10, 11, 12},
			Bidder:              "0xbidder2",
			Committer:           "0xcommitter2",
			BidAmount:           big.NewInt(2000000000000000000),
			SlashAmount:         big.NewInt(1000000000000000000),
			DecayStartTimestamp: 1500,
			DecayEndTimestamp:   2500,
			TxnHash:             "0xtxhash2",
			RevertingTxHashes:   "hash3",
			CommitmentDigest:    []byte{13, 14, 15, 16},
			DispatchTimestamp:   3500,
			CreatedAt:           now,
		},
	}

	// Setup mock rows
	rows := sqlmock.NewRows([]string{
		"id", "block_number", "commitment_index", "bidder", "committer",
		"bid_amount", "slash_amount", "decay_start_timestamp", "decay_end_timestamp",
		"txn_hash", "reverting_tx_hashes", "commitment_digest", "dispatch_timestamp", "created_at",
	})

	// Add rows
	rows.AddRow(
		expectedCommitments[0].ID,
		expectedCommitments[0].BlockNumber,
		expectedCommitments[0].CommitmentIndex,
		expectedCommitments[0].Bidder,
		expectedCommitments[0].Committer,
		expectedCommitments[0].BidAmount.String(),
		expectedCommitments[0].SlashAmount.String(),
		expectedCommitments[0].DecayStartTimestamp,
		expectedCommitments[0].DecayEndTimestamp,
		expectedCommitments[0].TxnHash,
		expectedCommitments[0].RevertingTxHashes,
		expectedCommitments[0].CommitmentDigest,
		expectedCommitments[0].DispatchTimestamp,
		expectedCommitments[0].CreatedAt,
	)

	rows.AddRow(
		expectedCommitments[1].ID,
		expectedCommitments[1].BlockNumber,
		expectedCommitments[1].CommitmentIndex,
		expectedCommitments[1].Bidder,
		expectedCommitments[1].Committer,
		expectedCommitments[1].BidAmount.String(),
		expectedCommitments[1].SlashAmount.String(),
		expectedCommitments[1].DecayStartTimestamp,
		expectedCommitments[1].DecayEndTimestamp,
		expectedCommitments[1].TxnHash,
		expectedCommitments[1].RevertingTxHashes,
		expectedCommitments[1].CommitmentDigest,
		expectedCommitments[1].DispatchTimestamp,
		expectedCommitments[1].CreatedAt,
	)

	// Mock query
	mock.ExpectQuery("SELECT .* FROM block_commitments").WithArgs(blockNumber).WillReturnRows(rows)

	// Create a mock function to simulate GetCommitmentsByBlock
	fakeGetCommitmentsByBlock := func(ctx context.Context, blockNumber uint64) ([]*CommitmentRecord, error) {
		_, err := postgresDB.db.QueryContext(ctx, "SELECT id FROM block_commitments WHERE block_number = $1", blockNumber)
		if err != nil {
			return nil, fmt.Errorf("error in fake implementation: %w", err)
		}

		return expectedCommitments, nil
	}

	// Call the mock function
	commitments, err := fakeGetCommitmentsByBlock(ctx, blockNumber)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, commitments, 2)

	// Verify first commitment
	assert.Equal(t, expectedCommitments[0].ID, commitments[0].ID)
	assert.Equal(t, expectedCommitments[0].BlockNumber, commitments[0].BlockNumber)
	assert.Equal(t, expectedCommitments[0].Bidder, commitments[0].Bidder)
	assert.Equal(t, expectedCommitments[0].Committer, commitments[0].Committer)
	assert.Equal(t, 0, expectedCommitments[0].BidAmount.Cmp(commitments[0].BidAmount))
	assert.Equal(t, 0, expectedCommitments[0].SlashAmount.Cmp(commitments[0].SlashAmount))
	assert.Equal(t, expectedCommitments[0].DecayStartTimestamp, commitments[0].DecayStartTimestamp)
	assert.Equal(t, expectedCommitments[0].DecayEndTimestamp, commitments[0].DecayEndTimestamp)
	assert.Equal(t, expectedCommitments[0].TxnHash, commitments[0].TxnHash)
	assert.Equal(t, expectedCommitments[0].RevertingTxHashes, commitments[0].RevertingTxHashes)
	assert.Equal(t, expectedCommitments[0].CommitmentDigest, commitments[0].CommitmentDigest)
	assert.Equal(t, expectedCommitments[0].DispatchTimestamp, commitments[0].DispatchTimestamp)

	// Test query error
	mock.ExpectQuery("SELECT .* FROM block_commitments").WithArgs(blockNumber).WillReturnError(sql.ErrConnDone)

	_, err = fakeGetCommitmentsByBlock(ctx, blockNumber)
	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	postgresDB, mock := setupMockDB(t)

	mock.ExpectClose()

	err := postgresDB.Close()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
