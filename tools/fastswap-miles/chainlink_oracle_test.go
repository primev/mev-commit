package main

import (
	"math/big"
	"testing"
)

func TestDeriveEthInputSurplusEth_TypicalETHtoUSDC(t *testing.T) {
	// 1 ETH input → 3000 USDC out (userAmtOut) + 30 USDC surplus (1% of total).
	// USDC has 6 decimals, so:
	//   inputAmt   = 1e18 wei
	//   userAmtOut = 3000 × 1e6 = 3e9
	//   surplus    = 30 × 1e6 = 3e7
	// Expected surplus_eth = surplus × inputAmt / (userAmtOut + surplus)
	//                      = 3e7 × 1e18 / 3.03e9
	//                      ≈ 9.901e15 wei (~0.0099 ETH)
	inputAmt := new(big.Int).Mul(big.NewInt(1), big.NewInt(1_000_000_000_000_000_000))
	userAmtOut := big.NewInt(3_000_000_000) // 3000 USDC
	surplus := big.NewInt(30_000_000)       // 30 USDC

	got := deriveEthInputSurplusEth(inputAmt, userAmtOut, surplus)

	wantApprox := new(big.Int).Mul(big.NewInt(9_900_000_000), big.NewInt(1_000_000)) // 9.9e15
	tolerance := new(big.Int).Mul(big.NewInt(10_000_000), big.NewInt(1_000_000))     // 1e13 (0.001 ETH)

	diff := new(big.Int).Sub(got, wantApprox)
	diff.Abs(diff)
	if diff.Cmp(tolerance) > 0 {
		t.Errorf("surplus_eth = %s, want ~%s ± %s", got, wantApprox, tolerance)
	}
}

func TestDeriveEthInputSurplusEth_AtTheTwoPercentSlippageCap(t *testing.T) {
	// userAmtOut = 98 (uniswap × 0.98), surplus = 2 → exactly the cap shape.
	// 1 ETH bought 100 token total. surplus_eth = 2/100 of inputAmt = 0.02 ETH.
	inputAmt := big.NewInt(1_000_000_000_000_000_000) // 1 ETH
	userAmtOut := big.NewInt(98)
	surplus := big.NewInt(2)

	got := deriveEthInputSurplusEth(inputAmt, userAmtOut, surplus)

	want := big.NewInt(20_000_000_000_000_000) // 0.02 ETH
	if got.Cmp(want) != 0 {
		t.Errorf("surplus_eth = %s, want %s", got, want)
	}
}

func TestDeriveEthInputSurplusEth_ZeroDenominatorReturnsNil(t *testing.T) {
	got := deriveEthInputSurplusEth(big.NewInt(1e18), big.NewInt(0), big.NewInt(0))
	if got != nil {
		t.Errorf("expected nil for zero denominator, got %s", got)
	}
}

func TestDeriveEthInputSurplusEth_NilInputsReturnNil(t *testing.T) {
	if got := deriveEthInputSurplusEth(nil, big.NewInt(1), big.NewInt(1)); got != nil {
		t.Errorf("nil inputAmt: expected nil, got %s", got)
	}
	if got := deriveEthInputSurplusEth(big.NewInt(1), nil, big.NewInt(1)); got != nil {
		t.Errorf("nil userAmtOut: expected nil, got %s", got)
	}
	if got := deriveEthInputSurplusEth(big.NewInt(1), big.NewInt(1), nil); got != nil {
		t.Errorf("nil surplus: expected nil, got %s", got)
	}
}

func TestScaleChainlinkAnswer_USDCFeed18Decimals(t *testing.T) {
	// Realistic USDC/ETH at $3000 ETH:
	//   1 USDC ≈ 0.000333 ETH → answer = 333_333_333_333_333 (3.33e14, 18 decimals)
	//   1 USDC raw = 1e6 (USDC has 6 decimals)
	//   Expected: 1 USDC → 0.000333 ETH = 3.33e14 wei
	answer := big.NewInt(333_333_333_333_333) // ~1/3000 ETH per USDC, 1e18 scale
	surplus := big.NewInt(1_000_000)          // 1 USDC raw

	got := scaleChainlinkAnswer(surplus, answer, 6, 18)

	want := big.NewInt(333_333_333_333_333) // back to the same magnitude — 0.000333 ETH wei
	if got.Cmp(want) != 0 {
		t.Errorf("scaleChainlinkAnswer = %s, want %s", got, want)
	}
}

func TestScaleChainlinkAnswer_WBTCFeed18Decimals(t *testing.T) {
	// WBTC has 8 decimals. Suppose WBTC/ETH = 20 (1 BTC = 20 ETH).
	//   answer = 20 × 1e18
	//   1 WBTC raw = 1e8
	//   Expected: 1 WBTC → 20 ETH = 20e18 wei
	answer := new(big.Int).Mul(big.NewInt(20), big.NewInt(1_000_000_000_000_000_000))
	surplus := big.NewInt(100_000_000) // 1 WBTC raw

	got := scaleChainlinkAnswer(surplus, answer, 8, 18)

	want := new(big.Int).Mul(big.NewInt(20), big.NewInt(1_000_000_000_000_000_000))
	if got.Cmp(want) != 0 {
		t.Errorf("scaleChainlinkAnswer WBTC = %s, want %s", got, want)
	}
}

func TestScaleChainlinkAnswer_NonStandardFeedDecimals(t *testing.T) {
	// Exercises the negative-exponent branch of the scaler.
	//   surplus_raw    = 1e6
	//   answer         = 333
	//   token_decimals = 6
	//   feed_decimals  = 8
	//   exp = token_decimals + feed_decimals - 18 = -4 → multiply by 10^4
	//   result = 1e6 × 333 × 10^4 / 10^0 = 333 × 10^10 = 3.33e12
	answer := big.NewInt(333)
	surplus := big.NewInt(1_000_000)

	got := scaleChainlinkAnswer(surplus, answer, 6, 8)

	want := big.NewInt(3_330_000_000_000)
	if got.Cmp(want) != 0 {
		t.Errorf("scaleChainlinkAnswer non-std feed = %s, want %s", got, want)
	}
}

func TestScaleChainlinkAnswer_HighDecimalsTokenZeroSurplus(t *testing.T) {
	// Zero surplus → zero output regardless of rate.
	got := scaleChainlinkAnswer(big.NewInt(0), big.NewInt(1e18), 18, 18)
	if got.Sign() != 0 {
		t.Errorf("zero surplus should give zero, got %s", got)
	}
}

// Note: the live Chainlink Registry / ERC20 decimals calls require an L1 RPC
// client and are not unit-tested here. The math above plus whitelist gating
// (covered in tiers_test.go) covers everything that doesn't depend on the
// network layer. The Registry round-trip is exercised in deployment by the
// reconciliation monitor (which catches drift between estimates and reality)
// and by the warning-log path inside getChainlinkRate.
