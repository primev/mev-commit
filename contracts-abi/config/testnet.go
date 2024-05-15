package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
	BlockTracker           string
}

var TestnetContracts = Contracts{
	BidderRegistry:         "0x1E218818D409E0f00dfeBE8A960F7585d4fDff70",
	ProviderRegistry:       "0x0332388390d9df01cA3d26269f2B1Fc314deD9c0",
	PreconfCommitmentStore: "0x2Aff805aBdF1Fe79AfcF8B3a9B4B45ECcD6b6D6e",
	Oracle:                 "0x77A4FE615de28fdf0bF68D9B9ba773A32b5C7630",
	BlockTracker:           "0x042744D8cF66d8455350D43F9e09CA73b5C0CB94",
}
