package sim

import (
	"github.com/ethereum/go-ethereum/common"
)

// Swap event signatures (topic0) used to detect DEX trades.
// These match the signatures used in rethsim.
var swapEventSignatures = map[common.Hash]string{
	// Uniswap V2 / SushiSwap Swap
	common.HexToHash("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"): "uniswap_v2_swap",
	// Uniswap V3 Swap
	common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"): "uniswap_v3_swap",
	// Uniswap V4 Swap
	common.HexToHash("0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f"): "uniswap_v4_swap",
	// MetaMask Swap Router
	common.HexToHash("0xbeee1e6e7fe307ddcf84b0a16137a4430ad5e2480fc4f4a8e250ab56ccd7630d"): "metamask_swap_router",
	// Fluid DEX
	common.HexToHash("0xfbce846c23a724e6e61161894819ec46c90a8d3dd96e90e7342c6ef49ffb539c"): "fluid_swap",
	// Curve TokenExchange
	common.HexToHash("0x56d0661e240dfb199ef196e16e6f42473990366314f0226ac978f7be3cd9ee83"): "curve_finance_swap",
	// Curve TokenExchange (tricrypto)
	common.HexToHash("0x143f1f8e861fbdeddd5b46e844b7d3ac7b86a122f36e8c463859ee6811b1f29c"): "curve_tricrypto_swap",
	// Curve TokenExchangeUnderlying
	common.HexToHash("0xd013ca23e77a65003c2c659c5442c00c805371b7fc1ebd4c206c41d1536bd90b"): "curve_token_exchange_underlying",
	// Curve StableSwap NG
	common.HexToHash("0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140"): "curve_stableswap_ng_swap",
	// Balancer V2 Swap
	common.HexToHash("0x2170c741c41531aec20e7c107c24eecfdd15e69c9bb0a8dd37b1840b9e0b207b"): "balancer_swap",
	// Balancer LOG_SWAP
	common.HexToHash("0x908fb5ee8f16c6bc9bc3690973819f32a4d4b10188134543c88706e0e1d43378"): "balancer_log_swap",
	// 1inch Aggregation Router V6
	common.HexToHash("0xfec331350fce78ba658e082a71da20ac9f8d798a99b3c79681c8440cbfe77e07"): "oneinch_aggregation_router_v6",

	// KyberSwap
	common.HexToHash("0xd6d4f5681c246c9f42c203e287975af1601f8df8035a9251f79aab5c8f09e2f8"): "kyberswap_swap",
	// PancakeSwap
	common.HexToHash("0x19b47279256b2a23a1665c810c8d55a1758940ee09377d4f8d26497a3577dc83"): "pancakeswap_swap",
	// DODO DODOSwap
	common.HexToHash("0xc2c0245e056d5fb095f04cd6373bc770802ebd1e6c918eb78fdef843cdb37b0f"): "dodoswap_swap",
	// DODO V2 SellBaseToken
	common.HexToHash("0xd8648b6ac54162763c86fd54bf2005af8ecd2f9cb273a5775921fd7f91e17b2d"): "dodo_v2_sell_base",
	// DODO V2 BuyBaseToken
	common.HexToHash("0xe93ad76094f247c0dafc1c61adc2187de1ac2738f7a3b49cb20b2263420251a3"): "dodo_v2_buy_base",
	// 0x Fill
	common.HexToHash("0x66a2bd850864ab5023bc4b90695fd068817db0a38bd599f6288473d20c46609f"): "zerox_fill",
	// 0x LimitOrderFilled
	common.HexToHash("0x50ae27db8b3385e989ce5067ad2962b57e8748968e2725be92d3624c8b345468"): "zerox_limit_order_filled",
	// 0x RfqOrderFilled
	common.HexToHash("0xd0c86be71d80e5d6536dc4729336b1ab10801cb568e3d9ab3da19852cfa9a0c8"): "zerox_rfq_order_filled",
	// 0x v4 OrderFilled
	common.HexToHash("0xc5feaae7fb097ff5dbe52a871af34429b2a5e29fe7256bbe9311e83df9f24d95"): "zerox_v4_order_filled",
	// 0x TransformedERC20
	common.HexToHash("0x0f6672f78a59ba8e5e5b5d38df3ebc67f3c792e2c9259b8d97d7f00dd78ba1b3"): "zerox_transformed_erc20",
	// 0x ERC20BridgeTransfer
	common.HexToHash("0x349fc08071558d8e3aa92dec9396e4e9f2dfecd6bb9065759d1932e7da43b8a9"): "zerox_erc20_bridge_transfer",
	// Paraswap Swapped
	common.HexToHash("0x6782190c91d4a7e8ad2a867deed6ec0a970cab8ff137ae2bd4abd92b3810f4d3"): "paraswap_swapped",
	// CoW Protocol Settlement
	common.HexToHash("0x40338ce1a7c49204f0099533b1e9a7ee0a3d261f84974ab7af36105b8c4e9db4"): "cow_settlement",
	// CoW Protocol Trade
	common.HexToHash("0xa07a543ab8a018198e99ca0184c93fe9050a79400a0a723441f84de1d972cc17"): "cow_trade",
	// Atlas SolverCallExecuted
	common.HexToHash("0x93485dcd31a905e3ffd7b012abe3438fa8fa77f98ddc9f50e879d3fa7ccdc324"): "solver_call_executed",
}

// DetectSwapsFromLogs scans logs for known swap events.
// Returns true if any swap was detected, along with the list of swap types found.
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
