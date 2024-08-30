// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

interface IProviderRegistry {

    /// @dev Event emitted when a provider is registered
    event ProviderRegistered(address indexed provider, uint256 stakedAmount, bytes blsPublicKey);

    /// @dev Event emitted when funds are deposited
    event FundsDeposited(address indexed provider, uint256 amount);

    /// @dev Event emitted when funds are slashed
    event FundsSlashed(address indexed provider, uint256 amount);

    /// @dev Event emitted when withdrawal is requested
    event Unstake(address indexed provider, uint256 timestamp);

    /// @dev Event emitted when withdrawal is completed
    event Withdraw(address indexed provider, uint256 amount);

    /// @dev Event emitted when the withdrawal delay is updated
    event WithdrawalDelayUpdated(uint256 newWithdrawalDelay);

    /// @dev Event emitted when the penalty fee recipient is updated
    event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient);

    /// @dev Event emitted when the fee payout period in blocks is updated
    event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks);

    function registerAndStake(bytes calldata blsPublicKey) external payable;

    function stake() external payable;

    function slash(
        uint256 amt,
        address provider,
        address payable bidder,
        uint256 residualBidPercentAfterDecay
    ) external;
    
    function isProviderValid(address committerAddress) external view;
}
