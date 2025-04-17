package config

type Config struct {
	FetchIntervalSec       int      `json:"fetch_interval_sec"`
	TrackMissed            bool     `json:"track_missed"`
	BeaconNodeURL          string   `json:"beacon_node_url"`
	EthereumRPCURL         string   `json:"ethereum_rpc_url"`
	ValidatorOptInContract string   `json:"contract_address"`
	RelayURLs              []string `json:"relay_urls"`
	SlackWebhookURL        string   `json:"slack_webhook_url"`
	DashboardApiUrl        string   `json:"dashboard_api_url"`
	DB                     DBConfig `json:"db"`

	LookbackMonths    int `json:"lookback_months"`
	MaxEpochsPerBatch int `json:"max_epochs_per_batch"`
}

type DBConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}
