// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

import {Test} from "forge-std/Test.sol";
import {Oracle} from "../../contracts/core/Oracle.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../../contracts/utils/WindowFromBlockNumber.sol";
import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import {MockBLSVerify} from "../precompiles/BLSVerifyPreCompileMockTest.sol";
import {IPreconfManager} from "../../contracts/interfaces/IPreconfManager.sol";

contract OracleTest is Test {
    using ECDSA for bytes32;
    address public owner;
    Oracle public oracle;
    PreconfManager public preconfManager;
    uint256 public feePercent;
    uint256 public minStake;
    address public feeRecipient;
    ProviderRegistry public providerRegistry;
    uint256 public testNumber;
    uint64 public testNumber2;
    BidderRegistry public bidderRegistry;
    BlockTracker public blockTracker;
    uint64 public dispatchTimestampTesting;
    bytes public sharedSecretKey;
    uint256[] public zkProof;
    bytes public constant validBLSPubkey =
        hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    bytes[] public validBLSPubkeys = [validBLSPubkey];
    bytes public dummyBLSSignature =
        hex"bbbbbbbbb1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2";
    uint256 public constant withdrawalDelay = 24 hours; // 24 hours
    uint256 public constant protocolFeePayoutPeriodBlocks = 100;

    struct CommitmentParamsSimple {
        uint64 bid;
        uint64 blockNumber;
        string txnHash;
        string revertingTxHashes;
        uint256 bidderPk;
        uint256 signerPk; // used for signing the commitment digest
        address provider; // also used as the committer
        uint64 dispatchTimestamp;
        uint256 slashAmt;
        uint256[] zkProof;
    }

    struct UnopenedCommitmentParams {
        address committerAddress;
        uint64 bid;
        uint64 blockNumber;
        string txnHash;
        string revertingTxHashes;
        uint64 decayStartTimestamp;
        uint64 decayEndTimestamp;
        uint256 bidderPk;
        uint256 signerPk;
        uint64 dispatchTimestamp;
        uint256 slashAmt;
        uint256[] zkProof;
    }

    // Events to match against
    event BlockDataRequested(uint256 blockNumber);
    event BlockDataReceived(
        string[] txnList,
        uint256 blockNumber,
        string blockBuilderName
    );
    event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash);
    event FundsRetrieved(
        bytes32 indexed commitmentDigest,
        uint256 window,
        uint256 amount
    );

    function setUp() public {
        address BLS_VERIFY_ADDRESS = address(0xf0);
        bytes memory code = type(MockBLSVerify).creationCode;
        vm.etch(BLS_VERIFY_ADDRESS, code);

        testNumber = 2;
        testNumber2 = 2;
        sharedSecretKey = bytes("0xsecret");
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

        feePercent = 10 * 1e16; // 10%
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);

        address proxy = Upgrades.deployUUPSProxy(
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
        providerRegistry = ProviderRegistry(payable(proxy));

        address ownerInstance = 0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3;

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(
                BlockTracker.initialize,
                (ownerInstance, ownerInstance)
            )
        );
        blockTracker = BlockTracker(payable(blockTrackerProxy));

        vm.startPrank(ownerInstance);
        blockTracker.setProviderRegistry(address(providerRegistry));
        vm.stopPrank();

        address proxy3 = Upgrades.deployUUPSProxy(
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
        bidderRegistry = BidderRegistry(payable(proxy3));

        address proxy4 = Upgrades.deployUUPSProxy(
            "PreconfManager.sol",
            abi.encodeCall(
                PreconfManager.initialize,
                (
                    address(providerRegistry),
                    address(bidderRegistry),
                    feeRecipient,
                    address(this),
                    address(blockTracker),
                    500
                )
            )
        );
        preconfManager = PreconfManager(payable(proxy4));

        vm.deal(ownerInstance, 5 ether);
        vm.startPrank(ownerInstance);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 2 ether}(window + 1);

        address oracleProxy = Upgrades.deployUUPSProxy(
            "Oracle.sol",
            abi.encodeCall(
                Oracle.initialize,
                (
                    address(preconfManager),
                    address(blockTracker),
                    ownerInstance,
                    ownerInstance
                )
            )
        );
        oracle = Oracle(payable(oracleProxy));

        vm.stopPrank();

        preconfManager.updateOracleContract(address(oracle));
        bidderRegistry.setPreconfManager(address(preconfManager));
        providerRegistry.setPreconfManager(address(preconfManager));

        dispatchTimestampTesting = 1000;
        vm.warp(1010);
    }

    function test_process_commitment_payment_payout() public {
        string
            memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(
            WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2
        );
        uint64 bid = 2;
        uint256 slashAmt = 0;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

        // Construct and store the commitment.
        CommitmentParamsSimple memory params = CommitmentParamsSimple({
            bid: bid,
            blockNumber: blockNumber,
            txnHash: txn,
            revertingTxHashes: revertingTxHashes,
            bidderPk: bidderPk,
            signerPk: providerPk,
            provider: provider,
            dispatchTimestamp: dispatchTimestampTesting,
            slashAmt: slashAmt,
            zkProof: zkProof
        });
        bytes32 index = constructAndStoreCommitment(params);

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        oracle.processBuilderCommitmentForBlockNumber(
            index,
            blockNumber,
            provider,
            false,
            50 * bidderRegistry.PRECISION()
        );
        vm.stopPrank();

        // Verify the bidder received the correct payout.
        assertEq(bidderRegistry.getProviderAmount(provider), (bid * 50) / 100);
    }

    function test_process_commitment_slash() public {
        string
            memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(
            WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2
        );
        uint64 bid = 200;
        uint256 slashAmt = 0;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

        CommitmentParamsSimple memory params = CommitmentParamsSimple({
            bid: bid,
            blockNumber: blockNumber,
            txnHash: txn,
            revertingTxHashes: revertingTxHashes,
            bidderPk: bidderPk,
            signerPk: providerPk,
            provider: provider,
            dispatchTimestamp: dispatchTimestampTesting,
            slashAmt: slashAmt,
            zkProof: zkProof
        });
        bytes32 index = constructAndStoreCommitment(params);

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index,
            blockNumber,
            provider,
            true,
            50 * bidderRegistry.PRECISION()
        );
        vm.stopPrank();

        // Verify the provider’s stake decreased appropriately.
        assertEq(
            providerRegistry.getProviderStake(provider) + ((bid * 55) / 100),
            250 ether
        );
    }

    function test_process_commitment_slash_and_reward() public {
        string
            memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(
            WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2
        );
        uint64 bid = 100;
        uint256 slashAmt = 0;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        uint256 residualAfterDecay = 50;

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

        CommitmentParamsSimple memory params1 = CommitmentParamsSimple({
            bid: bid,
            blockNumber: blockNumber,
            txnHash: txn1,
            revertingTxHashes: revertingTxHashes,
            bidderPk: bidderPk,
            signerPk: providerPk,
            provider: provider,
            dispatchTimestamp: dispatchTimestampTesting,
            slashAmt: slashAmt,
            zkProof: zkProof
        });
        bytes32 index1 = constructAndStoreCommitment(params1);

        CommitmentParamsSimple memory params2 = CommitmentParamsSimple({
            bid: bid,
            blockNumber: blockNumber,
            txnHash: txn2,
            revertingTxHashes: revertingTxHashes,
            bidderPk: bidderPk,
            signerPk: providerPk,
            provider: provider,
            dispatchTimestamp: dispatchTimestampTesting,
            slashAmt: slashAmt,
            zkProof: zkProof
        });
        bytes32 index2 = constructAndStoreCommitment(params2);

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index1,
            blockNumber,
            provider,
            true,
            bidderRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, false);
        oracle.processBuilderCommitmentForBlockNumber(
            index2,
            blockNumber,
            provider,
            false,
            50 * providerRegistry.PRECISION()
        );
        vm.stopPrank();

        assertEq(
            providerRegistry.getProviderStake(provider),
            250 ether - ((bid * 110) / 100)
        );
        assertEq(
            bidderRegistry.getProviderAmount(provider),
            (((bid * (providerRegistry.ONE_HUNDRED_PERCENT() - feePercent)) /
                providerRegistry.ONE_HUNDRED_PERCENT()) * residualAfterDecay) /
                100
        );
    }

    function test_process_commitment_slash_multiple() public {
        string
            memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        string
            memory txn3 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d10";
        string
            memory txn4 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d11";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = 201;
        uint64 bid = 5;
        uint256 slashAmt = 0;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = WindowFromBlockNumber.getWindowFromBlockNumber(
            blockNumber
        );
        vm.startPrank(bidder);
        bidderRegistry.depositForWindow{value: 250 ether}(window);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(
            CommitmentParamsSimple({
                bid: bid,
                blockNumber: blockNumber,
                txnHash: txn1,
                revertingTxHashes: revertingTxHashes,
                bidderPk: bidderPk,
                signerPk: providerPk,
                provider: provider,
                dispatchTimestamp: dispatchTimestampTesting,
                slashAmt: slashAmt,
                zkProof: zkProof
            })
        );
        bytes32 index2 = constructAndStoreCommitment(
            CommitmentParamsSimple({
                bid: bid,
                blockNumber: blockNumber,
                txnHash: txn2,
                revertingTxHashes: revertingTxHashes,
                bidderPk: bidderPk,
                signerPk: providerPk,
                provider: provider,
                dispatchTimestamp: dispatchTimestampTesting,
                slashAmt: slashAmt,
                zkProof: zkProof
            })
        );
        bytes32 index3 = constructAndStoreCommitment(
            CommitmentParamsSimple({
                bid: bid,
                blockNumber: blockNumber,
                txnHash: txn3,
                revertingTxHashes: revertingTxHashes,
                bidderPk: bidderPk,
                signerPk: providerPk,
                provider: provider,
                dispatchTimestamp: dispatchTimestampTesting,
                slashAmt: slashAmt,
                zkProof: zkProof
            })
        );
        bytes32 index4 = constructAndStoreCommitment(
            CommitmentParamsSimple({
                bid: bid,
                blockNumber: blockNumber,
                txnHash: txn4,
                revertingTxHashes: revertingTxHashes,
                bidderPk: bidderPk,
                signerPk: providerPk,
                provider: provider,
                dispatchTimestamp: dispatchTimestampTesting,
                slashAmt: slashAmt,
                zkProof: zkProof
            })
        );

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index1,
            blockNumber,
            provider,
            true,
            bidderRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index2,
            blockNumber,
            provider,
            true,
            bidderRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index3, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index3,
            blockNumber,
            provider,
            true,
            bidderRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index4, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index4,
            blockNumber,
            provider,
            true,
            bidderRegistry.ONE_HUNDRED_PERCENT()
        );
        vm.stopPrank();
        assertEq(
            providerRegistry.getProviderStake(provider),
            250 ether - bid * 4
        );
        assertEq(bidderRegistry.getProviderAmount(provider), 0);
    }

    function test_process_commitment_reward_multiple() public {
        // In this test we now pass decayStartTimeStamp and decayEndTimeStamp when opening commitments.
        string[] memory txnHashes = new string[](4);
        txnHashes[
            0
        ] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        txnHashes[
            1
        ] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        txnHashes[
            2
        ] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d10";
        txnHashes[
            3
        ] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d11";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(
            WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2
        );
        uint64 bid = 5;
        uint256 slashAmt = 0;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = blockTracker.getCurrentWindow();
        vm.startPrank(bidder);
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        for (uint256 i = 0; i < validBLSPubkeys.length; i++) {
            providerRegistry.addVerifiedBLSKey(
                validBLSPubkeys[i],
                dummyBLSSignature
            );
        }
        vm.stopPrank();

        bytes32[] memory commitments = new bytes32[](4);
        bytes[] memory bidSignatures = new bytes[](4);
        bytes[] memory commitmentSignatures = new bytes[](4);
        for (uint256 i = 0; i < commitments.length; ++i) {
            UnopenedCommitmentParams
                memory paramsUnopened = UnopenedCommitmentParams({
                    committerAddress: provider,
                    bid: bid,
                    blockNumber: blockNumber,
                    txnHash: txnHashes[i],
                    revertingTxHashes: revertingTxHashes,
                    decayStartTimestamp: 10,
                    decayEndTimestamp: 20,
                    bidderPk: bidderPk,
                    signerPk: providerPk,
                    dispatchTimestamp: dispatchTimestampTesting,
                    slashAmt: slashAmt,
                    zkProof: zkProof
                });
            (
                commitments[i],
                bidSignatures[i],
                commitmentSignatures[i]
            ) = constructAndStoreUnopenedCommitment(paramsUnopened);
        }

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, validBLSPubkey);
        vm.stopPrank();

        for (uint256 i = 0; i < commitments.length; ++i) {
            vm.startPrank(provider);
            // Open the commitment with the two extra decay parameters.
            preconfManager.openCommitment(
                IPreconfManager.OpenCommitmentParams(
                    commitments[i],
                    bid,
                    slashAmt,
                    blockNumber,
                    10, // decayStartTimeStamp
                    20, // decayEndTimeStamp
                    txnHashes[i],
                    revertingTxHashes,
                    bidSignatures[i],
                    zkProof
                )
            );
            vm.stopPrank();
        }
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        for (uint256 i = 0; i < commitments.length; ++i) {
            vm.expectEmit(true, false, false, true);
            emit CommitmentProcessed(commitments[i], false);
            oracle.processBuilderCommitmentForBlockNumber(
                commitments[i],
                blockNumber,
                provider,
                false,
                bidderRegistry.ONE_HUNDRED_PERCENT()
            );
        }
        vm.stopPrank();
        assertEq(providerRegistry.getProviderStake(provider), 250 ether);
        assertEq(bidderRegistry.getProviderAmount(provider), 4 * bid);
    }

    function constructAndStoreCommitment(
        CommitmentParamsSimple memory params
    ) public returns (bytes32 commitmentIndex) {
        bytes32 bidHash = getBidHash(
            params.txnHash,
            params.revertingTxHashes,
            params.bid,
            params.blockNumber,
            params.slashAmt,
            params.zkProof
        );
        bytes memory bidSignature = getBidSignature(params.bidderPk, bidHash);
        bytes32 commitmentDigest = getCommitmentDigest(bidHash, bidSignature);
        bytes memory commitmentSignature = getCommitmentSignature(
            params.signerPk,
            commitmentDigest
        );
        bytes32 unopenedCommitmentIndex = storeUnopenedCommitment(
            params.provider,
            commitmentDigest,
            commitmentSignature,
            params.dispatchTimestamp
        );
        recordBlockData(validBLSPubkey, params.blockNumber);
        commitmentIndex = openCommitment(
            params.provider,
            unopenedCommitmentIndex,
            params.bid,
            params.blockNumber,
            params.txnHash,
            params.revertingTxHashes,
            bidSignature,
            params.slashAmt,
            params.zkProof
        );
        return commitmentIndex;
    }

    function constructAndStoreUnopenedCommitment(
        UnopenedCommitmentParams memory params
    )
        public
        returns (
            bytes32 commitmentIndex,
            bytes memory bidSignature,
            bytes memory commitmentSignature
        )
    {
        bytes32 bidHash = preconfManager.getBidHash(
            IPreconfManager.OpenCommitmentParams(
                hex"",
                params.bid,
                params.slashAmt,
                params.blockNumber,
                params.decayStartTimestamp,
                params.decayEndTimestamp,
                params.txnHash,
                params.revertingTxHashes,
                hex"",
                params.zkProof
            )
        );
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(params.bidderPk, bidHash);
        bidSignature = abi.encodePacked(r, s, v);
        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            bidHash,
            bidSignature,
            params.zkProof
        );
        (v, r, s) = vm.sign(params.signerPk, commitmentDigest);
        commitmentSignature = abi.encodePacked(r, s, v);
        vm.startPrank(params.committerAddress);
        commitmentIndex = preconfManager.storeUnopenedCommitment(
            commitmentDigest,
            commitmentSignature,
            params.dispatchTimestamp
        );
        vm.stopPrank();
        return (commitmentIndex, bidSignature, commitmentSignature);
    }

    function getBidHash(
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 bid,
        uint64 blockNumber,
        uint256 slashAmt,
        uint256[] memory zkproof
    ) public view returns (bytes32) {
        return
            preconfManager.getBidHash(
                IPreconfManager.OpenCommitmentParams(
                    hex"",
                    bid,
                    slashAmt,
                    blockNumber,
                    10,
                    20,
                    txnHash,
                    revertingTxHashes,
                    hex"",
                    zkproof
                )
            );
    }

    function getBidSignature(
        uint256 bidderPk,
        bytes32 bidHash
    ) public pure returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        return abi.encodePacked(r, s, v);
    }

    function getCommitmentDigest(
        bytes32 bidHash,
        bytes memory bidSignature
    ) public view returns (bytes32) {
        return preconfManager.getPreConfHash(bidHash, bidSignature, zkProof);
    }

    function getCommitmentSignature(
        uint256 signerPk,
        bytes32 commitmentDigest
    ) public pure returns (bytes memory) {
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPk, commitmentDigest);
        return abi.encodePacked(r, s, v);
    }

    function storeUnopenedCommitment(
        address provider,
        bytes32 commitmentDigest,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp
    ) public returns (bytes32) {
        vm.startPrank(provider);
        bytes32 unopenedCommitmentIndex = preconfManager
            .storeUnopenedCommitment(
                commitmentDigest,
                commitmentSignature,
                dispatchTimestamp
            );
        vm.stopPrank();
        return unopenedCommitmentIndex;
    }

    function recordBlockData(
        bytes memory blsPubKey,
        uint64 blockNumber
    ) public {
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        blockTracker.recordL1Block(blockNumber, blsPubKey);
        vm.stopPrank();
    }

    function openCommitment(
        address provider,
        bytes32 unopenedCommitmentIndex,
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes memory bidSignature,
        uint256 slashAmt,
        uint256[] memory zkproof
    ) public returns (bytes32) {
        vm.startPrank(provider);
        bytes32 commitmentIndex = preconfManager.openCommitment(
            IPreconfManager.OpenCommitmentParams(
                unopenedCommitmentIndex,
                bid,
                slashAmt,
                blockNumber,
                10, // decayStartTimeStamp
                20, // decayEndTimeStamp
                txnHash,
                revertingTxHashes,
                bidSignature,
                zkproof
            )
        );
        vm.stopPrank();
        return commitmentIndex;
    }

    function _bytesToHexString(
        bytes memory _bytes
    ) internal pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(_bytes.length * 2);
        for (uint256 i = 0; i < _bytes.length; ++i) {
            _string[i * 2] = HEXCHARS[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }
}
