// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title StipendDistributorStorage
/// @notice Storage layout for StipendDistributor
abstract contract StipendDistributorStorage {
    /// @dev Address authorized to grant stipends.
    address public stipendManager;

    /// @dev Default recipient per operator (used when no pubkey-specific override exists).
    mapping(address operator => address recipient) public defaultRecipient;

    /// @dev Recipient override by BLS pubkey hash (keccak256(pubkey)).
    mapping(address operator => mapping(bytes32 keyhash => address recipient)) public operatorKeyOverrides;

    /// @dev Accrued and claimed amounts per (operator, recipient).
    mapping(address operator => mapping(address recipient => uint256 amount)) public accrued;
    mapping(address operator => mapping(address recipient => uint256 amount)) public claimed;

    /// @dev Operator → recipient → delegate → isAuthorized
    mapping(address operator => mapping(address recipient => mapping(address delegate => bool))) public claimDelegate;
    
    // === Storage gap for future upgrades ===
    uint256[40] private __gap;
}