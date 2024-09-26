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

    /// @dev Flat fee (wei) paid to relayer on destination chain upon transfer finalization.
    /// This must be greater than what relayer will pay per tx.
    uint256 public finalizationFee;

    /// @dev The counterparty's finalization fee (wei), included for UX purposes
    uint256 public counterpartyFee;

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#storage-gaps
    uint256[48] private __gap;
}
