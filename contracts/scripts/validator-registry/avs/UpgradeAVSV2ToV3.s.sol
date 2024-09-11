// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.20;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {MevCommitAVSV2} from "../../../contracts/validator-registry/avs/MevCommitAVSV2.sol";
import {MevCommitAVSV3} from "../../../contracts/validator-registry/avs/MevCommitAVSV3.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract UpgradeAVS is Script {
    function _runTests(MevCommitAVSV2 avs) internal {
        bytes memory iowarsKey1 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254b3f";
        bytes memory iowarsKey2 = hex"aca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
        
        bool optedIn1;
        bool optedIn2;
        try avs.isValidatorOptedIn(iowarsKey1) returns (bool result) {
            optedIn1 = result;
            console.log("isValidatorOptedIn:", optedIn1);
        } catch {
            console.log("isValidatorOptedIn reverted (expected for pre-upgrade)");
        }
        try avs.isValidatorOptedIn(iowarsKey2) returns (bool result) {
            optedIn2 = result;
            console.log("isValidatorOptedIn:", optedIn2);
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
        MevCommitAVSV2 existingMevCommitAVSProxy = MevCommitAVSV2(payable(0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4));

        console.log("Pre-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);

        vm.startBroadcast();
        Upgrades.upgradeProxy(
            address(existingMevCommitAVSProxy), 
            "MevCommitAVSV3.sol", 
            ""
        );
        console.log("Upgraded to MevCommitAVSV3");
        vm.stopBroadcast();

        console.log("Post-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);
    }
}

contract UpgradeAnvil is UpgradeAVS {
    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        MevCommitAVSV2 existingMevCommitAVSProxy = MevCommitAVSV2(payable(0x5FC8d32690cc91D4c39d9d3abcBD16989F875707));

        console.log("Pre-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);

        vm.startBroadcast();
        Upgrades.upgradeProxy(
            address(existingMevCommitAVSProxy), 
            "MevCommitAVSV3.sol", 
            "" 
        );
        console.log("Upgraded to MevCommitAVSV3");
        vm.stopBroadcast();

        console.log("Post-upgrade tests:");
        _runTests(existingMevCommitAVSProxy);
    }
}
