package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/primev/mev-commit/cl/types"
)

// StateManager interface for accessing block build state
type StateManager interface {
	GetBlockBuildState(ctx context.Context) types.BlockBuildState
}

// PayloadServer provides HTTP API for member nodes to fetch payloads
type PayloadServer struct {
	logger       *slog.Logger
	stateManager StateManager
	payloadRepo  types.PayloadRepository
	server       *http.Server
}

// NewPayloadServer creates a new payload API server
func NewPayloadServer(
	addr string,
	stateManager StateManager,
	payloadRepo types.PayloadRepository,
	logger *slog.Logger,
) *PayloadServer {
	mux := http.NewServeMux()

	ps := &PayloadServer{
		logger:       logger.With("component", "PayloadServer"),
		stateManager: stateManager,
		payloadRepo:  payloadRepo,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	// Register endpoints
	mux.HandleFunc("/api/v1/payload/latest", ps.handleGetLatestPayload)
	mux.HandleFunc("/api/v1/payload/since/", ps.handleGetPayloadsSince)
	mux.HandleFunc("/api/v1/payload/height/", ps.handleGetPayloadByHeight)
	mux.HandleFunc("/api/v1/health", ps.handleHealth)

	return ps
}

// Start starts the HTTP server
func (ps *PayloadServer) Start(ctx context.Context) error {
	ps.logger.Info("Starting payload API server", "addr", ps.server.Addr)

	// Start server in goroutine
	go func() {
		if err := ps.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ps.logger.Error("Payload API server error", "error", err)
		}
	}()

	// Wait for context cancellation to shutdown
	<-ctx.Done()
	return ps.Stop()
}

// Stop gracefully stops the HTTP server
func (ps *PayloadServer) Stop() error {
	ps.logger.Info("Stopping payload API server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ps.server.Shutdown(ctx)
}

// convertToResponse converts types.PayloadInfo to PayloadResponse
func convertToResponse(payload *types.PayloadInfo) PayloadResponse {
	return PayloadResponse{
		PayloadID:        payload.PayloadID,
		ExecutionPayload: payload.ExecutionPayload,
		BlockHeight:      payload.BlockHeight,
		Timestamp:        payload.InsertedAt.Unix(),
	}
}

// handleGetLatestPayload returns the latest payload from the current block build state
func (ps *PayloadServer) handleGetLatestPayload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ps.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Try to get from repository first if available
	if ps.payloadRepo != nil {
		if payload, err := ps.payloadRepo.GetLatestPayload(ctx); err == nil && payload != nil {
			response := convertToResponse(payload)
			ps.writeJSON(w, response, http.StatusOK)
			ps.logger.Debug(
				"Served latest payload from repository",
				"payload_id", payload.PayloadID,
				"height", payload.BlockHeight,
			)
			return
		}
	}

	// Fallback to state manager
	state := ps.stateManager.GetBlockBuildState(ctx)

	if state.PayloadID == "" || state.ExecutionPayload == "" {
		ps.writeError(w, "No payload available", http.StatusNotFound)
		return
	}

	response := PayloadResponse{
		PayloadID:        state.PayloadID,
		ExecutionPayload: state.ExecutionPayload,
		BlockHeight:      0, // We don't have height from state manager
		Timestamp:        time.Now().Unix(),
	}

	ps.writeJSON(w, response, http.StatusOK)
	ps.logger.Debug("Served latest payload from state", "payload_id", state.PayloadID)
}

// handleGetPayloadsSince returns payloads since a given block height
func (ps *PayloadServer) handleGetPayloadsSince(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ps.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ps.payloadRepo == nil {
		ps.writeError(w, "Payload repository not available", http.StatusServiceUnavailable)
		return
	}

	// Extract height from URL path: /api/v1/payload/since/{height}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		ps.writeError(w, "Block height required", http.StatusBadRequest)
		return
	}

	heightStr := pathParts[5]
	height, err := strconv.ParseUint(heightStr, 10, 64)
	if err != nil {
		ps.writeError(w, "Invalid block height", http.StatusBadRequest)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	payloads, err := ps.payloadRepo.GetPayloadsSince(ctx, height, limit+1) // +1 to check if there are more
	if err != nil {
		ps.logger.Error(
			"Failed to get payloads since height",
			"height", height,
			"error", err,
		)
		ps.writeError(w, "Failed to retrieve payloads", http.StatusInternalServerError)
		return
	}

	// Check if there are more payloads
	hasMore := len(payloads) > limit
	if hasMore {
		payloads = payloads[:limit] // Remove the extra payload
	}

	// Convert to response format
	responsePayloads := make([]PayloadResponse, len(payloads))
	var nextHeight uint64
	for i, payload := range payloads {
		responsePayloads[i] = convertToResponse(&payload)
		if i == len(payloads)-1 {
			nextHeight = payload.BlockHeight + 1
		}
	}

	response := PayloadListResponse{
		Payloads:   responsePayloads,
		HasMore:    hasMore,
		NextHeight: nextHeight,
		TotalCount: len(responsePayloads),
	}

	ps.writeJSON(w, response, http.StatusOK)
	ps.logger.Debug(
		"Served payloads since height",
		"since_height", height,
		"count", len(responsePayloads),
		"has_more", hasMore,
		"next_height", nextHeight,
	)
}

// handleGetPayloadByHeight returns a specific payload by block height
func (ps *PayloadServer) handleGetPayloadByHeight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ps.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ps.payloadRepo == nil {
		ps.writeError(w, "Payload repository not available", http.StatusServiceUnavailable)
		return
	}

	// Extract height from URL path: /api/v1/payload/height/{height}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		ps.writeError(w, "Block height required", http.StatusBadRequest)
		return
	}

	heightStr := pathParts[5]
	height, err := strconv.ParseUint(heightStr, 10, 64)
	if err != nil {
		ps.writeError(w, "Invalid block height", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	payload, err := ps.payloadRepo.GetPayloadByHeight(ctx, height)
	if err != nil {
		if err == sql.ErrNoRows {
			ps.writeError(w, "Payload not found", http.StatusNotFound)
		} else {
			ps.logger.Error(
				"Failed to get payload by height",
				"height", height,
				"error", err,
			)
			ps.writeError(w, "Failed to retrieve payload", http.StatusInternalServerError)
		}
		return
	}

	response := convertToResponse(payload)
	ps.writeJSON(w, response, http.StatusOK)
	ps.logger.Debug(
		"Served payload by height",
		"height", height,
		"payload_id", payload.PayloadID,
	)
}

// handleHealth returns server health status
func (ps *PayloadServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ps.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		ps.logger.Error(
			"Failed to write health response",
			"error", err,
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// writeJSON writes a JSON response
func (ps *PayloadServer) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		ps.logger.Error(
			"Failed to encode JSON response",
			"error", err,
		)
	}
}

// writeError writes an error response
func (ps *PayloadServer) writeError(w http.ResponseWriter, message string, statusCode int) {
	response := ErrorResponse{
		Error:   message,
		Code:    statusCode,
		Message: message,
	}

	ps.writeJSON(w, response, statusCode)
}
