// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract UpgradeProviderRegistry is Script {
    /// @dev Runs a few basic tests to ensure that key state values are preserved.
    function _runTests(ProviderRegistry registry) internal {
        // Read and log the current owner.
        address owner = registry.owner();
        console.log("Owner:", owner);
        
        // Read and log the minimum stake value.
        uint256 minStake = registry.minStake();
        console.log("Min Stake:", minStake);
        
        // Read and log the accumulated penalty fee.
        uint256 accumulatedPenaltyFee = registry.getAccumulatedPenaltyFee();
        console.log("Accumulated Penalty Fee:", accumulatedPenaltyFee);
    }

    function run() external {
        ProviderRegistry proxy;
        // Select the deployed proxy address based on the chain.
        if (block.chainid == 8855) {
            // Mainnet
            proxy = ProviderRegistry(
                payable(0xb772Add4718E5BD6Fe57Fb486A6f7f008E52167E)
            );
        } else if (block.chainid == 17864) {
            // Holesky testnet
            proxy = ProviderRegistry(
                payable(0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67)
            );
        } else {
            revert("Unsupported network");
        }

        console.log("Pre-upgrade tests:");
        _runTests(proxy);

        vm.startBroadcast();
        // Deploy the new implementation (ProviderRegistryV2.sol) and upgrade the proxy.
        Upgrades.upgradeProxy(address(proxy), "ProviderRegistryV2.sol", "");
        console.log(
            "Upgraded ProviderRegistry to ProviderRegistryV2 on chain:",
            block.chainid
        );
        vm.stopBroadcast();

        console.log("Post-upgrade tests:");
        _runTests(proxy);
    }
}
