// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;
import "forge-std/Script.sol";
import "contracts/BidderRegistry.sol";
import "contracts/ProviderRegistry.sol";
import "contracts/PreConfirmations.sol";
import "contracts/Oracle.sol";
import "contracts/Whitelist.sol";

// Deploy scripts should inherit this contract if they deploy using create2 deterministic addrs.
contract Create2Deployer {
    address constant create2Proxy = 0x4e59b44847b379578588920cA78FbF26c0B4956C;
    address constant expectedDeployer = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;

    function checkCreate2Deployed() internal view {
        require(isContractDeployed(create2Proxy), "Create2 proxy needs to be deployed. See https://github.com/primevprotocol/deterministic-deployment-proxy");
    }

    function checkDeployer() internal view {
        if (msg.sender != expectedDeployer) {
            console.log("Warning: deployer is not expected address of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266. Contracts addresses will not match documentation");
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

        checkCreate2Deployed();
        checkDeployer();

        // Replace these with your contract's constructor parameters
        uint256 minStake = 1 ether;
        address feeRecipient = address(0x68bC10674b265f266b4b1F079Fa06eF4045c3ab9);
        uint16 feePercent = 2;
        uint256 nextRequestedBlockNumber = 4958905;

        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        BidderRegistry bidderRegistry = new BidderRegistry{salt: salt}(minStake, feeRecipient, feePercent, msg.sender);
        console.log("BidderRegistry deployed to:", address(bidderRegistry));

        ProviderRegistry providerRegistry = new ProviderRegistry{salt: salt}(minStake, feeRecipient, feePercent, msg.sender);
        console.log("ProviderRegistry deployed to:", address(providerRegistry));

        PreConfCommitmentStore preConfCommitmentStore = new PreConfCommitmentStore{salt: salt}(address(providerRegistry), address(bidderRegistry), feeRecipient, msg.sender);
        console.log("PreConfCommitmentStore deployed to:", address(preConfCommitmentStore));

        providerRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));
        console.log("ProviderRegistry updated with PreConfCommitmentStore address:", address(preConfCommitmentStore));

        bidderRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));
        console.log("BidderRegistry updated with PreConfCommitmentStore address:", address(preConfCommitmentStore));

        Oracle oracle = new Oracle{salt: salt}(address(preConfCommitmentStore), nextRequestedBlockNumber, msg.sender);
        console.log("Oracle deployed to:", address(oracle));

        preConfCommitmentStore.updateOracle(address(oracle));
        console.log("PreConfCommitmentStore updated with Oracle address:", address(oracle));

        vm.stopBroadcast();
    }
}

// Deploys whitelist contract and adds HypERC20 to whitelist
contract DeployWhitelist is Script, Create2Deployer {
    function run() external {

        console.log("Warning: DeployWhitelist is deprecated and only for backwards compatibility with hyperlane");

        address expectedWhiteListAddr = 0xcf59aDa3C5FBa545Cc50FB9AEAe83D37b46F6E1B;
        if (isContractDeployed(expectedWhiteListAddr)) {
            console.log("Whitelist already deployed to:", expectedWhiteListAddr);
            return;
        }

        vm.startBroadcast();

        checkCreate2Deployed();
        checkDeployer();

        address hypERC20Addr = vm.envAddress("HYP_ERC20_ADDR");
        require(hypERC20Addr != address(0), "Address to whitelist not provided");

        // Forge deploy with salt uses create2 proxy from https://github.com/primevprotocol/deterministic-deployment-proxy
        bytes32 salt = 0x8989000000000000000000000000000000000000000000000000000000000000;

        Whitelist whitelist = new Whitelist{salt: salt}(msg.sender);
        console.log("Whitelist deployed to:", address(whitelist));

        whitelist.addToWhitelist(address(hypERC20Addr));
        console.log("Whitelist updated with hypERC20 address:", address(hypERC20Addr));

        vm.stopBroadcast();
    }
}
