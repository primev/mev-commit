// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {BlockHeightOccurrence} from "../utils/Occurrence.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";

interface IMevCommitAVS {

    /// @notice Struct representing MevCommitAVS registration info for an operator
    struct OperatorRegistrationInfo {
        /// @notice Whether the operator is registered with MevCommitAVS
        bool exists;
        /// @notice Block height at which the operator possibly requested deregistration
        BlockHeightOccurrence.Occurrence deregRequestOccurrence;
    }

    /// @notice Struct representing MevCommitAVS registration info for a validator
    struct ValidatorRegistrationInfo {
        /// @notice Whether the validator is registered with MevCommitAVS
        bool exists;
        /// @notice Address of the pod owner for the validator
        address podOwner;
        /// @notice Block height at which the validator was possibly frozen
        BlockHeightOccurrence.Occurrence freezeOccurrence;
        /// @notice Block height at which the validator possibly requested deregistration
        BlockHeightOccurrence.Occurrence deregRequestOccurrence;
    }

    /// @notice Struct representing MevCommitAVS registration info for a LST restaker
    struct LSTRestakerRegistrationInfo {
        /// @notice Whether the LST restaker is registered with MevCommitAVS
        bool exists;
        /// @notice Address of validator(s) chosen by the LST restaker, which equally represent the restaker
        bytes[] chosenValidators;
        /// @notice Total number of validators chosen by the LST restaker, where attribution is split evenly
        uint256 numChosen;
        /// @notice Block height at which the LST restaker possibly requested deregistration
        BlockHeightOccurrence.Occurrence deregRequestOccurrence;
    }

    /// @notice Emmitted when an operator is registered with MevCommitAVS
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when a deregistration request is made for an operator
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered from MevCommitAVS
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when a validator is registered with MevCommitAVS
    event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a deregistration request is made for a validator
    event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is deregistered from MevCommitAVS
    event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a LST restaker registers (chooses a validator) with MevCommitAVS
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker);

    /// @notice Emmitted when a deregistration request is made by an LST restaker
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker);

    /// @notice Emmitted when a LST restaker is deregistered from MevCommitAVS
    /// @dev numChosen is the total number of validators chosen by the LST restaker, where attribution is split evenly.
    event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker);

    /// @notice Emmitted when a validator is frozen by the oracle
    event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner);

    /// @notice Emmitted when a validator is unfrozen
    event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner);

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

    error OperatorNotRegistered(address operator);
    error SenderIsRegisteredOperator();
    error ValidatorNotRegistered(bytes valPubKey);
    error ValidatorIsRegistered(bytes valPubKey);
    error LstRestakerNotRegistered();
    error LstRestakerIsRegistered();
    error SenderIsNotFreezeOracle();
    error SenderIsNotEigenCoreOperator();
    error SenderIsNotSpecifiedOperator(address operator);
    error SenderNotPodOwnerOrOperator(address podOwner);
    error SenderNotPodOwnerOrOperatorOfValidator(bytes valPubKey);
    error NoDelegationToRegisteredOperator();
    error ValidatorNotFrozen(bytes valPubKey);
    error UnfreezeFeeRequired(uint256 requiredFee);
    error UnfreezeTransferFailed();
    error RefundFailed();
    error OperatorDeregAlreadyRequested();
    error NoPodExists(address podOwner);
    error ValidatorNotActiveWithEigenCore(bytes valPubKey);
    error ValidatorDeregAlreadyRequested();
    error FrozenValidatorCannotDeregister();
    error DeregistrationNotRequested();
    error DeregistrationTooSoon();
    error NeedChosenValidators();
    error NoEigenStrategyDeposits();
    error DeregistrationAlreadyRequested();
    error ValidatorAlreadyFrozen();
    error UnfreezeTooSoon();

    /// @dev Registers an operator with the MevCommitAVS.
    function registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature) external;

    /// @dev Allows an operator to request deregistration from the MevCommitAVS.
    function requestOperatorDeregistration(address operator) external;

    /// @dev Allows an operator to deregister from the MevCommitAVS.
    function deregisterOperator(address operator) external;

    /// @dev Registers sets of validator pubkeys associated to one or more pod owners.
    function registerValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external;

    /// @dev Allows a validator to request deregistration from the MevCommitAVS.
    function requestValidatorsDeregistration(bytes[] calldata valPubKeys) external;

    /// @dev Allows a validator to deregister from the MevCommitAVS.
    function deregisterValidators(bytes[] calldata valPubKeys) external;

    /// @dev Registers sender as an LST restaker with chosen validators.
    function registerLSTRestaker(bytes[] calldata chosenValidators) external;

    /// @dev Allows an LST restaker to request deregistration from the MevCommitAVS.
    function requestLSTRestakerDeregistration() external;

    /// @dev Allows an LST restaker to deregister from the MevCommitAVS.
    function deregisterLSTRestaker() external;

    /// @dev Allows the freeze oracle account to freeze validators which disobey the mev-commit protocol.
    function freeze(bytes[] calldata valPubKeys) external;

    /// @dev Allows any account to unfreeze validators which have been frozen, for a fee.
    function unfreeze(bytes[] calldata valPubKeys) external payable;

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external;

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external;

    /// @dev Sets the AVS directory, restricted to contract owner.
    function setAVSDirectory(IAVSDirectory avsDirectory_) external;

    /// @dev Sets the strategy manager, restricted to contract owner.
    function setStrategyManager(IStrategyManager strategyManager_) external;

    /// @dev Sets the delegation manager, restricted to contract owner.
    function setDelegationManager(IDelegationManager delegationManager_) external;

    /// @dev Sets the EigenPod manager, restricted to contract owner.
    function setEigenPodManager(IEigenPodManager eigenPodManager_) external;

    /// @dev Sets the restakeable strategies, restricted to contract owner.
    function setRestakeableStrategies(address[] calldata restakeableStrategies_) external;

    /// @dev Sets the freeze oracle account, restricted to contract owner.
    function setFreezeOracle(address freezeOracle_) external;

    /// @dev Sets the unfreeze fee, restricted to contract owner.
    function setUnfreezeFee(uint256 unfreezeFee_) external;

    /// @dev Sets the unfreeze receiver, restricted to contract owner.
    function setUnfreezeReceiver(address unfreezeReceiver_) external;

    /// @dev Sets the unfreeze period in blocks, restricted to contract owner.
    function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) external;

    /// @dev Sets the operator deregistration period in blocks, restricted to contract owner.
    function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) external;

    /// @dev Sets the validator deregistration period in blocks, restricted to contract owner.
    function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) external;

    /// @dev Sets the LST restaker deregistration period in blocks, restricted to contract owner.
    function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) external;

    /// @dev Updates the eigenlayer metadata URI, restricted to contract owner.
    function updateMetadataURI(string memory metadataURI_) external;

    /// @dev Checks if a validator is opted-in.
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);

    /// @dev Returns operator registration info.
    function getOperatorRegInfo(address operator) external view returns (OperatorRegistrationInfo memory);

    /// @dev Returns validator registration info.
    function getValidatorRegInfo(bytes calldata valPubKey) external view returns (ValidatorRegistrationInfo memory);

    /// @dev Returns LST restaker registration info.
    function getLSTRestakerRegInfo(address lstRestaker) external view returns (LSTRestakerRegistrationInfo memory);

    /// @dev Returns the address of AVS directory.
    function avsDirectory() external view returns (address);
}
