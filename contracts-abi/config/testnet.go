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
	BidderRegistry:         "0xd44adA804c53a7eE42145f752daD8fBa4a521D50",
	ProviderRegistry:       "0x1714b4E6b60FE34f0ec39e428150944D80af7E63",
	PreconfCommitmentStore: "0x56f6A527B07Dc9980dE4609F06887cB498903A0D",
	Oracle:                 "0xc5958e569556b54B25DfE4ad6B3CD4690a5db039",
	BlockTracker:           "0xF4C89c9851A613a447470CCe6866923C6e14F041",
}
