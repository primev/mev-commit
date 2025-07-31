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
	MevCommitMiddleware  string
	L1Gateway            string
}

var DefaultsContracts = map[string]Contracts{
	MainnetChainID.String(): {
		PreconfManager:   MevCommitChainContracts.PreconfManager,
		BlockTracker:     MevCommitChainContracts.BlockTracker,
		ProviderRegistry: MevCommitChainContracts.ProviderRegistry,
		BidderRegistry:   MevCommitChainContracts.BidderRegistry,
		Oracle:           MevCommitChainContracts.Oracle,
	},
	TestnetChainID.String(): {
		PreconfManager:   TestnetContracts.PreconfManager,
		BlockTracker:     TestnetContracts.BlockTracker,
		ProviderRegistry: TestnetContracts.ProviderRegistry,
		BidderRegistry:   TestnetContracts.BidderRegistry,
		Oracle:           TestnetContracts.Oracle,
	},
}

var DefaultsL1Contracts = map[string]L1Contracts{
	MainnetChainID.String(): {
		ValidatorOptInRouter: EthereumContracts.ValidatorOptInRouter,
	},
	TestnetChainID.String(): {
		ValidatorOptInRouter: HoodieContracts.ValidatorOptInRouter,
	},
}
