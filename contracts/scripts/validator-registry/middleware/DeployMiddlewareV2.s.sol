// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable no-console

import {MevCommitMiddlewareV2} from "../../../contracts/validator-registry/middleware/MevCommitMiddlewareV2.sol";
import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";

contract DeployMiddlewareV2 is Script {
    function run() external {
        vm.startBroadcast();
        MevCommitMiddlewareV2 newImplementation = new MevCommitMiddlewareV2();
        console.log("Deployed new implementation contract at address: ", address(newImplementation));
        vm.stopBroadcast();
    }
}
