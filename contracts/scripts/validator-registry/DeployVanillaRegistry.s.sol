// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {VanillaRegistry} from "../../contracts/validator-registry/VanillaRegistry.sol";
import {console} from "forge-std/console.sol";
import {MainnetConstants} from "../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployVanillaRegistry(
        uint256 minStake,
        address slashOracle,
        address slashReceiver,
        uint256 unstakePeriodBlocks,
        uint256 payoutPeriodBlocks,
        address owner
    ) public returns (address) {
        console.log("Deploying VanillaRegistry on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "VanillaRegistry.sol",
            abi.encodeCall(
                VanillaRegistry.initialize,
                (minStake, slashOracle, slashReceiver, unstakePeriodBlocks, payoutPeriodBlocks, owner)
            )
        );
        console.log("VanillaRegistry UUPS proxy deployed to:", address(proxy));
        VanillaRegistry vanillaRegistry = VanillaRegistry(payable(proxy));
        console.log("VanillaRegistry owner:", vanillaRegistry.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    uint256 constant public MIN_STAKE = 1 ether;
    address constant public SLASH_ORACLE = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public SLASH_RECEIVER = MainnetConstants.COMMITMENT_HOLDINGS_MULTISIG;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 50000; // 50k * 12s ~= 1 week, which suffices for short-term manual slashing.
    uint256 constant public PAYOUT_PERIOD_BLOCKS = 12000; // ~ 1 day

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();
        deployVanillaRegistry(MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, PAYOUT_PERIOD_BLOCKS, OWNER);
        vm.stopBroadcast();
    }
}

contract DeployHolesky is BaseDeploy {
    uint256 constant public MIN_STAKE = 0.0001 ether; // 10k vals = 1 ETH cost
    address constant public SLASH_ORACLE = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;
    address constant public SLASH_RECEIVER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer
    uint256 constant public PAYOUT_PERIOD = 10000; // 10k * 12s = 1.39 days

    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        vm.startBroadcast();
        deployVanillaRegistry(MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, PAYOUT_PERIOD, OWNER);
        vm.stopBroadcast();
    }
}

contract DeployAnvil is BaseDeploy {
    uint256 constant public MIN_STAKE = 1 ether;
    address constant public SLASH_ORACLE = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address constant public SLASH_RECEIVER = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 100;
    uint256 constant public PAYOUT_PERIOD = 200;
    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        vm.startBroadcast();
        deployVanillaRegistry(MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, PAYOUT_PERIOD, msg.sender);
        vm.stopBroadcast();
    }
}

contract DeployHoodi is BaseDeploy {
    uint256 constant public MIN_STAKE = 0.0001 ether; // 10k vals = 1 ETH cost
    address constant public SLASH_ORACLE = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public SLASH_RECEIVER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    uint256 constant public UNSTAKE_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer
    uint256 constant public PAYOUT_PERIOD = 10000; // 10k * 12s = 1.39 days

    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;

    function run() external {
        require(block.chainid == 23118, "must deploy on Hoodi");
        vm.startBroadcast();
        deployVanillaRegistry(MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD_BLOCKS, PAYOUT_PERIOD, OWNER);
        vm.stopBroadcast();
    }
}
