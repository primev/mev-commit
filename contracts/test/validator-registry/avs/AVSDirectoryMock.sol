// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

// Eigenlayer core does not define their own mock for AVSDirectory.sol, hence we define our own.
contract AVSDirectoryMock is IAVSDirectory, Test {
    mapping(address => bool) public isOperatorRegistered;
    address public avs_;

    bytes32 public constant OPERATOR_AVS_REGISTRATION_TYPEHASH =
        keccak256("OperatorAVSRegistration(address operator,address avs,bytes32 salt,uint256 expiry)");

    bytes32 public constant DOMAIN_TYPEHASH =
        keccak256("EIP712Domain(string name,uint256 chainId,address verifyingContract)");

    function registerOperator(address operator) external {
        isOperatorRegistered[operator] = true;
    }

    function deregisterOperator(address operator) external {
        isOperatorRegistered[operator] = false;
    }

    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external override {
        require(operator != address(0), "Operator required");
        require(keccak256(operatorSignature.signature) != keccak256(bytes("")), "Signature required");
        require(operatorSignature.salt != bytes32(0), "Salt required");
        require(operatorSignature.expiry != 0, "Expiry required");

        bytes32 operatorRegistrationDigestHash = calculateOperatorAVSRegistrationDigestHash({
            operator: operator,
            avs: avs_,
            salt: operatorSignature.salt,
            expiry: operatorSignature.expiry
        });

        // solhint-disable-next-line reason-string
        require(
            ECDSA.recover(operatorRegistrationDigestHash, operatorSignature.signature) == operator,
                "EIP1271SignatureUtils.checkSignature_EIP1271: signature not from signer"
        );

        isOperatorRegistered[operator] = true;
        emit OperatorAVSRegistrationStatusUpdated(operator, msg.sender, OperatorAVSRegistrationStatus.REGISTERED);
    }

    function setAVS(address _avs) external {
        avs_ = _avs;
    }

    function deregisterOperatorFromAVS(address operator) external override {
        isOperatorRegistered[operator] = false;
        emit OperatorAVSRegistrationStatusUpdated(operator, msg.sender, OperatorAVSRegistrationStatus.UNREGISTERED);
    }

    function updateAVSMetadataURI(string calldata metadataURI) external override {
        emit AVSMetadataURIUpdated(msg.sender, metadataURI);
    }

    function isRegisteredOperator(address operator) external view returns (bool) {
        return isOperatorRegistered[operator];
    }

    function operatorSaltIsSpent(address operator, bytes32 salt) external pure override returns (bool) {
        require(operator != address(0), "Operator cannot be zero");
        require(salt != bytes32(0), "Salt cannot be zero");
        return false;
    }

    function calculateOperatorAVSRegistrationDigestHash(
        address operator,
        address avs,
        bytes32 salt,
        uint256 expiry
    ) public view returns (bytes32) {
        bytes32 structHash = keccak256(
            abi.encode(OPERATOR_AVS_REGISTRATION_TYPEHASH, operator, avs, salt, expiry)
        );
        bytes32 digestHash = keccak256(
            abi.encodePacked("\x19\x01", domainSeparator(), structHash)
        );
        return digestHash;
    }

    function domainSeparator() public view returns (bytes32) {
        return keccak256(abi.encode(DOMAIN_TYPEHASH, keccak256(bytes("EigenLayer")), block.chainid, address(this)));
    }
}
