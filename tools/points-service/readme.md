
# MEV-Commit Points Service

This repository implements an Ethereum event listener and points accrual service for validators across different registries (Vanilla, Middleware/Symbiotic, and AVS/EigenLayer). It also exposes a simple HTTP API (see [`api.go`](./api.go)) that provides an external points interface—primarily for Symbiotic usage ([details here](https://symbioticfi.notion.site/External-Points-API-Specification-16981c079c178059b6b7e9740cf987a8)).

---

## Vanilla Registry Events

The Vanilla Registry emits a `Staked` event when a validator is registered through direct ETH staking:

```solidity
event Staked(
    address indexed sender,     // The account that called stake()
    address withdrawalAddress,  // The address that can withdraw/unstake 
    bytes pubkey,              // The validator's BLS public key
    uint256 amount             // Amount of ETH staked
);
```

---

## Middleware Registry Events

The Middleware Registry emits a `ValRecordAdded` event when a validator is registered through Symbiotic integration:

```solidity
event ValRecordAdded(
    bytes pubkey,              // The validator's BLS public key
    address operator,          // The operator registering the validator
    address vault,             // The vault securing the validator
    uint256 position           // Position in the validator set
);
```

---

## AVS Registry Events

The AVS Registry emits `ValidatorRegistered` and `LSTRestakerRegistered` (for restakers) events when a validator is registered through EigenLayer integration:

```solidity
event ValidatorRegistered(
    bytes pubkey,              // The validator's BLS public key
    address podOwner           // The eigenpod owner address
);

event LSTRestakerRegistered(
    bytes pubkey,              // The chosen validator's BLS public key
    uint256 numChosen,         // Total number of validators chosen by this LST restaker
    address lstRestaker        // Address of the LST restaker
);
```

---

## Overview

- **Event Listening**: In [`main.go`](./main.go), the service subscribes to on-chain events from each registry and updates a local SQLite database to track validator opt-ins/opt-outs.
- **Points Accrual**: The logic for computing points is found in `computePointsForMonths()`. Tests are in [`points_test.go`](./points_test.go).
- **HTTP API**: A basic REST API for querying points, found in [`api.go`](./api.go). Symbiotic usage details are [here](https://symbioticfi.notion.site/External-Points-API-Specification-16981c079c178059b6b7e9740cf987a8).



## Contributing

- Open an issue or PR for any bugs or new features.
- For questions, reach out to the maintainers.