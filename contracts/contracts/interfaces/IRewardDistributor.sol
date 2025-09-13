// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title IStipendDistributor
/// @notice Interface for stipend distribution and claims.
interface IRewardDistributor {

    struct Distribution {
        address operator;
        address recipient;
        uint128 amount;
    }

    /// @dev Pack both counters into a single slot for each asset.
    struct RewardData {
        uint128 accrued;
        uint128 claimed;
    }

    // -------- Events --------
    /// @dev Emitted when the oracle address is updated.
    event RewardManagerSet(address indexed rewardManager);

    /// @dev Emitted when stipends are granted.
    event ETHGranted(address indexed operator, address indexed recipient, uint256 indexed amount);
    event TokensGranted(address indexed operator, address indexed recipient, uint256 indexed amount);
    event RewardsBatchGranted(uint256 indexed tokenID, uint256 indexed amount);
    /// @dev Emitted when rewards are claimed by a recipient for an operator.
    event ETHRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount);
    event TokenRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount);

    /// @dev Emitted when a recipient mapping is overridden for a specific pubkey.
    event RecipientSet(address indexed operator, bytes pubkey, address indexed recipient);

    /// @dev Emitted when an operator sets/updates their global override recipient.
    event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient);

    /// @dev Emitted when an operator sets/updates a claim delegate for a given recipient.
    event ClaimDelegateSet(address indexed operator, address indexed recipient, address indexed delegate, bool status);

    /// @dev Emitted when accrued rewards are migrated from one recipient to another for an operator.
    event RewardsMigrated(uint256 tokenID, address indexed operator, address indexed from, address indexed to, uint128 amount);

    /// @dev Emitted when accrued rewards are reclaimed by the owner.
    event RewardsReclaimed(uint256 indexed tokenID, address indexed operator, address indexed recipient, uint256 amount);

    /// @dev Emitted when the reward token address is updated.
    event RewardTokenSet(address indexed rewardToken, uint256 indexed tokenID);

    // -------- Errors --------
    error NotOwnerOrRewardManager();
    error InvalidRewardToken();
    error ZeroAddress();
    error InvalidTokenID();
    error InvalidBLSPubKeyLength();
    error InvalidRecipient();
    error InvalidOperator();
    error InvalidClaimDelegate();
    error LengthMismatch();
    error NoClaimableRewards(address operator, address recipient);
    error RewardsTransferFailed(address recipient);
    error IncorrectPaymentAmount(uint256 received, uint256 expected);

    // -------- Externals --------
    /// @notice Initialize the proxy.
    function initialize(address owner, address rewardManager) external;

    /// @notice Grant ETH rewards to multiple (operator, recipient) pairs.
    function grantETHRewards(Distribution[] calldata rewardList) external payable;
 
    /// @notice Grant token rewards to multiple (operator, recipient) pairs.
    function grantTokenRewards(Distribution[] calldata rewardList, uint256 tokenID) external payable;

    /// @notice Claim rewards for the caller (as operator) to specific recipients.
    function claimRewards(address[] calldata recipients, uint256 tokenID) external;

    /// @notice Claim rewards on behalf of an operator to specific recipients (must be delegated).
    function claimOnbehalfOfOperator(address operator, address[] calldata recipients, uint256 tokenID) external;

    /// @notice Override recipient for a list of BLS pubkeys in a registry.
    function overrideRecipientByPubkey(bytes[] calldata pubkeys, address recipient) external;

    /// @notice Set the caller's global override recipient for any non-overridden keys.
    function setOperatorGlobalOverride(address recipient) external;

    /// @notice Allow or revoke a delegate to claim for a given recipient of the caller (operator).
    function setClaimDelegate(address delegate, address recipient, bool status) external;

    /// @notice Migrate unclaimed rewards from one recipient to another for the caller (operator).
    function migrateExistingRewards(address from, address to, uint256 tokenID) external;

    /// @notice Pause / Unpause admin controls.
    function reclaimStipendsToOwner(address[] calldata operators, address[] calldata recipients, uint256 tokenID) external;
    function pause() external;
    function unpause() external;
    function setRewardManager(address _rewardManager) external;
    function setRewardToken(address _rewardToken, uint256 _id) external;
    function getKeyRecipient(address operator, bytes calldata pubkey) external view returns (address);
    function getPendingRewards(address operator, address recipient, uint256 tokenID) external view returns (uint128);

}