// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";

contract MockInstantSlasher is MockEntity {
    constructor(uint64 type_) MockEntity(type_) {}
}
