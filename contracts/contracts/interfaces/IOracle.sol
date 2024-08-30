// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

interface IOracle {

    /// @dev Event emitted when the oracle account is set.
    event OracleAccountSet(
        address indexed oldOracleAccount,
        address indexed newOracleAccount
    );

    /// @dev Event emitted when a commitment is processed.
    event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash);

    receive() external payable;

    fallback() external payable;

    function initialize(
        address preConfContract_,
        address blockTrackerContract_,
        address oracleAccount_,
        address owner_
    ) external;

    function processBuilderCommitmentForBlockNumber(
        bytes32 commitmentIndex,
        uint256 blockNumber,
        address builder,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) external;

    function setOracleAccount(address newOracleAccount) external;
}
