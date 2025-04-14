// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {MockEntity} from "./MockEntity.sol";
import {MockDelegator} from "./MockDelegator.sol";

contract MockInstantSlasher is MockEntity {
    constructor(uint64 type_, MockDelegator mockDelegator_, bool isBurnerHook_) MockEntity(type_) {
        mockDelegator = MockDelegator(mockDelegator_);
        _isBurnerHook = isBurnerHook_;
    }
    MockDelegator public mockDelegator;
    bool public _isBurnerHook;

    mapping(address operator => uint256 slashedAmount) public slashedAmounts;

    error InvalidSubnetwork();
    error InvalidOperator();
    error InvalidAmount();
    error InvalidInfractionTimestamp();
    error InvalidData();
    error InsufficientStake();

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
        uint256 stake = mockDelegator.stake(subnetwork, operator);
        require(stake >= amount, InsufficientStake());
        mockDelegator.setStake(operator, stake - amount);
        return amount;
    }

    function isBurnerHook() external view returns (bool) {
        return _isBurnerHook;
    }
}
