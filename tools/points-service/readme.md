
## Vanilla Registry Events

The Vanilla Registry emits a `Staked` event when a validator is registered through direct ETH staking:

```solidity
event Staked(
    address indexed sender,     // The account that called stake()
    address withdrawalAddress,  // The address that can withdraw/unstake 
    bytes pubkey,              // The validator's BLS public key
    uint256 amount            // Amount of ETH staked
);
```


## Middleware Registry Events

The Middleware Registry emits a `ValRecordAdded` event when a validator is registered through Symbiotic integration:

```solidity
event ValRecordAdded(
    bytes pubkey,              // The validator's BLS public key
    address operator,          // The operator registering the validator
    address vault,             // The vault securing the validator
    uint256 position          // Position in the validator set
);
```


## AVS Registry Events

The AVS Registry emits a `ValidatorRegistered` and `LSTRestakerRegistered` (for restakers) event when a validator is registered through EigenLayer integration:

```solidity
event ValidatorRegistered(
    bytes pubkey,              // The validator's BLS public key
    address podOwner           // The eigenpod owner address
);
```

```solidity 
event LSTRestakerRegistered(
    bytes pubkey,              // The chosen validator's BLS public key
    uint256 numChosen,         // Total number of validators chosen by this LST restaker
    address lstRestaker        // Address of the LST restaker
);
```
