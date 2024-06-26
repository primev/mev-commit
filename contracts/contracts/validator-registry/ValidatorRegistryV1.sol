// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {IValidatorRegistryV1} from "../interfaces/IValidatorRegistryV1.sol";
import {ValidatorRegistryV1Storage} from "./ValidatorRegistryV1Storage.sol";
import {EventHeightLib} from "../utils/EventHeight.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title Validator Registry v1
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit 
/// via simply staking ETH outside what's staked with the beacon chain.
contract ValidatorRegistryV1 is IValidatorRegistryV1, ValidatorRegistryV1Storage,
    OwnableUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    /// @dev Modifier to confirm a validator record exists for all provided BLS pubkeys.
    modifier onlyExistentValidatorRecords(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedValidators[blsPubKeys[i]].exists, "Validator record must exist");
        }
        _;
    }

    /// @dev Modifier to confirm a validator record does not exist for all provided BLS pubkeys.
    modifier onlyNonExistentValidatorRecords(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(!stakedValidators[blsPubKeys[i]].exists, "Validator record must NOT exist");
        }
        _;
    }

    /// @dev Modifier to confirm the sender is the withdrawal address for all provided BLS pubkeys.
    modifier onlyWithdrawalAddress(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedValidators[blsPubKeys[i]].withdrawalAddress == msg.sender, "Only withdrawal address can call this function");
        }
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are valid length.
    modifier onlyValidBLSPubKeys(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(blsPubKeys[i].length == 48, "Invalid BLS public key length. Must be 48 bytes");
        }
        _;
    }

    /// @dev Modifier to confirm the sender is the oracle account.
    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "Only slashing oracle account can call this function");
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Initializes the contract with the provided parameters.
    function initialize(
        uint256 _minStake, 
        uint256 _slashAmount,
        address _slashOracle,
        address _slashReceiver,
        uint256 _unstakePeriodBlocks, 
        address _owner
    ) external initializer {
        _setMinStake(_minStake);
        _setSlashAmount(_slashAmount);
        _setSlashOracle(_slashOracle);
        _setSlashReceiver(_slashReceiver);
        _setUnstakePeriodBlocks(_unstakePeriodBlocks);
        __Ownable_init(_owner);
    }

    /*
     * @dev implements _authorizeUpgrade from UUPSUpgradeable to enable only
     * the owner to upgrade the implementation contract.
     */
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /* 
     * @dev Stakes ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The validator BLS public keys to stake.
     */
    function stake(bytes[] calldata blsPubKeys) external payable
        onlyNonExistentValidatorRecords(blsPubKeys) onlyValidBLSPubKeys(blsPubKeys) whenNotPaused() {
        _stake(blsPubKeys, msg.sender);
    }

    /* 
     * @dev Stakes ETH on behalf of one or multiple validators via their BLS pubkey,
     * and specifies an address other than msg.sender to be the withdrawal address.
     * @param blsPubKeys The validator BLS public keys to stake.
     * @param withdrawalAddress The address to receive the staked ETH.
     */
    function delegateStake(bytes[] calldata blsPubKeys, address withdrawalAddress) external payable
        onlyNonExistentValidatorRecords(blsPubKeys) onlyValidBLSPubKeys(blsPubKeys) onlyOwner {
        require(withdrawalAddress != address(0), "Withdrawal address must be set");
        _stake(blsPubKeys, withdrawalAddress);
    }

    /* 
     * @dev Adds ETH to the staked balance of one or multiple validators via their BLS pubkey.
     * @dev A staking entry must already exist for each provided BLS pubkey.
     * @param blsPubKeys The BLS public keys to add stake to.
     */
    function addStake(bytes[] calldata blsPubKeys) external payable 
        onlyExistentValidatorRecords(blsPubKeys) whenNotPaused() {
        _addStake(blsPubKeys);
    }

    /* 
     * @dev Unstakes ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to unstake.
     */
    function unstake(bytes[] calldata blsPubKeys) external 
        onlyExistentValidatorRecords(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) whenNotPaused() {
        _unstake(blsPubKeys);
    }

    /* 
     * @dev Withdraws ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to withdraw.
     */
    function withdraw(bytes[] calldata blsPubKeys) external
        onlyExistentValidatorRecords(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) whenNotPaused() {
        _withdraw(blsPubKeys);
    }

    /* 
     * @dev Allows oracle to slash some portion of stake for one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to slash.
     */
    function slash(bytes[] calldata blsPubKeys) external
        onlyExistentValidatorRecords(blsPubKeys) onlySlashOracle whenNotPaused() {
        _slash(blsPubKeys);
    }

    /// @dev Enables the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Enables the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @dev Enables the owner to set the minimum stake parameter.
    function setMinStake(uint256 newMinStake) external onlyOwner {
        _setMinStake(newMinStake);
    }

    /// @dev Enables the owner to set the slash amount parameter.
    function setSlashAmount(uint256 newSlashAmount) external onlyOwner {
        _setSlashAmount(newSlashAmount);
    }

    /// @dev Enables the owner to set the slash oracle parameter.
    function setSlashOracle(address newSlashOracle) external onlyOwner {
        _setSlashOracle(newSlashOracle);
    }

    /// @dev Enables the owner to set the slash receiver parameter.
    function setSlashReceiver(address newSlashReceiver) external onlyOwner {
        _setSlashReceiver(newSlashReceiver);
    }
    
    /// @dev Enables the owner to set the unstake period parameter.
    function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) external onlyOwner {
        _setUnstakePeriodBlocks(newUnstakePeriodBlocks);
    }

    /*
     * @dev Internal function to stake ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The validator BLS public keys to stake.
     * @param withdrawalAddress The address to receive the staked ETH.
     */
    function _stake(bytes[] calldata blsPubKeys, address withdrawalAddress) internal {
        require(blsPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / blsPubKeys.length;
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            bytes calldata pubKey = blsPubKeys[i];
            stakedValidators[pubKey] = StakedValidator({
                exists: true,
                balance: splitAmount,
                withdrawalAddress: withdrawalAddress,
                unstakeHeight: EventHeightLib.EventHeight({ exists: false, blockHeight: 0 })
            });
            emit Staked(msg.sender, withdrawalAddress, pubKey, splitAmount);
        }
    }

    /* 
     * @dev Internal function to add ETH to the staked balance of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to add stake to.
     */
    function _addStake(bytes[] calldata blsPubKeys) internal {
        require(blsPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / blsPubKeys.length;
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            bytes calldata pubKey = blsPubKeys[i];
            stakedValidators[pubKey].balance += splitAmount;
            emit StakeAdded(msg.sender, stakedValidators[pubKey].withdrawalAddress,
                pubKey, splitAmount, stakedValidators[pubKey].balance);
        }
    }

    /* 
     * @dev Internal function to unstake ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to unstake.
     */
    function _unstake(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            _unstakeSingle(blsPubKeys[i]);
        }
    }

    /* 
     * @dev Internal function to unstake ETH on behalf of one validator via their BLS pubkey.
     * This function is neccessary for slashing. 
     * @param pubKey The single BLS public key to unstake.
     */
    function _unstakeSingle(bytes calldata pubKey) internal {
        require(!_isUnstaking(pubKey), "Unstake must NOT be initiated for validator");
        EventHeightLib.set(stakedValidators[pubKey].unstakeHeight, block.number);
        emit Unstaked(msg.sender, stakedValidators[pubKey].withdrawalAddress,
            pubKey, stakedValidators[pubKey].balance);
    }

    /* 
     * @dev Internal function to withdraw ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to withdraw.
     */
    function _withdraw(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            bytes calldata pubKey = blsPubKeys[i];
            require(_isUnstaking(pubKey), "Unstake must be initiated before withdrawal");
            require(block.number >= stakedValidators[pubKey].unstakeHeight.blockHeight + unstakePeriodBlocks,
                "withdrawal not allowed yet. Blocks requirement not met.");
            uint256 balance = stakedValidators[pubKey].balance;
            address withdrawalAddress = stakedValidators[pubKey].withdrawalAddress;
            delete stakedValidators[pubKey];
            payable(withdrawalAddress).transfer(balance);
            emit StakeWithdrawn(msg.sender, withdrawalAddress, pubKey, balance);
        }
    }

    /* 
     * @dev Internal function to slash ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to slash.
     */
    function _slash(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            bytes calldata pubKey = blsPubKeys[i];
            require(stakedValidators[pubKey].balance >= slashAmount,
                "Validator balance must be greater than or equal to slash amount");
            stakedValidators[pubKey].balance -= slashAmount;
            payable(slashReceiver).transfer(slashAmount);
            if (_isUnstaking(pubKey)) {
                // If validator is already unstaking, reset their unstake block number
                EventHeightLib.set(stakedValidators[pubKey].unstakeHeight, block.number);
            } else {
                _unstakeSingle(pubKey);
            }
            emit Slashed(msg.sender, slashReceiver,
                stakedValidators[pubKey].withdrawalAddress, pubKey, slashAmount);
        }
    }

    /// @dev Internal function to set the minimum stake parameter.
    function _setMinStake(uint256 newMinStake) internal {
        require(newMinStake > 0, "Minimum stake must be greater than 0");
        minStake = newMinStake;
        emit MinStakeSet(msg.sender, newMinStake);
    }

    /// @dev Internal function to set the slash amount parameter.
    function _setSlashAmount(uint256 newSlashAmount) internal {
        require(newSlashAmount >= 0, "Slash amount must be positive or 0");
        require(newSlashAmount <= minStake, "Slash amount must be less than or equal to minimum stake");
        slashAmount = newSlashAmount;
        emit SlashAmountSet(msg.sender, newSlashAmount);
    }

    /// @dev Internal function to set the slash oracle parameter.
    function _setSlashOracle(address newSlashOracle) internal {
        require(newSlashOracle != address(0), "Slash oracle must be set");
        slashOracle = newSlashOracle;
        emit SlashOracleSet(msg.sender, newSlashOracle);
    }

    /// @dev Internal function to set the slash receiver parameter.
    function _setSlashReceiver(address newSlashReceiver) internal {
        require(newSlashReceiver != address(0), "Slash receiver must be set");
        slashReceiver = newSlashReceiver;
        emit SlashReceiverSet(msg.sender, newSlashReceiver);
    }

    /// @dev Internal function to set the unstake period parameter.
    function _setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) internal {
        require(newUnstakePeriodBlocks > 0, "Unstake period must be greater than 0");
        unstakePeriodBlocks = newUnstakePeriodBlocks;
        emit UnstakePeriodBlocksSet(msg.sender, newUnstakePeriodBlocks);
    }

    /// @dev Returns true if a validator is considered "opted-in" to mev-commit via this registry.
    function isValidatorOptedIn(bytes calldata valBLSPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valBLSPubKey);
    }

    /// @dev Returns stored staked validator struct for a given BLS pubkey.
    function getStakedValidator(bytes calldata valBLSPubKey) external view returns (StakedValidator memory) {
        return stakedValidators[valBLSPubKey];
    }

    /// @dev Returns the staked amount for a given BLS pubkey.
    function getStakedAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        return stakedValidators[valBLSPubKey].balance;
    }

    /// @dev Returns true if a validator is currently unstaking.
    function isUnstaking(bytes calldata valBLSPubKey) external view returns (bool) {
        return _isUnstaking(valBLSPubKey);
    }

    /// @dev Returns the number of blocks remaining until an unstaking validator can withdraw their staked ETH.
    function getBlocksTillWithdrawAllowed(bytes calldata valBLSPubKey) external view returns (uint256) {
        require(_isUnstaking(valBLSPubKey), "Unstake must be initiated to check withdrawal eligibility");
        uint256 blocksSinceUnstakeInitiated = block.number - stakedValidators[valBLSPubKey].unstakeHeight.blockHeight;
        return blocksSinceUnstakeInitiated > unstakePeriodBlocks ? 0 : unstakePeriodBlocks - blocksSinceUnstakeInitiated;
    }

    /// @dev Internal function to check if a validator is considered "opted-in" to mev-commit via this registry.
    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal view returns (bool) {
        return !_isUnstaking(valBLSPubKey) && stakedValidators[valBLSPubKey].balance >= minStake;
    }

    /// @dev Internal function to check if a validator is currently unstaking.
    function _isUnstaking(bytes calldata valBLSPubKey) internal view returns (bool) {
        return stakedValidators[valBLSPubKey].unstakeHeight.exists;
    }

    /// @dev Fallback function to revert all calls, ensuring no unintended interactions.
    fallback() external payable {
        revert("Invalid call");
    }

    /// @dev Receive function is disabled for this contract to prevent unintended interactions.
    receive() external payable {
        revert("Invalid call");
    }
}
