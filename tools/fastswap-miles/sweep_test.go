package main

import (
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	fastsettlement "github.com/primev/mev-commit/contracts-abi/clients/FastSettlementV3"
)

func newTestEvent() *fastsettlement.Fastsettlementv3IntentExecuted {
	return &fastsettlement.Fastsettlementv3IntentExecuted{
		User:        common.HexToAddress("0xabc"),
		InputToken:  common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // USDC
		OutputToken: common.HexToAddress("0x0000000000000000000000000000000000000000"), // ETH
		InputAmt:    big.NewInt(100_000_000),
		UserAmtOut:  big.NewInt(50_000_000_000_000_000),
		Surplus:     big.NewInt(500_000_000_000_000),
		Raw: types.Log{
			TxHash:      common.HexToHash("0xdead"),
			BlockNumber: 12345,
		},
	}
}

// TestInsertEvent_BackfillsNullMetadataOnExistingRow verifies the critical
// idempotency guarantee: when a row with this tx_hash already exists in
// fastswap_miles, insertEvent must NOT issue an INSERT (which would UPSERT
// under StarRocks PK semantics, wiping derived columns and causing the
// 2026-04-16 double-credit incident). Instead it issues a COALESCE-only
// UPDATE that fills in NULL gas_cost or block_timestamp from a later rescan
// while preserving every derived column (processed, miles, surplus_eth,
// net_profit_eth, bid_cost).
func TestInsertEvent_BackfillsNullMetadataOnExistingRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer func() { _ = db.Close() }()

	txHash := "0xdead"
	blockTS := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM mevcommit_57173.fastswap_miles WHERE tx_hash = ?)",
	)).WithArgs(txHash).WillReturnRows(
		sqlmock.NewRows([]string{"exists"}).AddRow(true),
	)
	// COALESCE UPDATE that only touches the two fields that can legitimately
	// be NULL from a prior RPC-failed insert. No INSERT must fire.
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE mevcommit_57173.fastswap_miles\n"+
			"SET gas_cost = COALESCE(gas_cost, ?),\n"+
			"    block_timestamp = COALESCE(block_timestamp, ?)\n"+
			"WHERE tx_hash = ?",
	)).WithArgs("1000", blockTS, txHash).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := insertEvent(db, txHash, 12345, &blockTS, newTestEvent(), big.NewInt(1000), "eth_weth"); err != nil {
		t.Fatalf("insertEvent returned error on existing row: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}

// TestInsertEvent_InsertsWhenRowDoesNotExist verifies that insertEvent still
// inserts fresh rows. The idempotency check must not break the base case.
func TestInsertEvent_InsertsWhenRowDoesNotExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer func() { _ = db.Close() }()

	txHash := "0xdead"
	blockTS := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM mevcommit_57173.fastswap_miles WHERE tx_hash = ?)",
	)).WithArgs(txHash).WillReturnRows(
		sqlmock.NewRows([]string{"exists"}).AddRow(false),
	)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO mevcommit_57173.fastswap_miles")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := insertEvent(db, txHash, 12345, &blockTS, newTestEvent(), big.NewInt(1000), "eth_weth"); err != nil {
		t.Fatalf("insertEvent returned error on new row: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}

// TestInsertEvent_PropagatesExistenceCheckError verifies that a DB error on
// the SELECT EXISTS returns an error rather than falling through to INSERT —
// failing closed preserves the idempotency guarantee under DB trouble.
func TestInsertEvent_PropagatesExistenceCheckError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer func() { _ = db.Close() }()

	txHash := "0xdead"
	blockTS := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT EXISTS(SELECT 1 FROM mevcommit_57173.fastswap_miles WHERE tx_hash = ?)",
	)).WithArgs(txHash).WillReturnError(errForceTest)
	// No INSERT expected.

	err = insertEvent(db, txHash, 12345, &blockTS, newTestEvent(), big.NewInt(1000), "eth_weth")
	if err == nil {
		t.Fatalf("expected error from existence check, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock expectations: %v", err)
	}
}

// errForceTest is a sentinel used to force an error path in tests.
var errForceTest = sqlmockErr("forced test error")

type sqlmockErr string

func (e sqlmockErr) Error() string { return string(e) }
