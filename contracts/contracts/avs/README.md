# Mev-commit AVS Implementation

## Overview

The `MevCommitAVS` contract will be deployed to L1 to act a tie-in to the eigenlayer core contracts, enabling validators to opt-in to the mev-commit protocol via restaking. It serves as the next iteration of our validator registry. Notion has a more detailed protocol design doc, whereas this doc is specific to the implementation.

[MevCommitAVS.sol](./MevCommitAVS.sol)
[MevCommitAVSStorage.sol](./MevCommitAVSStorage.sol)

## Operator registration

Operators will not yet be assigned concrete tasks as a part of our AVS, however they are nonetheless able to register with our AVS to abide by the `IAVSDirectory.registerOperatorToAVS` and `IAVSDirectory.deregisterOperatorFromAVS` functions that any AVS must implement. Registration is simple for operators and only requires providing a valid signature. To deregister, operators must first `requestOperatorDeregistration`, wait a configurable amount of blocks, then call `deregisterOperator`, with no staking required.

Operators simply act as a (required) placeholder for now, but future iterations can assign them concrete tasks as further discussed in the _Open Questions_ section.

## Validator Opt-in

Recall that a native-restaking enabled validator opting-in to mev-commit requires two steps:

1. The validator must delegate their native stake to an Operator who's registered with the mev-commit AVS.
2. The validator must separately *register* with the mev-commit AVS, confirming their attestation to follow the rules of the protocol.

Since each eigenpod owner account possibly represents one or many restaked validators, any number of validator public keys can be registered, as long as an associated eigenpod owner `podOwner` is provided:

```solidity
function registerValidatorsByPodOwner(bytes[] calldata valPubKeys, address podOwner);
```

This function stores relevant state and ensures that the provided pubkeys are indeed actively restaked with `podOwner`'s eigenPod. Note two entities are able to register validator pub keys in this way:

1. The eigenpod owner account itself.
2. An operator account, so long as the relevant eigenpod is delegated to that operator.

Deregistration requires calling `requestValidatorsDeregistration`, waiting a configurable amount of blocks, then calling `deregisterValidators`. These functions are similarly callable by the eigenpod owner OR delegated operator.

## LST Restaker Registration

LST restakers are also able to register with our avs by:

1. Delegating to to an Operator who's registered with the mev-commit AVS.
2. Calling `registerLSTRestaker` with a chosen validator pubkey, taking on the reward or freeze risk of the chosen validator.

```solidity
function registerLSTRestaker(bytes calldata chosenValidator) onlyProperlyDelegatedLSTRestaker();
```

Deregistration requires calling `requestLSTRestakerDeregistration`, waiting a configurable amount of blocks, then calling `deregisterLSTRestaker`. These functions are callable by the LST restaker or delegated operator.

## Freezing
... replaces slashing for now

## Design Intentions

When looking through this design doc one may ask _why do validators and LST restakers have to delegate to an Operator through the eigenlayer core contracts, AND separately register with the AVS contract?_

The answer is that *validators* are the entities that enable credible commitments in our protocol. It would be challenging to reward or slash/freeze entirely through Operator delegation, in that an Operator can potentially represent thousands of validators from different organizations, home-staking setups, etc.

Further, we need some sort of explicit mechanism for validators to confirm on-chain that they wish to participate in our protocol via restaking, and generate additional revenue at the risk of being slashed. Eigenlayer's current design does not offer this.

## Open Questions

* Will upcoming Eigenlayer upgrades allow for rewarding and slashing validators directly? Or will Operators be the only rewardable/slashable entities? If the latter is true, this design will need to change.


// Write about how v2 (or future version with more decentralization) will 
// give operators the task of doing the pubkey relaying to the mev-commit chain. 
// That is the off-chain process is replaced by operators, who all look for the 
// valset lists posted to some DA layer (eigenDA?), and then race/attest to post
// this to the mev-commit chain. The operator accounts could be auto funded on our chain. 
// Slashing operators in this scheme would require social intervention as it could
// be pretty clear off chain of malicous actions and/or malicious off-chain validation
// of eigenpod conditions, delegation conditions, etc. 

// TODO: Whitelist is now just operators! Every large org seems to have its own operator.
// Note this can be what "operators do" for now. ie. they have the ability to opt-in their users. 
// But we still allow home stakers to opt-in themselves too. 
// Make it very clear that part 2 of opt-in is neccessary to explicitly communicate to 
// the opter-inner that they must follow the relay connection requirement. Otherwise delegators may be 
// blindly frozen. When opting in as a part of step 2, the sender should be running the validators
// its opting in (st. relay requirement is met).

// TODO: overall gas optimization
// TODO: order of funcs, finish interfaces, comments for everything etc.
// TODO: use tests from other PR? 
// TODO: test upgradability before Holesky deploy
// TODO: Note and document everything from https://docs.eigenlayer.xyz/eigenlayer/avs-guides/avs-dashboard-onboarding
// TODO: Confirm all setters are present and in right order, confirm interface is fully populated
// TODO: Non reentrant or is this not relevant? 
// TODO: Decide how multisig will tie into this contract, likely use gnosis safe? 
// TODO: worth adding validator blacklist?

// TODO; open questions section in design doc.. do we need to have operators? Or will eigen support rewards/slashing directly to stakers?
// where stakers could be eigenpod owners themselves or LST restakers.
// Make sure to ask/address whether we'll be able to slash validators specifically! Not an entire operator group. 

// TODO: Look into edge cases around validators being like nah dawg I hit the 10 limit but never got to delegate myself!
// ^ solution to above is likely to always allow self LST delegation from podOwner, limit is only for rando acconts

// TODO: Look into edge cases around validators being like nah dawg I hit the 10 limit but never got to delegate myself!
// ^ solution to above is likely to always allow self LST delegation from podOwner, limit is only for rando acconts

// TODO: see if possible to launch with rewards, you're using correct branch
// TODO: add to docs the importance of association of operators/restakers to valAddrs
