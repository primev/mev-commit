// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {FeePayout} from "../utils/FeePayout.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";

abstract contract BidderRegistryStorage {
    using FeePayout for FeePayout.Tracker;

    /// @dev For improved precision
    uint256 constant public PRECISION = 1e16;
    uint256 constant public ONE_HUNDRED_PERCENT = 100 * PRECISION;

    /// @dev Address of the preconfManager contract
    address public preconfManager;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint256 public feePercent;

    /// @dev Block tracker contract
    IBlockTracker public blockTrackerContract;

    /// Struct enabling automatic protocol fee payouts
    FeePayout.TimestampTracker public protocolFeeTracker;

    /// @dev Mapping from commitment digest for a bid, to its BidState
    mapping(bytes32 => IBidderRegistry.BidState) public bidPayment;

    /// @dev Funds rewarded to providers for fulfilling commitments
    mapping(address => uint256) public providerAmount;

    /// @dev Bidder withdrawal period in milliseconds (mev-commit chain uses ms timestamps)
    /// @dev This period should be greater than the worst case scenario amount of time it'd take for a newly opened bid to be settled.
    uint256 public bidderWithdrawalPeriodMs;

    /// @dev Mapping from bidder address to deposits for specific providers
    mapping(address bidder => mapping(address provider => IBidderRegistry.Deposit deposit)) public deposits;

    /// @dev Address of the deposit manager implementation contract
    address public depositManagerImpl;

    /// Hash of EIP-7702 stub (0xef0100‖impl) for the deposit manager implementation contract
    bytes32 public depositManagerHash;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
