// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";

contract MockInstantSlasher is MockEntity {
    constructor(uint64 type_) MockEntity(type_) {}

    mapping(address operator => uint256 slashedAmount) public slashedAmounts;

    error InvalidSubnetwork();
    error InvalidOperator();
    error InvalidAmount();
    error InvalidInfractionTimestamp();
    error InvalidData();

    function slash(
        bytes32 subnetwork,
        address operator,
        uint256 amount,
        uint48 infractionTimestamp,
        bytes memory data
    ) external returns (uint256 slashedAmount) {
        require(subnetwork != bytes32(0), InvalidSubnetwork());
        require(operator != address(0), InvalidOperator());
        require(amount != 0, InvalidAmount());
        require(infractionTimestamp != 0, InvalidInfractionTimestamp());
        require(data.length == 0, InvalidData());
        slashedAmounts[operator] += amount;
        return amount;
    }
}
