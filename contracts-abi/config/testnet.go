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
	BidderRegistry:         "0xC383f487A93eD43a0dC85f9f63069ce37767f875",
	ProviderRegistry:       "0x1Ca3e5398228b4E6dA1f8b7f176b2a797b9B33bf",
	PreconfCommitmentStore: "0x2Aa7b59E07A92908570b79f68addf2915e5586e5",
	Oracle:                 "0xAaA1e7B1CE47910FA3CE79FC18c57cC12322991E",
	BlockTracker:           "0xABa711EdEC215972d8Da0C65DDc3E092754e691c",
}
