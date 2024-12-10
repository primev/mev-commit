// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// Temporary stand-in until the mev-commit middleware and symbiotic core, are ready for mainnet.
contract AlwaysFalseMevCommitMiddleware {
    function isValidatorOptedIn(bytes calldata) external pure returns (bool) {
        return false;
    }
}
