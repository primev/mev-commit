// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IPreconfManager} from "../interfaces/IPreconfManager.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";

abstract contract OracleStorage {
    /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;

    /// @dev Reference to the PreconfManager contract interface.
    IPreconfManager internal _preconfManager;

    /// @dev Reference to the BlockTracker contract interface.
    IBlockTracker internal _blockTrackerContract;

    /// @dev Reference to the ProviderRegistry contract interface.
    IProviderRegistry internal _providerRegistry;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
