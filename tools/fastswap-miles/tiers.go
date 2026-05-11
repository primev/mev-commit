package main

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Tier classifies a token by price-stability profile. The tier drives sweep
// cadence, the gas-cap percentile applied during sweep timing, and the
// fallback behavior for the upfront cost estimate.
type Tier int

const (
	TierStable Tier = iota
	TierBlueChip
	TierVolatile
)

func (t Tier) String() string {
	switch t {
	case TierStable:
		return "stable"
	case TierBlueChip:
		return "bluechip"
	case TierVolatile:
		return "volatile"
	default:
		return "unknown"
	}
}

// GasCapPercentile returns the percentile of recent L1 gas observations above
// which sweep attempts are deferred. Stable tokens are most selective (only
// sweep at the cheapest 25% of recent gas); volatile tokens sweep at almost
// any gas (only blocked at the top 25%). The lookback window for the
// percentile computation is the cadence period (or 6h for volatile).
func (t Tier) GasCapPercentile() int {
	switch t {
	case TierStable:
		return 25
	case TierBlueChip:
		return 50
	case TierVolatile:
		return 75
	default:
		return 75
	}
}

// tokenConfig holds the per-token sweep parameters.
type tokenConfig struct {
	Tier Tier

	// SweepCadence is the target maximum interval between sweep attempts.
	// During this window, the sweep loop deliberately waits for the batch to
	// grow. After cadence elapses, the loop attempts to sweep subject to
	// profitability and gas-cap checks. Force-sweep happens at 1.5×cadence.
	//
	// A zero value means "no cadence floor" — try every cycle. Used for
	// volatile tokens where price risk dominates batching benefit; force-sweep
	// for those is a fixed 6h since last attempt (handled in the scheduler).
	SweepCadence time.Duration
}

// L1 mainnet token addresses for the configured set.
var (
	usdcAddr  = common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	usdtAddr  = common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	daiAddr   = common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f")
	wbtcAddr  = common.HexToAddress("0x2260fac5e5542a773aa44fbcfedf7c193bc2c599")
	arbAddr   = common.HexToAddress("0xb50721bcf8d664c30412cfbc6cf7a15145234ad1")
	linkAddr  = common.HexToAddress("0x514910771af9ca656af840dff83e8264ecf986ca")
	compAddr  = common.HexToAddress("0xc00e94cb662c3520282e6f5717214004a7f26888")
	uniAddr   = common.HexToAddress("0x1f9840a85d5af5bf1d1762f925bdaddc4201f984")
	sushiAddr = common.HexToAddress("0x6b3595068778dd592e39a122f4f5a5cf09c90fe2")
	inchAddr  = common.HexToAddress("0x111111111117dc0aa78b770fa6a738034120c302")
	yfiAddr   = common.HexToAddress("0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9")
	pepeAddr  = common.HexToAddress("0x6982508145454ce325ddbe47a25d4ec3d2311933")
)

// tokenConfigs maps known L1 token addresses to their sweep configuration.
// Unknown addresses fall through to defaultTokenConfig.
//
// Cadence values are calibrated against the Apr 13-27 stable-volume window:
//   - USDC (~31 swaps/day): daily cadence yields ~30-row batches.
//   - USDT (~17 swaps/day): 48h yields ~34-row batches.
//   - DAI  (~8 swaps/day):  48h yields ~15-row batches (compromise between
//     the daily-default preference and the value of larger batches).
//   - Blue chips: 24h with smaller batches.
//   - Volatile / unknowns: zero cadence (every cycle eligible) with 6h
//     force-sweep, because price risk dominates batching benefit.
var tokenConfigs = map[common.Address]tokenConfig{
	usdcAddr:  {TierStable, 24 * time.Hour},
	usdtAddr:  {TierStable, 48 * time.Hour},
	daiAddr:   {TierStable, 48 * time.Hour},
	wbtcAddr:  {TierBlueChip, 24 * time.Hour},
	arbAddr:   {TierBlueChip, 24 * time.Hour},
	linkAddr:  {TierBlueChip, 24 * time.Hour},
	compAddr:  {TierBlueChip, 24 * time.Hour},
	uniAddr:   {TierBlueChip, 24 * time.Hour},
	sushiAddr: {TierBlueChip, 24 * time.Hour},
	inchAddr:  {TierBlueChip, 24 * time.Hour},
	yfiAddr:   {TierBlueChip, 24 * time.Hour},
	pepeAddr:  {TierVolatile, 0},
}

// defaultTokenConfig is used when an output token is not in tokenConfigs.
// Conservative defaults: treat as volatile, no cadence floor.
var defaultTokenConfig = tokenConfig{
	Tier:         TierVolatile,
	SweepCadence: 0,
}

// volatileForceSweepInterval is the time since last sweep attempt at which a
// volatile token's gas cap is dropped (profitability still required).
const volatileForceSweepInterval = 6 * time.Hour

// forceSweepInterval returns the time-since-last-sweep at which the gas cap
// is dropped (profitability remains the only check). For stable/bluechip
// tokens this is 1.5× cadence; for volatile tokens it is a fixed 6h.
func (c tokenConfig) forceSweepInterval() time.Duration {
	if c.SweepCadence == 0 {
		return volatileForceSweepInterval
	}
	return c.SweepCadence + c.SweepCadence/2
}

// gasCapLookback returns the lookback window for computing the gas cap
// percentile. For stable/bluechip tokens this matches the cadence period
// (so "cheap relative to gas during the period we'd actually consider
// sweeping"). For volatile tokens it is 6h.
func (c tokenConfig) gasCapLookback() time.Duration {
	if c.SweepCadence == 0 {
		return volatileForceSweepInterval
	}
	return c.SweepCadence
}

// lookupTokenConfig returns the configuration for a token address, case
// insensitive. Unknown addresses get defaultTokenConfig.
func lookupTokenConfig(addr common.Address) tokenConfig {
	if cfg, ok := tokenConfigs[addr]; ok {
		return cfg
	}
	// tokenConfigs is keyed by checksum address from common.HexToAddress; the
	// caller may pass a lowercased value. Normalize to be safe.
	if cfg, ok := tokenConfigs[common.HexToAddress(strings.ToLower(addr.Hex()))]; ok {
		return cfg
	}
	return defaultTokenConfig
}

// isWhitelisted reports whether the token is explicitly listed in
// tokenConfigs (rather than falling through to defaultTokenConfig). Used to
// gate upfront miles awarding: only whitelisted output tokens are eligible,
// because an attacker who mints their own token and controls its on-chain
// liquidity could otherwise extract miles for surplus that has no realizable
// ETH value. Unknown tokens defer to sweep time, where the realized swap
// result is the source of truth (and where attacker tokens never sweep
// because Barter can't quote them and the profitability gate blocks).
func isWhitelisted(addr common.Address) bool {
	if _, ok := tokenConfigs[addr]; ok {
		return true
	}
	if _, ok := tokenConfigs[common.HexToAddress(strings.ToLower(addr.Hex()))]; ok {
		return true
	}
	return false
}
