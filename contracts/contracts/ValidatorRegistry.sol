// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {EnumerableMap} from "./utils/EnumerableMap.sol";

/// @title Validator Registry
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit via staking. 
/// @dev Slashing is not yet implemented for this contract, hence it is upgradable to incorporate slashing in the future.
/// @dev This contract is meant to be deployed via a proxy contract.
contract ValidatorRegistry is OwnableUpgradeable, ReentrancyGuardUpgradeable {

    uint256 public minStake;
    uint256 public unstakePeriodBlocks;

    function initialize(
        uint256 _minStake, 
        uint256 _unstakePeriodBlocks, 
        address _owner
    ) external initializer {
        require(_minStake > 0, "Minimum stake must be greater than 0");
        require(_unstakePeriodBlocks > 0, "Unstake period must be greater than 0");
        minStake = _minStake;
        unstakePeriodBlocks = _unstakePeriodBlocks;
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    using EnumerableMap for EnumerableMap.BytesToUint256Map;
    /// @dev Enumerable mapping of validator bls public key to staked balance. 
    EnumerableMap.BytesToUint256Map internal stakedBalances;

    /// @dev Mapping of validator bls public key to EOA stake originator. 
    mapping(bytes => address) public stakeOriginators;

    /// @dev Mapping of bls public key to block number of unstake initiation block.
    mapping(bytes => uint256) public unstakeBlockNums;

    /// @dev Mapping of bls public key to balance of currently unstaking ether.
    mapping(bytes => uint256) public unstakingBalances;

    event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount);
    event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount);

    function stake(bytes[] calldata valBLSPubKeys) external payable {

        require(valBLSPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / valBLSPubKeys.length;
        require(splitAmount >= minStake, "Split amount must meet the minimum requirement");

        for (uint256 i = 0; i < valBLSPubKeys.length; i++) {
            _validateBLSPubKey(valBLSPubKeys[i]);
            require(unstakeBlockNums[valBLSPubKeys[i]] == 0, "validator cannot be staked with in-progress unstake process");

            bool exists;
            uint256 value;
            (exists, value) = stakedBalances.tryGet(valBLSPubKeys[i]);
            require(!exists, "Validator already staked");

            stakedBalances.set(valBLSPubKeys[i], splitAmount);
            stakeOriginators[valBLSPubKeys[i]] = msg.sender;
            emit Staked(msg.sender, valBLSPubKeys[i], splitAmount);
        }
    }

    function unstake(bytes[] calldata blsPubKeys) external nonReentrant {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            _validateBLSPubKey(blsPubKeys[i]);
            require(unstakeBlockNums[blsPubKeys[i]] == 0, "Unstake already initiated for validator");

            bool exists;
            uint256 balance;
            (exists, balance) = stakedBalances.tryGet(blsPubKeys[i]);
            require(exists, "Validator not staked");
            require(balance >= minStake, "No staked balance over min stake");

            require(stakeOriginators[blsPubKeys[i]] == msg.sender, "Not authorized to unstake validator. Must be stake originator"); 

            bool removed = stakedBalances.remove(blsPubKeys[i]);
            require(removed, "Failed to remove staked balance");

            unstakeBlockNums[blsPubKeys[i]] = block.number;
            unstakingBalances[blsPubKeys[i]] = balance;

            emit Unstaked(msg.sender, blsPubKeys[i], balance);
        }
    }

    function withdraw(bytes[] calldata blsPubKeys) external nonReentrant {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            _validateBLSPubKey(blsPubKeys[i]);

            require(unstakeBlockNums[blsPubKeys[i]] > 0, "Unstake must be initiated before withdrawal");
            require(unstakingBalances[blsPubKeys[i]] >= minStake, "No unstaking balance over min stake");

            require(stakeOriginators[blsPubKeys[i]] == msg.sender , "Not authorized to withdraw stake. Must be stake originator");
            require(block.number >= unstakeBlockNums[blsPubKeys[i]] + unstakePeriodBlocks, "withdrawal not allowed yet. Blocks requirement not met.");

            stakeOriginators[blsPubKeys[i]] = address(0);
            unstakeBlockNums[blsPubKeys[i]] = 0;

            uint256 balance = unstakingBalances[blsPubKeys[i]];
            unstakingBalances[blsPubKeys[i]] = 0;
            payable(msg.sender).transfer(balance);

            emit StakeWithdrawn(msg.sender, blsPubKeys[i], balance);
        }
    }

    function _validateBLSPubKey(bytes calldata valBLSPubKey) internal pure {
        require(valBLSPubKey.length == 48, "Invalid BLS public key length. Must be 48 bytes");
    }

    function getStakedAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        bool exists;
        uint256 balance;
        (exists, balance) = stakedBalances.tryGet(valBLSPubKey);
        return exists ? balance : 0;
    }

    function isStaked(bytes calldata valBLSPubKey) external view returns (bool) {
        bool exists;
        (exists,) = stakedBalances.tryGet(valBLSPubKey);
        return exists;
    }

    function getUnstakingAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        return unstakingBalances[valBLSPubKey];
    }

    function getNumberOfStakedValidators() external view returns (uint256) {
        return stakedBalances.length();
    }

    /// @dev Returns an array of staked validator BLS pubkeys within the specified range.
    /// For now we leave out a function for obtaining all validator pubkeys at once, to prevent possible DOS scenarios.
    /// Insead we require users to query in batches, with knowledge from getNumberOfStakedValidators().
    ///
    /// TODO: Practical race conditions between user and chain, caused by state transitions during batch queries, 
    /// need to be explored further.
    function getStakedValidators(uint256 start, uint256 end) external view returns (bytes[] memory) {
        require(start < end, "Invalid range");
        require(end <= stakedBalances.length(), "Range exceeds staked balances length");
        require(end - start <= 1000, "Range must be less than or equal to 1000");

        bytes[] memory keys = new bytes[](end - start);
        for (uint256 i = start; i < end; i++) {
            (bytes memory key,) = stakedBalances.at(i);
            keys[i - start] = key;
        }
        return keys;
    }
}
