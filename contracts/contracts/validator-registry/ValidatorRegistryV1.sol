// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title Validator Registry v1
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit 
/// via simply staking ETH outside what's staked with the beacon chain.
/// @dev Slashing is not yet implemented for this contract, hence it is upgradable to incorporate slashing in the future.
/// @dev This contract is meant to be deployed via UUPS proxy contract on mainnet.
contract ValidatorRegistryV1 is OwnableUpgradeable, UUPSUpgradeable {

    /// @dev Index tracking changes in the set of staked (opted-in) validators.
    /// This enables optimistic locking for batch queries.
    uint256 public stakedValsetVersion;

    /// @dev Minimum stake required for validators. 
    uint256 public minStake;
    
    uint256 public slashAmount; // TODO: test

    address public slashOracle; // TODO: test

    address public slashReceiver; // TODO: test

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
        uint256 _slashAmount,
        address _slashOracle,
        address _slashReceiver,
        uint256 _unstakePeriodBlocks, 
        address _owner
    ) external initializer {
        require(_minStake > 0, "Minimum stake must be greater than 0");
        require(_slashAmount > 0, "Slash amount must be greater than 0");
        require(_slashAmount <= _minStake, "Slash amount must be less than or equal to minimum stake");
        require(_slashOracle != address(0), "Slash oracle must be set");
        require(_slashReceiver != address(0), "Slash receiver must be set");
        require(_unstakePeriodBlocks > 0, "Unstake period must be greater than 0");

        minStake = _minStake;
        slashAmount = _slashAmount;
        slashOracle = _slashOracle;
        slashReceiver = _slashReceiver;
        unstakePeriodBlocks = _unstakePeriodBlocks;

        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Mapping of validator bls public key to staked balance.
    /// @dev Values in this mapping should always be greater than or equal to minStake.
    mapping(bytes => uint256) internal stakedBalances;

    /// @dev Mapping of validator bls public key to EOA withdrawal address. 
    mapping(bytes => address) public withdrawalAddresses;

    /// @dev Mapping of bls public key to block number of unstake initiation block.
    mapping(bytes => uint256) public unstakeBlockNums;

    /// @dev Mapping of bls public key to balance of currently unstaking ether.
    /// @dev Values in this mapping can be any positive value depending on amount staked and amount possibly slashed.
    mapping(bytes => uint256) public unstakingBalances;

    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    modifier onlyHasStakedBalance(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedBalances[blsPubKeys[i]] > 0, "Validator must have staked balance");
        }
        _;
    }

    modifier onlyHasUnstakingBalance(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(unstakingBalances[blsPubKeys[i]] > 0, "Validator must have unstaking balance");
        }
        _;
    }

    modifier onlyWithdrawalAddress(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(withdrawalAddresses[blsPubKeys[i]] == msg.sender, "Only withdrawal address can call this function");
        }
        _;
    }

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "Only slashing oracle account can call this function");
        _;
    }

    function stake(bytes[] calldata valBLSPubKeys) external payable {
        _stake(valBLSPubKeys, msg.sender);
    }

    function delegateStake(bytes[] calldata valBLSPubKeys, address withdrawalAddress) external payable onlyOwner {
        _stake(valBLSPubKeys, withdrawalAddress);
    }

    function _stake(bytes[] calldata valBLSPubKeys, address withdrawalAddress) internal {

        require(valBLSPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / valBLSPubKeys.length;
        require(splitAmount >= minStake, "Split amount must meet the minimum requirement");

        for (uint256 i = 0; i < valBLSPubKeys.length; i++) {

            _validateBLSPubKey(valBLSPubKeys[i]);
            require(unstakeBlockNums[valBLSPubKeys[i]] == 0, "validator cannot be staked with in-progress unstake process");
            require(stakedBalances[valBLSPubKeys[i]] == 0, "Validator already staked");

            stakedBalances[valBLSPubKeys[i]] = splitAmount;
            withdrawalAddresses[valBLSPubKeys[i]] = withdrawalAddress;
            emit Staked(msg.sender, withdrawalAddress, valBLSPubKeys[i], splitAmount);
        }
        ++stakedValsetVersion;
    }

    function unstake(bytes[] calldata blsPubKeys) external onlyHasStakedBalance(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) {
        _unstake(blsPubKeys);
    }

    function _unstake(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {

            _validateBLSPubKey(blsPubKeys[i]);
            require(unstakeBlockNums[blsPubKeys[i]] == 0, "Unstake already initiated for validator");

            uint256 balance = stakedBalances[blsPubKeys[i]];
            delete stakedBalances[blsPubKeys[i]];

            unstakeBlockNums[blsPubKeys[i]] = block.number;
            unstakingBalances[blsPubKeys[i]] = balance;

            emit Unstaked(msg.sender, withdrawalAddresses[blsPubKeys[i]], blsPubKeys[i], balance);
        }
        ++stakedValsetVersion;
    }

    function withdraw(bytes[] calldata blsPubKeys) external onlyHasUnstakingBalance(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) {
        _withdraw(blsPubKeys);
    }

    function _withdraw(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {

            _validateBLSPubKey(blsPubKeys[i]);
            require(unstakeBlockNums[blsPubKeys[i]] > 0, "Unstake must be initiated before withdrawal");
            require(block.number >= unstakeBlockNums[blsPubKeys[i]] + unstakePeriodBlocks,
                "withdrawal not allowed yet. Blocks requirement not met.");

            address withdrawalAddress = withdrawalAddresses[blsPubKeys[i]];
            withdrawalAddresses[blsPubKeys[i]] = address(0);
            unstakeBlockNums[blsPubKeys[i]] = 0;

            uint256 balance = unstakingBalances[blsPubKeys[i]];
            unstakingBalances[blsPubKeys[i]] = 0;
            payable(withdrawalAddress).transfer(balance);

            emit StakeWithdrawn(msg.sender, withdrawalAddress, blsPubKeys[i], balance);
        }
        // No need to increment stakedValsetVersion here, as stakedBalances map is not modified.
    }

    // TODO: test
    function slash(bytes[] calldata blsPubKeys) external onlyHasStakedBalance(blsPubKeys) onlySlashOracle {
        _slash(blsPubKeys);
    }

    // TODO: test
    function _slash(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            stakedBalances[blsPubKeys[i]] -= slashAmount;
            payable(slashReceiver).transfer(slashAmount);
            emit Slashed(msg.sender, slashReceiver, withdrawalAddresses[blsPubKeys[i]], blsPubKeys[i], slashAmount);
            _unstake(blsPubKeys);
        }
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

    // TODO: aggregator contract that exposes an isStaked func that'll call this and AVS contracts
}
