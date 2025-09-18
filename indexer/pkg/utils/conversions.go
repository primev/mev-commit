// pkg/utils/conversions.go
package utils

func BlockNumberToSlot(blockNumber int64) int64 {
	// Ethereum mainnet merge happened at slot 4700013 (block 15537394)
	const MERGE_BLOCK = 15537394
	const MERGE_SLOT = 4700013

	if blockNumber < MERGE_BLOCK {
		return 0 // Pre-merge blocks don't have valid slots
	}

	// Post-merge: roughly 1 slot per block
	return MERGE_SLOT + (blockNumber - MERGE_BLOCK)
}
