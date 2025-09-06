// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.26;

abstract contract TransientReentrancyGuard {
    bytes32 private constant _SLOT = keccak256("primev.reentrancy.guard.transient");

    modifier nonReentrant() {
        bytes32 slot = _SLOT;
        assembly ("memory-safe") {
            if tload(slot) { revert(0, 0) }
            tstore(slot, 1)
        }
        _;
        assembly ("memory-safe") {
            tstore(slot, 0)
        }
    }
}