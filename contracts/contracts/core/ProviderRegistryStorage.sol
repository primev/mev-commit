// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {FeePayout} from "../utils/FeePayout.sol";

abstract contract ProviderRegistryStorage {
    using FeePayout for FeePayout.Tracker;

    /// @dev For improved precision
    uint256 public constant PRECISION = 1e25;
    uint256 public constant PERCENT = 100 * PRECISION;

    /// @dev Minimum stake required for registration
    uint256 public minStake;

    /// @dev Address of the preconf manager
    address public preconfManager;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Configurable withdrawal delay in milliseconds
    uint256 public withdrawalDelay;

    /// Struct enabling automatic penalty fee payouts
    FeePayout.Tracker public penaltyFeeTracker;

    /// @dev Mapping from provider address to whether they are registered or not
    mapping(address => bool) public providerRegistered;

    /// @dev Mapping from a provider's EOA address to their BLS public key
    mapping(address => bytes) public eoaToBlsPubkey;

    /// @dev Mapping from provider addresses to their staked amount
    mapping(address => uint256) public providerStakes;

    /// @dev Mapping of provider to withdrawal request timestamp
    mapping(address => uint256) public withdrawalRequests;

    /// @dev Mapping from bidder to provider slashed amount
    mapping(address => uint256) public bidderSlashedAmount;
    
    /// @dev Maps BLS public keys to their corresponding block builder addresses
    mapping(bytes => address) public blockBuilderBLSKeyToAddress;
    
    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[47] private __gap;
}
