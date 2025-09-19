// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IRegistry} from "../interfaces/IRegistry.sol";

/// @title ValidatorOptInHubStorage
/// @notice Storage components of the ValidatorOptInHub contract.
contract ValidatorOptInHubStorage {

    IRegistry[] public registries;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}