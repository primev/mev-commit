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

contract RocketMinipoolRegistryTest is Test {
    // actors
    address internal owner      = makeAddr("owner");
    address internal oracle     = makeAddr("freezeOracle");
    address internal node       = makeAddr("node");
    address internal withdrawal = makeAddr("withdrawal");
    address internal stranger   = makeAddr("stranger");
    address internal receiver   = makeAddr("receiver");

    // system
    RocketStorageMock internal storageMock;
    MinipoolMock internal mp1;
    RocketMinipoolRegistry internal reg;

    // params
    uint256 internal fee = 0.01 ether;
    uint64  internal period = 3 days;

    // sample 48-byte pubkeys (96 hex chars)
    bytes internal pk1 = hex"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
    bytes internal pk2 = hex"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbaaaaaaaabbbbb";

    function setUp() public {
        deal(address(this), 100 ether); // for unfreeze tests

        // mocks
        storageMock = new RocketStorageMock();
        mp1 = new MinipoolMock(node);

        // wire lookups
        storageMock.setMinipoolForPubkey(pk1, address(mp1));
        storageMock.setNodeWithdrawalAddress(node, withdrawal);

        Options memory opts;
        opts.unsafeSkipAllChecks = true; // or opts.unsafeSkipStorageCheck = true

        address proxy = Upgrades.deployUUPSProxy(
            "RocketMinipoolRegistry.sol",
            abi.encodeCall(
                RocketMinipoolRegistry.initialize,
                (owner, oracle, receiver, address(storageMock), fee, period)
            ),
            opts // ← IMPORTANT: pass opts here
        );
        reg = RocketMinipoolRegistry(payable(proxy));
    }

    // helpers
    function _one(bytes memory pk) internal pure returns (bytes[] memory a) {
        a = new bytes[](1);
        a[0] = pk;
    }

    // ---------- initializer, setters, pause ----------
    function test_Initialize_And_Setters_And_Pause() public {
        assertEq(reg.unfreezeFee(), fee);
        assertEq(reg.freezeOracle(), oracle);
        assertEq(reg.unfreezeReceiver(), receiver);
        assertEq(address(reg.rocketStorage()), address(storageMock));
        assertEq(reg.deregistrationPeriod(), period);

        vm.prank(owner); reg.setUnfreezeFee(2 ether);              assertEq(reg.unfreezeFee(), 2 ether);
        vm.prank(owner); reg.setFreezeOracle(stranger);            assertEq(reg.freezeOracle(), stranger);
        vm.prank(owner); reg.setUnfreezeReceiver(withdrawal);      assertEq(reg.unfreezeReceiver(), withdrawal);
        vm.prank(owner); reg.setRocketStorage(address(storageMock));assertEq(address(reg.rocketStorage()), address(storageMock));
        vm.prank(owner); reg.setDeregistrationPeriod(1 days);      assertEq(reg.deregistrationPeriod(), 1 days);

        vm.prank(owner); reg.pause();
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        reg.registerValidators(_one(pk1));

        vm.prank(owner); reg.unpause();
        vm.prank(node); reg.registerValidators(_one(pk1)); // ok
    }

    // ---------- register ----------
    function test_Register_ByNode_Succeeds() public {
        vm.expectEmit(false, true, false, true);
        emit IRocketMinipoolRegistry.ValidatorRegistered(pk1, node);
        vm.prank(node);
        reg.registerValidators(_one(pk1));

        // ✅ read struct, then fields
        IRocketMinipoolRegistry.ValidatorRegistration memory info = reg.getValidatorRegInfo(pk1);
        assertEq(info.freezeTimestamp, 0);
        assertEq(info.exists, true);

        assertTrue(reg.isValidatorRegistered(pk1));
        assertTrue(reg.isValidatorOptedIn(pk1));

        (address n, address w) = reg.getValidOperatorsForKey(pk1);
        assertEq(n, node);
        assertEq(w, withdrawal);
    }

    function test_Register_ByWithdrawal_Succeeds() public {
        vm.prank(withdrawal);
        reg.registerValidators(_one(pk1));
        assertTrue(reg.isValidatorRegistered(pk1));
    }

    function test_Register_InvalidKeyLength_Reverts() public {
        bytes[] memory bad = new bytes[](1);
        bad[0] = hex"01";
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.InvalidBLSPubKeyLength.selector, 48, 1));
        reg.registerValidators(bad);
    }

    function test_Register_MinipoolNotActive_Reverts() public {
        mp1.setStatus(MinipoolStatus.Withdrawable); // anything != Staking
        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.MinipoolNotActive.selector, pk1));
        reg.registerValidators(_one(pk1));
    }

    function test_Register_Twice_Reverts() public {
        vm.prank(node); reg.registerValidators(_one(pk1));
        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.ValidatorAlreadyRegistered.selector, pk1));
        reg.registerValidators(_one(pk1));
    }

    // ---------- freeze / unfreeze ----------
    function test_Freeze_OnlyOracle() public {
        vm.prank(node); reg.registerValidators(_one(pk1));

        vm.expectRevert(IRocketMinipoolRegistry.OnlyFreezeOracle.selector);
        reg.freeze(_one(pk1)); // caller = this

        vm.prank(oracle);
        reg.freeze(_one(pk1));
        IRocketMinipoolRegistry.ValidatorRegistration memory info = reg.getValidatorRegInfo(pk1);
        assertTrue(info.freezeTimestamp != 0);

        vm.prank(oracle);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.ValidatorAlreadyFrozen.selector, pk1));
        reg.freeze(_one(pk1));
    }

    function test_Unfreeze_RequiresFee_And_Refunds() public {
        // setup: register & freeze
        vm.prank(node); reg.registerValidators(_one(pk1));
        vm.prank(oracle); reg.freeze(_one(pk1));

        // underpay -> revert with required fee
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.UnfreezeFeeRequired.selector, fee));
        reg.unfreeze{value: fee - 1}(_one(pk1));

        // exact pay -> receiver gets fee, registry keeps nothing, sender net-spend = fee
        uint256 rBefore = receiver.balance;
        uint256 sBefore = address(this).balance;
        uint256 regBefore = address(reg).balance;

        reg.unfreeze{value: fee}(_one(pk1));

        assertEq(receiver.balance, rBefore + fee, "receiver should get exact fee");
        assertEq(address(this).balance, sBefore - fee, "sender net-spend should equal fee");
        assertEq(address(reg).balance, regBefore, "registry should not retain ETH");

        // re-freeze, then overpay -> refund everything above fee
        vm.prank(oracle); reg.freeze(_one(pk1));

        uint256 extra = 1 ether;
        rBefore = receiver.balance;
        sBefore = address(this).balance;
        regBefore = address(reg).balance;

        reg.unfreeze{value: fee + extra}(_one(pk1));

        assertEq(receiver.balance, rBefore + fee, "receiver should still get only fee");
        assertEq(address(this).balance, sBefore - fee, "sender net-spend still fee (extra refunded)");
        assertEq(address(reg).balance, regBefore, "registry balance unchanged");
    }

    function test_OwnerUnfreeze_NoFee() public {
        vm.prank(node); reg.registerValidators(_one(pk1));
        vm.prank(oracle); reg.freeze(_one(pk1));

        vm.prank(owner);
        reg.ownerUnfreeze(_one(pk1));
        IRocketMinipoolRegistry.ValidatorRegistration memory info = reg.getValidatorRegInfo(pk1);
        assertEq(info.freezeTimestamp, 0);
    }

    // ---------- deregistration flow ----------
    function test_RequestDereg_Then_Finalize() public {
        vm.prank(node); reg.registerValidators(_one(pk1));

        vm.expectEmit(false, true, false, true);
        emit IRocketMinipoolRegistry.ValidatorDeregistrationRequested(pk1, node);
        vm.prank(node);
        reg.requestValidatorDeregistration(_one(pk1));

        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.DeregRequestAlreadyExists.selector, pk1));
        reg.requestValidatorDeregistration(_one(pk1));

        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.DeregistrationTooSoon.selector, pk1));
        reg.deregisterValidators(_one(pk1));

        vm.prank(oracle); reg.freeze(_one(pk1));
        vm.warp(block.timestamp + period + 1);

        vm.prank(node);
        vm.expectRevert(abi.encodeWithSelector(IRocketMinipoolRegistry.FrozenValidatorCannotDeregister.selector, pk1));
        reg.deregisterValidators(_one(pk1));

        reg.unfreeze{value: fee}(_one(pk1));
        vm.prank(node);
        reg.deregisterValidators(_one(pk1));

        assertFalse(reg.isValidatorRegistered(pk1));
        IRocketMinipoolRegistry.ValidatorRegistration memory info2 = reg.getValidatorRegInfo(pk1);
        // after dereg, nodeAddress can be zeroed; rely on the isValidatorRegistered guard instead
        assertEq(info2.exists, false);
        assertEq(reg.getEligibleTimeForDeregistration(pk1), 0);
    }

    // ---------- getters / helpers ----------
    function test_Getters_Work() public {
        vm.prank(node); reg.registerValidators(_one(pk1));
        assertEq(reg.getNodeAddressFromPubkey(pk1), node);
        assertEq(reg.getMinipoolFromPubkey(pk1), address(mp1));

        (address n, address w) = reg.getValidOperatorsForKey(pk1);
        assertEq(n, node);
        assertEq(w, withdrawal);

        assertTrue(reg.isMinipoolActive(address(mp1)));

        // if your isOperatorValidForKey uses msg.sender, these pass; if it takes (addr, key), swap accordingly.
        vm.prank(node);      assertTrue(reg.isOperatorValidForKey(node, pk1));
        vm.prank(withdrawal);assertTrue(reg.isOperatorValidForKey(withdrawal,pk1));
        vm.prank(stranger);  assertFalse(reg.isOperatorValidForKey(stranger,pk1));
    }

    function test_IsValidatorOptedIn_TruthTable() public {
        // not registered
        assertFalse(reg.isValidatorOptedIn(pk1));

        // registered + active
        vm.prank(node); reg.registerValidators(_one(pk1));
        assertTrue(reg.isValidatorOptedIn(pk1));

        // frozen -> false
        vm.prank(oracle); reg.freeze(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));
        reg.unfreeze{value: fee}(_one(pk1));

        // dereg requested -> false
        vm.prank(node); reg.requestValidatorDeregistration(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));

        // complete dereg clears
        vm.warp(block.timestamp + period + 1);
        vm.prank(node); reg.deregisterValidators(_one(pk1));
        assertFalse(reg.isValidatorOptedIn(pk1));

        // inactive minipool -> false
        vm.prank(node); reg.registerValidators(_one(pk1));
        MinipoolMock(address(mp1)).setStatus(MinipoolStatus.Withdrawable);
        assertFalse(reg.isValidatorOptedIn(pk1));
    }

    // receive refunds
    receive() external payable {}
}
