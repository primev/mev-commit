// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import {ProviderRegistry} from "../contracts/ProviderRegistry.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {PreConfCommitmentStore} from "../contracts/PreConfCommitmentStore.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

contract ProviderRegistryTest is Test {
    uint256 testNumber;
    ProviderRegistry internal providerRegistry;
    uint16 internal feePercent;
    uint256 internal minStake;
    address internal provider;
    address internal feeRecipient;
    BidderRegistry bidderRegistry;
    PreConfCommitmentStore preConfCommitmentStore;
    BlockTracker blockTracker;
    uint256 blocksPerWindow;
    event ProviderRegistered(address indexed provider, uint256 stakedAmount);

    function setUp() public {
        testNumber = 42;
        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, 
            (minStake, 
            feeRecipient, 
            feePercent, 
            address(this))) 
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, 
            (address(this), blocksPerWindow)) 
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, 
            (minStake, 
            feeRecipient, 
            feePercent, 
            address(this), 
            address(blockTracker),
            blocksPerWindow)) 
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

    function test_VerifyInitialContractState() public view {
        assertEq(providerRegistry.minStake(), 1e18 wei);
        assertEq(providerRegistry.feeRecipient(), feeRecipient);
        assertEq(providerRegistry.feePercent(), feePercent);
        assertEq(providerRegistry.preConfirmationsContract(), address(0));
        assertEq(providerRegistry.providerRegistered(provider), false);
    }

    function testFail_ProviderStakeAndRegisterMinStake() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}();
    }

    function test_ProviderStakeAndRegister() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectEmit(true, false, false, true);

        emit ProviderRegistered(provider, 1e18 wei);

        providerRegistry.registerAndStake{value: 1e18 wei}();

        bool isProviderRegistered = providerRegistry.providerRegistered(
            provider
        );
        assertEq(isProviderRegistered, true);

        uint256 providerStakeStored = providerRegistry.checkStake(provider);
        assertEq(providerStakeStored, 1e18 wei);
    }

    function testFail_ProviderStakeAndRegisterAlreadyRegistered() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
        vm.expectRevert(bytes(""));
        providerRegistry.registerAndStake{value: 1 wei}();
    }

    function testFail_receive() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function testFail_fallback() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        vm.expectRevert(bytes(""));
        (bool success, ) = address(providerRegistry).call{value: 1 wei}("");
        require(success, "Couldn't transfer to provider");
    }

    function test_SetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.prank(address(this));
        providerRegistry.setNewFeeRecipient(newRecipient);

        assertEq(providerRegistry.feeRecipient(), newRecipient);
    }

    function testFail_SetNewFeeRecipient() public {
        address newRecipient = vm.addr(2);
        vm.expectRevert(bytes(""));
        providerRegistry.setNewFeeRecipient(newRecipient);
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

    function test_shouldSlashProvider() public {
        providerRegistry.setPreconfirmationsContract(address(this));
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();
        address bidder = vm.addr(4);

        providerRegistry.slash(1 ether, provider, payable(bidder),100);

        assertEq(providerRegistry.bidderAmount(bidder), 900000000000000000 wei);
        assertEq(providerRegistry.feeRecipientAmount(), 100000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 1 ether);
    }

    function test_shouldSlashProviderWithoutFeeRecipient() public {
        vm.prank(address(this));
        providerRegistry.setNewFeeRecipient(address(0));
        providerRegistry.setPreconfirmationsContract(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();
        address bidder = vm.addr(4);

        providerRegistry.slash(1 ether, provider, payable(bidder),100);

        assertEq(providerRegistry.bidderAmount(bidder), 900000000000000000 wei);
        assertEq(providerRegistry.providerStakes(provider), 1 ether);
    }

    function testFail_shouldRetrieveFundsNotPreConf() public {
        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();
        address bidder = vm.addr(4);
        vm.expectRevert(bytes(""));
        providerRegistry.slash(1 ether, provider, payable(bidder),100);
    }

    function testFail_shouldRetrieveFundsGreaterThanStake() public {
        vm.prank(address(this));
        providerRegistry.setPreconfirmationsContract(address(this));

        vm.deal(provider, 3 ether);
        vm.prank(provider);
        providerRegistry.registerAndStake{value: 2 ether}();
        address bidder = vm.addr(4);
        vm.expectRevert(bytes(""));
        vm.prank(address(this));

        providerRegistry.slash(3 ether, provider, payable(bidder),100);
    }

    function test_FeeRecipientAmount() public {
        providerRegistry.setNewFeeRecipient(vm.addr(6));
        vm.deal(provider, 3 ether);
        vm.prank(provider);

        providerRegistry.registerAndStake{value: 2 ether}();
        providerRegistry.setPreconfirmationsContract(address(this));
        providerRegistry.slash(1e18 wei, provider, payable(provider),50);
        assertEq(
            providerRegistry.feeRecipientAmount(),
            5e16 wei,
            "FeeRecipientAmount should match"
        );
        providerRegistry.withdrawFeeRecipientAmount();
        assertEq(
            providerRegistry.feeRecipientAmount(),
            0,
            "FeeRecipientAmount should be zero after withdrawal"
        );
    }

    function test_WithdrawBidderAmount() public {
        address bidder = vm.addr(7);
        vm.deal(bidder, 3 ether);
        vm.prank(bidder);
        providerRegistry.registerAndStake{value: 2 ether}();

        providerRegistry.setPreconfirmationsContract(address(this));
        providerRegistry.slash(1e18 wei, bidder, payable(bidder),100);
        vm.prank(bidder);
        providerRegistry.withdrawBidderAmount(bidder);
        assertEq(
            providerRegistry.bidderAmount(bidder),
            0,
            "BidderAmount should be zero after withdrawal"
        );
    }

    function test_WithdrawStakedAmountWithoutFeeRecipient() public {
        providerRegistry.setNewFeeRecipient(address(0));
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        vm.prank(address(preConfCommitmentStore));
        providerRegistry.slash(1e18 wei, newProvider, payable(newProvider),100);
        vm.prank(newProvider);
        providerRegistry.withdrawStakedAmount(payable(newProvider));
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

    function testFail_WithdrawStakedAmountUnauthorized() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
        vm.expectRevert(bytes(""));
        providerRegistry.withdrawStakedAmount(payable(vm.addr(12)));
    }

    function test_RegisterAndStake() public {
        address newProvider = vm.addr(5);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
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

    function test_WithdrawStakedAmount() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
        providerRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
        vm.prank(address(preConfCommitmentStore));
        providerRegistry.slash(1e18 wei, newProvider, payable(newProvider),100);
        vm.prank(newProvider);
        providerRegistry.withdrawStakedAmount(payable(newProvider));
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

    function testFail_WithdrawStakedAmountWithoutCommitments() public {
        address newProvider = vm.addr(8);
        vm.deal(newProvider, 3 ether);
        vm.prank(newProvider);
        providerRegistry.registerAndStake{value: 2e18 wei}();
        vm.expectRevert("Provider Commitments still pending");
        providerRegistry.withdrawStakedAmount(payable(newProvider));
    }
}
