// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {IPreconfManager} from "../../contracts/interfaces/IPreconfManager.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../../contracts/utils/WindowFromBlockNumber.sol";
import {IProviderRegistry} from "../../contracts/interfaces/IProviderRegistry.sol";

contract PreconfManagerTest is Test {
    struct TestCommitment {
        uint256 bidAmt;
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
    PreconfManager public preconfManager;
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
            0xc311dfce5df35601ec4b562dfdf048e22cb66373fb6ba5160e83dac3d72f0d2b,
            0x6ebb8f592c9e75ea8c4a2403884e97237ce3a559da30461317391b388f440eac,
            hex"77731700031fe79fba2dae5614bc44af07167f39a42ae5c1b4e136035870d0fb6ee98e239bc27fe289b47d9dc9cd461de2f0a50528bc8f24c2290fb4286821a41b",
            hex"0cc9c50fb1fd57db6f31226db5b97af3537c3f1f9699a94e3702e6db58131fc462d5126bda70f0f41fa812f215304dd4b766457430bb8fa229e0ae3ad71368631c",
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
            "PreconfManager.sol",
            abi.encodeCall(
                PreconfManager.initialize,
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
        preconfManager = PreconfManager(payable(preconfStoreProxy));

        // Sets fake block timestamp
        vm.warp(500);
        bidderRegistry.setPreconfManager(address(preconfManager));
    }

    function test_GetBidHash1() public {
        // Step 1: Prepare the test commitment data
        PreconfManager.CommitmentParams
            memory testCommitment = IPreconfManager.CommitmentParams({
                txnHash: "0xkartik",
                revertingTxHashes: "0xkartik",
                bidAmt: 2,
                blockNumber: 2,
                decayStartTimeStamp: 10,
                decayEndTimeStamp: 20,
                sharedSecretKey: bytes("0xsecret"),
                bidHash: hex"c311dfce5df35601ec4b562dfdf048e22cb66373fb6ba5160e83dac3d72f0d2b",
                bidSignature: hex"77731700031fe79fba2dae5614bc44af07167f39a42ae5c1b4e136035870d0fb6ee98e239bc27fe289b47d9dc9cd461de2f0a50528bc8f24c2290fb4286821a41b",
                commitmentSignature: hex"5b3000290d4f347b94146eb37f66d5368aed18fb8713bf78620abe40ae3de7f635f7ed161801c31ea10e736d88e6fd2a2286bbd59385161dd24c9fefd2568f341b"
            });
        // Step 2: Calculate the bid hash using the getBidHash function
        bytes32 bidHash = preconfManager.getBidHash(
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.bidAmt,
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
        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.bidAmt,
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
        assertEq(preconfManager.oracleContract(), oracleContract);
        assertEq(
            address(preconfManager.providerRegistry()),
            address(providerRegistry)
        );
        assertEq(
            address(preconfManager.bidderRegistry()),
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
        bytes32 commitmentIndex = preconfManager
            .storeUnopenedCommitment(
                commitmentDigest,
                commitmentSignature,
                1000
            );

        // Step 3: Verify the results
        // a. Check that the commitment index is correctly generated and not zero
        assert(commitmentIndex != bytes32(0));

        // b. Retrieve the commitment by index and verify its properties
        PreconfManager.UnopenedCommitment
            memory commitment = preconfManager.getUnopenedCommitment(
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
        // Calculate the minimum valid timestamp for dispatching the commitment
        uint256 minTime = block.timestamp - preconfManager.commitmentDispatchWindow();

        vm.expectRevert(
            abi.encodeWithSelector(
                IPreconfManager.InvalidDispatchTimestamp.selector,
                minTime,
                _testCommitmentAliceBob.dispatchTimestamp
            )
        );

        preconfManager.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp
        );
    }

    function test_StoreCommitmentFailureDueToTimestampValidationWithNewWindow() public {
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

        vm.prank(preconfManager.owner());
        preconfManager.updateCommitmentDispatchWindow(200);

        vm.warp(201 + _testCommitmentAliceBob.dispatchTimestamp);

        // Calculate the minimum valid timestamp for dispatching the commitment
        uint256 minTime = block.timestamp - preconfManager.commitmentDispatchWindow();

        vm.expectRevert(
            abi.encodeWithSelector(
                IPreconfManager.InvalidDispatchTimestamp.selector,
                minTime,
                _testCommitmentAliceBob.dispatchTimestamp
            )
        );
        preconfManager.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp
        );
    }

    function test_UpdateOracle() public {
        address newOracle = address(0x123);
        preconfManager.updateOracleContract(newOracle);
        assertEq(preconfManager.oracleContract(), newOracle);
    }

    function test_UpdateProviderRegistry() public {
        preconfManager.updateProviderRegistry(feeRecipient);
        assertEq(
            address(preconfManager.providerRegistry()),
            feeRecipient
        );
    }

    function test_UpdateBidderRegistry() public {
        preconfManager.updateBidderRegistry(feeRecipient);
        assertEq(
            address(preconfManager.bidderRegistry()),
            feeRecipient
        );
    }

    function test_GetBidHash2() public view {
        bytes32 bidHash = preconfManager.getBidHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bidAmt,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp
        );
        assertEq(bidHash, _testCommitmentAliceBob.bidDigest);
    }

    function test_GetCommitmentDigest() public {
        (, uint256 bidderPk) = makeAddrAndKey("alice");

        bytes32 bidHash = preconfManager.getBidHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bidAmt,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp
        );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bytes memory signature = abi.encodePacked(r, s, v);
        bytes memory sharedSecretKey = bytes("0xsecret");
        bytes32 preConfHash = preconfManager.getPreConfHash(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bidAmt,
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
            _testCommitmentAliceBob.bidAmt,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature
        );

        (address committer, ) = makeAddrAndKey("bob");

        // Step 2: Store the commitment
        bytes32 unopenedIndex = storeCommitment(
            committer,
            _testCommitmentAliceBob.bidAmt,
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
            _testCommitmentAliceBob.bidAmt,
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
            _testCommitmentAliceBob.bidAmt,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.revertingTxHashes,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.sharedSecretKey
        );

        string memory commitmentTxnHash = preconfManager
            .getTxnHashFromCommitment(index);
        assertEq(commitmentTxnHash, _testCommitmentAliceBob.txnHash);
    }

    function verifyCommitmentNotUsed(
        string memory txnHash,
        string memory revertingTxHashes,
        uint256 bidAmt,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature
    ) public view returns (bytes32) {
        bytes32 bidHash = preconfManager.getBidHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );
        bytes memory sharedSecretKey = abi.encodePacked(keccak256("0xsecret"));
        bytes32 preConfHash = preconfManager.getPreConfHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        (, bool isSettled, , , , , , , , , , , , , ) = preconfManager
            .openedCommitments(preConfHash);
        assertEq(isSettled, false);

        return bidHash;
    }

    function storeCommitment(
        address committer,
        uint256 bidAmt,
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
        bytes32 bidHash = preconfManager.getBidHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );

        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
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

        bytes32 commitmentIndex = preconfManager
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
        uint256 bidAmt,
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
        bytes32 commitmentIndex = preconfManager.openCommitment(
            unopenedCommitmentIndex,
            bidAmt,
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
        uint256 bidAmt,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) public view {
        PreconfManager.OpenedCommitment
            memory commitment = preconfManager.getCommitment(index);

        PreconfManager.CommitmentParams
            memory commitmentParams = IPreconfManager.CommitmentParams({
                txnHash: txnHash,
                revertingTxHashes: revertingTxHashes,
                bidAmt: bidAmt,
                blockNumber: blockNumber,
                decayStartTimeStamp: decayStartTimestamp,
                decayEndTimeStamp: decayEndTimestamp,
                bidHash: commitment.bidHash,
                bidSignature: bidSignature,
                commitmentSignature: commitmentSignature,
                sharedSecretKey: sharedSecretKey
            });

        (, address committerAddress) = preconfManager
            .verifyPreConfCommitment(commitmentParams);

        assertNotEq(committerAddress, address(0));
        assertEq(commitment.bidAmt, bidAmt, "Stored bid should match input bid");
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
            _testCommitmentAliceBob.bidAmt,
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
            _testCommitmentAliceBob.bidAmt,
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
        PreconfManager.UnopenedCommitment
            memory storedCommitment = preconfManager
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
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );

            bytes32 preConfHash = preconfManager.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been set before
            (, bool isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bidAmt,
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
            providerRegistry.setPreconfManager(address(preconfManager));
            uint256 blockNumber = 2;
            blockTracker.addBuilderAddress("test", committer);
            blockTracker.recordL1Block(blockNumber, "test");
            bytes32 index = openCommitment(
                committer,
                unopenedIndex,
                _testCommitmentAliceBob.bidAmt,
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
            preconfManager.initiateSlash(index);

            (, isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been deleted
            assert(isSettled == true);

            assertEq(
                bidderRegistry.lockedFunds(bidder, depositWindow),
                2 ether - _testCommitmentAliceBob.bidAmt
            );
            assertEq(bidderRegistry.providerAmount(committer), 0 ether);
            assertEq(bidder.balance, 3 ether + _testCommitmentAliceBob.bidAmt + 2); // +2 is the slashed funds from provider
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
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preconfManager.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (, bool isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bidAmt,
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
                _testCommitmentAliceBob.bidAmt,
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
            preconfManager.initiateReward(index, 100);

            (, isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isSettled == true);
            // commitmentDigest value is internal to contract and not asserted
            assertEq(
                bidderRegistry.lockedFunds(bidder, depositWindow),
                2 ether - _testCommitmentAliceBob.bidAmt
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
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preconfManager.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (, bool isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);
            (address committer, ) = makeAddrAndKey("bob");

            bytes32 unopenedIndex = storeCommitment(
                committer,
                _testCommitmentAliceBob.bidAmt,
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
                _testCommitmentAliceBob.bidAmt,
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
            preconfManager.initiateReward(index, 0);

            (, isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isSettled == true);
            // commitmentDigest value is internal to contract and not asserted

            assertEq(
                bidderRegistry.lockedFunds(bidder, window),
                2 ether - _testCommitmentAliceBob.bidAmt
            );
            assertEq(bidderRegistry.providerAmount(committer), 0 ether);
            assertEq(bidder.balance, 3 ether + _testCommitmentAliceBob.bidAmt);
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
        vm.expectRevert(
            abi.encodeWithSelector(
                IProviderRegistry.InsufficientStake.selector,
                0,
                1e18 // min stake
            )
        );
        preconfManager.storeUnopenedCommitment(
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
        vm.expectRevert(
            abi.encodeWithSelector(IProviderRegistry.PendingWithdrawalRequest.selector, committer)
        );
        preconfManager.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            1000
        );
    }

    function test_OpenCommitmentWithDuplicateTxnHash() public {
        // Set up the initial commitment data
        TestCommitment memory testCommitment = _testCommitmentAliceBob;

        // Set up the initial commitment
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);

        depositForBidder(bidder, testCommitment.blockNumber);

        (address committer, uint256 committerPk) = makeAddrAndKey("bob");
        vm.deal(committer, 11 ether);

        // Store and open the first commitment
        bytes32 unopenedIndex1 = storeFirstCommitment(committer, testCommitment);
        blockTracker.addBuilderAddress("test", committer);
        blockTracker.recordL1Block(testCommitment.blockNumber, "test");

        openFirstCommitment(
            bidder,
            unopenedIndex1,
            testCommitment
        );

        // Verify that the first commitment is processed
        assertTrue(
            preconfManager.processedTxnHashes(testCommitment.txnHash),
            "First txnHash should be marked as processed"
        );

        // Prepare and store the second commitment with the same txnHash
        TestCommitment memory testCommitment2 = prepareSecondCommitment(
            bidderPk,
            committerPk,
            testCommitment
        );

        bytes32 unopenedIndex2 = storeSecondCommitment(committer, testCommitment2);

        blockTracker.addBuilderAddress("test2", committer);
        blockTracker.recordL1Block(testCommitment2.blockNumber, "test2");

        // Attempt to open the second commitment with the same txnHash
        vm.prank(bidder);
        vm.expectRevert(
            abi.encodeWithSelector(
                IPreconfManager.TxnHashAlreadyProcessed.selector,
                testCommitment2.txnHash
            )
        );
        preconfManager.openCommitment(
            unopenedIndex2,
            testCommitment2.bidAmt,
            testCommitment2.blockNumber,
            testCommitment2.txnHash,
            testCommitment2.revertingTxHashes,
            testCommitment2.decayStartTimestamp,
            testCommitment2.decayEndTimestamp,
            testCommitment2.bidSignature,
            testCommitment2.commitmentSignature,
            testCommitment2.sharedSecretKey
        );
    }

    function depositForBidder(address bidder, uint64 blockNumber) internal returns (uint256) {
        vm.prank(bidder);
        uint256 depositWindow = WindowFromBlockNumber.getWindowFromBlockNumber(
            blockNumber,
            blocksPerWindow
        );
        bidderRegistry.depositForWindow{value: 2 ether}(depositWindow);
        return depositWindow;
    }

    function registerProvider(address committer, bytes memory blsPubkey) internal {
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 10 ether}(blsPubkey);
        vm.stopPrank();
    }

    function storeFirstCommitment(
        address committer,
        TestCommitment memory testCommitment
    ) internal returns (bytes32) {
        return storeCommitment(
            committer,
            testCommitment.bidAmt,
            testCommitment.blockNumber,
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.decayStartTimestamp,
            testCommitment.decayEndTimestamp,
            testCommitment.bidSignature,
            testCommitment.commitmentSignature,
            testCommitment.dispatchTimestamp,
            testCommitment.sharedSecretKey
        );
    }

    function openFirstCommitment(
        address bidder,
        bytes32 unopenedIndex,
        TestCommitment memory testCommitment
    ) internal returns (bytes32) {
        return openCommitment(
            bidder,
            unopenedIndex,
            testCommitment.bidAmt,
            testCommitment.blockNumber,
            testCommitment.txnHash,
            testCommitment.revertingTxHashes,
            testCommitment.decayStartTimestamp,
            testCommitment.decayEndTimestamp,
            testCommitment.bidSignature,
            testCommitment.commitmentSignature,
            testCommitment.sharedSecretKey
        );
    }

    function prepareSecondCommitment(
        uint256 bidderPk,
        uint256 committerPk,
        TestCommitment memory testCommitment
    ) internal view returns (TestCommitment memory) {
        TestCommitment memory testCommitment2 = testCommitment;

        // Update the fields for the second commitment
        testCommitment2.bidAmt += 1;
        testCommitment2.blockNumber += 1;
        testCommitment2.decayStartTimestamp += 1;
        testCommitment2.decayEndTimestamp += 1;
        testCommitment2.dispatchTimestamp += 1;

        // Recompute bidHash and bidSignature
        bytes32 bidHash2 = preconfManager.getBidHash(
            testCommitment2.txnHash,
            testCommitment2.revertingTxHashes,
            testCommitment2.bidAmt,
            testCommitment2.blockNumber,
            testCommitment2.decayStartTimestamp,
            testCommitment2.decayEndTimestamp
        );

        testCommitment2.bidDigest = bidHash2;
        testCommitment2.bidSignature = signHash(bidderPk, bidHash2);

        // Recompute commitmentDigest and commitmentSignature
        bytes32 commitmentDigest2 = preconfManager.getPreConfHash(
            testCommitment2.txnHash,
            testCommitment2.revertingTxHashes,
            testCommitment2.bidAmt,
            testCommitment2.blockNumber,
            testCommitment2.decayStartTimestamp,
            testCommitment2.decayEndTimestamp,
            bidHash2,
            _bytesToHexString(testCommitment2.bidSignature),
            _bytesToHexString(testCommitment2.sharedSecretKey)
        );

        testCommitment2.commitmentDigest = commitmentDigest2;
        testCommitment2.commitmentSignature = signHash(committerPk, commitmentDigest2);

        return testCommitment2;
    }

    function storeSecondCommitment(
        address committer,
        TestCommitment memory testCommitment2
    ) internal returns (bytes32) {
        vm.startPrank(committer);
        bytes32 unopenedIndex = preconfManager.storeUnopenedCommitment(
            testCommitment2.commitmentDigest,
            testCommitment2.commitmentSignature,
            testCommitment2.dispatchTimestamp
        );
        vm.stopPrank();
        return unopenedIndex;
    }

    function signHash(uint256 privateKey, bytes32 hash) internal pure returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKey, hash);
        return abi.encodePacked(r, s, v);
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
