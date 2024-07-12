// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

interface IValidatorOptInRouter {
    /// @notice Initializes the contract with the validator registry and mev-commit AVS contracts.
    function initialize(
        address _validatorRegistry,
        address _mevCommitAVS,
        address _owner
    ) external;

    /// @notice Allows the owner to set the validator registry V1 contract.
    function setValidatorRegistryV1(address _validatorRegistry) external;

    /// @notice Allows the owner to set the mev-commit AVS contract.
    function setMevCommitAVS(address _mevCommitAVS) external;

    /// @notice Returns an array of bools indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory);
}
