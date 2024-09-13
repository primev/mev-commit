// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";

contract MockVetoSlasher is MockEntity {
    address private _resolver;
    uint256 private _vetoDuration;

    constructor(uint64 type_, address resolver_, uint256 vetoDuration_) MockEntity(type_) {
        _resolver = resolver_;
        _vetoDuration = vetoDuration_;
    }

    function setResolver(address resolver_) external {
        _resolver = resolver_;
    }

    function setVetoDuration(uint256 vetoDuration_) external {
        _vetoDuration = vetoDuration_;
    }

    function resolver(bytes32, bytes memory) external view returns (address) {
        return _resolver;
    }

    function vetoDuration() external view returns (uint256) {
        return _vetoDuration;
    }
}
