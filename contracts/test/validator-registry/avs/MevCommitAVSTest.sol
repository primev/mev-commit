// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {StrategyManagerMock} from "eigenlayer-contracts/src/test/mocks/StrategyManagerMock.sol";
import {DelegationManagerMock} from "eigenlayer-contracts/src/test/mocks/DelegationManagerMock.sol";
import {EigenPodManagerMock} from "eigenlayer-contracts/src/test/mocks/EigenPodManagerMock.sol";
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
        vm.expectRevert("sender must be operator or MevCommitAVS owner");
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
        vm.prank(owner);
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
        vm.expectRevert("sender must be operator or MevCommitAVS owner");
        vm.prank(otherAcct);
        mevCommitAVS.deregisterOperator(operator);

        vm.expectRevert("operator must have requested deregistration");
        vm.prank(owner);
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
}
