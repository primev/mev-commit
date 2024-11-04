// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

contract GatewayStorage {
    /// @dev index for tracking transfer initiations.
    /// Also total number of transfers initiated from this gateway.
    uint256 public transferInitiatedIdx;

    /// @dev index for tracking transfer finalizations.
    /// Also total number of transfers finalized on this gateway.
    uint256 public transferFinalizedIdx;

    /// @dev Address of relayer account. 
    address public relayer;

    /// @dev The finalization fee (wei) paid to the relayer by this contract upon transfer finalization.
    /// This must be greater on average, over time, than what the relayer will pay per finalizeTransfer tx.
    /// @notice When setting this value, ensure the same value is set as the `counterpartyFee` in the counterparty contract.
    uint256 public finalizationFee;

    /// @dev The finalization fee (wei) of the counterparty gateway contract, included for UX purposes.
    /// @notice When setting this value, ensure the same value is set as the `finalizationFee` in the counterparty contract.
    uint256 public counterpartyFee;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
