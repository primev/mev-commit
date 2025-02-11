package validatorapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- Mocks ---

// MockValidatorRouterContract implements ValidatorRouterContract.
type MockValidatorRouterContract struct {
	ExpectedCalls map[string]interface{}
}

func (m *MockValidatorRouterContract) AreValidatorsOptedIn(_ *bind.CallOpts, valBLSPubKeys [][]byte) ([]validatoroptinrouter.IValidatorOptInRouterOptInStatus, error) {
	results := make([]validatoroptinrouter.IValidatorOptInRouterOptInStatus, len(valBLSPubKeys))
	for i, key := range valBLSPubKeys {
		if m.ExpectedCalls[string(key)] == nil {
			return nil, errors.New("unexpected call")
		}
		results[i] = m.ExpectedCalls[string(key)].(validatoroptinrouter.IValidatorOptInRouterOptInStatus)
	}
	return results, nil
}

// MockNotifier records published notifications.
type MockNotifier struct {
	mu            sync.Mutex
	Notifications []*notifications.Notification
	NotifyCh      chan *notifications.Notification
}

func NewMockNotifier() *MockNotifier {
	return &MockNotifier{
		NotifyCh: make(chan *notifications.Notification, 10),
	}
}

func (m *MockNotifier) Notify(n *notifications.Notification) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Notifications = append(m.Notifications, n)
	m.NotifyCh <- n
}

// --- Override timing globals for tests ---
func init() {
	SlotDuration = 100 * time.Millisecond
	EpochSlots = 4
	EpochDuration = SlotDuration * time.Duration(EpochSlots) // 400ms
	NotifyOffset = 10 * time.Millisecond
	FetchOffset = 20 * time.Millisecond
}

func TestGetValidators(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/beacon/genesis":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":{"genesis_time":"1672531200"}}`)
		case "/eth/v1/validator/duties/proposer/123":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"},{"pubkey":"0xfedcba0987654321","slot":"2"}]}`)
		case "/eth/v1/beacon/states/head/finality_checkpoints":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":{"previous_justified":{"epoch":"100"},"current_justified":{"epoch":"101"},"finalized":{"epoch":"100"}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: false,
			},
		},
	}
	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	notifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	service, err := NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	ctx := context.Background()
	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	resp, err := service.GetValidators(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if item, ok := resp.Items[1]; !ok || item.BLSKey != "0x1234567890abcdef" || !item.IsOptedIn {
		t.Fatalf("unexpected slot 1: %v", item)
	}
	if item, ok := resp.Items[2]; !ok || item.BLSKey != "0xfedcba0987654321" || item.IsOptedIn {
		t.Fatalf("unexpected slot 2: %v", item)
	}
}

func TestGetValidators_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/eth/v1/beacon/genesis" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":{"genesis_time":"1672531200"}}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	mockValidatorRouter := &MockValidatorRouterContract{}
	notifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	service, err := NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier)
	if err != nil {
		t.Fatalf("expected no error during initialization, got %v", err)
	}

	ctx := context.Background()
	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	_, err = service.GetValidators(ctx, req)
	if err == nil || status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestGetValidators_EpochZero(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/beacon/states/head/finality_checkpoints":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":{"previous_justified":{"epoch":"100"},"current_justified":{"epoch":"101"},"finalized":{"epoch":"100"}}}`)
		case "/eth/v1/validator/duties/proposer/102":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"},{"pubkey":"0xfedcba0987654321","slot":"2"}]}`)
		case "/eth/v1/beacon/genesis":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"data":{"genesis_time":"1672531200"}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: false,
			},
		},
	}
	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	notifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	ctx := context.Background()
	service, err := NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 0}
	resp, err := service.GetValidators(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if item, ok := resp.Items[1]; !ok || item.BLSKey != "0x1234567890abcdef" || !item.IsOptedIn {
		t.Fatalf("unexpected slot 1: %v", item)
	}
	if item, ok := resp.Items[2]; !ok || item.BLSKey != "0xfedcba0987654321" || item.IsOptedIn {
		t.Fatalf("unexpected slot 2: %v", item)
	}
}

func TestNewService_FetchGenesisTime(t *testing.T) {
	genesisTimeStr := "1672531200"
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/eth/v1/beacon/genesis" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"data":{"genesis_time":"%s"}}`, genesisTimeStr)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer mockServer.Close()

	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	mockValidatorRouter := &MockValidatorRouterContract{ExpectedCalls: map[string]interface{}{}}
	notifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)

	svc, err := NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	val := svc.genesisTime
	expected := time.Unix(1672531200, 0)
	if !val.Equal(expected) {
		t.Errorf("expected genesis time %v, got %v", expected, val)
	}
}

func TestScheduleNotificationForSlot(t *testing.T) {
	mockNotifier := NewMockNotifier()
	now := time.Now()
	genesisTime := now.Add(100*time.Millisecond + NotifyOffset - SlotDuration)
	svc := &Service{
		apiURL:          "http://dummy",
		validatorRouter: nil,
		logger:          util.NewTestLogger(os.Stdout),
		metrics:         newMetrics(),
		optsGetter:      func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil },
		notifier:        mockNotifier,
		genesisTime:     genesisTime,
	}

	slotInfo := &validatorapiv1.SlotInfo{
		BLSKey:    "0x1234",
		IsOptedIn: true,
	}

	svc.scheduleNotificationForSlot(1, 1, slotInfo)

	select {
	case n := <-mockNotifier.NotifyCh:
		if n == nil {
			t.Fatal("expected notification, got nil")
		}
		val, ok := n.Value()["slot"]
		if !ok {
			t.Fatal("expected slot in notification value")
		}
		if val != uint64(1) {
			t.Errorf("expected slot 1, got %v", val)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for notification")
	}
}

func TestProcessEpoch(t *testing.T) {
	dutiesJSON := `{"data":[
        {"pubkey":"0x1234567890abcdef","slot":"1"},
        {"pubkey":"0xfedcba0987654321","slot":"2"}
    ]}`

	mux := http.NewServeMux()
	// Handle genesis requests.
	mux.HandleFunc("/eth/v1/beacon/genesis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":{"genesis_time":"1672531200"}}`)
	})
	// Handle proposer duties.
	mux.HandleFunc("/eth/v1/validator/duties/proposer/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, dutiesJSON)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: false,
			},
		},
	}

	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	mockNotifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	ctx := context.Background()

	svc, err := NewService(ts.URL, mockValidatorRouter, logger, optsGetter, mockNotifier)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	now := time.Now()
	svc.genesisTime = now.Add(100*time.Millisecond + NotifyOffset - SlotDuration)

	svc.processEpoch(ctx, 10)

	select {
	case n := <-mockNotifier.NotifyCh:
		if n == nil {
			t.Fatal("expected notification, got nil")
		}
		slotVal, ok := n.Value()["slot"]
		if !ok {
			t.Fatal("expected slot in notification value")
		}
		if slotVal != uint64(1) {
			t.Errorf("expected slot 1, got %v", slotVal)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for notification from processEpoch")
	}
}

func TestStart(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/eth/v1/beacon/genesis", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":{"genesis_time":"1672531200"}}`)
	})
	mux.HandleFunc("/eth/v1/validator/duties/proposer/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return a duty with slot "4"
		fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"4"}]}`)
	})
	mux.HandleFunc("/eth/v1/beacon/states/head/finality_checkpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Return values so that current epoch becomes 0.
		fmt.Fprint(w, `{"data":{"previous_justified":{"epoch":"0"},"current_justified":{"epoch":"0"},"finalized":{"epoch":"0"}}}`)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
		},
	}
	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	mockNotifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	svc, err := NewService(ts.URL, mockValidatorRouter, logger, optsGetter, mockNotifier)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// Set genesisTime so that currentEpoch becomes 0 and nextEpoch = 1.
	svc.genesisTime = time.Now().Add(-350 * time.Millisecond)

	// Start the epoch cron job using the new Start method.
	doneChan := svc.Start(ctx)

	// Expect that within 2 seconds, a notification for slot "4" is published.
	select {
	case n := <-mockNotifier.NotifyCh:
		if n == nil {
			t.Fatal("expected notification, got nil")
		}
		slotVal, ok := n.Value()["slot"]
		if !ok || slotVal != uint64(4) {
			t.Errorf("expected slot 4, got %v", n.Value()["slot"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for notification from Start")
	}

	// Cancel context and wait for the cron job to finish.
	cancel()
	<-doneChan
}
