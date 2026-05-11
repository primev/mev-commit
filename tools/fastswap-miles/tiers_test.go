package main

import (
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestTierString(t *testing.T) {
	cases := []struct {
		tier Tier
		want string
	}{
		{TierStable, "stable"},
		{TierBlueChip, "bluechip"},
		{TierVolatile, "volatile"},
		{Tier(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.tier.String(); got != c.want {
			t.Errorf("Tier(%d).String() = %q, want %q", c.tier, got, c.want)
		}
	}
}

func TestTierGasCapPercentile(t *testing.T) {
	cases := []struct {
		tier Tier
		want int
	}{
		{TierStable, 25},
		{TierBlueChip, 50},
		{TierVolatile, 75},
	}
	for _, c := range cases {
		if got := c.tier.GasCapPercentile(); got != c.want {
			t.Errorf("Tier(%v).GasCapPercentile() = %d, want %d", c.tier, got, c.want)
		}
	}
}

func TestLookupTokenConfig_KnownStable(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr)
	if cfg.Tier != TierStable {
		t.Errorf("USDC tier = %v, want stable", cfg.Tier)
	}
	if cfg.SweepCadence != 24*time.Hour {
		t.Errorf("USDC cadence = %v, want 24h", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_DAICadence(t *testing.T) {
	// DAI cadence was tightened from 5d to 48h after backtesting.
	cfg := lookupTokenConfig(daiAddr)
	if cfg.SweepCadence != 48*time.Hour {
		t.Errorf("DAI cadence = %v, want 48h", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_KnownBlueChip(t *testing.T) {
	cfg := lookupTokenConfig(wbtcAddr)
	if cfg.Tier != TierBlueChip {
		t.Errorf("WBTC tier = %v, want bluechip", cfg.Tier)
	}
	if cfg.SweepCadence != 24*time.Hour {
		t.Errorf("WBTC cadence = %v, want 24h", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_KnownVolatileNoCadence(t *testing.T) {
	cfg := lookupTokenConfig(pepeAddr)
	if cfg.Tier != TierVolatile {
		t.Errorf("PEPE tier = %v, want volatile", cfg.Tier)
	}
	if cfg.SweepCadence != 0 {
		t.Errorf("PEPE cadence = %v, want zero (no cadence floor)", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_UnknownDefaultsVolatileNoCadence(t *testing.T) {
	addr := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	cfg := lookupTokenConfig(addr)
	if cfg != defaultTokenConfig {
		t.Errorf("unknown addr config = %+v, want defaultTokenConfig %+v", cfg, defaultTokenConfig)
	}
	if cfg.Tier != TierVolatile {
		t.Errorf("default tier = %v, want volatile", cfg.Tier)
	}
	if cfg.SweepCadence != 0 {
		t.Errorf("default cadence = %v, want zero (no cadence floor)", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_CaseInsensitive(t *testing.T) {
	upper := common.HexToAddress(strings.ToUpper(usdcAddr.Hex()))
	lower := common.HexToAddress(strings.ToLower(usdcAddr.Hex()))
	if lookupTokenConfig(upper).Tier != TierStable {
		t.Errorf("upper-case USDC not recognized")
	}
	if lookupTokenConfig(lower).Tier != TierStable {
		t.Errorf("lower-case USDC not recognized")
	}
}

func TestForceSweepInterval(t *testing.T) {
	// Stable / bluechip force-sweep at 1.5×cadence.
	usdc := lookupTokenConfig(usdcAddr)
	if got := usdc.forceSweepInterval(); got != 36*time.Hour {
		t.Errorf("USDC force-sweep interval = %v, want 36h", got)
	}
	usdt := lookupTokenConfig(usdtAddr)
	if got := usdt.forceSweepInterval(); got != 72*time.Hour {
		t.Errorf("USDT force-sweep interval = %v, want 72h", got)
	}

	// Volatile force-sweep at fixed 6h regardless of cadence.
	pepe := lookupTokenConfig(pepeAddr)
	if got := pepe.forceSweepInterval(); got != 6*time.Hour {
		t.Errorf("PEPE force-sweep interval = %v, want 6h", got)
	}
}

func TestGasCapLookback(t *testing.T) {
	usdc := lookupTokenConfig(usdcAddr)
	if got := usdc.gasCapLookback(); got != 24*time.Hour {
		t.Errorf("USDC gas-cap lookback = %v, want 24h", got)
	}
	dai := lookupTokenConfig(daiAddr)
	if got := dai.gasCapLookback(); got != 48*time.Hour {
		t.Errorf("DAI gas-cap lookback = %v, want 48h", got)
	}
	pepe := lookupTokenConfig(pepeAddr)
	if got := pepe.gasCapLookback(); got != 6*time.Hour {
		t.Errorf("PEPE gas-cap lookback = %v, want 6h (volatile fallback)", got)
	}
}

func TestAllStableConfigsConsistent(t *testing.T) {
	stables := []common.Address{usdcAddr, usdtAddr, daiAddr}
	for _, a := range stables {
		cfg := lookupTokenConfig(a)
		if cfg.Tier != TierStable {
			t.Errorf("%s tier = %v, want stable", a.Hex(), cfg.Tier)
		}
		if cfg.SweepCadence == 0 {
			t.Errorf("%s cadence is zero; stable tokens require a cadence floor", a.Hex())
		}
	}
}

func TestIsWhitelisted(t *testing.T) {
	whitelistedTokens := []common.Address{
		usdcAddr, usdtAddr, daiAddr,
		wbtcAddr, arbAddr, linkAddr, compAddr, uniAddr, sushiAddr, inchAddr, yfiAddr,
		pepeAddr,
	}
	for _, a := range whitelistedTokens {
		if !isWhitelisted(a) {
			t.Errorf("%s should be whitelisted", a.Hex())
		}
	}

	// Random unknown token must not be whitelisted (otherwise the
	// attacker-token defense is broken).
	random := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if isWhitelisted(random) {
		t.Errorf("unknown token %s incorrectly whitelisted", random.Hex())
	}

	// Case-insensitive lookup.
	upper := common.HexToAddress(strings.ToUpper(usdcAddr.Hex()))
	if !isWhitelisted(upper) {
		t.Errorf("upper-case USDC not recognized as whitelisted")
	}
}

func TestAllBlueChipConfigsConsistent(t *testing.T) {
	blueChips := []common.Address{wbtcAddr, arbAddr, linkAddr, compAddr, uniAddr, sushiAddr, inchAddr, yfiAddr}
	for _, a := range blueChips {
		cfg := lookupTokenConfig(a)
		if cfg.Tier != TierBlueChip {
			t.Errorf("%s tier = %v, want bluechip", a.Hex(), cfg.Tier)
		}
		if cfg.SweepCadence != 24*time.Hour {
			t.Errorf("%s cadence = %v, want 24h", a.Hex(), cfg.SweepCadence)
		}
	}
}
