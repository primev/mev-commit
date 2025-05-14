# Reward Manager

The reward manager contract allows mev-commit providers (usually L1 block builders) to send mev-boost and/or mev-commit rewards to an L1 smart contract, instead of paying proposers directly. This design enables future use-cases of the mev-commit protocol.

To pay a proposer, the mev-commit provider calls `payProposer` with the reward set as msg.value. `payProposer` only accepts a validator's BLS pubkey as an argument. The reward contract will attempt to map a pubkey to it's associated reward receiver address, checking all three methods of validator opt-in to mev-commit. So long as the provided pubkey is valid and represents a validator who's currently opted-in to mev-commit, a valid receiver address be found.

## What is a receiver address?

* For vanilla opted-in valiators, the receiver is the address that originally called `stake`
* For symbiotic opted-in validators, the receiver is the operator address
* For eigenlayer opted-in validators, the receiver is the validator's Eigenpod owner

## Overriding the receiver address

Receiver addresses have the ability to set an override address which will accumulate or be transferred rewards instead of the receiver address. Custom reward splitting logic can be implemented by the override address.

## Auto Claim

Receive addresses have the ability to enable and disable auto-claim. When auto-claim is enabled, rewards will automatically be transferred to the receiver or override address during `payProposer`. If an auto-claim transfer fails, the relevant receiver address may be blacklisted from auto-claim. Auto-claim can only be enabled and disabled by the receiver, NOT its override address.

## Manual Claim

To manually claim rewards, call `claimRewards`. This will transfer all available rewards to the calling address. Note manual claims should be made by the override address if set. If no override address is set, the receiver claims rewards.
