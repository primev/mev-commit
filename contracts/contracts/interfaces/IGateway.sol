// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/**
 * @title IGateway
 * @dev Interface for the Gateway contract in the standard bridge.
 */
interface IGateway {
    /**
     * @dev Emitted when a cross-chain transfer is initiated.
     * @param sender Address initiating the transfer.
     * @param recipient Address receiving the tokens.
     * @param amount Ether being transferred in wei.
     * @param transferIdx Current index of this gateway.
     * @param counterpartyFinalizationFee The finalization fee (wei) paid to the relayer by the counterparty contract.
     */
    event TransferInitiated(
        address indexed sender,
        address indexed recipient,
        uint256 amount,
        uint256 indexed transferIdx,
        uint256 counterpartyFinalizationFee
    );

    /**
     * @dev Emitted when a transfer is finalized.
     * @param recipient Address receiving the tokens.
     * @param amount Ether being transferred in wei.
     * @param counterpartyIdx Index of counterparty gateway when transfer was initiated.
     */
    event TransferFinalized(
        address indexed recipient,
        uint256 amount,
        uint256 indexed counterpartyIdx
    );

    event CounterpartyFinalizationFeeSet(uint256 counterpartyFinalizationFee);
    event RelayerSet(address indexed relayer);

    error SenderNotRelayer(address sender, address relayer);
    error AmountTooSmall(uint256 amount, uint256 counterpartyFinalizationFee);
    error InvalidCounterpartyIndex(uint256 counterpartyIdx, uint256 transferFinalizedIdx);
    error CounterpartyFinalizationFeeTooSmall(uint256 _counterpartyFinalizationFee);
    error RelayerCannotBeZeroAddress();

    /**
     * @notice Initiates a cross-chain transfer.
     * @param _recipient Address to receive the tokens.
     * @param _amount Amount of Ether to transfer in wei.
     * @return returnIdx The index of the initiated transfer.
     */
    function initiateTransfer(address _recipient, uint256 _amount)
        external
        payable
        returns (uint256 returnIdx);

    /**
     * @notice Finalizes a cross-chain transfer.
     * @param _recipient Address to receive the tokens.
     * @param _amount Amount of Ether to transfer in wei.
     * @param _counterpartyIdx Index of the counterparty transfer.
     * @param _finalizationFee The finalization fee (wei) paid to the relayer by the this contract.
     */
    function finalizeTransfer(
        address _recipient,
        uint256 _amount,
        uint256 _counterpartyIdx,
        uint256 _finalizationFee
    ) external;

    /**
     * @notice Pauses the contract, preventing certain functions from being called.
     */
    function pause() external;

    /**
     * @notice Unpauses the contract, allowing previously paused functions to be called.
     */
    function unpause() external;
}
