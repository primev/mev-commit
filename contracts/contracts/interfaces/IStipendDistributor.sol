// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title IStipendDistributor
/// @notice Interface for stipend distribution and claims.
interface IStipendDistributor {

    struct Stipend {
        address operator;
        address recipient;
        uint256 amount;
    }

    // -------- Events --------
    /// @dev Emitted when the oracle address is updated.
    event StipendManagerSet(address indexed stipendManager);

    /// @dev Emitted when stipends are granted.
    event StipendsGranted(address indexed operator, address indexed recipient, uint256 amount);

    /// @dev Emitted when rewards are claimed by a recipient for an operator.
    event RewardsClaimed(address indexed operator, address indexed recipient, uint256 amount);

    /// @dev Emitted when a recipient mapping is overridden for a specific pubkey.
    event RecipientSet(address indexed operator, bytes pubkey, address indexed recipient);

    /// @dev Emitted when an operator sets/updates their global override recipient.
    event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient);

    /// @dev Emitted when an operator sets/updates a claim delegate for a given recipient.
    event ClaimDelegateSet(address indexed operator, address indexed recipient, address indexed delegate, bool status);

    /// @dev Emitted when accrued rewards are migrated from one recipient to another for an operator.
    event RewardsMigrated(address indexed from, address indexed to, uint256 amount);

    /// @dev Emitted when accrued rewards are reclaimed by the owner.
    event StipendsReclaimed(address indexed operator, address indexed recipient, uint256 amount);

    // -------- Errors --------
    error NotOwnerOrStipendManager();
    error ZeroAddress();
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
    function initialize(address owner, address stipendManager) external;

    function grantStipends(Stipend[] calldata stipends) external payable;

    /// @notice Claim rewards for the caller (as operator) to specific recipients.
    function claimRewards(address[] calldata recipients) external;

    /// @notice Claim rewards on behalf of an operator to specific recipients (must be delegated).
    function claimOnbehalfOfOperator(address operator, address[] calldata recipients) external;

    /// @notice Override recipient for a list of BLS pubkeys in a registry.
    function overrideRecipientByPubkey(bytes[] calldata pubkeys, address recipient) external;

    /// @notice Set the caller's global override recipient for any non-overridden keys.
    function setOperatorGlobalOverride(address recipient) external;

    /// @notice Allow or revoke a delegate to claim for a given recipient of the caller (operator).
    function setClaimDelegate(address delegate, address recipient, bool status) external;

    /// @notice Migrate unclaimed rewards from one recipient to another for the caller (operator).
    function migrateExistingRewards(address from, address to) external;

    /// @notice Pause / Unpause admin controls.
    function pause() external;
    function unpause() external;
    function setStipendManager(address _stipendManager) external;
    function getKeyRecipient(address operator, bytes calldata pubkey) external view returns (address);
    function getPendingRewards(address operator, address recipient) external view returns (uint256);

}