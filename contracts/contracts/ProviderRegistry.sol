// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import {PreConfCommitmentStore} from "./PreConfCommitmentStore.sol";
import {IProviderRegistry} from "./interfaces/IProviderRegistry.sol";

/// @title Provider Registry
/// @author Kartik Chopra
/// @notice This contract is for provider registry and staking.
contract ProviderRegistry is
    IProviderRegistry,
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    ReentrancyGuardUpgradeable
{
    /// @dev For improved precision
    uint256 public constant PRECISION = 10 ** 25;
    uint256 public constant PERCENT = 100 * PRECISION;

    /// @dev Minimum stake required for registration
    uint256 public minStake;

    /// @dev Amount assigned to feeRecipient
    uint256 public feeRecipientAmount;

    /// @dev Amount assigned to protocol fee
    uint256 public protocolFeeAmount;

    /// @dev Address of the pre-confirmations contract
    address public preConfirmationsContract;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Fee recipient
    address public feeRecipient;

    /// @dev Configurable withdrawal delay in milliseconds
    uint256 public withdrawalDelay;

    /// @dev Mapping from provider address to whether they are registered or not
    mapping(address => bool) public providerRegistered;
    /// @dev Mapping from a provider's EOA address to their BLS public key
    mapping(address => bytes) public eoaToBlsPubkey;
    /// @dev Mapping from provider addresses to their staked amount
    mapping(address => uint256) public providerStakes;
    /// @dev Mapping of provider to withdrawal request timestamp
    mapping(address => uint256) public withdrawalRequests;

    /// @dev Event emitted when a provider is registered
    event ProviderRegistered(address indexed provider, uint256 stakedAmount, bytes blsPublicKey);
    /// @dev Event emitted when funds are deposited
    event FundsDeposited(address indexed provider, uint256 amount);
    /// @dev Event emitted when funds are slashed
    event FundsSlashed(address indexed provider, uint256 amount);
    /// @dev Event emitted when withdrawal is requested
    event Unstake(address indexed provider, uint256 timestamp);
    /// @dev Event emitted when withdrawal is completed
    event Withdraw(address indexed provider, uint256 amount);
    /// @dev Event emitted when the withdrawal delay is updated
    event WithdrawalDelayUpdated(uint256 newWithdrawalDelay);

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
     * @param _feeRecipient The address that receives fee
     * @param _feePercent The fee percentage for protocol
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _withdrawalDelay The withdrawal delay in seconds.
     */
    function initialize(
        uint256 _minStake,
        address _feeRecipient,
        uint16 _feePercent,
        address _owner,
        uint256 _withdrawalDelay
    ) external initializer {
        minStake = _minStake;
        feeRecipient = _feeRecipient;
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
        require(providerRegistered[msg.sender], "Provider not registered");
        providerStakes[msg.sender] += msg.value;
        emit FundsDeposited(msg.sender, msg.value);
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
        require(providerStakes[provider] >= residualAmt, "Insufficient funds to slash");
        providerStakes[provider] -= residualAmt;

        uint256 feeAmt = (residualAmt * uint256(feePercent) * PRECISION) / PERCENT;
        uint256 amtMinusFee = residualAmt - feeAmt;

        if (feeRecipient != address(0)) {
            feeRecipientAmount += feeAmt;
        } else {
            protocolFeeAmount += feeAmt;
        }

        (bool success, ) = payable(bidder).call{value: amtMinusFee}("");
        require(success, "Transfer to bidder failed");

        emit FundsSlashed(provider, amtMinusFee);
    }

    /**
     * @notice Sets the new fee recipient
     * @dev onlyOwner restriction
     * @param newFeeRecipient The address to transfer the slashed funds to.
     */
    function setNewFeeRecipient(address newFeeRecipient) external onlyOwner {
        feeRecipient = newFeeRecipient;
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

    /**
     * @dev Reward funds to the fee receipt.
     */
    function withdrawFeeRecipientAmount() external nonReentrant {
        feeRecipientAmount = 0;
        (bool successFee, ) = feeRecipient.call{value: feeRecipientAmount}("");
        require(successFee, "fee recipient transfer failed");
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

        uint256 providerPendingCommitmentsCount = PreConfCommitmentStore(
            payable(preConfirmationsContract)
        ).commitmentsCount(msg.sender);

        require(providerPendingCommitmentsCount == 0, "provider commitments are pending");

        (bool success, ) = msg.sender.call{value: providerStake}("");
        require(success, "stake transfer failed");

        emit Withdraw(msg.sender, providerStake);
    }

    /**
     * @dev Withdraw protocol fee.
     * @param treasuryAddress The address of the treasury.
     */
    function withdrawProtocolFee(
        address payable treasuryAddress
    ) external onlyOwner nonReentrant {
        uint256 _protocolFeeAmount = protocolFeeAmount;
        protocolFeeAmount = 0;
        require(_protocolFeeAmount != 0, "insufficient protocol fee amount");
        (bool success, ) = treasuryAddress.call{value: _protocolFeeAmount}("");
        require(success, "transfer failed");
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

    /**
     * @dev Register and stake function for providers.
     * @param blsPublicKey The BLS public key of the provider.
     * The validity of this key must be verified manually off-chain.
     */
    function registerAndStake(bytes calldata blsPublicKey) public payable {
        require(!providerRegistered[msg.sender], "Provider already registered");
        require(msg.value >= minStake, "Insufficient stake");
        require(blsPublicKey.length == 48, "Invalid BLS public key length");
        
        eoaToBlsPubkey[msg.sender] = blsPublicKey;
        providerStakes[msg.sender] = msg.value;
        providerRegistered[msg.sender] = true;
        emit ProviderRegistered(msg.sender, msg.value, blsPublicKey);
    }

    /// @dev Ensure the provider's balance is greater than minStake and no pending withdrawal
    function isProviderValid(address provider) public view {
        require(providerStakes[provider] >= minStake, "Insufficient stake");
        require(withdrawalRequests[provider] == 0, "Pending withdrawal request");
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
