// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import { EventHeightLib } from "../utils/EventHeight.sol";

/// @title IValidatorRegistryV1
/// @notice Interface for the ValidatorRegistryV1 contract.
contract IValidatorRegistryV1 {

    /// @dev Event emitted when a validator is staked.
    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when ETH is added to the staked balance a validator. 
    event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance);

    /// @dev Event emitted when a validator is unstaked.
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator's stake is withdrawn.
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator is slashed.
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when the min stake parameter is set.
    event MinStakeSet(address indexed msgSender, uint256 newMinStake);

    /// @dev Event emitted when the slash amount parameter is set.
    event SlashAmountSet(address indexed msgSender, uint256 newSlashAmount);

    /// @dev Event emitted when the slash oracle parameter is set.
    event SlashOracleSet(address indexed msgSender, address newSlashOracle);

    /// @dev Event emitted when the slash receiver parameter is set.
    event SlashReceiverSet(address indexed msgSender, address newSlashReceiver);

    /// @dev Event emitted when the unstake period blocks parameter is set.
    event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks);

    /// @dev Struct representing a validator staked with the registry.
    struct StakedValidator {
        bool exists;
        uint256 balance;
        address withdrawalAddress;
        EventHeightLib.EventHeight unstakeHeight;
    }
}
