// SPDX-License-Identifier: BSL 1.1

pragma solidity 0.8.26;

contract AlwaysFalseRegistry {
    function isValidatorOptedIn(bytes calldata) external pure returns (bool) {
        return false;
    }
}