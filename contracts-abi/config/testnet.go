package config

type Contracts struct {
	BidderRegistry   string
	ProviderRegistry string
	PreconfManager   string
	Oracle           string
	BlockTracker     string
}

type L1Contracts struct {
	ValidatorOptInRouter string
	VanillaRegistry      string
	MevCommitAVS         string
}

var TestnetContracts = Contracts{
	// If these addresses change for a testnet deployment,
	// please also update snippets/testnet-addresses.mdx
	// in https://github.com/primev/mev-commit-docs
	BidderRegistry:   "0x7ffa86fF89489Bca72Fec2a978e33f9870B2Bd25",
	ProviderRegistry: "0x4FC9b98e1A0Ff10de4c2cf294656854F1d5B207D",
	PreconfManager:   "0xCAC68D97a56b19204Dd3dbDC103CB24D47A825A3",
	Oracle:           "0x6856Eb630C79D491886E104D328834643B3F69E3",
	BlockTracker:     "0x2eEbF31f5c932D51556E70235FB98bB2237d065c",
}

var HoleskyContracts = L1Contracts{
	ValidatorOptInRouter: "0xCae46e1013D33587180Db5933Abd75D977c2d7ab",
	VanillaRegistry:      "0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf",
	MevCommitAVS:         "0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4",
}
