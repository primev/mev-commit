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

// TODO: Allow owner account to opt out any hash in case keys are lost
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
        bytes calldata pubkey,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external {
        address operator = msg.sender;
        require(delegationManager.isOperator(operator), "sender must be an operator");
        require(!_isOperatorRegistered(operator), "operator must not already be registered with MevCommitAVS");
        eigenAVSDirectory.registerOperatorToAVS(operator, operatorSignature);
        emit OperatorAdded(operator);
    }

    function deregisterOperator(address operator) external {
        require(msg.sender == operator || msg.sender == owner(), "sender must be operator or owner");
        require(_isOperatorRegistered(operator), "operator must already be registered with MevCommitAVS");
        eigenAVSDirectory.deregisterOperatorFromAVS(operator);
        emit OperatorRemoved(operator);
    }


    // TODO: Incorporate the below enum into this contract. 
    // Note opt-in status will truly reside on mev-commit chain. However,
    // if one of the calling EOAs tries to opt-in validators that are not a part
    // of their eigenpod OR are not delegated to a registered operator, the off-chain
    // process post a tx that REJECTS the val set. 
    // (will also need to consider duplicate opt-in. Likely freeze in that scenario too)
    // (also got to decide if we freeze whitelisted account in duplicate scenario, if actual
    // (will likely have to enforce that a whitelisted address has to opt in entire
    // eigenpods, not certain vals. This way we could tell eigenpods with dup req,
    // "sorry, a whitelisted account has already opted in the vals from your pod)
    // eigenpod opted-in right before them.. Likely not. But this can be impl detail of oracle)
    // For both REJECTED AND FROZEN valsets, the user has to pay a fee to the oracle
    // to opt-in once again. Note can prob consolidate frozen and rejected.


    // TODO: Oracle reject function

    // TODO: Oracle Freeze function


    
    // TODO: Now we freeze by group of validators! Not vals themselves. This simplifies implementation. 
    // Note for v1 no freeze duration will need to exist. Any account can unfreeze itself immediately for fee. 

    // TODO: No requirement will exist that hash has to include all vals from pod. Adds to much complexity
    // Instead let a user know on frontend (before sending tx) that they cannot opt-in a key that's already opted in. 

    // TODO: Whitelist is now just operators! Every large org seems to have its own operator.
    // Note this can be what "operators do" for now. ie. they have the ability to opt-in their users. 
    // But we still allow home stakers to opt-in themselves too. 
    // Make it very clear that part 2 of opt-in is neccessary to explicitly communicate to 
    // the opter-inner that they must follow the relay connection requirement. Otherwise delegators may be 
    // blindly frozen. When opting in as a part of step 2, the sender should be running the validators
    // its opting in (st. relay requirement is met).

    // TODO: Determine if we have a function to let some accounts bypass pod requirement.
    // Update yes operators can do this.

    function storeValidatorsByPods(bytes[][] calldata valPubKeys, address[] calldata podAddrs) external {
        for (uint256 i = 0; i < podAddrs.length; i++) {
            _storeValidatorsByPod(valPubKeys[i], podAddrs[i]);
        }
    }

    function storeValidatorsByPod(bytes[] calldata valPubKeys, address podAddr) external {
        _storeValidatorsByPod(valPubKeys, podAddr);
    }

    function _storeValidatorsByPod(bytes[] calldata valPubKeys, address podAddr) internal {
        IEigenPod pod = IEigenPod(podAddr);
        address delegatedTo = delegationManager.delegatedTo(podAddr); // TODO: Confirm delegatedTo takes pod address
        require(msg.sender == pod.podOwner() || msg.sender == delegatedTo, "sender must be pod owner or operator that a pod is delegated to");
        for (uint256 i = 0; i < valPubKeys.length; i++) {
            require(pod.validatorPubkeyToInfo(valPubKeys[i]).status == IEigenPod.VALIDATOR_STATUS.ACTIVE, "validator must be active under pod");
            require(validatorRecords[valPubKeys[i]].status == ValidatorStatus.NULL, "validator record must not already exist");
            validatorRecords[valPubKeys[i]] = ValidatorRecord({
                status: ValidatorStatus.STORED,
                podAddress: podAddr
            });
        }
    }


    // TODO: Implement "request withdraw" type deal with block enforcement.
    // Use existing code and tests from other PR. Also likely will need to define a FSM struct. 



    // TODO: Pull whitelisting CRUD from other PR

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (ValidatorStatus) {
        // TODO: Check two things: validator is stored in this contract AND
        // Validator is active with eigenlayer (still associated to same pod).
        // Also check that validator still has delegation to registered operator. 
        // TODO: Anything else to check? Make sure no edge cases in eigenlayer state changes that happen async. 
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


