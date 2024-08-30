// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IVanillaRegistry} from "../interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";

/// @title ValidatorOptInRouterStorage
/// @notice Storage components of the ValidatorOptInRouter contract.
contract ValidatorOptInRouterStorage {

    /// @notice The address of the vanilla registry contract.
    IVanillaRegistry public vanillaRegistry;

    /// @notice The address of the mev-commit AVS contract.
    IMevCommitAVS public mevCommitAVS;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
