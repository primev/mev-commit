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

// PointsAPI holds references needed to handle API requests
type PointsAPI struct {
	logger *slog.Logger
	db     *sql.DB
	ps     *PointsService
}

// NewPointsAPI creates a new PointsAPI instance
func NewPointsAPI(logger *slog.Logger, db *sql.DB, ps *PointsService) *PointsAPI {
	return &PointsAPI{logger: logger, db: db, ps: ps}
}

// StartAPIServer starts the HTTP server with all routes
func (p *PointsAPI) StartAPIServer(ctx context.Context, addr string) error {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/health", p.HealthCheck).Methods("GET")
	r.HandleFunc("/last_block", p.GetLastBlock).Methods("GET")
	r.HandleFunc("/{receiver_type}/{receiver_address}", p.GetPointsForAddress).Methods("GET")
	r.HandleFunc("/stats", p.GetTotalPointsStats).Methods("GET")
	r.HandleFunc("/all", p.GetAllPoints).Methods("GET")

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Graceful shutdown
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

// HealthCheck endpoint checks if the API is "healthy"
// GET /health
func (p *PointsAPI) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Ensure both the points routine and event subscription are running
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

// GetLastBlock returns the last processed block number
// GET /last_block
func (p *PointsAPI) GetLastBlock(w http.ResponseWriter, r *http.Request) {
	block, err := p.ps.LastBlock()
	if err != nil {
		http.Error(w, "Failed to fetch last block", http.StatusInternalServerError)
		return
	}
	resp := map[string]uint64{"last_block_number": block}
	writeJSON(w, resp, http.StatusOK)
}

// GetPointsForAddress returns points for a specific address and block number
// GET /{receiver_type}/{receiver_address}?block_number=...
func (p *PointsAPI) GetPointsForAddress(w http.ResponseWriter, r *http.Request) {
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

	// Example DB logic: gather points from the 'validator_records' table
	rows, err := p.db.Query(`
		SELECT vault, registry_type, points_accumulated
		FROM validator_records
		WHERE (pubkey = ? OR adder = ?)
		  AND opted_in_block <= ?
		  AND (opted_out_block IS NULL OR opted_out_block > ?)
	`, receiverAddr, receiverAddr, blockNum, blockNum)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pointsArray []map[string]interface{}
	for rows.Next() {
		var vaultAddr, registryType string
		var pointsAccumulated int64
		if scanErr := rows.Scan(&vaultAddr, &registryType, &pointsAccumulated); scanErr != nil {
			p.logger.Error("scan error", "error", scanErr)
			continue
		}
		// For demonstration, map registryType to "network_address"/"vault_address" usage
		pItem := map[string]interface{}{
			"points": float64(pointsAccumulated),
		}
		if registryType == "vanilla" || registryType == "eigenlayer" {
			pItem["network_address"] = registryType
			pItem["vault_address"] = vaultAddr
		} else {
			pItem["network_address"] = ""
			pItem["vault_address"] = vaultAddr
		}
		pointsArray = append(pointsArray, pItem)
	}
	if len(pointsArray) == 0 {
		http.Error(w, "no points data found for address", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"address":      receiverAddr,
		"receiver":     receiverType,
		"block_number": blockNum,
		"points":       pointsArray,
	}
	writeJSON(w, resp, http.StatusOK)
}

// GetTotalPointsStats returns aggregated points stats for a given block
// GET /stats?receiver_type=...&block_number=...
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
	receiverType := r.URL.Query().Get("receiver_type")

	// Filter by registry_type if given
	query := `
		SELECT SUM(points_accumulated), 
		       COUNT(DISTINCT CASE WHEN registry_type = 'vanilla' THEN pubkey END) as stakers,
		       COUNT(DISTINCT CASE WHEN registry_type = 'symbiotic' THEN pubkey END) as networks,
		       COUNT(DISTINCT CASE WHEN registry_type = 'eigenlayer' THEN pubkey END) as operators
		FROM validator_records
		WHERE opted_in_block <= ?
		  AND (opted_out_block IS NULL OR opted_out_block > ?)
	`
	args := []interface{}{blockNum, blockNum}
	if receiverType != "" {
		// Minimal mapping from external "staker"/"network"/"operator" to our internal registry_type
		var rt string
		switch receiverType {
		case "staker":
			rt = "vanilla"
		case "network":
			rt = "symbiotic"
		case "operator":
			rt = "eigenlayer"
		default:
			http.Error(w, "invalid receiver_type", http.StatusBadRequest)
			return
		}
		query += " AND registry_type = ?"
		args = append(args, rt)
	}

	var totalPoints, stakers, networks, operators sql.NullInt64
	err = p.db.QueryRow(query, args...).Scan(&totalPoints, &stakers, &networks, &operators)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"total_points": totalPoints.Int64,
		"stakers":      stakers.Int64,
		"networks":     networks.Int64,
		"operators":    operators.Int64,
	}
	writeJSON(w, resp, http.StatusOK)
}

// GetAllPoints returns paginated list of all points for the specified block
// GET /all?offset=...&limit=...&receiver_type=...&block_number=...
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

	// Build query and args
	query := `
		SELECT pubkey, adder, vault, registry_type, points_accumulated
		FROM validator_records
		WHERE opted_in_block <= ?
		  AND (opted_out_block IS NULL OR opted_out_block > ?)
	`
	args := []interface{}{blockNum, blockNum}

	if receiverType != "" {
		// Minimal mapping from external "staker"/"network"/"operator" to our internal registry_type
		var rt string
		switch receiverType {
		case "staker":
			rt = "vanilla"
		case "network":
			rt = "symbiotic"
		case "operator":
			rt = "eigenlayer"
		default:
			http.Error(w, "invalid receiver_type", http.StatusBadRequest)
			return
		}
		query += " AND registry_type = ?"
		args = append(args, rt)
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := p.db.Query(query, args...)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var pubkey, adder, vault, registryType string
		var points int64
		if scanErr := rows.Scan(&pubkey, &adder, &vault, &registryType, &points); scanErr != nil {
			p.logger.Error("scan error", "error", scanErr)
			continue
		}

		// Convert internal registry_type -> external "receiver"
		var receiver string
		switch registryType {
		case "vanilla":
			receiver = "staker"
		case "symbiotic":
			receiver = "network"
		case "eigenlayer":
			receiver = "operator"
		default:
			receiver = "unknown"
		}

		item := map[string]interface{}{
			"address":      pubkey, // or adder, depending on your usage
			"receiver":     receiver,
			"block_number": blockNum,
			"network_address": func() string {
				if registryType == "vanilla" || registryType == "eigenlayer" {
					return registryType
				}
				return ""
			}(),
			"vault_address": vault,
			"points":        float64(points),
		}
		result = append(result, item)
	}
	writeJSON(w, result, http.StatusOK)
}

// Helper to write JSON responses
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON: %v", err)
	}
}
