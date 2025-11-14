// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import {GasTankDepositor} from "../../contracts/core/GasTankDepositor.sol";

contract GasTankDepositorTest is Test {
    uint256 public constant ALICE_PK = uint256(0xA11CE);
    uint256 public constant BOB_PK = uint256(0xB0B);
    uint256 public constant RPC_SERVICE_PK = uint256(0x1234567890);
    uint256 public constant MAXIMUM_DEPOSIT = 0.01 ether;
    address public constant ZERO_ADDRESS = address(0);
    bytes32 public constant EMPTY_CODEHASH = keccak256("");

    address public alice = payable(vm.addr(ALICE_PK));
    address public bob = payable(vm.addr(BOB_PK));
    address public rpcService = payable(vm.addr(RPC_SERVICE_PK));

    GasTankDepositor private _gasTankDepositorImpl;

    function setUp() public {
        vm.deal(alice, 10 ether);
        vm.deal(bob, 10 ether);
        vm.deal(rpcService, 10 ether);
        _gasTankDepositorImpl = new GasTankDepositor(rpcService, MAXIMUM_DEPOSIT);
    }

    //=======================TESTS=======================

    function testSetsDelegationCodeAtAddress() public {
        // Initial code is empty
        assertEq(alice.code.length, 0);

        // Set delegation as the GasTankDepositor
        _signAndAttachDelegation(address(_gasTankDepositorImpl), ALICE_PK);

        assertEq(alice.codehash, _delegateCodeHash(address(_gasTankDepositorImpl)));
        assertEq(alice.code.length, 23);
    }

    function testRemovesDelegationCodeAtAddress() public {
        // Set delegation as the GasTankDepositor
        _signAndAttachDelegation(address(_gasTankDepositorImpl), ALICE_PK);
        assertEq(alice.codehash, _delegateCodeHash(address(_gasTankDepositorImpl)));
        assertEq(alice.code.length, 23);

        // Remove delegation
        _signAndAttachDelegation(ZERO_ADDRESS, ALICE_PK);
        assertEq(alice.codehash, EMPTY_CODEHASH);
        assertEq(alice.code.length, 0);
    }

    //=======================TESTS FOR IMPROPER CALLS TO THE GAS TANK MANAGER=======================

    function testFallbackRevert() public {
        bytes memory badData =
            abi.encodeWithSelector(GasTankDepositor.recoverFunds.selector, address(0x55555), 1 ether, 1 ether, 1 ether);
        vm.prank(alice);
        (bool success,) = address(_gasTankDepositorImpl).call{value: 1 ether}(badData);
        assertFalse(success);
    }

    function testReceiveNoRevert() public {
        uint256 beforeBalance = alice.balance;
        vm.prank(bob);
        (bool success,) = address(alice).call{value: 1 ether}("");
        assertTrue(success);
        uint256 afterBalance = alice.balance;
        assertEq(afterBalance, beforeBalance + 1 ether, "balance not increased");
    }

    function testFundsSentDirectlyToDelegateAddress() public {
        uint256 beforeBalance = address(_gasTankDepositorImpl).balance;

        vm.prank(bob);
        (bool success,) = address(_gasTankDepositorImpl).call{value: 1 ether}("");
        assertTrue(success);

        uint256 afterBalance = address(_gasTankDepositorImpl).balance;
        assertEq(afterBalance, beforeBalance + 1 ether, "balance not increased");
    }

    function testWithdrawsFundsDirectlyFromDelegateAddress() public {
        uint256 gasTankBeforeBalance = address(_gasTankDepositorImpl).balance;
        uint256 depositAmount = 1 ether;

        vm.prank(bob);
        (bool success,) = address(_gasTankDepositorImpl).call{value: depositAmount}("");
        assertTrue(success);

        uint256 gasTankAfterBalance = address(_gasTankDepositorImpl).balance;
        assertEq(gasTankAfterBalance, gasTankBeforeBalance + depositAmount, "balance not increased");

        vm.prank(rpcService);
        uint256 rpcServiceBeforeBalance = rpcService.balance;
        _gasTankDepositorImpl.recoverFunds();
        uint256 rpcServiceAfterBalance = rpcService.balance;

        assertEq(address(_gasTankDepositorImpl).balance, 0, "funds not drained");
        assertEq(rpcServiceAfterBalance, rpcServiceBeforeBalance + depositAmount, "balance not recovered");
    }

    function testRevertsWhenRecoverFundsIsCalledByUnknownCaller() public {
        vm.prank(bob);
        vm.expectRevert(abi.encodeWithSelector(GasTankDepositor.NotRPCService.selector, bob));
        _gasTankDepositorImpl.recoverFunds();
    }

    //=======================TESTS FOR FUNDING THE GAS TANK=======================

    function testRpcServiceFundsMaximumDeposit() public {
        _delegate();

        uint256 rpcBalanceBefore = rpcService.balance;
        _expectGasTankFunded(rpcService, MAXIMUM_DEPOSIT);

        vm.prank(rpcService);
        GasTankDepositor(payable(alice)).fundGasTank();

        assertEq(rpcService.balance, rpcBalanceBefore + MAXIMUM_DEPOSIT, "rpc balance not increased");
    }

    function testRpcServiceFundRevertsWhenCallerNotRpcService() public {
        _delegate();

        vm.prank(alice);
        vm.expectRevert(abi.encodeWithSelector(GasTankDepositor.NotRPCService.selector, alice));
        GasTankDepositor(payable(alice)).fundGasTank();
    }

    function testRpcServiceFundRevertsWhenInsufficientBalance() public {
        vm.deal(alice, MAXIMUM_DEPOSIT - 1);
        _delegate();

        vm.prank(rpcService);
        vm.expectRevert(
            abi.encodeWithSelector(GasTankDepositor.InsufficientFunds.selector, MAXIMUM_DEPOSIT - 1, MAXIMUM_DEPOSIT)
        );
        GasTankDepositor(payable(alice)).fundGasTank();
    }

    function testEOAFundsGasTank() public {
        uint256 amount = 1 ether;
        _delegate();

        uint256 rpcBalanceBefore = rpcService.balance;
        _expectGasTankFunded(alice, amount);

        vm.prank(alice);
        GasTankDepositor(payable(alice)).fundGasTank(amount);

        assertEq(rpcService.balance, rpcBalanceBefore + amount, "rpc balance not increased");
    }

    function testEOAFundRevertsWhenCallerNotEOA() public {
        _delegate();

        vm.prank(rpcService);
        vm.expectRevert(abi.encodeWithSelector(GasTankDepositor.NotThisEOA.selector, rpcService, alice));
        GasTankDepositor(payable(alice)).fundGasTank(MAXIMUM_DEPOSIT);
    }

    function testEOAFundRevertsWhenInsufficientBalance() public {
        vm.deal(alice, MAXIMUM_DEPOSIT - 1);
        _delegate();

        vm.prank(alice);
        vm.expectRevert(
            abi.encodeWithSelector(GasTankDepositor.InsufficientFunds.selector, MAXIMUM_DEPOSIT - 1, MAXIMUM_DEPOSIT)
        );
        GasTankDepositor(payable(alice)).fundGasTank(MAXIMUM_DEPOSIT);
    }

    //=======================HELPERS=======================

    function _delegate() internal {
        _signAndAttachDelegation(address(_gasTankDepositorImpl), ALICE_PK);
    }

    function _expectGasTankFunded(address caller, uint256 amount) internal {
        vm.expectEmit(true, true, true, true);
        emit GasTankDepositor.GasTankFunded(alice, caller, amount);
    }

    function _signAndAttachDelegation(address contractAddress, uint256 pk) internal {
        vm.prank(alice);
        vm.signAndAttachDelegation(contractAddress, pk);
        vm.stopPrank();
    }

    function _delegateCodeHash(address contractAddress) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(hex"ef0100", contractAddress));
    }
}
