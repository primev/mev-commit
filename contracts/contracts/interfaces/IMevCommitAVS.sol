// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

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

    /// @notice Enum for LST restaker registration status with MevCommitAVS
    enum LSTRestakerRegistrationStatus {
        // LST restaker is not registered (hasn't chosen a validator) with MevCommitAVS
        NOT_REGISTERED,
        // LST restaker is registered (has chosen a validator) with MevCommitAVS
        REGISTERED,
        // LST restaker has requested deregistration with MevCommitAVS
        REQ_DEREGISTRATION
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

    /// @notice Struct representing MevCommitAVS registration info for a LST restaker
    struct LSTRestakerRegistrationInfo {
        // Status of the LST restaker's registration with MevCommitAVS
        LSTRestakerRegistrationStatus status;
        // Address of the validator chosen by the LST restaker, to represent the restaker's delegation
        bytes chosenValidator;
        // Height at which the LST restaker requested deregistration. Only non-zero if status is REQ_DEREGISTRATION
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

    /// @notice Emmitted when a LST restaker registers (chooses a validator) with MevCommitAVS
    event LSTRestakerRegistered(bytes indexed chosenValidator, address indexed lstRestaker);

    /// @notice Emmitted when a deregistration request is made by an LST restaker
    event LSTRestakerDeregistrationRequested(bytes indexed chosenValidator, address indexed lstRestaker);

    /// @notice Emmitted when a LST restaker is deregistered from MevCommitAVS
    event LSTRestakerDeregistered(bytes indexed chosenValidator, address indexed lstRestaker);

    /// @notice Emmitted when a validator is frozen by the oracle
    event ValidatorFrozen(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is unfrozen
    event ValidatorUnfrozen(bytes indexed validatorPubKey, address indexed podOwner);

    /// @notice Emitted when the AVS directory is set
    event AVSDirectorySet(address indexed avsDirectory);

    /// @notice Emitted when the strategy manager is set
    event StrategyManagerSet(address indexed strategyManager);

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

    /// @notice Emitted when the operator deregistration period is set
    event OperatorDeregistrationPeriodBlocksSet(uint256 operatorDeregistrationPeriodBlocks);

    /// @notice Emitted when the validator deregistration period is set
    event ValidatorDeregistrationPeriodBlocksSet(uint256 validatorDeregistrationPeriodBlocks);

    /// @notice Emitted when the LST restaker deregistration period is set
    event LSTRestakerDeregistrationPeriodBlocksSet(uint256 lstRestakerDeregistrationPeriodBlocks);

    /// @notice Emitted when the max LST restakers per validator is set
    event MaxLSTRestakersPerValidatorSet(uint256 maxLSTRestakersPerValidator);

    function registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature) external;

    function requestOperatorDeregistration(address operator) external;

    function deregisterOperator(address operator) external;

    function registerValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external;

    function requestValidatorsDeregistration(bytes[] calldata valPubKeys) external;

    function deregisterValidators(bytes[] calldata valPubKeys) external;

    function freeze(bytes[] calldata valPubKey) external;

    function unfreeze(bytes calldata valPubKey) payable external;

    function pause() external;

    function unpause() external;

    function areValidatorsOptedIn(bytes[] calldata valPubKeys) external view returns (bool[] memory);

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);

    function avsDirectory() external view returns (address);
}