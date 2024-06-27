// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {StrategyManagerMock} from "eigenlayer-contracts/src/test/mocks/StrategyManagerMock.sol";
import {DelegationManagerMock} from "eigenlayer-contracts/src/test/mocks/DelegationManagerMock.sol";
import {EigenPodManagerMock} from "./EigenPodManagerMock.sol";
import {EigenPodMock} from "./EigenPodMock.sol";
import {IEigenPod} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {AVSDirectoryMock} from "./AVSDirectoryMock.sol";

contract MevCommitAVSTest is Test {
    MevCommitAVS public mevCommitAVS;

    address public owner;
    StrategyManagerMock public strategyManagerMock;
    DelegationManagerMock public delegationManagerMock;
    EigenPodManagerMock public eigenPodManagerMock;
    AVSDirectoryMock public avsDirectoryMock;
    address[] public restakeableStrategies;
    address public freezeOracle;
    uint256 public unfreezeFee;
    address public unfreezeReceiver;
    uint256 public unfreezePeriodBlocks;
    uint256 public operatorDeregPeriodBlocks;
    uint256 public validatorDeregPeriodBlocks;
    uint256 public lstRestakerDeregPeriodBlocks;
    string public metadataUrl;

    event OperatorRegistered(address indexed operator);
    event OperatorDeregistrationRequested(address indexed operator);
    event OperatorDeregistered(address indexed operator);
    event ValidatorRegistered(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed podOwner);
    event LSTRestakerRegistered(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);
    event LSTRestakerDeregistrationRequested(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);
    event LSTRestakerDeregistered(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);
    event ValidatorFrozen(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorUnfrozen(bytes indexed validatorPubKey, address indexed podOwner);
    event AVSDirectorySet(address indexed avsDirectory);
    event StrategyManagerSet(address indexed strategyManager);
    event DelegationManagerSet(address indexed delegationManager);
    event EigenPodManagerSet(address indexed eigenPodManager);
    event RestakeableStrategiesSet(address[] indexed restakeableStrategies);
    event FreezeOracleSet(address indexed freezeOracle);
    event UnfreezeFeeSet(uint256 unfreezeFee);
    event UnfreezeReceiverSet(address indexed unfreezeReceiver);
    event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks);
    event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks);
    event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks);
    event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks);

    function setUp() public {
        owner = address(0x123456);
        strategyManagerMock = new StrategyManagerMock();
        delegationManagerMock = new DelegationManagerMock();
        eigenPodManagerMock = new EigenPodManagerMock();
        avsDirectoryMock = new AVSDirectoryMock();
        restakeableStrategies = [address(0x1), address(0x2), address(0x3)];
        freezeOracle = address(0x5);
        unfreezeFee = 1 ether;
        unfreezeReceiver = address(0x6);
        unfreezePeriodBlocks = 100;
        operatorDeregPeriodBlocks = 200;
        validatorDeregPeriodBlocks = 300;
        lstRestakerDeregPeriodBlocks = 400;
        metadataUrl = "https://raw.githubusercontent.com/primev/mev-commit/main/static/avs-metadata.json";

        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitAVS.sol",
            abi.encodeCall(MevCommitAVS.initialize, (
                owner,
                delegationManagerMock, 
                eigenPodManagerMock, 
                strategyManagerMock, 
                avsDirectoryMock, 
                restakeableStrategies,
                freezeOracle,
                unfreezeFee,
                unfreezeReceiver,
                unfreezePeriodBlocks,
                operatorDeregPeriodBlocks,
                validatorDeregPeriodBlocks,
                lstRestakerDeregPeriodBlocks,
                metadataUrl
            ))
        );
        mevCommitAVS = MevCommitAVS(payable(proxy));
    }

    function testRegisterOperator() public {
        address operator = address(0x888);
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature = ISignatureUtils.SignatureWithSaltAndExpiry({
            signature: bytes("signature"),
            salt: bytes32("salt"),
            expiry: block.timestamp + 1 days
        });

        vm.expectRevert("sender must be an eigenlayer operator");
        vm.prank(operator);
        mevCommitAVS.registerOperator(operatorSignature);

        delegationManagerMock.setIsOperator(operator, true);

        vm.expectEmit(true, true, true, true);
        emit OperatorRegistered(operator);
        vm.prank(operator);
        mevCommitAVS.registerOperator(operatorSignature);

        IMevCommitAVS.OperatorRegistrationInfo memory reg = mevCommitAVS.getOperatorRegInfo(operator);
        assertTrue(reg.exists);
        assertFalse(reg.deregRequestHeight.exists);

        vm.expectRevert("sender must not be registered operator");
        vm.prank(operator);
        mevCommitAVS.registerOperator(operatorSignature);
    }

    function testRequestOperatorDeregistration() public {
        vm.roll(108);

        address operator = address(0x888);
        vm.expectRevert("operator must be registered");
        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);

        testRegisterOperator();

        address otherAcct = address(0x777);
        vm.expectRevert("sender must be operator");
        vm.prank(otherAcct);
        mevCommitAVS.requestOperatorDeregistration(operator);

        IMevCommitAVS.OperatorRegistrationInfo memory reg = mevCommitAVS.getOperatorRegInfo(operator);
        assertTrue(reg.exists);
        assertFalse(reg.deregRequestHeight.exists);

        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistrationRequested(operator);
        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);

        reg = mevCommitAVS.getOperatorRegInfo(operator);
        assertTrue(reg.exists);
        assertTrue(reg.deregRequestHeight.exists);
        assertEq(reg.deregRequestHeight.blockHeight, 108);

        vm.expectRevert("operator must not have already requested deregistration");
        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);
    }

    function testDeregisterOperator() public {
        vm.roll(11);

        address operator = address(0x888);
        vm.expectRevert("operator must be registered");
        vm.prank(operator);
        mevCommitAVS.deregisterOperator(operator);

        testRegisterOperator();

        address otherAcct = address(0x777);
        vm.expectRevert("sender must be operator");
        vm.prank(otherAcct);
        mevCommitAVS.deregisterOperator(operator);

        vm.expectRevert("operator must have requested deregistration");
        vm.prank(operator);
        mevCommitAVS.deregisterOperator(operator);

        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);

        vm.expectRevert("deregistration must happen at least operatorDeregPeriodBlocks after deregistration request height");
        vm.prank(operator);
        mevCommitAVS.deregisterOperator(operator);

        IMevCommitAVS.OperatorRegistrationInfo memory operatorRegInfo = mevCommitAVS.getOperatorRegInfo(operator);
        assertTrue(operatorRegInfo.exists);
        assertTrue(operatorRegInfo.deregRequestHeight.exists);
        assertEq(operatorRegInfo.deregRequestHeight.blockHeight, 11);

        avsDirectoryMock.registerOperator(operator);

        vm.roll(11 + operatorDeregPeriodBlocks);

        vm.expectEmit(true, true, true, true);
        emit OperatorDeregistered(operator);
        vm.prank(operator);
        mevCommitAVS.deregisterOperator(operator);
        assertFalse(avsDirectoryMock.isRegisteredOperator(operator));

        operatorRegInfo = mevCommitAVS.getOperatorRegInfo(operator);
        assertFalse(operatorRegInfo.exists);
        assertFalse(operatorRegInfo.deregRequestHeight.exists);
        assertEq(operatorRegInfo.deregRequestHeight.blockHeight, 0);
    }

    function testRegisterValidatorsByPodOwners() public {
        vm.roll(55);

        address operator = address(0x888);
        address podOwner = address(0x420);
        ISignatureUtils.SignatureWithExpiry memory sig = ISignatureUtils.SignatureWithExpiry({
            signature: bytes("signature"),
            expiry: 10
        });
        vm.prank(podOwner);
        delegationManagerMock.delegateTo(operator, sig, bytes32("salt"));

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");
        bytes[][] memory arrayValPubkeys = new bytes[][](1);
        arrayValPubkeys[0] = valPubkeys;
        address[] memory podOwners = new address[](1);
        podOwners[0] = podOwner;

        address otherAcct = address(0x777);
        vm.expectRevert("sender must be podOwner or delegated operator");
        vm.prank(otherAcct);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);

        vm.expectRevert("delegated operator must be registered with MevCommitAVS");
        vm.prank(podOwner);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);

        testRegisterOperator();

        EigenPodMock mockPod = new EigenPodMock();
        mockPod.setMockValidatorInfo(valPubkeys[0], IEigenPod.ValidatorInfo({
            validatorIndex: 1,
            restakedBalanceGwei: 1,
            mostRecentBalanceUpdateTimestamp: 1,
            status: IEigenPod.VALIDATOR_STATUS.INACTIVE
        }));
        eigenPodManagerMock.setMockPod(podOwner, mockPod);

        vm.expectRevert("validator must be active under pod");
        vm.prank(podOwner);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);

        mockPod.setMockValidatorInfo(valPubkeys[0], IEigenPod.ValidatorInfo({
            validatorIndex: 1,
            restakedBalanceGwei: 1,
            mostRecentBalanceUpdateTimestamp: 1,
            status: IEigenPod.VALIDATOR_STATUS.ACTIVE
        }));

        mockPod.setMockValidatorInfo(valPubkeys[1], IEigenPod.ValidatorInfo({
            validatorIndex: 2,
            restakedBalanceGwei: 1,
            mostRecentBalanceUpdateTimestamp: 1,
            status: IEigenPod.VALIDATOR_STATUS.ACTIVE
        }));

        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo0 = mevCommitAVS.getValidatorRegInfo(valPubkeys[0]);
        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo1 = mevCommitAVS.getValidatorRegInfo(valPubkeys[1]);
        assertFalse(regInfo0.exists);
        assertFalse(regInfo1.exists);
        assertEq(regInfo0.podOwner, address(0));
        assertEq(regInfo1.podOwner, address(0));
        assertFalse(regInfo0.freezeHeight.exists);
        assertFalse(regInfo1.freezeHeight.exists);
        assertFalse(regInfo0.deregRequestHeight.exists);
        assertFalse(regInfo1.deregRequestHeight.exists);

        vm.expectEmit(true, true, true, true);
        emit ValidatorRegistered(valPubkeys[0], podOwner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorRegistered(valPubkeys[1], podOwner);
        vm.prank(podOwner);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);

        regInfo0 = mevCommitAVS.getValidatorRegInfo(valPubkeys[0]);
        regInfo1 = mevCommitAVS.getValidatorRegInfo(valPubkeys[1]);
        assertTrue(regInfo0.exists);
        assertTrue(regInfo1.exists);
        assertEq(regInfo0.podOwner, podOwner);
        assertEq(regInfo1.podOwner, podOwner);
        assertFalse(regInfo0.freezeHeight.exists);
        assertFalse(regInfo1.freezeHeight.exists);
        assertFalse(regInfo0.deregRequestHeight.exists);
        assertFalse(regInfo1.deregRequestHeight.exists);

        vm.expectRevert("validator must not be registered");
        vm.prank(podOwner);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);
    }

    function testRequestValidatorsDeregistration() public {

        address operator = address(0x888);
        address podOwner = address(0x420);

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");

        vm.expectRevert("validator must be registered");
        vm.prank(podOwner);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys);

        testRegisterValidatorsByPodOwners();
        vm.roll(103);

        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(valPubkeys[0], podOwner);
        bytes[] memory valPubkeys2 = new bytes[](1);
        valPubkeys2[0] = bytes("valPubkey1");
        vm.prank(podOwner);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys2);

        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistrationRequested(valPubkeys[1], podOwner);
        bytes[] memory valPubkeys3 = new bytes[](1);
        valPubkeys3[0] = bytes("valPubkey2");
        vm.prank(operator);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys3);

        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo0 = mevCommitAVS.getValidatorRegInfo(valPubkeys[0]);
        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo1 = mevCommitAVS.getValidatorRegInfo(valPubkeys[1]);
        assertTrue(regInfo0.exists);
        assertTrue(regInfo1.exists);
        assertEq(regInfo0.podOwner, podOwner);
        assertEq(regInfo1.podOwner, podOwner);
        assertTrue(regInfo0.deregRequestHeight.exists);
        assertTrue(regInfo1.deregRequestHeight.exists);
        assertEq(regInfo0.deregRequestHeight.blockHeight, 103);
        assertEq(regInfo1.deregRequestHeight.blockHeight, 103);
        assertFalse(regInfo0.freezeHeight.exists);
        assertFalse(regInfo1.freezeHeight.exists);

        vm.expectRevert("validator must not have already requested deregistration");
        vm.prank(podOwner);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys);
    }

    function testDeregisterValidator() public {

        address operator = address(0x888);
        address podOwner = address(0x420);
        bytes[] memory valPubkeys = new bytes[](1);
        valPubkeys[0] = bytes("valPubkey1");
        vm.expectRevert("validator must be registered");
        vm.prank(podOwner);
        mevCommitAVS.deregisterValidators(valPubkeys);

        testRegisterValidatorsByPodOwners();

        address otherAcct = address(0x777);
        vm.expectRevert("sender must be podOwner or delegated operator of validator");
        vm.prank(otherAcct);
        mevCommitAVS.deregisterValidators(valPubkeys);

        vm.expectRevert("validator must have requested deregistration");
        vm.prank(podOwner);
        mevCommitAVS.deregisterValidators(valPubkeys);

        bytes[] memory valPubkeys2 = new bytes[](2);
        valPubkeys2[0] = bytes("valPubkey1");
        valPubkeys2[1] = bytes("valPubkey2");
        vm.prank(podOwner);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys2);

        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo0 = mevCommitAVS.getValidatorRegInfo(valPubkeys2[0]);
        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo1 = mevCommitAVS.getValidatorRegInfo(valPubkeys2[1]);
        assertTrue(regInfo0.deregRequestHeight.exists);
        assertTrue(regInfo1.deregRequestHeight.exists);

        vm.expectRevert("deregistration must happen at least validatorDeregPeriodBlocks after deregistration request height");
        vm.prank(operator);
        mevCommitAVS.deregisterValidators(valPubkeys2);

        vm.roll(103 + validatorDeregPeriodBlocks);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistered(valPubkeys2[0], podOwner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregistered(valPubkeys2[1], podOwner);
        vm.prank(operator);
        mevCommitAVS.deregisterValidators(valPubkeys2);

        regInfo0 = mevCommitAVS.getValidatorRegInfo(valPubkeys2[0]);
        regInfo1 = mevCommitAVS.getValidatorRegInfo(valPubkeys2[1]);
        assertFalse(regInfo0.exists);
        assertFalse(regInfo1.exists);
    }

    function testRegisterLSTRestaker() public {

        address operator = address(0x888);
        address lstRestaker = address(0x34443);
        address otherAcct = address(0x777);
        bytes[] memory chosenVals = new bytes[](0);
        vm.expectRevert("sender must be delegated to an operator that is registered with MevCommitAVS");
        vm.prank(otherAcct);
        mevCommitAVS.registerLSTRestaker(chosenVals);

        testRegisterValidatorsByPodOwners();

        ISignatureUtils.SignatureWithExpiry memory sig = ISignatureUtils.SignatureWithExpiry({
            signature: bytes("signature"),
            expiry: 10
        });
        vm.prank(lstRestaker);
        delegationManagerMock.delegateTo(operator, sig, bytes32("salt"));

        vm.expectRevert("LST restaker must choose at least one validator");
        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(chosenVals);

        bytes[] memory chosenVals2 = new bytes[](2);
        chosenVals2[0] = bytes("valPubkey1");
        chosenVals2[1] = bytes("valPubkey2");
        vm.expectRevert("LST restaker must have deposited into at least one strategy");
        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(chosenVals2);

        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators.length, 0);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).numChosen, 0);
        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);

        strategyManagerMock.setStakerStrategyListLengthReturnValue(3);
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerRegistered(chosenVals2[0], 2, lstRestaker);
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerRegistered(chosenVals2[1], 2, lstRestaker);
        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(chosenVals2);

        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators.length, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).numChosen, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[0], chosenVals2[0]);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[1], chosenVals2[1]);
        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);
    
        vm.expectRevert("sender must not be registered LST restaker");
        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(chosenVals2);
    }

    function testRequestLSTRestakerDeregistration() public {
        address otherAcct = address(0x777);
        vm.expectRevert("sender must be registered LST restaker");
        vm.prank(otherAcct);
        mevCommitAVS.requestLSTRestakerDeregistration();

        testRegisterLSTRestaker();
        vm.roll(177);

        address lstRestaker = address(0x34443);
        bytes[] memory chosenVals = new bytes[](2);
        chosenVals[0] = bytes("valPubkey1");
        chosenVals[1] = bytes("valPubkey2");
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerDeregistrationRequested(chosenVals[0], 2, lstRestaker);
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerDeregistrationRequested(chosenVals[1], 2, lstRestaker);
        vm.prank(lstRestaker);
        mevCommitAVS.requestLSTRestakerDeregistration();

        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators.length, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).numChosen, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[0], chosenVals[0]);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[1], chosenVals[1]);
        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.blockHeight, 177);

        vm.expectRevert("LST restaker must not have already requested deregistration");
        vm.prank(lstRestaker);
        mevCommitAVS.requestLSTRestakerDeregistration();
    }

    function testDeregisterLSTRestaker() public {

        address otherAcct = address(0x777);
        vm.expectRevert("sender must be registered LST restaker");
        vm.prank(otherAcct);
        mevCommitAVS.deregisterLSTRestaker();

        testRegisterLSTRestaker();

        vm.roll(302);

        address lstRestaker = address(0x34443);
        vm.expectRevert("LST restaker must have requested deregistration");
        vm.prank(lstRestaker);
        mevCommitAVS.deregisterLSTRestaker();

        vm.prank(lstRestaker);
        mevCommitAVS.requestLSTRestakerDeregistration();

        vm.expectRevert("deregistration must happen at least lstRestakerDeregPeriodBlocks after deregistration request height");
        vm.prank(lstRestaker);
        mevCommitAVS.deregisterLSTRestaker();

        vm.roll(302 + lstRestakerDeregPeriodBlocks);

        bytes[] memory chosenVals = new bytes[](2);
        chosenVals[0] = bytes("valPubkey1");
        chosenVals[1] = bytes("valPubkey2");

        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators.length, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).numChosen, 2);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[0], chosenVals[0]);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[1], chosenVals[1]);
        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.blockHeight, 302);

        vm.expectEmit(true, true, true, true);
        emit LSTRestakerDeregistered(chosenVals[0], 2, lstRestaker);
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerDeregistered(chosenVals[1], 2, lstRestaker);
        vm.prank(lstRestaker);
        mevCommitAVS.deregisterLSTRestaker();

        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators.length, 0);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).numChosen, 0);
        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);
    }

    function testFreeze() public {
        // TODO: test multiple vals frozen in one tx.
    }

    function testFrozenValidatorsCantDeregister() public {
    }

    function testFrozenValidatorDoesntAffectLSTRestakerDeregistration() public {
        // TODO: Confirm a chosen validator being frozen does not affect an LST restaker being able to deregister.
    }

    function testUnfreeze() public {
        // TODO: test scenario where validator was req deregistered before being frozen, and goes back to registerd after unfreeze.
        // TODO: Also confirm the unfreeze fee is fully given to reciever and not contract.
        // TODO: test multiple vals unfrozen in one tx.
    }

    function testPause() public {
        // TODO: list out all funcs that are affected
    }

    function testUnpause() public {
        // TODO: list out all funcs that are affected
    }

    function testSetters() public {
        // only owner for all
    }
}
