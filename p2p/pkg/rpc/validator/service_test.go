package validatorapi_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	validatorapi "github.com/primev/mev-commit/p2p/pkg/rpc/validator"
	"github.com/primev/mev-commit/x/epoch"
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
	ec := epoch.NewCalculator(
		0,
		epoch.SlotDuration,
		epoch.SlotsPerEpoch,
		0,
	)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)
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
	ec := epoch.NewCalculator(
		0,
		epoch.SlotDuration,
		epoch.SlotsPerEpoch,
		0,
	)

	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)

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
	ec := epoch.NewCalculator(
		0,
		epoch.SlotDuration,
		epoch.SlotsPerEpoch,
		0,
	)

	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)

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
	ec := epoch.NewCalculator(
		0,
		epoch.SlotDuration,
		epoch.SlotsPerEpoch,
		0,
	)

	svc := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)
	svc.Start(context.Background())

	val := svc.GenesisTime()
	expected := time.Unix(1672531200, 0)
	if !val.Equal(expected) {
		t.Errorf("expected genesis time %v, got %v", expected, val)
	}
}

func TestScheduleNotificationForSlot(t *testing.T) {
	mockNotifier := NewMockNotifier()
	notifyOffset := 1 * time.Second
	slotDuration := 12 * time.Second
	genesisTime := time.Now().Add(100*time.Millisecond + notifyOffset - slotDuration)
	ec := epoch.NewCalculator(
		0,
		slotDuration,
		3, // slots per epoch
		0,
	)

	svc := validatorapi.NewService("http://dummy", nil, util.NewTestLogger(os.Stdout), func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }, mockNotifier, notifyOffset, ec)
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
	t.Parallel()

	dutiesJSONEpoch10 := `{"data":[
        {"pubkey":"0x1234567890abcdef","slot":"0"},
		{"pubkey":"0x7777777777777777","slot":"1"},
        {"pubkey":"0xfedcba0987654321","slot":"2"}
    ]}`
	dutiesJSONEpoch11 := `{"data":[
		{"pubkey":"0x8888888888888888","slot":"3"},
		{"pubkey":"0x9999999999999999","slot":"4"}
    ]}`

	currentEpoch := uint64(0)

	mux := http.NewServeMux()
	// Handle proposer duties.
	mux.HandleFunc("/eth/v1/validator/duties/proposer/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path
		epochStr := path[strings.LastIndex(path, "/")+1:]
		epoch, err := strconv.Atoi(epochStr)
		if err != nil {
			http.Error(w, "Invalid epoch", http.StatusBadRequest)
			return
		}

		if epoch == int(currentEpoch) {
			fmt.Fprint(w, dutiesJSONEpoch10)
		} else if epoch == int(currentEpoch+1) {
			fmt.Fprint(w, dutiesJSONEpoch11)
		} else {
			t.Fatalf("unexpected epoch request: %d", epoch)
		}
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0x7777777777777777")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: false,
			},
			string(hexutil.MustDecode("0x8888888888888888")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
			string(hexutil.MustDecode("0x9999999999999999")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn: true,
			},
		},
	}

	optsGetter := func() (*bind.CallOpts, error) { return &bind.CallOpts{}, nil }
	mockNotifier := NewMockNotifier()
	logger := util.NewTestLogger(os.Stdout)
	ctx := context.Background()

	notifyOffset := 1 * time.Second
	slotDuration := 2 * time.Second
	ec := epoch.NewCalculator(
		0,
		slotDuration,
		3, // slots per epoch
		0,
	)

	svc := validatorapi.NewService(ts.URL, mockValidatorRouter, logger, optsGetter, mockNotifier, notifyOffset, ec)

	svc.SetGenesisTime(time.Now().Add(1*time.Second + 500*time.Millisecond - slotDuration))
	svc.SetProcessEpoch(ctx, currentEpoch)

	numValidatorOptedInNotifs := 0
	numEpochNotifs := 0

	timeout := time.After(2*time.Second + 4*slotDuration)

	// expect 2 validator opted in notifs (second slot in epoch 0 and first slot in epoch 1)
	// also expect 1 epoch notif
	for numValidatorOptedInNotifs < 2 || numEpochNotifs < 1 {
		select {
		case n := <-mockNotifier.NotifyCh:
			if n == nil {
				t.Fatal("expected notification, got nil")
			}
			if n.Topic() == notifications.TopicValidatorOptedIn {
				numValidatorOptedInNotifs++
				slotVal, ok := n.Value()["slot"]
				if !ok {
					t.Fatal("expected slot in notification value")
				}
				if slotVal != uint64(1) && slotVal != uint64(3) {
					t.Errorf("expected slot 1 or 3, got %v", slotVal)
				}
			} else if n.Topic() == notifications.TopicEpochValidatorsOptedIn {
				numEpochNotifs++
				slotsVal, ok := n.Value()["slots"]
				if !ok {
					t.Fatal("expected slots in notification value")
				}
				slotsSlice, ok := slotsVal.([]any)
				if !ok {
					t.Fatalf("expected slots to be a slice of maps, got %T: %v", slotsVal, slotsVal)
				}
				if len(slotsSlice) != 2 {
					t.Fatalf("expected 2 slots, got %d", len(slotsSlice))
				}
				slot1Found, slot2Found := false, false
				for _, slot := range slotsSlice {
					slotData, ok := slot.(map[string]interface{})
					if !ok {
						t.Fatalf("expected slot to be a map, got %T: %v", slot, slot)
					}
					slotVal, ok := slotData["slot"]
					if !ok {
						t.Fatal("expected slot in notification value")
					}
					if slotVal == uint64(0) {
						slot1Found = true
					}
					if slotVal == uint64(1) {
						slot2Found = true
					}
				}
				if !slot1Found || !slot2Found {
					t.Fatalf("expected slots 1 and 2 to be found in epoch notification value")
				}
			} else {
				t.Fatalf("unexpected notification topic: %s", n.Topic())
			}
		case <-timeout:
			t.Fatalf("timeout waiting for notifications: got %d validator opted in (expected 3) and %d epoch validators opted in (expected 1)",
				numValidatorOptedInNotifs, numEpochNotifs)
		}
	}
}

func TestStart(t *testing.T) {
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
	notifyOffset := 10 * time.Millisecond
	ec := epoch.NewCalculator(
		0,
		100*time.Millisecond, // slot duration
		4,                    // slots per epoch
		0,
	)

	svc := validatorapi.NewService(ts.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)

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

func TestStart_NoDuplicateNotifications(t *testing.T) {
	mux := http.NewServeMux()
	// Genesis endpoint returns current time.
	mux.HandleFunc("/eth/v1/beacon/genesis", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":{"genesis_time":"%d"}}`, now)
	})
	// Proposer duties always return a duty for slot "1" with a fixed pubkey.
	mux.HandleFunc("/eth/v1/validator/duties/proposer/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"}]}`)
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
	notifyOffset := 10 * time.Millisecond
	ec := epoch.NewCalculator(
		0,
		100*time.Millisecond, // slot duration
		3,                    // slots per epoch
		0,
	)

	svc := validatorapi.NewService(ts.URL, mockValidatorRouter, logger, optsGetter, notifier, notifyOffset, ec)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	doneChan := svc.Start(ctx)

	// Allow enough time so that at least two tick cycles occur.
	time.Sleep(1 * time.Second)

	// Stop the service.
	cancel()
	<-doneChan

	// Now inspect the notifications to ensure no duplicates for a given epoch.
	// Build maps that count how many notifications per epoch are sent for each topic.
	epochSummaryCount := make(map[uint64]int)
	epochSlotCount := make(map[uint64]int)
	notifier.mu.Lock()
	for _, n := range notifier.Notifications {
		if n.Topic() == notifications.TopicEpochValidatorsOptedIn {
			if epochVal, ok := n.Value()["epoch"]; ok {
				if ep, ok2 := epochVal.(uint64); ok2 {
					epochSummaryCount[ep]++
				}
			}
		}
		if n.Topic() == notifications.TopicValidatorOptedIn {
			if epochVal, ok := n.Value()["epoch"]; ok {
				if ep, ok2 := epochVal.(uint64); ok2 {
					epochSlotCount[ep]++
				}
			}
		}
	}
	notifier.mu.Unlock()

	// Assert that for each epoch, we only see one summary notification.
	for ep, count := range epochSummaryCount {
		if count > 1 {
			t.Errorf("duplicate epoch summary notifications for epoch %d: got %d", ep, count)
		}
	}
	// And assert that for each epoch, we only see one slot notification (per slot).
	for ep, count := range epochSlotCount {
		if count > 1 {
			t.Errorf("duplicate slot notifications for epoch %d: got %d", ep, count)
		}
	}
}
