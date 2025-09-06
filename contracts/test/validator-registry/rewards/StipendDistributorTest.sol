// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {StipendDistributor} from "../../../contracts/validator-registry/rewards/StipendDistributor.sol";
import {IStipendDistributor} from "../../../contracts/interfaces/IStipendDistributor.sol"; // events/types only

// helper registries & harnesses (same style/paths as your example tests)
import {VanillaRegistry} from "../../../contracts/validator-registry/VanillaRegistry.sol";
import {ValidatorOptInRouter} from "../../../contracts/validator-registry/ValidatorOptInRouter.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {VanillaRegistryTest} from "../VanillaRegistryTest.sol";
import {MevCommitAVSTest} from "../avs/MevCommitAVSTest.sol";
import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {MevCommitMiddlewareTestCont} from "../middleware/MevCommitMiddlewareTestCont.sol";
import {IVanillaRegistry} from "../../../contracts/interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../../../contracts/interfaces/IMevCommitAVS.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {PausableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/utils/PausableUpgradeable.sol";

contract StipendDistributorTest is Test {
    // system under test
    StipendDistributor internal distributor;

    // test registries used by tests
    VanillaRegistry public vanillaRegistry;
    VanillaRegistryTest public vanillaRegistryTest;
    MevCommitAVS public mevCommitAVS;
    MevCommitAVSTest public mevCommitAVSTest;
    MevCommitMiddleware public mevCommitMiddleware;
    MevCommitMiddlewareTestCont public mevCommitMiddlewareTest;

    // actors
    address internal owner;
    address internal oracle;
    address internal operator1;
    address internal operator2;
    address internal delegate1;
    address internal recipient1;
    address internal recipient2;
    address internal recipient3;

    // sample pubkey (48 bytes)
    bytes public samplePubkey1 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254b3f";

    // events from interface for expectEmit
    event RecipientSet(address indexed operator, bytes pubkey, uint256 indexed registryID, address indexed recipient);
    event StipendsGranted(address indexed operator, address indexed recipient, uint256 amount);
    event RewardsClaimed(address indexed operator, address indexed recipient, uint256 amount);
    event DefaultRecipientSet(address indexed operator, address indexed recipient);

    // setup: deploy registries + distributor and fund oracle for payable calls
        function setUp() public {
        // Test actors
        owner = address(0xA11CE);
        oracle = address(0x04AC1E);
        operator1 = address(0x111);
        operator2 = address(0x222);
        delegate1 = address(0xD311);
        recipient1 = address(0xAAA1);
        recipient2 = address(0xAAA2);
        recipient3 = address(0xAAA3);

        // Bring up helper test environments (they seed their internal state in setUp)
        vanillaRegistryTest = new VanillaRegistryTest();
        vanillaRegistryTest.setUp();
        vanillaRegistry = vanillaRegistryTest.validatorRegistry();

        mevCommitAVSTest = new MevCommitAVSTest();
        mevCommitAVSTest.setUp();
        mevCommitAVS = mevCommitAVSTest.mevCommitAVS();

        mevCommitMiddlewareTest = new MevCommitMiddlewareTestCont();
        mevCommitMiddlewareTest.setUp();
        mevCommitMiddleware = mevCommitMiddlewareTest.mevCommitMiddleware();

        // Deploy distributor proxy with registries
        StipendDistributor implementation = new StipendDistributor();
        bytes memory initData = abi.encodeCall(
            StipendDistributor.initialize,
            (owner, oracle, address(vanillaRegistry), address(mevCommitAVS), address(mevCommitMiddleware))
        );

        address proxy = address(new ERC1967Proxy(address(implementation), initData));
        distributor = StipendDistributor(payable(proxy));

        vm.deal(oracle, 1_000 ether); // for payable grant calls
    }

    // helper: grant three combos (op1→r1:1e, op1→r2:2e, op2→r3:3e)
    function _grantThreeCombos(
        address addr1,
        address addr2,
        address addr3,
        address op1,
        address op2
    ) internal {
        address[] memory operators = new address[](3);
        address[] memory receivers = new address[](3);
        uint256[] memory amounts = new uint256[](3);

        operators[0] = op1;
        receivers[0] = addr1;
        amounts[0] = 1 ether;

        operators[1] = op1;
        receivers[1] = addr2;
        amounts[1] = 2 ether;

        operators[2] = op2;
        receivers[2] = addr3;
        amounts[2] = 3 ether;

        vm.prank(oracle);
        distributor.grantStipends{value: amounts[0] + amounts[1] + amounts[2]}(operators, receivers, amounts);
    }

    // default recipient: set and read mapping
    function test_SetDefaultRecipient_setsMapping() public {
        // starts empty
        assertEq(distributor.defaultRecipient(operator1), address(0));

        // operator sets default
        vm.prank(operator1);
        distributor.setDefaultRecipient(recipient1);

        // mapping reflects default
        assertEq(distributor.defaultRecipient(operator1), recipient1);
    }

    // override by pubkey: same operator sets 3 keys → recipient2, then 2 keys → recipient3 (middleware registry id=2)
    function test_OverrideRecipientByPubkey_multipleBatches() public {
        // seed middleware validators for operator vm.addr(0x1117)
        mevCommitMiddlewareTest.test_registerValidators();
        address opFromMiddlewareTest = vm.addr(0x1117);

        // fetch 5 registered pubkeys
        bytes memory pubkey1 = mevCommitMiddlewareTest.sampleValPubkey1();
        bytes memory pubkey2 = mevCommitMiddlewareTest.sampleValPubkey2();
        bytes memory pubkey3 = mevCommitMiddlewareTest.sampleValPubkey3();
        bytes memory pubkey4 = mevCommitMiddlewareTest.sampleValPubkey4();
        bytes memory pubkey5 = mevCommitMiddlewareTest.sampleValPubkey5();

        // batch 1: 3 keys → recipient2
        bytes[] memory firstBatch = new bytes[](3);
        firstBatch[0] = pubkey1;
        firstBatch[1] = pubkey2;
        firstBatch[2] = pubkey3;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(firstBatch, 2, recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey1)), recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey2)), recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey3)), recipient2);

        // batch 2: 2 keys → recipient3
        bytes[] memory secondBatch = new bytes[](2);
        secondBatch[0] = pubkey4;
        secondBatch[1] = pubkey5;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(secondBatch, 2, recipient3);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey4)), recipient3);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey5)), recipient3);
    }

    // override by pubkey: reverts when caller isn't the registered operator
    function test_OverrideRecipientByPubkey_wrongOperator_reverts() public {
        // seed validators to establish key → operator mapping
        mevCommitMiddlewareTest.test_registerValidators();
        address rightfulOperator = vm.addr(0x1117);
        bytes memory pubkey = mevCommitMiddlewareTest.sampleValPubkey2();

        // different operator tries to override → revert
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = pubkey;
        vm.prank(operator2);
        vm.expectRevert();
        distributor.overrideRecipientByPubkey(pubs, 2, recipient1);

        // mapping unchanged
        assertEq(distributor.operatorKeyOverrides(rightfulOperator, keccak256(pubkey)), address(0));

        // rightful operator can set it
        vm.prank(rightfulOperator);
        distributor.overrideRecipientByPubkey(pubs, 2, recipient1);
        assertEq(distributor.operatorKeyOverrides(rightfulOperator, keccak256(pubkey)), recipient1);
    }

    // grantStipends: three combos accrue correctly (no claim here)
    function test_GrantStipends_threeCombos_setsAccrued() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        // accrued reflects grants
        assertEq(distributor.accrued(operator1, recipient1), 1 ether);
        assertEq(distributor.accrued(operator1, recipient2), 2 ether);
        assertEq(distributor.accrued(operator2, recipient3), 3 ether);
    }

    // claim: operator can claim; delegate can claim when authorized
    function test_Claim_byOperator_and_byDelegate() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        // operator1 claims 2e for recipient2
        address payable[] memory toClaim = new address payable[](1);
        toClaim[0] = payable(recipient2);
        uint256 r2Before = recipient2.balance;
        vm.prank(operator1);
        distributor.claimRewards(toClaim);
        assertEq(recipient2.balance, r2Before + 2 ether);

        // operator1 authorizes delegate for recipient1
        vm.prank(operator1);
        distributor.setClaimDelegate(delegate1, recipient1, true);

        // delegate claims 1e for recipient1
        address payable[] memory one = new address payable[](1);
        one[0] = payable(recipient1);
        uint256 r1Before = recipient1.balance;
        vm.prank(delegate1);
        distributor.claimOnbehalfOfOperator(operator1, one);
        assertEq(recipient1.balance, r1Before + 1 ether);
    }

    // claim: unauthorized caller cannot claim on behalf of another operator
    function test_ClaimOnBehalf_unauthorized_reverts() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address payable[] memory ask = new address payable[](1);
        ask[0] = payable(recipient3);

        // operator2 tries to claim as if for operator1 → revert
        vm.expectRevert();
        vm.prank(operator2);
        distributor.claimOnbehalfOfOperator(operator1, ask);
    }

    // pending rewards: increments on grant, clears on claim, and stacks across grants
    function test_PendingRewards_increment_and_clear() public {
        // 1) first grant (1e) to operator1→recipient1
        address[] memory ops = new address[](1);
        address[] memory recs = new address[](1);
        uint256[] memory amts = new uint256[](1);
        ops[0] = operator1;
        recs[0] = recipient1;
        amts[0] = 1 ether;
        vm.prank(oracle);
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);
        assertEq(distributor.accrued(operator1, recipient1), 1 ether);

        // claim pays 1e
        address payable[] memory list = new address payable[](1);
        list[0] = payable(recipient1);
        uint256 before = recipient1.balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(recipient1.balance, before + 1 ether);

        // immediate re-claim is no-op
        before = recipient1.balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(recipient1.balance, before);

        // 2) second grant (2e) → total accrued becomes 3e
        amts[0] = 2 ether;
        vm.prank(oracle);
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);
        assertEq(distributor.accrued(operator1, recipient1), 3 ether);

        // 3) third grant (3e) without claiming → total accrued becomes 6e
        amts[0] = 3 ether;
        vm.prank(oracle);
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);
        assertEq(distributor.accrued(operator1, recipient1), 6 ether);

        // claim now pays 5e (the unclaimed 2e + 3e)
        before = recipient1.balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(recipient1.balance, before + 5 ether);

        // re-claim still no-op
        before = recipient1.balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(recipient1.balance, before);
    }

    // getKeyRecipient: baseline → default → override (registry 0 routes to owning registry)
    function test_GetKeyRecipient_and_registry0_routing() public {
        // seed middleware so key belongs to operator vm.addr(0x1117) under registry id 2
        mevCommitMiddlewareTest.test_registerValidators();
        address opFromMiddlewareTest = vm.addr(0x1117);
        bytes memory key = mevCommitMiddlewareTest.sampleValPubkey2();

        // 1) baseline: no default/override → resolves to operator
        address rec0 = distributor.getKeyRecipient(key);
        assertEq(rec0, opFromMiddlewareTest, "registry 0 should resolve to owning operator");

        // 2) set default for operator → returns default
        vm.prank(opFromMiddlewareTest);
        distributor.setDefaultRecipient(recipient1);
        address rec1 = distributor.getKeyRecipient(key);
        assertEq(rec1, recipient1, "default recipient should be returned");

        // 3) set explicit override for this key → precedence over default
        bytes[] memory oneKey = new bytes[](1);
        oneKey[0] = key;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(oneKey, 2, recipient2);
        address rec2 = distributor.getKeyRecipient(key);
        assertEq(rec2, recipient2, "override should take precedence");
    }

    // pause: user funcs revert when paused; owner can pause/unpause; grant is blocked; unpause restores
    function test_Pause_allPausableFunctions() public {
        // works unpaused
        vm.prank(operator1);
        distributor.setDefaultRecipient(recipient1);

        // pause as owner
        vm.prank(owner);
        distributor.pause();
        assertTrue(distributor.paused());

        // pausable funcs revert when paused
        vm.prank(operator1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.setDefaultRecipient(recipient2);

        bytes[] memory pubs = new bytes[](1);
        pubs[0] = samplePubkey1;
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.overrideRecipientByPubkey(pubs, 3, recipient2);

        address[] memory ops = new address[](1);
        address[] memory recs = new address[](1);
        uint256[] memory amts = new uint256[](1);
        ops[0] = operator1;
        recs[0] = recipient1;
        amts[0] = 1 ether;
        vm.prank(oracle);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);

        address payable[] memory list = new address payable[](1);
        list[0] = payable(recipient1);
        vm.prank(operator1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.claimRewards(list);

        vm.prank(operator1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.setClaimDelegate(delegate1, recipient1, true);

        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.claimOnbehalfOfOperator(operator1, list);

        // unpause restores
        vm.prank(owner);
        distributor.unpause();
        vm.prank(operator1);
        distributor.setDefaultRecipient(recipient2);
    }

    // reentrancy: malicious recipient can't reenter claimRewards
    function test_ReentrancyGuard_onClaimRewards() public {
        // grant to a recipient that tries to reenter
        ReenteringRecipient attacker = new ReenteringRecipient();
        address[] memory ops = new address[](1);
        address[] memory recs = new address[](1);
        uint256[] memory amts = new uint256[](1);
        ops[0] = operator1;
        recs[0] = address(attacker);
        amts[0] = 1 ether;
        vm.prank(oracle);
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);

        // claim once → paid exactly once; inner call blocked by nonReentrant
        address payable[] memory list = new address payable[](1);
        list[0] = payable(address(attacker));
        uint256 before = address(attacker).balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(address(attacker).balance, before + 1 ether);
    }

    // avs path (registry id=1): default vs override precedence
    function test_AVS_Override_and_GetKeyRecipient() public {
        // seed avs validators (pod owner is 0x420 in harness)
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();
        address podOwner = address(0x420);
        bytes memory key = mevCommitAVSTest.sampleValPubkey2();

        // baseline → pod owner
        address base = distributor.getKeyRecipient(key);
        assertEq(base, podOwner, "avs baseline should return pod owner");

        // set default → returned
        vm.prank(podOwner);
        distributor.setDefaultRecipient(recipient1);
        address def = distributor.getKeyRecipient(key);
        assertEq(def, recipient1);

        // set override (id=1) → takes precedence
        bytes[] memory oneKey = new bytes[](1);
        oneKey[0] = key;
        vm.prank(podOwner);
        distributor.overrideRecipientByPubkey(oneKey, 1, recipient2);
        address over = distributor.getKeyRecipient(key);
        assertEq(over, recipient2);

        // mapping is scoped by operator
        assertEq(distributor.operatorKeyOverrides(podOwner, keccak256(key)), recipient2);
        assertEq(distributor.operatorKeyOverrides(operator2, keccak256(key)), address(0));
    }

    // vanilla path (registry id=3): default vs override precedence
    function test_Vanilla_Override_and_GetKeyRecipient() public {
        // seed a vanilla validator owned by vanillaRegistryTest.user1()
        vanillaRegistryTest.testSelfStake();
        address valOperator = vanillaRegistryTest.user1();
        bytes memory key = vanillaRegistryTest.user1BLSKey();

        // baseline → operator
        address base = distributor.getKeyRecipient(key);
        assertEq(base, valOperator, "vanilla baseline should return operator");

        // set default → returned
        vm.prank(valOperator);
        distributor.setDefaultRecipient(recipient1);
        address def = distributor.getKeyRecipient(key);
        assertEq(def, recipient1);

        // set override (id=3) → takes precedence
        bytes[] memory oneKey = new bytes[](1);
        oneKey[0] = key;
        vm.prank(valOperator);
        distributor.overrideRecipientByPubkey(oneKey, 3, recipient2);
        address over = distributor.getKeyRecipient(key);
        assertEq(over, recipient2);

        // mapping is scoped by operator
        assertEq(distributor.operatorKeyOverrides(valOperator, keccak256(key)), recipient2);
        assertEq(distributor.operatorKeyOverrides(operator2, keccak256(key)), address(0));
    }

    // wrong registry id: override must revert when ownership doesn't match
    function test_Override_wrongRegistryID_reverts() public {
        // seed middleware validators for operator vm.addr(0x1117)
        mevCommitMiddlewareTest.test_registerValidators();
        address mwOperator = vm.addr(0x1117);
        bytes memory key = mevCommitMiddlewareTest.sampleValPubkey2();
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = key;

        // wrong id (1 = avs) → revert
        vm.prank(mwOperator);
        vm.expectRevert();
        distributor.overrideRecipientByPubkey(pubs, 1, recipient1);

        // correct id (2 = middleware) → ok
        vm.prank(mwOperator);
        distributor.overrideRecipientByPubkey(pubs, 2, recipient1);
        assertEq(distributor.operatorKeyOverrides(mwOperator, keccak256(key)), recipient1);
    }

    // invalid pubkey length: revert
    function test_Override_invalidPubkeyLength_reverts() public {
        // length check happens first → caller doesn't matter
        bytes memory bad = hex"1234"; // 2 bytes, not 48
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = bad;
        vm.prank(operator1);
        vm.expectRevert();
        distributor.overrideRecipientByPubkey(pubs, 2, recipient1);
    }

    // only oracle can grant stipends
    function test_Grant_onlyOracle_revertsForOthers() public {
        address[] memory ops = new address[](1);
        address[] memory recs = new address[](1);
        uint256[] memory amts = new uint256[](1);
        ops[0] = operator1;
        recs[0] = recipient1;
        amts[0] = 1 ether;
        vm.deal(operator1, 10 ether);

        // non-oracle caller → revert
        vm.prank(operator1);
        vm.expectRevert();
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);
    }

    // wrong operator can't claim another operator's recipients
    function test_ClaimRewards_wrongOperator_reverts() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address payable[] memory list = new address payable[](1);
        list[0] = payable(recipient2);

        uint256 before = recipient2.balance;
        vm.prank(operator2);
        distributor.claimRewards(list);
        assertEq(recipient2.balance, before);
    }

    // grantStipends: arrays length mismatch reverts
    function test_Grant_arraysLengthMismatch_reverts() public {
        address[] memory ops = new address[](2);
        address[] memory recs = new address[](1);
        uint256[] memory amts = new uint256[](1);
        ops[0] = operator1;
        ops[1] = operator2;
        recs[0] = recipient1;
        amts[0] = 1 ether;

        vm.prank(oracle);
        vm.expectRevert();
        distributor.grantStipends{value: amts[0]}(ops, recs, amts);
    }

    // zero-address guards
    function test_SetDefaultRecipient_zero_reverts() public {
        vm.prank(operator1);
        vm.expectRevert();
        distributor.setDefaultRecipient(address(0));
    }

    function test_Override_zeroRecipient_reverts() public {
        mevCommitMiddlewareTest.test_registerValidators();
        address mwOperator = vm.addr(0x1117);
        bytes memory key = mevCommitMiddlewareTest.sampleValPubkey2();
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = key;
        vm.prank(mwOperator);
        vm.expectRevert();
        distributor.overrideRecipientByPubkey(pubs, 2, address(0));
    }

    // batch claim: multiple recipients in one call
    function test_Claim_batchMultipleRecipients() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address payable[] memory list = new address payable[](2);
        list[0] = payable(recipient1); // 1 ether
        list[1] = payable(recipient2); // 2 ether

        uint256 r1Before = recipient1.balance;
        uint256 r2Before = recipient2.balance;

        vm.prank(operator1);
        distributor.claimRewards(list);

        assertEq(recipient1.balance, r1Before + 1 ether);
        assertEq(recipient2.balance, r2Before + 2 ether);
    }
}

// recipient that attempts to re-enter claimRewards during payout
contract ReenteringRecipient {
    fallback() external payable {
        // try to re-enter claimRewards(address[])
        bytes memory data = abi.encodeWithSignature("claimRewards(address[])", _arr());
        (bool ok, ) = msg.sender.call(data); // blocked by nonReentrant
        ok; // silence warning
    }

    function _arr() internal view returns (address[] memory a) {
        a = new address[](1);
        a[0] = address(this);
    }
}
