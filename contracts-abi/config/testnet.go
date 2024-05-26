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
	BidderRegistry:         "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
	ProviderRegistry:       "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707",
	PreconfCommitmentStore: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
	Oracle:                 "0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e",
	BlockTracker:           "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
}
