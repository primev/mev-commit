// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable func-name-mixedcase

import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareTest} from "./MevCommitMiddlewareTest.sol";
import {MockVetoSlasher} from "./MockVetoSlasher.sol";
import {MockInstantSlasher} from "./MockInstantSlasher.sol";
import {MockDelegator} from "./MockDelegator.sol";

contract MevCommitMiddlewareTestCont is MevCommitMiddlewareTest {

    function setUp() public override {
        super.setUp();
    }

    function test_registerValidatorsVaultReverts() public {
        test_registerOperators();
        address operator1 = vm.addr(0x1117);

        bytes[][] memory blsPubkeys = new bytes[][](2);
        blsPubkeys[0] = new bytes[](2);
        blsPubkeys[0][0] = sampleValPubkey1;
        blsPubkeys[0][1] = sampleValPubkey2;
        blsPubkeys[1] = new bytes[](1);
        blsPubkeys[1][0] = sampleValPubkey3;

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotEntity.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        vm.prank(address(vault1));
        vaultFactoryMock.register();

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        uint160[] memory slashAmounts = new uint160[](2);
        slashAmounts[0] = 10;
        slashAmounts[1] = 20;

        uint64 networkRestakeDelegatorType = 0;

        mockDelegator1.setType(networkRestakeDelegatorType);
        mockDelegator2.setType(networkRestakeDelegatorType);

        uint64 instantSlasherType = 0;
        uint64 vetoSlasherType = 1;

        MockInstantSlasher mockSlasher1 = new MockInstantSlasher(instantSlasherType, mockDelegator1);
        uint256 vetoDuration = 5 hours;
        MockVetoSlasher mockSlasher2 = new MockVetoSlasher(vetoSlasherType, address(0), vetoDuration, mockDelegator2, address(mevCommitMiddleware));

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vault1.setEpochDuration(151 hours);
        vault2.setEpochDuration(151 hours + 5 hours);

        vm.prank(address(vault1));
        vaultFactoryMock.register();
        vm.prank(address(vault2));
        vaultFactoryMock.register();

        vm.prank(owner);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vm.prank(owner);
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultDeregRequestExists.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);
    }

    function test_registerValidatorsInvalidArrayLengths() public {
        test_registerOperators();

        bytes[][] memory blsPubkeys = new bytes[][](1);
        blsPubkeys[0] = new bytes[](1);
        blsPubkeys[0][0] = sampleValPubkey1;

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidArrayLengths.selector, 2, 1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);
    }

    function test_registerValidators() public {
        test_registerOperators();
        address operator1 = vm.addr(0x1117);

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.prank(address(vault1));
        vaultFactoryMock.register();
        vm.prank(address(vault2));
        vaultFactoryMock.register();

        uint160[] memory slashAmounts = new uint160[](2);
        slashAmounts[0] = 10;
        slashAmounts[1] = 20;

        uint64 networkRestakeDelegatorType = 0;

        mockDelegator1.setType(networkRestakeDelegatorType);
        mockDelegator2.setType(networkRestakeDelegatorType);

        uint64 instantSlasherType = 0;
        uint64 vetoSlasherType = 1;

        MockInstantSlasher mockSlasher1 = new MockInstantSlasher(instantSlasherType, mockDelegator1);
        uint256 vetoDuration = 5 hours;
        MockVetoSlasher mockSlasher2 = new MockVetoSlasher(vetoSlasherType, address(0), vetoDuration, mockDelegator2, address(mevCommitMiddleware));

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vault1.setEpochDuration(151 hours);
        vault2.setEpochDuration(151 hours + 5 hours);

        vm.prank(owner);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        // delegator 1 (associated with vault 1) allocates 29 stake to operator 1
        mockDelegator1.setStake(operator1, 29);

        uint256 potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1);
        assertEq(potentialSlashableVals, 2);

        // delegator 2 (associated with vault 2) allocates 55 stake to operator 1
        mockDelegator2.setStake(operator1, 55);

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1);
        assertEq(potentialSlashableVals, 2);

        bytes[][] memory blsPubkeys = new bytes[][](2);
        blsPubkeys[0] = new bytes[](3);
        blsPubkeys[0][0] = sampleValPubkey1;
        blsPubkeys[0][1] = sampleValPubkey2;
        blsPubkeys[0][2] = sampleValPubkey3;
        blsPubkeys[1] = new bytes[](3);
        blsPubkeys[1][0] = sampleValPubkey4;
        blsPubkeys[1][1] = sampleValPubkey5;
        blsPubkeys[1][2] = sampleValPubkey6;

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault1), operator1, 3, 2)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 10);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault1), operator1, 3, 1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 30);

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1);
        assertEq(potentialSlashableVals, 3);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault2), operator1, 3, 2)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator2.setStake(operator1, 19);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault2), operator1, 3, 0)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator2.setStake(operator1, 39);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault2), operator1, 3, 1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator2.setStake(operator1, 99);

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1);
        assertEq(potentialSlashableVals, 4);

        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            for (uint256 j = 0; j < blsPubkeys[i].length; j++) {
                IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(blsPubkeys[i][j]);
                assertEq(valRecord.vault, address(0));
                assertEq(valRecord.operator, address(0));
                assertFalse(valRecord.exists);
                assertFalse(valRecord.deregRequestOccurrence.exists);
            }
        }

        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey1, operator1, address(vault1), 1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey2, operator1, address(vault1), 2);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey3, operator1, address(vault1), 3);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey4, operator1, address(vault2), 1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey5, operator1, address(vault2), 2);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey6, operator1, address(vault2), 3);
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        IMevCommitMiddleware.ValidatorRecord memory valRecord1 = getValidatorRecord(sampleValPubkey1);
        assertEq(valRecord1.vault, address(vault1));
        assertEq(valRecord1.operator, operator1);
        assertTrue(valRecord1.exists);
        assertFalse(valRecord1.deregRequestOccurrence.exists);

        IMevCommitMiddleware.ValidatorRecord memory valRecord2 = getValidatorRecord(sampleValPubkey2);
        assertEq(valRecord2.vault, address(vault1));
        assertEq(valRecord2.operator, operator1);
        assertTrue(valRecord2.exists);
        assertFalse(valRecord2.deregRequestOccurrence.exists);

        IMevCommitMiddleware.ValidatorRecord memory valRecord3 = getValidatorRecord(sampleValPubkey3);
        assertEq(valRecord3.vault, address(vault1));
        assertEq(valRecord3.operator, operator1);
        assertTrue(valRecord3.exists);
        assertFalse(valRecord3.deregRequestOccurrence.exists);

        IMevCommitMiddleware.ValidatorRecord memory valRecord4 = getValidatorRecord(sampleValPubkey4);
        assertEq(valRecord4.vault, address(vault2));
        assertEq(valRecord4.operator, operator1);
        assertTrue(valRecord4.exists);
        assertFalse(valRecord4.deregRequestOccurrence.exists);

        IMevCommitMiddleware.ValidatorRecord memory valRecord5 = getValidatorRecord(sampleValPubkey5);
        assertEq(valRecord5.vault, address(vault2));
        assertEq(valRecord5.operator, operator1);
        assertTrue(valRecord5.exists);
        assertFalse(valRecord5.deregRequestOccurrence.exists);

        IMevCommitMiddleware.ValidatorRecord memory valRecord6 = getValidatorRecord(sampleValPubkey6);
        assertEq(valRecord6.vault, address(vault2));
        assertEq(valRecord6.operator, operator1);
        assertTrue(valRecord6.exists);
        assertFalse(valRecord6.deregRequestOccurrence.exists);

        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1), 1);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1), 2);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey3, address(vault1), operator1), 3);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1), 1);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1), 2);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1), 3);

        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey1));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey2));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey3));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey4));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey5));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey6));

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1);
        assertEq(potentialSlashableVals, 0);

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1);
        assertEq(potentialSlashableVals, 1);

        blsPubkeys = new bytes[][](1);
        blsPubkeys[0] = new bytes[](1);
        blsPubkeys[0][0] = sampleValPubkey3;

        vaults = new address[](1);
        vaults[0] = address(vault2);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorRecordAlreadyExists.selector, sampleValPubkey3)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        blsPubkeys[0][0] = sampleValPubkey7;
        vaults[0] = address(vault1);
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault1), operator1, 1, 0)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 10);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault1), operator1, 1, 0)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 39);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorsNotSlashable.selector,
                address(vault1), operator1, 1, 0)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 40);

        potentialSlashableVals = mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1);
        assertEq(potentialSlashableVals, 1);

        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey7, operator1, address(vault1), 4);
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        mockDelegator1.setStake(operator1, 11);

        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey1));
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey2));
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey3));
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey7));

        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey4));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey5));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey6));

        mockDelegator1.setStake(operator1, 51);

        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey1));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey2));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey3));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey7));

        mockDelegator1.setStake(operator1, 40);
    }

    function test_requestValidatorDeregistrationsMissingValidatorRecord() public { 
        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingValidatorRecord.selector,
                sampleValPubkey1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsOnlyOperator() public {
        test_registerValidators();
        address operator1 = vm.addr(0x1117);
        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.prank(vm.addr(0x9999999));
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OnlyOperator.selector, operator1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsOperatorNotEntity() public {
        test_registerValidators();
        address operator1 = vm.addr(0x1117);

        vm.prank(operator1);
        operatorRegistryMock.deregister();

        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotEntity.selector, operator1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsOperatorNotRegistered() public {
        test_registerValidators();

        address operator1 = vm.addr(0x1117);
        address[] memory operators = new address[](1);
        operators[0] = operator1;
        vm.prank(owner);
        mevCommitMiddleware.requestOperatorDeregistrations(operators);
        IMevCommitMiddleware.OperatorRecord memory operatorRecord = getOperatorRecord(operator1);
        assertTrue(operatorRecord.deregRequestOccurrence.exists);

        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorDeregRequestExists.selector, operator1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);

        vm.warp(block.timestamp + mevCommitMiddleware.slashPeriodSeconds() + 1);

        vm.prank(owner);
        mevCommitMiddleware.deregisterOperators(operators);
        operatorRecord = getOperatorRecord(operator1);
        assertFalse(operatorRecord.exists);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorNotRegistered.selector, operator1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsOperatorIsBlacklisted() public {
        test_registerValidators();

        address operator1 = vm.addr(0x1117);
        address[] memory operators = new address[](1);
        operators[0] = operator1;

        vm.prank(owner);
        mevCommitMiddleware.blacklistOperators(operators);

        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorIsBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsFromValidOperator() public {
        test_registerValidators();
        bytes[] memory blsPubkeys = getSixPubkeys();

        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(blsPubkeys[i]);
            assertTrue(valRecord.exists);
            assertFalse(valRecord.deregRequestOccurrence.exists);
            assertEq(valRecord.deregRequestOccurrence.timestamp, 0);
        }

        vm.warp(91);

        address operator1 = vm.addr(0x1117);
        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey1, operator1, 1);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey2, operator1, 2);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey3, operator1, 3);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey4, operator1, 1);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey5, operator1, 2);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey6, operator1, 3);
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);

        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(blsPubkeys[i]);
            assertTrue(valRecord.exists);
            assertTrue(valRecord.deregRequestOccurrence.exists);
            assertEq(valRecord.deregRequestOccurrence.timestamp, 91);
        }

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorDeregRequestExists.selector, sampleValPubkey1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_requestValidatorDeregistrationsFromContractOwner() public {
        test_registerValidators();
        bytes[] memory blsPubkeys = getSixPubkeys();

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey1, owner, 1);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey2, owner, 2);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey3, owner, 3);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey4, owner, 1);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey5, owner, 2);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(sampleValPubkey6, owner, 3);
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);

        vm.prank(owner);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorDeregRequestExists.selector, sampleValPubkey1)
        );
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);
    }

    function test_getPositionInValsetInvalidParameters() public {
        address vault = address(vault1);
        address operator = vm.addr(0x1117);
        bytes memory badKey = bytes("0x1234");
        uint256 position = mevCommitMiddleware.getPositionInValset(badKey, vault, operator);
        assertEq(position, 0);

        vault = address(0x12347);
        position = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, vault, operator);
        assertEq(position, 0);

        vault = address(vault1);
        operator = address(0x1234567890123456789012345678901234567890);
        position = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, vault, operator);
        assertEq(position, 0);

        test_registerValidators();

        vault = address(vault1);
        operator = vm.addr(0x1117);
        badKey = bytes("0x1234");
        position = mevCommitMiddleware.getPositionInValset(badKey, vault, operator);
        assertEq(position, 0);

        vault = address(0x12347);
        position = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, vault, operator);
        assertEq(position, 0);

        vault = address(vault1);
        operator = address(0x1234567890123456789012345678901234567890);
        position = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, vault, operator);
        assertEq(position, 0);
    }

    function test_deregisterValidatorsMissingValidatorRecord() public {
        bytes[] memory blsPubkeys = getSixPubkeys();
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingValidatorRecord.selector, sampleValPubkey1)
        );
        mevCommitMiddleware.deregisterValidators(blsPubkeys);
    }

    function test_deregisterValidatorsFromContractOwner() public {
        test_requestValidatorDeregistrationsFromValidOperator();
        assertEq(getValidatorRecord(sampleValPubkey1).deregRequestOccurrence.timestamp, 91);

        address operator1 = vm.addr(0x1117);
        bytes[] memory blsPubkeys = getSixPubkeys();

        vm.warp(91 + 1);
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorNotReadyToDeregister.selector, sampleValPubkey1, 92, 91)
        );
        mevCommitMiddleware.deregisterValidators(blsPubkeys);

        vm.warp(91 + 20);
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.ValidatorNotReadyToDeregister.selector, sampleValPubkey1, 111, 91)
        );
        mevCommitMiddleware.deregisterValidators(blsPubkeys);

        uint256 pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1);
        assertEq(pos1, 1);  
        uint256 pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1);
        assertEq(pos2, 2);
        uint256 pos3 = mevCommitMiddleware.getPositionInValset(sampleValPubkey3, address(vault1), operator1);
        assertEq(pos3, 3);
        uint256 pos4 = mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1);
        assertEq(pos4, 1);
        uint256 pos5 = mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1);
        assertEq(pos5, 2);
        uint256 pos6 = mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1);
        assertEq(pos6, 3);

        vm.warp(91+mevCommitMiddleware.slashPeriodSeconds()+1);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey1, owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey2, owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey3, owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey4, owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey5, owner);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey6, owner);
        mevCommitMiddleware.deregisterValidators(blsPubkeys);

        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(blsPubkeys[i]);
            assertFalse(valRecord.exists);
            uint256 pos = mevCommitMiddleware.getPositionInValset(blsPubkeys[i], address(vault1), operator1);
            assertEq(pos, 0);
        }
    }

    function test_deregisterValidatorsFromValidOperator() public {
        test_requestValidatorDeregistrationsFromValidOperator();
        assertEq(getValidatorRecord(sampleValPubkey1).deregRequestOccurrence.timestamp, 91);

        address operator1 = vm.addr(0x1117);
        bytes[] memory blsPubkeys = getSixPubkeys();

        vm.warp(91+mevCommitMiddleware.slashPeriodSeconds()+1);

        uint256 pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1);
        assertEq(pos1, 1);  
        uint256 pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1);
        assertEq(pos2, 2);
        uint256 pos3 = mevCommitMiddleware.getPositionInValset(sampleValPubkey3, address(vault1), operator1);
        assertEq(pos3, 3);
        uint256 pos4 = mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1);
        assertEq(pos4, 1);
        uint256 pos5 = mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1);
        assertEq(pos5, 2);
        uint256 pos6 = mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1);
        assertEq(pos6, 3);

        uint256 length = mevCommitMiddleware.valsetLength(address(vault1), operator1);
        assertEq(length, 4); // Recall sampleValPubkey7 was added to vault1 and never deregistered.
        length = mevCommitMiddleware.valsetLength(address(vault2), operator1);
        assertEq(length, 3);

        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey1, operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey2, operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey3, operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey4, operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey5, operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordDeleted(sampleValPubkey6, operator1);
        mevCommitMiddleware.deregisterValidators(blsPubkeys);

        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(blsPubkeys[i]);
            assertFalse(valRecord.exists);
            uint256 pos = mevCommitMiddleware.getPositionInValset(blsPubkeys[i], address(vault1), operator1);
            assertEq(pos, 0);
        }

        length = mevCommitMiddleware.valsetLength(address(vault1), operator1);
        assertEq(length, 1); // Recall sampleValPubkey7 was added to vault1 and never deregistered.
        length = mevCommitMiddleware.valsetLength(address(vault2), operator1);
        assertEq(length, 0); 
    }

    function test_deregisterValidatorsInvalidOperator() public {
        test_requestValidatorDeregistrationsFromContractOwner();

        address operator1 = vm.addr(0x1117);

        bytes memory badKey = bytes("0x1234");
        bytes[] memory badKeys = new bytes[](1);
        badKeys[0] = badKey;

        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingValidatorRecord.selector, badKey)
        );
        mevCommitMiddleware.deregisterValidators(badKeys);

        bytes[] memory blsPubkeys = getSixPubkeys();

        vm.warp(91+mevCommitMiddleware.slashPeriodSeconds()+1);

        vm.prank(vm.addr(0x99999998888));
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OnlyOperator.selector, operator1)
        );
        mevCommitMiddleware.deregisterValidators(blsPubkeys);
    }

    function test_valRegCycle() public {
        test_deregisterValidatorsFromValidOperator();

        address operator1 = vm.addr(0x1117);

        address[] memory vaults = new address[](1);
        vaults[0] = address(vault1);

        bytes[][] memory blsPubkeys = new bytes[][](1);
        blsPubkeys[0] = new bytes[](1);
        blsPubkeys[0][0] = sampleValPubkey3;

        uint256 length = mevCommitMiddleware.valsetLength(address(vault1), operator1);
        assertEq(length, 1);

        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey3, operator1, address(vault1), 2);
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        uint256 pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1);
        assertEq(pos1, 0);
        uint256 pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1);
        assertEq(pos2, 0);
        uint256 pos3 = mevCommitMiddleware.getPositionInValset(sampleValPubkey3, address(vault1), operator1);
        assertEq(pos3, 2); // Recall sampleValPubkey7 still exists in the vault1 valset.
        uint256 pos4 = mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1);
        assertEq(pos4, 0);
        uint256 pos5 = mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1);
        assertEq(pos5, 0);
        uint256 pos6 = mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1);
        assertEq(pos6, 0);

        length = mevCommitMiddleware.valsetLength(address(vault1), operator1);
        assertEq(length, 2);

        bytes memory pubkey = mevCommitMiddleware.pubkeyAtPositionInValset(0, address(vault1), operator1);
        assertEq(pubkey, bytes(""));

        pubkey = mevCommitMiddleware.pubkeyAtPositionInValset(1, address(vault1), operator1);
        assertEq(pubkey, sampleValPubkey7);

        pubkey = mevCommitMiddleware.pubkeyAtPositionInValset(2, address(vault1), operator1);
        assertEq(pubkey, sampleValPubkey3);

        pubkey = mevCommitMiddleware.pubkeyAtPositionInValset(3, address(vault1), operator1);
        assertEq(pubkey, bytes(""));
    }

    function test_slashValidatorsNotOracle() public {
        bytes[] memory blsPubkeys = getSixPubkeys();
        uint256[] memory timestamps = new uint256[](6);
        timestamps[0] = 100;
        timestamps[1] = 101;
        timestamps[2] = 102;
        timestamps[3] = 103;
        timestamps[4] = 104;
        timestamps[5] = 105;
        vm.prank(vm.addr(0x99999998888888));
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OnlySlashOracle.selector, slashOracle)
        );
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps);
    }

    function test_slashValidatorsInvalidArrayLengths() public {
        bytes[] memory blsPubkeys = getSixPubkeys();
        uint256[] memory timestamps = new uint256[](1);
        timestamps[0] = 100;
        vm.prank(slashOracle);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.InvalidArrayLengths.selector, 6, 1)
        );
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps);
    }

    function test_slashValidatorsValidatorDeregistered() public {
        test_registerValidators();
        bytes[] memory blsPubkeys = getSixPubkeys();

        address operator1 = vm.addr(0x1117);

        vm.prank(owner);
        mevCommitMiddleware.requestValDeregistrations(blsPubkeys);

        vm.warp(106);

        uint256[] memory timestamps = new uint256[](6);
        timestamps[0] = 100;
        timestamps[1] = 101;
        timestamps[2] = 102;
        timestamps[3] = 103;
        timestamps[4] = 104;
        timestamps[5] = 105;
        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey1, operator1, address(vault1), 10);
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps); // slash successful with req dereg

        vm.warp(block.timestamp + slashPeriodSeconds + 1);
        vm.prank(owner);
        mevCommitMiddleware.deregisterValidators(blsPubkeys);

        vm.prank(slashOracle);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingValidatorRecord.selector, sampleValPubkey1)
        );
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps);
    }

    function test_slashValidatorsVaultDeregistered() public { 
        test_registerValidators();
        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        address operator1 = vm.addr(0x1117);

        vm.prank(owner);
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vm.warp(106);

        bytes[] memory blsPubkeys = getSixPubkeys();
        uint256[] memory timestamps = new uint256[](6);
        timestamps[0] = 100;
        timestamps[1] = 101;
        timestamps[2] = 102;
        timestamps[3] = 103;
        timestamps[4] = 104;
        timestamps[5] = 105;
        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey1, operator1, address(vault1), 10);
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps); // slash successful with req dereg

        vm.warp(block.timestamp + slashPeriodSeconds + 1);
        vm.prank(owner);
        mevCommitMiddleware.deregisterVaults(vaults);

        vm.prank(slashOracle);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingVaultRecord.selector, address(vault1))
        );
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps);
    }

    function test_slashValidatorsOperatorDeregistered() public { 
        test_registerValidators();
        address[] memory operators = new address[](2);
        operators[0] = vm.addr(0x1117);
        operators[1] = vm.addr(0x1118);

        vm.prank(owner);
        mevCommitMiddleware.requestOperatorDeregistrations(operators);

        vm.warp(106);

        bytes[] memory blsPubkeys = getSixPubkeys();
        uint256[] memory timestamps = new uint256[](6);
        timestamps[0] = 100;
        timestamps[1] = 101;
        timestamps[2] = 102;
        timestamps[3] = 103;
        timestamps[4] = 104;
        timestamps[5] = 105;
        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey1, operators[0], address(vault1), 10);
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps); // slash successful with req dereg

        vm.warp(block.timestamp + slashPeriodSeconds + 1);
        vm.prank(owner);
        mevCommitMiddleware.deregisterOperators(operators);

        vm.prank(slashOracle);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.MissingOperatorRecord.selector, operators[0])
        );
        mevCommitMiddleware.slashValidators(blsPubkeys, timestamps);
    }

    function test_slashValidatorsSuccess() public { 
        test_registerValidators();
        bytes[] memory firstTwoBlsPubkeys = new bytes[](2);
        firstTwoBlsPubkeys[0] = sampleValPubkey1;
        firstTwoBlsPubkeys[1] = sampleValPubkey2;

        address operator1 = vm.addr(0x1117);

        uint256[] memory timestamps = new uint256[](2);
        timestamps[0] = 0;
        timestamps[1] = 101;

        vm.prank(slashOracle);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.CaptureTimestampMustBeNonZero.selector)
        );
        mevCommitMiddleware.slashValidators(firstTwoBlsPubkeys, timestamps);

        timestamps[0] = 100;

        // Oracle (acting correctly) will not slash for capture timestamps in the future.
        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.FutureTimestampDisallowed.selector, address(vault1), 100));
        vm.prank(slashOracle);
        mevCommitMiddleware.slashValidators(firstTwoBlsPubkeys, timestamps);

        vm.warp(200);

        IMevCommitMiddleware.ValidatorRecord memory valRecord1 = getValidatorRecord(sampleValPubkey1);
        IMevCommitMiddleware.ValidatorRecord memory valRecord2 = getValidatorRecord(sampleValPubkey2);
        assertTrue(valRecord1.exists);
        assertTrue(valRecord2.exists);
        assertFalse(valRecord1.deregRequestOccurrence.exists);
        assertFalse(valRecord2.deregRequestOccurrence.exists);

        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey1));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey2));

        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey1));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey2));

        assertEq(mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1), 0);

        uint256 pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1);
        assertEq(pos1, 1);
        uint256 pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1);
        assertEq(pos2, 2);

        assertEq(mevCommitMiddleware.valsetLength(address(vault1), operator1), 4);

        IMevCommitMiddleware.SlashRecord memory slashRecord = getSlashRecord(address(vault1), operator1, block.number);
        assertFalse(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 0);
        assertEq(slashRecord.numSlashed, 0);

        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey1, operator1, address(vault1), 10);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey2, operator1, address(vault1), 10);
        mevCommitMiddleware.slashValidators(firstTwoBlsPubkeys, timestamps); 

        valRecord1 = getValidatorRecord(sampleValPubkey1);
        valRecord2 = getValidatorRecord(sampleValPubkey2);
        assertTrue(valRecord1.exists);
        assertTrue(valRecord2.exists);
        assertTrue(valRecord1.deregRequestOccurrence.exists);
        assertTrue(valRecord2.deregRequestOccurrence.exists);

        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey1));
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey2));

        // Validators should no longer be slashable, since instant slasher decrements stake immediately.
        assertFalse(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey1));
        assertFalse(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey2));

        assertEq(mevCommitMiddleware.potentialSlashableValidators(address(vault1), operator1), 0);

        pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey1, address(vault1), operator1);
        assertEq(pos1, 4); // final position of first set
        pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey2, address(vault1), operator1);
        assertEq(pos2, 3); // second to final position of first set

        assertEq(mevCommitMiddleware.valsetLength(address(vault1), operator1), 4);

        slashRecord = getSlashRecord(address(vault1), operator1, block.number);
        assertTrue(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 4);
        assertEq(slashRecord.numSlashed, 2);
    }

    event ExecuteSlash(uint256 indexed slashIndex, uint256 slashedAmount);

    function test_slashValidatorsWithVetoSlasher() public { 
        test_slashValidatorsSuccess();
        address operator1 = vm.addr(0x1117);

        vm.roll(block.number + 20);
        vm.warp(300);

        bytes[] memory firstTwoBlsPubkeysFromVault2 = new bytes[](2);
        firstTwoBlsPubkeysFromVault2[0] = sampleValPubkey4;
        firstTwoBlsPubkeysFromVault2[1] = sampleValPubkey5;

        uint256[] memory timestamps = new uint256[](2);
        timestamps[0] = 201;
        timestamps[1] = 202;

        IMevCommitMiddleware.ValidatorRecord memory valRecord4 = getValidatorRecord(sampleValPubkey4);
        IMevCommitMiddleware.ValidatorRecord memory valRecord5 = getValidatorRecord(sampleValPubkey5);
        assertTrue(valRecord4.exists);
        assertTrue(valRecord5.exists);
        assertFalse(valRecord4.deregRequestOccurrence.exists);
        assertFalse(valRecord5.deregRequestOccurrence.exists);

        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey4));
        assertTrue(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey5));

        // get stake amount before slashing
        MockDelegator delegator = MockDelegator(vault2.delegator());
        uint256 allocatedStake = delegator.stake(bytes32("subnet"), operator1);
        assertEq(allocatedStake, 99);

        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey4));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey5));

        assertEq(mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1), 1);

        uint256 pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1);
        assertEq(pos1, 1);
        uint256 pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1);
        assertEq(pos2, 2);

        assertEq(mevCommitMiddleware.valsetLength(address(vault2), operator1), 3);

        IMevCommitMiddleware.SlashRecord memory slashRecord = getSlashRecord(address(vault2), operator1, block.number);
        assertFalse(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 0);
        assertEq(slashRecord.numSlashed, 0);

        assertEq(mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1), 1); 

        uint256 slashAmount = mevCommitMiddleware.getLatestSlashAmount(address(vault2));

        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey4, operator1, address(vault2), slashAmount);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey5, operator1, address(vault2), slashAmount);
        vm.expectEmit(true, true, true, true);
        address[] memory expectedVaults = new address[](2);
        expectedVaults[0] = address(vault2);
        expectedVaults[1] = address(vault2);
        address[] memory expectedOperators = new address[](2);
        expectedOperators[0] = operator1;
        expectedOperators[1] = operator1;
        uint256[] memory expectedPositions = new uint256[](2);
        expectedPositions[0] = 3;
        expectedPositions[1] = 2;
        emit ValidatorPositionsSwapped(firstTwoBlsPubkeysFromVault2, expectedVaults, expectedOperators, expectedPositions);
        mevCommitMiddleware.slashValidators(firstTwoBlsPubkeysFromVault2, timestamps); 

        valRecord4 = getValidatorRecord(sampleValPubkey4);
        valRecord5 = getValidatorRecord(sampleValPubkey5);
        assertTrue(valRecord4.exists);
        assertTrue(valRecord5.exists);
        assertTrue(valRecord4.deregRequestOccurrence.exists);
        assertTrue(valRecord5.deregRequestOccurrence.exists);

        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey4));
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey5));

        // total stake should be 59
        allocatedStake = delegator.stake(bytes32("subnet"), operator1);
        assertEq(allocatedStake, 59);

        // validator 4 should not be slashable, validator 5 should be since it's lower index
        assertFalse(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey4));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey5));

        assertEq(mevCommitMiddleware.potentialSlashableValidators(address(vault2), operator1), 0); 

        pos1 = mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1);
        assertEq(pos1, 3); // final position of second set
        pos2 = mevCommitMiddleware.getPositionInValset(sampleValPubkey5, address(vault2), operator1);
        assertEq(pos2, 2); // second to final position of second set

        assertEq(mevCommitMiddleware.valsetLength(address(vault2), operator1), 3);

        slashRecord = getSlashRecord(address(vault1), operator1, block.number);
        assertFalse(slashRecord.exists);

        slashRecord = getSlashRecord(address(vault1), operator1, block.number-20);
        assertTrue(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 4);
        assertEq(slashRecord.numSlashed, 2);

        slashRecord = getSlashRecord(address(vault2), operator1, block.number);
        assertTrue(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 3);
        assertEq(slashRecord.numSlashed, 2);

        assertEq(mevCommitMiddleware.getNumSlashableVals(address(vault2), operator1), 2); // TODO: revisit
        assertEq(mevCommitMiddleware.valsetLength(address(vault2), operator1), 3);

        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1), 1);

        // First two validators are slashable, last one is not
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey6));
        assertTrue(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey5));
        assertFalse(mevCommitMiddleware.isValidatorSlashable(sampleValPubkey4));

        // now slash validator 6 at a later block
        vm.roll(block.number + 20);

        bytes[] memory pubkeys = new bytes[](1);
        pubkeys[0] = sampleValPubkey6;
        timestamps = new uint256[](1);
        timestamps[0] = 203;
        vm.prank(slashOracle);
        vm.expectEmit(true, true, true, true);
        emit ValidatorSlashed(sampleValPubkey6, operator1, address(vault2), slashAmount);
        vm.expectEmit(true, true, true, true);
        expectedVaults = new address[](1);
        expectedVaults[0] = address(vault2);
        expectedOperators = new address[](1);
        expectedOperators[0] = operator1;
        expectedPositions = new uint256[](1);
        expectedPositions[0] = 3;
        emit ValidatorPositionsSwapped(pubkeys, expectedVaults, expectedOperators, expectedPositions);
        mevCommitMiddleware.slashValidators(pubkeys, timestamps);

        slashRecord = getSlashRecord(address(vault2), operator1, block.number);
        assertTrue(slashRecord.exists);
        assertEq(slashRecord.numRegistered, 3);
        assertEq(slashRecord.numSlashed, 1);

        assertFalse(mevCommitMiddleware.isValidatorOptedIn(sampleValPubkey6));

        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey4, address(vault2), operator1), 1);
        assertEq(mevCommitMiddleware.getPositionInValset(sampleValPubkey6, address(vault2), operator1), 3);
    }

    function test_operatorGreifingScenario() public { 
        test_registerValidators();

        address operator1 = vm.addr(0x1117);

        vm.prank(owner);
        address[] memory operators = new address[](1);
        operators[0] = operator1;
        mevCommitMiddleware.blacklistOperators(operators);

        // confirm operator cant opt in any more validators
        bytes[][] memory pubkeys = new bytes[][](1);
        pubkeys[0] = new bytes[](1);
        pubkeys[0][0] = sampleValPubkey8;
        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.OperatorIsBlacklisted.selector, operator1)
        );
        mevCommitMiddleware.registerValidators(pubkeys, operators);

        bytes[] memory allPubkeys = getSixPubkeys();
        bytes[] memory allPubkeysWith7 = new bytes[](7);
        for (uint256 i = 0; i < allPubkeys.length; i++) {
            allPubkeysWith7[i] = allPubkeys[i];
        }
        allPubkeysWith7[6] = sampleValPubkey7;

        vm.prank(owner);
        mevCommitMiddleware.requestValDeregistrations(allPubkeys);

        vm.warp(block.timestamp + slashPeriodSeconds + 1);

        vm.prank(owner);
        mevCommitMiddleware.deregisterValidators(allPubkeys);

        for (uint256 i = 0; i < allPubkeys.length; i++) {
            IMevCommitMiddleware.ValidatorRecord memory valRecord = getValidatorRecord(allPubkeys[i]);
            assertFalse(valRecord.exists);
        }
    }

    function test_updateSlashAmounts() public {

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, address(vault1)));
        mevCommitMiddleware.getSlashAmountAt(address(vault1), 0);

        assertEq(block.timestamp, 1);
        test_registerVaults();
        assertEq(block.timestamp, 1);

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.NoSlashAmountAtTimestamp.selector, address(vault2), 0));
        mevCommitMiddleware.getSlashAmountAt(address(vault2), 0);

        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 15);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 20);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 1), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 1), 20);

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.FutureTimestampDisallowed.selector, address(vault1), 500));
        mevCommitMiddleware.getSlashAmountAt(address(vault1), 500);

        vm.warp(333);
        assertEq(block.timestamp, 333);

        vm.prank(owner);
        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);
        uint160 newSlashAmount1 = 30;
        uint160 newSlashAmount2 = 40;
        uint160[] memory slashAmounts = new uint160[](2);
        slashAmounts[0] = newSlashAmount1;
        slashAmounts[1] = newSlashAmount2;
        mevCommitMiddleware.updateSlashAmounts(vaults, slashAmounts);

        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 30);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 40);

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.NoSlashAmountAtTimestamp.selector, address(vault2), 0));
        mevCommitMiddleware.getSlashAmountAt(address(vault2), 0);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 1), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 1), 20);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 332), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 332), 20);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 333), 30);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 333), 40);

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.FutureTimestampDisallowed.selector, address(vault1), 500));
        mevCommitMiddleware.getSlashAmountAt(address(vault1), 500);

        vm.warp(500);
        assertEq(block.timestamp, 500);

        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault1)), 30);
        assertEq(mevCommitMiddleware.getLatestSlashAmount(address(vault2)), 40);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 500), 30);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 500), 40);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 334), 30);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 334), 40);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 332), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 332), 20);

        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault1), 1), 15);
        assertEq(mevCommitMiddleware.getSlashAmountAt(address(vault2), 1), 20);

        vm.expectRevert(abi.encodeWithSelector(IMevCommitMiddleware.NoSlashAmountAtTimestamp.selector, address(vault1), 0));
        mevCommitMiddleware.getSlashAmountAt(address(vault1), 0);
    }

    function test_isValidatorOptedInBadKey() public view {
        bytes memory badKey = bytes("0x1234");
        assertFalse(mevCommitMiddleware.isValidatorOptedIn(badKey));
    }

    // For repeated use in requestValidatorDeregistrations and deregisterValidators tests
    function getSixPubkeys() public view returns (bytes[] memory) {
        bytes[] memory blsPubkeys = new bytes[](6);
        blsPubkeys[0] = sampleValPubkey1;
        blsPubkeys[1] = sampleValPubkey2;
        blsPubkeys[2] = sampleValPubkey3;
        blsPubkeys[3] = sampleValPubkey4;
        blsPubkeys[4] = sampleValPubkey5;
        blsPubkeys[5] = sampleValPubkey6;
        return blsPubkeys;
    }
}
