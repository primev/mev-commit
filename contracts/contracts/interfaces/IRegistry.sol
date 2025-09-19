// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IRegistry {
    /// @notice Returns an array of OptInStatus structs indicating whether each validator pubkey is opted in to mev-commit.
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool);
}