// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import "forge-std/Test.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import "eigenlayer-contracts/src/contracts/libraries/BeaconChainProofs.sol";

// Similar to eigenlayer core's EigenPodMock but their mocks don't use virtual functions.
contract EigenPodMock is IEigenPod, Test {
    mapping(bytes => ValidatorInfo) public validatorInfo;

    function setMockValidatorInfo(bytes memory validatorPubkey, ValidatorInfo memory info) external {
        validatorInfo[validatorPubkey] = info;
    }

    function validatorPubkeyToInfo(bytes calldata validatorPubkey)
        external view override returns(ValidatorInfo memory) {
        return validatorInfo[validatorPubkey];
    }

    function validatorPubkeyHashToInfo(bytes32 validatorPubkeyHash) external view override returns (ValidatorInfo memory) {}

    function initialize(address owner) external override {}

    function stake(bytes calldata pubkey, bytes calldata signature, bytes32 depositDataRoot) external payable override {}

    function withdrawRestakedBeaconChainETH(address recipient, uint256 amount) external override {}

    function startCheckpoint(bool revertIfNoBalance) external override {}

    function verifyCheckpointProofs(
        BeaconChainProofs.BalanceContainerProof calldata balanceContainerProof,
        BeaconChainProofs.BalanceProof[] calldata proofs
    ) external override {}

    function verifyWithdrawalCredentials(
        uint64 beaconTimestamp,
        BeaconChainProofs.StateRootProof calldata stateRootProof,
        uint40[] calldata validatorIndices,
        bytes[] calldata validatorFieldsProofs,
        bytes32[][] calldata validatorFields
    ) external override {}

    function verifyStaleBalance(
        uint64 beaconTimestamp,
        BeaconChainProofs.StateRootProof calldata stateRootProof,
        BeaconChainProofs.ValidatorProof calldata proof
    ) external override {}

    function recoverTokens(IERC20[] memory tokenList, uint256[] memory amountsToWithdraw, address recipient) external override {}

    function setProofSubmitter(address newProofSubmitter) external override {}

    function proofSubmitter() external view override returns (address) {}

    function withdrawableRestakedExecutionLayerGwei() external view override returns (uint64) {}

    function eigenPodManager() external view override returns (IEigenPodManager) {}

    function podOwner() external view override returns (address) {}

    function validatorStatus(bytes32 pubkeyHash) external view override returns(VALIDATOR_STATUS) {}

    function validatorStatus(bytes calldata validatorPubkey) external view override returns (VALIDATOR_STATUS){}

    function activeValidatorCount() external view override returns (uint256) {}

    function lastCheckpointTimestamp() external view override returns (uint64) {}

    function currentCheckpointTimestamp() external view override returns (uint64) {}

    function currentCheckpoint() external view override returns (Checkpoint memory) {}

    function checkpointBalanceExitedGwei(uint64) external view override returns (uint64) {}

    function getParentBlockRoot(uint64 timestamp) external view override returns (bytes32) {}
}
