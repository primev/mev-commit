// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";

// Eigenlayer core does not define their own mock for AVSDirectory.sol, hence we define our own.
contract AVSDirectoryMock is IAVSDirectory, Test {
    mapping(address => bool) public isOperatorRegistered;

    function registerOperator(address operator) external {
        isOperatorRegistered[operator] = true;
    }

    function deregisterOperator(address operator) external {
        isOperatorRegistered[operator] = false;
    }

    function isRegisteredOperator(address operator) external view returns (bool) {
        return isOperatorRegistered[operator];
    }

    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external override {
        require(operator != address(0), "Operator cannot be zero");
        require(operatorSignature.salt != bytes32(0), "Salt cannot be zero");
        isOperatorRegistered[operator] = true;
        emit OperatorAVSRegistrationStatusUpdated(operator, msg.sender, OperatorAVSRegistrationStatus.REGISTERED);
    }

    function deregisterOperatorFromAVS(address operator) external override {
        isOperatorRegistered[operator] = false;
        emit OperatorAVSRegistrationStatusUpdated(operator, msg.sender, OperatorAVSRegistrationStatus.UNREGISTERED);
    }

    function updateAVSMetadataURI(string calldata metadataURI) external override {
        emit AVSMetadataURIUpdated(msg.sender, metadataURI);
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
    ) external pure override returns (bytes32) {
        return keccak256(abi.encodePacked(operator, avs, salt, expiry));
    }

    function OPERATOR_AVS_REGISTRATION_TYPEHASH() external pure override returns (bytes32) {
        return keccak256("OperatorAVSRegistration(address operator,address avs,bytes32 salt,uint256 expiry)");
    }
}
