// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IValidatorOptInHub {

    /// @notice Emitted when a registry is added.
    event RegistryAdded(uint256 indexed index, address indexed registry);

    /// @notice Emitted when a registry is replaced.
    event RegistryReplaced(uint256 indexed index, address indexed oldRegistry, address indexed newRegistry);

    /// @notice Emitted when a registry is removed.
    event RegistryRemoved(uint256 indexed index, address indexed registry);
    

    error InvalidRegistry();
    error InvalidIndex();
    error ZeroAddress();
    error IndexRegistryMismatch();

    /// @notice Initializes the contract with the validator registry and mev-commit AVS contracts.
    function initialize(
        address[] calldata _registries,
        address _owner
    ) external;


    /// @notice Adds a registry to the contract.
    function addRegistry(address registry) external;

    /// @notice Replaces a registry with a new one.
    function updateRegistry(uint256 index, address oldRegistry, address newRegistry) external;

    /// @notice Removes a registry from the contract.
    function removeRegistry(uint256 index, address registry) external;

    /// @notice Returns an array of bool lists indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedInList(bytes[] calldata valBLSPubKeys) external view returns (bool[][] memory optInStatuses);

    /// @notice Returns a bool list indicating whether a validator pubkey is opted in to mev-commit with any of the registries.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory optInStatuses);

    /// @notice Returns a bool list indicating whether a validator pubkey is opted in to mev-commit.
    function isValidatorOptedInList(bytes calldata valPubKey) external view returns (bool[] memory optInStatus);

    /// @notice Returns a bool indicating whether a validator pubkey is opted in to mev-commit with any of the registries.
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);
}