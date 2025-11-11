// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title GasTankManager storage layout
/// @notice Keeps state variables isolated for upgradeable deployments.
abstract contract GasTankManagerStorage {
    /// @dev Minimum wei balance the provider should top up to
    uint256 public minDeposit;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
