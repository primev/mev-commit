// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {SettlementGateway} from "../../contracts/standard-bridge/SettlementGateway.sol";
import {Allocator} from "../../contracts/standard-bridge/Allocator.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IGateway} from "../../contracts/interfaces/IGateway.sol";
import {IAllocator} from "../../contracts/interfaces/IAllocator.sol";
import {RevertingReceiver} from "./RevertingReceiver.sol";
import {EventReceiver} from "./EventReceiver.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract SettlementGatewayTest is Test {

    SettlementGateway settlementGateway;
    Allocator allocator;

    address owner;
    address relayer;
    address bridgeUser;
    uint256 counterpartyFinalizationFee;

    function setUp() public {
        owner = address(this); // Original contract deployer as owner
        relayer = address(0x78);
        bridgeUser = address(0x101);
        counterpartyFinalizationFee = 0.1 ether;

        address allocatorProxy = Upgrades.deployUUPSProxy(
            "Allocator.sol",
            abi.encodeCall(Allocator.initialize, (owner))
        ); 
        allocator = Allocator(payable(allocatorProxy));
        
        address settlementGatewayProxy = Upgrades.deployUUPSProxy(
            "SettlementGateway.sol",
            abi.encodeCall(SettlementGateway.initialize, 
            (address(allocator), 
            owner, 
            relayer, 
            counterpartyFinalizationFee))
        );
        settlementGateway = SettlementGateway(payable(settlementGatewayProxy));

        vm.prank(owner);
        allocator.addToWhitelist(address(settlementGateway));
    }

    function test_ConstructorSetsVariablesCorrectly() public view {
        // Test if the constructor correctly initializes variables
        assertEq(settlementGateway.owner(), owner);
        assertEq(settlementGateway.relayer(), relayer);
        assertEq(settlementGateway.counterpartyFinalizationFee(), counterpartyFinalizationFee);
        assertEq(settlementGateway.allocatorAddr(), address(allocator));
    }

    // Expected event signature emitted in initiateTransfer()
    event TransferInitiated(
        address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx, uint256 counterpartyFinalizationFee);

    event TransferNeedsWithdrawal(address indexed recipient, uint256 amount);
    event TransferSuccess(address indexed recipient, uint256 amount);

    event CounterpartyFinalizationFeeSet(uint256 counterpartyFinalizationFee);

    function test_SetCounterpartyFinalizationFee() public {
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(888))
        );
        vm.prank(vm.addr(888));
        settlementGateway.setCounterpartyFinalizationFee(0.0005 ether);

        assertEq(settlementGateway.counterpartyFinalizationFee(), 0.1 ether);
        vm.expectEmit(true, true, true, true);
        emit CounterpartyFinalizationFeeSet(0.0005 ether);
        settlementGateway.setCounterpartyFinalizationFee(0.0005 ether);
        assertEq(settlementGateway.counterpartyFinalizationFee(), 0.0005 ether);
    }

    event RelayerSet(address indexed relayer);

    function test_SetRelayer() public {

        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(888))
        );
        vm.prank(vm.addr(888));
        settlementGateway.setRelayer(address(0x123));

        assertEq(settlementGateway.relayer(), address(0x78));

        vm.expectEmit(true, true, true, true);
        emit RelayerSet(address(0x12345));
        settlementGateway.setRelayer(address(0x12345));
        assertEq(settlementGateway.relayer(), address(0x12345));
    }

    function test_InitiateTransferSuccess() public {
        vm.deal(bridgeUser, 100 ether);
        uint256 amount = 7 ether;

        // Initial assertions
        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(address(allocator).balance, 0 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        // Set up expectation for event
        vm.expectEmit(true, true, true, true);
        emit TransferInitiated(bridgeUser, bridgeUser, amount, 1, counterpartyFinalizationFee);

        // Call function as bridgeUser
        vm.prank(bridgeUser);
        uint256 returnedIdx = settlementGateway.initiateTransfer{value: amount}(bridgeUser, amount);

        // Assertions after call
        assertEq(address(bridgeUser).balance, 93 ether);
        assertEq(address(allocator).balance, 7 ether);
        assertEq(address(settlementGateway).balance, 0 ether);

        assertEq(settlementGateway.transferInitiatedIdx(), 1);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
        assertEq(returnedIdx, 1); 
    }

    function test_InitiateTransferAmountTooSmallForCounterpartyFee() public {
        vm.deal(bridgeUser, 100 ether);
        vm.deal(address(settlementGateway), 1 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(address(settlementGateway).balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.AmountTooSmall.selector, 0.04 ether, 0.1 ether));

        vm.prank(bridgeUser);
        settlementGateway.initiateTransfer{value: 0.04 ether}(bridgeUser, 0.04 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(address(settlementGateway).balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_InitiateTransferUserInsufficientBalance() public {
        vm.deal(bridgeUser, 0.01 ether);

        assertEq(address(bridgeUser).balance, 0.01 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.prank(bridgeUser);
        // Foundry 1.2.x onward doesn't support vm.expectRevert() for catching EvmError: OutOfFunds. We'll use try/catch instead.
        try settlementGateway.initiateTransfer{value: 0.9 ether}(bridgeUser, 0.9 ether) {
            fail(); // Call should not succeed
        } catch {
            assertEq(address(bridgeUser).balance, 0.01 ether);
            assertEq(settlementGateway.transferInitiatedIdx(), 0);
            assertEq(settlementGateway.transferFinalizedIdx(), 1);
        }
    }

    function test_InitiateTransferValueMismatch() public {
        vm.deal(bridgeUser, 100 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(SettlementGateway.IncorrectEtherValueSent.selector, 0.8 ether, 0.9 ether));
        vm.prank(bridgeUser);
        settlementGateway.initiateTransfer{value: 0.8 ether}(bridgeUser, 0.9 ether);

        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    event TransferFinalized(
        address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx);

    function test_FinalizeTransferSuccess() public {
        uint256 amount = 2 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.05 ether;

        // Fund allocator and relayer
        vm.deal(address(allocator), 3 ether);
        vm.deal(relayer, 3 ether);

        // Initial assertions
        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(bridgeUser.balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        // Set up expectation for event
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(bridgeUser, amount-finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, true, true);
        emit TransferFinalized(bridgeUser, amount, counterpartyIdx);

        // Call function as relayer
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx, finalizationFee);

        // Final assertions
        assertEq(address(allocator).balance, 1 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(bridgeUser.balance, 1.95 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);

        // One more
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(bridgeUser, 0.45 ether);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, 0.05 ether);
        vm.expectEmit(true, true, true, true);
        emit TransferFinalized(bridgeUser, 0.5 ether, counterpartyIdx+1);
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, 0.5 ether, counterpartyIdx+1, finalizationFee);
    }

    function test_OnlyRelayerCanCallFinalizeTransfer() public {
        uint256 amount = 0.1 ether;
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 3 ether);
        vm.deal(relayer, 3 ether);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.SenderNotRelayer.selector, bridgeUser, relayer));
        vm.prank(bridgeUser);
        settlementGateway.finalizeTransfer(address(0x101), amount, 1, finalizationFee);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferAmountTooSmallForFinalizationFee() public {
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.AmountTooSmall.selector, 0.04 ether, 0.05 ether));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, 0.04 ether, 1, finalizationFee);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferInvalidCounterpartyIdx() public {
        uint256 amount = 0.1 ether;
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.InvalidCounterpartyIndex.selector, 7, 1));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, 7, finalizationFee);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferWithInsufficientContractBalance() public {
        uint256 amount = 0.1 ether;
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 0.04 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 0.04 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IAllocator.InsufficientContractBalance.selector, 0.04 ether, 0.1 ether - finalizationFee));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, 1, finalizationFee);

        assertEq(address(allocator).balance, 0.04 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferRevertingReceiver() public {
        uint256 amount = 2 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 3 ether);
        vm.deal(relayer, 3 ether);

        RevertingReceiver revertingReceiver = new RevertingReceiver();
        revertingReceiver.setShouldRevert(true);
        address receiver = address(revertingReceiver);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectEmit(true, true, false, false);
        emit TransferNeedsWithdrawal(receiver, amount - finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferFinalized(receiver, amount, counterpartyIdx);

        vm.prank(relayer);
        settlementGateway.finalizeTransfer(receiver, amount, counterpartyIdx, finalizationFee);

        assertEq(address(allocator).balance, 2.95 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);

        // Any account can retry the transfer, but user is responsible for receiver functioning properly.
        vm.expectRevert(abi.encodeWithSelector(IAllocator.TransferFailed.selector, receiver));
        vm.prank(vm.addr(88));
        allocator.withdraw(receiver);

        revertingReceiver.setShouldRevert(false);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(receiver, amount - finalizationFee);
        vm.prank(vm.addr(99));
        allocator.withdraw(receiver);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(receiver.balance, 1.95 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);
    }

    function test_finalizeTransferEventReceiver() public {
        uint256 amount = 2 ether;
        uint256 counterpartyIdx = 1;
        uint256 finalizationFee = 0.05 ether;
        vm.deal(address(allocator), 3 ether);
        vm.deal(relayer, 3 ether);

        EventReceiver eventReceiver = new EventReceiver();
        address receiver = address(eventReceiver);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
        assertEq(allocator.transferredFundsNeedingWithdrawal(receiver), 0);

        // Too much gas usage in receiver, so manual withdrawal is needed.
        vm.expectEmit(true, true, false, false);
        emit TransferNeedsWithdrawal(receiver, amount - finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(relayer, finalizationFee);
        vm.expectEmit(true, true, false, false);
        emit TransferFinalized(receiver, amount, counterpartyIdx);

        vm.prank(relayer);
        settlementGateway.finalizeTransfer(receiver, amount, counterpartyIdx, finalizationFee);

        assertEq(address(allocator).balance, 2.95 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(receiver.balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);

        vm.expectEmit(true, true, false, false);
        emit TransferSuccess(receiver, 1.95 ether);
        vm.prank(vm.addr(888));
        allocator.withdraw(receiver);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(receiver.balance, 1.95 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);
    }

    function test_WithdrawNoFundsNeedingWithdrawal() public {
        vm.expectRevert(abi.encodeWithSelector(IAllocator.NoFundsNeedingWithdrawal.selector, vm.addr(999)));
        vm.prank(vm.addr(888));
        allocator.withdraw(vm.addr(999));
    }
}
