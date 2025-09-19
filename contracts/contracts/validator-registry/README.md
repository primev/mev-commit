# Validator Registry Design doc

Validators are able to _opt-in to mev-commit_ in one of three ways:

1. Restaking with the `MevCommitAVS` contract.
2. Restaking with the `MevCommitMiddleware` contract.
3. Simple staking with the `VanillaRegistry` contract.

The `ValidatorOptInHub` contract acts as a query router between all three solutions, allowing any actor to query whether a group of validator pubkeys is opted-in to mev-commit. This is an updated version of the ValidatorOptInRouter contract for ease of use and to add future registry support.

## Mev-commit AVS - Eigenlayer Restaking Solution

For more details on the `MevCommitAVS` contract, please refer to the [MevCommitAVS README](avs/README.md).

## Mev-commit Middleware - Symbiotic Restaking Solution

For more details on the `MevCommitMiddleware` contract, please refer to the [MevCommitMiddleware README](middleware/README.md).

## Vanilla Registry - Simple Staking Solution

The vanilla registry allows validators to _opt-in to mev-commit_ by staking native ETH directly with the contract. This stake is separate from a validator's 32 ETH already staked with the beacon chain. 

### Staking

Staking involves an account depositing ETH into the contract on behalf of one or more validator BLS pubkeys. Validator pubkeys are only verified by length, and not verified as a pubkey residing from an active validator on the beacon chain. Therefore stake associated with a non-active or otherwise invalid validator pubkey **may be slashed by the oracle to prevent spam**.

For the `stake` function, the account which stakes each validator pubkey is the withdrawal address for that validator. The transaction calling `stake` must attach `minStake` ETH for each validator pubkey to be staked. The `delegateStake` function allows only the contract owner to stake on behalf of other specified withdrawal accounts.

### Unstaking

Unstaking involves the withdrawal account for a validator pubkey calling `Unstake` for that validator. This transaction does not move funds. It does mark the validator as no longer "opted-in", and starts the process for a later withdrawal.

### Withdrawals

After a validator has been unstaked, and `unstakePeriodBlocks` amount of blocks have passed, the withdrawal account for a validator can call `withdraw`. This will transfer the validator's ETH back to their withdrawal address.

The owner account is also able to call `withdraw` for any validator pubkey who's already unstaked. This gives the owner the ability to forcefully delete records for validators that have either been slashed, or explicitly unstaked by the user.

### Slashing

Note the permissioned oracle account for this contract can slash any validator that proposes a block which does not deliver preconfs from the mev-commit network. This corresponds to some configurable portion of the validator's stake being slashed (immediately sent to the contracts' `slashReceiver`).

Further, slashing automatically unstakes the relevant validator pubkey. If the relevant validator was already unstaking, the `unstakePeriodBlocks` timer is reset, and this period must be fully elapsed again before non-slashed funds are withdrawable.

### Configuration of `unstakePeriodBlocks`

`unstakePeriodBlocks` must be set such that a mev-commit bidder knows all _currently opted-in_ block proposers from the current epoch, and next epoch, must deliver commitments, or are guaranteed slashable.

For similar reasoning to the `slashPeriodSeconds` configuration in the [middleware README](middleware/README.md#configuration-of-slashperiodseconds), the minimum value for `unstakePeriodBlocks` is:

`6 L1 epochs` + `oracleProcessingPeriod`

A recommended value to assume for `oracleProcessingPeriod` is 60 minutes, although longer is always beneficial as a buffer for oracle transaction inclusion.
