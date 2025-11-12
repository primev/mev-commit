// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {Options}  from "openzeppelin-foundry-upgrades/Options.sol";

import {RocketMinipoolRegistry} from "../../../contracts/validator-registry/rocketpool/RocketMinipoolRegistry.sol";
import {IRocketMinipoolRegistry} from "../../../contracts/interfaces/IRocketMinipoolRegistry.sol";
import {MinipoolStatus} from "rocketpool/contracts/types/MinipoolStatus.sol";

// ---------- Mocks ----------

contract RocketStorageMock {
    mapping(bytes32 => address) internal _addr;
    mapping(address => address) internal _withdrawal;

    function setMinipoolForPubkey(bytes calldata pk, address mp) external {
        _addr[keccak256(abi.encodePacked("validator.minipool", pk))] = mp;
    }

    function setNodeWithdrawalAddress(address node, address withdrawal) external {
        _withdrawal[node] = withdrawal;
    }

    // registry reads these
    function getAddress(bytes32 key) external view returns (address) {
        return _addr[key];
    }

    function getNodeWithdrawalAddress(address node) external view returns (address) {
        return _withdrawal[node];
    }
}

contract MinipoolMock {
    address public node;
    MinipoolStatus public status;

    constructor(address _node) {
        node = _node;
        status = MinipoolStatus.Staking; // default to active
    }

    function getNodeAddress() external view returns (address) {
        return node;
    }

    function getStatus() external view returns (MinipoolStatus) {
        return status;
    }

    function setStatus(MinipoolStatus s) external {
        status = s;
    }
}

// ---------- Tests ----------

contract RocketMinipoolRegistrySigTest is Test {
    // actors
    address internal owner;
    address internal oracle;
    address internal node;        // derived from nodeSk
    address internal withdrawal;
    address internal receiver;
    address internal stranger;

    // signer key
    uint256 internal nodeSk;

    // system
    RocketStorageMock internal storageMock;
    MinipoolMock internal mp1;
    RocketMinipoolRegistry internal reg;

    // params
    uint256 internal fee = 0.01 ether;
    uint64  internal period = 3 days;

    // sample 48-byte pubkeys (96 hex chars)
    bytes internal pk1 = hex"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

    // --- EIP-712 constants (must match contract) ---
    bytes32 internal constant EIP712_DOMAIN_TYPEHASH =
        keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");
    bytes32 internal constant NAME_HASH = keccak256(bytes("RocketMinipoolRegistry"));
    bytes32 internal constant VERSION_HASH = keccak256(bytes("1"));

    bytes32 internal constant REGISTER_TYPEHASH   =
        keccak256("Register(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");
    bytes32 internal constant DEREG_REQ_TYPEHASH  =
        keccak256("DeregRequest(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");
    bytes32 internal constant DEREG_TYPEHASH      =
        keccak256("Deregister(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");

    // helpers
    function _one(bytes memory pk) internal pure returns (bytes[] memory a) {
        a = new bytes[](1);
        a[0] = pk;
    }

    function _domainSeparator() internal view returns (bytes32) {
        return keccak256(abi.encode(
            EIP712_DOMAIN_TYPEHASH,
            NAME_HASH,
            VERSION_HASH,
            block.chainid,
            address(reg)
        ));
    }

    function _pubkeysHash(bytes[] memory pks) internal pure returns (bytes32) {
        // Concatenate 48-byte pubkeys and hash (matches contract intent)
        bytes memory c;
        for (uint256 i = 0; i < pks.length; ++i) {
            // tests rely on the production modifier enforcing 48-byte length; keep it simple here
            c = bytes.concat(c, pks[i]);
        }
        return keccak256(c);
    }

    function _sign(bytes32 structHash) internal view returns (bytes memory) {
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _domainSeparator(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(nodeSk, digest);
        return abi.encodePacked(r, s, v);
    }

    function setUp() public {
        deal(address(this), 100 ether); // for unfreeze tests

        // choose a real private key for node so we can sign
        nodeSk = 0xA11CE;
        node = vm.addr(nodeSk);

        // other actors
        owner      = makeAddr("owner");
        oracle     = makeAddr("freezeOracle");
        withdrawal = makeAddr("withdrawal");
        receiver   = makeAddr("receiver");
        stranger   = makeAddr("stranger");

        // mocks
        storageMock = new RocketStorageMock();
        mp1 = new MinipoolMock(node);

        // wire lookups
        storageMock.setMinipoolForPubkey(pk1, address(mp1));
        storageMock.setNodeWithdrawalAddress(node, withdrawal);

        Options memory opts;
        opts.unsafeSkipAllChecks = true;

        address proxy = Upgrades.deployUUPSProxy(
            "RocketMinipoolRegistry.sol",
            abi.encodeCall(
                RocketMinipoolRegistry.initialize,
                (owner, oracle, receiver, address(storageMock), fee, period)
            ),
            opts
        );
        reg = RocketMinipoolRegistry(payable(proxy));
    }

    // ---------- access: only withdrawal for non-sig paths ----------
    function test_Register_ByWithdrawal_Succeeds_NodeReverts() public {
        // node should NOT be allowed anymore
        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.NotWithdrawalAddress.selector, withdrawal));
        reg.registerValidators(_one(pk1));

        // withdrawal can call
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));
        assertTrue(reg.isValidatorRegistered(pk1));
    }

    function test_RequestDereg_ByWithdrawal_Succeeds_NodeReverts() public {
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));

        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.NotWithdrawalAddress.selector, withdrawal));
        reg.requestValidatorDeregistration(_one(pk1));

        vm.prank(withdrawal);
        reg.requestValidatorDeregistration(_one(pk1));
    }

    function test_Deregister_ByWithdrawal_Succeeds_NodeReverts() public {
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));

        vm.prank(withdrawal);
        reg.requestValidatorDeregistration(_one(pk1));

        vm.warp(block.timestamp + period + 1);

        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.NotWithdrawalAddress.selector, withdrawal));
        reg.deregisterValidators(_one(pk1));

        vm.prank(withdrawal);
        reg.deregisterValidators(_one(pk1));
        assertFalse(reg.isValidatorRegistered(pk1));
    }

    // ---------- WithSig happy paths ----------
    function test_Register_WithSig_ByIntendedExecutor_Works() public {
        // build struct hash with executor bound to withdrawal (who will call)
        uint256 nonce = 0; // fresh contract: per-node nonce starts at 0
        uint256 deadline = block.timestamp + 1 days;
        bytes32 pkHash = _pubkeysHash(_one(pk1));

        bytes32 structHash = keccak256(abi.encode(
            REGISTER_TYPEHASH,
            pkHash,
            withdrawal,
            nonce,
            deadline
        ));

        bytes memory sig = _sign(structHash);

        vm.prank(withdrawal);
        reg.registerValidatorsWithSig(_one(pk1), sig, deadline);
        assertTrue(reg.isValidatorRegistered(pk1));
    }

    function test_RequestDereg_WithSig_Works() public {
        // first register via withdrawal (non-sig path) to set state
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));

        uint256 nonce = 0; // first signature use for this node
        uint256 deadline = block.timestamp + 1 days;
        bytes32 pkHash = _pubkeysHash(_one(pk1));

        bytes32 structHash = keccak256(abi.encode(
            DEREG_REQ_TYPEHASH,
            pkHash,
            withdrawal,
            nonce,
            deadline
        ));
        bytes memory sig = _sign(structHash);

        vm.prank(withdrawal);
        reg.requestValidatorDeregistrationWithSig(_one(pk1), sig, deadline);
    }

    function test_Deregister_WithSig_Works() public {
        // setup: registered, requested, period elapsed
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));

        // Dereg request first (non-sig or sig). Use non-sig via withdrawal:
        vm.prank(withdrawal);
        reg.requestValidatorDeregistration(_one(pk1));
        vm.warp(block.timestamp + period + 1);

        uint256 nonce = 0; // still 0: requestDereg non-sig path doesn't consume node nonce
        uint256 deadline = block.timestamp + 1 days;
        bytes32 pkHash = _pubkeysHash(_one(pk1));

        bytes32 structHash = keccak256(abi.encode(
            DEREG_TYPEHASH,
            pkHash,
            withdrawal,
            nonce,
            deadline
        ));
        bytes memory sig = _sign(structHash);

        vm.prank(withdrawal);
        reg.deregisterValidatorsWithSig(_one(pk1), sig, deadline);
        assertFalse(reg.isValidatorRegistered(pk1));
    }

    // ---------- WithSig failures ----------
    function test_WithSig_InvalidSignature_Reverts() public {
        // correct message, but signed by the wrong key
        uint256 deadline = block.timestamp + 1 days;
        bytes32 pkHash = _pubkeysHash(_one(pk1));
        bytes32 structHash = keccak256(abi.encode(
            REGISTER_TYPEHASH,
            pkHash,
            withdrawal,
            0,
            deadline
        ));

        // sign with a different key
        uint256 badSk = 0xBEEF;
        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", _domainSeparator(), structHash));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(badSk, digest);
        bytes memory badSig = abi.encodePacked(r, s, v);

        vm.prank(withdrawal);
        vm.expectRevert(IRocketMinipoolRegistry.InvalidSignature.selector);
        reg.registerValidatorsWithSig(_one(pk1), badSig, deadline);
    }

    function test_WithSig_ExecutorMismatch_Reverts() public {
        // node signs for executor = withdrawal
        uint256 nonce = 0;
        uint256 deadline = block.timestamp + 1 days;
        bytes32 pkHash = _pubkeysHash(_one(pk1));
        bytes32 structHash = keccak256(abi.encode(
            REGISTER_TYPEHASH,
            pkHash,
            withdrawal,
            nonce,
            deadline
        ));
        bytes memory sig = _sign(structHash);

        // call from a different address -> digest won't match (bound to msg.sender)
        vm.prank(makeAddr("not-withdrawal"));
        vm.expectRevert(IRocketMinipoolRegistry.InvalidSignature.selector);
        reg.registerValidatorsWithSig(_one(pk1), sig, deadline);
    }

    function test_WithSig_ExpiredDeadline_Reverts() public {
        uint256 nonce = 0;
        uint256 deadline = block.timestamp - 1; // already expired
        bytes32 pkHash = _pubkeysHash(_one(pk1));
        bytes32 structHash = keccak256(abi.encode(
            REGISTER_TYPEHASH,
            pkHash,
            withdrawal,
            nonce,
            deadline
        ));
        bytes memory sig = _sign(structHash);

        vm.prank(withdrawal);
        vm.expectRevert(IRocketMinipoolRegistry.ExpiredSignature.selector);
        reg.registerValidatorsWithSig(_one(pk1), sig, deadline);
    }

    // ---------- Behavior parity checks ----------
    function test_IsValidatorOptedIn_ReflectsState() public {
        // not registered
        assertFalse(reg.isValidatorOptedIn(pk1));

        // registered + active
        vm.prank(withdrawal); reg.registerValidators(_one(pk1));
        assertTrue(reg.isValidatorOptedIn(pk1));

        // frozen -> false
        vm.prank(oracle); reg.freeze(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));
        reg.unfreeze{value: fee}(_one(pk1));

        // dereg requested -> false
        vm.prank(withdrawal); reg.requestValidatorDeregistration(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));

        // complete dereg clears
        vm.warp(block.timestamp + period + 1);
        vm.prank(withdrawal); reg.deregisterValidators(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));

        // inactive minipool -> false
        vm.prank(withdrawal); reg.registerValidators(_one(pk1));
        MinipoolMock(address(mp1)).setStatus(MinipoolStatus.Withdrawable);
        assertFalse(reg.isValidatorOptedIn(pk1));
    }

    // receive refunds
    receive() external payable {}
}
