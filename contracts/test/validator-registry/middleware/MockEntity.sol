// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

contract MockEntity {
    uint64 private _type;

    constructor(uint64 type_) {
        _type = type_;
    }

    function setType(uint64 type_) external {
        _type = type_;
    }

    function TYPE() external view returns (uint64) {
        return _type;
    }
}
