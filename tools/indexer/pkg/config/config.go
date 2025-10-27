package config

import (
	"time"
)

type Relay struct {
	Relay_id int64
	Name     string
	Tag      string
	URL      string
}

var RelaysDefault = []Relay{
	{Relay_id: 1, Name: "Titan", Tag: "titan-relay", URL: "https://regional.titanrelay.xyz"},
	{Relay_id: 2, Name: "Aestus", Tag: "aestus-relay", URL: "https://aestus.live"},
	{Relay_id: 3, Name: "Bloxroute Max Profit", Tag: "bloxroute-max-profit-relay", URL: "https://bloxroute.max-profit.blxrbdn.com"},
	{Relay_id: 4, Name: "Bloxroute Regulated", Tag: "bloxroute-regulated-relay", URL: "https://bloxroute.regulated.blxrbdn.com"},
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
	RPCURL           string
	BeaconBase       string
	BeaconchaAPIKey  string
	BeaconchaRPS     int
}
