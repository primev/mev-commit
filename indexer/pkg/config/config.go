package config

import (
	"os"
	"strconv"
	"time"
)

type Relay struct {
	Name string
	Tag  string
	URL  string
}

var RelaysDefault = []Relay{
	{Name: "Titan", Tag: "titan-relay", URL: "https://regional.titanrelay.xyz"},
	{Name: "Aestus", Tag: "aestus-relay", URL: "https://aestus.live"},
	{Name: "Bloxroute Max Profit", Tag: "bloxroute-max-profit-relay", URL: "https://bloxroute.max-profit.blxrbdn.com"},
	{Name: "Bloxroute Regulated", Tag: "bloxroute-regulated-relay", URL: "https://bloxroute.regulated.blxrbdn.com"},
}

const beaconBase = "https://beaconcha.in/api/v1"

type Config struct {
	BlockTick        time.Duration
	ValidatorWait    time.Duration
	BackfillEvery    time.Duration
	BackfillLookback int64
	BackfillBatch    int
	MaxRetries       int
	BaseRetryDelay   time.Duration
	HTTPTimeout      time.Duration
	OptInContract    string
	EtherscanKey     string
	InfuraRPC        string
	BeaconBase       string
}

func LoadConfig() *Config {
	return &Config{
		BlockTick:        getenvDuration("BLOCK_INTERVAL", 12*time.Second),
		ValidatorWait:    getenvDuration("VALIDATOR_DELAY", 1500*time.Millisecond),
		BackfillEvery:    getenvDuration("BACKFILL_EVERY", 5*time.Minute),
		BackfillLookback: int64(getenvInt("BACKFILL_LOOKBACK_SLOTS", 512)),
		BackfillBatch:    getenvInt("BACKFILL_BATCH", 50),
		MaxRetries:       getenvInt("MAX_RETRIES", 3),
		BaseRetryDelay:   getenvDuration("BASE_RETRY_DELAY", 1*time.Second),
		HTTPTimeout:      getenvDuration("HTTP_TIMEOUT", 15*time.Second),
		OptInContract:    getenv("OPT_IN_CONTRACT", "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64"),
		EtherscanKey:     os.Getenv("ETHERSCAN_API_KEY"),
		InfuraRPC:        os.Getenv("INFURA_RPC"),
	}
}

var config *Config

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
