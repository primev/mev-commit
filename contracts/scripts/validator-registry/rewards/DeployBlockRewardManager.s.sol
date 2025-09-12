// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {BlockRewardManager} from "../../../contracts/validator-registry/rewards/BlockRewardManager.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployBlockRewardManager(
        address owner,
        uint256 rewardsPctBps,
        address treasury
    ) public returns (address) {
        console.log("Deploying BlockRewardManager on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "BlockRewardManager.sol",
            abi.encodeCall(
                BlockRewardManager.initialize,
                (owner, rewardsPctBps, payable(treasury))
            )
        );
        console.log("BlockRewardManager UUPS proxy deployed to:", address(proxy));
        BlockRewardManager rewardsV2 = BlockRewardManager(payable(proxy));
        console.log("BlockRewardManager owner:", rewardsV2.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    //address public TREASURY;
    uint256 constant public REWARDS_PCT_BPS = 0;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();
        //deploy call here
        vm.stopBroadcast();
    }
}

contract DeployHoodi is BaseDeploy {
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public TREASURY = 0x1623fE21185c92BB43bD83741E226288B516134a;
    uint256 constant public REWARDS_PCT_BPS = 0;
    
    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");

        vm.startBroadcast();
        deployBlockRewardManager(
            OWNER,
            REWARDS_PCT_BPS,
            TREASURY
        );
        vm.stopBroadcast();
    }
}
