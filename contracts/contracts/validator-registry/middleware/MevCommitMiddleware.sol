// SPDX-License-Identifier: BSL 1.1

pragma solidity 0.8.26;

import {TimestampOccurrence} from "../../utils/Occurrence.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareStorage} from "./MevCommitMiddlewareStorage.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";
import {IVault} from "symbiotic-core/interfaces/vault/IVault.sol";
import {IVaultStorage} from "symbiotic-core/interfaces/vault/IVaultStorage.sol";
import {IBaseDelegator} from "symbiotic-core/interfaces/delegator/IBaseDelegator.sol";
import {IEntity} from "symbiotic-core/interfaces/common/IEntity.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Subnetwork} from "symbiotic-core/contracts/libraries/Subnetwork.sol";
import {ISlasher} from "symbiotic-core/interfaces/slasher/ISlasher.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {Errors} from "../../utils/Errors.sol";
import {IVetoSlasher} from "symbiotic-core/interfaces/slasher/IVetoSlasher.sol";

/// @notice This contracts serve as an entrypoint for L1 validators
/// to *opt-in* to mev-commit, ie. attest to the rules of mev-commit,
/// at the risk of funds being slashed. 
contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    using EnumerableSet for EnumerableSet.BytesSet;

    /// @notice Only the slash oracle account can call functions marked with this modifier.
    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, OnlySlashOracle(slashOracle));
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are valid length.
    modifier onlyValidBLSPubKeys(bytes[][] calldata blsPubKeys) {
        uint256 len = blsPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            bytes[] calldata innerArray = blsPubKeys[i];
            uint256 len2 = innerArray.length;
            for (uint256 j = 0; j < len2; ++j) {
                require(innerArray[j].length == 48, IMevCommitMiddleware.InvalidBLSPubKeyLength(
                    48, innerArray[j].length));
            }
        }
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the middleware contract.
    /// @param _networkRegistry Symbiotic core network registry contract.
    /// @param _operatorRegistry Symbiotic core operator registry contract.
    /// @param _vaultFactory Symbiotic core vault factory contract.
    /// @param _network Address of the mev-commit network EOA.
    /// @param _slashPeriodSeconds Oracle slashing must be invoked within `slashPeriodSeconds` of any event causing a validator to transition from *opted-in* to **not** *opted-in*.
    /// @param _slashOracle Address of the mev-commit oracle.
    /// @param _owner Contract owner address.
    function initialize(
        IRegistry _networkRegistry,
        IRegistry _operatorRegistry,
        IRegistry _vaultFactory,
        address _network,
        uint256 _slashPeriodSeconds,
        address _slashOracle,
        address _owner
    ) public initializer {
        _setNetworkRegistry(_networkRegistry);
        _setOperatorRegistry(_operatorRegistry);
        _setVaultFactory(_vaultFactory);
        _setNetwork(_network);
        _setSlashPeriodSeconds(_slashPeriodSeconds);
        _setSlashOracle(_slashOracle);
        __Pausable_init();
        __UUPSUpgradeable_init();
        require(_owner != address(0), ZeroAddressNotAllowed());
        __Ownable_init(_owner);
    }

    /// @dev Receive function to prevent unintended contract interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to prevent unintended contract interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Register operators, restricted to contract owner.
    /// @param operators Addresses of the operators to register.
    function registerOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _registerOperator(operators[i]);
        }
    }

    /// @notice Request operator deregistrations, restricted to contract owner.
    /// @param operators Addresses of the operators to request deregistrations for.
    function requestOperatorDeregistrations(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestOperatorDeregistration(operators[i]);
        }
    }

    /// @notice Deregisters operators, restricted to contract owner.
    /// @param operators Addresses of the operators to deregister.
    function deregisterOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterOperator(operators[i]);
        }
    }

    /// @notice Blacklists operators, restricted to contract owner.
    /// @param operators Addresses of the operators to blacklist.
    function blacklistOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _blacklistOperator(operators[i]);
        }
    }

    /// @notice Unblacklists operators, restricted to contract owner.
    /// @param operators Addresses of the operators to unblacklist.
    function unblacklistOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _unblacklistOperator(operators[i]);
        }
    }

    /// @notice Registers vaults, restricted to contract owner.
    /// @param vaults Addresses of the vaults to register.
    /// @param slashAmounts Corresponding slash amounts for each vault.
    function registerVaults(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, InvalidArrayLengths(vLen, slashAmounts.length));
        for (uint256 i = 0; i < vLen; ++i) {
            _registerVault(vaults[i], slashAmounts[i]);
        }
    }

    /// @notice Updates the slash amounts for vaults, restricted to contract owner.
    /// @param vaults Addresses of the vaults to update.
    /// @param slashAmounts Corresponding slash amounts for each vault.
    function updateSlashAmounts(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, InvalidArrayLengths(vLen, slashAmounts.length));
        for (uint256 i = 0; i < vLen; ++i) {
            _updateSlashAmount(vaults[i], slashAmounts[i]);
        }
    }

    /// @notice Requests vault deregistrations, restricted to contract owner.
    /// @param vaults Addresses of the vaults to request deregistrations for.
    function requestVaultDeregistrations(address[] calldata vaults) external onlyOwner {
        uint256 len = vaults.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestVaultDeregistration(vaults[i]);
        }
    }

    /// @notice Deregisters vaults, restricted to contract owner.
    /// @param vaults Addresses of the vaults to deregister.
    function deregisterVaults(address[] calldata vaults) external onlyOwner {
        uint256 len = vaults.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterVault(vaults[i]);
        }
    }

    /// @notice Registers validators via their BLS public key and vault which will secure them.
    /// @dev This function is callable by any delegated operator on behalf of a vault.
    /// @param blsPubkeys BLS public keys of the validators to register.
    /// @param vaults Addresses of vaults which will secure groups of validators.
    function registerValidators(bytes[][] calldata blsPubkeys,
        address[] calldata vaults) external whenNotPaused onlyValidBLSPubKeys(blsPubkeys) {
        uint256 vaultLen = vaults.length;
        require(vaultLen == blsPubkeys.length, InvalidArrayLengths(vaultLen, blsPubkeys.length));
        address operator = msg.sender;
        _checkOperator(operator);
        for (uint256 i = 0; i < vaultLen; ++i) {
            address vault = vaults[i];
            _checkVault(vault);
            uint256 potentialSlashableVals = _potentialSlashableVals(vault, operator);
            bytes[] calldata pubkeyArray = blsPubkeys[i];
            uint256 keyLen = pubkeyArray.length;
            require(keyLen <= potentialSlashableVals,
                ValidatorsNotSlashable(vault, operator, keyLen, potentialSlashableVals));
            for (uint256 j = 0; j < keyLen; ++j) {
                _addValRecord(pubkeyArray[j], vault, operator);
            }
        }
    }

    /// @notice Requests deregistrations for validators, restricted to contract owner,
    /// or the (still registered and non-blacklisted) operator of the validator pubkey.
    /// @param blsPubkeys BLS public keys of the validators to request deregistrations for.
    function requestValDeregistrations(bytes[] calldata blsPubkeys) external whenNotPaused {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestValDeregistration(blsPubkeys[i]);
        }
    }

    /// @dev Deletes validator records, callable by contract owner,
    /// or the (still registered and non-blacklisted) operator of the validator pubkey.
    /// @notice This function allows the contract owner to combat a greifing scenario where an operator
    /// registers a validator pubkey that it does not control, own, or otherwise manage.
    function deregisterValidators(bytes[] calldata blsPubkeys) external {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterValidator(blsPubkeys[i]);
        }
    }

    /// @dev Slashes validators and marks them for deregistration.
    /// @param blsPubkeys The L1 validator BLS public keys to slash.
    /// @param infractionTimestamps The block.timestamps for blocks during which each infraction occurred.
    function slashValidators(bytes[] calldata blsPubkeys, uint256[] calldata infractionTimestamps) external onlySlashOracle {
        uint256 len = blsPubkeys.length;
        require(len == infractionTimestamps.length, InvalidArrayLengths(len, infractionTimestamps.length));

        address[] memory swappedOperators = new address[](len);
        address[] memory swappedVaults = new address[](len);
        uint256[] memory newPositions = new uint256[](len);

        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubkey = blsPubkeys[i];
            // These and other checks in _slashValidator will succeed if current tx executes within
            // slashPeriodSeconds of validator being marked as "not opted-in",
            // OR relevant validator/vault/operator has not fully deregistered yet.
            ValidatorRecord storage valRecord = validatorRecords[pubkey];
            require(valRecord.exists, MissingValidatorRecord(pubkey));
            // Store slash record if it doesn't already exist. To ensure desirable ordering, _getNumSlashableVals should 
            // intentionally be computed once for the slash record, as collateral that is slashed later in this function
            // affects the metric.
            SlashRecord storage slashRecord = slashRecords[valRecord.vault][valRecord.operator][block.number];
            if (!slashRecord.exists) {
                uint256 numSlashableRegistered = _getNumSlashableRegisteredVals(valRecord.vault, valRecord.operator);
                require(numSlashableRegistered != 0, ValidatorsNotSlashable(valRecord.vault, valRecord.operator, len, numSlashableRegistered));
                slashRecords[valRecord.vault][valRecord.operator][block.number] = SlashRecord({
                    exists: true,
                    numSlashed: 0,
                    numInitSlashableRegistered: numSlashableRegistered
                });
            }
            // Swap about to be slashed pubkey with last slashable pubkey in valset.
            uint256 newPosition = slashRecord.numInitSlashableRegistered - slashRecord.numSlashed; // 1-indexed
            _vaultAndOperatorToValset[valRecord.vault][valRecord.operator].swapWithPosition(pubkey, newPosition);
            swappedVaults[i] = valRecord.vault;
            swappedOperators[i] = valRecord.operator;
            newPositions[i] = newPosition;

            ++slashRecord.numSlashed;
            _slashValidator(blsPubkeys[i], infractionTimestamps[i], valRecord);
        }
        emit ValidatorPositionsSwapped(blsPubkeys, swappedVaults, swappedOperators, newPositions);
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner { _pause(); }

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external onlyOwner { _unpause(); }

    /// @dev Sets the network registry, restricted to contract owner.
    function setNetworkRegistry(IRegistry _networkRegistry) external onlyOwner {
        _setNetworkRegistry(_networkRegistry);
    }

    /// @dev Sets the operator registry, restricted to contract owner.
    function setOperatorRegistry(IRegistry _operatorRegistry) external onlyOwner {
        _setOperatorRegistry(_operatorRegistry);
    }

    /// @dev Sets the vault factory, restricted to contract owner.
    function setVaultFactory(IRegistry _vaultFactory) external onlyOwner {
        _setVaultFactory(_vaultFactory);
    }

    /// @dev Sets the network address, restricted to contract owner.
    function setNetwork(address _network) external onlyOwner {
        _setNetwork(_network);
    }

    /// @dev Sets the slash period in seconds, restricted to contract owner.
    function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) external onlyOwner {
        _setSlashPeriodSeconds(slashPeriodSeconds_);
    }

    /// @dev Sets the slash oracle, restricted to contract owner.
    function setSlashOracle(address slashOracle_) external onlyOwner {
        _setSlashOracle(slashOracle_);
    }

    /// @notice Queries if a validator is opted-in to mev-commit through a vault.
    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    /// @notice Queries if a validator is slashable.
    function isValidatorSlashable(bytes calldata blsPubkey) external view returns (bool) {
        ValidatorRecord storage record = validatorRecords[blsPubkey];
        require(record.exists, MissingValRecord(blsPubkey));
        _checkVault(record.vault);
        _checkOperator(record.operator);
        return _isValidatorSlashable(blsPubkey, record.vault, record.operator);
    }

    /// @return Number of potential new validators that could be registered and be slashable.
    function potentialSlashableValidators(address vault, address operator) external view returns (uint256) {
        return _potentialSlashableVals(vault, operator);
    }

    /// @notice Queries the one-indexed position of a validator's BLS pubkey in its valset.
    /// @return 0 if the blsPubkey is not in the valset.
    function getPositionInValset(bytes calldata blsPubkey, address vault, address operator) external view returns (uint256) {
        return _getPositionInValset(blsPubkey, vault, operator);
    }

    /// @return Number of validators that could be slashable according to vault stake.
    function getNumSlashableVals(address vault, address operator) external view returns (uint256) {
        return _getNumSlashableVals(vault, operator);
    }

    /// @notice Queries the BLS pubkey at a given one-indexed position in the valset for a vault and operator.
    /// @return An empty bytes array if the index is out of bounds or the valset is empty.
    function pubkeyAtPositionInValset(uint256 index, address vault, address operator) external view returns (bytes memory) {
        if (index == 0 || _vaultAndOperatorToValset[vault][operator].length() < index) {
            return new bytes(0);
        }
        return _vaultAndOperatorToValset[vault][operator].at(index - 1);
    }

    /// @return Length of the valset for a given vault and operator.
    function valsetLength(address vault, address operator) external view returns (uint256) {
        return _vaultAndOperatorToValset[vault][operator].length();
    }

    function _setOperatorRecord(address operator) internal {
        operatorRecords[operator] = OperatorRecord({
            exists: true,
            deregRequestOccurrence: TimestampOccurrence.Occurrence({
                exists: false,
                timestamp: 0
            }),
            isBlacklisted: false
        });
    }

    function _registerOperator(address operator) internal {
        require(!operatorRecords[operator].exists, OperatorAlreadyRegistered(operator));
        require(operatorRegistry.isEntity(operator), OperatorNotEntity(operator));
        _setOperatorRecord(operator);
        emit OperatorRegistered(operator);
    }

    function _requestOperatorDeregistration(address operator) internal {
        OperatorRecord storage record = operatorRecords[operator];
        require(record.exists, OperatorNotRegistered(operator));
        require(!record.isBlacklisted, OperatorIsBlacklisted(operator));
        require(!record.deregRequestOccurrence.exists, OperatorDeregRequestExists(operator));
        TimestampOccurrence.captureOccurrence(record.deregRequestOccurrence);
        emit OperatorDeregistrationRequested(operator);
    }

    function _deregisterOperator(address operator) internal {
        OperatorRecord storage record = operatorRecords[operator];
        require(record.exists, OperatorNotRegistered(operator));
        require(_isOperatorReadyToDeregister(operator), OperatorNotReadyToDeregister(
            operator, block.timestamp, record.deregRequestOccurrence.timestamp));
        require(!record.isBlacklisted, OperatorIsBlacklisted(operator));
        delete operatorRecords[operator];
        emit OperatorDeregistered(operator);
    }

    function _blacklistOperator(address operator) internal {
        OperatorRecord storage record = operatorRecords[operator];
        if (!record.exists) {
            _setOperatorRecord(operator);
        }
        require(!record.isBlacklisted, OperatorAlreadyBlacklisted(operator));
        record.isBlacklisted = true;
        emit OperatorBlacklisted(operator);
    }

    function _unblacklistOperator(address operator) internal {
        OperatorRecord storage record = operatorRecords[operator];
        require(record.exists, OperatorNotRegistered(operator));
        require(record.isBlacklisted, OperatorNotBlacklisted(operator));
        record.isBlacklisted = false;
        emit OperatorUnblacklisted(operator);
    }

    function _setVaultRecord(address vault, uint256 slashAmount) internal {
        vaultRecords[vault] = VaultRecord({
            exists: true,
            deregRequestOccurrence: TimestampOccurrence.Occurrence({
                exists: false,
                timestamp: 0
            }),
            slashAmount: slashAmount
        });
    }

    function _registerVault(address vault, uint256 slashAmount) internal {
        require(!vaultRecords[vault].exists, VaultAlreadyRegistered(vault));
        require(vaultFactory.isEntity(vault), VaultNotEntity(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));

        IEntity delegator = IEntity(IVault(vault).delegator());
        if (delegator.TYPE() == _FULL_RESTAKE_DELEGATOR_TYPE) {
            revert FullRestakeDelegatorNotSupported(vault);
        } else if (delegator.TYPE() != _NETWORK_RESTAKE_DELEGATOR_TYPE) {
            revert UnknownDelegatorType(vault, delegator.TYPE());
        }

        IVaultStorage vaultContract = IVaultStorage(vault);
        uint256 vaultEpochDurationSeconds = vaultContract.epochDuration();

        address slasher = IVault(vault).slasher();
        require(slasher != address(0), SlasherNotSetForVault(vault));
        uint256 slasherType = IEntity(slasher).TYPE();
        if (slasherType == _VETO_SLASHER_TYPE) {
            IVetoSlasher vetoSlasher = IVetoSlasher(slasher);
            // For veto slashers, incorporate that veto duration will eat into vault's epoch duration.
            /// @dev vetoDuration must be less than epochDuration as enforced by VetoSlasher.sol.
            vaultEpochDurationSeconds -= vetoSlasher.vetoDuration();
            require(vetoSlasher.resolver(_getSubnetwork(), new bytes(0)) == address(0),
                VetoSlasherMustHaveZeroResolver(vault));
        } else if (slasherType != _INSTANT_SLASHER_TYPE) {
            revert UnknownSlasherType(vault, slasherType);
        }

        require(vaultEpochDurationSeconds > slashPeriodSeconds,
            InvalidVaultEpochDuration(vault, vaultEpochDurationSeconds, slashPeriodSeconds));

        _setVaultRecord(vault, slashAmount);
        emit VaultRegistered(vault, slashAmount);
    }

    function _updateSlashAmount(address vault, uint256 slashAmount) internal {
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));
        record.slashAmount = slashAmount;
        emit VaultSlashAmountUpdated(vault, slashAmount);
    }

    function _requestVaultDeregistration(address vault) internal {
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        require(!record.deregRequestOccurrence.exists, VaultDeregRequestExists(vault));
        TimestampOccurrence.captureOccurrence(record.deregRequestOccurrence);
        emit VaultDeregistrationRequested(vault);
    }

    function _deregisterVault(address vault) internal {
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        require(_isVaultReadyToDeregister(vault), VaultNotReadyToDeregister(vault, block.timestamp,
            record.deregRequestOccurrence.timestamp));
        delete vaultRecords[vault];
        emit VaultDeregistered(vault);
    }

    function _addValRecord(bytes calldata blsPubkey, address vault, address operator) internal {
        require(!validatorRecords[blsPubkey].exists, ValidatorRecordAlreadyExists(blsPubkey));
        validatorRecords[blsPubkey] = ValidatorRecord({
            exists: true,
            deregRequestOccurrence: TimestampOccurrence.Occurrence({
                exists: false,
                timestamp: 0
            }),
            vault: vault,
            operator: operator
        });
        _vaultAndOperatorToValset[vault][operator].add(blsPubkey);
        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        emit ValRecordAdded(blsPubkey, operator, vault, position);
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        ValidatorRecord storage record = validatorRecords[blsPubkey];
        require(record.exists, MissingValidatorRecord(blsPubkey));
        require(!record.deregRequestOccurrence.exists, ValidatorDeregRequestExists(blsPubkey));
        if (msg.sender != owner()) {
            _checkCallingOperator(record.operator);
        }
        TimestampOccurrence.captureOccurrence(record.deregRequestOccurrence);
        uint256 position = _getPositionInValset(blsPubkey, record.vault, record.operator);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, position);
    }

    function _deregisterValidator(bytes calldata blsPubkey) internal {
        ValidatorRecord storage record = validatorRecords[blsPubkey];
        require(record.exists, MissingValidatorRecord(blsPubkey));
        require(_isValidatorReadyToDeregister(blsPubkey), ValidatorNotReadyToDeregister(
            blsPubkey, block.timestamp, record.deregRequestOccurrence.timestamp));
        if (msg.sender != owner()) {
            _checkCallingOperator(record.operator);
        }
        bool removed = _vaultAndOperatorToValset[record.vault][record.operator].remove(blsPubkey);
        require(removed, ValidatorNotRemovedFromValset(blsPubkey, record.vault, record.operator));
        delete validatorRecords[blsPubkey];
        emit ValRecordDeleted(blsPubkey, msg.sender);
    }

    /// @dev Slashes a validator and marks it for deregistration.
    /// @param blsPubkey The L1 validator BLS public key to slash.
    /// @param infractionTimestamp The block.timestamp for the block during which the infraction occurred.
    /// @dev Operator and vault are not deregistered for the validator's infraction,
    /// so as to avoid opting-out large groups of validators at once.
    function _slashValidator(bytes calldata blsPubkey, uint256 infractionTimestamp, ValidatorRecord storage valRecord) internal {
        VaultRecord storage vaultRecord = vaultRecords[valRecord.vault];
        require(vaultRecord.exists, MissingVaultRecord(valRecord.vault));
        OperatorRecord storage operatorRecord = operatorRecords[valRecord.operator];
        require(operatorRecord.exists, MissingOperatorRecord(valRecord.operator));

        require(infractionTimestamp != 0, InfractionTimestampMustBeNonZero());

        require(_isValidatorSlashable(blsPubkey, valRecord.vault, valRecord.operator),
            ValidatorNotSlashable(blsPubkey, valRecord.vault, valRecord.operator));

        // Slash amount is enforced as non-zero in _registerVault.
        uint256 amount = vaultRecord.slashAmount;

        address slasher = IVault(valRecord.vault).slasher();
        uint256 slasherType = IEntity(slasher).TYPE();
        if (slasherType == _VETO_SLASHER_TYPE) {
            uint256 slashIndex = IVetoSlasher(slasher).requestSlash(
                _getSubnetwork(), valRecord.operator, amount, SafeCast.toUint48(infractionTimestamp), new bytes(0));
            emit ValidatorSlashRequested(blsPubkey, valRecord.operator, valRecord.vault, slashIndex);
        } else if (slasherType == _INSTANT_SLASHER_TYPE) {
            uint256 slashedAmount = ISlasher(slasher).slash(
                _getSubnetwork(), valRecord.operator, amount, SafeCast.toUint48(infractionTimestamp), new bytes(0));
            emit ValidatorSlashed(blsPubkey, valRecord.operator, valRecord.vault, slashedAmount);
        }

        // If validator has not already requested deregistration,
        // do so to mark them as no longer opted-in.
        if (!valRecord.deregRequestOccurrence.exists) {
            TimestampOccurrence.captureOccurrence(valRecord.deregRequestOccurrence);
        }
    }

    /// @dev Internal function to set the network registry.
    function _setNetworkRegistry(IRegistry _networkRegistry) internal {
        require(_networkRegistry != IRegistry(address(0)), ZeroAddressNotAllowed());
        networkRegistry = _networkRegistry;
        emit NetworkRegistrySet(address(_networkRegistry));
    }

    /// @dev Internal function to set the operator registry.
    function _setOperatorRegistry(IRegistry _operatorRegistry) internal {
        require(_operatorRegistry != IRegistry(address(0)), ZeroAddressNotAllowed());
        operatorRegistry = _operatorRegistry;
        emit OperatorRegistrySet(address(_operatorRegistry));
    }

    /// @dev Internal function to set the vault factory.
    function _setVaultFactory(IRegistry _vaultFactory) internal {
        require(_vaultFactory != IRegistry(address(0)), ZeroAddressNotAllowed());
        vaultFactory = _vaultFactory;
        emit VaultFactorySet(address(_vaultFactory));
    }

    /// @dev Internal function to set the network address, which must have registered with the NETWORK_REGISTRY.
    function _setNetwork(address _network) internal {
        require(_network != address(0), ZeroAddressNotAllowed());
        require(networkRegistry.isEntity(_network), NetworkNotEntity(_network));
        network = _network;
        emit NetworkSet(_network);
    }

    /// @dev Internal function to set the slash period in seconds.
    function _setSlashPeriodSeconds(uint256 slashPeriodSeconds_) internal {
        require(slashPeriodSeconds_ != 0, ZeroUintNotAllowed());
        slashPeriodSeconds = slashPeriodSeconds_;
        emit SlashPeriodSecondsSet(slashPeriodSeconds_);
    }

    /// @dev Internal function to set the slash oracle.
    function _setSlashOracle(address slashOracle_) internal {
        require(slashOracle_ != address(0), ZeroAddressNotAllowed());
        slashOracle = slashOracle_;
        emit SlashOracleSet(slashOracle_);
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _checkOperator(address operator) internal view {
        require(operatorRegistry.isEntity(operator), OperatorNotEntity(operator));
        OperatorRecord storage record = operatorRecords[operator];
        require(record.exists, OperatorNotRegistered(operator));
        require(!record.deregRequestOccurrence.exists, OperatorDeregRequestExists(operator));
        require(!record.isBlacklisted, OperatorIsBlacklisted(operator));
    }

    function _checkCallingOperator(address operator) internal view {
        require(msg.sender == operator, OnlyOperator(operator));
        _checkOperator(operator);
    }

    function _checkVault(address vault) internal view {
        require(vaultFactory.isEntity(vault), VaultNotEntity(vault));
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        require(!record.deregRequestOccurrence.exists, VaultDeregRequestExists(vault));
    }

    /// @dev Returns the one-indexed position of the blsPubkey in the set.
    /// @return 0 if the blsPubkey is not in the set.
    function _getPositionInValset(bytes calldata blsPubkey,
        address vault, address operator) internal view returns (uint256) {
        return _vaultAndOperatorToValset[vault][operator].position(blsPubkey);
    }

    function _isValidatorReadyToDeregister(bytes calldata blsPubkey) internal view returns (bool) {
        ValidatorRecord storage record = validatorRecords[blsPubkey];
        return record.deregRequestOccurrence.exists && 
            block.timestamp > slashPeriodSeconds + record.deregRequestOccurrence.timestamp;
    }

    function _isOperatorReadyToDeregister(address operator) internal view returns (bool) {
        OperatorRecord storage record = operatorRecords[operator];
        return record.deregRequestOccurrence.exists && 
            block.timestamp > slashPeriodSeconds + record.deregRequestOccurrence.timestamp;
    }

    function _isVaultReadyToDeregister(address vault) internal view returns (bool) {
        VaultRecord storage record = vaultRecords[vault];
        return record.deregRequestOccurrence.exists && 
            block.timestamp > slashPeriodSeconds + record.deregRequestOccurrence.timestamp;
    }

    function _getSubnetwork() internal view returns (bytes32) {
        return Subnetwork.subnetwork(network, _SUBNETWORK_ID);
    }

    /// @return Number of validators that could be slashable according to vault stake.
    function _getNumSlashableVals(address vault, address operator) internal view returns (uint256) {
        IBaseDelegator delegator = IBaseDelegator(IVault(vault).delegator());
        uint256 allocatedStake = delegator.stake(_getSubnetwork(), operator);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        return allocatedStake / slashAmount;
    }

    /// @return Number of validators that are both slashable and registered.
    function _getNumSlashableRegisteredVals(address vault, address operator) internal view returns (uint256) {
        uint256 slashableVals = _getNumSlashableVals(vault, operator);
        uint256 numRegistered = _vaultAndOperatorToValset[vault][operator].length();
        return slashableVals < numRegistered ? slashableVals : numRegistered;
    }

    function _isValidatorSlashable(bytes calldata blsPubkey, address vault, address operator) internal view returns (bool) {
        uint256 slashableVals = _getNumSlashableVals(vault, operator);
        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        require(position != 0, ValidatorNotInValset(blsPubkey, vault, operator));
        return position <= slashableVals; // position is 1-indexed
    }

    /// @return Number of potential new validators that could be registered and be slashable.
    function _potentialSlashableVals(address vault, address operator) internal view returns (uint256) {
        uint256 slashableVals = _getNumSlashableVals(vault, operator);
        uint256 numRegistered = _vaultAndOperatorToValset[vault][operator].length();
        if (slashableVals < numRegistered) {
            return 0;
        }
        return slashableVals - numRegistered;
    }
    
    function _isValidatorOptedIn(bytes calldata blsPubkey) internal view returns (bool) {
        ValidatorRecord storage valRecord = validatorRecords[blsPubkey];
        if (!valRecord.exists) {
            return false;
        }
        if (valRecord.deregRequestOccurrence.exists) {
            return false;
        }
        VaultRecord storage vaultRecord = vaultRecords[valRecord.vault];
        if (!vaultRecord.exists) {
            return false;
        }
        if (vaultRecord.deregRequestOccurrence.exists) {
            return false;
        }
        if (!vaultFactory.isEntity(valRecord.vault)) {
            return false;
        }
        OperatorRecord storage operatorRecord = operatorRecords[valRecord.operator];
        if (!operatorRecord.exists) {
            return false;
        }
        if (operatorRecord.deregRequestOccurrence.exists) {
            return false;
        }
        if (operatorRecord.isBlacklisted) {
            return false;
        }
        if (!operatorRegistry.isEntity(valRecord.operator)) {
            return false;
        }
        if (!_isValidatorSlashable(blsPubkey, valRecord.vault, valRecord.operator)) {
            return false;
        }
        return true;
    }
}
