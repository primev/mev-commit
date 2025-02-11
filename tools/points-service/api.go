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
	r := mux.NewRouter()
	r.HandleFunc("/health", p.HealthCheck).Methods("GET")
	r.HandleFunc("/last_block", p.GetLastBlock).Methods("GET")
	r.HandleFunc("/{receiver_type}/{receiver_address}", p.RecomputePointsForAddress).Methods("GET")
	r.HandleFunc("/stats", p.GetTotalPointsStats).Methods("GET")
	r.HandleFunc("/all", p.GetAllPoints).Methods("GET")

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
	w.Write([]byte("OK\n"))
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

	var totalPoints int64
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
