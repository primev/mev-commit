    // SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

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

// TODO: confirm all this still conforms to README
// make all public functions similarily have an internal func
contract MevCommitAVS is IMevCommitAVS, MevCommitAVSStorage,
    OwnableUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    
    modifier onlyRegisteredOperator(address operator) {
        require(operatorRegistrations[operator].exists, "operator must be registered");
        _;
    }

    modifier onlyNonRegisteredOperator() {
        require(!operatorRegistrations[msg.sender].exists, "sender must not be registered operator");
        _;
    }

    modifier onlyRegisteredValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(validatorRegistrations[valPubKeys[i]].exists, "validator must be registered");
        }
        _;
    }

    modifier onlyRegisteredValidator(bytes calldata valPubKey) {
        require(validatorRegistrations[valPubKey].exists, "validator must be registered");
        _;
    }

    modifier onlyNonRegisteredValidators(bytes[] calldata valPubKeys) {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(!validatorRegistrations[valPubKeys[i]].exists, "validator must not be registered");
        }
        _;
    }

    modifier onlyRegisteredLstRestaker() {
        require(lstRestakerRegistrations[msg.sender].exists, "sender must be registered LST restaker");
        _;
    }

    modifier onlyNonRegisteredLstRestaker() {
        require(!lstRestakerRegistrations[msg.sender].exists, "sender must not be registered LST restaker");
        _;
    }

    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, "sender must be freeze oracle");
        _;
    }

    modifier onlyEigenCoreOperator() {
        require(_delegationManager.isOperator(msg.sender), "sender must be an eigenlayer operator");
        _;
    }
    
    modifier onlyOperatorOrContractOwner(address operator) {
        require(msg.sender == operator || msg.sender == owner(), "sender must be operator or MevCommitAVS owner");
        _;
    }

    modifier onlyPodOwnerOrOperator(address podOwner) {
        address delegatedOperator = _delegationManager.delegatedTo(podOwner);
        require(msg.sender == podOwner || msg.sender == delegatedOperator, 
            "sender must be podOwner or delegated operator");
        require(operatorRegistrations[delegatedOperator].exists,
            "delegated operator must be registered with MevCommitAVS");
        _;
    }

    modifier onlyPodOwnerOperatorOrContractOwner(bytes calldata valPubKey) {
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        require(msg.sender == podOwner ||
            msg.sender == _delegationManager.delegatedTo(podOwner) ||
            msg.sender == owner(),
            "sender must be podOwner, delegated operator, or MevCommitAVS owner");
        _;
    }

    modifier onlySenderWithRegisteredOperator() {
        address delegatedOperator = _delegationManager.delegatedTo(msg.sender);
        require(operatorRegistrations[delegatedOperator].exists,
            "sender must be delegated to an operator that is registered with MevCommitAVS");
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address owner_,
        IDelegationManager delegationManager_,
        IEigenPodManager eigenPodManager_,
        IStrategyManager strategyManager_,
        IAVSDirectory avsDirectory_,
        address[] calldata restakeableStrategies_,
        address freezeOracle_,
        uint256 unfreezeFee_,
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

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    function registerOperator (
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external onlyNonRegisteredOperator() onlyEigenCoreOperator() whenNotPaused() {
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

    function requestOperatorDeregistration(address operator
    ) external onlyRegisteredOperator(operator) onlyOperatorOrContractOwner(operator) whenNotPaused() {
        EventHeightLib.set(operatorRegistrations[operator].deregRequestHeight, block.number);
        emit OperatorDeregistrationRequested(operator);
    }

    function deregisterOperator(address operator
    ) external onlyRegisteredOperator(operator) onlyOperatorOrContractOwner(operator) whenNotPaused() {
        require(operatorRegistrations[operator].deregRequestHeight.exists, "operator must have requested deregistration");
        require(block.number >= operatorRegistrations[operator].deregRequestHeight.blockHeight + operatorDeregPeriodBlocks,
            "deregistration must happen at least operatorDeregPeriodBlocks after deregistration request height");
        _eigenAVSDirectory.deregisterOperatorFromAVS(operator);
        delete operatorRegistrations[operator];
        emit OperatorDeregistered(operator);
    }

    function registerValidatorsByPodOwners(
        bytes[][] calldata valPubKeys,
        address[] calldata podOwners
    ) external whenNotPaused() {
        for (uint256 i = 0; i < podOwners.length; i++) {
            _registerValidatorsByPodOwner(valPubKeys[i], podOwners[i]);
        }
    }

    function requestValidatorsDeregistration(bytes[] calldata valPubKeys)
        external onlyRegisteredValidators(valPubKeys) whenNotPaused() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _requestValidatorDeregistration(valPubKeys[i]);
        }
    }

    function deregisterValidators(bytes[] calldata valPubKeys)
        external onlyRegisteredValidators(valPubKeys) whenNotPaused() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _deregisterValidator(valPubKeys[i]);
        }
    }

    function registerLSTRestaker(bytes[] calldata chosenValidators)
        external onlyNonRegisteredLstRestaker() onlySenderWithRegisteredOperator() {
        for (uint256 i = 0; i < chosenValidators.length; i++) {
            require(_isValidatorOptedIn(chosenValidators[i]), "chosen validator must be opted in");
        }
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

    function requestLSTRestakerDeregistration() external onlyRegisteredLstRestaker() {
        LSTRestakerRegistrationInfo storage reg = lstRestakerRegistrations[msg.sender];
        EventHeightLib.set(reg.deregRequestHeight, block.number);
        for (uint256 i = 0; i < reg.numChosen; i++) {
            emit LSTRestakerDeregistrationRequested(reg.chosenValidators[i], reg.numChosen, msg.sender);
        }
    }

    function deregisterLSTRestaker() external onlyRegisteredLstRestaker() {
        LSTRestakerRegistrationInfo storage reg = lstRestakerRegistrations[msg.sender];
        require(reg.deregRequestHeight.exists, "LST restaker must have requested deregistration");
        require(block.number >= reg.deregRequestHeight.blockHeight + lstRestakerDeregPeriodBlocks,
            "deregistration must happen at least lstRestakerDeregPeriodBlocks after deletion request height");
        for (uint256 i = 0; i < reg.numChosen; i++) {
            emit LSTRestakerDeregistered(reg.chosenValidators[i], reg.numChosen, msg.sender);
        }
        delete lstRestakerRegistrations[msg.sender];
    }

    function freeze(bytes[] calldata valPubKeys) external
        onlyRegisteredValidators(valPubKeys) onlyFreezeOracle() whenNotPaused() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _freeze(valPubKeys[i]);
        }
    }

    // TODO: send fee to specified account
    // Write about how we restore the validator to REGISTERED since it has explicitly paid the protocol,
    // and signaled intention to act correctly again. This is unlike the slashing functionality from ValidatorRegistryV1,
    // where slashed validators are immediately unstaked, and must withdraw, then stake again to be once again "opted-in".
    function unfreeze(bytes calldata valPubKey) payable external onlyRegisteredValidator(valPubKey) whenNotPaused() {
        require(validatorRegistrations[valPubKey].freezeHeight.exists, "validator must be frozen");
        require(block.number >= validatorRegistrations[valPubKey].freezeHeight.blockHeight + unfreezePeriodBlocks,
            "unfreeze must be happen at least unfreezePeriodBlocks after freeze height");
        require(msg.value >= unfreezeFee, "sender must pay unfreeze fee");
        EventHeightLib.del(validatorRegistrations[valPubKey].freezeHeight);
        emit ValidatorUnfrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    function setAVSDirectory(IAVSDirectory avsDirectory_) external onlyOwner {
        _setAVSDirectory(avsDirectory_);
    }

    function setStrategyManager(IStrategyManager strategyManager_) external onlyOwner {
        _setStrategyManager(strategyManager_);
    }

    function setDelegationManager(IDelegationManager delegationManager_) external onlyOwner {
        _setDelegationManager(delegationManager_);
    }

    function setEigenPodManager(IEigenPodManager eigenPodManager_) external onlyOwner {
        _setEigenPodManager(eigenPodManager_);
    }

    function setRestakeableStrategies(address[] calldata restakeableStrategies_) external onlyOwner {
        _setRestakeableStrategies(restakeableStrategies_);
    }

    function setFreezeOracle(address freezeOracle_) external onlyOwner {
        _setFreezeOracle(freezeOracle_);
    }

    function setUnfreezeFee(uint256 unfreezeFee_) external onlyOwner {
        _setUnfreezeFee(unfreezeFee_);
    }

    function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) external onlyOwner {
        _setUnfreezePeriodBlocks(unfreezePeriodBlocks_);
    }

    function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) external onlyOwner {
        _setOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_);
    }

    function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) external onlyOwner {
        _setValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_);
    }

    function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) external onlyOwner {
        _setLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_);
    }

    function updateMetadataURI(string memory metadataURI_) external onlyOwner {
        _updateMetadataURI(metadataURI_);
    }

    function _registerValidatorsByPodOwner(
        bytes[] calldata valPubKeys,
        address podOwner
    ) internal onlyNonRegisteredValidators(valPubKeys) onlyPodOwnerOrOperator(podOwner)  {
        IEigenPod pod = _eigenPodManager.getPod(podOwner);
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE,
                "validator must be active under pod");
            _registerValidator(valPubKeys[i], podOwner);
        }
    }

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

    function _requestValidatorDeregistration(bytes calldata valPubKey) internal onlyPodOwnerOperatorOrContractOwner(valPubKey) {
        EventHeightLib.set(validatorRegistrations[valPubKey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function _deregisterValidator(bytes calldata valPubKey) internal onlyPodOwnerOperatorOrContractOwner(valPubKey) {
        require(validatorRegistrations[valPubKey].deregRequestHeight.exists,
            "validator must have requested deregistration");
        require(block.number >= validatorRegistrations[valPubKey].deregRequestHeight.blockHeight + validatorDeregPeriodBlocks,
            "deregistration must happen at least validatorDeregPeriodBlocks after deletion request height");
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        delete validatorRegistrations[valPubKey];
        emit ValidatorDeregistered(valPubKey, podOwner);
    }

    function _freeze(bytes calldata valPubKey) internal {
        EventHeightLib.set(validatorRegistrations[valPubKey].freezeHeight, block.number);
        EventHeightLib.del(validatorRegistrations[valPubKey].deregRequestHeight);
        emit ValidatorFrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function _setAVSDirectory(IAVSDirectory avsDirectory_) internal {
        _eigenAVSDirectory = avsDirectory_;
        emit AVSDirectorySet(address(_eigenAVSDirectory));
    }

    function _setStrategyManager(IStrategyManager strategyManager_) internal {
        _strategyManager = strategyManager_;
        emit StrategyManagerSet(address(strategyManager_));
    }

    function _setDelegationManager(IDelegationManager delegationManager_) internal {
        _delegationManager = delegationManager_;
        emit DelegationManagerSet(address(delegationManager_));
    }

    function _setEigenPodManager(IEigenPodManager eigenPodManager_) internal {
        _eigenPodManager = eigenPodManager_;
        emit EigenPodManagerSet(address(eigenPodManager_));
    }

    function _setRestakeableStrategies(address[] calldata restakeableStrategies_) internal {
        restakeableStrategies = restakeableStrategies_;
        emit RestakeableStrategiesSet(restakeableStrategies);
    }

    function _setFreezeOracle(address _freezeOracle) internal {
        freezeOracle = _freezeOracle;
        emit FreezeOracleSet(_freezeOracle);
    }

    function _setUnfreezeFee(uint256 _unfreezeFee) internal {
        unfreezeFee = _unfreezeFee;
        emit UnfreezeFeeSet(_unfreezeFee);
    }

    function _setUnfreezePeriodBlocks(uint256 _unfreezePeriodBlocks) internal {
        unfreezePeriodBlocks = _unfreezePeriodBlocks;
        emit UnfreezePeriodBlocksSet(_unfreezePeriodBlocks);
    }
    
    function _setOperatorDeregPeriodBlocks(uint256 _operatorDeregPeriodBlocks) internal {
        operatorDeregPeriodBlocks = _operatorDeregPeriodBlocks;
        emit OperatorDeregPeriodBlocksSet(_operatorDeregPeriodBlocks);
    }

    function _setValidatorDeregPeriodBlocks(uint256 _validatorDeregPeriodBlocks) internal {
        validatorDeregPeriodBlocks = _validatorDeregPeriodBlocks;
        emit ValidatorDeregPeriodBlocksSet(_validatorDeregPeriodBlocks);
    }

    function _setLstRestakerDeregPeriodBlocks(uint256 _lstRestakerDeregPeriodBlocks) internal {
        lstRestakerDeregPeriodBlocks = _lstRestakerDeregPeriodBlocks;
        emit LSTRestakerDeregPeriodBlocksSet(_lstRestakerDeregPeriodBlocks);
    }

    function _updateMetadataURI(string memory _metadataURI) internal {
        _eigenAVSDirectory.updateAVSMetadataURI(_metadataURI);
    }

    function getRestakeableStrategies() external view returns (address[] memory) {
        return _getRestakeableStrategies();
    }

    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
        if (!operatorRegistrations[operator].exists) {
            return new address[](0);
        }
        return _getRestakeableStrategies();
    }

    function areValidatorsOptedIn(bytes[] calldata valPubKeys) external view returns (bool[] memory) {
        bool[] memory result = new bool[](valPubKeys.length);
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            result[i] = _isValidatorOptedIn(valPubKeys[i]);
        }
        return result;
    }

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valPubKey);
    }

    function getLSTRestakerRegInfo(address lstRestaker) 
        external view returns (LSTRestakerRegistrationInfo memory) {
        return lstRestakerRegistrations[lstRestaker];
    }

    function avsDirectory() external view returns (address) {
        return address(_eigenAVSDirectory);
    }

    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        bool isValRegistered = validatorRegistrations[valPubKey].exists;
        IEigenPod pod = _eigenPodManager.getPod(validatorRegistrations[valPubKey].podOwner);
        bool isValActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address delegatedOperator = _delegationManager.delegatedTo(validatorRegistrations[valPubKey].podOwner);
        bool isOperatorRegistered = operatorRegistrations[delegatedOperator].exists;
        return isValRegistered && isValActive && isOperatorRegistered;
    }

    function _getRestakeableStrategies() internal view returns (address[] memory) {
        return restakeableStrategies;
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}
