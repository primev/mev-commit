// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.26;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {IRewardsManagerV2} from "../../interfaces/IRewardsManagerV2.sol";
import {RewardsManagerV2Storage} from "./RewardsManagerV2Storage.sol";
import {Errors} from "../../utils/Errors.sol";
import {TransientReentrancyGuard} from "../../utils/TransientReentrancyGuard.sol";


contract RewardsManagerV2 is 
    Initializable, 
    Ownable2StepUpgradeable, 
    TransientReentrancyGuard, 
    RewardsManagerV2Storage, 
    IRewardsManagerV2, 
    UUPSUpgradeable 
{
    uint256 constant BPS_DENOMINATOR = 10_000;

    constructor() {
        _disableInitializers();
    }

    // -------- Initializer --------
    function initialize(address initialOwner, uint256 rewardsPctBps, address payable treasury) external initializer override {
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();
        _setRewardsPctBps(rewardsPctBps);
        _setTreasury(treasury);
    }

    // -------- Proposer payment (EL rewards routed through this contract) --------
    function payProposer(address payable feeRecipient) external payable nonReentrant {
        uint256 totalAmt = msg.value;
        uint256 bps = rewardsPctBps;
        if (bps == 0) {
            (bool success, ) = feeRecipient.call{value: totalAmt}("");
            require(success, ProposerTransferFailed(feeRecipient, totalAmt));
            emit ProposerPaid(feeRecipient, totalAmt, 0);
        } else {
            uint256 amtForRewards = totalAmt * bps / BPS_DENOMINATOR;
            uint256 proposerAmt = totalAmt - amtForRewards;
            unchecked { toTreasury += amtForRewards; }
            (bool success, ) = feeRecipient.call{value: proposerAmt}("");
            require(success, ProposerTransferFailed(feeRecipient, proposerAmt));
            emit ProposerPaid(feeRecipient, proposerAmt, amtForRewards);
        }
    }

    // -------- Owner Functions--------

    function withdrawToTreasury() external nonReentrant onlyOwner {
        require(treasury != address(0), TreasuryIsZero());
        require(toTreasury > 0, NoFundsToWithdraw());
        uint256 treasuryAmt = toTreasury;
        toTreasury = 0;
        treasury.call{value: treasuryAmt}("");
        emit TreasuryWithdrawn(treasuryAmt);
    }

    function setRewardsPctBps(uint256 rewardsPctBps) external onlyOwner {
        _setRewardsPctBps(rewardsPctBps);
    }

    function setTreasury(address payable treasury) external onlyOwner {
        _setTreasury(treasury);
    }

    // -------- Internal --------

    function _setTreasury(address payable _treasury) internal {
        treasury = _treasury;
        emit TreasurySet(treasury);
    }

    function _setRewardsPctBps(uint256 _rewardsPctBps) internal {
        require (_rewardsPctBps <= 2500, RewardsPctTooHigh());
        rewardsPctBps = _rewardsPctBps;
        emit RewardsPctBpsSet(rewardsPctBps);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // -------- Receive/Fallback (explicitly disabled) --------
    receive() external payable { revert Errors.InvalidReceive(); }
    fallback() external payable { revert Errors.InvalidFallback(); }
}
