// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/**
 * @title IGasTankManager
 * @dev Interface for GasTankManager
 */
interface IGasTankManager {
    /// @dev Event to log successful update of the minimum deposit requirement.
    event MinimumDepositSet(uint256 indexed minDeposit);

    /// @dev Event to log successful gas tank top up.
    event GasTankToppedUp(address indexed smartAccount, address indexed caller, uint256 indexed amount);

    /// @notice Raised when a call expected from the contract itself is not.
    /// @param msgSender Address that invoked the call.
    /// @param thisAddress Address of the contract guard.
    error NotThisEOA(address msgSender, address thisAddress);

    /// @notice Raised when an invalid caller is detected.
    /// @param caller Address that invoked the call.
    error NotValidCaller(address caller);

    /// @notice Raised when an invalid amount is detected.
    error InvalidAmount();

    /// @notice Raised when forwarding funds to the provider reverts.
    /// @param rpcProvider Destination address for the transfer.
    /// @param transferAmount Amount of wei attempted.
    error GasTankTopUpFailed(address rpcProvider, uint256 transferAmount);

    /// @notice Raised when a gas tank balance already meets the minimum deposit requirement.
    /// @param currentBalance Gas tank balance reported by the provider.
    /// @param minDeposit Minimum deposit requirement.
    error GasTankBalanceIsSufficient(uint256 currentBalance, uint256 minDeposit);

    /// @notice Raised when the provider address is not set.
    /// @param provider Provider address.
    error ProviderNotSet(address provider);

    /// @notice Updates the minimum deposit requirement.
    /// @param minDeposit New minimum balance expressed in wei.
    /// @dev Only the owner can call this function.
    function setMinimumDeposit(uint256 minDeposit) external;

    /// @notice RPC requested top-up of the gas tank.
    /// @dev Always transfers the minimum deposit requirement.
    function topUpGasTank() external;

    /// @notice Allows anyone to fund the gas tank.
    function fundGasTank() external payable;

    /// @notice Initializes the gas tank with the minimum deposit.
    /// @dev Only the EOA can call this function.
    function initializeGasTank() external;
}
