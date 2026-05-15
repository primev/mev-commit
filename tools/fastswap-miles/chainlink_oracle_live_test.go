//go:build integration

package main

import (
	"context"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Live integration test against Ethereum mainnet. Probes the Chainlink Feed
// Registry contract for every whitelisted token, reports coverage + the
// actual surplus_eth value computed for a synthetic 1-unit surplus.
//
// Run with:
//
//	MAINNET_RPC_URL=https://eth.llamarpc.com go test -tags=integration -v -run TestPriceOracle_Live ./tools/fastswap-miles/...
//
// Skipped in normal test runs because it makes real RPC calls. Build tag
// `integration` keeps it out of `go test ./...` by default.

func TestPriceOracle_Live_ChainlinkCoverage(t *testing.T) {
	rpcURL := os.Getenv("MAINNET_RPC_URL")
	if rpcURL == "" {
		t.Skip("MAINNET_RPC_URL not set; skipping live test")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	weth := common.HexToAddress(defaultWETH)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	o, err := newPriceOracle(client, weth, logger)
	if err != nil {
		t.Fatalf("newPriceOracle: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cases := []struct {
		name     string
		token    common.Address
		decimals int
	}{
		{"USDC", usdcAddr, 6},
		{"USDT", usdtAddr, 6},
		{"DAI", daiAddr, 18},
		{"WBTC", wbtcAddr, 8},
		{"ARB", arbAddr, 18},
		{"LINK", linkAddr, 18},
		{"COMP", compAddr, 18},
		{"UNI", uniAddr, 18},
		{"SUSHI", sushiAddr, 18},
		{"1INCH", inchAddr, 18},
		{"YFI", yfiAddr, 18},
		{"PEPE", pepeAddr, 18},
	}

	// We synthesize an ERC20-input swap to force the Chainlink path
	// (event-derivation only kicks in for ETH/WETH input). The actual
	// inputAmt and userAmtOut don't matter because the Chainlink branch
	// only uses surplus + the token rate.
	inputToken := usdtAddr
	inputAmt := big.NewInt(0)
	userAmtOut := big.NewInt(0)

	t.Log("=== Per-token Chainlink Feed Registry coverage ===")
	covered := 0
	deferred := 0
	rateLimitedRetries := []int{}
	for i, c := range cases {
		// 1 token's worth of surplus in raw units: 10^decimals.
		surplus := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(c.decimals)), nil)

		// Throttle to stay under Infura free-tier limits (~10 rps).
		time.Sleep(300 * time.Millisecond)

		ethValue, eligible, source := o.PriceSurplusEth(ctx, inputToken, c.token, inputAmt, userAmtOut, surplus)
		if !eligible {
			deferred++
			t.Logf("  %-6s %s  DEFERRED  source=%s",
				c.name, c.token.Hex(), source)
			// Schedule retry — we want to differentiate "feed truly missing"
			// (revert) from "rate limited" so the user sees real coverage.
			rateLimitedRetries = append(rateLimitedRetries, i)
			continue
		}
		covered++
		ethFloat := weiToEth(ethValue)
		t.Logf("  %-6s %s  COVERED   source=%s  1 token = %.10f ETH (%s wei)",
			c.name, c.token.Hex(), source, ethFloat, ethValue.String())
	}
	t.Logf("=== First-pass coverage: %d covered, %d deferred ===", covered, deferred)

	if len(rateLimitedRetries) == 0 {
		return
	}

	t.Log("")
	t.Log("=== Retry pass for deferred tokens (slower, fresh oracle to clear cached failures) ===")
	// Build a fresh oracle so cached-failure state (if any) is cleared.
	o2, err := newPriceOracle(client, weth, logger)
	if err != nil {
		t.Fatalf("newPriceOracle (retry): %v", err)
	}
	finalCovered := covered
	finalDeferred := 0
	for _, i := range rateLimitedRetries {
		c := cases[i]
		surplus := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(c.decimals)), nil)

		// Bigger delay on retry to avoid rate-limit noise.
		time.Sleep(2 * time.Second)

		ethValue, eligible, source := o2.PriceSurplusEth(ctx, inputToken, c.token, inputAmt, userAmtOut, surplus)
		if !eligible {
			finalDeferred++
			t.Logf("  %-6s %s  STILL DEFERRED  source=%s (likely real coverage gap)",
				c.name, c.token.Hex(), source)
			continue
		}
		finalCovered++
		ethFloat := weiToEth(ethValue)
		t.Logf("  %-6s %s  RECOVERED       source=%s  1 token = %.10f ETH",
			c.name, c.token.Hex(), source, ethFloat)
	}
	t.Logf("=== Final coverage: %d covered, %d deferred (out of %d) ===",
		finalCovered, finalDeferred, len(cases))
}

// Exercises every realistic swap shape end-to-end through PriceSurplusEth
// against mainnet. Logs the surplus_eth that would be used for miles
// awarding for each case so it can be eyeballed before deploy.
func TestPriceOracle_Live_AllSwapShapes(t *testing.T) {
	rpcURL := os.Getenv("MAINNET_RPC_URL")
	if rpcURL == "" {
		t.Skip("MAINNET_RPC_URL not set; skipping live test")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	weth := common.HexToAddress(defaultWETH)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	o, _ := newPriceOracle(client, weth, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cases := []struct {
		name       string
		inputTok   common.Address
		outputTok  common.Address
		inputAmt   *big.Int
		userAmtOut *big.Int
		surplus    *big.Int
		expect     string // expected source
	}{
		{
			name:       "ETH→USDC (typical, 1% surplus)",
			inputTok:   zeroAddr,
			outputTok:  usdcAddr,
			inputAmt:   big.NewInt(1_000_000_000_000_000_000), // 1 ETH
			userAmtOut: big.NewInt(3_000_000_000),             // 3000 USDC
			surplus:    big.NewInt(30_000_000),                // 30 USDC
			expect:     "event_derived",
		},
		{
			name:       "WETH→WBTC (event-derived since WBTC has no Chainlink)",
			inputTok:   weth,
			outputTok:  wbtcAddr,
			inputAmt:   new(big.Int).Mul(big.NewInt(10), big.NewInt(1_000_000_000_000_000_000)),
			userAmtOut: big.NewInt(40_000_000), // ~0.4 BTC at 8 decimals
			surplus:    big.NewInt(800_000),    // ~0.008 BTC
			expect:     "event_derived",
		},
		{
			name:       "USDC→USDT (Chainlink path)",
			inputTok:   usdcAddr,
			outputTok:  usdtAddr,
			inputAmt:   big.NewInt(3_000_000_000), // 3000 USDC
			userAmtOut: big.NewInt(2_990_000_000), // 2990 USDT (typical)
			surplus:    big.NewInt(10_000_000),    // 10 USDT
			expect:     "chainlink",
		},
		{
			name:       "USDC→WBTC (defers: WBTC not in Registry)",
			inputTok:   usdcAddr,
			outputTok:  wbtcAddr,
			inputAmt:   big.NewInt(50_000_000_000), // 50K USDC
			userAmtOut: big.NewInt(80_000_000),     // 0.8 WBTC
			surplus:    big.NewInt(1_600_000),      // 0.016 WBTC
			expect:     "deferred:no_chainlink",
		},
		{
			name:       "USDC→AttackerToken (defers: not whitelisted)",
			inputTok:   usdcAddr,
			outputTok:  common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			inputAmt:   big.NewInt(1_000_000_000),
			userAmtOut: big.NewInt(100),
			surplus:    big.NewInt(5),
			expect:     "deferred:not_whitelisted",
		},
		{
			name:       "ETH→AttackerToken (defers: not whitelisted, defends attack vector)",
			inputTok:   zeroAddr,
			outputTok:  common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			inputAmt:   big.NewInt(1_000_000_000_000_000_000),
			userAmtOut: big.NewInt(98),
			surplus:    big.NewInt(2),
			expect:     "deferred:not_whitelisted",
		},
	}

	for _, c := range cases {
		// Pacing for shared RPC.
		time.Sleep(400 * time.Millisecond)

		ethValue, eligible, source := o.PriceSurplusEth(
			ctx, c.inputTok, c.outputTok, c.inputAmt, c.userAmtOut, c.surplus,
		)
		if source != c.expect {
			t.Errorf("%s: source = %s, want %s", c.name, source, c.expect)
		}
		if eligible {
			t.Logf("  %s -> %s, surplus_eth = %.10f ETH",
				c.name, source, weiToEth(ethValue))
		} else {
			t.Logf("  %s -> %s (deferred, correct)", c.name, source)
		}
	}
}

// Sanity-check that ETH-leg event derivation needs no RPC at all (it's
// pure event-data math). This complements the unit test by running it
// through PriceSurplusEth with a real client wired up.
func TestPriceOracle_Live_EthLegNoRPCNeeded(t *testing.T) {
	rpcURL := os.Getenv("MAINNET_RPC_URL")
	if rpcURL == "" {
		t.Skip("MAINNET_RPC_URL not set; skipping live test")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("dial RPC: %v", err)
	}
	defer client.Close()

	weth := common.HexToAddress(defaultWETH)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	o, _ := newPriceOracle(client, weth, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1 ETH → 3000 USDC user out + 30 USDC surplus.
	// Expected surplus_eth ≈ 30 / 3030 ≈ 0.0099 ETH
	inputAmt := big.NewInt(1_000_000_000_000_000_000)
	userAmtOut := big.NewInt(3_000_000_000)
	surplus := big.NewInt(30_000_000)

	ethValue, eligible, source := o.PriceSurplusEth(ctx, zeroAddr, usdcAddr, inputAmt, userAmtOut, surplus)
	if !eligible {
		t.Fatalf("ETH→USDC must be eligible; got source=%s", source)
	}
	if source != "event_derived" {
		t.Errorf("source = %s, want event_derived", source)
	}
	t.Logf("ETH→USDC 1 ETH / 3000+30 USDC: surplus_eth = %.6f ETH (%s wei)",
		weiToEth(ethValue), ethValue.String())
}
