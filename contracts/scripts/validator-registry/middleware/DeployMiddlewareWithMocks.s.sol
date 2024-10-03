// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable no-console
// solhint-disable one-contract-per-file

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {MockDelegator} from "../../../test/validator-registry/middleware/MockDelegator.sol";
import {MockVault} from "../../../test/validator-registry/middleware/MockVault.sol";
import {RegistryMock} from "../../../test/validator-registry/middleware/RegistryMock.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";

contract DeployMiddlewareWithMocks is Script {
    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        vm.startBroadcast();

        RegistryMock networkRegistryMock = new RegistryMock();
        RegistryMock operatorRegistryMock = new RegistryMock();
        RegistryMock vaultFactoryMock = new RegistryMock();

        uint256 slashPeriodSeconds = 150;
        address network = msg.sender;
        address slashOracle = msg.sender;
        address owner = msg.sender;

        networkRegistryMock.register();

        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitMiddleware.sol",
            abi.encodeCall(MevCommitMiddleware.initialize, (
                IRegistry(networkRegistryMock), 
                IRegistry(operatorRegistryMock), 
                IRegistry(vaultFactoryMock), 
                network, 
                slashPeriodSeconds,
                slashOracle,
                owner
            ))
        );
        MevCommitMiddleware mevCommitMiddleware = MevCommitMiddleware(payable(proxy));
        console.log("MevCommitMiddleware deployed to:", address(mevCommitMiddleware));

        MockDelegator mockDelegator1 = new MockDelegator(15);
        MockDelegator mockDelegator2 = new MockDelegator(16);
        MockVault vault1 = new MockVault(address(mockDelegator1), address(0), 10);
        MockVault vault2 = new MockVault(address(mockDelegator2), address(0), 10);

        console.log("MockDelegator 1 deployed to:", address(mockDelegator1));
        console.log("MockDelegator 2 deployed to:", address(mockDelegator2));
        console.log("MockVault 1 deployed to:", address(vault1));
        console.log("MockVault 2 deployed to:", address(vault2));

        vm.stopBroadcast();
    }
}
