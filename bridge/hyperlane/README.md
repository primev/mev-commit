# Hyperlane Warp Route Bridge

This directory houses a [hyperlane warp route](https://docs.hyperlane.xyz/docs/protocol/warp-routes) between Sepolia and the mev-commit chain, with two validators attesting to mev-commit chain state, and a message relayer.

## Bridge CLI

The bridge cli is built out as a shell script that interacts with bridging contracts on both L1 and the mev-commit chain. The cli must first be initialized with relevant contract addresses, chain IDs, and RPC endpoints. The cli user can then bridge in either direction accordingly, to any destination account. `cli.sh` requires both [foundry](https://book.getfoundry.sh/getting-started/installation) and `jq` to be installed on the host. 

We encourage anyone using the bridge cli to understand the underlying shell script they're executing. It's essentially a simple wrapper around some foundry commands that invoke bridging txes.

### Quickstart

To use the `cli.sh` in bridging ether to or from the mev-commit chain, first make the script executable:

```bash
chmod +x cli.sh
```

Optionally, move the script to a folder in your `PATH` similar to:

```bash
sudo mv cli.sh /usr/local/bin/bridge-cli
```

Next we'll initialize bridge client parameters. Note all following commands display confirmation prompts. Use [primev docs](https://docs.primev.xyz/mev-commit-chain) to obtain relevant arguments. Router arguments are addresses of deployed hyperlane router contracts for each chain. Executing this command will save a `.bridge_config` json in the working directory:

```bash
bridge-cli init <L1 Router> <mev-commit chain Router> <L1 Chain ID> <mev-commit chain ID> <L1 URL> <MEV-Commit URL>
```

Once initialized, bridge ether to the mev-commit chain with

```bash
bridge-cli bridge-to-mev-commit <amount in wei> <dest_addr> <private_key>
```

Remember to bridge enough ether such that fees to bridge back to L1 can be paid! Bridge ether back to L1 with

```bash
bridge-cli bridge-to-l1 <amount in wei> <dest_addr> <private_key>
```

Note support for keystore and hardware wallets will be added later.

## Versions

Agents are built from https://github.com/primevprotocol/hyperlane-monorepo, using the rust/build.sh script from that repo. Docker images for the agents are optimized for amd64 architecture, and may need to be compiled on a powerful machine.

Hyperlane contracts and CLI are built from custom fork of their monorepo https://github.com/primevprotocol/hyperlane-monorepo.

## Contract deployer

Address:    `0x82b941824b43F33e417be1E92273A3020a0D760c`

Note if the relayer is emitting errors related to unexpected contract routing, try redeploying the hyperlane contracts using a new key pair. It's likely the current deployments are clashing with previous deployments on Sepolia.

To properly set a new hyperlane deployer:
* Generate a new key pair (ex: `cast wallet new`)
* Send or [mine](https://sepolia-faucet.pk910.de/) some Sepolia ETH to `Address`
* replace `Address` above for book keeping
* replace `CONTRACT_DEPLOYER_PRIVATE_KEY` in `.env`
* allocate funds to `Address` in the allocs field of `genesis.json`

Note the deployer of [primev contracts](https://github.com/primevprotocol/contracts) can be a separate account.

## Validator Accounts (same keys as POA signers)

### Node1

Address:     `0xd9cd8E5DE6d55f796D980B818D350C0746C25b97`

### Node2

Address:     `0x788EBABe5c3dD422Ef92Ca6714A69e2eabcE1Ee4`

## Relayer

Address:     `0x0DCaa27B9E4Db92F820189345792f8eC5Ef148F6`

## User emulators

There are 5 emulator services that simulate EOA's bridging to/from the mev-commit chain. Use the Makefile to start them. 

Note all these accounts must be funded with Sepolia ether and enough mev-commit chain ether to pay for gas.

Emulator 1 Address: `0x04F713A0b687c84D4F66aCd1423712Af6F852B78`
Emulator 2 Address: `0x4E2D04c65C399Eb27B3E3ADA06110BCd47b5a506`
Emulator 3 Address: `0x7AEe7AD6b2EAd96532D84D20358Db0e697f060Cd`
Emulator 4 Address: `0x765235CDda5FC6a620Fea2208A333a97CEDA2E1d`
Emulator 5 Address: `0x163c7bD4C3B815B06503D8E8B5906519C319EA6f`

## Starter .env file
To get a standard starter .env file from primev internal development, [click here.](https://www.notion.so/Private-keys-and-env-for-settlement-layer-245a4f3f4fe040a7b72a6be91131d9c2?pvs=4). Note this repo is being actively developed and required .env variables may change.

Environment variables that must be specified in the .env file or as command line arguments:

```
HYPERLANE_DEPLOYER_PRIVATE_KEY=0xpk1
NODE1_PRIVATE_KEY=0xpk2
NODE2_PRIVATE_KEY=0xpk3
RELAYER_PRIVATE_KEY=0xpk4
NEXT_PUBLIC_WALLET_CONNECT_ID=0xcId
AGENT_BASE_IMAGE=image
SETTLEMENT_RPC_URL=https://url
PUBLIC_SETTLEMENT_RPC_URL=https://url
SEPOLIA_RPC_URL=https://url
DD_API_KEY=432
DD_APP_KEY=808
EMULATOR1_PRIVATE_KEY=0xpk5
EMULATOR2_PRIVATE_KEY=0xpk6
EMULATOR3_PRIVATE_KEY=0xpk7
EMULATOR4_PRIVATE_KEY=0xpk8
EMULATOR5_PRIVATE_KEY=0xpk9
```
