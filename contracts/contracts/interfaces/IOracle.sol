// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IOracle {

    /// @dev Event emitted when the oracle account is set.
    event OracleAccountSet(
        address indexed oldOracleAccount,
        address indexed newOracleAccount
    );

    /// @dev Event emitted when the preconf manager is set.
    event PreconfManagerSet(address indexed newPreconfManager);

    /// @dev Event emitted when the block tracker is set.
    event BlockTrackerSet(address indexed newBlockTracker);

    /// @dev Event emitted when a commitment is processed.
    event CommitmentProcessed(bytes32 indexed commitmentDigest, bool isSlash, uint256 residualBidPercentAfterDecay);

    /// @dev Error emitted when the sender is not the oracle account
    error NotOracleAccount(address sender, address oracleAccount);

    /// @dev Error emitted when the builder is not the block winner
    error BuilderNotBlockWinner(address blockWinner, address builder);

    /// @dev Error emitted when the residual bid percent after decay exceeds max amount
    error ResidualBidPercentAfterDecayExceedsMax(uint256 residualBidPercentAfterDecay);

    receive() external payable;

    fallback() external payable;

    function initialize(
        address preConfContract_,
        address blockTrackerContract_,
        address oracleAccount_,
        address owner_
    ) external;

    function processBuilderCommitmentForBlockNumber(
        bytes32 commitmentDigest,
        uint256 blockNumber,
        address builder,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) external;

    function setOracleAccount(address newOracleAccount) external;
}
