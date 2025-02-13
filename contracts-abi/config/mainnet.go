package config

var MevCommitChainContracts = Contracts{
	BidderRegistry:    "0xC973D09e51A20C9Ab0214c439e4B34Dbac52AD67",
	ProviderRegistry:  "0xb772Add4718E5BD6Fe57Fb486A6f7f008E52167E",
	PreconfManager:    "0x9fF03b7Ca0767f069e7AA811E383752267cc47Ec",
	Oracle:            "0xa1aaCA1e4583dB498D47f3D5901f2B2EB49Bd8f6",
	BlockTracker:      "0x0DA2a367C51f2a34465ACd6AE5d8A48385e9cB03",
	SettlementGateway: "0x138c60599946280e5a2DCc1f553B8f0cC0554E03",
}

var EthereumContracts = L1Contracts{
	ValidatorOptInRouter: "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64",
	VanillaRegistry:      "0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9",
	MevCommitAVS:         "0xBc77233855e3274E1903771675Eb71E602D9DC2e",
	L1Gateway:            "0xDBf24cafF1470a6D08bF2FF2c6875bafC60Cf881",
	MevCommitMiddleware:  "0x21fD239311B050bbeE7F32850d99ADc224761382",
}
