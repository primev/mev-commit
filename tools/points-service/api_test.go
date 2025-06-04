package main

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllPointsNoResult(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close database: %v", err)
		}
	}()

	if _, err := db.Exec(createTableValidatorRecordsQuery); err != nil {
		t.Fatalf("failed to create validator_records table: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	api := PointsAPI{logger: logger, db: db}

	req, err := http.NewRequest("GET", "all?block_number=21831800&limit=100&offset=900", nil)
	if err != nil {
		t.Fatalf("failed creating request: %v", err)
	}
	rr := httptest.NewRecorder()
	api.GetAllPoints(rr, req)
	bodyBytes, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("failed reading response body: %v", err)
	}
	expected := "[]\n"
	if string(bodyBytes) != expected {
		t.Errorf("expected response to be %q, got %q", expected, string(bodyBytes))
	}
}
