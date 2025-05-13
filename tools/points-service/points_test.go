package main

import (
	"database/sql"
	"io"
	"log/slog"
	"testing"
)

// TestComputePointsForMonths verifies totalPoints & preSixMonthPoints for various month scenarios.
func TestComputePointsForMonths(t *testing.T) {
	testCases := []struct {
		name    string
		months  int64
		wantTot int64
		wantPre int64
	}{
		{
			name:    "0 months => (0,0)",
			months:  0,
			wantTot: 0,
			wantPre: 0,
		},
		{
			name:    "1 month => (1000,1000)",
			months:  1,
			wantTot: 1000,
			wantPre: 1000,
		},
		{
			name:    "2 months => (2270,2000)",
			months:  2,
			wantTot: 2270,
			wantPre: 2000,
		},
		{
			name:    "5 months => (7670,5000)",
			months:  5,
			wantTot: 7670,
			wantPre: 5000,
		},
		{
			name:    "6 months => (10000,6000)",
			months:  6,
			wantTot: 10000,
			wantPre: 10000,
		},
		{
			name:    "7 months => (11983,6000)",
			months:  7,
			wantTot: 11983,
			wantPre: 11000,
		},
		{
			name:    "11 months => (25317,6000)",
			months:  11,
			wantTot: 25317,
			wantPre: 15000,
		},
		{
			name:    "12 months => (30000,6000)",
			months:  12,
			wantTot: 30000,
			wantPre: 30000,
		},
		{
			name:    "13 months => (34683,6000)",
			months:  13,
			wantTot: 34683,
			wantPre: 31000,
		},
		{
			name:    "14 months => (39367,6000)",
			months:  14,
			wantTot: 39367,
			wantPre: 32000,
		},
		{
			name:    "17 months => (53417,6000)",
			months:  17,
			wantTot: 53417,
			wantPre: 35000,
		},
		{
			name:    "18 months => (58100,6000)",
			months:  18,
			wantTot: 58100,
			wantPre: 58100,
		},
		{
			name:    "19 months => (58100,6000) beyond 18 => cap",
			months:  19,
			wantTot: 58100,
			wantPre: 58100,
		},
	}

	for _, tc := range testCases {
		blocksActive := tc.months * blocksInOneMonth
		gotTot, gotPre := computePointsForMonths(blocksActive)

		if gotTot != tc.wantTot || gotPre != tc.wantPre {
			t.Errorf("%s: months=%d => got (tot=%d, pre=%d), want (tot=%d, pre=%d)",
				tc.name, tc.months, gotTot, gotPre, tc.wantTot, tc.wantPre)
		}
	}
}

func TestManualPointsEntry(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	//nolint:errcheck
	defer db.Close()

	if _, err := db.Exec(createTableValidatorRecordsQuery); err != nil {
		t.Fatalf("failed to create validator_records table: %v", err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	insertOptIn(db, logger, "0x123", "0x456", "vanilla", "staked", 100)

	if err := insertManualValRecord(db, "0x12345", "0x45678", 90); err != nil {
		t.Fatalf("failed to insert manual val record: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM validator_records").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 records, got %d", count)
	}

	var registryType, eventType, optedOutBlock sql.NullString
	err = db.QueryRow("SELECT registry_type, event_type, opted_out_block FROM validator_records WHERE pubkey = '0x12345'").Scan(&registryType, &eventType, &optedOutBlock)
	if err != nil {
		t.Fatalf("failed to query registry_type and event_type: %v", err)
	}
	if registryType.Valid || eventType.Valid || optedOutBlock.Valid {
		t.Errorf("expected registry_type, event_type and opted_out_block to be NULL, got %v, %v and %v",
			registryType, eventType, optedOutBlock)
	}

	// manual opt out should only be allowed for manually inserted records
	err = insertManualValRecord(db, "0x123", "0x456", 200)
	if err == nil {
		t.Fatalf("expected error for manual opt out on non-manually inserted record")
	}

	err = insertManualOptOut(db, logger, "0x12345", "0x45678", 2227)
	if err != nil {
		t.Fatalf("failed to insert manual opt out: %v", err)
	}

	// opt out should now exist
	var optedOutBlockAfter sql.NullString
	err = db.QueryRow("SELECT opted_out_block FROM validator_records WHERE pubkey = '0x12345'").Scan(&optedOutBlockAfter)
	if err != nil {
		t.Fatalf("failed to query opted_out_block: %v", err)
	}
	if !optedOutBlockAfter.Valid {
		t.Errorf("expected opted_out_block to be non-NULL, got NULL")
	}
	if optedOutBlockAfter.String != "2227" {
		t.Errorf("expected opted_out_block to be 2227, got %s", optedOutBlockAfter.String)
	}

	// opt out replacement should be allowed
	err = insertManualOptOut(db, logger, "0x12345", "0x45678", 4001)
	if err != nil {
		t.Fatalf("expected no error for manual opt out, got %v", err)
	}
	var optedOutBlockAfter2 sql.NullString
	err = db.QueryRow("SELECT opted_out_block FROM validator_records WHERE pubkey = '0x12345'").Scan(&optedOutBlockAfter2)
	if err != nil {
		t.Fatalf("failed to query opted_out_block: %v", err)
	}
	if optedOutBlockAfter2.String != "4001" {
		t.Errorf("expected opted_out_block to be 4001, got %s", optedOutBlockAfter2.String)
	}
}
