// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {VanillaRegistry} from "../../../contracts/validator-registry/VanillaRegistry.sol";
import {ValidatorOptInRouter} from "../../../contracts/validator-registry/ValidatorOptInRouter.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {VanillaRegistryTest} from "../VanillaRegistryTest.sol";
import {MevCommitAVSTest} from "../avs/MevCommitAVSTest.sol";
import {IValidatorOptInRouter} from "../../../contracts/interfaces/IValidatorOptInRouter.sol";
import {IVanillaRegistry} from "../../../contracts/interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../../../contracts/interfaces/IMevCommitAVS.sol";
import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {MevCommitMiddlewareTestCont} from "../middleware/MevCommitMiddlewareTestCont.sol";
import {RewardManager} from "../../../contracts/validator-registry/rewards/RewardManager.sol";
import {IRewardManager} from "../../../contracts/interfaces/IRewardManager.sol";
import {PausableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/utils/PausableUpgradeable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

contract RewardManagerTest is Test {
    RewardManager public rewardManager;

    VanillaRegistry public vanillaRegistry;
    VanillaRegistryTest public vanillaRegistryTest;
    MevCommitAVS public mevCommitAVS;
    MevCommitAVSTest public mevCommitAVSTest;
    MevCommitMiddleware public mevCommitMiddleware;
    MevCommitMiddlewareTestCont public mevCommitMiddlewareTest;

    address public owner;
    address public user1;
    address public user2;
    address public user3;
    address public user4;
    address public user5;

    bytes public sampleValPubkey1 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254b3f";
    bytes public sampleValPubkey2 = hex"aca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes public sampleValPubkey3 = hex"cca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes public sampleValPubkey4 = hex"dca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes public sampleValPubkey5 = hex"eca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";

    event VanillaRegistrySet(address indexed newVanillaRegistry);
    event MevCommitAVSSet(address indexed newMevCommitAVS);
    event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware);
    event AutoClaimGasLimitSet(uint256 autoClaimGasLimit);
    event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount);
    event OrphanedRewardsClaimed(address indexed toPay, uint256 amount);
    event RemovedFromAutoClaimBlacklist(address indexed addr);
    event AutoClaimEnabled(address indexed caller);
    event AutoClaimDisabled(address indexed caller);
    event PaymentStored(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount);
    event RewardsClaimed(address indexed toPay, uint256 amount);
    event OverrideAddressSet(address indexed receiver, address indexed overrideAddress);
    event OverrideAddressRemoved(address indexed receiver);
    event AutoClaimed(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount);
    event AutoClaimTransferFailed(address indexed provider, address indexed receiver, address indexed toPay);

    function setUp() public {
        owner = address(0x123456);
        user1 = address(0x123);
        user2 = address(0x456);
        user3 = address(0x789);
        user4 = address(0xabc);
        user5 = address(0xdef);

        vanillaRegistryTest = new VanillaRegistryTest();
        vanillaRegistryTest.setUp();
        vanillaRegistry = vanillaRegistryTest.validatorRegistry();

        mevCommitAVSTest = new MevCommitAVSTest();
        mevCommitAVSTest.setUp();
        mevCommitAVS = mevCommitAVSTest.mevCommitAVS();

        mevCommitMiddlewareTest = new MevCommitMiddlewareTestCont();
        mevCommitMiddlewareTest.setUp();
        mevCommitMiddleware = mevCommitMiddlewareTest.mevCommitMiddleware();

        uint256 autoClaimGasLimit = 50000;

        address rewardManagerProxy = Upgrades.deployUUPSProxy(
            "RewardManager.sol",
            abi.encodeCall(RewardManager.initialize,
            (
                address(vanillaRegistry),
                address(mevCommitAVS),
                address(mevCommitMiddleware),
                autoClaimGasLimit,
                owner
            ))
        );
        rewardManager = RewardManager(payable(rewardManagerProxy));
    }

    function testRMSetters() public {
        IVanillaRegistry newRegistry = new VanillaRegistry();
        vm.prank(owner);
        vm.expectEmit();
        emit VanillaRegistrySet(address(newRegistry));
        rewardManager.setVanillaRegistry(address(newRegistry));

        IMevCommitAVS newAVS = new MevCommitAVS();
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitAVSSet(address(newAVS));
        rewardManager.setMevCommitAVS(address(newAVS));

        IMevCommitMiddleware newMiddleware = new MevCommitMiddleware();
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitMiddlewareSet(address(newMiddleware));
        rewardManager.setMevCommitMiddleware(address(newMiddleware));

        uint256 newAutoClaimGasLimit = 79000;
        vm.prank(owner);
        vm.expectEmit();
        emit AutoClaimGasLimitSet(newAutoClaimGasLimit);
        rewardManager.setAutoClaimGasLimit(newAutoClaimGasLimit);
        assertEq(rewardManager.autoClaimGasLimit(), 79000);
    }

    function testRMPause() public {

        assertEq(rewardManager.paused(), false);
        vm.prank(user1);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, user1)
        );
        rewardManager.pause();
        assertEq(rewardManager.paused(), false);

        vm.prank(owner);
        rewardManager.pause();
        assertEq(rewardManager.paused(), true);

        // payProposer should still work when paused
        vm.deal(user1, 1 ether);
        vm.prank(user1);
        vm.expectEmit();
        emit OrphanedRewardsAccumulated(user1, sampleValPubkey1, 1 ether);
        rewardManager.payProposer{value: 1 ether}(sampleValPubkey1);

        // User functions should not work when paused
        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.enableAutoClaim(true);

        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.disableAutoClaim();

        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.overrideReceiver(user2, false);

        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.removeOverrideAddress();

        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.claimRewards();

        // Owner functions should still work
        vm.prank(owner);
        bytes[] memory pubkeys = new bytes[](1);
        pubkeys[0] = sampleValPubkey1;
        vm.expectEmit();
        emit OrphanedRewardsClaimed(user2, 1 ether);
        rewardManager.claimOrphanedRewards(pubkeys, user2);

        vm.prank(owner);
        vm.expectEmit();
        emit RemovedFromAutoClaimBlacklist(user1);
        rewardManager.removeFromAutoClaimBlacklist(user1);

        VanillaRegistry newRegistry = new VanillaRegistry();
        vm.prank(owner);
        vm.expectEmit();
        emit VanillaRegistrySet(address(newRegistry));
        rewardManager.setVanillaRegistry(address(newRegistry));

        IMevCommitAVS newAVS = new MevCommitAVS();
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitAVSSet(address(newAVS));
        rewardManager.setMevCommitAVS(address(newAVS));

        IMevCommitMiddleware newMiddleware = new MevCommitMiddleware();
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitMiddlewareSet(address(newMiddleware));
        rewardManager.setMevCommitMiddleware(address(newMiddleware));
        
        uint256 newAutoClaimGasLimit = 79000;
        vm.prank(owner);
        vm.expectEmit();
        emit AutoClaimGasLimitSet(newAutoClaimGasLimit);
        rewardManager.setAutoClaimGasLimit(newAutoClaimGasLimit);

        assertEq(rewardManager.paused(), true);
        vm.prank(user1);
        vm.expectRevert(
            abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, user1)
        );
        rewardManager.unpause();
        assertEq(rewardManager.paused(), true);

        vm.prank(owner);
        rewardManager.unpause();
        assertEq(rewardManager.paused(), false);

        // User functions work once again
        vm.expectEmit();
        emit AutoClaimEnabled(user1);
        vm.prank(user1);
        rewardManager.enableAutoClaim(false);
    }

    function testPayProposerNoEthPayable() public {
        vm.prank(user1);
        vm.expectRevert(IRewardManager.NoEthPayable.selector);
        rewardManager.payProposer(sampleValPubkey1);
    }

    function testOrphanedRewards() public {

        vm.deal(user1, 4 ether);

        vm.prank(user1);
        vm.expectEmit();
        emit OrphanedRewardsAccumulated(user1, sampleValPubkey1, 2 ether);
        rewardManager.payProposer{value: 2 ether}(sampleValPubkey1);
        assertEq(rewardManager.orphanedRewards(sampleValPubkey1), 2 ether);

        vm.prank(user1);
        vm.expectEmit();
        emit OrphanedRewardsAccumulated(user1, sampleValPubkey2, 1 ether);
        rewardManager.payProposer{value: 1 ether}(sampleValPubkey2);
        assertEq(rewardManager.orphanedRewards(sampleValPubkey2), 1 ether);

        assertEq(user3.balance, 0 ether);
        bytes[] memory pubkeys = new bytes[](2);
        pubkeys[0] = sampleValPubkey1;
        pubkeys[1] = sampleValPubkey2;
        vm.prank(owner);
        vm.expectEmit();
        emit OrphanedRewardsClaimed(user3, 3 ether);
        rewardManager.claimOrphanedRewards(pubkeys, user3);
        assertEq(rewardManager.orphanedRewards(sampleValPubkey1), 0);
        assertEq(rewardManager.orphanedRewards(sampleValPubkey2), 0);
        assertEq(user3.balance, 3 ether);
    }

    function testPayProposerVanillaValidator() public {
        vanillaRegistryTest.testSelfStake();

        address vanillaTestUser = vanillaRegistryTest.user1();
        bytes memory vanillaTestUserPubkey = vanillaRegistryTest.user1BLSKey();

        vm.deal(user2, 4 ether);
        vm.expectEmit();
        emit PaymentStored(user2, vanillaTestUser, vanillaTestUser, 4 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 4 ether}(vanillaTestUserPubkey);

        uint256 balanceBefore = vanillaTestUser.balance;
        vm.prank(vanillaTestUser);
        vm.expectEmit();
        emit RewardsClaimed(vanillaTestUser, 4 ether);
        rewardManager.claimRewards();
        assertEq(vanillaTestUser.balance, balanceBefore + 4 ether);
    }

    function testPayProposerMiddlewareValidator() public {
        mevCommitMiddlewareTest.test_registerValidators();

        address operatorFromMiddlewareTest = vm.addr(0x1117);
        bytes memory pubkey2 = mevCommitMiddlewareTest.sampleValPubkey2();

        vm.deal(user3, 7 ether);
        vm.expectEmit();
        emit PaymentStored(user3, operatorFromMiddlewareTest, operatorFromMiddlewareTest, 7 ether);
        vm.prank(user3);
        rewardManager.payProposer{value: 7 ether}(pubkey2);

        uint256 balanceBefore = operatorFromMiddlewareTest.balance;
        vm.prank(operatorFromMiddlewareTest);
        vm.expectEmit();
        emit RewardsClaimed(operatorFromMiddlewareTest, 7 ether);
        rewardManager.claimRewards();
        assertEq(operatorFromMiddlewareTest.balance, balanceBefore + 7 ether);
    }

    function testPayProposerAVSValidator() public {
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();

        address podOwnerFromAVSTest = address(0x420);
        bytes memory pubkey = mevCommitAVSTest.sampleValPubkey2();

        vm.deal(user4, 10 ether);
        vm.expectEmit();
        emit PaymentStored(user4, podOwnerFromAVSTest, podOwnerFromAVSTest, 10 ether);
        vm.prank(user4);
        rewardManager.payProposer{value: 10 ether}(pubkey);

        uint256 balanceBefore = podOwnerFromAVSTest.balance;
        vm.prank(podOwnerFromAVSTest);
        vm.expectEmit();
        emit RewardsClaimed(podOwnerFromAVSTest, 10 ether);
        rewardManager.claimRewards();
        assertEq(podOwnerFromAVSTest.balance, balanceBefore + 10 ether);
    }

    function testOverrideClaimAddress() public {
        mevCommitMiddlewareTest.test_registerValidators();
        address operatorFromMiddlewareTest = vm.addr(0x1117);
        bytes memory pubkey2 = mevCommitMiddlewareTest.sampleValPubkey2();

        address overrideAddr = vm.addr(0x999999977777777);
        vm.prank(operatorFromMiddlewareTest);
        vm.expectEmit();
        emit OverrideAddressSet(operatorFromMiddlewareTest, overrideAddr);
        rewardManager.overrideReceiver(overrideAddr, false);

        vm.deal(user3, 2 ether);
        vm.expectEmit();
        emit PaymentStored(user3, operatorFromMiddlewareTest, overrideAddr, 2 ether);
        vm.prank(user3);
        rewardManager.payProposer{value: 2 ether}(pubkey2);

        vm.prank(operatorFromMiddlewareTest);
        vm.expectEmit();
        emit OverrideAddressRemoved(operatorFromMiddlewareTest);
        rewardManager.removeOverrideAddress();

        vm.deal(user3, 4 ether);
        vm.expectEmit();
        emit PaymentStored(user3, operatorFromMiddlewareTest, operatorFromMiddlewareTest, 4 ether);
        vm.prank(user3);
        rewardManager.payProposer{value: 4 ether}(pubkey2);

        vm.prank(operatorFromMiddlewareTest);
        vm.expectEmit();
        emit RewardsClaimed(operatorFromMiddlewareTest, 4 ether);
        rewardManager.claimRewards();

        vm.prank(overrideAddr);
        vm.expectEmit();
        emit RewardsClaimed(overrideAddr, 2 ether);
        rewardManager.claimRewards();
    }

    function testMigrateRewardsDuringOverride() public {
        vanillaRegistryTest.testSelfStake();
        address vanillaTestUser = vanillaRegistryTest.user1();
        bytes memory vanillaTestUserPubkey = vanillaRegistryTest.user1BLSKey();

        vm.deal(user2, 4 ether);
        vm.expectEmit();
        emit PaymentStored(user2, vanillaTestUser, vanillaTestUser, 4 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 4 ether}(vanillaTestUserPubkey);

        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 4 ether);
        assertEq(user4.balance, 0 ether);

        vm.expectEmit();
        emit OverrideAddressSet(vanillaTestUser, user4);
        vm.prank(vanillaTestUser);
        rewardManager.overrideReceiver(user4, true);

        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 0 ether);
        assertEq(rewardManager.unclaimedRewards(user4), 4 ether);

        vm.prank(user4);
        vm.expectEmit();
        emit RewardsClaimed(user4, 4 ether);
        rewardManager.claimRewards();
        assertEq(user4.balance, 4 ether);

        vm.deal(user2, 9 ether);
        vm.expectEmit();
        emit PaymentStored(user2, vanillaTestUser, user4, 9 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 9 ether}(vanillaTestUserPubkey);

        assertEq(rewardManager.unclaimedRewards(user4), 9 ether);
        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 0 ether);
        
        vm.expectEmit();
        emit OverrideAddressRemoved(vanillaTestUser);
        vm.prank(vanillaTestUser);
        rewardManager.removeOverrideAddress();

        assertEq(rewardManager.unclaimedRewards(user4), 9 ether);
        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 0 ether);
        
        // Rewards must be claimed manually from the override address, even if that override address is removed
        uint256 balanceBefore = user4.balance;
        vm.prank(user4);
        vm.expectEmit();
        emit RewardsClaimed(user4, 9 ether);
        rewardManager.claimRewards();
        assertEq(user4.balance, balanceBefore + 9 ether);
    }

    function testAutoClaim() public { 
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();

        address podOwnerFromAVSTest = address(0x420);
        bytes memory pubkey = mevCommitAVSTest.sampleValPubkey2();

        vm.deal(user4, 10 ether);
        vm.expectEmit();
        emit PaymentStored(user4, podOwnerFromAVSTest, podOwnerFromAVSTest, 10 ether);
        vm.prank(user4);
        rewardManager.payProposer{value: 10 ether}(pubkey);

        uint256 balanceBefore = podOwnerFromAVSTest.balance;
        vm.expectEmit();
        emit RewardsClaimed(podOwnerFromAVSTest, 10 ether);
        vm.expectEmit();
        emit AutoClaimEnabled(podOwnerFromAVSTest);
        vm.prank(podOwnerFromAVSTest);
        rewardManager.enableAutoClaim(true);
        assertEq(podOwnerFromAVSTest.balance, balanceBefore + 10 ether);

        balanceBefore = podOwnerFromAVSTest.balance;
        vm.deal(user5, 11 ether);
        vm.expectEmit();
        emit AutoClaimed(user5, podOwnerFromAVSTest, podOwnerFromAVSTest, 11 ether);
        vm.prank(user5);
        rewardManager.payProposer{value: 11 ether}(pubkey);
        assertEq(podOwnerFromAVSTest.balance, balanceBefore + 11 ether);
        
        vm.expectEmit();
        emit AutoClaimDisabled(podOwnerFromAVSTest);
        vm.prank(podOwnerFromAVSTest);
        rewardManager.disableAutoClaim();

        balanceBefore = podOwnerFromAVSTest.balance;
        vm.deal(user5, 12 ether);
        vm.expectEmit();
        emit PaymentStored(user5, podOwnerFromAVSTest, podOwnerFromAVSTest, 12 ether);
        vm.prank(user5);
        rewardManager.payProposer{value: 12 ether}(pubkey);
        assertEq(podOwnerFromAVSTest.balance, balanceBefore);
    }

    function testOverrideDoesntAffectAutoClaim() public { 
        vanillaRegistryTest.testSelfStake();
        address vanillaTestUser = vanillaRegistryTest.user1();
        bytes memory vanillaTestUserPubkey = vanillaRegistryTest.user1BLSKey();

        vm.prank(vanillaTestUser);
        vm.expectEmit();
        emit AutoClaimEnabled(vanillaTestUser);
        rewardManager.enableAutoClaim(true);

        uint256 balanceBefore = vanillaTestUser.balance;
        vm.deal(user2, 4 ether);
        vm.expectEmit();
        emit AutoClaimed(user2, vanillaTestUser, vanillaTestUser, 4 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 4 ether}(vanillaTestUserPubkey);
        assertEq(vanillaTestUser.balance, balanceBefore + 4 ether);

        address overrideAddr = vm.addr(0x999999911111111);
        vm.prank(vanillaTestUser);
        vm.expectEmit();
        emit OverrideAddressSet(vanillaTestUser, overrideAddr);
        rewardManager.overrideReceiver(overrideAddr, false);

        balanceBefore = overrideAddr.balance;
        vm.deal(user2, 5 ether);
        vm.expectEmit();
        emit AutoClaimed(user2, vanillaTestUser, overrideAddr, 5 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 5 ether}(vanillaTestUserPubkey);
        assertEq(overrideAddr.balance, balanceBefore + 5 ether);

        vm.prank(vanillaTestUser);
        vm.expectEmit();
        emit AutoClaimDisabled(vanillaTestUser);
        rewardManager.disableAutoClaim();

        vm.deal(user2, 6 ether);
        vm.expectEmit();
        emit PaymentStored(user2, vanillaTestUser, overrideAddr, 6 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 6 ether}(vanillaTestUserPubkey);
    }

    function testAutoClaimBlacklist() public {
        CanRevert canRevert = new CanRevert();

        address[] memory stakers = new address[](1);
        stakers[0] = address(canRevert);
        vm.prank(vanillaRegistry.owner());
        vanillaRegistry.whitelistStakers(stakers);

        vm.deal(address(canRevert), 9 ether);
        assertEq(address(canRevert).balance, 9 ether);
        bytes[] memory validators = new bytes[](1);
        validators[0] = sampleValPubkey4;
        vm.startPrank(address(canRevert));
        vanillaRegistryTest.validatorRegistry().stake{value: 9 ether}(validators);
        vm.stopPrank();
        assertTrue(vanillaRegistryTest.validatorRegistry().isValidatorOptedIn(sampleValPubkey4));

        vm.expectEmit();
        emit AutoClaimEnabled(address(canRevert));
        vm.prank(address(canRevert));
        rewardManager.enableAutoClaim(false);

        assertFalse(rewardManager.autoClaimBlacklist(address(canRevert)));
        assertTrue(rewardManager.autoClaim(address(canRevert)));

        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 0 ether);
        vm.deal(user4, 2 ether);
        vm.expectEmit();
        emit AutoClaimTransferFailed(user4, address(canRevert), address(canRevert));
        vm.prank(user4);
        rewardManager.payProposer{value: 2 ether}(sampleValPubkey4);
        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 2 ether);

        assertTrue(rewardManager.autoClaimBlacklist(address(canRevert)));
        assertFalse(rewardManager.autoClaim(address(canRevert)));

        vm.deal(user4, 3 ether);
        vm.expectEmit();
        emit PaymentStored(user4, address(canRevert), address(canRevert), 3 ether);
        vm.prank(user4);
        rewardManager.payProposer{value: 3 ether}(sampleValPubkey4);
        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 5 ether);

        canRevert.setRevertOnReceive(false);

        vm.prank(address(canRevert));
        vm.expectEmit();
        emit RewardsClaimed(address(canRevert), 5 ether);
        vm.expectEmit();
        emit AutoClaimEnabled(address(canRevert));
        rewardManager.enableAutoClaim(true);

        // Auto claim should not work, blacklist is still active
        vm.deal(user4, 4 ether);
        vm.expectEmit();
        emit PaymentStored(user4, address(canRevert), address(canRevert), 4 ether);
        vm.prank(user4);
        rewardManager.payProposer{value: 4 ether}(sampleValPubkey4);
        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 4 ether);

        vm.prank(owner);
        vm.expectEmit();
        emit RemovedFromAutoClaimBlacklist(address(canRevert));
        rewardManager.removeFromAutoClaimBlacklist(address(canRevert));
        assertFalse(rewardManager.autoClaimBlacklist(address(canRevert)));

        uint256 balanceBefore = address(canRevert).balance;
        vm.deal(user4, 19 ether);
        vm.expectEmit();
        emit AutoClaimed(user4, address(canRevert), address(canRevert), 19 ether);
        vm.prank(user4);
        rewardManager.payProposer{value: 19 ether}(sampleValPubkey4);
        assertEq(address(canRevert).balance, balanceBefore + 19 ether);

        // User still has unclaimed rewards from blacklisted period, this is fine. User can still claim those manually
        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 4 ether);

        balanceBefore = address(canRevert).balance;
        vm.expectEmit();
        emit RewardsClaimed(address(canRevert), 4 ether);
        vm.prank(address(canRevert));
        rewardManager.claimRewards();
        assertEq(address(canRevert).balance, balanceBefore + 4 ether);
        assertEq(rewardManager.unclaimedRewards(address(canRevert)), 0 ether);
    }
}

contract CanRevert {
    bool public revertOnReceive = true;
    receive() external payable {
        if (revertOnReceive) {
            revert("AlwaysReverts");
        }
    }
    function setRevertOnReceive(bool revertOnReceive_) external {
        revertOnReceive = revertOnReceive_;
    }
}
