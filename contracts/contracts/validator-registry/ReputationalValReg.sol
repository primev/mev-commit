// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";

// ReputationalValReg manages the reputational opt-in for mev-commit validators. 
// This contract is meant to be deployed on L1. Future contracts will implement 
// other types of opt-in including restaked opt-in and simple stake opt-in. 
//
// TODO: Consider separating out contract owner, and account that manages the whitelist. This depends how exactly upgrades will work.
// TODO: Determine need for reentrancy guard. Also determine if certain functions need to be external vs public for future integration.
// TODO: Hash out and test upgrade process before deployment.
contract ReputationalValReg is OwnableUpgradeable, UUPSUpgradeable {

    uint256 constant FUNC_ARG_ARRAY_LIMIT = 100;

    uint256 public maxConsAddrsPerEOA;
    uint256 public minFreezeBlocks;
    uint256 public unfreezeFee;

    enum State { NotWhitelisted, Active, Frozen }

    struct WhitelistedEOAInfo {
        State state;
        uint numConsAddrsStored;
        uint256 freezeHeight;
    }

    // Mapping of whitelisted EOAs to their info struct
    mapping(address => WhitelistedEOAInfo) public whitelistedEOAs;

    // List of stored validator consensus addresses with O(1) lookup indexed by consensus address. 
    // These addresses were at some point stored by a whitelisted EOA.
    // 
    // This mapping is intentionally not enumerable,
    // since actors should only need to query the 32 relevant proposers for an epoch at a time.
    // If for some reason an actor desires the full set of stored validator cons addrs,
    // they could construct the set offchain via events.
    mapping(bytes => address) public storedConsAddrs;

    event WhitelistedEOAAdded(address indexed eoa);
    event WhitelistedEOADeleted(address indexed eoa);
    event EOAFrozen(address indexed eoa);
    event EOAUnfrozen(address indexed eoa);
    event ConsAddrStored(bytes consAddr, address indexed eoa);
    event ConsAddrDeleted(bytes consAddr, address indexed eoa);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // TODO: Determine upgrade logic and test process
    }

    function initialize(
        address _owner,
        uint256 _maxConsAddrsPerEOA,
        uint256 _minFreezeBlocks,
        uint256 _unfreezeFee
    ) external initializer {
        __Ownable_init(_owner);
        maxConsAddrsPerEOA = _maxConsAddrsPerEOA;
        minFreezeBlocks = _minFreezeBlocks;
        unfreezeFee = _unfreezeFee;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function addWhitelistedEOA(address eoa) external onlyOwner {
        require(eoa != address(0), "Invalid address");
        require(whitelistedEOAs[eoa].state == State.NotWhitelisted, "EOA must not already be whitelisted");
        whitelistedEOAs[eoa] = WhitelistedEOAInfo({
            state: State.Active,
            numConsAddrsStored: 0,
            freezeHeight: 0
        });
        emit WhitelistedEOAAdded(eoa);
    }

    function deleteWhitelistedEOA(address eoa) external {
        require(msg.sender == owner() || msg.sender == eoa, "Only owner or EOA itself can delete whitelisted EOA");
        require(whitelistedEOAs[eoa].state != State.NotWhitelisted, "EOA must be whitelisted");
        delete whitelistedEOAs[eoa];
        emit WhitelistedEOADeleted(eoa);
    }

    function freeze(address eoa) onlyOwner external {
        require(whitelistedEOAs[eoa].state == State.Active, "EOA must be active");
        whitelistedEOAs[eoa].state = State.Frozen;
        whitelistedEOAs[eoa].freezeHeight = block.number;
        emit EOAFrozen(eoa);
    }

    function unfreeze() external payable {
        require(whitelistedEOAs[msg.sender].state == State.Frozen, "sender must be frozen");
        require(block.number >= whitelistedEOAs[msg.sender].freezeHeight + minFreezeBlocks, "Freeze period has not elapsed");
        require(msg.value >= unfreezeFee, "Insufficient unfreeze fee");
        whitelistedEOAs[msg.sender].state = State.Active;
        whitelistedEOAs[msg.sender].freezeHeight = 0;
        emit EOAUnfrozen(msg.sender);
    }

    function storeConsAddrs(bytes[] memory consAddrs) external {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        require(whitelistedEOAs[msg.sender].state != State.NotWhitelisted, "sender must be whitelisted");
        for (uint i = 0; i < consAddrs.length; i++) {
            require(storedConsAddrs[consAddrs[i]] == address(0), "Consensus address is already stored");
            require(whitelistedEOAs[msg.sender].numConsAddrsStored < maxConsAddrsPerEOA, "EOA must not store more than max allowed cons addrs");
            storedConsAddrs[consAddrs[i]] = msg.sender;
            whitelistedEOAs[msg.sender].numConsAddrsStored++;
            emit ConsAddrStored(consAddrs[i], msg.sender);
        }
    }

    function deleteConsAddrs(bytes[] memory consAddrs) external {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        for (uint i = 0; i < consAddrs.length; i++) {
            require(whitelistedEOAs[msg.sender].state != State.NotWhitelisted, "sender must be whitelisted");
            require(storedConsAddrs[consAddrs[i]] == msg.sender, "Consensus address must be stored by sender");
            _deleteConsAddr(msg.sender, consAddrs[i]);
            whitelistedEOAs[msg.sender].numConsAddrsStored--;
        }
    }

    function deleteConsAddrsFromNonWhitelistedEOAs(bytes[] memory consAddrs) external {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        for (uint i = 0; i < consAddrs.length; i++) {
            address eoa = storedConsAddrs[consAddrs[i]];
            require(eoa != address(0), "Consensus address must be stored");
            require(whitelistedEOAs[eoa].state == State.NotWhitelisted, "EOA who originally stored cons addr must not be whitelisted");
            _deleteConsAddr(eoa, consAddrs[i]);
        }
    }

    function isEOAWhitelisted(address eoa) external view returns (bool) {
        return whitelistedEOAs[eoa].state != State.NotWhitelisted;
    }

    function areValidatorsOptedIn(bytes[] memory consAddrs) external view returns (bool[] memory) {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        bool[] memory results = new bool[](consAddrs.length);
        for (uint i = 0; i < consAddrs.length; i++) {
            results[i] = _isValidatorOptedIn(consAddrs[i]);
        }
        return results;
    }

    function _deleteConsAddr(address eoa, bytes memory consAddr) internal {
        delete storedConsAddrs[consAddr];
        emit ConsAddrDeleted(consAddr, eoa);
    }

    function _isValidatorOptedIn(bytes memory consAddr) internal view returns (bool) {
        address eoa = storedConsAddrs[consAddr];
        bool isConsAddrStored = eoa != address(0);
        bool isEoaActive = whitelistedEOAs[eoa].state == State.Active;
        return isConsAddrStored && isEoaActive;
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}

