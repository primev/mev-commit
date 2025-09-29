// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {RocketMinipoolRegistry} from "../../contracts/validator-registry/rocketpool/RocketMinipoolRegistry.sol";
import {console} from "forge-std/console.sol";
import {MainnetConstants} from "../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployRocketMinipoolRegistry(
        address owner,
        address freezeOracle,
        address unfreezeReceiver,
        address rocketStorage,
        uint256 unfreezeFee,
        uint256 deregistrationPeriod
    ) public returns (address) {
        console.log("Deploying RocketMinipoolRegistry on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "RocketMinipoolRegistry.sol",
            abi.encodeCall(
                RocketMinipoolRegistry.initialize,
                (owner, freezeOracle, unfreezeReceiver, rocketStorage, unfreezeFee, deregistrationPeriod)
            )
        );
        console.log("RocketMinipoolRegistry UUPS proxy deployed to:", address(proxy));
        RocketMinipoolRegistry rocketMinipoolRegistry = RocketMinipoolRegistry(payable(proxy));
        console.log("RocketMinipoolRegistry owner:", rocketMinipoolRegistry.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public FREEZE_ORACLE = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public UNFREEZE_RECEIVER = MainnetConstants.COMMITMENT_HOLDINGS_MULTISIG;
    uint256 constant public UNFREEZE_FEE = 1 ether;
    uint256 constant public DEREGISTRATION_PERIOD = 86400; // 86400/60s = 1 day
    address constant public ROCKET_STORAGE = 0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();
        deployRocketMinipoolRegistry(OWNER, FREEZE_ORACLE, UNFREEZE_RECEIVER, ROCKET_STORAGE, UNFREEZE_FEE, DEREGISTRATION_PERIOD);
        vm.stopBroadcast();
    }
}


contract DeployHoodi is BaseDeploy {
    address constant public FREEZE_ORACLE = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public UNFREEZE_RECEIVER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    uint256 constant public UNFREEZE_FEE = 0.01 ether;
    uint256 constant public DEREGISTRATION_PERIOD = 300; // 300/60s = 5 minutes
    address constant public ROCKET_STORAGE = 0x594Fb75D3dc2DFa0150Ad03F99F97817747dd4E1;

    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;

    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");
        vm.startBroadcast();
        deployRocketMinipoolRegistry(OWNER, FREEZE_ORACLE, UNFREEZE_RECEIVER, ROCKET_STORAGE, UNFREEZE_FEE, DEREGISTRATION_PERIOD);
        vm.stopBroadcast();
    }
}
