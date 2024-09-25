// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IVanillaRegistry} from "../interfaces/IVanillaRegistry.sol";
import {FeePayout} from "../utils/FeePayout.sol";

/// @title VanillaRegistryStorage
/// @notice Storage components of the VanillaRegistry contract.
contract VanillaRegistryStorage { 

    /// @dev Minimum stake required for validators. 
    uint256 public minStake;
    
    /// @dev Amount of ETH to slash per validator pubkey when a slash is invoked.
    uint256 public slashAmount;

    /// @dev Permissioned account that is able to invoke slashes.
    address public slashOracle; 

    /// @dev Number of blocks required between unstake initiation and withdrawal.
    uint256 public unstakePeriodBlocks;

    /// @dev Struct enabling automatic slashing funds payouts
    FeePayout.Tracker public slashingFundsTracker;

    /// @dev Mapping of BLS pubkeys to stored staked validator structs. 
    mapping(bytes => IVanillaRegistry.StakedValidator) public stakedValidators;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
