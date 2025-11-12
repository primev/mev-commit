// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Errors} from "../utils/Errors.sol";

/// @title GasTankManager
/// @notice Coordinates on-demand ETH Transfers to the RPC Service for EOA custodial gas tanks.
/// @dev This contract implicitly trusts the RPC_SERVICE address.
contract GasTankManager {
    address public immutable RPC_SERVICE;
    uint256 public immutable MINIMUM_DEPOSIT;

    event FundsRecovered(address indexed owner, uint256 indexed amount);
    event GasTankFunded(address indexed smartAccount, address indexed caller, uint256 indexed amount);

    error FailedToRecoverFunds(address owner, uint256 amount);
    error NotValidCaller(address caller);
    error InvalidAmount();
    error FailedToFundGasTank(address rpcProvider, uint256 transferAmount);
    error RPCServiceNotSet(address provider);
    error NotRPCService(address caller);
    error InsufficientFunds(uint256 currentBalance, uint256 requiredBalance);
    error NotThisEOA(address msgSender, address thisAddress);
    error MinimumDepositNotMet(uint256 amountToTransfer, uint256 minimumDeposit);

    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    modifier onlyRPCService() {
        require(msg.sender == RPC_SERVICE, NotRPCService(msg.sender));
        _;
    }

    constructor(address rpcService, uint256 _minDeposit) {
        require(rpcService != address(0), RPCServiceNotSet(rpcService));
        require(_minDeposit > 0, MinimumDepositNotMet(0, _minDeposit));
        RPC_SERVICE = rpcService;
        MINIMUM_DEPOSIT = _minDeposit;
    }

    receive() external payable { /* ETH transfers allowed. */ }

    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Recovers funds inadvertently sent to this contract directly.
    function recoverFunds() external onlyRPCService {
        uint256 balance = address(this).balance;

        (bool success,) = RPC_SERVICE.call{value: balance}("");
        require(success, FailedToRecoverFunds(RPC_SERVICE, balance));

        emit FundsRecovered(RPC_SERVICE, balance);
    }

    /// @notice Transfers ETH from the EOA's balance to the Gas RPC Service.
    /// @param _amount The amount to fund the gas tank with.
    /// @dev Only the EOA can call this function.
    function fundGasTank(uint256 _amount) external onlyThisEOA {
        require(_amount >= MINIMUM_DEPOSIT, MinimumDepositNotMet(_amount, MINIMUM_DEPOSIT));
        _fundGasTank(_amount);
    }

    /// @notice Transfers the minimum deposit amount of ETH from the EOA's balance to the Gas RPC Service.
    /// @dev Only the RPC Service can call this function.
    function fundGasTank() external onlyRPCService {
        _fundGasTank(MINIMUM_DEPOSIT);
    }

    /// @dev `fundGasTank` Internal function to fund the gas tank.
    function _fundGasTank(uint256 _amountToTransfer) internal {
        require(address(this).balance >= _amountToTransfer, InsufficientFunds(address(this).balance, _amountToTransfer));

        (bool success,) = RPC_SERVICE.call{value: _amountToTransfer}("");
        if (!success) {
            revert FailedToFundGasTank(RPC_SERVICE, _amountToTransfer);
        }

        emit GasTankFunded(address(this), msg.sender, _amountToTransfer);
    }
}
