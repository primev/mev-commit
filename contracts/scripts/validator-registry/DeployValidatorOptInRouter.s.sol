// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.25;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ValidatorOptInRouter} from "../../contracts/validator-registry/ValidatorOptInRouter.sol";

contract BaseDeploy is Script {
    function deployValidatorOptInRouter(
        address vanillaRegistry,
        address mevCommitAVS,
        address owner
    ) public returns (address) {
        console.log("Deploying ValidatorOptInRouter on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorOptInRouter.sol",
            abi.encodeCall(
                ValidatorOptInRouter.initialize,
                (vanillaRegistry, mevCommitAVS, owner)
            )
        );
        console.log("ValidatorOptInRouter UUPS proxy deployed to:", address(proxy));
        ValidatorOptInRouter router = ValidatorOptInRouter(payable(proxy));
        console.log("ValidatorOptInRouter owner:", router.owner());
        return proxy;
    }
}

contract DeployHolesky is BaseDeploy {
    address constant public VANILLA_REGISTRY = 0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf;
    address constant public MEV_COMMIT_AVS = 0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4;

    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");

        vm.startBroadcast();
        deployValidatorOptInRouter(
            VANILLA_REGISTRY,
            MEV_COMMIT_AVS,
            OWNER
        );
        vm.stopBroadcast();
    }
}
