// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;


interface IRewardManager {
    event AutoClaimed(address indexed provider, address indexed toPay, uint256 amount);
    event PaymentStored(address indexed provider, address indexed toPay, uint256 amount);
    
    event ProposerNotFound(bytes indexed pubkey);

    event VanillaRegistrySet(address indexed oldVanillaRegistry, address indexed newVanillaRegistry);
    event MevCommitAVSSet(address indexed oldMevCommitAVS, address indexed newMevCommitAVS);
    event MevCommitMiddlewareSet(address indexed oldMevCommitMiddleware, address indexed newMevCommitMiddleware);
    event AutoClaimGasLimitSet(uint256 autoClaimGasLimit);

    event AutoClaimTransferFailed(address toPay);

    event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount);

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
}
