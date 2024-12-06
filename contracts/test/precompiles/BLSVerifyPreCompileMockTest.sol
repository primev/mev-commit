// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";

contract BLSVerifyPreCompileMock is Test {
    address constant BLS_VERIFY = address(0xf0);

    function setUp() public {
        // Deploy the mock contract to the precompile address
        bytes memory code = type(MockBLSVerify).creationCode;
        vm.etch(BLS_VERIFY, code);
    }

    function testVerifySignature() public {
        // Store the values as constants or immutable variables since they're fixed test data
        bytes memory pubKeyData = hex"b67a5148a03229926e34b190af81a82a81c4df66831c98c03a139778418dd09a3b542ced0022620d19f35781ece6dc36";
        bytes32 messageData = keccak256(abi.encodePacked("test message"));
        bytes memory signatureData = hex"bbbbbbbbb1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2";

        // Call verifySignature with the abi.encodePacked data to convert to calldata
        bool result = this.verifySignature(
            pubKeyData,
            messageData,
            signatureData
        );
        assertTrue(result, "Signature verification should succeed.");
    }
    /**
     * @dev Verifies a BLS signature using the precompile
     * @param pubKey The public key (48 bytes G1 point)
     * @param message The message hash (32 bytes)
     * @param signature The signature (96 bytes G2 point)
     * @return success True if verification succeeded
     */
    function verifySignature(
        bytes calldata pubKey,
        bytes32 message,
        bytes calldata signature
    ) public view returns (bool) {
        // Input validation
        require(pubKey.length == 48, "Public key must be 48 bytes");
        require(signature.length == 96, "Signature must be 96 bytes");

        // Concatenate inputs in required format:
        // [pubkey (48 bytes) | message (32 bytes) | signature (96 bytes)]
        bytes memory input = bytes.concat(
            pubKey,
            message,
            signature
        );

        // Call precompile
        (bool success, bytes memory result) = BLS_VERIFY.staticcall(input);
        
        // Check if call was successful
        if (!success) {
            return false;
        }

        // If we got a result back and it's not empty, verification succeeded
        return result.length > 0;
    }
}

contract MockBLSVerify {
    function VERIFY(bytes calldata) external pure returns (bool) {
        return true;
    }
}
