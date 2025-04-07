// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";

contract UpgradeOracleImplementation is Script {
    function run() external {
        vm.startBroadcast();

        // Retrieve Oracle proxy address from environment variables
        address oracleProxyAddress = vm.envAddress("ORACLE_PROXY");

        console.log("Current Oracle proxy address:", oracleProxyAddress);

        // Upgrade the Oracle implementation
        Upgrades.upgradeProxy(
            oracleProxyAddress,
            "Oracle.sol",
            "" // No initialization needed for upgrade
        );

        console.log("Oracle implementation upgraded");

        vm.stopBroadcast();
    }
}
