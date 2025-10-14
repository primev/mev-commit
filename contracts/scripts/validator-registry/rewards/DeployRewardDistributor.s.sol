// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {RewardDistributor} from "../../../contracts/validator-registry/rewards/RewardDistributor.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployRewardDistributor(
        address owner,
        address rewardManager
    ) public returns (address) {
        console.log("Deploying RewardDistributor on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "RewardDistributor.sol",
            abi.encodeCall(
                RewardDistributor.initialize,
                (owner, rewardManager)
            )
        );
        console.log("RewardDistributor UUPS proxy deployed to:", address(proxy));
        RewardDistributor rewardDistributor = RewardDistributor(payable(proxy));
        console.log("RewardDistributor owner:", rewardDistributor.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    // address constant public REWARD_MANAGER

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();
        //deploy call here
        vm.stopBroadcast();
    }
}

contract DeployHoodi is BaseDeploy {
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public REWARD_MANAGER = 0x1623fE21185c92BB43bD83741E226288B516134a;
 
    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");

        vm.startBroadcast();
        deployRewardDistributor(
            OWNER,
            REWARD_MANAGER
        );
        vm.stopBroadcast();
    }
}
