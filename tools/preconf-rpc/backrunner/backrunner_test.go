package backrunner_test

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/backrunner"
	"github.com/primev/mev-commit/x/util"
)

var txnsResp = []byte(`{"success":true,"data":{"totalRecords":1,"limit":10,"offset":0,"records":[{"chainId":1,"amount":"111444335163840","userAmount":"0","userAddress":"","userPercent":"0","createdAt":"2025-12-08T18:08:09.508Z","revenueType":"Backrun","bundleId":"0xf16f498f05b85cc93d6f498f05b8","bundleHashes":["0xd92eadf1fc432cbfc8db9b06d1a809e3a826666e66052daf28d29ea0417e6965","0xa3d8155e77cc46237e007e7a1274ca277209c47f27bae4405c74f01bb14673ec"],"signalTxHash":""}],"lastRecordTime":"2025-12-08T18:08:09.508Z"}}`)

type swapInfo struct {
	txnHash     string
	blockNumber int64
	builders    []string
}

type mockStore struct {
	mtx      sync.Mutex
	swapInfo map[common.Hash]swapInfo
	rewards  map[common.Hash]*big.Int
}

func (m *mockStore) AddSwapInfo(ctx context.Context, txnHash common.Hash, blockNumber int64, builders []string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.swapInfo[txnHash] = swapInfo{
		txnHash:     txnHash.Hex(),
		blockNumber: blockNumber,
		builders:    builders,
	}
	return nil
}

func (m *mockStore) GetStartHintForRewards(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockStore) UpdateSwapReward(ctx context.Context, reward *big.Int, bundle []string) (bool, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, b := range bundle {
		bundleHash := common.HexToHash(b)
		if _, exists := m.swapInfo[bundleHash]; !exists {
			continue
		}
		m.rewards[bundleHash] = reward
		delete(m.swapInfo, bundleHash)
		return true, nil
	}
	return false, nil
}

func (m *mockStore) GetReward(bundleHash common.Hash) (*big.Int, bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	reward, exists := m.rewards[bundleHash]
	return reward, exists
}

func (m *mockStore) GetSwapInfo(bundleHash common.Hash) (swapInfo, bool) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	info, exists := m.swapInfo[bundleHash]
	return info, exists
}

func (m *mockStore) GetSwapRewardee(ctx context.Context, bundle []string) (common.Address, common.Hash, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for _, b := range bundle {
		bundleHash := common.HexToHash(b)
		if _, exists := m.rewards[bundleHash]; exists {
			return common.HexToAddress("0xRewardeeAddress"), bundleHash, nil
		}
	}
	return common.Address{}, common.Hash{}, errors.New("bundle not found")
}

type pointsEntry struct {
	userID          common.Address
	transactionHash common.Hash
	mevRevenue      *big.Int
}

type mockPointsTracker struct {
	mu      sync.Mutex
	entries []pointsEntry
}

func (m *mockPointsTracker) AssignPoints(ctx context.Context, userID common.Address, transactionHash common.Hash, mevRevenue *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, pointsEntry{
		userID:          userID,
		transactionHash: transactionHash,
		mevRevenue:      mevRevenue,
	})
	return nil
}

func (m *mockPointsTracker) GetEntries() []pointsEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.entries
}

func TestBackrun(t *testing.T) {
	waitForResp := make(chan struct{})
	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/transactions":
				<-waitForResp
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(txnsResp)
			case "/rpc":
				var req map[string]interface{}
				defer func() { _ = r.Body.Close() }()
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				if req["method"] != "eth_sendBundle" {
					http.Error(w, "method not supported", http.StatusBadRequest)
					return
				}
				if req["jsonrpc"] != "2.0" {
					http.Error(w, "bad version", http.StatusBadRequest)
					return
				}
				if params, ok := req["params"].([]any); !ok || len(params) != 1 {
					http.Error(w, "bad params", http.StatusBadRequest)
					return
				} else {
					if p, ok := params[0].(map[string]any); !ok {
						http.Error(w, "bad param type", http.StatusBadRequest)
					} else {
						if _, ok := p["txs"].([]any); !ok {
							http.Error(w, "bad bundle param", http.StatusBadRequest)
							return
						}
						if _, ok := p["blockNumber"].(string); !ok {
							http.Error(w, "bad blockNumber param", http.StatusBadRequest)
							return
						}
						if _, ok := p["trustedBuilders"].([]any); !ok {
							http.Error(w, "bad trustedBuilders param", http.StatusBadRequest)
							return
						}
					}
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0xbundlehash"}`))
			default:
				http.NotFound(w, r)
			}
		}),
	)

	defer srv.Close()

	st := &mockStore{
		swapInfo: make(map[common.Hash]swapInfo),
		rewards:  make(map[common.Hash]*big.Int),
	}

	commitments := []*bidderapiv1.Commitment{
		{
			BlockNumber:     12345678,
			TxHashes:        []string{"0xa3d8155e77cc46237e007e7a1274ca277209c47f27bae4405c74f01bb14673ec"},
			ProviderAddress: "0x2445e5e28890De3e93F39fCA817639c470F4d3b9",
		},
		{
			BlockNumber:     12345678,
			TxHashes:        []string{"0xa3d8155e77cc46237e007e7a1274ca277209c47f27bae4405c74f01bb14673ec"},
			ProviderAddress: "0xB3998135372F1eE16Cb510af70ed212b5155Af62",
		},
	}

	pts := &mockPointsTracker{}

	runner, err := backrunner.New(
		"apiKey",
		srv.URL,
		srv.URL+"/rpc",
		st,
		pts,
		util.NewTestLogger(os.Stdout),
	)
	if err != nil {
		t.Fatalf("failed to create backrunner: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := runner.Start(ctx)

	err = runner.Backrun(ctx, "0xdeadbeef", commitments)
	if err != nil {
		t.Fatalf("failed to backrun: %v", err)
	}

	sInfo, exists := st.GetSwapInfo(common.HexToHash("0xa3d8155e77cc46237e007e7a1274ca277209c47f27bae4405c74f01bb14673ec"))
	if !exists {
		t.Fatalf("swap info not found in store")
	}

	if sInfo.blockNumber != 12345678 {
		t.Fatalf("unexpected block number: got %v, want %v", sInfo.blockNumber, 12345678)
	}
	if len(sInfo.builders) != 2 {
		t.Fatalf("unexpected builders length: got %v, want %v", len(sInfo.builders), 2)
	}

	close(waitForResp)

	for {
		if reward, exists := st.GetReward(common.HexToHash("0xa3d8155e77cc46237e007e7a1274ca277209c47f27bae4405c74f01bb14673ec")); exists {
			expectedReward := big.NewInt(111444335163840)
			expectedReward = new(big.Int).Div(new(big.Int).Mul(expectedReward, big.NewInt(90)), big.NewInt(100)) // 90% to user
			if reward.Cmp(expectedReward) != 0 {
				t.Fatalf("unexpected reward: got %v, want %v", reward, expectedReward)
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(pts.GetEntries()) != 1 {
		t.Fatalf("unexpected points entries length: got %v, want %v", len(pts.entries), 1)
	}

	cancel()
	<-done
}
