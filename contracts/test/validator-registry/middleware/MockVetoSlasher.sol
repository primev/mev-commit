// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";

contract MockVetoSlasher is MockEntity {
    address private _resolver;
    uint256 private _vetoDuration;
    mapping(address operator => uint256 slashedAmount) public slashedAmounts;
    uint256 private _slashIndex;

    constructor(uint64 type_, address resolver_, uint256 vetoDuration_) MockEntity(type_) {
        _resolver = resolver_;
        _vetoDuration = vetoDuration_;
    }

    error InvalidSubnetwork();
    error InvalidOperator();
    error InvalidAmount();
    error InvalidInfractionTimestamp();
    error InvalidData();

    function requestSlash(
        bytes32 subnetwork,
        address operator,
        uint256 amount,
        uint48 infractionTimestamp,
        bytes memory data
    ) external returns (uint256 slashIndex) {
        require(subnetwork != bytes32(0), InvalidSubnetwork());
        require(operator != address(0), InvalidOperator());
        require(amount != 0, InvalidAmount());
        require(infractionTimestamp != 0, InvalidInfractionTimestamp());
        require(data.length == 0, InvalidData());
        slashedAmounts[operator] += amount;
        return _slashIndex++;
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
