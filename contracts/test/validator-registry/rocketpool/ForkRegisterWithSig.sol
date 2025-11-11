// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import "forge-std/console.sol";

import {RocketMinipoolRegistry} from "../../../contracts/validator-registry/rocketpool/RocketMinipoolRegistry.sol";

contract MinipoolShim {
    address public immutable nodeAddr;
    constructor(address _node) {
        nodeAddr = _node;
    }

    function getNodeAddress() external view returns (address) {
        return nodeAddr;
    }

    // Return Staking (2) to behave as active if caller checks status
    function getStatus() external pure returns (uint8) {
        return 2;
    }
}

// contract ExecutorProxy {
//     function execRegister(address registry, bytes[] calldata pks, bytes calldata sig, uint256 deadline) external {
//         RocketMinipoolRegistry(registry).registerValidatorsWithSig(pks, sig, deadline);
//     }
// }

contract ForkRegisterWithSig is Test {
    // Replace these constants with your values (or set as environment variables)
    address constant HOODI_REGISTRY  = 0x0694bFD12dcBC9165e91A8C2843bc1fe9d3f3DD0;
    address constant TARGET_MINIPOOL = 0x862e2909Dafb8dc16417B48dBBe688714737bBf5;
    address constant NODE_ADDR       = 0x1623fE21185c92BB43bD83741E226288B516134a;
    // The 48-byte validator pubkey (as bytes). Provide exact bytes.
    bytes constant VAL_PUBKEY = hex"b0d7f0d83e3bbaefacd7b05596f0e8f4c520d96f494b46180b9bcbb139a866bb89a63c92c86365658b1411d22ada016c";

    // EIP-712 typehash (must match on-chain contract)
    bytes32 internal constant REGISTER_TYPEHASH =
        keccak256("Register(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");

    string internal constant NAME = "RocketMinipoolRegistry";
    string internal constant VERSION = "1";
    bytes32 internal constant EIP712_DOMAIN_TYPEHASH =
        keccak256("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)");

    function test_fork_register_with_sig() public {
        // This test is intended to run on a mainnet (hoodi) fork.
        // Make sure you run: forge test -vvvv --fork-url <ANVIL_FORK_URL> -m test_fork_register_with_sig

        // 1) Etch the minipool shim at the target minipool address so getNodeAddress() returns NODE_ADDR
        MinipoolShim shim = new MinipoolShim(NODE_ADDR);
        bytes memory shimCode = address(shim).code;
        vm.etch(TARGET_MINIPOOL, shimCode);
        console.log("Etched shim at minipool:", TARGET_MINIPOOL);

        // 2) Point to deployed registry on the fork
        RocketMinipoolRegistry reg = RocketMinipoolRegistry(payable(HOODI_REGISTRY));

        // 3) Compute pubkeysHash exactly as the contract expects (concatenate 48-byte pubkeys)
        bytes[] memory pks = new bytes[](1);
        pks[0] = VAL_PUBKEY;
        bytes32 pkHash = reg.pubkeysHash(pks);
        console.logBytes32(pkHash);

        // 4) Determine executor (the address that will call the tx). We'll use address(this) by default.
        address executor = address(this);

        // Optionally override executor with an env var "EXECUTOR"
        string memory execMaybe = vm.envOr("EXECUTOR", string(""));
        if (bytes(execMaybe).length != 0) {
            executor = vm.envAddress("EXECUTOR");
        }
        console.log("Executor:", executor);
        

        // 5) Get nonce for the node from the on-chain registry
        uint256 nonce = reg.nonces(NODE_ADDR);
        console.logUint(nonce);

        // 6) Prepare deadline
        uint256 deadline = block.timestamp + 1 days;

        // 7) Build struct hash and domain separator, then the digest
        bytes32 structHash = keccak256(abi.encode(
            REGISTER_TYPEHASH,
            pkHash,
            executor,
            nonce,
            deadline
        ));

        bytes32 domainSeparator = keccak256(abi.encode(
            EIP712_DOMAIN_TYPEHASH,
            keccak256(bytes(NAME)),
            keccak256(bytes(VERSION)),
            block.chainid,
            address(reg)
        ));

        bytes32 digest = keccak256(abi.encodePacked("\x19\x01", domainSeparator, structHash));
        console.logBytes32(digest);

        // 8) Try to sign the digest using an env-provided private key (NODE_PK). If present, we'll sign and call.
        string memory skHex = vm.envOr("NODE_PK", string(""));
        if (bytes(skHex).length == 0) {
            // Not provided: print instructions for offline signing.
            console.log("No NODE_PK env var provided. To complete flow:");
            console.log(" - Sign the digest printed above with the private key controlling NODE_ADDR.");
            console.log(" - Provide the r||s||v as bytes in the call below, or set NODE_PK to the hex private key.");
            return;
        }

        // If NODE_PK provided as hex (with or without 0x), read as uint
        uint256 nodeSk = vm.envUint("NODE_PK");
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(nodeSk, digest);
        bytes memory sig = abi.encodePacked(r, s, v);
        console.logBytes(sig);

        // 9) Call from the intended executor; ensure the account exists locally to avoid fork lookups
        reg.registerValidatorsWithSig(pks, sig, deadline);

        console.log("Called registerValidatorsWithSig on registry:", HOODI_REGISTRY);
    }
}
