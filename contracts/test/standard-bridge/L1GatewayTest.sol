// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "../../contracts/standard-bridge/L1Gateway.sol";

contract L1GatewayTest is Test {
    L1Gateway l1Gateway;
    address owner;
    address relayer;
    address bridgeUser;
    uint256 finalizationFee;
    uint256 counterpartyFee;

    function setUp() public {
        owner = address(this); // Original contract deployer as owner
        relayer = address(0x78);
        bridgeUser = address(0x101);
        finalizationFee = 0.1 ether;
        counterpartyFee = 0.05 ether;
        l1Gateway = new L1Gateway(owner, relayer, finalizationFee, counterpartyFee);
    }

    function test_ConstructorSetsVariablesCorrectly() public {
        assertEq(l1Gateway.owner(), owner);
        assertEq(l1Gateway.relayer(), relayer);
        assertEq(l1Gateway.finalizationFee(), finalizationFee);
        assertEq(l1Gateway.counterpartyFee(), counterpartyFee);
    }

    // Expected event signature emitted in initiateTransfer()
    event TransferInitiated(
        address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx);

    function test_InitiateTransferSuccess() public {
        vm.deal(bridgeUser, 100 ether);
        uint256 amount = 7 ether;

        // Initial assertions
        assertEq(address(bridgeUser).balance, 100 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        // Set up expectation for event
        vm.expectEmit(true, true, true, true);
        emit TransferInitiated(bridgeUser, bridgeUser, amount, 1);

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

        vm.expectRevert("Amount must cover counterpartys finalization fee");
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

        vm.expectRevert("Incorrect Ether value sent");
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
        emit TransferFinalized(bridgeUser, amount, counterpartyIdx);

        // Call function as relayer
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx);

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

        vm.expectRevert("Only relayer can call this function");
        vm.prank(bridgeUser);
        l1Gateway.finalizeTransfer(address(0x101), amount, 1);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    // This scenario shouldn't be possible since initiateTransfer() should have prevented it.
    function test_FinalizeTranferAmountTooSmallForFinalizationFee() public {
        uint256 amount = 0.09 ether;
        vm.deal(address(l1Gateway), 1 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert("Amount must cover finalization fee");
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(address(0x101), amount, 1);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferInvalidCounterpartyIdx() public {
        uint256 amount = 0.1 ether;
        vm.deal(address(l1Gateway), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);

        vm.expectRevert("Invalid counterparty index. Transfers must be relayed FIFO");
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(address(0x101), amount, 2);

        assertEq(address(l1Gateway).balance, 1 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferWithInsufficientContractBalance() public {
        uint256 amount = 4 ether;
        uint256 counterpartyIdx = 1; // First transfer idx
        vm.deal(address(l1Gateway), 0.09 ether);
        assertEq(address(l1Gateway).balance, 0.09 ether);

        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
        
        vm.expectRevert("Insufficient contract balance");
        vm.prank(relayer);
        l1Gateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx);

        assertEq(l1Gateway.transferInitiatedIdx(), 0);
        assertEq(l1Gateway.transferFinalizedIdx(), 1);
    }
}
