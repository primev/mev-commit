// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {EventHeightLib} from "../utils/EventHeight.sol";

interface IMevCommitAVS {

    /// @notice Struct representing MevCommitAVS registration info for an operator
    struct OperatorRegistrationInfo {
        /// @notice Whether the operator is registered with MevCommitAVS
        bool exists;
        /// @notice Height at which the operator possibly requested deregistration
        EventHeightLib.EventHeight deregRequestHeight;
    }

    /// @notice Struct representing MevCommitAVS registration info for a validator
    struct ValidatorRegistrationInfo {
        /// @notice Whether the validator is registered with MevCommitAVS
        bool exists;
        /// @notice Address of the pod owner for the validator
        address podOwner;
        /// @notice Height at which the validator was possibly frozen
        EventHeightLib.EventHeight freezeHeight;
        /// @notice Height at which the validator possibly requested deregistration
        EventHeightLib.EventHeight deregRequestHeight;
    }

    /// @notice Struct representing MevCommitAVS registration info for a LST restaker
    struct LSTRestakerRegistrationInfo {
        /// @notice Whether the LST restaker is registered with MevCommitAVS
        bool exists;
        /// @notice Address of validator(s) chosen by the LST restaker, which equally represent the restaker
        bytes[] chosenValidators;
        /// @notice Total number of validators chosen by the LST restaker, where attribution is split evenly
        uint256 numChosen;
        /// @notice Height at which the LST restaker possibly requested deregistration
        EventHeightLib.EventHeight deregRequestHeight;
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
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerRegistered(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);

    /// @notice Emmitted when a deregistration request is made by an LST restaker
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerDeregistrationRequested(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);

    /// @notice Emmitted when a LST restaker is deregistered from MevCommitAVS
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerDeregistered(bytes indexed chosenValidator, uint256 numChosen, address indexed lstRestaker);

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

    /// @notice Emitted when the unfreeze receiver is set
    event UnfreezeReceiverSet(address indexed unfreezeReceiver);

    /// @notice Emitted when the unfreeze period is set
    event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks);

    /// @notice Emitted when the operator deregistration period is set
    event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks);

    /// @notice Emitted when the validator deregistration period is set
    event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks);

    /// @notice Emitted when the LST restaker deregistration period is set
    event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks);

    /// @notice Emitted when the max LST restakers per validator is set
    event MaxLSTRestakersPerValidatorSet(uint256 maxLSTRestakersPerValidator);

    function registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature) external;
    function requestOperatorDeregistration(address operator) external;
    function deregisterOperator(address operator) external;
    function registerValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external;
    function requestValidatorsDeregistration(bytes[] calldata valPubKeys) external;
    function deregisterValidators(bytes[] calldata valPubKeys) external;
    function registerLSTRestaker(bytes[] calldata chosenValidators) external;
    function requestLSTRestakerDeregistration() external;
    function deregisterLSTRestaker() external;
    function freeze(bytes[] calldata valPubKeys) external;
    function unfreeze(bytes[] calldata valPubKeys) payable external;
    function pause() external;
    function unpause() external;
    function areValidatorsOptedIn(bytes[] calldata valPubKeys) external view returns (bool[] memory);
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);
    function getLSTRestakerRegInfo(address lstRestaker) external view returns (LSTRestakerRegistrationInfo memory);
    function avsDirectory() external view returns (address);
}