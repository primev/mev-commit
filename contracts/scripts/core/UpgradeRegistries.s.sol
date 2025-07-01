// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";
import {BidderRegistryV2} from "../../contracts/core/BidderRegistryV2.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";
import {RegistryUpgradeLib} from "../../contracts/core/RegistryUpgradeLib.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";

contract UpgradeRegistries is Script {
    function run() external {
        vm.startBroadcast();

        address bidderRegistryProxyAddress = vm.envAddress("BIDDER_REGISTRY_PROXY");
        address providerRegistryProxyAddress = vm.envAddress("PROVIDER_REGISTRY_PROXY");
        console.log("BidderRegistry proxy address:", bidderRegistryProxyAddress);
        console.log("ProviderRegistry proxy address:", providerRegistryProxyAddress);
        uint256 newPayoutPeriodInMs = 1 hours * 1000; // 1 hour on mev-commit chain (ms timestamps)

        // Should be consistent with RegistryUpgradeLib.sol (used in unit tests).
        // That lib is un-usable here because delegate calls have a different msg.sender.
        BidderRegistry brv1 = BidderRegistry(payable(bidderRegistryProxyAddress));
        ProviderRegistry prv1 = ProviderRegistry(payable(providerRegistryProxyAddress));
        brv1.pause();
        prv1.pause();
        console.log("V1 contracts paused");

        brv1.manuallyWithdrawProtocolFee();
        prv1.manuallyWithdrawPenaltyFee();
        console.log("Fees manually withdrawn from v1 contracts");

        (address oldProtocolFeeRecipient,
        uint256 oldProtocolFeeAccumulatedAmount,
        uint256 oldProtocolFeeLastPayoutBlock,
        uint256 oldProtocolFeePayoutPeriodInBlocks) = brv1.protocolFeeTracker();
        console.log("V1 protocol fee recipient:", oldProtocolFeeRecipient);
        console.log("V1 protocol fee accumulated amount:", oldProtocolFeeAccumulatedAmount);
        console.log("V1 protocol fee last payout block:", oldProtocolFeeLastPayoutBlock);
        console.log("V1 protocol fee payout period in blocks:", oldProtocolFeePayoutPeriodInBlocks);
        
        (address oldPenaltyFeeRecipient,
        uint256 oldPenaltyFeeAccumulatedAmount,
        uint256 oldPenaltyFeeLastPayoutBlock,
        uint256 oldPenaltyFeePayoutPeriodInBlocks) = prv1.penaltyFeeTracker();
        console.log("V1 penalty fee recipient:", oldPenaltyFeeRecipient);
        console.log("V1 penalty fee accumulated amount:", oldPenaltyFeeAccumulatedAmount);
        console.log("V1 penalty fee last payout block:", oldPenaltyFeeLastPayoutBlock);
        console.log("V1 penalty fee payout period in blocks:", oldPenaltyFeePayoutPeriodInBlocks);

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
        brv2.setNewProtocolFeeRecipient(oldProtocolFeeRecipient);
        prv2.setNewPenaltyFeeRecipient(oldPenaltyFeeRecipient);
        console.log("V1 recipients have been set in V2");

        brv2.manuallyWithdrawProtocolFee();
        prv2.manuallyWithdrawPenaltyFee();
        console.log("Fees manually withdrawn from V2 contracts to properly set last payout timestamp");

        brv2.setNewFeePayoutPeriod(newPayoutPeriodInMs);
        prv2.setFeePayoutPeriod(newPayoutPeriodInMs);
        console.log("V2 payout periods in ms have been set");

        (address newProtocolFeeRecipient,
        uint256 newProtocolFeeAccumulatedAmount,
        uint256 newProtocolFeeLastPayoutTimestamp,
        uint256 newProtocolFeePayoutPeriodInMs) = brv2.protocolFeeTimestampTracker();
        console.log("V2 protocol fee recipient:", newProtocolFeeRecipient);
        console.log("V2 protocol fee accumulated amount:", newProtocolFeeAccumulatedAmount);
        console.log("V2 protocol fee last payout timestamp:", newProtocolFeeLastPayoutTimestamp);
        console.log("V2 protocol fee payout period in ms:", newProtocolFeePayoutPeriodInMs);

        (address newPenaltyFeeRecipient,
        uint256 newPenaltyFeeAccumulatedAmount,
        uint256 newPenaltyFeeLastPayoutTimestamp,
        uint256 newPenaltyFeePayoutPeriodInMs) = prv2.penaltyFeeTimestampTracker();
        console.log("V2 penalty fee recipient:", newPenaltyFeeRecipient);
        console.log("V2 penalty fee accumulated amount:", newPenaltyFeeAccumulatedAmount);
        console.log("V2 penalty fee last payout timestamp:", newPenaltyFeeLastPayoutTimestamp);
        console.log("V2 penalty fee payout period in ms:", newPenaltyFeePayoutPeriodInMs);

        brv2.unpause();
        prv2.unpause();
        console.log("V2 contracts unpaused");

        vm.stopBroadcast();
    }
}
