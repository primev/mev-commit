// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.20;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {MevCommitAVSV2} from "../../../contracts/validator-registry/avs/MevCommitAVSV2.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract UpgradeAVS is Script {
    function _runTests(MevCommitAVS avs) internal {
        bytes memory invalidKey = hex"89898989";
        
        bool optedIn;
        try avs.isValidatorOptedIn(invalidKey) returns (bool result) {
            optedIn = result;
            console.log("isValidatorOptedIn:", optedIn);
        } catch {
            console.log("isValidatorOptedIn reverted (expected for pre-upgrade)");
        }

        address owner = avs.owner();
        console.log("owner:", owner);
        uint256 valDeregPeriodBlocks = avs.validatorDeregPeriodBlocks();
        console.log("validatorDeregPeriodBlocks:", valDeregPeriodBlocks);
    }
}

contract UpgradeHolesky is UpgradeAVS {
    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        MevCommitAVS existingMevCommitAVSProxy = MevCommitAVS(payable(0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4));

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
