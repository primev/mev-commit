package mockrelay

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MockRelay implements RelayQuerier interface
type MockRelay struct {
	client     *ethclient.Client
	blsKey     string
	blsKeyLock sync.RWMutex
}

// NewMockRelay creates a new mock relay with the given L1 RPC URL and BLS key
func NewMockRelay(l1RpcURL string, initialBLSKey string) (*MockRelay, error) {
	client, err := ethclient.Dial(l1RpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to L1: %w", err)
	}

	return &MockRelay{
		client: client,
		blsKey: initialBLSKey,
	}, nil
}

// GetBlockHashForNumber fetches the block hash for a given block number from L1
func (m *MockRelay) GetBlockHashForNumber(ctx context.Context, blockNumber uint64) (common.Hash, error) {
	block, err := m.client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get block: %w", err)
	}
	return block.Hash(), nil
}

// UpdateBLSKey updates the BLS key used by the relay
func (m *MockRelay) UpdateBLSKey(newKey string) error {
	// Validate BLS key format (48 bytes hex)
	newKey = common.TrimHexPrefix(newKey)
	keyBytes, err := hex.DecodeString(newKey)
	if err != nil {
		return fmt.Errorf("invalid BLS key format: %w", err)
	}
	if len(keyBytes) != 48 {
		return fmt.Errorf("BLS key must be 48 bytes")
	}

	m.blsKeyLock.Lock()
	defer m.blsKeyLock.Unlock()
	m.blsKey = newKey
	return nil
}

// GetBLSKey returns the current BLS key
func (m *MockRelay) GetBLSKey() string {
	m.blsKeyLock.RLock()
	defer m.blsKeyLock.RUnlock()
	return m.blsKey
}

// StartServer starts the HTTP server for the mock relay
func (m *MockRelay) StartServer(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/relay/v1/data/bidtraces/proposer_payload_delivered", m.HandleBidTraceRequest)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}

// HandleBidTraceRequest handles the /relay/v1/data/bidtraces/proposer_payload_delivered endpoint
func (m *MockRelay) HandleBidTraceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blockNumberStr := r.URL.Query().Get("block_number")
	if blockNumberStr == "" {
		http.Error(w, "Missing block_number parameter", http.StatusBadRequest)
		return
	}

	blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid block_number", http.StatusBadRequest)
		return
	}

	blockHash, err := m.GetBlockHashForNumber(r.Context(), blockNumber)
	if err != nil {
		http.Error(w, "Failed to get block hash", http.StatusInternalServerError)
		return
	}

	response := []map[string]interface{}{
		{
			"block_number":   blockNumberStr,
			"block_hash":     blockHash.Hex(),
			"builder_pubkey": "0x" + m.GetBLSKey(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
