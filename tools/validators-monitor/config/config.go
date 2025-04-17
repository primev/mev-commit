package config

import "log/slog"

type Config struct {
	Logger                 *slog.Logger
	FetchIntervalSec       int      `json:"fetch_interval_sec"`
	TrackMissed            bool     `json:"track_missed"`
	BeaconNodeURL          string   `json:"beacon_node_url"`
	EthereumRPCURL         string   `json:"ethereum_rpc_url"`
	ValidatorOptInContract string   `json:"contract_address"`
	RelayURLs              []string `json:"relay_urls"`
	SlackWebhookURL        string   `json:"slack_webhook_url"`
	DashboardApiUrl        string   `json:"dashboard_api_url"`
}
