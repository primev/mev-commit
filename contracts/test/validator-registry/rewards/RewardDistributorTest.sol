// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {RewardDistributor} from "../../../contracts/validator-registry/rewards/RewardDistributor.sol";
import {IRewardDistributor} from "../../../contracts/interfaces/IRewardDistributor.sol";
// Minimal mintable ERC20 for tests
contract ERC20Mintable is IERC20 {
    string public name;
    string public symbol;
    uint8 public immutable decimals = 18;

    uint256 public override totalSupply;
    mapping(address => uint256) public override balanceOf;
    mapping(address => mapping(address => uint256)) public override allowance;

    constructor(string memory tokenName, string memory tokenSymbol) {
        name = tokenName;
        symbol = tokenSymbol;
    }

    function transfer(address to, uint256 amount) external override returns (bool) {
        _move(msg.sender, to, amount);
        return true;
    }

    function approve(address spender, uint256 amount) external override returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address from, address to, uint256 amount) external override returns (bool) {
        uint256 allowedAmount = allowance[from][msg.sender];
        require(allowedAmount >= amount, "allowance");
        allowance[from][msg.sender] = allowedAmount - amount;
        _move(from, to, amount);
        return true;
    }

    function mint(address to, uint256 amount) external {
        totalSupply += amount;
        balanceOf[to] += amount;
        emit Transfer(address(0), to, amount);
    }

    function _move(address from, address to, uint256 amount) internal {
        require(balanceOf[from] >= amount, "balance");
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        emit Transfer(from, to, amount);
    }
}

// Name chosen to match `--match-contract RewardDistributor`
contract RewardDistributorTest is Test {
    RewardDistributor internal rewardDistributor;

    // Roles / actors
    address internal contractOwner      = address(0xA11CE);
    address internal rewardManager      = address(0xB0B);
    address internal operatorAlpha      = address(0xA0A);
    address internal operatorBeta       = address(0xB0B0);
    address internal claimDelegateAlpha = address(0xD1);

    // Recipients
    address internal recipientOne       = address(0x111);
    address internal recipientTwo       = address(0x222);
    address internal recipientThree     = address(0x333);

    // Tokens
    ERC20Mintable internal rewardTokenOne; // tokenId = 1
    ERC20Mintable internal rewardTokenTwo; // tokenId = 2

    // 48-byte BLS pubkeys
    bytes internal pubkeyOne = hex"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";
    bytes internal pubkeyTwo = hex"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb";

    // Convenience
    function _getPending(address operator, address recipient, uint256 tokenId) internal view returns (uint128) {
        return rewardDistributor.getPendingRewards(operator, recipient, tokenId);
    }

    function setUp() public {
        // Deploy implementation & proxy, then initialize
        RewardDistributor implementation = new RewardDistributor();
        bytes memory initializerData = abi.encodeWithSelector(
            RewardDistributor.initialize.selector,
            contractOwner,
            rewardManager
        );
        ERC1967Proxy proxy = new ERC1967Proxy(address(implementation), initializerData);
        rewardDistributor = RewardDistributor(payable(address(proxy)));

        // Deploy and register tokens
        rewardTokenOne = new ERC20Mintable("TokenOne", "T1");
        rewardTokenTwo = new ERC20Mintable("TokenTwo", "T2");

        vm.startPrank(contractOwner);
        rewardDistributor.setRewardToken(address(rewardTokenOne), 1);
        rewardDistributor.setRewardToken(address(rewardTokenTwo), 2);
        vm.stopPrank();

        // Fund manager with ERC20s (pulled during grant)
        rewardTokenOne.mint(rewardManager, 1_000 ether);
        rewardTokenTwo.mint(rewardManager, 500 ether);

        // Fund ETH
        vm.deal(rewardManager, 1_000 ether);
        vm.deal(operatorAlpha, 1 ether);
        vm.deal(operatorBeta, 1 ether);
        vm.deal(contractOwner, 1_000 ether);
    }

    // ───────────────────────── Helpers

    function _grantETHRewards(address caller, RewardDistributor.Distribution[] memory distributions) internal {
        uint256 totalGrantAmount = 0;
        for (uint256 i = 0; i < distributions.length; ++i) {
            totalGrantAmount += distributions[i].amount;
        }
        vm.prank(caller);
        rewardDistributor.grantETHRewards{value: totalGrantAmount}(distributions);
    }

    function _grantTokenRewards(address caller, RewardDistributor.Distribution[] memory distributions, uint256 tokenId) internal {
        uint256 totalGrantAmount = 0;
        for (uint256 i = 0; i < distributions.length; ++i) {
            totalGrantAmount += distributions[i].amount;
        }
        vm.startPrank(caller);
        IERC20(rewardDistributor.rewardTokens(tokenId)).approve(address(rewardDistributor), totalGrantAmount);
        rewardDistributor.grantTokenRewards(distributions, tokenId);
        vm.stopPrank();
    }

    function _distribution(address operator, address recipient, uint128 amount)
        internal
        pure
        returns (RewardDistributor.Distribution memory entry)
    {
        entry.operator  = operator;
        entry.recipient = recipient;
        entry.amount    = amount;
    }

    // ───────────────────────── Grants: ETH

    function test_grantETHRewards_accruesAndPartitionsByOperatorRecipient() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](3);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 10 ether);
        distributions[1] = _distribution(operatorAlpha, recipientTwo, 5 ether);
        distributions[2] = _distribution(operatorBeta,  recipientOne, 7 ether);

        _grantETHRewards(rewardManager, distributions);

        assertEq(_getPending(operatorAlpha, recipientOne, 0), 10 ether);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 5 ether);
        assertEq(_getPending(operatorBeta,  recipientOne, 0), 7 ether);
        assertEq(_getPending(operatorBeta,  recipientTwo, 0), 0);
    }

    function test_grantETHRewards_revertsOnMismatchedMsgValue() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);

        vm.prank(rewardManager);
        vm.expectRevert(); // IncorrectPaymentAmount
        rewardDistributor.grantETHRewards{value: 0}(distributions);
    }

    function test_grantETHRewards_onlyOwnerOrManager() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);

        vm.expectRevert(); // NotOwnerOrRewardManager
        rewardDistributor.grantETHRewards{value: 1 ether}(distributions);

        vm.prank(contractOwner);
        rewardDistributor.grantETHRewards{value: 1 ether}(distributions); // ok
    }

    // ───────────────────────── Grants: Token

    function test_grantTokenRewards_accruesAndPullsFromCaller_tokenId1() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](2);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 100 ether);
        distributions[1] = _distribution(operatorAlpha, recipientTwo, 50 ether);

        _grantTokenRewards(rewardManager, distributions, 1);

        assertEq(rewardTokenOne.balanceOf(address(rewardDistributor)), 150 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 100 ether);
        assertEq(_getPending(operatorAlpha, recipientTwo, 1), 50 ether);
    }

    function test_grantTokenRewards_revertsIfTokenNotRegistered() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);

        vm.startPrank(rewardManager);
        rewardTokenOne.approve(address(rewardDistributor), type(uint256).max);
        vm.expectRevert(); // InvalidRewardToken
        rewardDistributor.grantTokenRewards(distributions, 9_999);
        vm.stopPrank();
    }

    // ───────────────────────── Claims: operator self-claim

    function test_claimRewards_byOperator_ETH_transfersAndZeroesPending() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](2);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
        distributions[1] = _distribution(operatorAlpha, recipientTwo, 3 ether);
        _grantETHRewards(contractOwner, distributions);

        address[] memory recipientsToClaim = new address[](2);
        recipientsToClaim[0] = recipientOne;
        recipientsToClaim[1] = recipientTwo;

        uint256 recipientOneBalanceBefore = recipientOne.balance;
        uint256 recipientTwoBalanceBefore = recipientTwo.balance;

        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);

        assertEq(recipientOne.balance, recipientOneBalanceBefore + 2 ether);
        assertEq(recipientTwo.balance, recipientTwoBalanceBefore + 3 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 0);
    }

    function test_claimRewards_byOperator_Token_transfersAndZeroesPending() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](2);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 20 ether);
        distributions[1] = _distribution(operatorAlpha, recipientTwo, 10 ether);
        _grantTokenRewards(rewardManager, distributions, 1);

        uint256 recipientOneTokenBefore = rewardTokenOne.balanceOf(recipientOne);
        uint256 recipientTwoTokenBefore = rewardTokenOne.balanceOf(recipientTwo);

        address[] memory recipientsToClaim = new address[](2);
        recipientsToClaim[0] = recipientOne;
        recipientsToClaim[1] = recipientTwo;

        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 1);

        assertEq(rewardTokenOne.balanceOf(recipientOne), recipientOneTokenBefore + 20 ether);
        assertEq(rewardTokenOne.balanceOf(recipientTwo), recipientTwoTokenBefore + 10 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 1), 0);
    }

    function test_claimRewards_zeroAmountNoop_noRevert() public {
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientThree;

        vm.prank(operatorBeta);
        rewardDistributor.claimRewards(recipientsToClaim, 0); // should not revert
    }

    // ───────────────────────── Delegated claim

    function test_claimOnBehalf_requiresPerRecipientAuthorization() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        _grantETHRewards(rewardManager, distributions);

        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        vm.expectRevert(); // InvalidClaimDelegate
        vm.prank(claimDelegateAlpha);
        rewardDistributor.claimOnbehalfOfOperator(operatorAlpha, recipientsToClaim, 0);

        vm.prank(operatorAlpha);
        rewardDistributor.setClaimDelegate(claimDelegateAlpha, recipientOne, true);

        uint256 recipientOneBalanceBefore = recipientOne.balance;
        vm.prank(claimDelegateAlpha);
        rewardDistributor.claimOnbehalfOfOperator(operatorAlpha, recipientsToClaim, 0);
        assertEq(recipientOne.balance, recipientOneBalanceBefore + 1 ether);

        vm.prank(operatorAlpha);
        rewardDistributor.setClaimDelegate(claimDelegateAlpha, recipientOne, false);

        vm.expectRevert(); // InvalidClaimDelegate
        vm.prank(claimDelegateAlpha);
        rewardDistributor.claimOnbehalfOfOperator(operatorAlpha, recipientsToClaim, 0);
    }

    // ───────────────────────── Recipient resolution & overrides

    function test_getKeyRecipient_precedence_perKey_over_global_over_operator() public {
        // Default fallback: operator itself
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), operatorAlpha);

        // Global override
        vm.prank(operatorAlpha);
        rewardDistributor.setOperatorGlobalOverride(recipientTwo);
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), recipientTwo);

        // Per-key override beats global
        bytes[] memory pubkeysToOverride = new bytes[](1);
        pubkeysToOverride[0] = pubkeyOne;

        vm.prank(operatorAlpha);
        rewardDistributor.overrideRecipientByPubkey(pubkeysToOverride, recipientOne);
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), recipientOne);

        // Another key still resolves to global
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyTwo), recipientTwo);
    }

    function test_getKeyRecipient_revertsOnInvalidPubkeyLength() public {
        bytes memory invalidLengthPubkey = hex"01"; // not 48 bytes
        vm.expectRevert(); // InvalidBLSPubKeyLength
        rewardDistributor.getKeyRecipient(operatorAlpha, invalidLengthPubkey);
    }

    // ───────────────────────── Admin & pause

    function test_onlyOwner_canSetRewardToken_andRejectsTokenIdZero() public {
        vm.expectRevert(); // onlyOwner
        rewardDistributor.setRewardToken(address(rewardTokenOne), 9);
        vm.startPrank(contractOwner);
        vm.expectRevert(); // InvalidTokenID
        rewardDistributor.setRewardToken(address(rewardTokenOne), 0);
        vm.stopPrank();
    }

    function test_onlyOwner_canSetRewardManager() public {
        address newRewardManager = address(0xDEAD);

        vm.expectRevert(); // onlyOwner
        rewardDistributor.setRewardManager(newRewardManager);

        vm.startPrank(contractOwner);
        vm.expectRevert(); // ZeroAddress
        rewardDistributor.setRewardManager(address(0));
        rewardDistributor.setRewardManager(newRewardManager);
        vm.stopPrank();
    }

    function test_pause_blocksMutatingEndpoints_unpauseRestores() public {
        vm.prank(contractOwner);
        rewardDistributor.pause();

        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);

        vm.expectRevert(); vm.prank(rewardManager);
        rewardDistributor.grantETHRewards{value: 1 ether}(distributions);

        vm.expectRevert(); vm.prank(rewardManager);
        rewardDistributor.grantTokenRewards(distributions, 1);

        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        vm.expectRevert(); vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);

        vm.prank(contractOwner);
        rewardDistributor.unpause();

        vm.prank(rewardManager);
        rewardDistributor.grantETHRewards{value: 1 ether}(distributions); // ok after unpause
    }

    // ───────────────────────── Multi-token sanity

    function test_secondToken_tokenId2_isIndependentOfTokenId1() public {
        // Grant tokenId 1
        RewardDistributor.Distribution[] memory distributionsToken1 = new RewardDistributor.Distribution[](1);
        distributionsToken1[0] = _distribution(operatorAlpha, recipientOne, 5 ether);
        _grantTokenRewards(rewardManager, distributionsToken1, 1);

        // Grant tokenId 2
        RewardDistributor.Distribution[] memory distributionsToken2 = new RewardDistributor.Distribution[](2);
        distributionsToken2[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        distributionsToken2[1] = _distribution(operatorAlpha, recipientTwo, 2 ether);
        _grantTokenRewards(rewardManager, distributionsToken2, 2);

        // Pending are independent per token
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 5 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 2), 1 ether);
        assertEq(_getPending(operatorAlpha, recipientTwo, 2), 2 ether);

        // Claim only tokenId 2
        address[] memory recipientsToClaim = new address[](2);
        recipientsToClaim[0] = recipientOne;
        recipientsToClaim[1] = recipientTwo;

        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 2);

        // tokenId 1 still pending
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 5 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 2), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 2), 0);
    }

    // 1) Multiple ETH grants accumulate into the same (operator, recipient) bucket.
    function test_multipleGrants_accumulatePending_ETH() public {
        RewardDistributor.Distribution[] memory firstBatch = new RewardDistributor.Distribution[](2);
        firstBatch[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        firstBatch[1] = _distribution(operatorAlpha, recipientTwo, 2 ether);
        _grantETHRewards(rewardManager, firstBatch);

        RewardDistributor.Distribution[] memory secondBatch = new RewardDistributor.Distribution[](2);
        secondBatch[0] = _distribution(operatorAlpha, recipientOne, 3 ether);
        secondBatch[1] = _distribution(operatorAlpha, recipientTwo, 4 ether);
        _grantETHRewards(rewardManager, secondBatch);

        assertEq(_getPending(operatorAlpha, recipientOne, 0), 4 ether);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 6 ether);
    }

    // 2) Partial claim leaves remainder; next claim pays the remainder.
    function test_partialClaim_leavesRemainder_ETH() public {
    
        // app for deterministic partial:
        // Grant 2 first
        RewardDistributor.Distribution[] memory firstGrant = new RewardDistributor.Distribution[](1);
        firstGrant[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
        _grantETHRewards(rewardManager, firstGrant);

        // Claim now (drains 2)
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;
        uint256 r1Before = recipientOne.balance;
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);
        assertEq(recipientOne.balance, r1Before + 2 ether);

        // Grant additional 3, ensure only the remainder is pending
        RewardDistributor.Distribution[] memory secondGrant = new RewardDistributor.Distribution[](1);
        secondGrant[0] = _distribution(operatorAlpha, recipientOne, 3 ether);
        _grantETHRewards(rewardManager, secondGrant);
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 3 ether);

        // Claim again, should transfer 3
        r1Before = recipientOne.balance;
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);
        assertEq(recipientOne.balance, r1Before + 3 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);
    }

    // 3) Double-claim after full claim is a no-op (no revert, no transfer).
    function test_doubleClaim_afterFullClaim_noop_ETH() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        _grantETHRewards(rewardManager, distributions);

        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        // First claim pays 1 ether
        uint256 r1Before = recipientOne.balance;
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);
        assertEq(recipientOne.balance, r1Before + 1 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);

        // Second claim is a no-op
        r1Before = recipientOne.balance;
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);
        assertEq(recipientOne.balance, r1Before); // unchanged
    }

    // 4) Operator isolation: same recipient cannot be claimed by the wrong operator.
    function test_operatorIsolation_sameRecipient_cannotCrossClaim_ETH() public {
        // Grant to (alpha, r1) and (beta, r1)
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](2);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
        distributions[1] = _distribution(operatorBeta,  recipientOne, 3 ether);
        _grantETHRewards(rewardManager, distributions);

        // OperatorAlpha claims: only its 2 ether should transfer
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;
        uint256 before = recipientOne.balance;
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);
        assertEq(recipientOne.balance, before + 2 ether);

        // OperatorBeta still has 3 pending
        assertEq(_getPending(operatorBeta, recipientOne, 0), 3 ether);
    }

    // 5) Delegate auth is bound to the operator as well; cannot claim for a different operator.
    function test_delegateCannotClaimForDifferentOperator_evenIfAuthorizedForRecipient() public {
        // Authorize delegate for (operatorAlpha, recipientOne)
        vm.prank(operatorAlpha);
        rewardDistributor.setClaimDelegate(claimDelegateAlpha, recipientOne, true);

        // Grant to (operatorBeta, recipientOne)
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorBeta, recipientOne, 1 ether);
        _grantETHRewards(rewardManager, distributions);

        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        // Delegate cannot claim for operatorBeta
        vm.expectRevert();
        vm.prank(claimDelegateAlpha);
        rewardDistributor.claimOnbehalfOfOperator(operatorBeta, recipientsToClaim, 0);
    }

    // 6) Token grants also require owner/manager permissions (mirror of ETH).
    function test_grantTokenRewards_onlyOwnerOrManager() public {
        // Prepare a simple distribution
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 10 ether);

        // 1) Unauthorized caller (this contract) → revert
        vm.expectRevert(); // NotOwnerOrRewardManager (or equivalent)
        rewardDistributor.grantTokenRewards(distributions, 1);

        // 2) Owner path → must have balance + allowance
        vm.startPrank(contractOwner);
        // Give owner enough token #1 so transferFrom(owner → distributor) can succeed
        rewardTokenOne.mint(contractOwner, 10 ether);
        rewardTokenOne.approve(address(rewardDistributor), 10 ether);
        rewardDistributor.grantTokenRewards(distributions, 1);
        vm.stopPrank();

        assertEq(_getPending(operatorAlpha, recipientOne, 1), 10 ether);

        // 3) Reward manager path → already funded in setUp(); just approve and call
        RewardDistributor.Distribution[] memory distributions2 = new RewardDistributor.Distribution[](1);
        distributions2[0] = _distribution(operatorAlpha, recipientTwo, 1 ether);

        vm.startPrank(rewardManager);
        IERC20(rewardDistributor.rewardTokens(1)).approve(address(rewardDistributor), 1 ether);
        rewardDistributor.grantTokenRewards(distributions2, 1);
        vm.stopPrank();

        assertEq(_getPending(operatorAlpha, recipientTwo, 1), 1 ether);
    }

    // 7) Token grants without sufficient allowance should revert.
    function test_grantTokenRewards_withoutAllowance_reverts() public {
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 5 ether);

        // No approve done → transferFrom should fail inside grant
        vm.expectRevert();
        vm.prank(rewardManager);
        rewardDistributor.grantTokenRewards(distributions, 1);
    }

    // 8) Per-key override with multiple keys updates each independently.
    function test_overrideRecipientByPubkey_multipleKeys_updatesEach() public {
        // Ensure globals are not interfering
        vm.prank(operatorAlpha);
        rewardDistributor.setOperatorGlobalOverride(recipientOne);

        // Build two distinct 48-byte keys
        bytes[] memory keys = new bytes[](2);
        keys[0] = hex"101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f";
        keys[1] = hex"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"; // 48 bytes of 0xaa

        // Override both to recipientTwo
        vm.prank(operatorAlpha);
        rewardDistributor.overrideRecipientByPubkey(keys, recipientTwo);

        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, keys[0]), recipientTwo);
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, keys[1]), recipientTwo);
    }

    // 9) Global override applies to keys without per-key overrides, and changing it updates resolution.
    function test_setOperatorGlobalOverride_updatesAllKeysWithoutPerKey() public {
        // No per-key override set → default to operator
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), operatorAlpha);

        // Set global override → resolution updates
        vm.prank(operatorAlpha);
        rewardDistributor.setOperatorGlobalOverride(recipientTwo);
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), recipientTwo);

        // Switch global again → resolution follows
        vm.prank(operatorAlpha);
        rewardDistributor.setOperatorGlobalOverride(recipientThree);
        assertEq(rewardDistributor.getKeyRecipient(operatorAlpha, pubkeyOne), recipientThree);
    }

    // 10) Cross-token independence with mixed grants: claim ETH only; ERC20 remains.
    function test_crossTokenIndependence_mixedGrants_claimOnlyOneToken() public {
        // Grant: ETH 2 to (alpha,r1), ERC20(1) 3 to (alpha,r1)
        {
            RewardDistributor.Distribution[] memory ethBatch = new RewardDistributor.Distribution[](1);
            ethBatch[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
            _grantETHRewards(rewardManager, ethBatch);

            RewardDistributor.Distribution[] memory tknBatch = new RewardDistributor.Distribution[](1);
            tknBatch[0] = _distribution(operatorAlpha, recipientOne, 3 ether);
            _grantTokenRewards(rewardManager, tknBatch, 1);
        }

        // Claim ETH only
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;
        uint256 r1Before = recipientOne.balance;

        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 0);

        // ETH cleared, token still pending
        assertEq(recipientOne.balance, r1Before + 2 ether);
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 3 ether);
    }

    // MIGRATION: move accrued from one recipient bucket to another for the CALLER (operator)

    function test_migrateExistingRewards_happyPath_ETH() public {
        // Accrue (operatorAlpha → recipientOne) 4 ETH
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 4 ether);
        _grantETHRewards(rewardManager, distributions);

        assertEq(_getPending(operatorAlpha, recipientOne, 0), 4 ether);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 0);

        // OperatorAlpha migrates its own accrued from recipientOne → recipientTwo
        vm.prank(operatorAlpha);
        rewardDistributor.migrateExistingRewards(recipientOne, recipientTwo, 0);

        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 4 ether);
    }

    function test_migrateExistingRewards_happyPath_Token() public {
        // Accrue (operatorAlpha → recipientOne) 7 tokens (id=1)
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 7 ether);
        _grantTokenRewards(rewardManager, distributions, 1);

        // Migrate by the operator (caller)
        vm.prank(operatorAlpha);
        rewardDistributor.migrateExistingRewards(recipientOne, recipientTwo, 1);

        assertEq(_getPending(operatorAlpha, recipientOne, 1), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 1), 7 ether);
    }

    function test_migrateExistingRewards_revert_ifNoClaimableRewardsForCaller() public {
        // Accrue to operatorAlpha, but attempt migration from operatorBeta (caller)
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        _grantETHRewards(rewardManager, distributions);

        vm.expectRevert(); // NoClaimableRewards(msg.sender, from)
        vm.prank(operatorBeta);
        rewardDistributor.migrateExistingRewards(recipientOne, recipientTwo, 0);
    }

    function test_migrateExistingRewards_revert_zeroRecipient_orSameRecipient() public {
        // Accrue a small amount so the "has rewards" guard passes
        RewardDistributor.Distribution[] memory distributions = new RewardDistributor.Distribution[](1);
        distributions[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        _grantETHRewards(rewardManager, distributions);

        // to == address(0)
        vm.expectRevert(); // ZeroAddress()
        vm.prank(operatorAlpha);
        rewardDistributor.migrateExistingRewards(recipientOne, address(0), 0);

        // to == from
        vm.expectRevert(); // InvalidRecipient()
        vm.prank(operatorAlpha);
        rewardDistributor.migrateExistingRewards(recipientOne, recipientOne, 0);
    }

    // ───────────────────────────────────────────────────────────────────────────
    // OWNER RECLAIM: pull multiple buckets back to the owner

    function test_reclaimStipendsToOwner_happyPath_ETH_and_Token() public {
        // Accrue mixed (ETH + token) for two buckets of the same operator
        {
            RewardDistributor.Distribution[] memory ethBatch = new RewardDistributor.Distribution[](2);
            ethBatch[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
            ethBatch[1] = _distribution(operatorAlpha, recipientTwo, 3 ether);
            _grantETHRewards(rewardManager, ethBatch);

            RewardDistributor.Distribution[] memory tokenBatch = new RewardDistributor.Distribution[](2);
            tokenBatch[0] = _distribution(operatorAlpha, recipientOne, 5 ether);
            tokenBatch[1] = _distribution(operatorAlpha, recipientTwo, 7 ether);
            _grantTokenRewards(rewardManager, tokenBatch, 1);
        }

        address[] memory operatorList = new address[](2);
        address[] memory recipientList = new address[](2);
        operatorList[0] = operatorAlpha;
        operatorList[1] = operatorAlpha;
        recipientList[0] = recipientOne;
        recipientList[1] = recipientTwo;

        // Reclaim ETH buckets to owner
        uint256 ownerEthBefore = contractOwner.balance;
        vm.prank(contractOwner);
        rewardDistributor.reclaimStipendsToOwner(operatorList, recipientList, 0);

        assertEq(contractOwner.balance, ownerEthBefore + 5 ether); // 2 + 3
        assertEq(_getPending(operatorAlpha, recipientOne, 0), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 0), 0);

        // Reclaim token buckets to owner
        uint256 ownerTokenBefore = rewardTokenOne.balanceOf(contractOwner);
        vm.prank(contractOwner);
        rewardDistributor.reclaimStipendsToOwner(operatorList, recipientList, 1);

        assertEq(rewardTokenOne.balanceOf(contractOwner), ownerTokenBefore + 12 ether); // 5 + 7
        assertEq(_getPending(operatorAlpha, recipientOne, 1), 0);
        assertEq(_getPending(operatorAlpha, recipientTwo, 1), 0);
    }

    function test_reclaimStipendsToOwner_revert_lengthMismatch() public {
        address[] memory operatorList = new address[](2);
        address[] memory recipientList = new address[](1);
        operatorList[0] = operatorAlpha;
        operatorList[1] = operatorBeta;
        recipientList[0] = recipientOne;

        vm.expectRevert(); // LengthMismatch()
        vm.prank(contractOwner);
        rewardDistributor.reclaimStipendsToOwner(operatorList, recipientList, 0);
    }

    function test_reclaimStipendsToOwner_revert_noClaimableAcrossSet() public {
        // Ensure no pending in the targeted (operator, recipient) pairs
        address[] memory operatorList = new address[](1);
        address[] memory recipientList = new address[](1);
        operatorList[0] = operatorAlpha;
        recipientList[0] = recipientOne;

        vm.expectRevert(); // NoClaimableRewards(owner, owner)
        vm.prank(contractOwner);
        rewardDistributor.reclaimStipendsToOwner(operatorList, recipientList, 0);
    }

    // ───────────────────────────────────────────────────────────────────────────
    // CLAIM: invalid inputs

    function test_claimRewards_revert_invalidOperatorZero_viaOnBehalf() public {
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        // claimOnbehalfOfOperator passes operator arg → zero should revert InvalidOperator()
        vm.expectRevert();
        vm.prank(claimDelegateAlpha);
        rewardDistributor.claimOnbehalfOfOperator(address(0), recipientsToClaim, 0);
    }

    function test_claimRewards_revert_invalidTokenId() public {
        // InvalidRewardToken() in _claimRewards should guard this
        address[] memory recipientsToClaim = new address[](1);
        recipientsToClaim[0] = recipientOne;

        vm.expectRevert();
        vm.prank(operatorAlpha);
        rewardDistributor.claimRewards(recipientsToClaim, 9_999);
    }

    // ───────────────────────────────────────────────────────────────────────────
    // EVENTS + CONSOLIDATION

    function test_events_emitted_onGrants_and_setters() public {
        // ETH grant emits ETHGranted for each item and a value check
        RewardDistributor.Distribution[] memory ethBatch = new RewardDistributor.Distribution[](1);
        ethBatch[0] = _distribution(operatorAlpha, recipientOne, 1 ether);

        vm.expectEmit(true, true, true, true);
        emit IRewardDistributor.ETHGranted(operatorAlpha, recipientOne, 1 ether);
        vm.prank(rewardManager);
        rewardDistributor.grantETHRewards{value: 1 ether}(ethBatch);

        // Token grant emits TokensGranted and RewardsBatchGranted with total
        RewardDistributor.Distribution[] memory tokenBatch = new RewardDistributor.Distribution[](2);
        tokenBatch[0] = _distribution(operatorAlpha, recipientOne, 2 ether);
        tokenBatch[1] = _distribution(operatorAlpha, recipientTwo, 3 ether);

        vm.startPrank(rewardManager);
        IERC20(rewardDistributor.rewardTokens(1)).approve(address(rewardDistributor), 5 ether);
        vm.expectEmit(true, true, true, true);
        emit IRewardDistributor.TokensGranted(operatorAlpha, recipientOne, 2 ether);
        vm.expectEmit(true, true, true, true);
        emit IRewardDistributor.TokensGranted(operatorAlpha, recipientTwo, 3 ether);
        vm.expectEmit(false, false, false, true);
        emit IRewardDistributor.RewardsBatchGranted(0, 5 ether);
        rewardDistributor.grantTokenRewards(tokenBatch, 1);
        vm.stopPrank();

        // Setters
        vm.expectEmit(true, true, true, true);
        emit IRewardDistributor.RewardManagerSet(address(0xBEEF));
        vm.prank(contractOwner);
        rewardDistributor.setRewardManager(address(0xBEEF));

        vm.expectEmit(true, true, true, true);
        emit IRewardDistributor.RewardTokenSet(address(rewardTokenTwo), 42);
        vm.prank(contractOwner);
        rewardDistributor.setRewardToken(address(rewardTokenTwo), 42);
    }

    function test_grant_duplicateEntries_consolidates_ETH_and_Token() public {
        // ETH: three entries for same (operatorAlpha, recipientOne)
        RewardDistributor.Distribution[] memory ethBatch = new RewardDistributor.Distribution[](3);
        ethBatch[0] = _distribution(operatorAlpha, recipientOne, 1 ether);
        ethBatch[1] = _distribution(operatorAlpha, recipientOne, 2 ether);
        ethBatch[2] = _distribution(operatorAlpha, recipientOne, 4 ether);
        _grantETHRewards(rewardManager, ethBatch);

        assertEq(_getPending(operatorAlpha, recipientOne, 0), 7 ether);

        // Token: two entries for same (operatorAlpha, recipientOne) on token 1
        RewardDistributor.Distribution[] memory tknBatch = new RewardDistributor.Distribution[](2);
        tknBatch[0] = _distribution(operatorAlpha, recipientOne, 3 ether);
        tknBatch[1] = _distribution(operatorAlpha, recipientOne, 5 ether);
        _grantTokenRewards(rewardManager, tknBatch, 1);

        assertEq(_getPending(operatorAlpha, recipientOne, 1), 8 ether);
    }
}
