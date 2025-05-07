// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {VanillaRegistryStorage} from "./VanillaRegistryStorage.sol";
import {MevCommitAVSStorage} from "./avs/MevCommitAVSStorage.sol";
import {MevCommitMiddlewareStorage} from "./middleware/MevCommitMiddlewareStorage.sol";

contract RewardManagerStorage {
    VanillaRegistryStorage internal _vanillaRegistry;
    MevCommitAVSStorage internal _mevCommitAVS;
    MevCommitMiddlewareStorage internal _mevCommitMiddleware;

    uint256 public autoClaimGasLimit;

    mapping(address addr => bool enabled) public autoClaim;
    
    mapping(address addr => uint256 amount) public rewards;

    mapping(bytes pubkey => uint256 amount) public orphanedRewards;
    
    mapping(address delegator => address overrideAddress) public overrideClaimAddresses;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
