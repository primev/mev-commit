// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import {PreConfCommitmentStore} from "../contracts/PreConfCommitmentStore.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/BlockTracker.sol";
import "forge-std/console.sol";

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
    bytes public constant validBLSPubkey = hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";

    function setUp() public {
        _testCommitmentAliceBob = TestCommitment(
            2,
            2,
            "0xkartik",
            "0xkartik",
            10,
            20,
            0x56c06a13be335eba981b780ea45dff258a7c429d0e9d993235ef2d3a7e435df8,
            0x9d954942ad3f6cb41ccd029869be7b28036270b4754665a3783c2d6bf0ef7d08,
            hex"2e7df27808c72d7d5b2543bb63b06c0ae2144e021593b8d2a7cca6a3fb2d9c4b1a82dd2a07266de9364d255bdb709476ad96b826ec855efb528eaff66682997e1c",
            hex"682a1bc8659a0d2616f579ca0e0944d4e21911e65561b7d9ae7dba36c5a5b699248f1c93e8e6f52492ee6c3cb0f34f20f5a97fb2653c5a67200e1de9254389e11c",
            15,
            bytes("0xsecret")
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        blocksPerWindow = 10;

        address providerRegistryProxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (minStake, feeRecipient, feePercent, address(this))
            )
        );
        providerRegistry = ProviderRegistry(payable(providerRegistryProxy));

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(
                BlockTracker.initialize,
                (address(this), blocksPerWindow)
            )
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(
                BidderRegistry.initialize,
                (
                    minStake,
                    feeRecipient,
                    feePercent,
                    address(this),
                    address(blockTracker),
                    blocksPerWindow
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
                    feeRecipient, // Oracle
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

    function test_Initialize() public view {
        assertEq(preConfCommitmentStore.oracle(), feeRecipient);
        assertEq(
            address(preConfCommitmentStore.providerRegistry()),
            address(providerRegistry)
        );
        assertEq(
            address(preConfCommitmentStore.bidderRegistry()),
            address(bidderRegistry)
        );
    }

    function test_storeEncryptedCommitment() public {
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

        // Step 2: Store the commitment
        bytes32 commitmentIndex = preConfCommitmentStore
            .storeEncryptedCommitment(
                commitmentDigest,
                commitmentSignature,
                1000
            );

        // Step 3: Verify the results
        // a. Check that the commitment index is correctly generated and not zero
        assert(commitmentIndex != bytes32(0));

        // b. Retrieve the commitment by index and verify its properties
        PreConfCommitmentStore.EncrPreConfCommitment
            memory commitment = preConfCommitmentStore.getEncryptedCommitment(
                commitmentIndex
            );

        // c. Assertions to verify the stored commitment matches the input
        assertEq(commitment.commiter, committer);
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

        preConfCommitmentStore.storeEncryptedCommitment(
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
        vm.warp(200 + _testCommitmentAliceBob.dispatchTimestamp);
        vm.expectRevert("Invalid dispatch timestamp");
        preConfCommitmentStore.storeEncryptedCommitment(
            commitmentDigest,
            commitmentSignature,
            _testCommitmentAliceBob.dispatchTimestamp
        );
    }

    function test_UpdateOracle() public {
        preConfCommitmentStore.updateOracle(feeRecipient);
        assertEq(preConfCommitmentStore.oracle(), feeRecipient);
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

    function test_GetBidHash() public view {
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
        console.logBytes(signature);
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
        console.logBytes(signature);
    }

    function _bytes32ToHexString(
        bytes32 _bytes32
    ) internal pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(64);
        for (uint8 i = 0; i < 32; i++) {
            _string[i * 2] = HEXCHARS[uint8(_bytes32[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes32[i] & 0x0f)];
        }
        return string(_string);
    }

    function test_StoreCommitment() public {
        (address bidder, ) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);
        vm.prank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(2);

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

        (address commiter, ) = makeAddrAndKey("bob");
        vm.deal(commiter, 5 ether);

        // Step 2: Store the commitment
        bytes32 encryptedIndex = storeCommitment(
            commiter,
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
        blockTracker.addBuilderAddress("test", address(this));
        blockTracker.recordL1Block(2, "test");

        // Step 4: Open the commitment
        bytes32 index = openCommitment(
            bidder,
            encryptedIndex,
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

        (bool isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
            .commitments(preConfHash);
        assertEq(isUsed, false);

        return bidHash;
    }

    function storeCommitment(
        address commiter,
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

        bytes32 commitmentHash = preConfCommitmentStore.getPreConfHash(
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
        vm.startPrank(commiter);
        bytes32 commitmentIndex = preConfCommitmentStore
            .storeEncryptedCommitment(
                commitmentHash,
                commitmentSignature,
                dispatchTimestamp
            );
        vm.stopPrank();
        return commitmentIndex;
    }

    function openCommitment(
        address msgSender,
        bytes32 encryptedCommitmentIndex,
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
            encryptedCommitmentIndex,
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
        PreConfCommitmentStore.PreConfCommitment
            memory commitment = preConfCommitmentStore.getCommitment(index);

        (, address commiterAddress) = preConfCommitmentStore
            .verifyPreConfCommitment(
                txnHash,
                revertingTxHashes,
                bid,
                blockNumber,
                decayStartTimestamp,
                decayEndTimestamp,
                commitment.bidHash,
                bidSignature,
                commitmentSignature,
                sharedSecretKey
            );

        assertNotEq(commiterAddress, address(0));
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
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(window);
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
        (address commiter, ) = makeAddrAndKey("bob");
        vm.deal(commiter, 5 ether);
        bytes32 commitmentIndex = storeCommitment(
            commiter,
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
        PreConfCommitmentStore.EncrPreConfCommitment
            memory storedCommitment = preConfCommitmentStore
                .getEncryptedCommitment(commitmentIndex);

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
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(
                depositWindow
            );

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
            (bool isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(isUsed == false);
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);
            bytes32 encryptedIndex = storeCommitment(
                commiter,
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
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}(validBLSPubkey);
            uint256 blockNumber = 2;
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(blockNumber, "test");
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
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
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateSlash(index, 100);

            (isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been deleted
            assert(isUsed == true);
        }
        // commitmentHash value is internal to contract and not asserted
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
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(
                depositWindow
            );

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
            (bool isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(isUsed == false);
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);
            bytes32 encryptedIndex = storeCommitment(
                commiter,
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
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}(validBLSPubkey);
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                "test"
            );
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
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
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateReward(index, 100);

            (isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been marked as used
            assert(isUsed == true);
            // commitmentHash value is internal to contract and not asserted
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
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(
                depositWindow
            );

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
            (bool isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(isUsed == false);
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);

            bytes32 encryptedIndex = storeCommitment(
                commiter,
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
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}(validBLSPubkey);
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                "test"
            );
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
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
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateReward(index, 0);

            (isUsed, , , , , , , , , , , , , ,) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been marked as used
            assert(isUsed == true);
            // commitmentHash value is internal to contract and not asserted

            assert(bidderRegistry.lockedFunds(bidder, window) == 2 ether);
            assert(bidderRegistry.providerAmount(commiter) == 0 ether);
        }
    }

    function _bytesToHexString(
        bytes memory _bytes
    ) internal pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(_bytes.length * 2);
        for (uint256 i = 0; i < _bytes.length; i++) {
            _string[i * 2] = HEXCHARS[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }
}
