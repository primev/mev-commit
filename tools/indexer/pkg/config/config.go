package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

type Relay struct {
	Relay_id int64

	URL string
}

var RelaysDefault = []Relay{
	{Relay_id: 1, URL: "https://regional.titanrelay.xyz"},
	{Relay_id: 2, URL: "https://aestus.live"},
	{Relay_id: 3, URL: "https://bloxroute.max-profit.blxrbdn.com"},
	{Relay_id: 4, URL: "https://bloxroute.regulated.blxrbdn.com"},
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
	AlchemyRPC       string
	BeaconBase       string
	RelayData        bool
	RelaysJSON       string
	RatedAPIKey      string
	QuickNodeBase    string
}
