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
// TODO: Talk to taylor about offchain db if needed to enable full list via events, include kant in convo
// TODO: generate go bindings for this contract, adapt existing binding creation script for v1 registry.
contract ReputationValReg is OwnableUpgradeable, UUPSUpgradeable {

    uint256 constant FUNC_ARG_ARRAY_LIMIT = 100;

    uint256 public maxConsAddrsPerEOA;
    uint256 public minFreezeBlocks;
    uint256 public unfreezeFee;

    enum State { NotWhitelisted, Active, Frozen }

    struct WhitelistedEOAInfo {
        State state;
        uint numConsAddrsStored;
        uint256 freezeHeight;
        string moniker; // Neccessary for delegators to identify the EOA
    }

    // Mapping of whitelisted EOAs to their info struct
    mapping(address => WhitelistedEOAInfo) public whitelistedEOAs;

    // Set of stored validator consensus addresses with O(1) lookup, mapped to the whitelisted EOA
    // that originally stored the consensus address. The 'address' value may or may not correspond to 
    // an actively whitelisted EOA.
    // 
    // This mapping is intentionally not enumerable,
    // since actors should only need to query the 32 relevant proposers for an epoch at a time.
    // If for some reason an actor desires the full set of stored validator cons addrs,
    // they could construct the set offchain via events.
    mapping(bytes => address) public storedConsAddrs;

    event WhitelistedEOAAdded(address indexed eoa, string moniker);
    event WhitelistedEOADeleted(address indexed eoa, string moniker);
    event EOAFrozen(address indexed eoa, string moniker);
    event EOAUnfrozen(address indexed eoa, string moniker);
    event ConsAddrStored(bytes consAddr, address indexed eoa, string moniker);
    event ConsAddrDeleted(bytes consAddr, address indexed eoa, string moniker);

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

    function addWhitelistedEOA(address eoa, string memory moniker) external onlyOwner {
        require(eoa != address(0), "Invalid address");
        require(!isEOAWhitelisted(eoa), "EOA must not already be whitelisted");
        whitelistedEOAs[eoa] = WhitelistedEOAInfo({
            state: State.Active,
            numConsAddrsStored: 0,
            freezeHeight: 0,
            moniker: moniker
        });
        emit WhitelistedEOAAdded(eoa, moniker);
    }

    function deleteWhitelistedEOA(address eoa) external {
        require(msg.sender == owner() || msg.sender == eoa, "Only owner or EOA itself can delete whitelisted EOA");
        require(isEOAWhitelisted(eoa), "EOA must be whitelisted");
        string memory moniker = whitelistedEOAs[eoa].moniker;
        delete whitelistedEOAs[eoa];
        emit WhitelistedEOADeleted(eoa, moniker);
    }

    function freeze(bytes memory validatorConsAddr) onlyOwner external {
        address eoa = storedConsAddrs[validatorConsAddr];
        require(eoa != address(0), "Validator consensus address must be stored");
        require(whitelistedEOAs[eoa].state == State.Active, "EOA representing validator must be active");
        whitelistedEOAs[eoa].state = State.Frozen;
        whitelistedEOAs[eoa].freezeHeight = block.number;
        emit EOAFrozen(eoa, whitelistedEOAs[eoa].moniker);
    }

    function unfreeze() external payable {
        require(whitelistedEOAs[msg.sender].state == State.Frozen, "Sender must be frozen");
        require(block.number >= whitelistedEOAs[msg.sender].freezeHeight + minFreezeBlocks, "Freeze period has not elapsed");
        require(msg.value >= unfreezeFee, "Insufficient unfreeze fee");
        whitelistedEOAs[msg.sender].state = State.Active;
        whitelistedEOAs[msg.sender].freezeHeight = 0;
        emit EOAUnfrozen(msg.sender, whitelistedEOAs[msg.sender].moniker);
    }

    function storeConsAddrs(bytes[] memory consAddrs) external {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        require(isEOAWhitelisted(msg.sender), "Sender must be whitelisted");
        for (uint i = 0; i < consAddrs.length; i++) {
            require(storedConsAddrs[consAddrs[i]] == address(0), "Duplicate consensus address is already stored");
            require(whitelistedEOAs[msg.sender].numConsAddrsStored < maxConsAddrsPerEOA, "EOA must not store more than max allowed cons addrs");
            storedConsAddrs[consAddrs[i]] = msg.sender;
            whitelistedEOAs[msg.sender].numConsAddrsStored++;
            emit ConsAddrStored(consAddrs[i], msg.sender, whitelistedEOAs[msg.sender].moniker);
        }
    }

    function deleteConsAddrs(bytes[] memory consAddrs) external {
        require(consAddrs.length <= FUNC_ARG_ARRAY_LIMIT, "Too many cons addrs in request. Try batching");
        for (uint i = 0; i < consAddrs.length; i++) {
            address eoa = storedConsAddrs[consAddrs[i]];
            require(eoa != address(0), "Consensus address must be stored");
            if (isEOAWhitelisted(eoa)) {
                require(eoa == msg.sender, "Consensus address must be originally stored by sender");
                whitelistedEOAs[eoa].numConsAddrsStored--;
                _deleteConsAddr(consAddrs[i], eoa, whitelistedEOAs[eoa].moniker);
            } else {
                _deleteConsAddr(consAddrs[i], eoa, "non-whitelisted EOA");
            }
        }
    }

    function isEOAWhitelisted(address eoa) public view returns (bool) {
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

    function getWhitelistedEOAInfo(address eoa) external view returns (ReputationValReg.State, uint, uint256, string memory) {
        return (whitelistedEOAs[eoa].state, whitelistedEOAs[eoa].numConsAddrsStored,
            whitelistedEOAs[eoa].freezeHeight, whitelistedEOAs[eoa].moniker);
    }

    function _deleteConsAddr(bytes memory consAddr, address eoa, string memory moniker) internal {
        delete storedConsAddrs[consAddr];
        emit ConsAddrDeleted(consAddr, eoa, moniker);
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

