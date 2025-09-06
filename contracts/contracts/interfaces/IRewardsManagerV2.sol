// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.26;


interface IRewardsManagerV2 {
    /// @dev Config mode: 0=unset, 1=primaryOnly (100% to feeRecipient), 2=withSecondary (secBps to `secondary`)
    

    error RewardsPctTooHigh();
    error TreasuryIsZero();
    error NoFundsToWithdraw();
    error ProposerTransferFailed(address feeRecipient, uint256 amount);

    // -------- Views --------
 

    /// @notice Builders/relays call this to route EL rewards *through* this contract.
    /// If no config / not opted in / operator changed, pays 100% to `feeRecipient`.
    function payProposer(address payable feeRecipient) external payable;

    function withdrawToTreasury() external;

    function setRewardsPctBps(uint256 rewardsPctBps) external;

    function setTreasury(address payable treasury) external;
    
    // -------- Admin --------
    function initialize(address initialOwner, uint256 rewardsPctBps, address payable treasury) external;

    // -------- Events --------

    /// @notice Emitted for each proposer payment routed by this contract
    event ProposerPaid(
        address indexed feeRecipient,
        uint256 indexed proposerAmt,
        uint256 indexed rewardAmt
    );
    /// @notice Emitted when the treasury is withdrawn
    event TreasuryWithdrawn(uint256 indexed treasuryAmt);
    /// @notice Emitted when the rewards pct is set
    event RewardsPctBpsSet(uint256 indexed rewardsPctBps);
    /// @notice Emitted when the treasury is set
    event TreasurySet(address indexed treasury);



}