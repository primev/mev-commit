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

    /// @dev Minimum stake required for validators. 
    uint256 public minStake;
    
    /// @dev Amount of ETH to slash per validator pubkey when a slash is invoked.
    uint256 public slashAmount;

    /// @dev Permissioned account that is able to invoke slashes.
    address public slashOracle; 

    /// @dev Account to receive all slashed ETH.
    address public slashReceiver;

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
        require(_slashAmount >= 0, "Slash amount must be positive or 0");
        require(_slashAmount <= _minStake, "Slash amount must be less than or equal to minimum stake");
        require(_slashOracle != address(0), "Slash oracle must be set");
        require(_slashReceiver != address(0), "Slash receiver must be set");
        require(_unstakePeriodBlocks > 0, "Unstake period must be greater than 0");
        require(_owner != address(0), "Owner must be set");

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

    struct StakedValidator {
        uint256 balance;
        address withdrawalAddress;
        uint256 unstakeBlockNum;
    }

    mapping(bytes => StakedValidator) public stakedValidators;

    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    modifier onlyHasStakingBalance(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedValidators[blsPubKeys[i]].balance > 0, "Validator must have staked balance");
        }
        _;
    }

    modifier onlyWithdrawalAddress(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedValidators[blsPubKeys[i]].withdrawalAddress == msg.sender, "Only withdrawal address can call this function");
        }
        _;
    }

    modifier onlyValidBLSPubKeys(bytes[] calldata blsPubKeys) {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(blsPubKeys[i].length == 48, "Invalid BLS public key length. Must be 48 bytes");
        }
        _;
    }

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "Only slashing oracle account can call this function");
        _;
    }

    function stake(bytes[] calldata valBLSPubKeys)
        external payable onlyValidBLSPubKeys(valBLSPubKeys) {
        _stake(valBLSPubKeys, msg.sender);
    }

    function delegateStake(bytes[] calldata valBLSPubKeys, address withdrawalAddress)
        external payable onlyOwner onlyValidBLSPubKeys(valBLSPubKeys) {
        _stake(valBLSPubKeys, withdrawalAddress);
    }

    function _stake(bytes[] calldata valBLSPubKeys, address withdrawalAddress) internal {

        require(valBLSPubKeys.length > 0, "There must be at least one recipient");
        uint256 splitAmount = msg.value / valBLSPubKeys.length;
        require(splitAmount >= minStake, "Split amount must meet the minimum requirement");

        for (uint256 i = 0; i < valBLSPubKeys.length; i++) {

            require(
                stakedValidators[valBLSPubKeys[i]].balance == 0 &&
                stakedValidators[valBLSPubKeys[i]].withdrawalAddress == address(0) &&
                stakedValidators[valBLSPubKeys[i]].unstakeBlockNum == 0,
                "Validator staking record must be empty"
            );

            stakedValidators[valBLSPubKeys[i]] = StakedValidator({
                balance: splitAmount,
                withdrawalAddress: withdrawalAddress,
                unstakeBlockNum: 0
            });
            emit Staked(msg.sender, withdrawalAddress, valBLSPubKeys[i], splitAmount);
        }
    }

    function unstake(bytes[] calldata blsPubKeys) external 
        onlyHasStakingBalance(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) {
        _unstake(blsPubKeys);
    }

    function _unstake(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            require(stakedValidators[blsPubKeys[i]].unstakeBlockNum == 0, "Unstake already initiated for validator");
            stakedValidators[blsPubKeys[i]].unstakeBlockNum = block.number;
            emit Unstaked(msg.sender, stakedValidators[blsPubKeys[i]].withdrawalAddress,
                blsPubKeys[i], stakedValidators[blsPubKeys[i]].balance);
        }
    }

    function withdraw(bytes[] calldata blsPubKeys) external
        onlyHasStakingBalance(blsPubKeys) onlyWithdrawalAddress(blsPubKeys) {
        _withdraw(blsPubKeys);
    }

    function _withdraw(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {

            require(stakedValidators[blsPubKeys[i]].unstakeBlockNum > 0, "Unstake must be initiated before withdrawal");
            require(block.number >= stakedValidators[blsPubKeys[i]].unstakeBlockNum + unstakePeriodBlocks,
                "withdrawal not allowed yet. Blocks requirement not met.");

            uint256 balance = stakedValidators[blsPubKeys[i]].balance;
            address withdrawalAddress = stakedValidators[blsPubKeys[i]].withdrawalAddress;
            delete stakedValidators[blsPubKeys[i]];

            payable(withdrawalAddress).transfer(balance);

            emit StakeWithdrawn(msg.sender, withdrawalAddress, blsPubKeys[i], balance);
        }
    }

    // TODO: test
    function slash(bytes[] calldata blsPubKeys) external onlySlashOracle {
        _slash(blsPubKeys);
    }

    // TODO: test
    function _slash(bytes[] calldata blsPubKeys) internal {
        for (uint256 i = 0; i < blsPubKeys.length; i++) {
            StakedValidator storage stakedValidator = stakedValidators[blsPubKeys[i]];
            require(stakedValidator.balance >= slashAmount,
                "Validator balance must be greater than or equal to slash amount");

            stakedValidator.balance -= slashAmount;
            payable(slashReceiver).transfer(slashAmount);
            if (_isUnstaking(blsPubKeys[i])) {
                // If validator is already unstaking, reset their unstake block number
                stakedValidator.unstakeBlockNum = block.number;
            } else {
                _unstake(blsPubKeys);
            }
            emit Slashed(msg.sender, slashReceiver, stakedValidator.withdrawalAddress, blsPubKeys[i], slashAmount);
        }
    }

    function getStakedValidator(bytes calldata valBLSPubKey) external view returns (StakedValidator memory) {
        return stakedValidators[valBLSPubKey];
    }

    function getStakedAmount(bytes calldata valBLSPubKey) external view returns (uint256) {
        return stakedValidators[valBLSPubKey].balance;
    }

    function isValidatorOptedIn(bytes calldata valBLSPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valBLSPubKey);
    }

    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal view returns (bool) {
        return !_isUnstaking(valBLSPubKey) && stakedValidators[valBLSPubKey].balance >= minStake;
    }

    function isUnstaking(bytes calldata valBLSPubKey) external view returns (bool) {
        return _isUnstaking(valBLSPubKey);
    }

    function _isUnstaking(bytes calldata valBLSPubKey) internal view returns (bool) {
        return stakedValidators[valBLSPubKey].unstakeBlockNum > 0;
    }

    function getBlocksTillWithdrawAllowed(bytes calldata valBLSPubKey) external view returns (uint256) {
        require(_isUnstaking(valBLSPubKey), "Unstake must be initiated to check withdrawal eligibility");
        uint256 blocksSinceUnstakeInitiated = block.number - stakedValidators[valBLSPubKey].unstakeBlockNum;
        return blocksSinceUnstakeInitiated > unstakePeriodBlocks ? 0 : unstakePeriodBlocks - blocksSinceUnstakeInitiated;
    }
    // TODO: aggregator contract that exposes an isStaked func that'll call this and AVS contracts
}
