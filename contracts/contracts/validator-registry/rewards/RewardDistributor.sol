// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRewardDistributor} from "../../interfaces/IRewardDistributor.sol";
import {RewardDistributorStorage} from "./RewardDistributorStorage.sol";
import {Errors} from "../../utils/Errors.sol";

contract RewardDistributor is IRewardDistributor, RewardDistributorStorage,
    Ownable2StepUpgradeable, ReentrancyGuardUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    using SafeERC20 for IERC20;

    modifier onlyOwnerOrRewardManager() {
        require(msg.sender == rewardManager || msg.sender == owner(), NotOwnerOrRewardManager());
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
        address _owner,
        address _rewardManager
    ) external initializer {
        __Ownable_init(_owner);
        __ReentrancyGuard_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
        _setRewardManager(_rewardManager);
    }

    /// @param rewardList Array of ETH Distributions.
    function grantETHRewards(Distribution[] calldata rewardList) external payable nonReentrant whenNotPaused onlyOwnerOrRewardManager {
        uint256 len = rewardList.length;
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < len; ++i) {
            totalAmount += rewardList[i].amount;
            rewardData[rewardList[i].operator][rewardList[i].recipient][0].accrued += rewardList[i].amount;
            emit ETHGranted(rewardList[i].operator, rewardList[i].recipient, rewardList[i].amount);
        }
        emit RewardsBatchGranted(0, totalAmount);
        require(msg.value == totalAmount, IncorrectPaymentAmount(msg.value, totalAmount));
    }

    /// @param rewardList Array of token Distributions.
    function grantTokenRewards(Distribution[] calldata rewardList, uint256 tokenID) external payable nonReentrant whenNotPaused onlyOwnerOrRewardManager {
        uint256 len = rewardList.length;
        uint256 totalAmount = 0;
        address rewardToken = rewardTokens[tokenID];
        require(rewardToken != address(0), InvalidRewardToken());
        for (uint256 i = 0; i < len; ++i) {
            totalAmount += rewardList[i].amount;
            rewardData[rewardList[i].operator][rewardList[i].recipient][tokenID].accrued += rewardList[i].amount;
            emit TokensGranted(rewardList[i].operator, rewardList[i].recipient, rewardList[i].amount);
        }
        emit RewardsBatchGranted(tokenID, totalAmount);
        IERC20(rewardToken).safeTransferFrom(msg.sender, address(this), totalAmount);
    }

    /// @notice Claim rewards for the caller (as operator) to specific recipients.
    /// @param recipients List of recipients to claim rewards for.
    /// @param tokenID The ID of the token to claim rewards for. 0 for ETH.
    function claimRewards(address[] calldata recipients, uint256 tokenID) external whenNotPaused nonReentrant {
        _claimRewards(msg.sender, recipients, tokenID);
    }

    /// @notice Claim rewards on behalf of an operator to specific recipients (must be delegated).
    /// @param operator Operator to claim rewards for.
    /// @param recipients List of recipients to claim rewards for.
    /// @param tokenID The ID of the token to claim rewards for. 0 for ETH.
    function claimOnbehalfOfOperator(address operator, address[] calldata recipients, uint256 tokenID) external whenNotPaused nonReentrant {
        uint256 len = recipients.length;
        for (uint256 i = 0; i < len; ++i) {
            require(claimDelegate[operator][recipients[i]][msg.sender], InvalidClaimDelegate());
        }
        _claimRewards(operator, recipients, tokenID);
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
    /// @param tokenID The ID of the token to migrate rewards for.
    function migrateExistingRewards(address from, address to, uint256 tokenID) external whenNotPaused nonReentrant {
        require(to != address(0), ZeroAddress());
        require(to != from, InvalidRecipient());
        require(tokenID == 0 || rewardTokens[tokenID] != address(0), InvalidRewardToken());
        uint128 claimableAmt = getPendingRewards(msg.sender, from, tokenID);
        require(claimableAmt > 0, NoClaimableRewards(msg.sender, from));
        rewardData[msg.sender][from][tokenID].accrued -= claimableAmt;
        rewardData[msg.sender][to][tokenID].accrued += claimableAmt;
        emit RewardsMigrated(tokenID, msg.sender, from, to, claimableAmt);
    }

    /// @dev Allows the owner to reclaim stipends that were incorrectly granted or unable to be claimed by an operator.
    function reclaimStipendsToOwner(address[] calldata operators, address[] calldata recipients, uint256 tokenID) external onlyOwner {
        require(tokenID == 0 || rewardTokens[tokenID] != address(0), InvalidRewardToken());
        address _owner = owner();
        uint256 toWithdraw = 0;
        uint256 len = operators.length;
        require(len == recipients.length, LengthMismatch());
        for (uint256 i = 0; i < len; ++i) {
            address operator = operators[i];
            address recipient = recipients[i];
            uint128 claimableAmt = getPendingRewards(operator, recipient, tokenID);
            rewardData[operator][recipient][tokenID].accrued -= claimableAmt;
            toWithdraw += claimableAmt;
            emit RewardsReclaimed(tokenID, operator, recipient, claimableAmt);
        }
        require(toWithdraw > 0, NoClaimableRewards(_owner, _owner));
        _transferFunds(_owner, _owner, toWithdraw, tokenID);
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
    function setRewardManager(address _rewardManager) external onlyOwner {
        _setRewardManager(_rewardManager);
    }

    /// @dev Allows the owner to set a reward token address for a given id.
    function setRewardToken(address _rewardToken, uint256 _id) external onlyOwner {
        _setRewardToken(_rewardToken, _id);
    }

    // Retreives the recipient for an operator's registered key
    function getKeyRecipient(address operator, bytes calldata pubkey) external view returns (address) {
        require(pubkey.length == 48, InvalidBLSPubKeyLength());
        require(operator != address(0), InvalidOperator());
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

    function getPendingRewards(address operator, address recipient, uint256 tokenID) public view returns (uint128) {
        return rewardData[operator][recipient][tokenID].accrued - rewardData[operator][recipient][tokenID].claimed;
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @dev Allows a reward recipient to claim their rewards.
    function _claimRewards(address operator, address[] calldata recipients, uint256 tokenID) internal {
        require(operator != address(0), InvalidOperator());
        require(tokenID == 0 || rewardTokens[tokenID] != address(0), InvalidRewardToken());
        uint256 len = recipients.length;
        uint128[] memory claimAmounts = new uint128[](len);
        for (uint256 i = 0; i < len; ++i) {
            address recipient = recipients[i];
            claimAmounts[i] = getPendingRewards(operator, recipient, tokenID);
            rewardData[operator][recipient][tokenID].claimed += claimAmounts[i];
        }
        for (uint256 i = 0; i < len; ++i) {
            address recipient = recipients[i];
            if (claimAmounts[i] > 0) {
                _transferFunds(operator, recipient, claimAmounts[i], tokenID);
            }
        }
    }

    function _transferFunds(address operator, address recipient, uint256 amount, uint256 tokenID) internal {
        if (tokenID == 0) {
            (bool success, ) = payable(recipient).call{value: amount}("");
            require(success, RewardsTransferFailed(recipient));
            emit ETHRewardsClaimed(operator, recipient, amount);
        } else {
            IERC20(rewardTokens[tokenID]).safeTransfer(recipient, amount);
            emit TokenRewardsClaimed(operator, recipient, amount);
        }
    }

    function _setRewardManager(address _rewardManager) internal {
        require(_rewardManager != address(0), ZeroAddress());
        rewardManager = _rewardManager;
        emit RewardManagerSet(_rewardManager);
    }

    function _setRewardToken(address _rewardToken, uint256 _id) internal {
        require(_id != 0, InvalidTokenID());
        rewardTokens[_id] = _rewardToken;
        emit RewardTokenSet(_rewardToken, _id);
    }
}
