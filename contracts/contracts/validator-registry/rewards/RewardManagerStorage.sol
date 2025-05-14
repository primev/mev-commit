// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {VanillaRegistryStorage} from "../VanillaRegistryStorage.sol";
import {MevCommitAVSStorage} from "../avs/MevCommitAVSStorage.sol";
import {MevCommitMiddlewareStorage} from "../middleware/MevCommitMiddlewareStorage.sol";

contract RewardManagerStorage {

    /// Storage reference to the VanillaRegistry contract.
    VanillaRegistryStorage internal _vanillaRegistry;

    /// Storage reference to the MevCommitAVS contract.
    MevCommitAVSStorage internal _mevCommitAVS;

    /// Storage reference to the MevCommitMiddleware contract.
    MevCommitMiddlewareStorage internal _mevCommitMiddleware;

    /// @dev Gas limit for forwarded auto-claim calls.
    uint256 public autoClaimGasLimit;

    /// @dev Receiver addresses with auto-claim enabled.
    mapping(address receiver => bool enabled) public autoClaim;

    /// @dev Receiver addresses which are blacklisted from auto-claim.
    mapping(address receiver => bool blacklisted) public autoClaimBlacklist;
    
    /// @dev Unclaimed rewards by (receiver or override) address.
    mapping(address addr => uint256 amount) public unclaimedRewards;

    /// @dev Orphaned rewards by validator pubkey.
    mapping(bytes pubkey => uint256 amount) public orphanedRewards;

    /// @dev Override addresses associated to receivers.
    mapping(address receiver => address overrideAddress) public overrideAddresses;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
