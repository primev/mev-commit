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
    event PaymentStored(address indexed provider, address indexed toPay, uint256 amount);
    event RewardsClaimed(address indexed toPay, uint256 amount);
    event OverrideClaimAddressSet(address indexed msgSender, address indexed newClaimAddress);
    event OverrideClaimAddressRemoved(address indexed msgSender);
    event AutoClaimed(address indexed provider, address indexed toPay, uint256 amount);

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
        rewardManager.overrideClaimAddress(user2, false);

        vm.prank(user1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        rewardManager.removeOverriddenClaimAddress(false);

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
        emit PaymentStored(user2, vanillaTestUser, 4 ether);
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
        emit PaymentStored(user3, operatorFromMiddlewareTest, 7 ether);
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
        emit PaymentStored(user4, podOwnerFromAVSTest, 10 ether);
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
        emit OverrideClaimAddressSet(operatorFromMiddlewareTest, overrideAddr);
        rewardManager.overrideClaimAddress(overrideAddr, false);

        vm.deal(user3, 2 ether);
        vm.expectEmit();
        emit PaymentStored(user3, overrideAddr, 2 ether);
        vm.prank(user3);
        rewardManager.payProposer{value: 2 ether}(pubkey2);

        vm.prank(operatorFromMiddlewareTest);
        vm.expectEmit();
        emit OverrideClaimAddressRemoved(operatorFromMiddlewareTest);
        rewardManager.removeOverriddenClaimAddress(false);

        vm.deal(user3, 4 ether);
        vm.expectEmit();
        emit PaymentStored(user3, operatorFromMiddlewareTest, 4 ether);
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
        emit PaymentStored(user2, vanillaTestUser, 4 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 4 ether}(vanillaTestUserPubkey);

        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 4 ether);
        assertEq(user4.balance, 0 ether);

        vm.expectEmit();
        emit OverrideClaimAddressSet(vanillaTestUser, user4);
        vm.prank(vanillaTestUser);
        rewardManager.overrideClaimAddress(user4, true);

        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 0 ether);
        assertEq(rewardManager.unclaimedRewards(user4), 4 ether);

        vm.prank(user4);
        vm.expectEmit();
        emit RewardsClaimed(user4, 4 ether);
        rewardManager.claimRewards();
        assertEq(user4.balance, 4 ether);

        vm.deal(user2, 9 ether);
        vm.expectEmit();
        emit PaymentStored(user2, user4, 9 ether);
        vm.prank(user2);
        rewardManager.payProposer{value: 9 ether}(vanillaTestUserPubkey);

        assertEq(rewardManager.unclaimedRewards(user4), 9 ether);
        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 0 ether);
        
        vm.expectEmit();
        emit OverrideClaimAddressRemoved(vanillaTestUser);
        vm.prank(vanillaTestUser);
        rewardManager.removeOverriddenClaimAddress(true);

        assertEq(rewardManager.unclaimedRewards(user4), 0 ether);
        assertEq(rewardManager.unclaimedRewards(vanillaTestUser), 9 ether);
        
        uint256 balanceBefore = vanillaTestUser.balance;
        vm.prank(vanillaTestUser);
        vm.expectEmit();
        emit RewardsClaimed(vanillaTestUser, 9 ether);
        rewardManager.claimRewards();
        assertEq(vanillaTestUser.balance, balanceBefore + 9 ether);
    }

    function testAutoClaim() public { 
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();

        address podOwnerFromAVSTest = address(0x420);
        bytes memory pubkey = mevCommitAVSTest.sampleValPubkey2();

        vm.deal(user4, 10 ether);
        vm.expectEmit();
        emit PaymentStored(user4, podOwnerFromAVSTest, 10 ether);
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
        emit AutoClaimed(user5, podOwnerFromAVSTest, 11 ether);
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
        emit PaymentStored(user5, podOwnerFromAVSTest, 12 ether);
        vm.prank(user5);
        rewardManager.payProposer{value: 12 ether}(pubkey);
        assertEq(podOwnerFromAVSTest.balance, balanceBefore);
    }

    function testAutoClaimBlacklist() public {
        // blacklist and then removeFromAutoClaimBlacklist
    }
}
