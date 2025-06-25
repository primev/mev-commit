package rpcserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	defaultTimeout     = 5 * time.Second
	defaultMaxBodySize = 30 * 1024 * 1024 // 30 MB

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

type JSONRPCServer struct {
	rwLock   sync.RWMutex
	methods  map[string]methodHandler
	proxyURL string
	logger   *slog.Logger
}

func NewJSONRPCServer(proxyURL string, logger *slog.Logger) *JSONRPCServer {
	return &JSONRPCServer{
		proxyURL: proxyURL,
		methods:  make(map[string]methodHandler),
		logger:   logger,
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

	s.logger.Info("Received JSON-RPC request", "method", r.Method)

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

	s.logger.Info("Processing JSON-RPC request", "method", req.Method, "id", req.ID)
	defer s.logger.Info("Finished processing JSON-RPC request", "method", req.Method, "id", req.ID)

	s.rwLock.RLock()
	handler, ok := s.methods[req.Method]
	s.rwLock.RUnlock()
	if !ok {
		s.proxyRequest(w, body)
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
		s.proxyRequest(w, body)
		return
	case resp == nil:
		s.writeError(w, req.ID, CodeCustomError, "No response")
		return
	}

	s.writeResponse(w, req.ID, &resp)
}

func (s *JSONRPCServer) writeResponse(w http.ResponseWriter, id any, result *json.RawMessage) {
	s.logger.Info("Writing JSON-RPC response", "id", id, "result", result)
	response := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
		Error:   nil,
	}
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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write error response", http.StatusInternalServerError)
		return
	}
}

func (s *JSONRPCServer) proxyRequest(w http.ResponseWriter, body []byte) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	req, err := http.NewRequest(http.MethodPost, s.proxyURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	s.logger.Info("Proxying request", "url", s.proxyURL, "body", string(body))
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to execute proxy request", http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	rdr := io.LimitReader(resp.Body, defaultMaxBodySize)
	n, err := io.Copy(w, rdr)
	if err != nil {
		http.Error(w, "Failed to copy proxy response", http.StatusInternalServerError)
		return
	}
	if n == 0 {
		http.Error(w, "Empty response from proxy", http.StatusInternalServerError)
		return
	}
}
