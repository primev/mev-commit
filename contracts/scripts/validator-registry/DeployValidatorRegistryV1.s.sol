// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.20;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ValidatorRegistryV1} from "../../contracts/validator-registry/ValidatorRegistryV1.sol";
import {console} from "forge-std/console.sol";

contract BaseDeploy is Script {
    function deployValidatorRegistry(
        uint256 minStake,
        uint256 slashAmount,
        address slashOracle,
        address slashReceiver,
        uint256 unstakePeriodBlocks,
        address owner
    ) public returns (address) {
        console.log("Deploying ValidatorRegistryV1 on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorRegistryV1.sol",
            abi.encodeCall(
                ValidatorRegistryV1.initialize,
                (minStake, slashAmount, slashOracle, slashReceiver, unstakePeriodBlocks, owner)
            )
        );
        console.log("ValidatorRegistryV1 UUPS proxy deployed to:", address(proxy));
        ValidatorRegistryV1 validatorRegistry = ValidatorRegistryV1(payable(proxy));
        console.log("ValidatorRegistryV1 owner:", validatorRegistry.owner());
        return proxy;
    }
}

contract DeployHolesky is BaseDeploy {
    uint256 constant public MIN_STAKE = 0.0001 ether; // 10k vals = 1 ETH cost
    uint256 constant public SLASH_AMOUNT = 0.00003 ether; 
    address constant public SLASH_ORACLE = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;
    address constant public SLASH_RECEIVER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        vm.startBroadcast();
        deployValidatorRegistry(MIN_STAKE, SLASH_AMOUNT,
           SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, msg.sender);
        vm.stopBroadcast();
    }
}

contract DeployAnvil is BaseDeploy {
    uint256 constant public MIN_STAKE = 3 ether;
    uint256 constant public SLASH_AMOUNT = 1 ether;
    address constant public SLASH_ORACLE = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address constant public SLASH_RECEIVER = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 100;

    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        vm.startBroadcast();
        deployValidatorRegistry(MIN_STAKE, SLASH_AMOUNT,
            SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, msg.sender);
        vm.stopBroadcast();
    }
}
