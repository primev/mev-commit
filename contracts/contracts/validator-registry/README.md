# Validator Registry v1 Design doc

The v1 validator registry allows validators to _opt-in to mev-commit_ by staking ETH directly with the contract. This stake is separate from what the validator has already staked with the beacon chain. 

This _simple staking_ model serves as an alternative to a validator opting-in to mev-commit via restaking ETH and/or LSTs.

## Staking

Staking involves an account depositing ETH into the contract on behalf of one or more validator BLS pubkeys. Validator pubkeys are only verified by length, and not verified as a pubkey residing from an active validator on the beacon chain. Therefore stake associated with a non-active or otherwise invalid validator pubkey **can be slashed by the oracle to prevent spam**.

For the `stake` function, the account which stakes each validator pubkey is the withdrawal address for that validator. The `delegateStake` function allows only the contract owner to stake on behalf of other specified withdrawal accounts.

## Unstaking

Unstaking involves the withdrawal account for a validator pubkey calling `Unstake` for that validator. This transaction does not move funds. It does mark the validator as no longer "opted-in", and starts the process for a later withdrawal.

## Withdrawals

After a validator has been unstaked, and `unstakePeriodBlocks` amount of blocks have passed, the withdrawal account for a validator can call `withdraw`. This will transfer the validator's ETH back to their withdrawal address.

## Slashing

Note the permissioned oracle account for this contract can slash any validator that proposes a block which does not deliver preconfs from the mev-commit network. This corresponds to some configurable portion of the validator's stake being slashed (immediately sent to the contracts' `slashReceiver`).

Further, slashing automatically unstakes the relevant validator pubkey. If the relevant validator was already unstaking, the `unstakePeriodBlocks` timer is reset, and this period must be fully elapsed again before non-slashed funds are withdrawable.
