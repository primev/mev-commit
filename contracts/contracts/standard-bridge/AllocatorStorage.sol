// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

contract AllocatorStorage {
    /// @dev Mapping of whitelisted addresses which can mint native ETH on the mev-commit chain.
    mapping(address => bool) public whitelistedAddresses;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
