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


// Write about how v2 (or future version with more decentralization) will 
// give operators the task of doing the pubkey relaying to the mev-commit chain. 
// That is the off-chain process is replaced by operators, who all look for the 
// valset lists posted to some DA layer (eigenDA?), and then race/attest to post
// this to the mev-commit chain. The operator accounts could be auto funded on our chain. 
// Slashing operators in this scheme would require social intervention as it could
// be pretty clear off chain of malicous actions and/or malicious off-chain validation
// of eigenpod conditions, delegation conditions, etc. 
contract MevCommitAVS is MevCommitAVSStorage, OwnableUpgradeable, UUPSUpgradeable {

    IDelegationManager internal delegationManager;
    IEigenPodManager internal eigenPodManager;
    IAVSDirectory internal eigenAVSDirectory;
    
    event OperatorAdded(address indexed operator);
    event OperatorRemoved(address indexed operator);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    function initialize(
        address _owner,
        IDelegationManager _delegationManager,
        IEigenPodManager _eigenPodManager,
        IAVSDirectory _avsDirectory
    ) external initializer {
        delegationManager = _delegationManager;
        eigenPodManager = _eigenPodManager;
        _avsDirectory = _avsDirectory;
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function registerOperator(
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external {
        address operator = msg.sender;
        require(delegationManager.isOperator(operator), "sender must be an operator");
        require(!_isOperatorRegistered(operator), "operator must not already be registered with MevCommitAVS");
        eigenAVSDirectory.registerOperatorToAVS(operator, operatorSignature);
        operators[operator] = true;
        emit OperatorAdded(operator);
    }

    function deregisterOperator(address operator) external {
        require(msg.sender == operator || msg.sender == owner(), "sender must be operator or owner");
        require(_isOperatorRegistered(operator), "operator must already be registered with MevCommitAVS");
        eigenAVSDirectory.deregisterOperatorFromAVS(operator);
        delete operators[operator];
        emit OperatorRemoved(operator);
    }

    // TODO: Oracle Freeze function
    // For FROZEN valsets, the user has to pay a fee to the oracle to opt-in once again. 

    // TODO: Whitelist is now just operators! Every large org seems to have its own operator.
    // Note this can be what "operators do" for now. ie. they have the ability to opt-in their users. 
    // But we still allow home stakers to opt-in themselves too. 
    // Make it very clear that part 2 of opt-in is neccessary to explicitly communicate to 
    // the opter-inner that they must follow the relay connection requirement. Otherwise delegators may be 
    // blindly frozen. When opting in as a part of step 2, the sender should be running the validators
    // its opting in (st. relay requirement is met).

    function storeValidatorsByPodOwners(bytes[][] calldata valPubKeys, address[] calldata podOwners) external {
        for (uint256 i = 0; i < podOwners.length; i++) {
            _storeValidatorsByPodOwner(valPubKeys[i], podOwners[i]);
        }
    }

    function _storeValidatorsByPodOwner(bytes[] calldata valPubKeys, address podOwner) internal {
        address operator = delegationManager.delegatedTo(podOwner);
        require(operator != address(0), "operator must be set for pod owner");
        require(_isOperatorRegistered(operator),
            "delegated operator must be registered with MevCommitAVS");
        require(msg.sender == podOwner || msg.sender == operator,
            "sender must be pod owner or delegated operator");
        IEigenPod pod = eigenPodManager.getPod(podOwner);

        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE,
                "validator must be active under pod");
            require(validatorRecords[valPubKeys[i]].status == VALIDATOR_RECORD_STATUS.NULL,
                "record must not already exist for relevant validator");
            validatorRecords[valPubKeys[i]] = ValidatorRecord({
                status: VALIDATOR_RECORD_STATUS.STORED,
                podOwner: podOwner
            });
        }
    }

    function deleteValidators(bytes[] calldata valPubKeys) external {
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(validatorRecords[valPubKeys[i]].status == VALIDATOR_RECORD_STATUS.STORED,
                "validator record must be stored");
            address podOwner = validatorRecords[valPubKeys[i]].podOwner;
            address operator = delegationManager.delegatedTo(podOwner);
            require(msg.sender == owner() || msg.sender == podOwner || msg.sender == operator,
                "sender must be MevCommitAVS owner, pod owner, or delegated operator");
            _deleteValidator(valPubKeys[i]);
        }
    }

    function _deleteValidator(bytes calldata valPubKey) internal {
        delete validatorRecords[valPubKey];
    }

    // TODO: Implement "request withdraw" type deal with block enforcement.
    // Use existing code and tests from other PR. Also likely will need to define a FSM struct. 

    // TODO: Pull some CRUD from other PR? 

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool) {
        bool isStored = validatorRecords[valPubKey].status == VALIDATOR_RECORD_STATUS.STORED;
        IEigenPod pod = eigenPodManager.getPod(validatorRecords[valPubKey].podOwner);
        bool isActive = pod.validatorPubkeyToInfo(valPubKey).status == IEigenPod.VALIDATOR_STATUS.ACTIVE;
        address operator = delegationManager.delegatedTo(validatorRecords[valPubKey].podOwner);
        bool isOperatorRegistered = _isOperatorRegistered(operator);
        return isStored && isActive && isOperatorRegistered;
    }

    function _isOperatorRegistered(address operator) internal view returns (bool) {
        return operators[operator];
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


