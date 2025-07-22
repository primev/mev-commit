// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.29;

import {Script} from "forge-std/Script.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {Oracle} from "../../contracts/core/Oracle.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";

/**
 * @notice This script deploys a new upgradeable PreconfManager contract using the UUPS proxy pattern.
 *         It then updates the PreconfManager address in the BidderRegistry, ProviderRegistry, and Oracle.
 *
 * Expected environment variables:
 *  - BIDDER_REGISTRY_PROXY: Address of the BidderRegistry proxy.
 *  - PROVIDER_REGISTRY_PROXY: Address of the ProviderRegistry proxy.
 *  - ORACLE_PROXY: Address of the Oracle proxy.
 *  - BLOCK_TRACKER_PROXY: Address of the BlockTracker proxy.
 */
contract DeployNewPreconfManagerUpgradeable is Script {
    function run() external {
        vm.startBroadcast();

        // Retrieve dependent contract proxy addresses from environment variables.
        address bidderRegistryAddress = vm.envAddress("BIDDER_REGISTRY_PROXY");
        address providerRegistryAddress = vm.envAddress(
            "PROVIDER_REGISTRY_PROXY"
        );
        address oracleAddress = vm.envAddress("ORACLE_PROXY");
        address blockTrackerAddress = vm.envAddress("BLOCK_TRACKER_PROXY");

        // Set initialization parameters.
        uint64 commitmentDispatchWindow = 2000;

        // Deploy a new PreconfManager instance via UUPS proxy.
        address newPreconfManagerProxy = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(
                PreconfManager.initialize,
                (
                    providerRegistryAddress,
                    bidderRegistryAddress,
                    oracleAddress,
                    msg.sender,
                    blockTrackerAddress,
                    commitmentDispatchWindow
                )
            )
        );
        PreconfManager newPreconfManager = PreconfManager(
            payable(newPreconfManagerProxy)
        );
        console.log(
            "New PreconfManager (via UUPS proxy) deployed at:",
            address(newPreconfManager)
        );

        // Update dependent contracts with the new PreconfManager address.
        BidderRegistry(payable(bidderRegistryAddress)).setPreconfManager(
            address(newPreconfManager)
        );
        console.log("BidderRegistry updated with new PreconfManager");

        ProviderRegistry(payable(providerRegistryAddress)).setPreconfManager(
            address(newPreconfManager)
        );
        console.log("ProviderRegistry updated with new PreconfManager");

        Oracle(payable(oracleAddress)).setPreconfManager(
            address(newPreconfManager)
        );
        console.log("Oracle updated with new PreconfManager");

        vm.stopBroadcast();
    }
}
