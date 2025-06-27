// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";
import {BidderRegistryV2} from "../../contracts/core/BidderRegistryV2.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";
import {RegistryUpgradeLib} from "../../contracts/core/RegistryUpgradeLib.sol";

contract UpgradeRegistries is Script {
    function run() external {
        vm.startBroadcast();

        address bidderRegistryProxyAddress = vm.envAddress("BIDDER_REGISTRY_PROXY");
        address providerRegistryProxyAddress = vm.envAddress("PROVIDER_REGISTRY_PROXY");
        console.log("BidderRegistry proxy address:", bidderRegistryProxyAddress);
        console.log("ProviderRegistry proxy address:", providerRegistryProxyAddress);
        uint256 newPayoutPeriodInMs = 1 hours * 1000; // 1 hour on mev-commit chain (ms timestamps)
        RegistryUpgradeLib.upgradeRegistries(
            bidderRegistryProxyAddress,
            providerRegistryProxyAddress,
            newPayoutPeriodInMs
        );

        vm.stopBroadcast();
    }
}
