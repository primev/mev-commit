# Block Reward Manager

`BlockRewardManager` should be used by mev-commit providers (builders) to pay a validator’s **fee recipient**.

## How it works

- To pay a proposer, call:
  ```solidity
  payProposer(address payable feeRecipient)
  ```
  (funds provided to function via msg.value)
- `feeRecipient` must be the validator’s **execution-layer fee recipient** for the block you’re paying.
- Payment is immediately forwarded to the fee recipient address. If a protocol fee is enabled, a small percentage of payment is reserved in the contract for mev-commit participant rewards. This fee will initially be switched off.

## Usage example

**Solidity (from a builder integration):**
```solidity
IBlockRewardManager(brm).payProposer{value: reward}(feeRecipient);
```



# Reward Distributor — Overview (ETH-focused)

For further details, see the Reward Distributor Design Doc: https://www.notion.so/primev/RewardDistributor-Design-2696865efd6f80b2a4f0e6b8fc3ab0c4

`RewardDistributor` tracks and pays **ETH stipends** to operator-defined recipients based on validator-key participation. Operators map their validator BLS pubkeys to payout addresses (“recipients”) and may authorize delegates to claim on their behalf.

> **Note:** The contract also supports granting **ERC20 token rewards** (single-token per `tokenID`) and will utilize this in the future. In the near term, **ETH (tokenID `0`)** is the primary path.

---

## Setting recipients

- **Global default (applies to all keys unless overridden):**  
  `setOperatorGlobalOverride(address recipient)`  
  Sets a default recipient for the operator’s keys.

- **Per-key override (takes precedence over the default):**  
  `overrideRecipientByPubkey(bytes[] pubkeys, address recipient)`  
  Assigns a recipient for one or more BLS pubkeys (each must be 48 bytes).  
  **Precedence:** per-key override → global override → fallback to the operator address.

- **Resolve the active recipient for a key:**  
  `getKeyRecipient(address operator, bytes pubkey) → address`  
  Returns the payout address considering per-key overrides, global override, or operator fallback.

- **Migrate unclaimed accruals between recipients (for the calling operator):**  
  `migrateExistingRewards(address from, address to, uint256 tokenID)`  
  Moves **unclaimed** rewards for `(msg.sender, from, tokenID)` into `(msg.sender, to, tokenID)`.  
  Use `tokenID = 0` for ETH; future token IDs correspond to configured ERC20s.

---

## Delegation (optional)

- **Authorize a delegate to claim for a specific recipient:**  
  `setClaimDelegate(address delegate, address recipient, bool status)`  
  When `status = true`, `delegate` may claim for `(operator, recipient)`; set `false` to revoke.

- **Delegate claim (on behalf of operator):**  
  `claimOnbehalfOfOperator(address operator, address payable[] recipients, uint256 tokenID)`  
  Delegate must be authorized **per recipient** by that operator.

---

## Rewards & claiming (ETH-first)

1. **Accrual off-chain; batched grants on-chain.**  
   A RewardManager service monitors blocks won by mev-commit–registered validators, resolves recipients via `getKeyRecipient`, and tallies `(operator, recipient)` off-chain. At period end, it submits **batched grants**:

   - **ETH grants (primary path):**  
     `grantETHRewards(Distribution[] distributions)`  
     The transaction `msg.value` must equal the sum of `distributions[i].amount`.

   - **Token grants (future use):**  
     `grantTokenRewards(Distribution[] distributions, uint256 tokenID)`  
     Pulls tokens from `msg.sender` via `transferFrom`. Requires prior `approve`.

   A `Distribution` item packs: `{operator, recipient, amount}`.

2. **Claim by operator (pull-to-recipient):**  
   `claimRewards(address payable[] recipients, uint256 tokenID)`  
   For each recipient listed, transfers the **full pending** amount for that `(operator, recipient, tokenID)` bucket to the recipient.  
   - ETH: use `tokenID = 0`.  
   - Tokens: use the configured nonzero `tokenID` (future use).

3. **Get pending rewards:**  
   `getPendingRewards(address operator, address recipient, uint256 tokenID) → uint128`  
   Computed as `accrued - claimed` for that bucket.

> **Isolation guarantee:** Balances are **strictly partitioned** by `(operator, recipient, tokenID)`. Multiple operators can share a recipient, but each operator can only claim their own bucket for that recipient.

---

## Reclaim by owner (administrative)

- **Owner can reclaim accrued rewards to itself:**  
  `reclaimStipendsToOwner(address[] operators, address[] recipients, uint256 tokenID)`  
  Sums claimable amounts across the provided pairs, transfers them to the **owner**, and zeroes the accruals.  
  Requirements: arrays must be equal length; total claimable must be nonzero.

---

## Access control & safety

- **Grant permissions:** Only the **owner** or the **reward manager** may call `grantETHRewards` / `grantTokenRewards`.  
- **Pause:** Owner can `pause()`/`unpause()` to block mutating endpoints (grants, claims, and—if configured—delegation/override changes).  
- **Zero-address checks & input validation:** Functions validate parameters (e.g., nonzero addresses, 48-byte pubkeys, registered token IDs, consistent array lengths).

---

## Events (key ones)

- **ETH grants:** `ETHGranted(address indexed operator, address indexed recipient, uint256 indexed amount)`  
- **ETH claims:** `ETHRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)`
- **Operator Reward Migrations** `RewardsMigrated(uint256 tokenID, address indexed operator, address indexed from, address indexed to, uint128 amount)`

- **(Future) token grants:** Implementation emits analogous token-grant events and batch totals where applicable.  
- **Admin updates:** Events are emitted for reward manager and token registration changes.

> Event names above match the interface (`IRewardDistributor`) for ETH. Token event names may vary depending on your current implementation; update this section if you finalize them.

---

## Typical ETH stipend flow

1. Operator sets a **global default** recipient and optional **per-key overrides**.  
2. Over the period, off-chain tally adds up ETH stipends per `(operator, recipient)`.  
3. After the period, RewardManager calls `grantETHRewards([...])` with the **consolidated totals**.  
4. Operator (or an authorized delegate) calls `claimRewards([recipients], 0)` to transfer ETH to each listed recipient.

---

## Notes on future token support

- **Registration:** Owner maps an ERC20 to a nonzero `tokenID` via `setRewardToken(address token, uint256 tokenID)`.  
- **Grants:** Use `grantTokenRewards(distributions, tokenID)`; caller must hold tokens and `approve` the distributor.  
- **Claims:** Same claim APIs as ETH but pass the nonzero `tokenID`. Buckets remain isolated per `tokenID`.

---
