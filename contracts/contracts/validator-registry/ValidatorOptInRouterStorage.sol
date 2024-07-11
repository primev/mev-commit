// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {IValidatorRegistryV1} from "../interfaces/IValidatorRegistryV1.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";

/// @title ValidatorOptInRouterStorage
/// @notice Storage components of the ValidatorOptInRouter contract.
contract ValidatorOptInRouterStorage {
    /// @notice The address of the V1 validator registry contract.
    IValidatorRegistryV1 public validatorRegistryV1;
    /// @notice The address of the mev-commit AVS contract.
    IMevCommitAVS public mevCommitAVS;
}
