// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";

abstract contract MevCommitMiddlewareStorage {

    /// @notice The only subnetwork ID for mev-commit middleware. Ie. mev-commit doesn't implement multiple subnets.
    uint96 internal constant _SUBNETWORK_ID = 1;

    /// @notice Enum TYPE for Symbiotic core NetworkRestakeDelegator.
    uint64 internal constant _NETWORK_RESTAKE_DELEGATOR_TYPE = 0;

    /// @notice Enum TYPE for Symbiotic core FullRestakeDelegator.
    uint64 internal constant _FULL_RESTAKE_DELEGATOR_TYPE = 1;

    /// @notice Enum TYPE for Symbiotic core InstantSlasher.
    uint64 internal constant _INSTANT_SLASHER_TYPE = 0;

    /// @notice Enum TYPE for Symbiotic core VetoSlasher.
    uint64 internal constant _VETO_SLASHER_TYPE = 1;

    /// @notice Minimum veto duration of 60 minutes for any vault.
    /// @dev This is enforced because veto duration is repurposed as the phase in which the oracle can feasibly call `executeSlash`,
    /// after initially requesting a slash.
    uint256 internal constant _MIN_VETO_DURATION = 60 minutes;

    /// @notice Symbiotic core network registry.
    IRegistry public networkRegistry;

    /// @notice Symbiotic core operator registry.
    IRegistry public operatorRegistry;

    /// @notice Symbiotic core vault factory.
    IRegistry public vaultFactory;

    /// @notice The network address, which must have registered with the NETWORK_REGISTRY.
    address public network;

    /// @dev A period in seconds during which the mev-commit oracle can invoke slashing.
    /// @notice This serves as the deregistration period for all of validator, operator, and vault records.
    /// @notice This also serves as the number of seconds that a registered Vault's epochDuration must be greater than.
    uint256 public slashPeriodSeconds;

    /// @notice Address of the mev-commit slash oracle.
    address public slashOracle;

    /// @notice Mapping of a validator's BLS public key to its validator record.
    mapping(bytes blsPubkey => IMevCommitMiddleware.ValidatorRecord) public validatorRecords;

    /// @notice Mapping of an operator's address to its operator record.
    mapping(address operatorAddress => IMevCommitMiddleware.OperatorRecord) public operatorRecords;

    /// @notice Mapping of a vault's address to its vault record.
    mapping(address vaultAddress => IMevCommitMiddleware.VaultRecord) public vaultRecords;

    /// @notice Mapping of a vault and operator to block number to slash record.
    mapping(address vault => mapping(address operator => mapping(uint256 blockNumber => IMevCommitMiddleware.SlashRecord))) public slashRecords;

    /// @notice Mapping of a vault to its representative operator, to a set of validator BLS public keys being secured
    /// by the vault.
    mapping(address vault => mapping(address operator => EnumerableSet.BytesSet)) internal _vaultAndOperatorToValset;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
