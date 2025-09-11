// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import {IBlockRewardManager} from "../../interfaces/IBlockRewardManager.sol";
import {BlockRewardManagerStorage} from "./BlockRewardManagerStorage.sol";
import {Errors} from "../../utils/Errors.sol";

contract BlockRewardManager is 
    Initializable, 
    Ownable2StepUpgradeable, 
    ReentrancyGuardUpgradeable, 
    BlockRewardManagerStorage, 
    IBlockRewardManager, 
    UUPSUpgradeable 
{
    uint256 constant _BPS_DENOMINATOR = 10_000;

    modifier onlyOwnerOrTreasury() {
        require(msg.sender == owner() || msg.sender == treasury, OnlyOwnerOrTreasury());
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // -------- Receive/Fallback (explicitly disabled) --------
    receive() external payable { revert Errors.InvalidReceive(); }
    fallback() external payable { revert Errors.InvalidFallback(); }

    // -------- Initializer --------
    function initialize(address initialOwner, uint256 rewardsPctBps, address treasury) external initializer override {
        __Ownable_init(initialOwner);
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();
        _setRewardsPctBps(rewardsPctBps);
        _setTreasury(treasury);
    }

    // -------- Proposer payment (EL rewards routed through this contract) --------
    function payProposer(address payable feeRecipient) external payable {
        uint256 totalAmt = msg.value;
        uint256 bps = rewardsPctBps;
        if (bps == 0) {
            (bool success, ) = feeRecipient.call{value: totalAmt}("");
            require(success, ProposerTransferFailed(feeRecipient, totalAmt)); //revert if transfer fails
            emit ProposerPaid(feeRecipient, totalAmt, 0);
        } else {
            uint256 amtForRewards = totalAmt * bps / _BPS_DENOMINATOR;
            uint256 proposerAmt = totalAmt - amtForRewards;
            toTreasury += amtForRewards;
            (bool success, ) = feeRecipient.call{value: proposerAmt}("");
            require(success, ProposerTransferFailed(feeRecipient, proposerAmt)); //revert if transfer fails
            emit ProposerPaid(feeRecipient, proposerAmt, amtForRewards);
        }
    }

    function withdrawToTreasury() external onlyOwnerOrTreasury {
        require(toTreasury > 0, NoFundsToWithdraw());
        uint256 treasuryAmt = toTreasury;
        toTreasury = 0;
        treasury.call{value: treasuryAmt}(""); //Treasury will not revert
        emit TreasuryWithdrawn(treasuryAmt);
    }

    function setRewardsPctBps(uint256 rewardsPctBps) external onlyOwner {
        _setRewardsPctBps(rewardsPctBps);
    }

    function setTreasury(address treasury) external onlyOwner {
        _setTreasury(treasury);
    }
    
    function _setTreasury(address _treasury) internal {
        require(_treasury != address(0), TreasuryIsZero());
        treasury = payable(_treasury);
        emit TreasurySet(_treasury);
    }

    function _setRewardsPctBps(uint256 _rewardsPctBps) internal {
        require (_rewardsPctBps <= 2500, RewardsPctTooHigh());
        rewardsPctBps = _rewardsPctBps;
        emit RewardsPctBpsSet(_rewardsPctBps);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
