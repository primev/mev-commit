// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {IMevCommitAVS} from "../../interfaces/IMevCommitAVS.sol";
import {MevCommitAVSStorage} from "./MevCommitAVSStorage.sol";
import {EventHeightLib} from "../../utils/EventHeight.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IEigenPod} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";

/// @title MevCommitAVS
/// @notice This contract serves as the entrypoint for operators, validators and LST restakers to register with
/// the mev-commit protocol via an eigenlayer AVS.
contract MevCommitAVS is IMevCommitAVS, MevCommitAVSStorage,
    OwnableUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    
    /// @dev Modifier to ensure the provided operator is registered with MevCommitAVS.
    modifier onlyRegisteredOperator(address operator) {
        require(operatorRegistrations[operator].exists, "operator must be registered");
        _;
    }

    /// @dev Modifier to ensure the sender is not a registered operator with MevCommitAVS.
    modifier onlyNonRegisteredOperator() {
        require(!operatorRegistrations[msg.sender].exists, "sender must not be registered operator");
        _;
    }

    /// @dev Modifier to ensure all provided validators are registered with MevCommitAVS.
    modifier onlyRegisteredValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(validatorRegistrations[valPubKeys[i]].exists, "validator must be registered");
        }
        _;
    }

    /// @dev Modifier to ensure all provided validators are not registered with MevCommitAVS.
    modifier onlyNonRegisteredValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(!validatorRegistrations[valPubKeys[i]].exists, "validator must not be registered");
        }
        _;
    }

    /// @dev Modifier to ensure the sender is a registered LST restaker with MevCommitAVS.
    modifier onlyRegisteredLstRestaker() {
        require(lstRestakerRegistrations[msg.sender].exists, "sender must be registered LST restaker");
        _;
    }

    /// @dev Modifier to ensure the sender is not a registered LST restaker with MevCommitAVS.
    modifier onlyNonRegisteredLstRestaker() {
        require(!lstRestakerRegistrations[msg.sender].exists, "sender must not be registered LST restaker");
        _;
    }

    /// @dev Modifier to ensure the sender is the MevCommitAVS freeze oracle account.
    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, "sender must be freeze oracle");
        _;
    }

    /// @dev Modifier to ensure the sender is registered as an operator with the eigenlayer core contracts.
    modifier onlyEigenCoreOperator() {
        require(_delegationManager.isOperator(msg.sender), "sender must be an eigenlayer operator");
        _;
    }
    
    /// @dev Modifier to ensure the sender is the given operator 
    modifier onlyOperator(address operator) {
        require(msg.sender == operator, "sender must be operator");
        _;
    }

    /// @dev Modifier to ensure the sender is either the given pod owner, 
    /// or the delegated operator for the given pod owner.
    modifier onlyPodOwnerOrOperator(address podOwner) {
        require(msg.sender == podOwner || msg.sender == _delegationManager.delegatedTo(podOwner), 
            "sender must be podOwner or delegated operator");
        _;
    }

    /// @dev Modifier to ensure the sender is either the pod owner or operator of all the given validators.
    modifier onlyPodOwnerOrOperatorOfValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            IMevCommitAVS.ValidatorRegistrationInfo memory regInfo = validatorRegistrations[valPubKeys[i]];
            require(msg.sender == regInfo.podOwner || msg.sender == _delegationManager.delegatedTo(regInfo.podOwner),
                "sender must be podOwner or delegated operator of validator");
        }
        _;
    }

    /// @dev Modifier to ensure the sender is delegated to a registered operator.
    modifier onlySenderWithRegisteredOperator() {
        address delegatedOperator = _delegationManager.delegatedTo(msg.sender);
        require(operatorRegistrations[delegatedOperator].exists,
            "sender must be delegated to an operator that is registered with MevCommitAVS");
        _;
    }

    /// @dev Modifier to ensure all provided validators are frozen.
    modifier onlyFrozenValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(validatorRegistrations[valPubKeys[i]].freezeHeight.exists, "validator must be frozen");
        }
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Initializes the contract with provided parameters.
    function initialize(
        address owner_,
        IDelegationManager delegationManager_,
        IEigenPodManager eigenPodManager_,
        IStrategyManager strategyManager_,
        IAVSDirectory avsDirectory_,
        address[] calldata restakeableStrategies_,
        address freezeOracle_,
        uint256 unfreezeFee_,
        address unfreezeReceiver_,
        uint256 unfreezePeriodBlocks_,
        uint256 operatorDeregPeriodBlocks_,
        uint256 validatorDeregPeriodBlocks_,
        uint256 lstRestakerDeregPeriodBlocks_,
        string calldata metadataURI_
    ) external initializer {
        _setDelegationManager(delegationManager_);
        _setEigenPodManager(eigenPodManager_);
        _setStrategyManager(strategyManager_);
        _setAVSDirectory(avsDirectory_);
        _setRestakeableStrategies(restakeableStrategies_);
        _setFreezeOracle(freezeOracle_);
        _setUnfreezeFee(unfreezeFee_);
        _setUnfreezeReceiver(unfreezeReceiver_);
        _setUnfreezePeriodBlocks(unfreezePeriodBlocks_);
        _setOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_);
        _setValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_);
        _setLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_);
        if (bytes(metadataURI_).length > 0) {
            _updateMetadataURI(metadataURI_);
        }
        __Ownable_init(owner_);
        __UUPSUpgradeable_init();
        __Pausable_init();
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    /// @dev Registers an operator with the MevCommitAVS.
    function registerOperator (
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external whenNotPaused() onlyNonRegisteredOperator() onlyEigenCoreOperator() {
        _registerOperator(operatorSignature);
    }

    /// @dev Allows an operator to request deregistration from the MevCommitAVS.
    function requestOperatorDeregistration(address operator
    ) external whenNotPaused() onlyRegisteredOperator(operator) onlyOperator(operator) {
        _requestOperatorDeregistration(operator);
    }

    /// @dev Allows an operator to deregister from the MevCommitAVS.
    function deregisterOperator(address operator
    ) external whenNotPaused() onlyRegisteredOperator(operator) onlyOperator(operator) {
        _deregisterOperator(operator);
    }

    /// @dev Registers sets of validator pubkeys associated to one or more pod owners.
    /// @notice The underlying _registerValidatorsByPodOwner enforces the sender is either
    /// the provided pod owner, or the delegated operator for each pod owner.
    function registerValidatorsByPodOwners(
        bytes[][] calldata valPubKeys,
        address[] calldata podOwners
    ) external whenNotPaused() {
        for (uint256 i = 0; i < podOwners.length; i++) {
            _registerValidatorsByPodOwner(valPubKeys[i], podOwners[i]);
        }
    }

    /// @dev Allows a validator to request deregistration from the MevCommitAVS.
    /// @notice For each validator the underlying _requestValidatorDeregistration enforces the sender is either
    /// the podOwner, delegated operator, or the contract owner.
    function requestValidatorsDeregistration(bytes[] calldata valPubKeys)
        external whenNotPaused() onlyRegisteredValidators(valPubKeys) onlyPodOwnerOrOperatorOfValidators(valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _requestValidatorDeregistration(valPubKeys[i]);
        }
    }

    /// @dev Allows a validator to deregister from the MevCommitAVS.
    /// @notice For each validator the underlying _deregisterValidator enforces the sender is either
    /// the podOwner, delegated operator, or the contract owner.
    function deregisterValidators(bytes[] calldata valPubKeys)
        external whenNotPaused() onlyRegisteredValidators(valPubKeys) onlyPodOwnerOrOperatorOfValidators(valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _deregisterValidator(valPubKeys[i]);
        }
    }

    /// @dev Registers sender as an LST restaker with chosen validators.
    function registerLSTRestaker(bytes[] calldata chosenValidators)
        external whenNotPaused() onlyNonRegisteredLstRestaker() onlySenderWithRegisteredOperator() {
        _registerLSTRestaker(chosenValidators);
    }

    /// @dev Allows an LST restaker to request deregistration from the MevCommitAVS.
    function requestLSTRestakerDeregistration() external whenNotPaused() onlyRegisteredLstRestaker() {
        _requestLSTRestakerDeregistration();
    }

    /// @dev Allows an LST restaker to deregister from the MevCommitAVS.
    function deregisterLSTRestaker() external whenNotPaused() onlyRegisteredLstRestaker() {
        _deregisterLSTRestaker();
    }

    /// @dev Allows the freeze oracle account to freeze validators which disobey the mev-commit protocol.
    function freeze(bytes[] calldata valPubKeys) external
        whenNotPaused() onlyRegisteredValidators(valPubKeys) onlyFreezeOracle() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _freeze(valPubKeys[i]);
        }
    }

    /// @dev Allows any account to unfreeze validators which have been frozen, for a fee.
    function unfreeze(bytes[] calldata valPubKey) payable external 
        whenNotPaused() onlyRegisteredValidators(valPubKey) onlyFrozenValidators(valPubKey) {
        uint256 requiredFee = unfreezeFee * valPubKey.length;
        require(msg.value >= requiredFee,
            "sender must pay at least the unfreeze fee for each validator");
        for (uint256 i = 0; i < valPubKey.length; i++) {
            _unfreeze(valPubKey[i]);
            payable(unfreezeReceiver).transfer(unfreezeFee);
        }
        uint256 excessFee = msg.value - requiredFee;
        if (excessFee > 0) {
            payable(msg.sender).transfer(excessFee);
        }
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @dev Sets the AVS directory, restricted to contract owner.
    function setAVSDirectory(IAVSDirectory avsDirectory_) external onlyOwner {
        _setAVSDirectory(avsDirectory_);
    }

    /// @dev Sets the strategy manager, restricted to contract owner.
    function setStrategyManager(IStrategyManager strategyManager_) external onlyOwner {
        _setStrategyManager(strategyManager_);
    }

    /// @dev Sets the delegation manager, restricted to contract owner.
    function setDelegationManager(IDelegationManager delegationManager_) external onlyOwner {
        _setDelegationManager(delegationManager_);
    }

    /// @dev Sets the EigenPod manager, restricted to contract owner.
    function setEigenPodManager(IEigenPodManager eigenPodManager_) external onlyOwner {
        _setEigenPodManager(eigenPodManager_);
    }

    /// @dev Sets the restakeable strategies, restricted to contract owner.
    function setRestakeableStrategies(address[] calldata restakeableStrategies_) external onlyOwner {
        _setRestakeableStrategies(restakeableStrategies_);
    }

    /// @dev Sets the freeze oracle account, restricted to contract owner.
    function setFreezeOracle(address freezeOracle_) external onlyOwner {
        _setFreezeOracle(freezeOracle_);
    }

    /// @dev Sets the unfreeze fee, restricted to contract owner.
    function setUnfreezeFee(uint256 unfreezeFee_) external onlyOwner {
        _setUnfreezeFee(unfreezeFee_);
    }

    /// @dev Sets the unfreeze receiver, restricted to contract owner.
    function setUnfreezeReceiver(address unfreezeReceiver_) external onlyOwner {
        _setUnfreezeReceiver(unfreezeReceiver_);
    }

    /// @dev Sets the unfreeze period in blocks, restricted to contract owner.
    function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) external onlyOwner {
        _setUnfreezePeriodBlocks(unfreezePeriodBlocks_);
    }

    /// @dev Sets the operator deregistration period inblocks, restricted to contract owner.
    function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) external onlyOwner {
        _setOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_);
    }

    /// @dev Sets the validator deregistration period in blocks, restricted to contract owner.
    function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) external onlyOwner {
        _setValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_);
    }

    /// @dev Sets the LST restaker deregistration period in blocks, restricted to contract owner.
    function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) external onlyOwner {
        _setLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_);
    }

    /// @dev Updates the eigenlayer metadata URI, restricted to contract owner.
    function updateMetadataURI(string memory metadataURI_) external onlyOwner {
        _updateMetadataURI(metadataURI_);
    }

    /// @dev Internal function to register an operator.
    function _registerOperator(ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature) internal {
        _eigenAVSDirectory.registerOperatorToAVS(msg.sender, operatorSignature);
        operatorRegistrations[msg.sender] = OperatorRegistrationInfo({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            })
        });
        emit OperatorRegistered(msg.sender);
    }

    /// @dev Internal function to request deregistration of an operator.
    function _requestOperatorDeregistration(address operator) internal {
        require(!operatorRegistrations[operator].deregRequestHeight.exists,
            "operator must not have already requested deregistration");
        EventHeightLib.set(operatorRegistrations[operator].deregRequestHeight, block.number);
        emit OperatorDeregistrationRequested(operator);
    }

    /// @dev Internal function to deregister an operator.
    function _deregisterOperator(address operator) internal {
        require(operatorRegistrations[operator].deregRequestHeight.exists, "operator must have requested deregistration");
        require(block.number >= operatorRegistrations[operator].deregRequestHeight.blockHeight + operatorDeregPeriodBlocks,
            "deregistration must happen at least operatorDeregPeriodBlocks after deregistration request height");
        _eigenAVSDirectory.deregisterOperatorFromAVS(operator);
        delete operatorRegistrations[operator];
        emit OperatorDeregistered(operator);
    }

    /// @dev Internal function to register validators by their pod owner.
    /// @notice Invalid pubkeys should not correspond to VALIDATOR_STATUS.ACTIVE due to validations in EigenPod.sol
    function _registerValidatorsByPodOwner(
        bytes[] calldata valPubKeys,
        address podOwner
    ) internal onlyNonRegisteredValidators(valPubKeys) onlyPodOwnerOrOperator(podOwner)  {
        address operator = _delegationManager.delegatedTo(podOwner);
        require(operatorRegistrations[operator].exists,
            "delegated operator must be registered with MevCommitAVS");
        require(!operatorRegistrations[operator].deregRequestHeight.exists,
            "delegated operator must not have requested deregistration");
        IEigenPod pod = _eigenPodManager.getPod(podOwner);
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE,
                "validator must be active under pod");
            _registerValidator(valPubKeys[i], podOwner);
        }
    }

    /// @dev Internal function to register a validator.
    function _registerValidator(bytes calldata valPubKey, address podOwner) internal {
        validatorRegistrations[valPubKey] = ValidatorRegistrationInfo({
            exists: true,
            podOwner: podOwner,
            freezeHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            })
        });
        emit ValidatorRegistered(valPubKey, podOwner);
    }

    /// @dev Internal function to request deregistration of a validator.
    function _requestValidatorDeregistration(bytes calldata valPubKey) internal {
        require(!validatorRegistrations[valPubKey].deregRequestHeight.exists,
            "validator must not have already requested deregistration");
        EventHeightLib.set(validatorRegistrations[valPubKey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    /// @dev Internal function to deregister a validator.
    function _deregisterValidator(bytes calldata valPubKey) internal {
        require(!validatorRegistrations[valPubKey].freezeHeight.exists, "frozen validator cannot deregister");
        require(validatorRegistrations[valPubKey].deregRequestHeight.exists,
            "validator must have requested deregistration");
        require(block.number >= validatorRegistrations[valPubKey].deregRequestHeight.blockHeight + validatorDeregPeriodBlocks,
            "deregistration must happen at least validatorDeregPeriodBlocks after deregistration request height");
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        delete validatorRegistrations[valPubKey];
        emit ValidatorDeregistered(valPubKey, podOwner);
    }

    /// @dev Internal function to register an LST restaker.
    function _registerLSTRestaker(bytes[] calldata chosenValidators) internal {
        require(chosenValidators.length > 0, "LST restaker must choose at least one validator");
        uint256 stratLen = _strategyManager.stakerStrategyListLength(msg.sender);
        require(stratLen > 0, "LST restaker must have deposited into at least one strategy");
        lstRestakerRegistrations[msg.sender] = LSTRestakerRegistrationInfo({
            exists: true,
            chosenValidators: chosenValidators,
            numChosen: chosenValidators.length,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            })
        });
        for (uint256 i = 0; i < chosenValidators.length; i++) {
            emit LSTRestakerRegistered(chosenValidators[i], chosenValidators.length, msg.sender);
        }
    }

    /// @dev Internal function to request deregistration of an LST restaker.
    function _requestLSTRestakerDeregistration() internal {
        LSTRestakerRegistrationInfo storage reg = lstRestakerRegistrations[msg.sender];
        require(!reg.deregRequestHeight.exists, "LST restaker must not have already requested deregistration");
        EventHeightLib.set(reg.deregRequestHeight, block.number);
        for (uint256 i = 0; i < reg.numChosen; i++) {
            emit LSTRestakerDeregistrationRequested(reg.chosenValidators[i], reg.numChosen, msg.sender);
        }
    }

    /// @dev Internal function to deregister an LST restaker.
    function _deregisterLSTRestaker() internal {
        LSTRestakerRegistrationInfo storage reg = lstRestakerRegistrations[msg.sender];
        require(reg.deregRequestHeight.exists, "LST restaker must have requested deregistration");
        require(block.number >= reg.deregRequestHeight.blockHeight + lstRestakerDeregPeriodBlocks,
            "deregistration must happen at least lstRestakerDeregPeriodBlocks after deregistration request height");
        for (uint256 i = 0; i < reg.numChosen; i++) {
            emit LSTRestakerDeregistered(reg.chosenValidators[i], reg.numChosen, msg.sender);
        }
        delete lstRestakerRegistrations[msg.sender];
    }

    /// @dev Internal function to freeze a validator.
    function _freeze(bytes calldata valPubKey) internal {
        require(!validatorRegistrations[valPubKey].freezeHeight.exists, "validator must not already be frozen");
        EventHeightLib.set(validatorRegistrations[valPubKey].freezeHeight, block.number);
        emit ValidatorFrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    /// @dev Internal function to unfreeze a validator.
    function _unfreeze(bytes calldata valPubKey) internal {
        require(block.number >= validatorRegistrations[valPubKey].freezeHeight.blockHeight + unfreezePeriodBlocks,
            "unfreeze must happen at least unfreezePeriodBlocks after freeze height");
        EventHeightLib.del(validatorRegistrations[valPubKey].freezeHeight);
        emit ValidatorUnfrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    /// @dev Internal function to set the AVS directory.
    function _setAVSDirectory(IAVSDirectory avsDirectory_) internal {
        _eigenAVSDirectory = avsDirectory_;
        emit AVSDirectorySet(address(_eigenAVSDirectory));
    }

    /// @dev Internal function to set the strategy manager.
    function _setStrategyManager(IStrategyManager strategyManager_) internal {
        _strategyManager = strategyManager_;
        emit StrategyManagerSet(address(strategyManager_));
    }

    /// @dev Internal function to set the delegation manager.
    function _setDelegationManager(IDelegationManager delegationManager_) internal {
        _delegationManager = delegationManager_;
        emit DelegationManagerSet(address(delegationManager_));
    }

    /// @dev Internal function to set the EigenPod manager.
    function _setEigenPodManager(IEigenPodManager eigenPodManager_) internal {
        _eigenPodManager = eigenPodManager_;
        emit EigenPodManagerSet(address(eigenPodManager_));
    }

    /// @dev Internal function to set the restakeable strategies.
    function _setRestakeableStrategies(address[] calldata restakeableStrategies_) internal {
        restakeableStrategies = restakeableStrategies_;
        emit RestakeableStrategiesSet(restakeableStrategies);
    }

    /// @dev Internal function to set the freeze oracle account.
    function _setFreezeOracle(address _freezeOracle) internal {
        freezeOracle = _freezeOracle;
        emit FreezeOracleSet(_freezeOracle);
    }

    /// @dev Internal function to set the unfreeze fee.
    function _setUnfreezeFee(uint256 _unfreezeFee) internal {
        unfreezeFee = _unfreezeFee;
        emit UnfreezeFeeSet(_unfreezeFee);
    }

    /// @dev Internal function to set the unfreeze receiver.
    function _setUnfreezeReceiver(address _unfreezeReceiver) internal {
        unfreezeReceiver = _unfreezeReceiver;
        emit UnfreezeReceiverSet(_unfreezeReceiver);
    }

    /// @dev Internal function to set the unfreeze period in blocks.
    function _setUnfreezePeriodBlocks(uint256 _unfreezePeriodBlocks) internal {
        unfreezePeriodBlocks = _unfreezePeriodBlocks;
        emit UnfreezePeriodBlocksSet(_unfreezePeriodBlocks);
    }
    
    /// @dev Internal function to set the operator deregistration period in blocks.
    function _setOperatorDeregPeriodBlocks(uint256 _operatorDeregPeriodBlocks) internal {
        operatorDeregPeriodBlocks = _operatorDeregPeriodBlocks;
        emit OperatorDeregPeriodBlocksSet(_operatorDeregPeriodBlocks);
    }

    /// @dev Internal function to set the validator deregistration period in blocks.
    function _setValidatorDeregPeriodBlocks(uint256 _validatorDeregPeriodBlocks) internal {
        validatorDeregPeriodBlocks = _validatorDeregPeriodBlocks;
        emit ValidatorDeregPeriodBlocksSet(_validatorDeregPeriodBlocks);
    }

    /// @dev Internal function to set the LST restaker deregistration period in blocks.
    function _setLstRestakerDeregPeriodBlocks(uint256 _lstRestakerDeregPeriodBlocks) internal {
        lstRestakerDeregPeriodBlocks = _lstRestakerDeregPeriodBlocks;
        emit LSTRestakerDeregPeriodBlocksSet(_lstRestakerDeregPeriodBlocks);
    }

    /// @dev Internal function to update the eigenlayer metadata URI.
    function _updateMetadataURI(string memory _metadataURI) internal {
        _eigenAVSDirectory.updateAVSMetadataURI(_metadataURI);
    }

    /// @dev Returns the list of restakeable strategies.
    function getRestakeableStrategies() external view returns (address[] memory) {
        return _getRestakeableStrategies();
    }

    /// @dev Returns the restakeable strategies for a given operator.
    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
        if (!operatorRegistrations[operator].exists) {
            return new address[](0);
        }
        return _getRestakeableStrategies();
    }

    /// @dev Checks if a validator is opted-in.
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valPubKey);
    }

    /// @dev Returns operator registration info.
    function getOperatorRegInfo(address operator) external view returns (OperatorRegistrationInfo memory) {
        return operatorRegistrations[operator];
    }

    /// @dev Returns validator registration info.
    function getValidatorRegInfo(bytes calldata valPubKey) external view returns (ValidatorRegistrationInfo memory) {
        return validatorRegistrations[valPubKey];
    }

    /// @dev Returns LST restaker registration info.
    function getLSTRestakerRegInfo(address lstRestaker) 
        external view returns (LSTRestakerRegistrationInfo memory) {
        return lstRestakerRegistrations[lstRestaker];
    }

    /// @dev Returns the address of AVS directory.
    function avsDirectory() external view returns (address) {
        return address(_eigenAVSDirectory);
    }

    /// @dev Internal function to check if a validator is opted-in.
    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        bool isValRegistered = validatorRegistrations[valPubKey].exists;
        bool isFrozen = validatorRegistrations[valPubKey].freezeHeight.exists;
        bool isValDeregRequested = validatorRegistrations[valPubKey].deregRequestHeight.exists;
        IEigenPod pod = _eigenPodManager.getPod(validatorRegistrations[valPubKey].podOwner);
        bool isValActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address delegatedOperator = _delegationManager.delegatedTo(validatorRegistrations[valPubKey].podOwner);
        bool isOperatorRegistered = operatorRegistrations[delegatedOperator].exists;
        bool isOperatorDeregRequested = operatorRegistrations[delegatedOperator].deregRequestHeight.exists;

        return isValRegistered && !isFrozen && !isValDeregRequested && isValActive
            && isOperatorRegistered && !isOperatorDeregRequested;
    }

    /// @dev Internal function to get the list of restakeable strategies.
    function _getRestakeableStrategies() internal view returns (address[] memory) {
        return restakeableStrategies;
    }

    /// @dev Fallback function to prevent unintended contract interactions.
    fallback() external payable {
        revert("Invalid call");
    }

    /// @dev Receive function to prevent unintended contract interactions.
    receive() external payable {
        revert("Invalid call");
    }
}
