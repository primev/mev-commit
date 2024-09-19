// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {MockEntity} from "./MockEntity.sol";
import {MockDelegator} from "./MockDelegator.sol";

contract MockVetoSlasher is MockEntity {
    address private _resolver;
    uint256 private _vetoDuration;
    mapping(address operator => uint256 slashedAmount) public slashedAmounts;
    address[] public slashedOperators;
    uint256 private _slashIndex;
    MockDelegator public mockDelegator;

    constructor(uint64 type_, address resolver_, uint256 vetoDuration_, MockDelegator mockDelegator_) MockEntity(type_) {
        _resolver = resolver_;
        _vetoDuration = vetoDuration_;
        mockDelegator = MockDelegator(mockDelegator_);
    }

    error InvalidSubnetwork();
    error InvalidOperator();
    error InvalidAmount();
    error InvalidInfractionTimestamp();
    error InvalidData();
    error InsufficientStake();
    error InvalidHints();

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
        slashedOperators.push(operator);
        return _slashIndex++;
    }

    function executeSlash(
        uint256 slashIndex,
        bytes calldata hints
    ) external returns (uint256 slashedAmount) {
        require(hints.length == 0, InvalidHints());
        address operator = slashedOperators[slashIndex];
        uint256 amount = slashedAmounts[operator];
        uint256 stake = mockDelegator.stake(bytes32("subnet"), operator);
        require(stake >= amount, InsufficientStake());
        mockDelegator.setStake(operator, stake - amount);
        slashedAmounts[operator] = 0;
        slashedOperators[slashIndex] = address(0);
        return amount;
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
