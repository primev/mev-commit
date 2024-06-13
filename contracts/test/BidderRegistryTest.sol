// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";

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
    event BidderRegistered(address indexed bidder, uint256 stakedAmount, uint256 windowNumber);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;
        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (address(this), blocksPerWindow))
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, (minStake, feeRecipient, feePercent, address(this), address(blockTracker), blocksPerWindow))
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        bidder = vm.addr(1);
        vm.deal(bidder, 1000 ether);
        vm.deal(address(this), 1000 ether);
    }

    function test_VerifyInitialContractState() public view {
        assertEq(bidderRegistry.minDeposit(), 1e18 wei);
        assertEq(bidderRegistry.feeRecipient(), feeRecipient);
        assertEq(bidderRegistry.feePercent(), feePercent);
        assertEq(bidderRegistry.preConfirmationsContract(), address(0));
        assertEq(bidderRegistry.bidderRegistered(bidder), false);
    }

    function testFail_BidderStakeAndRegisterMinStake() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        bidderRegistry.depositForSpecificWindow{value: 1 wei}(2);
    }

    function test_BidderStakeAndRegister() public {
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;

        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 1 ether, nextWindow);

        bidderRegistry.depositForSpecificWindow{value: 1 ether}(nextWindow);

        bool isBidderRegistered = bidderRegistry.bidderRegistered(bidder);
        assertEq(isBidderRegistered, true);

        uint256 bidderStakeStored = bidderRegistry.getDeposit(bidder, nextWindow);
        assertEq(bidderStakeStored, 1 ether);

        // For the second deposit, calculate the new next window
        currentWindow = blockTracker.getCurrentWindow();
        nextWindow = currentWindow + 1;

        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 2 ether, nextWindow);

        bidderRegistry.depositForSpecificWindow{value: 1 ether}(nextWindow);

        uint256 bidderStakeStored2 = bidderRegistry.getDeposit(bidder, nextWindow);
        assertEq(bidderStakeStored2, 2 ether);
    }

    function testFail_BidderStakeAndRegisterAlreadyRegistered() public {
        vm.prank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 2e18 wei}(2);
        vm.expectRevert(bytes(""));
        bidderRegistry.depositForSpecificWindow{value: 1 wei}(2);
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
        bidderRegistry.depositForSpecificWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.OpenBid(bidID, 1 ether, bidder, blockNumber);

        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.feeRecipientAmount();

        assertEq(providerAmount, 900000000000000000);
        assertEq(feeRecipientAmount, 100000000000000000);
        assertEq(bidderRegistry.getFeeRecipientAmount(), 100000000000000000);
        
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
        bidderRegistry.depositForSpecificWindow{value: 64 ether}(nextWindow);

        address provider = vm.addr(4);
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        bytes32 bidID = keccak256("1234");
        bidderRegistry.OpenBid(bidID, 1 ether, bidder, blockNumber);
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
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(nextWindow);
        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        bytes32 bidID = keccak256("1234");
        bidderRegistry.OpenBid(bidID, 1 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
    }

    function testFail_shouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        bidderRegistry.setPreconfirmationsContract(address(this));

        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        uint64 blockNumber = 66;
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(nextWindow);

        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));
        bytes32 bidID = keccak256("1234");
        bidderRegistry.OpenBid(bidID, 3 ether, bidder, blockNumber);
        bidderRegistry.retrieveFunds(nextWindow, bidID, payable(provider),100);
    }

    function test_withdrawFeeRecipientAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        vm.prank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 64 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.OpenBid(bidID, 1 ether, bidder, blockNumber);
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
        bidderRegistry.depositForSpecificWindow{value: 128 ether}(nextWindow);
        address provider = vm.addr(4);
        uint256 balanceBefore = address(provider).balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.OpenBid(bidID, 2 ether, bidder, blockNumber);
        
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
        bidderRegistry.depositForSpecificWindow{value: 5 ether}(nextWindow);
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
        bidderRegistry.depositForSpecificWindow{value: 128 ether}(nextWindow);
        uint256 balanceBefore = address(bidder).balance;
        bytes32 bidID = keccak256("1234");
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");

        bidderRegistry.OpenBid(bidID, 2 ether, bidder, blockNumber);
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
        bidderRegistry.depositForSpecificWindow{value: 5 ether}(nextWindow);
        vm.prank(bidderRegistry.owner());
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
    }

    function testFail_withdrawProtocolFeeNotOwner() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        vm.prank(bidder);
        uint256 currentWindow = blockTracker.getCurrentWindow();
        uint256 nextWindow = currentWindow + 1;
        bidderRegistry.depositForSpecificWindow{value: 5 ether}(nextWindow);
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
    }
}
