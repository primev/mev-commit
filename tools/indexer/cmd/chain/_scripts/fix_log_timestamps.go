package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := flag.String("dsn", "", "Database DSN (required)")
	dryRun := flag.Bool("dry-run", false, "Dry run mode (don't update database)")
	chunkSize := flag.Int("chunk-size", 10000, "Number of logs to update per chunk")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("--dsn is required")
	}

	log.Printf("Starting log timestamp fix")
	log.Printf("Chunk size: %d", *chunkSize)
	log.Printf("Dry run: %v", *dryRun)

	// Connect to database
	db, err := sql.Open("mysql", *dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Get total count of logs
	var totalLogs int64
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM logs").Scan(&totalLogs); err != nil {
		log.Fatalf("Failed to count logs: %v", err)
	}
	log.Printf("Total logs to update: %d", totalLogs)

	if *dryRun {
		log.Printf("DRY RUN - Checking if update would work...")

		// Test query to see if blocks have correct timestamps
		var sampleBlock int64
		var sampleTimestamp int64
		err := db.QueryRowContext(ctx, "SELECT number, timestamp FROM blocks ORDER BY number ASC LIMIT 1").Scan(&sampleBlock, &sampleTimestamp)
		if err != nil {
			log.Fatalf("Failed to sample blocks table: %v", err)
		}
		log.Printf("Sample from blocks table: block %d has timestamp %d", sampleBlock, sampleTimestamp)

		// Check if timestamp looks correct (should be in milliseconds, between 2020-2030)
		if sampleTimestamp < 1577836800000 || sampleTimestamp > 1893456000000 {
			log.Printf("WARNING: Sample timestamp looks suspicious. Make sure you ran fix-timestamps on blocks table first!")
		} else {
			log.Printf("Sample timestamp looks good!")
		}

		log.Printf("DRY RUN complete - no changes made")
		return
	}

	// Update logs table from blocks table
	// We'll do this in one UPDATE statement with a JOIN
	log.Printf("Updating logs.block_timestamp from blocks.timestamp...")

	updateQuery := `
		UPDATE logs
		INNER JOIN blocks ON logs.block_number = blocks.number
		SET logs.block_timestamp = blocks.timestamp
	`

	startTime := time.Now()
	result, err := db.ExecContext(ctx, updateQuery)
	if err != nil {
		log.Fatalf("Failed to update logs: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Warning: could not get rows affected: %v", err)
		rowsAffected = -1
	}

	elapsed := time.Since(startTime)

	log.Printf("========================================")
	log.Printf("Log timestamp fix completed!")
	log.Printf("Rows affected: %d", rowsAffected)
	log.Printf("Duration: %v", elapsed)
	if elapsed.Seconds() > 0 {
		log.Printf("Rate: %.2f logs/sec", float64(rowsAffected)/elapsed.Seconds())
	}
	log.Printf("========================================")
}
