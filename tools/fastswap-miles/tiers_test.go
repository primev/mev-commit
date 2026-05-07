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

func TestLookupTokenConfig_KnownStable(t *testing.T) {
	cfg := lookupTokenConfig(usdcAddr)
	if cfg.Tier != TierStable {
		t.Errorf("USDC tier = %v, want stable", cfg.Tier)
	}
	if cfg.SweepCadence != 24*time.Hour {
		t.Errorf("USDC cadence = %v, want 24h", cfg.SweepCadence)
	}
	if cfg.CostEstimatePctile != 40 {
		t.Errorf("USDC pctile = %d, want 40", cfg.CostEstimatePctile)
	}
	if cfg.ExpectedBatchSize != 30 {
		t.Errorf("USDC batch = %d, want 30", cfg.ExpectedBatchSize)
	}
}

func TestLookupTokenConfig_KnownBlueChip(t *testing.T) {
	cfg := lookupTokenConfig(wbtcAddr)
	if cfg.Tier != TierBlueChip {
		t.Errorf("WBTC tier = %v, want bluechip", cfg.Tier)
	}
	if cfg.ExpectedBatchSize != 3 {
		t.Errorf("WBTC batch = %d, want 3", cfg.ExpectedBatchSize)
	}
}

func TestLookupTokenConfig_KnownVolatile(t *testing.T) {
	cfg := lookupTokenConfig(pepeAddr)
	if cfg.Tier != TierVolatile {
		t.Errorf("PEPE tier = %v, want volatile", cfg.Tier)
	}
	if cfg.SweepCadence != 6*time.Hour {
		t.Errorf("PEPE cadence = %v, want 6h", cfg.SweepCadence)
	}
}

func TestLookupTokenConfig_UnknownDefaultsVolatile(t *testing.T) {
	addr := common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	cfg := lookupTokenConfig(addr)
	if cfg != defaultTokenConfig {
		t.Errorf("unknown addr config = %+v, want defaultTokenConfig %+v", cfg, defaultTokenConfig)
	}
	if cfg.Tier != TierVolatile {
		t.Errorf("default tier = %v, want volatile", cfg.Tier)
	}
}

func TestLookupTokenConfig_CaseInsensitive(t *testing.T) {
	// Same address, different cases.
	upper := common.HexToAddress(strings.ToUpper(usdcAddr.Hex()))
	lower := common.HexToAddress(strings.ToLower(usdcAddr.Hex()))
	if lookupTokenConfig(upper).Tier != TierStable {
		t.Errorf("upper-case USDC not recognized")
	}
	if lookupTokenConfig(lower).Tier != TierStable {
		t.Errorf("lower-case USDC not recognized")
	}
}

func TestAllStableConfigsConsistent(t *testing.T) {
	stables := []common.Address{usdcAddr, usdtAddr, daiAddr}
	for _, a := range stables {
		cfg := lookupTokenConfig(a)
		if cfg.Tier != TierStable {
			t.Errorf("%s tier = %v, want stable", a.Hex(), cfg.Tier)
		}
		if cfg.CostEstimatePctile != 40 {
			t.Errorf("%s pctile = %d, want 40", a.Hex(), cfg.CostEstimatePctile)
		}
		if cfg.ExpectedBatchSize != 30 {
			t.Errorf("%s batch = %d, want 30", a.Hex(), cfg.ExpectedBatchSize)
		}
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
