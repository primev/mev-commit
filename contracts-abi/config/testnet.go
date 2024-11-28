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

var TestnetContracts = Contracts{
	// If these addresses change for a testnet deployment,
	// please also update snippets/testnet-addresses.mdx
	// in https://github.com/primev/mev-commit-docs
	BidderRegistry:    "0x948eCD70FaeF6746A30a00F30f8b9fB2659e4062",
	ProviderRegistry:  "0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67",
	PreconfManager:    "0xa254D1A10777e358B0c2e945343664c7309A0D9d",
	Oracle:            "0xCd27C2Dc26d37Bb17686F709Db438D3Dc546437C",
	BlockTracker:      "0x0b3b6Cf113959214E313d6Ad37Ad56831acb1776",
	SettlementGateway: "0xFaF6F0d4bbc7bC33a4b403b274aBb82d0E794202",
}

var HoleskyContracts = L1Contracts{
	ValidatorOptInRouter: "0x251Fbc993f58cBfDA8Ad7b0278084F915aCE7fc3",
	VanillaRegistry:      "0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf",
	MevCommitAVS:         "0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4",
	L1Gateway:            "0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67",
}
