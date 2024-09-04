// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";

abstract contract MevCommitMiddlewareStorage {

    IRegistry public networkRegistry;

    IRegistry public operatorRegistry;

    IRegistry public vaultFactory;

    /// @notice The network address, which must have registered with the NETWORK_REGISTRY.
    address public network;

    /// @dev A period in blocks during which the mev-commit oracle can invoke slashing.
    /// @notice This serves as the deregistration period for all of validator, operator, and vault records.
    /// @notice This also serves as the number of blocks that a registered Vault's epochDuration must be greater than.
    uint256 public slashPeriodBlocks;

    address public slashOracle;

    mapping(bytes blsPubkey => IMevCommitMiddleware.ValidatorRecord) public validatorRecords;

    mapping(address operatorAddress => IMevCommitMiddleware.OperatorRecord) public operatorRecords;

    mapping(address vaultAddress => IMevCommitMiddleware.VaultRecord) public vaultRecords;

    mapping(address vault =>
        mapping(address operator => EnumerableSet.BytesSet)) internal _vaultAndOperatorToValset;
}
