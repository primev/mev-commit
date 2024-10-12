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

contract BridgeBase {
    // Amount of ETH which must be allocated only to the contract deployer on mev-commit chain genesis.
    uint256 public constant DEPLOYER_GENESIS_ALLOCATION = type(uint256).max - 10 ether;

    // Amount of ETH to initially fund the relayer account on both chains.
    uint256 public constant RELAYER_INITIAL_FUNDING = 1 ether;

    // Amount of ETH required by the contract deployer to setup all bridge and core contract state,
    // AND initially fund the relayer, all on mev-commit chain.
    // This amount of ETH must be initially locked in the L1 gateway contract to ensure a 1:1 peg
    // between mev-commit chain ETH and L1 ETH.
    uint256 public constant MEV_COMMIT_CHAIN_SETUP_COST = 1 ether + RELAYER_INITIAL_FUNDING;

    error RelayerAddressNotSet(address addr);
    error FailedToSendETHToRelayer(address addr);
}

contract DeploySettlementGatewayWithAllocator is Script, BridgeBase {
    error DeployerMustHaveGenesisAllocation(uint256 balance, uint256 expected);

    function run() external {

        vm.startBroadcast();

        address relayerAddr = vm.envAddress("RELAYER_ADDRESS");
        require(relayerAddr != address(0), RelayerAddressNotSet(relayerAddr));

        require(address(msg.sender).balance == DEPLOYER_GENESIS_ALLOCATION, 
            DeployerMustHaveGenesisAllocation(address(msg.sender).balance, DEPLOYER_GENESIS_ALLOCATION));

        address allocatorProxy = Upgrades.deployUUPSProxy(
            "Allocator.sol",
            abi.encodeCall(Allocator.initialize, (msg.sender))
        );
        Allocator allocator = Allocator(payable(allocatorProxy));
        console.log("Allocator deployed to:", address(allocator));

        // TODO: fund allocator w/ DEPLOYER_GENESIS_ALLOCATION - MEV_COMMIT_CHAIN_SETUP_COST.

        address sgProxy = Upgrades.deployUUPSProxy(
            "SettlementGateway.sol",
            abi.encodeCall(SettlementGateway.initialize,
            (allocatorProxy,
            msg.sender, // Owner
            relayerAddr,
            1, 1)) // Fees set to 1 wei for now. TODO: Change this in the revamp PR.
        );
        SettlementGateway gateway = SettlementGateway(payable(sgProxy));
        console.log("Standard bridge gateway for settlement chain deployed to:",
            address(gateway));

        allocator.addToWhitelist(address(gateway));
        console.log("Settlement gateway has been whitelisted");

        string memory jsonOutput = string.concat(
            "{'settlement_gateway_addr': '",
            Strings.toHexString(address(gateway)),
            "', 'allocator_addr': '",
            Strings.toHexString(address(allocator)),
            "'}"
        );
        console.log("JSON_DEPLOY_ARTIFACT:", jsonOutput); 

        (bool success, ) = payable(relayerAddr).call{value: RELAYER_INITIAL_FUNDING}("");
        require(success, FailedToSendETHToRelayer(relayerAddr));

        vm.stopBroadcast();
    }
}

contract DeployL1Gateway is Script, BridgeBase {
    error DeployerMustHaveEnoughFunds(uint256 balance, uint256 expected);

    function run() external {

        vm.startBroadcast();

        address relayerAddr = vm.envAddress("RELAYER_ADDRESS");
        require(relayerAddr != address(0), RelayerAddressNotSet(relayerAddr));

        // Caller needs funds to cover mev-commit chain setup cost AND initial relayer funding on L1.
        require(address(msg.sender).balance == MEV_COMMIT_CHAIN_SETUP_COST + RELAYER_INITIAL_FUNDING,
            DeployerMustHaveEnoughFunds(address(msg.sender).balance, MEV_COMMIT_CHAIN_SETUP_COST + RELAYER_INITIAL_FUNDING));

        address l1gProxy = Upgrades.deployUUPSProxy(
            "L1Gateway.sol",
            abi.encodeCall(L1Gateway.initialize,
            // TODO: Owner on prod L1 will need to be primev multisig. Account for this somehow.
            (msg.sender, // Owner
            relayerAddr,
            1, 1)) // Fees set to 1 wei for now. TODO: Change this in the revamp PR.
        );
        L1Gateway gateway = L1Gateway(payable(l1gProxy));
        console.log("Standard bridge gateway for l1 deployed to:", address(gateway));

        // TODO: transfer CHAIN_SETUP_COST ETH to L1 gateway to ensure 1:1 peg.
        
        string memory jsonOutput = string.concat(
            "{'l1_gateway_addr': '",
            Strings.toHexString(address(gateway)),
            "'}"
        );
        console.log("JSON_DEPLOY_ARTIFACT:", jsonOutput);

        (bool success, ) = payable(relayerAddr).call{value: RELAYER_INITIAL_FUNDING}("");
        require(success, FailedToSendETHToRelayer(relayerAddr));

        vm.stopBroadcast();
    }
}
