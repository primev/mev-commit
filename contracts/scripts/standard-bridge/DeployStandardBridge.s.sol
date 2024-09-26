// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {SettlementGateway} from "../../contracts/standard-bridge/SettlementGateway.sol";
import {L1Gateway} from "../../contracts/standard-bridge/L1Gateway.sol";
import {Allocator} from "../../contracts/standard-bridge/Allocator.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";

contract DeploySettlementGateway is Script {
    function run() external {

        vm.startBroadcast();

        address relayerAddr = vm.envAddress("RELAYER_ADDR");
        require(relayerAddr != address(0), "RELAYER_ADDR is not set");
        address whitelistAddr = vm.envAddress("WHITELIST_ADDR");
        require(whitelistAddr != address(0), "WHITELIST_ADDR is not set");

        address sgProxy = Upgrades.deployUUPSProxy(
            "SettlementGateway.sol",
            abi.encodeCall(SettlementGateway.initialize,
            (whitelistAddr,
            msg.sender, // Owner
            relayerAddr,
            1, 1)) // Fees set to 1 wei for now
        );
        SettlementGateway gateway = SettlementGateway(payable(sgProxy));
        console.log("Standard bridge gateway for settlement chain deployed to:",
            address(gateway));
        address allocatorProxy = Upgrades.deployUUPSProxy(
            "Allocator.sol",
            abi.encodeCall(Allocator.initialize, (msg.sender))
        );
        Allocator allocator = Allocator(payable(allocatorProxy));
        console.log("Allocator deployed to:", address(allocator));

        allocator.addToWhitelist(address(gateway));
        console.log("Settlement gateway has been whitelisted. Gateway contract address:", address(gateway));

        string memory jsonOutput = string.concat(
            "{'settlement_gateway_addr': '",
            Strings.toHexString(address(gateway)),
            "', 'whitelist_addr': '",
            Strings.toHexString(address(allocator)),
            "'}"
        );
        console.log("JSON_DEPLOY_ARTIFACT:", jsonOutput); 

        vm.stopBroadcast();
    }
}

contract DeployL1Gateway is Script {
    function run() external {

        vm.startBroadcast();

        address relayerAddr = vm.envAddress("RELAYER_ADDR");

        address l1gProxy = Upgrades.deployUUPSProxy(
            "L1Gateway.sol",
            abi.encodeCall(L1Gateway.initialize,
            (msg.sender, // Owner
            relayerAddr,
            1, 1)) // Fees set to 1 wei for now
        );
        L1Gateway gateway = L1Gateway(payable(l1gProxy));
        console.log("Standard bridge gateway for l1 deployed to:",
            address(gateway));
        
        string memory jsonOutput = string.concat(
            "{'l1_gateway_addr': '",
            Strings.toHexString(address(gateway)),
            "'}"
        );
        console.log("JSON_DEPLOY_ARTIFACT:", jsonOutput);

        vm.stopBroadcast();
    }
}
