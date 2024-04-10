// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract ValidatorRegistry is OwnableUpgradeable {

    uint256 public minStake;
    uint256 public unstakePeriodBlocks;

    constructor(uint256 _minStake, uint256 _unstakePeriodBlocks) {
        require(_minStake > 0, "Minimum stake must be greater than 0");
        require(_unstakePeriodBlocks > 0, "Unstake period must be greater than 0");
        minStake = _minStake;
        unstakePeriodBlocks = _unstakePeriodBlocks;
    }

    mapping(address => uint256) public stakedBalances;
    mapping(address => address) public stakeOriginators;
    mapping(address => uint256) public unstakeBlockNums;

    event SelfStaked(address indexed staker, uint256 amount);
    event SplitStaked(address indexed staker, address[] recipients, uint256 totalAmount);
    event Unstaked(address indexed staker, uint256 amount);
    event StakeWithdrawn(address indexed staker, uint256 amount);

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
}
