// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {RegistryMock} from "./RegistryMock.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {TimestampOccurrence} from "../../../contracts/utils/Occurrence.sol";

contract MevCommitMiddlewareTest is Test {

    RegistryMock public networkRegistryMock;
    RegistryMock public operatorRegistryMock;
    RegistryMock public vaultFactoryMock;
    address public network;
    uint256 public slashPeriodSeconds;
    address public slashOracle;
    address public owner;

    MevCommitMiddleware public mevCommitMiddleware;

    event OperatorRegistered(address indexed operator);
    event OperatorDeregistrationRequested(address indexed operator);
    event OperatorDeregistered(address indexed operator);
    event OperatorBlacklisted(address indexed operator);
    event ValRecordAdded(bytes blsPubkey, address indexed msgSender, uint256 indexed position);
    event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position);
    event ValRecordDeleted(bytes blsPubkey, address indexed msgSender);
    event VaultRegistered(address indexed vault, uint256 slashAmount);
    event VaultSlashAmountUpdated(address indexed vault, uint256 slashAmount);
    event VaultDeregistrationRequested(address indexed vault);
    event VaultDeregistered(address indexed vault);
    event ValidatorSlashed(bytes blsPubkey, address indexed operator, uint256 indexed position);
    event NetworkRegistrySet(address networkRegistry);
    event OperatorRegistrySet(address operatorRegistry);
    event VaultFactorySet(address vaultFactory);
    event NetworkSet(address network);
    event SlashPeriodSecondsSet(uint256 slashPeriodSeconds);
    event SlashOracleSet(address slashOracle);

    function setUp() public {
        networkRegistryMock = new RegistryMock();
        operatorRegistryMock = new RegistryMock();
        vaultFactoryMock = new RegistryMock();

        network = vm.addr(0x1);
        slashPeriodSeconds = 150;
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
                slashPeriodSeconds,
                slashOracle,
                owner
            ))
        );
        mevCommitMiddleware = MevCommitMiddleware(payable(proxy));
    }

    function getOperatorRecord(address operator) internal view
        returns (IMevCommitMiddleware.OperatorRecord memory) {
        (TimestampOccurrence.Occurrence memory occurrence, bool exists, bool isBlacklisted) =
            mevCommitMiddleware.operatorRecords(operator);
        return IMevCommitMiddleware.OperatorRecord(occurrence, exists, isBlacklisted);
    }

    function getVaultRecord(address vault) internal view
        returns (IMevCommitMiddleware.VaultRecord memory) {
        (bool exists, TimestampOccurrence.Occurrence memory occurrence, uint256 slashAmount) =
            mevCommitMiddleware.vaultRecords(vault);
        return IMevCommitMiddleware.VaultRecord(exists, occurrence, slashAmount);
    }

    function getValidatorRecord(bytes memory blsPubkey) internal view
        returns (IMevCommitMiddleware.ValidatorRecord memory) {
        (address vault, address operator, bool exists, TimestampOccurrence.Occurrence memory occurrence) =
            mevCommitMiddleware.validatorRecords(blsPubkey);
        return IMevCommitMiddleware.ValidatorRecord(vault, operator, exists, occurrence);
    }

    function test_setters() public {
        assertEq(address(mevCommitMiddleware.networkRegistry()), address(networkRegistryMock));
        assertEq(address(mevCommitMiddleware.operatorRegistry()), address(operatorRegistryMock));
        assertEq(address(mevCommitMiddleware.vaultFactory()), address(vaultFactoryMock));
        assertEq(mevCommitMiddleware.network(), network);
        assertEq(mevCommitMiddleware.slashPeriodSeconds(), slashPeriodSeconds);
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

        uint256 newSlashPeriodSeconds = 204;
        vm.prank(owner);
        mevCommitMiddleware.setSlashPeriodSeconds(newSlashPeriodSeconds);
        assertEq(mevCommitMiddleware.slashPeriodSeconds(), newSlashPeriodSeconds);

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

    function test_registerOperators() public {
        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, operator1)
        );
        mevCommitMiddleware.registerOperators(operators);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotEntity.selector, operator1)
        );
        mevCommitMiddleware.registerOperators(operators);

        vm.prank(operator1);
        operatorRegistryMock.register();

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotEntity.selector, operator2)
        );
        mevCommitMiddleware.registerOperators(operators);

        vm.prank(operator2);
        operatorRegistryMock.register();

        vm.expectEmit(true, true, true, true);
        emit OperatorRegistered(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorRegistered(operator2);
        vm.prank(owner);
        mevCommitMiddleware.registerOperators(operators);

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.deregRequestOccurrence.exists, false);
        assertEq(operatorRecord2.deregRequestOccurrence.exists, false);
        assertEq(operatorRecord1.isBlacklisted, false);
        assertEq(operatorRecord2.isBlacklisted, false);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorAlreadyRegistered.selector, operator1)
        );
        mevCommitMiddleware.registerOperators(operators);
    }

    function test_requestOperatorDeregistrations() public {

        vm.warp(44);

        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, operator1)
        );
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        test_registerOperators();

        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistrationRequested(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistrationRequested(operator2);
        vm.prank(owner);
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.deregRequestOccurrence.exists, true);
        assertEq(operatorRecord2.deregRequestOccurrence.exists, true);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorDeregRequestExists.selector, operator1)
        );
        mevCommitMiddleware.requestOperatorDeregistrations(operators);
    }

    function test_deregisterOperators() public {

        vm.warp(10);

        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, operator1)
        );
        mevCommitMiddleware.deregisterOperators(operators);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.deregisterOperators(operators);

        test_registerOperators();

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotReadyToDeregister.selector,
            operator1, block.timestamp, 0)
        );
        mevCommitMiddleware.deregisterOperators(operators);

        vm.prank(owner);
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        vm.warp(10 + mevCommitMiddleware.slashPeriodSeconds() + 1);

        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistered(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistered(operator2);
        vm.prank(owner);
        mevCommitMiddleware.deregisterOperators(operators);

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, false);
        assertEq(operatorRecord2.exists, false);
        assertEq(operatorRecord1.deregRequestOccurrence.exists, false);
        assertEq(operatorRecord2.deregRequestOccurrence.exists, false);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.deregisterOperators(operators);
    }

    function test_operatorRegistrationCycle() public {
        test_deregisterOperators();
        operatorRegistryMock = new RegistryMock();
        vm.prank(owner);
        mevCommitMiddleware.setOperatorRegistry(IRegistry(address(operatorRegistryMock)));
        test_registerOperators();
    }

    function test_blacklistNonRegisteredOperators() public {
        address operator1 = vm.addr(0x133333);
        address operator2 = vm.addr(0x133334);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        vm.prank(vm.addr(0x11888));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x11888))
        );
        mevCommitMiddleware.blacklistOperators(operators);

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, false);
        assertEq(operatorRecord2.exists, false);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator2);
        mevCommitMiddleware.blacklistOperators(operators);

        operatorRecord1 = getOperatorRecord(operator1);
        operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.isBlacklisted, true);
        assertEq(operatorRecord2.isBlacklisted, true);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorAlreadyBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.blacklistOperators(operators);
    }

    function test_blacklistRegisteredOperators() public {
        test_registerOperators();

        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.isBlacklisted, false);
        assertEq(operatorRecord2.isBlacklisted, false);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator2);
        mevCommitMiddleware.blacklistOperators(operators);

        operatorRecord1 = getOperatorRecord(operator1);
        operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.isBlacklisted, true);
        assertEq(operatorRecord2.isBlacklisted, true);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorAlreadyBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.blacklistOperators(operators);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorAlreadyRegistered.selector, operator1)
        );
        mevCommitMiddleware.registerOperators(operators);
    }

    function test_blacklistReqDeregisteredOperators() public {
        test_requestOperatorDeregistrations();

        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.deregRequestOccurrence.exists, true);
        assertEq(operatorRecord2.deregRequestOccurrence.exists, true);
        assertEq(operatorRecord1.isBlacklisted, false);
        assertEq(operatorRecord2.isBlacklisted, false);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorBlacklisted(operator2);
        mevCommitMiddleware.blacklistOperators(operators);

        operatorRecord1 = getOperatorRecord(operator1);
        operatorRecord2 = getOperatorRecord(operator2);
        assertEq(operatorRecord1.exists, true);
        assertEq(operatorRecord2.exists, true);
        assertEq(operatorRecord1.deregRequestOccurrence.exists, true);
        assertEq(operatorRecord2.deregRequestOccurrence.exists, true);
        assertEq(operatorRecord1.isBlacklisted, true);
        assertEq(operatorRecord2.isBlacklisted, true);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorAlreadyBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.blacklistOperators(operators);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorIsBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        vm.warp(66);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotReadyToDeregister.selector, operator1, block.timestamp, 44)
        );
        mevCommitMiddleware.deregisterOperators(operators);

        vm.warp(block.timestamp + mevCommitMiddleware.slashPeriodSeconds() + 1);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorIsBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.deregisterOperators(operators);
    }
}
