// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;
import "forge-std/Script.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/PreConfCommitmentStore.sol";
import "../contracts/Oracle.sol";
import "../contracts/Whitelist.sol";
import "../contracts/validator-registry/ValidatorRegistry.sol";
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

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (msg.sender))
        );
        BlockTracker blockTracker = BlockTracker(payable(blockTrackerProxy));
        console.log("BlockTracker deployed to:", address(blockTracker));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, (minStake, feeRecipient, feePercent, msg.sender, address(blockTracker)))
        );
        BidderRegistry bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));
        console.log("BidderRegistry deployed to:", address(bidderRegistry));

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, (minStake, feeRecipient, feePercent, msg.sender))
        );
        ProviderRegistry providerRegistry = ProviderRegistry(payable(providerRegistryProxy));
        console.log("ProviderRegistry deployed to:", address(providerRegistry));

        address preconfCommitmentStoreProxy = Upgrades.deployUUPSProxy(
            "PreConfCommitmentStore.sol",
            abi.encodeCall(PreConfCommitmentStore.initialize, (address(providerRegistry), address(bidderRegistry), feeRecipient, msg.sender, address(blockTracker), commitmentDispatchWindow))
        );
        PreConfCommitmentStore preConfCommitmentStore = PreConfCommitmentStore(payable(preconfCommitmentStoreProxy));
        console.log("PreConfCommitmentStore deployed to:", address(preConfCommitmentStore));

        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        console.log(
            "ProviderRegistry updated with PreConfCommitmentStore address:",
            address(preConfCommitmentStore)
        );

        bidderRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        console.log(
            "BidderRegistry updated with PreConfCommitmentStore address:",
            address(preConfCommitmentStore)
        );

        address oracleProxy = Upgrades.deployUUPSProxy(
            "Oracle.sol",
            abi.encodeCall(Oracle.initialize, (address(preConfCommitmentStore), address(blockTracker), msg.sender))
        );
        Oracle oracle = Oracle(payable(oracleProxy));
        console.log("Oracle deployed to:", address(oracle));

        preConfCommitmentStore.updateOracle(address(oracle));
        console.log(
            "PreConfCommitmentStore updated with Oracle address:",
            address(oracle)
        );

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

// Deploys ValidatorRegistry contract via UUPS proxy
contract DeployValidatorRegistry is Script {
    function run() external {
        vm.startBroadcast();

        // 7000 blocks @ 200ms per block = 23.3 min. This allows two L1 epochs (finalization time) + settlement buffer,
        // to pass between validator unstake initiation and withdrawal.
        uint256 unstakePeriodBlocks = 7000;

        // Can later be upgraded with https://docs.openzeppelin.com/upgrades-plugins/1.x/api-foundry-upgrades#Upgrades-upgradeProxy-address-string-bytes-
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorRegistry.sol",
            abi.encodeCall(
                ValidatorRegistry.initialize,
                (3 ether, unstakePeriodBlocks, msg.sender)
            )
        );
        console.log(
            "ValidatorRegistry UUPS proxy deployed to:",
            address(proxy)
        );

        ValidatorRegistry validatorRegistry = ValidatorRegistry(payable(proxy));
        console.log("ValidatorRegistry owner:", validatorRegistry.owner());

        vm.stopBroadcast();
    }
}
