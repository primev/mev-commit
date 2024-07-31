// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import {ProviderRegistry} from "../contracts/ProviderRegistry.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {PreConfCommitmentStore} from "../contracts/PreConfCommitmentStore.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract ProviderRegistryTest is Test {
    uint256 public testNumber;
    ProviderRegistry public providerRegistry;
    uint16 public feePercent;
    uint256 public minStake;
    address public provider;
    address public feeRecipient;
    BidderRegistry public bidderRegistry;
    PreConfCommitmentStore public preConfCommitmentStore;
    BlockTracker public blockTracker;
    uint256 public blocksPerWindow;
    uint256 public withdrawalDelay;
    bytes public validBLSPubkey = hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    uint256 public protocolFeePayoutPeriodBlocks;
    event ProviderRegistered(address indexed provider, uint256 stakedAmount, bytes blsPublicKey);
    event WithdrawalRequested(address indexed provider, uint256 timestamp);
    event WithdrawalCompleted(address indexed provider, uint256 amount);
    event FeeTransfer(uint256 amount, address recipient);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;
        withdrawalDelay = 24 * 3600; // 24 hours
        protocolFeePayoutPeriodBlocks = 100;
        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, 
            (minStake, 
            feeRecipient, 
            feePercent, 
            address(this),
            withdrawalDelay,
            protocolFeePayoutPeriodBlocks
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
            protocolFeePayoutPeriodBlocks)) 
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));
        
        address preconfStoreProxy = Upgrades.deployUUPSProxy(
            "PreConfCommitmentStore.sol",
            abi.encodeCall(PreConfCommitmentStore.initialize, 
            (address(providerRegistry), // Provider Registry
            address(bidderRegistry), // User Registry
            address(blockTracker), // Block Tracker
            feeRecipient, // Oracle
            address(this),
            500,
            blocksPerWindow))
        );
        preConfCommitmentStore = PreConfCommitmentStore(payable(preconfStoreProxy));

        provider = vm.addr(1);
        vm.deal(provider, 100 ether);
        vm.deal(address(this), 100 ether);
    }

    function testVerifyInitialContractState() public view {
        assertEq(providerRegistry.minStake(), 1e18 wei);
        assertEq(feePercent, feePercent);
        assertEq(withdrawalDelay, withdrawalDelay);
        assertEq(providerRegistry.feePercent(), feePercent);
        assertEq(providerRegistry.preConfirmationsContract(), address(0));
        assertEq(providerRegistry.providerRegistered(provider), false);
        (address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks) = bidderRegistry.protocolFeeTracker();
        assertEq(recipient, feeRecipient);
        assertEq(payoutPeriodBlocks, protocolFeePayoutPeriodBlocks);
        assertEq(lastPayoutBlock, block.number);
        assertEq(accumulatedAmount, 0);
    }

    function testFailProviderStakeAndRegisterMinStake() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}(validBLSPubkey);
    }

    function testFailProviderStakeAndRegisterInvalidBLSKey() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert("Invalid BLS public key length");
        bytes memory blsPublicKey = abi.encodePacked(uint256(134));
        providerRegistry.registerAndStake{value: 1 wei}(blsPublicKey);
    }

    function testProviderStakeAndRegister() public {
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

    function testFailProviderStakeAndRegisterAlreadyRegistered() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}(validBLSPubkey);
    }

    function testFailReceive() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function testFailFallback() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function testSetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        providerRegistry.setNewProtocolFeeRecipient(newRecipient);
        (address recipient, , ,) = providerRegistry.protocolFeeTracker();
        assertEq(recipient, newRecipient);
    }

    function testFailSetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert(bytes(""));
        providerRegistry.setNewProtocolFeeRecipient(newRecipient);
    }

    function testSetNewFeePercent() public {
        vm.prank(address(this));
        providerRegistry.setNewFeePercent(uint16(25));

        assertEq(providerRegistry.feePercent(), uint16(25));
    }

    function testFailSetNewFeePercent() public {
        vm.expectRevert(bytes(""));
        providerRegistry.setNewFeePercent(uint16(25));
    }

    function testSetPreConfContract() public {
        vm.prank(address(this));
        address newPreConfContract = vm.addr(3);
        providerRegistry.setPreconfirmationsContract(newPreConfContract);

        assertEq(
            providerRegistry.preConfirmationsContract(),
            newPreConfContract
        );
    }

    function testFailSetPreConfContract() public {
        vm.prank(address(this));
        vm.expectRevert(bytes(""));
        providerRegistry.setPreconfirmationsContract(address(0));
    }

    function testShouldSlashProvider() public {
        providerRegistry.setPreconfirmationsContract(address(this));
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);

        vm.expectCall(bidder, 900000000000000000 wei, new bytes(0));
        providerRegistry.slash(1 ether, provider, payable(bidder), 100);

        assertEq(bidder.balance, 900000000000000000 wei);
        assertEq(providerRegistry.getAccumulatedProtocolFee(), 100000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 1 ether);
    }

    function testShouldSlashProviderWithoutFeeRecipient() public {
        vm.prank(address(this));
        providerRegistry.setNewProtocolFeeRecipient(address(0));
        providerRegistry.setPreconfirmationsContract(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);

        vm.expectCall(bidder, 900000000000000000 wei, new bytes(0));
        providerRegistry.slash(1 ether, provider, payable(bidder), 100);

        assertEq(bidder.balance, 900000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 1 ether);
    }

    function testFailShouldRetrieveFundsNotPreConf() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        address bidder = vm.addr(4);
        vm.expectRevert(bytes(""));
        providerRegistry.slash(1 ether, provider, payable(bidder),100);
    }

    function testFailShouldRetrieveFundsGreaterThanStake() public {
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

    function testProtocolFeeBehavior() public {
        providerRegistry.setNewProtocolFeeRecipient(vm.addr(6));
        vm.deal(provider, 3 ether);
        vm.prank(provider);

        address bidder = vm.addr(4);

        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(address(this));
        providerRegistry.slash(1e18 wei, provider, payable(bidder), 50);
        assertEq(
            providerRegistry.getAccumulatedProtocolFee(),
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
            providerRegistry.getAccumulatedProtocolFee(),
            0,
            "Accumulated protocol fee should be zero"
        );
        assertEq(
            vm.addr(6).balance,
            1e17 wei,
            "FeeRecipient should have received 1e17 wei"
        );
    }

    function testWithdrawStakedAmountWithoutFeeRecipient() public {
        providerRegistry.setNewProtocolFeeRecipient(address(0));
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        vm.prank(address(preConfCommitmentStore));
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
            2e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function testFailWithdrawStakedAmountUnauthorized() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.expectRevert(bytes(""));
        address wrongNewProvider = vm.addr(12);
        vm.prank(wrongNewProvider);
        providerRegistry.withdraw();
    }

    function testRegisterAndStake() public {
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

    function testFailWithdrawStakedAmountWithoutCommitments() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.unstake();
        vm.warp(block.timestamp + 24 hours); // Move forward in time
        vm.expectRevert("Provider Commitments still pending");
        providerRegistry.withdraw();
    }

    function testRequestWithdrawal() public {
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

    function testWithdrawStakedAmount() public {
        address newProvider = vm.addr(8);
        address bidder = vm.addr(9);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        vm.prank(address(preConfCommitmentStore));
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
            2e18 wei,
            "Provider's balance should increase by staked amount"
        );
    }

    function testWithdrawStakedAmountBefore24Hours() public {
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

    function testWithdrawStakedAmountWithoutRequest() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}(validBLSPubkey);
        vm.prank(newProvider);
        vm.expectRevert("No unstake request");
        providerRegistry.withdraw();
    }
}
