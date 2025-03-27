// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";
import {L1Gateway} from "../../contracts/standard-bridge/L1Gateway.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IGateway} from "../../contracts/interfaces/IGateway.sol";
import {RevertingReceiver} from "./RevertingReceiver.sol";
import {EventReceiver} from "./EventReceiver.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract L1GatewayTest is Test {
    L1Gateway l1Gateway;
    address owner;
    address relayer;
    address bridgeUser;
    uint256 counterpartyFinalizationFee;

    function setUp() public {
        owner = address(this); // Original contract deployer as owner
        relayer = address(0x78);
        bridgeUser = address(0x101);
        counterpartyFinalizationFee = 0.05 ether;

        address l1GatewayProxy = Upgrades.deployUUPSProxy(
            "L1Gateway.sol",
            abi.encodeCall(L1Gateway.initialize,
            (owner, 
            relayer, 
            counterpartyFinalizationFee))); 
        l1Gateway = L1Gateway(payable(l1GatewayProxy));
    }

    function test_ConstructorSetsVariablesCorrectly() public view {
        assertEq(l1Gateway.owner(), owner);
        assertEq(l1Gateway.relayer(), relayer);
        assertEq(l1Gateway.counterpartyFinalizationFee(), counterpartyFinalizationFee);
    }
    event CounterpartyFinalizationFeeSet(uint256 counterpartyFinalizationFee);

    function test_SetCounterpartyFinalizationFee() public {
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(888))
        );
        vm.prank(vm.addr(888));
        l1Gateway.setCounterpartyFinalizationFee(0.0005 ether);

        assertEq(l1Gateway.counterpartyFinalizationFee(), 0.05 ether);
        vm.expectEmit(true, true, true, true);
        emit CounterpartyFinalizationFeeSet(0.0005 ether);
        l1Gateway.setCounterpartyFinalizationFee(0.0005 ether);
        assertEq(l1Gateway.counterpartyFinalizationFee(), 0.0005 ether);
    }

    event RelayerSet(address indexed relayer);

    function test_SetRelayer() public {

        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(888))
        );
        vm.prank(vm.addr(888));
        l1Gateway.setRelayer(address(0x123));

        assertEq(l1Gateway.relayer(), address(0x78));

        vm.expectEmit(true, true, true, true);
        emit RelayerSet(address(0x12345));
        l1Gateway.setRelayer(address(0x12345));
        assertEq(l1Gateway.relayer(), address(0x12345));
    }

    // Expected event signature emitted in initiateTransfer()
    event TransferInitiated(
        address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx, uint256 counterpartyFinalizationFee);

    event TransferNeedsWithdrawal(address indexed recipient, uint256 amount);
    event TransferSuccess(address indexed recipient, uint256 amount);

    function test_InitiateTransferSuccess() public {
        vm.deal(bridgeUser, 100 ether);
        uint256 amount = 7 ether;

        // Initial assertions
        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        // Set up expectation for event
        vm.expectEmit(true, true, true, true);
        emit TransferInitiated(bridgeUser, bridgeUser, amount, 1, counterpartyFinalizationFee);

        // Call function as bridgeUser
        vm.prank(bridgeUser);
        uint256 returnedIdx = l1Gateway.initiateTransfer{value: amount}(bridgeUser, amount);

        // Assertions after call
        assertEq(address(bridgeUser).balance, 93 ether); 
        assertEq(l1Gateway.transferInitiatedIdx(), 1);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
        assertEq(returnedIdx, 1);
    }

    function test_InitiateTransferAmountTooSmallForCounterpartyFee() public {
        vm.deal(bridgeUser, 100 ether);
        vm.deal(address(l1Gateway), 1 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.AmountTooSmall.selector, 0.04 ether, 0.05 ether));
        vm.prank(bridgeUser);
        l1Gateway.initiateTransfer{value: 0.04 ether}(bridgeUser, 0.04 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_InitiateTransferUserInsufficientBalance() public {
        vm.deal(bridgeUser, 0.01 ether);

        assertEq(address(bridgeUser).balance, 0.01 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert();
        vm.prank(bridgeUser);
        l1Gateway.initiateTransfer{value: 0.9 ether}(bridgeUser, 0.9 ether);

        assertEq(address(bridgeUser).balance, 0.01 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_InitiateTransferValueMismatch() public {
        vm.deal(bridgeUser, 100 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(L1Gateway.IncorrectEtherValueSent.selector, 0.8 ether, 0.9 ether));
        vm.prank(bridgeUser);
        l1Gateway.initiateTransfer{value: 0.8 ether}(bridgeUser, 0.9 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    event TransferFinalized(
        address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx);

    function test_FinalizeTransferSuccess() public {
        uint256 amount = 4 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.1 ether;

        // Fund gateway and relayer
        vm.deal(address(l1Gateway), 5 ether);
        vm.deal(relayer, 5 ether);

        // Initial assertions
        assertEq(address(l1Gateway).balance, 5 ether);
        assertEq(relayer.balance, 5 ether);
        assertEq(bridgeUser.balance, 0 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        // Set up expectation for event
        vm.expectEmit(true, true, true, true);
        emit TransferSuccess(bridgeUser, amount - finalizationFee);
        vm.expectEmit(true, true, true, true);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, true, true);
        emit TransferFinalized(bridgeUser, amount, counterpartyIdx);

        // Call function as relayer
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx, finalizationFee);

        // Finalization fee is 0.1 ether
        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(relayer.balance, 5.1 ether);
        assertEq(bridgeUser.balance, 3.9 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 2);
    }

    function test_OnlyRelayerCanCallFinalizeTransfer() public {
        uint256 amount = 0.1 ether;
        vm.deal(address(l1Gateway), 1 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.SenderNotRelayer.selector, bridgeUser, relayer));
        vm.prank(bridgeUser);
        l1Gateway.finalizeTransfer(address(0x101), amount, 1, 0.05 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    // This scenario shouldn't be possible since initiateTransfer() should have prevented it.
    function test_FinalizeTranferAmountTooSmallForFinalizationFee() public {
        uint256 amount = 0.09 ether;
        uint256 finalizationFee = 0.1 ether;
        vm.deal(address(l1Gateway), 1 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.AmountTooSmall.selector, 0.09 ether, 0.1 ether));
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(address(0x101), amount, 1, finalizationFee);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferInvalidCounterpartyIdx() public {
        uint256 amount = 0.1 ether;
        uint256 finalizationFee = 0.1 ether;
        vm.deal(address(l1Gateway), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.InvalidCounterpartyIndex.selector, 2, 1));
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(address(0x101), amount, 2, finalizationFee);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferWithInsufficientContractBalance() public {
        uint256 amount = 4 ether;
        uint256 counterpartyIdx = 1; // First transfer idx
        uint256 finalizationFee = 0.1 ether;
        vm.deal(address(l1Gateway), 0.09 ether);
        assertEq(address(l1Gateway).balance, 0.09 ether);

        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
        
        vm.expectRevert(abi.encodeWithSelector(L1Gateway.InsufficientContractBalance.selector, 0.09 ether, 4 ether - finalizationFee));
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx, finalizationFee);

        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferRevertingReceiver() public {
        uint256 amount = 4 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.1 ether;

        vm.deal(address(l1Gateway), 5 ether);
        vm.deal(relayer, 5 ether);

        RevertingReceiver revertingReceiver = new RevertingReceiver();
        revertingReceiver.setShouldRevert(true);
        address receiver = address(revertingReceiver);

        assertEq(address(l1Gateway).balance, 5 ether);
        assertEq(relayer.balance, 5 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), 0);

        vm.expectEmit(true, true, false, false);
        emit TransferNeedsWithdrawal(receiver, amount - finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferFinalized(receiver, amount, counterpartyIdx);

        vm.prank(relayer);
        l1Gateway.finalizeTransfer(receiver, amount, counterpartyIdx, finalizationFee);

        assertEq(address(l1Gateway).balance, 5 ether - finalizationFee);
        assertEq(relayer.balance, 5 ether + finalizationFee);
        assertEq(receiver.balance, 0 ether);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), amount - finalizationFee);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 2);

        // Any account can retry the transfer, but user is responsible for receiver functioning properly.
        vm.expectRevert(abi.encodeWithSelector(L1Gateway.TransferFailed.selector, receiver));
        vm.prank(vm.addr(88));
        l1Gateway.withdraw(receiver);

        revertingReceiver.setShouldRevert(false);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(receiver, amount - finalizationFee);
        vm.prank(vm.addr(99));
        l1Gateway.withdraw(receiver);

        assertEq(address(l1Gateway).balance, 5 ether - amount); // User pays for finalization fee
        assertEq(relayer.balance, 5.1 ether);
        assertEq(receiver.balance, 3.9 ether);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), 0);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 2);
    }

    function test_finalizeTransferEventReceiver() public {
        uint256 amount = 4 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.1 ether;

        vm.deal(address(l1Gateway), 5 ether);
        vm.deal(relayer, 5 ether);

        EventReceiver eventReceiver = new EventReceiver();
        address receiver = address(eventReceiver);

        assertEq(address(l1Gateway).balance, 5 ether);
        assertEq(relayer.balance, 5 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), 0);

        // Too much gas usage in receiver, so manual withdrawal is needed.
        vm.expectEmit(true, true, false, false);
        emit TransferNeedsWithdrawal(receiver, amount - finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferFinalized(receiver, amount, counterpartyIdx);

        vm.prank(relayer);
        l1Gateway.finalizeTransfer(receiver, amount, counterpartyIdx, finalizationFee);

        assertEq(address(l1Gateway).balance, 5 ether - finalizationFee);
        assertEq(relayer.balance, 5 ether + finalizationFee);
        assertEq(receiver.balance, 0 ether);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), amount - finalizationFee);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 2);

        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(receiver, amount - finalizationFee);
        vm.prank(vm.addr(888));
        l1Gateway.withdraw(receiver);

        assertEq(address(l1Gateway).balance, 5 ether - amount);
        assertEq(relayer.balance, 5.1 ether);
        assertEq(receiver.balance, 3.9 ether);
        assertEq(l1Gateway.transferredFundsNeedingWithdrawal(receiver), 0);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 2);
    }

    function test_WithdrawNoFundsNeedingWithdrawal() public {
        vm.expectRevert(abi.encodeWithSelector(L1Gateway.NoFundsNeedingWithdrawal.selector, vm.addr(999)));
        vm.prank(vm.addr(888));
        l1Gateway.withdraw(vm.addr(999));
    }
}
