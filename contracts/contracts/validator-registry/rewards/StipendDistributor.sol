// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IRewardManager} from "../../interfaces/IRewardManager.sol";
import {RewardManagerStorage} from "./RewardManagerStorage.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IStipendDistributor} from "../../interfaces/IStipendDistributor.sol";
import {StipendDistributorStorage} from "./StipendDistributorStorage.sol";
import {TransientReentrancyGuard} from "../../utils/TransientReentrancyGuard.sol";
import {Errors} from "../../utils/Errors.sol";

import {VanillaRegistryStorage} from "../VanillaRegistryStorage.sol";
import {MevCommitAVSStorage} from "../avs/MevCommitAVSStorage.sol";
import {MevCommitMiddlewareStorage} from "../middleware/MevCommitMiddlewareStorage.sol";


contract StipendDistributor is IStipendDistributor, StipendDistributorStorage,
    Ownable2StepUpgradeable, TransientReentrancyGuard, PausableUpgradeable, UUPSUpgradeable {


    modifier onlyOwnerOrOracle() {
        require(msg.sender == oracle || msg.sender == owner(), NotOwnerOrOracle());
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    constructor() {
        _disableInitializers();
    }

    /// @dev Receive function is disabled prevent misc transfers.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function disabled to prevent misc transfers.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @dev Initializes the RewardManager contract.
    function initialize(
        address owner,
        address oracle,
        address vanillaRegistry,
        address mevCommitAVS,
        address mevCommitMiddleware
    ) external initializer {
        __Ownable_init(owner);
        __Pausable_init();
        __UUPSUpgradeable_init();
        _setOracle(oracle);
        _setVanillaRegistry(vanillaRegistry);
        _setMevCommitAVS(mevCommitAVS);
        _setMevCommitMiddleware(mevCommitMiddleware);
    }

    /// @dev Grant stipends to multiple (operator, recipient) pairs.
    /// @param operators Array of operator addresses.
    /// @param recipients Array of recipient addresses of the corresponding operator.
    /// @param amounts Array of stipend amounts.
    function grantStipends(address[] calldata operators, address[] calldata recipients, uint256[] calldata amounts) external payable nonReentrant whenNotPaused onlyOwnerOrOracle {
        require(operators.length == amounts.length && operators.length == recipients.length, LengthMismatch());
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            accrued[operators[i]][recipients[i]] += amounts[i];
            emit StipendsGranted(operators[i], recipients[i], amounts[i]);
        }
    }

    /// @notice Allows an operator to claim their rewards for specified recipients.
    function claimRewards(address payable[] calldata recipients) external whenNotPaused nonReentrant {
        _claimRewards(msg.sender, recipients);
    }

    /// @notice Claims rewards accrued by an operator to a specific recipient. Must be authorized by the specified operator.
    /// @dev Caller must be an authorized delegate for every (operator → recipient) pair.
    function claimOnbehalfOfOperator(address operator, address payable[] calldata recipients) external whenNotPaused nonReentrant {
        uint256 len = recipients.length;
        for (uint256 i = 0; i < len; ++i) {
            require(claimDelegate[operator][recipients[i]][msg.sender], InvalidClaimDelegate());
        }
        _claimRewards(operator, recipients);
    }

    /// @notice Allows an operator to set the recipient for a list of pubkeys.
    /// @dev If operator is no longer valid at the time of stipend distribution, the recipient will not receive the stipend.
    ///      If the key has a new operator that has not updated the key's recipient, the new operator will receive the stipend.
    /// @param pubkeys List of pubkeys to set the recipient for.
    /// @param registryID Registry in which the pubkeys are registered.
    /// @param recipient Recipient to set for the pubkeys.
    function overrideRecipientByPubkey(bytes[] calldata pubkeys, uint256 registryID, address recipient) external whenNotPaused nonReentrant {
        require(recipient != address(0), ZeroAddress());
        for (uint256 i = 0; i < pubkeys.length; ++i) {
            bytes calldata pubkey = pubkeys[i];
            require(pubkey.length == 48, InvalidBLSPubKeyLength(48, pubkey.length));
            require(msg.sender == findOperator(pubkey, registryID), InvalidRecipient());
            bytes32 pkHash = keccak256(pubkey);
            operatorKeyOverrides[msg.sender][pkHash] = recipient;  
            emit RecipientSet(msg.sender, pubkey, registryID, recipient);
        }
    }

    /// @dev Allows an operator to set a default recipient for all non-overridden keys.
    ///      If a recipient is set for a specific key, it will override the default recipient.
    /// @param recipient Default recipient to set for all non-overridden keys of the operator.
    function setDefaultRecipient(address recipient) external whenNotPaused nonReentrant {
        require(recipient != address(0), ZeroAddress());
        defaultRecipient[msg.sender] = recipient;
        emit DefaultRecipientSet(msg.sender, recipient);
    }

    /// @dev Allows an operator to set a delegate to claim rewards for one of their recipients.
    function setClaimDelegate(address delegate, address recipient, bool status) external whenNotPaused nonReentrant {
        claimDelegate[msg.sender][recipient][delegate] = status;
        emit ClaimDelegateSet(msg.sender, recipient, delegate, status);
    }

    /// @dev Allows an operator to migrate unclaimed recipient rewards to a different address.
    function migrateExistingRewards(address from, address to) external whenNotPaused nonReentrant {
        uint256 claimableAmt = accrued[msg.sender][from] - claimed[msg.sender][from];
        require(claimableAmt > 0, NoClaimableRewards(from));
        require(to != address(0), ZeroAddress());
        claimed[msg.sender][from] += claimableAmt;
        accrued[msg.sender][to] += claimableAmt;
        emit RewardsMigrated(from, to, claimableAmt);
    }
    
    /// @dev Enables the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Enables the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /// @dev Allows the owner to set the oracle address.
    function setOracle(address _oracle) external onlyOwner {
        _setOracle(_oracle);
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

    // --- Getters ---
 
    // Retreives the recipient for an operator's registered key
    function getKeyRecipient(bytes calldata pubkey) external view returns (address) {
        require(pubkey.length == 48, InvalidBLSPubKeyLength(48, pubkey.length));
        bytes32 pkHash = keccak256(pubkey);
        address registeredOperator = findOperator(pubkey, 0);
        // If the key is not registered, return 0
        if (registeredOperator == address(0)) {
            return address(0);
        }
        // Individual key overrides take priority over the default recipient
        if (operatorKeyOverrides[registeredOperator][pkHash] != address(0)) {
            return operatorKeyOverrides[registeredOperator][pkHash];
        }
        // If no key override, return the default recipient
        address defaultOverride = defaultRecipient[registeredOperator];
        if (defaultOverride != address(0)) {
            return defaultOverride;
        }
        // If no default override, return the operator
        return registeredOperator;
    }

    function getPendingRewards(address operator, address recipient) public view returns (uint256) {
        return accrued[operator][recipient] - claimed[operator][recipient];
    }

    /// @dev Finds the operator address for a given validator pubkey based on the registry ID.
    // A registry id of 0 can be used to check all registries
    function findOperator(bytes calldata pubkey, uint256 registryID) public view returns (address) {
        if (registryID == 1) {
            (bool existsAvs,address podOwner,,) = _mevCommitAVS.validatorRegistrations(pubkey);
            if (existsAvs && podOwner != address(0)) {
                return podOwner;
            }
        } else if (registryID == 2) {
            (,address operatorAddr,bool existsMiddleware,) = _mevCommitMiddleware.validatorRecords(pubkey);
            if (existsMiddleware && operatorAddr != address(0)) {
                return operatorAddr;
            }
        } else if (registryID == 3) {
            (bool existsVanilla,address vanillaWithdrawalAddr,,) = _vanillaRegistry.stakedValidators(pubkey);
            if (existsVanilla && vanillaWithdrawalAddr != address(0)) {
                return vanillaWithdrawalAddr;
            }
        } else if (registryID == 0) {
            for (uint256 i = 1; i < 4; ++i) {
                address operator = findOperator(pubkey, i);
                if (operator != address(0)) {
                    return operator;
                }
            }
        }
        return address(0);
    }

    // --- Internal Functions ---

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @dev Allows a reward recipient to claim their rewards.
    function _claimRewards(address operator, address payable[] calldata recipients) internal {
        require(operator != address(0), InvalidOperator());
        uint256 len = recipients.length;
        uint256[] memory claimAmounts = new uint256[](len);
        for (uint256 i = 0; i < len; ++i) {
            address recipient = recipients[i];
            claimAmounts[i] = getPendingRewards(operator, recipient);
            claimed[operator][recipient] += claimAmounts[i];
        }
        for (uint256 i = 0; i < len; ++i) {
            address recipient = recipients[i];
            if (claimAmounts[i] > 0) {
                (bool success, ) = recipient.call{value: claimAmounts[i]}("");
                require(success, RewardsTransferFailed(recipient));
                emit RewardsClaimed(operator, recipient, claimAmounts[i]);
            }
        }
    }

    function _setOracle(address _oracle) internal {
        oracle = _oracle;
        emit OracleSet(oracle);
    }

    function _setVanillaRegistry(address vanillaRegistry) internal {
        _vanillaRegistry = VanillaRegistryStorage(vanillaRegistry);
        emit VanillaRegistrySet(vanillaRegistry);
    }

    function _setMevCommitAVS(address mevCommitAVS) internal {
        _mevCommitAVS = MevCommitAVSStorage(mevCommitAVS);
        emit MevCommitAVSSet(mevCommitAVS);
    }

    function _setMevCommitMiddleware(address mevCommitMiddleware) internal {
        _mevCommitMiddleware = MevCommitMiddlewareStorage(mevCommitMiddleware);
        emit MevCommitMiddlewareSet(mevCommitMiddleware);
    }

}
