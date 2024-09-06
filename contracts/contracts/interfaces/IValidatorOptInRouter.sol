// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IVanillaRegistry} from "./IVanillaRegistry.sol";
import {IMevCommitAVS} from "./IMevCommitAVS.sol";
import {IMevCommitMiddleware} from "./IMevCommitMiddleware.sol";

interface IValidatorOptInRouter {

    /// @notice Emitted when the vanilla registry contract is set.
    event VanillaRegistrySet(address oldContract, address newContract);

    /// @notice Emitted when the mev-commit AVS contract is set.
    event MevCommitAVSSet(address oldContract, address newContract);

    /// @notice Emitted when the mev-commit middleware contract is set.
    event MevCommitMiddlewareSet(address oldContract, address newContract);

    /// @notice Initializes the contract with the vanilla registry and mev-commit AVS contracts.
    function initialize(
        address _vanillaRegistry,
        address _mevCommitAVS,
        address _mevCommitMiddleware,
        address _owner
    ) external;

    /// @notice Allows the owner to set the vanilla registry contract.
    function setVanillaRegistry(IVanillaRegistry _vanillaRegistry) external;

    /// @notice Allows the owner to set the mev-commit AVS contract.
    function setMevCommitAVS(IMevCommitAVS _mevCommitAVS) external;

    /// @notice Allows the owner to set the mev-commit middleware contract.
    function setMevCommitMiddleware(IMevCommitMiddleware _mevCommitMiddleware) external;

    /// @notice Returns an array of bools indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory);
}
