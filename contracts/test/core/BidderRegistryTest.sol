// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {IBidderRegistry} from "../../contracts/interfaces/IBidderRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../../contracts/utils/WindowFromBlockNumber.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {DepositManager} from "../../contracts/core/DepositManager.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract BidderRegistryTest is Test {
    uint256 public testNumber;
    BidderRegistry public bidderRegistry;
    uint256 public feePercent;
    uint256 public minStake;
    address public bidder;
    address public feeRecipient;
    uint256 public feePayoutPeriodMs;
    uint256 public bidderWithdrawalPeriodMs;
    BlockTracker public blockTracker;
    ProviderRegistry public providerRegistry;

    event BidderDeposited(address indexed bidder, address indexed provider, uint256 indexed depositedAmount);
    event WithdrawalRequested(address indexed bidder, address indexed provider, uint256 indexed withdrawalRequestTimestamp);
    event BidderWithdrawal(address indexed bidder, address indexed provider, uint256 indexed withdrawalAmount, uint256 escrowedAmount);

    event FeeTransfer(uint256 amount, address indexed recipient);
    event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient);
    event FeePayoutPeriodUpdated(uint256 indexed newFeePayoutPeriod);

    function setUp() public {
        testNumber = 42;
        feePercent = 10 * 1e16;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        feePayoutPeriodMs = 10000;
        bidderWithdrawalPeriodMs = 5000;
        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (address(this), address(this)))
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize,
            (feeRecipient,
            feePercent,
            address(this),
            address(blockTracker),
            feePayoutPeriodMs,
            bidderWithdrawalPeriodMs))
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        uint256 depositManagerMinBalance = 0.01 ether;
        DepositManager depositManager = new DepositManager(bidderRegistryProxy, depositManagerMinBalance);
        bidderRegistry.setDepositManagerImpl(address(depositManager));

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (minStake, feeRecipient, feePercent, address(this), 24 hours , 5 hours )
            )
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        vm.startPrank(address(this));
        blockTracker.setProviderRegistry(address(providerRegistry));
        vm.stopPrank();
        bidder = vm.addr(1);
        vm.deal(bidder, 1000 ether);
        vm.deal(address(this), 1000 ether);
    }

    function test_VerifyInitialContractState() public view {
        (address recipient, uint256 accumulatedAmount, uint256 lastPayoutTimestamp, uint256 payoutPeriodMs) = bidderRegistry.protocolFeeTracker();
        assertEq(recipient, feeRecipient);
        assertEq(payoutPeriodMs, feePayoutPeriodMs);
        assertEq(lastPayoutTimestamp, block.timestamp);
        assertEq(accumulatedAmount, 0);
        assertEq(bidderRegistry.feePercent(), feePercent);
        assertEq(bidderRegistry.preconfManager(), address(0));
        assertEq(bidderRegistry.bidderWithdrawalPeriodMs(), bidderWithdrawalPeriodMs);
    }

    function test_BidderStakeAndRegister() public {
        address provider1 = vm.addr(2);

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);
        emit BidderDeposited(bidder, provider1, 1 ether);
        bidderRegistry.depositAsBidder{value: 1 ether}(provider1);
        uint256 bidderStakeStored = bidderRegistry.getDeposit(bidder, provider1);
        assertEq(bidderStakeStored, 1 ether);

        address provider2 = vm.addr(3);
        vm.expectEmit(true, false, false, true);
        emit BidderDeposited(bidder, provider2, 2 ether);
        bidderRegistry.depositAsBidder{value: 2 ether}(provider2);
        uint256 bidderStakeStored2 = bidderRegistry.getDeposit(bidder, provider2);
        assertEq(bidderStakeStored2, 2 ether);
    }

    function test_TwoDeposits() public {
        address provider1 = vm.addr(2);
        address provider2 = vm.addr(3);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 2e18 wei}(provider1);
        bidderRegistry.depositAsBidder{value: 1 wei}(provider2);
    }

    function test_BidderRegistryReceive() public {
        vm.prank(bidder);
        vm.expectRevert();
        payable(address(bidderRegistry)).transfer(1 wei);
    }

    function test_BidderRegistryFallback() public {
        vm.prank(bidder);
        bytes memory data = abi.encode(1, 2);
        (bool success, ) = address(bidderRegistry).call{value: 1 wei}(data);
        require(!success, "should revert");
    }

    function test_SetNewProtocolFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit ProtocolFeeRecipientUpdated(newRecipient);
        bidderRegistry.setNewProtocolFeeRecipient(newRecipient);
        (address recipient, , , ) = bidderRegistry.protocolFeeTracker();
        assertEq(recipient, newRecipient);
    }

    function test_RevertWhen_SetNewProtocolFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(vm.addr(1));
        vm.expectRevert();
        bidderRegistry.setNewProtocolFeeRecipient(newRecipient);
    }

    function test_SetNewFeePayoutPeriodBlocks() public {
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit FeePayoutPeriodUpdated(890);
        bidderRegistry.setNewFeePayoutPeriod(890);
        (, , , uint256 payoutPeriodBlocks) = bidderRegistry.protocolFeeTracker();
        assertEq(payoutPeriodBlocks, 890);
    }

    function test_RevertWhen_SetNewFeePayoutPeriodBlocks() public {
        vm.prank(vm.addr(1));
        vm.expectRevert();
        bidderRegistry.setNewFeePayoutPeriod(83424);
    }

    function test_SetNewFeePercent() public {
        vm.prank(address(this));
        bidderRegistry.setNewFeePercent(uint16(25));
        assertEq(bidderRegistry.feePercent(), uint16(25));
    }

    function test_RevertWhen_SetNewFeePercent() public {
        vm.prank(vm.addr(1));
        vm.expectRevert();
        bidderRegistry.setNewFeePercent(uint16(25));
    }

    function test_SetDepositManagerImpl() public {
        address depositManager = vm.addr(3);
        vm.expectEmit(true, true, true, true);
        emit IBidderRegistry.DepositManagerImplUpdated(depositManager);
        vm.prank(address(this));
        bidderRegistry.setDepositManagerImpl(depositManager);
        assertEq(bidderRegistry.depositManagerImpl(), depositManager);
        assertEq(bidderRegistry.depositManagerHash(), keccak256(abi.encodePacked(hex"ef0100", depositManager)));
    }

    function test_RevertWhen_SetDepositManagerImpl() public {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, vm.addr(888)));
        vm.prank(vm.addr(888));
        bidderRegistry.setDepositManagerImpl(vm.addr(3));
    }

    function test_OpenBid_RevertWhen_DepositManagerNotSet() public {
        vm.prank(address(this));
        bidderRegistry.setDepositManagerImpl(address(0));
        assertEq(bidderRegistry.depositManagerImpl(), address(0));
        address preconfManager = bidderRegistry.preconfManager();
        vm.expectRevert(abi.encodeWithSelector(IBidderRegistry.DepositManagerNotSet.selector));
        vm.prank(preconfManager);
        bidderRegistry.openBid(keccak256("1234"), 1 ether, bidder, vm.addr(4));
    }

    function test_SetPreconfManager() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        bidderRegistry.setPreconfManager(newPreConfContract);
        assertEq(bidderRegistry.preconfManager(), newPreConfContract);
    }

    function test_RevertWhen_SetPreconfManager() public {
        vm.prank(vm.addr(1));
        address newPreConfContract = vm.addr(3);
        vm.expectRevert();
        bidderRegistry.setPreconfManager(newPreConfContract);
    }

    function test_ConvertFundsToProviderReward() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));

        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 64 ether}(provider);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        uint256 depositBefore = bidderRegistry.getDeposit(bidder, provider);

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);

        uint256 depositAfter = bidderRegistry.getDeposit(bidder, provider);
        assertEq(depositAfter, depositBefore-1 ether, "deposit should be reduced by bid amount since no top-up happened");

        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 900000000000000000);
        assertEq(feeRecipientAmount, 100000000000000000);
        assertEq(bidderRegistry.getDeposit(bidder, provider), 63 ether);
    }

    function test_ConvertFundsToProviderRewardWithDecay() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));

        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 64 ether}(provider);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);

        uint256 bidderBalance = bidder.balance;

        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), 50 * bidderRegistry.PRECISION());
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 450000000000000000);
        assertEq(feeRecipientAmount, 50000000000000000);
        
        assertEq(bidder.balance, bidderBalance + 500000000000000000);
        assertEq(bidderRegistry.getDeposit(bidder, provider), 63 ether);
    }

    function test_shouldReturnFunds() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));

        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 64 ether}(provider);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        uint256 bidderBalance = bidder.balance;

        assertEq(bidderRegistry.getDeposit(bidder, provider), 64 ether);
        assertEq(bidderRegistry.getEscrowedAmount(bidder, provider), 0);

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);

        assertEq(bidderRegistry.getDeposit(bidder, provider), 63 ether);
        assertEq(bidderRegistry.getEscrowedAmount(bidder, provider), 1 ether);

        bidderRegistry.unlockFunds(provider, commitmentDigest);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 0);
        assertEq(feeRecipientAmount, 0);
        
        assertEq(bidder.balance, bidderBalance + 1 ether);
        assertEq(bidderRegistry.getDeposit(bidder, provider), 63 ether);
    }

    function test_RevertWhen_ConvertFundsToProviderRewardNotPreConf() public {
        vm.prank(bidder);
        address provider = vm.addr(4);
        bidderRegistry.depositAsBidder{value: 2 ether}(provider);
        bytes32 commitmentDigest = keccak256("1234");
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);
        vm.prank(vm.addr(1));
        uint256 residualBidAfterDecay = bidderRegistry.ONE_HUNDRED_PERCENT();
        vm.expectRevert();
        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), residualBidAfterDecay);
    }

    function test_withdrawProviderAmount() public {
        bidderRegistry.setPreconfManager(address(this));
        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 128 ether}(provider);
        uint256 balanceBefore = address(provider).balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(commitmentDigest, 2 ether, bidder, provider);
        
        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        bidderRegistry.withdrawProviderAmount(payable(provider));
        uint256 balanceAfter = address(provider).balance;
        assertEq(balanceAfter - balanceBefore, 1800000000000000000);
        assertEq(bidderRegistry.providerAmount(provider), 0);
    }

    function test_RevertWhen_WithdrawProviderAmount() public {
        bidderRegistry.setPreconfManager(address(this));
        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 5 ether}(provider);
        vm.expectRevert();
        bidderRegistry.withdrawProviderAmount(payable(provider));
    }

    function test_DepositAsBidder() public {
        address provider1 = vm.addr(2);
        address provider2 = vm.addr(3);
        address provider3 = vm.addr(4);
        uint256 depositAmount = 3 ether;

        address[] memory providers = new address[](3);
        providers[0] = provider1;
        providers[1] = provider2;
        providers[2] = provider3;

        vm.startPrank(bidder);
        for (uint256 i = 0; i < 3; ++i) {

            vm.expectEmit(true, false, false, true);
            emit BidderDeposited(bidder, providers[i], depositAmount / 3);
            bidderRegistry.depositAsBidder{value: depositAmount / 3}(providers[i]);

            uint256 lockedFunds = bidderRegistry.getDeposit(bidder, providers[i]);
            assertEq(lockedFunds, depositAmount / 3);

            uint256 escrowedFunds = bidderRegistry.getEscrowedAmount(bidder, providers[i]);
            assertEq(escrowedFunds, 0);
        }
    }

    function test_WithdrawAsBidder() public {
        address provider1 = vm.addr(2);
        address provider2 = vm.addr(3);
        address provider3 = vm.addr(4);
        uint256 depositAmount = 3 ether;

        address[] memory providers = new address[](3);
        providers[0] = provider1;
        providers[1] = provider2;
        providers[2] = provider3;

        vm.startPrank(bidder);
        for (uint16 i = 0; i < 3; ++i) {
            vm.expectEmit(true, false, false, true);
            emit BidderDeposited(bidder, providers[i], depositAmount / 3);
            bidderRegistry.depositAsBidder{value: depositAmount / 3}(providers[i]);

            uint256 lockedFunds = bidderRegistry.getDeposit(bidder, providers[i]);
            assertEq(lockedFunds, depositAmount / 3);

            uint256 escrowedFunds = bidderRegistry.getEscrowedAmount(bidder, providers[i]);
            assertEq(escrowedFunds, 0);
        }
        vm.stopPrank();
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW*3 + 2);
        blockTracker.recordL1Block(blockNumber, "test");

        vm.expectEmit(true, false, false, true);
        emit WithdrawalRequested(bidder, provider1, block.timestamp);
        vm.expectEmit(true, false, false, true);
        emit WithdrawalRequested(bidder, provider2, block.timestamp);
        vm.expectEmit(true, false, false, true);
        emit WithdrawalRequested(bidder, provider3, block.timestamp);
        vm.prank(bidder);
        bidderRegistry.requestWithdrawalsAsBidder(providers);

        vm.warp(block.timestamp + bidderRegistry.bidderWithdrawalPeriodMs() + 1);

        vm.expectEmit(true, false, false, true);
        emit BidderWithdrawal(bidder, provider1, 1 ether, 0);
        vm.expectEmit(true, false, false, true);
        emit BidderWithdrawal(bidder, provider2, 1 ether, 0);
        vm.expectEmit(true, false, false, true);
        emit BidderWithdrawal(bidder, provider3, 1 ether, 0);
        vm.prank(bidder);
        bidderRegistry.withdrawAsBidder(providers);

        for (uint16 i = 0; i < 3; ++i) {
            uint256 lockedFunds = bidderRegistry.getDeposit(bidder, providers[i]);
            assertEq(lockedFunds, 0);

            uint256 escrowedFunds = bidderRegistry.getEscrowedAmount(bidder, providers[i]);
            assertEq(escrowedFunds, 0);
        }
    }

    function test_OpenBid_TransferExcessBid() public {
        bytes32 commitmentDigest = keccak256("commitment");
        uint256 bidAmt = 5 ether;
        address testBidder = vm.addr(2);
        
        vm.deal(testBidder, 10 ether);

        bidderRegistry.setPreconfManager(address(this));
        
        address provider = vm.addr(4);
        vm.prank(testBidder);
        bidderRegistry.depositAsBidder{value: 4 ether}(provider);
        
        uint256 maxBidAmt = bidderRegistry.getDeposit(testBidder, provider);
        uint256 usedAmount = bidderRegistry.getEscrowedAmount(testBidder, provider);
        uint256 availableAmount = maxBidAmt > usedAmount ? maxBidAmt - usedAmount : 0;

        assertEq(availableAmount, 4 ether);

        // open a bid that exceeds the available amount
        assertEq(bidAmt, 5 ether);
        vm.prank(address(this));
        bidderRegistry.openBid(commitmentDigest, bidAmt, testBidder, provider);

        uint256 depositAfter = bidderRegistry.getDeposit(testBidder, provider);
        assertEq(depositAfter, 0, "remaining deposit should be 0, since no top-up happened");
        
        (address storedBidder, uint256 storedBidAmt, IBidderRegistry.State storedState) = bidderRegistry.bidPayment(commitmentDigest);
        assertEq(storedBidder, testBidder);
        assertEq(storedBidAmt, 4 ether);
        assertEq(uint(storedState), uint(IBidderRegistry.State.PreConfirmed));
    }

    function test_ProtocolFeePayout() public {
        (, , uint256 lastPayoutTimestamp,) = bidderRegistry.protocolFeeTracker();
        uint256 defaultStartTimestamp = 1;
        assertEq(lastPayoutTimestamp, 1);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 0);
        vm.warp(defaultStartTimestamp + 10000 + 1); // roll past protocol fee payout period

        bidderRegistry.setPreconfManager(address(this));
        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 64 ether}(provider);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);
        vm.expectEmit(true, true, true, true);
        emit FeeTransfer(100000000000000000, feeRecipient);
        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 balanceAfter = feeRecipient.balance;
        assertEq(balanceAfter - balanceBefore, 100000000000000000);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 0);
    }

    function test_ProtocolFeeAccumulation() public {
        bidderRegistry.setPreconfManager(address(this));
        address provider = vm.addr(4);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 64 ether}(provider);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, provider);
        bidderRegistry.convertFundsToProviderReward(commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 balanceAfter = feeRecipient.balance;
        assertEq(balanceAfter - balanceBefore, 0);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 100000000000000000);
    }

    // Regression test for https://cantina.xyz/code/4ee8716d-3e0e-4f59-b90d-aa56bf3b484c/findings/8
    function test_OpenBidWithExcessExploit() public {
        address aliceBidder = vm.addr(7);
        address bobBidder = vm.addr(8);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 1);

        vm.deal(aliceBidder, 10 ether);
        vm.deal(bobBidder, 10 ether);

        bidderRegistry.setPreconfManager(address(this));

        address provider = vm.addr(4);
        vm.prank(aliceBidder);
        bidderRegistry.depositAsBidder{value: 2 ether}(provider);
        vm.prank(bobBidder);
        bidderRegistry.depositAsBidder{value: 2 ether}(provider);

        uint256 aliceBalanceBefore = aliceBidder.balance;
        uint256 bobBalanceBefore = bobBidder.balance;
        uint256 registryBalanceBefore = address(bidderRegistry).balance;

        assertEq(aliceBalanceBefore, 8 ether, "Alice balance BEFORE");
        assertEq(bobBalanceBefore, 8 ether, "Bob balance BEFORE");
        assertEq(registryBalanceBefore, 4 ether, "BidderRegistry balance BEFORE");

        assertEq(bidderRegistry.getDeposit(aliceBidder, provider), 2 ether);
        assertEq(bidderRegistry.getDeposit(bobBidder, provider), 2 ether);
        assertEq(bidderRegistry.getEscrowedAmount(aliceBidder, provider), 0);
        assertEq(bidderRegistry.getEscrowedAmount(bobBidder, provider), 0);

        uint256 maxBid = bidderRegistry.getDeposit(aliceBidder, provider);
        vm.startPrank(address(this));
        bidderRegistry.openBid(
            keccak256("commitment1"),
            maxBid,
            aliceBidder,
            provider
        );

        assertEq(bidderRegistry.getDeposit(aliceBidder, provider), 0); // no top-up happened
        assertEq(bidderRegistry.getEscrowedAmount(aliceBidder, provider), 2 ether);
        assertEq(bidderRegistry.getDeposit(bobBidder, provider), 2 ether);
        assertEq(bidderRegistry.getEscrowedAmount(bobBidder, provider), 0);

        bidderRegistry.openBid(
            keccak256("commitment2"),
            maxBid,
            aliceBidder,
            provider
        );
        bidderRegistry.openBid(
            keccak256("commitment3"),
            maxBid,
            aliceBidder,
            provider
        );
        vm.stopPrank();

        blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW * 2 + 1);
        blockTracker.recordL1Block(blockNumber, "test");
        address[] memory providers = new address[](1);
        providers[0] = provider;
        vm.prank(aliceBidder);
        bidderRegistry.requestWithdrawalsAsBidder(providers);
        vm.warp(block.timestamp + bidderRegistry.bidderWithdrawalPeriodMs()+1);
        vm.prank(aliceBidder);
        bidderRegistry.withdrawAsBidder(providers);

        uint256 aliceBalanceAfter = aliceBidder.balance;
        uint256 bobBalanceAfter = bobBidder.balance;
        uint256 registryBalanceAfter = address(bidderRegistry).balance;

        assertEq(aliceBalanceAfter, 8 ether, "Alice balance AFTER");
        assertEq(bobBalanceAfter, 8 ether, "Bob balance AFTER");
        assertEq(registryBalanceAfter, 4 ether, "BidderRegistry balance AFTER");

        assertEq(bidderRegistry.getDeposit(aliceBidder, provider), 0);
        assertEq(bidderRegistry.getEscrowedAmount(aliceBidder, provider), 2 ether);
        assertEq(bidderRegistry.getDeposit(bobBidder, provider), 2 ether);
        assertEq(bidderRegistry.getEscrowedAmount(bobBidder, provider), 0 ether);
    }

    function test_RevertWhen_DepositAsBidder_ZeroAmount() public {
        address provider = vm.addr(4);
        vm.startPrank(bidder);
        vm.expectRevert(IBidderRegistry.DepositAmountIsZero.selector);
        bidderRegistry.depositAsBidder{value: 0 ether}(provider);   
    }

    function test_OpenBid_NoTopUp_WrongCodeHash() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = payable(vm.addr(alicePK));
        vm.deal(alice, 10 ether);

        IncorrectBidderContract incorrectContract = new IncorrectBidderContract();
        vm.signAndAttachDelegation(address(incorrectContract), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 2 ether}(bob);
        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, depositBefore-1 ether, "deposit should be reduced by 1 ether, since no top-up happened");

        assertEq(alice.codehash, keccak256(abi.encodePacked(hex"ef0100", address(incorrectContract))));
        assertEq(alice.code.length, 23);
    }

    function test_OpenBid_NoTopUp_WithdrawalRequestExists() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = payable(vm.addr(alicePK));
        vm.deal(alice, 10 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 2 ether}(bob);

        address[] memory providers = new address[](1);
        providers[0] = bob;
        vm.prank(alice);
        bidderRegistry.requestWithdrawalsAsBidder(providers);

        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        vm.expectEmit(true, true, true, true);
        emit DepositManager.WithdrawalRequestExists(bob);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, depositBefore-1 ether, "deposit should be reduced by 1 ether, since no top-up happened");
        assertEq(bidderRegistry.getEscrowedAmount(alice, bob), 1 ether);
    }

    function test_OpenBid_NoTopUp_TargetDepositDoesNotExist() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = payable(vm.addr(alicePK));
        vm.deal(alice, 10 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 2 ether}(bob);
        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        vm.expectEmit(true, true, true, true);
        emit DepositManager.TargetDepositDoesNotExist(bob);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, depositBefore-1 ether, "deposit should be reduced by 1 ether, since no top-up happened");
    }

    function test_OpenBid_NoTopUp_CurrentDepositIsSufficient() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = vm.addr(alicePK);
        vm.deal(alice, 10 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);

        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(bob, 0.5 ether);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 2 ether}(bob);
        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositBefore, 2 ether, "deposit should be 2 ether");

        vm.expectEmit(true, true, true, true);
        emit DepositManager.CurrentDepositIsSufficient(bob, 1 ether, 0.5 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, 1 ether, "deposit should be reduced by 1 ether, since no top-up happened");
        assertEq(DepositManager(payable(alice)).targetDeposits(bob), 0.5 ether);

        vm.expectEmit(true, true, true, true);
        emit DepositManager.CurrentDepositIsSufficient(bob, 0.5 ether, 0.5 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment2"), 0.5 ether, alice, bob);

        uint256 depositAfter2 = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter2, 0.5 ether, "deposit should be 0.5 ether");
        assertEq(DepositManager(payable(alice)).targetDeposits(bob), 0.5 ether);
    }

    function test_OpenBid_NoTopUp_CurrentBalanceAtOrBelowMin() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = vm.addr(alicePK);
        vm.deal(alice, 1.01 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(bob, 1 ether);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 1 ether}(bob);
        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositBefore, 1 ether, "deposit should be 1 ether");

        vm.expectEmit(true, true, true, true);
        emit DepositManager.CurrentBalanceAtOrBelowMin(bob, 0.01 ether, 0.01 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, 0 ether, "deposit should be 0 ether, since no top-up happened");
        assertEq(bidderRegistry.getEscrowedAmount(alice, bob), 1 ether);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 0.001 ether}(bob);

        assertEq(alice.balance, 0.009 ether);
        vm.expectEmit(true, true, true, true);
        emit DepositManager.CurrentBalanceAtOrBelowMin(bob, 0.009 ether, 0.01 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment2"), 0.001 ether, alice, bob);

        uint256 depositAfter2 = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter2, 0 ether, "deposit should be 0 ether, since no top-up happened");
        assertEq(bidderRegistry.getEscrowedAmount(alice, bob), 1.001 ether);
    }

    function test_OpenBid_TopUpReduced() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = vm.addr(alicePK);
        vm.deal(alice, 2 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(bob, 1.5 ether);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 1.5 ether}(bob);
        uint256 depositBefore = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositBefore, 1.5 ether, "deposit should be 1.5 ether");
        assertEq(alice.balance, 0.5 ether, "alice should have 0.5 ether");
        assertEq(DepositManager(payable(alice)).minBalance(), 0.01 ether);

        vm.expectEmit(true, true, true, true);
        emit DepositManager.TopUpReduced(bob, 0.49 ether, 1.5 ether); // available = 0.5 - minBalance
        vm.expectEmit(true, true, true, true);
        emit DepositManager.DepositToppedUp(bob, 0.49 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1.5 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, 0.49 ether, "deposit should be 0.49 ether");
        assertEq(bidderRegistry.getEscrowedAmount(alice, bob), 1.5 ether);
        assertEq(DepositManager(payable(alice)).targetDeposits(bob), 1.5 ether);

        vm.deal(alice, 10 ether);
        assertEq(alice.balance, 10 ether, "alice should have 10 ether");

        vm.expectEmit(true, true, true, true);
        emit DepositManager.DepositToppedUp(bob, 1.5 ether); // 0.49 - 0.49 + full top-up amount
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment2"), 0.49 ether, alice, bob);
        uint256 depositAfterPart2 = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfterPart2, 1.5 ether, "deposit after part 2 should be 1.5 ether since that's the full target deposit");

        assertEq(bidderRegistry.getEscrowedAmount(alice, bob), 1.5 ether + 0.49 ether);
        assertEq(alice.balance, 8.50 ether);
    }

    function test_OpenBid_NormalTopUp() public {
        uint256 alicePK = uint256(0xA11CE);
        address alice = vm.addr(alicePK);
        vm.deal(alice, 2 ether);
        vm.signAndAttachDelegation(address(bidderRegistry.depositManagerImpl()), alicePK);

        address bob = vm.addr(8);
        vm.prank(alice);
        DepositManager(payable(alice)).setTargetDeposit(bob, 1.5 ether);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 1.5 ether}(bob);

        assertEq(alice.balance, 0.5 ether, "alice should have 0.5 ether");

        vm.expectEmit(true, true, true, true);
        emit DepositManager.DepositToppedUp(bob, 0.25 ether);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment3"), 0.25 ether, alice, bob);

        uint256 depositAfter = bidderRegistry.getDeposit(alice, bob);
        assertEq(depositAfter, 1.5 ether, "deposit should be 1.5 ether");
        assertEq(alice.balance, 0.25 ether, "alice should have 0.25 ether");
    }

    function test_OpenBid_GracefulTopUpFailure() public {
        AlwaysRevertsDepositManager alwaysRevertsDepositManager = new AlwaysRevertsDepositManager();
        vm.prank(bidderRegistry.owner());
        bidderRegistry.setDepositManagerImpl(address(alwaysRevertsDepositManager));

        uint256 alicePK = uint256(0xA11CE);
        address alice = vm.addr(alicePK);
        vm.deal(alice, 10 ether);
        vm.signAndAttachDelegation(address(alwaysRevertsDepositManager), alicePK);

        address provider = vm.addr(4);

        vm.prank(alice);
        bidderRegistry.depositAsBidder{value: 1 ether}(provider);

        vm.expectEmit(true, true, true, true);
        emit IBidderRegistry.TopUpFailed(alice, provider);
        vm.prank(bidderRegistry.preconfManager());
        bidderRegistry.openBid(keccak256("commitment"), 1 ether, alice, provider);
    }
}

contract IncorrectBidderContract {
    event somethingBadHappened(address provider);
    function topUpDeposit(address provider) public {
        emit somethingBadHappened(provider);
        revert("control flow should not reach here");
    }
}

contract AlwaysRevertsDepositManager {
    error RevertErr(address provider);
    function topUpDeposit(address provider) public pure {
        revert RevertErr(provider);
    }
}
