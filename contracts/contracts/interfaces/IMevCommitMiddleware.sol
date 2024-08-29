// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {EventHeightLib} from "../utils/EventHeight.sol";

interface IMevCommitMiddleware {

    struct ValidatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        /// @notice The vault holding slashable stake which represents the validator.
        address vault;
    }

    // TODO: Confirm we don't need to store/track slashing of operators. This must be handled somewhere tho.
    struct OperatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        bool isBlacklisted;
    }

    struct VaultRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        // TODO: For now, a single operator can register multiple vaults,
        // a single vault can collateralize only one operator.
        // A single vault can collateralize multiple validators.
        // Evaluate how this compares to intended usage of Symbiotic.
        address operator;
    }

    /// @notice Emmitted when an operator is registered
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when an operator requests deregistration
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when an operator is blacklisted
    event OperatorBlacklisted(address indexed operator);

    /// @notice Emmitted when a validator record is added to state
    event ValRecordAdded(bytes indexed blsPubkey, address indexed operator,
        uint256 indexed position);

    /// @notice Emmitted when validator deregistration is requested
    event ValidatorDeregistrationRequested(bytes indexed blsPubkey, address indexed operator,
        uint256 indexed position);

    /// @notice Emmitted when a validator record is deleted by the contract owner
    event ValRecordDeleted(bytes indexed blsPubkey, address indexed operator);

    /// @notice Emmitted when a vault record is added
    event VaultRegistered(address indexed vault, address indexed operator);

    /// @notice Emmitted when a vault deregistration is requested
    event VaultDeregistrationRequested(address indexed vault);

    /// @notice Emmitted when a vault is deregistered
    event VaultDeregistered(address indexed vault);
    
    /// @notice Emmitted when a validator is slashed
    event ValidatorSlashed(bytes indexed blsPubkey, address indexed operator, uint256 indexed position);

    /// @notice Emmitted when the operator deregistration period in blocks is set
    event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks);

    /// @notice Emmitted when the validator deregistration period in blocks is set
    event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks);

    /// @notice Emmitted when the vault deregistration period in blocks is set
    event VaultDeregPeriodBlocksSet(uint256 vaultDeregPeriodBlocks);

    /// @notice Emmitted when the slash oracle is set
    event SlashOracleSet(address slashOracle);
}
