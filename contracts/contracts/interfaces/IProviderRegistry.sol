// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

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

    /// @dev Event emitted when the min stake is updated
    event MinStakeUpdated(uint256 indexed newMinStake);

    /// @dev Event emitted when the preconf manager is updated
    event PreconfManagerUpdated(address indexed newPreconfManager);

    /// @dev Event emitted when the fee percent is updated
    event FeePercentUpdated(uint16 indexed newFeePercent);

    /// @dev Event emitted when there are insufficient funds to slash
    event InsufficientFundsToSlash(
        address indexed provider,
        uint256 providerStake,
        uint256 residualAmount,
        uint256 penaltyFee
    );

    /// @dev Event emitted when bidder withdraws slashed funds
    event BidderWithdrawal(address bidder, uint256 amount);

    /// @dev Event emitted when transfer to bidder fails
    event TransferToBidderFailed(address bidder, uint256 amount);

    error NotPreconfContract(address sender, address preconfManager);
    error NoStakeToWithdraw(address sender);
    error UnstakeRequestExists(address sender);
    error NoUnstakeRequest(address sender);
    error DelayNotPassed(uint256 withdrawalRequestTimestamp, uint256 withdrawalDelay, uint256 currentBlockTimestamp);
    error PreconfManagerNotSet();
    error ProviderCommitmentsPending(address sender, uint256 numPending);
    error StakeTransferFailed(address sender, uint256 amount);
    error ProviderAlreadyRegistered(address sender);
    error InsufficientStake(uint256 stake, uint256 minStake);
    error InvalidBLSPublicKeyLength(uint256 length, uint256 expectedLength);
    error ProviderNotRegistered(address sender);
    error PendingWithdrawalRequest(address sender);
    error BidderAmountIsZero(address sender);
    error BidderWithdrawalTransferFailed(address sender, uint256 amount);
    
    function registerAndStake(bytes calldata blsPublicKey) external payable;

    function stake() external payable;

    function slash(
        uint256 amt,
        address provider,
        address payable bidder
    ) external;
    
    function isProviderValid(address committerAddress) external view;
}
