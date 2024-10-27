// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {Oracle} from "../../contracts/core/Oracle.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {BidderRegistry} from "../../contracts/core/BidderRegistry.sol";
import {BlockTracker} from "../../contracts/core/BlockTracker.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {WindowFromBlockNumber} from "../../contracts/utils/WindowFromBlockNumber.sol";
import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";

contract OracleTest is Test {
    using ECDSA for bytes32;
    address public owner;
    Oracle public oracle;
    PreconfManager public preconfManager;
    uint16 public feePercent;
    uint256 public minStake;
    address public feeRecipient;
    ProviderRegistry public providerRegistry;
    uint256 public testNumber;
    uint64 public testNumber2;
    BidderRegistry public bidderRegistry;
    BlockTracker public blockTracker;
    TestCommitment internal _testCommitmentAliceBob;
    uint64 public dispatchTimestampTesting;
    bytes public sharedSecretKey;
    bytes public constant validBLSPubkey = hex"80000cddeec66a800e00b0ccbb62f12298073603f5209e812abbac7e870482e488dd1bbe533a9d44497ba8b756e1e82b";
    uint256 public constant withdrawalDelay = 24 * 3600; // 24 hours
    uint256 public constant protocolFeePayoutPeriodBlocks = 100;
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
        testNumber = 2;
        testNumber2 = 2;
        sharedSecretKey = bytes("0xsecret");
        _testCommitmentAliceBob = TestCommitment(
            2,
            2,
            "0xkartik",
            10,
            20,
            0xa0327970258c49b922969af74d60299a648c50f69a2d98d6ab43f32f64ac2100,
            0x54c118e537dd7cf63b5388a5fc8322f0286a978265d0338b108a8ca9d155dccc,
            hex"876c1216c232828be9fabb14981c8788cebdf6ed66e563c4a2ccc82a577d052543207aeeb158a32d8977736797ae250c63ef69a82cd85b727da21e20d030fb311b",
            hex"ec0f11f77a9e96bb9c2345f031a5d12dca8d01de8a2e957cf635be14802f9ad01c6183688f0c2672639e90cc2dce0662d9bea3337306ca7d4b56dd80326aaa231b",
            1000
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);

        address proxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(
                ProviderRegistry.initialize,
                (minStake, feeRecipient, feePercent, address(this), withdrawalDelay, protocolFeePayoutPeriodBlocks)
            )
        );
        providerRegistry = ProviderRegistry(payable(proxy));

        address ownerInstance = 0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3;

        address blockTrackerProxy = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (ownerInstance, ownerInstance))
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

        // We set the system time to 1010 and dispatchTimestamps for testing to 1000
        dispatchTimestampTesting = 1000;
        vm.warp(1010);
    }

    function test_process_commitment_payment_payout() public {
        string
            memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        uint64 bid = 2;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}(validBLSPubkey);
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        oracle.processBuilderCommitmentForBlockNumber(
            index,
            blockNumber,
            provider,
            false,
            50
        );
        vm.stopPrank();
        assertEq(
            bidderRegistry.getProviderAmount(provider),
            (bid * (50)) / 100
        );
    }

    function test_process_commitment_slash() public {
        string
            memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string
            memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = 200;
        uint64 bid = 200;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}(validBLSPubkey);
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index,
            blockNumber,
            provider,
            true,
            50
        );
        vm.stopPrank();
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
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        uint64 bid = 100;
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
        providerRegistry.registerAndStake{value: 250 ether}(validBLSPubkey);
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn1,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );
        bytes32 index2 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn2,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index1,
            blockNumber,
            provider,
            true,
            100
        );

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, false);
        oracle.processBuilderCommitmentForBlockNumber(
            index2,
            blockNumber,
            provider,
            false,
            50
        );
        vm.stopPrank();
        assertEq(providerRegistry.getProviderStake(provider), 250 ether - ((bid*110)/100));
        assertEq(
            bidderRegistry.getProviderAmount(provider),
            (((bid * (100 - feePercent)) / 100) * residualAfterDecay) / 100
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
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = WindowFromBlockNumber.getWindowFromBlockNumber(blockNumber);
        vm.startPrank(bidder);
        bidderRegistry.depositForWindow{value: 250 ether}(window);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}(validBLSPubkey);
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn1,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );
        bytes32 index2 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn2,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );
        bytes32 index3 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn3,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );
        bytes32 index4 = constructAndStoreCommitment(
            bid,
            blockNumber,
            txn4,
            revertingTxHashes,
            bidderPk,
            providerPk,
            provider,
            dispatchTimestampTesting
        );

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index1,
            blockNumber,
            provider,
            true,
            100
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index2,
            blockNumber,
            provider,
            true,
            100
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index3, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index3,
            blockNumber,
            provider,
            true,
            100
        );
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index4, true);
        oracle.processBuilderCommitmentForBlockNumber(
            index4,
            blockNumber,
            provider,
            true,
            100
        );
        vm.stopPrank();
        assertEq(providerRegistry.getProviderStake(provider), 250 ether - bid * 4);
        assertEq(bidderRegistry.getProviderAmount(provider), 0);
    }

    function test_process_commitment_reward_multiple() public {
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
        string memory revertingTxHashes = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d12";
        uint64 blockNumber = uint64(WindowFromBlockNumber.BLOCKS_PER_WINDOW + 2);
        uint64 bid = 5;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = blockTracker.getCurrentWindow();
        vm.startPrank(bidder);
        bidderRegistry.depositForWindow{value: 250 ether}(window + 1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}(validBLSPubkey);
        vm.stopPrank();

        bytes32[] memory commitments = new bytes32[](4);
        bytes[] memory bidSignatures = new bytes[](4);
        bytes[] memory commitmentSignatures = new bytes[](4);
        for (uint256 i = 0; i < commitments.length; ++i) {
            (
                commitments[i],
                bidSignatures[i],
                commitmentSignatures[i]
            ) = constructAndStoreUnopenedCommitment(
                provider,
                bid,
                blockNumber,
                txnHashes[i],
                revertingTxHashes,
                10,
                20,
                bidderPk,
                providerPk,
                dispatchTimestampTesting
            );
        }

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, validBLSPubkey);
        vm.stopPrank();

        for (uint256 i = 0; i < commitments.length; ++i) {
            vm.startPrank(provider);
            preconfManager.openCommitment(
                commitments[i],
                bid,
                blockNumber,
                txnHashes[i],
                revertingTxHashes,
                10,
                20,
                bidSignatures[i],
                sharedSecretKey
            );
            vm.stopPrank();
        }
        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        for (uint256 i = 0; i < commitments.length; ++i) {
            vm.expectEmit(true, false, false, true);
            emit CommitmentProcessed(commitments[i], false);
            oracle.processBuilderCommitmentForBlockNumber(
                commitments[i],
                blockNumber,
                provider,
                false,
                100
            );
        }
        vm.stopPrank();
        assertEq(providerRegistry.getProviderStake(provider), 250 ether);
        assertEq(bidderRegistry.getProviderAmount(provider), 4 * bid);
    }

    function constructAndStoreCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint256 bidderPk,
        uint256 signerPk,
        address provider,
        uint64 dispatchTimestamp
    ) public returns (bytes32 commitmentIndex) {
        bytes32 bidHash = getBidHash(txnHash, revertingTxHashes, bid, blockNumber);
        bytes memory bidSignature = getBidSignature(bidderPk, bidHash);
        
        bytes32 commitmentDigest = getCommitmentDigest(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            bidHash,
            bidSignature
        );
        
        bytes memory commitmentSignature = getCommitmentSignature(
            signerPk,
            commitmentDigest
        );
        

        bytes32 unopenedCommitmentIndex = storeUnopenedCommitment(
            provider,
            commitmentDigest,
            commitmentSignature,
            dispatchTimestamp
        );
        recordBlockData(provider, blockNumber);

        commitmentIndex = openCommitment(
            provider,
            unopenedCommitmentIndex,
            bid,
            blockNumber,
            txnHash,
            revertingTxHashes,
            bidSignature
        );
        return commitmentIndex;
    }

    function getBidHash(
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 bid,
        uint64 blockNumber
    ) public view returns (bytes32) {
        return
            preconfManager.getBidHash(
                txnHash,
                revertingTxHashes,
                bid,
                blockNumber,
                10,
                20
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
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 bid,
        uint64 blockNumber,
        bytes32 bidHash,
        bytes memory bidSignature
    ) public view returns (bytes32) {
        return
            preconfManager.getPreConfHash(
                txnHash,
                revertingTxHashes,
                bid,
                blockNumber,
                10,
                20,
                bidHash,
                bidSignature,
                sharedSecretKey
            );
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

    function recordBlockData(address provider, uint64 blockNumber) public {
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);

        blockTracker.recordL1Block(blockNumber, validBLSPubkey);
        vm.stopPrank();
    }

    function openCommitment(
        address provider,
        bytes32 unopenedCommitmentIndex,
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes memory bidSignature
    ) public returns (bytes32) {
        vm.startPrank(provider);
        bytes32 commitmentIndex = preconfManager.openCommitment(
            unopenedCommitmentIndex,
            bid,
            blockNumber,
            txnHash,
            revertingTxHashes,
            10,
            20,
            bidSignature,
            sharedSecretKey
        );
        vm.stopPrank();
        return commitmentIndex;
    }

    function constructAndStoreUnopenedCommitment(
        address committerAddress,
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        uint256 bidderPk,
        uint256 signerPk,
        uint64 dispatchTimestamp
    )
        public
        returns (
            bytes32 commitmentIndex,
            bytes memory bidSignature,
            bytes memory commitmentSignature
        )
    {
        bytes32 bidHash = preconfManager.getBidHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bidSignature = abi.encodePacked(r, s, v);

        bytes32 commitmentDigest = preconfManager.getPreConfHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            bidSignature,
            sharedSecretKey
        );

        (v, r, s) = vm.sign(signerPk, commitmentDigest);
        commitmentSignature = abi.encodePacked(r, s, v);
        vm.startPrank(committerAddress);
        bytes32 unopenedCommitmentIndex = preconfManager
            .storeUnopenedCommitment(
                commitmentDigest,
                commitmentSignature,
                dispatchTimestamp
            );
        vm.stopPrank();
        return (unopenedCommitmentIndex, bidSignature, commitmentSignature);
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
