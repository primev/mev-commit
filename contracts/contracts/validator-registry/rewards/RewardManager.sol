// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Errors} from "../../utils/Errors.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";
import {RewardManagerStorage} from "./RewardManagerStorage.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {VanillaRegistryStorage} from "../VanillaRegistryStorage.sol";
import {MevCommitAVSStorage} from "../avs/MevCommitAVSStorage.sol";
import {MevCommitMiddlewareStorage} from "../middleware/MevCommitMiddlewareStorage.sol";

contract RewardManager is IRewardManager, RewardManagerStorage,
    Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    modifier onlyValidBLSPubKey(bytes calldata pubkey) {
        require(pubkey.length == 48, InvalidBLSPubKeyLength(48, pubkey.length));
        _;
    }

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
        __Ownable_init(owner);
        __ReentrancyGuard_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    /// @dev Enables the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Enables the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @dev Allows providers to pay the opted-in proposer for a block. 
    /// @notice It is assumed the validator pubkey being paid is opted-in to mev-commit.
    /// Otherwise the rewards are accumulated as "orphaned" and must be handled by the owner.
    function payProposer(bytes calldata pubkey) external payable onlyValidBLSPubKey(pubkey) nonReentrant { // Intentionally don't allow pausing.
        require(msg.value != 0, NoEthPayable());
        address receiver = _findReceiver(pubkey);
        if (receiver == address(0)) {
            orphanedRewards[pubkey] += msg.value;
            emit OrphanedRewardsAccumulated(msg.sender, pubkey, msg.value);
            return;
        }
        address toPay = receiver;
        address overrideAddr = overrideAddresses[receiver];
        if (overrideAddr != address(0)) {
            toPay = overrideAddr;
        }
        if (autoClaim[receiver] && !autoClaimBlacklist[receiver]) {
            (bool success, ) = payable(toPay).call{value: msg.value, gas: autoClaimGasLimit}("");
            if (!success) {
                autoClaim[receiver] = false;
                autoClaimBlacklist[receiver] = true;
                unclaimedRewards[toPay] += msg.value;
                emit AutoClaimTransferFailed(msg.sender, receiver, toPay);
                return;
            }
            emit AutoClaimed(msg.sender, receiver, toPay, msg.value);
        } else {
            unclaimedRewards[toPay] += msg.value;
            emit PaymentStored(msg.sender, receiver, toPay, msg.value);
        }
    }

    /// @dev Enables auto-claim for a receiver address.
    /// @param claimExistingRewards If true, existing rewards will be claimed atomically before enabling auto-claim.
    function enableAutoClaim(bool claimExistingRewards) external whenNotPaused nonReentrant {
        if (claimExistingRewards) { _claimRewards(); }
        autoClaim[msg.sender] = true;
        emit AutoClaimEnabled(msg.sender);
    }

    /// @dev Disables auto-claim for a receiver address.
    function disableAutoClaim() external whenNotPaused {
        autoClaim[msg.sender] = false;
        emit AutoClaimDisabled(msg.sender);
    }

    /// @dev Allows any receiver address to set an override address for their rewards.
    /// @param migrateExistingRewards If true, existing msg.sender rewards will be migrated atomically to the new claim address.
    function overrideReceiver(address overrideAddress, bool migrateExistingRewards) external whenNotPaused {
        if (migrateExistingRewards) { _migrateRewards(msg.sender, overrideAddress); }
        require(overrideAddress != address(0) && overrideAddress != msg.sender, InvalidAddress());
        overrideAddresses[msg.sender] = overrideAddress;
        emit OverrideAddressSet(msg.sender, overrideAddress);
    }

    /// @dev Removes the override address for a receiver.
    /// @param migrateExistingRewards If true, existing rewards for the overridden address will be migrated atomically to the msg.sender.
    function removeOverrideAddress(bool migrateExistingRewards) external whenNotPaused {
        address toBeRemoved = overrideAddresses[msg.sender];
        require(toBeRemoved != address(0), NoOverriddenAddressToRemove());
        if (migrateExistingRewards) { _migrateRewards(toBeRemoved, msg.sender); }
        overrideAddresses[msg.sender] = address(0);
        emit OverrideAddressRemoved(msg.sender);
    }

    /// @dev Allows a reward recipient to claim their rewards.
    function claimRewards() external whenNotPaused nonReentrant { _claimRewards(); }

    /// @dev Allows the owner to claim orphaned rewards to appropriate addresses.
    function claimOrphanedRewards(bytes[] calldata pubkeys, address toPay) external onlyOwner nonReentrant {
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
        require(success, OrphanedRewardsClaimFailed());
        emit OrphanedRewardsClaimed(toPay, totalAmount);
    }

    /// @dev Allows the owner to remove an address from the auto claim blacklist.
    function removeFromAutoClaimBlacklist(address addr) external onlyOwner {
        autoClaimBlacklist[addr] = false;
        emit RemovedFromAutoClaimBlacklist(addr);
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

    function _claimRewards() internal {
        uint256 amount = unclaimedRewards[msg.sender];
        if (amount == 0) {
            emit NoRewards(msg.sender);
            return;
        }
        unclaimedRewards[msg.sender] = 0;
        (bool success, ) = payable(msg.sender).call{value: amount}("");
        require(success, RewardsClaimFailed());
        emit RewardsClaimed(msg.sender, amount);
    }

    /// @dev DANGER: This function should ONLY be called from overrideClaimAddress or removeOverriddenClaimAddress
    /// with careful attention to parameter order.
    function _migrateRewards(address from, address to) internal {
        uint256 amount = unclaimedRewards[from];
        if (amount == 0) {
            emit NoRewards(from);
            return;
        }
        unclaimedRewards[from] = 0;
        unclaimedRewards[to] += amount;
        emit RewardsMigrated(from, to, amount);
    }

    function _setVanillaRegistry(address vanillaRegistry) internal {
        require(vanillaRegistry != address(0), InvalidAddress());
        _vanillaRegistry = VanillaRegistryStorage(vanillaRegistry);
        emit VanillaRegistrySet(vanillaRegistry);
    }

    function _setMevCommitAVS(address mevCommitAVS) internal {
        require(mevCommitAVS != address(0), InvalidAddress());
        _mevCommitAVS = MevCommitAVSStorage(mevCommitAVS);
        emit MevCommitAVSSet(mevCommitAVS);
    }

    function _setMevCommitMiddleware(address mevCommitMiddleware) internal {
        require(mevCommitMiddleware != address(0), InvalidAddress());
        _mevCommitMiddleware = MevCommitMiddlewareStorage(mevCommitMiddleware);
        emit MevCommitMiddlewareSet(mevCommitMiddleware);
    }

    function _setAutoClaimGasLimit(uint256 limit) internal {
        require(limit != 0, InvalidAutoClaimGasLimit());
        autoClaimGasLimit = limit;
        emit AutoClaimGasLimitSet(limit);
    }

    /// @dev Finds the receiver address for a given validator pubkey,
    /// corresponding to state that'd exist in each type of registry if the pubkey were opted-in through that registry.
    /// @notice Zero address is returned if the pubkey is not opted-in to mev-commit.
    function _findReceiver(bytes calldata pubkey) internal view returns (address) {
        (,address operatorAddr,bool existsMiddleware,) = _mevCommitMiddleware.validatorRecords(pubkey);
        if (existsMiddleware && operatorAddr != address(0)) {
            return operatorAddr;
        }
        (bool existsVanilla,address vanillaWithdrawalAddr,,) = _vanillaRegistry.stakedValidators(pubkey);
        if (existsVanilla && vanillaWithdrawalAddr != address(0)) {
            return vanillaWithdrawalAddr;
        }
        (bool existsAvs,address podOwner,,) = _mevCommitAVS.validatorRegistrations(pubkey);
        if (existsAvs && podOwner != address(0)) {
            return podOwner;
        }
        return address(0);
    }
}
