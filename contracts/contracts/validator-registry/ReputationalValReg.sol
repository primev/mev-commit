// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

// To be deployed on L1 implementing reputational opt-in. 

// Optimize this contract around reputational opt-in, then future contract aggregates this with other types 

// TODO: separate out contract owner, and account that manages the whitelist
contract ReputationalValReg is OwnableUpgradeable, ReentrancyGuardUpgradeable {

    // TODO: add initializer and other init logic, also config params

    // TODO: Confirm this FSM makes point calcuations easy enough
    enum State { NotWhitelisted, Active, Frozen }
    struct WhitelistedEOAInfo {
        State state;
        uint numConsAddrsStored;
        uint256 freezeHeight;
        string organizationName; // TODO: Consider removing for gas purposes. This info could be stored offchain.
    }
    mapping(address => WhitelistedEOAInfo) private whitelistedEOAs;

    // List of stored validator consensus addresses with O(1) lookup indexed by consensus address. 
    // These addresses were at some point stored by a whitelisted EOA.
    // 
    // This mapping is intentionally not enumerable,
    // since actors should only need to query the 32 relevant proposers for an epoch at a time.
    // If for some reason an actor desires the full list of store validator cons addrs,
    // they could construct the set offchain via events.
    mapping(bytes => address) public storedConsAddrs;

    function isWhitelistedEOA(address eoa) external view returns (bool) {
        return whitelistedEOAs[eoa].state != State.NotWhitelisted;
    }

    function isValidatorOptedIn(bytes memory consAddr) public view returns (bool) {
        address eoa = storedConsAddrs[consAddr];
        bool isConsAddrStored = eoa != address(0);
        bool isEoaActive = whitelistedEOAs[eoa].state == State.Active;
        return isConsAddrStored && isEoaActive;
    }

    function areValidatorsOptedIn(bytes[] memory consAddrs) external view returns (bool[] memory) {
        bool[] memory results = new bool[](consAddrs.length);
        for (uint i = 0; i < consAddrs.length; i++) {
            results[i] = isValidatorOptedIn(consAddrs[i]);
        }
        return results;
    }

    event ConsAddrStored(bytes consAddr, address indexed eoa); // TODO: Index consAddr too?
    function storeConsAddrs(bytes[] memory consAddrs) external {
        require(whitelistedEOAs[msg.sender].state != State.NotWhitelisted, "sender must be whitelisted");
        for (uint i = 0; i < consAddrs.length; i++) {
            require(storedConsAddrs[consAddrs[i]] == address(0), "Consensus address is already stored");
            // TODO: Make configurable
            require(whitelistedEOAs[msg.sender].numConsAddrsStored < 10000, "EOA must not store more than 10k cons addrs");
            storedConsAddrs[consAddrs[i]] = msg.sender;
            whitelistedEOAs[msg.sender].numConsAddrsStored++;
            emit ConsAddrStored(consAddrs[i], msg.sender);
        }
    }

    function deleteConsAddrs(bytes[] memory consAddrs) external {
        for (uint i = 0; i < consAddrs.length; i++) {
            require(whitelistedEOAs[msg.sender].state != State.NotWhitelisted, "sender must be whitelisted");
            require(storedConsAddrs[consAddrs[i]] == msg.sender, "Consensus address must be stored by sender");
            _deleteConsAddr(msg.sender, consAddrs[i]);
            whitelistedEOAs[msg.sender].numConsAddrsStored--;
        }
    }

    function deleteConsAddrsFromInactiveEOAs(bytes[] memory consAddrs) external onlyOwner {
        for (uint i = 0; i < consAddrs.length; i++) {
            address eoa = storedConsAddrs[consAddrs[i]];
            require(eoa != address(0), "Consensus address must be stored");
            require(whitelistedEOAs[eoa].state == State.NotWhitelisted, "EOA who originally stored cons addr must be inactive");
            _deleteConsAddr(eoa, consAddrs[i]);
        }
    }

    event ConsAddrDeleted(bytes consAddr, address indexed eoa); // TODO: Index consAddr too?
    function _deleteConsAddr(address eoa, bytes memory consAddr) internal {
        delete storedConsAddrs[consAddr];
        emit ConsAddrDeleted(consAddr, eoa);
    }

    function addWhitelistedEOA(address eoa, string memory organizationName) external onlyOwner {
        require(eoa != address(0), "Invalid address");
        require(bytes(organizationName).length > 0, "Organization name cannot be empty");
        require(whitelistedEOAs[eoa].state == State.NotWhitelisted, "EOA must not already be whitelisted");
        whitelistedEOAs[eoa] = WhitelistedEOAInfo({
            state: State.Active,
            numConsAddrsStored: 0,
            freezeHeight: 0,
            organizationName: organizationName
        });
    }

    function deleteWhitelistedEOA(address eoa) external {
        require(msg.sender == owner() || msg.sender == eoa, "Only owner or EOA itself can delete whitelisted EOA");
        require(whitelistedEOAs[eoa].state != State.NotWhitelisted, "EOA must be whitelisted");
        delete whitelistedEOAs[eoa];
    }

    function freeze(address eoa) onlyOwner external {
        require(whitelistedEOAs[eoa].state == State.Active, "EOA must be active");
        whitelistedEOAs[eoa].state = State.Frozen;
        whitelistedEOAs[eoa].freezeHeight = block.number;
    }

    function unfreeze() external payable {
        require(whitelistedEOAs[msg.sender].state == State.Frozen, "sender must be frozen");
        // TODO: make configurable
        require(block.number > whitelistedEOAs[msg.sender].freezeHeight + 1000, "EOA must have been frozen for at least 1000 blocks");
        // TODO: make configurable
        require(msg.value >= 10 ether, "10 ether must be sent with txn");
        // TODO: confirm eth received by contract
        whitelistedEOAs[msg.sender].state = State.Active;
        whitelistedEOAs[msg.sender].freezeHeight = 0;
    }
}
