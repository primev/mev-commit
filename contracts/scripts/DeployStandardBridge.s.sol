// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;
import "forge-std/Script.sol";
import {SettlementGateway} from "../contracts/standard-bridge/SettlementGateway.sol";
import {L1Gateway} from "../contracts/standard-bridge/L1Gateway.sol";
import {Whitelist} from "../contracts/Whitelist.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

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
        SettlementGateway gateway = SettlementGateway(sgProxy);
        console.log("Standard bridge gateway for settlement chain deployed to:",
            address(gateway));
        address whitelistProxy = Upgrades.deployUUPSProxy(
            "Whitelist.sol",
            abi.encodeCall(Whitelist.initialize, (msg.sender))
        );
        Whitelist whitelist = Whitelist(payable(whitelistProxy));
        console.log("Whitelist deployed to:", address(whitelist));

        whitelist.addToWhitelist(address(gateway));
        console.log("Settlement gateway has been whitelisted. Gateway contract address:", address(gateway));

        string memory jsonOutput = string.concat(
            '{"settlement_gateway_addr": "',
            Strings.toHexString(address(gateway)),
            '", "whitelist_addr": "',
            Strings.toHexString(address(whitelist)),
            '"}'
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
            '{"l1_gateway_addr": "',
            Strings.toHexString(address(gateway)),
            '"}'
        );
        console.log("JSON_DEPLOY_ARTIFACT:", jsonOutput);

        vm.stopBroadcast();
    }
}
