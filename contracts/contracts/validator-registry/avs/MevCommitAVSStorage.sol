// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;


import {IMevCommitAVS} from "../../interfaces/IMevCommitAVS.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";

abstract contract MevCommitAVSStorage {

    /// @notice reference to eigenlayer core delegation manager contract
    IDelegationManager internal _delegationManager;

    /// @notice reference to eigenlayer core eigenpod manager contract
    IEigenPodManager internal _eigenPodManager;

    /// @notice reference to eigenlayer core strategy manager contract
    IStrategyManager internal _strategyManager;

    /// @notice reference to eigenlayer core AVS directory contract
    IAVSDirectory internal _eigenAVSDirectory;

    /// @notice Mapping of operator addresses to their registration info
    mapping(address => IMevCommitAVS.OperatorRegistrationInfo) public operatorRegistrations;

    /// @notice Mapping of validator pubkeys to their registration info
    mapping(bytes => IMevCommitAVS.ValidatorRegistrationInfo) public validatorRegistrations;

    /// @notice Mapping of LST restaker address to their registration info
    mapping(address => IMevCommitAVS.LSTRestakerRegistrationInfo) public lstRestakerRegistrations;

    /// @notice List of restakeable strategy addresses
    address[] public restakeableStrategies;

    /// @notice Address of the oracle responsible for freezing validators.
    address public freezeOracle;

    /// @notice Fee required to unfreeze a validator.
    uint256 public unfreezeFee;

    /// @notice Address that will receive unfreeze fees from frozen validators.
    address public unfreezeReceiver;

    /**
     * @notice Number of blocks a validator must remain frozen before it can be unfrozen.
     * This is param is optional to allow frozen validators to pay the fee immediately.
     */
    uint256 public unfreezePeriodBlocks;

    /// @notice Number of blocks an operator must wait after requesting deregistration before it is finalized.
    uint256 public operatorDeregPeriodBlocks;

    /// @notice Number of blocks a validator must wait after requesting deregistration before it is finalized.
    uint256 public validatorDeregPeriodBlocks;

    /// @notice Number of blocks a LST restaker must wait after requesting deregistration before it is finalized.
    uint256 public lstRestakerDeregPeriodBlocks;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
