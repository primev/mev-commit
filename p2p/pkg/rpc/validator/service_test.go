package validatorapi_test

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
	validatorapi "github.com/primev/mev-commit/p2p/pkg/rpc/validator"
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
	notifyOffset := 1 * time.Second
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset)
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
	notifyOffset := 1 * time.Second
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset)

	ctx := context.Background()
	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	_, err := service.GetValidators(ctx, req)
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
	notifyOffset := 1 * time.Second
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset)

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
	notifyOffset := 1 * time.Second
	svc := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset)
	svc.Start(context.Background())

	val := svc.GenesisTime()
	expected := time.Unix(1672531200, 0)
	if !val.Equal(expected) {
		t.Errorf("expected genesis time %v, got %v", expected, val)
	}
}

func TestScheduleNotificationForSlot(t *testing.T) {
	mockNotifier := NewMockNotifier()
	now := time.Now()
	genesisTime := now.Add(100*time.Millisecond + validatorapi.NotifyOffset - validatorapi.SlotDuration)
	notifyOffset := 1 * time.Second
	svc := validatorapi.NewService("http://dummy", nil, util.NewTestLogger(os.Stdout), func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }, mockNotifier, notifyOffset)
	svc.SetGenesisTime(genesisTime)

	slotInfo := &validatorapiv1.SlotInfo{
		BLSKey:    "0x1234",
		IsOptedIn: true,
	}

	svc.ScheduleNotificationForSlot(1, 1, slotInfo)

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

	notifyOffset := 1 * time.Second
	svc := validatorapi.NewService(ts.URL, mockValidatorRouter, logger, optsGetter, mockNotifier, notifyOffset)

	now := time.Now()
	svc.SetGenesisTime(now.Add(100*time.Millisecond + validatorapi.NotifyOffset - validatorapi.SlotDuration))

	svc.SetProcessEpoch(ctx, 10, time.Now().Unix())

	select {
	case n := <-mockNotifier.NotifyCh:
		if n == nil {
			t.Fatal("expected notification, got nil")
		}
		if n.Topic() == notifications.TopicValidatorOptedIn {
			slotVal, ok := n.Value()["slot"]
			if !ok {
				t.Fatal("expected slot in notification value")
			}
			if slotVal != uint64(1) {
				t.Errorf("expected slot 1, got %v", slotVal)
			}
		} else if n.Topic() == notifications.TopicEpochValidatorsOptedIn {
			slotsVal, ok := n.Value()["slots"]
			if !ok {
				t.Fatal("expected slots in notification value")
			}
			slotsSlice, ok := slotsVal.([]any)
			if !ok {
				t.Fatalf("expected slots to be a slice of maps, got %T: %v", slotsVal, slotsVal)
			}
			if len(slotsSlice) != 1 {
				t.Fatalf("expected 1 slot, got %d", len(slotsSlice))
			}
			if slotData, ok := slotsSlice[0].(map[string]interface{}); !ok {
				t.Fatalf("expected slot data to be a map, got %T: %v", slotsSlice[0], slotsSlice[0])
			} else {
				if slotVal, ok := slotData["slot"]; !ok || slotVal != uint64(1) {
					t.Fatalf("expected slot 1, got %v", slotVal)
				}
			}
		} else {
			t.Fatalf("unexpected notification topic: %s", n.Topic())
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for notification from processEpoch")
	}
}

func TestStart(t *testing.T) {
	validatorapi.SetTestTimings(100*time.Millisecond, 4, 10*time.Millisecond, 20*time.Millisecond)
	mux := http.NewServeMux()
	// Return the current Unix timestamp (truncated to seconds) as genesis time.
	mux.HandleFunc("/eth/v1/beacon/genesis", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now() // Use current time so that s.genesisTime is nearly "now".
		w.Header().Set("Content-Type", "application/json")
		// Truncated to whole seconds.
		fmt.Fprintf(w, `{"data":{"genesis_time":"%d"}}`, now.Unix())
	})
	// For any proposer duties request, always return a duty with slot "50".
	// (This ensures that the notification is scheduled for
	mux.HandleFunc("/eth/v1/validator/duties/proposer/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"50"}]}`)
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
	notifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	notifyOffset := 1 * time.Second
	svc := validatorapi.NewService(ts.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	doneChan := svc.Start(ctx)

	// Wait for the notification.
	// With a slot of 50 and SlotDuration = 100ms, the notification is scheduled for roughly 5s after s.genesisTime.
	// Hence we wait up to 7 seconds.
	select {
	case n := <-notifier.NotifyCh:
		if n == nil {
			t.Fatal("expected notification, got nil")
		}
		if n.Topic() == notifications.TopicValidatorOptedIn {
			val := n.Value()
			// We expect the duty returned by the proposer duties endpoint (with slot "50") to be used.
			slotVal, ok := val["slot"]
			if !ok {
				t.Fatal("notification missing slot")
			}
			if slotVal != uint64(50) {
				t.Errorf("expected slot 50, got %v", slotVal)
			}
			blsKeyVal, ok := val["bls_key"]
			if !ok {
				t.Fatal("notification missing bls_key")
			}
			if blsKeyVal != "0x1234567890abcdef" {
				t.Errorf("expected bls_key 0x1234567890abcdef, got %v", blsKeyVal)
			}
			epochVal, ok := val["epoch"]
			if !ok {
				t.Fatal("notification missing epoch")
			}
			// We may not know the exact epoch number but it should be > 0.
			if ep, ok := epochVal.(uint64); !ok || ep == 0 {
				t.Errorf("expected epoch > 0, got %v", epochVal)
			}
		} else if n.Topic() == notifications.TopicEpochValidatorsOptedIn {
			slotsVal, ok := n.Value()["slots"]
			if !ok {
				t.Fatal("expected slots in notification value")
			}
			slotsSlice, ok := slotsVal.([]any)
			if !ok {
				t.Fatalf("expected slots to be a slice of maps, got %T: %v", slotsVal, slotsVal)
			}
			if len(slotsSlice) != 1 {
				t.Fatalf("expected 1 slot, got %d", len(slotsSlice))
			}
			if slotData, ok := slotsSlice[0].(map[string]interface{}); !ok {
				t.Fatalf("expected slot data to be a map, got %T: %v", slotsSlice[0], slotsSlice[0])
			} else {
				if slotVal, ok := slotData["slot"]; !ok || slotVal != uint64(50) {
					t.Fatalf("expected slot 50, got %v", slotVal)
				}
			}
		}
	case <-time.After(7 * time.Second):
		t.Fatal("timeout waiting for notification")
	}

	cancel()
	select {
	case <-doneChan:
		// Service stopped gracefully.
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for service to stop")
	}
}
