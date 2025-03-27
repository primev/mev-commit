// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {TimestampOccurrence} from "../../utils/Occurrence.sol";

contract MevCommitBappStorage {

    struct ValidatorRecord {
        bool exists;
        address registrar;
        TimestampOccurrence.Occurrence freezeOccurrence;
        TimestampOccurrence.Occurrence deregRequestOccurrence;
    }

    uint256 public unfreezePeriod;
    uint256 public unfreezeFee;
    address public unfreezeReceiver;

    mapping(address => bool) public isWhitelisted;

    mapping(address => bool) public isOptedIn;

    mapping(bytes blsPubkey => ValidatorRecord) public validatorRecords;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
