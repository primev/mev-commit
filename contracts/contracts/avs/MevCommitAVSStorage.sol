// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;


// TODO: see if you can simplify these enums/structs (do you need enum for ex)
abstract contract MevCommitAVSStorage {

    struct OperatorRegistrationInfo {
        OPERATOR_REGISTRATION_STATUS status;
        uint256 deregistrationRequestHeight;
    }

    // TODO: Put in interface
    enum OPERATOR_REGISTRATION_STATUS {
        NOT_REGISTERED,
        REGISTERED,
        REQ_DEREGISTRATION
    }

    mapping(address => OperatorRegistrationInfo) public operatorRegistrations;

    // TODO: Put in interface
    struct ValidatorRegistrationInfo {
        VALIDATOR_REGISTRATION_STATUS status;
        address podOwner;
        uint256 freezeHeight;
        uint256 deregistrationRequestHeight;
    }

    // TODO: Put in interface
    enum VALIDATOR_REGISTRATION_STATUS{
        NOT_REGISTERED,
        REGISTERED,
        REQ_DEREGISTRATION,
        FROZEN
    }
    mapping(bytes => ValidatorRegistrationInfo) public validatorRegistrations;
}
