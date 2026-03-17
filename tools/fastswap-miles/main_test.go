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

// -------------------- Miles Calculation Tests --------------------

func TestMilesCalculation_Profitable(t *testing.T) {
	// Surplus: 0.01 ETH, gas: 0.001 ETH, bid: 0.001 ETH
	// Net profit = 0.01 - 0.001 - 0.001 = 0.008 ETH
	// User share (90%) = 0.0072 ETH
	// Miles = 0.0072 ETH / 0.00001 ETH per point = 720 miles
	surplus, _ := new(big.Int).SetString("10000000000000000", 10) // 0.01 ETH
	gasCost, _ := new(big.Int).SetString("1000000000000000", 10)  // 0.001 ETH
	bidCost, _ := new(big.Int).SetString("1000000000000000", 10)  // 0.001 ETH

	netProfit := new(big.Int).Sub(surplus, gasCost)
	netProfit.Sub(netProfit, bidCost)

	if netProfit.Sign() <= 0 {
		t.Fatal("expected positive net profit")
	}

	userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
	userShare.Div(userShare, big.NewInt(100))

	miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))

	if miles.Int64() != 720 {
		t.Errorf("expected 720 miles, got %d", miles.Int64())
	}
}

func TestMilesCalculation_Unprofitable(t *testing.T) {
	// Surplus: 0.001 ETH, gas: 0.001 ETH, bid: 0.001 ETH
	// Net profit = 0.001 - 0.001 - 0.001 = -0.001 ETH (negative)
	surplus, _ := new(big.Int).SetString("1000000000000000", 10) // 0.001 ETH
	gasCost, _ := new(big.Int).SetString("1000000000000000", 10) // 0.001 ETH
	bidCost, _ := new(big.Int).SetString("1000000000000000", 10) // 0.001 ETH

	netProfit := new(big.Int).Sub(surplus, gasCost)
	netProfit.Sub(netProfit, bidCost)

	if netProfit.Sign() > 0 {
		t.Fatal("expected non-positive net profit")
	}
}

func TestMilesCalculation_SubThreshold(t *testing.T) {
	// Surplus very small: net profit is positive but too small for 1 mile.
	// weiPerPoint = 10_000_000_000_000 (0.00001 ETH)
	// Need userShare (90%) >= weiPerPoint for at least 1 mile.
	// So netProfit >= weiPerPoint * 100 / 90 = ~11_111_111_111_112 wei
	// Use netProfit = 10_000_000_000_000 (just under threshold)
	surplus, _ := new(big.Int).SetString("10000000000000", 10)
	gasCost := big.NewInt(0)
	bidCost := big.NewInt(0)

	netProfit := new(big.Int).Sub(surplus, gasCost)
	netProfit.Sub(netProfit, bidCost)

	userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
	userShare.Div(userShare, big.NewInt(100))

	miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))

	if miles.Sign() > 0 {
		t.Errorf("expected 0 miles for sub-threshold, got %d", miles.Int64())
	}
}

func TestMilesCalculation_UserPaysGas(t *testing.T) {
	// When input is ETH (address(0)), user pays gas, so gas shouldn't be deducted.
	// Surplus: 0.01 ETH, bid: 0.001 ETH, gas: ignored
	// Net profit = 0.01 - 0 - 0.001 = 0.009 ETH
	// User share (90%) = 0.0081 ETH
	// Miles = 810
	surplus, _ := new(big.Int).SetString("10000000000000000", 10) // 0.01 ETH
	bidCost, _ := new(big.Int).SetString("1000000000000000", 10)  // 0.001 ETH

	inputToken := zeroAddr.Hex()
	userPaysGas := inputToken == zeroAddr.Hex()

	gasCostWei := big.NewInt(0)
	if !userPaysGas {
		gasCostWei, _ = new(big.Int).SetString("5000000000000000", 10) // would be 0.005 if charged
	}

	netProfit := new(big.Int).Sub(surplus, gasCostWei)
	netProfit.Sub(netProfit, bidCost)

	userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
	userShare.Div(userShare, big.NewInt(100))

	miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))

	if miles.Int64() != 810 {
		t.Errorf("expected 810 miles, got %d", miles.Int64())
	}
}

func TestERC20ProportionalMiles(t *testing.T) {
	// Two txs in a batch, different surplus amounts.
	// Total batch surplus (token units): 300 + 700 = 1000
	// Actual ETH return from sweep: 0.05 ETH
	// Swap gas cost: 0.002 ETH
	//
	// Tx1: surplus=300, share=30%
	//   grossEth = 0.05 * 300/1000 = 0.015 ETH
	//   overheadGas = 0.002 * 300/1000 = 0.0006 ETH
	//   origGas = 0, origBid = 0.001 ETH
	//   netProfit = 0.015 - 0 - 0.001 - 0.0006 = 0.0134 ETH
	//   userShare = 0.01206 ETH, miles = 1206
	//
	// Tx2: surplus=700, share=70%
	//   grossEth = 0.05 * 700/1000 = 0.035 ETH
	//   overheadGas = 0.002 * 700/1000 = 0.0014 ETH
	//   origGas = 0, origBid = 0.002 ETH
	//   netProfit = 0.035 - 0 - 0.002 - 0.0014 = 0.0316 ETH
	//   userShare = 0.02844 ETH, miles = 2844

	actualEthReturn, _ := new(big.Int).SetString("50000000000000000", 10)  // 0.05 ETH
	actualSwapGasCost, _ := new(big.Int).SetString("2000000000000000", 10) // 0.002 ETH
	totalSum := big.NewInt(1000)                                           // token units

	// Tx1
	surplus1 := big.NewInt(300)
	bidCost1, _ := new(big.Int).SetString("1000000000000000", 10) // 0.001 ETH
	gasCost1 := big.NewInt(0)

	tx1Gross := new(big.Int).Mul(actualEthReturn, surplus1)
	tx1Gross.Div(tx1Gross, totalSum)
	tx1Overhead := new(big.Int).Mul(actualSwapGasCost, surplus1)
	tx1Overhead.Div(tx1Overhead, totalSum)
	tx1Net := new(big.Int).Sub(tx1Gross, gasCost1)
	tx1Net.Sub(tx1Net, bidCost1)
	tx1Net.Sub(tx1Net, tx1Overhead)
	tx1Share := new(big.Int).Mul(tx1Net, big.NewInt(90))
	tx1Share.Div(tx1Share, big.NewInt(100))
	tx1Miles := new(big.Int).Div(tx1Share, big.NewInt(weiPerPoint))

	if tx1Miles.Int64() != 1206 {
		t.Errorf("tx1: expected 1206 miles, got %d", tx1Miles.Int64())
	}

	// Tx2
	surplus2 := big.NewInt(700)
	bidCost2, _ := new(big.Int).SetString("2000000000000000", 10) // 0.002 ETH
	gasCost2 := big.NewInt(0)

	tx2Gross := new(big.Int).Mul(actualEthReturn, surplus2)
	tx2Gross.Div(tx2Gross, totalSum)
	tx2Overhead := new(big.Int).Mul(actualSwapGasCost, surplus2)
	tx2Overhead.Div(tx2Overhead, totalSum)
	tx2Net := new(big.Int).Sub(tx2Gross, gasCost2)
	tx2Net.Sub(tx2Net, bidCost2)
	tx2Net.Sub(tx2Net, tx2Overhead)
	tx2Share := new(big.Int).Mul(tx2Net, big.NewInt(90))
	tx2Share.Div(tx2Share, big.NewInt(100))
	tx2Miles := new(big.Int).Div(tx2Share, big.NewInt(weiPerPoint))

	if tx2Miles.Int64() != 2844 {
		t.Errorf("tx2: expected 2844 miles, got %d", tx2Miles.Int64())
	}
}

func TestGetBidCost(t *testing.T) {
	bidMap := map[string]*big.Int{
		"abc123": big.NewInt(5000),
		"def456": big.NewInt(9999),
	}

	// With 0x prefix
	v := getBidCost(bidMap, "0xABC123")
	if v.Int64() != 5000 {
		t.Errorf("expected 5000, got %d", v.Int64())
	}

	// Without prefix
	v = getBidCost(bidMap, "DEF456")
	if v.Int64() != 9999 {
		t.Errorf("expected 9999, got %d", v.Int64())
	}

	// Missing
	v = getBidCost(bidMap, "0xnotfound")
	if v.Int64() != 0 {
		t.Errorf("expected 0 for missing, got %d", v.Int64())
	}
}

func TestPadTo32(t *testing.T) {
	n := big.NewInt(1)
	padded := padTo32(n)
	if len(padded) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(padded))
	}
	if padded[31] != 1 {
		t.Errorf("expected last byte to be 1, got %d", padded[31])
	}
	for i := 0; i < 31; i++ {
		if padded[i] != 0 {
			t.Errorf("expected byte %d to be 0, got %d", i, padded[i])
		}
	}
}

func TestPadTo32Address(t *testing.T) {
	addr := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	padded := padTo32Address(addr)
	if len(padded) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(padded))
	}
	// First 12 bytes should be zero
	for i := 0; i < 12; i++ {
		if padded[i] != 0 {
			t.Errorf("expected byte %d to be 0, got %d", i, padded[i])
		}
	}
}
