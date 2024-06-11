// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

abstract contract MevCommitAVSStorage {
    mapping(address => bool) public operators;

    // TODO: Put in interface
    struct ValidatorRecord {
        ValidatorStatus status;
        address podAddress;
    }

    // TODO: Put in interface
    enum ValidatorStatus{
        NULL,
        STORED,
        FROZEN,
        WITHDRAWING
    }
    mapping(bytes => ValidatorRecord) public validatorRecords;
}
