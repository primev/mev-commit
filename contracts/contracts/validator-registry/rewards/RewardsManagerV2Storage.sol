// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {IRewardsManagerV2} from "../../interfaces/IRewardsManagerV2.sol";

abstract contract RewardsManagerV2Storage {

    uint256 public toTreasury;
    uint256 public rewardsPctBps;
    address payable public treasury;

    uint256[42] private __gap; // reserve slots for upgrades

}