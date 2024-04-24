package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
	BlockTracker           string
}

var TestnetContracts = Contracts{
	BidderRegistry:         "0xded9029fC3789ED393D62686c0c0f9dfA92aA2f6",
	ProviderRegistry:       "0xFA19327bDBf2632aAB7C77e61DC69DbC872d5AC1",
	PreconfCommitmentStore: "0x1F8989fAd5f0538D794Fd9fa15d50942F305f367",
	Oracle:                 "0x1cB85eC90320Ef25FB4F991E41392f518980e53a",
	BlockTracker:           "0xCB4AA84C916BB891cBF43320e0c97C3d4329Cec7",
}
