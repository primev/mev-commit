// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {VanillaRegistryV2} from "../../contracts/validator-registry/VanillaRegistryV2.sol";
import {console} from "forge-std/console.sol";
import {VanillaRegistryStorageV2} from "../../contracts/validator-registry/VanillaRegistryStorageV2.sol";
import {MevCommitAVS} from "../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {MevCommitAVSStorage} from "../../contracts/validator-registry/avs/MevCommitAVSStorage.sol";

contract GetVanillaRegistryParams is Script {
    function run() external view {
        console.log("Getting params for VanillaRegistry on chain:", block.chainid);
        address vanillaRegAddr = 0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9;

        address owner = VanillaRegistryV2(payable(vanillaRegAddr)).owner();
        console.log("VanillaRegistry owner:", owner);
        bool isPaused = VanillaRegistryV2(payable(vanillaRegAddr)).paused();
        console.log("VanillaRegistry isPaused:", isPaused);
        uint256 minStake = VanillaRegistryStorageV2(payable(vanillaRegAddr)).minStake();
        console.log("VanillaRegistry minStake:", minStake);
        address slashOracle = VanillaRegistryStorageV2(payable(vanillaRegAddr)).slashOracle();
        console.log("VanillaRegistry slashOracle:", slashOracle);
        uint256 unstakePeriodBlocks = VanillaRegistryStorageV2(payable(vanillaRegAddr)).unstakePeriodBlocks();
        console.log("VanillaRegistry unstakePeriodBlocks:", unstakePeriodBlocks);
        uint256 accumulatedFunds = VanillaRegistryV2(payable(vanillaRegAddr)).getAccumulatedSlashingFunds();
        console.log("VanillaRegistry accumulatedFunds:", accumulatedFunds);
    }
}

contract GetMevCommitAVSParams is Script {
    function run() external view {
        console.log("Getting params for MevCommitAVS on chain:", block.chainid);
        address avsAddr = 0xBc77233855e3274E1903771675Eb71E602D9DC2e;

        address avsOwner = MevCommitAVS(payable(avsAddr)).owner();
        console.log("MevCommitAVS owner:", avsOwner);
        bool isPaused = MevCommitAVS(payable(avsAddr)).paused();
        console.log("MevCommitAVS isPaused:", isPaused);
        address[] memory restakeableStrategies = MevCommitAVS(payable(avsAddr)).getRestakeableStrategies();
        uint256 len = restakeableStrategies.length;
        console.log("MevCommitAVS restakeableStrategies length:", len);
        for (uint256 i = 0; i < len; ++i) {
            console.log("MevCommitAVS restakeableStrategy:", restakeableStrategies[i]);
        }
        address freezeOracle = MevCommitAVSStorage(payable(avsAddr)).freezeOracle();
        console.log("MevCommitAVS freezeOracle:", freezeOracle);
        uint256 unfreezeFee = MevCommitAVSStorage(payable(avsAddr)).unfreezeFee();
        console.log("MevCommitAVS unfreezeFee:", unfreezeFee);
        address unfreezeReceiver = MevCommitAVSStorage(payable(avsAddr)).unfreezeReceiver();
        console.log("MevCommitAVS unfreezeReceiver:", unfreezeReceiver);
        uint256 unfreezePeriodBlocks = MevCommitAVSStorage(payable(avsAddr)).unfreezePeriodBlocks();
        console.log("MevCommitAVS unfreezePeriodBlocks:", unfreezePeriodBlocks);
        uint256 operatorDeregPeriodBlocks = MevCommitAVSStorage(payable(avsAddr)).operatorDeregPeriodBlocks();
        console.log("MevCommitAVS operatorDeregPeriodBlocks:", operatorDeregPeriodBlocks);
        uint256 validatorDeregPeriodBlocks = MevCommitAVSStorage(payable(avsAddr)).validatorDeregPeriodBlocks();
        console.log("MevCommitAVS validatorDeregPeriod:", validatorDeregPeriodBlocks);
        uint256 lstRestakerDeregPeriodBlocks = MevCommitAVSStorage(payable(avsAddr)).lstRestakerDeregPeriodBlocks();
        console.log("MevCommitAVS lstRestakerDeregPeriodBlocks:", lstRestakerDeregPeriodBlocks);
    }
}
