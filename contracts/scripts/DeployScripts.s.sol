// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.20;

import {Script} from "forge-std/Script.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {ProviderRegistry} from "../contracts/ProviderRegistry.sol";
import {PreConfCommitmentStore} from "../contracts/PreConfCommitmentStore.sol";
import {Oracle} from "../contracts/Oracle.sol";
import {Whitelist} from "../contracts/Whitelist.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";
import {console} from "forge-std/console.sol";

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

        address oracleKeystoreAddress = vm.envAddress("ORACLE_KEYSTORE_ADDRESS");
        require(oracleKeystoreAddress != address(0), "Oracle keystore address not provided");

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (blocksPerWindow, oracleKeystoreAddress, msg.sender))
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
            abi.encodeCall(Oracle.initialize, (address(preConfCommitmentStore), address(blockTracker), oracleKeystoreAddress, msg.sender))
        );
        Oracle oracle = Oracle(payable(oracleProxy));
        console.log("Oracle:", address(oracle));

        preConfCommitmentStore.updateOracle(address(oracle));
        console.log("PreConfCommitmentStoreWithOracle:", address(oracle));

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
            "hypERC20 addr not provided"
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
