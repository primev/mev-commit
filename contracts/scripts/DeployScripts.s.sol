// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;
import "forge-std/Script.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/PreConfirmations.sol";
import "../contracts/Oracle.sol";
import "../contracts/Whitelist.sol";
import "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import "../contracts/ValidatorRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import "../contracts/BlockTracker.sol";

// Deploy scripts should inherit this contract if they deploy using create2 deterministic addrs.
contract Create2Deployer {
    address constant _CREATE2_PROXY =
        0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant _EXPECTED_DEPLOYER =
        0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;

    function _checkCreate2Deployed() internal view {
        require(
            isContractDeployed(_CREATE2_PROXY),
            "Create2 proxy needs to be deployed. See https://github.com/primevprotocol/deterministic-deployment-proxy"
        );
    }

    function _checkDeployer() internal view {
        if (msg.sender != _EXPECTED_DEPLOYER) {
            console.log(
                "Warning: deployer is not expected address of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266. Contracts addresses will not match documentation"
            );
        }
    }

    function isContractDeployed(address addr) public view returns (bool) {
        uint size;
        assembly {
            size := extcodesize(addr)
        }
        return size > 0;
    }
}

// Deploys core contracts
contract DeployScript is Script, Create2Deployer {
    function run() external {
        vm.startBroadcast();

        _checkCreate2Deployed();
        _checkDeployer();

        // Replace these with your contract's constructor parameters
        uint256 minStake = 1 ether;
        address feeRecipient = address(
            0x68bC10674b265f266b4b1F079Fa06eF4045c3ab9
        );
        uint16 feePercent = 2;
        uint64 commitmentDispatchWindow = 2000;

        uint256 blocksPerWindow = 10;
        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        BlockTracker blockTracker = new BlockTracker{salt: salt}(msg.sender);
        console.log("BlockTracker deployed to:", address(blockTracker));

        BidderRegistry bidderRegistry = new BidderRegistry{salt: salt}(
            minStake,
            feeRecipient,
            feePercent,
            msg.sender,
            address(blockTracker)
        );
        console.log("BidderRegistry deployed to:", address(bidderRegistry));

        ProviderRegistry providerRegistry = new ProviderRegistry{salt: salt}(
            minStake,
            feeRecipient,
            feePercent,
            msg.sender
        );
        console.log("ProviderRegistry deployed to:", address(providerRegistry));

        PreConfCommitmentStore preConfCommitmentStore = new PreConfCommitmentStore{
                salt: salt
            }(
                address(providerRegistry),
                address(bidderRegistry),
                feeRecipient,
                msg.sender,
                address(blockTracker),
                commitmentDispatchWindow
            );
        console.log(
            "PreConfCommitmentStore deployed to:",
            address(preConfCommitmentStore)
        );

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

        Oracle oracle = new Oracle{salt: salt}(
            address(preConfCommitmentStore),
            address(blockTracker),
            msg.sender
        );
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
contract DeployWhitelist is Script, Create2Deployer {
    function run() external {
        console.log(
            "Warning: DeployWhitelist is deprecated and only for backwards compatibility with hyperlane"
        );

        address expectedWhiteListAddr = 0x57508f0B0f3426758F1f3D63ad4935a7c9383620;
        if (isContractDeployed(expectedWhiteListAddr)) {
            console.log(
                "Whitelist already deployed to:",
                expectedWhiteListAddr
            );
            return;
        }

        vm.startBroadcast();

        _checkCreate2Deployed();
        _checkDeployer();

        address hypERC20Addr = vm.envAddress("HYP_ERC20_ADDR");
        require(
            hypERC20Addr != address(0),
            "Address to whitelist not provided"
        );

        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        Whitelist whitelist = new Whitelist{salt: salt}(msg.sender);
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

        ValidatorRegistry validatorRegistry = ValidatorRegistry(proxy);
        console.log("ValidatorRegistry owner:", validatorRegistry.owner());

        vm.stopBroadcast();
    }
}
