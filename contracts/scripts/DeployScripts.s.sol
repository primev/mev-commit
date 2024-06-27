// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;
import "forge-std/Script.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/PreConfCommitmentStore.sol";
import "../contracts/Oracle.sol";
import "../contracts/Whitelist.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import "../contracts/BlockTracker.sol";

// Deploys core contracts
contract DeployScript is Script {
    function run() external {
        vm.startBroadcast();

        // Replace these with your contract's constructor parameters
        uint256 minStake = 1 ether;
        address feeRecipient = address(
            0x68bC10674b265f266b4b1F079Fa06eF4045c3ab9
        );
        uint16 feePercent = 2;
        uint64 commitmentDispatchWindow = 2000;
        uint256 blocksPerWindow = 10;

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (msg.sender, blocksPerWindow))
        );
        BlockTracker blockTracker = BlockTracker(payable(blockTrackerProxy));
        console.log("BlockTracker:", address(blockTracker));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, (feeRecipient, feePercent, msg.sender, address(blockTracker), blocksPerWindow))
        );
        BidderRegistry bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));
        console.log("BidderRegistry:", address(bidderRegistry));

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, (minStake, feeRecipient, feePercent, msg.sender))
        );
        ProviderRegistry providerRegistry = ProviderRegistry(payable(providerRegistryProxy));
        console.log("ProviderRegistry:", address(providerRegistry));

        address preconfCommitmentStoreProxy = Upgrades.deployUUPSProxy(
            "PreConfCommitmentStore.sol",
            abi.encodeCall(PreConfCommitmentStore.initialize, (address(providerRegistry), address(bidderRegistry), feeRecipient, msg.sender, address(blockTracker), commitmentDispatchWindow, blocksPerWindow))
        );
        PreConfCommitmentStore preConfCommitmentStore = PreConfCommitmentStore(payable(preconfCommitmentStoreProxy));
        console.log("PreConfCommitmentStore:", address(preConfCommitmentStore));

        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        console.log("ProviderRegistryWithPreConfCommitmentStore:", address(preConfCommitmentStore));

        bidderRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        console.log("BidderRegistryWithPreConfCommitmentStore:", address(preConfCommitmentStore));

        address oracleProxy = Upgrades.deployUUPSProxy(
            "Oracle.sol",
            abi.encodeCall(Oracle.initialize, (address(preConfCommitmentStore), address(blockTracker), msg.sender))
        );
        Oracle oracle = Oracle(payable(oracleProxy));
        console.log("Oracle:", address(oracle));

        preConfCommitmentStore.updateOracle(address(oracle));
        console.log("PreConfCommitmentStoreWithOracle:", address(oracle));

        vm.stopBroadcast();
    }
}

contract TransferOwnership is Script {
    function run() external {
        vm.startBroadcast();

        address oracleKeystoreAddress = vm.envAddress("ORACLE_KEYSTORE_ADDRESS");
        require(oracleKeystoreAddress != address(0), "Oracle keystore address not provided");

        address blockTrackerProxy = vm.envAddress("BLOCK_TRACKER_ADDRESS");
        require(blockTrackerProxy != address(0), "Block tracker not provided");
        BlockTracker blockTracker = BlockTracker(payable(blockTrackerProxy));
        blockTracker.transferOwnership(oracleKeystoreAddress);
        console.log("BlockTracker owner:", blockTracker.owner());

        address oracleProxy = vm.envAddress("ORACLE_ADDRESS");
        require(oracleProxy != address(0), "Oracle proxy not provided");
        Oracle oracle = Oracle(payable(oracleProxy));
        oracle.transferOwnership(oracleKeystoreAddress);
        console.log("Oracle owner:", oracle.owner());

        vm.stopBroadcast();
    }
}

// Deploys whitelist contract and adds HypERC20 to whitelist
contract DeployWhitelist is Script {
    function run() external {
        console.log(
            "Warning: DeployWhitelist is deprecated and only for backwards compatibility with hyperlane"
        );

        vm.startBroadcast();

        address hypERC20Addr = vm.envAddress("HYP_ERC20_ADDR");
        require(
            hypERC20Addr != address(0),
            "Address to whitelist not provided"
        );

        address whitelistProxy = Upgrades.deployUUPSProxy(
            "Whitelist.sol",
            abi.encodeCall(Whitelist.initialize, (msg.sender))
        );
        Whitelist whitelist = Whitelist(payable(whitelistProxy));
        console.log("Whitelist deployed to:", address(whitelist));

        whitelist.addToWhitelist(address(hypERC20Addr));
        console.log(
            "Whitelist updated with hypERC20 address:",
            address(hypERC20Addr)
        );

        vm.stopBroadcast();
    }
}
