package config

import (
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
