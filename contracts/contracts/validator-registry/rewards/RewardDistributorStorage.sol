// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IRewardDistributor} from "../../interfaces/IRewardDistributor.sol";

/// @title RewardDistributorStorage
/// @notice Storage layout for RewardDistributor
abstract contract RewardDistributorStorage {
    /// @dev Address authorized to grant ETH and token rewards.
    address public rewardManager;
    mapping(uint256 id => address token) public rewardTokens;

    /// @dev Default recipient per operator (used when no pubkey-specific override exists).
    mapping(address operator => address recipient) public operatorGlobalOverride;
    /// @dev Recipient override by BLS pubkey hash (keccak256(pubkey)).
    mapping(address operator => mapping(bytes32 keyhash => address recipient)) public operatorKeyOverrides;

    /// @dev Accrued and claimed amounts per (operator, recipient).
    mapping(address operator => mapping(address recipient => mapping(uint256 tokenID => IRewardDistributor.RewardData))) public rewardData;

    /// @dev Operator → recipient → delegate → isAuthorized
    mapping(address operator => mapping(address recipient => mapping(address delegate => bool))) public claimDelegate;

    // === Storage gap for future upgrades ===
    uint256[48] private __gap;
}