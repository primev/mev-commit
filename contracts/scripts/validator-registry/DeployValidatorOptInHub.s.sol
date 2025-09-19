// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ValidatorOptInHub} from "../../contracts/validator-registry/ValidatorOptInHub.sol";
import {MainnetConstants} from "../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployValidatorOptInHub(
        address[] memory registries,
        address owner
    ) public returns (address) {
        console.log("Deploying ValidatorOptInHub on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorOptInHub.sol",
            abi.encodeCall(
                ValidatorOptInHub.initialize,
                (registries, owner)
            )
        );
        console.log("ValidatorOptInHub UUPS proxy deployed to:", address(proxy));
        ValidatorOptInHub hub = ValidatorOptInHub(payable(proxy));
        console.log("ValidatorOptInHub owner:", hub.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    address constant public MEV_COMMIT_MIDDLEWARE = 0x21fD239311B050bbeE7F32850d99ADc224761382;
    address constant public MEV_COMMIT_AVS = 0xBc77233855e3274E1903771675Eb71E602D9DC2e;
    address constant public VANILLA_REGISTRY = 0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9;
    address[] public REGISTRIES = [MEV_COMMIT_MIDDLEWARE, MEV_COMMIT_AVS, VANILLA_REGISTRY];
    
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        deployValidatorOptInHub(
            REGISTRIES,
            OWNER
        );
        vm.stopBroadcast();
    }
}

contract DeployHoodi is BaseDeploy {
    address constant public MEV_COMMIT_MIDDLEWARE = 0x8E847EC4a36c8332652aB3b2B7D5c54dE29c7fde;
    address constant public MEV_COMMIT_AVS = 0xdF8649d298ad05f019eE4AdBD6210867B8AB225F;
    address constant public VANILLA_REGISTRY = 0x536F0792c5D5Ed592e67a9260606c85F59C312F0;
    
    address[] public REGISTRIES = [MEV_COMMIT_MIDDLEWARE, MEV_COMMIT_AVS, VANILLA_REGISTRY];

    //This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;

    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");

        vm.startBroadcast();
        deployValidatorOptInHub(
            REGISTRIES,
            OWNER
        );
        vm.stopBroadcast();
    }
}
