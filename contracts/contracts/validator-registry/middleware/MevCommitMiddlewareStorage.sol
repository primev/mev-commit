// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";

abstract contract MevCommitMiddlewareStorage {

    uint256 public validatorDeregPeriodBlocks;

    uint256 public operatorDeregPeriodBlocks;

    uint256 public vaultDeregPeriodBlocks;

    address public slashOracle;

    mapping(bytes blsPubkey => IMevCommitMiddleware.ValidatorRecord) public validatorRecords;

    mapping(address operatorAddress => IMevCommitMiddleware.OperatorRecord) public operatorRecords;

    mapping(address vaultAddress => IMevCommitMiddleware.VaultRecord) public vaultRecords;

    mapping(address vault => EnumerableSet.BytesSet) internal _vaultToValidatorSet;
}
