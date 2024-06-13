// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;


import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";

abstract contract MevCommitAVSStorage {
    /// @notice Mapping of operator addresses to their registration info
    mapping(address => IMevCommitAVS.OperatorRegistrationInfo) public operatorRegistrations;

    /// @notice Mapping of validator pubkeys to their registration info
    mapping(bytes => IMevCommitAVS.ValidatorRegistrationInfo) public validatorRegistrations;
}
