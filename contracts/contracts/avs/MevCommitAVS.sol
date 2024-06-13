// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {MevCommitAVSStorage} from "./MevCommitAVSStorage.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IEigenPod} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";

// TODO: overall gas optimization
// TODO: order of funcs, finish interfaces, comments for everything etc.
// TODO: use tests from other PR? 
// TODO: test upgradability before Holesky deploy
// TODO: Decide of LST delegation is v1 or next version. See chooseValidator in doc
// TODO: Note and document everything from https://docs.eigenlayer.xyz/eigenlayer/avs-guides/avs-dashboard-onboarding
// TODO: Confirm all setters are present and in right order, confirm interface is fully populated
contract MevCommitAVS is IMevCommitAVS, MevCommitAVSStorage, OwnableUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    IDelegationManager internal _delegationManager;
    IEigenPodManager internal _eigenPodManager;
    IAVSDirectory internal _eigenAVSDirectory;

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }
    function pause() external onlyOwner { _pause(); }
    function unpause() external onlyOwner { _unpause(); }

    function initialize(
        address owner_,
        IDelegationManager delegationManager_,
        IEigenPodManager eigenPodManager_,
        IAVSDirectory avsDirectory_,
        address[] calldata restakeableStrategies_,
        address freezeOracle_,
        uint256 unfreezeFee_,
        uint256 unfreezePeriodBlocks_,
        uint256 operatorDeregistrationPeriodBlocks_,
        uint256 validatorDeregistrationPeriodBlocks_,
        string calldata metadataURI_
    ) external initializer {
        _setDelegationManager(delegationManager_);
        _setEigenPodManager(eigenPodManager_);
        _setAVSDirectory(avsDirectory_);
        _setRestakeableStrategies(restakeableStrategies_);
        _setFreezeOracle(freezeOracle_);
        _setUnfreezeFee(unfreezeFee_);
        _setUnfreezePeriodBlocks(unfreezePeriodBlocks_);
        _setOperatorDeregistrationPeriodBlocks(operatorDeregistrationPeriodBlocks_);
        _setValidatorDeregistrationPeriodBlocks(validatorDeregistrationPeriodBlocks_);

        if (bytes(metadataURI_).length > 0) {
            _eigenAVSDirectory.updateAVSMetadataURI(metadataURI_);
        }

        __Ownable_init(owner_);
        __UUPSUpgradeable_init();
        __Pausable_init();
    }

    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, "sender must be freeze oracle");
        _;
    }

    modifier onlyEigenlayerRegisteredOperator() {
        require(_delegationManager.isOperator(msg.sender), "sender must be an eigenlayer operator");
        _;
    }
    
    modifier onlyOperatorDeregistrar(address operator) {
        require(msg.sender == operator || msg.sender == owner(), "sender must be operator or MevCommitAVS owner");
        _;
    }

    modifier onlyValidatorRegistrarWithOperatorRegistered(address podOwner) {
        address delegatedOperator = _delegationManager.delegatedTo(podOwner);
        require(msg.sender == podOwner || msg.sender == delegatedOperator, 
            "sender must be podOwner or delegated operator");
        require(operatorRegistrations[delegatedOperator].status == OperatorRegistrationStatus.REGISTERED,
            "delegated operator must be registered with MevCommitAVS");
        _;
    }

    modifier onlyValidatorDeregistrar(bytes calldata valPubKey) {
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        require(msg.sender == podOwner ||
            msg.sender == _delegationManager.delegatedTo(podOwner) ||
            msg.sender == owner(),
            "sender must be podOwner, delegated operator, or MevCommitAVS owner");
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function registerOperator (
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external onlyEigenlayerRegisteredOperator() whenNotPaused() {
        require(operatorRegistrations[msg.sender].status == OperatorRegistrationStatus.NOT_REGISTERED,
            "operator must not already be registered with MevCommitAVS");
        _eigenAVSDirectory.registerOperatorToAVS(msg.sender, operatorSignature);
        operatorRegistrations[msg.sender] = OperatorRegistrationInfo({
            status: OperatorRegistrationStatus.REGISTERED,
            deregistrationRequestHeight: 0
        });
        emit OperatorRegistered(msg.sender);
    }

    function requestOperatorDeregistration(address operator
    ) external onlyOperatorDeregistrar(operator) whenNotPaused() {
        require(operatorRegistrations[operator].status == OperatorRegistrationStatus.REGISTERED,
            "operator must be registered with MevCommitAVS");
        operatorRegistrations[operator].status = OperatorRegistrationStatus.REQ_DEREGISTRATION;
        operatorRegistrations[operator].deregistrationRequestHeight = block.number;
        emit OperatorDeregistrationRequested(operator);
    }

    function deregisterOperator(address operator
    ) external onlyOperatorDeregistrar(operator) whenNotPaused() {
        require(operatorRegistrations[operator].status == OperatorRegistrationStatus.REQ_DEREGISTRATION,
            "operator must have requested deregistration");
        require(block.number >= operatorRegistrations[operator].deregistrationRequestHeight + operatorDeregistrationPeriodBlocks,
            "deregistration must happen at least operatorDeregistrationPeriodBlocks after deregistration request height");
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

    function _registerValidatorsByPodOwner(
        bytes[] calldata valPubKeys,
        address podOwner
    ) internal onlyValidatorRegistrarWithOperatorRegistered(podOwner) {
        IEigenPod pod = _eigenPodManager.getPod(podOwner);
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE,
                "validator must be active under pod");
            require(validatorRegistrations[valPubKeys[i]].status == ValidatorRegistrationStatus.NOT_REGISTERED,
                "validator must not already be registered");
            _registerValidator(valPubKeys[i], podOwner);
        }
    }

    function _registerValidator(bytes calldata valPubKey, address podOwner) internal {
        validatorRegistrations[valPubKey] = ValidatorRegistrationInfo({
            status: ValidatorRegistrationStatus.REGISTERED,
            podOwner: podOwner,
            freezeHeight: 0,
            deregistrationRequestHeight: 0
        });
        emit ValidatorRegistered(valPubKey, podOwner);
    }

    function requestValidatorsDeregistration(
        bytes[] calldata valPubKeys
    ) external whenNotPaused() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _requestValidatorDeregistration(valPubKeys[i]);
        }
    }
    
    function _requestValidatorDeregistration(bytes calldata valPubKey) internal onlyValidatorDeregistrar(valPubKey) {
        require(validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.REGISTERED,
            "validator must be currently registered");
        validatorRegistrations[valPubKey].status = ValidatorRegistrationStatus.REQ_DEREGISTRATION;
        validatorRegistrations[valPubKey].deregistrationRequestHeight = block.number;
        emit ValidatorDeregistrationRequested(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function deregisterValidators(
        bytes[] calldata valPubKeys
    ) external whenNotPaused() {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _deregisterValidator(valPubKeys[i]);
        }
    }

    function _deregisterValidator(bytes calldata valPubKey) internal onlyValidatorDeregistrar(valPubKey) {
        require(validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.REQ_DEREGISTRATION,
            "validator must have requested deregistration");
        require(block.number >= validatorRegistrations[valPubKey].deregistrationRequestHeight + validatorDeregistrationPeriodBlocks,
            "deletion must happen at least validatorDeregistrationPeriodBlocks after deletion request height");
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        delete validatorRegistrations[valPubKey];
        emit ValidatorDeregistered(valPubKey, podOwner);
    }

    // TODO: Implement LST restakers having to choose validator. 
    // There will be additional storage field for this. 
    // Make sure to check supported strategies for this. 

    // function chooseValidator() // includes choosing new val
    // function removeValidatorChoice()

    function freeze(
        bytes calldata valPubKey
    ) external onlyFreezeOracle() whenNotPaused() {
        require(validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.REGISTERED ||
            validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.REQ_DEREGISTRATION,
            "validator must be registered or requested for deregistration");
        validatorRegistrations[valPubKey].status = ValidatorRegistrationStatus.FROZEN;
        validatorRegistrations[valPubKey].freezeHeight = block.number;
        validatorRegistrations[valPubKey].deregistrationRequestHeight = 0;
        emit ValidatorFrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function unfreeze(bytes calldata valPubKey
    ) payable external whenNotPaused() {
        require(validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.FROZEN,
            "validator must be frozen");
        require(block.number >= validatorRegistrations[valPubKey].freezeHeight + unfreezePeriodBlocks,
            "unfreeze must be happen at least unfreezePeriodBlocks after freeze height");
        require(msg.value >= unfreezeFee, "sender must pay unfreeze fee");
        validatorRegistrations[valPubKey].status = ValidatorRegistrationStatus.REGISTERED;
        validatorRegistrations[valPubKey].freezeHeight = 0;
        emit ValidatorUnfrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
        if (operatorRegistrations[operator].status != OperatorRegistrationStatus.REGISTERED) {
            return new address[](0);
        }
        return _getRestakeableStrategies();
    }

    function getRestakeableStrategies() external view returns (address[] memory) {
        return _getRestakeableStrategies();
    }

    function _getRestakeableStrategies() internal view returns (address[] memory) {
        return restakeableStrategies;
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

    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        bool isRegistered = validatorRegistrations[valPubKey].status == ValidatorRegistrationStatus.REGISTERED;
        IEigenPod pod = _eigenPodManager.getPod(validatorRegistrations[valPubKey].podOwner);
        bool isActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address delegatedOperator = _delegationManager.delegatedTo(validatorRegistrations[valPubKey].podOwner);
        bool isOperatorRegistered = operatorRegistrations[delegatedOperator].status == OperatorRegistrationStatus.REGISTERED;
        return isRegistered && isActive && isOperatorRegistered;
    }

    /// @notice Returns eigenlayer AVS directory contract address to abide by IServiceManager interface.
    function avsDirectory() external view returns (address) {
        return address(_eigenAVSDirectory);
    }

    function updateMetadataURI(string memory metadataURI) external onlyOwner {
        _eigenAVSDirectory.updateAVSMetadataURI(metadataURI);
    }

    function setAVSDirectory(IAVSDirectory avsDirectory_) external onlyOwner {
        _setAVSDirectory(avsDirectory_);
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

    function setOperatorDeregistrationPeriodBlocks(uint256 operatorDeregistrationPeriodBlocks_) external onlyOwner {
        _setOperatorDeregistrationPeriodBlocks(operatorDeregistrationPeriodBlocks_);
    }

    function setValidatorDeregistrationPeriodBlocks(uint256 validatorDeregistrationPeriodBlocks_) external onlyOwner {
        _setValidatorDeregistrationPeriodBlocks(validatorDeregistrationPeriodBlocks_);
    }

    function _setAVSDirectory(IAVSDirectory avsDirectory_) private {
        _eigenAVSDirectory = avsDirectory_;
        emit AVSDirectorySet(address(_eigenAVSDirectory));
    }

    function _setDelegationManager(IDelegationManager delegationManager_) private {
        _delegationManager = delegationManager_;
        emit DelegationManagerSet(address(delegationManager_));
    }

    function _setEigenPodManager(IEigenPodManager eigenPodManager_) private {
        _eigenPodManager = eigenPodManager_;
        emit EigenPodManagerSet(address(eigenPodManager_));
    }

    function _setRestakeableStrategies(address[] calldata restakeableStrategies_) private {
        restakeableStrategies = restakeableStrategies_;
        emit RestakeableStrategiesSet(restakeableStrategies);
    }

    function _setFreezeOracle(address _freezeOracle) private {
        freezeOracle = _freezeOracle;
        emit FreezeOracleSet(_freezeOracle);
    }

    function _setUnfreezeFee(uint256 _unfreezeFee) private {
        unfreezeFee = _unfreezeFee;
        emit UnfreezeFeeSet(_unfreezeFee);
    }

    function _setUnfreezePeriodBlocks(uint256 _unfreezePeriodBlocks) private {
        unfreezePeriodBlocks = _unfreezePeriodBlocks;
        emit UnfreezePeriodBlocksSet(_unfreezePeriodBlocks);
    }
    
    function _setOperatorDeregistrationPeriodBlocks(uint256 _operatorDeregistrationPeriodBlocks) private {
        operatorDeregistrationPeriodBlocks = _operatorDeregistrationPeriodBlocks;
        emit OperatorDeregistrationPeriodBlocksSet(_operatorDeregistrationPeriodBlocks);
    }

    function _setValidatorDeregistrationPeriodBlocks(uint256 _validatorDeregistrationPeriodBlocks) private {
        validatorDeregistrationPeriodBlocks = _validatorDeregistrationPeriodBlocks;
        emit ValidatorDeregistrationPeriodBlocksSet(_validatorDeregistrationPeriodBlocks);
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}
