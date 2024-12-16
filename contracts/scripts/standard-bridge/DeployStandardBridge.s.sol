// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {SettlementGateway} from "../../contracts/standard-bridge/SettlementGateway.sol";
import {L1Gateway} from "../../contracts/standard-bridge/L1Gateway.sol";
import {Allocator} from "../../contracts/standard-bridge/Allocator.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";

contract BridgeBase is Script {
    // Amount of ETH which must be allocated only to the contract deployer on mev-commit chain genesis.
    uint256 public constant DEPLOYER_GENESIS_ALLOCATION = 5 ether;

    // Amount of ETH to initially fund the relayer account on both chains.
    uint256 public constant RELAYER_INITIAL_FUNDING = 1 ether;

    // Amount of ETH to initially fund the oracle account on L1 chain.
    uint256 public constant ORACLE_INITIAL_FUNDING = 1 ether;

    // Amount of ETH required by the contract deployer to initialize all bridge and core contract state,
    // AND initially fund the relayer, all on mev-commit chain.
    // This amount of ETH must be initially locked in the L1 gateway contract to ensure a 1:1 peg
    // between mev-commit chain ETH and L1 ETH.
    uint256 public constant MEV_COMMIT_CHAIN_SETUP_COST = 1 ether + RELAYER_INITIAL_FUNDING + ORACLE_INITIAL_FUNDING;

    // Amount of ETH required on L1 to initialize the L1 gateway, make transfer calls, and initially fund the relayer on L1.
    uint256 public constant L1_SETUP_COST = 1 ether + RELAYER_INITIAL_FUNDING;

    error RelayerAddressNotSet(address addr);
    error L1FinalizationFeeNotSet(uint256 fee);
    error SettlementFinalizationFeeNotSet(uint256 fee);
    error FailedToSendETHToRelayer(address addr);

    function _getRelayerAddress() internal view returns (address relayerAddr) {
        relayerAddr = vm.envAddress("RELAYER_ADDRESS");
        require(relayerAddr != address(0), RelayerAddressNotSet(relayerAddr));
    }

    function _getL1FinalizationFee() internal view returns (uint256 l1FinalizationFee) {
        l1FinalizationFee = vm.envUint("L1_FINALIZATION_FEE");
        require(l1FinalizationFee != 0, L1FinalizationFeeNotSet(l1FinalizationFee));
    }

    function _getSettlementFinalizationFee() internal view returns (uint256 settlementFinalizationFee) {
        settlementFinalizationFee = vm.envUint("SETTLEMENT_FINALIZATION_FEE");
        require(settlementFinalizationFee != 0, SettlementFinalizationFeeNotSet(settlementFinalizationFee));
    }
}

contract DeploySettlementGateway is BridgeBase {
    error DeployerMustHaveGenesisAllocation(uint256 balance, uint256 expected);
    error FailedToFundAllocator(address allocator);

    function run() external {

        vm.startBroadcast();

        address relayerAddr = _getRelayerAddress();
        uint256 l1FinalizationFee = _getL1FinalizationFee();

        require(address(msg.sender).balance >= DEPLOYER_GENESIS_ALLOCATION,
            DeployerMustHaveGenesisAllocation(address(msg.sender).balance, DEPLOYER_GENESIS_ALLOCATION));

        address allocatorProxy = Upgrades.deployUUPSProxy(
            "Allocator.sol",
            abi.encodeCall(Allocator.initialize,
                (msg.sender)) // Owner
        );
        Allocator allocator = Allocator(payable(allocatorProxy));
        console.log("_Allocator:", address(allocator));

        (bool success, ) = payable(address(allocator)).call{value: DEPLOYER_GENESIS_ALLOCATION - MEV_COMMIT_CHAIN_SETUP_COST}("");
        require(success, FailedToFundAllocator(address(allocator)));

        address sgProxy = Upgrades.deployUUPSProxy(
            "SettlementGateway.sol",
            abi.encodeCall(SettlementGateway.initialize,
                (allocatorProxy,
                    msg.sender, // Owner
                    relayerAddr,
                    l1FinalizationFee)) // SettlementGateway._counterpartyFinalizationFee
        );
        SettlementGateway settlementGateway = SettlementGateway(payable(sgProxy));
        console.log("SettlementGateway:", address(settlementGateway));

        allocator.addToWhitelist(address(settlementGateway));

        (success, ) = payable(relayerAddr).call{value: RELAYER_INITIAL_FUNDING}("");
        require(success, FailedToSendETHToRelayer(relayerAddr));

        vm.stopBroadcast();
    }
}

contract DeployL1Gateway is BridgeBase {
    error L1OwnerAddressNotSet(address addr);
    error DeployerMustHaveEnoughFunds(uint256 balance, uint256 expected);
    error FailedToFundL1Gateway(address gateway);

    function run() external {

        vm.startBroadcast();

        address owner = _getL1OwnerAddress(); // On mainnet, this must be the primev multisig.
        address relayerAddr = _getRelayerAddress();
        uint256 settlementFinalizationFee = _getSettlementFinalizationFee();

        // Caller needs funds to lock ETH w.r.t mev-commit chain setup cost, and ETH for L1 setup cost.
        require(address(msg.sender).balance >= MEV_COMMIT_CHAIN_SETUP_COST + L1_SETUP_COST,
            DeployerMustHaveEnoughFunds(address(msg.sender).balance, MEV_COMMIT_CHAIN_SETUP_COST + L1_SETUP_COST));

        address l1gProxy = Upgrades.deployUUPSProxy(
            "L1Gateway.sol",
            abi.encodeCall(L1Gateway.initialize,
                (owner, // Owner
                    relayerAddr,
                    settlementFinalizationFee)) // L1Gateway._counterpartyFinalizationFee
        );
        L1Gateway l1Gateway = L1Gateway(payable(l1gProxy));
        console.log("L1Gateway:", address(l1Gateway));

        (bool success, ) = payable(address(l1Gateway)).call{value: MEV_COMMIT_CHAIN_SETUP_COST}("");
        require(success, FailedToFundL1Gateway(address(l1Gateway)));

        (success, ) = payable(relayerAddr).call{value: RELAYER_INITIAL_FUNDING}("");
        require(success, FailedToSendETHToRelayer(relayerAddr));

        vm.stopBroadcast();
    }

    function _getL1OwnerAddress() internal view returns (address ownerAddr) {
        ownerAddr = vm.envAddress("L1_OWNER_ADDRESS");
        require(ownerAddr != address(0), L1OwnerAddressNotSet(ownerAddr));
    }
}
