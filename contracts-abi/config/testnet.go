package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
}

var TestnetContracts = Contracts{
	BidderRegistry:         "0x02CcEcB19c6D7EFe583C8b97022cB4b4C0B65608",
	ProviderRegistry:       "0x070cE6161AD79a3BC7aEa222FdfC6AD171Ca83F3",
	PreconfCommitmentStore: "0x4DfF34f74aE5C48a5050eb54e7cEDAb9DEF03715",
	Oracle:                 "0x6bD961a0c15057983F6b457319187782fc98FcAc",
}
