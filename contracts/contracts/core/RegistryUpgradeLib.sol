// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";
import {BidderRegistryV2} from "../../contracts/core/BidderRegistryV2.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";

library RegistryUpgradeLib {
    function upgradeRegistries(
        address bidderRegistryProxyAddress,
        address providerRegistryProxyAddress,
        uint256 newPayoutPeriodInMs
    ) external {
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
        brv2.setNewFeePayoutPeriod(newPayoutPeriodInMs);
        prv2.setFeePayoutPeriod(newPayoutPeriodInMs);
        console.log("Payout periods updated to 1 hour in ms");

        brv2.manuallyWithdrawProtocolFee();
        prv2.manuallyWithdrawPenaltyFee();
        console.log("Fees manually withdrawn to overwrite last payout timestamp");
    }
}
