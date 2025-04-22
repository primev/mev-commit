// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {Oracle} from "../../contracts/core/Oracle.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {console} from "forge-std/console.sol";

contract DeployCore is Script {

    // Amount of ETH to initially fund the oracle account on L1 chain.
    uint256 public constant ORACLE_INITIAL_FUNDING = 0.5 ether;
    uint256 public constant PERCENT_MULTIPLIER = 1e16;
    
    error FailedToSendETHToOracle(address addr);

    function run() external {
        vm.startBroadcast();

        uint256 minStake = 10 ether;
        address protocolFeeRecipient = address(
            0xfA0B0f5d298d28EFE4d35641724141ef19C05684 // Placeholder for now, L1 preconf.eth address
        );
        uint256 feePercent = 2 * PERCENT_MULTIPLIER; // 2%
        uint256 providerPenaltyPercent = 5 * PERCENT_MULTIPLIER; // 5%
        uint64 commitmentDispatchWindow = 300;
        uint256 withdrawalDelay = 24 hours  * 1000; // 24 hours in milliseconds
        uint256 protocolFeePayoutPeriodBlocks = 5 hours ; // 1 hour with 200ms blocks
        address oracleKeystoreAddress = vm.envAddress("ORACLE_KEYSTORE_ADDRESS");
        require(oracleKeystoreAddress != address(0), "missing Oracle keystore address");

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize,
            (oracleKeystoreAddress, // oracleAccount_ param
            msg.sender)) // owner_ param
        );
        BlockTracker blockTracker = BlockTracker(payable(blockTrackerProxy));
        console.log("BlockTracker:", address(blockTracker));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize,
            (protocolFeeRecipient, // _protocolFeeRecipient param
            feePercent, // _feePercent param
            msg.sender, // _owner param
            address(blockTracker), // _blockTracker param
            protocolFeePayoutPeriodBlocks)) // _protocolFeePayoutPeriodBlocks param
        );
        BidderRegistry bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));
        console.log("BidderRegistry:", address(bidderRegistry));

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize,
            (minStake, // _minStake param
            protocolFeeRecipient, // _protocolFeeRecipient param
            providerPenaltyPercent, // _feePercent param
            msg.sender, // _owner param
            withdrawalDelay, // _withdrawalDelay param
            protocolFeePayoutPeriodBlocks)) // _protocolFeePayoutPeriodBlocks param
        );
        ProviderRegistry providerRegistry = ProviderRegistry(payable(providerRegistryProxy));
        console.log("ProviderRegistry:", address(providerRegistry));

        address preconfCommitmentStoreProxy = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(PreconfManager.initialize,
            (address(providerRegistry), // _providerRegistry param
            address(bidderRegistry), // _bidderRegistry param
            address(0x0), // _oracleContract param, updated later in script. 
            msg.sender, // _owner param
            address(blockTracker), // _blockTracker param
            commitmentDispatchWindow)) // _commitmentDispatchWindow param
        );
        PreconfManager preconfManager = PreconfManager(payable(preconfCommitmentStoreProxy));
        console.log("PreconfManager:", address(preconfManager));

        providerRegistry.setPreconfManager(address(preconfManager));
        console.log("_ProviderRegistryWithPreconfManager:", address(preconfManager));

        bidderRegistry.setPreconfManager(address(preconfManager));
        console.log("_BidderRegistryWithPreconfManager:", address(preconfManager));

        address oracleProxy = Upgrades.deployUUPSProxy(
            "Oracle.sol",
            abi.encodeCall(Oracle.initialize,
            (address(preconfManager), // preConfContract_ param
            address(blockTracker), // blockTrackerContract_ param
            oracleKeystoreAddress, // oracleAcount_ param
            msg.sender)) // owner_ param
        );
        Oracle oracle = Oracle(payable(oracleProxy));
        console.log("Oracle:", address(oracle));

        preconfManager.updateOracleContract(address(oracle));
        console.log("_PreconfManagerWithOracle:", address(oracle));

        (bool success, ) = payable(oracleKeystoreAddress).call{value: ORACLE_INITIAL_FUNDING}("");
        require(success, FailedToSendETHToOracle(oracleKeystoreAddress));
        
        blockTracker.setProviderRegistry(address(providerRegistry));

        vm.stopBroadcast();
    }
}
