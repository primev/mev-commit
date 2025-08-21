// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {IPreconfManager} from "../../contracts/interfaces/IPreconfManager.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {IProviderRegistry} from "../../contracts/interfaces/IProviderRegistry.sol";
import {MockBLSVerify} from "../precompiles/BLSVerifyPreCompileMockTest.sol";
import {DepositManager} from "../../contracts/core/DepositManager.sol";

contract PreconfManagerTest is Test {
    struct TestCommitment {
        uint256 bidAmt;
        uint256 slashAmt;
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
        uint256[] zkProof;
        bytes bidOptions;
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
    bytes public validBLSPubkey2 =
        hex"90000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82c";
    bytes public validBLSPubkey3 =
        hex"a0000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82d";
    bytes[] public validMultiBLSPubkeys = [
        validBLSPubkey,
        validBLSPubkey2,
        validBLSPubkey3
    ];
    bytes public dummyBLSSignature =
        hex"bbbbbbbbb1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2";
    uint256 public withdrawalDelay;
    uint256 public protocolFeePayoutPeriodBlocks;
    uint256 public bidderWithdrawalPeriodMs;
    uint256[] zkProof;
    address public oracleContract;

    function setUp() public {
        address BLS_VERIFY_ADDRESS = address(0xf0);
        bytes memory code = type(MockBLSVerify).creationCode;
        vm.etch(BLS_VERIFY_ADDRESS, code);

        zkProof = new uint256[](8);
        zkProof[0] = 1;
        zkProof[1] = 2;
        zkProof[2] = 1;
        zkProof[3] = 2;
        zkProof[4] = 1;
        zkProof[5] = 2;
        zkProof[
            6
        ] = 6086978054802466984342546804805639539210350900596926314244208215637076231818;
        zkProof[
            7
        ] = 15801264817036808237903858940451635549338013499819108029453995970938732263800;

        _testCommitmentAliceBob = TestCommitment(
            2,
            0,
            2,
            "0xkartik",
            "0xkartik",
            10,
            20,
            0x8d1f669e1d55329ba0dc133fba063c06c8ae146b8e815732f9951930c807ff7f,
            0x87d7e787de6386cba19d3d5680a8feaa5378c46f1c5e13c622ffcdb354485d23,
            hex"aeed5b345d04360c6ad52d4fb4fce32eec8a552f87686afb39ceea04f9fd1a782b180e4eef5e02af77015292840c541e2681c8e165b44be1d8276aba7211bde21b",
            hex"a0b508b09c6942d73b8feb4feb308ea0b753e14e32cff231f75348f25feb07b02b97743c4c9c493825ec6730e2bc24a513d09a2996a20beac2013f134c047cd71b",
            15,
            zkProof,
            hex""
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);
        withdrawalDelay = 24 hours; // 24 hours
        protocolFeePayoutPeriodBlocks = 100;
        bidderWithdrawalPeriodMs = 10000;
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
                    protocolFeePayoutPeriodBlocks,
                    bidderWithdrawalPeriodMs
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

        uint256 depositManagerMinBalance = 0.01 ether;
        DepositManager depositManager = new DepositManager(address(bidderRegistry), depositManagerMinBalance);
        bidderRegistry.setDepositManagerImpl(address(depositManager));

        // Sets fake block timestamp
        vm.warp(500);
        bidderRegistry.setPreconfManager(address(preconfManager));

        provider = vm.addr(10);
    }

    function test_GetBidHash1() public {
        // Step 1: Prepare the test commitment data
        PreconfManager.CommitmentParams memory testCommitment = IPreconfManager
            .CommitmentParams({
                txnHash: "0xkartik",
                revertingTxHashes: "0xkartik",
                bidAmt: 2,
                slashAmt: 0,
                blockNumber: 2,
                decayStartTimeStamp: 10,
                decayEndTimeStamp: 20,
                bidHash: hex"447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf",
                bidSignature: hex"aeed5b345d04360c6ad52d4fb4fce32eec8a552f87686afb39ceea04f9fd1a782b180e4eef5e02af77015292840c541e2681c8e165b44be1d8276aba7211bde21b",
                commitmentSignature: hex"026b7694e7eaeca9f77718b127e33e20588825820ecc939d751ad2bd21bbd78b71685e2c3f3f76eb37ce8e67843089effd731e93463b8f935cbbf52add269a6d1c",
                zkProof: zkProof
            });

        // Step 2: Calculate the bid hash using the getBidHash function
        bytes32 bidHash = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                testCommitment.bidAmt,
                testCommitment.slashAmt,
                testCommitment.blockNumber,
                testCommitment.decayStartTimeStamp,
                testCommitment.decayEndTimeStamp,
                testCommitment.txnHash,
                testCommitment.revertingTxHashes,
                hex"",
                testCommitment.zkProof,
                hex""
            )
        );

        // Add a alice private key and console log the key
        (, uint256 alicePk) = makeAddrAndKey("alice");

        // Make a signature on the bid hash
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(alicePk, bidHash);
        bytes memory bidSignature = abi.encodePacked(r, s, v);
        // Step 3: Calculate the commitment hash using the getPreConfHash function
        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            bidHash,
            bidSignature,
            zkProof
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
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 1 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

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
            IPreconfManager.OpenCommitmentParams(
                hex"",
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                hex"",
                _testCommitmentAliceBob.zkProof,
                hex""
            )
        );
        assertEq(bidHash, _testCommitmentAliceBob.bidDigest);
    }

    function test_GetCommitmentDigest() public {
        (, uint256 bidderPk) = makeAddrAndKey("alice");

        bytes32 bidHash = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                _testCommitmentAliceBob.bidAmt,
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.blockNumber,
                _testCommitmentAliceBob.decayStartTimestamp,
                _testCommitmentAliceBob.decayEndTimestamp,
                _testCommitmentAliceBob.txnHash,
                _testCommitmentAliceBob.revertingTxHashes,
                hex"",
                _testCommitmentAliceBob.zkProof,
                hex""
            )
        );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bytes memory signature = abi.encodePacked(r, s, v);

        assertEq(signature, _testCommitmentAliceBob.bidSignature);

        bytes32 preConfHash = preconfManager.getPreConfHash(
            bidHash,
            signature,
            zkProof
        );
        assertEq(preConfHash, _testCommitmentAliceBob.commitmentDigest);

        (, uint256 providerPk) = makeAddrAndKey("bob");
        (v, r, s) = vm.sign(providerPk, preConfHash);
        signature = abi.encodePacked(r, s, v);

        assertEq(signature, _testCommitmentAliceBob.commitmentSignature);
    }

    function test_StoreCommitment() public {
        (address committer, ) = makeAddrAndKey("bob");
        (address bidder, ) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 2 ether}(committer);

        verifyCommitmentNotUsed(_testCommitmentAliceBob);

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
            _testCommitmentAliceBob.slashAmt,
            _testCommitmentAliceBob.zkProof,
            _testCommitmentAliceBob.bidOptions
        );

        // Step 3: Record the block
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
            _testCommitmentAliceBob.slashAmt,
            _testCommitmentAliceBob.zkProof,
            _testCommitmentAliceBob.bidOptions
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
            _testCommitmentAliceBob.bidDigest,
            _testCommitmentAliceBob.bidSignature,
            _testCommitmentAliceBob.commitmentSignature,
            _testCommitmentAliceBob.slashAmt,
            _testCommitmentAliceBob.zkProof
        );
    }

    function verifyCommitmentNotUsed(
        TestCommitment memory c
    ) public view returns (bytes32) {
        bytes32 bidHash = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                c.bidAmt,
                c.slashAmt,
                c.blockNumber,
                c.decayStartTimestamp,
                c.decayEndTimestamp,
                c.txnHash,
                c.revertingTxHashes,
                hex"",
                c.zkProof,
                hex""
            )
        );

        bytes32 preConfHash = preconfManager.getPreConfHash(
            bidHash,
            c.bidSignature,
            c.zkProof
        );

        (, bool isSettled, , , , , , , , , , , , ) = preconfManager
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
        uint256 slashAmt,
        uint256[] memory zkproof,
        bytes memory bidOptions
    ) public returns (bytes32) {
        bytes32 bidHash = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                bidAmt,
                slashAmt,
                blockNumber,
                decayStartTimestamp,
                decayEndTimestamp,
                txnHash,
                revertingTxHashes,
                hex"",
                zkproof,
                bidOptions
            )
        );

        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            bidHash,
            bidSignature,
            zkproof
        );
        vm.deal(committer, 11 ether);
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 10 ether}();
        for (uint256 i = 0; i < validMultiBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validMultiBLSPubkeys[i],
                dummyBLSSignature
            );
        }

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
        uint256 slashAmt,
        uint256[] memory zkproof,
        bytes memory bidOptions
    ) public returns (bytes32) {
        vm.prank(msgSender);
        bytes32 commitmentIndex = preconfManager.openCommitment(
            IPreconfManager.OpenCommitmentParams(
                unopenedCommitmentIndex,
                bidAmt,
                slashAmt,
                blockNumber,
                decayStartTimestamp,
                decayEndTimestamp,
                txnHash,
                revertingTxHashes,
                bidSignature,
                zkproof,
                bidOptions
            )
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
        bytes32 bidHash,
        bytes memory bidSignature,
        bytes memory commitmentSignature,
        uint256 slashAmt,
        uint256[] memory zkproof
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
                bidHash: bidHash,
                bidSignature: bidSignature,
                commitmentSignature: commitmentSignature,
                slashAmt: slashAmt,
                zkProof: zkproof
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
            commitment.commitmentSignature,
            commitmentSignature,
            "Stored commitmentSignature should match input commitmentSignature"
        );
    }

    function test_GetCommitment() public {
        (address bidder, ) = makeAddrAndKey("alice");
        vm.deal(bidder, 5 ether);
        vm.prank(bidder);
        bidderRegistry.depositAsBidder{value: 2 ether}(provider);
        // Step 1: Verify that the commitment has not been used before
        verifyCommitmentNotUsed(_testCommitmentAliceBob);
        // Step 2: Store the commitment
        (address committer, ) = makeAddrAndKey("bob");
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
            _testCommitmentAliceBob.slashAmt,
            _testCommitmentAliceBob.zkProof,
            _testCommitmentAliceBob.bidOptions
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
            (address committer, ) = makeAddrAndKey("bob");
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            bidderRegistry.depositAsBidder{value: 2 ether}(committer);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(_testCommitmentAliceBob);

            bytes32 preConfHash = preconfManager.getPreConfHash(
                bidHash,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.zkProof
            );

            // Verify that the commitment has not been set before
            (, bool isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);

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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
            );
            uint256 oneHundredPercent = providerRegistry.ONE_HUNDRED_PERCENT();

            vm.prank(oracleContract);
            preconfManager.initiateSlash(index, oneHundredPercent);

            (, isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been deleted
            assert(isSettled == true);

            assertEq(
                bidderRegistry.getDeposit(bidder, committer),
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
            (address committer, ) = makeAddrAndKey("bob");
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            bidderRegistry.depositAsBidder{value: 2 ether}(committer);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(_testCommitmentAliceBob);
            bytes32 preConfHash = preconfManager.getPreConfHash(
                bidHash,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.zkProof
            );

            // Verify that the commitment has not been used before
            (, bool isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);

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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
            );
            vm.prank(oracleContract);
            preconfManager.initiateReward(index, 100);

            (, isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isSettled == true);
            // commitmentDigest value is internal to contract and not asserted
            assertEq(
                bidderRegistry.getDeposit(bidder, committer),
                2 ether - _testCommitmentAliceBob.bidAmt
            );
        }
    }

    function test_InitiateRewardFullyDecayed() public {
        // Assuming you have a stored commitment
        {
            (address committer, ) = makeAddrAndKey("bob");
            (address bidder, ) = makeAddrAndKey("alice");
            vm.deal(bidder, 5 ether);
            vm.prank(bidder);
            bidderRegistry.depositAsBidder{value: 2 ether}(committer);

            // Step 1: Verify that the commitment has not been used before
            bytes32 bidHash = verifyCommitmentNotUsed(_testCommitmentAliceBob);
            bytes32 preConfHash = preconfManager.getPreConfHash(
                bidHash,
                _testCommitmentAliceBob.bidSignature,
                _testCommitmentAliceBob.zkProof
            );

            // Verify that the commitment has not been used before
            (, bool isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(preConfHash);
            assert(isSettled == false);

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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
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
                _testCommitmentAliceBob.slashAmt,
                _testCommitmentAliceBob.zkProof,
                _testCommitmentAliceBob.bidOptions
            );
            vm.prank(oracleContract);
            preconfManager.initiateReward(index, 0);

            (, isSettled, , , , , , , , , , , , ) = preconfManager
                .openedCommitments(index);
            // Verify that the commitment has been marked as used
            assert(isSettled == true);
            // commitmentDigest value is internal to contract and not asserted

            assertEq(
                bidderRegistry.getDeposit(bidder, committer),
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
        vm.startPrank(committer);
        providerRegistry.registerAndStake{value: 2 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

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

        bidderRegistry.depositAsBidder{value: 2 ether}(provider);

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
            abi.encode(
                testCommitment.txnHash,
                bidder,
                testCommitment.blockNumber
            )
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
            IPreconfManager.OpenCommitmentParams(
                unopenedIndex2,
                testCommitment2.bidAmt,
                testCommitment2.slashAmt,
                testCommitment2.blockNumber,
                testCommitment2.decayStartTimestamp,
                testCommitment2.decayEndTimestamp,
                testCommitment2.txnHash,
                testCommitment2.revertingTxHashes,    
                testCommitment2.bidSignature,
                testCommitment2.zkProof,
                testCommitment2.bidOptions
            )
        );
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
                testCommitment.slashAmt,
                testCommitment.zkProof,
                testCommitment.bidOptions
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
                testCommitment.slashAmt,
                testCommitment.zkProof,
                testCommitment.bidOptions
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
        testCommitment2.decayStartTimestamp += 1;
        testCommitment2.decayEndTimestamp += 1;
        testCommitment2.dispatchTimestamp += 1;

        // Recompute bidHash and bidSignature
        bytes32 bidHash2 = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                testCommitment2.bidAmt,
                testCommitment2.slashAmt,
                testCommitment2.blockNumber,
                testCommitment2.decayStartTimestamp,
                testCommitment2.decayEndTimestamp,
                testCommitment2.txnHash,
                testCommitment2.revertingTxHashes,
                hex"",
                testCommitment2.zkProof,
                testCommitment2.bidOptions
            )
        );

        testCommitment2.bidDigest = bidHash2;
        testCommitment2.bidSignature = signHash(bidderPk, bidHash2);

        // Recompute commitmentDigest and commitmentSignature
        bytes32 commitmentDigest2 = preconfManager.getPreConfHash(
            bidHash2,
            testCommitment2.bidSignature,
            testCommitment2.zkProof
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
