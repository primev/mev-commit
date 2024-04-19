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

    bytes public constant user1BLSKey = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public constant user2BLSKey = hex"a5c99dfdfc69791937ac5efc5d33316cd4e0698be24ef149bbc18f0f25ad92e5e11aafd39701dcdab6d3205ad38c307b";

    event Staked(address indexed staker, bytes valBLSPubKey, uint256 amount);
    event Unstaked(address indexed staker, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed staker, bytes valBLSPubKey, uint256 amount);

    function setUp() public {
        owner = address(this);

        user1 = address(0x123);
        user2 = address(0x456);

        assertEq(user1BLSKey.length, 48);
        assertEq(user2BLSKey.length, 48);
        
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
        vm.deal(user1, 9 ether);
        assertEq(address(user1).balance, 9 ether);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);

        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);

        vm.stopPrank();

        assertEq(address(user1).balance, 8 ether);
        assertEq(validatorRegistry.stakedBalances(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testMultiStake() public {
        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        uint256 totalAmount = 2 ether;
        vm.deal(address(this), 3 ether);
        assertEq(address(this).balance, 3 ether);

        vm.expectEmit(true, true, true, true);
        emit Staked(address(this), validators[0], 1 ether);
        emit Staked(address(this), validators[1], 1 ether);
        validatorRegistry.stake{value: totalAmount}(validators);

        assertEq(address(this).balance, 1 ether);
        assertEq(validatorRegistry.stakedBalances(user1BLSKey), 1 ether);
        assertEq(validatorRegistry.stakedBalances(user2BLSKey), 1 ether);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testUnstakeInsufficientFunds() public {
        vm.startPrank(user2);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;

        assertEq(validatorRegistry.stakedBalances(user2BLSKey), 0);
        vm.expectRevert("No balance to unstake");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
        assertEq(validatorRegistry.stakedBalances(user2BLSKey), 0);
    }

    function testUnauthorizedUnstake() public {
        testSelfStake();
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user2);
        vm.expectRevert("Not authorized to unstake validator. Must be stake originator");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testWithdrawBeforeUnstake() public {
        testSelfStake();
        vm.roll(500);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectRevert("Unstake must be initiated before withdrawal");
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
    }

    function testAlreadyUnstaked() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        emit Unstaked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert("Unstake already initiated for validator");
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Unstake already initiated for validator");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testStakeWhenAlreadyUnstaking() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        emit Unstaked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        assertFalse(validatorRegistry.isStaked(user1BLSKey));
        assertTrue(validatorRegistry.unstakeBlockNums(user1BLSKey) > 0);
        assertTrue(validatorRegistry.stakedBalances(user1BLSKey) == MIN_STAKE);

        vm.startPrank(user1);
        vm.expectRevert("validator cannot be staked with in-progress unstake process");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("validator cannot be staked with in-progress unstake process");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.deal(user2, 10 ether);
        vm.startPrank(user2);
        vm.expectRevert("validator cannot be staked with in-progress unstake process");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        // Withdraw then try again
        assertEq(address(user1).balance, 8 ether);
        vm.startPrank(user1);
        emit StakeWithdrawn(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
        assertEq(address(user1).balance, 9 ether);

        vm.startPrank(user1);
        emit Staked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testUnstakeWaitThenWithdraw() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        assertEq(address(user1).balance, 8 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        // still has staked balance until withdrawal, but not considered "staked"
        assertEq(validatorRegistry.stakedBalances(user1BLSKey), MIN_STAKE);
        assertFalse(validatorRegistry.isStaked(user1BLSKey));
        assertEq(address(user1).balance, 8 ether);

        uint256 blockWaitPeriod = 11;
        vm.roll(block.number + blockWaitPeriod);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        assertFalse(validatorRegistry.isStaked(user1BLSKey), "User1 should not be considered staked after withdrawal");
        assertEq(address(user1).balance, 9 ether, "User1 should have all 9 ether after withdrawal");

        assertEq(validatorRegistry.stakedBalances(user1BLSKey), 0, "User1s staked balance should be 0 after withdrawal");
        assertTrue(validatorRegistry.stakeOriginators(user1BLSKey) == address(0), "User1s stake originator should be reset after withdrawal");
        assertTrue(validatorRegistry.unstakeBlockNums(user1BLSKey) == 0, "User1s unstake block number should be reset after withdrawal");
    }

    // To sanity check that relevant state for an account is reset s.t. they could stake again in future
    function testStakingCycle() public {
        testUnstakeWaitThenWithdraw();

        // Reset user1 balance for next cycle
        vm.prank(user1);
        (bool sent, ) = user2.call{value: 9 ether}("");
        require(sent, "Failed to send Ether");

        testUnstakeWaitThenWithdraw();
    }
}
