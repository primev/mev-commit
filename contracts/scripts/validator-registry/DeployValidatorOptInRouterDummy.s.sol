// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {ValidatorOptInRouterDummy} from "../../contracts/validator-registry/ValidatorOptInRouterDummy.sol";

contract BaseDeploy is Script {
    function deployValidatorOptInRouterDummy(
        address vanillaRegistry,
        address mevCommitAVS,
        address mevCommitMiddleware,
        address owner
    ) public returns (address) {
        console.log("Deploying ValidatorOptInRouterDummy on chain:", block.chainid);
        ValidatorOptInRouterDummy router = new ValidatorOptInRouterDummy();
        router.initialize(vanillaRegistry, mevCommitAVS, mevCommitMiddleware, owner);
        console.log("ValidatorOptInRouterDummy deployed to:", address(router));
        return address(router);
    }
}

contract DeployHolesky is BaseDeploy {
    address constant public VANILLA_REGISTRY = address(0);
    address constant public MEV_COMMIT_AVS = address(0);
    address constant public MEV_COMMIT_MIDDLEWARE = address(0);
    address constant public OWNER = address(0);

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        vm.startBroadcast();
        deployValidatorOptInRouterDummy(
            VANILLA_REGISTRY,
            MEV_COMMIT_AVS, 
            MEV_COMMIT_MIDDLEWARE,
            OWNER
        );
        vm.stopBroadcast();
    }
}
