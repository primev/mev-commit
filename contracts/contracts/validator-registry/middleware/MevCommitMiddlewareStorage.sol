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

    uint256 public validatorDeregPeriodBlocks;

    uint256 public operatorDeregPeriodBlocks;

    uint256 public vaultDeregPeriodBlocks;

    address public slashOracle;

    mapping(bytes blsPubkey => IMevCommitMiddleware.ValidatorRecord) public validatorRecords;

    mapping(address operatorAddress => IMevCommitMiddleware.OperatorRecord) public operatorRecords;

    mapping(address vaultAddress => IMevCommitMiddleware.VaultRecord) public vaultRecords;

    mapping(address vault => EnumerableSet.BytesSet) internal _vaultToValidatorSet;
}
