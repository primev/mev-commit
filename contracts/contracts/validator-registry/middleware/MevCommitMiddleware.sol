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
import {Checkpoints} from "@openzeppelin/contracts/utils/structs/Checkpoints.sol";
import {IBurnerRouter} from "symbiotic-burners/interfaces/router/IBurnerRouter.sol";
import {IBaseSlasher} from "symbiotic-core/interfaces/slasher/IBaseSlasher.sol";

/// @notice This contracts serve as an entrypoint for L1 validators
/// to *opt-in* to mev-commit, ie. attest to the rules of mev-commit,
/// at the risk of funds being slashed. 
contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    using EnumerableSet for EnumerableSet.BytesSet;
    using Checkpoints for Checkpoints.Trace160;

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
    /// @param _burnerRouterFactory Symbiotic core burner router factory contract.
    /// @param _network Address of the mev-commit network EOA.
    /// @param _slashPeriodSeconds Oracle slashing must be invoked within `slashPeriodSeconds` of any event causing a validator to transition from *opted-in* to **not** *opted-in*.
    /// @param _slashOracle Address of the mev-commit oracle.
    /// @param _slashReceiver Address of the mev-commit slash receiver.
    /// @param _minBurnerRouterDelay Minimum burner router delay.
    /// @param _owner Contract owner address.
    function initialize(
        IRegistry _networkRegistry,
        IRegistry _operatorRegistry,
        IRegistry _vaultFactory,
        IRegistry _delegatorFactory,
        IRegistry _slasherFactory,
        IRegistry _burnerRouterFactory,
        address _network,
        uint256 _slashPeriodSeconds,
        address _slashOracle,
        address _slashReceiver,
        uint256 _minBurnerRouterDelay,
        address _owner
    ) public initializer {
        _setNetworkRegistry(_networkRegistry);
        _setOperatorRegistry(_operatorRegistry);
        _setVaultFactory(_vaultFactory);
        _setDelegatorFactory(_delegatorFactory);
        _setSlasherFactory(_slasherFactory);
        _setBurnerRouterFactory(_burnerRouterFactory);
        _setNetwork(_network);
        _setSlashPeriodSeconds(_slashPeriodSeconds);
        _setSlashOracle(_slashOracle);
        _setSlashReceiver(_slashReceiver);
        _setMinBurnerRouterDelay(_minBurnerRouterDelay);
        __Pausable_init();
        __UUPSUpgradeable_init();
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
    function registerVaults(address[] calldata vaults, uint160[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, InvalidArrayLengths(vLen, slashAmounts.length));
        for (uint256 i = 0; i < vLen; ++i) {
            _registerVault(vaults[i], slashAmounts[i]);
        }
    }

    /// @notice Updates the slash amounts for vaults, restricted to contract owner.
    /// @param vaults Addresses of the vaults to update.
    /// @param slashAmounts Corresponding slash amounts for each vault.
    function updateSlashAmounts(address[] calldata vaults, uint160[] calldata slashAmounts) external onlyOwner {
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
            require(_validateVaultBurnerAgainstOperator(vault, operator), InvalidVaultBurnerConsideringOperator(vault, operator));
            uint256 potentialSlashableVals = _potentialSlashableVals(vault, operator);
            bytes[] calldata pubkeyArray = blsPubkeys[i];
            uint256 numKeys = pubkeyArray.length;
            // This check exists for UX, in that the vault should have enough collateral staked prior to validator registration.
            require(numKeys <= potentialSlashableVals, ValidatorsNotSlashable(vault, operator, numKeys, potentialSlashableVals));
            for (uint256 j = 0; j < numKeys; ++j) {
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
    function deregisterValidators(bytes[] calldata blsPubkeys) external whenNotPaused {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterValidator(blsPubkeys[i]);
        }
    }

    /// @dev Slashes validators and marks them for deregistration.
    /// @param blsPubkeys The L1 validator BLS public keys to slash.
    /// @param captureTimestamps block.timestamps of the latest finalized block that the blsPubkey was queried as "opted-in" by the oracle.
    function slashValidators(bytes[] calldata blsPubkeys, uint256[] calldata captureTimestamps) external onlySlashOracle whenNotPaused {
        uint256 len = blsPubkeys.length;
        require(len == captureTimestamps.length, InvalidArrayLengths(len, captureTimestamps.length));

        address[] memory swappedOperators = new address[](len);
        address[] memory swappedVaults = new address[](len);
        uint256[] memory newPositions = new uint256[](len);

        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubkey = blsPubkeys[i];
            // These and other checks in _slashValidator are guaranteed to succeed if current tx executes within
            // slashPeriodSeconds of a validator's captureTimestamp (defined in README.md).
            ValidatorRecord storage valRecord = validatorRecords[pubkey];
            require(valRecord.exists, MissingValidatorRecord(pubkey));

            SlashRecord storage slashRecord = slashRecords[valRecord.vault][valRecord.operator][block.number];
            if (!slashRecord.exists) {
                uint256 numRegistered = _vaultAndOperatorToValset[valRecord.vault][valRecord.operator].length();
                require(numRegistered != 0, NoRegisteredValidators(valRecord.vault, valRecord.operator));
                slashRecords[valRecord.vault][valRecord.operator][block.number] = SlashRecord({
                    exists: true,
                    numSlashed: 0,
                    numRegistered: numRegistered
                });
            }
            // Swap about to be slashed pubkey with last pubkey in registered valset.
            uint256 newPosition = slashRecord.numRegistered - slashRecord.numSlashed; // 1-indexed
            _vaultAndOperatorToValset[valRecord.vault][valRecord.operator].swapWithPosition(pubkey, newPosition);
            swappedVaults[i] = valRecord.vault;
            swappedOperators[i] = valRecord.operator;
            newPositions[i] = newPosition;

            ++slashRecord.numSlashed;
            _slashValidator(blsPubkeys[i], captureTimestamps[i], valRecord);
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

    /// @dev Sets the delegator factory, restricted to contract owner.
    function setDelegatorFactory(IRegistry _delegatorFactory) external onlyOwner {
        _setDelegatorFactory(_delegatorFactory);
    }

    /// @dev Sets the slasher factory, restricted to contract owner.
    function setSlasherFactory(IRegistry _slasherFactory) external onlyOwner {
        _setSlasherFactory(_slasherFactory);
    }

    /// @dev Sets the burner router factory, restricted to contract owner.
    function setBurnerRouterFactory(IRegistry _burnerRouterFactory) external onlyOwner {
        _setBurnerRouterFactory(_burnerRouterFactory);
    }

    /// @dev Sets the network address, restricted to contract owner.
    function setNetwork(address _network) external onlyOwner {
        _setNetwork(_network);
    }

    /// @dev Sets the slash period in seconds, restricted to contract owner.
    ///
    /// @dev WARNING: Before the owner calls this function, ALL registered vaults must be validated
    /// or deregistered, to conform to the new slashPeriodSeconds.
    ///
    /// @dev The owner will need to use VaultRegistered, VaultDeregistrationRequested, and VaultDeregistered events
    /// to determine the registered vault set, call wouldVaultBeValidWith for each vault, and deregister any
    /// vaults that are not valid. Only then can this function be called.
    function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) external onlyOwner {
        _setSlashPeriodSeconds(slashPeriodSeconds_);
    }

    /// @dev Sets the slash oracle, restricted to contract owner.
    function setSlashOracle(address slashOracle_) external onlyOwner {
        _setSlashOracle(slashOracle_);
    }

    /// @dev Sets the slash receiver, restricted to contract owner.
    /// @dev In practice, this function should never be called.
    /// @dev It exists for emergency scenarios only, with the assumption that changing the slash receiver
    /// would invalidate all existing vaults, and they would need to update their burnerRouters.
    function setSlashReceiver(address slashReceiver_) external onlyOwner {
        _setSlashReceiver(slashReceiver_);
    }

    /// @dev Sets the minimum burner router delay, restricted to contract owner.
    /// @dev Calling this function may invalidate existing registered vaults. This exists for emergency scenarios only.
    function setMinBurnerRouterDelay(uint256 minBurnerRouterDelay_) external onlyOwner {
        _setMinBurnerRouterDelay(minBurnerRouterDelay_);
    }

    /// @dev Checks if a vault would be valid with a given slashPeriodSeconds.
    /// @return True if the vault would be valid, reverts otherwise.
    function wouldVaultBeValidWith(address vault, uint256 potentialSLashPeriodSeconds) external view returns (bool) {
        _validateVaultParams(vault, potentialSLashPeriodSeconds);
        return true;
    }

    /// @notice Queries if a validator is opted-in to mev-commit through a vault.
    /// @dev The oracle must continuously call this function for upcoming proposers, in order to maintain 
    /// the most recent (finalized) block timestamp that a validator was queried as "opted-in", see `captureTimestamp` in README.md.
    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    /// @notice Queries if a validator is slashable.
    function isValidatorSlashable(bytes calldata blsPubkey) external view returns (bool) {
        ValidatorRecord storage record = validatorRecords[blsPubkey];
        if (!record.exists) {
            return false;
        }
        VaultRecord storage vaultRecord = vaultRecords[record.vault];
        if (!vaultRecord.exists) {
            return false;
        }
        OperatorRecord storage operatorRecord = operatorRecords[record.operator];
        if (!operatorRecord.exists) {
            return false;
        }
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

    function getLatestSlashAmount(address vault) external view returns (uint160) {
        return _getLatestSlashAmount(vault);
    }

    function getSlashAmountAt(address vault, uint256 blockTimestamp) external view returns (uint160) {
        return _getSlashAmountAt(vault, blockTimestamp);
    }

    function isVaultBurnerValid(address vault) external view returns (bool) {
        return _validateVaultBurner(vault);
    }

    function isVaultBurnerValidAgainstOperator(address vault, address operator) external view returns (bool) {
        return _validateVaultBurnerAgainstOperator(vault, operator);
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

    function _setVaultRecord(address vault, uint160 slashAmount) internal {
        vaultRecords[vault] = VaultRecord({
            exists: true,
            deregRequestOccurrence: TimestampOccurrence.Occurrence({
                exists: false,
                timestamp: 0
            }),
            slashAmountHistory: Checkpoints.Trace160({
                _checkpoints: new Checkpoints.Checkpoint160[](0)
            })
        });
        vaultRecords[vault].slashAmountHistory.push(SafeCast.toUint96(block.timestamp), slashAmount);
    }

    function _registerVault(address vault, uint160 slashAmount) internal {
        require(!vaultRecords[vault].exists, VaultAlreadyRegistered(vault));
        require(vaultFactory.isEntity(vault), VaultNotEntity(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));
        _validateVaultParams(vault, slashPeriodSeconds);
        require(_validateVaultBurner(vault), InvalidVaultBurner(vault));
        _setVaultRecord(vault, slashAmount);
        emit VaultRegistered(vault, slashAmount);
    }

    function _updateSlashAmount(address vault, uint160 slashAmount) internal {
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));
        record.slashAmountHistory.push(SafeCast.toUint96(block.timestamp), slashAmount);
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
        bool success = _vaultAndOperatorToValset[vault][operator].add(blsPubkey);
        require(success, FailedToAddValidatorToValset(blsPubkey, vault, operator)); // This error would indicate state corruption.
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
    /// @param captureTimestamp block.timestamp of the most recent finalized L1 block that the blsPubkey was queried as "opted-in" by the oracle.
    /// @dev This function is guaranteed to succeed if current tx executes within slashPeriodSeconds 
    /// of the provided captureTimestamp, AND the captureTimestamp was correctly computed as defined in README.md.
    /// @dev Operator and vault are not deregistered for the validator's infraction, so as to avoid opting-out large groups of validators at once.
    function _slashValidator(bytes calldata blsPubkey, uint256 captureTimestamp, ValidatorRecord storage valRecord) internal {
        VaultRecord storage vaultRecord = vaultRecords[valRecord.vault];
        require(vaultRecord.exists, MissingVaultRecord(valRecord.vault));
        OperatorRecord storage operatorRecord = operatorRecords[valRecord.operator];
        require(operatorRecord.exists, MissingOperatorRecord(valRecord.operator));

        require(captureTimestamp != 0, CaptureTimestampMustBeNonZero());

        // Slash amount is enforced as non-zero in _registerVault.
        uint160 amount = _getSlashAmountAt(valRecord.vault, captureTimestamp);

        address slasher = IVault(valRecord.vault).slasher();
        uint256 slasherType = IEntity(slasher).TYPE();
        uint256 slashedAmount;
        if (slasherType == _VETO_SLASHER_TYPE) {
            IVetoSlasher vetoSlasher = IVetoSlasher(slasher);
            uint256 slashIndex = vetoSlasher.requestSlash(
                _getSubnetwork(), valRecord.operator, amount, SafeCast.toUint48(captureTimestamp), new bytes(0));
            // Since resolver = address(0), slash can be executed immediately.
            slashedAmount = vetoSlasher.executeSlash(slashIndex, new bytes(0));
        } else if (slasherType == _INSTANT_SLASHER_TYPE) {
            slashedAmount = ISlasher(slasher).slash(
                _getSubnetwork(), valRecord.operator, amount, SafeCast.toUint48(captureTimestamp), new bytes(0));
        }
        emit ValidatorSlashed(blsPubkey, valRecord.operator, valRecord.vault, slashedAmount);

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

    /// @dev Internal function to set the delegator factory.
    function _setDelegatorFactory(IRegistry _delegatorFactory) internal {
        require(_delegatorFactory != IRegistry(address(0)), ZeroAddressNotAllowed());
        delegatorFactory = _delegatorFactory;
        emit DelegatorFactorySet(address(_delegatorFactory));
    }

    /// @dev Internal function to set the slasher factory.
    function _setSlasherFactory(IRegistry _slasherFactory) internal {
        require(_slasherFactory != IRegistry(address(0)), ZeroAddressNotAllowed());
        slasherFactory = _slasherFactory;
        emit SlasherFactorySet(address(_slasherFactory));
    }

    /// @dev Internal function to set the burner router factory.
    function _setBurnerRouterFactory(IRegistry _burnerRouterFactory) internal {
        require(_burnerRouterFactory != IRegistry(address(0)), ZeroAddressNotAllowed());
        burnerRouterFactory = _burnerRouterFactory;
        emit BurnerRouterFactorySet(address(_burnerRouterFactory));
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

    /// @dev Internal function to set the slash receiver.
    function _setSlashReceiver(address slashReceiver_) internal {
        require(slashReceiver_ != address(0), ZeroAddressNotAllowed());
        slashReceiver = slashReceiver_;
        emit SlashReceiverSet(slashReceiver_);
    }

    /// @dev Internal function to set the minimum burner router delay.
    function _setMinBurnerRouterDelay(uint256 minBurnerRouterDelay_) internal {
        require(minBurnerRouterDelay_ != 0, ZeroUintNotAllowed());
        minBurnerRouterDelay = minBurnerRouterDelay_;
        emit MinBurnerRouterDelaySet(minBurnerRouterDelay_);
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _validateVaultParams(address vault, uint256 slashPeriodSeconds) internal view {
        address delegator = IVault(vault).delegator();

        require(delegatorFactory.isEntity(delegator), DelegatorNotEntity(delegator, address(delegatorFactory)));

        uint256 delegatorType = IEntity(delegator).TYPE();
        if (delegatorType == _FULL_RESTAKE_DELEGATOR_TYPE) {
            revert FullRestakeDelegatorNotSupported(vault);
        // Only two delegator types are supported, network-restake and operator-specific.
        } else if (delegatorType != _NETWORK_RESTAKE_DELEGATOR_TYPE && delegatorType != _OPERATOR_SPECIFIC_DELEGATOR_TYPE) {
            revert UnknownDelegatorType(vault, delegatorType);
        }

        IVaultStorage vaultContract = IVaultStorage(vault);
        uint256 vaultEpochDurationSeconds = vaultContract.epochDuration();

        address slasher = IVault(vault).slasher();
        require(slasher != address(0), SlasherNotSetForVault(vault));

        require(slasherFactory.isEntity(slasher), SlasherNotEntity(slasher, address(slasherFactory)));

        uint256 slasherType = IEntity(slasher).TYPE();
        if (slasherType == _VETO_SLASHER_TYPE) {
            IVetoSlasher vetoSlasher = IVetoSlasher(slasher);
            uint256 vetoDuration = vetoSlasher.vetoDuration();
            // Incorporate that veto duration will eat into portion of the epoch that oracle can feasibly request slashes.
            vaultEpochDurationSeconds -= vetoDuration; /// @dev No underflow possible, vetoDuration must be less than epochDuration as enforced by VetoSlasher.sol.
            // Veto slasher must have a zero resolver s.t. slash can be executed immediately.
            require(vetoSlasher.resolver(_getSubnetwork(), new bytes(0)) == address(0), VetoSlasherMustHaveZeroResolver(vault));
        } else if (slasherType != _INSTANT_SLASHER_TYPE) {
            revert UnknownSlasherType(vault, slasherType);
        }

        require(IBaseSlasher(slasher).isBurnerHook(), BurnerHookNotSetForVault(vault));

        require(vaultEpochDurationSeconds > slashPeriodSeconds,
            InvalidVaultEpochDuration(vault, vaultEpochDurationSeconds, slashPeriodSeconds));
    }

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

    function _validateVaultBurner(address vault) internal view returns (bool) {
        address burner = IVault(vault).burner();
        if (!burnerRouterFactory.isEntity(burner)) {
            return false;
        }
        IBurnerRouter burnerRouter = IBurnerRouter(burner);
        if (burnerRouter.networkReceiver(network) != slashReceiver) {
            return false;
        }
        return burnerRouter.delay() >= minBurnerRouterDelay;
    }

    function _validateVaultBurnerAgainstOperator(address vault, address operator) internal view returns (bool) {
        bool isValidBurner = _validateVaultBurner(vault);
        if (!isValidBurner) {
            return false;
        }
        address burner = IVault(vault).burner();
        IBurnerRouter burnerRouter = IBurnerRouter(burner);
        address operatorNetworkReceiver = burnerRouter.operatorNetworkReceiver(network, operator);
        // Operator network receiver must be slashReceiver or null, but it cannot override a correct network receiver.
        return operatorNetworkReceiver == address(0) || operatorNetworkReceiver == slashReceiver;
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

    /// @return Number of validators that are slashable given stake in the vault at the current block.timestamp.
    function _getNumSlashableVals(address vault, address operator) internal view returns (uint256) {
        IBaseDelegator delegator = IBaseDelegator(IVault(vault).delegator());
        uint256 allocatedStake = delegator.stake(_getSubnetwork(), operator); // Uses current block.timestamp, contrary to stakeAt().
        uint160 slashAmount = vaultRecords[vault].slashAmountHistory.latest();
        return allocatedStake / slashAmount;
    }

    function _isValidatorSlashable(bytes calldata blsPubkey, address vault, address operator) internal view returns (bool) {
        uint256 slashableVals = _getNumSlashableVals(vault, operator);
        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        require(position != 0, ValidatorNotInValset(blsPubkey, vault, operator));
        return position <= slashableVals; // position is 1-indexed
    }

    /// @return Number of validators that could be slashable, given the current stake in the vault.
    function _potentialSlashableVals(address vault, address operator) internal view returns (uint256) {
        uint256 slashableVals = _getNumSlashableVals(vault, operator);
        uint256 numRegistered = _vaultAndOperatorToValset[vault][operator].length();
        if (slashableVals < numRegistered) {
            return 0;
        }
        return slashableVals - numRegistered;
    }

    function _getLatestSlashAmount(address vault) internal view returns (uint160 amount) {
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        amount = record.slashAmountHistory.latest();
        require(amount != 0, NoSlashAmountAtTimestamp(vault, block.timestamp));
        return amount;
    }

    function _getSlashAmountAt(address vault, uint256 timestamp) internal view returns (uint160 amount) {
        require(timestamp <= block.timestamp, FutureTimestampDisallowed(vault, timestamp));
        VaultRecord storage record = vaultRecords[vault];
        require(record.exists, VaultNotRegistered(vault));
        amount = record.slashAmountHistory.upperLookup(SafeCast.toUint96(timestamp));
        require(amount != 0, NoSlashAmountAtTimestamp(vault, timestamp));
        return amount;
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
        if (!_validateVaultBurnerAgainstOperator(valRecord.vault, valRecord.operator)) {
            return false;
        }
        return true;
    }
}
