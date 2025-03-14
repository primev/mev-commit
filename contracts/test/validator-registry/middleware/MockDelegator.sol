// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {MockEntity} from "./MockEntity.sol";

contract MockDelegator is MockEntity {

    mapping(address => uint256) private _stakes;

    error SubnetworkCannotBeEmpty();

    constructor(uint64 type_) MockEntity(type_) { }

    function setStake(address operator, uint256 stake_) external {
        _stakes[operator] = stake_;
    }

    function stake(bytes32 subnetwork, address operator) external view returns (uint256) {
        if (subnetwork == bytes32("")) {
            revert SubnetworkCannotBeEmpty();
        }
        return _stakes[operator];
    }
}
