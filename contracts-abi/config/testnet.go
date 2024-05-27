package config

type Contracts struct {
	BidderRegistry         string
	ProviderRegistry       string
	PreconfCommitmentStore string
	Oracle                 string
	BlockTracker           string
}

var TestnetContracts = Contracts{
	// If these addresses change for a testnet deployment,
	// please also update snippets/testnet-addresses.mdx
	// in https://github.com/primev/mev-commit-docs
	BidderRegistry:         "0x8021D1a949E124E39F9FFc812F39054e1521dAed",
	ProviderRegistry:       "0x9C108C6F4A6350B42B4Cf0384249A0BB2F57e12D",
	PreconfCommitmentStore: "0xb4FFB66e8a913F6Bcc7FbF4d27BC2b1545A0811b",
	Oracle:                 "0xfC9a86CcC74F63AFE1c1d4A1aC25B2e248b234F6",
	BlockTracker:           "0x2816c159E9F0807a5121e39dA4969de3b6D8d25e",
}
