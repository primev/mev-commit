// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

/**
 * @dev Gateway contract for standard bridge. 
 */
abstract contract Gateway is OwnableUpgradeable {   
    
    // @dev index for tracking transfer initiations.
    // Also total number of transfers initiated from this gateway.
    uint256 public transferInitiatedIdx;

    // @dev index for tracking transfer finalizations.
    // Also total number of transfers finalized on this gateway.
    uint256 public transferFinalizedIdx;

    // @dev Address of relayer account. 
    address public relayer;

    // @dev Flat fee (wei) paid to relayer on destination chain upon transfer finalization.
    // This must be greater than what relayer will pay per tx.
    uint256 public finalizationFee;

    // The counterparty's finalization fee (wei), included for UX purposes
    uint256 public counterpartyFee;

    function initiateTransfer(address _recipient, uint256 _amount
    ) external payable returns (uint256 returnIdx) {
        require(_amount >= counterpartyFee, "Amount must cover counterpartys finalization fee");
        _decrementMsgSender(_amount);
        ++transferInitiatedIdx;
        emit TransferInitiated(msg.sender, _recipient, _amount, transferInitiatedIdx);
        return transferInitiatedIdx;
    }
    // @dev where _decrementMsgSender is implemented by inheriting contract.
    function _decrementMsgSender(uint256 _amount) internal virtual;

    modifier onlyRelayer() {
        require(msg.sender == relayer, "Only relayer can call this function");
        _;
    }

    function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx
    ) external onlyRelayer {
        require(_amount >= finalizationFee, "Amount must cover finalization fee");
        require(_counterpartyIdx == transferFinalizedIdx, "Invalid counterparty index. Transfers must be relayed FIFO");
        uint256 amountAfterFee = _amount - finalizationFee;
        _fund(amountAfterFee, _recipient);
        _fund(finalizationFee, relayer);
        ++transferFinalizedIdx;
        emit TransferFinalized(_recipient, _amount, _counterpartyIdx);
    }
    // @dev where _fund is implemented by inheriting contract.
    function _fund(uint256 _amount, address _toFund) internal virtual;

    /**
     * @dev Emitted when a cross chain transfer is initiated.
     * @param sender Address initiating the transfer. Indexed for efficient filtering.
     * @param recipient Address receiving the tokens. Indexed for efficient filtering.
     * @param amount Ether being transferred in wei.
     * @param transferIdx Current index of this gateway.
     */
    event TransferInitiated(
        address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx);

    /**
     * @dev Emitted when a transfer is finalized.
     * @param recipient Address receiving the tokens. Indexed for efficient filtering.
     * @param amount Ether being transferred in wei.
     * @param counterpartyIdx Index of counterpary gateway when transfer was initiated.
     */
    event TransferFinalized(
        address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx);
}
