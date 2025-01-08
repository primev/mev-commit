// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IBidderRegistry {
    struct OpenedCommitment {
        string txnHash;
        uint256 bidAmt;
        uint64 blockNumber;
        string bidHash;
        string bidSignature;
        string commitmentDigest;
        string commitmentSignature;
    }

    struct BidState {
        address bidder;
        uint256 bidAmt;
        State state;
    }

    enum State {
        Undefined,
        PreConfirmed,
        Withdrawn
    }

    /// @dev Event emitted when a bidder is registered with their deposited amount
    event BidderRegistered(
        address indexed bidder,
        uint256 indexed depositedAmount,
        uint256 indexed windowNumber
    );

    /// @dev Event emitted when funds are retrieved from a bidder's deposit for a commitment that was not successfull
    event FundsRetrieved(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        uint256 indexed window,
        uint256 amount
    );

    /// @dev Event emitted when funds are rewarded to a provider from a bidder's deposit for successfully carrying out a commitment
    event FundsRewarded(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        address indexed provider,
        uint256 window,
        uint256 amount
    );

    /// @dev Event emitted when left over funds are returned to a bidder after a commitment is processed
    event LeftOverFundsReturned(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        uint256 window,
        uint256 amount
    );

    /// @dev Event emitted when a bidder withdraws their deposit
    event BidderWithdrawal(
        address indexed bidder,
        uint256 indexed window,
        uint256 indexed amount
    );

    /// @dev Event emitted when a bidder's bid amount exceeds the available amount for a commitment
    event BidAmountExceedsAvailableAmount(bytes32 indexed commitmentDigest, uint256 bidAmt, uint256 availableAmount);

    /// @dev Event emitted when a bidder's bid amount is used for a commitment
    event BidAmountUsed(bytes32 indexed commitmentDigest, uint256 bidAmt, uint64 blockNumber);

    /// @dev Event emitted when the preconfManager is updated
    event PreconfManagerUpdated(address indexed newPreconfManager);

    /// @dev Event emitted when the fee percent is updated
    event FeePercentUpdated(uint256 indexed newFeePercent);

    /// @dev Event emitted when the block tracker is updated
    event BlockTrackerUpdated(address indexed newBlockTracker);

    /// @dev Event emitted when the fee payout period in blocks is updated
    event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks);

    /// @dev Event emitted when the protocol fee recipient is updated
    event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient);

    /// @dev Event emitted when transfer to bidder fails
    event TransferToBidderFailed(bytes32 indexed commitmentDigest, address indexed bidder, uint256 amount);

    /// @dev Event emitted when the protocol fee is deducted from the bidder's deposit for successfully carrying out a commitment
    event ProtocolFeeTransferred(bytes32 indexed commitmentDigest, uint256 amount, address feeRecipient);

    /// @dev Error emitted when the sender is not the preconfManager
    error SenderIsNotPreconfManager(address sender, address preconfManager);

    /// @dev Error emitted when the bid is not preconfirmed
    error BidNotPreConfirmed(bytes32 commitmentDigest, State actualState, State expectedState);

    /// @dev Error emitted when the withdraw after window settled
    error WithdrawAfterWindowSettled(uint256 window, uint256 currentWindow);

    /// @dev Error emitted when the transfer to the provider fails
    error TransferToProviderFailed(address provider, uint256 amount);

    /// @dev Error emitted when the provider amount is zero
    error ProviderAmountIsZero(address provider);

    /// @dev Error emitted when the only bidder can withdraw
    error OnlyBidderCanWithdraw(address sender, address bidder);

    /// @dev Error emitted when the bidder tries to deposit 0 amount
    error DepositAmountIsZero();

    /// @dev Error emitted when the window is not settled
    error WindowNotSettled();

    /// @dev Error emitted when withdrawal transfer failed
    error BidderWithdrawalTransferFailed(address bidder, uint256 amount);

    function openBid(
        bytes32 commitmentDigest,
        uint256 bidAmt,
        address bidder,
        uint64 blockNumber
    ) external returns (uint256);

    function depositForWindow(uint256 window) external payable;

    function retrieveFunds(
        uint256 windowToSettle,
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external;

    function unlockFunds(uint256 windowToSettle, bytes32 bidID) external;

    function getDeposit(
        address bidder,
        uint256 window
    ) external view returns (uint256);
}
