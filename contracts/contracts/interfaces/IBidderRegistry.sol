// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

interface IBidderRegistry {
    struct OpenedCommitment {
        string txnHash;
        uint256 bid;
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

    /// @dev Event emitted when funds are retrieved from a bidder's deposit
    event FundsRetrieved(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        uint256 indexed window,
        uint256 amount
    );

    /// @dev Event emitted when funds are retrieved from a bidder's deposit
    event FundsRewarded(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        address indexed provider,
        uint256 window,
        uint256 amount
    );

    /// @dev Event emitted when a bidder withdraws their deposit
    event BidderWithdrawal(
        address indexed bidder,
        uint256 indexed window,
        uint256 indexed amount
    );

    /// @dev Event emitted when the protocol fee recipient is updated
    event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient);

    /// @dev Event emitted when the fee payout period in blocks is updated
    event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks);

    function openBid(
        bytes32 commitmentDigest,
        uint256 bid,
        address bidder,
        uint64 blockNumber
    ) external;

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
