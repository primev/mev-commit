// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

abstract contract BlockTrackerStorage {
    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;
    
    uint256 public currentWindow;

    uint256 public blocksPerWindow;

    // Mapping from block number to the winner's address
    mapping(uint256 => address) public blockWinners;

     /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;
}
