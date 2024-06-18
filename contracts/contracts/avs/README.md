# Mev-commit AVS Implementation

## Overview

The `MevCommitAVS` contract(s) will be deployed on L1 to act as a tie-in to the eigenlayer core contracts, enabling validators to opt-in to the mev-commit protocol via restaking. It serves as the next iteration of our validator registry. Notion has a more detailed protocol design doc, whereas this doc is specific to the implementation. See the following files for current implementation:

* [IMevCommitAVS.sol](../interfaces/IMevCommitAVS.sol)
* [MevCommitAVS.sol](./MevCommitAVS.sol)
* [MevCommitAVSStorage.sol](./MevCommitAVSStorage.sol)

## Operator registration

Operators will not yet be assigned concrete tasks as a part of our AVS, however they are nonetheless able to register with our AVS to abide by the `IAVSDirectory.registerOperatorToAVS` and `IAVSDirectory.deregisterOperatorFromAVS` functions that any AVS must implement. Registration is simple for Operators and only requires providing a valid signature. To deregister, Operators must first `requestOperatorDeregistration`, wait a configurable amount of blocks, then call `deregisterOperator`. No staking is required of Operators.

Operators mainly serve the purpose of (optionally) being able to register validators on their behalf, if the relevant validator is delegated to them. Future iterations of our AVS can assign Operators required oracle tasks, as further discussed in the _Future Upgrades_ section.

## Validator Opt-in

Recall that a native-restaking enabled validator opting-in to mev-commit requires two steps:

1. The validator must delegate their native stake to an Operator who's registered with the mev-commit AVS.
2. The validator must separately *register* with the mev-commit AVS, confirming their attestation to follow the rules of the protocol.

Multiple validator public keys can be registered at once, alongside their associated eigenpod owner `podOwner` address. Note each eigenpod owner account can represent one or many restaked validators:

```solidity
function registerValidatorsByPodOwner(bytes[] calldata valPubKeys, address podOwner);
```

This function stores relevant state and ensures that the provided pubkeys are indeed actively restaked with `podOwner`'s eigenPod. Note two entities are able to register validator pub keys in this way:

1. The eigenpod owner account itself.
2. An Operator account, so long as the relevant eigenpod is delegated to that Operator.

Note if an Operator is registering pubkeys on behalf of validators, it's expected that the Operator manages those validators itself, or represents the validators to an extent that the Operator can realistically attest to the validator following the rules of mev-commit (staking-as-a-service providers for example). This trustful relationship between validators and their delegated Operator piggybacks off already agreed upon trust assumptions with eigenlayer delegation.

Deregistration requires calling `requestValidatorsDeregistration`, waiting a configurable amount of blocks, then calling `deregisterValidators`. These functions are similarly callable by the eigenpod owner OR delegated operator.

## LST Restaker Registration

LST restakers are also able to register with our avs by:

1. Delegating to an Operator who's registered with the mev-commit AVS.
2. Calling `registerLSTRestaker` with a chosen validator pubkey, taking on the freeze risk of the chosen validator.

```solidity
function registerLSTRestaker(bytes calldata chosenValidator) onlyProperlyDelegatedLSTRestaker();
```

LST restakers who register in this way will receive points or rewards when their chosen validator correctly follows the protocol.

Deregistration requires the restaker calling `requestLSTRestakerDeregistration`, waiting a configurable amount of blocks, then calling `deregisterLSTRestaker`.

## Freezing

A permissioned oracle account is able to `freeze` any registered validator for acting maliciously against agreements to the mev-commit protocol:

```solidity
function freeze(bytes calldata valPubKey) external onlyFreezeOracle();
```

To exit the frozen state, a configurable unfreeze period must first pass. Then any account can call `unfreeze`:

```solidity
function unfreeze(bytes calldata valPubKey) payable external;
```

where a minimum of `unfreezeFee` must be included in the transaction. If the validator was in the `REQUESTED_DEREGISTRATION` state prior to being frozen, the validator will be returned to the `REGISTERED` state. That is, a validator must *not* be frozen for a full deregistration period, before it's able to deregister.

Freezing is the mechanism that punishes a validator prior to eigenlayer core contracts having slashing. For now freezing corresponds to a public, reputational slash for the validator (and relevant LST restakers), and a lack of potential points accrual.

## Design Intentions

When looking through this design doc one may ask, _why do validators and LST restakers have to delegate to an Operator through the eigenlayer core contracts, AND separately register with the AVS contract?_

The answer is that **stakers** are the L1 entities that enable credible commitments through our AVS, not Operators. It would be challenging to slash/freeze entirely through Operators, in that an Operator can potentially represent thousands of validators and LST restakers from different organizations, home-staking setups, etc.

Further, we need some sort of explicit mechanism for **stakers** (not Operators) to attest to following the rules of mev-commit, at the risk of being slashed. Eigenlayer's current design does not offer this on a per-AVS basis.

## Future Upgrades 

* Operators could be given the task of replacing the oracle service that currently freezes (or will eventually slash) validators. This could rely on honest Operator majority, or a multi-tier slashing system where EIGEN holders are able to slash Operators, while the Operator set is able to slash validators.
* Operators could further be required to be validators of the existing mev-commit chain, an evm sidechain which manages preconf settlement. 

## Open Questions

* Will upcoming Eigenlayer upgrades allow for slashing stakers (validators or LST restakers) directly? Or will Operators be the only slashable entities? If the latter is true, our AVS design will need to drastically change, and may necessitate more complexity.
