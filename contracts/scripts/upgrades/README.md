# Generic Upgrade Scripts

## Overview

This directory contains generic upgrade scripts for upgrading proxy contracts. These scripts are designed to be used exclusively through the `l1-upgrade-cli.sh` command-line interface.

There are two main scripts:

- **`GenericUpgrade.s.sol`** - For direct upgrades (EOA-owned contracts)
- **`GenericMultisigUpgrade.s.sol`** - For multisig upgrades (deploys implementation only)

Both scripts support multiple chains and automatically enforce chain ID checks for safety.

## Usage

### Via CLI (Recommended)

Both scripts are invoked through `l1-upgrade-cli.sh`:

**Direct Upgrade (EOA-owned contracts):**
```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --proxy-address 0x1234... \
  --chain mainnet \
  --keystore
```

**Multisig Upgrade (multisig-owned contracts):**
```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --chain mainnet \
  --multisig \
  --keystore
```

The CLI script will:
1. Validate the upgrade safety using `validate-upgrade.sh`
2. Find the contract files automatically
3. Set the required environment variables
4. Call the appropriate contract variant from the selected script
5. Execute the upgrade/deployment using forge

<details>
<summary><strong>GenericUpgrade.s.sol - Direct Upgrades</strong></summary>

### Overview

`GenericUpgrade.s.sol` performs complete proxy upgrades for contracts owned by a single EOA (Externally Owned Account). It handles both deployment of the new implementation and the upgrade transaction in a single operation.

### When to Use

- ✅ Contract is owned by a single EOA wallet
- ✅ You want to perform the upgrade in one transaction
- ✅ You have access to the wallet that owns the proxy or proxy admin

### How It Works

1. **Environment Variables**: The script reads configuration from environment variables:
   - `PROXY_ADDRESS` - Address of the proxy to upgrade
   - `OLD_CONTRACT_NAME` - Name of the old contract (for logging)
   - `NEW_CONTRACT_NAME` - Name of the new contract (for logging)
   - `NEW_CONTRACT_PATH` - Filename of the new contract (e.g., `MevCommitAVSV2.sol`)

2. **Upgrade Process**: Uses OpenZeppelin's `Upgrades.upgradeProxy()` to:
   - Deploy the new implementation contract
   - Upgrade the proxy to point to the new implementation
   - Validate storage layout compatibility (if annotation is present)

3. **Contract Path Format**: The `NEW_CONTRACT_PATH` should be just the filename (e.g., `MevCommitAVSV2.sol`), not a full path. The library will find the artifact in `out/` automatically.

### Contract Variants

- **`UpgradeContract`** - Generic variant (works on any chain)
- **`UpgradeContractAnvil`** - Anvil local testing (chain ID 31337)
- **`UpgradeContractHolesky`** - Holesky testnet (chain ID 17000)
- **`UpgradeContractHoodi`** - Hoodi testnet (chain ID 560048)
- **`UpgradeContractMainnet`** - Mainnet (chain ID 1)

### Direct Usage (Not Recommended)

While you can call the script directly with forge, it's not recommended as you'll need to manually set all environment variables:

```bash
# Required environment variables
export PROXY_ADDRESS="0x1234..."
export OLD_CONTRACT_NAME="MevCommitAVS"
export NEW_CONTRACT_NAME="MevCommitAVSV2"
export NEW_CONTRACT_PATH="MevCommitAVSV2.sol"
export RPC_URL="..."
export SENDER="..."

forge script scripts/upgrades/GenericUpgrade.s.sol:UpgradeContractMainnet \
  --rpc-url $RPC_URL \
  --sender $SENDER \
  --broadcast
```

### Examples

**Anvil (Local Testing):**
```bash
export RPC_URL="http://127.0.0.1:8545"

./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --proxy-address 0xc5a5c42992decbae36851359345fe25997f5c42d \
  --chain anvil \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --skip-validation  # Optional: skip validation for faster local testing
```

**Mainnet:**
```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --proxy-address 0x1234567890123456789012345678901234567890 \
  --chain mainnet \
  --keystore \
  --priority-gas-price 2000000000 \
  --with-gas-price 5000000000
```

**Holesky Testnet:**
```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --proxy-address 0x5678901234567890123456789012345678901234 \
  --chain holesky \
  --ledger
```

</details>

<details>
<summary><strong>GenericMultisigUpgrade.s.sol - Multisig Upgrades</strong></summary>

### Overview

`GenericMultisigUpgrade.s.sol` deploys **only** the implementation contract for multisig-owned proxy upgrades. After deployment, you must use your multisig UI (e.g., Safe wallet) to call `upgradeToAndCall()` on the proxy contract.

### When to Use

- ✅ Contract is owned by a multisig wallet
- ✅ You need to deploy the implementation separately before multisig approval
- ✅ You want to prepare the implementation address for multisig transaction submission

### How It Works

1. **Environment Variables**: The script reads configuration from environment variables:
   - `NEW_CONTRACT_NAME` - Name of the new contract (e.g., "MevCommitAVSV2")
   - `NEW_CONTRACT_PATH` - Path to the new contract file (e.g., "MevCommitAVSV2.sol") - used for logging only

2. **Deployment Process**: Uses `vm.deployCode()` to:
   - Deploy the new implementation contract bytecode
   - **Does NOT** call any constructor or initializer
   - **Does NOT** upgrade the proxy (that's done via multisig)

3. **Output**: The script outputs the implementation address, which you then use in your multisig UI to call `upgradeToAndCall(implementation, callData)`.

### Contract Variants

- **`DeployMultisigImpl`** - Generic variant (works on any chain)
- **`DeployMultisigImplAnvil`** - Anvil local testing (chain ID 31337)
- **`DeployMultisigImplHolesky`** - Holesky testnet (chain ID 17000)
- **`DeployMultisigImplHoodi`** - Hoodi testnet (chain ID 560048)
- **`DeployMultisigImplMainnet`** - Mainnet (chain ID 1)

### Direct Usage (Not Recommended)

While you can call the script directly with forge, it's not recommended as you'll need to manually set all environment variables:

```bash
# Required environment variables
export NEW_CONTRACT_NAME="MevCommitAVSV2"
export NEW_CONTRACT_PATH="MevCommitAVSV2.sol"
export RPC_URL="..."
export SENDER="..."

forge script scripts/upgrades/GenericMultisigUpgrade.s.sol:DeployMultisigImplMainnet \
  --rpc-url $RPC_URL \
  --sender $SENDER \
  --broadcast \
  --verify
```

### Examples

**Anvil (Local Testing):**
```bash
export RPC_URL="http://127.0.0.1:8545"

./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --chain anvil \
  --multisig \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --skip-validation  # Optional: skip validation for faster local testing
```

**Mainnet:**
```bash
./l1-upgrade-cli.sh upgrade \
  --old-contract MevCommitAVS \
  --new-contract MevCommitAVSV2 \
  --chain mainnet \
  --multisig \
  --keystore
```

**After Deployment:**

Once the implementation is deployed, use your multisig UI to complete the upgrade:

1. Copy the implementation address from the script output
2. In your multisig UI (e.g., Safe wallet), create a new transaction to the proxy contract
3. Call function: `upgradeToAndCall(implementation, callData)`
   - `implementation`: The address shown in the script output
   - `callData`: Empty bytes (`"0x"`) if no initialization needed, or encoded function call if needed
4. Submit the transaction through your multisig workflow

### Important Notes

- The deployer address can be any account (doesn't need to be the multisig)
- The implementation contract serves only as a blueprint and has no ownership
- **Manual validation required**: The multisig workflow bypasses automatic safety checks. Always manually validate using `validate-upgrade.sh` before deploying

</details>

## Prerequisites

Before running any upgrade:

1. **Validation**: Always validate the upgrade using `validate-upgrade.sh`:
   ```bash
   ./validate-upgrade.sh --contract MevCommitAVSV2 --reference MevCommitAVS
   ```

2. **Contract Annotation**: Ensure the new contract has the proper annotation:
   ```solidity
   /// @custom:oz-upgrades-from MevCommitAVS
   contract MevCommitAVSV2 is ...
   ```

3. **Build**: Contracts must be compiled:
   ```bash
   forge clean && forge build
   ```

## Important Notes

- **Always use the CLI**: The `l1-upgrade-cli.sh` script handles all the setup, validation, and path construction automatically
- **Never skip validation on mainnet**: Always run validation before upgrading on mainnet
- **Contract filename only**: The contract path should be just the filename (e.g., `MevCommitAVSV2.sol`), the CLI extracts this automatically
- **Chain-specific variants**: Each chain has its own contract variant that enforces chain ID checks for safety
- **Multisig validation**: For multisig upgrades, manual validation is critical since automatic checks are bypassed

<details>
<summary><strong>Upgrade Checklist</strong></summary>

> **After completing any upgrade, immediately record it in the "Upgrade History" table in the main README.md.**

### Phase 1: Implementing the Feat/Fix

**Initial Implementation:**
- [ ] Feat/fix implemented and merged to `main` branch
- [ ] Upgrade branch created from appropriate release branch (see Current Deployments in README.md)
- [ ] Currently deployed contract implementation copied to new versioned file (e.g., `MevCommitAVSV2.sol`)
- [ ] New implementation contract updated with feat/fix from main (cherry-pick or manual re-implementation)

**Storage Contract Versioning (if applicable):**
- [ ] Storage contract changes identified
- [ ] New storage contract created with incremented version (e.g., `MevCommitAVSStorageV2.sol`)
- [ ] New implementation contract inherits from new storage contract
- [ ] Old storage contract moved to `contracts/upgrades/[feature-folder]/`

**Contract Organization:**
- [ ] Old contract moved to `contracts/upgrades/[feature-folder]/`
- [ ] `/// @custom:oz-upgrades-from` annotation added to new contract
- [ ] Reference contract properly defined (e.g., `/// @custom:oz-upgrades-from MevCommitAVS`)

**Testing:**
- [ ] Tests updated on upgrade branch
- [ ] `setUp()` function includes upgrade from previous to new implementation
- [ ] Regression tests from main ported to upgrade branch
- [ ] All tests passing on upgrade branch

**Validation & Safety (CRITICAL):**
- [ ] Contracts built: `forge clean && forge build`
- [ ] ✅ **UPGRADE VALIDATION PASSED**: `npx @openzeppelin/upgrades-core validate --contract ContractV2 --reference ContractV1` - **DO NOT PROCEED WITHOUT THIS**
- [ ] Validation failures addressed or upgrade deemed not possible/appropriate

**ABI Changes (if applicable):**
- [ ] ABI changes identified and documented
- [ ] New ABI file and go bindings generated (if needed)
- [ ] Considered avoiding ABI changes if possible

### Phase 2: Deployment Preparation

**Branch Management:**
- [ ] Upgrade branch reviewed and merged into release branch
- [ ] Latest commit from release branch tagged (tag will populate Upgrade History table)
- [ ] Tag verified and documented

**Environment Setup:**
- [ ] `RPC_URL` environment variable set for target chain
- [ ] Wallet type selected (keystore/ledger/trezor/private-key)

**For Keystore Wallet:**
- [ ] `KEYSTORES` environment variable set
- [ ] `KEYSTORE_PASSWORD` environment variable set
- [ ] `SENDER` environment variable set
- [ ] `ETHERSCAN_API_KEY` set (optional, for verification)

**For Ledger/Trezor Wallet:**
- [ ] `HD_PATHS` environment variable set
- [ ] `SENDER` environment variable set
- [ ] Hardware wallet connected and ready

**For Private Key (Anvil/Local Testing):**
- [ ] Private key prepared (or using default anvil key)
- [ ] `SENDER` set (optional)

**Deployment Parameters:**
- [ ] Proxy address confirmed (from Current Deployments table)
- [ ] Old contract name confirmed
- [ ] New contract name confirmed
- [ ] Chain selected (`mainnet`, `holesky`, `hoodi`, or `anvil`)
- [ ] Gas price parameters set (if needed): `--priority-gas-price` and/or `--with-gas-price`
- [ ] Backup and rollback plan documented

### Phase 3: Testing on Anvil (Recommended)

**Local Testing:**
- [ ] Anvil chain running locally
- [ ] Local upgrade tested using `l1-upgrade-cli.sh`:
  ```bash
  ./l1-upgrade-cli.sh upgrade \
    --old-contract ContractV1 \
    --new-contract ContractV2 \
    --proxy-address <anvil_proxy_address> \
    --chain anvil \
    --private-key <key> \
    --skip-validation  # Optional for faster local testing
  ```
- [ ] Pre-upgrade state verified
- [ ] Post-upgrade state verified
- [ ] Functionality tested after upgrade

### Phase 4: Production Deployment

**For EOA-owned Contracts (using l1-upgrade-cli.sh):**
- [ ] Production environment variables verified
- [ ] RPC URL chain ID matches target chain
- [ ] Upgrade command prepared:
  ```bash
  ./l1-upgrade-cli.sh upgrade \
    --old-contract ContractV1 \
    --new-contract ContractV2 \
    --proxy-address <proxy_address> \
    --chain <chain> \
    <wallet_option> \
    [--priority-gas-price <price>] \
    [--with-gas-price <price>]
  ```
- [ ] Validation will run automatically (or `--skip-validation` explicitly approved for non-mainnet)
- [ ] Upgrade executed successfully
- [ ] Upgrade transaction verified on block explorer
- [ ] Contract verification completed (if ETHERSCAN_API_KEY provided)

**For Multisig-owned Contracts:**
- [ ] ✅ **UPGRADE VALIDATION PASSED**: `npx @openzeppelin/upgrades-core validate --contract ContractV2 --reference ContractV1` - **MANUAL VALIDATION REQUIRED**
- [ ] Implementation deployment command prepared:
  ```bash
  ./l1-upgrade-cli.sh upgrade \
    --old-contract ContractV1 \
    --new-contract ContractV2 \
    --chain <chain> \
    --multisig \
    <wallet_option>
  ```
- [ ] New implementation contract deployed
- [ ] Implementation address documented
- [ ] Multisig UI prepared with `upgradeToAndCall(newImplAddr, callData)` transaction
- [ ] All multisig signers informed and ready
- [ ] Upgrade transaction submitted via multisig
- [ ] Multisig transaction executed successfully

### Phase 5: Post-Deployment

**Documentation:**
- [ ] Upgrade recorded in "Upgrade History" table in README.md
- [ ] Timestamp (UTC) recorded
- [ ] Upgrade tag recorded
- [ ] Any notes or observations documented

**Verification:**
- [ ] Proxy address still points to correct contract
- [ ] New implementation contract address verified
- [ ] Contract functions tested post-upgrade
- [ ] Monitoring and alerts verified (if applicable)

</details>

## Related Documentation

- [l1-upgrade-cli.sh](../l1-upgrade-cli.sh) - The CLI script that uses these upgrade scripts
- [Contract Upgrade Structure](../../contracts/upgrades/docs/CONTRACT_UPGRADE_STRUCTURE.md) - Documentation on contract versioning structure
- [validate-upgrade.sh](../../validate-upgrade.sh) - Validation script for upgrade safety
