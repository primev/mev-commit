// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {PreconfManager} from "./PreconfManager.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {ProviderRegistryStorage} from "./ProviderRegistryStorage.sol";
import {FeePayout} from "../utils/FeePayout.sol";

/// @title Provider Registry
/// @author Kartik Chopra
/// @notice This contract is for provider registry and staking.
contract ProviderRegistry is
    IProviderRegistry,
    ProviderRegistryStorage,
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable
{
    /**
     * @dev Modifier to restrict a function to only be callable by the pre-confirmations contract.
     */
    modifier onlyPreConfirmationEngine() {
        require(
            msg.sender == preConfirmationsContract,
            "sender is not preconf contract"
        );
        _;
    }

    /**
     * @dev Initializes the contract with a minimum stake requirement.
     * @param _minStake The minimum stake required for provider registration.
     * @param _penaltyFeeRecipient The address that accumulates penalty fees
     * @param _feePercent The fee percentage for penalty
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _withdrawalDelay The withdrawal delay in milliseconds.
     * @param _penaltyFeePayoutPeriodBlocks The min number of blocks between penalty fee payouts
     */
    function initialize(
        uint256 _minStake,
        address _penaltyFeeRecipient,
        uint16 _feePercent,
        address _owner,
        uint256 _withdrawalDelay,
        uint256 _penaltyFeePayoutPeriodBlocks
    ) external initializer {
        FeePayout.init(penaltyFeeTracker, _penaltyFeeRecipient, _penaltyFeePayoutPeriodBlocks);
        minStake = _minStake;
        feePercent = _feePercent;
        withdrawalDelay = _withdrawalDelay;
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Receive function is disabled for this contract to prevent unintended interactions.
     * Should be removed from here in case the registerAndStake function becomes more complex
     */
    receive() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Sets the pre-confirmations contract address. Can only be called by the owner.
     * @param contractAddress The address of the pre-confirmations contract.
     */
    function setPreconfirmationsContract(
        address contractAddress
    ) external onlyOwner {
        require(
            preConfirmationsContract == address(0),
            "preconf contract already set"
        );
        preConfirmationsContract = contractAddress;
    }

    /**
     * @dev Stake more funds into the provider's stake.
     */
    function stake() external payable {
        _stake(msg.sender);
    }

    /// @dev Delegate stake to a provider.
    function delegateStake(address provider) external payable {
        _stake(provider);
    }

    /**
     * @dev Slash funds from the provider and send the slashed amount to the bidder.
     * @dev reenterancy not necessary but still putting here for precaution.
     * @dev Note we slash the funds irrespective of decay, this is to prevent timing games.
     * @param amt The amount to slash from the provider's stake.
     * @param provider The address of the provider.
     * @param bidder The address to transfer the slashed funds to.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function slash(
        uint256 amt,
        address provider,
        address payable bidder,
        uint256 residualBidPercentAfterDecay
    ) external nonReentrant onlyPreConfirmationEngine {
        uint256 residualAmt = (amt * residualBidPercentAfterDecay * PRECISION) / PERCENT;
        uint256 penaltyFee = (residualAmt * uint256(feePercent) * PRECISION) / PERCENT;
        require(providerStakes[provider] >= residualAmt + penaltyFee, "Insufficient funds to slash");
        providerStakes[provider] -= residualAmt + penaltyFee;

        penaltyFeeTracker.accumulatedAmount += penaltyFee;
        if (FeePayout.isPayoutDue(penaltyFeeTracker)) {
            FeePayout.transferToRecipient(penaltyFeeTracker);
        }

        (bool success, ) = payable(bidder).call{value: residualAmt}("");
        require(success, "Transfer to bidder failed");

        emit FundsSlashed(provider, residualAmt + penaltyFee);
    }

    /**
     * @notice Sets a new penalty fee recipient
     * @dev onlyOwner restriction
     * @param newFeeRecipient The address of the new penalty fee recipient
     */
    function setNewPenaltyFeeRecipient(address newFeeRecipient) external onlyOwner {
        penaltyFeeTracker.recipient = newFeeRecipient;
        emit PenaltyFeeRecipientUpdated(newFeeRecipient);
    }

    /**
     * @notice Sets the new fee recipient
     * @dev onlyOwner restriction
     * @param newFeePercent this is the new fee percent
     */
    function setNewFeePercent(uint16 newFeePercent) external onlyOwner {
        feePercent = newFeePercent;
    }

    /// @dev Sets the withdrawal delay. Can only be called by the owner.
    /// @param _withdrawalDelay The new withdrawal delay in milliseconds 
    /// as mev-commit chain is running with milliseconds.
    function setWithdrawalDelay(uint256 _withdrawalDelay) external onlyOwner {
        withdrawalDelay = _withdrawalDelay;
        emit WithdrawalDelayUpdated(_withdrawalDelay);
    }

    /// @dev Sets the fee payout period in blocks
    /// @param _feePayoutPeriodBlocks The new fee payout period in blocks
    function setFeePayoutPeriodBlocks(uint256 _feePayoutPeriodBlocks) external onlyOwner {
        penaltyFeeTracker.payoutPeriodBlocks = _feePayoutPeriodBlocks;
        emit FeePayoutPeriodBlocksUpdated(_feePayoutPeriodBlocks);
    }

    /// @dev Requests unstake of the staked amount.
    function unstake() external {
        require(providerStakes[msg.sender] != 0, "No stake to withdraw");
        require(withdrawalRequests[msg.sender] == 0, "Unstake request exists");
        withdrawalRequests[msg.sender] = block.timestamp;
        emit Unstake(msg.sender, block.timestamp);
    }

    /// @dev Completes the withdrawal of the staked amount.
    function withdraw() external nonReentrant {
        require(withdrawalRequests[msg.sender] != 0, "No unstake request");
        require(block.timestamp >= withdrawalRequests[msg.sender] + withdrawalDelay, "Delay has not passed");

        uint256 providerStake = providerStakes[msg.sender];
        providerStakes[msg.sender] = 0;
        withdrawalRequests[msg.sender] = 0;
        require(providerStake != 0, "Provider Staked Amount is zero");
        require(preConfirmationsContract != address(0), "preconf contract not set");

        uint256 providerPendingCommitmentsCount = PreconfManager(
            payable(preConfirmationsContract)
        ).commitmentsCount(msg.sender);

        require(providerPendingCommitmentsCount == 0, "provider commitments are pending");

        (bool success, ) = msg.sender.call{value: providerStake}("");
        require(success, "stake transfer failed");

        emit Withdraw(msg.sender, providerStake);
    }

    /**
     * @dev Manually withdraws accumulated penalty fees to the recipient
     * to cover the edge case that oracle doesn't slash/reward, and funds still need to be withdrawn.
     */
    function manuallyWithdrawPenaltyFee() external onlyOwner {
        FeePayout.transferToRecipient(penaltyFeeTracker);
    }

    /**
     * @dev Get provider staked amount.
     * @param provider The address of the provider.
     * @return The staked amount for the provider.
     */
    function getProviderStake(address provider) external view returns (uint256) {
        return providerStakes[provider];
    }

    /// @dev Returns the BLS public key corresponding to a provider's staked EOA address.
    function getBLSKey(address provider) external view returns (bytes memory) {
        return eoaToBlsPubkey[provider];
    }

    /// @return penaltyFee amount not yet transferred to recipient
    function getAccumulatedPenaltyFee() external view returns (uint256) {
        return penaltyFeeTracker.accumulatedAmount;
    }

    /**
     * @dev Register and stake function for providers.
     * @param blsPublicKey The BLS public key of the provider.
     * The validity of this key must be verified manually off-chain.
     */
    function registerAndStake(bytes calldata blsPublicKey) public payable {
        _registerAndStake(msg.sender, blsPublicKey);
    }

    /**
     * @dev Register and stake on behalf of a provider.
     * @param provider Address of the provider.
     * @param blsPublicKey BLS public key of the provider.
     */
    function delegateRegisterAndStake(address provider, bytes calldata blsPublicKey) public payable {
        _registerAndStake(provider, blsPublicKey);
    }

    /// @dev Ensure the provider's balance is greater than minStake and no pending withdrawal
    function isProviderValid(address provider) public view {
        require(providerStakes[provider] >= minStake, "Insufficient stake");
        require(withdrawalRequests[provider] == 0, "Pending withdrawal request");
    }

    function _stake(address provider) internal {
        require(providerRegistered[provider], "Provider not registered");
        providerStakes[provider] += msg.value;
        emit FundsDeposited(provider, msg.value);
    }

    function _registerAndStake(address provider, bytes calldata blsPublicKey) internal {
        require(!providerRegistered[provider], "Provider already registered");
        require(msg.value >= minStake, "Insufficient stake");
        require(blsPublicKey.length == 48, "Invalid BLS public key length");
        
        eoaToBlsPubkey[provider] = blsPublicKey;
        providerStakes[provider] = msg.value;
        providerRegistered[provider] = true;
        emit ProviderRegistered(provider, msg.value, blsPublicKey);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
