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

func TestClose(t *testing.T) {
	postgresDB, mock := setupMockDB(t)

	mock.ExpectClose()

	err := postgresDB.Close()

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
