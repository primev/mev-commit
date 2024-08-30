// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";

/**
 * @dev Gateway contract for standard bridge. 
 */
abstract contract Gateway is Ownable2StepUpgradeable, UUPSUpgradeable, PausableUpgradeable {   
    
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

    modifier onlyRelayer() {
        require(msg.sender == relayer, "sender is not relayer");
        _;
    }

    function initiateTransfer(address _recipient, uint256 _amount
    ) external payable whenNotPaused returns (uint256 returnIdx) {
        require(_amount >= counterpartyFee, "Amount too small");
        _decrementMsgSender(_amount);
        ++transferInitiatedIdx;
        emit TransferInitiated(msg.sender, _recipient, _amount, transferInitiatedIdx);
        return transferInitiatedIdx;
    }

    function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx
    ) external onlyRelayer whenNotPaused {
        require(_amount >= finalizationFee, "Amount too small");
        require(_counterpartyIdx == transferFinalizedIdx, "Invalid counterparty index");
        uint256 amountAfterFee = _amount - finalizationFee;
        _fund(amountAfterFee, _recipient);
        _fund(finalizationFee, relayer);
        ++transferFinalizedIdx;
        emit TransferFinalized(_recipient, _amount, _counterpartyIdx);
    }

    /// @dev Allows owner to pause the contract.
    function pause() external onlyOwner { _pause(); }

    /// @dev Allows owner to unpause the contract.
    function unpause() external onlyOwner { _unpause(); }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // @dev where _decrementMsgSender is implemented by inheriting contract.
    function _decrementMsgSender(uint256 _amount) internal virtual;

    // @dev where _fund is implemented by inheriting contract.
    function _fund(uint256 _amount, address _toFund) internal virtual;
}
