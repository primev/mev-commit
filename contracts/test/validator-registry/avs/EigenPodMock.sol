// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";

// Similar to eigenlayer core's EigenPodMock but their mocks don't use virtual functions.
contract EigenPodMock is IEigenPod, Test {
    mapping(bytes => ValidatorInfo) public validatorInfo;

    function setMockValidatorInfo(bytes memory validatorPubkey, ValidatorInfo memory info) external {
        validatorInfo[validatorPubkey] = info;
    }

    function validatorPubkeyToInfo(bytes memory validatorPubkey)
        external view virtual override returns(ValidatorInfo memory) {
        return validatorInfo[validatorPubkey];
    }

    function MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR() external view returns(uint64) {}

    function nonBeaconChainETHBalanceWei() external view returns(uint256) {}

    function withdrawableRestakedExecutionLayerGwei() external view returns(uint64) {}

    function initialize(address owner) external {}

    function stake(bytes calldata pubkey, bytes calldata signature, bytes32 depositDataRoot) external payable {}

    function withdrawRestakedBeaconChainETH(address recipient, uint256 amount) external {}

    function eigenPodManager() external view returns (IEigenPodManager) {}

    function podOwner() external view returns (address) {}

    function hasRestaked() external view returns (bool) {}

    function mostRecentWithdrawalTimestamp() external view returns (uint64) {}

    function validatorPubkeyHashToInfo(bytes32 validatorPubkeyHash) external view returns (ValidatorInfo memory) {}

    function provenWithdrawal(bytes32 validatorPubkeyHash, uint64 slot) external view returns (bool) {}

    function validatorStatus(bytes32 pubkeyHash) external view returns(VALIDATOR_STATUS) {}

    function verifyWithdrawalCredentials(
        uint64 oracleTimestamp,
        BeaconChainProofs.StateRootProof calldata stateRootProof,
        uint40[] calldata validatorIndices,
        bytes[] calldata withdrawalCredentialProofs,
        bytes32[][] calldata validatorFields
    ) external {}
    
    function verifyBalanceUpdates(
        uint64 oracleTimestamp,
        uint40[] calldata validatorIndices,
        BeaconChainProofs.StateRootProof calldata stateRootProof,
        bytes[] calldata validatorFieldsProofs,
        bytes32[][] calldata validatorFields
    ) external {}

    function verifyAndProcessWithdrawals(
        uint64 oracleTimestamp,
        BeaconChainProofs.StateRootProof calldata stateRootProof,
        BeaconChainProofs.WithdrawalProof[] calldata withdrawalProofs,
        bytes[] calldata validatorFieldsProofs,
        bytes32[][] calldata validatorFields,
        bytes32[][] calldata withdrawalFields
    ) external {}

    function activateRestaking() external {}

    function withdrawBeforeRestaking() external {}

    function withdrawNonBeaconChainETHBalanceWei(address recipient, uint256 amountToWithdraw) external {}

    function recoverTokens(IERC20[] memory tokenList, uint256[] memory amountsToWithdraw, address recipient) external {}

    function validatorStatus(bytes calldata pubkey) external view returns (VALIDATOR_STATUS){}
}
