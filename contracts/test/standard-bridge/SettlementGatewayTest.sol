// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {SettlementGateway} from "../../contracts/standard-bridge/SettlementGateway.sol";
import {Allocator} from "../../contracts/standard-bridge/Allocator.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IGateway} from "../../contracts/interfaces/IGateway.sol";
import {IAllocator} from "../../contracts/interfaces/IAllocator.sol";

contract SettlementGatewayTest is Test {

    SettlementGateway settlementGateway;
    Allocator allocator;

    address owner;
    address relayer;
    address bridgeUser;
    uint256 finalizationFee;
    uint256 counterpartyFee;

    function setUp() public {
        owner = address(this); // Original contract deployer as owner
        relayer = address(0x78);
        bridgeUser = address(0x101);
        finalizationFee = 0.05 ether;
        counterpartyFee = 0.1 ether;

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
            finalizationFee, 
            counterpartyFee))
        );
        settlementGateway = SettlementGateway(payable(settlementGatewayProxy));

        vm.prank(owner);
        allocator.addToWhitelist(address(settlementGateway));
    }

    function test_ConstructorSetsVariablesCorrectly() public view {
        // Test if the constructor correctly initializes variables
        assertEq(settlementGateway.owner(), owner);
        assertEq(settlementGateway.relayer(), relayer);
        assertEq(settlementGateway.finalizationFee(), finalizationFee);
        assertEq(settlementGateway.counterpartyFee(), counterpartyFee);
        assertEq(settlementGateway.allocatorAddr(), address(allocator));
    }

    // Expected event signature emitted in initiateTransfer()
    event TransferInitiated(
        address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx);

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
        emit TransferInitiated(bridgeUser, bridgeUser, amount, 1);

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

        vm.expectRevert();
        vm.prank(bridgeUser);
        settlementGateway.initiateTransfer{value: 0.9 ether}(bridgeUser, 0.9 ether);

        assertEq(address(bridgeUser).balance, 0.01 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
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
        vm.expectEmit(true, true, true, true);
        emit TransferFinalized(bridgeUser, amount, counterpartyIdx);

        // Call function as relayer
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, counterpartyIdx);

        // Final assertions
        assertEq(address(allocator).balance, 1 ether);
        assertEq(relayer.balance, 3.05 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(bridgeUser.balance, 1.95 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 2);

        // One more
        vm.expectEmit(true, true, true, true);
        emit TransferFinalized(bridgeUser, 0.5 ether, counterpartyIdx+1);
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, 0.5 ether, counterpartyIdx+1);
    }

    function test_OnlyRelayerCanCallFinalizeTransfer() public {
        uint256 amount = 0.1 ether;
        vm.deal(address(allocator), 3 ether);
        vm.deal(relayer, 3 ether);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.SenderNotRelayer.selector, bridgeUser, relayer));
        vm.prank(bridgeUser);
        settlementGateway.finalizeTransfer(address(0x101), amount, 1);

        assertEq(address(allocator).balance, 3 ether);
        assertEq(relayer.balance, 3 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferAmountTooSmallForFinalizationFee() public {
        vm.deal(address(allocator), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.AmountTooSmall.selector, 0.04 ether, 0.05 ether));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, 0.04 ether, 1);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferInvalidCounterpartyIdx() public {
        uint256 amount = 0.1 ether;
        vm.deal(address(allocator), 1 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IGateway.InvalidCounterpartyIndex.selector, 7, 1));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, 7);

        assertEq(address(allocator).balance, 1 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }

    function test_FinalizeTransferWithInsufficientContractBalance() public {
        uint256 amount = 0.1 ether;
        vm.deal(address(allocator), 0.04 ether);
        vm.deal(relayer, 1 ether);

        assertEq(address(allocator).balance, 0.04 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);

        vm.expectRevert(abi.encodeWithSelector(IAllocator.InsufficientContractBalance.selector, 0.04 ether, 0.1 ether - finalizationFee));
        vm.prank(relayer);
        settlementGateway.finalizeTransfer(bridgeUser, amount, 1);

        assertEq(address(allocator).balance, 0.04 ether);
        assertEq(address(settlementGateway).balance, 0 ether);
        assertEq(relayer.balance, 1 ether);
        assertEq(settlementGateway.transferInitiatedIdx(), 0);
        assertEq(settlementGateway.transferFinalizedIdx(), 1);
    }
}
