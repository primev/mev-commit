// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

import { EventHeightLib } from "../utils/EventHeight.sol";

/// @title IVanillaRegistry
/// @notice Interface for the VanillaRegistry contract for validators.
interface IVanillaRegistry {

    /// @dev Struct representing a validator staked with the registry.
    struct StakedValidator {
        bool exists;
        address withdrawalAddress;
        uint256 balance;
        EventHeightLib.EventHeight unstakeHeight;
    }

    /// @dev Event emitted when a validator is staked.
    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when ETH is added to the staked balance a validator. 
    event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance);

    /// @dev Event emitted when a validator is unstaked.
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator's stake is withdrawn.
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator is slashed.
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when the min stake parameter is set.
    event MinStakeSet(address indexed msgSender, uint256 newMinStake);

    /// @dev Event emitted when the slash amount parameter is set.
    event SlashAmountSet(address indexed msgSender, uint256 newSlashAmount);

    /// @dev Event emitted when the slash oracle parameter is set.
    event SlashOracleSet(address indexed msgSender, address newSlashOracle);

    /// @dev Event emitted when the slash receiver parameter is set.
    event SlashReceiverSet(address indexed msgSender, address newSlashReceiver);

    /// @dev Event emitted when the unstake period blocks parameter is set.
    event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks);

    /// @dev Initializes the contract with the provided parameters.
    function initialize(
        uint256 _minStake, 
        uint256 _slashAmount,
        address _slashOracle,
        address _slashReceiver,
        uint256 _unstakePeriodBlocks, 
        address _owner
    ) external;

    /* 
     * @dev Stakes ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The validator BLS public keys to stake.
     */
    function stake(bytes[] calldata blsPubKeys) external payable;

    /* 
     * @dev Stakes ETH on behalf of one or multiple validators via their BLS pubkey,
     * and specifies an address other than msg.sender to be the withdrawal address.
     * @param blsPubKeys The validator BLS public keys to stake.
     * @param withdrawalAddress The address to receive the staked ETH.
     */
    function delegateStake(bytes[] calldata blsPubKeys, address withdrawalAddress) external payable;

    /* 
     * @dev Adds ETH to the staked balance of one or multiple validators via their BLS pubkey.
     * @dev A staking entry must already exist for each provided BLS pubkey.
     * @param blsPubKeys The BLS public keys to add stake to.
     */
    function addStake(bytes[] calldata blsPubKeys) external payable;

    /* 
     * @dev Unstakes ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to unstake.
     */
    function unstake(bytes[] calldata blsPubKeys) external;

    /* 
     * @dev Withdraws ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to withdraw.
     */
    function withdraw(bytes[] calldata blsPubKeys) external;

    /* 
     * @dev Allows oracle to slash some portion of stake for one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to slash.
     */
    function slash(bytes[] calldata blsPubKeys) external;

    /// @dev Enables the owner to pause the contract.
    function pause() external;

    /// @dev Enables the owner to unpause the contract.
    function unpause() external;

    /// @dev Enables the owner to set the minimum stake parameter.
    function setMinStake(uint256 newMinStake) external;

    /// @dev Enables the owner to set the slash amount parameter.
    function setSlashAmount(uint256 newSlashAmount) external;

    /// @dev Enables the owner to set the slash oracle parameter.
    function setSlashOracle(address newSlashOracle) external;

    /// @dev Enables the owner to set the slash receiver parameter.
    function setSlashReceiver(address newSlashReceiver) external;
    
    /// @dev Enables the owner to set the unstake period parameter.
    function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) external;

    /// @dev Returns true if a validator is considered "opted-in" to mev-commit via this registry.
    function isValidatorOptedIn(bytes calldata valBLSPubKey) external view returns (bool);

    /// @dev Returns stored staked validator struct for a given BLS pubkey.
    function getStakedValidator(bytes calldata valBLSPubKey) external view returns (StakedValidator memory);

    /// @dev Returns the staked amount for a given BLS pubkey.
    function getStakedAmount(bytes calldata valBLSPubKey) external view returns (uint256);

    /// @dev Returns true if a validator is currently unstaking.
    function isUnstaking(bytes calldata valBLSPubKey) external view returns (bool);

    /// @dev Returns the number of blocks remaining until an unstaking validator can withdraw their staked ETH.
    function getBlocksTillWithdrawAllowed(bytes calldata valBLSPubKey) external view returns (uint256);
}
