// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Test} from "forge-std/Test.sol";

import {PreConfCommitmentStore} from "../contracts/PreConfCommitmentStore.sol";
import {IPreConfCommitmentStore} from "../contracts/interfaces/IPreConfCommitmentStore.sol";
import {ProviderRegistry} from "../contracts/ProviderRegistry.sol";
import {BidderRegistry} from "../contracts/BidderRegistry.sol";
import {BlockTracker} from "../contracts/BlockTracker.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../contracts/utils/WindowFromBlockNumber.sol";

contract TestPreConfCommitmentStore is Test {
    struct TestCommitment {
        uint256 bid;
        uint64 blockNumber;
        string txnHash;
        string revertingTxHashes;
        uint64 decayStartTimestamp;
        uint64 decayEndTimestamp;
        bytes32 bidDigest;
        bytes32 commitmentDigest;
        bytes bidSignature;
        bytes commitmentSignature;
        uint64 dispatchTimestamp;
        bytes sharedSecretKey;
    }

    TestCommitment internal _testCommitmentAliceBob;
    PreConfCommitmentStore public preConfCommitmentStore;
    uint16 public feePercent;
    uint256 public minStake;
    address public provider;
    address public feeRecipient;
    ProviderRegistry public providerRegistry;
    BlockTracker public blockTracker;
    uint256 public blocksPerWindow;
    BidderRegistry public bidderRegistry;
    bytes public validBLSPubkey =
        hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    uint256 public withdrawalDelay;
    uint256 public protocolFeePayoutPeriodBlocks;
    address public oracleContract;
    function setUp() public {
        _testCommitmentAliceBob = TestCommitment(
            2,
            2,
            "0xkartik",
            "0xkartik",
            10,
            20,
            0x9890bcda118cfabed02ff3b9d05a54dca5310e9ace3b05f259f4731f58ad0900,
            0x47662cacd554166e7a1835d52ccc33ffe9d3be42e07c316e5828cd7c3584d954,
            hex"f9b66c6d57dac947a3aa2b37010df745592cf57f907d437767bc0af6d44b3dc1112168e4cab311d6dfddf7f58c0d07bb95403fca2cc48d4450e088cf9ee894c81b",
            hex"85402f77f929e8e82f559eb2eb5d32f7bba6f39f4f7273e29a888233996c949955678d601fc8c26edbd862f79548e2112f5a7442403189c15515ec43cf3fea221b",
            15,
            bytes("0xsecret")
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;
        withdrawalDelay = 24 * 3600; // 24 hours
        protocolFeePayoutPeriodBlocks = 100;
        oracleContract = address(0x6793);
        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (
                    minStake,
                    feeRecipient,
                    feePercent,
                    address(this),
                    withdrawalDelay,
                    protocolFeePayoutPeriodBlocks
                )
            )
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(
                BlockTracker.initialize,
                (blocksPerWindow, address(this), address(this))
            )
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(
                BidderRegistry.initialize,
                (
                    feeRecipient,
                    feePercent,
                    address(this),
                    address(blockTracker),
                    blocksPerWindow,
                    protocolFeePayoutPeriodBlocks
                )
            )
        );
        bidderRegistry = BidderRegistry(payable(bidderRegistryProxy));

        address preconfStoreProxy = Upgrades.deployUUPSProxy(
            "PreConfCommitmentStore.sol",
            abi.encodeCall(
                PreConfCommitmentStore.initialize,
                (
                    address(providerRegistry), // Provider Registry
                    address(bidderRegistry), // User Registry
                    oracleContract, // Oracle
                    address(this),
                    address(blockTracker), // Block Tracker
                    500,
                    blocksPerWindow
                )
            ) // Commitment Dispatch Window
        );
        preConfCommitmentStore = PreConfCommitmentStore(
            payable(preconfStoreProxy)
        );

        // Sets fake block timestamp
        vm.warp(500);
        bidderRegistry.setPreconfirmationsContract(
            address(preConfCommitmentStore)
        );
    }

    function test_GetBidHash1() public {
        // Step 1: Prepare the test commitment data
        PreConfCommitmentStore.CommitmentParams
            memory testCommitment = IPreConfCommitmentStore.CommitmentParams({
                txnHash: "0xkartik",
                revertingTxHashes: "0xkartik",
                bid: 2,
                blockNumber: 2,
                decayStartTimeStamp: 10,
                decayEndTimeStamp: 20,
                sharedSecretKey: bytes("0xsecret"),
                bidHash: hex"9890bcda118cfabed02ff3b9d05a54dca5310e9ace3b05f259f4731f58ad0900",
                bidSignature: hex"c2ab6e530f6b09337e53e1192857fa10017cdb488cf2a07e0aa4457571492b8c6bff93cbda4e003336656b4ecf8ff46bd1d408b310acdf07be4925a1a8fee4471c",
                commitmentSignature: hex"5b3000290d4f347b94146eb37f66d5368aed18fb8713bf78620abe40ae3de7f635f7ed161801c31ea10e736d88e6fd2a2286bbd59385161dd24c9fefd2568f341b"
            });
        // Step 2: Calculate the bid hash using the getBidHash function
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.bid,
            testCommitment.blockNumber,
            testCommitment.decayStartTimeStamp,
            testCommitment.decayEndTimeStamp
        );

        // Add a alice private key and console log the key
        (, uint256 alicePk) = makeAddrAndKey("alice");

        // Make a signature on the bid hash
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(alicePk, bidHash);
        bytes memory bidSignature = abi.encodePacked(r, s, v);

        // Step 3: Calculate the commitment hash using the getPreConfHash function
        bytes32 commitmentDigest = preConfCommitmentStore.getPreConfHash(
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.bid,
            testCommitment.blockNumber,
            testCommitment.decayStartTimeStamp,
            testCommitment.decayEndTimeStamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(testCommitment.sharedSecretKey)
        );

        // Step 4: Verify the bid hash is correctly generated and not zero
        assert(bidHash != bytes32(0));

        // Step 5: Verify the commitment hash is correctly generated and not zero
        assert(commitmentDigest != bytes32(0));
    }

    function test_Initialize() public view {
        assertEq(preConfCommitmentStore.oracleContract(), oracleContract);
        assertEq(
            address(preConfCommitmentStore.providerRegistry()),
            address(providerRegistry)
        );
        assertEq(
            address(preConfCommitmentStore.bidderRegistry()),
            address(bidderRegistry)
        );
    }

    function test_StoreUnopenedCommitment() public {
        // Step 1: Prepare the commitment information and signature
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        // Optional: Ensure the committer has enough ETH if needed for the operation
        vm.deal(committer, 1 ether);
        vm.prank(committer);
        providerRegistry.registerAndStake{value: 1 ether}(validBLSPubkey);

        // Step 2: Store the commitment
        vm.prank(committer);
        bytes32 commitmentIndex = preConfCommitmentStore
            .storeUnopenedCommitment(
                commitmentDigest,
                commitmentSignature,
                1000
            );

        // Step 3: Verify the results
        // a. Check that the commitment index is correctly generated and not zero
        assert(commitmentIndex != bytes32(0));

        // b. Retrieve the commitment by index and verify its properties
        PreConfCommitmentStore.UnopenedCommitment
            memory commitment = preConfCommitmentStore.getUnopenedCommitment(
                commitmentIndex
            );

        // c. Assertions to verify the stored commitment matches the input
        assertEq(commitment.committer, committer);
        assertEq(commitment.commitmentDigest, commitmentDigest);
        assertEq(commitment.commitmentSignature, commitmentSignature);
    }

    function test_StoreCommitmentFailureDueToTimestampValidation() public {
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        vm.deal(committer, 1 ether);
        vm.prank(committer);

        vm.warp(1000);
        vm.expectRevert("Invalid dispatch timestamp");

        preConfCommitmentStore.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp
        );
    }

    function test_StoreCommitmentFailureDueToTimestampValidationWithNewWindow()
        public
    {
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        assertNotEq(committer, address(0));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        vm.prank(preConfCommitmentStore.owner());
        preConfCommitmentStore.updateCommitmentDispatchWindow(200);
        vm.warp(201 + _testCommitmentAliceBob.dispatchTimestamp);
        vm.expectRevert("Invalid dispatch timestamp");
        preConfCommitmentStore.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp
        );
    }

    function test_UpdateOracle() public {
        address newOracle = address(0x123);
        preConfCommitmentStore.updateOracleContract(newOracle);
        assertEq(preConfCommitmentStore.oracleContract(), newOracle);
    }

    function test_UpdateProviderRegistry() public {
        preConfCommitmentStore.updateProviderRegistry(feeRecipient);
        assertEq(
            address(preConfCommitmentStore.providerRegistry()),
            feeRecipient
        );
    }

    function test_UpdateBidderRegistry() public {
        preConfCommitmentStore.updateBidderRegistry(feeRecipient);
        assertEq(
            address(preConfCommitmentStore.bidderRegistry()),
            feeRecipient
        );
    }

    function test_GetBidHash2() public view {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp
        );
        assertEq(bidHash, _testCommitmentAliceBob.bidDigest);
    }

    function test_GetCommitmentDigest() public {
        (, uint256 bidderPk) = makeAddrAndKey("alice");

        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp
        );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        bytes memory sharedSecretKey = bytes("0xsecret");
        bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            bidHash,
            _bytesToHexString(signature),
            _bytesToHexString(sharedSecretKey)
        );
        assertEq(preConfHash, _testCommitmentAliceBob.commitmentDigest);

        (, uint256 providerPk) = makeAddrAndKey("bob");
        (v, r, s) = vm.sign(providerPk, preConfHash);
        signature = abi.encodePacked(r, s, v);
    }

    function test_StoreCommitment() public {
        (address bidder, ) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 2 ether}(2);

        // Step 1: Verify that the commitment has not been used before
        verifyCommitmentNotUsed(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature
        );

        (address committer, ) = makeAddrAndKey("bob");

        // Step 2: Store the commitment
        bytes32 unopenedIndex = storeCommitment(
            committer,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp,
            _testCommitmentAliceBob.sharedSecretKey
        );

        // Step 3: Move to the next window
        blockTracker.addBuilderAddress("test", committer);
        blockTracker.recordL1Block(2, "test");

        // Step 4: Open the commitment
        bytes32 index = openCommitment(
            bidder,
            unopenedIndex,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.sharedSecretKey
        );

        // Step 5: Verify the stored commitment
        verifyStoredCommitment(
            index,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.sharedSecretKey
        );

        string memory commitmentTxnHash = preConfCommitmentStore
            .getTxnHashFromCommitment(index);
        assertEq(commitmentTxnHash, _testCommitmentAliceBob.txnHash);
    }

    function verifyCommitmentNotUsed(
        string memory txnHash,
        string memory revertingTxHashes,
        uint256 bid,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature
    ) public view returns (bytes32) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );
        bytes memory sharedSecretKey = abi.encodePacked(keccak256("0xsecret"));
        bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        (, bool isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
            .openedCommitments(preConfHash);
        assertEq(isUsed, false);

        return bidHash;
    }

    function storeCommitment(
        address committer,
        uint256 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp,
        bytes memory sharedSecretKey
    ) public returns (bytes32) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );

        bytes32 commitmentDigest = preConfCommitmentStore.getPreConfHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );
        vm.deal(committer, 11 ether);
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 10 ether}(validBLSPubkey);

        bytes32 commitmentIndex = preConfCommitmentStore
            .storeUnopenedCommitment(
                commitmentDigest,
                commitmentSignature,
                dispatchTimestamp
            );
        vm.stopPrank();
        return commitmentIndex;
    }

    function openCommitment(
        address msgSender,
        bytes32 unopenedCommitmentIndex,
        uint256 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) public returns (bytes32) {
        vm.prank(msgSender);
        bytes32 commitmentIndex = preConfCommitmentStore.openCommitment(
            unopenedCommitmentIndex,
            bid,
            blockNumber,
            txnHash,
            revertingTxHashes,
            decayStartTimestamp,
            decayEndTimestamp,
            bidSignature,
            commitmentSignature,
            sharedSecretKey
        );

        return commitmentIndex;
    }

    function verifyStoredCommitment(
        bytes32 index,
        uint256 bid,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) public view {
        PreConfCommitmentStore.OpenedCommitment
            memory commitment = preConfCommitmentStore.getCommitment(index);

        PreConfCommitmentStore.CommitmentParams
            memory commitmentParams = IPreConfCommitmentStore.CommitmentParams({
                txnHash: txnHash,
                revertingTxHashes: revertingTxHashes,
                bid: bid,
                blockNumber: blockNumber,
                decayStartTimeStamp: decayStartTimestamp,
                decayEndTimeStamp: decayEndTimestamp,
                bidHash: commitment.bidHash,
                bidSignature: bidSignature,
                commitmentSignature: commitmentSignature,
                sharedSecretKey: sharedSecretKey
            });

        (, address committerAddress) = preConfCommitmentStore
            .verifyPreConfCommitment(commitmentParams);

        assertNotEq(committerAddress, address(0));
        assertEq(commitment.bid, bid, "Stored bid should match input bid");
        assertEq(
            commitment.blockNumber,
            blockNumber,
            "Stored blockNumber should match input blockNumber"
        );
        assertEq(
            commitment.txnHash,
            txnHash,
            "Stored txnHash should match input txnHash"
        );
        assertEq(
            commitment.bidSignature,
            bidSignature,
            "Stored bidSignature should match input bidSignature"
        );
        assertEq(
            commitment.commitmentSignature,
            commitmentSignature,
            "Stored commitmentSignature should match input commitmentSignature"
        );
    }

    function test_GetCommitment() public {
        (address bidder, ) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);
        uint256 window = WindowFromBlockNumber.getWindowFromBlockNumber(
            _testCommitmentAliceBob.blockNumber,
            blocksPerWindow
        );
        vm.prank(bidder);
        bidderRegistry.depositForWindow{value: 2 ether}(window);
        // Step 1: Verify that the commitment has not been used before
        verifyCommitmentNotUsed(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature
        );
        // Step 2: Store the commitment
        (address committer, ) = makeAddrAndKey("bob");
        providerRegistry.registerAndStake{value: 10 ether}(validBLSPubkey);
        bytes32 commitmentIndex = storeCommitment(
            committer,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp,
            _testCommitmentAliceBob.sharedSecretKey
        );
        PreConfCommitmentStore.UnopenedCommitment
            memory storedCommitment = preConfCommitmentStore
                .getUnopenedCommitment(commitmentIndex);

        assertEq(
            storedCommitment.commitmentDigest,
            _testCommitmentAliceBob.commitmentDigest
        );
        assertEq(
            storedCommitment.commitmentSignature,
            _testCommitmentAliceBob.commitmentSignature
        );
    }

    function test_InitiateSlash() public {
        // Assuming you have a stored commitment
        {
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            uint256 depositWindow = WindowFromBlockNumber
                .getWindowFromBlockNumber(
                    _testCommitmentAliceBob.blockNumber,
                    blocksPerWindow
                );
            bidderRegistry.depositForWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );

            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been set before
            (, bool isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(preConfHash);
            assert(isUsed == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.dispatchTimestamp,
                _testCommitmentAliceBob.sharedSecretKey
            );
            providerRegistry.setPreconfirmationsContract(
                address(preConfCommitmentStore)
            );
            uint256 blockNumber = 2;
            blockTracker.addBuilderAddress("test", committer);
            blockTracker.recordL1Block(blockNumber, "test");
            bytes32 index = openCommitment(
                committer,
                unopenedIndex,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            vm.prank(oracleContract);
            preConfCommitmentStore.initiateSlash(index, 100);

            (, isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(index);
            // Verify that the commitment has been deleted
            assert(isUsed == true);

            assertEq(
                bidderRegistry.lockedFunds(bidder, depositWindow),
                2 ether - _testCommitmentAliceBob.bid
            );
            assertEq(bidderRegistry.providerAmount(committer), 0 ether);
            assertEq(bidder.balance, 3 ether + _testCommitmentAliceBob.bid + 2); // +2 is the slashed funds from provider
        }
        // commitmentDigest value is internal to contract and not asserted
    }

    function test_InitiateReward() public {
        // Assuming you have a stored commitment
        {
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            uint256 depositWindow = WindowFromBlockNumber
                .getWindowFromBlockNumber(
                    _testCommitmentAliceBob.blockNumber,
                    blocksPerWindow
                );
            bidderRegistry.depositForWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (, bool isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(preConfHash);
            assert(isUsed == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.dispatchTimestamp,
                _testCommitmentAliceBob.sharedSecretKey
            );
            blockTracker.addBuilderAddress("test", committer);
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                "test"
            );
            bytes32 index = openCommitment(
                committer,
                unopenedIndex,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            vm.prank(oracleContract);
            preConfCommitmentStore.initiateReward(index, 100);

            (, isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isUsed == true);
            // commitmentDigest value is internal to contract and not asserted
            assertEq(
                bidderRegistry.lockedFunds(bidder, depositWindow),
                2 ether - _testCommitmentAliceBob.bid
            );
        }
    }

    function test_InitiateRewardFullyDecayed() public {
        // Assuming you have a stored commitment
        {
            (address bidder, ) = makeAddrAndKey("alice");
            uint256 depositWindow = WindowFromBlockNumber
                .getWindowFromBlockNumber(
                    _testCommitmentAliceBob.blockNumber,
                    blocksPerWindow
                );
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            bidderRegistry.depositForWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (, bool isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(preConfHash);
            assert(isUsed == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.dispatchTimestamp,
                _testCommitmentAliceBob.sharedSecretKey
            );
            blockTracker.addBuilderAddress("test", committer);
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                "test"
            );
            bytes32 index = openCommitment(
                committer,
                unopenedIndex,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            uint256 window = blockTracker.getCurrentWindow();
            vm.prank(oracleContract);
            preConfCommitmentStore.initiateReward(index, 0);

            (, isUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isUsed == true);
            // commitmentDigest value is internal to contract and not asserted

            assertEq(
                bidderRegistry.lockedFunds(bidder, window),
                2 ether - _testCommitmentAliceBob.bid
            );
            assertEq(bidderRegistry.providerAmount(committer), 0 ether);
            assertEq(bidder.balance, 3 ether + _testCommitmentAliceBob.bid);
        }
    }

    function test_StoreUnopenedCommitmentInsufficientStake() public {
        // Step 1: Prepare the commitment information and signature
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        // Step 2: Attempt to store the commitment and expect it to fail due to insufficient stake
        vm.prank(committer);
        vm.expectRevert("Insufficient stake");
        preConfCommitmentStore.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            1000
        );
    }

    function test_StoreUnopenedCommitmentPendingWithdrawal() public {
        // Step 1: Prepare the commitment information and signature
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        // Ensure the committer has enough ETH for the required stake
        vm.deal(committer, 2 ether);
        vm.prank(committer);
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkey);

        // Request a withdrawal to create a pending withdrawal request
        vm.prank(committer);
        providerRegistry.unstake();

        // Step 2: Attempt to store the commitment and expect it to fail due to pending withdrawal request
        vm.prank(committer);
        vm.expectRevert("Pending withdrawal request");
        preConfCommitmentStore.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            1000
        );
    }

    function _bytesToHexString(
        bytes memory _bytes
    ) internal pure returns (string memory) {
        bytes memory hexChars = "0123456789abcdef";
        bytes memory _string = new bytes(_bytes.length * 2);
        for (uint256 i = 0; i < _bytes.length; ++i) {
            _string[i * 2] = hexChars[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = hexChars[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }

    function _bytes32ToHexString(
        bytes32 _bytes32
    ) internal pure returns (string memory) {
        bytes memory hexChars = "0123456789abcdef";
        bytes memory _string = new bytes(64);
        for (uint8 i = 0; i < 32; ++i) {
            _string[i * 2] = hexChars[uint8(_bytes32[i] >> 4)];
            _string[1 + i * 2] = hexChars[uint8(_bytes32[i] & 0x0f)];
        }
        return string(_string);
    }
}
