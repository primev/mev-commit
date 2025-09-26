// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IRocketMinipoolRegistry {

    struct ValidatorRegistration {
        bool exists;
        uint64 deregTimestamp;
        uint64 freezeTimestamp;
    }

    // ----------------- EVENTS -----------------

    /// @notice Emitted when a validator is registered.
    event ValidatorRegistered(bytes indexed validatorPubKey, address indexed nodeAddress);

    /// @notice Emitted when a validator deregistration request is made.
    event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed nodeAddress);

    /// @notice Emitted when a validator is deregistered.
    event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed nodeAddress);

    /// @notice Emitted when a validator is frozen.
    event ValidatorFrozen(bytes indexed validatorPubKey);

    /// @notice Emitted when a validator is unfrozen.
    event ValidatorUnfrozen(bytes indexed validatorPubKey);

    // ----------------- ERRORS -----------------

    error ValidatorAlreadyRegistered(bytes validatorPubkey);

    error ValidatorAlreadyFrozen(bytes validatorPubkey);

    error ValidatorNotFrozen(bytes validatorPubkey);

    error ValidatorNotRegistered(bytes validatorPubkey);

    error DeregRequestAlreadyExists(bytes validatorPubkey);

    error FrozenValidatorCannotDeregister(bytes validatorPubkey);

    error DeregistrationTooSoon(bytes validatorPubkey);

    error ValidatorDeregistrationNotExpired(bytes validatorPubkey);

    error NotMinipoolOperator(bytes validatorPubkey);

    error MinipoolNotActive(bytes validatorPubkey);

    error NoMinipoolForKey(bytes validatorPubkey);

    error DeregRequestDoesNotExist(bytes validatorPubkey);

    error OnlyFreezeOracle();

    error InvalidBLSPubKeyLength(uint256 expectedLength, uint256 actualLength);

    error UnfreezeFeeRequired(uint256 requiredFee);

    error UnfreezeTransferFailed();

    error RefundFailed();

    error ZeroParam();

    // ----------------- FUNCTIONS -----------------

    /// @notice Registers validators with a minipool.
    function registerValidators(bytes[] calldata validatorPubkeys) external;

    /// @notice Requests deregistration for validators.
    function requestValidatorDeregistration(bytes[] calldata validatorPubkeys) external;

    /// @notice Deregisters validators.
    function deregisterValidators(bytes[] calldata validatorPubkeys) external;

    /// @notice Freezes validators.
    function freeze(bytes[] calldata validatorPubkeys) external;

    /// @notice Unfreezes validators.
    function unfreeze(bytes[] calldata validatorPubkeys) external payable;

    /// @notice Returns the node address for a validator.
    function getNodeAddressFromPubkey(bytes calldata validatorPubkey) external view returns (address);

    /// @notice Returns the minipool for a validator.
    function getMinipoolFromPubkey(bytes calldata validatorPubkey) external view returns (address);

    /// @notice Returns the node address for a validator's minipool.
    function getNodeAddressFromMinipool(address minipool) external view returns (address);

    /// @notice Returns both the node address and the withdrawal address of the key's minipool, as these addresses both have minipool permissions.
    function getValidOperatorsForKey(bytes calldata validatorPubkey) external view returns (address, address);

    /// @notice Returns the time at which a validator can be deregistered.
    function getEligibleTimeForDeregistration(bytes calldata validatorPubkey) external view returns (uint64);

    /// @notice Returns the validator registration info.
    function getValidatorRegInfo(bytes calldata validatorPubkey) external view returns (ValidatorRegistration memory);

    /// @notice Checks if a validator is opted-in.  
    function isValidatorOptedIn(bytes calldata validatorPubkey) external view returns (bool);

    /// @notice Checks if an operator is valid to interact on behalf of a validator.
    function isOperatorValidForKey(address operator, bytes calldata validatorPubkey) external view returns (bool);
}