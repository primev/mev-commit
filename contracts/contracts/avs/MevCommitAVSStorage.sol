// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;


import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";

abstract contract MevCommitAVSStorage {

    /// @notice Mapping of operator addresses to their registration info
    mapping(address => IMevCommitAVS.OperatorRegistrationInfo) public operatorRegistrations;

    /// @notice Mapping of validator pubkeys to their registration info
    mapping(bytes => IMevCommitAVS.ValidatorRegistrationInfo) public validatorRegistrations;

    /// @notice List of restakeable strategy addresses
    address[] public restakeableStrategies;

    /// @notice Address of the oracle responsible for freezing validators.
    address public freezeOracle;

    /// @notice Fee required to unfreeze a validator.
    uint256 public unfreezeFee;

    /**
     * @notice Number of blocks a validator must remain frozen before it can be unfrozen.
     * This is param is optional to allow frozen validators to pay the fee immediately.
     */
    uint256 public unfreezePeriodBlocks;

    /// @notice Number of blocks an operator must wait after requesting deregistration before it is finalized.
    uint256 public operatorDeregistrationPeriodBlocks;

    /// @notice Number of blocks a validator must wait after requesting deregistration before it is finalized.
    uint256 public validatorDeregistrationPeriodBlocks;
}
