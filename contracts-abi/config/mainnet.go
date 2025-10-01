package config

import "math/big"

var MainnetChainID = big.NewInt(57173)

var MevCommitChainContracts = Contracts{
	BidderRegistry:    "0x145a9f4cbae2ec281f417195ea3464dbd04289a2",
	ProviderRegistry:  "0xeb6d22309062a86fa194520344530874221ef48c",
	PreconfManager:    "0x2ee9e88f57a7db801e114a4df7a99eb7257871e2",
	Oracle:            "0x37a037d2423221f403cfa146f5fb962e19582d90",
	BlockTracker:      "0x5d64b933739558101f9359e2750acc228f0cb64f",
	SettlementGateway: "0x21f5f1142200a515248a2eef5b0654581c7f2b46",
}

var EthereumContracts = L1Contracts{
	ValidatorOptInRouter: "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64",
	VanillaRegistry:      "0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9",
	MevCommitAVS:         "0xBc77233855e3274E1903771675Eb71E602D9DC2e",
	L1Gateway:            "0x5d64b933739558101f9359e2750acc228f0cb64f",
	MevCommitMiddleware:  "0x21fD239311B050bbeE7F32850d99ADc224761382",
}
