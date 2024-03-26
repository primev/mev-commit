# Bridge to mev-commit chain

This repository houses a purpose-built lock/mint bridge implementation, including testing and cli tools, for a bridge between L1 ethereum and the mev-commit chain.

## Standard bridge

The [standard](./standard) directory houses a purpose-built lock/mint bridge implementation between L1 ethereum and the mev-commit chain. Its architecture prioritizes simplicity and understandability. Focusing solely on bridging native ether between the two chains, and disallowing message censorship. For more information see [here](./standard/bridge-v1/README.md). The standard bridge implementation is actively being stress tested and shows promising results. 

## Hyperlane warp route

The [hyperlane](./hyperlane) directory houses config and docker infra for a [hyperlane warp route](https://docs.hyperlane.xyz/docs/protocol/warp-routes) between L1 ethereum and the mev-commit chain. Including two validators attesting to mev-commit chain state, and a message relayer. This directory is deprecated until stress testing can confirm the relayer implementation is able to relay a decent throughput of messages with 0% drop rate. 
