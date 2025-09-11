// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IStipendDistributor} from "../../interfaces/IStipendDistributor.sol";
import {StipendDistributorStorage} from "./StipendDistributorStorage.sol";
import {Errors} from "../../utils/Errors.sol";

contract StipendDistributor is IStipendDistributor, StipendDistributorStorage,
    Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    modifier onlyOwnerOrStipendManager() {
        require(msg.sender == stipendManager || msg.sender == owner(), NotOwnerOrStipendManager());
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
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
        address stipendManager
    ) external initializer {
        __Ownable_init(owner);
        __ReentrancyGuard_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        _setStipendManager(stipendManager);
    }

    /// @dev Grant stipends to multiple (operator, recipient) pairs.
    /// @param stipends Array of stipends.
    function grantStipends(Stipend[] calldata stipends) external payable nonReentrant whenNotPaused onlyOwnerOrStipendManager {
        uint256 len = stipends.length;
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < len; ++i) {
            totalAmount += stipends[i].amount;
            accrued[stipends[i].operator][stipends[i].recipient] += stipends[i].amount;
            emit StipendsGranted(stipends[i].operator, stipends[i].recipient, stipends[i].amount);
        }
        require(msg.value == totalAmount, IncorrectPaymentAmount(msg.value, totalAmount));
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
    /// @param pubkeys List of pubkeys to set the recipient for.
    /// @param recipient Recipient to set for the pubkeys.
    function overrideRecipientByPubkey(bytes[] calldata pubkeys, address recipient) external whenNotPaused nonReentrant {
        require(recipient != address(0), ZeroAddress());
        uint256 len = pubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubkey = pubkeys[i];
            require(pubkey.length == 48, InvalidBLSPubKeyLength());
            bytes32 pkHash = keccak256(pubkey);
            operatorKeyOverrides[msg.sender][pkHash] = recipient;
            emit RecipientSet(msg.sender, pubkey, recipient);
        }
    }

    /// @dev Allows an operator to set a default recipient for all non-overridden keys.
    ///      If a recipient is set for a specific key, it will override the default recipient.
    /// @param recipient Default recipient to set for all non-overridden keys of the operator.
    function setOperatorGlobalOverride(address recipient) external whenNotPaused nonReentrant {
        require(recipient != address(0), ZeroAddress());
        operatorGlobalOverride[msg.sender] = recipient;
        emit OperatorGlobalOverrideSet(msg.sender, recipient);
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

    /// @dev Allows the owner to set the stipend manager address.
    function setStipendManager(address _stipendManager) external onlyOwner {
        _setStipendManager(_stipendManager);
    }

    // Retreives the recipient for an operator's registered key
    function getKeyRecipient(address operator, bytes calldata pubkey) external view returns (address) {
        require(pubkey.length == 48, InvalidBLSPubKeyLength());
        bytes32 pkHash = keccak256(pubkey);
        // Individual key overrides take priority over the default recipient
        if (operatorKeyOverrides[operator][pkHash] != address(0)) {
            return operatorKeyOverrides[operator][pkHash];
        }
        // If no key override, return the default recipient
        address defaultOverride = operatorGlobalOverride[operator];
        if (defaultOverride != address(0)) {
            return defaultOverride;
        }
        // If no default override, return the operator
        return operator;
    }

    function getPendingRewards(address operator, address recipient) public view returns (uint256) {
        return accrued[operator][recipient] - claimed[operator][recipient];
    }

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

    function _setStipendManager(address _stipendManager) internal {
        require(_stipendManager != address(0), ZeroAddress());
        stipendManager = _stipendManager;
        emit StipendManagerSet(_stipendManager);
    }
}
