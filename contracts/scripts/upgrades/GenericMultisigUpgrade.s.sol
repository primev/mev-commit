// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {Options} from "openzeppelin-foundry-upgrades/Options.sol";

/**
 * @title GenericMultisigUpgrade
 * @notice Generic script to deploy implementation contracts for multisig upgrades.
 * 
 * This script deploys ONLY the implementation contract without initializing it.
 * After deployment, use your multisig UI to call upgradeToAndCall() on the proxy.
 * 
 * This script reads configuration from environment variables:
 *  - NEW_CONTRACT_NAME: Name of the new contract (used for contract name resolution, e.g., "MevCommitAVSV2")
 *  - NEW_CONTRACT_PATH: Path to the new contract file (e.g., "MevCommitAVSV2.sol") - used for logging only
 * 
 * ⚠️  IMPORTANT: Run validation before deploying!
 * ./validate-upgrade.sh --contract <NEW_CONTRACT_NAME> --reference <OLD_CONTRACT_NAME>
 * 
 * Usage:
 *  forge script scripts/upgrades/GenericMultisigDeploy.s.sol:DeployMultisigImplMainnet --rpc-url <RPC_URL> --sender <DEPLOYER_ADDRESS> --broadcast --verify -vvvv
 * 
 * Note: The deployer address can be any account (doesn't need to be the multisig).
 * The implementation contract serves only as a blueprint and has no ownership.
 */
contract GenericMultisigUpgrade is Script {
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

        // Get contract name and path from environment variables
        string memory newContractName = vm.envOr("NEW_CONTRACT_NAME", string(""));
        string memory newContractPath = vm.envOr("NEW_CONTRACT_PATH", string(""));
        string memory oldContractName = vm.envOr("OLD_CONTRACT_NAME", string(""));

        console.log(string.concat(_CYAN, "Deploying implementation contract on chain:", _RESET), block.chainid);
        console.log(string.concat(_CYAN, "Contract name:", _RESET), newContractName);
        if (bytes(newContractPath).length > 0) {
            console.log(string.concat(_CYAN, "Contract path:", _RESET), newContractPath);
        }
        console.log(string.concat(_CYAN, "Deployer address:", _RESET), msg.sender);

        // Construct contract identifiers for validation and deployment
        // Format: ContractName.sol (OpenZeppelin will resolve to fully qualified name)
        string memory contractIdentifier = bytes(newContractPath).length > 0 
            ? newContractPath 
            : string.concat(newContractName, ".sol");
        
        // Validate upgrade safety with OpenZeppelin Foundry Upgrades
        if (bytes(oldContractName).length > 0) {
            console.log("");
            console.log(string.concat(_YELLOW, "Validating upgrade safety with OpenZeppelin Foundry Upgrades...", _RESET));
            Options memory opts;
            // Construct reference contract identifier (just filename, OpenZeppelin will resolve it)
            string memory referenceIdentifier = string.concat(oldContractName, ".sol");
            opts.referenceContract = referenceIdentifier;
            Upgrades.validateUpgrade(contractIdentifier, opts);
            console.log(string.concat(_BRIGHT_GREEN, "[PASS] OpenZeppelin Foundry Upgrades verification passed!", _RESET));
            console.log("");
        }

        // Deploy the implementation contract using vm.deployCode()
        // This deploys the contract bytecode without calling any constructor/initializer
        address implementation = vm.deployCode(contractIdentifier);
        
        if (implementation == address(0)) {
            revert("Failed to deploy implementation");
        }

        console.log("");
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));
        console.log(string.concat(_BRIGHT_GREEN, "Implementation contract deployed successfully!", _RESET));
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));
        console.log(string.concat(_CYAN, "Contract name:", _RESET), newContractName);
        console.log(string.concat(_BRIGHT_CYAN, "Implementation address:", _RESET), implementation);
        console.log("");
        console.log(string.concat(_YELLOW, "Next steps:", _RESET));
        console.log(string.concat(_YELLOW, "1. Copy the implementation address above", _RESET));
        console.log(string.concat(_YELLOW, "2. Use your multisig UI (e.g., Safe wallet) to call upgradeToAndCall() on the proxy", _RESET));
        console.log(string.concat(_YELLOW, "3. Function: upgradeToAndCall(implementation, callData)", _RESET));
        console.log(string.concat(_YELLOW, "   - implementation: the address shown above", _RESET));
        console.log(string.concat(_YELLOW, "   - callData: empty bytes (\"0x\") if no initialization needed", _RESET));
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));
        console.log("");

        vm.stopBroadcast();
    }

    /**
     * @notice Alternative entry point (directly called by the user) that accepts contract name as parameter
     * @param newContractName The name of the new contract (e.g., "MevCommitAVSV2")
     */
    function run(string calldata newContractName) public {
        vm.startBroadcast();

        string memory newContractPath = vm.envOr("NEW_CONTRACT_PATH", string(""));
        string memory oldContractName = vm.envOr("OLD_CONTRACT_NAME", string(""));

        console.log(string.concat(_CYAN, "Deploying implementation contract on chain:", _RESET), block.chainid);
        console.log(string.concat(_CYAN, "Contract name:", _RESET), newContractName);
        if (bytes(newContractPath).length > 0) {
            console.log(string.concat(_CYAN, "Contract path:", _RESET), newContractPath);
        }
        console.log(string.concat(_CYAN, "Deployer address:", _RESET), msg.sender);

        // Construct contract identifier for validation and deployment
        // Format: ContractName.sol (OpenZeppelin will resolve to fully qualified name)
        string memory contractIdentifier = bytes(newContractPath).length > 0 
            ? newContractPath 
            : string.concat(newContractName, ".sol");
        
        // Validate upgrade safety with OpenZeppelin Foundry Upgrades
        if (bytes(oldContractName).length > 0) {
            console.log("");
            console.log(string.concat(_YELLOW, "Validating upgrade safety with OpenZeppelin Foundry Upgrades...", _RESET));
            Options memory opts;
            string memory referenceIdentifier = string.concat(oldContractName, ".sol");
            opts.referenceContract = referenceIdentifier;
            Upgrades.validateUpgrade(contractIdentifier, opts);
            console.log(string.concat(_BRIGHT_GREEN, "[PASS] OpenZeppelin Foundry Upgrades verification passed!", _RESET));
            console.log("");
        }

        address implementation = vm.deployCode(contractIdentifier);
        
        if (implementation == address(0)) {
            revert("Failed to deploy implementation");
        }

        console.log("");
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));
        console.log(string.concat(_BRIGHT_GREEN, "Implementation contract deployed successfully!", _RESET));
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));
        console.log(string.concat(_CYAN, "Contract name:", _RESET), newContractName);
        console.log(string.concat(_BRIGHT_CYAN, "Implementation address:", _RESET), implementation);
        console.log(string.concat(_BRIGHT_GREEN, "========================================", _RESET));

        vm.stopBroadcast();
    }
}

/**
 * @notice Anvil-specific variant for local testing
 */
contract DeployMultisigImplAnvil is GenericMultisigUpgrade {
    function run() public override {
        require(block.chainid == 31337, "must deploy on anvil");
        super.run();
    }
}

/**
 * @notice Holesky-specific variant for testnet
 */
contract DeployMultisigImplHolesky is GenericMultisigUpgrade {
    function run() public override {
        require(block.chainid == 17000, "must deploy on Holesky");
        super.run();
    }
}

/**
 * @notice Hoodi-specific variant for testnet
 */
contract DeployMultisigImplHoodi is GenericMultisigUpgrade {
    function run() public override {
        require(block.chainid == 560048, "must deploy on Hoodi");
        super.run();
    }
}

/**
 * @notice Mainnet-specific variant
 */
contract DeployMultisigImplMainnet is GenericMultisigUpgrade {
    function run() public override {
        require(block.chainid == 1, "must deploy on Mainnet");
        super.run();
    }
}

/**
 * @notice Generic deploy contract (default, works on any chain)
 */
contract DeployMultisigImpl is GenericMultisigUpgrade {
    function run() public override {
        super.run();
    }
}

