package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
}

var TestnetContracts = Contracts{
	BidderRegistry:         "0x02CcEcB19c6D7EFe583C8b97022cB4b4C0B65608",
	ProviderRegistry:       "0x48D4521ac3537256042568Bc5BCD77be10c094cc",
	PreconfCommitmentStore: "0x5dBbe982640B3029A8A3422efe9dDD3faAc8A591",
	Oracle:                 "0xD4821CBFAeC31BE139cC2D5B0874460f684B237B",
}
