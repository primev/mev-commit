// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/// @title Validator Registry
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit via staking. 
/// @dev Slashing is not yet implemented for this contract, hence it is upgradable to incorporate slashing in the future.
/// @dev This contract is meant to be deployed via a proxy contract.
contract ValidatorRegistry is OwnableUpgradeable {

    // TODO: add reentrancy gaurd

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

    mapping(address => uint256) public stakedBalances;
    mapping(address => address) public stakeOriginators;
    mapping(address => uint256) public unstakeBlockNums;

    event SelfStaked(address indexed txOriginator, uint256 amount);
    event SplitStaked(address indexed txOriginator, address[] recipients, uint256 totalAmount);
    event Unstaked(address indexed txOriginator, uint256 amount);
    event StakeWithdrawn(address indexed txOriginator, uint256 amount);

    function selfStake() external payable {
        require(msg.value >= minStake, "Stake amount must meet the minimum requirement");
        require(stakedBalances[msg.sender] == 0, "Already staked");

        stakedBalances[msg.sender] += msg.value;
        stakeOriginators[msg.sender] = msg.sender;

        emit SelfStaked(msg.sender, msg.value);
    }

    function splitStake(address[] calldata recipients) external payable {
        require(recipients.length > 0, "There must be at least one recipient");

        uint256 splitAmount = msg.value / recipients.length;
        require(splitAmount >= minStake, "Split amount must meet the minimum requirement");

        for (uint256 i = 0; i < recipients.length; i++) {
            require(stakedBalances[recipients[i]] == 0, "Recipient already staked");
            stakedBalances[recipients[i]] += splitAmount;
            stakeOriginators[recipients[i]] = msg.sender;
        }

        emit SplitStaked(msg.sender, recipients, msg.value);
    }

    function unstake(address[] calldata fromAddrs) external {
        for (uint256 i = 0; i < fromAddrs.length; i++) {
            require(stakedBalances[fromAddrs[i]] > 0, "No balance to unstake");
            require(stakeOriginators[fromAddrs[i]] == msg.sender || fromAddrs[i] == msg.sender, "Not authorized to unstake. Must be stake originator or EOA whos staked");

            unstakeBlockNums[fromAddrs[i]] = block.number;
            emit Unstaked(msg.sender, stakedBalances[fromAddrs[i]]);
        }
    }

    function withdraw(address[] calldata fromAddrs) external {
        for (uint256 i = 0; i < fromAddrs.length; i++) {
            require(stakedBalances[fromAddrs[i]] > 0, "No staked balance to withdraw");
            require(stakeOriginators[fromAddrs[i]] == msg.sender || fromAddrs[i] == msg.sender, "Not authorized to withdraw. Must be stake originator or EOA whos staked");
            require(block.number >= unstakeBlockNums[fromAddrs[i]] + unstakePeriodBlocks, "withdrawal not allowed yet. Blocks requirement not met.");

            uint256 amount = stakedBalances[fromAddrs[i]];
            stakedBalances[fromAddrs[i]] -= amount;
            (bool sent, ) = msg.sender.call{value: amount}("");
            require(sent, "Failed to withdraw stake");
            stakeOriginators[fromAddrs[i]] = address(0);

            emit StakeWithdrawn(msg.sender, amount);
        }
    }

    function isStaked(address staker) external view returns (bool) {
        return stakedBalances[staker] >= minStake;
    }

    function getStakedAmount(address staker) external view returns (uint256) {
        return stakedBalances[staker];
    }
}
