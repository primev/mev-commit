package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PointsAPI struct {
	logger *slog.Logger
	db     *sql.DB
	ps     *PointsService
}

func NewPointsAPI(logger *slog.Logger, db *sql.DB, ps *PointsService) *PointsAPI {
	return &PointsAPI{logger: logger, db: db, ps: ps}
}

func (p *PointsAPI) StartAPIServer(ctx context.Context, addr string) error {
	r := http.NewServeMux()
	r.HandleFunc("GET /health", p.HealthCheck)
	r.HandleFunc("GET /last_block", p.GetLastBlock)
	r.HandleFunc("GET /{receiver_type}/{receiver_address}", p.RecomputePointsForAddress)
	r.HandleFunc("GET /stats", p.GetTotalPointsStats)
	r.HandleFunc("GET /all", p.GetAllPoints)

	// Personal API
	r.HandleFunc("GET /{address}", p.GetAnyPointsForAddress)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			p.logger.Error("HTTP server shutdown error", "error", err)
		}
	}()

	p.logger.Info("Starting External Points API", slog.String("addr", addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (p *PointsAPI) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if !p.ps.IsPointsRoutineRunning() {
		http.Error(w, "Points routine not running", http.StatusServiceUnavailable)
		return
	}
	if !p.ps.IsSubscriptionActive() {
		http.Error(w, "Event subscription not active", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK\n")); err != nil {
		p.logger.Error("failed to write response", "error", err)
	}
}

func (p *PointsAPI) GetLastBlock(w http.ResponseWriter, r *http.Request) {
	block, err := p.ps.LastBlock()
	if err != nil {
		http.Error(w, "Failed to fetch last block", http.StatusInternalServerError)
		return
	}
	resp := map[string]uint64{"last_block_number": block}
	writeJSON(w, resp, http.StatusOK)
}

func (p *PointsAPI) RecomputePointsForAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiverType := vars["receiver_type"]
	receiverAddr := vars["receiver_address"]

	blockNumStr := r.URL.Query().Get("block_number")
	if blockNumStr == "" {
		http.Error(w, "block_number query param is required", http.StatusBadRequest)
		return
	}
	blockNum, err := strconv.ParseUint(blockNumStr, 10, 64)
	if err != nil || blockNum == 0 {
		http.Error(w, "invalid block_number parameter", http.StatusBadRequest)
		return
	}

	if receiverType != "operator" {
		resp := map[string]interface{}{
			"address":      receiverAddr,
			"receiver":     receiverType,
			"block_number": blockNum,
			"points":       float64(0),
		}
		writeJSON(w, resp, http.StatusOK)
		return
	}

	totalPoints, err := p.calculatePointsForSymbioticOperator(receiverAddr, blockNum)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	if totalPoints == 0 {
		http.Error(w, "no points data found for address", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"address":      receiverAddr,
		"receiver":     receiverType,
		"block_number": blockNum,
		"points":       float64(totalPoints),
	}
	writeJSON(w, resp, http.StatusOK)
}

func (p *PointsAPI) calculatePointsForSymbioticOperator(receiverAddr string, blockNum uint64) (int64, error) {

	// Get count of unique pubkeys for this operator
	var uniquePubkeys int
	err := p.db.QueryRow(`
		SELECT COUNT(DISTINCT pubkey)
		FROM validator_records
		WHERE registry_type = 'symbiotic'
		  AND adder = ?
		  AND opted_in_block <= ?
	`, receiverAddr, blockNum).Scan(&uniquePubkeys)
	if err != nil {
		p.logger.Error("failed to get unique pubkey count", "error", err)
		return 0, err
	}

	// Sign-up Bonus of 1000
	totalPoints := int64(uniquePubkeys) * 1000

	rows, err := p.db.Query(`
		SELECT vault, registry_type, opted_in_block, opted_out_block
		FROM validator_records
		WHERE registry_type = 'symbiotic'
		  AND (pubkey = ? OR adder = ?)
		  AND opted_in_block <= ?
	`, receiverAddr, receiverAddr, blockNum)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			vaultAddr     string
			regType       string
			optedInBlock  uint64
			optedOutBlock sql.NullInt64
		)
		if scanErr := rows.Scan(&vaultAddr, &regType, &optedInBlock, &optedOutBlock); scanErr != nil {
			p.logger.Error("scan error", "error", scanErr)
			continue
		}

		effectiveEnd := blockNum
		if optedOutBlock.Valid && optedOutBlock.Int64 <= int64(blockNum) {
			effectiveEnd = uint64(optedOutBlock.Int64)
		}

		var blocksActive int64
		if effectiveEnd > optedInBlock {
			blocksActive = int64(effectiveEnd - optedInBlock)
		}
		if blocksActive < 0 {
			blocksActive = 0
		}

		recomputedPoints, optedOutPoints := computePointsForMonths(blocksActive)
		if optedOutBlock.Valid {
			totalPoints += optedOutPoints
		} else {
			totalPoints += recomputedPoints
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return 0, rowsErr
	}
	return totalPoints, nil
}
func (p *PointsAPI) GetTotalPointsStats(w http.ResponseWriter, r *http.Request) {
	blockNumStr := r.URL.Query().Get("block_number")
	if blockNumStr == "" {
		http.Error(w, "block_number query param is required", http.StatusBadRequest)
		return
	}
	blockNum, err := strconv.ParseUint(blockNumStr, 10, 64)
	if err != nil || blockNum == 0 {
		http.Error(w, "invalid block_number parameter", http.StatusBadRequest)
		return
	}

	rows, err := p.db.Query(`
		SELECT adder, opted_in_block, opted_out_block
		FROM validator_records
		WHERE registry_type = 'symbiotic'
		  AND opted_in_block <= ?
	`, blockNum)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var totalPoints int64
	operatorSet := make(map[string]struct{})

	for rows.Next() {
		var (
			adder         string
			optedInBlock  uint64
			optedOutBlock sql.NullInt64
		)
		if scanErr := rows.Scan(&adder, &optedInBlock, &optedOutBlock); scanErr != nil {
			p.logger.Error("scan error", "error", scanErr)
			continue
		}

		operatorSet[adder] = struct{}{}

		effectiveEnd := blockNum
		if optedOutBlock.Valid && optedOutBlock.Int64 <= int64(blockNum) {
			effectiveEnd = uint64(optedOutBlock.Int64)
		}

		var blocksActive int64
		if effectiveEnd > optedInBlock {
			blocksActive = int64(effectiveEnd - optedInBlock)
		}
		if blocksActive < 0 {
			blocksActive = 0
		}

		recomputedPoints, optedOutPoints := computePointsForMonths(blocksActive)
		if optedOutBlock.Valid {
			totalPoints += optedOutPoints
		} else {
			totalPoints += recomputedPoints
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "DB iteration error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"total_points": totalPoints,
		"stakers":      0,
		"networks":     0,
		"operators":    len(operatorSet),
	}
	writeJSON(w, resp, http.StatusOK)
}

func (p *PointsAPI) GetAllPoints(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")
	blockNumStr := r.URL.Query().Get("block_number")

	if offsetStr == "" || limitStr == "" || blockNumStr == "" {
		http.Error(w, "offset, limit, and block_number are required", http.StatusBadRequest)
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil || offset < 0 {
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}
	blockNum, err := strconv.ParseUint(blockNumStr, 10, 64)
	if err != nil || blockNum == 0 {
		http.Error(w, "invalid block_number", http.StatusBadRequest)
		return
	}

	receiverType := r.URL.Query().Get("receiver_type")
	// If receiver_type is not "operator", return empty array
	if receiverType != "operator" && receiverType != "" {
		writeJSON(w, []interface{}{}, http.StatusOK)
		return
	}

	query := `
		SELECT DISTINCT adder, vault
		FROM validator_records
		WHERE registry_type = 'symbiotic'
		  AND opted_in_block <= ?
		  AND (opted_out_block IS NULL OR opted_out_block > ?)
		ORDER BY adder
		LIMIT ? OFFSET ?
	`
	rows, err := p.db.Query(query, blockNum, blockNum, limit, offset)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var operatorAddr, vaultAddr string
		if scanErr := rows.Scan(&operatorAddr, &vaultAddr); scanErr != nil {
			p.logger.Error("scan error", "error", scanErr)
			continue
		}

		// Calculate aggregated points for this operator address
		totalPoints, calcErr := p.calculatePointsForSymbioticOperator(operatorAddr, blockNum)
		if calcErr != nil {
			p.logger.Error("failed to compute points", "address", operatorAddr, "error", calcErr)
			continue
		}

		result = append(result, map[string]interface{}{
			"address":         operatorAddr,
			"receiver":        "operator",
			"block_number":    blockNum,
			"network_address": "0x9101eda106A443A0fA82375936D0D1680D5a64F5",
			"vault_address":   vaultAddr,
			"points":          float64(totalPoints),
		})
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "DB iteration error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, result, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}

// GetAnyPointsForAddress sums up points for the given address:
//   - points_accumulated if `opted_out_block` is NULL
//   - pre_cliff_points if `opted_out_block` is NOT NULL
func (p *PointsAPI) GetAnyPointsForAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adder := vars["address"]

	// Summation query uses a CASE expression for each record
	const q = `
        SELECT COALESCE(
          SUM(
            CASE WHEN opted_out_block IS NULL 
                 THEN points_accumulated 
                 ELSE pre_cliff_points 
            END
          ), 
        0)
        FROM validator_records
        WHERE adder = ?
    `
	var totalPoints int64
	if err := p.db.QueryRow(q, adder).Scan(&totalPoints); err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		p.logger.Error("GetAnyPointsForAddress query error", "error", err)
		return
	}

	// Get count of unique pubkeys for this address and multiply by 1000
	var pubkeyBonus int64
	const pubkeyQuery = `
		SELECT COUNT(DISTINCT pubkey) * 1000 
		FROM validator_records 
		WHERE adder = ?`
	if err := p.db.QueryRow(pubkeyQuery, adder).Scan(&pubkeyBonus); err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		p.logger.Error("GetAnyPointsForAddress pubkey bonus query error", "error", err)
		return
	}
	totalPoints += pubkeyBonus

	resp := map[string]interface{}{
		"address":      adder,
		"total_points": totalPoints,
	}
	writeJSON(w, resp, http.StatusOK)
}
