// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {TimestampOccurrence} from "../utils/Occurrence.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Checkpoints} from "@openzeppelin/contracts/utils/structs/Checkpoints.sol";

interface IMevCommitMiddleware {

    /// @notice Struct representing a registered operator.
    struct OperatorRecord {
        /// @notice A possible occurrence of a deregistration request.
        TimestampOccurrence.Occurrence deregRequestOccurrence;
        /// @notice Whether this operator record exists.
        bool exists;
        /// @notice Whether this operator is blacklisted.
        bool isBlacklisted;
    }

    /// @notice Struct representing a registered vault.
    struct VaultRecord {
        /// @notice Whether this vault record exists.
        bool exists;
        /// @notice A possible occurrence of a deregistration request.
        TimestampOccurrence.Occurrence deregRequestOccurrence;
        /// @notice The slash amount per validator, relevant to this vault.
        Checkpoints.Trace160 slashAmountHistory;
    }

    /// @notice Struct representing a registered validator.
    struct ValidatorRecord {
        /// @notice The vault holding slashable stake which represents the validator.
        address vault;
        /// @notice The operator which registered this validator pubkey with a vault.
        address operator;
        /// @notice Whether this validator record exists.
        bool exists;
        /// @notice A possible occurrence of a deregistration request.
        TimestampOccurrence.Occurrence deregRequestOccurrence;
    }

    struct SlashRecord {
        /// @notice Whether this slash record exists.
        bool exists;
        /// @notice The number of validators slashed for this vault and operator at the current block.
        uint256 numSlashed;
        /// @notice The number of validators that are registered for this vault and operator at the current block.
        /// @dev This is computed once upon slash record creation to ensure desirable valset ordering.
        uint256 numRegistered;
    }

    /// @notice Emmitted when an operator is registered
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when an operator requests deregistration
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when an operator is blacklisted
    event OperatorBlacklisted(address indexed operator);

    /// @notice Emmitted when an operator is unblacklisted
    event OperatorUnblacklisted(address indexed operator);

    /// @notice Emmitted when a vault record is added
    event VaultRegistered(address indexed vault, uint160 slashAmount);

    /// @notice Emmitted when a vault slash amount is updated
    event VaultSlashAmountUpdated(address indexed vault, uint160 slashAmount);

    /// @notice Emmitted when a vault deregistration is requested
    event VaultDeregistrationRequested(address indexed vault);

    /// @notice Emmitted when a vault is deregistered
    event VaultDeregistered(address indexed vault);

    /// @notice Emmitted when a validator record is added to state
    /// @dev The position is one-indexed.
    event ValRecordAdded(bytes blsPubkey, address indexed operator, address indexed vault, uint256 indexed position);

    /// @notice Emmitted when validator deregistration is requested
    event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position);

    /// @notice Emmitted when a validator record is deleted by the contract owner
    event ValRecordDeleted(bytes blsPubkey, address indexed msgSender);

    /// @notice Emmitted when a validator is slashed from an instant slasher
    event ValidatorSlashed(bytes blsPubkey, address indexed operator, address indexed vault, uint256 slashedAmount);

    /// @notice Emmitted when the network registry is set
    event NetworkRegistrySet(address networkRegistry);

    /// @notice Emmitted when the operator registry is set
    event OperatorRegistrySet(address operatorRegistry);

    /// @notice Emmitted when the vault factory is set
    event VaultFactorySet(address vaultFactory);

    /// @notice Emmitted when the burner router factory is set
    event BurnerRouterFactorySet(address burnerRouterFactory);

    /// @notice Emmitted when the network is set
    event NetworkSet(address network);

    /// @notice Emmitted when the slash period in seconds is set
    event SlashPeriodSecondsSet(uint256 slashPeriodSeconds);

    /// @notice Emmitted when the slash period in blocks is set
    event SlashPeriodBlocksSet(uint256 slashPeriodBlocks);

    /// @notice Emmitted when the slash oracle is set
    event SlashOracleSet(address slashOracle);

    /// @notice Emmitted when the slash receiver is set
    event SlashReceiverSet(address slashReceiver);

    /// @notice Emmitted when the minimum burner router delay is set
    event MinBurnerRouterDelaySet(uint256 minBurnerRouterDelay);

    /// @notice Emmitted when validator positions are swapped as a part of slashing
    /// @dev Each array index corresponds to a swap instance. ie. all lists should be of equal length.
    event ValidatorPositionsSwapped(bytes[] blsPubkeys, address[] vaults, address[] operators, uint256[] newPositions);

    error OnlySlashOracle(address slashOracle);

    error OnlyOperator(address operator);

    error InvalidArrayLengths(uint256 vaultLen, uint256 pubkeyLen);

    error ValidatorsNotSlashable(address vault, address operator, uint256 numRequested, uint256 potentialSlashableVals);

    error NoRegisteredValidators(address vault, address operator);

    error MissingValRecord(bytes blsPubkey);

    error OperatorAlreadyRegistered(address operator);

    error OperatorNotEntity(address operator);

    error OperatorNotRegistered(address operator);

    error OperatorDeregRequestExists(address operator);

    error OperatorIsBlacklisted(address operator);

    error OperatorNotReadyToDeregister(address operator, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error OperatorAlreadyBlacklisted(address operator);

    error OperatorNotBlacklisted(address operator);

    error ValidatorRecordAlreadyExists(bytes blsPubkey);

    error MissingValidatorRecord(bytes blsPubkey);

    error ValidatorDeregRequestExists(bytes blsPubkey);

    error ValidatorNotReadyToDeregister(bytes blsPubkey, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error VaultAlreadyRegistered(address vault);

    error MissingVaultRecord(address vault);

    error VaultNotEntity(address vault);

    error VaultNotReadyToDeregister(address vault, uint256 currentTimestamp, uint256 deregRequestTimestamp);

    error FailedToAddValidatorToValset(bytes blsPubkey, address vault, address operator);

    error SlashAmountMustBeNonZero(address vault);

    error InvalidVaultEpochDuration(address vault, uint256 vaultEpochDurationSec, uint256 slashPeriodSec);

    error FullRestakeDelegatorNotSupported(address vault);

    error UnknownDelegatorType(address vault, uint256 delegatorType);

    error SlasherNotSetForVault(address vault);

    error VetoSlasherMustHaveZeroResolver(address vault);

    error VetoDurationTooShort(address vault, uint256 vetoDuration);

    error UnknownSlasherType(address vault, uint256 slasherType);

    error OnlyVetoSlashersRequireExecution(address vault, uint256 slasherType);

    error VaultNotRegistered(address vault);

    error VaultDeregRequestExists(address vault);

    error InvalidVaultBurner(address vault);

    error InvalidVaultBurnerConsideringOperator(address vault, address operator);

    error ValidatorNotInValset(bytes blsPubkey, address vault, address operator);

    error NoSlashAmountAtTimestamp(address vault, uint256 timestamp);

    error FutureTimestampDisallowed(address vault, uint256 timestamp);

    error VaultDeregNotRequested(address vault);

    error ZeroAddressNotAllowed();

    error NetworkNotEntity(address network);

    error ZeroUintNotAllowed();

    error MissingOperatorRecord(address operator);

    error InvalidBLSPubKeyLength(uint256 expectedLength, uint256 actualLength);

    error CaptureTimestampMustBeNonZero();

    error ValidatorNotRemovedFromValset(bytes blsPubkey, address vault, address operator);

    error ValidatorNotSlashable(bytes blsPubkey, address vault, address operator);

    /// @notice Registers multiple operators.
    function registerOperators(address[] calldata operators) external;

    /// @notice Requests deregistration for multiple operators.
    function requestOperatorDeregistrations(address[] calldata operators) external;

    /// @notice Deregisters multiple operators.
    function deregisterOperators(address[] calldata operators) external;

    /// @notice Blacklists multiple operators.
    function blacklistOperators(address[] calldata operators) external;

    /// @notice Unblacklists multiple operators.
    function unblacklistOperators(address[] calldata operators) external;

    /// @notice Registers multiple vaults with corresponding slash amounts.
    function registerVaults(address[] calldata vaults, uint160[] calldata slashAmounts) external;

    /// @notice Updates slash amounts for multiple vaults.
    function updateSlashAmounts(address[] calldata vaults, uint160[] calldata slashAmounts) external;

    /// @notice Requests deregistration for multiple vaults.
    function requestVaultDeregistrations(address[] calldata vaults) external;

    /// @notice Deregisters multiple vaults.
    function deregisterVaults(address[] calldata vaults) external;

    /// @notice Registers validators via their BLS public key and vault which will secure them.
    function registerValidators(bytes[][] calldata blsPubkeys, address[] calldata vaults) external;

    /// @notice Requests deregistration for multiple validators.
    function requestValDeregistrations(bytes[] calldata blsPubkeys) external;

    /// @notice Deregisters multiple validators.
    function deregisterValidators(bytes[] calldata blsPubkeys) external;

    /// @notice Slashes multiple validators with their respective infraction timestamps.
    function slashValidators(bytes[] calldata blsPubkeys, uint256[] calldata infractionTimestamps) external;

    /// @notice Pauses the contract.
    function pause() external;

    /// @notice Unpauses the contract.
    function unpause() external;

    /// @notice Sets the network registry.
    function setNetworkRegistry(IRegistry _networkRegistry) external;

    /// @notice Sets the operator registry.
    function setOperatorRegistry(IRegistry _operatorRegistry) external;

    /// @notice Sets the vault factory.
    function setVaultFactory(IRegistry _vaultFactory) external;

    /// @notice Sets the network address.
    function setNetwork(address _network) external;

    /// @notice Sets the slash period in seconds.
    function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) external;

    /// @notice Sets the slash oracle address.
    function setSlashOracle(address slashOracle_) external;

    /// @notice Checks if a validator is opted in.
    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool);

    /// @notice Checks if a validator is slashable.
    function isValidatorSlashable(bytes calldata blsPubkey) external view returns (bool);

    /// @notice Returns the potential number of slashable validators for a given vault and operator.
    function potentialSlashableValidators(address vault, address operator) external view returns (uint256);

    /// @notice Returns the one-indexed position of a blsPubkey in its valset.
    /// @param blsPubkey The BLS public key of the validator.
    /// @param vault The address of the vault.
    /// @param operator The address of the operator.
    /// @return The position in the valset or 0 if not present.
    function getPositionInValset(bytes calldata blsPubkey, address vault, address operator) external view returns (uint256);

    /// @return Number of validators that could be slashable according to vault stake.
    function getNumSlashableVals(address vault, address operator) external view returns (uint256);

    /// @notice Queries the BLS pubkey at a given one-indexed position in the valset for a vault and operator.
    /// @return An empty bytes array if the index is out of bounds or the valset is empty.
    function pubkeyAtPositionInValset(uint256 index, address vault, address operator) external view returns (bytes memory);

    /// @return Length of the valset for a given vault and operator.
    function valsetLength(address vault, address operator) external view returns (uint256);
}
