// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable func-name-mixedcase

import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareTest} from "./MevCommitMiddlewareTest.sol";
import {MockVetoSlasher} from "./MockVetoSlasher.sol";
import {MockInstantSlasher} from "./MockInstantSlasher.sol";

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

        uint256[] memory slashAmounts = new uint256[](2);
        slashAmounts[0] = 10;
        slashAmounts[1] = 20;

        mockDelegator1.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());
        mockDelegator2.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());

        MockInstantSlasher mockSlasher1 = new MockInstantSlasher(mevCommitMiddleware.INSTANT_SLASHER_TYPE());
        MockVetoSlasher mockSlasher2 = new MockVetoSlasher(mevCommitMiddleware.VETO_SLASHER_TYPE(), address(0), 5);

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vault1.setEpochDuration(151);
        vault2.setEpochDuration(151 + 5);

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

        uint256[] memory slashAmounts = new uint256[](2);
        slashAmounts[0] = 10;
        slashAmounts[1] = 20;

        mockDelegator1.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());
        mockDelegator2.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());

        MockInstantSlasher mockSlasher1 = new MockInstantSlasher(mevCommitMiddleware.INSTANT_SLASHER_TYPE());
        MockVetoSlasher mockSlasher2 = new MockVetoSlasher(mevCommitMiddleware.VETO_SLASHER_TYPE(), address(0), 5);

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vault1.setEpochDuration(151);
        vault2.setEpochDuration(151 + 5);

        vm.prank(owner);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        // delegator 1 (associated with vault 1) allocates 29 stake to operator 1
        mockDelegator1.setStake(operator1, 29);

        // delegator 2 (associated with vault 2) allocates 55 stake to operator 1
        mockDelegator2.setStake(operator1, 55);

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

        vm.prank(operator1);
        vm.expectEmit(true, true, true, true);
        emit ValRecordAdded(sampleValPubkey7, operator1, address(vault1), 4);
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);
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

    // For repeated use in requestValidatorDeregistrations tests
    function getSixPubkeys() internal view returns (bytes[] memory) {
        bytes[] memory blsPubkeys = new bytes[](6);
        blsPubkeys[0] = sampleValPubkey1;
        blsPubkeys[1] = sampleValPubkey2;
        blsPubkeys[2] = sampleValPubkey3;
        blsPubkeys[3] = sampleValPubkey4;
        blsPubkeys[4] = sampleValPubkey5;
        blsPubkeys[5] = sampleValPubkey6;
        return blsPubkeys;
    }

    // Test dereg functions are valid from contract owner or fully valid operator

    // TODO: val reg cycle
}
