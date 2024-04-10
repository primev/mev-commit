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

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit SelfStaked(user1, MIN_STAKE);
        validatorRegistry.selfStake{value: MIN_STAKE}();
        vm.stopPrank();

        assertEq(validatorRegistry.stakedBalances(user1), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1));
    }

    function testSplitStake() public {
        address[] memory recipients = new address[](2);
        recipients[0] = user1;
        recipients[1] = user2;

        uint256 totalAmount = 2 ether;
        vm.deal(address(this), totalAmount);

        vm.expectEmit(true, true, true, true);
        emit SplitStaked(address(this), recipients, totalAmount);
        validatorRegistry.splitStake{value: totalAmount}(recipients);

        assertEq(validatorRegistry.stakedBalances(user1), 1 ether);
        assertEq(validatorRegistry.stakedBalances(user2), 1 ether);
        assertTrue(validatorRegistry.isStaked(user1));
        assertTrue(validatorRegistry.isStaked(user2));
    }

    function testFailUnstakeInsufficientFunds() public {
        vm.startPrank(user2);
        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user2;

        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();
    }

    function testUnauthorizedUnstake() public {
        uint256 stakeAmount = 1 ether;
        vm.deal(user1, stakeAmount);

        vm.startPrank(user1);
        validatorRegistry.selfStake{value: stakeAmount}();
        vm.stopPrank();
        assertTrue(validatorRegistry.isStaked(user1));

        vm.startPrank(user2);
        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;
        vm.expectRevert("Not authorized to unstake. Must be stake originator or EOA whos staked");
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();
    }

    function testUnstakeWaitThenWithdraw() public {
        testSelfStake();

        address[] memory fromAddrs = new address[](1);
        fromAddrs[0] = user1;

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, MIN_STAKE);
        validatorRegistry.unstake(fromAddrs);
        vm.stopPrank();

        // still has stake until withdrawal
        assertEq(validatorRegistry.stakedBalances(user1), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1));

        uint256 blockWaitPeriod = 11;
        vm.roll(block.number + blockWaitPeriod);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, MIN_STAKE);
        validatorRegistry.withdraw(fromAddrs);
        vm.stopPrank();

        assertEq(validatorRegistry.stakedBalances(user1), 0, "User1's staked balance should be 0 after withdrawal");
        assertFalse(validatorRegistry.isStaked(user1), "User1 should not be considered staked after withdrawal");
    }
}
