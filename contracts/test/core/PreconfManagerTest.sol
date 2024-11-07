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
    BidderRegistry public bidderRegistry;
    bytes public validBLSPubkey =
        hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    bytes[] public validBLSPubkeys = [validBLSPubkey];
    bytes public validBLSPubkey2 = hex"90000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82c";
    bytes public validBLSPubkey3 = hex"a0000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82d";
    bytes[] public validMultiBLSPubkeys = [validBLSPubkey, validBLSPubkey2, validBLSPubkey3];
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
            0x447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf,
            0xa7f6241be0c5055f054fcbe03d98a1920f0ab874039474401323d8d95930a076,
            hex"5cd1f790192a0ab79661c48f39e77937a6de537ccf6b428682583d13d30294cb113cea12822f821c064c9db918667bf74490535b35b4ef4f28f5d67b133ec22e1b",
            hex"026b7694e7eaeca9f77718b127e33e20588825820ecc939d751ad2bd21bbd78b71685e2c3f3f76eb37ce8e67843089effd731e93463b8f935cbbf52add269a6d1c",
            15,
            bytes("0xsecret")
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        withdrawalDelay = 24  hours ; // 24 hours
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
                (address(this), address(this))
            )
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));
vm.prank(address(this));
        blockTracker.setProviderRegistry(address(providerRegistry));
        address bidderRegistryProxy = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(
                BidderRegistry.initialize,
                (
                    feeRecipient,
                    feePercent,
                    address(this),
                    address(blockTracker),
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
                    500
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
        PreconfManager.CommitmentParams memory testCommitment = IPreconfManager
            .CommitmentParams({
                txnHash: "0xkartik",
                revertingTxHashes: "0xkartik",
                bidAmt: 2,
                blockNumber: 2,
                decayStartTimeStamp: 10,
                decayEndTimeStamp: 20,
                sharedSecretKey: bytes("0xsecret"),
                bidHash: hex"447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf",
                bidSignature: hex"5cd1f790192a0ab79661c48f39e77937a6de537ccf6b428682583d13d30294cb113cea12822f821c064c9db918667bf74490535b35b4ef4f28f5d67b133ec22e1b",
                commitmentSignature: hex"026b7694e7eaeca9f77718b127e33e20588825820ecc939d751ad2bd21bbd78b71685e2c3f3f76eb37ce8e67843089effd731e93463b8f935cbbf52add269a6d1c"
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
            bidSignature,
            testCommitment.sharedSecretKey
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
        providerRegistry.registerAndStake{value: 1 ether}(validBLSPubkeys);

        // Step 2: Store the commitment
        vm.prank(committer);
        bytes32 commitmentIndex = preconfManager.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            1000
        );

        // Step 3: Verify the results
        // a. Check that the commitment index is correctly generated and not zero
        assert(commitmentIndex != bytes32(0));

        // b. Retrieve the commitment by index and verify its properties
        PreconfManager.UnopenedCommitment memory commitment = preconfManager
            .getUnopenedCommitment(commitmentIndex);

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
        uint256 minTime = block.timestamp -
            preconfManager.commitmentDispatchWindow();

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

        vm.prank(preconfManager.owner());
        preconfManager.updateCommitmentDispatchWindow(200);

        vm.warp(201 + _testCommitmentAliceBob.dispatchTimestamp);

        // Calculate the minimum valid timestamp for dispatching the commitment
        uint256 minTime = block.timestamp -
            preconfManager.commitmentDispatchWindow();

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
        assertEq(address(preconfManager.providerRegistry()), feeRecipient);
    }

    function test_UpdateBidderRegistry() public {
        preconfManager.updateBidderRegistry(feeRecipient);
        assertEq(address(preconfManager.bidderRegistry()), feeRecipient);
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
            signature,
            sharedSecretKey
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
        blockTracker.recordL1Block(2, validBLSPubkey2);

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
            bidSignature,
            sharedSecretKey
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
            bidSignature,
            sharedSecretKey
        );
        vm.deal(committer, 11 ether);
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 10 ether}(validMultiBLSPubkeys);

        bytes32 commitmentIndex = preconfManager.storeUnopenedCommitment(
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
        PreconfManager.OpenedCommitment memory commitment = preconfManager
            .getCommitment(index);

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

        (, address committerAddress) = preconfManager.verifyPreConfCommitment(
            commitmentParams
        );

        assertNotEq(committerAddress, address(0));
        assertEq(
            commitment.bidAmt,
            bidAmt,
            "Stored bid should match input bid"
        );
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
            _testCommitmentAliceBob.blockNumber
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
        providerRegistry.registerAndStake{value: 10 ether}(validBLSPubkeys);
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
            memory storedCommitment = preconfManager.getUnopenedCommitment(
                commitmentIndex
            );

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
                    _testCommitmentAliceBob.blockNumber
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
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.sharedSecretKey
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
            blockTracker.recordL1Block(blockNumber, validBLSPubkey);
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
                _testCommitmentAliceBob.sharedSecretKey
            );
            vm.prank(oracleContract);
            preconfManager.initiateSlash(index, 100);

            (, isSettled, , , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been deleted
            assert(isSettled == true);

            assertEq(
                bidderRegistry.lockedFunds(bidder, depositWindow),
                2 ether - _testCommitmentAliceBob.bidAmt
            );
            assertEq(bidderRegistry.providerAmount(committer), 0 ether);
            assertEq(
                bidder.balance,
                3 ether + _testCommitmentAliceBob.bidAmt + 2
            ); // +2 is the slashed funds from provider
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
                    _testCommitmentAliceBob.blockNumber
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
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.sharedSecretKey
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
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                validBLSPubkey
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
                    _testCommitmentAliceBob.blockNumber
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
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.sharedSecretKey
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
            blockTracker.recordL1Block(
                _testCommitmentAliceBob.blockNumber,
                validBLSPubkey
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
        providerRegistry.registerAndStake{value: 2 ether}(validBLSPubkeys);

        // Request a withdrawal to create a pending withdrawal request
        vm.prank(committer);
        providerRegistry.unstake();

        // Step 2: Attempt to store the commitment and expect it to fail due to pending withdrawal request
        vm.prank(committer);
        vm.expectRevert(
            abi.encodeWithSelector(
                IProviderRegistry.PendingWithdrawalRequest.selector,
                committer
            )
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
        bytes32 unopenedIndex1 = storeFirstCommitment(
            committer,
            testCommitment
        );

        blockTracker.recordL1Block(testCommitment.blockNumber, validBLSPubkey);

        openFirstCommitment(bidder, unopenedIndex1, testCommitment);

        bytes32 txnHashAndBidder = keccak256(
            abi.encode(testCommitment.txnHash, bidder)
        );
        // Verify that the first commitment is processed
        assertTrue(
            preconfManager.processedTxnHashes(txnHashAndBidder),
            "First txnHash should be marked as processed"
        );

        // Prepare and store the second commitment with the same txnHash
        TestCommitment memory testCommitment2 = prepareSecondCommitment(
            bidderPk,
            committerPk,
            testCommitment
        );

        bytes32 unopenedIndex2 = storeSecondCommitment(
            committer,
            testCommitment2
        );

        blockTracker.recordL1Block(testCommitment2.blockNumber, validBLSPubkey);

        // Attempt to open the second commitment with the same txnHash
        vm.prank(bidder);
        vm.expectRevert(
            abi.encodeWithSelector(
                IPreconfManager.TxnHashAlreadyProcessed.selector,
                testCommitment2.txnHash,
                bidder
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
            testCommitment2.sharedSecretKey
        );
    }

    function depositForBidder(
        address bidder,
        uint64 blockNumber
    ) internal returns (uint256) {
        vm.prank(bidder);
        uint256 depositWindow = WindowFromBlockNumber.getWindowFromBlockNumber(
            blockNumber
        );
        bidderRegistry.depositForWindow{value: 2 ether}(depositWindow);
        return depositWindow;
    }

    function storeFirstCommitment(
        address committer,
        TestCommitment memory testCommitment
    ) internal returns (bytes32) {
        return
            storeCommitment(
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
        return
            openCommitment(
                bidder,
                unopenedIndex,
                testCommitment.bidAmt,
                testCommitment.blockNumber,
                testCommitment.txnHash,
                testCommitment.revertingTxHashes,
                testCommitment.decayStartTimestamp,
                testCommitment.decayEndTimestamp,
                testCommitment.bidSignature,
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
            testCommitment2.bidSignature,
            testCommitment2.sharedSecretKey
        );

        testCommitment2.commitmentDigest = commitmentDigest2;
        testCommitment2.commitmentSignature = signHash(
            committerPk,
            commitmentDigest2
        );

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

    function signHash(
        uint256 privateKey,
        bytes32 hash
    ) internal pure returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privateKey, hash);
        return abi.encodePacked(r, s, v);
    }
}