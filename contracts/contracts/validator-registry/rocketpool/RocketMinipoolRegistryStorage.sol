// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {RocketStorageInterface} from "rocketpool/contracts/interface/RocketStorageInterface.sol";
import {IRocketMinipoolRegistry} from "../../interfaces/IRocketMinipoolRegistry.sol";

abstract contract RocketMinipoolRegistryStorage {
    RocketStorageInterface public rocketStorage;

    /// @notice Number of seconds a validator must wait after requesting deregistration before it can be finalized.
    uint64 public deregistrationPeriod = 1 days;

    uint256 public unfreezeFee;
    address public freezeOracle;
    address public unfreezeReceiver;

    mapping(bytes => IRocketMinipoolRegistry.ValidatorRegistration) public validatorRegistrations;

}