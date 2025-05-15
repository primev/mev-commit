package rpcserver

import (
	"context"
	"encoding/json"
	"io"
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

type jsonRPCRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      any               `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      any              `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError    `json:"error,omitempty"`
}

type methodHandler func(ctx context.Context, params ...json.RawMessage) (json.RawMessage, error)

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *any   `json:"data,omitempty"`
}

type JSONRPCServer struct {
	rwLock   sync.RWMutex
	methods  map[string]methodHandler
	proxyURL string
}

func NewJSONRPCServer(proxyURL string) *JSONRPCServer {
	return &JSONRPCServer{
		proxyURL: proxyURL,
		methods:  make(map[string]methodHandler),
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
	defer r.Body.Close()

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

	s.rwLock.RLock()
	handler, ok := s.methods[req.Method]
	s.rwLock.RUnlock()
	if !ok {
		s.proxyRequest(w, r)
		return
	}

	resp, err := handler(r.Context(), req.Params...)
	if err != nil {
		s.writeError(w, req.ID, CodeCustomError, err.Error())
		return
	}

	s.writeResponse(w, req.ID, &resp)
}

func (s *JSONRPCServer) writeResponse(w http.ResponseWriter, id any, result *json.RawMessage) {
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

func (s *JSONRPCServer) proxyRequest(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	req, err := http.NewRequest(r.Method, s.proxyURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to execute proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	rdr := io.LimitReader(resp.Body, defaultMaxBodySize)
	respBuf, err := io.ReadAll(rdr)
	if err != nil {
		http.Error(w, "Failed to read proxy response", http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(respBuf); err != nil {
		http.Error(w, "Failed to write proxy response", http.StatusInternalServerError)
		return
	}
}
