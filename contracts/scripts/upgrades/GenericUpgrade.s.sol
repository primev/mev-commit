// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {Options} from "openzeppelin-foundry-upgrades/Options.sol";

/**
 * @title GenericUpgrade
 * @notice Generic upgrade script that can upgrade any contract.
 * 
 * This script reads configuration from environment variables:
 *  - OLD_CONTRACT_NAME: Name of the old contract (for logging)
 *  - NEW_CONTRACT_NAME: Name of the new contract (for logging)
 *  - NEW_CONTRACT_PATH: Path to the new contract file (e.g., "validator-registry/avs/MevCommitAVSV2.sol")
 *  - PROXY_ADDRESS: Address of the proxy contract to upgrade
 * 
 * ⚠️  IMPORTANT: Run validation before running this script!
 * ./validate-upgrade.sh --contract <NEW_CONTRACT_NAME> --reference <OLD_CONTRACT_NAME>
 * 
 * Usage:
 *  forge script scripts/upgrades/GenericUpgrade.s.sol:UpgradeContract --rpc-url <RPC_URL> --sender <OWNER_ADDRESS> --broadcast --verify -vvvv
 * 
 * IMPORTANT: Always use --sender <OWNER_ADDRESS> flag with the address that owns the proxy or proxy admin.
 */
contract GenericUpgrade is Script {
    // ANSI color codes
    string private constant _RESET = "\x1b[0m";
    string private constant _GREEN = "\x1b[32m";
    string private constant _BRIGHT_GREEN = "\x1b[92m";
    string private constant _YELLOW = "\x1b[33m";
    string private constant _CYAN = "\x1b[36m";
    string private constant _BRIGHT_CYAN = "\x1b[96m";
    string private constant _BLUE = "\x1b[34m";

    function run() public virtual {
        vm.startBroadcast();

        // Get proxy address from environment variable
        address proxyAddress = vm.envOr("PROXY_ADDRESS", address(0));
        
        if (proxyAddress == address(0)) {
            revert("PROXY_ADDRESS must be set");
        }

        // Get contract names and path from environment variables
        string memory oldContractName = vm.envOr("OLD_CONTRACT_NAME", string("OldContract"));
        string memory newContractName = vm.envOr("NEW_CONTRACT_NAME", string("NewContract"));
        string memory newContractPath = vm.envOr("NEW_CONTRACT_PATH", string(""));
        
        if (bytes(newContractPath).length == 0) {
            revert("NEW_CONTRACT_PATH must be set");
        }

        console.log(string.concat(_CYAN, "Upgrading contract on chain:", _RESET), block.chainid);
        console.log(string.concat(_CYAN, "Old contract:", _RESET), oldContractName);
        console.log(string.concat(_CYAN, "New contract:", _RESET), newContractName);
        console.log(string.concat(_CYAN, "New contract path:", _RESET), newContractPath);
        console.log(string.concat(_CYAN, "Proxy address:", _RESET), proxyAddress);
        console.log(string.concat(_CYAN, "Upgrader address:", _RESET), msg.sender);

        // Validate upgrade safety with OpenZeppelin Foundry Upgrades
        if (bytes(oldContractName).length > 0) {
            console.log("");
            console.log(string.concat(_YELLOW, "Validating upgrade safety with OpenZeppelin Foundry Upgrades...", _RESET));
            Options memory opts;
            // Construct reference contract identifier (just filename, OpenZeppelin will resolve it)
            string memory referenceIdentifier = string.concat(oldContractName, ".sol");
            opts.referenceContract = referenceIdentifier;
            Upgrades.validateUpgrade(newContractPath, opts);
            console.log(string.concat(_BRIGHT_GREEN, "[PASS] OpenZeppelin Foundry Upgrades verification passed!", _RESET));
            console.log("");
        }

        // Upgrade to new contract
        // No function call during upgrade by default
        // If you need to call a function during upgrade, modify this script or use a custom upgrade script
        Upgrades.upgradeProxy(
            proxyAddress,
            newContractPath,
            ""
        );

        vm.stopBroadcast();
    }

    /**
     * @notice Alternative entry point that accepts proxy address and contract path as parameters
     * @param proxyAddress The address of the proxy contract to upgrade
     * @param newContractPath The path to the new contract file
     */
    function run(address proxyAddress, string calldata newContractPath) public {
        vm.startBroadcast();

        string memory oldContractName = vm.envOr("OLD_CONTRACT_NAME", string("OldContract"));
        string memory newContractName = vm.envOr("NEW_CONTRACT_NAME", string("NewContract"));

        console.log(string.concat(_CYAN, "Upgrading contract on chain:", _RESET), block.chainid);
        console.log(string.concat(_CYAN, "Old contract:", _RESET), oldContractName);
        console.log(string.concat(_CYAN, "New contract:", _RESET), newContractName);
        console.log(string.concat(_CYAN, "New contract path:", _RESET), newContractPath);
        console.log(string.concat(_CYAN, "Proxy address:", _RESET), proxyAddress);
        console.log(string.concat(_CYAN, "Upgrader address:", _RESET), msg.sender);

        // Validate upgrade safety with OpenZeppelin Foundry Upgrades
        if (bytes(oldContractName).length > 0) {
            console.log("");
            console.log(string.concat(_YELLOW, "Validating upgrade safety with OpenZeppelin Foundry Upgrades...", _RESET));
            Options memory opts;
            string memory referenceIdentifier = string.concat(oldContractName, ".sol");
            opts.referenceContract = referenceIdentifier;
            Upgrades.validateUpgrade(newContractPath, opts);
            console.log(string.concat(_BRIGHT_GREEN, "[PASS] OpenZeppelin Foundry Upgrades verification passed!", _RESET));
            console.log("");
        }

        // Upgrade to new contract
        Upgrades.upgradeProxy(
            proxyAddress,
            newContractPath,
            ""
        );

        vm.stopBroadcast();
    }
}

/**
 * @notice Anvil-specific variant for local testing
 */
contract UpgradeContractAnvil is GenericUpgrade {
    function run() public override {
        require(block.chainid == 31337, "must upgrade on anvil");
        super.run();
    }
}

/**
 * @notice Holesky-specific variant for testnet
 */
contract UpgradeContractHolesky is GenericUpgrade {
    function run() public override {
        require(block.chainid == 17000, "must upgrade on Holesky");
        super.run();
    }
}

/**
 * @notice Hoodi-specific variant for testnet
 */
contract UpgradeContractHoodi is GenericUpgrade {
    function run() public override {
        require(block.chainid == 560048, "must upgrade on Hoodi");
        super.run();
    }
}

/**
 * @notice Mainnet-specific variant
 */
contract UpgradeContractMainnet is GenericUpgrade {
    function run() public override {
        require(block.chainid == 1, "must upgrade on Mainnet");
        super.run();
    }
}

/**
 * @notice Generic upgrade contract (default, works on any chain)
 */
contract UpgradeContract is GenericUpgrade {
    function run() public override {
        super.run();
    }
}

