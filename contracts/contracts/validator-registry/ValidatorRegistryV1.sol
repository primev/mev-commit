// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title Validator Registry v1
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit 
/// via simply staking ETH outside what's staked with the beacon chain.
/// @dev Slashing is not yet implemented for this contract, hence it is upgradable to incorporate slashing in the future.
/// @dev This contract is meant to be deployed via UUPS proxy contract on mainnet.
contract ValidatorRegistryV1 is OwnableUpgradeable, ReentrancyGuardUpgradeable, UUPSUpgradeable {

    /// @dev Index tracking changes in the set of staked (opted-in) validators.
    /// This enables optimistic locking for batch queries.
    uint256 public stakedValsetVersion;

    /// @dev Minimum stake required for validators. 
    uint256 public minStake;

    /// @dev Number of blocks required between unstake initiation and withdrawal.
    uint256 public unstakePeriodBlocks;

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Receive function is disabled for this contract to prevent unintended interactions.
     */
    receive() external payable {
        revert("Invalid call");
    }

    function _authorizeUpgrade(address) internal override onlyOwner {}

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

    /// @dev Mapping of validator bls public key to staked balance. 
    mapping(bytes => uint256) internal stakedBalances;

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
        _stake(valBLSPubKeys, msg.sender);
    }

    function delegateStake(bytes[] calldata valBLSPubKeys, address stakeOriginator) external payable onlyOwner {
        _stake(valBLSPubKeys, stakeOriginator);
    }

    function _stake(bytes[] calldata valBLSPubKeys, address stakeOriginator) internal {

        require(valBLSPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / valBLSPubKeys.length;
        require(splitAmount >= minStake, "Split amount must meet the minimum requirement");

        for (uint256 i = 0; i < valBLSPubKeys.length; i++) {
            _validateBLSPubKey(valBLSPubKeys[i]);
            require(unstakeBlockNums[valBLSPubKeys[i]] == 0, "validator cannot be staked with in-progress unstake process");

            require(stakedBalances[valBLSPubKeys[i]] == 0, "Validator already staked");

            stakedBalances[valBLSPubKeys[i]] = splitAmount;
            stakeOriginators[valBLSPubKeys[i]] = stakeOriginator;
            emit Staked(stakeOriginator, valBLSPubKeys[i], splitAmount);
        }
        ++stakedValsetVersion;
    }

    function unstake(bytes[] calldata blsPubKeys) external nonReentrant {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            _validateBLSPubKey(blsPubKeys[i]);
            require(unstakeBlockNums[blsPubKeys[i]] == 0, "Unstake already initiated for validator");

            require(stakedBalances[blsPubKeys[i]] >= minStake, "No staked balance over min stake");

            require(stakeOriginators[blsPubKeys[i]] == msg.sender || owner() == msg.sender,
                "Not authorized to unstake validator. Must be stake originator or owner"); 

            uint256 balance = stakedBalances[blsPubKeys[i]];
            delete stakedBalances[blsPubKeys[i]];

            unstakeBlockNums[blsPubKeys[i]] = block.number;
            unstakingBalances[blsPubKeys[i]] = balance;

            emit Unstaked(msg.sender, blsPubKeys[i], balance);
        }
        ++stakedValsetVersion;
    }

    function withdraw(bytes[] calldata blsPubKeys) external nonReentrant {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            _validateBLSPubKey(blsPubKeys[i]);

            require(unstakeBlockNums[blsPubKeys[i]] > 0, "Unstake must be initiated before withdrawal");
            require(unstakingBalances[blsPubKeys[i]] >= minStake, "No unstaking balance over min stake");

            require(stakeOriginators[blsPubKeys[i]] == msg.sender || owner() == msg.sender,
                "Not authorized to withdraw stake. Must be stake originator or owner");
            require(block.number >= unstakeBlockNums[blsPubKeys[i]] + unstakePeriodBlocks, "withdrawal not allowed yet. Blocks requirement not met.");

            stakeOriginators[blsPubKeys[i]] = address(0);
            unstakeBlockNums[blsPubKeys[i]] = 0;

            uint256 balance = unstakingBalances[blsPubKeys[i]];
            unstakingBalances[blsPubKeys[i]] = 0;
            payable(msg.sender).transfer(balance);

            emit StakeWithdrawn(msg.sender, blsPubKeys[i], balance);
        }
        // No need to increment stakedValsetVersion here, as stakedBalances map is not modified.
    }

    function _validateBLSPubKey(bytes calldata valBLSPubKey) internal pure {
        require(valBLSPubKey.length == 48, "Invalid BLS public key length. Must be 48 bytes");
    }

    function getStakedAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        return stakedBalances[valBLSPubKey];
    }

    function isStaked(bytes calldata valBLSPubKey) external view returns (bool) {
        return stakedBalances[valBLSPubKey] >= minStake;
    }

    function getUnstakingAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        return unstakingBalances[valBLSPubKey];
    }

    function getBlocksTillWithdrawAllowed(bytes calldata valBLSPubKey) external view returns (uint256) {
        require(unstakeBlockNums[valBLSPubKey] > 0, "Unstake must be initiated to check withdrawal eligibility");
        uint256 blocksSinceUnstakeInitiated = block.number - unstakeBlockNums[valBLSPubKey];
        return blocksSinceUnstakeInitiated > unstakePeriodBlocks ? 0 : unstakePeriodBlocks - blocksSinceUnstakeInitiated;
    }

    function getStakedValsetVersion() external view returns (uint256) {
        return stakedValsetVersion;
    }

    // TODO: Slashing
    
    // TODO: aggregator contract that exposes an isStaked func that'll call this and AVS contracts
}
