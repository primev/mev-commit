// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import {GasTankDepositor} from "../../contracts/core/GasTankDepositor.sol";
import {Errors} from "../../contracts/utils/Errors.sol";

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

    //=======================TESTS FOR CONSTRUCTOR=======================

    function testConstructorRevertsWhenRpcServiceIsZeroAddress() public {
        vm.expectRevert(abi.encodeWithSelector(GasTankDepositor.RPCServiceNotSet.selector, ZERO_ADDRESS));
        new GasTankDepositor(ZERO_ADDRESS, MAXIMUM_DEPOSIT);
    }

    function testConstructorRevertsWhenMaximumDepositIsZero() public {
        vm.expectRevert(abi.encodeWithSelector(GasTankDepositor.MaximumDepositNotMet.selector, 0, 0));
        new GasTankDepositor(rpcService, 0);
    }

    function testConstructorSetsVariables() public view {
        assertEq(_gasTankDepositorImpl.RPC_SERVICE(), rpcService);
        assertEq(_gasTankDepositorImpl.MAXIMUM_DEPOSIT(), MAXIMUM_DEPOSIT);
        assertEq(_gasTankDepositorImpl.GAS_TANK_ADDRESS(), address(_gasTankDepositorImpl));
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

    //=======================TESTS FOR RECEIVE AND FALLBACK=======================

    function testFallbackRevert() public {
        bytes memory badData = abi.encodeWithSelector(bytes4(keccak256("invalidFunction()")));
        vm.prank(alice);
        (bool success,) = address(_gasTankDepositorImpl).call{value: 1 ether}(badData);
        assertFalse(success);
    }

    function testFundsSentDirectlyToDelegateAddress() public {
        vm.prank(bob);
        (bool success, bytes memory data) = address(_gasTankDepositorImpl).call{value: 1 ether}("");
        assertFalse(success);
        bytes4 selector;
        assembly {
            selector := mload(add(data, 0x20))
        }
        assertEq(selector, Errors.InvalidReceive.selector);
    }

    function testFundsSentDirectlyToEOAAddressWithDelegation() public {
        _delegate();

        uint256 beforeBalance = alice.balance;
        vm.prank(bob);
        (bool success,) = alice.call{value: 1 ether}("");
        assertTrue(success);
        uint256 afterBalance = alice.balance;
        assertEq(afterBalance, beforeBalance + 1 ether, "balance not increased");
    }

    function testFundsSentDirectlyToEOAAddressWithoutDelegation() public {
        uint256 beforeBalance = alice.balance;
        vm.prank(bob);
        (bool success,) = address(alice).call{value: 1 ether}("");
        assertTrue(success);
        uint256 afterBalance = alice.balance;
        assertEq(afterBalance, beforeBalance + 1 ether, "balance not increased");
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
