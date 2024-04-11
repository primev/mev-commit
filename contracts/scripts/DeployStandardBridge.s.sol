// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;
import "forge-std/Script.sol";
import {Create2Deployer} from "../scripts/DeployScripts.s.sol";
import {SettlementGateway} from "../contracts/standard-bridge/SettlementGateway.sol";
import {L1Gateway} from "../contracts/standard-bridge/L1Gateway.sol";
import {Whitelist} from "../contracts/Whitelist.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

contract DeploySettlementGateway is Script, Create2Deployer {
    function run() external {

        vm.startBroadcast();

        checkCreate2Deployed();
        checkDeployer();

        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        address expectedWhitelistAddr = 0x57508f0B0f3426758F1f3D63ad4935a7c9383620;
        if (isContractDeployed(expectedWhitelistAddr)) {
            console.log("Whitelist must not be deployed to execute DeploySettlementGateway script. Exiting...");
            return;
        }

        address relayerAddr = vm.envAddress("RELAYER_ADDR");
        SettlementGateway gateway = new SettlementGateway{salt: salt}(
            expectedWhitelistAddr,
            msg.sender, // Owner
            relayerAddr,
            1, 1); // Fees set to 1 wei for now
        console.log("Standard bridge gateway for settlement chain deployed to:",
            address(gateway));

        Whitelist whitelist = new Whitelist{salt: salt}(msg.sender);
        console.log("Whitelist deployed to:", address(whitelist));

        if (!isContractDeployed(expectedWhitelistAddr)) {
            console.log("Whitelist not deployed to expected address:", expectedWhitelistAddr);
            return;
        }

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

contract DeployL1Gateway is Script, Create2Deployer {
    function run() external {

        vm.startBroadcast();

        checkCreate2Deployed();
        checkDeployer();

        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        address relayerAddr = vm.envAddress("RELAYER_ADDR");

        L1Gateway gateway = new L1Gateway{salt: salt}(
            msg.sender, // Owner
            relayerAddr,
            1, 1); // Fees set to 1 wei for now
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
