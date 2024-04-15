package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
}

var TestnetContracts = Contracts{
	BidderRegistry:         "0x02CcEcB19c6D7EFe583C8b97022cB4b4C0B65608",
	ProviderRegistry:       "0xF69451b49598F11c63956bAD5E27f55114200753",
	PreconfCommitmentStore: "0x86281283DA6D9e3987A55Aa702140fAB4dC71B27",
	Oracle:                 "0x6a9C99178af0c99D0d99510f5643A34C622d5d94",
}
