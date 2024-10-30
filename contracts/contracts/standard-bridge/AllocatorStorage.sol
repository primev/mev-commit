// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

contract AllocatorStorage {
    /// @dev Mapping of whitelisted addresses which can mint native ETH on the mev-commit chain.
    mapping(address => bool) public whitelistedAddresses;

    /// @dev Mapping that tracks funds which had a failed transfer to the recipient and need manual withdrawal.
    mapping(address recipient => uint256 amount) public transferredFundsNeedingWithdrawal;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
