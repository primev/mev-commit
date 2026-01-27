package sim

import (
	"github.com/ethereum/go-ethereum/common"
)

// Swap event signatures from rethsim (topic0)
// These are the exact signatures used in rethsim/src/main.rs
var swapEventSignatures = map[common.Hash]string{
	// Uniswap V2: Sync(uint112,uint112) - emitted on every swap
	common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"): "uniswap_v2_swap",
	// Uniswap V3: Swap(address,address,int256,int256,uint160,uint128,int24)
	common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"): "uniswap_v3_swap",
	// Uniswap V4: Swap(PoolId,address,int128,int128,uint160,uint128,int24,uint24)
	common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f"): "uniswap_v4_swap",
	// MetaMask Swap Router
	common.HexToHash("0xbeee1e6e7fe307ddcf84b0a16137a4430ad5e2480fc4f4a8e250ab56ccd7630d"): "metamask_swap_router",
	// Fluid DEX
	common.HexToHash("0xfbce846c23a724e6e61161894819ec46c90a8d3dd96e90e7342c6ef49ffb539c"): "fluid_swap",
	// Curve: TokenExchange
	common.HexToHash("0x56d0661e240dfb199ef196e16e6f42473990366314f0226ac978f7be3cd9ee83"): "curve_finance_swap",
	// Curve: TokenExchange (tricrypto pools)
	common.HexToHash("0x143f1f8e861fbdeddd5b46e844b7d3ac7b86a122f36e8c463859ee6811b1f29c"): "curve_tricrypto_swap",
	// Curve: StableSwap NG
	common.HexToHash("0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140"): "curve_stableswap_ng_swap",
	// Balancer V2: Swap(bytes32,address,address,uint256,uint256)
	common.HexToHash("0x2170c741c41531aec20e7c107c24eecfdd15e69c9bb0a8dd37b1840b9e0b207b"): "balancer_swap",
	// Balancer: LOG_SWAP
	common.HexToHash("0x908fb5ee8f16c6bc9bc3690973819f32a4d4b10188134543c88706e0e1d43378"): "balancer_log_swap",
	// 1inch Aggregation Router V6
	common.HexToHash("0xfec331350fce78ba658e082a71da20ac9f8d798a99b3c79681c8440cbfe77e07"): "oneinch_aggregation_router_v6",
	// SushiSwap: Swap (same signature as Uniswap V2 Swap event)
	common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"): "sushiswap_swap",
	// KyberSwap
	common.HexToHash("0xd6d4f5681c246c9f42c203e287975af1601f8df8035a9251f79aab5c8f09e2f8"): "kyberswap_swap",
	// PancakeSwap
	common.HexToHash("0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"): "pancakeswap_swap",
	// DODO
	common.HexToHash("0xc2c0245e056d5fb095f04cd6373bc770802ebd1e6c918eb78fdef843cdb37b0f"): "dodoswap_swap",
}

// DetectSwapsFromLogs checks if logs contain swap events.
// Returns whether a swap was detected and the list of swap kinds found.
func DetectSwapsFromLogs(logs []TraceLog) (bool, []string) {
	var swapKinds []string
	seen := make(map[string]bool)

	for _, log := range logs {
		if len(log.Topics) > 0 {
			if swapType, ok := swapEventSignatures[log.Topics[0]]; ok {
				if !seen[swapType] {
					swapKinds = append(swapKinds, swapType)
					seen[swapType] = true
				}
			}
		}
	}

	return len(swapKinds) > 0, swapKinds
}
