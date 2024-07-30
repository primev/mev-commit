// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";
import {IBidderRegistry} from "../contracts/interfaces/IBidderRegistry.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract BidderRegistryTest is Test {
    uint256 testNumber;
    BidderRegistry internal bidderRegistry;
    uint16 internal feePercent;
    uint256 internal minStake;
    address internal bidder;
    address internal feeRecipient;
    uint256 blocksPerWindow;
    BlockTracker internal blockTracker;

    /// @dev Event emitted when a bidder is registered with their staked amount
    event BidderRegistered(address indexed bidder, uint256 indexed stakedAmount, uint256 indexed windowNumber);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;
        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (blocksPerWindow, address(this), address(this)))
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, (feeRecipient, feePercent, address(this), address(blockTracker), blocksPerWindow))
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        bidder = vm.addr(1);
        vm.deal(bidder, 1000 ether);
        vm.deal(address(this), 1000 ether);
    }

    function test_VerifyInitialContractState() public view {
        assertEq(bidderRegistry.feeRecipient(), feeRecipient);
        assertEq(bidderRegistry.feePercent(), feePercent);
        assertEq(bidderRegistry.preConfirmationsContract(), address(0));
        assertEq(bidderRegistry.bidderRegistered(bidder), false);
    }

    function test_BidderStakeAndRegister() public {
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 1 ether, nextWindow);

        bidderRegistry.depositForWindow{value: 1 ether}(nextWindow);

        bool isBidderRegistered = bidderRegistry.bidderRegistered(bidder);
        assertEq(isBidderRegistered, true);

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

    function testFail_receive() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(bidderRegistry).call{value: 1 wei}("");
        require(success, "couldn't transfer to bidder");
    }

    function testFail_fallback() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(bidderRegistry).call{value: 1 wei}("");
        require(success, "couldn't transfer to bidder");
    }

    function test_SetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        bidderRegistry.setNewFeeRecipient(newRecipient);

        assertEq(bidderRegistry.feeRecipient(), newRecipient);
    }

    function testFail_SetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert(bytes(""));
        bidderRegistry.setNewFeeRecipient(newRecipient);
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

    function test_SetPreConfContract() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        bidderRegistry.setPreconfirmationsContract(newPreConfContract);

        assertEq(bidderRegistry.preConfirmationsContract(), newPreConfContract);
    }

    function testFail_SetPreConfContract() public {
        vm.prank(address(this));
        vm.expectRevert(bytes(""));
        bidderRegistry.setPreconfirmationsContract(address(0));
    }

    function test_shouldRetrieveFunds() public {
        bytes32 bidID = keccak256("1234");
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);

        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider), 100);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.feeRecipientAmount();

        assertEq(providerAmount, 900000000000000000);
        assertEq(feeRecipientAmount, 100000000000000000);
        assertEq(bidderRegistry.getFeeRecipientAmount(), 100000000000000000);
        
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function test_shouldRetrieveFundsWithDecay() public {
        bytes32 bidID = keccak256("1234");
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);

        uint256 bidderBalance = bidder.balance;

        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider), 50);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.feeRecipientAmount();

        assertEq(providerAmount, 450000000000000000);
        assertEq(feeRecipientAmount, 50000000000000000);
        assertEq(bidderRegistry.getFeeRecipientAmount(), 50000000000000000);
        
        assertEq(bidder.balance, bidderBalance + 500000000000000000);
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function test_shouldReturnFunds() public {
        bytes32 bidID = keccak256("1234");
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        uint256 bidderBalance = bidder.balance;

        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);

        bidderRegistry.unlockFunds(nextWindow, bidID);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.feeRecipientAmount();

        assertEq(providerAmount, 0);
        assertEq(feeRecipientAmount, 0);
        
        assertEq(bidder.balance, bidderBalance + 1 ether);
        assertEq(bidderRegistry.lockedFunds(bidder, nextWindow), 63 ether);
    }

    function test_shouldRetrieveFundsWithoutFeeRecipient() public {
        vm.prank(address(this));
        uint256 feerecipientValueBefore = bidderRegistry.feeRecipientAmount();

        bidderRegistry.setNewFeeRecipient(address(0));
        bidderRegistry.setPreconfirmationsContract(address(this));

        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);

        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bytes32 bidID = keccak256("1234");
        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider), 100);

        uint256 feerecipientValueAfter = bidderRegistry.feeRecipientAmount();
        uint256 providerAmount = bidderRegistry.providerAmount(provider);

        assertEq(providerAmount, 900000000000000000);
        assertEq(feerecipientValueAfter, feerecipientValueBefore);

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
        bytes32 bidID = keccak256("1234");
        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
    }

    function testFail_shouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        bidderRegistry.setPreconfirmationsContract(address(this));

        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        uint64 blockNumber = 66;
        bidderRegistry.depositForWindow{value: 2 ether}(nextWindow);

        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));
        bytes32 bidID = keccak256("1234");
        bidderRegistry.openBid(bidID, 3 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
    }

    function test_withdrawFeeRecipientAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(bidID, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
        bidderRegistry.withdrawFeeRecipientAmount();
        uint256 balanceAfter = feeRecipient.balance;
        assertEq(balanceAfter - balanceBefore, 100000000000000000);
        assertEq(bidderRegistry.feeRecipientAmount(), 0);
        assertEq(bidderRegistry.getFeeRecipientAmount(), 0);
    }

    function testFail_withdrawFeeRecipientAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.withdrawFeeRecipientAmount();
    }

    function test_withdrawProviderAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 128 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = address(provider).balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(bidID, 2 ether, bidder, blockNumber);
        
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider), 100);
        bidderRegistry.withdrawProviderAmount(payable(provider));
        uint256 balanceAfter = address(provider).balance;
        assertEq(balanceAfter - balanceBefore, 1800000000000000000);
        assertEq(bidderRegistry.providerAmount(provider), 0);
    }

    function testFail_withdrawProviderAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 5 ether}(nextWindow);
        address provider = vm.addr(4);
        bidderRegistry.withdrawProviderAmount(payable(provider));
    }

    function test_withdrawProtocolFee() public {
        address provider = vm.addr(4);
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 128 ether}(nextWindow);
        uint256 balanceBefore = address(bidder).balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.openBid(bidID, 2 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider), 100);
        vm.prank(bidderRegistry.owner());
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
        uint256 balanceAfter = address(bidder).balance;
        assertEq(balanceAfter - balanceBefore, 200000000000000000);
        assertEq(bidderRegistry.protocolFeeAmount(), 0);
    }

    function testFail_withdrawProtocolFee() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        bidderRegistry.depositForWindow{value: 5 ether}(nextWindow);
        vm.prank(bidderRegistry.owner());
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
    }

    function testFail_withdrawProtocolFeeNotOwner() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        bidderRegistry.depositForWindow{value: 5 ether}(nextWindow);
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
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
            assertEq(maxBidPerBlock, depositAmount / (windows.length * blocksPerWindow));
        }

        bool isBidderRegistered = bidderRegistry.bidderRegistered(bidder);
        assertEq(isBidderRegistered, true);
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
            assertEq(maxBid, (depositAmount / windows.length) / blocksPerWindow);
        }
        vm.stopPrank();
        uint64 blockNumber = uint64(blocksPerWindow*3 + 2);
        blockTracker.recordL1Block(blockNumber, "test");

        vm.startPrank(bidder);
        bidderRegistry.withdrawFromWindows(windows);

        for (uint16 i = 0; i < windows.length; ++i) {
            uint256 lockedFunds = bidderRegistry.lockedFunds(bidder, currentWindow + i);
            assertEq(lockedFunds, 0);

            uint256 maxBid = bidderRegistry.maxBidPerBlock(bidder, currentWindow + i);
            assertEq(maxBid, 0);
        }
    }

    function test_openBid_transferExcessBid() public {
        bytes32 commitmentDigest = keccak256("commitment");
        uint256 bid = 3 ether;
        address testBidder = vm.addr(2);
        uint64 blockNumber = uint64(blocksPerWindow + 1);
        
        // Deal some ETH to the test bidder
        vm.deal(testBidder, 10 ether);

        // Simulate the pre-confirmations contract
        bidderRegistry.setPreconfirmationsContract(address(this));
        
        // Deposit some funds for the next window
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(testBidder);
        bidderRegistry.depositForWindow{value: 4 ether}(nextWindow);
        
        // Ensure the used amount is less than the max bid per block
        uint256 maxBid = bidderRegistry.maxBidPerBlock(testBidder, nextWindow);
        uint256 usedAmount = bidderRegistry.usedFunds(testBidder, blockNumber);
        uint256 availableAmount = maxBid > usedAmount ? maxBid - usedAmount : 0;
        
        // Open a bid that exceeds the available amount
        vm.prank(address(this));
        bidderRegistry.openBid(commitmentDigest, bid, testBidder, blockNumber);
        
        // Verify that the excess bid was transferred back to the test bidder
        uint256 expectedBid = availableAmount;
        uint256 testBidderBalance = testBidder.balance;
        assertEq(testBidderBalance, 10 ether - 4 ether + (bid - expectedBid));
        
        // Verify the bid state
        (address storedBidder, uint256 storedBidAmt, IBidderRegistry.State storedState) = bidderRegistry.bidPayment(commitmentDigest);
        assertEq(storedBidder, testBidder);
        assertEq(storedBidAmt, expectedBid);
        assertEq(uint(storedState), uint(IBidderRegistry.State.PreConfirmed));
    }
}
