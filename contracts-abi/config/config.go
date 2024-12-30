package config

type Contracts struct {
	BidderRegistry    string
	ProviderRegistry  string
	PreconfManager    string
	Oracle            string
	BlockTracker      string
	SettlementGateway string
}

type L1Contracts struct {
	ValidatorOptInRouter string
	VanillaRegistry      string
	MevCommitAVS         string
	L1Gateway            string
}
