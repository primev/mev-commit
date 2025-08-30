// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {LidoV3Middleware} from "../../../contracts/validator-registry/lido/LidoV3Middleware.sol";
import {MockVaultHub} from "./mocks/MockVaultHub.sol";
import {MockStakingVault} from "./mocks/MockStakingVault.sol";

contract LidoV3MiddlewareTest is Test {
    // env constants
    address constant OWNER             = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant UNFREEZE_RECEIVER = 0x1623fE21185c92BB43bD83741E226288B516134a;

    // periods: seconds (per your correction)
    uint256 constant DEREG_PERIOD     = 600;
    uint256 constant UNFREEZE_PERIOD  = 600;

    // econ
    uint256 constant SLASH_AMOUNT_WEI = 10; // capacity = totalValue / 10 wei
    uint256 constant UNFREEZE_FEE_WEI = 0;

    // actors
    address operator = address(0xa11ce);
    address bob    = address(0xb0bb1);

    // sut
    MockVaultHub hub;
    MockStakingVault stakingVault;
    LidoV3Middleware mw;

    address vaultAddr;
    address hubAddr;

    function setUp() public {
        // deploy mocks
        hub = new MockVaultHub();
        hubAddr = address(hub);

        stakingVault = new MockStakingVault(operator);
        vaultAddr = address(stakingVault);

        // deploy implementation
        LidoV3Middleware impl = new LidoV3Middleware();

        // initializer (note the two extra args: slashAmount, unfreezeFee)
        bytes memory initData = abi.encodeWithSelector(
            LidoV3Middleware.initialize.selector,
            OWNER,
            hubAddr,
            UNFREEZE_RECEIVER,
            DEREG_PERIOD,
            UNFREEZE_PERIOD,
            SLASH_AMOUNT_WEI,
            UNFREEZE_FEE_WEI
        );

        // deploy proxy
        ERC1967Proxy proxy = new ERC1967Proxy(address(impl), initData);
        mw = LidoV3Middleware(address(proxy));

        // baseline: connected & large value
        hub.setConnected(vaultAddr, true);
        hub.setTotalValue(vaultAddr, 1_000_000 ether);

        // whitelist operator
        vm.prank(OWNER);
        mw.setWhitelist(operator, true);
    }

    // helper to fabricate 48-byte “BLS pubkeys”
    function _pk(bytes32 seed) internal pure returns (bytes memory pk) {
        pk = new bytes(48);
        assembly {
            mstore(add(pk, 32), seed)
            mstore(add(pk, 64), seed)
        }
    }

    function _keys(uint256 n) internal pure returns (bytes[] memory arr) {
        arr = new bytes[](n);
        for (uint256 i = 0; i < n; ++i) {
            arr[i] = _pk(bytes32(i + 1));
        }
    }

    // -----------------------
    // core tests
    // -----------------------

    function test_OwnerAndHubConfigured() public {
        assertEq(mw.owner(), OWNER, "owner");
        assertEq(mw.vaultHub(), hubAddr, "hub");
        assertEq(mw.unfreezePeriod(), UNFREEZE_PERIOD);
        assertEq(mw.deregistrationPeriod(), DEREG_PERIOD);
        assertEq(mw.slashAmount(), SLASH_AMOUNT_WEI);
        assertEq(mw.unfreezeFee(), UNFREEZE_FEE_WEI);
    }

    function test_WhitelistBlocksNonWhitelisted() public {
        bytes[] memory keys = _keys(2);

        vm.prank(bob);
        vm.expectRevert(LidoV3Middleware.NotWhitelisted.selector);
        mw.registerValidators(vaultAddr, keys);
    }

    function test_WhitelistedOperatorCanRegister() public {
        bytes[] memory keys = _keys(2);

        vm.prank(operator);
        mw.registerValidators(vaultAddr, keys);

        // both keys opted-in
        assertTrue(mw.isValidatorOptedIn(keys[0]));
        assertTrue(mw.isValidatorOptedIn(keys[1]));

        // capacity still huge
        assertGt(mw.remainingRegistrable(vaultAddr), 0);
    }

    function test_RevertIfVaultNotConnected() public {
        hub.setConnected(vaultAddr, false);
        bytes[] memory keys = _keys(1);

        vm.prank(operator);
        vm.expectRevert(LidoV3Middleware.InvalidVault.selector);
        mw.registerValidators(vaultAddr, keys);
    }

    function test_RevertIfNodeOperatorMismatch() public {
        stakingVault.setNodeOperator(address(0xDEAD));
        bytes[] memory keys = _keys(1);

        vm.prank(operator);
        vm.expectRevert(abi.encodeWithSelector(LidoV3Middleware.NodeOperatorMismatch.selector, address(0xDEAD), operator));
        mw.registerValidators(vaultAddr, keys);
    }

    function test_RevertIfCapacityExceeded() public {
        // totalValue=15 wei, slash=10 wei -> floor(1.5)=1
        hub.setTotalValue(vaultAddr, 15);

        bytes[] memory two = _keys(2);
        vm.prank(operator);
        vm.expectRevert(abi.encodeWithSelector(LidoV3Middleware.CapacityExceeded.selector, 2, 1));
        mw.registerValidators(vaultAddr, two);

        bytes[] memory one = _keys(1);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, one);
        assertTrue(mw.isValidatorOptedIn(one[0]));
    }

    function test_InvalidBLSPubKeyLength() public {
        bytes[] memory bad = _keys(1);
        bad[0] = new bytes(47); // not 48 bytes

        vm.prank(operator);
        vm.expectRevert(abi.encodeWithSelector(LidoV3Middleware.InvalidBLSPubKeyLength.selector, 48, 47));
        mw.registerValidators(vaultAddr, bad);
    }

    function test_FreezeAndUnfreeze() public {
        // register 1 key
        bytes[] memory one = _keys(1);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, one);
        assertTrue(mw.isValidatorOptedIn(one[0]));

        // owner freezes → optedIn becomes false
        vm.prank(OWNER);
        mw.freezeValidators(one);
        assertFalse(mw.isValidatorOptedIn(one[0]));

        // set receiver to a plain EOA and fee = 1 wei
        address recv = address(0xA11CE);
        vm.prank(OWNER);
        mw.setUnfreezeReceiver(recv);
        vm.prank(OWNER);
        mw.setUnfreezeFee(1);

        // too soon → revert
        vm.expectRevert(LidoV3Middleware.UnfreezeTooSoon.selector);
        mw.unfreeze{value: 1}(one);

        // after period → success; exact fee delivered
        uint256 before = recv.balance;
        vm.warp(block.timestamp + UNFREEZE_PERIOD + 1);
        mw.unfreeze{value: 1}(one);
        assertEq(recv.balance, before + 1, "receiver got fee");
        assertTrue(mw.isValidatorOptedIn(one[0]), "unfrozen opted in again");
    }

    function test_Unfreeze_InsufficientFee_Reverts() public {
        // register & freeze
        bytes[] memory one = _keys(1);
        vm.prank(operator); mw.registerValidators(vaultAddr, one);
        vm.prank(OWNER);    mw.freezeValidators(one);

        // set receiver to any EOA; fee = 2 wei
        vm.prank(OWNER); mw.setUnfreezeReceiver(address(0xB0B));
        vm.prank(OWNER); mw.setUnfreezeFee(2);

        vm.warp(block.timestamp + UNFREEZE_PERIOD + 1);

        vm.expectRevert(abi.encodeWithSelector(LidoV3Middleware.UnfreezeFeeRequired.selector, uint256(2)));
        mw.unfreeze{value: 1}(one);
    }

    function test_Unfreeze_Refunds_Excess() public {
        // register & freeze
        bytes[] memory one = _keys(1);
        vm.prank(operator); mw.registerValidators(vaultAddr, one);
        vm.prank(OWNER);    mw.freezeValidators(one);

        // receiver: plain EOA
        address recv = address(0xC0FFEE);
        vm.prank(OWNER); mw.setUnfreezeReceiver(recv);
        vm.prank(OWNER); mw.setUnfreezeFee(3); // need 3 wei

        vm.warp(block.timestamp + UNFREEZE_PERIOD + 1);

        // call unfreeze from a funded EOA so refund succeeds
        address payer = address(0xD00D);
        vm.deal(payer, 100);
        uint256 payerBefore = payer.balance;
        uint256 recvBefore  = recv.balance;

        vm.prank(payer);
        mw.unfreeze{value: 10}(one); // pays 10, fee=3, refund=7 to payer (EOA can receive ETH)

        assertEq(recv.balance,  recvBefore + 3, "receiver got exact fee");
        assertEq(payer.balance, payerBefore - 3, "only fee retained; refund returned");
    }

    function test_Unfreeze_TransferFailed_Reverts() public {
        // register & freeze
        bytes[] memory one = _keys(1);
        vm.prank(operator); mw.registerValidators(vaultAddr, one);
        vm.prank(OWNER);    mw.freezeValidators(one);

        // set receiver to THIS contract, which has NO payable receive/fallback → low-level .call with value will revert
        vm.prank(OWNER); mw.setUnfreezeReceiver(address(this));
        vm.prank(OWNER); mw.setUnfreezeFee(1);

        vm.warp(block.timestamp + UNFREEZE_PERIOD + 1);
        vm.expectRevert(LidoV3Middleware.TransferFailed.selector);
        mw.unfreeze{value: 1}(one);
    }

    function test_Deregister_Flow_WithTiming() public {
        // register 2 keys
        bytes[] memory two = _keys(2);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, two);

        // request dereg → optedIn becomes false
        vm.prank(operator);
        mw.requestDeregistrations(two);
        assertFalse(mw.isValidatorOptedIn(two[0]));
        assertFalse(mw.isValidatorOptedIn(two[1]));

        // too early to deregister → revert
        vm.expectRevert(LidoV3Middleware.UnfreezeTooSoon.selector);
        vm.prank(operator);
        mw.deregisterValidators(vaultAddr, two);

        // after window → success; per-vault count decreases by 2
        uint256 before = mw.vaultRegisteredCount(vaultAddr);
        vm.warp(block.timestamp + DEREG_PERIOD + 1);
        vm.prank(operator);
        mw.deregisterValidators(vaultAddr, two);

        uint256 after_ = mw.vaultRegisteredCount(vaultAddr);
        assertEq(before - 2, after_, "count decreased");
        // Not registered anymore
        assertFalse(mw.isValidatorOptedIn(two[0]));
        assertFalse(mw.isValidatorOptedIn(two[1]));
    }

    function test_AdminSetters_OnlyOwner() public {
        // non-owner setWhitelist
        vm.prank(bob);
        vm.expectRevert(
            abi.encodeWithSelector(
                OwnableUpgradeable.OwnableUnauthorizedAccount.selector,
                bob
            )
        );
        mw.setWhitelist(bob, true);

        // owner actions
        vm.prank(OWNER);
        mw.setSlashAmount(123);
        assertEq(mw.slashAmount(), 123);

        vm.prank(OWNER);
        mw.setUnfreezeFee(5);
        assertEq(mw.unfreezeFee(), 5);

        vm.prank(OWNER);
        mw.setUnfreezePeriod(777);
        assertEq(mw.unfreezePeriod(), 777);

        vm.prank(OWNER);
        mw.setDeregistrationPeriod(888);
        assertEq(mw.deregistrationPeriod(), 888);
    }

    function test_OptedIn_False_IfRegistrarRemovedFromWhitelist() public {
        bytes[] memory one = _keys(1);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, one);
        assertTrue(mw.isValidatorOptedIn(one[0]));

        // remove operator from whitelist → becomes false
        vm.prank(OWNER);
        mw.setWhitelist(operator, false);
        assertFalse(mw.isValidatorOptedIn(one[0]));
    }

    function test_SwitchVaultHub_BlocksNewRegistrations() public {
        // allow one successful registration under current hub
        bytes[] memory one = _keys(1);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, one);

        // switch to a fresh hub mock that does NOT know our vault
        MockVaultHub newHub = new MockVaultHub();
        vm.prank(OWNER);
        mw.setVaultHub(address(newHub));

        // any new registration should fail InvalidVault
        bytes[] memory another = _keys(1);
        vm.prank(operator);
        vm.expectRevert(LidoV3Middleware.InvalidVault.selector);
        mw.registerValidators(vaultAddr, another);
    }

    function test_CapacityMath_RemainingRegistrable() public {
        // capacity = floor(totalValue / slashAmount)
        // make it small: 105 wei / 10 wei = 10 keys
        hub.setTotalValue(vaultAddr, 105);

        // none registered yet
        assertEq(mw.maxRegistrableByVault(vaultAddr), 10);
        assertEq(mw.remainingRegistrable(vaultAddr), 10);

        // register 3
        bytes[] memory three = _keys(3);
        vm.prank(operator);
        mw.registerValidators(vaultAddr, three);

        assertEq(mw.vaultRegisteredCount(vaultAddr), 3);
        assertEq(mw.remainingRegistrable(vaultAddr), 7);
    }

}
