// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {MevCommitAVSStorage} from "./MevCommitAVSStorage.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
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
// TODO: Do all this strategy param stuff, and decide of LST delegation is v1 or next version. See chooseValidator in doc
// TODO: Note and document everything from https://docs.eigenlayer.xyz/eigenlayer/avs-guides/avs-dashboard-onboarding
contract MevCommitAVS is IMevCommitAVS, MevCommitAVSStorage, OwnableUpgradeable, UUPSUpgradeable {

    IDelegationManager internal delegationManager;
    IEigenPodManager internal eigenPodManager;
    IAVSDirectory internal eigenAVSDirectory;

    address public freezeOracle;
    uint256 public unfreezeFee;
    uint256 public unfreezePeriodBlocks; // Optional, as we can allow frozen validators to pay fee immediately
    uint256 public operatorDeregistrationPeriodBlocks;
    uint256 public validatorDeregistrationPeriodBlocks;

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    function initialize(
        address _owner,
        IDelegationManager _delegationManager,
        IEigenPodManager _eigenPodManager,
        IAVSDirectory _avsDirectory,
        address _freezeOracle,
        uint256 _unfreezeFee,
        uint256 _unfreezePeriodBlocks,
        uint256 _operatorDeregistrationPeriodBlocks,
        uint256 _validatorDeregistrationPeriodBlocks,
        string calldata metadataURI_
    ) external initializer {
        _setDelegationManager(_delegationManager);
        _setEigenPodManager(_eigenPodManager);
        _setAVSDirectory(_avsDirectory);
        _setFreezeOracle(_freezeOracle);
        _setUnfreezeFee(_unfreezeFee);
        _setUnfreezePeriodBlocks(_unfreezePeriodBlocks);
        _setOperatorDeregistrationPeriodBlocks(_operatorDeregistrationPeriodBlocks);
        _setValidatorDeregistrationPeriodBlocks(_validatorDeregistrationPeriodBlocks);

        if (bytes(metadataURI_).length > 0) {
            _avsDirectory.updateAVSMetadataURI(metadataURI_);
        }

        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, "sender must be freeze oracle");
        _;
    }

    modifier onlyEigenlayerRegisteredOperator() {
        require(delegationManager.isOperator(msg.sender), "sender must be an eigenlayer operator");
        _;
    }
    
    modifier onlyOperatorDeregistrar(address operator) {
        require(msg.sender == operator || msg.sender == owner(), "sender must be operator or MevCommitAVS owner");
        _;
    }

    modifier onlyValidatorRegistrarWithOperatorRegistered(address podOwner) {
        address delegatedOperator = delegationManager.delegatedTo(podOwner);
        require(msg.sender == podOwner || msg.sender == delegatedOperator, 
            "sender must be podOwner or delegated operator");
        require(operatorRegistrations[delegatedOperator].status == OPERATOR_REGISTRATION_STATUS.REGISTERED,
            "delegated operator must be registered with MevCommitAVS");
        _;
    }

    modifier onlyValidatorDeregistrar(bytes calldata valPubKey) {
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        require(msg.sender == podOwner ||
            msg.sender == delegationManager.delegatedTo(podOwner) ||
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
    ) external onlyEigenlayerRegisteredOperator() {
        require(operatorRegistrations[msg.sender].status == OPERATOR_REGISTRATION_STATUS.NOT_REGISTERED,
            "operator must not already be registered with MevCommitAVS");
        eigenAVSDirectory.registerOperatorToAVS(msg.sender, operatorSignature);
        operatorRegistrations[msg.sender] = OperatorRegistrationInfo({
            status: OPERATOR_REGISTRATION_STATUS.REGISTERED,
            deregistrationRequestHeight: 0
        });
        emit OperatorRegistered(msg.sender);
    }

    function requestOperatorDeregistration(address operator) external onlyOperatorDeregistrar(operator) {
        require(operatorRegistrations[operator].status == OPERATOR_REGISTRATION_STATUS.REGISTERED,
            "operator must be registered with MevCommitAVS");
        operatorRegistrations[operator].status = OPERATOR_REGISTRATION_STATUS.REQ_DEREGISTRATION;
        operatorRegistrations[operator].deregistrationRequestHeight = block.number;
        emit OperatorDeregistrationRequested(operator);
    }

    function deregisterOperator(address operator) external onlyOperatorDeregistrar(operator) {
        require(operatorRegistrations[operator].status == OPERATOR_REGISTRATION_STATUS.REQ_DEREGISTRATION,
            "operator must have requested deregistration");
        require(block.number >= operatorRegistrations[operator].deregistrationRequestHeight + operatorDeregistrationPeriodBlocks,
            "deregistration must happen at least operatorDeregistrationPeriodBlocks after deregistration request height");
        eigenAVSDirectory.deregisterOperatorFromAVS(operator);
        delete operatorRegistrations[operator];
        emit OperatorDeregistered(operator);
    }

    function registerValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external {
        for (uint256 i = 0; i < podOwners.length; i++) {
            _registerValidatorsByPodOwner(valPubKeys[i], podOwners[i]);
        }
    }

    function _registerValidatorsByPodOwner(
        bytes[] calldata valPubKeys,
        address podOwner
    ) internal onlyValidatorRegistrarWithOperatorRegistered(podOwner) {
        IEigenPod pod = eigenPodManager.getPod(podOwner);
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE,
                "validator must be active under pod");
            require(validatorRegistrations[valPubKeys[i]].status == VALIDATOR_REGISTRATION_STATUS.NOT_REGISTERED,
                "validator must not already be registered");
            _registerValidator(valPubKeys[i], podOwner);
        }
    }

    function _registerValidator(bytes calldata valPubKey, address podOwner) internal {
        validatorRegistrations[valPubKey] = ValidatorRegistrationInfo({
            status: VALIDATOR_REGISTRATION_STATUS.REGISTERED,
            podOwner: podOwner,
            freezeHeight: 0,
            deregistrationRequestHeight: 0
        });
        emit ValidatorRegistered(valPubKey, podOwner);
    }

    function requestValidatorsDeregistration(bytes[] calldata valPubKeys) external {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _requestValidatorDeregistration(valPubKeys[i]);
        }
    }
    
    function _requestValidatorDeregistration(bytes calldata valPubKey) internal onlyValidatorDeregistrar(valPubKey) {
        require(validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REGISTERED,
            "validator must be currently registered");
        validatorRegistrations[valPubKey].status = VALIDATOR_REGISTRATION_STATUS.REQ_DEREGISTRATION;
        validatorRegistrations[valPubKey].deregistrationRequestHeight = block.number;
        emit ValidatorDeregistrationRequested(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function deregisterValidators(bytes[] calldata valPubKeys) external {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            _deregisterValidator(valPubKeys[i]);
        }
    }

    function _deregisterValidator(bytes calldata valPubKey) internal onlyValidatorDeregistrar(valPubKey) {
        require(validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REQ_DEREGISTRATION,
            "validator must have requested deregistration");
        require(block.number >= validatorRegistrations[valPubKey].deregistrationRequestHeight + validatorDeregistrationPeriodBlocks,
            "deletion must happen at least validatorDeregistrationPeriodBlocks after deletion request height");
        address podOwner = validatorRegistrations[valPubKey].podOwner;
        delete validatorRegistrations[valPubKey];
        emit ValidatorDeregistered(valPubKey, podOwner);
    }

    function freeze(bytes calldata valPubKey) onlyFreezeOracle external {
        require(validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REGISTERED ||
            validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REQ_DEREGISTRATION,
            "validator must be registered or requested for deregistration");
        validatorRegistrations[valPubKey].status = VALIDATOR_REGISTRATION_STATUS.FROZEN;
        validatorRegistrations[valPubKey].freezeHeight = block.number;
        validatorRegistrations[valPubKey].deregistrationRequestHeight = 0;
        emit ValidatorFrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function unfreeze(bytes calldata valPubKey) payable external {
        require(validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.FROZEN,
            "validator must be frozen");
        require(block.number >= validatorRegistrations[valPubKey].freezeHeight + unfreezePeriodBlocks,
            "unfreeze must be happen at least unfreezePeriodBlocks after freeze height");
        require(msg.value >= unfreezeFee, "sender must pay unfreeze fee");
        validatorRegistrations[valPubKey].status = VALIDATOR_REGISTRATION_STATUS.REGISTERED;
        validatorRegistrations[valPubKey].freezeHeight = 0;
        emit ValidatorUnfrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    // function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {
    //     if (operatorRegistrations[operator].status != OPERATOR_REGISTRATION_STATUS.REGISTERED) {
    //         return new address[](0);
    //     }
    //     return _getRestakeableStrategies();
    // }

    // function getRestakeableStrategies() external view returns (address[] memory) {
    //     return _getRestakeableStrategies();
    // }

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
        bool isRegistered = validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REGISTERED;
        IEigenPod pod = eigenPodManager.getPod(validatorRegistrations[valPubKey].podOwner);
        bool isActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address delegatedOperator = delegationManager.delegatedTo(validatorRegistrations[valPubKey].podOwner);
        bool isOperatorRegistered = operatorRegistrations[delegatedOperator].status == OPERATOR_REGISTRATION_STATUS.REGISTERED;
        return isRegistered && isActive && isOperatorRegistered;
    }

    /// @notice Returns eigenlayer AVS directory contract address to abide by IServiceManager interface.
    function avsDirectory() external view returns (address) {
        return address(eigenAVSDirectory);
    }

    function setMetadataURI(string memory metadataURI) external onlyOwner {
        eigenAVSDirectory.updateAVSMetadataURI(metadataURI);
    }

    function setAVSDirectory(IAVSDirectory _avsDirectory) external onlyOwner {
        _setAVSDirectory(_avsDirectory);
    }

    function setDelegationManager(IDelegationManager _delegationManager) external onlyOwner {
        _setDelegationManager(_delegationManager);
    }

    function setEigenPodManager(IEigenPodManager _eigenPodManager) external onlyOwner {
        _setEigenPodManager(_eigenPodManager);
    }

    function setFreezeOracle(address _freezeOracle) external onlyOwner {
        _setFreezeOracle(_freezeOracle);
    }

    function setUnfreezeFee(uint256 _unfreezeFee) external onlyOwner {
        _setUnfreezeFee(_unfreezeFee);
    }

    function setUnfreezePeriodBlocks(uint256 _unfreezePeriodBlocks) external onlyOwner {
        _setUnfreezePeriodBlocks(_unfreezePeriodBlocks);
    }

    function setOperatorDeregistrationPeriodBlocks(uint256 _operatorDeregistrationPeriodBlocks) external onlyOwner {
        _setOperatorDeregistrationPeriodBlocks(_operatorDeregistrationPeriodBlocks);
    }

    function setValidatorDeregistrationPeriodBlocks(uint256 _validatorDeregistrationPeriodBlocks) external onlyOwner {
        _setValidatorDeregistrationPeriodBlocks(_validatorDeregistrationPeriodBlocks);
    }

    function _setAVSDirectory(IAVSDirectory _avsDirectory) private {
        eigenAVSDirectory = _avsDirectory;
        emit AVSDirectorySet(address(_avsDirectory));
    }

    function _setDelegationManager(IDelegationManager _delegationManager) private {
        delegationManager = _delegationManager;
        emit DelegationManagerSet(address(_delegationManager));
    }

    function _setEigenPodManager(IEigenPodManager _eigenPodManager) private {
        eigenPodManager = _eigenPodManager;
        emit EigenPodManagerSet(address(_eigenPodManager));
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
