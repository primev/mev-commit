// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Errors} from "../../utils/Errors.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";
import {RewardManagerStorage} from "./RewardManagerStorage.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {VanillaRegistryStorage} from "../VanillaRegistryStorage.sol";
import {MevCommitAVSStorage} from "../avs/MevCommitAVSStorage.sol";
import {MevCommitMiddlewareStorage} from "../middleware/MevCommitMiddlewareStorage.sol";

contract RewardManager is IRewardManager, RewardManagerStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Receive function is disabled for this contract to prevent unintended interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to revert all calls, ensuring no unintended interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @dev Initializes the RewardManager contract.
    function initialize(
        address vanillaRegistry,
        address mevCommitAVS,
        address mevCommitMiddleware,
        uint256 autoClaimGasLimit,
        address owner
    ) external initializer {
        _setVanillaRegistry(vanillaRegistry);
        _setMevCommitAVS(mevCommitAVS);
        _setMevCommitMiddleware(mevCommitMiddleware);
        _setAutoClaimGasLimit(autoClaimGasLimit);
        __Pausable_init();
        __UUPSUpgradeable_init();
        __Ownable_init(owner);
    }

    /// @dev Enables the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Enables the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @dev Allows providers to pay opted-in proposers. 
    function payProposer(bytes calldata pubkey) external payable { // Intentionally don't allow pausing.
        address toPay = _findAddrToPay(pubkey);
        if (toPay == address(0)) {
            orphanedRewards[pubkey] += msg.value;
            emit OrphanedRewardsAccumulated(msg.sender, pubkey, msg.value);
            return;
        }
        if (overrideClaimAddresses[toPay] != address(0)) {
            toPay = overrideClaimAddresses[toPay];
        }
        if (autoClaim[toPay]) {
            (bool success, ) = payable(toPay).call{value: msg.value, gas: autoClaimGasLimit}("");
            if (!success) {
                autoClaim[toPay] = false; // AutoClaim disabled after first failed transfer
                unclaimedRewards[toPay] += msg.value;
                emit AutoClaimTransferFailed(toPay);
                return;
            }
            emit AutoClaimed(msg.sender, toPay, msg.value);
        } else {
            unclaimedRewards[toPay] += msg.value;
            emit PaymentStored(msg.sender, toPay, msg.value);
        }
    }

    /// @dev Enables auto-claim for a reward recipient.
    function enableAutoClaim() external whenNotPaused {
        autoClaim[msg.sender] = true;
        emit AutoClaimEnabled(msg.sender);
    }

    /// @dev Disables auto-claim for a reward recipient.
    function disableAutoClaim() external whenNotPaused {
        autoClaim[msg.sender] = false;
        emit AutoClaimDisabled(msg.sender);
    }

    /// @dev Allows the any reward recipient to delegate their rewards to another address.
    function overrideClaimAddress(address newClaimAddress) external whenNotPaused {
        require(newClaimAddress != address(0) && newClaimAddress != msg.sender, InvalidAddress());
        overrideClaimAddresses[msg.sender] = newClaimAddress;
        emit OverrideClaimAddressSet(msg.sender, newClaimAddress);
    }

    /// @dev Removes the override claim address for a reward recipient.
    function removeOverriddenClaimAddress() external whenNotPaused {
        overrideClaimAddresses[msg.sender] = address(0);
        emit OverrideClaimAddressRemoved(msg.sender);
    }

    /// @dev Allows a reward recipient to claim their rewards.
    function claimRewards() external whenNotPaused {
        uint256 amount = unclaimedRewards[msg.sender];
        require(amount > 0, NoRewardsToClaim());
        unclaimedRewards[msg.sender] = 0;
        (bool success, ) = payable(msg.sender).call{value: amount}("");
        if (!success) {
            revert RewardsClaimFailed();
        }
        emit RewardsClaimed(msg.sender, amount);
    }

    /// @dev Allows the owner to claim orphaned rewards to appropriate addresses.
    function claimOrphanedRewards(bytes[] calldata pubkeys, address toPay) external onlyOwner {
        uint256 totalAmount = 0;
        uint256 len = pubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubkey = pubkeys[i];
            uint256 amount = orphanedRewards[pubkey];
            require(amount > 0, NoOrphanedRewards());
            orphanedRewards[pubkey] = 0;
            totalAmount += amount;
        }
        (bool success, ) = payable(toPay).call{value: totalAmount}("");
        if (!success) {
            revert OrphanedRewardsClaimFailed();
        }
    }

    /// @dev Allows the owner to set the vanilla registry address.
    function setVanillaRegistry(address vanillaRegistry) external onlyOwner {
        _setVanillaRegistry(vanillaRegistry);
    }

    /// @dev Allows the owner to set the mev commit avs address.
    function setMevCommitAVS(address mevCommitAVS) external onlyOwner {
        _setMevCommitAVS(mevCommitAVS);
    }

    /// @dev Allows the owner to set the mev commit middleware address.
    function setMevCommitMiddleware(address mevCommitMiddleware) external onlyOwner {
        _setMevCommitMiddleware(mevCommitMiddleware);
    }

    /// @dev Allows the owner to set the auto claim gas limit.
    function setAutoClaimGasLimit(uint256 autoClaimGasLimit) external onlyOwner {
        _setAutoClaimGasLimit(autoClaimGasLimit);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _setVanillaRegistry(address vanillaRegistry) internal {
        require(vanillaRegistry != address(0), InvalidAddress());
        _vanillaRegistry = VanillaRegistryStorage(vanillaRegistry);
        emit VanillaRegistrySet(address(vanillaRegistry), vanillaRegistry);
    }

    function _setMevCommitAVS(address mevCommitAVS) internal {
        require(mevCommitAVS != address(0), InvalidAddress());
        _mevCommitAVS = MevCommitAVSStorage(mevCommitAVS);
        emit MevCommitAVSSet(address(mevCommitAVS), mevCommitAVS);
    }

    function _setMevCommitMiddleware(address mevCommitMiddleware) internal {
        require(mevCommitMiddleware != address(0), InvalidAddress());
        _mevCommitMiddleware = MevCommitMiddlewareStorage(mevCommitMiddleware);
        emit MevCommitMiddlewareSet(address(mevCommitMiddleware), mevCommitMiddleware);
    }

    function _setAutoClaimGasLimit(uint256 autoClaimGasLimit) internal {
        require(autoClaimGasLimit > 0, InvalidAutoClaimGasLimit());
        autoClaimGasLimit = autoClaimGasLimit;
        emit AutoClaimGasLimitSet(autoClaimGasLimit);
    }

    function _findAddrToPay(bytes calldata pubkey) internal view returns (address) {
        (,address operatorAddr,,) = _mevCommitMiddleware.validatorRecords(pubkey);
        if (operatorAddr != address(0)) {
            return operatorAddr;
        }
        (,address vanillaWithdrawalAddr,,) = _vanillaRegistry.stakedValidators(pubkey);
        if (vanillaWithdrawalAddr != address(0)) {
            return vanillaWithdrawalAddr;
        }
        (,address podOwner,,) = _mevCommitAVS.validatorRegistrations(pubkey);
        if (podOwner != address(0)) {
            return podOwner;
        }
        return address(0);
    }
}
