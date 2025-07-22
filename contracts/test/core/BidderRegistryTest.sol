// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

import {Test} from "forge-std/Test.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {IBidderRegistry} from "../../contracts/interfaces/IBidderRegistry.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../../contracts/utils/WindowFromBlockNumber.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";

contract BidderRegistryTest is Test {
    uint256 public testNumber;
    BidderRegistry public bidderRegistry;
    uint256 public feePercent;
    uint256 public minStake;
    address public bidder;
    address public feeRecipient;
    uint256 public feePayoutPeriodMs;
    BlockTracker public blockTracker;
    ProviderRegistry public providerRegistry;

    /// @dev Event emitted when a bidder is registered with their staked amount
    event BidderRegistered(address indexed bidder, uint256 indexed stakedAmount, uint256 indexed windowNumber);

    event FeeTransfer(uint256 amount, address indexed recipient);
    event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient);
    event FeePayoutPeriodUpdated(uint256 indexed newFeePayoutPeriod);

    function setUp() public {
        testNumber = 42;
        feePercent = 10 * 1e16;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        feePayoutPeriodMs = 10000;
        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (address(this), address(this)))
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, (feeRecipient, feePercent, address(this), address(blockTracker), feePayoutPeriodMs))
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

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
    }

    function test_BidderStakeAndRegister() public {
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 1 ether, nextWindow);

        bidderRegistry.depositForWindow{value: 1 ether}(nextWindow);

        uint256 bidderStakeStored = bidderRegistry.getDeposit(bidder, nextWindow);
        assertEq(bidderStakeStored, 1 ether);

        // For the second deposit, calculate the new next window
        currentWindow = blockTracker.getCurrentWindow();
        nextWindow = currentWindow + 1;

        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 2 ether, nextWindow);

        bidderRegistry.depositForWindow{value: 1 ether}(nextWindow);

        uint256 bidderStakeStored2 = bidderRegistry.getDeposit(bidder, nextWindow);
        assertEq(bidderStakeStored2, 2 ether);
    }

    function testFail_BidderStakeAndRegisterAlreadyRegistered() public {
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 2e18 wei}(2);
        vm.expectRevert(bytes(""));
        bidderRegistry.depositForWindow{value: 1 wei}(2);
    }

    function testFail_Receive() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(bidderRegistry).call{value: 1 wei}("");
        require(success, "couldn't transfer to bidder");
    }

    function testFail_Fallback() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(bidderRegistry).call{value: 1 wei}("");
        require(success, "couldn't transfer to bidder");
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

    function testFail_SetNewProtocolFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert(bytes(""));
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

    function testFail_SetNewFeePayoutPeriodBlocks() public {
        vm.expectRevert(bytes(""));
        bidderRegistry.setNewFeePayoutPeriod(83424);
    }

    function test_SetNewFeePercent() public {
        vm.prank(address(this));
        bidderRegistry.setNewFeePercent(uint16(25));
        assertEq(bidderRegistry.feePercent(), uint16(25));
    }

    function testFail_SetNewFeePercent() public {
        vm.expectRevert(bytes(""));
        bidderRegistry.setNewFeePercent(uint16(25));
    }

    function test_SetPreconfManager() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        bidderRegistry.setPreconfManager(newPreConfContract);
        assertEq(bidderRegistry.preconfManager(), newPreConfContract);
    }

    function testFail_SetPreconfManager() public {
        vm.prank(address(this));
        vm.expectRevert(bytes(""));
        bidderRegistry.setPreconfManager(address(0));
    }

    function test_shouldRetrieveFunds() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);

        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 900000000000000000);
        assertEq(feeRecipientAmount, 100000000000000000);
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function test_shouldRetrieveFundsWithDecay() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);

        uint256 bidderBalance = bidder.balance;

        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), 50 * bidderRegistry.PRECISION());
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 450000000000000000);
        assertEq(feeRecipientAmount, 50000000000000000);
        
        assertEq(bidder.balance, bidderBalance + 500000000000000000);
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function test_shouldReturnFunds() public {
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        uint256 bidderBalance = bidder.balance;

        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);

        bidderRegistry.unlockFunds(nextWindow, commitmentDigest);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.getAccumulatedProtocolFee();

        assertEq(providerAmount, 0);
        assertEq(feeRecipientAmount, 0);
        
        assertEq(bidder.balance, bidderBalance + 1 ether);
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function testFail_shouldRetrieveFundsNotPreConf() public {
        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        uint64 blockNumber = 66;
        bidderRegistry.depositForWindow{value: 2 ether}(nextWindow);
        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
    }

    function testFail_shouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        bidderRegistry.setPreconfManager(address(this));

        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        uint64 blockNumber = 66;
        bidderRegistry.depositForWindow{value: 2 ether}(nextWindow);

        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));
        bytes32 commitmentDigest = keccak256("1234");
        bidderRegistry.openBid(commitmentDigest, 3 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
    }

    function test_withdrawProviderAmount() public {
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 128 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = address(provider).balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(commitmentDigest, 2 ether, bidder, blockNumber);
        
        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        bidderRegistry.withdrawProviderAmount(payable(provider));
        uint256 balanceAfter = address(provider).balance;
        assertEq(balanceAfter - balanceBefore, 1800000000000000000);
        assertEq(bidderRegistry.providerAmount(provider), 0);
    }

    function testFail_withdrawProviderAmount() public {
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 5 ether}(nextWindow);
        address provider = vm.addr(4);
        bidderRegistry.withdrawProviderAmount(payable(provider));
    }

    function test_DepositForWindows() public {
        uint256[] memory windows = new uint256[](3);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        for (uint256 i = 0; i < windows.length; ++i) {
            windows[i] = currentWindow + i;
        }
        uint256 depositAmount = 3 ether;

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);
        for (uint256 i = 0; i < windows.length; ++i) {
            emit BidderRegistered(bidder, depositAmount / windows.length, windows[i]);
        }

        bidderRegistry.depositForWindows{value: depositAmount}(windows);
        for (uint256 i = 0; i < windows.length; ++i) {
            uint256 lockedFunds = bidderRegistry.lockedFunds(bidder, windows[i]);
            assertEq(lockedFunds, depositAmount / windows.length);

            uint256 maxBidPerBlock = bidderRegistry.maxBidPerBlock(bidder, windows[i]);
            assertEq(maxBidPerBlock, depositAmount / (windows.length * WindowFromBlockNumber.BLOCKS_PER_WINDOW));
        }
    }

    function test_WithdrawFromWindows() public {
        uint256[] memory windows = new uint256[](3);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        for (uint256 i = 0; i < windows.length; ++i) {
            windows[i] = currentWindow + i;
        }
        uint256 depositAmount = minStake * windows.length;

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);
        for (uint16 i = 0; i < windows.length; ++i) {
            emit BidderRegistered(bidder, depositAmount / windows.length, currentWindow + i);
        }

        bidderRegistry.depositForWindows{value: depositAmount}(windows);

        for (uint16 i = 0; i < windows.length; ++i) {
            uint256 lockedFunds = bidderRegistry.lockedFunds(bidder, currentWindow + i);
            assertEq(lockedFunds, depositAmount / windows.length);

            uint256 maxBid = bidderRegistry.maxBidPerBlock(bidder, currentWindow + i);
            assertEq(maxBid, (depositAmount / windows.length) / WindowFromBlockNumber.BLOCKS_PER_WINDOW);
        }
        vm.stopPrank();
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW*3 + 2);
        blockTracker.recordL1Block(blockNumber, "test");

        vm.startPrank(bidder);
        bidderRegistry.withdrawFromWindows(windows);

        for (uint16 i = 0; i < windows.length; ++i) {
            uint256 lockedFunds = bidderRegistry.lockedFunds(bidder, currentWindow + i);
            assertEq(lockedFunds, 0);

            uint256 maxBid = bidderRegistry.maxBidPerBlock(bidder, currentWindow + i);
            assertEq(maxBid, 0);

            (uint256 startBlock, uint256 endBlock) = WindowFromBlockNumber.getBlockNumbersFromWindow(currentWindow + i);
            for (uint256 blockNum = startBlock; blockNum <= endBlock; ++blockNum) {
                uint256 usedFunds = bidderRegistry.usedFunds(bidder, uint64(blockNum));
                assertEq(usedFunds, 0);
            }
        }
    }

    function test_OpenBidtransferExcessBid() public {
        bytes32 commitmentDigest = keccak256("commitment");
        uint256 bidAmt = 3 ether;
        address testBidder = vm.addr(2);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 1);
        
        // Deal some ETH to the test bidder
        vm.deal(testBidder, 10 ether);

        // Simulate the pre-confirmations contract
        bidderRegistry.setPreconfManager(address(this));
        
        // Deposit some funds for the next window
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(testBidder);
        bidderRegistry.depositForWindow{value: 4 ether}(nextWindow);
        
        // Ensure the used amount is less than the max bid per block
        uint256 maxBidAmt = bidderRegistry.maxBidPerBlock(testBidder, nextWindow);
        uint256 usedAmount = bidderRegistry.usedFunds(testBidder, blockNumber);
        uint256 availableAmount = maxBidAmt > usedAmount ? maxBidAmt - usedAmount : 0;
        
        // Open a bid that exceeds the available amount
        vm.prank(address(this));
        bidderRegistry.openBid(commitmentDigest, bidAmt, testBidder, blockNumber);
        
        // Verify that the excess bid was transferred back to the test bidder
        uint256 expectedBidAmt = availableAmount;
        
        // Verify the bid state
        (address storedBidder, uint256 storedBidAmt, IBidderRegistry.State storedState) = bidderRegistry.bidPayment(commitmentDigest);
        assertEq(storedBidder, testBidder);
        assertEq(storedBidAmt, expectedBidAmt);
        assertEq(uint(storedState), uint(IBidderRegistry.State.PreConfirmed));
    }

    function test_ProtocolFeePayout() public {
        (, , uint256 lastPayoutTimestamp,) = bidderRegistry.protocolFeeTracker();
        uint256 defaultStartTimestamp = 1;
        assertEq(lastPayoutTimestamp, 1);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 0);
        vm.warp(defaultStartTimestamp + 10000 + 1); // roll past protocol fee payout period

        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);
        vm.expectEmit(true, true, true, true);
        emit FeeTransfer(100000000000000000, feeRecipient);
        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 balanceAfter = feeRecipient.balance;
        assertEq(balanceAfter - balanceBefore, 100000000000000000);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 0);
    }

    function test_ProtocolFeeAccumulation() public {
        bidderRegistry.setPreconfManager(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 commitmentDigest = keccak256("1234");
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bidderRegistry.openBid(commitmentDigest, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, commitmentDigest, payable(provider), bidderRegistry.ONE_HUNDRED_PERCENT());
        uint256 balanceAfter = feeRecipient.balance;
        assertEq(balanceAfter - balanceBefore, 0);
        assertEq(bidderRegistry.getAccumulatedProtocolFee(), 100000000000000000);
    }

    function test_OpenBidWithExcessExploit() public {
        address aliceBidder = vm.addr(7);
        address bodBidder = vm.addr(8);
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 1);

        //1)  Deal some ETH to the Alice and Bob
        vm.deal(aliceBidder, 10 ether);
        vm.deal(bodBidder, 10 ether);

        //2) Simulate the pre-confirmations contract
        bidderRegistry.setPreconfManager(address(this));

        //3) Deposit some funds for the next window
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(aliceBidder);
        bidderRegistry.depositForWindow{value: 2 ether}(nextWindow);
        vm.prank(bodBidder);
        bidderRegistry.depositForWindow{value: 2 ether}(nextWindow);

        // Capture balances before the exploit
        uint256 aliceBalanceBefore = aliceBidder.balance;
        uint256 bobBalanceBefore = bodBidder.balance;
        uint256 registryBalanceBefore = address(bidderRegistry).balance;

        // Expected balances before the exploit
        assertEq(aliceBalanceBefore, 8 ether, "Alice balance BEFORE");
        assertEq(bobBalanceBefore, 8 ether, "Bob balance BEFORE");
        assertEq(registryBalanceBefore, 4 ether, "BidderRegistry balance BEFORE");

        uint256 maxBid = bidderRegistry.maxBidPerBlock(aliceBidder, nextWindow);
        //4) Alice open bids at maxBid multiple times
        vm.startPrank(address(this));
        // First bid works fine. maxBid is depleted from lockedFunds
        bidderRegistry.openBid(
            keccak256("commitment1"),
            maxBid,
            aliceBidder,
            blockNumber
        );
        // Second bid start the stealing show. maxBid is being refunded (by the excess logic) and lockedFunds is NOT depleted.
        // This is effectively stealing from poor Bob.
        bidderRegistry.openBid(
            keccak256("commitment2"),
            maxBid,
            aliceBidder,
            blockNumber
        );
        // Thrid bid continue the stealing show, exactly behaving as the second bid.
        bidderRegistry.openBid(
            keccak256("commitment3"),
            maxBid,
            aliceBidder,
            blockNumber
        );

        // And Alice could do this until she fully drain the bidderRegistry contract, effectively stealing from all the bidders.
        vm.stopPrank();

        //5) Alice withdraw her locked funds (which will be intact minus maxBid as being spent in the first bid)
        blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW * 2 + 1);
        blockTracker.recordL1Block(blockNumber, "test");

        uint256[] memory windows = new uint256[](1);
        windows[0] = nextWindow;
        vm.prank(aliceBidder);
        bidderRegistry.withdrawFromWindows(windows);

        // Capture balances after the exploit
        uint256 aliceBalanceAfter = aliceBidder.balance;
        uint256 bobBalanceAfter = bodBidder.balance;
        uint256 registryBalanceAfter = address(bidderRegistry).balance;

        // Expected balances after the exploit
        assertEq(aliceBalanceAfter, 9.8 ether, "Alice balance AFTER");
        assertEq(bobBalanceAfter, 8 ether, "Bob balance AFTER");
        assertEq(registryBalanceAfter, 2.2 ether, "BidderRegistry balance AFTER");
    }

    function test_BidderStakeAndRegisteratZero() public {
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.startPrank(bidder);

        vm.expectRevert(IBidderRegistry.DepositAmountIsZero.selector);
        bidderRegistry.depositForWindow{value: 0 ether}(nextWindow);   
    }
}
