// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from"forge-std/Test.sol";
import {VanillaRegistry} from"../../contracts/validator-registry/VanillaRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IVanillaRegistry} from "../../contracts/interfaces/IVanillaRegistry.sol";
import {OwnableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract VanillaRegistryTest is Test {
    VanillaRegistry public validatorRegistry;
    address public owner;
    address public user1;
    address public user2;

    uint256 public constant MIN_STAKE = 1 ether;
    uint256 public constant UNSTAKE_PERIOD = 10;
    uint256 public constant PAYOUT_PERIOD = 20;
    address public constant SLASH_ORACLE = address(0x78888);
    address public constant SLASH_RECEIVER = address(0x78886);

    bytes public user1BLSKey = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public user2BLSKey = hex"a5c99dfdfc69791937ac5efc5d33316cd4e0698be24ef149bbc18f0f25ad92e5e11aafd39701dcdab6d3205ad38c307b";
    bytes public user3BLSKey = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";

    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance);
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount);
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    event MinStakeSet(address indexed owner, uint256 minStake);
    event SlashOracleSet(address indexed owner, address slashOracle);
    event SlashReceiverSet(address indexed owner, address slashReceiver);
    event UnstakePeriodBlocksSet(address indexed owner, uint256 unstakePeriodBlocks);

    event FeeTransfer(uint256 amount, address indexed recipient);

    event WithdrawalAddressBlacklisted(address indexed msgSender, address indexed withdrawalAddress);
    event WithdrawalAddressUnblacklisted(address indexed msgSender, address indexed withdrawalAddress);

    function setUp() public {
        owner = address(this);
        user1 = address(0x123);
        user2 = address(0x456);

        assertEq(user1BLSKey.length, 48);
        assertEq(user2BLSKey.length, 48);
        
        address proxy = Upgrades.deployUUPSProxy(
            "VanillaRegistry.sol",
            abi.encodeCall(VanillaRegistry.initialize, (MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD, PAYOUT_PERIOD, owner))
        );
        validatorRegistry = VanillaRegistry(payable(proxy));
    }

    function testSecondInitialize() public {
        vm.prank(owner);
        vm.expectRevert();
        validatorRegistry.initialize(MIN_STAKE, SLASH_ORACLE, SLASH_RECEIVER, UNSTAKE_PERIOD, PAYOUT_PERIOD, owner);
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

        uint256 totalAmount = 6 ether;
        vm.deal(user1, 7 ether);
        assertEq(user1.balance, 7 ether);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1, user1BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit Staked(user1, user1, user2BLSKey, 3 ether);
        validatorRegistry.stake{value: totalAmount}(validators);
        vm.stopPrank();

        assertEq(user1.balance, 1 ether);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 3 ether);
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 3 ether);
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

        vm.deal(user2, 10 ether);
        assertEq(user2.balance, 10 ether);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.StakeTooLowForNumberOfKeys.selector, MIN_STAKE/2, MIN_STAKE));
        validatorRegistry.stake{value: MIN_STAKE/2}(validators);

        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.prank(user1);
        validatorRegistry.stake{value: MIN_STAKE}(validators);

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.StakeTooLowForNumberOfKeys.selector, 0, 1));
        validatorRegistry.addStake{value: 0}(validators);

        vm.prank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector, user1, user2));
        validatorRegistry.addStake{value: MIN_STAKE/2}(validators);

        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit StakeAdded(user1, user1, user1BLSKey, MIN_STAKE/2, 3*MIN_STAKE/2);
        validatorRegistry.addStake{value: MIN_STAKE/2}(validators);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), 3*MIN_STAKE/2);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    function testUnstakeInsufficientFunds() public {
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);

        vm.startPrank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustExist.selector, user2BLSKey));
        validatorRegistry.unstake(validators);
        vm.stopPrank();
        assertEq(validatorRegistry.getStakedAmount(user2BLSKey), 0);
    }

    function testUnauthorizedUnstake() public {
        testSelfStake();
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.startPrank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector,
        user1, // actual from validator record
        user2)); // expected from msg.sender
        validatorRegistry.unstake(validators);
        vm.stopPrank();
    }

    function testUnathorizedMultiUnstake() public {
        testSelfStake();
        bytes[] memory validators = new bytes[](1);
        validators[0] = user2BLSKey;
        vm.deal(user2, MIN_STAKE);
        vm.prank(user2);
        validatorRegistry.stake{value: MIN_STAKE}(validators);

        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user2);

        bytes[] memory bothValidators = new bytes[](2);
        bothValidators[0] = user1BLSKey;
        bothValidators[1] = user2BLSKey;
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector,
        user2, // actual from (second) validator record
        user1)); // expected from msg.sender
        validatorRegistry.unstake(bothValidators);
    }

    function testUnauthorizedWithdraw() public {
        testSelfStake();

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user1);
        validatorRegistry.unstake(validators);

        vm.prank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector, 
        user1, // actual from validator record
        user2)); // expected from msg.sender
        validatorRegistry.withdraw(validators);
    }

    function testUnathorizedMultiWithdraw() public {
        testUnathorizedMultiUnstake(); // Use setup where two validators are staked from different withdrawal addresses

        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user2);


        bytes[] memory val1 = new bytes[](1);
        val1[0] = user1BLSKey;
        vm.prank(user1);
        validatorRegistry.unstake(val1);

        bytes[] memory val2 = new bytes[](1);
        val2[0] = user2BLSKey;
        vm.prank(user2);
        validatorRegistry.unstake(val2);

        assertTrue(validatorRegistry.isUnstaking(user1BLSKey));
        assertTrue(validatorRegistry.isUnstaking(user2BLSKey));

        vm.roll(2000);

        bytes[] memory bothValidators = new bytes[](2);
        bothValidators[0] = user1BLSKey;
        bothValidators[1] = user2BLSKey;
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector,
        user2, // actual from validator record
        user1)); // expected from msg.sender
        validatorRegistry.withdraw(bothValidators);
    }

    function testWithdrawBeforeUnstake() public {
        testSelfStake();
        vm.roll(500);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.startPrank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.MustUnstakeToWithdraw.selector));
        validatorRegistry.withdraw(validators);
        vm.stopPrank();
    }

    function testMultiWithdraw() public {
        testMultiStake();

        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        vm.startPrank(user1);
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        assertEq(address(user1).balance, 1 ether);

        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1, user1BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(user1, user1, user2BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit TotalStakeWithdrawn(user1, user1, 6 ether);
        validatorRegistry.withdraw(validators);
        vm.stopPrank();

        assertEq(address(user1).balance, 7 ether);
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
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorCannotBeUnstaking.selector, user1BLSKey));
        validatorRegistry.unstake(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorCannotBeUnstaking.selector, user1BLSKey));
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
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, block.number);
        assertEq(validatorRegistry.getStakedAmount(user1BLSKey), MIN_STAKE);

        vm.startPrank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustNotExist.selector, user1BLSKey));
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.roll(500);

        vm.startPrank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustNotExist.selector, user1BLSKey));
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        vm.deal(user2, 10 ether);
        vm.startPrank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustNotExist.selector, user1BLSKey));
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
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, block.number);
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
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 0, "User1s unstake block number should be reset after withdrawal");
    }

    function testSlashWithoutEnoughStake() public {
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustExist.selector, user1BLSKey));
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(SLASH_ORACLE);
        validatorRegistry.slash(validators, true);

        vm.deal(user1, 2 ether);
        vm.startPrank(user1);
        uint256 stakeAmount = MIN_STAKE+1;
        validatorRegistry.stake{value: stakeAmount}(validators);
        vm.stopPrank();

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.slash(validators, true);

        vm.expectRevert(IVanillaRegistry.NotEnoughBalanceToSlash.selector);
        vm.prank(SLASH_ORACLE);
        validatorRegistry.slash(validators, true);
    }

    function testUnauthorizedSlash() public {
        testSelfStake();

        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.SenderIsNotSlashOracle.selector, user2, SLASH_ORACLE));
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user2);
        validatorRegistry.slash(validators, true);
    }

    function testSlashingStakedValidator() public {
        testSelfStake();

        assertEq(address(user1).balance, 8 ether);
        assertEq(address(SLASH_RECEIVER).balance, 0);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 1 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 0);
        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.roll(11);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(SLASH_ORACLE, user1, user1BLSKey, 1 ether);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 1 ether);
        validatorRegistry.slash(validators, true);

        assertEq(address(user1).balance, 8.0 ether);

        assertEq(address(SLASH_RECEIVER).balance, 0 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 1 ether);
        assertFalse(validatorRegistry.isSlashingPayoutDue());

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 11);
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
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 11);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        vm.roll(22);

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 1 ether);
        validatorRegistry.slash(validators, false);

        finalAssertions(); // See directly below
    }

    // Split final assertions into own func to avoid stack overflow
    function finalAssertions() public view {
        assertEq(address(user1).balance, 8 ether);

        assertEq(address(SLASH_RECEIVER).balance, 0 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 1 ether);
        assertTrue(validatorRegistry.isSlashingPayoutDue());

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 0 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        // Unstake occurrence should not be updated for already unstaked validators
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 11);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
    }

    // solhint-disable-next-line ordering
    function testBatchedSlashing() public {
        testMultiStake();
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 3 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 3 ether);

        vm.roll(14);

        bytes[] memory vals = new bytes[](1);
        vals[0] = user1BLSKey;
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, 3 ether);
        validatorRegistry.unstake(vals);
        vm.stopPrank();

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 3 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 14);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 3 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).unstakeOccurrence.blockHeight, 0);
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));

        vm.roll(78);

        assertTrue(validatorRegistry.isSlashingPayoutDue());
        assertEq(address(SLASH_RECEIVER).balance, 0 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 0 ether);

        bytes[] memory toSlash = new bytes[](2);
        toSlash[0] = user1BLSKey;
        toSlash[1] = user2BLSKey;
        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user1BLSKey, 1 ether);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user1, user2BLSKey, 1 ether);
        validatorRegistry.slash(toSlash, true);

        assertFalse(validatorRegistry.isSlashingPayoutDue());
        assertEq(address(SLASH_RECEIVER).balance, 2 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 0 ether);

        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).balance, 2 ether);
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).withdrawalAddress, user1);
        // Unstake occurrence should not be updated for already unstaked validators
        assertEq(validatorRegistry.getStakedValidator(user1BLSKey).unstakeOccurrence.blockHeight, 14);
        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));

        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).balance, 2 ether);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).withdrawalAddress, user1);
        assertEq(validatorRegistry.getStakedValidator(user2BLSKey).unstakeOccurrence.blockHeight, 78);
        assertFalse(validatorRegistry.isValidatorOptedIn(user2BLSKey));
    }

    function testManualPayout() public { 
        testBatchedSlashing();

        vm.roll(10000);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user3BLSKey;
        address user3 = vm.addr(0x23333);
        vm.deal(user3, 10 ether);
        vm.startPrank(user3);
        vm.expectEmit(true, true, true, true);
        emit Staked(user3, user3, user3BLSKey, MIN_STAKE);
        validatorRegistry.stake{value: MIN_STAKE}(validators);
        vm.stopPrank();

        assertTrue(validatorRegistry.isSlashingPayoutDue());
        assertEq(address(SLASH_RECEIVER).balance, 2 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 0 ether);

        vm.prank(SLASH_ORACLE);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(SLASH_ORACLE, user3, user3BLSKey, MIN_STAKE);
        vm.expectEmit(true, true, true, true);
        emit Slashed(SLASH_ORACLE, SLASH_RECEIVER, user3, user3BLSKey, MIN_STAKE);
        validatorRegistry.slash(validators, false);

        assertTrue(validatorRegistry.isSlashingPayoutDue());
        assertEq(address(SLASH_RECEIVER).balance, 2 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 1 ether);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit FeeTransfer(1 ether, SLASH_RECEIVER);
        validatorRegistry.manuallyTransferSlashingFunds();

        assertFalse(validatorRegistry.isSlashingPayoutDue());
        assertEq(address(SLASH_RECEIVER).balance, 3 ether);
        assertEq(validatorRegistry.getAccumulatedSlashingFunds(), 0 ether);
    }

    function testGetBlocksTillWithdrawAllowed() public {
        testSelfStake();

        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.MustUnstakeToWithdraw.selector));
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
        vm.expectRevert(IVanillaRegistry.WithdrawingTooSoon.selector);
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
        vm.expectRevert(IVanillaRegistry.WithdrawingTooSoon.selector);
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
        (address recipient,,,) = validatorRegistry.slashingFundsTracker();
        assertEq(recipient, user2);

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
        vm.expectRevert(
            abi.encodeWithSelector(IVanillaRegistry.ValidatorCannotBeUnstaking.selector, user1BLSKey)
        );
        validatorRegistry.addStake{value: MIN_STAKE}(validators);
    }

    function testPrecisionLossPrevention() public {
        vm.prank(owner);
        validatorRegistry.setMinStake(1 wei);

        bytes[] memory validators = new bytes[](90);
        for (uint256 i = 0; i < 90; ++i) {
            validators[i] = user1BLSKey;
            validators[i][0] = bytes1(uint8(i + 1));
        }
        vm.deal(user1, 100 wei);

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.StakeTooLowForNumberOfKeys.selector, 80, 90));
        validatorRegistry.stake{value: 80 wei}(validators);

        assertEq(address(validatorRegistry).balance, 0);

        vm.prank(user1);
        validatorRegistry.stake{value: 100 wei}(validators);

        for (uint256 i = 0; i < 89; ++i) {
            assertEq(validatorRegistry.getStakedAmount(validators[i]), 1 wei);
            assertTrue(validatorRegistry.isValidatorOptedIn(validators[i]));
        }
        uint256 expectedFinalStake = 100 wei - (89 * 1 wei);
        assertEq(validatorRegistry.getStakedAmount(validators[89]), expectedFinalStake);
        assertTrue(validatorRegistry.isValidatorOptedIn(validators[89]));

        assertEq(user1.balance, 0);
        assertEq(address(validatorRegistry).balance, 100 wei);

        vm.deal(user1, 100 wei);

        vm.prank(user1);
        validatorRegistry.addStake{value: 100 wei}(validators);

        for (uint256 i = 0; i < 89; ++i) {
            assertEq(validatorRegistry.getStakedAmount(validators[i]), 2 wei);
        }
        expectedFinalStake = 2 * expectedFinalStake;
        assertEq(validatorRegistry.getStakedAmount(validators[89]), expectedFinalStake);

        assertEq(user1.balance, 0);
        assertEq(address(validatorRegistry).balance, 200 wei);
    }

    function testBlacklistWithdrawalAddresses() public { 
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user1));
        validatorRegistry.blacklistWithdrawalAddresses(new address[](0));

        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = address(0);
        vm.expectRevert(IVanillaRegistry.WithdrawalAddressMustBeSet.selector);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        withdrawalAddresses[0] = owner;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.OwnerCantBlacklistSelf.selector, owner));
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        withdrawalAddresses[0] = user2;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user2);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AlreadyBlacklisted.selector, user2));
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);
    }

    function testUnblacklistWithdrawalAddresses() public {
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user1));
        validatorRegistry.blacklistWithdrawalAddresses(new address[](0));

        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = user2;
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.NotBlacklisted.selector, user2));
        validatorRegistry.unblacklistWithdrawalAddresses(withdrawalAddresses);

        testBlacklistWithdrawalAddresses();

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressUnblacklisted(owner, user2);
        validatorRegistry.unblacklistWithdrawalAddresses(withdrawalAddresses);
    }

    function testBlacklistedAddrCantStake() public {
        testBlacklistWithdrawalAddresses();

        vm.deal(user2, MIN_STAKE);
        vm.prank(user2);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AddressIsBlacklisted.selector, user2));
        validatorRegistry.stake{value: MIN_STAKE}(new bytes[](0));
    }

    function testOwnerCantDelegateStakeToBlacklistedAddr() public {
        testBlacklistWithdrawalAddresses();

        vm.deal(owner, MIN_STAKE);
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AddressIsBlacklisted.selector, user2));
        validatorRegistry.delegateStake(new bytes[](0), user2);
    }

    function testBlacklistedAddrCantAddStake() public {
        testSelfStake();

        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = user1;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        vm.deal(user1, 2*MIN_STAKE);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AddressIsBlacklisted.selector, user1));
        validatorRegistry.addStake(validators);
    }

    function testBlacklistedAddrCantUnstake() public {
        testSelfStake();

        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = user1;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AddressIsBlacklisted.selector, user1));
        validatorRegistry.unstake(validators);
    }

    function testBlacklistedAddrCantWithdraw() public {
        testSelfStake();

        vm.prank(user1);
        bytes[] memory validators = new bytes[](1);
        validators[0] = user1BLSKey;
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, MIN_STAKE);
        validatorRegistry.unstake(validators);

        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = user1;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.AddressIsBlacklisted.selector, user1));
        validatorRegistry.withdraw(validators);
    }

    function testBlacklistingWithRegisteredValidators() public {
        testMultiStake();

        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));

        vm.prank(owner);
        address[] memory withdrawalAddresses = new address[](1);
        withdrawalAddresses[0] = user1;
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(withdrawalAddresses);

        assertFalse(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertFalse(validatorRegistry.isValidatorOptedIn(user2BLSKey));

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressUnblacklisted(owner, user1);
        validatorRegistry.unblacklistWithdrawalAddresses(withdrawalAddresses);

        assertTrue(validatorRegistry.isValidatorOptedIn(user1BLSKey));
        assertTrue(validatorRegistry.isValidatorOptedIn(user2BLSKey));
    }

    function testUnstakeViaBlacklist() public { 
        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustExist.selector, user1BLSKey));
        validatorRegistry.unstakeViaBlacklist(validators, user1);

        testMultiStake();

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user1));
        validatorRegistry.unstakeViaBlacklist(validators, user1);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.NotBlacklisted.selector, user1));
        validatorRegistry.unstakeViaBlacklist(validators, user1);

        address[] memory toBlacklist = new address[](2);
        toBlacklist[0] = user1;
        toBlacklist[1] = user2;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(toBlacklist);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawalAddressMismatch.selector, user1, user2));
        validatorRegistry.unstakeViaBlacklist(validators, user2);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(owner, user1, user1BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(owner, user1, user2BLSKey, 3 ether);
        validatorRegistry.unstakeViaBlacklist(validators, user1);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorCannotBeUnstaking.selector, user1BLSKey));
        validatorRegistry.unstakeViaBlacklist(validators, user1);
    }

    function testWithdrawViaBlacklist() public { 
        bytes[] memory validators = new bytes[](2);
        validators[0] = user1BLSKey;
        validators[1] = user2BLSKey;

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustExist.selector, user1BLSKey));
        validatorRegistry.withdrawViaBlacklist(validators, user1);

        testMultiStake();
        vm.prank(user1);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user1BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit Unstaked(user1, user1, user2BLSKey, 3 ether);
        validatorRegistry.unstake(validators);

        vm.prank(user1);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user1));
        validatorRegistry.withdrawViaBlacklist(validators, user1);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.NotBlacklisted.selector, user1));
        validatorRegistry.withdrawViaBlacklist(validators, user1);

        address[] memory toBlacklist = new address[](2);
        toBlacklist[0] = user1;
        toBlacklist[1] = user2;
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit WithdrawalAddressBlacklisted(owner, user1);
        validatorRegistry.blacklistWithdrawalAddresses(toBlacklist);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.WithdrawingTooSoon.selector));
        validatorRegistry.withdrawViaBlacklist(validators, user1);

        vm.roll(200);

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(owner, user1, user1BLSKey, 3 ether);
        vm.expectEmit(true, true, true, true);
        emit StakeWithdrawn(owner, user1, user2BLSKey, 3 ether);
        validatorRegistry.withdrawViaBlacklist(validators, user1);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IVanillaRegistry.ValidatorRecordMustExist.selector, user1BLSKey));
        validatorRegistry.withdrawViaBlacklist(validators, user1);
    }

    function testGriefScenarioAlleviated() public {
        // TODO: use withdrawl address change func
    }

    function testIsValidatorOptedInInvalidInput() public view {
        bytes memory invalidBLSPubKey = hex"0000000004";
        assertFalse(validatorRegistry.isValidatorOptedIn(invalidBLSPubKey));
    }
}
