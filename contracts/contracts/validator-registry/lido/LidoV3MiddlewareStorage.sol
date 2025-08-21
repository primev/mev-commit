// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {TimestampOccurrence} from "../../utils/Occurrence.sol";

/**
 * @dev Storage layout for LidoV3Middleware (upgrade-safe).
 * Keep ordering and types stable across upgrades.
 */
abstract contract LidoV3MiddlewareStorage {
    struct ValidatorRecord {
        bool exists;
        address registrar; // operator who registered this pubkey
        TimestampOccurrence.Occurrence freezeOccurrence;
        TimestampOccurrence.Occurrence deregRequestOccurrence;
    }

    // --- Config ---
    address public vaultHub;              // canonical VaultHub to verify connections
    uint256 public slashAmount;           // denominator for capacity (ETH / slashAmount)
    uint256 public unfreezeFee;           // fee per validator to unfreeze
    address public unfreezeReceiver;      // where unfreeze fees go
    uint256 public unfreezePeriod;        // min seconds from freeze until unfreeze allowed
    uint256 public deregistrationPeriod;  // min seconds from request until deregister allowed

    // --- Operator gating ---
    mapping(address => bool) public isWhitelisted; // whitelisted operators

    // mapping(pubkey -> record). Keys must be 48-byte BLS pubkeys.
    mapping(bytes => ValidatorRecord) public validatorRecords;
    // How many validators each vault has registered here
    mapping(address vault => uint256) public vaultRegisteredCount;

    // --- Storage gap for upgrades ---
    uint256[48] private __gap;
}
