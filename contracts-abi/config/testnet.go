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
	BidderRegistry:   "0x401B3287364f95694c43ACA3252831cAc02e5C41",
	ProviderRegistry: "0xf4F10e18244d836311508917A3B04694D88999Dd",
	PreconfManager:   "0x9433bCD9e89F923ce587f7FA7E39e120E93eb84D",
	Oracle:           "0x0a3ad886AEfd3bA877bcB23E171e0e2a375806a0",
	BlockTracker:     "0x7538F3AaA07dA1990486De21A0B438F55e9639e4",
}

var HoleskyContracts = L1Contracts{
	ValidatorOptInRouter: "0x251Fbc993f58cBfDA8Ad7b0278084F915aCE7fc3",
	VanillaRegistry:      "0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf",
	MevCommitAVS:         "0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4",
}
