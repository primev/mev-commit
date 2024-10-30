// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {FeePayout} from "../utils/FeePayout.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";

abstract contract BidderRegistryStorage {
    using FeePayout for FeePayout.Tracker;

    /// @dev For improved precision
    uint256 constant public PRECISION = 1e25;
    uint256 constant public PERCENT = 100 * PRECISION;

    /// @dev Address of the preconfManager contract
    address public preconfManager;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Block tracker contract
    IBlockTracker public blockTrackerContract;

    /// Struct enabling automatic protocol fee payouts
    FeePayout.Tracker public protocolFeeTracker;

    /// @dev Mapping for if bidder is registered
    mapping(address => bool) public bidderRegistered;

    // Mapping from bidder addresses and window numbers to their locked funds
    mapping(address => mapping(uint256 => uint256)) public lockedFunds;

    // Mapping from bidder addresses and blocks to their used funds
    mapping(address => mapping(uint64 => uint256)) public usedFunds;

    /// Mapping from bidder addresses and window numbers to their funds per window
    mapping(address => mapping(uint256 => uint256)) public maxBidPerBlock;

    /// @dev Mapping from bidder addresses to their locked amount based on commitmentDigest
    mapping(bytes32 => IBidderRegistry.BidState) public bidPayment;

    /// @dev Amount assigned to bidders
    mapping(address => uint256) public providerAmount;

    /// @dev Amount assigned to bidders
    /// Not used anymore, still here bcs of upgradeability
    uint256 public blocksPerWindow;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
