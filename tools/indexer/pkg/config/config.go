package config

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
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

func ResolveRelays(c *cli.Context) ([]Relay, error) {
	s := strings.TrimSpace(c.String("relays-json"))
	if s == "" {
		return RelaysDefault, nil
	}
	var v []Relay
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, fmt.Errorf("invalid --relays-json: %w", err)
	}
	if len(v) == 0 {
		return nil, fmt.Errorf("--relays-json provided but empty")
	}
	return v, nil
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
	RelayData        bool
	RelaysJSON       string
}
