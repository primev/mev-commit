// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";

contract BidderRegistryTest is Test {
    uint256 testNumber;
    BidderRegistry internal bidderRegistry;
    uint16 internal feePercent;
    uint256 internal minStake;
    address internal bidder;
    address internal feeRecipient;

    /// @dev Event emitted when a bidder is registered with their staked amount
    event BidderRegistered(address indexed bidder, uint256 stakedAmount);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);

        bidderRegistry = new BidderRegistry(minStake, feeRecipient, feePercent, address(this));

        bidder = vm.addr(1);
        vm.deal(bidder, 100 ether);
        vm.deal(address(this), 100 ether);
    }

    function test_VerifyInitialContractState() public {
        assertEq(bidderRegistry.minAllowance(), 1e18 wei);
        assertEq(bidderRegistry.feeRecipient(), feeRecipient);
        assertEq(bidderRegistry.feePercent(), feePercent);
        assertEq(bidderRegistry.preConfirmationsContract(), address(0));
        assertEq(bidderRegistry.bidderRegistered(bidder), false);
    }

    function testFail_BidderStakeAndRegisterMinStake() public {
        vm.prank(bidder);
        vm.expectRevert(bytes(""));
        bidderRegistry.prepay{value: 1 wei}();
    }

    function test_BidderStakeAndRegister() public {
        vm.startPrank(bidder);
        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 1 ether);

        bidderRegistry.prepay{value: 1 ether}();

        bool isBidderRegistered = bidderRegistry.bidderRegistered(bidder);
        assertEq(isBidderRegistered, true);

        uint256 bidderStakeStored = bidderRegistry.getAllowance(bidder);
        assertEq(bidderStakeStored, 1 ether);

        vm.expectEmit(true, false, false, true);

        emit BidderRegistered(bidder, 2 ether);

        bidderRegistry.prepay{value: 1 ether}();

        uint256 bidderStakeStored2 = bidderRegistry.getAllowance(bidder);
        assertEq(bidderStakeStored2, 2 ether);

    }

    function testFail_BidderStakeAndRegisterAlreadyRegistered() public {
        vm.prank(bidder);
        bidderRegistry.prepay{value: 2e18 wei}();
        vm.expectRevert(bytes(""));
        bidderRegistry.prepay{value: 1 wei}();
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
        vm.prank(bidder);
        bidderRegistry.prepay{value: 2 ether}();
        address provider = vm.addr(4);
        bidderRegistry.LockBidFunds(bidID, 1 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider),100);
        uint256 providerAmount = bidderRegistry.providerAmount(provider);
        uint256 feeRecipientAmount = bidderRegistry.feeRecipientAmount();

        assertEq(providerAmount, 900000000000000000);
        assertEq(feeRecipientAmount, 100000000000000000);
        assertEq(bidderRegistry.getFeeRecipientAmount(), 100000000000000000);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 1 ether);
    }

    function test_shouldRetrieveFundsWithoutFeeRecipient() public {
        vm.prank(address(this));
        uint256 feerecipientValueBefore = bidderRegistry.feeRecipientAmount();

        bidderRegistry.setNewFeeRecipient(address(0));
        bidderRegistry.setPreconfirmationsContract(address(this));

        vm.prank(bidder);
        bidderRegistry.prepay{value: 2 ether}();
        address provider = vm.addr(4);
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 1 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider),100);

        uint256 feerecipientValueAfter = bidderRegistry.feeRecipientAmount();
        uint256 providerAmount = bidderRegistry.providerAmount(provider);

        assertEq(providerAmount, 900000000000000000);
        assertEq(feerecipientValueAfter, feerecipientValueBefore);

        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 1 ether);
    }

    function testFail_shouldRetrieveFundsNotPreConf() public {
        vm.prank(bidder);
        bidderRegistry.prepay{value: 2 ether}();
        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 1 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider),100);
    }

    function testFail_shouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        bidderRegistry.setPreconfirmationsContract(address(this));

        vm.prank(bidder);
        bidderRegistry.prepay{value: 2 ether}();

        address provider = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 3 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider),100);
    }

    function test_withdrawFeeRecipientAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 2 ether}();
        address provider = vm.addr(4);
        uint256 balanceBefore = feeRecipient.balance;
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 1 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider),100);
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
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        address provider = vm.addr(4);
        uint256 balanceBefore = address(provider).balance;
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 2 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider), 100);
        bidderRegistry.withdrawProviderAmount(payable(provider));
        uint256 balanceAfter = address(provider).balance;
        assertEq(balanceAfter - balanceBefore, 1800000000000000000);
        assertEq(bidderRegistry.providerAmount(provider), 0);
    }

    function testFail_withdrawProviderAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        address provider = vm.addr(4);
        bidderRegistry.withdrawProviderAmount(payable(provider));
    }

    function test_withdrawStakedAmount() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        uint256 balanceBefore = address(bidder).balance;
        vm.prank(bidder);
        bidderRegistry.withdrawPrepaidAmount(payable(bidder));
        uint256 balanceAfter = address(bidder).balance;
        assertEq(balanceAfter - balanceBefore, 5 ether);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 0);
    }

    function testFail_withdrawStakedAmountNotOwner() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        bidderRegistry.withdrawPrepaidAmount(payable(bidder));
    }

    function testFail_withdrawStakedAmountStakeZero() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        vm.prank(bidder);
        bidderRegistry.withdrawPrepaidAmount(payable(bidder));
    }

    function test_withdrawProtocolFee() public {
        address provider = vm.addr(4);
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        uint256 balanceBefore = address(bidder).balance;
        bytes32 bidID = keccak256("1234");
        bidderRegistry.LockBidFunds(bidID, 2 ether, bidder);
        bidderRegistry.retrieveFunds(bidID, payable(provider), 100);
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
        bidderRegistry.prepay{value: 5 ether}();
        vm.prank(bidderRegistry.owner());
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
    }

    function testFail_withdrawProtocolFeeNotOwner() public {
        bidderRegistry.setPreconfirmationsContract(address(this));
        bidderRegistry.setNewFeeRecipient(address(0));
        vm.prank(bidder);
        bidderRegistry.prepay{value: 5 ether}();
        bidderRegistry.withdrawProtocolFee(payable(address(bidder)));
    }
}
