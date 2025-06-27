// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {FeePayout} from "../utils/FeePayout.sol";

abstract contract ProviderRegistryStorageV2 {
    using FeePayout for FeePayout.Tracker;

    /// @dev For improved precision
    uint256 public constant PRECISION = 1e16;
    uint256 public constant ONE_HUNDRED_PERCENT = 100 * PRECISION;

    /// @dev Minimum stake required for registration
    uint256 public minStake;

    /// @dev Address of the preconf manager
    address public preconfManager;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint256 public feePercent;

    /// @dev Configurable withdrawal delay in milliseconds
    uint256 public withdrawalDelay;

    /// Struct enabling automatic penalty fee payouts
    FeePayout.TimestampTracker public penaltyFeeTracker;

    /// @dev Mapping from provider address to whether they are registered or not
    mapping(address => bool) public providerRegistered;

    /// @dev Mapping from provider addresses to their staked amount
    mapping(address => uint256) public providerStakes;

    /// @dev Mapping of provider to withdrawal request timestamp
    mapping(address => uint256) public withdrawalRequests;

    /// @dev Mapping from bidder to provider slashed amount
    mapping(address => uint256) public bidderSlashedAmount;

   /// @dev Maps BLS public key to their corresponding block builder address
    mapping(bytes => address) public blockBuilderBLSKeyToAddress;

    /// @dev Mapping from a provider's EOA address to their BLS public keys
    mapping(address => bytes[]) public eoaToBlsPubkeys;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
