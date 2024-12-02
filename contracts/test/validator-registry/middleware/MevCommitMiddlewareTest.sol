// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable func-name-mixedcase

import {Test} from "forge-std/Test.sol";
import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {RegistryMock} from "./RegistryMock.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {TimestampOccurrence} from "../../../contracts/utils/Occurrence.sol";
import {MockVault} from "./MockVault.sol";
import {MockVetoSlasher} from "./MockVetoSlasher.sol";
import {MockInstantSlasher} from "./MockInstantSlasher.sol";
import {MockDelegator} from "./MockDelegator.sol";
import {MockBurnerRouter} from "./MockBurnerRouter.sol";
import {Checkpoints} from "@openzeppelin/contracts/utils/structs/Checkpoints.sol";

contract MevCommitMiddlewareTest is Test {

    RegistryMock public networkRegistryMock;
    RegistryMock public operatorRegistryMock;
    RegistryMock public vaultFactoryMock;
    RegistryMock public burnerRouterFactoryMock;
    address public network;
    uint256 public slashPeriodSeconds;
    address public slashOracle;
    address public slashReceiver;
    uint256 public minBurnerRouterDelay;
    address public owner;

    MevCommitMiddleware public mevCommitMiddleware;

    MockDelegator public mockDelegator1;
    MockDelegator public mockDelegator2;
    MockVault public vault1;
    MockVault public vault2;
    MockBurnerRouter public mockBurnerRouter;
    MockBurnerRouter public mockBurnerRouter2;

    bytes public sampleValPubkey1 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254b3f";
    bytes public sampleValPubkey2 = hex"aca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes public sampleValPubkey3 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254777";

    bytes public sampleValPubkey4 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254888";
    bytes public sampleValPubkey5 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254999";
    bytes public sampleValPubkey6 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254aaa";
    bytes public sampleValPubkey7 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254bbb";

    bytes public sampleValPubkey8 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254ccc";

    event OperatorRegistered(address indexed operator);
    event OperatorDeregistrationRequested(address indexed operator);
    event OperatorDeregistered(address indexed operator);
    event OperatorBlacklisted(address indexed operator);
    event OperatorUnblacklisted(address indexed operator);
    event ValRecordAdded(bytes blsPubkey, address indexed operator, address indexed vault, uint256 indexed position);
    event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position);
    event ValRecordDeleted(bytes blsPubkey, address indexed msgSender);
    event VaultRegistered(address indexed vault, uint160 slashAmount);
    event VaultSlashAmountUpdated(address indexed vault, uint160 slashAmount);
    event VaultDeregistrationRequested(address indexed vault);
    event VaultDeregistered(address indexed vault);
    event ValidatorSlashed(bytes blsPubkey, address indexed operator, address indexed vault, uint256 slashedAmount);
    event NetworkRegistrySet(address networkRegistry);
    event OperatorRegistrySet(address operatorRegistry);
    event VaultFactorySet(address vaultFactory);
    event NetworkSet(address network);
    event SlashPeriodSecondsSet(uint256 slashPeriodSeconds);
    event SlashOracleSet(address slashOracle);
    event ValidatorPositionsSwapped(bytes[] blsPubkeys, address[] vaults, address[] operators, uint256[] newPositions);

    function setUp() public virtual {
        networkRegistryMock = new RegistryMock();
        operatorRegistryMock = new RegistryMock();
        vaultFactoryMock = new RegistryMock();
        burnerRouterFactoryMock = new RegistryMock();

        network = vm.addr(0x1);
        slashPeriodSeconds = 150 hours;
        slashOracle = vm.addr(0x2);
        slashReceiver = vm.addr(0x3);
        minBurnerRouterDelay = 2 days;
        owner = vm.addr(0x4);

        // Network addr must be registered with the network registry
        vm.prank(network);
        networkRegistryMock.register();

        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitMiddleware.sol",
            abi.encodeCall(MevCommitMiddleware.initialize, (
                IRegistry(networkRegistryMock), 
                IRegistry(operatorRegistryMock), 
                IRegistry(vaultFactoryMock), 
                IRegistry(burnerRouterFactoryMock),
                network, 
                slashPeriodSeconds,
                slashOracle,
                slashReceiver,
                minBurnerRouterDelay,
                owner
            ))
        );
        mevCommitMiddleware = MevCommitMiddleware(payable(proxy));

        mockDelegator1 = new MockDelegator(15);
        mockDelegator2 = new MockDelegator(16);
        uint48 epochDuration = 10 hours;
        address emptyBurner = address(0);
        vault1 = new MockVault(address(mockDelegator1), address(0), emptyBurner, epochDuration);
        vault2 = new MockVault(address(mockDelegator2), address(0), emptyBurner, epochDuration);
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

    function test_unblacklistOperators() public {

        vm.prank(vm.addr(0x1121));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x1121))
        );
        mevCommitMiddleware.unblacklistOperators(new address[](0));

        address operator1 = vm.addr(0x1117);
        address operator2 = vm.addr(0x1118);
        address[] memory operators = new address[](2);
        operators[0] = operator1;
        operators[1] = operator2;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.unblacklistOperators(operators);

        test_registerOperators();

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.unblacklistOperators(operators);

        vm.prank(owner);
        mevCommitMiddleware.blacklistOperators(operators);

        IMevCommitMiddleware.OperatorRecord memory operatorRecord1 = getOperatorRecord(operator1);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord2 = getOperatorRecord(operator2);
        assertTrue(operatorRecord1.exists);
        assertTrue(operatorRecord2.exists);
        assertTrue(operatorRecord1.isBlacklisted);
        assertTrue(operatorRecord2.isBlacklisted);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit OperatorUnblacklisted(operator1);
        vm.expectEmit(true, true, true, true);
        emit OperatorUnblacklisted(operator2);
        mevCommitMiddleware.unblacklistOperators(operators);

        operatorRecord1 = getOperatorRecord(operator1);
        operatorRecord2 = getOperatorRecord(operator2);
        assertTrue(operatorRecord1.exists);
        assertTrue(operatorRecord2.exists);
        assertFalse(operatorRecord1.isBlacklisted);
        assertFalse(operatorRecord2.isBlacklisted);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.unblacklistOperators(operators);
    }

    function test_registerVaults() public {
        vm.prank(vm.addr(0x1121));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x1121))
        );
        mevCommitMiddleware.registerVaults(new address[](0), new uint160[](0));

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        uint160[] memory slashAmounts = new uint160[](1);
        slashAmounts[0] = 20;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidArrayLengths.selector, 2, 1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        slashAmounts = new uint160[](2);
        slashAmounts[0] = 0;
        slashAmounts[1] = 20;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotEntity.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vm.prank(address(vault1));
        vaultFactoryMock.register();
        vm.prank(address(vault2));
        vaultFactoryMock.register();

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.SlashAmountMustBeNonZero.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        slashAmounts[0] = 15;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.UnknownDelegatorType.selector, vault1, 15)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        uint64 fullRestakeDelegatorType = 1;

        mockDelegator1.setType(fullRestakeDelegatorType);
        mockDelegator2.setType(fullRestakeDelegatorType);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.FullRestakeDelegatorNotSupported.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        uint64 networkRestakeDelegatorType = 0;
        uint64 operatorSpecificDelegatorType = 2;

        mockDelegator1.setType(networkRestakeDelegatorType);
        mockDelegator2.setType(operatorSpecificDelegatorType);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.SlasherNotSetForVault.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        uint256 vetoDuration = 5 hours;
        MockVetoSlasher mockSlasher1 = new MockVetoSlasher(77, address(77), vetoDuration, mockDelegator1, address(mevCommitMiddleware));
        MockInstantSlasher mockSlasher2 = new MockInstantSlasher(88, mockDelegator2);

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.UnknownSlasherType.selector, vault1, 77)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        uint64 vetoSlasherType = 1;

        mockSlasher1.setType(vetoSlasherType);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VetoSlasherMustHaveZeroResolver.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        mockSlasher1.setResolver(address(0));

        assertEq(10 hours, vault1.epochDuration());
        
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultEpochDuration.selector, vault1,
            10 hours - vetoDuration, 150 hours)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        MockVault(vault1).setEpochDuration(151 hours);
        MockVault(vault2).setEpochDuration(151 hours);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultEpochDuration.selector, vault1,
            146 hours, 150 hours)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        MockVault(vault2).setEpochDuration(157 hours);
        MockVault(vault1).setEpochDuration(157 hours);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultBurner.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        mockBurnerRouter = new MockBurnerRouter(15 minutes);
        mockBurnerRouter2 = new MockBurnerRouter(15 minutes);

        vm.prank(address(mockBurnerRouter));
        burnerRouterFactoryMock.register();
        vm.prank(address(mockBurnerRouter2));
        burnerRouterFactoryMock.register();

        vault1.setBurner(address(mockBurnerRouter));
        vault2.setBurner(address(mockBurnerRouter2));

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultBurner.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        mockBurnerRouter.setNetworkReceiver(network, vm.addr(0x1143242));
        mockBurnerRouter2.setNetworkReceiver(network, vm.addr(0x1143243));

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultBurner.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        mockBurnerRouter.setNetworkReceiver(network, slashReceiver);
        mockBurnerRouter2.setNetworkReceiver(network, slashReceiver);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultBurner.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        mockBurnerRouter.setDelay(3 days);
        mockBurnerRouter2.setDelay(3 days);

        // No more errors relevant to vault 1

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.UnknownSlasherType.selector, vault2, 88)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.UnknownSlasherType.selector, vault2, 88)
        );
        mevCommitMiddleware.wouldVaultBeValidWith(address(vault2), 150 hours);

        uint64 instantSlasherType = 0;

        mockSlasher2.setType(instantSlasherType);

        assertTrue(mevCommitMiddleware.wouldVaultBeValidWith(address(vault1), 150 hours));
        assertTrue(mevCommitMiddleware.wouldVaultBeValidWith(address(vault2), 150 hours));

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultRegistered(address(vault1), 15);
        vm.expectEmit(true, true, true, true);
        emit VaultRegistered(address(vault2), 20);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        IMevCommitMiddleware.VaultRecord memory vaultRecord1 = getVaultRecord(address(vault1));
        IMevCommitMiddleware.VaultRecord memory vaultRecord2 = getVaultRecord(address(vault2));
        assertTrue(vaultRecord1.exists);
        assertTrue(vaultRecord2.exists);
        assertFalse(vaultRecord1.deregRequestOccurrence.exists);
        assertFalse(vaultRecord2.deregRequestOccurrence.exists);
        uint160 latestSlashAmount1 = mevCommitMiddleware.getLatestSlashAmount(address(vault1));
        uint160 latestSlashAmount2 = mevCommitMiddleware.getLatestSlashAmount(address(vault2));
        assertEq(latestSlashAmount1, 15);
        assertEq(latestSlashAmount2, 20);
        assertEq(latestSlashAmount1, mevCommitMiddleware.getSlashAmountAt(address(vault1), block.timestamp));
        assertEq(latestSlashAmount2, mevCommitMiddleware.getSlashAmountAt(address(vault2), block.timestamp));

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultAlreadyRegistered.selector, vault1)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vaults = new address[](1);
        vaults[0] = address(vault2);
        slashAmounts = new uint160[](1);
        slashAmounts[0] = 88;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultAlreadyRegistered.selector, vault2)
        );
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidVaultEpochDuration.selector, vault1, 152 hours, 24 days)
        );
        mevCommitMiddleware.wouldVaultBeValidWith(address(vault1), 24 days);
    }

    function test_updateSlashAmount() public {
        vm.prank(vm.addr(0x1121));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x1121))
        );
        mevCommitMiddleware.updateSlashAmounts(new address[](0), new uint160[](0));

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        uint160[] memory slashAmounts = new uint160[](1);
        slashAmounts[0] = 777;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidArrayLengths.selector, 2, 1)
        );
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        slashAmounts = new uint160[](2);
        slashAmounts[0] = 0;
        slashAmounts[1] = 999;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault1)
        );
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        test_registerVaults();

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.SlashAmountMustBeNonZero.selector, vault1)
        );
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        slashAmounts[0] = 888;
        slashAmounts[1] = 0;

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.SlashAmountMustBeNonZero.selector, vault2)
        );
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        slashAmounts[1] = 999;

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultSlashAmountUpdated(address(vault1), 888);
        vm.expectEmit(true, true, true, true);
        emit VaultSlashAmountUpdated(address(vault2), 999);
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 888);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 999);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), block.timestamp), 888);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), block.timestamp), 999);

        slashAmounts[0] = 3333;
        slashAmounts[1] = 4444;

        vaults = new address[](2);
        vaults[0] = address(vault2);
        vaults[1] = address(vault1);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultSlashAmountUpdated(address(vault2), 3333);
        vm.expectEmit(true, true, true, true);
        emit VaultSlashAmountUpdated(address(vault1), 4444);
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 4444);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 3333);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), block.timestamp), 4444);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), block.timestamp), 3333);
    }

    function test_requestVaultDeregistrations() public {

        vm.prank(vm.addr(0x1121));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x1121))
        );
        mevCommitMiddleware.requestVaultDeregistrations(new address[](0));

        address[] memory vaults = new address[](1);
        vaults[0] = address(vault1);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault1)
        );
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vaults[0] = address(vault2);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault2)
        );
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vm.warp(888);

        test_registerVaults();

        IMevCommitMiddleware.VaultRecord memory vaultRecord1 = getVaultRecord(address(vault1));
        IMevCommitMiddleware.VaultRecord memory vaultRecord2 = getVaultRecord(address(vault2));
        assertTrue(vaultRecord1.exists);
        assertTrue(vaultRecord2.exists);
        assertFalse(vaultRecord1.deregRequestOccurrence.exists);
        assertFalse(vaultRecord2.deregRequestOccurrence.exists);
        assertEq(vaultRecord1.deregRequestOccurrence.timestamp, 0);
        assertEq(vaultRecord2.deregRequestOccurrence.timestamp, 0);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 15);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 20);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), block.timestamp), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), block.timestamp), 20);

        vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.warp(999);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultDeregistrationRequested(address(vault1));
        vm.expectEmit(true, true, true, true);
        emit VaultDeregistrationRequested(address(vault2));
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vaultRecord1 = getVaultRecord(address(vault1));
        vaultRecord2 = getVaultRecord(address(vault2));
        assertTrue(vaultRecord1.exists);
        assertTrue(vaultRecord2.exists);
        assertTrue(vaultRecord1.deregRequestOccurrence.exists);
        assertTrue(vaultRecord2.deregRequestOccurrence.exists);
        assertEq(vaultRecord1.deregRequestOccurrence.timestamp, 999);
        assertEq(vaultRecord2.deregRequestOccurrence.timestamp, 999);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 15);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 20);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 999), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 999), 20);

        vaults = new address[](1);
        vaults[0] = address(vault2);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultDeregRequestExists.selector, vault2)
        );
        mevCommitMiddleware.requestVaultDeregistrations(vaults);
    }

    function test_deregisterVaults() public {
        vm.prank(vm.addr(0x1121));
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(0x1121))
        );
        mevCommitMiddleware.deregisterVaults(new address[](0));

        address[] memory vaults = new address[](1);
        vaults[0] = address(vault2);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault2)
        );
        mevCommitMiddleware.deregisterVaults(vaults);

        test_requestVaultDeregistrations();

        assertEq(block.timestamp, 999);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotReadyToDeregister.selector, vault2, 999, 999)
        );
        mevCommitMiddleware.deregisterVaults(vaults);

        vm.warp(1001);
        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotReadyToDeregister.selector, vault2, 1001, 999)
        );
        mevCommitMiddleware.deregisterVaults(vaults);

        vm.warp(1001 + mevCommitMiddleware.slashPeriodSeconds() + 1);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultDeregistered(address(vault2));
        mevCommitMiddleware.deregisterVaults(vaults);

        IMevCommitMiddleware.VaultRecord memory vaultRecord1 = getVaultRecord(address(vault1));
        IMevCommitMiddleware.VaultRecord memory vaultRecord2 = getVaultRecord(address(vault2));
        assertTrue(vaultRecord1.exists);
        assertFalse(vaultRecord2.exists);
        
        vaults = new address[](2);
        vaults[1] = address(vault2);
        vaults[0] = address(vault1);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault2)
        );
        mevCommitMiddleware.deregisterVaults(vaults);

        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault2)
        );
        mevCommitMiddleware.deregisterVaults(vaults);

        vaults = new address[](1);
        vaults[0] = address(vault1);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultDeregistered(address(vault1));
        mevCommitMiddleware.deregisterVaults(vaults);

        vaultRecord1 = getVaultRecord(address(vault1));
        vaultRecord2 = getVaultRecord(address(vault2));
        assertFalse(vaultRecord1.exists);
        assertFalse(vaultRecord2.exists);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault1)
        );
        mevCommitMiddleware.deregisterVaults(vaults);
    }

    function test_vaultRegCycle() public {
        test_deregisterVaults();

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        uint160[] memory slashAmounts = new uint160[](2);
        slashAmounts[0] = 88888;
        slashAmounts[1] = 99999;

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit VaultRegistered(address(vault1), 88888);
        vm.expectEmit(true, true, true, true);
        emit VaultRegistered(address(vault2), 99999);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);
    }

    function test_registerValidatorsOperatorReverts() public {

        address operator1 = vm.addr(0x1117);

        bytes[][] memory blsPubkeys = new bytes[][](1);
        blsPubkeys[0] = new bytes[](1);
        blsPubkeys[0][0] = hex"000004444444";

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidBLSPubKeyLength.selector, 48, 6)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        blsPubkeys[0][0] = sampleValPubkey1;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidArrayLengths.selector, 2, 1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        blsPubkeys = new bytes[][](2);
        blsPubkeys[0] = new bytes[](2);
        blsPubkeys[0][0] = sampleValPubkey1;
        blsPubkeys[0][1] = sampleValPubkey2;
        blsPubkeys[1] = new bytes[](1);
        blsPubkeys[1][0] = sampleValPubkey3;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotEntity.selector, operator1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        vm.prank(operator1);
        operatorRegistryMock.register();

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        address[] memory operators = new address[](1);
        operators[0] = operator1;

        vm.prank(owner);
        mevCommitMiddleware.blacklistOperators(operators);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorIsBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        vm.prank(owner);
        mevCommitMiddleware.unblacklistOperators(operators);

        vm.prank(owner);
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorDeregRequestExists.selector, operator1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);
    }

    function getOperatorRecord(address operator) public view
        returns (IMevCommitMiddleware.OperatorRecord memory) {
        (TimestampOccurrence.Occurrence memory occurrence, bool exists, bool isBlacklisted) =
            mevCommitMiddleware.operatorRecords(operator);
        return IMevCommitMiddleware.OperatorRecord(occurrence, exists, isBlacklisted);
    }

    function getVaultRecord(address vault) public view
        returns (IMevCommitMiddleware.VaultRecord memory) {
        (bool exists, TimestampOccurrence.Occurrence memory occurrence, Checkpoints.Trace160 memory slashAmountHistory) =
            mevCommitMiddleware.vaultRecords(vault);
        return IMevCommitMiddleware.VaultRecord(exists, occurrence, slashAmountHistory);
    }

    function getValidatorRecord(bytes memory blsPubkey) public view
        returns (IMevCommitMiddleware.ValidatorRecord memory) {
        (address vault, address operator, bool exists, TimestampOccurrence.Occurrence memory occurrence) =
            mevCommitMiddleware.validatorRecords(blsPubkey);
        return IMevCommitMiddleware.ValidatorRecord(vault, operator, exists, occurrence);
    }

    function getSlashRecord(address vault, address operator, uint256 blockNumber) public view
        returns (IMevCommitMiddleware.SlashRecord memory) {
        (bool exists, uint256 numSlashed, uint256 numInitSlashable) = mevCommitMiddleware.slashRecords(vault, operator, blockNumber);
        return IMevCommitMiddleware.SlashRecord(exists, numSlashed, numInitSlashable);
    }
}
