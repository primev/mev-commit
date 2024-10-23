// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

abstract contract BlockTrackerStorage {
    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;
    
    uint256 public currentWindow;

    uint256 public blocksPerWindow;

    // Mapping from block number to the winner's BLS key
    mapping(uint256 => bytes) public blockWinners;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[49] private __gap;
}
