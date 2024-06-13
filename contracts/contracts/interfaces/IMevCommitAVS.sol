// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";

interface IMevCommitAVS {

    /// @notice Enum for operator registration status with MevCommitAVS
    enum OperatorRegistrationStatus {
        // Operator is not registered with MevCommitAVS
        NOT_REGISTERED,
        // Operator is registered with MevCommitAVS
        REGISTERED, 
        // Operator has requested deregistration with MevCommitAVS
        REQ_DEREGISTRATION
    }

    /// @notice Enum for validator registration status with MevCommitAVS
    enum ValidatorRegistrationStatus {
        // Validator is not registered with MevCommitAVS
        NOT_REGISTERED,
        // Validator is registered with MevCommitAVS
        REGISTERED,
        // Validator has requested deregistration with MevCommitAVS
        REQ_DEREGISTRATION,
        // Validator is frozen by MevCommitAVS
        FROZEN
    }

    /// @notice Struct representing MevCommitAVS registration info for an operator
    struct OperatorRegistrationInfo {
        // Status of the operator's registration with MevCommitAVS
        OperatorRegistrationStatus status;
        // Height at which the operator requested deregistration. Only non-zero if status is REQ_DEREGISTRATION
        uint256 deregistrationRequestHeight;
    }

    /// @notice Struct representing MevCommitAVS registration info for a validator
    struct ValidatorRegistrationInfo {
        // Status of the validator's registration with MevCommitAVS
        ValidatorRegistrationStatus status;
        // Address of the pod owner of the validator
        address podOwner;
        // Height at which the validator was frozen. Only non-zero if status is FROZEN
        uint256 freezeHeight;
        // Height at which the validator requested deregistration. Only non-zero if status is REQ_DEREGISTRATION
        uint256 deregistrationRequestHeight;
    }

    /// @notice Emmitted when an operator is registered with MevCommitAVS
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when a deregistration request is made for an operator
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered from MevCommitAVS
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when a validator is registered with MevCommitAVS
    event ValidatorRegistered(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a deregistration request is made for a validator
    event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is deregistered from MevCommitAVS
    event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is frozen by the oracle
    event ValidatorFrozen(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is unfrozen
    event ValidatorUnfrozen(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emitted when the AVS directory is set
    event AVSDirectorySet(address indexed avsDirectory);

    /// @notice Emitted when the delegation manager is set
    event DelegationManagerSet(address indexed delegationManager);

    /// @notice Emitted when the EigenPod manager is set
    event EigenPodManagerSet(address indexed eigenPodManager);

    /// @notice Emitted when the restakeable strategies are set
    event RestakeableStrategiesSet(address[] indexed restakeableStrategies);

    /// @notice Emitted when the freeze oracle is set
    event FreezeOracleSet(address indexed freezeOracle);

    /// @notice Emitted when the unfreeze fee is set
    event UnfreezeFeeSet(uint256 unfreezeFee);

    /// @notice Emitted when the unfreeze period is set
    event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks);

    /// @notice Emitted when the operator deregistration period blocks are set
    event OperatorDeregistrationPeriodBlocksSet(uint256 operatorDeregistrationPeriodBlocks);

    /// @notice Emitted when the validator deregistration period blocks are set
    event ValidatorDeregistrationPeriodBlocksSet(uint256 validatorDeregistrationPeriodBlocks);

    function registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature) external;

    function requestOperatorDeregistration(address operator) external;

    function deregisterOperator(address operator) external;

    function registerValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external;

    function requestValidatorsDeregistration(bytes[] calldata valPubKeys) external;

    function deregisterValidators(bytes[] calldata valPubKeys) external;

    function freeze(bytes calldata valPubKey) external;

    function unfreeze(bytes calldata valPubKey) payable external;

    function areValidatorsOptedIn(bytes[] calldata valPubKeys) external view returns (bool[] memory);

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);

    function avsDirectory() external view returns (address);
}