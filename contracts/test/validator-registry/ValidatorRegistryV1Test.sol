// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../contracts/validator-registry/ValidatorRegistryV1.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract ValidatorRegistryV1Test is Test {
    ValidatorRegistryV1 public validatorRegistry;
    address public owner;
    address public user1;
    address public user2;

    uint256 public constant MIN_STAKE = 1 ether;
    uint256 public constant SLASH_AMOUNT = 0.1 ether;
    uint256 public constant UNSTAKE_PERIOD = 10;
    address public constant SLASH_ORACLE = address(0x78888);
    address public constant SLASH_RECEIVER = address(0x78886);

    bytes public constant user1BLSKey = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public constant user2BLSKey = hex"a5c99dfdfc69791937ac5efc5d33316cd4e0698be24ef149bbc18f0f25ad92e5e11aafd39701dcdab6d3205ad38c307b";

    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    function setUp() public {
        owner = address(this);
        user1 = address(0x123);
        user2 = address(0x456);

        assertEq(user1BLSKey.length, 48);
        assertEq(user2BLSKey.length, 48);
        
        address proxy = Upgrades.deployUUPSProxy(
            "ValidatorRegistryV1.sol",
            abi.encodeCall(ValidatorRegistryV1.initialize, (MIN_STAKE, SLASH_AMOUNT, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD, owner))
        );
        validatorRegistry = ValidatorRegistryV1(payable(proxy));
    }

    function testSecondInitialize() public {
        vm.prank(owner);
        vm.expectRevert();
        validatorRegistry.initialize(MIN_STAKE, SLASH_AMOUNT, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD, owner);
        vm.stopPrank();
    }

    function testSelfStake() public {
        vm.deal(user1, 9 ether);
        assertEq(address(user1).balance, 9 ether);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);

        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);

        vm.stopPrank();

        assertEq(address(user1).balance, 8 ether);
        assertEq(validatorRegistry.withdrawalAddresses(user1BLSKey), user1);
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
        emit Staked(user1, user1, user1BLSKey, 1 ether);
        emit Staked(user1, user2, user2BLSKey, 1 ether);
        validatorRegistry.stake{value: totalAmount}(validators);
        vm.stopPrank();

        assertEq(user1.balance, 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 1 ether);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
    }

    function testDelegateStake() public {
        vm.deal(owner, 9 ether);
        assertEq(address(owner).balance, 9 ether);

        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        vm.startPrank(owner);

        vm.expectEmit(true, true, true, true);
        emit Staked(owner, user1, user1BLSKey, MIN_STAKE);
        vm.expectEmit(true, true, true, true);
        emit Staked(owner, user1, user2BLSKey, MIN_STAKE);
        validatorRegistry.delegateStake{value: 2*MIN_STAKE}(validators, user1); // Both validators are opted-in on user1's behalf

        vm.stopPrank();

        assertEq(address(owner).balance, 7 ether);
        
        assertEq(validatorRegistry.withdrawalAddresses(user1BLSKey), user1);
        assertEq(validatorRegistry.withdrawalAddresses(user2BLSKey), user1);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isStaked(user1BLSKey));
        assertTrue(validatorRegistry.isStaked(user2BLSKey));
    }

    function testUnstakeInsufficientFunds() public {
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);

        vm.startPrank(user2);
        vm.expectRevert("Validator must have staked balance");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);
    }

    function testUnauthorizedUnstake() public {
        testSelfStake();
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user2);
        vm.expectRevert("Only withdrawal address can call this function");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testWithdrawBeforeUnstake() public {
        testSelfStake();
        vm.roll(500);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectRevert("Validator must have unstaking balance");
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
    }

    function testAlreadyUnstaked() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert("Validator must have staked balance");
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Validator must have staked balance");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testStakeWhenAlreadyUnstaking() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
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
        emit StakeWithdrawn(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
        assertEq(address(user1).balance, 9 ether);

        vm.startPrank(user1);
        emit Staked(user1, user1, user1BLSKey, MIN_STAKE);
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
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
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
        emit StakeWithdrawn(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        assertFalse(validatorRegistry.isStaked(user1BLSKey), "User1 should not be considered staked after withdrawal");
        assertEq(address(user1).balance, 9 ether, "User1 should have all 9 ether after withdrawal");

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 0, "User1s staked balance should be 0 after withdrawal");
        assertTrue(validatorRegistry.withdrawalAddresses(user1BLSKey) == address(0), "User1s withdrawal address should be reset after withdrawal");
        assertTrue(validatorRegistry.unstakeBlockNums(user1BLSKey) == 0, "User1s unstake block number should be reset after withdrawal");
        assertTrue(validatorRegistry.getUnstakingAmount(user1BLSKey) == 0, "User1s unstaking balance should be reset after withdrawal");
    }
   
    function testGetBlocksTillWithdrawAllowed() public {
        testSelfStake();

        vm.expectRevert("Unstake must be initiated to check withdrawal eligibility");
        validatorRegistry.getBlocksTillWithdrawAllowed(user2BLSKey);

        assertEq(block.number, 1);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user1);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        uint256 blocksTillWithdraw = uint256(validatorRegistry.getBlocksTillWithdrawAllowed(user1BLSKey));
        assertEq(blocksTillWithdraw, 10);

        vm.roll(6);
        assertEq(block.number, 6);

        blocksTillWithdraw = uint256(validatorRegistry.getBlocksTillWithdrawAllowed(user1BLSKey));
        assertEq(blocksTillWithdraw, 5);

        vm.roll(10);
        assertEq(block.number, 10);

        blocksTillWithdraw = uint256(validatorRegistry.getBlocksTillWithdrawAllowed(user1BLSKey));
        assertEq(blocksTillWithdraw, 1);

        vm.startPrank(user1);
        vm.expectRevert("withdrawal not allowed yet. Blocks requirement not met.");
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        vm.roll(11);
        assertEq(block.number, 11);

        blocksTillWithdraw = uint256(validatorRegistry.getBlocksTillWithdrawAllowed(user1BLSKey));
        assertEq(blocksTillWithdraw, 0);

        vm.roll(17);
        blocksTillWithdraw = uint256(validatorRegistry.getBlocksTillWithdrawAllowed(user1BLSKey));
        assertEq(blocksTillWithdraw, 0);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
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
