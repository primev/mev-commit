// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

import {StipendDistributor} from "../../../contracts/validator-registry/rewards/StipendDistributor.sol";
import {IStipendDistributor} from "../../../contracts/interfaces/IStipendDistributor.sol"; // events/types only

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {PausableUpgradeable} from "openzeppelin-contracts-upgradeable/contracts/utils/PausableUpgradeable.sol";

contract StipendDistributorTest is Test {
    // system under test
    StipendDistributor internal distributor;

    // actors
    address internal owner;
    address internal stipendManager;
    address internal operator1;
    address internal operator2;
    address internal delegate1;
    address internal recipient1;
    address internal recipient2;
    address internal recipient3;

    // sample 48-byte pubkeys
    bytes internal pubkey1 = hex"b61a6e5f09217278efc7ddad4dc4b0553b2c076d4a5fef6509c233a6531c99146347193467e84eb5ca921af1b8254b3f";
    bytes internal pubkey2 = hex"aca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes internal pubkey3 = hex"cca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes internal pubkey4 = hex"dca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
    bytes internal pubkey5 = hex"eca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";

    // events from interface for expectEmit
    event RecipientSet(address indexed operator, bytes pubkey, uint256 indexed registryID, address indexed recipient);
    event StipendsGranted(address indexed operator, address indexed recipient, uint256 amount);
    event RewardsClaimed(address indexed operator, address indexed recipient, uint256 amount);
    event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient);
    event StipendsReclaimed(address indexed operator, address indexed recipient, uint256 amount);


    // setup: deploy registries + distributor and fund stipendManager for payable calls
        function setUp() public {
        // Test actors
        owner = address(0xA11CE);
        stipendManager = address(0x04AC1E);
        operator1 = address(0x111);
        operator2 = address(0x222);
        delegate1 = address(0xD311);
        recipient1 = address(0xAAA1);
        recipient2 = address(0xAAA2);
        recipient3 = address(0xAAA3);

        // Deploy distributor proxy
        StipendDistributor implementation = new StipendDistributor();
        bytes memory initData = abi.encodeCall(
            StipendDistributor.initialize,
            (owner, stipendManager)
        );

        address proxy = address(new ERC1967Proxy(address(implementation), initData));
        distributor = StipendDistributor(payable(proxy));

        vm.deal(stipendManager, 1_000 ether); // for payable grant calls
    }

    // helper: grant three combos (op1→r1:1e, op1→r2:2e, op2→r3:3e)
    function _grantThreeCombos(
        address addr1,
        address addr2,
        address addr3,
        address op1,
        address op2
    ) internal {
        IStipendDistributor.Stipend[] memory stipends = new IStipendDistributor.Stipend[](3);

        stipends[0].operator = op1;
        stipends[0].recipient = addr1;
        stipends[0].amount = 1 ether;

        stipends[1].operator = op1;
        stipends[1].recipient = addr2;
        stipends[1].amount = 2 ether;

        stipends[2].operator = op2;
        stipends[2].recipient = addr3;
        stipends[2].amount = 3 ether;

        vm.prank(stipendManager);
        distributor.grantStipends{value: stipends[0].amount + stipends[1].amount + stipends[2].amount}(stipends);
    }

    // default recipient: set and read mapping
    function test_SetOperatorGlobalOverride_setsMapping() public {
        // starts empty
        assertEq(distributor.operatorGlobalOverride(operator1), address(0));

        // operator sets default
        vm.prank(operator1);
        distributor.setOperatorGlobalOverride(recipient1);

        // mapping reflects default
        assertEq(distributor.operatorGlobalOverride(operator1), recipient1);
    }

    // override by pubkey: same operator sets 3 keys → recipient2, then 2 keys → recipient3 (middleware registry id=2)
    function test_OverrideRecipientByPubkey_multipleBatches() public {
        address opFromMiddlewareTest = vm.addr(0x1117);

        // batch 1: 3 keys → recipient2
        bytes[] memory firstBatch = new bytes[](3);
        firstBatch[0] = pubkey1;
        firstBatch[1] = pubkey2;
        firstBatch[2] = pubkey3;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(firstBatch, recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey1)), recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey2)), recipient2);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey3)), recipient2);

        // batch 2: 2 keys → recipient3
        bytes[] memory secondBatch = new bytes[](2);
        secondBatch[0] = pubkey4;
        secondBatch[1] = pubkey5;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(secondBatch, recipient3);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey4)), recipient3);
        assertEq(distributor.operatorKeyOverrides(opFromMiddlewareTest, keccak256(pubkey5)), recipient3);
    }

    // override by pubkey: reverts when caller isn't the registered operator
    function test_OverrideRecipientByPubkey_wrongOperator_reverts() public {
        address rightfulOperator = vm.addr(0x1117);
        // rightful operator can set it
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = pubkey1;
        vm.prank(rightfulOperator);
        distributor.overrideRecipientByPubkey(pubs, recipient1);
        assertEq(distributor.operatorKeyOverrides(rightfulOperator, keccak256(pubkey1)), recipient1);
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
        address[] memory toClaim = new address[](1);
        toClaim[0] = recipient2;
        uint256 r2Before = recipient2.balance;
        vm.prank(operator1);
        distributor.claimRewards(toClaim);
        assertEq(recipient2.balance, r2Before + 2 ether);

        // operator1 authorizes delegate for recipient1
        vm.prank(operator1);
        distributor.setClaimDelegate(delegate1, recipient1, true);

        // delegate claims 1e for recipient1
        address[] memory one = new address[](1);
        one[0] = recipient1;
        uint256 r1Before = recipient1.balance;
        vm.prank(delegate1);
        distributor.claimOnbehalfOfOperator(operator1, one);
        assertEq(recipient1.balance, r1Before + 1 ether);
    }

    // claim: unauthorized caller cannot claim on behalf of another operator
    function test_ClaimOnBehalf_unauthorized_reverts() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address[] memory ask = new address[](1);
        ask[0] = recipient3;

        // operator2 tries to claim as if for operator1 → revert
        vm.expectRevert();
        vm.prank(operator2);
        distributor.claimOnbehalfOfOperator(operator1, ask);
    }

    // pending rewards: increments on grant, clears on claim, and stacks across grants
    function test_PendingRewards_increment_and_clear() public {
        // 1) first grant (1e) to operator1→recipient1
        IStipendDistributor.Stipend[] memory stipends = new IStipendDistributor.Stipend[](1);
        stipends[0].operator = operator1;
        stipends[0].recipient = recipient1;
        stipends[0].amount = 1 ether;
        vm.prank(stipendManager);
        distributor.grantStipends{value: stipends[0].amount}(stipends);
        assertEq(distributor.accrued(operator1, recipient1), 1 ether);

        // claim pays 1e
        address[] memory list = new address[](1);
        list[0] = recipient1;
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
        stipends[0].amount = 2 ether;
        vm.prank(stipendManager);
        distributor.grantStipends{value: stipends[0].amount}(stipends);
        assertEq(distributor.accrued(operator1, recipient1), 3 ether);

        // 3) third grant (3e) without claiming → total accrued becomes 6e
        stipends[0].amount = 3 ether;
        vm.prank(stipendManager);
        distributor.grantStipends{value: stipends[0].amount}(stipends);
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
        address opFromMiddlewareTest = vm.addr(0x1117);

        // 1) baseline: no default/override → resolves to operator
        address rec0 = distributor.getKeyRecipient(opFromMiddlewareTest, pubkey1);
        assertEq(rec0, opFromMiddlewareTest, "registry 0 should resolve to owning operator");

        // 2) set default for operator → returns default
        vm.prank(opFromMiddlewareTest);
        distributor.setOperatorGlobalOverride(recipient1);
        address rec1 = distributor.getKeyRecipient(opFromMiddlewareTest, pubkey1);
        assertEq(rec1, recipient1, "default recipient should be returned");

        // 3) set explicit override for this key → precedence over default
        bytes[] memory oneKey = new bytes[](1);
        oneKey[0] = pubkey1;
        vm.prank(opFromMiddlewareTest);
        distributor.overrideRecipientByPubkey(oneKey, recipient2);
        address rec2 = distributor.getKeyRecipient(opFromMiddlewareTest, pubkey1);
        assertEq(rec2, recipient2, "override should take precedence");
    }

    // pause: user funcs revert when paused; owner can pause/unpause; grant is blocked; unpause restores
    function test_Pause_allPausableFunctions() public {
        // works unpaused
        vm.prank(operator1);
        distributor.setOperatorGlobalOverride(recipient1);

        // pause as owner
        vm.prank(owner);
        distributor.pause();
        assertTrue(distributor.paused());

        // pausable funcs revert when paused
        vm.prank(operator1);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.setOperatorGlobalOverride(recipient2);

        bytes[] memory pubs = new bytes[](1);
        pubs[0] = pubkey1;
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.overrideRecipientByPubkey(pubs, recipient2);

        IStipendDistributor.Stipend[] memory stipends = new IStipendDistributor.Stipend[](1);
        stipends[0].operator = operator1;
        stipends[0].recipient = recipient1;
        stipends[0].amount = 1 ether;
        vm.prank(stipendManager);
        vm.expectRevert(PausableUpgradeable.EnforcedPause.selector);
        distributor.grantStipends{value: stipends[0].amount}(stipends);

        address[] memory list = new address[](1);
        list[0] = recipient1;
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
        distributor.setOperatorGlobalOverride(recipient2);
    }

    // reentrancy: malicious recipient can't reenter claimRewards
    function test_ReentrancyGuard_onClaimRewards() public {
        // grant to a recipient that tries to reenter
        ReenteringRecipient attacker = new ReenteringRecipient();
        IStipendDistributor.Stipend[] memory stipends = new IStipendDistributor.Stipend[](1);
        stipends[0].operator = operator1;
        stipends[0].recipient = address(attacker);
        stipends[0].amount = 1 ether;
        vm.prank(stipendManager);
        distributor.grantStipends{value: stipends[0].amount}(stipends);

        // claim once → paid exactly once; inner call blocked by nonReentrant
        address[] memory list = new address[](1);
        list[0] = address(attacker);
        uint256 before = address(attacker).balance;
        vm.prank(operator1);
        distributor.claimRewards(list);
        assertEq(address(attacker).balance, before + 1 ether);
    }

    function test_OverrideByPubkey() public {
        address mwOperator = vm.addr(0x1117);
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = pubkey1;

        vm.prank(mwOperator);
        distributor.overrideRecipientByPubkey(pubs, recipient1);
        assertEq(distributor.operatorKeyOverrides(mwOperator, keccak256(pubkey1)), recipient1);
    }

    function test_OverrideByPubkeyFailsOnInvalidPubkeyLength() public {
        bytes memory bad = hex"1234"; // 2 bytes, not 48
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = bad;
        vm.prank(operator1);
        vm.expectRevert(IStipendDistributor.InvalidBLSPubKeyLength.selector);
        distributor.overrideRecipientByPubkey(pubs, recipient1);
    }

    // only stipendManager can grant stipends
    function test_Grant_onlystipendManager_revertsForOthers() public {
        IStipendDistributor.Stipend[] memory stipends = new IStipendDistributor.Stipend[](1);
        stipends[0].operator = operator1;
        stipends[0].recipient = recipient1;
        stipends[0].amount = 1 ether;
        vm.deal(operator1, 10 ether);

        // non-stipendManager caller → revert
        vm.prank(operator1);
        vm.expectRevert();
        distributor.grantStipends{value: stipends[0].amount}(stipends);
    }

    // wrong operator can't claim another operator's recipients
    function test_ClaimRewards_wrongOperator_reverts() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address[] memory list = new address[](1);
        list[0] = payable(recipient2);

        uint256 before = recipient2.balance;
        vm.prank(operator2);
        distributor.claimRewards(list);
        assertEq(recipient2.balance, before);
    }

    // zero-address guards
    function test_SetOperatorGlobalOverride_zero_reverts() public {
        vm.prank(operator1);
        vm.expectRevert();
            distributor.setOperatorGlobalOverride(address(0));
    }

    function test_Override_zeroRecipient_reverts() public {
        address mwOperator = vm.addr(0x1117);
        bytes[] memory pubs = new bytes[](1);
        pubs[0] = pubkey1;
        vm.prank(mwOperator);
        vm.expectRevert();
        distributor.overrideRecipientByPubkey(pubs, address(0));
    }

    // batch claim: multiple recipients in one call
    function test_Claim_batchMultipleRecipients() public {
        _grantThreeCombos(recipient1, recipient2, recipient3, operator1, operator2);

        address[] memory list = new address[](2);
        list[0] = payable(recipient1); // 1 ether
        list[1] = payable(recipient2); // 2 ether

        uint256 r1Before = recipient1.balance;
        uint256 r2Before = recipient2.balance;

        vm.prank(operator1);
        distributor.claimRewards(list);

        assertEq(recipient1.balance, r1Before + 1 ether);
        assertEq(recipient2.balance, r2Before + 2 ether);
    }

    // Helper: grant arbitrary pairs in one go (prank as stipendManager and fund it)
    function _grantPairs(address[] memory ops, address[] memory recs, uint256[] memory amts) internal {
        uint256 len = ops.length;
        assertEq(len, recs.length, "setup: length mismatch");
        assertEq(len, amts.length, "setup: length mismatch");

        IStipendDistributor.Stipend[] memory s = new IStipendDistributor.Stipend[](len);
        uint256 total = 0;
        for (uint256 i = 0; i < len; ++i) {
            s[i] = IStipendDistributor.Stipend({ operator: ops[i], recipient: recs[i], amount: amts[i] });
            total += amts[i];
        }

        vm.deal(stipendManager, stipendManager.balance + total);
        vm.prank(stipendManager);
        distributor.grantStipends{value: total}(s);
    }

    function test_reclaimGrantsToOwner() public {
        address opA = makeAddr("opA");
        address rcA = makeAddr("rcA");
        address opB = makeAddr("opB");
        address rcB = makeAddr("rcB");

        address[] memory ops = new address[](2);
        address[] memory recs = new address[](2);
        uint256[] memory amts = new uint256[](2);
        ops[0] = opA; recs[0] = rcA; amts[0] = 5 ether;
        ops[1] = opB; recs[1] = rcB; amts[1] = 7 ether;

        _grantPairs(ops, recs, amts);

        address ownerAddr = distributor.owner();
        uint256 ownerBefore = ownerAddr.balance;
        uint256 contractBefore = address(distributor).balance;

        // (optional) check events
        vm.expectEmit(true, true, false, true);
        emit StipendsReclaimed(opA, rcA, 5 ether);
        vm.expectEmit(true, true, false, true);
        emit StipendsReclaimed(opB, rcB, 7 ether);

        vm.prank(ownerAddr);
        distributor.reclaimGrantsToOwner(ops, recs);

        assertEq(ownerAddr.balance, ownerBefore + 12 ether, "owner did not receive reclaimed ETH");
        assertEq(address(distributor).balance, contractBefore - 12 ether, "contract balance mismatch");
        assertEq(distributor.getPendingRewards(opA, rcA), 0, "pending not cleared for A");
        assertEq(distributor.getPendingRewards(opB, rcB), 0, "pending not cleared for B");
    }

    function test_reclaimGrantsToOwner_RespectsClaimedAndReclaimsOnlyUnclaimed() public {
        address operator = makeAddr("op");
        address recipient = makeAddr("rc");
        address ownerAddr = distributor.owner();

        // ---------- setup: grant #1 (5 ether) ----------
        IStipendDistributor.Stipend[] memory s1 = new IStipendDistributor.Stipend[](1);
        s1[0] = IStipendDistributor.Stipend({operator: operator, recipient: recipient, amount: 5 ether});

        vm.deal(ownerAddr, ownerAddr.balance + 5 ether);
        vm.prank(ownerAddr);
        distributor.grantStipends{value: 5 ether}(s1);

        // operator fully claims grant #1
        address[] memory recs = new address[](1);
        recs[0] = recipient;
        uint256 rcBefore = recipient.balance;
        uint256 contractBeforeClaim = address(distributor).balance;

        vm.prank(operator);
        distributor.claimRewards(recs);

        assertEq(recipient.balance, rcBefore + 5 ether, "recipient did not receive claim #1");
        assertEq(address(distributor).balance, contractBeforeClaim - 5 ether, "contract balance mismatch after claim #1");
        assertEq(distributor.getPendingRewards(operator, recipient), 0, "pending should be zero after claim #1");

        // ---------- part A: reclaim when fully claimed -> revert & no payout ----------
        uint256 ownerBeforeReclaimA = ownerAddr.balance;
        uint256 contractBeforeReclaimA = address(distributor).balance;

        address[] memory opsA = new address[](1);
        address[] memory recsA = new address[](1);
        opsA[0] = operator;
        recsA[0] = recipient;

        vm.prank(ownerAddr);
        vm.expectRevert(abi.encodeWithSignature("NoClaimableRewards(address,address)", ownerAddr, ownerAddr));
        distributor.reclaimGrantsToOwner(opsA, recsA);

        // balances unchanged
        assertEq(ownerAddr.balance, ownerBeforeReclaimA, "owner balance changed on failed reclaim");
        assertEq(address(distributor).balance, contractBeforeReclaimA, "contract balance changed on failed reclaim");

        // ---------- grant #2 (9 ether) ----------
        IStipendDistributor.Stipend[] memory s2 = new IStipendDistributor.Stipend[](1);
        s2[0] = IStipendDistributor.Stipend({operator: operator, recipient: recipient, amount: 9 ether});

        vm.deal(ownerAddr, ownerAddr.balance + 9 ether);
        vm.prank(ownerAddr);
        distributor.grantStipends{value: 9 ether}(s2);

        assertEq(distributor.getPendingRewards(operator, recipient), 9 ether, "pending should equal grant #2");

        // ---------- part B: reclaim pulls only unclaimed (9 ether) ----------
        uint256 ownerBeforeReclaimB = ownerAddr.balance;
        uint256 contractBeforeReclaimB = address(distributor).balance;

        address[] memory opsB = new address[](1);
        address[] memory recsB = new address[](1);
        opsB[0] = operator;
        recsB[0] = recipient;

        vm.prank(ownerAddr);
        distributor.reclaimGrantsToOwner(opsB, recsB);

        assertEq(ownerAddr.balance, ownerBeforeReclaimB + 9 ether, "owner did not receive only unclaimed amount");
        assertEq(address(distributor).balance, contractBeforeReclaimB - 9 ether, "contract balance mismatch after reclaim");
        assertEq(distributor.getPendingRewards(operator, recipient), 0, "pending should be zero after reclaim");
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
