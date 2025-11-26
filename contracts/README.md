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
| L1Gateway             | `0x5d64B933739558101F9359E2750ACC228f0CB64F` | `release/v1.2.x`    | `v1.2.0`                |
| MevCommitMiddleware   | `0x21fD239311B050bbeE7F32850d99ADc224761382` | `release/v1.1.x`    | `v1.1.0-contracts`         |
| ValidatorOptInHub     | `0x1059155bD5854224bF58e43CD3EEC6B07b4F30Ad` | `release/v1.2.x`    | `v1.2.3`         |
| BlockRewardsManager   | `0x531D19cFEB3e685745DA1e1897485E9F351e7Aa0` | `release/v1.2.x`    | `v1.2.3`         |
| RewardDistributor     | `0xccf10d9911c2e1b2d589cbc8a3878d1b138aa2c2` | `release/v1.2.x`    | `v1.2.3`         |


### Upgrade History

| Timestamp (UTC)             | Contract            | New Impl Version      | Upgrade Tag       |
|-----------------------------|---------------------|-----------------------|-------------------|
| Mar-12-2025 03:33:35 AM UTC | MevCommitMiddleware | MevCommitMiddlewareV2 | v1.1.0-middleware |

## Mev-commit Chain Contract Changelog

This changelog tracks contract deployments on the mev-commit chain. This changelog is only valid from the `main` branch.

### Current Deployments

| Contract              | Proxy Address                                | Release Branch      | Initial Tag           |
|-----------------------|----------------------------------------------|---------------------|--------------------------|
| BidderRegistry        | `0x145a9f4cbae2ec281f417195ea3464dbd04289a2` | `release/v1.2.x`    | `v1.2.0`      |
| ProviderRegistry      | `0xeb6d22309062a86fa194520344530874221ef48c` | `release/v1.2.x`    | `v1.2.0`      |
| PreconfManager        | `0x2ee9e88f57a7db801e114a4df7a99eb7257871e2` | `release/v1.2.x`    | `v1.2.0`      |
| Oracle                | `0x37a037d2423221f403cfa146f5fb962e19582d90` | `release/v1.2.x`    | `v1.2.0`      |
| BlockTracker          | `0x5d64b933739558101f9359e2750acc228f0cb64f` | `release/v1.2.x`    | `v1.2.0`      |
| SettlementGateway     | `0x21f5f1142200a515248a2eef5b0654581c7f2b46` | `release/v1.2.x`    | `v1.2.0`      |



### Upgrade history related to current deployments

| Timestamp (UTC)             | Contract            | New Impl Version      | Upgrade Tag       |
|-----------------------------|---------------------|-----------------------|-------------------|
| Oct-06-2025 10:30:00 PM UTC | ProviderRegistry    | ProviderRegistryV2    | v1.2.1            |

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
| RocketMinipoolRegistry| `0x30d478b02918c4b11731efc4868d848d551e79b2` | `d994896a8b17b131ebda8bbae1a43ddd55c906c1` in 'rocketpool-registry'      |
| Reputational VanillaRegistry| `0xf65fea786014a2c4442df57a086298d76bd23abe` | `9ffa6c82f1e516f89c321cd576c3d777f0d5d261` in 'main'      |


### Upgrade History

| Timestamp (UTC)             | Contract            | New Impl Version      | Commmit       |
|-----------------------------|---------------------|-----------------------|-------------------|
| N/A |  |  |  |


## Hoodi ValidatorOptInHub Registry Indexes

| Registry Index              | Registry Type                                |  Proxy Address           |
|-----------------------|----------------------------------------------|---------------------|
| 0                           | Vanilla Registry         | `0x536F0792c5D5Ed592e67a9260606c85F59C312F0` |
| 1                           | MevCommitAVS             | `0xdF8649d298ad05f019eE4AdBD6210867B8AB225F` |
| 2                           | MevCommitMiddleware      | `0x8E847EC4a36c8332652aB3b2B7D5c54dE29c7fde` |
| 3                           | Vanilla Lido             | `0xEfd6333907fc73c1ac3167D843488B8899bac91b` |
| 4                           | RocketMinipoolRegistry   | `0x82018E031cE9B4040752544BC63086b6FECFB07B` |
| 5                           | Vanilla Reputational     | `0x1f6391d28d062acb8edb96ec94752045ca4815e5` |


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

Contract upgrades are not always possible, as there are [strict limitations as enforced by Solidity](https://docs.openzeppelin.com/upgrades-plugins/writing-upgradeable#modifying-your-contracts). When a contract feat/fix cannot be implemented as a contract upgrade, simply PR the changes into main, and release/deploy a new contract instance as needed.

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

Now [define the reference contract](https://docs.openzeppelin.com/upgrades-plugins/api-core#define-reference-contracts) for the upgrade right above the new contract implementation. E.g `/// @custom:oz-upgrades-from MevCommitAVS`.

Finally, build the contracts and validate the upgrade.

```bash
forge clean && forge build
npx @openzeppelin/upgrades-core validate --contract MevCommitAVSV2 --reference MevCommitAVS
```

If validation fails, you'll need to address whether a contract upgrade is still possible/appropriate.

**Note:** The `l1-upgrade-cli.sh` script automatically runs validation before executing the upgrade (unless `--skip-validation` is used), so you don't need to run it manually if using the CLI. However, it's recommended to validate during development to catch issues early.

### Note on ABI changes

If the feat/fix to the implementation contract results in a change to the contract's ABI, you'll need to generate new ABI file and go binding within the upgrade branch. This can be handled on a case-by-case basis, as generating (for example) a `MevCommitAVSV2.abi` in the upgrade branch, that is equivalent to the `MevCommitAVS.abi` in the main branch, could be difficult to maintain.

If possible it's recommended to avoid changing a contract's ABI for an upgrade.

### Deployment

Once your "upgrade branch" is reviewed and merged into the release branch, tag the latest commit from the release branch. This tag will be used to populate the **Upgrade History** table above.

The upgrade is performed via the [l1-upgrade-cli.sh](./l1-upgrade-cli.sh) script. This CLI tool automates the upgrade process, handles validation, and supports multiple wallet types and chains.

#### Using l1-upgrade-cli.sh (Recommended)

The `l1-upgrade-cli.sh` script provides a streamlined way to upgrade contracts. It automatically:
- Validates upgrade safety (unless `--skip-validation` is used)
- Finds contract files automatically
- Handles chain-specific configurations
- Supports multiple wallet types (keystore, ledger, trezor, private-key)

**Basic Usage:**

```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --proxy-address 0x1234... \
  --chain mainnet \
  --keystore
```



#### Manual Upgrade Scripts (Advanced)

For custom upgrade logic or multisig scenarios, you can create manual upgrade scripts. See example below:

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

In this example, no function call is made during the upgrade. However [these examples](https://docs.openzeppelin.com/upgrades-plugins/foundry/foundry-upgrades#examples) demonstrate how to make a function call during the upgrade.

It's encouraged to test your upgrade process using anvil, then use the `l1-upgrade-cli.sh` to invoke the upgrade on Holesky/mainnet.

### Note on multisig vs EOA

**For EOA-owned contracts:**

The `l1-upgrade-cli.sh` script handles the complete upgrade process for contracts owned by a single EOA. You can use keystore, ledger, trezor, or private-key wallet options as described above.

**For multisig-owned contracts:**

Contracts owned by a multisig require a different approach since the multisig must approve the upgrade transaction. The process is:

1. **Validate the upgrade** (critical step):
   ```bash
   ./validate-upgrade.sh --contract MevCommitAVSV2 --reference MevCommitAVS
   ```

2. **Deploy the new implementation contract** from any account. You can use a simple forge script:

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

3. **Call `upgradeToAndCall(newImplAddr, callData)`** directly on the proxy contract using your multisig UI (e.g., Safe wallet).

Ownership of the implementation contract itself is irrelevant in this scenario, as the implementation is deployed without calling its initializer or setting any state. The implementation contract serves only as a blueprint for state transition functionality.

**Important:** The multisig workflow bypasses the automatic safety checks that `l1-upgrade-cli.sh` performs. Therefore, it's essential to manually validate the upgrade using `validate-upgrade.sh` before proceeding. See the [validation section](#implementing-the-featfix) above for details.
