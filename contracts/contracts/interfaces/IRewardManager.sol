// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;


interface IRewardManager {
    event AutoClaimed(address indexed provider, address indexed toPay, uint256 amount);
    event PaymentStored(address indexed provider, address indexed toPay, uint256 amount);
    event ProposerNotFound(bytes indexed pubkey);
    event VanillaRegistrySet(address indexed newVanillaRegistry);
    event MevCommitAVSSet(address indexed newMevCommitAVS);
    event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware);
    event AutoClaimGasLimitSet(uint256 autoClaimGasLimit);
    event AutoClaimTransferFailed(address toPay);
    event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount);
    event OrphanedRewardsClaimed(address indexed toPay, uint256 amount);
    event RemovedFromAutoClaimBlacklist(address indexed addr);
    event OverrideClaimAddressSet(address indexed provider, address indexed newClaimAddress);
    event OverrideClaimAddressRemoved(address indexed provider);
    event AutoClaimEnabled(address indexed caller);
    event AutoClaimDisabled(address indexed caller);
    event RewardsClaimed(address indexed claimer, uint256 amount);
    error InvalidAddress();
    error InvalidAutoClaimGasLimit();
    error RewardsClaimFailed();
    error NoRewardsToClaim();
    error NoOrphanedRewards();
    error OrphanedRewardsClaimFailed();
    /// @dev Allows providers to pay opted-in proposers.
    function payProposer(bytes calldata pubkey) external payable;
    /// @dev Enables auto-claim for a reward recipient.
    function enableAutoClaim() external;
    /// @dev Disables auto-claim for a reward recipient.
    function disableAutoClaim() external;
    /// @dev Allows a reward recipient to delegate their rewards to another address.
    function overrideClaimAddress(address newClaimAddress) external;
    /// @dev Removes the override claim address for a reward recipient.
    function removeOverriddenClaimAddress() external;
    /// @dev Allows a reward recipient to claim their rewards.
    function claimRewards() external;
    /// @dev Allows the owner to claim orphaned rewards to appropriate addresses.
    function claimOrphanedRewards(bytes[] calldata pubkeys, address toPay) external;
    /// @dev Allows the owner to set the vanilla registry address.
    function setVanillaRegistry(address vanillaRegistry) external;
    /// @dev Allows the owner to set the mev commit avs address.
    function setMevCommitAVS(address mevCommitAVS) external;
    /// @dev Allows the owner to set the mev commit middleware address.
    function setMevCommitMiddleware(address mevCommitMiddleware) external;
    /// @dev Allows the owner to set the auto claim gas limit.
    function setAutoClaimGasLimit(uint256 autoClaimGasLimit) external;
}
