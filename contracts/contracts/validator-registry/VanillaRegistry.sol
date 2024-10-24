// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IVanillaRegistry} from "../interfaces/IVanillaRegistry.sol";
import {VanillaRegistryStorage} from "./VanillaRegistryStorage.sol";
import {BlockHeightOccurrence} from "../utils/Occurrence.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Errors} from "../utils/Errors.sol";
import {FeePayout} from "../utils/FeePayout.sol";

/// @title Vanilla Registry
/// @notice Logic contract enabling L1 validators to opt-in to mev-commit 
/// via simply staking ETH outside what's staked with the beacon chain.
contract VanillaRegistry is IVanillaRegistry, VanillaRegistryStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    /// @dev Modifier to confirm a validator record exists for all provided BLS pubkeys.
    modifier onlyExistentValidatorRecords(bytes[] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            require(stakedValidators[blsPubKeys[i]].exists, IVanillaRegistry.ValidatorRecordMustExist(blsPubKeys[i]));
        }
        _;
    }

    /// @dev Modifier to confirm a validator record does not exist for all provided BLS pubkeys.
    modifier onlyNonExistentValidatorRecords(bytes[] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            require(!stakedValidators[blsPubKeys[i]].exists, IVanillaRegistry.ValidatorRecordMustNotExist(blsPubKeys[i]));
        }
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are NOT unstaking.
    modifier onlyNotUnstaking(bytes[] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            require(!_isUnstaking(blsPubKeys[i]), IVanillaRegistry.ValidatorCannotBeUnstaking(blsPubKeys[i]));
        }
        _;
    }

    /// @dev Modifier to confirm the sender is the withdrawal address for all provided BLS pubkeys.
    modifier onlyWithdrawalAddress(bytes[] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            IVanillaRegistry.StakedValidator storage validator = stakedValidators[blsPubKeys[i]];
            require(validator.withdrawalAddress == msg.sender,
                IVanillaRegistry.SenderIsNotWithdrawalAddress(msg.sender, validator.withdrawalAddress));
        }
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are valid length.
    modifier onlyValidBLSPubKeys(bytes[] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            require(blsPubKeys[i].length == 48, IVanillaRegistry.InvalidBLSPubKeyLength(48, blsPubKeys[i].length));
        }
        _;
    }

    /// @dev Modifier to confirm the sender is the oracle account.
    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, IVanillaRegistry.SenderIsNotSlashOracle(msg.sender, slashOracle));
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
        address _slashOracle,
        address _slashReceiver,
        uint256 _unstakePeriodBlocks, 
        uint256 _slashingPayoutPeriodBlocks,
        address _owner
    ) external initializer {
        _setMinStake(_minStake);
        _setSlashOracle(_slashOracle);
        _setUnstakePeriodBlocks(_unstakePeriodBlocks);
        FeePayout.init(slashingFundsTracker, _slashReceiver, _slashingPayoutPeriodBlocks);
        __Ownable_init(_owner);
    }

    /// @dev Receive function is disabled for this contract to prevent unintended interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to revert all calls, ensuring no unintended interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

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
        require(withdrawalAddress != address(0), IVanillaRegistry.WithdrawalAddressMustBeSet());
        _stake(blsPubKeys, withdrawalAddress);
    }

    /* 
     * @dev Adds ETH to the staked balance of one or multiple validators via their BLS pubkey.
     * @dev A staking entry must already exist for each provided BLS pubkey.
     * @param blsPubKeys The BLS public keys to add stake to.
     */
    function addStake(bytes[] calldata blsPubKeys) external payable 
        onlyExistentValidatorRecords(blsPubKeys) onlyWithdrawalAddress(blsPubKeys)
            onlyNotUnstaking(blsPubKeys) whenNotPaused() {
        _addStake(blsPubKeys);
    }

    /* 
     * @dev Unstakes ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to unstake.
     */
    function unstake(bytes[] calldata blsPubKeys) external 
        onlyExistentValidatorRecords(blsPubKeys) onlyWithdrawalAddress(blsPubKeys)
            onlyNotUnstaking(blsPubKeys) whenNotPaused() {
        _unstake(blsPubKeys);
    }

    /// @dev Allows owner to withdraw ETH on behalf of one or multiple validators via their BLS pubkey.
    /// @param blsPubKeys The BLS public keys to withdraw.
    /// @dev msg.sender must be the withdrawal address for every provided validator pubkey as enforced in _withdraw.
    function withdraw(bytes[] calldata blsPubKeys) external
        onlyExistentValidatorRecords(blsPubKeys) whenNotPaused() {
        uint256 totalAmount = _withdraw(blsPubKeys, msg.sender);
        if (totalAmount != 0) {
            (bool success, ) = msg.sender.call{value: totalAmount}("");
            require(success, IVanillaRegistry.WithdrawalFailed());
        }
    }

    /// @dev Allows owner to withdraw ETH on behalf of one or multiple validators via their BLS pubkey.
    /// @param blsPubKeys The BLS public keys to withdraw.
    /// @param withdrawalAddress The address to receive the staked ETH.
    /// @dev withdrawalAddress must be the withdrawal address for every provided validator pubkeyas enforced in _withdraw.
    function forceWithdrawalAsOwner(bytes[] calldata blsPubKeys, address withdrawalAddress) external
        onlyExistentValidatorRecords(blsPubKeys) onlyOwner {
        uint256 totalAmount = _withdraw(blsPubKeys, withdrawalAddress);
        if (totalAmount != 0) {
            forceWithdrawnFunds[withdrawalAddress] += totalAmount;
        }
    }

    /// @dev Allows a withdrawal address to claim any ETH that was force withdrawn by the owner.
    function claimForceWithdrawnFunds() external {
        uint256 amountToClaim = forceWithdrawnFunds[msg.sender];
        require(amountToClaim != 0, IVanillaRegistry.NoFundsToWithdraw());
        forceWithdrawnFunds[msg.sender] = 0;
        (bool success, ) = msg.sender.call{value: amountToClaim}("");
        require(success, IVanillaRegistry.WithdrawalFailed());
    }

    /// @dev Allows oracle to slash some portion of stake for one or multiple validators via their BLS pubkey.
    /// @param blsPubKeys The BLS public keys to slash.
    /// @param payoutIfDue Whether to payout slashed funds to receiver if the payout period is due.
    function slash(bytes[] calldata blsPubKeys, bool payoutIfDue) external
        onlyExistentValidatorRecords(blsPubKeys) onlySlashOracle whenNotPaused() {
        _slash(blsPubKeys, payoutIfDue);
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

    /// @dev Enables the owner to set the slashing payout period parameter.
    function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) external onlyOwner {
        _setSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks);
    }

    /// @dev Enables the owner to manually transfer slashing funds.
    function manuallyTransferSlashingFunds() external onlyOwner {
        FeePayout.transferToRecipient(slashingFundsTracker);
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
        require(_isUnstaking(valBLSPubKey), IVanillaRegistry.MustUnstakeToWithdraw());
        uint256 blocksSinceUnstakeInitiated = block.number - stakedValidators[valBLSPubKey].unstakeOccurrence.blockHeight;
        return blocksSinceUnstakeInitiated > unstakePeriodBlocks ? 0 : unstakePeriodBlocks - blocksSinceUnstakeInitiated;
    }

    /// @dev Returns true if the slashing payout period is due.
    function isSlashingPayoutDue() external view returns (bool) {
        return FeePayout.isPayoutDue(slashingFundsTracker);
    }

    function getAccumulatedSlashingFunds() external view returns (uint256) {
        return slashingFundsTracker.accumulatedAmount;
    }

    /*
     * @dev implements _authorizeUpgrade from UUPSUpgradeable to enable only
     * the owner to upgrade the implementation contract.
     */
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @dev Internal function that splits msg.value stake to apply an action for each validator.
    function _splitStakeAndApplyAction(
        bytes[] calldata blsPubKeys,
        address withdrawalAddress,
        function(bytes calldata, uint256, address) internal action
    ) internal {
        require(blsPubKeys.length != 0, IVanillaRegistry.AtLeastOneRecipientRequired());
        uint256 baseStakeAmount = msg.value / blsPubKeys.length;
        uint256 lastStakeAmount = msg.value - (baseStakeAmount * (blsPubKeys.length - 1));
        uint256 numKeys = blsPubKeys.length;
        for (uint256 i = 0; i < numKeys; ++i) {
            uint256 stakeAmount = (i == numKeys - 1) ? lastStakeAmount : baseStakeAmount;
            action(blsPubKeys[i], stakeAmount, withdrawalAddress);
        }
    }

    /*
     * @dev Internal function to stake ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The validator BLS public keys to stake.
     * @param withdrawalAddress The address to receive the staked ETH.
     */
    function _stake(bytes[] calldata blsPubKeys, address withdrawalAddress) internal {
        // At least minStake must be staked for each pubkey.
        require(msg.value >= minStake * blsPubKeys.length, IVanillaRegistry.StakeTooLowForNumberOfKeys(msg.value, minStake * blsPubKeys.length));
        _splitStakeAndApplyAction(blsPubKeys, withdrawalAddress, _stakeAction);
    }

    /// @dev Internal function that creates a staked validator record and emits a Staked event.
    function _stakeAction(bytes calldata pubKey, uint256 stakeAmount, address withdrawalAddress) internal {
        stakedValidators[pubKey] = StakedValidator({
            exists: true,
            balance: stakeAmount,
            withdrawalAddress: withdrawalAddress,
            unstakeOccurrence: BlockHeightOccurrence.Occurrence({ exists: false, blockHeight: 0 })
        });
        emit Staked(msg.sender, withdrawalAddress, pubKey, stakeAmount);
    }

    /* 
     * @dev Internal function to add ETH to the staked balance of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to add stake to.
     */
    function _addStake(bytes[] calldata blsPubKeys) internal {
        // At least 1 wei must be added for each pubkey.
        require(msg.value >= blsPubKeys.length, IVanillaRegistry.StakeTooLowForNumberOfKeys(msg.value, blsPubKeys.length));
        _splitStakeAndApplyAction(blsPubKeys, address(0), _addStakeAction);
    }

    /// @dev Internal function that adds stake to an already existing validator record, emitting a StakeAdded event.
    function _addStakeAction(bytes calldata pubKey, uint256 stakeAmount, address) internal {
        IVanillaRegistry.StakedValidator storage validator = stakedValidators[pubKey];
        validator.balance += stakeAmount;
        emit StakeAdded(msg.sender, validator.withdrawalAddress, pubKey, stakeAmount, validator.balance);
    }

    /* 
     * @dev Internal function to unstake ETH on behalf of one or multiple validators via their BLS pubkey.
     * @param blsPubKeys The BLS public keys to unstake.
     */
    function _unstake(bytes[] calldata blsPubKeys) internal {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _unstakeSingle(blsPubKeys[i]);
        }
    }

    /* 
     * @dev Internal function to unstake ETH on behalf of one validator via their BLS pubkey.
     * This function is necessary for slashing. 
     * @param pubKey The single BLS public key to unstake.
     */
    function _unstakeSingle(bytes calldata pubKey) internal {
        IVanillaRegistry.StakedValidator storage validator = stakedValidators[pubKey];
        BlockHeightOccurrence.captureOccurrence(validator.unstakeOccurrence);
        emit Unstaked(msg.sender, validator.withdrawalAddress, pubKey, validator.balance);
    }


    /// @dev Internal function to withdraw ETH on behalf of one or multiple validators via their BLS pubkey.
    /// @dev This function also deletes the validator record, and therefore serves a purpose even if no withdawable funds exist.
    /// @param blsPubKeys The BLS public keys to withdraw.
    /// @param expectedWithdrawalAddress The expected withdrawal address for every provided validator. 
    /// @return totalAmount The total amount of ETH withdrawn, to be handled by calling function.
    /// @dev msg.sender must be contract owner, or the withdrawal address for every provided validator.
    function _withdraw(bytes[] calldata blsPubKeys, address expectedWithdrawalAddress) internal returns (uint256) {
        uint256 len = blsPubKeys.length;
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubKey = blsPubKeys[i];
            IVanillaRegistry.StakedValidator storage validator = stakedValidators[pubKey];
            require(_isUnstaking(pubKey), IVanillaRegistry.MustUnstakeToWithdraw());
            require(block.number > validator.unstakeOccurrence.blockHeight + unstakePeriodBlocks,
                IVanillaRegistry.WithdrawingTooSoon());
            require(validator.withdrawalAddress == expectedWithdrawalAddress,
                IVanillaRegistry.WithdrawalAddressMismatch(validator.withdrawalAddress, expectedWithdrawalAddress));
            uint256 balance = validator.balance;
            totalAmount += balance;
            delete stakedValidators[pubKey];
            emit StakeWithdrawn(msg.sender, expectedWithdrawalAddress, pubKey, balance);
        }
        emit TotalStakeWithdrawn(msg.sender, expectedWithdrawalAddress, totalAmount);
        return totalAmount;
    }

    /// @dev Internal function to slash minStake worth of ETH on behalf of one or multiple validators via their BLS pubkey.
    /// @param blsPubKeys The BLS public keys to slash.
    /// @param payoutIfDue Whether to payout slashed funds to receiver if the payout period is due.
    function _slash(bytes[] calldata blsPubKeys, bool payoutIfDue) internal {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubKey = blsPubKeys[i];
            IVanillaRegistry.StakedValidator storage validator = stakedValidators[pubKey];
            require(validator.balance >= minStake, IVanillaRegistry.NotEnoughBalanceToSlash());
            if (!_isUnstaking(pubKey)) { 
                _unstakeSingle(pubKey);
            }
            validator.balance -= minStake;
            slashingFundsTracker.accumulatedAmount += minStake;
            bool isLastEntry = i == len - 1;
            if (payoutIfDue && FeePayout.isPayoutDue(slashingFundsTracker) && isLastEntry) {
                FeePayout.transferToRecipient(slashingFundsTracker);
            }
            emit Slashed(msg.sender, slashingFundsTracker.recipient, validator.withdrawalAddress, pubKey, minStake);
        }
    }

    /// @dev Internal function to set the minimum stake parameter.
    function _setMinStake(uint256 newMinStake) internal {
        require(newMinStake != 0, IVanillaRegistry.MinStakeMustBePositive());
        minStake = newMinStake;
        emit MinStakeSet(msg.sender, newMinStake);
    }

    /// @dev Internal function to set the slash oracle parameter.
    function _setSlashOracle(address newSlashOracle) internal {
        require(newSlashOracle != address(0), IVanillaRegistry.SlashOracleMustBeSet());
        slashOracle = newSlashOracle;
        emit SlashOracleSet(msg.sender, newSlashOracle);
    }

    /// @dev Internal function to set the slash receiver parameter.
    function _setSlashReceiver(address newSlashReceiver) internal {
        require(newSlashReceiver != address(0), IVanillaRegistry.SlashReceiverMustBeSet());
        slashingFundsTracker.recipient = newSlashReceiver;
        emit SlashReceiverSet(msg.sender, newSlashReceiver);
    }

    /// @dev Internal function to set the unstake period parameter.
    function _setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) internal {
        require(newUnstakePeriodBlocks != 0, IVanillaRegistry.UnstakePeriodMustBePositive());
        unstakePeriodBlocks = newUnstakePeriodBlocks;
        emit UnstakePeriodBlocksSet(msg.sender, newUnstakePeriodBlocks);
    }

    /// @dev Internal function to set the slashing payout period parameter in blocks.
    function _setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) internal {
        require(newSlashingPayoutPeriodBlocks != 0, IVanillaRegistry.SlashingPayoutPeriodMustBePositive());
        slashingFundsTracker.payoutPeriodBlocks = newSlashingPayoutPeriodBlocks;
        emit SlashingPayoutPeriodBlocksSet(msg.sender, newSlashingPayoutPeriodBlocks);
    }

    /// @dev Internal function to check if a validator is considered "opted-in" to mev-commit via this registry.
    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal view returns (bool) {
        return !_isUnstaking(valBLSPubKey) && stakedValidators[valBLSPubKey].balance >= minStake;
    }

    /// @dev Internal function to check if a validator is currently unstaking.
    function _isUnstaking(bytes calldata valBLSPubKey) internal view returns (bool) {
        return stakedValidators[valBLSPubKey].unstakeOccurrence.exists;
    }
}
