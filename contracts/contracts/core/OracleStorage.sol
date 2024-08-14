// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {IPreconfManager} from "../interfaces/IPreconfManager.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";

abstract contract OracleStorage {
    /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;

    /// @dev Reference to the PreconfManager contract interface.
    IPreconfManager internal _preConfContract;

    /// @dev Reference to the BlockTracker contract interface.
    IBlockTracker internal _blockTrackerContract;
}
