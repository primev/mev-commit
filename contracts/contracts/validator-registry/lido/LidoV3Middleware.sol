// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {LidoV3MiddlewareStorage} from "./LidoV3MiddlewareStorage.sol";
import {TimestampOccurrence} from "../../utils/Occurrence.sol";

interface IVaultHub {
    function isVaultConnected(address vault) external view returns (bool);
    function totalValue(address vault) external view returns (uint256);
}

interface IStakingVault {
    function nodeOperator() external view returns (address);
}

contract LidoV3Middleware is
    Initializable,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable,
    LidoV3MiddlewareStorage
{
    // ---- Constants ----
    uint256 private constant BLS_PUBKEY_LEN = 48;

    // ---- Events ----
    event VaultHubUpdated(address indexed newHub);
    event WhitelistUpdated(address indexed operator, bool allowed);
    event SlashAmountUpdated(uint256 newSlashAmount);
    event UnfreezeFeeUpdated(uint256 newUnfreezeFee);
    event UnfreezeReceiverUpdated(address indexed newReceiver);
    event UnfreezePeriodUpdated(uint256 newPeriod);
    event DeregistrationPeriodUpdated(uint256 newPeriod);

    event ValidatorRegistered(bytes indexed pubkey, address indexed registrar, address indexed vault);
    event ValidatorDeregistrationRequested(bytes indexed pubkey, address indexed registrar);
    event ValidatorDeregistered(bytes indexed pubkey, address indexed deregistrar);
    event ValidatorFrozen(bytes indexed pubkey, address indexed freezer);
    event ValidatorUnfrozen(bytes indexed pubkey, address indexed unfreezer);

    // ---- Errors ----
    error ZeroAddress();
    error NotWhitelisted();
    error InvalidVault();
    error HubDoesNotOwnVault();
    error NodeOperatorMismatch(address expected, address actual);
    error CapacityExceeded(uint256 want, uint256 have);
    error InvalidBLSPubKeyLength(uint256 expected, uint256 actual);
    error ValidatorAlreadyRegistered(bytes pubkey);
    error ValidatorNotRegistered(bytes pubkey);
    error ValidatorAlreadyRequestedDeregistration(bytes pubkey);
    error ValidatorNotFrozen(bytes pubkey);
    error UnfreezeTooSoon();
    error UnfreezeFeeRequired(uint256 requiredFee);
    error TransferFailed();

    // ---- Initializer ----
    /**
     * @param _owner            initial owner (Ownable2StepUpgradeable)
     * @param _vaultHub         canonical VaultHub (Hoodi: 0x26b92f0fdfeBAf43E5Ea5b5974EeBee95F17Fe08)
     * @param _unfreezeReceiver receiver of unfreeze fees (can be owner same as _owner)
     * @param _deregPeriod      seconds from request to allow deregistration
     * @param _unfreezePeriod   seconds after freeze until unfreeze is permitted
     */
    function initialize(
        address _owner,
        address _vaultHub,
        address _unfreezeReceiver,
        uint256 _deregPeriod,
        uint256 _unfreezePeriod,
        uint256 _slashAmount,
        uint256 _unfreezeFee
    ) external initializer {
        if (_owner == address(0) || _vaultHub == address(0) || _unfreezeReceiver == address(0)) revert ZeroAddress();

        __Ownable_init(_owner);
        __Ownable2Step_init();
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();

        vaultHub = _vaultHub;
        slashAmount = _slashAmount;
        unfreezeFee = _unfreezeFee;
        unfreezeReceiver = _unfreezeReceiver;
        deregistrationPeriod = _deregPeriod;
        unfreezePeriod = _unfreezePeriod;

        emit VaultHubUpdated(_vaultHub);
        emit SlashAmountUpdated(slashAmount);
        emit UnfreezeFeeUpdated(unfreezeFee);
        emit UnfreezeReceiverUpdated(_unfreezeReceiver);
        emit DeregistrationPeriodUpdated(_deregPeriod);
        emit UnfreezePeriodUpdated(_unfreezePeriod);
    }

    // ---- Setters ----

    function setVaultHub(address newHub) external onlyOwner {
        if (newHub == address(0)) revert ZeroAddress();
        vaultHub = newHub;
        emit VaultHubUpdated(newHub);
    }

    function setWhitelist(address operator, bool allowed) external onlyOwner {
        if (operator == address(0)) revert ZeroAddress();
        isWhitelisted[operator] = allowed;
        emit WhitelistUpdated(operator, allowed);
    }

    function setSlashAmount(uint256 newSlashAmount) external onlyOwner {
        // allow zero? probably not, but testnet—you decide.
        slashAmount = newSlashAmount;
        emit SlashAmountUpdated(newSlashAmount);
    }

    function setUnfreezeFee(uint256 newFee) external onlyOwner {
        unfreezeFee = newFee;
        emit UnfreezeFeeUpdated(newFee);
    }

    function setUnfreezeReceiver(address newReceiver) external onlyOwner {
        if (newReceiver == address(0)) revert ZeroAddress();
        unfreezeReceiver = newReceiver;
        emit UnfreezeReceiverUpdated(newReceiver);
    }

    function setUnfreezePeriod(uint256 newPeriod) external onlyOwner {
        unfreezePeriod = newPeriod;
        emit UnfreezePeriodUpdated(newPeriod);
    }

    function setDeregistrationPeriod(uint256 newPeriod) external onlyOwner {
        deregistrationPeriod = newPeriod;
        emit DeregistrationPeriodUpdated(newPeriod);
    }

    // ---- View ----

    function isValidatorOptedIn(bytes calldata valBLSPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valBLSPubKey);
    }

    function maxRegistrableByVault(address vault) public view returns (uint256) {
        // Capacity = floor( totalValue / slashAmount )
        uint256 tv = IVaultHub(vaultHub).totalValue(vault);
        if (slashAmount == 0) return 0; // avoid div by zero; owner should set non-zero
        return tv / slashAmount;
    }

    function remainingRegistrable(address vault) public view returns (uint256) {
        uint256 cap = maxRegistrableByVault(vault);
        uint256 used = vaultRegisteredCount[vault];
        return cap > used ? (cap - used) : 0;
    }

    // ---- Core logic ----

    /**
     * @notice Register validator BLS pubkeys for a specific Lido V3 vault.
     *         - caller must be whitelisted
     *         - vault must be connected to our configured VaultHub
     *         - vault's nodeOperator() must equal msg.sender
     *         - capacity check: (registered + new) <= floor(totalValue(vault) / slashAmount)
     */
    function registerValidators(address vault, bytes[] calldata blsPubKeys) external {
        if (!isWhitelisted[msg.sender]) revert NotWhitelisted();
        if (!_isValidVault(vault)) revert InvalidVault();

        // Node operator must be the caller
        address operator = IStakingVault(vault).nodeOperator();
        if (operator != msg.sender) revert NodeOperatorMismatch(operator, msg.sender);

        uint256 n = blsPubKeys.length;
        if (n == 0) revert CapacityExceeded(0, remainingRegistrable(vault)); // keep a consistent revert

        uint256 remaining = remainingRegistrable(vault);
        if (n > remaining) revert CapacityExceeded(n, remaining);

        for (uint256 i = 0; i < n; ++i) {
            bytes calldata pk = blsPubKeys[i];
            if (pk.length != BLS_PUBKEY_LEN) revert InvalidBLSPubKeyLength(BLS_PUBKEY_LEN, pk.length);
            LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pk];
            if (r.exists) revert ValidatorAlreadyRegistered(pk);

            validatorRecords[pk] = LidoV3MiddlewareStorage.ValidatorRecord({
                exists: true,
                registrar: msg.sender,
                freezeOccurrence: TimestampOccurrence.Occurrence({exists: false, timestamp: 0}),
                deregRequestOccurrence: TimestampOccurrence.Occurrence({exists: false, timestamp: 0})
            });

            emit ValidatorRegistered(pk, msg.sender, vault);
        }

        // bump per-vault count
        vaultRegisteredCount[vault] += n;
    }

    /**
     * @notice Operators request deregistration. Completion can happen after `deregistrationPeriod`.
     */
    function requestDeregistrations(bytes[] calldata blsPubKeys) external {
        if (!isWhitelisted[msg.sender]) revert NotWhitelisted();

        uint256 n = blsPubKeys.length;
        for (uint256 i = 0; i < n; ++i) {
            bytes calldata pk = blsPubKeys[i];
            if (pk.length != BLS_PUBKEY_LEN) revert InvalidBLSPubKeyLength(BLS_PUBKEY_LEN, pk.length);

            LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pk];
            if (!r.exists) revert ValidatorNotRegistered(pk);
            // Optional: only the original registrar can request dereg
            if (r.registrar != msg.sender) revert NotWhitelisted(); // reuse error to avoid extra bytecode

            if (r.deregRequestOccurrence.exists) revert ValidatorAlreadyRequestedDeregistration(pk);
            TimestampOccurrence.captureOccurrence(r.deregRequestOccurrence);

            emit ValidatorDeregistrationRequested(pk, msg.sender);
        }
    }

    /**
     * @notice Complete deregistration after the delay. Only the original registrar can complete.
     * @param vault The vault this batch was originally accounted toward (to reduce its count).
     */
    function deregisterValidators(address vault, bytes[] calldata blsPubKeys) external {
        if (!isWhitelisted[msg.sender]) revert NotWhitelisted();

        uint256 removed = 0;
        uint256 n = blsPubKeys.length;

        for (uint256 i = 0; i < n; ++i) {
            bytes calldata pk = blsPubKeys[i];
            if (pk.length != BLS_PUBKEY_LEN) revert InvalidBLSPubKeyLength(BLS_PUBKEY_LEN, pk.length);

            LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pk];
            if (!r.exists) revert ValidatorNotRegistered(pk);
            if (r.registrar != msg.sender) revert NotWhitelisted();
            if (!r.deregRequestOccurrence.exists) revert ValidatorAlreadyRequestedDeregistration(pk);
            // Wait out the deregistration window
            if (block.timestamp < r.deregRequestOccurrence.timestamp + deregistrationPeriod) revert UnfreezeTooSoon();

            delete validatorRecords[pk];
            unchecked { ++removed; }

            emit ValidatorDeregistered(pk, msg.sender);
        }

        if (removed > 0) {
            // reduce per-vault count (saturating)
            uint256 prev = vaultRegisteredCount[vault];
            vaultRegisteredCount[vault] = removed > prev ? 0 : (prev - removed);
        }
    }

    // ---- Freeze / Unfreeze ----

    function freezeValidators(bytes[] calldata blsPubKeys) external onlyOwner {
        uint256 n = blsPubKeys.length;
        for (uint256 i = 0; i < n; ++i) {
            bytes calldata pk = blsPubKeys[i];
            if (pk.length != BLS_PUBKEY_LEN) revert InvalidBLSPubKeyLength(BLS_PUBKEY_LEN, pk.length);

            LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pk];
            if (!r.exists) revert ValidatorNotRegistered(pk);

            TimestampOccurrence.captureOccurrence(r.freezeOccurrence);
            emit ValidatorFrozen(pk, msg.sender);
        }
    }

    /**
     * @notice Anyone can unfreeze for a fee, after `unfreezePeriod` from freeze time.
     */
    function unfreeze(bytes[] calldata blsPubKeys) external payable nonReentrant {
        uint256 n = blsPubKeys.length;
        uint256 required = unfreezeFee * n;
        if (msg.value < required) revert UnfreezeFeeRequired(required);

        for (uint256 i = 0; i < n; ++i) {
            bytes calldata pk = blsPubKeys[i];
            if (pk.length != BLS_PUBKEY_LEN) revert InvalidBLSPubKeyLength(BLS_PUBKEY_LEN, pk.length);

            LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pk];
            if (!r.exists) revert ValidatorNotRegistered(pk);
            if (!r.freezeOccurrence.exists) revert ValidatorNotFrozen(pk);
            if (block.timestamp < r.freezeOccurrence.timestamp + unfreezePeriod) revert UnfreezeTooSoon();

            TimestampOccurrence.del(r.freezeOccurrence);
            emit ValidatorUnfrozen(pk, msg.sender);
        }

        if (required > 0) {
            (bool ok1, ) = payable(unfreezeReceiver).call{value: required}("");
            if (!ok1) revert TransferFailed();
        }
        // refund dust
        uint256 excess = msg.value - required;
        if (excess != 0) {
            (bool ok2, ) = payable(msg.sender).call{value: excess}("");
            if (!ok2) revert TransferFailed();
        }
    }

    // ---- Internals ----

    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _isValidVault(address vault) internal view returns (bool) {
        // must be connected to our configured VaultHub
        if (!IVaultHub(vaultHub).isVaultConnected(vault)) return false;
        return true;
    }

    function _isValidatorOptedIn(bytes calldata pubkey) internal view returns (bool) {
        LidoV3MiddlewareStorage.ValidatorRecord storage r = validatorRecords[pubkey];
        return 
            r.exists && 
            isWhitelisted[r.registrar] &&
            !r.freezeOccurrence.exists && 
            !r.deregRequestOccurrence.exists;
    }

    // storage gap reserved in *storage* contract
}
