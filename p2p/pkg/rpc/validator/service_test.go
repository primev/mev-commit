package validatorapi_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
	validatorapi "github.com/primev/mev-commit/p2p/pkg/rpc/validator"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func TestGetValidators(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":[{"pubkey":"0x1234567890abcdef","slot":"1"},{"pubkey":"0xfedcba0987654321","slot":"2"}]}`)
	}))
	defer mockServer.Close()

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn:    true,
				IsAvsOptedIn:        false,
				IsMiddlewareOptedIn: false,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn:    false,
				IsAvsOptedIn:        false,
				IsMiddlewareOptedIn: false,
			},
		},
	}

	optsGetter := func() (*bind.CallOpts, error) {
		return &bind.CallOpts{}, nil
	}

	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter)

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 123}
	resp, err := service.GetValidators(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	slot1 := resp.Items[1]
	if slot1 == nil || slot1.BLSKey != "0x1234567890abcdef" || !slot1.IsOptedIn {
		t.Fatalf("unexpected slot1: %v", slot1)
	}

	slot2 := resp.Items[2]
	if slot2 == nil || slot2.BLSKey != "0xfedcba0987654321" || slot2.IsOptedIn {
		t.Fatalf("unexpected slot2: %v", slot2)
	}
}

func TestGetValidators_HTTPError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	optsGetter := func() (*bind.CallOpts, error) {
		return &bind.CallOpts{}, nil
	}

	mockValidatorRouter := &MockValidatorRouterContract{}
	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter)

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

	mockValidatorRouter := &MockValidatorRouterContract{
		ExpectedCalls: map[string]interface{}{
			string(hexutil.MustDecode("0x1234567890abcdef")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn:    true,
				IsAvsOptedIn:        false,
				IsMiddlewareOptedIn: false,
			},
			string(hexutil.MustDecode("0xfedcba0987654321")): validatoroptinrouter.IValidatorOptInRouterOptInStatus{
				IsVanillaOptedIn:    false,
				IsAvsOptedIn:        false,
				IsMiddlewareOptedIn: false,
			},
		},
	}

	optsGetter := func() (*bind.CallOpts, error) {
		return &bind.CallOpts{}, nil
	}

	logger := util.NewTestLogger(os.Stdout)
	service := validatorapi.NewService(mockServer.URL, mockValidatorRouter, logger, optsGetter)

	req := &validatorapiv1.GetValidatorsRequest{Epoch: 0}
	resp, err := service.GetValidators(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}

	slot1 := resp.Items[1]
	if slot1 == nil || slot1.BLSKey != "0x1234567890abcdef" || !slot1.IsOptedIn {
		t.Fatalf("unexpected slot1: %v", slot1)
	}

	slot2 := resp.Items[2]
	if slot2 == nil || slot2.BLSKey != "0xfedcba0987654321" || slot2.IsOptedIn {
		t.Fatalf("unexpected slot2: %v", slot2)
	}
}
