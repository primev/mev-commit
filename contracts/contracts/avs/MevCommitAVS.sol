// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

// Path forward: Get DUNE chart working and likely go with eigenpod opt-in + whitelist
// (as long as there's evidence that eigenpods are quite dispersed / diverse)
// Otherwise just stick to eigen opt-in

// AddtoWhitelist(account)
// RemoveFromwhitelist instead of using "EOA"

// Still build out everything with "hash to L1" idea.. For now. 

/// Get started by just developing contracts, with associated doc at top-level which fully explains
// idea to external team. 

import {MevCommitAVSStorage} from "./MevCommitAVSStorage.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IEigenPod} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";

// TODO: documentation and overall gas optimization
// TODO: modifiers, pattern match prod contracts, order of funcs, interfaces etc.
// TODO: use tests from other PR? 
contract MevCommitAVS is MevCommitAVSStorage, OwnableUpgradeable, UUPSUpgradeable {

    IDelegationManager internal delegationManager;
    IEigenPodManager internal eigenPodManager;
    IAVSDirectory internal eigenAVSDirectory;

    address public freezeOracle;
    uint256 public unfreezeFee;
    uint256 public minUnfreezeBlocks; // Optional, as we can allow frozen validators to pay fee immediately
    uint256 public operatorDeregistrationPeriodBlocks;
    uint256 public validatorDeregistrationPeriodBlocks;

    // TODO: address if val pub keys need to be indexed
    event OperatorRegistered(address indexed operator);
    event OperatorDeregistrationRequested(address indexed operator);
    event OperatorDeregistered(address indexed operator);
    event ValidatorRegistered(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorFrozen(bytes indexed validatorPubKey, address indexed podOwner);
    event ValidatorUnfrozen(bytes indexed validatorPubKey, address indexed podOwner);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    function initialize(
        address _owner,
        address _freezeOracle,
        uint256 _unfreezeFee,
        uint256 _minUnfreezeBlocks,
        uint256 _operatorDeregistrationPeriodBlocks,
        uint256 _validatorDeregistrationPeriodBlocks,
        IDelegationManager _delegationManager,
        IEigenPodManager _eigenPodManager,
        IAVSDirectory _avsDirectory
    ) external initializer {
        delegationManager = _delegationManager;
        eigenPodManager = _eigenPodManager;
        _avsDirectory = _avsDirectory;
        freezeOracle = _freezeOracle;
        unfreezeFee = _unfreezeFee;
        minUnfreezeBlocks = _minUnfreezeBlocks;
        operatorDeregistrationPeriodBlocks = _operatorDeregistrationPeriodBlocks;
        validatorDeregistrationPeriodBlocks = _validatorDeregistrationPeriodBlocks;
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
        // If validator was requested for deregistration, they must become unfrozen and request again 
        validatorRegistrations[valPubKey].deregistrationRequestHeight = 0;
        emit ValidatorFrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
    }

    function unfreeze(bytes calldata valPubKey) payable external {
        require(validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.FROZEN,
            "validator must be frozen");
        require(block.number >= validatorRegistrations[valPubKey].freezeHeight + minUnfreezeBlocks,
            "unfreeze must be happen at least minUnfreezeBlocks after freeze height");
        require(msg.value >= unfreezeFee, "sender must pay unfreeze fee");
        validatorRegistrations[valPubKey].status = VALIDATOR_REGISTRATION_STATUS.REGISTERED;
        validatorRegistrations[valPubKey].freezeHeight = 0;
        emit ValidatorUnfrozen(valPubKey, validatorRegistrations[valPubKey].podOwner);
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
        bool isRegistered = validatorRegistrations[valPubKey].status == VALIDATOR_REGISTRATION_STATUS.REGISTERED;
        IEigenPod pod = eigenPodManager.getPod(validatorRegistrations[valPubKey].podOwner);
        bool isActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address delegatedOperator = delegationManager.delegatedTo(validatorRegistrations[valPubKey].podOwner);
        bool isOperatorRegistered = operatorRegistrations[delegatedOperator].status == OPERATOR_REGISTRATION_STATUS.REGISTERED;
        return isRegistered && isActive && isOperatorRegistered;
    }

    function avsDirectory() external view returns (address) {
        return address(eigenAVSDirectory);
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}
