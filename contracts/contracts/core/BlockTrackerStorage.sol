// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

abstract contract BlockTrackerStorage {
    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;
    
    uint256 public currentWindow;

    uint256 public blocksPerWindow;

    // Mapping from block number to the winner's address
    mapping(uint256 => address) public blockWinners;

     /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
