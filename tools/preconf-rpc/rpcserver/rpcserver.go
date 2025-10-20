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
		cache:  cache,
		logger: logger,
	}, nil
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

	s.logger.Debug("Received JSON-RPC request", "method", r.Method)

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
		s.writeError(w, nil, CodeParseError, "Failed to parse request")
		return
	}

	if req.JSONRPC != "2.0" {
		s.writeError(w, nil, CodeInvalidRequest, "Invalid JSON-RPC version")
		return
	}

	start := time.Now()
	defer func() {
		s.logger.Info("Request processing time", "method", req.Method, "id", req.ID, "duration", time.Since(start))
	}()

	if cacheMethods[req.Method] {
		if stubbed, resp := maybeStubERC20Meta(req.Method, req.Params); stubbed {
			s.writeResponse(w, req.ID, &resp)
			return
		}
		key := cacheKey(req.Method, req.Params)
		if entry, ok := s.cache.Get(key); ok && time.Now().Before(entry.until) {
			s.logger.Debug("Cache hit", "method", req.Method, "id", req.ID)
			s.writeResponse(w, req.ID, &entry.data)
			return
		}
	}

	handleProxy := func() {
		out, statusCode, err := s.proxyRequest(r.Context(), body)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}
		var resp jsonRPCResponse
		if err := json.Unmarshal(out, &resp); err != nil {
			http.Error(w, "Failed to parse proxy response", http.StatusInternalServerError)
			return
		}
		if resp.Error != nil {
			s.writeError(w, req.ID, resp.Error.Code, resp.Error.Message)
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
		return
	}

	s.writeResponse(w, req.ID, &resp)
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

// short-circuit a few very common ERC-20 metadata calls
// selectors: symbol 0x95d89b41, decimals 0x313ce567, name 0x06fdde03
func maybeStubERC20Meta(method string, params []any) (stubbed bool, resp json.RawMessage) {
	if method != "eth_call" || len(params) == 0 {
		return false, nil
	}
	// parse minimal: [{"to":"0x..","data":"0x...."}, <blockTag?>]
	callObj, ok := params[0].(map[string]any)
	if !ok {
		return false, nil
	}
	data, _ := callObj["data"].(string)
	switch strings.ToLower(data) {
	case "0x95d89b41": // symbol()
		// return ABI-encoded string "TOKEN"
		enc := "0x" +
			"0000000000000000000000000000000000000000000000000000000000000020" + // offset
			"0000000000000000000000000000000000000000000000000000000000000005" + // len
			"544f4b454e000000000000000000000000000000000000000000000000000000" // "TOKEN"
		return true, json.RawMessage(`"` + enc + `"`)
	case "0x313ce567": // decimals()
		enc := "0x" + "0000000000000000000000000000000000000000000000000000000000000012" // 18
		return true, json.RawMessage(`"` + enc + `"`)
	case "0x06fdde03": // name()
		enc := "0x" +
			"0000000000000000000000000000000000000000000000000000000000000020" +
			"0000000000000000000000000000000000000000000000000000000000000005" +
			"546f6b656e000000000000000000000000000000000000000000000000000000"
		return true, json.RawMessage(`"` + enc + `"`)
	default:
		return false, nil
	}
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
