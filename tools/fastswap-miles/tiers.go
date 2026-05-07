package main

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Tier classifies a token by price-stability profile. The tier drives sweep
// cadence, the percentile of recent sweep gas used to estimate per-user cost
// at miles-awarding time, and the gas cap above which sweeps are deferred.
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

// tokenConfig holds the per-token sweep parameters. Values are tuned by tier
// and refined per token where the realized data justifies a different default.
type tokenConfig struct {
	Tier               Tier
	SweepCadence       time.Duration // target maximum interval between sweeps
	CostEstimatePctile int           // percentile of recent sweep gas to use as upfront cost estimate
	ExpectedBatchSize  int           // assumed batch size for per-user cost dilution
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

// tokenConfigs maps known L1 token addresses (lowercased hex) to their sweep
// configuration. Lookup is case-insensitive via lookupTokenConfig. Unknown
// tokens fall through to defaultTokenConfig.
//
// Cadence values are derived from the Apr 13-27 stable-volume window:
//   - USDC (~31 swaps/day): daily yields ~30-row batches.
//   - USDT (~17 swaps/day): every 2d yields ~34-row batches.
//   - DAI  (~8 swaps/day, mostly 1-3/day): every 5d yields ~38-row batches.
//   - Blue chips (3-15 swaps/day): daily is reasonable; batch sizes are smaller.
//   - Volatile (low/sporadic volume): 6h to limit price-risk exposure.
//
// CostEstimatePctile and ExpectedBatchSize together set how much "buffer" the
// protocol keeps when paying out miles upfront. Stables use a low percentile
// (p40) and a generous batch-size assumption (30) because their realized cost
// is consistent and batches reliably exceed 30; the difference is protocol
// upside. Volatiles use p75 and batch-size 1 — the worst-case assumption.
var tokenConfigs = map[common.Address]tokenConfig{
	usdcAddr:  {TierStable, 24 * time.Hour, 40, 30},
	usdtAddr:  {TierStable, 48 * time.Hour, 40, 30},
	daiAddr:   {TierStable, 120 * time.Hour, 40, 30},
	wbtcAddr:  {TierBlueChip, 24 * time.Hour, 50, 3},
	arbAddr:   {TierBlueChip, 24 * time.Hour, 50, 3},
	linkAddr:  {TierBlueChip, 24 * time.Hour, 50, 3},
	compAddr:  {TierBlueChip, 24 * time.Hour, 50, 3},
	uniAddr:   {TierBlueChip, 24 * time.Hour, 50, 3},
	sushiAddr: {TierBlueChip, 24 * time.Hour, 50, 3},
	inchAddr:  {TierBlueChip, 24 * time.Hour, 50, 3},
	yfiAddr:   {TierBlueChip, 24 * time.Hour, 50, 3},
	pepeAddr:  {TierVolatile, 6 * time.Hour, 75, 1},
}

// defaultTokenConfig is used when an output token is not in tokenConfigs.
// Conservative defaults: treat as volatile, sweep every 6h, assume size-1
// batches at p75 of recent costs.
var defaultTokenConfig = tokenConfig{
	Tier:               TierVolatile,
	SweepCadence:       6 * time.Hour,
	CostEstimatePctile: 75,
	ExpectedBatchSize:  1,
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
