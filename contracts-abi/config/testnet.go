package config

import "math/big"

var TestnetChainID = big.NewInt(141414)

var TestnetContracts = Contracts{
	BidderRegistry:    "0xf960f272e6Cb64b2BDe29c2174ce6BEA998Aa067",
	ProviderRegistry:  "0xF62064943680487916558743c26D928Fb162cf5d",
	PreconfManager:    "0x911B08cb805E737DE3ea4E2326CD951E9EfCe39A",
	Oracle:            "0xedD0CCe91ac7e39E3Dda071881311858EC6d0085",
	BlockTracker:      "0x535F5204cFc8A52297dFA3CBe572869b6294f88E",
	SettlementGateway: "0xFaF6F0d4bbc7bC33a4b403b274aBb82d0E794202",
}

var HoodiContracts = L1Contracts{
	ValidatorOptInRouter: "0xa380ba6d6083a4Cb2a3B62b0a81Ea8727861c13e",
	VanillaRegistry:      "0x536f0792c5d5ed592e67a9260606c85f59c312f0",
	MevCommitAVS:         "0xdF8649d298ad05f019eE4AdBD6210867B8AB225F",
	MevCommitMiddleware:  "0x8E847EC4a36c8332652aB3b2B7D5c54dE29c7fde",
	L1Gateway:            "0x0b3b6cf113959214e313d6ad37ad56831acb1776",
}
