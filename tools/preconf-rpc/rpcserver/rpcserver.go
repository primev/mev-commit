package rpcserver

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultTimeout     = 5 * time.Second
	defaultMaxBodySize = 30 * 1024 * 1024 // 30 MB
	cacheSize          = 10000

	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeCustomError    = -32000
)

type JSONErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *JSONErr) Error() string {
	return e.Message
}

func NewJSONErr(code int, message string) *JSONErr {
	return &JSONErr{
		Code:    code,
		Message: message,
	}
}

type jsonRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      any              `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError    `json:"error,omitempty"`
}

type methodHandler func(ctx context.Context, params ...interface{}) (json.RawMessage, bool, error)

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *any   `json:"data,omitempty"`
}

var cacheMethods = map[string]bool{
	"eth_call":         true,
	"eth_getCode":      true,
	"eth_getStorageAt": true,
	"eth_getLogs":      true,
	"net_version":      true,
}

type cacheEntry struct {
	until time.Time
	data  json.RawMessage
}

func cacheKey(method string, params []any) string {
	b, _ := json.Marshal(params)
	h := sha1.Sum(append([]byte(method), b...))
	return string(h[:])
}

type JSONRPCServer struct {
	rwLock     sync.RWMutex
	methods    map[string]methodHandler
	proxyURL   string
	httpClient *http.Client
	cache      *lru.Cache[string, cacheEntry]
	metrics    *metrics
	logger     *slog.Logger
}

func NewJSONRPCServer(proxyURL string, logger *slog.Logger) (*JSONRPCServer, error) {
	cache, err := lru.New[string, cacheEntry](cacheSize)
	if err != nil {
		return nil, err
	}
	return &JSONRPCServer{
		proxyURL: proxyURL,
		methods:  make(map[string]methodHandler),
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment,
				MaxIdleConns:        256,
				MaxIdleConnsPerHost: 256,
				IdleConnTimeout:     90 * time.Second,
				ForceAttemptHTTP2:   true,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
			},
			Timeout: 15 * time.Second,
		},
		cache:   cache,
		metrics: newMetrics(),
		logger:  logger,
	}, nil
}

func (s *JSONRPCServer) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		s.metrics.methodSuccessCounts,
		s.metrics.methodFailureCounts,
		s.metrics.methodSuccessDurations,
		s.metrics.methodFailureDurations,
		s.metrics.proxyMethodSuccessCounts,
		s.metrics.proxyMethodFailureCounts,
		s.metrics.proxyMethodSuccessDurations,
		s.metrics.proxyMethodFailureDurations,
	}
}

func (s *JSONRPCServer) RegisterHandler(method string, handler methodHandler) {
	s.rwLock.Lock()
	s.methods[method] = handler
	s.rwLock.Unlock()
}

func (s *JSONRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, defaultMaxBodySize)
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, nil, CodeInvalidRequest, "Failed to read request body")
		return
	}

	var req jsonRPCRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		s.logger.Error("Failed to parse JSON-RPC request", "error", err, "body", string(body))
		s.writeError(w, nil, CodeParseError, "Failed to parse request")
		return
	}

	if req.JSONRPC != "2.0" {
		s.writeError(w, nil, CodeInvalidRequest, "Invalid JSON-RPC version")
		return
	}

	start := time.Now()

	if cacheMethods[req.Method] {
		key := cacheKey(req.Method, req.Params)
		if entry, ok := s.cache.Get(key); ok && time.Now().Before(entry.until) {
			s.logger.Debug("Cache hit", "method", req.Method, "id", req.ID)
			s.writeResponse(w, req.ID, &entry.data)
			s.metrics.methodSuccessCounts.WithLabelValues(req.Method).Inc()
			s.metrics.methodSuccessDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
			return
		}
	}

	handleProxy := func() {
		out, statusCode, err := s.proxyRequest(r.Context(), body)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			s.metrics.proxyMethodFailureCounts.WithLabelValues(req.Method).Inc()
			s.metrics.proxyMethodFailureDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
			return
		}
		var resp jsonRPCResponse
		if err := json.Unmarshal(out, &resp); err != nil {
			http.Error(w, "Failed to parse proxy response", http.StatusInternalServerError)
			s.metrics.proxyMethodFailureCounts.WithLabelValues(req.Method).Inc()
			s.metrics.proxyMethodFailureDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
			return
		}
		if resp.Error != nil {
			s.writeError(w, req.ID, resp.Error.Code, resp.Error.Message)
			s.metrics.proxyMethodFailureCounts.WithLabelValues(req.Method).Inc()
			s.metrics.proxyMethodFailureDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
			return
		}
		if cacheMethods[req.Method] && resp.Result != nil {
			key := cacheKey(req.Method, req.Params)
			_ = s.cache.Add(key, cacheEntry{
				until: time.Now().Add(pickTTL(req.Method, *resp.Result)),
				data:  *resp.Result,
			})
			s.logger.Debug("Cache store", "method", req.Method, "id", req.ID)
		}
		s.writeResponse(w, req.ID, resp.Result)
		s.metrics.proxyMethodSuccessCounts.WithLabelValues(req.Method).Inc()
		s.metrics.proxyMethodSuccessDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
	}

	s.rwLock.RLock()
	handler, ok := s.methods[req.Method]
	s.rwLock.RUnlock()
	if !ok {
		handleProxy()
		return
	}

	resp, proxy, err := handler(r.Context(), req.Params...)
	switch {
	case err != nil:
		defer func() {
			s.metrics.methodFailureCounts.WithLabelValues(req.Method).Inc()
			s.metrics.methodFailureDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
		}()
		var jsonErr *JSONErr
		if ok := errors.As(err, &jsonErr); ok {
			// If the error is a JSONErr, we can use it directly.
			s.writeError(w, req.ID, jsonErr.Code, jsonErr.Message)
			return
		}
		s.writeError(w, req.ID, CodeCustomError, err.Error())
		return
	case proxy:
		handleProxy()
		return
	case resp == nil:
		s.writeError(w, req.ID, CodeCustomError, "No response")
		s.metrics.methodFailureCounts.WithLabelValues(req.Method).Inc()
		s.metrics.methodFailureDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
		return
	}

	s.writeResponse(w, req.ID, &resp)
	s.metrics.methodSuccessCounts.WithLabelValues(req.Method).Inc()
	s.metrics.methodSuccessDurations.WithLabelValues(req.Method).Observe(float64(time.Since(start).Milliseconds()))
}

func (s *JSONRPCServer) writeResponse(w http.ResponseWriter, id any, result *json.RawMessage) {
	s.logger.Debug("Writing JSON-RPC response", "id", id, "result", result)
	response := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
		Error:   nil,
	}
	setCorsHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (s *JSONRPCServer) writeError(w http.ResponseWriter, id any, code int, message string) {
	s.logger.Error("JSON-RPC error", "id", id, "code", code, "message", message)
	response := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  nil,
		Error: &jsonRPCError{
			Code:    code,
			Message: message,
			Data:    nil,
		},
	}
	setCorsHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write error response", http.StatusInternalServerError)
		return
	}
}

func (s *JSONRPCServer) proxyRequest(ctx context.Context, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.proxyURL, bytes.NewReader(body))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create proxy request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	s.logger.Debug("Proxying request", "url", s.proxyURL, "body", string(body))
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("proxy request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("proxy request returned status %d", resp.StatusCode)
	}

	out, err := io.ReadAll(io.LimitReader(resp.Body, defaultMaxBodySize))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read proxy response: %w", err)
	}

	return out, resp.StatusCode, nil
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func pickTTL(method string, params json.RawMessage) time.Duration {
	switch method {
	case "net_version":
		return 24 * time.Hour
	case "eth_getCode":
		return 24 * time.Hour
	case "eth_call":
		// if block tag provided and hex number → immutable
		if strings.HasSuffix(string(params), "\"") { // cheap check
			if strings.Contains(string(params), "\"0x") && !strings.Contains(string(params), "\"latest\"") {
				return 24 * time.Hour
			}
		}
		return 1 * time.Millisecond
	default:
		return 2 * time.Second
	}
}
