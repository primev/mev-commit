// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../contracts/Oracle.sol";
import "../contracts/PreConfirmations.sol";
import "../contracts/interfaces/IPreConfirmations.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/BidderRegistry.sol";

contract OracleTest is Test {
    address internal owner;
    using ECDSA for bytes32;
    Oracle internal oracle;
    PreConfCommitmentStore internal preConfCommitmentStore;
    uint16 internal feePercent;
    uint256 internal minStake;
    address internal feeRecipient;
    ProviderRegistry internal providerRegistry;
    uint256 testNumber;
    uint64 testNumber2;
    BidderRegistry internal bidderRegistry;
    TestCommitment internal _testCommitmentAliceBob;
    uint64 internal dispatchTimestamp;

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
    event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName);
    event CommitmentProcessed(bytes32 commitmentHash, bool isSlash);
    event FundsRetrieved(bytes32 indexed commitmentDigest, uint256 amount);

    function setUp() public {
        testNumber = 2;
        testNumber2 = 2;

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

        providerRegistry = new ProviderRegistry(
            minStake,
            feeRecipient,
            feePercent,
            address(this)
        );
        bidderRegistry = new BidderRegistry(minStake, feeRecipient, feePercent, address(this));
        preConfCommitmentStore = new PreConfCommitmentStore(
            address(providerRegistry), // Provider Registry
            address(bidderRegistry), // User Registry
            feeRecipient, // Oracle
            address(this),
            500
        );

        address ownerInstance = 0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3;
        vm.deal(ownerInstance, 5 ether);
        vm.startPrank(ownerInstance);
        bidderRegistry.prepay{value: 2 ether}();
        
        oracle = new Oracle(address(preConfCommitmentStore), 2, ownerInstance);
        oracle.addBuilderAddress("mev builder", ownerInstance);
        vm.stopPrank();

        preConfCommitmentStore.updateOracle(address(oracle));
        bidderRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));
        providerRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));

        // We set the system time to 1010 and dispatchTimestamps for testing to 1000
        dispatchTimestamp = 1000;
        vm.warp(1010);

    }

    function test_MultipleBlockBuildersRegistred() public {
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        (address builder1,) = makeAddrAndKey("k builder");
        (address builder2,) = makeAddrAndKey("primev builder");
        (address builder3,) = makeAddrAndKey("titan builder");
        (address builder4,) = makeAddrAndKey("zk builder");


        oracle.addBuilderAddress("k builder", builder1);
        oracle.addBuilderAddress("primev builder", builder2);
        oracle.addBuilderAddress("titan builder", builder3);
        oracle.addBuilderAddress("zk builder", builder4);

        assertEq(oracle.blockBuilderNameToAddress("k builder"), builder1);
        assertEq(oracle.blockBuilderNameToAddress("primev builder"), builder2);
        assertEq(oracle.blockBuilderNameToAddress("titan builder"), builder3);
        assertEq(oracle.blockBuilderNameToAddress("zk builder"), builder4);
    }

    function test_builderUnidentified() public {
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        // Unregistered Builder
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("bob");

        (address builder3,) = makeAddrAndKey("titan builder");
        (address builder4,) = makeAddrAndKey("zk builder");

        oracle.addBuilderAddress("titan builder", builder3);
        oracle.addBuilderAddress("zk builder", builder4);

        assertEq(oracle.blockBuilderNameToAddress("titan builder"), builder3);
        assertEq(oracle.blockBuilderNameToAddress("zk builder"), builder4);
        vm.stopPrank();

        vm.deal(bidder, 1000 ether);
        vm.deal(provider, 1000 ether);

        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 commitmentIndex = constructAndStoreCommitment(
            _testCommitmentAliceBob.bid,
            _testCommitmentAliceBob.blockNumber,
            _testCommitmentAliceBob.txnHash,
            _testCommitmentAliceBob.decayStartTimestamp,
            _testCommitmentAliceBob.decayEndTimestamp,
            bidderPk,
            providerPk,
            _testCommitmentAliceBob.dispatchTimestamp
        );

        string[] memory txnList = new string[](1);
        txnList[0] = string(abi.encodePacked(keccak256("0xkartik")));
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        oracle.processBuilderCommitmentForBlockNumber(commitmentIndex, _testCommitmentAliceBob.blockNumber, "k builder", false, 50);
        vm.stopPrank();
        assertEq(bidderRegistry.getProviderAmount(provider), 0);
        assertEq(providerRegistry.checkStake(provider), 250 ether);
    }

    function test_process_commitment_payment_payout() public {
        string memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        uint64 blockNumber = 200;
        uint64 bid = 2;
        string memory blockBuilderName = "kartik builder";
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(bid, blockNumber, txn, 10, 20, bidderPk, providerPk, dispatchTimestamp);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        oracle.addBuilderAddress(blockBuilderName, provider);

        oracle.processBuilderCommitmentForBlockNumber(index, blockNumber, blockBuilderName, false, 50);
        vm.stopPrank();
        assertEq(bidderRegistry.getProviderAmount(provider), bid*(50)/100);

    }


    function test_process_commitment_slash() public {
        string memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        uint64 blockNumber = 200;
        uint64 bid = 200;
        string memory blockBuilderName = "kartik builder";
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(bid, blockNumber, txn, 10, 20, bidderPk, providerPk, dispatchTimestamp);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        oracle.addBuilderAddress(blockBuilderName, provider);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index, true);
        oracle.processBuilderCommitmentForBlockNumber(index, blockNumber, blockBuilderName, true,50);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider) + ((bid * 50)/100), 250 ether);
    }


    function test_process_commitment_slash_and_reward() public {
        string memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        uint64 blockNumber = 201;
        uint64 bid = 100;
        string memory blockBuilderName = "kartik builder";
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        uint256 residualAfterDecay = 50;

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(bid, blockNumber, txn1, 10, 20, bidderPk, providerPk,dispatchTimestamp);
        bytes32 index2 = constructAndStoreCommitment(bid, blockNumber, txn2, 10, 20, bidderPk, providerPk,dispatchTimestamp);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        oracle.addBuilderAddress(blockBuilderName, provider);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(index1, blockNumber, blockBuilderName, true,100);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, false);
        oracle.processBuilderCommitmentForBlockNumber(index2, blockNumber, blockBuilderName, false,50);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider), 250 ether - bid);
        assertEq(bidderRegistry.getProviderAmount(provider), (bid * (100 - feePercent) /100) * residualAfterDecay /100 );
    }


    function test_process_commitment_slash_multiple() public {
        string memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        string memory txn3 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d10";
        string memory txn4 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d11";
        uint64 blockNumber = 201;
        uint64 bid = 5;
        string memory blockBuilderName = "kartik builder";
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(bid, blockNumber, txn1, 10, 20, bidderPk, providerPk,dispatchTimestamp);
        bytes32 index2 = constructAndStoreCommitment(bid, blockNumber, txn2, 10, 20, bidderPk, providerPk,dispatchTimestamp);
        bytes32 index3 = constructAndStoreCommitment(bid, blockNumber, txn3, 10, 20, bidderPk, providerPk,dispatchTimestamp);
        bytes32 index4 = constructAndStoreCommitment(bid, blockNumber, txn4, 10, 20, bidderPk, providerPk,dispatchTimestamp);


        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        oracle.addBuilderAddress(blockBuilderName, provider);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(index1, blockNumber, blockBuilderName, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, true);
        oracle.processBuilderCommitmentForBlockNumber(index2, blockNumber, blockBuilderName, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index3, true);
        oracle.processBuilderCommitmentForBlockNumber(index3, blockNumber, blockBuilderName, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index4, true);
        oracle.processBuilderCommitmentForBlockNumber(index4, blockNumber, blockBuilderName, true,100);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider), 250 ether - bid*4);
        assertEq(bidderRegistry.getProviderAmount(provider), 0);
    }

    function test_process_commitment_reward_multiple() public {
        string memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        string memory txn3 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d10";
        string memory txn4 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d11";
        uint64 blockNumber = 201;
        uint64 bid = 5;
        string memory blockBuilderName = "kartik builder";
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(bid, blockNumber, txn1, 10, 20, bidderPk, providerPk ,dispatchTimestamp);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 250 ether - bid);
        bytes32 index2 = constructAndStoreCommitment(bid, blockNumber, txn2, 10, 20, bidderPk, providerPk ,dispatchTimestamp);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 250 ether - 2*bid);
        bytes32 index3 = constructAndStoreCommitment(bid, blockNumber, txn3, 10, 20, bidderPk, providerPk, dispatchTimestamp);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 250 ether - 3*bid);
        bytes32 index4 = constructAndStoreCommitment(bid, blockNumber, txn4, 10, 20, bidderPk, providerPk, dispatchTimestamp);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 250 ether - 4*bid);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        oracle.addBuilderAddress(blockBuilderName, provider);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, false);
        oracle.processBuilderCommitmentForBlockNumber(index1, blockNumber, blockBuilderName, false,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, false);
        oracle.processBuilderCommitmentForBlockNumber(index2, blockNumber, blockBuilderName, false,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index3, false);
        oracle.processBuilderCommitmentForBlockNumber(index3, blockNumber, blockBuilderName, false,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index4, false);
        oracle.processBuilderCommitmentForBlockNumber(index4, blockNumber, blockBuilderName, false,100);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider), 250 ether);
        assertEq(bidderRegistry.getProviderAmount(provider), 4*bid);
    }


    function test_process_commitment_and_return() public {
        string memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        uint64 blockNumber = 200;
        uint64 bid = 2;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        bidderRegistry.prepay{value: 250 ether }();
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(bid, blockNumber, txn, 10, 20, bidderPk, providerPk, dispatchTimestamp);
        PreConfCommitmentStore.PreConfCommitment memory commitment = preConfCommitmentStore.getCommitment(index);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        bytes32[] memory commitments = new bytes32[](1);
        commitments[0] = commitment.commitmentHash;

        vm.expectEmit(true, false, false, true);
        emit FundsRetrieved(commitment.commitmentHash, bid);
        oracle.unlockFunds(commitments);
        
        
        assertEq(providerRegistry.checkStake(provider) , 250 ether);
        assertEq(bidderRegistry.bidderPrepaidBalances(bidder), 250 ether);
    }


    /**
    constructAndStoreCommitment is a helper function to construct and store a commitment
     */
    function constructAndStoreCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        uint256 bidderPk,
        uint256 signerPk,
        uint64 dispatchTimestamp
    ) public returns (bytes32 commitmentIndex) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );


        (uint8 v,bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bytes memory bidSignature = abi.encodePacked(r, s, v);

        bytes32 commitmentHash = preConfCommitmentStore.getPreConfHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp,
            bidHash,
            _bytesToHexString(bidSignature)
        );

        (v,r,s) = vm.sign(signerPk, commitmentHash);
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        commitmentIndex = preConfCommitmentStore.storeCommitment(
            bid,
            blockNumber,
            txnHash,
            decayStartTimestamp,
            decayEndTimestamp,
            bidSignature,
            commitmentSignature,
            dispatchTimestamp
        );

        return commitmentIndex;
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
