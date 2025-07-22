// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;


interface IRewardManager {
    event AutoClaimed(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount);
    event PaymentStored(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount);
    event ProposerNotFound(bytes indexed pubkey);
    event VanillaRegistrySet(address indexed newVanillaRegistry);
    event MevCommitAVSSet(address indexed newMevCommitAVS);
    event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware);
    event AutoClaimGasLimitSet(uint256 autoClaimGasLimit);
    event AutoClaimTransferFailed(address indexed provider, address indexed receiver, address indexed toPay);
    event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount);
    event OrphanedRewardsClaimed(address indexed toPay, uint256 amount);
    event RemovedFromAutoClaimBlacklist(address indexed addr);
    event OverrideAddressSet(address indexed receiver, address indexed overrideAddress);
    event OverrideAddressRemoved(address indexed receiver);
    event AutoClaimEnabled(address indexed receiver);
    event AutoClaimDisabled(address indexed receiver);
    event RewardsClaimed(address indexed msgSender, uint256 amount);
    event RewardsMigrated(address indexed from, address indexed to, uint256 amount);
    event NoRewards(address addr);
    error NoEthPayable();
    error InvalidAddress();
    error NoOverriddenAddressToRemove();
    error InvalidAutoClaimGasLimit();
    error RewardsClaimFailed();
    error NoOrphanedRewards();
    error OrphanedRewardsClaimFailed();
    error InvalidBLSPubKeyLength(uint256 expectedLength, uint256 actualLength);

    /// @dev Allows providers to pay the opted-in proposer for a block. 
    /// @notice It is assumed the validator pubkey being paid is opted-in to mev-commit.
    /// Otherwise the rewards are accumulated as "orphaned" and must be handled by the owner.
    function payProposer(bytes calldata pubkey) external payable;

    /// @dev Enables auto-claim for a receiver address.
    /// @param claimExistingRewards If true, existing rewards will be claimed atomically before enabling auto-claim.
    function enableAutoClaim(bool claimExistingRewards) external;

    /// @dev Disables auto-claim for a reward recipient.
    function disableAutoClaim() external;

    /// @dev Allows any receiver address to set an override address for their rewards.
    /// @param migrateExistingRewards If true, existing msg.sender rewards will be migrated atomically to the new claim address.
    /// @notice Onus is on the calling address to ensure the override address does not revert upon receiving eth transfers.
    function overrideReceiver(address overrideAddress, bool migrateExistingRewards) external;

    /// @dev Removes the override address for a receiver.
    function removeOverrideAddress() external;

    /// @dev Allows a reward recipient to claim their rewards.
    function claimRewards() external;

    /// @dev Allows the owner to claim orphaned rewards to appropriate addresses.
    function claimOrphanedRewards(bytes[] calldata pubkeys, address toPay) external;

    /// @dev Allows the owner to remove an address from the auto claim blacklist.
    function removeFromAutoClaimBlacklist(address addr) external;

    /// @dev Allows the owner to set the vanilla registry address.
    function setVanillaRegistry(address vanillaRegistry) external;

    /// @dev Allows the owner to set the mev commit avs address.
    function setMevCommitAVS(address mevCommitAVS) external;

    /// @dev Allows the owner to set the mev commit middleware address.
    function setMevCommitMiddleware(address mevCommitMiddleware) external;

    /// @dev Allows the owner to set the auto claim gas limit.
    function setAutoClaimGasLimit(uint256 autoClaimGasLimit) external;
}
