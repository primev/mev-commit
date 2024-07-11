// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(operator);
        mevCommitAVS.registerOperator(operatorSignature);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(operator);
        mevCommitAVS.deregisterOperator(operator);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(otherAcct);
        mevCommitAVS.registerValidatorsByPodOwners(arrayValPubkeys, podOwners);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(podOwner);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(podOwner);
        mevCommitAVS.deregisterValidators(valPubkeys);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(chosenVals);
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(otherAcct);
        mevCommitAVS.requestLSTRestakerDeregistration();
        vm.prank(owner);
        mevCommitAVS.unpause();

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

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(otherAcct);
        mevCommitAVS.deregisterLSTRestaker();
        vm.prank(owner);
        mevCommitAVS.unpause();

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
        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);
        vm.prank(owner);
        mevCommitAVS.unpause();

        vm.expectRevert("validator must be registered");
        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);

        testRegisterValidatorsByPodOwners();

        vm.roll(403);

        address operator = address(0x888);
        bytes[] memory valPubkeys2 = new bytes[](1);
        valPubkeys2[0] = bytes("valPubkey1");
        vm.prank(operator);
        mevCommitAVS.requestValidatorsDeregistration(valPubkeys2);

        address podOwner = address(0x420);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).exists);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).podOwner, podOwner);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).podOwner, podOwner);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).deregRequestHeight.exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).deregRequestHeight.blockHeight, 403);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).deregRequestHeight.exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);

        vm.roll(461);

        vm.expectEmit(true, true, true, true);
        emit ValidatorFrozen(valPubkeys[0], podOwner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorFrozen(valPubkeys[1], podOwner);
        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).exists);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).podOwner, podOwner);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).podOwner, podOwner); 
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.blockHeight, 461);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.blockHeight, 461);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).deregRequestHeight.exists);
        assertEq(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).deregRequestHeight.blockHeight, 403);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).deregRequestHeight.exists);

        vm.expectRevert("validator must not already be frozen");
        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);
    }

    function testFrozenValidatorsCantDeregister() public {
        testFreeze();

        bytes[] memory valPubkeys = new bytes[](1);
        valPubkeys[0] = bytes("valPubkey1");
        IMevCommitAVS.ValidatorRegistrationInfo memory regInfo = mevCommitAVS.getValidatorRegInfo(valPubkeys[0]);
        assertTrue(regInfo.deregRequestHeight.exists);

        vm.roll(block.number + validatorDeregPeriodBlocks);

        address podOwner = address(0x420);
        vm.expectRevert("frozen validator cannot deregister");
        vm.prank(podOwner);
        mevCommitAVS.deregisterValidators(valPubkeys);
    }

    function testFrozenValidatorDoesntAffectLSTRestaker() public {
        testFreeze();
        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);

        address lstRestaker = address(0x34443);
        address operator = address(0x888);
        ISignatureUtils.SignatureWithExpiry memory sig = ISignatureUtils.SignatureWithExpiry({
            signature: bytes("signature"),
            expiry: 10
        });
        vm.prank(lstRestaker);
        delegationManagerMock.delegateTo(operator, sig, bytes32("salt"));
        strategyManagerMock.setStakerStrategyListLengthReturnValue(3);

        vm.prank(lstRestaker);
        mevCommitAVS.registerLSTRestaker(valPubkeys);

        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[0], valPubkeys[0]);
        assertEq(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).chosenValidators[1], valPubkeys[1]);

        vm.prank(lstRestaker);
        mevCommitAVS.requestLSTRestakerDeregistration();

        assertTrue(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).deregRequestHeight.exists);

        vm.roll(block.number + lstRestakerDeregPeriodBlocks);

        vm.prank(lstRestaker);
        mevCommitAVS.deregisterLSTRestaker();

        assertFalse(mevCommitAVS.getLSTRestakerRegInfo(lstRestaker).exists);
    }

    function testUnfreeze() public {

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");

        address newAccount = address(0x333333333);

        vm.prank(owner);
        mevCommitAVS.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        vm.prank(newAccount);
        mevCommitAVS.unfreeze(valPubkeys);
        vm.prank(owner);
        mevCommitAVS.unpause();

        vm.expectRevert("validator must be registered");
        vm.prank(newAccount);
        mevCommitAVS.unfreeze(valPubkeys);

        testRegisterValidatorsByPodOwners();

        vm.expectRevert("validator must be frozen");
        vm.prank(newAccount);
        mevCommitAVS.unfreeze(valPubkeys);

        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).exists);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).exists);
        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);

        vm.expectRevert("sender must pay at least the unfreeze fee for each validator");
        vm.prank(newAccount);
        mevCommitAVS.unfreeze(valPubkeys);

        vm.deal(newAccount, 2 * unfreezeFee);

        uint256 singleUnfreezeFee = unfreezeFee;
        vm.expectRevert("sender must pay at least the unfreeze fee for each validator");
        vm.prank(newAccount);
        mevCommitAVS.unfreeze{value: singleUnfreezeFee}(valPubkeys);

        uint256 doubleUnfreezeFee = unfreezeFee * 2;
        vm.expectRevert("unfreeze must happen at least unfreezePeriodBlocks after freeze height");
        vm.prank(newAccount);
        mevCommitAVS.unfreeze{value: doubleUnfreezeFee}(valPubkeys);

        vm.roll(block.number + unfreezePeriodBlocks);

        assertEq(address(mevCommitAVS).balance, 0);
        assertEq(address(newAccount).balance, unfreezeFee * 2);
        assertEq(address(unfreezeReceiver).balance, 0);

        address podOwner = address(0x420);
        vm.expectEmit(true, true, true, true);
        emit ValidatorUnfrozen(valPubkeys[0], podOwner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorUnfrozen(valPubkeys[1], podOwner);
        vm.prank(newAccount);
        mevCommitAVS.unfreeze{value: doubleUnfreezeFee}(valPubkeys);

        assertEq(address(mevCommitAVS).balance, 0);
        assertEq(address(newAccount).balance, 0);
        assertEq(address(unfreezeReceiver).balance, unfreezeFee * 2);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);
    }

    function testSetters() public {
        IAVSDirectory newAVSDirectory = new AVSDirectoryMock();
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setAVSDirectory(newAVSDirectory);
        vm.expectEmit(true, true, true, true);
        emit AVSDirectorySet(address(newAVSDirectory));
        vm.prank(owner);
        mevCommitAVS.setAVSDirectory(newAVSDirectory);

        IStrategyManager newStrategyManager = IStrategyManager(address(0x8));
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setStrategyManager(newStrategyManager);
        vm.expectEmit(true, true, true, true);
        emit StrategyManagerSet(address(newStrategyManager));
        vm.prank(owner);
        mevCommitAVS.setStrategyManager(newStrategyManager);

        IDelegationManager newDelegationManager = IDelegationManager(address(0x9));
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setDelegationManager(newDelegationManager);
        vm.expectEmit(true, true, true, true);
        emit DelegationManagerSet(address(newDelegationManager));
        vm.prank(owner);
        mevCommitAVS.setDelegationManager(newDelegationManager);

        IEigenPodManager newEigenPodManager = IEigenPodManager(address(0xA));
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setEigenPodManager(newEigenPodManager);
        vm.expectEmit(true, true, true, true);
        emit EigenPodManagerSet(address(newEigenPodManager));
        vm.prank(owner);
        mevCommitAVS.setEigenPodManager(newEigenPodManager);

        address[] memory newRestakeableStrategies = new address[](2);
        newRestakeableStrategies[0] = address(0xB);
        newRestakeableStrategies[1] = address(0xC);
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setRestakeableStrategies(newRestakeableStrategies);
        vm.expectEmit(true, true, true, true);
        emit RestakeableStrategiesSet(newRestakeableStrategies);
        vm.prank(owner);
        mevCommitAVS.setRestakeableStrategies(newRestakeableStrategies);

        address newFreezeOracle = address(0xD);
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setFreezeOracle(newFreezeOracle);
        vm.expectEmit(true, true, true, true);
        emit FreezeOracleSet(newFreezeOracle);
        vm.prank(owner);
        mevCommitAVS.setFreezeOracle(newFreezeOracle);
        assertEq(mevCommitAVS.freezeOracle(), newFreezeOracle);

        uint256 newUnfreezeFee = 2 ether;
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setUnfreezeFee(newUnfreezeFee);
        vm.expectEmit(true, true, true, true);
        emit UnfreezeFeeSet(newUnfreezeFee);
        vm.prank(owner);
        mevCommitAVS.setUnfreezeFee(newUnfreezeFee);
        assertEq(mevCommitAVS.unfreezeFee(), newUnfreezeFee);

        address newUnfreezeReceiver = address(0xE);
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setUnfreezeReceiver(newUnfreezeReceiver);
        vm.expectEmit(true, true, true, true);
        emit UnfreezeReceiverSet(newUnfreezeReceiver);
        vm.prank(owner);
        mevCommitAVS.setUnfreezeReceiver(newUnfreezeReceiver);
        assertEq(mevCommitAVS.unfreezeReceiver(), newUnfreezeReceiver);

        uint256 newUnfreezePeriodBlocks = 200;
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setUnfreezePeriodBlocks(newUnfreezePeriodBlocks);
        vm.expectEmit(true, true, true, true);
        emit UnfreezePeriodBlocksSet(newUnfreezePeriodBlocks);
        vm.prank(owner);
        mevCommitAVS.setUnfreezePeriodBlocks(newUnfreezePeriodBlocks);
        assertEq(mevCommitAVS.unfreezePeriodBlocks(), newUnfreezePeriodBlocks);

        uint256 newOperatorDeregPeriodBlocks = 300;
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setOperatorDeregPeriodBlocks(newOperatorDeregPeriodBlocks);
        vm.expectEmit(true, true, true, true);
        emit OperatorDeregPeriodBlocksSet(newOperatorDeregPeriodBlocks);
        vm.prank(owner);
        mevCommitAVS.setOperatorDeregPeriodBlocks(newOperatorDeregPeriodBlocks);
        assertEq(mevCommitAVS.operatorDeregPeriodBlocks(), newOperatorDeregPeriodBlocks);

        uint256 newValidatorDeregPeriodBlocks = 400;
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setValidatorDeregPeriodBlocks(newValidatorDeregPeriodBlocks);
        vm.expectEmit(true, true, true, true);
        emit ValidatorDeregPeriodBlocksSet(newValidatorDeregPeriodBlocks);
        vm.prank(owner);
        mevCommitAVS.setValidatorDeregPeriodBlocks(newValidatorDeregPeriodBlocks);
        assertEq(mevCommitAVS.validatorDeregPeriodBlocks(), newValidatorDeregPeriodBlocks);

        uint256 newLstRestakerDeregPeriodBlocks = 500;
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.setLstRestakerDeregPeriodBlocks(newLstRestakerDeregPeriodBlocks);
        vm.expectEmit(true, true, true, true);
        emit LSTRestakerDeregPeriodBlocksSet(newLstRestakerDeregPeriodBlocks);
        vm.prank(owner);
        mevCommitAVS.setLstRestakerDeregPeriodBlocks(newLstRestakerDeregPeriodBlocks);
        assertEq(mevCommitAVS.lstRestakerDeregPeriodBlocks(), newLstRestakerDeregPeriodBlocks);

        string memory newMetadataURI = "https://new-metadata-uri.com";
        vm.prank(address(0x1));
        vm.expectRevert();
        mevCommitAVS.updateMetadataURI(newMetadataURI);
        vm.prank(owner);
        mevCommitAVS.updateMetadataURI(newMetadataURI);
    }

    function testValidatorsWithReqDeregisteredOperatorsCannotRegister() public {
        testRegisterOperator();

        address operator = address(0x888);
        vm.prank(operator);
        mevCommitAVS.requestOperatorDeregistration(operator);

        address podOwner = address(0x420);
        vm.prank(podOwner);
        ISignatureUtils.SignatureWithExpiry memory sig = ISignatureUtils.SignatureWithExpiry({
            signature: bytes("signature"),
            expiry: 10
        });
        delegationManagerMock.delegateTo(operator, sig, bytes32("salt"));

        bytes[][] memory valPubkeys = new bytes[][](1);
        bytes[] memory inner = new bytes[](2);
        inner[0] = bytes("valPubkey1");
        inner[1] = bytes("valPubkey2");
        valPubkeys[0] = inner;

        address[] memory podOwners = new address[](1);
        podOwners[0] = podOwner;
        vm.prank(podOwner);
        vm.expectRevert("delegated operator must not have requested deregistration");
        mevCommitAVS.registerValidatorsByPodOwners(valPubkeys, podOwners);

        vm.prank(operator);
        vm.expectRevert("delegated operator must not have requested deregistration");
        mevCommitAVS.registerValidatorsByPodOwners(valPubkeys, podOwners);
    }

    function testUnfreezeExcessFeeReturned() public {
        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");

        address newAccount = address(0x333333333);

        testRegisterValidatorsByPodOwners();

        vm.prank(freezeOracle);
        mevCommitAVS.freeze(valPubkeys);

        uint256 tripleUnfreezeFee = unfreezeFee * 3;
        vm.deal(newAccount, tripleUnfreezeFee);

        vm.roll(block.number + unfreezePeriodBlocks);

        assertEq(address(newAccount).balance, tripleUnfreezeFee);
        assertEq(address(mevCommitAVS).balance, 0);
        assertEq(address(unfreezeReceiver).balance, 0);

        address podOwner = address(0x420);
        vm.expectEmit(true, true, true, true);
        emit ValidatorUnfrozen(valPubkeys[0], podOwner);
        vm.expectEmit(true, true, true, true);
        emit ValidatorUnfrozen(valPubkeys[1], podOwner);
        vm.prank(newAccount);
        mevCommitAVS.unfreeze{value: tripleUnfreezeFee}(valPubkeys);

        assertEq(address(newAccount).balance, unfreezeFee);
        assertEq(address(mevCommitAVS).balance, 0);
        assertEq(address(unfreezeReceiver).balance, 2 * unfreezeFee);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[0]).freezeHeight.exists);

        assertTrue(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).exists);
        assertFalse(mevCommitAVS.getValidatorRegInfo(valPubkeys[1]).freezeHeight.exists);
    }
}
