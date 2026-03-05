package main

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestWeiToEth(t *testing.T) {
	tests := []struct {
		name     string
		wei      string
		expected float64
	}{
		{"1 ETH", "1000000000000000000", 1.0},
		{"0.5 ETH", "500000000000000000", 0.5},
		{"0 ETH", "0", 0.0},
		{"Tiny fraction", "10000000000", 0.00000001},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, _ := new(big.Int).SetString(tc.wei, 10)
			result := weiToEth(val)
			if result != tc.expected {
				t.Errorf("expected %f, got %f", tc.expected, result)
			}
		})
	}
}

func TestCallBarter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/swap" {
			t.Errorf("expected path /swap, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", r.Header.Get("Authorization"))
		}

		resp := BarterResponse{
			To:        common.HexToAddress("0x123"),
			GasLimit:  "50000",
			Value:     "0",
			Data:      "0xabc",
			MinReturn: "1000",
		}
		resp.Route.OutputAmount = "1050"
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := barterRequest{
		Source:     "src",
		Target:     "dst",
		SellAmount: "100",
	}

	resp, err := callBarter(ctx, ts.Client(), ts.URL, "test-key", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.MinReturn != "1000" {
		t.Errorf("expected MinReturn 1000, got %s", resp.MinReturn)
	}
	if resp.Route.OutputAmount != "1050" {
		t.Errorf("expected Route.OutputAmount 1050, got %s", resp.Route.OutputAmount)
	}
}

func TestSubmitToFuel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer fuel-key" {
			t.Errorf("expected Bearer fuel-key, got %s", r.Header.Get("Authorization"))
		}

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Assert parts of the JSON body
		if req["name"] != "fast-swap-surplus" {
			t.Errorf("expected name=fast-swap-surplus, got %v", req["name"])
		}

		args := req["args"].(map[string]any)
		val := args["value"].(map[string]any)
		if val["amount"] != "150" {
			t.Errorf("expected amount=150, got %v", val["amount"])
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user := common.HexToAddress("0x999")
	txHash := common.HexToHash("0xabc")
	miles := big.NewInt(150)

	err := submitToFuel(ctx, ts.Client(), ts.URL, "fuel-key", user, txHash, miles)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
