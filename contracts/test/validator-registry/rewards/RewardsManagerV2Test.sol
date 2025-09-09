// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {RewardsManagerV2} from "../../../contracts/validator-registry/rewards/RewardsManagerV2.sol";
import {IRewardsManagerV2} from "../../../contracts/interfaces/IRewardsManagerV2.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract RewardsManagerV2Test is Test {
    RewardsManagerV2 internal rewardsManager;

    address internal ownerAddress;
    address payable internal treasuryAddress;
    address internal payerOne;
    address internal payerTwo;
    address internal feeRecipientOne;
    address internal feeRecipientTwo;

    // Events mirrored from V2 (for expectEmit)
    event ProposerPaid(address indexed feeRecipient, uint256 indexed proposerAmt, uint256 indexed rewardAmt);
    event TreasuryWithdrawn(uint256 indexed treasuryAmt);
    event RewardsPctBpsSet(uint256 indexed rewardsPctBps);
    event TreasurySet(address indexed treasury);

    function setUp() public {
        ownerAddress = address(0xA11CE);
        treasuryAddress = payable(address(0x12345));
        payerOne = address(0xBEEF1);
        payerTwo = address(0xBEEF2);
        feeRecipientOne = address(0xFEE01);
        feeRecipientTwo = address(0xFEE02);

        vm.deal(payerOne, 100 ether);
        vm.deal(payerTwo, 100 ether);

        uint256 initialRewardsPctBps = 1500; // 15%

        RewardsManagerV2 implementation = new RewardsManagerV2();
        bytes memory initData = abi.encodeCall(
            RewardsManagerV2.initialize,
            (ownerAddress, initialRewardsPctBps, treasuryAddress)
        );

        address proxy = address(new ERC1967Proxy(address(implementation), initData));
        rewardsManager = RewardsManagerV2(payable(proxy));
    }
    
    // initialize
    function test_Initialize_setsOwnerBpsTreasury() public {
        address ownerAfterInit = rewardsManager.owner();
        assertEq(ownerAfterInit, ownerAddress);

        uint256 bpsAfterInit = rewardsManager.rewardsPctBps();
        assertEq(bpsAfterInit, 1500);

        address treasuryAfterInit = rewardsManager.treasury();
        assertEq(treasuryAfterInit, treasuryAddress);

        uint256 toTreasuryAfterInit = rewardsManager.toTreasury();
        assertEq(toTreasuryAfterInit, 0);
    }
    
    // owner-only setters and bounds
    function test_SetTreasury_onlyOwner_and_emits() public {
        address payable newTreasury = payable(address(0x56789));

        vm.prank(ownerAddress);
        vm.expectEmit();
        emit TreasurySet(newTreasury);
        rewardsManager.setTreasury(newTreasury);

        address treasuryAfterSet = rewardsManager.treasury();
        assertEq(treasuryAfterSet, newTreasury);

        vm.prank(payerOne);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, payerOne));
        rewardsManager.setTreasury(payable(address(123)));
    }

    // updating rewards pct
    function test_SetRewardsPctBps_onlyOwner_and_bounds() public {
        vm.prank(ownerAddress);
        vm.expectEmit();
        emit RewardsPctBpsSet(2000);
        rewardsManager.setRewardsPctBps(2000);

        uint256 bpsAfterUpdate = rewardsManager.rewardsPctBps();
        assertEq(bpsAfterUpdate, 2000);

        vm.prank(ownerAddress);
        vm.expectRevert(IRewardsManagerV2.RewardsPctTooHigh.selector);
        rewardsManager.setRewardsPctBps(2501);

        vm.prank(payerOne);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, payerOne));
        rewardsManager.setRewardsPctBps(1000);
    }
    
    // payProposer with bps=0: all to recipient
    function test_PayProposer_bpsZero_allToRecipient() public {
        vm.prank(ownerAddress);
        vm.expectEmit();
        emit RewardsPctBpsSet(0);
        rewardsManager.setRewardsPctBps(0);

        uint256 transferAmount = 5 ether;
        uint256 recipientBalanceBefore = feeRecipientOne.balance;
        uint256 toTreasuryBefore = rewardsManager.toTreasury();

        vm.prank(payerOne);
        vm.expectEmit();
        emit ProposerPaid(feeRecipientOne, transferAmount, 0);
        rewardsManager.payProposer{value: transferAmount}(payable(feeRecipientOne));

        uint256 recipientBalanceAfter = feeRecipientOne.balance;
        assertEq(recipientBalanceAfter, recipientBalanceBefore + transferAmount);

        uint256 toTreasuryAfter = rewardsManager.toTreasury();
        assertEq(toTreasuryAfter, toTreasuryBefore);
    }
    
    // payProposer with bps>0: split and accrue treasury
    function test_PayProposer_withBps_splits_andAccruesTreasury() public {
        vm.prank(ownerAddress);
        rewardsManager.setRewardsPctBps(1500);

        uint256 transferAmount = 10 ether;
        uint256 rewardPortion = (transferAmount * 1500) / 10_000; // 1.5e
        uint256 proposerPortion = transferAmount - rewardPortion;  // 8.5e

        uint256 recipientBalanceBefore = feeRecipientTwo.balance;
        uint256 toTreasuryBefore = rewardsManager.toTreasury();

        vm.prank(payerTwo);
        vm.expectEmit();
        emit ProposerPaid(feeRecipientTwo, proposerPortion, rewardPortion);
        rewardsManager.payProposer{value: transferAmount}(payable(feeRecipientTwo));

        uint256 recipientBalanceAfter = feeRecipientTwo.balance;
        assertEq(recipientBalanceAfter, recipientBalanceBefore + proposerPortion);

        uint256 toTreasuryAfter = rewardsManager.toTreasury();
        assertEq(toTreasuryAfter, toTreasuryBefore + rewardPortion);
    }
    
    // withdraw to treasury
    function test_WithdrawToTreasury_transfers_and_resets() public {
        vm.prank(ownerAddress);
        rewardsManager.setRewardsPctBps(2000); // 20%

        uint256 transferAmount = 5 ether;
        uint256 expectedRewardPortion = (transferAmount * 2000) / 10_000; // 1e

        vm.prank(payerOne);
        rewardsManager.payProposer{value: transferAmount}(payable(feeRecipientOne));

        uint256 toTreasuryBeforeWithdraw = rewardsManager.toTreasury();
        assertEq(toTreasuryBeforeWithdraw, expectedRewardPortion);

        uint256 treasuryBalanceBefore = treasuryAddress.balance;

        vm.prank(ownerAddress);
        vm.expectEmit();
        emit TreasuryWithdrawn(expectedRewardPortion);
        rewardsManager.withdrawToTreasury();

        uint256 treasuryBalanceAfter = treasuryAddress.balance;
        assertEq(treasuryBalanceAfter, treasuryBalanceBefore + expectedRewardPortion);

        uint256 toTreasuryAfterWithdraw = rewardsManager.toTreasury();
        assertEq(toTreasuryAfterWithdraw, 0);
    }

    // withdraw to treasury only owner
    function test_WithdrawToTreasury_onlyOwner() public {
        vm.prank(payerOne);
        vm.expectRevert(abi.encodeWithSelector(IRewardsManagerV2.OnlyOwnerOrTreasury.selector));
        rewardsManager.withdrawToTreasury();
    }

    function test_setTreasury_revertsIfTreasuryZero() public {
        vm.prank(ownerAddress);
        vm.expectRevert(IRewardsManagerV2.TreasuryIsZero.selector);
        rewardsManager.setTreasury(payable(address(0)));
    }

    // revert when no funds to withdraw
    function test_WithdrawToTreasury_revertsIfNoFunds() public {
        vm.prank(ownerAddress);
        vm.expectRevert(IRewardsManagerV2.NoFundsToWithdraw.selector);
        rewardsManager.withdrawToTreasury();
    }
    
    // receive/fallback revert
    function test_Receive_and_Fallback_revert() public {
        vm.expectRevert();
        (bool successReceive, ) = address(rewardsManager).call{value: 1}("");
        successReceive;

        vm.expectRevert();
        (bool successFallback, ) = address(rewardsManager).call(abi.encodeWithSignature("nonexistentFunction()"));
        successFallback;
    }
}

/// @dev Simple recipient that rejects any ETH transfers.
contract RejectingRecipient {
    receive() external payable {
        revert();
    }
}
