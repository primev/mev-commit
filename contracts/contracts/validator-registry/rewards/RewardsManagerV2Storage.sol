// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

abstract contract RewardsManagerV2Storage {

    uint256 public toTreasury;
    uint256 public rewardsPctBps;
    address payable public treasury;

    uint256[42] private __gap; // reserve slots for upgrades

}