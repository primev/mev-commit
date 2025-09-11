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



# Stipend Distributor — Overview

`StipendDistributor` pays periodic (e.g., weekly) stipends to operator-defined recipients based on validator-key participation. Operators map their validator BLS pubkeys to payout addresses (“recipients”) and may authorize delegates to claim on their behalf. For more details, see the [design doc](https://www.notion.so/primev/StipendDistributor-Design-2696865efd6f80b2a4f0e6b8fc3ab0c4).

## Setting recipients

- **Global default (applies to all keys unless overridden):**
  ```solidity
  setOperatorGlobalOverride(address recipient)
  ```
  Sets a default recipient for the operator’s keys.

- **Per-key override (takes precedence over the default):**
  ```solidity
  overrideRecipientByPubkey(bytes[] calldata pubkeys, address recipient)
  ```
  Assigns a specific recipient for one or more BLS pubkeys (48-byte).

- **(Optional) Migrate unclaimed accruals between addresses:**
  ```solidity
  migrateExistingRewards(address from, address to)
  ```
  Moves **unclaimed** stipend accrued to `from` over to `to` for the calling operator.

## Delegation (optional)

- **Allow a delegate to claim for a given recipient:**
  ```solidity
  setClaimDelegate(address delegate, address recipient, bool status)
  ```
  When `status = true`, `delegate` can claim stipends for the `(operator → recipient)` pair; set `false` to revoke.

## Rewards & claiming

1. **Accrual:** Each distribution period (e.g., weekly), stipends are granted to `(operator, recipient)` pairs in proportion to validator-key participation recorded for that period.
2. **Claim by operator (pull to recipients):**
   ```solidity
   claimRewards(address payable[] calldata recipients)
   ```
   Transfers accrued amounts for the listed `recipients` to those addresses.
3. **Claim by delegate (on behalf of operator):**
   ```solidity
   claimOnbehalfOfOperator(address operator, address payable[] calldata recipients)
   ```
   Authorized delegates can trigger transfers for the specified `operator` to the listed `recipients`.

## Typical flow

1. Operator sets a **global default** recipient and (optionally) **per-key overrides**.
2. Over each period, keys that participate accrue stipends to their mapped recipients.
3. After the period, the **operator** or an **authorized delegate** calls the claim function to pay out recipients.
