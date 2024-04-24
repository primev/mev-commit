// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";

import {PreConfCommitmentStore} from "../contracts/PreConfirmations.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/BlockTracker.sol";
import "forge-std/console.sol";

contract TestPreConfCommitmentStore is Test {
    struct TestCommitment {
        uint64 bid;
        uint64 blockNumber;
        string txnHash;
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
    PreConfCommitmentStore internal preConfCommitmentStore;
    uint16 internal feePercent;
    uint256 internal minStake;
    address internal provider;
    address internal feeRecipient;
    ProviderRegistry internal providerRegistry;
    BlockTracker internal blockTracker;

    BidderRegistry internal bidderRegistry;

    function setUp() public {
        _testCommitmentAliceBob = TestCommitment(
            2,
            2,
            "0xkartik",
            10,
            20,
            0xa0327970258c49b922969af74d60299a648c50f69a2d98d6ab43f32f64ac2100,
            0x65618f8f9e46b8f0790c621ca2989cfe4c949594a4a3a81261baa682e8883840,
            hex"876c1216c232828be9fabb14981c8788cebdf6ed66e563c4a2ccc82a577d052543207aeeb158a32d8977736797ae250c63ef69a82cd85b727da21e20d030fb311b",
            hex"bfea9167927707ae7586ed3bba8565999f8b7ad874b2dd4f175caf81084c0d0a17f9599daf5b3f2773757408aa4b44875c95df0f4150cfb295f95273e1fefdd01b",
            15,
            bytes("0xsecret")
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        providerRegistry = new ProviderRegistry(
            minStake,
            feeRecipient,
            feePercent,
            address(this)
        );
        blockTracker = new BlockTracker(address(this));
        bidderRegistry = new BidderRegistry(
            minStake,
            feeRecipient,
            feePercent,
            address(this),
            address(blockTracker)
        );

        preConfCommitmentStore = new PreConfCommitmentStore(
            address(providerRegistry), // Provider Registry
            address(bidderRegistry), // User Registry
            feeRecipient, // Oracle
            address(this),
            address(blockTracker), // Block Tracker
            500
        );

        // Sets fake block timestamp
        vm.warp(16);
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
            .storeEncryptedCommitment(commitmentDigest, commitmentSignature, 1000);

        // Step 3: Verify the results
        // a. Check that the commitment index is correctly generated and not zero
        assert(commitmentIndex != bytes32(0));

        // b. Retrieve the commitment by index and verify its properties
        PreConfCommitmentStore.EncrPreConfCommitment
            memory commitment = preConfCommitmentStore.getEncryptedCommitment(
                commitmentIndex
            );

        // c. Assertions to verify the stored commitment matches the input
        assertEq(commitment.commitmentUsed, false);
        assertEq(commitment.commiter, committer);
        assertEq(commitment.commitmentDigest, commitmentDigest);
        assertEq(commitment.commitmentSignature, commitmentSignature);
    }


    function test_StoreCommitmentFailureDueToTimestampValidation() public {
        // bytes32 bidHash = preConfCommitmentStore.getBidHash(
        //     _testCommitmentAliceBob.txnHash,
        //     _testCommitmentAliceBob.bid,
        //     _testCommitmentAliceBob.blockNumber,
        //     _testCommitmentAliceBob.decayStartTimestamp,
        //     _testCommitmentAliceBob.decayEndTimestamp
        // );
        // (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        // // Wallet memory kartik = vm.createWallet('test wallet');
        // (uint8 v,bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        // bytes memory signature = abi.encodePacked(r, s, v);

        // vm.deal(bidder, 200000 ether);
        // vm.prank(bidder);
        // bidderRegistry.depositForSpecificWindow{value: 2 ether}(2);

        // (bytes32 digest, address recoveredAddress) =  preConfCommitmentStore.verifyBid(
        //     _testCommitmentAliceBob.bid, 
        //     _testCommitmentAliceBob.blockNumber, 
        //     _testCommitmentAliceBob.decayStartTimestamp, 
        //     _testCommitmentAliceBob.decayEndTimestamp, 
        //     _testCommitmentAliceBob.txnHash, 
        //     signature);
        
        // assertEq(bidder, recoveredAddress);
        // assertEq(digest, bidHash);
        // vm.warp(1000);
        // vm.expectRevert("Invalid dispatch timestamp, block.timestamp - dispatchTimestamp < commitment_dispatch_window");
        // preConfCommitmentStore.storeCommitment(
        //     _testCommitmentAliceBob.bid,
        //     _testCommitmentAliceBob.blockNumber,
        //     _testCommitmentAliceBob.txnHash,
        //     _testCommitmentAliceBob.decayStartTimestamp,
        //     _testCommitmentAliceBob.decayEndTimestamp,
        //     signature,
        //     _testCommitmentAliceBob.commitmentSignature,
        //     _testCommitmentAliceBob.dispatchTimestamp
        // );
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

        vm.deal(committer, 1 ether);
        vm.prank(committer);

        vm.warp(1000);
        vm.expectRevert("Invalid dispatch timestamp, block.timestamp - dispatchTimestamp < commitment_dispatch_window");

        preConfCommitmentStore
            .storeEncryptedCommitment(commitmentDigest, commitmentSignature, _testCommitmentAliceBob.dispatchTimestamp);
    }

    function test_StoreCommitmentFailureDueToTimestampValidationWithNewWindow() public {
        bytes32 commitmentDigest = keccak256(
            abi.encodePacked("commitment data")
        );
        (address committer, uint256 committerPk) = makeAddrAndKey("committer");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(
            committerPk,
            commitmentDigest
        );
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        vm.prank(preConfCommitmentStore.owner());
        preConfCommitmentStore.updateCommitmentDispatchWindow(200);
        vm.warp(200 + _testCommitmentAliceBob.dispatchTimestamp);
        vm.expectRevert("Invalid dispatch timestamp, block.timestamp - dispatchTimestamp < commitment_dispatch_window");
        preConfCommitmentStore
            .storeEncryptedCommitment(commitmentDigest, commitmentSignature, _testCommitmentAliceBob.dispatchTimestamp);
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

    function test_GetBidHash() public {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            _testCommitmentAliceBob.txnHash,
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
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature
        );

        // Step 2: Store the commitment
        bytes32 encryptedIndex = storeCommitment(
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
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
        uint64 bid,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature
    ) public returns (bytes32) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );
        bytes memory sharedSecretKey = abi.encodePacked(keccak256("0xsecret"));
        bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        (bool commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
            .commitments(preConfHash);
        assertEq(commitmentUsed, false);

        return bidHash;
    }

    function storeCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp,
        bytes memory sharedSecretKey
    ) internal returns (bytes32) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );

        bytes32 commitmentHash = preConfCommitmentStore.getPreConfHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        bytes32 commitmentIndex = preConfCommitmentStore
            .storeEncryptedCommitment(commitmentHash, commitmentSignature, dispatchTimestamp);

        return commitmentIndex;
    }

    function openCommitment(
        address msgSender,
        bytes32 encryptedCommitmentIndex,
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) internal returns (bytes32) {
        vm.prank(msgSender);
        bytes32 commitmentIndex = preConfCommitmentStore.openCommitment(
            encryptedCommitmentIndex,
            bid,
            blockNumber,
            txnHash,
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
        uint64 bid,
        uint64 blockNumber,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        string memory txnHash,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) public {
        PreConfCommitmentStore.PreConfCommitment
            memory commitment = preConfCommitmentStore.getCommitment(index);

        (, address commiterAddress) = preConfCommitmentStore
            .verifyPreConfCommitment(
                txnHash,
                bid,
                blockNumber,
                decayStartTimestamp,
                decayEndTimestamp,
                commitment.bidHash,
                bidSignature,
                commitmentSignature,
                sharedSecretKey
            );

        bytes32[] memory commitments = preConfCommitmentStore
            .getCommitmentsByCommitter(commiterAddress);

        assert(commitments.length >= 1);

        assertEq(
            commitment.commitmentUsed,
            false,
            "Commitment should have been marked as used"
        );
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
        vm.prank(bidder);
        uint256 window = blockTracker.getWindowFromBlockNumber(_testCommitmentAliceBob.blockNumber);
        vm.prank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(window);
        // Step 1: Verify that the commitment has not been used before
        verifyCommitmentNotUsed(
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            _testCommitmentAliceBob.bidSignature
        );
        // Step 2: Store the commitment
        bytes32 commitmentIndex = storeCommitment(
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
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
            uint256 depositWindow = blockTracker.getWindowFromBlockNumber(_testCommitmentAliceBob.blockNumber);
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );

            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (bool commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(commitmentUsed == false);
            bytes32 encryptedIndex = storeCommitment(
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
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
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}();
            uint256 blockNumber = 2;
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(blockNumber, "test");
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateSlash(index, 100);

            (commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been marked as used
            assert(commitmentUsed == true);
        }
        // commitmentHash value is internal to contract and not asserted
    }

    function test_InitiateReward() public {
        // Assuming you have a stored commitment
        {
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            uint256 depositWindow = blockTracker.getWindowFromBlockNumber(_testCommitmentAliceBob.blockNumber);
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (bool commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(commitmentUsed == false);
            bytes32 encryptedIndex = storeCommitment(
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.dispatchTimestamp,
                _testCommitmentAliceBob.sharedSecretKey
            );
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}();
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(_testCommitmentAliceBob.blockNumber, "test");
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateReward(index, 100);

            (commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been marked as used
            assert(commitmentUsed == true);
            // commitmentHash value is internal to contract and not asserted
        }
    }

    function test_InitiateRewardFullyDecayed() public {
        // Assuming you have a stored commitment
        {
            (address bidder, ) = makeAddrAndKey("alice");
            uint64 blockNumber = 66;
            uint256 depositWindow = blockTracker.getWindowFromBlockNumber(blockNumber);
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            bidderRegistry.depositForSpecificWindow{value: 2 ether}(depositWindow);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature
            );
            bytes32 preConfHash = preConfCommitmentStore.getPreConfHash(
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.bid,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                bidHash,
                _bytesToHexString(_testCommitmentAliceBob.bidSignature),
                _bytesToHexString(_testCommitmentAliceBob.sharedSecretKey)
            );

            // Verify that the commitment has not been used before
            (bool commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(preConfHash);
            assert(commitmentUsed == false);
            bytes32 encryptedIndex = storeCommitment(
                _testCommitmentAliceBob.bid,
                blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.dispatchTimestamp,
                _testCommitmentAliceBob.sharedSecretKey
            );
            (address commiter, ) = makeAddrAndKey("bob");
            vm.deal(commiter, 5 ether);
            vm.prank(commiter);
            providerRegistry.registerAndStake{value: 4 ether}();
            blockTracker.addBuilderAddress("test", commiter);
            blockTracker.recordL1Block(blockNumber, "test");
            bytes32 index = openCommitment(
                commiter,
                encryptedIndex,
                _testCommitmentAliceBob.bid,
                blockNumber,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.commitmentSignature,
                _testCommitmentAliceBob.sharedSecretKey
            );
            uint256 window = blockTracker.getCurrentWindow();
            vm.prank(feeRecipient);
            preConfCommitmentStore.initiateReward(index, 0);

            (commitmentUsed, , , , , , , , , , , , , ) = preConfCommitmentStore
                .commitments(index);
            // Verify that the commitment has been marked as used
            assert(commitmentUsed == true);
            // commitmentHash value is internal to contract and not asserted

            assert(bidderRegistry.lockedFunds(bidder, window) == 2 ether);
            assert(bidderRegistry.providerAmount(commiter) == 0 ether);
        }
    }

    function _bytesToHexString(
        bytes memory _bytes
    ) public pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(_bytes.length * 2);
        for (uint256 i = 0; i < _bytes.length; i++) {
            _string[i * 2] = HEXCHARS[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }
}