// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IVanillaRegistryV2} from "../interfaces/IVanillaRegistryV2.sol";
import {FeePayout} from "../utils/FeePayout.sol";

/// @title VanillaRegistryStorageV2
/// @notice Storage components of the VanillaRegistryV2 contract.
contract VanillaRegistryStorageV2 { 

    /// @dev Minimum stake required for validators, also used as the slash amount.
    uint256 public minStake;
    
    /// @dev Permissioned account that is able to invoke slashes.
    address public slashOracle; 

    /// @dev Number of blocks required between unstake initiation and withdrawal.
    uint256 public unstakePeriodBlocks;

    /// @dev Struct enabling automatic slashing funds payouts
    FeePayout.Tracker public slashingFundsTracker;

    /// @dev Mapping of BLS pubkeys to stored staked validator structs. 
    mapping(bytes => IVanillaRegistryV2.StakedValidator) public stakedValidators;

    /// @dev Mapping of withdrawal addresses to claimable ETH that was force withdrawn by the owner.
    mapping(address withdrawalAddress => uint256 amountToClaim) public forceWithdrawnFunds;

    /// @dev Mapping of staker addresses to whether they are whitelisted.
    mapping(address staker => bool whitelisted) public whitelistedStakers;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
