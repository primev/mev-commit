// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IProviderRegistry} from "../../contracts/interfaces/IProviderRegistry.sol";
import {MockBLSVerify} from "../precompiles/BLSVerifyPreCompileMockTest.sol";
import {DepositManager} from "../../contracts/core/DepositManager.sol";

contract ReputationProviderTest is Test {
    ProviderRegistry public providerRegistry;
    BidderRegistry public bidderRegistry;
    PreconfManager public preconfManager;
    BlockTracker public blockTracker;

    uint256 public minStake;
    uint256 public reputationMinStake;
    uint256 public feePercent;
    uint256 public withdrawalDelay;
    address public feeRecipient;
    address public owner;

    address public standardProvider;
    address public reputationProvider;
    address public randomUser;

    event ProviderRegistered(address indexed provider, uint256 stakedAmount);
    event ReputationRegistrationRequested(address indexed provider, uint256 stakedAmount);
    event ReputationRegistrationCancelled(address indexed provider, uint256 returnedAmount);
    event ReputationProviderApproved(address indexed provider, address indexed approver);
    event ReputationProviderRemoved(address indexed provider, address indexed approver);
    event Unstake(address indexed provider, uint256 timestamp);

    function setUp() public {
        address BLS_VERIFY_ADDRESS = address(0xf0);
        bytes memory code = type(MockBLSVerify).creationCode;
        vm.etch(BLS_VERIFY_ADDRESS, code);

        owner = address(this);
        minStake = 5 ether;
        reputationMinStake = 0.01 ether;
        feePercent = 10 * 1e16;
        feeRecipient = vm.addr(9);
        withdrawalDelay = 24 hours;

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (minStake, feeRecipient, feePercent, owner, withdrawalDelay, 10000)
            )
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));
        providerRegistry.setReputationMinStake(reputationMinStake);

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (owner, owner))
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(
                BidderRegistry.initialize,
                (feeRecipient, feePercent, owner, address(blockTracker), 10000, 10000)
            )
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        DepositManager depositManager = new DepositManager(address(bidderRegistry), 0.01 ether);
        bidderRegistry.setDepositManagerImpl(address(depositManager));

        address preconfManagerProxy = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(
                PreconfManager.initialize,
                (address(providerRegistry), address(bidderRegistry), address(blockTracker), feeRecipient, owner, 500)
            )
        );
        preconfManager = PreconfManager(payable(preconfManagerProxy));
        providerRegistry.setPreconfManager(address(preconfManager));

        standardProvider = vm.addr(1);
        reputationProvider = vm.addr(2);
        randomUser = vm.addr(3);

        vm.deal(standardProvider, 100 ether);
        vm.deal(reputationProvider, 100 ether);
        vm.deal(randomUser, 100 ether);

        // Register a standard provider
        vm.prank(standardProvider);
        providerRegistry.registerAndStake{value: 5 ether}();
    }

    // =========== Request Registration ===========

    function test_RequestReputationRegistration() public {
        vm.prank(reputationProvider);
        vm.expectEmit(true, false, false, true);
        emit ReputationRegistrationRequested(reputationProvider, 0.01 ether);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        assertEq(providerRegistry.pendingReputationStake(reputationProvider), 0.01 ether);
        assertFalse(providerRegistry.providerRegistered(reputationProvider));
    }

    function test_RevertWhen_RequestWithInsufficientStake() public {
        vm.prank(reputationProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.InsufficientReputationStake.selector, 0.001 ether, reputationMinStake
        ));
        providerRegistry.requestReputationRegistration{value: 0.001 ether}();
    }

    function test_RevertWhen_RequestWhenAlreadyRegistered() public {
        vm.prank(standardProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.ProviderAlreadyRegistered.selector, standardProvider
        ));
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
    }

    function test_RevertWhen_RequestWhenPendingExists() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(reputationProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.PendingReputationRequestExists.selector, reputationProvider
        ));
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
    }

    function test_RevertWhen_RequestWithReputationMinStakeNotSet() public {
        // Deploy fresh registry without setting reputationMinStake
        address freshProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (minStake, feeRecipient, feePercent, owner, withdrawalDelay, 10000)
            )
        );
        ProviderRegistry freshRegistry = ProviderRegistry(payable(freshProxy));

        vm.prank(reputationProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.InsufficientReputationStake.selector, 0, 0
        ));
        freshRegistry.requestReputationRegistration{value: 0.01 ether}();
    }

    // =========== Cancel Registration ===========

    function test_CancelReputationRegistration() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        uint256 balanceBefore = reputationProvider.balance;

        vm.prank(reputationProvider);
        vm.expectEmit(true, false, false, true);
        emit ReputationRegistrationCancelled(reputationProvider, 0.01 ether);
        providerRegistry.cancelReputationRegistration();

        assertEq(providerRegistry.pendingReputationStake(reputationProvider), 0);
        assertEq(reputationProvider.balance, balanceBefore + 0.01 ether);
    }

    function test_RevertWhen_CancelWithNoPendingRequest() public {
        vm.prank(reputationProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.NoPendingReputationRequest.selector, reputationProvider
        ));
        providerRegistry.cancelReputationRegistration();
    }

    // =========== Approve Registration ===========

    function test_ApproveByOwner() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(owner);
        vm.expectEmit(true, false, false, true);
        emit ProviderRegistered(reputationProvider, 0.01 ether);
        providerRegistry.approveReputationRegistration(reputationProvider);

        assertTrue(providerRegistry.providerRegistered(reputationProvider));
        assertEq(providerRegistry.getProviderStake(reputationProvider), 0.01 ether);
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), owner);
        assertEq(providerRegistry.pendingReputationStake(reputationProvider), 0);
        // Owner should not have lock count incremented
        assertEq(providerRegistry.approverLockCount(owner), 0);
    }

    function test_ApproveByStandardProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        assertTrue(providerRegistry.providerRegistered(reputationProvider));
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), standardProvider);
        assertEq(providerRegistry.approverLockCount(standardProvider), 1);
    }

    function test_RevertWhen_ApproveByUnregisteredProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(randomUser);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.ProviderNotRegistered.selector, randomUser
        ));
        providerRegistry.approveReputationRegistration(reputationProvider);
    }

    function test_RevertWhen_ApproveByProviderWithInsufficientStake() public {
        // Register a provider, then slash them below minStake
        address weakProvider = vm.addr(10);
        vm.deal(weakProvider, 100 ether);
        vm.prank(weakProvider);
        providerRegistry.registerAndStake{value: 5 ether}();

        // Slash them down
        vm.prank(address(preconfManager));
        providerRegistry.slash(4.5 ether, weakProvider, payable(randomUser));

        // Now they have 0.05 ether stake (< 5 ether minStake)
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(weakProvider);
        vm.expectRevert();
        providerRegistry.approveReputationRegistration(reputationProvider);
    }

    function test_RevertWhen_ApproveWithNoPendingRequest() public {
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.NoPendingReputationRequest.selector, reputationProvider
        ));
        providerRegistry.approveReputationRegistration(reputationProvider);
    }

    // =========== isProviderValid ===========

    function test_IsProviderValidForReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Should not revert - reputation provider is valid with 0.01 ETH
        providerRegistry.isProviderValid(reputationProvider);
    }

    function test_AreProvidersValidIncludesReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Add a BLS key so areProvidersValid passes the hasBLSKey check
        bytes memory blsKey = hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
        vm.prank(owner);
        providerRegistry.overrideAddBLSKey(reputationProvider, blsKey);

        address[] memory providers = new address[](2);
        providers[0] = standardProvider;
        providers[1] = reputationProvider;

        bool[] memory results = providerRegistry.areProvidersValid(providers);
        // standardProvider has no BLS key in this test, so it would be false
        // reputationProvider has a BLS key and sufficient reputation stake
        assertTrue(results[1], "Reputation provider should be valid");
    }

    // =========== Remove Reputation Provider ===========

    function test_RemoveByApprover() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        assertEq(providerRegistry.approverLockCount(standardProvider), 1);

        vm.prank(standardProvider);
        vm.expectEmit(true, true, false, true);
        emit ReputationProviderRemoved(reputationProvider, standardProvider);
        providerRegistry.removeReputationProvider(reputationProvider);

        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);
        // Unstake should have been initiated
        assertEq(providerRegistry.withdrawalRequests(reputationProvider), block.timestamp);
    }

    function test_RemoveByOwnerWhenOwnerApproved() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        vm.prank(owner);
        providerRegistry.removeReputationProvider(reputationProvider);

        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
    }

    function test_RevertWhen_OwnerRemovesProviderApprovedByOther() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.NotApprover.selector, owner
        ));
        providerRegistry.removeReputationProvider(reputationProvider);
    }

    function test_RevertWhen_RandomUserRemoves() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        vm.prank(randomUser);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.NotApprover.selector, randomUser
        ));
        providerRegistry.removeReputationProvider(reputationProvider);
    }

    function test_RevertWhen_RemoveNonReputationProvider() public {
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.ProviderIsNotReputationProvider.selector, standardProvider
        ));
        providerRegistry.removeReputationProvider(standardProvider);
    }

    function test_RemoveFullySlashedReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Slash the reputation provider to 0
        vm.prank(address(preconfManager));
        providerRegistry.slash(0.01 ether, reputationProvider, payable(randomUser));

        assertEq(providerRegistry.getProviderStake(reputationProvider), 0);

        vm.prank(owner);
        providerRegistry.removeReputationProvider(reputationProvider);

        // Should be fully deregistered since stake is 0
        assertFalse(providerRegistry.providerRegistered(reputationProvider));
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
    }

    // =========== Approver Withdrawal Lock ===========

    function test_ApproverCannotUnstakeWithActiveReputationProviders() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        vm.prank(standardProvider);
        vm.expectRevert(abi.encodeWithSelector(
            IProviderRegistry.ApproverHasActiveReputationProviders.selector,
            standardProvider,
            1
        ));
        providerRegistry.unstake();
    }

    function test_ApproverCanUnstakeAfterRemovingReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Remove the reputation provider first
        vm.prank(standardProvider);
        providerRegistry.removeReputationProvider(reputationProvider);

        // Now unstake should work
        vm.prank(standardProvider);
        providerRegistry.unstake();
        assertEq(providerRegistry.withdrawalRequests(standardProvider), block.timestamp);
    }

    function test_OwnerNotLockedWhenApprovingReputationProviders() public {
        // Owner should never be locked
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        assertEq(providerRegistry.approverLockCount(owner), 0);
    }

    function test_ApproverLockCountWithMultipleReputationProviders() public {
        address repProvider2 = vm.addr(20);
        vm.deal(repProvider2, 100 ether);

        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(repProvider2);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(repProvider2);

        assertEq(providerRegistry.approverLockCount(standardProvider), 2);

        // Can't unstake with 2 active
        vm.prank(standardProvider);
        vm.expectRevert();
        providerRegistry.unstake();

        // Remove one
        vm.prank(standardProvider);
        providerRegistry.removeReputationProvider(reputationProvider);
        assertEq(providerRegistry.approverLockCount(standardProvider), 1);

        // Still can't unstake with 1 active
        vm.prank(standardProvider);
        vm.expectRevert();
        providerRegistry.unstake();

        // Remove the other
        vm.prank(standardProvider);
        providerRegistry.removeReputationProvider(repProvider2);
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);

        // Now can unstake
        vm.prank(standardProvider);
        providerRegistry.unstake();
    }

    // =========== Auto-Conversion ===========

    function test_AutoConvertWhenStakeReachesMinStake() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        assertEq(providerRegistry.approverLockCount(standardProvider), 1);

        // Stake up to minStake
        vm.prank(reputationProvider);
        providerRegistry.stake{value: 4.99 ether}();

        // Should be auto-converted - reputation status cleared, approver lock released
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);
        assertTrue(providerRegistry.providerRegistered(reputationProvider));
        assertEq(providerRegistry.getProviderStake(reputationProvider), 5 ether);

        // isProviderValid should still pass (now checked against minStake)
        providerRegistry.isProviderValid(reputationProvider);
    }

    function test_AutoConvertDoesNotTriggerBelowMinStake() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Stake but not enough to reach minStake
        vm.prank(reputationProvider);
        providerRegistry.stake{value: 1 ether}();

        // Should still be a reputation provider
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), standardProvider);
        assertEq(providerRegistry.approverLockCount(standardProvider), 1);
    }

    function test_AutoConvertReleasesOwnerApproverLock() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Owner lock count should be 0 (owner is exempt)
        assertEq(providerRegistry.approverLockCount(owner), 0);

        // Stake up to minStake
        vm.prank(reputationProvider);
        providerRegistry.stake{value: 4.99 ether}();

        // Should be auto-converted
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
        // Owner lock count should still be 0
        assertEq(providerRegistry.approverLockCount(owner), 0);
    }

    // =========== Withdraw Clears Reputation Status ===========

    function test_WithdrawClearsReputationStatus() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Reputation provider voluntarily unstakes
        vm.prank(reputationProvider);
        providerRegistry.unstake();

        vm.warp(block.timestamp + 24 hours);

        vm.prank(reputationProvider);
        providerRegistry.withdraw();

        assertFalse(providerRegistry.providerRegistered(reputationProvider));
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);
        assertEq(providerRegistry.getProviderStake(reputationProvider), 0);
    }

    function test_ReRegistrationRequiresNewApproval() public {
        // Full lifecycle: request -> approve -> unstake -> withdraw -> request again
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        vm.prank(reputationProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours);
        vm.prank(reputationProvider);
        providerRegistry.withdraw();

        // Should be able to request again
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        assertEq(providerRegistry.pendingReputationStake(reputationProvider), 0.01 ether);
    }

    // =========== Reputation Provider Cannot Approve Others ===========

    function test_ReputationProviderCannotApproveOthers() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Another provider requests
        address repProvider2 = vm.addr(20);
        vm.deal(repProvider2, 100 ether);
        vm.prank(repProvider2);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();

        // Reputation provider tries to approve — should fail because their stake < minStake
        vm.prank(reputationProvider);
        vm.expectRevert();
        providerRegistry.approveReputationRegistration(repProvider2);
    }

    // =========== Slash Reputation Provider ===========

    function test_SlashReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(owner);
        providerRegistry.approveReputationRegistration(reputationProvider);

        address bidder = vm.addr(30);
        uint256 bidderBefore = bidder.balance;

        vm.prank(address(preconfManager));
        providerRegistry.slash(0.005 ether, reputationProvider, payable(bidder));

        assertEq(bidder.balance - bidderBefore, 0.005 ether);
        // Stake reduced (minus slash + penalty fee)
        assertTrue(providerRegistry.getProviderStake(reputationProvider) < 0.01 ether);
    }

    // =========== Edge Cases ===========

    function test_DelegateStakeTriggersAutoConversion() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Someone else stakes on behalf of the reputation provider
        vm.prank(randomUser);
        providerRegistry.delegateStake{value: 4.99 ether}(reputationProvider);

        // Should be auto-converted
        assertEq(providerRegistry.reputationProviderApprover(reputationProvider), address(0));
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);
        assertEq(providerRegistry.getProviderStake(reputationProvider), 5 ether);
    }

    function test_RemoveAlreadyUnstakingReputationProvider() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.01 ether}();
        vm.prank(standardProvider);
        providerRegistry.approveReputationRegistration(reputationProvider);

        // Reputation provider starts unstaking themselves
        vm.prank(reputationProvider);
        providerRegistry.unstake();

        // Approver also removes them — should not overwrite the withdrawal request
        uint256 originalTimestamp = providerRegistry.withdrawalRequests(reputationProvider);
        vm.prank(standardProvider);
        providerRegistry.removeReputationProvider(reputationProvider);

        // Withdrawal request should remain the same (not overwritten)
        assertEq(providerRegistry.withdrawalRequests(reputationProvider), originalTimestamp);
        assertEq(providerRegistry.approverLockCount(standardProvider), 0);
    }

    function test_RequestWithExactReputationMinStake() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: reputationMinStake}();
        assertEq(providerRegistry.pendingReputationStake(reputationProvider), reputationMinStake);
    }

    function test_RequestWithMoreThanReputationMinStake() public {
        vm.prank(reputationProvider);
        providerRegistry.requestReputationRegistration{value: 0.5 ether}();
        assertEq(providerRegistry.pendingReputationStake(reputationProvider), 0.5 ether);
    }
}
