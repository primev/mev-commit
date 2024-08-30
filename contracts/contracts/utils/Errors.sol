// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

library Errors {
    /// @dev Custom error for invalid fallback calls.
    error InvalidFallback();

    /// @dev Custom error for invalid receive calls.
    error InvalidReceive();
}
