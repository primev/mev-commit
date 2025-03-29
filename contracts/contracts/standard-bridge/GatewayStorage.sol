// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

contract GatewayStorage {
    /// @dev index for tracking transfer initiations.
    /// Also total number of transfers initiated from this gateway.
    uint256 public transferInitiatedIdx;

    /// @dev index for tracking transfer finalizations.
    /// Also total number of transfers finalized on this gateway.
    uint256 public transferFinalizedIdx;

    /// @dev Address of relayer account. 
    address public relayer;

    /// @dev The finalization fee (wei) of the counterparty gateway contract,
    /// paid to the relayer by the counterparty contract upon transfer finalization.
    /// @notice This value must on average, over time, be greater than what the relayer will pay per finalizeTransfer tx on the counterparty chain.
    /// @notice Consequently, the value may be mutated by the contract owner.
    uint256 public counterpartyFinalizationFee;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
