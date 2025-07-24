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

contract ProviderRegistryTest is Test {
    uint256 public testNumber;
    ProviderRegistry public providerRegistry;
    uint256 public feePercent;
    uint256 public minStake;
    address public provider;
    address public feeRecipient;
    BidderRegistry public bidderRegistry;
    PreconfManager public preconfManager;
    BlockTracker public blockTracker;
    uint256 public withdrawalDelay;
    bytes public validBLSPubkey =
        hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    bytes public dummyBLSSignature =
        hex"bbbbbbbbb1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2";
    bytes[] public validBLSPubkeys = [validBLSPubkey];
    uint256 public penaltyFeePayoutPeriodMs;
    mapping(address => uint256) public commitmentsCount; // For use in test_RevertWhen_WithdrawStakedAmountWithoutCommitments
    event ProviderRegistered(address indexed provider, uint256 stakedAmount);
    event WithdrawalRequested(address indexed provider, uint256 timestamp);
    event WithdrawalCompleted(address indexed provider, uint256 amount);
    event FeeTransfer(uint256 amount, address indexed recipient);
    event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient);
    event FeePayoutPeriodUpdated(uint256 indexed newFeePayoutPeriod);
    event InsufficientFundsToSlash(
        address indexed provider,
        uint256 providerStake,
        uint256 residualAmt,
        uint256 penaltyFee,
        uint256 slashAmt
    );
    event FundsSlashed(address indexed provider, uint256 totalSlash);
    event TransferToBidderFailed(address indexed bidder, uint256 amount);

    function setUp() public {
        address BLS_VERIFY_ADDRESS = address(0xf0);
        bytes memory code = type(MockBLSVerify).creationCode;
        vm.etch(BLS_VERIFY_ADDRESS, code);

        testNumber = 42;
        feePercent = 10 * 1e16;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        withdrawalDelay = 24 hours; // 24 hours
        penaltyFeePayoutPeriodMs = 10000;
        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (
                    minStake,
                    feeRecipient,
                    feePercent,
                    address(this),
                    withdrawalDelay,
                    penaltyFeePayoutPeriodMs
                )
            )
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(
                BlockTracker.initialize,
                (address(this), address(this))
            )
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(
                BidderRegistry.initialize,
                (
                    feeRecipient,
                    feePercent,
                    address(this),
                    address(blockTracker),
                    penaltyFeePayoutPeriodMs
                )
            )
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        address preconfStoreProxy = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(
                PreconfManager.initialize,
                (
                    address(providerRegistry), // Provider Registry
                    address(bidderRegistry), // User Registry
                    address(blockTracker), // Block Tracker
                    feeRecipient, // Oracle
                    address(this),
                    500
                )
            )
        );
        preconfManager = PreconfManager(payable(preconfStoreProxy));

        provider = vm.addr(1);
        vm.deal(provider, 100 ether);
        vm.deal(address(this), 100 ether);
    }

    function test_VerifyInitialContractState() public view {
        assertEq(providerRegistry.minStake(), 1e18 wei);
        assertEq(feePercent, feePercent);
        assertEq(withdrawalDelay, withdrawalDelay);
        assertEq(providerRegistry.feePercent(), feePercent);
        assertEq(providerRegistry.preconfManager(), address(0));
        assertEq(providerRegistry.providerRegistered(provider), false);
        (
            address recipient,
            uint256 accumulatedAmount,
            uint256 lastPayoutTimestamp,
            uint256 payoutPeriodMs
        ) = bidderRegistry.protocolFeeTracker();
        assertEq(recipient, feeRecipient);
        assertEq(payoutPeriodMs, penaltyFeePayoutPeriodMs);
        assertEq(lastPayoutTimestamp, block.timestamp);
        assertEq(accumulatedAmount, 0);
    }

    function test_RevertWhen_ProviderStakeAndRegisterMinStake() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert();
        providerRegistry.registerAndStake{value: 1 wei}();
    }

    function test_ProviderStakeAndRegister() public {
        vm.deal(provider, 3 ether);
        vm.startPrank(provider);
        vm.expectEmit(true, false, false, true);

        emit ProviderRegistered(provider, 1e18 wei);

        providerRegistry.registerAndStake{value: 1e18 wei}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();
        bool isProviderRegistered = providerRegistry.providerRegistered(
            provider
        );
        assertEq(isProviderRegistered, true);

        uint256 providerStakeStored = providerRegistry.getProviderStake(
            provider
        );
        assertEq(providerStakeStored, 1e18 wei);

        // Check if BLS keys were correctly registered
        bytes[] memory storedBLSKeys = providerRegistry.getBLSKeys(provider);
        assertEq(
            storedBLSKeys.length,
            validBLSPubkeys.length,
            "BLS keys array length mismatch"
        );

        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            assertEq(storedBLSKeys[i], validBLSPubkeys[i], "BLS key mismatch");
            address storedProvider = providerRegistry.getEoaFromBLSKey(
                validBLSPubkeys[i]
            );
            assertEq(
                storedProvider,
                provider,
                "Provider address mismatch for BLS key"
            );
        }
    }

    function test_RevertWhen_ProviderStakeAndRegisterAlreadyRegistered() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        vm.expectRevert();
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 1 wei}();
    }

    function test_RevertWhen_Receive() public {
        vm.prank(provider);
        vm.expectRevert();
        payable(address(providerRegistry)).transfer(1 wei);
    }

    function test_RevertWhen_Fallback() public {
        vm.prank(provider);
        bytes memory data = abi.encode(1, 2);
        (bool success, ) = address(providerRegistry).call{value: 1 wei}(data);
        require(!success, "should revert");
    }

    function test_SetNewPenaltyFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit PenaltyFeeRecipientUpdated(newRecipient);
        providerRegistry.setNewPenaltyFeeRecipient(newRecipient);
        (address recipient, , , ) = providerRegistry.penaltyFeeTracker();
        assertEq(recipient, newRecipient);
    }

    function test_RevertWhen_SetNewPenaltyFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert();
        vm.prank(vm.addr(1));
        providerRegistry.setNewPenaltyFeeRecipient(newRecipient);
    }

    function test_SetNewFeePayoutPeriod() public {
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit FeePayoutPeriodUpdated(890);
        providerRegistry.setFeePayoutPeriod(890);
        (, , , uint256 payoutPeriodMs) = providerRegistry
            .penaltyFeeTracker();
        assertEq(payoutPeriodMs, 890);
    }

    function test_RevertWhen_SetNewFeePayoutPeriod() public {
        vm.expectRevert();
        vm.prank(vm.addr(1));
        providerRegistry.setFeePayoutPeriod(83424);
    }

    function test_SetNewFeePercent() public {
        vm.prank(address(this));
        providerRegistry.setNewFeePercent(25);

        assertEq(providerRegistry.feePercent(), 25);
    }

    function test_RevertWhen_SetNewFeePercent() public {
        vm.expectRevert();
        vm.prank(vm.addr(1));
        providerRegistry.setNewFeePercent(25);
    }

    function test_SetPreConfContract() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        providerRegistry.setPreconfManager(newPreConfContract);

        assertEq(providerRegistry.preconfManager(), newPreConfContract);
    }

    function test_RevertWhen_SetPreConfContract() public {
        vm.prank(vm.addr(1));
        vm.expectRevert();
        providerRegistry.setPreconfManager(address(0));
    }

    function test_ShouldSlashProvider() public {
        providerRegistry.setPreconfManager(address(this));
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();

        address bidder = vm.addr(4);

        vm.expectCall(bidder, 1000000000000000000 wei, new bytes(0));
        providerRegistry.slash(1 ether, 0, provider, payable(bidder), 1e18);

        assertEq(bidder.balance, 1000000000000000000 wei);
        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            100000000000000000 wei
        );
        assertEq(providerRegistry.providerStakes(provider), 0.9 ether);
    }

    function test_ShouldSlashProviderWithoutFeeRecipient() public {
        vm.prank(address(this));
        providerRegistry.setNewPenaltyFeeRecipient(address(0));
        providerRegistry.setPreconfManager(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();

        address bidder = vm.addr(4);

        vm.expectCall(bidder, 1000000000000000000 wei, new bytes(0));
        providerRegistry.slash(
            1 ether,
            0,
            provider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );

        assertEq(bidder.balance, 1000000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 0.9 ether);
    }

    function test_RevertWhen_ShouldRetrieveFundsNotPreConf() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();

        address bidder = vm.addr(4);
        uint256 residualBidPercentAfterDecay = providerRegistry.ONE_HUNDRED_PERCENT();
        vm.expectRevert();
        providerRegistry.slash(
            1 ether,
            0,
            provider,
            payable(bidder),
            residualBidPercentAfterDecay
        );
    }

    function test_ShouldRetrieveFundsWhenSlashIsGreaterThanStake() public {
        vm.prank(address(this));
        providerRegistry.setPreconfManager(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();

        address bidder = vm.addr(4);
        vm.prank(address(this));

        vm.expectEmit(true, true, true, true);
        emit InsufficientFundsToSlash(provider, 2 ether, 3 ether, 0.3 ether, 0);
        providerRegistry.slash(
            3 ether,
            0,
            provider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );

        assertEq(providerRegistry.getAccumulatedPenaltyFee(), 0);
        assertEq(providerRegistry.providerStakes(provider), 0 ether);
    }

    function test_ShouldRetrieveFundsWhenSlashIsGreaterThanStakePenaltyNotCovered()
        public
    {
        vm.prank(address(this));
        providerRegistry.setPreconfManager(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 3 ether}();

        address bidder = vm.addr(4);
        vm.prank(address(this));

        vm.expectEmit(true, true, true, true);
        emit InsufficientFundsToSlash(provider, 3 ether, 3 ether, 0.3 ether, 0);
        providerRegistry.slash(
            3 ether,
            0,
            provider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );

        assertEq(providerRegistry.getAccumulatedPenaltyFee(), 0);
        assertEq(providerRegistry.providerStakes(provider), 0 ether);
    }

    function test_ShouldRetrieveFundsWhenSlashIsGreaterThanStakePenaltyNotFullyCovered()
        public
    {
        vm.prank(address(this));
        providerRegistry.setPreconfManager(address(this));

        vm.deal(provider, 3.1 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 3.1 ether}();

        address bidder = vm.addr(4);
        vm.prank(address(this));

        vm.expectEmit(true, true, true, true);
        emit InsufficientFundsToSlash(
            provider,
            3.1 ether,
            3 ether,
            0.3 ether,
            0
        );
        providerRegistry.slash(
            3 ether,
            0,
            provider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );

        assertEq(providerRegistry.getAccumulatedPenaltyFee(), 0.1 ether);
        assertEq(providerRegistry.providerStakes(provider), 0 ether);
    }

    function test_PenaltyFeeBehavior() public {
        providerRegistry.setNewPenaltyFeeRecipient(vm.addr(6));
        vm.deal(provider, 3 ether);
        vm.prank(provider);

        address bidder = vm.addr(4);

        providerRegistry.registerAndStake{value: 2 ether}();

        providerRegistry.setPreconfManager(address(this));
        providerRegistry.slash(
            1e18 wei,
            0,
            provider,
            payable(bidder),
            50 * providerRegistry.PRECISION()
        );
        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            5e16 wei,
            "FeeRecipientAmount should match"
        );

        address newProvider = vm.addr(11);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2 ether}();

        vm.warp(block.timestamp + 10000 + 1); // roll past protocol fee payout period

        vm.expectEmit(true, true, true, true);
        emit FeeTransfer(1e17 wei, vm.addr(6));
        providerRegistry.slash(
            1e18 wei,
            0,
            newProvider,
            payable(bidder),
            50 * providerRegistry.PRECISION()
        );

        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            0,
            "Accumulated protocol fee should be zero"
        );
        assertEq(
            vm.addr(6).balance,
            1e17 wei,
            "FeeRecipient should have received 1e17 wei"
        );
    }

    function test_WithdrawStakedAmountWithoutFeeRecipient() public {
        providerRegistry.setNewPenaltyFeeRecipient(address(0));
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        uint256 percent = providerRegistry.ONE_HUNDRED_PERCENT();
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        providerRegistry.setPreconfManager(address(preconfManager));
        vm.prank(address(preconfManager));
        providerRegistry.slash(
            1e18 wei,
            0,
            newProvider,
            payable(bidder),
            percent
        );
        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.prank(newProvider);
        providerRegistry.withdraw();
        assertEq(
            providerRegistry.providerStakes(newProvider),
            0,
            "Provider's staked amount should be zero after withdrawal"
        );
        assertEq(
            newProvider.balance,
            1.9e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function test_RevertWhen_WithdrawStakedAmountUnauthorized() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        address wrongNewProvider = vm.addr(12);
        vm.prank(wrongNewProvider);
        vm.expectRevert();
        providerRegistry.unstake();

        vm.prank(wrongNewProvider);
        vm.expectRevert();
        providerRegistry.withdraw();
    }

    function test_RegisterAndStake() public {
        address newProvider = vm.addr(5);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        assertEq(
            providerRegistry.providerStakes(newProvider),
            2e18 wei,
            "Staked amount should match"
        );
        assertEq(
            providerRegistry.providerRegistered(newProvider),
            true,
            "Provider should be registered"
        );
    }

    function test_RevertWhen_WithdrawStakedAmountWithoutCommitments() public {
        providerRegistry.setPreconfManager(address(this));
        commitmentsCount[vm.addr(8)] = 2;

        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.expectRevert(abi.encodeWithSelector(IProviderRegistry.ProviderCommitmentsPending.selector, newProvider, 2));
        vm.prank(newProvider);
        providerRegistry.withdraw();
    }

    function test_RequestWithdrawal() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        vm.prank(newProvider);
        providerRegistry.unstake();
        assertEq(
            providerRegistry.withdrawalRequests(newProvider),
            block.timestamp,
            "Withdrawal request timestamp should match"
        );
    }

    function test_WithdrawStakedAmount() public {
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        uint256 percent = providerRegistry.ONE_HUNDRED_PERCENT();
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        providerRegistry.setPreconfManager(address(preconfManager));
        vm.prank(address(preconfManager));
        providerRegistry.slash(
            1e18 wei,
            0,
            newProvider,
            payable(bidder),
            percent
        );
        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.prank(newProvider);
        providerRegistry.withdraw();
        assertEq(
            providerRegistry.providerStakes(newProvider),
            0,
            "Provider's staked amount should be zero after withdrawal"
        );
        assertEq(
            newProvider.balance,
            1.9e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function test_WithdrawStakedAmountBefore24Hours() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 23 hours); // Move forward less than 24 hours
        vm.prank(newProvider);
        vm.expectRevert(
            abi.encodeWithSelector(
                IProviderRegistry.DelayNotPassed.selector,
                block.timestamp - 23 hours, // withdrawalRequestTimestamp
                24 hours, // withdrawalDelay
                block.timestamp // currentBlockTimestamp
            )
        );
        providerRegistry.withdraw();
    }

    function test_WithdrawStakedAmountWithoutRequest() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();

        vm.prank(newProvider);
        vm.expectRevert(
            abi.encodeWithSelector(
                IProviderRegistry.NoUnstakeRequest.selector,
                newProvider
            )
        );
        providerRegistry.withdraw();
    }

    function test_DelegateRegisterAndStake() public {
        address newProvider = vm.addr(7);

        vm.prank(address(this));
        providerRegistry.delegateRegisterAndStake{value: 2e18 wei}(newProvider);

        assertEq(
            providerRegistry.providerStakes(newProvider),
            2e18 wei,
            "Staked amount should match"
        );
        assertEq(
            providerRegistry.providerRegistered(newProvider),
            true,
            "Provider should be registered"
        );
    }

    /// @dev Test that if the bidder’s transfer fails, the slashed bidder amount is recorded.
    function test_SlashBidderSendFailure() public {
        providerRegistry.setPreconfManager(address(this));

        address testProvider = vm.addr(302);
        vm.deal(testProvider, 3 ether);
        vm.prank(testProvider);
        providerRegistry.registerAndStake{value: 2 ether}();

        // Simulate a bidder that always reverts on receiving ETH.
        // (Using vm.etch to deploy minimal bytecode that simply reverts.)
        address revertingBidder = vm.addr(303);
        // Minimal code that immediately reverts (opcode 0xFD is REVERT).
        vm.etch(revertingBidder, hex"6000fd");

        uint256 initialStake = providerRegistry.getProviderStake(testProvider);

        // With the parameters below:
        //   amt = 1 ether
        //   slashAmt = 0
        //   residualBidPercentAfterDecay = ONE_HUNDRED_PERCENT
        // we have:
        //   residualAmt = 1 ether,
        //   penaltyFee = (1 ether * feePercent / ONE_HUNDRED_PERCENT) (e.g. 0.1 ether for feePercent = 10*1e16),
        //   bidderPortion = 1 ether,
        //   totalSlash = 1 ether + penaltyFee.
        vm.prank(address(this));
        providerRegistry.slash(
            1 ether,
            0,
            testProvider,
            payable(revertingBidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );

        assertEq(
            providerRegistry.bidderSlashedAmount(revertingBidder),
            1 ether,
            "bidderSlashedAmount mismatch"
        );

        uint256 expectedPenaltyFee = (1 ether * providerRegistry.feePercent()) /
            providerRegistry.ONE_HUNDRED_PERCENT();
        uint256 expectedTotalSlash = 1 ether + expectedPenaltyFee;

        assertEq(
            providerRegistry.getProviderStake(testProvider),
            initialStake - expectedTotalSlash,
            "Provider stake not reduced correctly"
        );
    }

    /// @dev Test that a successful slash emits the FundsSlashed event with the correct parameters.
    function test_SlashEmitsFundsSlashedEvent() public {
        providerRegistry.setPreconfManager(address(this));

        address testProvider = vm.addr(304);
        vm.deal(testProvider, 3 ether);
        vm.prank(testProvider);
        providerRegistry.registerAndStake{value: 2 ether}();

        address bidder = vm.addr(305);
        uint256 expectedPenaltyFee = (1 ether * providerRegistry.feePercent()) /
            providerRegistry.ONE_HUNDRED_PERCENT();
        uint256 expectedTotalSlash = 1 ether + expectedPenaltyFee;

        // Expect the FundsSlashed event to be emitted with the proper parameters.
        vm.prank(address(this));
        vm.expectEmit(true, false, false, true);
        emit FundsSlashed(testProvider, expectedTotalSlash);
        providerRegistry.slash(
            1 ether,
            0,
            testProvider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );
    }

    /// @dev Test the branch where providerStake is less than bidderPortion.
    /// In this case, with:
    ///   - amt = 1 ether,
    ///   - slashAmt = 0.5 ether, and
    ///   - residualBidPercentAfterDecay = ONE_HUNDRED_PERCENT,
    /// we have:
    ///   residualAmt = 1 ether,
    ///   penaltyFee = (1 ether * feePercent)/ONE_HUNDRED_PERCENT (≈0.1 ether for feePercent = 10%),
    ///   bidderPortion = 1 ether + 0.5 ether = 1.5 ether, and
    ///   totalSlash = 1.5 + 0.1 = 1.6 ether.
    /// If we register the provider with a stake less than bidderPortion (e.g. 1.4 ether),
    /// the branch sets bidderPortion to the entire stake and penaltyFee to 0.
    function test_SlashInsufficientFunds_Branch1() public {
        providerRegistry.setPreconfManager(address(this));

        address testProvider = vm.addr(400);
        vm.deal(testProvider, 3 ether);
        vm.prank(testProvider);
        providerRegistry.registerAndStake{value: 1.4 ether}();

        address bidder = vm.addr(401);
        uint256 initialBidderBalance = bidder.balance;

        uint256 feePercentValue = providerRegistry.feePercent();
        uint256 residualAmt = 1 ether;
        uint256 expectedInitialPenaltyFee = (residualAmt * feePercentValue) /
            providerRegistry.ONE_HUNDRED_PERCENT(); // ≈0.1 ether
        uint256 slashAmt = 0.5 ether;

        vm.startPrank(address(this));
        vm.expectEmit(true, true, true, true);
        emit InsufficientFundsToSlash(
            testProvider,
            1.4 ether,
            residualAmt,
            expectedInitialPenaltyFee,
            slashAmt
        );
        providerRegistry.slash(
            1 ether, // amt
            slashAmt, // slashAmt = 0.5 ether
            testProvider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.stopPrank();

        // In this branch, the code sets:
        //   bidderPortion = providerStake (1.4 ether),
        //   penaltyFee = 0, and totalSlash = 1.4 ether.
        // Thus, provider stake should drop to zero.
        assertEq(
            providerRegistry.providerStakes(testProvider),
            0,
            "Provider stake should be 0"
        );

        // The bidder receives the entire (adjusted) bidderPortion (1.4 ether).
        assertEq(
            bidder.balance - initialBidderBalance,
            1.4 ether,
            "Bidder should receive 1.4 ether"
        );

        // The accumulated penalty fee should remain unchanged (0).
        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            0,
            "Penalty fee should be 0"
        );
    }

    /// @dev Test the branch where providerStake is at least bidderPortion but still less than totalSlash.
    /// For inputs:
    ///   - amt = 1 ether,
    ///   - slashAmt = 0.5 ether, and
    ///   - residualBidPercentAfterDecay = ONE_HUNDRED_PERCENT,
    /// we get:
    ///   residualAmt = 1 ether,
    ///   penaltyFee = (1 ether * feePercent)/ONE_HUNDRED_PERCENT (≈0.1 ether),
    ///   bidderPortion = 1 ether + 0.5 ether = 1.5 ether,
    ///   totalSlash = 1.6 ether.
    /// By registering a provider with 1.55 ether (which is ≥ bidderPortion but less than totalSlash),
    /// the code subtracts bidderPortion (1.5 ether), then finds leftover = 0.05 ether (which is less than penaltyFee),
    /// so it sets penaltyFee = leftover (0.05 ether) and totalSlash = 1.5 + 0.05 = 1.55 ether.
    function test_SlashInsufficientFunds_Branch2() public {
        providerRegistry.setPreconfManager(address(this));

        address testProvider = vm.addr(410);
        vm.deal(testProvider, 3 ether);
        vm.prank(testProvider);
        providerRegistry.registerAndStake{value: 1.55 ether}();

        address bidder = vm.addr(411);
        uint256 initialBidderBalance = bidder.balance;

        uint256 feePercentValue = providerRegistry.feePercent();
        uint256 residualAmt = 1 ether;
        uint256 expectedInitialPenaltyFee = (residualAmt * feePercentValue) /
            providerRegistry.ONE_HUNDRED_PERCENT(); // ≈0.1 ether
        uint256 slashAmt = 0.5 ether;

        vm.startPrank(address(this));
        vm.expectEmit(true, true, true, true);
        emit InsufficientFundsToSlash(
            testProvider,
            1.55 ether,
            residualAmt,
            expectedInitialPenaltyFee,
            slashAmt
        );
        providerRegistry.slash(
            1 ether, // amt
            slashAmt, // slashAmt = 0.5 ether
            testProvider,
            payable(bidder),
            providerRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.stopPrank();

        // In this branch:
        //   leftover = providerStake - bidderPortion = 1.55 ether - 1.5 ether = 0.05 ether.
        //   Since 0.05 ether < penaltyFee (0.1 ether), penaltyFee is adjusted to 0.05 ether.
        //   TotalSlash becomes bidderPortion + penaltyFee = 1.5 ether + 0.05 ether = 1.55 ether.
        // Thus, the provider's stake is reduced to zero.
        assertEq(
            providerRegistry.providerStakes(testProvider),
            0,
            "Provider stake should be 0"
        );

        // The bidder receives bidderPortion (1.5 ether).
        assertEq(
            bidder.balance - initialBidderBalance,
            1.5 ether,
            "Bidder should receive 1.5 ether"
        );

        // The penalty fee tracker should increase by the adjusted penalty fee (0.05 ether).
        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            0.05 ether,
            "Penalty fee should be 0.05 ether"
        );
    }
}
