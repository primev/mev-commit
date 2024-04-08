# Standard bridge

This document outlines multiple iteration plans for a simple lock and mint bridging protocol between L1 ethereum and the mev-commit chain.

## User Quickstart

First build the user cli for your machine from this directory:

```bash
make user_cli
```

Or use a pre-built binary from the `releases` tab.

Set proper environment variables for a cross-chain transfer, example below:

```bash 
export PRIVATE_KEY="0xe82a054e06f89598485134b4f2ce8a612ce7f7f7e14e650f9f20b30efddd0e57"
export LOG_LEVEL="debug"
export L1_RPC_URL="https://ethereum-holesky.publicnode.com"
export SETTLEMENT_RPC_URL="https://chainrpc.testnet.mev-commit.xyz"
export L1_CHAIN_ID="17000"
export SETTLEMENT_CHAIN_ID="17864"
export L1_CONTRACT_ADDR="0xceff0a364f63f621ff6a8b5ce56569ec6f3c6220"
export SETTLEMENT_CONTRACT_ADDR="0xf60f8e762a3fe90fd4d8c005872b6f6e12eda8ca"
```

To bridge ether from Holesky to the mev-commit chain, use:

```bash
./bin/user_cli bridge-to-settlement --amount $AMOUNT_IN_WEI --dest-addr $DEST_ADDR
```

Where `PRIVATE_KEY` corresponds to an account that's funded on Holesky. 

To bridge ether back to Holesky from the mev-commit chain, use:

```bash
./bin/user_cli bridge-to-l1 --amount $AMOUNT_IN_WEI --dest-addr $DEST_ADDR
```
Where `PRIVATE_KEY` corresponds to an account that's funded on the mev-commit chain.

## Relayer

To build and run the relayer from this directory:

```bash
make relayer
```

```bash
./bin/relayer start --config=example_config/relayer_config.yml
```

## Relayer with emulators

To run a containerized relayer with five user emulators that continuously bridge back and forth, use:

```bash
make up 
```

and to include a datadog agent:

```bash
make up-agent
```

with an `.env` file specified in this directory looking like: 
```bash
DD_API_KEY=<DATADOG_API_KEY>
DD_APP_KEY=<DATADOG_APP_KEY>
```

## V1 High level design

V1 is intended to be as simple as possible, avoiding an intermediary validation network, on-chain light clients (as is used with IBC), or merkle attestations of cross chain messages. 

The v1 standard bridge will be built around a single agent type that assumes both the relayer and validator role. This will be referred to as the relayer node from now on.

To bridge to mev-commit chain, the user initiates a transaction to the contract on L1, which locks their ether on L1. The transaction should submit the necessary information to complete a cross chain transfer of funds. Importantly this transaction will emit an event which is subscribed to by the relayer.

The relayer is configured with it's own full-node for both L1, and the mev-commit chain. This can be replaced with a trusted rpc endpoint for testing.

The relayer listens to, and processes events residing from the contract on L1. Events will be handled in FIFO ordering, and would result in the data being relayed to the mev-commit chain, where native ether is minted. The destination contract accepts relay transactions only from the relayer EOA. More complex or decentralized attestation can be added in v2. 

Note to bridge from the mev-commit chain back to L1, the same protocol is used. Except mev-commit chain ether is burned upon initiating a bridge operation, and ether is unlocked on L1 upon bridge completion. Therefore the relayer should be concurrently monitoring both chains for events.

## V2 Design

V2 of the standard bridge could include the following improvements:
- Incorporate merkle attestations of cross chain messages into the transfer protocol
- Multiple relayers attesting to transfers (multisig)
- Utilize side chain state commitments posted to L1 (requires work elsewhere)

### Merkle Attestations

From some initial research, it seems like multiple bridging projects rely on merkle proof data being relayed across chains. The main difference seems to be what validator set is responsible for attesting to or managing the canonical merkle root to verify against.

Inspiration:

* Cosmos' IBC uses an on-chain light client to enable proofs of inclusion for particular values at particular paths on a remote blockchains. Essentially, the canonical merkle root is managed by an on-chain light client. See more [here](https://github.com/cosmos/ibc/tree/main/spec/core/ics-002-client-semantics). Note IBC won't work for our use-case but provides some inspiration.
* Hyperlane's multisig bridging protocol has an off-chain validator set (multisig) that attest to merkle root checkpoints of the subset of state that contains outgoing bridge messages from a source chain. A relayer agent subscribes to events related to these messages, and submits metadata (including merkle proof) to the destination chain. The destination chain contract logic verifies the relayed merkle proof data against the root attested to by the validators. See relevant contracts [here](https://github.com/hyperlane-xyz/hyperlane-monorepo/blob/5b4af6bf1db93102d54f114b03079cc873c08249/solidity/contracts/isms/multisig/AbstractMultisigIsm.sol) and [here](https://github.com/hyperlane-xyz/hyperlane-monorepo/blob/5b4af6bf1db93102d54f114b03079cc873c08249/solidity/contracts/isms/multisig/AbstractMerkleRootMultisigIsm.sol).
* Polygon POS bridging seems to piggyback off their state sync mechanism. To bridge from L1 to the sidechain, sidechain validators simply listen to events on an L1 contract and pass along this data to the sidechain as a part of the [state sync mechanism](https://docs.polygon.technology/pos/architecture/bor/state-sync/). To bridge from the sidechain back to L1, a tx must first reside on the sidechain. After some period of time, this tx is checkpointed on L1 by the sidechain validators. Once checkpointing is done, the hash of the transaction created on the sidechain is submitted with a proof to the `RootChainManager` contract on L1. This contract validates the transaction and associated merkle proof against the checkpointed root hash. That is, canonical merkle roots to verify against are posted to L1. See more about polygon bridging layers [here](https://docs.polygon.technology/pos/how-to/bridging/).
* For our implementation we may want to use [eth-getproof](https://docs.alchemy.com/reference/eth-getproof).

One possible implementation is the full data from each `TransferInitiated` event is committed to an on-chain merkle tree (separate from the full ethereum state tree). Upon each `TransferInitiated` event an off-chain validator set attests to the merkle root of said sub-tree. The proof data is relayed by the relayer, and verified on the destination. It's worth exploring how this scheme compares to using eth-getproof and the full merkle tree.  

### Multiple Relayers

V2 improvements will include decentralization of the relayer role. A simple scheme could adapt [contracts](https://github.com/primevprotocol/contracts/tree/main/contracts/standard-bridge) to require `n` out of `m` relayers each with separate EOAs, to listen for `TransferInitiated` events and/or associated metadata. Then submit `FinalizeTransfer` txes on the destination chain. 


###  Side Chain State Commitments

The mev-commit chain may adopt a state commitment scheme where merkle roots of sidechain state are periodically posted to L1 blobs, post EIP 4844. In such a scenario the bridging protocol could be adapted to utilize these commitments.
