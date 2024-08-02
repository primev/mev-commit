package validatorapi_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	validatorapi "github.com/primev/mev-commit/p2p/pkg/rpc/validator"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockValidatorRegistryContract struct {
	ExpectedCalls map[string]interface{}
}

func (m *MockValidatorRegistryContract) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	key := string(valBLSPubKey)
	if result, ok := m.ExpectedCalls[key]; ok {
		return result.(bool), nil
	}
	return false, errors.New("not found")
}

func TestGetEpoch(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"previous_justified":{"epoch":"100"},"current_justified":{"epoch":"101"},"finalized":{"epoch":"100"}}}`)
	}))
	defer mockServer.Close()

	mockValidatorRegistry := &MockValidatorRegistryContract{}
	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRegistry, logger)

	resp, err := service.GetEpoch(context.Background(), &validatorapiv1.EmptyMessage{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.CurrentEpoch != 102 || resp.CurrentJustifiedEpoch != 101 || resp.PreviousJustifiedEpoch != 100 || resp.FinalizedEpoch != 100 {
		t.Fatalf("unexpected epoch response: %v", resp)
	}
}

func TestGetEpoch_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	mockValidatorRegistry := &MockValidatorRegistryContract{}
	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRegistry, logger)

	_, err := service.GetEpoch(context.Background(), &validatorapiv1.EmptyMessage{})
	if err == nil || status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestGetValidators(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"},{"pubkey":"0xfedcba0987654321","slot":"2"}]}`)
	}))
	defer mockServer.Close()

	mockValidatorRegistry := &MockValidatorRegistryContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): true,
			string(hexutil.MustDecode("0xfedcba0987654321")): false,
		},
	}

	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRegistry, logger)

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	resp, err := service.GetValidators(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	slot1 := resp.Items[1]
	if slot1 == nil || slot1.BLSKey != "0x1234567890abcdef" || !slot1.IsActive {
		t.Fatalf("unexpected slot1: %v", slot1)
	}

	slot2 := resp.Items[2]
	if slot2 == nil || slot2.BLSKey != "0xfedcba0987654321" || slot2.IsActive {
		t.Fatalf("unexpected slot2: %v", slot2)
	}
}

func TestGetValidators_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	mockValidatorRegistry := &MockValidatorRegistryContract{}
	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRegistry, logger)

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	_, err := service.GetValidators(context.Background(), req)
	if err == nil || status.Code(err) != codes.Internal {
		t.Fatalf("expected internal error, got %v", err)
	}
}

func TestGetValidators_EpochZero(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/eth/v1/beacon/states/head/finality_checkpoints" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":{"previous_justified":{"epoch":"100"},"current_justified":{"epoch":"101"},"finalized":{"epoch":"100"}}}`)
		} else if r.URL.Path == "/eth/v1/validator/duties/proposer/102" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"},{"pubkey":"0xfedcba0987654321","slot":"2"}]}`)
		}
	}))
	defer mockServer.Close()

	mockValidatorRegistry := &MockValidatorRegistryContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): true,
			string(hexutil.MustDecode("0xfedcba0987654321")): false,
		},
	}

	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRegistry, logger)

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 0}
	resp, err := service.GetValidators(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	slot1 := resp.Items[1]
	if slot1 == nil || slot1.BLSKey != "0x1234567890abcdef" || !slot1.IsActive {
		t.Fatalf("unexpected slot1: %v", slot1)
	}

	slot2 := resp.Items[2]
	if slot2 == nil || slot2.BLSKey != "0xfedcba0987654321" || slot2.IsActive {
		t.Fatalf("unexpected slot2: %v", slot2)
	}
}
