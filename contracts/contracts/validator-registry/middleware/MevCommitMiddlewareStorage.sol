// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

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

    /// @notice Enum TYPE for Symbiotic core OperatorSpecificDelegator.
    uint64 internal constant _OPERATOR_SPECIFIC_DELEGATOR_TYPE = 2;

    /// @notice Enum TYPE for Symbiotic core InstantSlasher.
    uint64 internal constant _INSTANT_SLASHER_TYPE = 0;

    /// @notice Enum TYPE for Symbiotic core VetoSlasher.
    uint64 internal constant _VETO_SLASHER_TYPE = 1;

    /// @notice Symbiotic core network registry.
    IRegistry public networkRegistry;

    /// @notice Symbiotic core operator registry.
    IRegistry public operatorRegistry;

    /// @notice Symbiotic core vault factory.
    IRegistry public vaultFactory;

    /// @notice Symbiotic core delegator factory.
    IRegistry public delegatorFactory;

    /// @notice Symbiotic core slasher factory.
    IRegistry public slasherFactory;

    /// @notice Symbiotic core burner router factory.
    IRegistry public burnerRouterFactory;

    /// @notice The network address, which must have registered with the NETWORK_REGISTRY.
    address public network;

    /// @dev A period in seconds during which the mev-commit oracle can invoke slashing.
    /// @notice This serves as the deregistration period for all of validator, operator, and vault records.
    /// @notice This also serves as the number of seconds that a registered Vault's epochDuration must be greater than.
    uint256 public slashPeriodSeconds;

    /// @notice Address of the mev-commit slash oracle.
    address public slashOracle;

    /// @notice Address of the mev-commit slash receiver.
    /// @dev This address should be immutable in practice, as it is used to validate every vault.
    /// @dev If this address is ever changed by the owner, all vaults would then need to update their burnerRouter. This is by-design.
    address public slashReceiver;

    /// @notice Minimum burner router delay.
    uint256 public minBurnerRouterDelay;

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
