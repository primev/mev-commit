// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {PreconfManager} from "./PreconfManager.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {ProviderRegistryStorage} from "./ProviderRegistryStorage.sol";
import {FeePayout} from "../utils/FeePayout.sol";
import {Errors} from "../utils/Errors.sol";

/// @title Provider Registry
/// @author Kartik Chopra
/// @notice This contract is for provider registry and staking.
contract ProviderRegistry is
    IProviderRegistry,
    ProviderRegistryStorage,
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable,
    PausableUpgradeable
{
    /**
     * @dev Modifier to restrict a function to only be callable by the preconf manager.
     */
    modifier onlyPreconfManager() {
        require(msg.sender == preconfManager, NotPreconfContract(msg.sender, preconfManager));
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
        uint256 _feePercent,
        address _owner,
        uint256 _withdrawalDelay,
        uint256 _penaltyFeePayoutPeriodBlocks
    ) external initializer {
        FeePayout.init(penaltyFeeTracker, _penaltyFeeRecipient, _penaltyFeePayoutPeriodBlocks);
        minStake = _minStake;
        feePercent = _feePercent;
        withdrawalDelay = _withdrawalDelay;
        __ReentrancyGuard_init();
        __Ownable_init(_owner);
        __Pausable_init();
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
        revert Errors.InvalidReceive();
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /**
     * @dev Stake more funds into the provider's stake.
     */
    function stake() external payable whenNotPaused {
        _stake(msg.sender);
    }

    /// @dev Delegate stake to a provider.
    function delegateStake(address provider) external payable whenNotPaused {
        _stake(provider);
    }

    /**
     * @dev Slash funds from the provider and send the slashed amount to the bidder.
     * @dev reenterancy not necessary but still putting here for precaution.
     * @dev Note we slash the funds taking into account the residual bid percent after decay.
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
    ) external nonReentrant onlyPreconfManager whenNotPaused {
        uint256 residualAmt = (amt * residualBidPercentAfterDecay) / ONE_HUNDRED_PERCENT;
        uint256 penaltyFee = (residualAmt * feePercent) / ONE_HUNDRED_PERCENT;
        uint256 providerStake = providerStakes[provider];

        if (providerStake < residualAmt + penaltyFee) {
            emit InsufficientFundsToSlash(provider, providerStake, residualAmt, penaltyFee);
            if (providerStake < residualAmt) {
                residualAmt = providerStake;
            }
            penaltyFee = providerStake - residualAmt;
        }
        providerStakes[provider] -= residualAmt + penaltyFee;

        penaltyFeeTracker.accumulatedAmount += penaltyFee;
        if (FeePayout.isPayoutDue(penaltyFeeTracker)) {
            FeePayout.transferToRecipient(penaltyFeeTracker);
        }

        if (!payable(bidder).send(amt)) {
            emit TransferToBidderFailed(bidder, amt);
            bidderSlashedAmount[bidder] += amt;
        }

        emit FundsSlashed(provider, residualAmt + penaltyFee);
    }

    /**
     * @dev Sets the minimum stake required for registration. Can only be called by the owner.
     * @param _minStake The new minimum stake amount.
     */
    function setMinStake(uint256 _minStake) external onlyOwner {
        minStake = _minStake;
        emit MinStakeUpdated(_minStake);
    }

    /**
     * @dev Sets the pre-confirmations contract address. Can only be called by the owner.
     * @param contractAddress The address of the pre-confirmations contract.
     */
    function setPreconfManager(address contractAddress) external onlyOwner {
        preconfManager = contractAddress;
        emit PreconfManagerUpdated(contractAddress);
    }

    /**
     * @notice Sets the new fee percent
     * @dev onlyOwner restriction
     * @param newFeePercent this is the new fee percent
     */
    function setNewFeePercent(uint256 newFeePercent) external onlyOwner {
        feePercent = newFeePercent;
        emit FeePercentUpdated(newFeePercent);
    }

    /// @dev Sets the withdrawal delay. Can only be called by the owner.
    /// @param _withdrawalDelay The new withdrawal delay in milliseconds 
    /// as mev-commit chain is running with milliseconds.
    function setWithdrawalDelay(uint256 _withdrawalDelay) external onlyOwner {
        withdrawalDelay = _withdrawalDelay;
        emit WithdrawalDelayUpdated(_withdrawalDelay);
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

    /// @dev Sets the fee payout period in blocks
    /// @param _feePayoutPeriodBlocks The new fee payout period in blocks
    function setFeePayoutPeriodBlocks(uint256 _feePayoutPeriodBlocks) external onlyOwner {
        penaltyFeeTracker.payoutPeriodBlocks = _feePayoutPeriodBlocks;
        emit FeePayoutPeriodBlocksUpdated(_feePayoutPeriodBlocks);
    }

    function overrideAddBLSKey(address provider, bytes calldata blsPublicKey) external onlyOwner {
        require(providerRegistered[provider], ProviderNotRegistered(provider));
        eoaToBlsPubkeys[provider].push(blsPublicKey);
        blockBuilderBLSKeyToAddress[blsPublicKey] = provider;
    }

    /// @dev Requests unstake of the staked amount.
    function unstake() external whenNotPaused {
        require(providerStakes[msg.sender] != 0, NoStakeToWithdraw(msg.sender));
        require(withdrawalRequests[msg.sender] == 0, UnstakeRequestExists(msg.sender));
        withdrawalRequests[msg.sender] = block.timestamp;
        emit Unstake(msg.sender, block.timestamp);
    }

    /// @dev Completes the withdrawal of the staked amount.
    function withdraw() external nonReentrant whenNotPaused {
        require(withdrawalRequests[msg.sender] != 0, NoUnstakeRequest(msg.sender));
        require(block.timestamp >= withdrawalRequests[msg.sender] + withdrawalDelay,
            DelayNotPassed(withdrawalRequests[msg.sender], withdrawalDelay, block.timestamp));

        uint256 providerStake = providerStakes[msg.sender];
        providerStakes[msg.sender] = 0;
        providerRegistered[msg.sender] = false;
        withdrawalRequests[msg.sender] = 0;
        require(preconfManager != address(0), PreconfManagerNotSet());

        uint256 providerPendingCommitmentsCount = PreconfManager(
            payable(preconfManager)
        ).commitmentsCount(msg.sender);

        require(providerPendingCommitmentsCount == 0, ProviderCommitmentsPending(msg.sender, providerPendingCommitmentsCount));

        (bool success, ) = msg.sender.call{value: providerStake}("");
        require(success, StakeTransferFailed(msg.sender, providerStake));

        emit Withdraw(msg.sender, providerStake);
    }

    /**
     * @dev Allows the bidder to withdraw the slashed amount.
     */
    function withdrawSlashedAmount() external nonReentrant whenNotPaused() {
        require(bidderSlashedAmount[msg.sender] != 0, BidderAmountIsZero(msg.sender));
        uint256 amount = bidderSlashedAmount[msg.sender];
        bidderSlashedAmount[msg.sender] = 0;
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, BidderWithdrawalTransferFailed(msg.sender, amount));

        emit BidderWithdrawSlashedAmount(msg.sender, amount);
    }

    /**
     * @dev Manually withdraws accumulated penalty fees to the recipient
     * to cover the edge case that oracle doesn't slash/reward, and funds still need to be withdrawn.
     */
    function manuallyWithdrawPenaltyFee() external onlyOwner {
        FeePayout.transferToRecipient(penaltyFeeTracker);
    }

    /// @dev Allows the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Allows the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Get provider staked amount.
     * @param provider The address of the provider.
     * @return The staked amount for the provider.
     */
    function getProviderStake(address provider) external view returns (uint256) {
        return providerStakes[provider];
    }

        /// @dev Returns the BLS public keys corresponding to a provider's staked EOA address.
    function getBLSKeys(address provider) external view returns (bytes[] memory) {
        return eoaToBlsPubkeys[provider];
    }

    /// @dev Returns the EOA address corresponding to a provider's BLS public key.
    function getEoaFromBLSKey(bytes calldata blsKey) external view returns (address) {
        return blockBuilderBLSKeyToAddress[blsKey];
    }

    /// @return penaltyFee amount not yet transferred to recipient
    function getAccumulatedPenaltyFee() external view returns (uint256) {
        return penaltyFeeTracker.accumulatedAmount;
    }

    /**
     * @dev Register and stake function for providers.
     * The validity of this key must be verified manually off-chain.
     */
    function registerAndStake() public payable whenNotPaused {
        _registerAndStake(msg.sender);
    }

    /**
    * @dev Register and stake on behalf of a provider.
    * @param provider Address of the provider.
    */
    function delegateRegisterAndStake(address provider) public payable whenNotPaused onlyOwner {
        _registerAndStake(provider);
    }

    /// @dev Ensure the provider's balance is greater than minStake and no pending withdrawal
    function isProviderValid(address provider) public view {
        uint256 providerStake = providerStakes[provider];
        require(providerStake >= minStake, InsufficientStake(providerStake, minStake));
        require(withdrawalRequests[provider] == 0, PendingWithdrawalRequest(provider));
    }

    function _stake(address provider) internal {
        require(providerRegistered[provider], ProviderNotRegistered(provider));
        require(withdrawalRequests[provider] == 0, PendingWithdrawalRequest(provider));
        providerStakes[provider] += msg.value;
        emit FundsDeposited(provider, msg.value);
    }

    /**
     * @dev Adds a verified BLS key to the provider's account.
     * @param blsPublicKey The BLS public key to be added.
     * @param signature The signature (96 bytes) used for verification.
     */
    function addVerifiedBLSKey(
        bytes calldata blsPublicKey,
        bytes calldata signature
    ) external {
        address provider = msg.sender;

        require(providerRegistered[provider], "Provider not registered");
        require(blsPublicKey.length == 48, "Public key must be 48 bytes");
        require(signature.length == 96, "Signature must be 96 bytes");


        bytes32 message = keccak256(abi.encodePacked(provider));

        // Verify the BLS signature
        bool isValid = verifySignature(blsPublicKey, message, signature);
        require(isValid, "Invalid BLS signature");

        // Add the BLS public key to the provider's account
        eoaToBlsPubkeys[provider].push(blsPublicKey);
        blockBuilderBLSKeyToAddress[blsPublicKey] = provider;

        emit BLSKeyAdded(provider, blsPublicKey);
    }

    /**
     * @dev Verifies a BLS signature using the precompile
     * @param pubKey The public key (48 bytes G1 point)
     * @param message The message hash (32 bytes)
     * @param signature The signature (96 bytes G2 point)
     * @return success True if verification succeeded
     */
    function verifySignature(
        bytes calldata pubKey,
        bytes32 message,
        bytes calldata signature
    ) public view returns (bool) {
        // Input validation
        require(pubKey.length == 48, "Public key must be 48 bytes");
        require(signature.length == 96, "Signature must be 96 bytes");

        // Concatenate inputs in required format:
        // [pubkey (48 bytes) | message (32 bytes) | signature (96 bytes)]
        bytes memory input = bytes.concat(
            pubKey,
            message,
            signature
        );

        // Call precompile
        (bool success, bytes memory result) = address(0xf0).staticcall(input);
        
        // Check if call was successful
        if (!success) {
            return false;
        }

        // If we got a result back and it's not empty, verification succeeded
        return result.length > 0;
    }

    event BLSKeyAdded(address indexed provider, bytes blsPublicKey);
 
    function _registerAndStake(address provider) internal {
        require(!providerRegistered[provider], ProviderAlreadyRegistered(provider));
        require(msg.value >= minStake, InsufficientStake(msg.value, minStake));

        providerStakes[provider] = msg.value;
        providerRegistered[provider] = true;
        emit ProviderRegistered(provider, msg.value);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
