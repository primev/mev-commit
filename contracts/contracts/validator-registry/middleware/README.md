# Mev-commit Network Middleware Implementation

The `MevCommitMiddleware` contracts serve as an entrypoint for L1 validators to *opt-in* to mev-commit, and attest to the rules of our protocol, at the risk of funds being slashed. 

# Background

Before diving into protocol specifics, check out [Symbiotic docs](https://docs.symbiotic.fi/category/core-modules) and familiarize yourself with Vaults, Networks, and Operators.

# TLDR

![Integration](./integration.png)

*Figure 1: Integration Overview*

For L1 validators to be *opted-in to mev-commit*, some collateral stake must be slashable for each validator, in case that validator acts against its protocol commitments. In the context of Symbiotic, Vaults allocate slashable ERC20 collateral to Operators that are registered with the mev-commit network. This collateral may be slashable by other networks, hence it is *restaked*.

Operators for the mev-commit network are responsible for bulk registering groups of L1 validator pubkeys to an associated vault. Every registered validator is represented by restaked collateral from a single Vault and Operator. Each Vault’s total collateral can be split up to secure/represent many validators in groups**.**

Our network middleware contract requires any Vault registering with the contract to have a delegator of the `NetworkRestakeDelegator` type or `OperatorSpecificDelegator` type. `FullRestakeDelegator` is disallowed due to its ability for vaults to reuse stake within the same network to multiple Operators. In other words, a single instance of `slashAmount` vault collateral can only be used to secure a single validator, through a single Operator. Vaults are still able to split their slashable collateral across multiple Operators.

# Steps for validator opt-in via symbiotic

1. Network EOA opts-in to Vault which will secure the validator(s).
    1. Delegator and slasher requirements can be validated off-chain. 
2. Vault manager sets a non-zero stake limit for mev-commit Network. 
3. Operator opts-in to both Vault and Network.
4. Vault manager allocates stake to relevant Operator for mev-commit Network. 
5. Network middleware contract owner registers both Operator and Vault with the Network. 
6. Operator bulk register groups of validator pubkeys associated to a Vault. 

# Design Detail

### The role of Operators

The Symbiotic protocol defines *Operators* as the entities running infrastructure for decentralized networks. Since the mev-commit oracle slashes on a per-L1-proposer basis, we need to explicitly associate Operators to individual validator pubkeys that the oracle would use for slashing logic. 

Operators’ main purpose will be to bulk register groups of validator pubkeys (each associated to a vault) that will represent the Operator. Vault collateral allocated to that Operator, will be slashable in case relevant validators acts against their agreements.

Any operator registering a validator pubkey agrees to the following:

<aside>
💡 A validator’s mev-boost client should ONLY connect to mev-commit opted-in relays, to ensure proposed blocks deliver commitments made. The Titan Holesky relay is the only supporting relay at this time. This list will be updated as more relays support the network.
</aside>

If any validator pubkey acts against this agreement as determined by the mev-commit oracle, Vault collateral allocated to the Operator who registered said pubkey will be slashed.

### Operator Registration and Blacklisting

Each Operator entity must be registered by the `MevCommitMiddleware` contract owner, and may be deregistered by the contract owner as needed. On mainnet the contract owner will be a Primev multisig, this multisig may need to execute regular transactions which register Operator(s).

Without an implicit, on-chain way to verify the association between L1 validator pubkeys and Operators through Symbiotic, Operators are **trusted by the owner** to register only pubkeys for which they have control over.

Any validator pubkey can only be mapped to a single Operator. So if an Operator registers a pubkey for which they do not own/manage (think “griefing”), the contract owner account reserves the right to *blacklist* Operators. The contract owner also reserves the right to *blacklist* Operators who register non-active or sybiled L1 validator pubkeys.

An on-chain dispute mechanism could eventually replace permissioned blacklisting, but is not worth targeting for v1.

Blacklisting consists of the owner account marking a particular Operator as blacklisted, regardless of that Operator’s previous state within the contract. Once blacklisted, all validator pubkeys registered by the Operator are no longer considered *opted-in.* Further, the contract owner has the ability to delete validator records associated with blacklisted Operators, thus allowing non-malicious Operators to register previously griefed validator pubkeys.

### **What defines a validator being "opted-in"?**

A validator pubkey being opted-in following registration is defined by the following criteria:

1. The validator pubkey must be fully registered
    1. Validator record stored. 
    2. De-registration request must not exist.
2. The validator pubkey’s associated Vault must be registered. 
    1. Vault record stored. 
    2. De-registration request must not exist.
    3. Vault is entity to core vault factory contract. 
3. The validator pubkey’s associated Operator must be registered
    1. Operator record stored. 
    2. De-registration request must not exist.
    3. The associated Operator must not be blacklisted. 
    4. Operator is entity to core operator registry contract. 
4. The validator must be *slashable* as defined in the next section.

Directly following a successful validator registration, all of these criteria will be true as enforced by the function. However, anyone of these criteria becoming false will result in the validator no longer being *opted-in* from the mev-commit protocol's perspective. For example, if slashable funds are withdrawn from a Vault collateralizing a group of validators, some of those validators will no longer be *slashable* and thus no longer *opted-in* (particularly the validators with the highest insertion index).

### Indexing within a valSet

When an Operator registers validator pubkeys in a group (recall a validator set, valSet, corresponds to a single Vault/Operator), the index of the pubkey within the valSet acts as a priority rating representing which validators should be collateralized by the Vault over others. A lower index corresponds to a higher priority. This priority index can be queried for a pubkey with respect to its valSet using `getPositionInValset`.

Position within a valSet is **not guaranteed to be preserved**. That is, various actions can change the priority index of a validator within a valSet. Concretely the following events correspond to ordering changes:

1. `ValRecordAdded` corresponds to a new validator being added to a valSet, with the highest priority index.
2. `ValidatorPositionsSwapped` corresponds to a swap of positions within the valSet. Slashed validators' indexes are swapped with the highest indexes in the valSet.
3. `ValRecordDeleted` corresponds to a validator being removed from a valSet, using the swap-and-pop removal method.

### What defines a validator being slashable?

The middleware contract must set a `slashAmount` for each ERC20/Vault. This much stake is required to define a single validator as being *slashable*. It must be large enough to disincentivize validators from acting maliciously against mev-commit. For a Vault ↔ valSet pair, `slashAmount * numVals` worth of collateral must be allocated to the registering Operator, to define all validators in the set as *slashable*. 

If less than `slashAmount * numVals` is allocated (restaked) to an Operator, only the first `allocatedStake / slashAmount` number of validators will be considered *slashable,* as determined by the index of each stored validator record. Once additional Vault collateral is deposited and/or allocated to the Operator, relevant registered validators become *slashable* again. Operators are able to deregister pubkeys if slashable collateral becomes too low in an existing Vault.

It's an operator's responsibility to monitor vault collateral, and make sure all registered validators are also slashable (therefore "opted-in"). This means if vault collateral is reduced to a value that does not define all validators as slashable, the operator must deregister validators of its choice, or implicitly accept that some quasi-random validators will no longer be "opted-in".

### Slash mechanics

The mev-commit oracle has a permissioned account which is exclusively given rights to slash certain validator pubkeys which act maliciously or incorrectly against mev-commit. On-chain, the mapped validator record is retrieved, and the associated Vault/Operator securing said pubkey will be slashed `slashAmount` as configured by the middleware contract owner.

The oracle must continuously monitor opted-in status of upcoming L1 proposers. In doing this, the oracle must keep track of a `captureTimestamp` for each opted-in, upcoming proposer. The `captureTimestamp` is the timestamp of the most recent queried **finalized** L1 block that the validator pubkey was queried as *opted-in* by the oracle.

For validators who proposed incorrectly as determined by the oracle, slashing must be executed by the oracle account on L1 within `slashPeriodSeconds` of a validator's relevant `captureTimestamp`. The oracle attaches a `captureTimestamp` to each validator pubkey when calling `slashValidators`. Concretely, `slashValidators` is guaranteed to succeed, if the oracle provides valid `captureTimestamp`s, and `slashValidators` executes on-chain before each `captureTimestamp + slashPeriodSeconds` timestamp.

`slashPeriodSeconds` also enforces a minimum amount of time that must elapse before validator records, operator records, or vault records can be deleted from the middleware contract’s state, as these records are essential to proper slashing.

`slashPeriodSeconds` is initially configured by the middleware contract owner, and can be mutated by the owner at any time.

### Instant vs Veto slashers

Read more about Symbiotic slashing [here](https://docs.symbiotic.fi/core-modules/vaults#slashing).

Vaults with instant slashers must have an `epochDuration` greater than `slashPeriodSeconds` to register with our middleware contract, ensuring collateral is slashable during the full slashing period. 

Vaults with veto slashers:

* must have an `epochDuration` greater than `slashPeriodSeconds` + `vetoDuration`, where `vetoDuration` is specified by the slasher. 
* require the resolver to be disabled via `address(0)`, since a permissioned oracle account invokes slashing, requiring only the most basic slashing interface.

The oracle has at most `slashPeriodSeconds` to request a slash via `slashValidators`. Both `VetoSlasher.requestSlash` and `VetoSlasher.executeSlash` get called during `slashValidators`. Upon the oracle successfully calling `slashValidators`, the middleware contract emits a `ValidatorSlashed` event.

Read more about Symbiotic veto slashing [here](https://docs.symbiotic.fi/core-modules/vaults#veto-slashing).

The described process considers Symbiotic's current design, where it's possible to execute a slash right after requesting one, if the resolver is disabled. See source below:

1. `requestSlash` must be called before `epochDuration` - `vetoDuration` (equal to `slashPeriodSeconds` for our network) has passed since `captureTimestamp`, [source](https://github.com/symbioticfi/core/blob/629b9faac2377a9eb9cfdc6362b30d1dc1ef48f2/src/contracts/slasher/VetoSlasher.sol#L93).
2. `executeSlash` must be called before `epochDuration` has passed since `captureTimestamp`, [source](https://github.com/symbioticfi/core/blob/629b9faac2377a9eb9cfdc6362b30d1dc1ef48f2/src/contracts/slasher/VetoSlasher.sol#L152).

### Burner Router Requirements

Any vault integrating with the mev-commit middleware contract must use a [burner router](https://docs.symbiotic.fi/modules/extensions/burners#burner-router). This is validated by the symbiotic core `BurnerRouterFactory` contract. Prior to vault and validator registration, the burner router must be configured as follows:

* `IBurnerRouter.networkReceiver()` must be set to `MevCommitMiddlewareStorage.slashReceiver`.
* `IBurnerRouter.delay()` must be greater than `MevCommitMiddlewareStorage.minBurnerRouterDelay`, a suggested value for the latter param is `2 days`.
* `IBurnerRouter.operatorNetworkReceiver()` must be disable by setting to `address(0)`, or set to `MevCommitMiddlewareStorage.slashReceiver`. Essentially this value must not override a valid network receiver.

Upon vault registration, and validator registration, these conditions are checked. If later on, these conditions are not met, all validators associated to a vault will no longer be opted-in, as enforced in `_isValidatorOptedIn`.

### Vault Trust Assumptions

For vaults that use burner routers, the following scenario is possible. The vault can create a network and specify its own receiver for that network, which will transfer funds to a vault-owned address. A validator represented by this vault could then act maliciously. As soon as the vault realizes/confirms that it has misbehaved w.r.t the mev-commit network, the vault would be slashed from its custom network, transferring vault collateral away from the vault, and thus avoiding losing funds to the mev-commit slash.

Due to this scenario, the mev-commit network must fully trust all vaults that integrate with the middleware contract. Note `registerVaults` and related functions are `onlyOwner`, allowing the owner to choose its trusted vault set.

To alleviate the described vault issue, the mev-commit middleware contract enforces a minimum burner router delay, allowing time for off-chain monitoring services to observe pending changes in vault state, and try to detect malicious behavior. Failed slashing attempts will also be monitored off-chain, and if a vault is found to be maliciously escaping slashing, they will be immediately deregistered, hopefully prior to enacting the malicious behavior with multiple validators. 

Further, the described scenario carries significant reputational risks. Large vaults engaging in such actions would essentially ruin their reputation and ability to continue operating in the Symbiotic ecosystem. Thus, we assume that well known vaults have little incentive to perform the attack.

A future upgrade could introduce a trusted, primev curated vault, that would only integrate with mev-commit and select other trusted networks.

### Configuration of slashPeriodSeconds

`slashPeriodSeconds` must be set such that a mev-commit bidder knows all _currently opted-in_ block proposers from the current epoch, and next epoch, must deliver commitments, or are guaranteed slashable.

Note _currently opted-in_ in this context, means the validator is opted-in with respect to the latest finalized L1 block state, in the worst case this is two L1 epochs ago.

Consider the following scenario to exemplify how `slashPeriodSeconds` must be set:

* Current block is start of epoch n.
* Bidder queries finalized opted-in status (from epoch n-2) for 64 upcoming validators who will propose in epoch n and n+1.
* Final proposer in epoch n+1 proposes an invalid block.
* Oracle observes the finalized validator infraction at the end of epoch n+3. 
* Oracle takes some amount of time to get its slash transaction included.
* The slash is initiated on-chain (either `slash` or `requestSlash` depending on slasher type).

Concretely, for validators to be slashable by the oracle, `slashPeriodSeconds` must be greater than:

`6 L1 epochs` + `oracleProcessingPeriod`

A recommended value to assume for `oracleProcessingPeriod` is 60 minutes, although depending on vault constraints, assuming a longer period could make the chances of oracle transaction inclusion more likely.

### Rewards

Vault depositors will receive points commensurate with associated validator pubkey(s) being *opted-in* over time. When a Vault is represented by multiple validators, attribution is split evenly between the validators.

Validator opt-in state can be queried as described above. This query offers concrete criteria that must be true for Vault depositors to accrue points/rewards over time from an associated validator.

Points/rewards would be computed off-chain, with heavy use of indexed events. The exact point weightings associated to different actors/events is yet to be finalized.
