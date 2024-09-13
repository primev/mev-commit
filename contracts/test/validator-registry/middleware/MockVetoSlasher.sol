// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";

contract MockVetoSlasher is MockEntity {
    address private _resolver;

    constructor(uint64 type_, address resolver_) MockEntity(type_) {
        _resolver = resolver_;
    }

    function setResolver(address resolver_) external {
        _resolver = resolver_;
    }

    function resolver(bytes32, bytes memory) external view returns (address) {
        return _resolver;
    }
}
