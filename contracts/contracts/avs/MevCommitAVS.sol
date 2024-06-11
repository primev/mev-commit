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

// Write about how v2 (or future version with more decentralization) will 
// give operators the task of doing the pubkey relaying to the mev-commit chain. 
// That is the off-chain process is replaced by operators, who all look for the 
// valset lists posted to some DA layer (eigenDA?), and then race/attest to post
// this to the mev-commit chain. The operator accounts could be auto funded on our chain. 
contract MevCommitAVS is MevCommitAVSStorage, OwnableUpgradeable, UUPSUpgradeable {

    IDelegationManager internal delegationManager;
    IEigenPodManager internal eigenPodManager;

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner { }

    function initialize(
        address _owner,
        IDelegationManager _delegationManager,
        IEigenPodManager _eigenPodManager
    ) external initializer {
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
        delegationManager = _delegationManager;
        eigenPodManager = _eigenPodManager;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // TODO: Incorporate the below enum into this contract. 
    // Note opt-in status will truly reside on mev-commit chain. However,
    // if one of the calling EOAs tries to opt-in validators that are not a part
    // of their eigenpod OR are not delegated to a registered operator, the off-chain
    // process post a tx that REJECTS the val set. 
    // (will also need to consider duplicate opt-in. Likely freeze in that scenario too)
    // (also got to decide if we freeze whitelisted account in duplicate scenario, if actual
    // eigenpod opted-in right before them.. Likely not. But this can be impl detail of oracle)
    // For both REJECTED AND FROZEN valsets, the user has to pay a fee to the oracle
    // to opt-in once again. Note can prob consolidate frozen and rejected.

    // TODO: also note that it once again may make sense to freeze by group of validators!
    // Make sure murat knows of this. 
    // Note freezing MUST exist on L1, so that users could pay back the oracle on that same chain. 
    // In this case. Mev-commit chain is just a place for cheap blockspace that we can post the full
    // validator list, akin to DA. 
    // TODO: Does this enum go in storage file? 
    enum ValSetStatus{
        EMPTY,
        ATTESTED,
        WITHDRAWING,
        REJECTED,
        FROZEN
    }

    // TODO: Oracle reject function

    // TODO: Oracle Freeze function

    function optInValSet(bytes32 valSetHash) external {
        require(valSetHash != bytes32(0), "valSetHash must be non-zero");
        require(eigenPodManager.hasPod(msg.sender), "Caller must have eigenpod");
        require(addressToValSetHash[msg.sender] == bytes32(0), "Caller must not already have opted in");
        addressToValSetHash[msg.sender] = valSetHash;
    }

    // TODO: Determine if we have a function to let some accounts bypass pod requirement.
    function optInValSetFromWhitelist(bytes32 valSetHash) external {
        require(valSetHash != bytes32(0), "valSetHash must be non-zero");
        require(whitelist[msg.sender], "Caller must be whitelisted");
        require(addressToValSetHash[msg.sender] == bytes32(0), "Caller must not already have opted in");
        addressToValSetHash[msg.sender] = valSetHash;
    }

    // TODO: Implement "request withdraw" type deal with block enforcement.
    // Use existing code and tests from other PR. Also likely will need to define a FSM struct. 

    function optOutValSet(bytes32 valSetHash) external {
        require(valSetHash != bytes32(0), "valSetHash must be non-zero");
        require(addressToValSetHash[msg.sender] == valSetHash, "Caller must have opted in with same hash");
        addressToValSetHash[msg.sender] = bytes32(0);
    }

    function optOutValSetFromWhitelist(bytes32 valSetHash) external {
        require(valSetHash != bytes32(0), "valSetHash must be non-zero");
        require(whitelist[msg.sender], "Caller must be whitelisted");
        require(addressToValSetHash[msg.sender] == valSetHash, "Caller must have opted in with same hash");
        addressToValSetHash[msg.sender] = bytes32(0);
    }

    // TODO: Pull whitelisting CRUD from other PR

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}


