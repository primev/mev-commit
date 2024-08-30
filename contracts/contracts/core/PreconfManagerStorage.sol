// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IPreconfManager} from "../interfaces/IPreconfManager.sol";

abstract contract PreconfManagerStorage {
    // Represents the dispatch window in milliseconds
    uint64 public commitmentDispatchWindow;

    /// @dev Address of the oracle contract
    address public oracleContract;

    /// @dev The number of blocks per window
    uint256 public blocksPerWindow;

    /// @dev Address of provider registry
    IProviderRegistry public providerRegistry;

    /// @dev Address of bidderRegistry
    IBidderRegistry public bidderRegistry;

    /// @dev Address of blockTracker
    IBlockTracker public blockTracker;

    /// @dev Mapping from provider to commitments count
    mapping(address => uint256) public commitmentsCount;

    /// @dev Commitment Hash -> Opened Commitemnt
    /// @dev Only stores valid commitments
    mapping(bytes32 => IPreconfManager.OpenedCommitment) public openedCommitments;

    /// @dev Unopened Commitment Hash -> Unopened Commitment
    /// @dev Only stores valid unopened commitments
    mapping(bytes32 => IPreconfManager.UnopenedCommitment) public unopenedCommitments;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
