// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract ProviderRegistryTest is Test {
    uint256 public testNumber;
    ProviderRegistry public providerRegistry;
    uint16 public feePercent;
    uint256 public minStake;
    address public provider;
    address public feeRecipient;
    BidderRegistry public bidderRegistry;
    PreconfManager public preconfManager;
    BlockTracker public blockTracker;
    uint256 public blocksPerWindow;
    uint256 public withdrawalDelay;
    bytes public validBLSPubkey = hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    uint256 public penaltyFeePayoutPeriodBlocks;
    event ProviderRegistered(address indexed provider, uint256 stakedAmount, bytes blsPublicKey);
    event WithdrawalRequested(address indexed provider, uint256 timestamp);
    event WithdrawalCompleted(address indexed provider, uint256 amount);
    event FeeTransfer(uint256 amount, address indexed recipient);
    event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient);
    event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;
        withdrawalDelay = 24 * 3600; // 24 hours
        penaltyFeePayoutPeriodBlocks = 100;
        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, 
            (minStake, 
            feeRecipient, 
            feePercent, 
            address(this),
            withdrawalDelay,
            penaltyFeePayoutPeriodBlocks
            )) 
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, 
            (blocksPerWindow, address(this), address(this))) 
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, 
            (feeRecipient, 
            feePercent, 
            address(this), 
            address(blockTracker),
            blocksPerWindow,
            penaltyFeePayoutPeriodBlocks)) 
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));
        
        address preconfStoreProxy = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(PreconfManager.initialize, 
            (address(providerRegistry), // Provider Registry
            address(bidderRegistry), // User Registry
            address(blockTracker), // Block Tracker
            feeRecipient, // Oracle
            address(this),
            500,
            blocksPerWindow))
        );
        preconfManager = PreconfManager(payable(preconfStoreProxy));

        provider = vm.addr(1);
        vm.deal(provider, 100 ether);
        vm.deal(address(this), 100 ether);
    }

    function test_VerifyInitialContractState() public view {
        assertEq(providerRegistry.minStake(), 1e18 wei);
        assertEq(feePercent, feePercent);
        assertEq(withdrawalDelay, withdrawalDelay);
        assertEq(providerRegistry.feePercent(), feePercent);
        assertEq(providerRegistry.preConfirmationsContract(), address(0));
        assertEq(providerRegistry.providerRegistered(provider), false);
        (address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks) = bidderRegistry.protocolFeeTracker();
        assertEq(recipient, feeRecipient);
        assertEq(payoutPeriodBlocks, penaltyFeePayoutPeriodBlocks);
        assertEq(lastPayoutBlock, block.number);
        assertEq(accumulatedAmount, 0);
    }

    function testFail_ProviderStakeAndRegisterMinStake() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}(validBLSPubkey);
    }

    function testFail_ProviderStakeAndRegisterInvalidBLSKey() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert("Invalid BLS public key length");
        bytes memory blsPublicKey = abi.encodePacked(uint256(134));
        providerRegistry.registerAndStake{value: 1 wei}(blsPublicKey);
    }

    function test_ProviderStakeAndRegister() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectEmit(true, false, false, true);

        emit ProviderRegistered(provider, 1e18 wei, validBLSPubkey);

        providerRegistry.registerAndStake{value: 1e18 wei}(validBLSPubkey);

        bool isProviderRegistered = providerRegistry.providerRegistered(
            provider
        );
        assertEq(isProviderRegistered, true);

        uint256 providerStakeStored = providerRegistry.getProviderStake(provider);
        assertEq(providerStakeStored, 1e18 wei);
    }

    function testFail_ProviderStakeAndRegisterAlreadyRegistered() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}(validBLSPubkey);
    }

    function testFail_Receive() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function testFail_Fallback() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function test_SetNewPenaltyFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit PenaltyFeeRecipientUpdated(newRecipient);
        providerRegistry.setNewPenaltyFeeRecipient(newRecipient);
        (address recipient, , ,) = providerRegistry.penaltyFeeTracker();
        assertEq(recipient, newRecipient);
    }

    function testFail_SetNewPenaltyFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert(bytes(""));
        providerRegistry.setNewPenaltyFeeRecipient(newRecipient);
    }

    function test_SetNewFeePayoutPeriodBlocks() public {
        vm.prank(address(this));
        vm.expectEmit(true, true, true, true);
        emit FeePayoutPeriodBlocksUpdated(890);
        providerRegistry.setFeePayoutPeriodBlocks(890);
        (, , , uint256 payoutPeriodBlocks) = providerRegistry.penaltyFeeTracker();
        assertEq(payoutPeriodBlocks, 890);
    }

    function testFail_SetNewFeePayoutPeriodBlocks() public {
        vm.expectRevert(bytes(""));
        providerRegistry.setFeePayoutPeriodBlocks(83424);
    }

    function test_SetNewFeePercent() public {
        vm.prank(address(this));
        providerRegistry.setNewFeePercent(uint16(25));

        assertEq(providerRegistry.feePercent(), uint16(25));
    }

    function testFail_SetNewFeePercent() public {
        vm.expectRevert(bytes(""));
        providerRegistry.setNewFeePercent(uint16(25));
    }

    function test_SetPreConfContract() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        providerRegistry.setPreconfirmationsContract(newPreConfContract);

        assertEq(
            providerRegistry.preConfirmationsContract(),
            newPreConfContract
        );
    }

    function testFail_SetPreConfContract() public {
        vm.prank(address(this));
        vm.expectRevert(bytes(""));
        providerRegistry.setPreconfirmationsContract(address(0));
    }

    function test_ShouldSlashProvider() public {
        providerRegistry.setPreconfirmationsContract(address(this));
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);

        vm.expectCall(bidder, 1000000000000000000 wei, new bytes(0));
        providerRegistry.slash(1 ether, provider, payable(bidder), 100);

        assertEq(bidder.balance, 1000000000000000000 wei);
        assertEq(providerRegistry.getAccumulatedPenaltyFee(), 100000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 0.9 ether);
    }

    function test_ShouldSlashProviderWithoutFeeRecipient() public {
        vm.prank(address(this));
        providerRegistry.setNewPenaltyFeeRecipient(address(0));
        providerRegistry.setPreconfirmationsContract(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);

        vm.expectCall(bidder, 1000000000000000000 wei, new bytes(0));
        providerRegistry.slash(1 ether, provider, payable(bidder), 100);

        assertEq(bidder.balance, 1000000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 0.9 ether);
    }

    function testFail_ShouldRetrieveFundsNotPreConf() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);
        vm.expectRevert(bytes(""));
        providerRegistry.slash(1 ether, provider, payable(bidder),100);
    }

    function testFail_ShouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        providerRegistry.setPreconfirmationsContract(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));

        providerRegistry.slash(3 ether, provider, payable(bidder), 100);
    }

    function test_PenaltyFeeBehavior() public {
        providerRegistry.setNewPenaltyFeeRecipient(vm.addr(6));
        vm.deal(provider, 3 ether);
        vm.prank(provider);

        address bidder = vm.addr(4);

        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(address(this));
        providerRegistry.slash(1e18 wei, provider, payable(bidder), 50);
        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            5e16 wei,
            "FeeRecipientAmount should match"
        );

        address newProvider = vm.addr(11);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);

        vm.roll(350); // roll past protocol fee payout period

        vm.expectEmit(true, true, true, true);
        emit FeeTransfer(1e17 wei, vm.addr(6));
        providerRegistry.slash(1e18 wei, newProvider, payable(bidder), 50);

        assertEq(
            providerRegistry.getAccumulatedPenaltyFee(),
            0,
            "Accumulated protocol fee should be zero"
        );
        assertEq(
            vm.addr(6).balance,
            1e17 wei,
            "FeeRecipient should have received 1e17 wei"
        );
    }

    function test_WithdrawStakedAmountWithoutFeeRecipient() public {
        providerRegistry.setNewPenaltyFeeRecipient(address(0));
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(
            address(preconfManager)
        );
        vm.prank(address(preconfManager));
        providerRegistry.slash(1e18 wei, newProvider, payable(bidder),100);
        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.prank(newProvider);
        providerRegistry.withdraw();
        assertEq(
            providerRegistry.providerStakes(newProvider),
            0,
            "Provider's staked amount should be zero after withdrawal"
        );
        assertEq(
            newProvider.balance,
            1.9e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function testFail_WithdrawStakedAmountUnauthorized() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.expectRevert(bytes(""));
        address wrongNewProvider = vm.addr(12);
        vm.prank(wrongNewProvider);
        providerRegistry.withdraw();
    }

    function test_RegisterAndStake() public {
        address newProvider = vm.addr(5);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        assertEq(
            providerRegistry.providerStakes(newProvider),
            2e18 wei,
            "Staked amount should match"
        );
        assertEq(
            providerRegistry.providerRegistered(newProvider),
            true,
            "Provider should be registered"
        );
    }

    function test_DelegateRegisterAndStake() public {
        address otherAccount = vm.addr(8);
        vm.deal(otherAccount, 3 ether);

        address newProvider = vm.addr(7);

        vm.prank(otherAccount);
        providerRegistry.delegateRegisterAndStake{value: 2e18 wei}(newProvider, validBLSPubkey);
        assertEq(
            providerRegistry.providerStakes(newProvider),
            2e18 wei,
            "Staked amount should match"
        );
        assertEq(
            providerRegistry.providerRegistered(newProvider),
            true,
            "Provider should be registered"
        );
    }

    function testFail_WithdrawStakedAmountWithoutCommitments() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.expectRevert("Provider Commitments still pending");
        providerRegistry.withdraw();
    }

    function test_RequestWithdrawal() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.prank(newProvider);
        providerRegistry.unstake();
        assertEq(
            providerRegistry.withdrawalRequests(newProvider),
            block.timestamp,
            "Withdrawal request timestamp should match"
        );
    }

    function test_WithdrawStakedAmount() public {
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(
            address(preconfManager)
        );
        vm.prank(address(preconfManager));
        providerRegistry.slash(1e18 wei, newProvider, payable(bidder),100);
        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.prank(newProvider);
        providerRegistry.withdraw();
        assertEq(
            providerRegistry.providerStakes(newProvider),
            0,
            "Provider's staked amount should be zero after withdrawal"
        );
        assertEq(
            newProvider.balance,
            1.9e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function test_WithdrawStakedAmountBefore24Hours() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.prank(newProvider);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 23 hours); // Move forward less than 24 hours
        vm.prank(newProvider);
        vm.expectRevert("Delay has not passed");
        providerRegistry.withdraw();
    }

    function test_WithdrawStakedAmountWithoutRequest() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.prank(newProvider);
        vm.expectRevert("No unstake request");
        providerRegistry.withdraw();
    }
}
