// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import { TimestampOccurrence } from "../utils/Occurrence.sol";

interface IBidderRegistry {
    enum State {
        Undefined,
        PreConfirmed,
        Withdrawn
    }

    struct OpenedCommitment {
        string txnHash;
        uint256 bidAmt;
        uint64 blockNumber;
        string bidHash;
        string bidSignature;
        string commitmentDigest;
        string commitmentSignature;
    }

    // Represents a bidder's deposit for a specific provider
    struct Deposit {
        // Whether a deposit exists
        bool exists;
        // Amount deposited for this provider, not yet associated with an opened bid
        uint256 availableAmount;
        // Cumulative amount escrowed toward bid(s) for this provider
        /// @dev This corresponds to funds from bids that have been opened, but not yet settled
        uint256 escrowedAmount;
        // Occurrence struct facilitating withdrawal request
        TimestampOccurrence.Occurrence withdrawalRequestOccurrence;
    }

    struct BidState {
        address bidder;
        uint256 bidAmt;
        State state;
    }

    /// @dev Event emitted when a bidder is registered with their deposited amount
    event BidderDeposited(
        address indexed bidder,
        address indexed provider,
        uint256 indexed depositedAmount
    );

    /// @dev Event emitted when a bidder requests a withdrawal from a specific provider
    event WithdrawalRequested(
        address indexed bidder,
        address indexed provider,
        uint256 availableAmount,
        uint256 escrowedAmount,
        uint256 indexed timestamp
    );

    /// @dev Event emitted when funds are unlocked from a bidder's escrowed deposit
    event FundsUnlocked(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        address indexed provider,
        uint256 amount
    );

    /// @dev Event emitted when funds are retrieved from a bidder's deposit
    event FundsRewarded(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        address indexed provider,
        uint256 amount
    );

    /// @dev Event emitted when a bidder withdraws their deposit
    event BidderWithdrawal(
        address indexed bidder,
        address indexed provider,
        uint256 indexed amountWithdrawn,
        uint256 amountStillEscrowed
    );

    /// @dev Event emitted when the deposit manager implementation is updated
    event DepositManagerImplUpdated(address indexed newDepositManagerImpl);

    /// @dev Event emitted when the preconfManager is updated
    event PreconfManagerUpdated(address indexed newPreconfManager);

    /// @dev Event emitted when the fee percent is updated
    event FeePercentUpdated(uint256 indexed newFeePercent);

    /// @dev Event emitted when the block tracker is updated
    event BlockTrackerUpdated(address indexed newBlockTracker);

    /// @dev Event emitted when the fee payout period is updated
    event FeePayoutPeriodUpdated(uint256 indexed newFeePayoutPeriod);

    /// @dev Event emitted when the protocol fee recipient is updated
    event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient);

    /// @dev Event emitted when transfer to bidder fails
    event TransferToBidderFailed(address bidder, uint256 amount);

    /// @dev Event emitted when a bidder's top-up instance fails during openBid
    event TopUpFailed(address indexed bidder, address indexed provider);

    /// @dev Event emitted when an opened bid amount is reduced due to the bidder not having enough funds deposited    
    event BidAmountReduced(address indexed bidder, address indexed provider, uint256 indexed newBidAmt);

    /// @dev Error emitted when the sender is not the preconfManager
    error SenderIsNotPreconfManager(address sender, address preconfManager);

    /// @dev Error emitted when the deposit manager implementation is not set
    error DepositManagerNotSet();

    /// @dev Error emitted when the bid is not preconfirmed
    error BidNotPreConfirmed(bytes32 commitmentDigest, State actualState, State expectedState);

    /// @dev Error emitted when the transfer to the provider fails
    error TransferToProviderFailed(address provider, uint256 amount);

    /// @dev Error emitted when the provider amount is zero
    error ProviderAmountIsZero(address provider);

    /// @dev Error emitted when the only bidder can withdraw
    error OnlyBidderCanWithdraw(address sender, address bidder);

    /// @dev Error emitted when the bidder tries to deposit 0 amount
    error DepositAmountIsZero();

    /// @dev Error emitted when no providers are given as an argument
    error NoProviders();

    /// @dev Error emitted when withdrawal transfer failed
    error BidderWithdrawalTransferFailed(address bidder, uint256 amount);

    /// @dev Error emitted when the bidder withdrawal period has not elapsed
    error WithdrawalPeriodNotElapsed(uint256 currentTimestampMs, uint256 withdrawalTimestampMs, uint256 withdrawalPeriodMs);

    /// @dev Error emitted when a deposit does not exist
    error DepositDoesNotExist(address bidder, address provider);

    /// @dev Error emitted when a withdrawal occurrence exists
    error WithdrawalOccurrenceExists(address bidder, address provider, uint256 requestTimestamp);

    /// @dev Error emitted when a withdrawal occurrence does not exist
    error WithdrawalOccurrenceDoesNotExist(address bidder, address provider);

    function openBid(
        bytes32 commitmentDigest,
        uint256 bidAmt,
        address bidder,
        address provider
    ) external returns (uint256);

    function depositAsBidder(address provider) external payable;

    function convertFundsToProviderReward(
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external;

    function unlockFunds(address provider, bytes32 commitmentDigest) external;

    function getDeposit(
        address bidder,
        address provider
    ) external view returns (uint256);

    function withdrawalRequestExists(
        address bidder,
        address provider
    ) external view returns (bool);
}
