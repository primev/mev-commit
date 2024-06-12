// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

abstract contract MevCommitAVSStorage {
    mapping(address => bool) public operators;

    // TODO: Put in interface
    struct ValidatorRecord {
        VALIDATOR_RECORD_STATUS status;
        address podOwner;
    }

    // TODO: Put in interface
    enum VALIDATOR_RECORD_STATUS{
        NULL,
        STORED,
        FROZEN,
        WITHDRAWING
    }
    mapping(bytes => ValidatorRecord) public validatorRecords;
}
