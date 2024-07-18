// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {IValidatorRegistryV1} from "./IValidatorRegistryV1.sol";
import {IMevCommitAVS} from "./IMevCommitAVS.sol";

interface IValidatorOptInRouter {
    /// @notice Initializes the contract with the validator registry and mev-commit AVS contracts.
    function initialize(
        address _validatorRegistry,
        address _mevCommitAVS,
        address _owner
    ) external;

    /// @notice Allows the owner to set the validator registry V1 contract.
    function setValidatorRegistryV1(IValidatorRegistryV1 _validatorRegistry) external;

    /// @notice Allows the owner to set the mev-commit AVS contract.
    function setMevCommitAVS(IMevCommitAVS _mevCommitAVS) external;

    /// @notice Returns an array of bools indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory);

    /// @notice Emitted when the validator registry V1 contract is set.
    event ValidatorRegistryV1Set(address oldContract, address newContract);

    /// @notice Emitted when the mev-commit AVS contract is set.
    event MevCommitAVSSet(address oldContract, address newContract);
}
