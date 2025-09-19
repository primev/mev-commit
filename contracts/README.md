# Contracts

## Changelogs

* Add a row to **Current Deployments** when a contract is first deployed from a release branch and tagged commit. Note the proxy address, release branch, and initial tag should not change.
* After any on-chain upgrade, append a row to **Upgrade History**. The `Upgrade Tag` contains the most up-to-date source code for that contract.
* If an existing contract is ever redeployed (ex. because its required changes are not possible via contract upgrade), simply replace the appropriate row in **Current Deployments** with the new proxy address, release branch, and initial tag.

## Mainnet (L1) Contract Changelog

This changelog tracks deployments of **Ethereum mainnet** contracts. This changelog is only valid from the `main` branch.

### Current Deployments

| Contract              | Proxy Address                                | Release Branch      | Initial Tag           |
|-----------------------|----------------------------------------------|---------------------|--------------------------|
| ValidatorOptInRouter  | `0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64` | `release/v1.1.x`    | `v1.1.0-contracts`      |
| VanillaRegistry       | `0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9` | `release/v1.1.x`    | `v1.1.0-contracts`      |
| MevCommitAVS          | `0xBc77233855e3274E1903771675Eb71E602D9DC2e` | `release/v1.1.x`    | `v1.1.0-contracts`      |
| L1Gateway             | `0xDBf24cafF1470a6D08bF2FF2c6875bafC60Cf881` | `release/v1.1.x`    | `v1.1.0-contracts`      |
| MevCommitMiddleware   | `0x21fD239311B050bbeE7F32850d99ADc224761382` | `release/v1.1.x`    | `v1.1.0-contracts`         |

### Upgrade History

| Timestamp (UTC)             | Contract            | New Impl Version      | Upgrade Tag       |
|-----------------------------|---------------------|-----------------------|-------------------|
| Mar-12-2025 03:33:35 AM UTC | MevCommitMiddleware | MevCommitMiddlewareV2 | v1.1.0-middleware |

## Core Contract Changelog

This changelog tracks "core" contract deployments on the mev-commit chain. This changelog is only valid from the `main` branch.

### Current Deployments

| Contract              | Proxy Address                                | Release Branch      | Initial Tag           |
|-----------------------|----------------------------------------------|---------------------|--------------------------|
| BidderRegistry        | `0xC973D09e51A20C9Ab0214c439e4B34Dbac52AD67` | `release/v1.1.x`    | `v1.1.0`      |
| ProviderRegistry      | `0xb772Add4718E5BD6Fe57Fb486A6f7f008E52167E` | `release/v1.1.x`    | `v1.1.0`      |
| PreconfManager        | `0x3761bF3932cD22d684A7485002E1424c3aCCD69c` | `release/v1.1.x`    | `v1.1.0`      |
| Oracle                | `0xa1aaCA1e4583dB498D47f3D5901f2B2EB49Bd8f6` | `release/v1.1.x`    | `v1.1.0`      |
| BlockTracker          | `0x0DA2a367C51f2a34465ACd6AE5d8A48385e9cB03` | `release/v1.1.x`    | `v1.1.0`      |

### Upgrade History

| Timestamp (UTC)             | Contract            | New Impl Version      | Upgrade Tag       |
|-----------------------------|---------------------|-----------------------|-------------------|
| April 7th 2025 | Oracle | OracleV2 | No tag, see commit `bc4ebddd70f23d58ba6f9b2e8701e7f45d89cf82` in `release/v1.1.x`. |
| July 22nd 2025, 20:10:39 UTC | BidderRegistry | BidderRegistryV2 | `v1.1.5` in `release/v1.1.x`. |
| July 22nd 2025, 20:10:39 UTC | ProviderRegistry | ProviderRegistryV2 | `v1.1.5` in `release/v1.1.x`. |

## Hoodi Testnet (L1) Contract Changelog

This changelog tracks deployments of **Hoodi Testnet** contracts. This changelog is only valid from the `main` branch.

### Current Deployments

| Contract              | Proxy Address                                |  Initial Commit           |
|-----------------------|----------------------------------------------|---------------------|
| ValidatorOptInHub     | `0x953c2a669493A126fd50E9f56306f254B4e49709` | `35664894728008c34cf6e24cbe77ce91f091144b` in 'optinrouter-new-registry-support'      |
| ValidatorOptInRouter  | `0xa380ba6d6083a4Cb2a3B62b0a81Ea8727861c13e` | `13cf068477e6efdbb5c4fe5ce53a11af30bf8b47` in 'main'      |
| VanillaRegistry       | `0x536f0792c5d5ed592e67a9260606c85f59c312f0` | `13cf068477e6efdbb5c4fe5ce53a11af30bf8b47` in 'main'      |
| MevCommitAVS          | `0xdF8649d298ad05f019eE4AdBD6210867B8AB225F` | `13cf068477e6efdbb5c4fe5ce53a11af30bf8b47` in 'main'      |
| MevCommitMiddleware   | `0x8E847EC4a36c8332652aB3b2B7D5c54dE29c7fde` | `13cf068477e6efdbb5c4fe5ce53a11af30bf8b47` in 'main'      |


### Upgrade History

| Timestamp (UTC)             | Contract            | New Impl Version      | Commmit       |
|-----------------------------|---------------------|-----------------------|-------------------|
| N/A |  |  |  |


## L1 Deployer CLI

> **After completing any L1 deployment, immediately record it in the “Current Deployments” table above.**

The `l1-deployer-cli.sh` enables production deployment of L1 contracts, with publishing of source code to etherscan (see [source code verification](https://info.etherscan.com/how-to-verify-contracts/)). This deployment workflow is decoupled from the core mev-commit chain contracts. This cli accepts keystore (not suggested), ledger, or trezor wallets. L1 contracts must be deployed from a tagged commit, that's part of a release branch. For mainnet it's not recommended to use the `--skip-release-verification` flag.

If contract deployment succeeds but etherscan verification fails, try running [forge verify-contract](https://book.getfoundry.sh/reference/forge/forge-verify-contract) directly with the deployment address. Ex:

```bash
forge verify-contract --watch --rpc-url $RPC_URL 0x4c31ad10617bb36e7749c686eedf6fef0fd2502e ValidatorOptInRouter
```

To avoid issues with etherscan verification, use a non-public RPC that can support rapid requests.

### Dependencies

- [Foundry suite](https://book.getfoundry.sh/getting-started/installation)
- [git](https://git-scm.com/downloads)
- [curl](https://everything.curl.dev/install/linux.html)
- [jq](https://stedolan.github.io/jq/download/)

## Contract Upgrades

> **After completing any upgrade, immediately record it in the “Upgrade History” table above.**

Contract upgrades are not always possible, as there are [strict limitations as enforced by Solidity](https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#modifying-your-contracts). When a contract feat/fix cannot be implemented as a contract upgrade, simply PR the changes into main, and release/deploy a new contract instance as needed.

See [#360](https://github.com/primev/mev-commit/pull/360) for a reference example of a complete contract upgrade.

### Implementing the feat/fix

The following instructions assist in upgrading an existing deployment of any of the contracts within this directory. Contracts committed to this repo should utilize [UUPS proxies](https://docs.openzeppelin.com/contracts/4.x/api/proxy#UUPSUpgradeable). 

First, implement and merge an appropriate feat/fix to the implementation contract using the main branch. This ensures any future (new) deployments of the contract will include the feat/fix along with changes to the contract's ABI.

Next, create a branch off the appropriate release branch (see [Current Deployments](#current-deployments)), we'll refer to this as the "upgrade branch". Copy/paste the currently deployed contract implementation to a new file (currently deployed implementation and new file should both reside within the upgrade branch). Name the new file with an incremented version number as a postfix of the original contract's filename. E.g. if the original contract's filename is `MevCommitAVS.sol`, then the new contract's filename should be `MevCommitAVSV2.sol`. 

Now update the new implementation contract with the feat/fix that was merged to main. You can try to utilize cherry-picking here, but may have to re-implement the feat/fix manually.

Make sure to update and run tests on the upgrade branch. The `setUp()` function for relevant test contract(s) should include an upgrade from the previous to the new implementation contract. Also any regression tests merged to main should be ported to the upgrade branch.

If the feat/fix required changes to the storage contract (see limitations above), you **must** also create a new storage contract with an incremented version number as a postfix of the original storage contract's filename. E.g. if the original storage contract's filename is `MevCommitAVSStorage.sol`, then the new storage contract's filename should be `MevCommitAVSStorageV2.sol`.

**If applicable, make sure the incremented implementation contract inherits from the incremented storage contract. This is necessary to ensure accurate upgrade safety validation.**

Example: `MevCommitAVSV2.sol` would inherit `MevCommitAVSV2Storage.sol`.

Now [define the reference contract](https://docs.openzeppelin.com/upgrades-plugins/1.x/api-core#define-reference-contracts) for the upgrade right above the new contract implementation. E.g `/// @custom:oz-upgrades-from MevCommitAVS`.

Finally, build the contracts and use [openzeppelin's cli](https://docs.openzeppelin.com/upgrades-plugins/1.x/api-core#usage) to validate the upgrade, similar to `npx @openzeppelin/upgrades-core validate --contract MevCommitAVSV2`. If this command fails, you'll need to address whether a contract upgrade is still possible/appropriate.

### Note on ABI changes

If the feat/fix to the implementation contract results in a change to the contract's ABI, you'll need to generate new ABI file and go binding within the upgrade branch. This can be handled on a case-by-case basis, as generating (for example) a `MevCommitAVSV2.abi` in the upgrade branch, that is equivalent to the `MevCommitAVS.abi` in the main branch, could be difficult to maintain.

If possible it's recommended to avoid changing a contract's ABI for an upgrade.

### Deployment

Once your "upgrade branch" is reviewed and merged into the release branch, tag the latest commit from the release branch. This tag will be used to populate the **Upgrade History** table above.

Invoking the upgrade involves creating a script in which a new implementation contract is deployed, then calling `upgradeToAndCall` on the proxy contract, passing in the address of the new implementation contract.

See example below

```solidity
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract UpgradeAVS is Script {
    function _runTests(MevCommitAVS avs) internal {
        // Define expected state and/or function behavior to execute before and after upgrade
    }
}

contract UpgradeAnvil is UpgradeAVS {
    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        MevCommitAVS existingMevCommitAVSProxy = MevCommitAVS(payable(0x5FC8d32690cc91D4c39d9d3abcBD16989F875707));

        console.log("Pre-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);

        vm.startBroadcast();
        Upgrades.upgradeProxy(
            address(existingMevCommitAVSProxy), 
            "MevCommitAVSV2.sol", 
            "" 
        );
        console.log("Upgraded to MevCommitAVSV2");
        vm.stopBroadcast();

        console.log("Post-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);
    }
}
```

In this example, no function call is made during the upgrade. However [these examples](https://docs.openzeppelin.com/upgrades-plugins/1.x/foundry-upgrades#examples) demonstrate how to make a function call during the upgrade.

It's encouraged to test your upgrade process using anvil, then use identical code to invoke the upgrade on Holesky/mainnet.

### Note on multisig vs EOA

The aforementioned process can be followed exactly for contracts that are owned by a single EOA, where the forge script can be run directly using a keystore.

For contracts that are owned by a multisig, simply deploy an uninitialized new implementation contract from any account using a forge script etc.,

```solidity
import {MevCommitAVSV2} from "../../../contracts/validator-registry/avs/MevCommitAVSV2.sol";
import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";

contract DeployNewImpl is Script {
    function run() external {
        vm.startBroadcast();
        MevCommitAVSV2 newImplementation = new MevCommitAVSV2();
        console.log("Deployed new implementation contract at address: ", address(newImplementation));
        vm.stopBroadcast();
    }
}
```

Then call the `upgradeToAndCall(newImplAddr, callData)` function directly on the proxy contract using your multisig UI (e.g Safe wallet).

Ownership of the implementation contract itself would be irrelevant in this scenario, as the implementation is deployed without calling its initializer, or setting any state. Ie. the implementation contract serves only as a blueprint for state transition functionality.

The multisig option bypasses safety checks that would otherwise happen by using a forge script in tandem with [OpenZeppelin Foundry Upgrades](https://github.com/OpenZeppelin/openzeppelin-foundry-upgrades). Therefore it's very important to use the `validate` command from above to ensure the upgrade is safe to proceed with.
