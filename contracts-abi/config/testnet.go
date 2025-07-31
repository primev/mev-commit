package config

import "math/big"

var TestnetChainID = big.NewInt(141414)

var TestnetContracts = Contracts{
	BidderRegistry:    "0x948eCD70FaeF6746A30a00F30f8b9fB2659e4062",
	ProviderRegistry:  "0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67",
	PreconfManager:    "0xa254D1A10777e358B0c2e945343664c7309A0D9d",
	Oracle:            "0xCd27C2Dc26d37Bb17686F709Db438D3Dc546437C",
	BlockTracker:      "0x0b3b6Cf113959214E313d6Ad37Ad56831acb1776",
	SettlementGateway: "0xFaF6F0d4bbc7bC33a4b403b274aBb82d0E794202",
}

var HoleskyContracts = L1Contracts{
	ValidatorOptInRouter: "0xa380ba6d6083a4Cb2a3B62b0a81Ea8727861c13e",
	VanillaRegistry:      "0x536f0792c5d5ed592e67a9260606c85f59c312f0",
	MevCommitAVS:         "0xdF8649d298ad05f019eE4AdBD6210867B8AB225F",
	MevCommitMiddleware:  "0x8E847EC4a36c8332652aB3b2B7D5c54dE29c7fde",
	L1Gateway:            "0x567f0f6d4f7A306c9824d5Ffd0E26f39682cDd7c",
}
