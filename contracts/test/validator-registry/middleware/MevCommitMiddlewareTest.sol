// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {RegistryMock} from "./RegistryMock.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract MevCommitMiddlewareTest is Test {

    RegistryMock public networkRegistryMock;
    RegistryMock public operatorRegistryMock;
    RegistryMock public vaultFactoryMock;
    address public network;
    uint256 public slashPeriodBlocks;
    address public slashOracle;
    address public owner;

    MevCommitMiddleware public mevCommitMiddleware;

    function setUp() public {
        networkRegistryMock = new RegistryMock();
        operatorRegistryMock = new RegistryMock();
        vaultFactoryMock = new RegistryMock();

        network = vm.addr(0x1);
        slashPeriodBlocks = 150;
        slashOracle = vm.addr(0x2);
        owner = vm.addr(0x3);

        // Network addr must be registered with the network registry
        vm.prank(network);
        networkRegistryMock.register();

        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitMiddleware.sol",
            abi.encodeCall(MevCommitMiddleware.initialize, (
                IRegistry(networkRegistryMock), 
                IRegistry(operatorRegistryMock), 
                IRegistry(vaultFactoryMock), 
                network, 
                slashPeriodBlocks,
                slashOracle,
                owner
            ))
        );
        mevCommitMiddleware = MevCommitMiddleware(payable(proxy));
    }

    function test_setters() public {
        assertEq(address(mevCommitMiddleware.networkRegistry()), address(networkRegistryMock));
        assertEq(address(mevCommitMiddleware.operatorRegistry()), address(operatorRegistryMock));
        assertEq(address(mevCommitMiddleware.vaultFactory()), address(vaultFactoryMock));
        assertEq(mevCommitMiddleware.network(), network);
        assertEq(mevCommitMiddleware.slashPeriodBlocks(), slashPeriodBlocks);
        assertEq(mevCommitMiddleware.slashOracle(), slashOracle);
        assertEq(mevCommitMiddleware.owner(), owner);

        IRegistry newNetworkRegistry = IRegistry(new RegistryMock());
        vm.prank(owner);
        mevCommitMiddleware.setNetworkRegistry(newNetworkRegistry);
        assertEq(address(mevCommitMiddleware.networkRegistry()), address(newNetworkRegistry));

        IRegistry newOperatorRegistry = IRegistry(vm.addr(0x1112));
        vm.prank(owner);
        mevCommitMiddleware.setOperatorRegistry(newOperatorRegistry);
        assertEq(address(mevCommitMiddleware.operatorRegistry()), address(newOperatorRegistry));

        IRegistry newVaultFactory = IRegistry(vm.addr(0x1113));
        vm.prank(owner);
        mevCommitMiddleware.setVaultFactory(newVaultFactory);
        assertEq(address(mevCommitMiddleware.vaultFactory()), address(newVaultFactory));

        // register NEW network with NEW network registry
        address newNetwork = vm.addr(0x1114);
        vm.prank(newNetwork);
        RegistryMock(address(newNetworkRegistry)).register(); 

        vm.prank(owner);
        mevCommitMiddleware.setNetwork(newNetwork);
        assertEq(mevCommitMiddleware.network(), newNetwork);

        uint256 newSlashPeriodBlocks = 204;
        vm.prank(owner);
        mevCommitMiddleware.setSlashPeriodBlocks(newSlashPeriodBlocks);
        assertEq(mevCommitMiddleware.slashPeriodBlocks(), newSlashPeriodBlocks);

        address newSlashOracle = vm.addr(0x1115);
        vm.prank(owner);
        mevCommitMiddleware.setSlashOracle(newSlashOracle);
        assertEq(mevCommitMiddleware.slashOracle(), newSlashOracle);

        address newOwner = vm.addr(0x1116);
        vm.prank(owner);
        mevCommitMiddleware.transferOwnership(newOwner);
        assertEq(mevCommitMiddleware.pendingOwner(), newOwner);

        vm.prank(newOwner);
        mevCommitMiddleware.acceptOwnership();
        assertEq(mevCommitMiddleware.owner(), newOwner);
    }
}
