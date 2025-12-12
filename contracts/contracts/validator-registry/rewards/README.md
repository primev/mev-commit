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


# Reward Distributor — Operator Guide

The `RewardDistributor` contract is used to **receive, track, and distribute additional ETH or token rewards** earned by operators with active validators.

When operators earn extra rewards (for example, from mev-commit participation or other incentive programs), those rewards are **granted to this contract** and then **claimed by operators to their specified fee recipients**.

 
> Rewards are granted on-chain to the `RewardDistributor`, tracked per operator and recipient, and later claimed by operators (or their delegates) to the correct payout addresses.

---

### Core Contract Functionality

- Hold **additional ETH or token rewards** earned by operators
- Allow operators to **define where rewards should be paid**
- Ensure rewards are **isolated per operator and recipient**
- Support **delegated claiming** if operators want others to claim on their behalf

The contract does **not** automatically push funds. Operators (or authorized delegates) explicitly claim rewards when ready.

---

### Core roles/concepts

- **Operator**  
  The address that controls one or more validators and interacts with the contract.

- **Recipient**  
  The address that ultimately receives rewards (often the validator’s fee recipient or another payout address).

- **Validator pubkey mapping**  
  Operators can map validator BLS pubkeys to specific recipients.


Rewards are tracked separately for each:

(operator, recipient, tokenID)

---

## Setting reward recipients

>If a recipient is not set, rewards earned by a validator are granted to the validator's operator.

### 1. Set a global default recipient

Applies to all validator keys unless overridden:

setOperatorGlobalOverride(address recipient)

If no per-key override exists, rewards for the operator’s validators will go here.

---

### 2. (Optional) Set per-validator overrides

If desired, an operator can set per-key recipients which override the operator's global default recipient:

overrideRecipientByPubkey(bytes[] pubkeys, address recipient)

- Each pubkey must be exactly **48 bytes**
- Can be used for one or many pubkeys
- Takes precedence over the global override

---

### 3. How recipients are resolved

For a given validator key, the active recipient is resolved as:

1. Per-key override  
2. Global override  
3. Operator address (fallback)

You can query this directly:

getKeyRecipient(address operator, bytes pubkey) → address

---

## Receiving rewards

Rewards are **granted to the contract**, not sent directly to recipients.

- Grants are performed by an authorized reward manager or the contract owner
- Rewards correspond to operator activity (e.g., active validators)
- Operators do **not** need to take action to receive grants

Once granted, rewards accumulate under the operator until claimed.

---

## Claiming rewards

### Claim as the operator

To transfer all pending rewards to one or more recipients:

claimRewards(address payable[] recipients, uint256 tokenID)

- ETH rewards: `tokenID = 0`
- Token rewards: use the configured nonzero `tokenID`
- Transfers the **full pending amount** for each listed recipient

---

## Delegating claims (optional)

Operators may allow another address to claim rewards on their behalf.

setClaimDelegate(address delegate, address recipient, bool status)

Delegation is configured **per recipient**.

---

## Safety and guarantees

- Rewards are **strictly isolated** by `(operator, recipient, tokenID)`
- Operators can only claim their own rewards
- Multiple operators may use the same recipient address safely
