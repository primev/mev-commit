// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.29;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ValidatorOptInRouter} from "../../contracts/validator-registry/ValidatorOptInRouter.sol";
import {MainnetConstants} from "../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployValidatorOptInRouter(
        address vanillaRegistry,
        address mevCommitAVS,
        address mevCommitMiddleware,
        address owner
    ) public returns (address) {
        console.log("Deploying ValidatorOptInRouter on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorOptInRouter.sol",
            abi.encodeCall(
                ValidatorOptInRouter.initialize,
                (vanillaRegistry, mevCommitAVS, mevCommitMiddleware, owner)
            )
        );
        console.log("ValidatorOptInRouter UUPS proxy deployed to:", address(proxy));
        ValidatorOptInRouter router = ValidatorOptInRouter(payable(proxy));
        console.log("ValidatorOptInRouter owner:", router.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public VANILLA_REGISTRY = 0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9;
    address constant public MEV_COMMIT_AVS = 0xBc77233855e3274E1903771675Eb71E602D9DC2e;
    address constant public MEV_COMMIT_MIDDLEWARE = 0x21fD239311B050bbeE7F32850d99ADc224761382;
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        deployValidatorOptInRouter(
            VANILLA_REGISTRY,
            MEV_COMMIT_AVS,
            MEV_COMMIT_MIDDLEWARE,
            OWNER
        );
        vm.stopBroadcast();
    }
}

contract DeployHolesky is BaseDeploy {
    address constant public VANILLA_REGISTRY = 0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf;
    address constant public MEV_COMMIT_AVS = 0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4;
    address constant public MEV_COMMIT_MIDDLEWARE = 0x0D5A6dd3Ba8C6385ecA623B56199b7FFC490792a;

    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");

        vm.startBroadcast();
        deployValidatorOptInRouter(
            VANILLA_REGISTRY,
            MEV_COMMIT_AVS,
            MEV_COMMIT_MIDDLEWARE,
            OWNER
        );
        vm.stopBroadcast();
    }
}
