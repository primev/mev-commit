// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../contracts/Oracle.sol";
import "../contracts/PreConfirmations.sol";
import "../contracts/interfaces/IPreConfirmations.sol";
import "../contracts/ProviderRegistry.sol";
import "../contracts/BidderRegistry.sol";
import "../contracts/BlockTracker.sol";

import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

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
    BlockTracker internal blockTracker;
    TestCommitment internal _testCommitmentAliceBob;
    bytes internal sharedSecretKey;

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
    }

    // Events to match against
    event BlockDataRequested(uint256 blockNumber);
    event BlockDataReceived(string[] txnList, uint256 blockNumber, string blockBuilderName);
    event CommitmentProcessed(bytes32 commitmentHash, bool isSlash);
    event FundsRetrieved(bytes32 indexed commitmentDigest, uint256 window, uint256 amount);

    function setUp() public {
        testNumber = 2;
        testNumber2 = 2;
        sharedSecretKey = abi.encodePacked(keccak256("0xsecret"));
        _testCommitmentAliceBob = TestCommitment(
            2,
            2,
            "0xkartik",
            10,
            20,
            0xa0327970258c49b922969af74d60299a648c50f69a2d98d6ab43f32f64ac2100,
            0x54c118e537dd7cf63b5388a5fc8322f0286a978265d0338b108a8ca9d155dccc,
            hex"876c1216c232828be9fabb14981c8788cebdf6ed66e563c4a2ccc82a577d052543207aeeb158a32d8977736797ae250c63ef69a82cd85b727da21e20d030fb311b",
            hex"ec0f11f77a9e96bb9c2345f031a5d12dca8d01de8a2e957cf635be14802f9ad01c6183688f0c2672639e90cc2dce0662d9bea3337306ca7d4b56dd80326aaa231b"
        );

        feePercent = 10;
        minStake = 1e18 wei;
        feeRecipient = vm.addr(9);

        address proxy = Upgrades.deployUUPSProxy(
            "ProviderRegistry.sol",
            abi.encodeCall(ProviderRegistry.initialize, (minStake, feeRecipient, feePercent, address(this))) 
        );
        providerRegistry = ProviderRegistry(payable(proxy));

        address ownerInstance = 0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3;

        address proxy2 = Upgrades.deployUUPSProxy(
            "BlockTracker.sol",
            abi.encodeCall(BlockTracker.initialize, (ownerInstance))
        );
        blockTracker = BlockTracker(payable(proxy2));

        address proxy3 = Upgrades.deployUUPSProxy(
            "BidderRegistry.sol",
            abi.encodeCall(BidderRegistry.initialize, 
            (minStake, 
            feeRecipient, 
            feePercent, 
            address(this), 
            address(blockTracker)))
        );    
        bidderRegistry = BidderRegistry(payable(proxy3));

        address proxy4 = Upgrades.deployUUPSProxy(
            "PreConfirmations.sol",
            abi.encodeCall(PreConfCommitmentStore.initialize, 
            (address(providerRegistry), 
            address(bidderRegistry), 
            address(blockTracker), 
            feeRecipient, // Oracle
            address(this))) // Owner
        );
        preConfCommitmentStore = PreConfCommitmentStore(payable(proxy4));
        
        vm.deal(ownerInstance, 5 ether);
        vm.startPrank(ownerInstance);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForSpecificWindow{value: 2 ether}(window+1);
        
        address proxy5 = Upgrades.deployUUPSProxy(
            "Oracle.sol",
            abi.encodeCall(Oracle.initialize, 
            (address(preConfCommitmentStore), 
            address(blockTracker), 
            ownerInstance))
        );
        oracle = Oracle(payable(proxy5)); 

        vm.stopPrank();

        preConfCommitmentStore.updateOracle(address(oracle));
        bidderRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));
        providerRegistry.setPreconfirmationsContract(address(preConfCommitmentStore));

    }

    function test_process_commitment_payment_payout() public {
        string memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        uint256 blocksPerWindow = blockTracker.blocksPerWindow();
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        uint64 bid = 2;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForSpecificWindow{value: 250 ether}(window+1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(bid, blockNumber, txn, 10, 20, bidderPk, providerPk, provider);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        oracle.processBuilderCommitmentForBlockNumber(index, blockNumber, provider, false, 50);
        vm.stopPrank();
        assertEq(bidderRegistry.getProviderAmount(provider), bid*(50)/100);

    }


    function test_process_commitment_slash() public {
        string memory txn = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        uint64 blockNumber = 200;
        uint64 bid = 200;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForSpecificWindow{value: 250 ether}(window+1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index = constructAndStoreCommitment(bid, blockNumber, txn, 10, 20, bidderPk, providerPk, provider);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index, true);
        oracle.processBuilderCommitmentForBlockNumber(index, blockNumber, provider, true,50);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider) + ((bid * 50)/100), 250 ether);
    }


    function test_process_commitment_slash_and_reward() public {
        string memory txn1 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        string memory txn2 = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        uint256 blocksPerWindow = blockTracker.blocksPerWindow();
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        uint64 bid = 100;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        uint256 residualAfterDecay = 50;

        vm.deal(bidder, 200000 ether);
        vm.startPrank(bidder);
        uint256 window = blockTracker.getCurrentWindow();
        bidderRegistry.depositForSpecificWindow{value: 250 ether}(window+1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(bid, blockNumber, txn1, 10, 20, bidderPk, providerPk, provider);
        bytes32 index2 = constructAndStoreCommitment(bid, blockNumber, txn2, 10, 20, bidderPk, providerPk, provider);

        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(index1, blockNumber, provider, true,100);

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, false);
        oracle.processBuilderCommitmentForBlockNumber(index2, blockNumber, provider, false,50);
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
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = blockTracker.getWindowFromBlockNumber(blockNumber);
        vm.startPrank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 250 ether}(window);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32 index1 = constructAndStoreCommitment(bid, blockNumber, txn1, 10, 20, bidderPk, providerPk, provider);
        bytes32 index2 = constructAndStoreCommitment(bid, blockNumber, txn2, 10, 20, bidderPk, providerPk, provider);
        bytes32 index3 = constructAndStoreCommitment(bid, blockNumber, txn3, 10, 20, bidderPk, providerPk, provider);
        bytes32 index4 = constructAndStoreCommitment(bid, blockNumber, txn4, 10, 20, bidderPk, providerPk, provider);


        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));

        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index1, true);
        oracle.processBuilderCommitmentForBlockNumber(index1, blockNumber, provider, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index2, true);
        oracle.processBuilderCommitmentForBlockNumber(index2, blockNumber, provider, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index3, true);
        oracle.processBuilderCommitmentForBlockNumber(index3, blockNumber, provider, true,100);
        vm.expectEmit(true, false, false, true);
        emit CommitmentProcessed(index4, true);
        oracle.processBuilderCommitmentForBlockNumber(index4, blockNumber, provider, true,100);
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider), 250 ether - bid*4);
        assertEq(bidderRegistry.getProviderAmount(provider), 0);
    }

    function test_process_commitment_reward_multiple() public {
        string[] memory txnHashes = new string[](4);
        txnHashes[0] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d08";
        txnHashes[1] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d09";
        txnHashes[2] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d10";
        txnHashes[3] = "0x6d9c53ad81249775f8c082b11ac293b2e19194ff791bd1c4fd37683310e90d11";
        uint256 blocksPerWindow = blockTracker.blocksPerWindow();
        uint64 blockNumber = uint64(blocksPerWindow + 2);
        uint64 bid = 5;
        (address bidder, uint256 bidderPk) = makeAddrAndKey("alice");
        (address provider, uint256 providerPk) = makeAddrAndKey("kartik");

        vm.deal(bidder, 200000 ether);
        uint256 window = blockTracker.getCurrentWindow();
        vm.startPrank(bidder);
        bidderRegistry.depositForSpecificWindow{value: 250 ether}(window+1);
        vm.stopPrank();

        vm.deal(provider, 200000 ether);
        vm.startPrank(provider);
        providerRegistry.registerAndStake{value: 250 ether}();
        vm.stopPrank();

        bytes32[] memory commitments = new bytes32[](4);
        bytes[] memory bidSignatures = new bytes[](4);
        bytes[] memory commitmentSignatures = new bytes[](4);
        for (uint i = 0; i < commitments.length; i++) {
            (commitments[i], bidSignatures[i], commitmentSignatures[i]) = constructAndStoreEncryptedCommitment(
                bid,
                blockNumber,
                txnHashes[i],
                10,
                20,
                bidderPk,
                providerPk
            );
        }

        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        vm.stopPrank();

        for (uint i = 0; i < commitments.length; i++) {
            vm.startPrank(provider);
            preConfCommitmentStore.openCommitment(
                commitments[i],
                bid,
                blockNumber,
                txnHashes[i],
                10,
                20,
                bidSignatures[i],
                commitmentSignatures[i],
                sharedSecretKey
            );
            vm.stopPrank();
        }
       
        vm.startPrank(address(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3));
        for (uint i = 0; i < commitments.length; i++) {
            vm.expectEmit(true, false, false, true);
            emit CommitmentProcessed(commitments[i], false);
            oracle.processBuilderCommitmentForBlockNumber(commitments[i], blockNumber, provider, false,100);
        }
        vm.stopPrank();
        assertEq(providerRegistry.checkStake(provider), 250 ether);
        assertEq(bidderRegistry.getProviderAmount(provider), 4*bid);
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
        address provider
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
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        (v,r,s) = vm.sign(signerPk, commitmentHash);
        bytes memory commitmentSignature = abi.encodePacked(r, s, v);

        bytes32 encryptedCommitmentIndex = preConfCommitmentStore.storeEncryptedCommitment(
            commitmentHash,
            commitmentSignature
        );
        vm.startPrank(0x6d503Fd50142C7C469C7c6B64794B55bfa6883f3);
        blockTracker.addBuilderAddress("test", provider);
        blockTracker.recordL1Block(blockNumber, "test");
        vm.stopPrank();
        vm.startPrank(provider);
        commitmentIndex = preConfCommitmentStore.openCommitment(
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
        vm.stopPrank();

        return commitmentIndex;
    }
    
    function constructAndStoreEncryptedCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        uint64 decayStartTimestamp,
        uint64 decayEndTimestamp,
        uint256 bidderPk,
        uint256 signerPk
    ) public returns (bytes32 commitmentIndex, bytes memory bidSignature, bytes memory commitmentSignature) {
        bytes32 bidHash = preConfCommitmentStore.getBidHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimestamp,
            decayEndTimestamp
        );


        (uint8 v,bytes32 r, bytes32 s) = vm.sign(bidderPk, bidHash);
        bidSignature = abi.encodePacked(r, s, v);

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

        (v,r,s) = vm.sign(signerPk, commitmentHash);
        commitmentSignature = abi.encodePacked(r, s, v);

        bytes32 encryptedCommitmentIndex = preConfCommitmentStore.storeEncryptedCommitment(
            commitmentHash,
            commitmentSignature
        );

        return (encryptedCommitmentIndex, bidSignature, commitmentSignature);
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
