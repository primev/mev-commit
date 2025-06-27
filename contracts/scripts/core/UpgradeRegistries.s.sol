// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";
import {BidderRegistryV2} from "../../contracts/core/BidderRegistryV2.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";

contract UpgradeRegistries is Script {
    function run() external {
        vm.startBroadcast();

        address bidderRegistryProxyAddress = vm.envAddress("BIDDER_REGISTRY_PROXY");
        address providerRegistryProxyAddress = vm.envAddress("PROVIDER_REGISTRY_PROXY");
        console.log("BidderRegistry proxy address:", bidderRegistryProxyAddress);
        console.log("ProviderRegistry proxy address:", providerRegistryProxyAddress);

        Upgrades.upgradeProxy(
            bidderRegistryProxyAddress,
            "BidderRegistryV2.sol",
            ""
        );
        Upgrades.upgradeProxy(
            providerRegistryProxyAddress,
            "ProviderRegistryV2.sol",
            ""
        );
        console.log("Registries upgraded to V2");

        BidderRegistryV2 brv2 = BidderRegistryV2(payable(bidderRegistryProxyAddress));
        ProviderRegistryV2 prv2 = ProviderRegistryV2(payable(providerRegistryProxyAddress));
        brv2.setNewFeePayoutPeriod(1 hours * 1000);
        prv2.setFeePayoutPeriod(1 hours * 1000);
        console.log("Payout periods updated to 1 hour in ms");

        brv2.manuallyWithdrawProtocolFee();
        prv2.manuallyWithdrawPenaltyFee();
        console.log("Fees manually withdrawn to overwrite last payout timestamp");

        vm.stopBroadcast();
    }
}
