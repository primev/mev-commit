// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {TimestampOccurrence} from "../utils/Occurrence.sol";

interface IMevCommitMiddleware {

    struct OperatorRecord {
        TimestampOccurrence.Occurrence deregRequestOccurrence;
        bool exists;
        bool isBlacklisted;
    }

    struct VaultRecord {
        bool exists;
        TimestampOccurrence.Occurrence deregRequestOccurrence;
        uint256 slashAmount;
    }

    struct ValidatorRecord {
        /// @notice The vault holding slashable stake which represents the validator.
        address vault;
        /// @notice The operator which registered this validator pubkey with a vault.
        address operator;
        bool exists;
        TimestampOccurrence.Occurrence deregRequestOccurrence;
    }

    /// @notice Emmitted when an operator is registered
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when an operator requests deregistration
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when an operator is blacklisted
    event OperatorBlacklisted(address indexed operator);

    /// @notice Emmitted when a vault record is added
    event VaultRegistered(address indexed vault, uint256 slashAmount);

    /// @notice Emmitted when a vault slash amount is updated
    event VaultSlashAmountUpdated(address indexed vault, uint256 slashAmount);

    /// @notice Emmitted when a vault deregistration is requested
    event VaultDeregistrationRequested(address indexed vault);

    /// @notice Emmitted when a vault is deregistered
    event VaultDeregistered(address indexed vault);

    /// @notice Emmitted when a validator record is added to state
    event ValRecordAdded(bytes blsPubkey, address indexed msgSender,
        uint256 indexed position);

    /// @notice Emmitted when validator deregistration is requested
    event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender,
        uint256 indexed position);

    /// @notice Emmitted when a validator record is deleted by the contract owner
    event ValRecordDeleted(bytes blsPubkey, address indexed msgSender);
    
    /// @notice Emmitted when a validator is slashed
    event ValidatorSlashed(bytes blsPubkey, address indexed operator, uint256 indexed position);

    /// @notice Emmitted when the network registry is set
    event NetworkRegistrySet(address networkRegistry);

    /// @notice Emmitted when the operator registry is set
    event OperatorRegistrySet(address operatorRegistry);

    /// @notice Emmitted when the vault factory is set
    event VaultFactorySet(address vaultFactory);

    /// @notice Emmitted when the network is set
    event NetworkSet(address network);

    /// @notice Emmitted when the slash period in seconds is set
    event SlashPeriodSecondsSet(uint256 slashPeriodSeconds);

    /// @notice Emmitted when the slash oracle is set
    event SlashPeriodBlocksSet(uint256 slashPeriodBlocks);

    /// @notice Emmitted when the slash oracle is set
    event SlashOracleSet(address slashOracle);

    error OnlySlashOracle(address slashOracle);

    error OnlyOperator(address operator);

    error InvalidArrayLengths(uint256 expectedLength, uint256 actualLength);

    error ValidatorsNotSlashable(address vault, address operator,
        uint256 numRequested, uint256 potentialSlashableVals);

    error MissingValRecord(bytes blsPubkey);

    error OperatorAlreadyRegistered(address operator);

    error OperatorNotEntity(address operator);

    error OperatorNotRegistered(address operator);

    error OperatorDeregRequestExists(address operator);

    error OperatorIsBlacklisted(address operator);

    error OperatorNotReadyToDeregister(address operator, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error OperatorAlreadyBlacklisted(address operator);

    error ValidatorRecordAlreadyExists(bytes blsPubkey);

    error MissingValidatorRecord(bytes blsPubkey);

    error ValidatorNotReadyToDeregister(bytes blsPubkey, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error VaultAlreadyRegistered(address vault);

    error MissingVaultRecord(address vault);

    error VaultNotEntity(address vault);

    error VaultNotReadyToDeregister(address vault, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error SlashAmountMustBeNonZero(address vault);

    error InvalidVaultEpochDuration(address vault, uint256 vaultEpochDurationSec, uint256 slashPeriodSec);

    error FullRestakeDelegatorNotSupported(address vault);

    error UnknownDelegatorType(address vault, uint256 delegatorType);

    error SlasherNotSetForVault(address vault);

    error VetoSlasherNotSupported(address vault);

    error UnknownSlasherType(address vault, uint256 slasherType);

    error VaultNotRegistered(address vault);

    error VaultDeregRequestExists(address vault);

    error VaultDeregNotRequested(address vault);

    error ZeroAddressNotAllowed();

    error NetworkNotEntity(address network);

    error MissingOperatorRecord(address operator);

    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool);

    // TODO: remaining functions
}
