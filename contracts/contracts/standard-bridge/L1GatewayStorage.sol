// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @dev Any storage variables defined in this contract must NOT override those defined in GatewayStorage.sol!
contract L1GatewayStorage {
    /// @dev Mapping that tracks funds which had a failed transfer to the recipient and need manual withdrawal.
    mapping(address recipient => uint256 amount) public transferredFundsNeedingWithdrawal;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __l1GatewayStorageGap;
}
