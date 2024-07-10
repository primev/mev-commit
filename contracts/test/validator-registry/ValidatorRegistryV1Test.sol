// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

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
    event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance);
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    event MinStakeSet(address indexed owner, uint256 minStake);
    event SlashAmountSet(address indexed owner, uint256 slashAmount);
    event SlashOracleSet(address indexed owner, address slashOracle);
    event SlashReceiverSet(address indexed owner, address slashReceiver);
    event UnstakePeriodBlocksSet(address indexed owner, uint256 unstakePeriodBlocks);

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
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
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
        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1, user2BLSKey, 1 ether);
        validatorRegistry.stake{value: totalAmount}(validators);
        vm.stopPrank();

        assertEq(user1.balance, 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 1 ether);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
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
        
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));
    }

    function testAddStake() public {
        vm.deal(user1, 10 ether);
        assertEq(user1.balance, 10 ether);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        validatorRegistry.stake{value: MIN_STAKE/2}(validators);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE/2);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeAdded(user1, user1, user1BLSKey, MIN_STAKE/2, MIN_STAKE);
        validatorRegistry.addStake{value: MIN_STAKE/2}(validators);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    function testUnstakeInsufficientFunds() public {
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);

        vm.startPrank(user2);
        vm.expectRevert("Validator record must exist");
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
        vm.expectRevert("Unstake must be initiated before withdrawal");
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
    }

    function testAlreadyUnstaked() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert("Validator must NOT be unstaking");
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Validator must NOT be unstaking");
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testStakeWhenAlreadyUnstaking() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, block.number);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);

        vm.startPrank(user1);
        vm.expectRevert("Validator record must NOT exist");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert("Validator record must NOT exist");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.deal(user2, 10 ether);
        vm.startPrank(user2);
        vm.expectRevert("Validator record must NOT exist");
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        // Withdraw then try again
        assertEq(address(user1).balance, 8 ether);
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
        assertEq(address(user1).balance, 9 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
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

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, MIN_STAKE);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, block.number);
        assertEq(address(user1).balance, 8 ether);

        uint256 blockWaitPeriod = 11;
        vm.roll(block.number + blockWaitPeriod);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertEq(address(user1).balance, 9 ether, "User1 should have all 9 ether after withdrawal");

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0, "User1s staked balance should be 0 after withdrawal");
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, address(0), "User1s withdrawal address should be reset after withdrawal");
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 0, "User1s unstake block number should be reset after withdrawal");
    }

    function testSlashWithoutEnoughStake() public {
        vm.expectRevert("Validator record must exist");
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(SLASH_ORACLE);
        validatorRegistry.slash(validators);

        vm.deal(user1, 1 ether);
        vm.startPrank(user1);
        validatorRegistry.stake{value: MIN_STAKE/2}(validators);
        vm.stopPrank();

        vm.prank(owner);
        validatorRegistry.setSlashAmount(MIN_STAKE/2);

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, MIN_STAKE/2);
        validatorRegistry.slash(validators);

        vm.expectRevert("Validator balance must be greater than or equal to slash amount");
        vm.prank(SLASH_ORACLE);
        validatorRegistry.slash(validators);
    }

    function testUnauthorizedSlash() public {
        testSelfStake();

        vm.expectRevert("Only slashing oracle account can call this function");
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user2);
        validatorRegistry.slash(validators);
    }

    function testSlashingStakedValidator() public {
        testSelfStake();

        assertEq(address(user1).balance, 8 ether);
        assertEq(address(SLASH_RECEIVER).balance, 0);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 0);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.roll(11);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(SLASH_ORACLE, user1, user1BLSKey, 0.9 ether);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 0.1 ether);
        validatorRegistry.slash(validators);

        assertEq(address(user1).balance, 8.0 ether);
        assertEq(address(SLASH_RECEIVER).balance, 0.1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0.9 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 11);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    function testSlashingUnstakingValidator() public {
        testSelfStake();

        vm.roll(11);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        assertEq(address(user1).balance, 8 ether);
        assertEq(address(SLASH_RECEIVER).balance, 0);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 11);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.roll(22);

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 0.1 ether);
        validatorRegistry.slash(validators);

        finalAssertions(); // See directly below
    }

    // Split final assertions into own func to avoid stack overflow
    function finalAssertions() public view {
        assertEq(address(user1).balance, 8 ether);
        assertEq(address(SLASH_RECEIVER).balance, 0.1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0.9 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 22);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    function testBatchedSlashing() public {
        testMultiStake();
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 1 ether);

        vm.roll(14);

        bytes[] memory vals = new bytes[](1);
        vals[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, 1 ether);
        validatorRegistry.unstake(vals);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 14);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).unstakeHeight.blockHeight, 0);
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));

        vm.roll(78);

        bytes[] memory toSlash = new bytes[](2);
        toSlash[0] = user1BLSKey;
        toSlash[1] = user2BLSKey;
        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 0.1 ether);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user2BLSKey, 0.1 ether);
        validatorRegistry.slash(toSlash);

        assertEq(address(SLASH_RECEIVER).balance, 0.2 ether);

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0.9 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeHeight.blockHeight, 78);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 0.9 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).unstakeHeight.blockHeight, 78);
        assertFalse(validatorRegistry.isValidatorOptedIn(user2BLSKey));
    }
   
    function testGetBlocksTillWithdrawAllowed() public {
        testSelfStake();

        vm.expectRevert("Unstake must be initiated to check withdrawal eligibility");
        validatorRegistry.getBlocksTillWithdrawAllowed(user2BLSKey);

        assertEq(block.number, 1);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
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

    function testOwnerChangesSlashAmountAfterStaking() public {
        testSelfStake();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 1 ether);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.prank(owner);
        validatorRegistry.setMinStake(10 ether);
        assertEq(validatorRegistry.minStake(), 10 ether);

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 1 ether);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.deal(user1, 9 ether);
        vm.startPrank(user1);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        validatorRegistry.addStake{value: 9 ether}(validators);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 10 ether);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    function testOwnerChangesUnstakingPeriodWhileValIsUnstaking() public {
        testSelfStake();

        vm.roll(25);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(30);
        vm.prank(user1);
        vm.expectRevert("withdrawal not allowed yet. Blocks requirement not met.");
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        vm.prank(owner);
        validatorRegistry.setUnstakePeriodBlocks(3);
        assertEq(validatorRegistry.unstakePeriodBlocks(), 3);

        vm.prank(user1);
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

    function testOnlyOwnerCanMutateParams() public {

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.setMinStake(17 ether);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit MinStakeSet(owner, 17 ether);
        validatorRegistry.setMinStake(17 ether);
        vm.stopPrank();
        assertEq(validatorRegistry.minStake(), 17 ether);

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.setSlashAmount(0.2 ether);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit SlashAmountSet(owner, 0.2 ether);
        validatorRegistry.setSlashAmount(0.2 ether);
        vm.stopPrank();
        assertEq(validatorRegistry.slashAmount(), 0.2 ether);

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.setSlashOracle(user2);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit SlashOracleSet(owner, user2);
        validatorRegistry.setSlashOracle(user2);
        vm.stopPrank();
        assertEq(validatorRegistry.slashOracle(), user2);

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.setSlashReceiver(user2);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit SlashReceiverSet(owner, user2);
        validatorRegistry.setSlashReceiver(user2);
        vm.stopPrank();
        assertEq(validatorRegistry.slashReceiver(), user2);

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.setUnstakePeriodBlocks(20);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit UnstakePeriodBlocksSet(owner, 20);
        validatorRegistry.setUnstakePeriodBlocks(20);
        vm.stopPrank();
        assertEq(validatorRegistry.unstakePeriodBlocks(), 20);
    }

    function testPauseable() public {
        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.pause();
        vm.stopPrank();

        vm.startPrank(owner);
        validatorRegistry.pause();
        vm.stopPrank();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert();
        validatorRegistry.unpause();
        vm.stopPrank();

        vm.startPrank(owner);
        validatorRegistry.unpause();
        vm.stopPrank();

        vm.startPrank(user1);
        vm.deal(user1, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    // Tests the edge case where a user stakes, unstakes, then attempts to adds stake again
    function testAddStakeWhileUnstaking() public {
        testSelfStake();
        vm.roll(11);

        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);

        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isUnstaking(user1BLSKey));

        vm.prank(user1);
        vm.expectRevert("Validator must NOT be unstaking");
        validatorRegistry.addStake{value: MIN_STAKE}(validators);
    }
}
