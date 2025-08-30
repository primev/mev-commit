// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import {DepositManager} from "../../contracts/core/DepositManager.sol";
import {IBidderRegistry} from "../../contracts/interfaces/IBidderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {Errors} from "../../contracts/utils/Errors.sol";

contract DepositManagerTest is Test {
    uint256 public constant ALICE_PK = uint256(0xA11CE);
    uint256 public constant BOB_PK = uint256(0xB0B);
    address public alice = payable(vm.addr(ALICE_PK));
    address public bob = payable(vm.addr(BOB_PK));

    DepositManager private depositManagerImpl;

    function setUp() public {
        depositManagerImpl = new DepositManager(address(new BidderRegistry()), 0.01 ether);
        vm.deal(alice, 10 ether);
        vm.deal(bob, 10 ether);
        vm.signAndAttachDelegation(address(depositManagerImpl), ALICE_PK);
        vm.signAndAttachDelegation(address(depositManagerImpl), BOB_PK);
    }

    function testCodeAtAddress() public {
        bytes32 expectedCodehash = keccak256(abi.encodePacked(hex"ef0100", address(depositManagerImpl)));
        assertEq(alice.codehash, expectedCodehash);
        assertEq(alice.code.length, 23);
        uint256 otherUserPk = uint256(0x1237777777);
        address otherUser = payable(vm.addr(otherUserPk));
        assertEq(otherUser.codehash, 0x0000000000000000000000000000000000000000000000000000000000000000);
        assertEq(otherUser.code.length, 0);
        vm.signAndAttachDelegation(address(depositManagerImpl), otherUserPk);
        assertEq(otherUser.codehash, expectedCodehash);
        assertEq(otherUser.code.length, 23);
    }

    function testFallbackRevert() public {
        bytes memory badData = abi.encodeWithSelector(
            DepositManager.setTargetDeposit.selector,
            address(0x55555),
            1 ether,
            1 ether,
            1 ether
        );
        vm.prank(alice);
        (bool success, ) = address(depositManagerImpl).call{value: 1 ether}(badData);
        assertFalse(success);
    }

    function testReceiveNoRevert() public {
        uint256 beforeBalance = alice.balance;
        vm.prank(bob);
        (bool success, ) = address(alice).call{value: 1 ether}("");
        assertTrue(success);
        uint256 afterBalance = alice.balance;
        assertEq(afterBalance, beforeBalance + 1 ether, "balance not increased");
    }

    function testSetTargetDeposit() public {
        address provider = vm.addr(0x55555);

        uint256 aliceTarget = 1 ether;
        uint256 bobTarget = 2 ether;

        vm.expectEmit(true, true, true, true);
        emit DepositManager.TargetDepositSet(provider, aliceTarget);
        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(provider, aliceTarget);

        uint256 aliceStored = DepositManager(payable(alice)).targetDeposits(provider);
        assertEq(aliceStored, aliceTarget, "target deposit not set for alice");
        uint256 bobStored = DepositManager(payable(bob)).targetDeposits(provider);
        assertEq(bobStored, 0, "target deposit not set for bob");

        vm.expectEmit(true, true, true, true);
        emit DepositManager.TargetDepositSet(provider, bobTarget);
        vm.prank(bob);
        DepositManager(payable(bob)).setTargetDeposit(provider, bobTarget);

        aliceStored = DepositManager(payable(alice)).targetDeposits(provider);
        assertEq(aliceStored, aliceTarget, "target deposit changed for alice");
        bobStored = DepositManager(payable(bob)).targetDeposits(provider);
        assertEq(bobStored, bobTarget, "target deposit not set for bob");

        uint256 newTarget = 3 ether;
        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, bob, alice));
        vm.prank(bob);
        DepositManager(payable(alice)).setTargetDeposit(provider, newTarget);

        vm.expectEmit(true, true, true, true);
        emit DepositManager.TargetDepositSet(provider, newTarget);
        vm.prank(bob);
        DepositManager(payable(bob)).setTargetDeposit(provider, newTarget);
        bobStored = DepositManager(payable(bob)).targetDeposits(provider);
        assertEq(bobStored, newTarget, "target deposit not set for bob");

        uint256 newTargetForNewProvider = 4 ether;
        address newProvider = vm.addr(0x123);

        vm.expectEmit(true, true, true, true);
        emit DepositManager.TargetDepositSet(newProvider, newTargetForNewProvider);
        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(newProvider, newTargetForNewProvider);

        uint256 aliceStoredForNewProvider = DepositManager(payable(alice)).targetDeposits(newProvider);
        assertEq(aliceStoredForNewProvider, newTargetForNewProvider, "target deposit not set for alice");
        aliceStored = DepositManager(payable(alice)).targetDeposits(provider);
        assertEq(aliceStored, aliceTarget, "target deposit changed for alice");
    }

    function testSetTargetDeposit_RevertNotThisEOA() public {
        address provider = vm.addr(0x55555);
        uint256 target = 7 ether;

        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, provider, alice));
        vm.prank(provider);
        DepositManager(payable(alice)).setTargetDeposit(provider, target);

        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, alice, bob));
        vm.prank(alice);
        DepositManager(payable(bob)).setTargetDeposit(provider, target);

        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, bob, alice));
        vm.prank(bob);
        DepositManager(payable(alice)).setTargetDeposit(provider, target);

        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, address(depositManagerImpl), alice));
        vm.prank(address(depositManagerImpl));
        DepositManager(payable(alice)).setTargetDeposit(provider, target);

        vm.expectRevert(abi.encodeWithSelector(DepositManager.NotThisEOA.selector, alice, address(depositManagerImpl)));
        vm.prank(alice);
        DepositManager(payable(address(depositManagerImpl))).setTargetDeposit(provider, target);
    }
}
