// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

abstract contract MevCommitAVSStorage {
    // Mapping from address who opted in a set of validators to the Keccak-256 hash over that set.
    mapping(address => bytes32) public addressToValSetHash;

    mapping(address => bool) public whitelist;

    // TODO: Determine if stored on L1 or our chain. 
    mapping(bytes => bool) public frozenSet;
}
