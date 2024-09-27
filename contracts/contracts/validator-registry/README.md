# Validator Registry Design doc

Validators are able to _opt-in to mev-commit_ in one of three ways:

1. Vanilla staking with the `VanillaRegistry` contract.
2. Native restaking with the Eigenlayer-integrated `MevCommitAVS` contract.
3. ERC20 restaking with the Symbiotic-integrated `MevCommitMiddleware` contract.

The `ValidatorOptInRouter` contract acts as a query router between all three solutions, allowing any actor to query whether a group of validator pubkeys is opted-in to mev-commit.

## Mev-commit AVS - Native Restaking Solution

For more details on the Mev-commit AVS, please refer to the [Mev-commit AVS README](avs/README.md).

## Mev-commit Middleware - ERC20 Restaking Solution

For more details on the Mev-commit Middleware, please refer to the [Mev-commit Middleware README](middleware/README.md).

## Vanilla Registry - Simple Staking Solution

The vanilla registry allows validators to _opt-in to mev-commit_ by staking native ETH directly with the contract. This stake is separate from a validator's 32 ETH already staked with the beacon chain. 

### Staking

Staking involves an account depositing ETH into the contract on behalf of one or more validator BLS pubkeys.

For the `stake` function, the account which stakes each validator pubkey is the withdrawal address for that validator. The `delegateStake` function allows only the contract owner to stake on behalf of other specified withdrawal accounts.

### Unstaking

Unstaking involves the withdrawal account for a validator pubkey calling `Unstake` for that validator. This transaction does not move funds. It does mark the validator as no longer "opted-in", and starts the process for a later withdrawal.

### Withdrawals

After a validator has been unstaked, and `unstakePeriodBlocks` amount of blocks have passed, the withdrawal account for a validator can call `withdraw`. This will transfer the validator's ETH back to their withdrawal address.

### Slashing

Note the permissioned oracle account for this contract can slash any validator that proposes a block which does not deliver preconfs from the mev-commit network. This corresponds to some configurable portion of the validator's stake being slashed (immediately sent to the contracts' `slashReceiver`).

Further, slashing automatically unstakes the relevant validator pubkey. If the relevant validator was already unstaking, the `unstakePeriodBlocks` timer is reset, and this period must be fully elapsed again before non-slashed funds are withdrawable.

### Blacklisting

Upon registration, validator pubkeys are only verified by length, and not verified as a pubkey residing from an active validator on the beacon chain. Without an implicit, on-chain way to verify the association between L1 validator pubkeys and a withdrawal address, withdrawal addresses are **trusted by the owner** to register only **active beacon chain validator** pubkeys for which **they have access to the private key**.

Any validator pubkey can only be mapped to a single withdrawal address. So if a withdrawal address registers a pubkey for which they do not own/manage (think “greifing”), the contract owner account reserves the right to *blacklist* that withdrawal address. The contract owner also reserves the right to *blacklist* withdrawal addresses who register non-active or sybiled L1 validator pubkeys. Discretion is left to the contract owner account to resolve social disputes off-chain.

An on-chain dispute mechanism could eventually replace permissioned blacklisting.

Blacklisting consists of the contract owner account marking particular withdrawal addresses as blacklisted. All validator pubkeys associated to blacklisted withdrawal address(es) are no longer considered *opted-in*. Further, the contract owner has the ability to unstake and withdraw validator pubkeys associated with blacklisted withdrawal address(es), thus allowing non-malicious withdrawal addresses to register previously greifed validator pubkeys.

Note when the contract owner withdraws on behalf of a blacklisted validator, funds are still transferred to the blacklisted address.
