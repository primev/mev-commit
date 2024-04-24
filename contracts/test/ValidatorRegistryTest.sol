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
        
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testMultiStake() public {
        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        uint256 totalAmount = 2 ether;
        vm.deal(user1, 3 ether);
        assertEq(user1.balance, 3 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1BLSKey, 1 ether);
        emit Staked(user1, user2BLSKey, 1 ether);
        validatorRegistry.stake{value: totalAmount}(validators);
        vm.stopPrank();

        assertEq(user1.balance, 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 1 ether);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testUnstakeInsufficientFunds() public {
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);

        vm.startPrank(user2);
        vm.expectRevert("Validator not staked");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);
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
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 0);
        assertEq(validatorRegistry.getUnstakingAmount(user1BLSKey), MIN_STAKE);

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

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 0);
        assertFalse(validatorRegistry.isStaked(user1BLSKey));
        assertEq(validatorRegistry.getUnstakingAmount(user1BLSKey), MIN_STAKE);
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

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 0, "User1s staked balance should be 0 after withdrawal");
        assertTrue(validatorRegistry.stakeOriginators(user1BLSKey) == address(0), "User1s stake originator should be reset after withdrawal");
        assertTrue(validatorRegistry.unstakeBlockNums(user1BLSKey) == 0, "User1s unstake block number should be reset after withdrawal");
        assertTrue(validatorRegistry.getUnstakingAmount(user1BLSKey) == 0, "User1s unstaking balance should be reset after withdrawal");
    }
    
    function testGetStakedValidators() public {
        testMultiStake();

        bytes[] memory validators = validatorRegistry.getStakedValidators(0, 2);
        assertEq(validators.length, 2);
        assertEq(validatorRegistry.getNumberOfStakedValidators(), 2);
        assertEq(validators[0], user1BLSKey);
        assertEq(validators[1], user2BLSKey);

        vm.deal(user1, 1000 ether);

        for (uint256 i = 0; i < 100; i++) {
            bytes memory key = new bytes(48);
            for (uint256 j = 0; j < 48; j++) {
                key[j] = bytes1(uint8(i));
            }
            bytes[] memory keys = new bytes[](1);
            keys[0] = key;
            vm.prank(user1);
            validatorRegistry.stake{value: MIN_STAKE}(keys);
            vm.stopPrank();
        }
        
        validators = validatorRegistry.getStakedValidators(0, 102);
        assertEq(validators.length, 102);
        assertEq(validatorRegistry.getNumberOfStakedValidators(), 102);

        assertEq(validators[0], user1BLSKey);
        assertEq(validators[1], user2BLSKey);

        validators = validatorRegistry.getStakedValidators(40, 60);
        assertEq(validators.length, 20);

        for (uint256 i = 0; i < 20; i++) {
            assertEq(validators[i].length, 48);
        }
    } 

    function testGetStakedValidatorsWithUnstakingInProgress() public {
        testMultiStake();

        uint256 numStakedValidators = validatorRegistry.getNumberOfStakedValidators();
        assertEq(numStakedValidators, 2);
        bytes[] memory validators = validatorRegistry.getStakedValidators(0, numStakedValidators);
        assertEq(validators.length, 2);

        bytes[] memory keys = new bytes[](1);
        keys[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(keys);
        assertTrue(validatorRegistry.unstakeBlockNums(user1BLSKey) > 0);
        vm.stopPrank();

        numStakedValidators = validatorRegistry.getNumberOfStakedValidators();
        assertEq(numStakedValidators, 1);
        validators = validatorRegistry.getStakedValidators(0, numStakedValidators);
        assertEq(validators.length, 1);

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1BLSKey, MIN_STAKE);
        keys = new bytes[](1);
        keys[0] = user1BLSKey;
        validatorRegistry.withdraw(keys);
        vm.stopPrank();
        numStakedValidators = validatorRegistry.getNumberOfStakedValidators();
        assertEq(numStakedValidators, 1);
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
