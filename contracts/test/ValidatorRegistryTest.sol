// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "../contracts/ValidatorRegistry.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract ValidatorRegistryTest is Test {
    ValidatorRegistry public validatorRegistry;
    address public owner;
    address public user1;
    address public user2;

    uint256 public constant MIN_STAKE = 1 ether;
    uint256 public constant UNSTAKE_PERIOD = 10;

    event SelfStaked(address indexed staker, uint256 amount);
    event SplitStaked(address indexed staker, address[] recipients, uint256 totalAmount);
    event Unstaked(address indexed staker, uint256 amount);
    event StakeWithdrawn(address indexed staker, uint256 amount);

    function setUp() public {
        owner = address(this);
        user1 = address(0x123);
        user2 = address(0x456);
        
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorRegistry.sol",
            abi.encodeCall(ValidatorRegistry.initialize, (MIN_STAKE, UNSTAKE_PERIOD, owner))
        );
        validatorRegistry = ValidatorRegistry(proxy);
    }

    function testSecondInitialize() public {
        vm.prank(owner);
        vm.expectRevert();
        validatorRegistry.initialize(MIN_STAKE, UNSTAKE_PERIOD, owner);
        vm.stopPrank();
    }

    function testSelfStake() public {
        vm.deal(user1, 10 ether);
        assertEq(address(user1).balance, 10 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit SelfStaked(user1, MIN_STAKE);
        validatorRegistry.selfStake{value: MIN_STAKE}();
        vm.stopPrank();

        assertEq(address(user1).balance, 9 ether);
        assertEq(validatorRegistry.stakedBalances(user1), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1));
    }

    function testSplitStake() public {
        address[] memory recipients = new address[](2);
        recipients[0] = user1;
        recipients[1] = user2;

        uint256 totalAmount = 2 ether;
        vm.deal(address(this), 3 ether);
        assertEq(address(this).balance, 3 ether);

        vm.expectEmit(true, true, true, true);
        emit SplitStaked(address(this), recipients, totalAmount);
        validatorRegistry.splitStake{value: totalAmount}(recipients);

        assertEq(address(this).balance, 1 ether);
        assertEq(validatorRegistry.stakedBalances(user1), 1 ether);
        assertEq(validatorRegistry.stakedBalances(user2), 1 ether);
        assertTrue(validatorRegistry.isStaked(user1));
        assertTrue(validatorRegistry.isStaked(user2));
    }

    function testUnstakeInsufficientFunds() public {
        vm.startPrank(user2);
        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user2;

        assertEq(validatorRegistry.stakedBalances(user2), 0);
        vm.expectRevert("No balance to unstake");
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();
        assertEq(validatorRegistry.stakedBalances(user2), 0);
    }

    function testUnauthorizedUnstake() public {
        uint256 stakeAmount = 1 ether;
        vm.deal(user1, stakeAmount);

        vm.startPrank(user1);
        validatorRegistry.selfStake{value: stakeAmount}();
        vm.stopPrank();
        assertTrue(validatorRegistry.isStaked(user1));
        assertEq(validatorRegistry.stakedBalances(user1), stakeAmount);

        vm.startPrank(user2);
        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;
        vm.expectRevert("Not authorized to unstake. Must be stake originator or EOA whos staked");
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();
    }

    function testWithdrawBeforeUnstake() public {
        testSelfStake();
        vm.roll(500);

        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;

        vm.startPrank(user1);
        vm.expectRevert("Unstake must be initiated before withdrawal");
        validatorRegistry.withdraw(fromAddrs);
        vm.stopPrank();
    }

    function testAlreadyUnstaked() public {
        testSelfStake();

        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;

        vm.startPrank(user1);
        emit Unstaked(user1, MIN_STAKE);
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert("Unstake already initiated");
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Unstake already initiated");
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();
    }

    function testStakeWhenAlreadyUnstaking() public {
        testSelfStake();

        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;

        vm.startPrank(user1);
        emit Unstaked(user1, MIN_STAKE);
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();

        assertFalse(validatorRegistry.isStaked(user1));
        assertTrue(validatorRegistry.unstakeBlockNums(user1) > 0);
        assertTrue(validatorRegistry.stakedBalances(user1) == MIN_STAKE);

        vm.startPrank(user1);
        vm.expectRevert("Address cannot be staked with in-progress unstake process");
        validatorRegistry.selfStake{value: MIN_STAKE}();
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Address cannot be staked with in-progress unstake process");
        validatorRegistry.selfStake{value: MIN_STAKE}();
        vm.stopPrank();

        // Withdraw then try again
        vm.startPrank(user1);
        emit StakeWithdrawn(user1, MIN_STAKE);
        validatorRegistry.withdraw(fromAddrs);
        vm.stopPrank();

        vm.startPrank(user1);
        emit SelfStaked(user1, MIN_STAKE);
        validatorRegistry.selfStake{value: MIN_STAKE}();
        vm.stopPrank();
        assertTrue(validatorRegistry.isStaked(user1));
    }

    function testUnstakeWaitThenWithdraw() public {
        testSelfStake();

        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;

        assertEq(address(user1).balance, 9 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, MIN_STAKE);
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();

        // still has staked balance until withdrawal, but not considered "staked"
        assertEq(validatorRegistry.stakedBalances(user1), MIN_STAKE);
        assertFalse(validatorRegistry.isStaked(user1));
        assertEq(address(user1).balance, 9 ether);

        uint256 blockWaitPeriod = 11;
        vm.roll(block.number + blockWaitPeriod);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, MIN_STAKE);
        validatorRegistry.withdraw(fromAddrs);
        vm.stopPrank();

        assertFalse(validatorRegistry.isStaked(user1), "User1 should not be considered staked after withdrawal");
        assertEq(address(user1).balance, 10 ether, "User1 should have all 10 ether after withdrawal");

        assertEq(validatorRegistry.stakedBalances(user1), 0, "User1's staked balance should be 0 after withdrawal");
        assertTrue(validatorRegistry.stakeOriginators(user1) == address(0), "User1's stake originator should be reset after withdrawal");
        assertTrue(validatorRegistry.unstakeBlockNums(user1) == 0, "User1's unstake block number should be reset after withdrawal");
    }

    // To sanity check that relevant state for an account is reset s.t. they could stake again in future
    function testStakingCycle() public {
        testUnstakeWaitThenWithdraw();

        // Reset user1 balance for next cycle
        vm.prank(user1);
        (bool sent, ) = user2.call{value: 10 ether}("");
        require(sent, "Failed to send Ether");

        testUnstakeWaitThenWithdraw();
    }
}
