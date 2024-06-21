// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ValidatorRegistryV1} from "../../contracts/validator-registry/ValidatorRegistryV1.sol";

contract Deploy is Script {
    uint256 constant UNSTAKE_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer
    uint256 constant MIN_STAKE = 0.0001 ether; // 10k vals = 1 ETH cost

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        vm.startBroadcast();
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorRegistryV1.sol",
            abi.encodeCall(
                ValidatorRegistryV1.initialize,
                (MIN_STAKE, UNSTAKE_PERIOD_BLOCKS, msg.sender)
            )
        );
        console.log("ValidatorRegistryV1 UUPS proxy deployed to:", address(proxy));
        ValidatorRegistryV1 validatorRegistry = ValidatorRegistryV1(payable(proxy));
        console.log("ValidatorRegistryV1 owner:", validatorRegistry.owner());
        vm.stopBroadcast();
    }
}
