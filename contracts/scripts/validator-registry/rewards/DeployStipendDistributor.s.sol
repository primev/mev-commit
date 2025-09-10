// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {StipendDistributor} from "../../../contracts/validator-registry/rewards/StipendDistributor.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployStipendDistributor(
        address owner,
        address stipendManager
    ) public returns (address) {
        console.log("Deploying StipendDistributor on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "StipendDistributor.sol",
            abi.encodeCall(
                StipendDistributor.initialize,
                (owner, stipendManager)
            )
        );
        console.log("StipendDistributor UUPS proxy deployed to:", address(proxy));
        StipendDistributor stipendDistributor = StipendDistributor(payable(proxy));
        console.log("StipendDistributor owner:", stipendDistributor.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public STIPEND_MANAGER = MainnetConstants.PRIMEV_TEAM_MULTISIG;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        deployStipendDistributor(
            OWNER,
            STIPEND_MANAGER
        );
        vm.stopBroadcast();
    }
}

contract DeployHoodi is BaseDeploy {
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public STIPEND_MANAGER = 0x1623fE21185c92BB43bD83741E226288B516134a;
 
    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");

        vm.startBroadcast();
        deployStipendDistributor(
            OWNER,
            STIPEND_MANAGER
        );
        vm.stopBroadcast();
    }
}
