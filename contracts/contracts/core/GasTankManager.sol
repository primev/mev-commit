// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IGasTankManager} from "../interfaces/IGasTankManager.sol";
import {GasTankManagerStorage} from "./GasTankManagerStorage.sol";
import {Errors} from "../utils/Errors.sol";

/// @title GasTankManager
/// @notice Coordinates on-demand ETH bridging so EOA gas tank balances stay funded for mev-commit transactions.
/// @dev The RPC provider manages EOA gas tank balances on the Mev Commit chain.
/// @dev Flow overview:
/// - EOA (L1)
///   - Authorizes and sets this contract as a delegate against its own EOA address. (ERC-7702)
///   - Calls `sendMinimumDeposit` to send the initial ETH funds to their gas tank on the Mev Commit chain.
/// - Provider (mev-commit)
///   - Triggers a top-up via `topUpGasTank` when the gas tank balance cannot cover the next mev-commit transaction. This transfer amount is always the difference
///   between the `minDeposit` and the current balance of this contract. The provider is the only one who can trigger a top-up.
///   - These funds are then transferred to the EOA gas tank on the Mev Commit chain.
///   - When a mev-commit transaction is made, the provider deducts the amount needed from the EOA's gas tank.
contract GasTankManager is IGasTankManager, GasTankManagerStorage, Ownable2StepUpgradeable, UUPSUpgradeable {
    /// @notice Restricts calls to those triggered internally via `onlyThisEOA`.
    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    /// @notice Locks the implementation upon deployment.
    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/writing-upgradeable#initializing-the-implementation-contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Accepts direct ETH deposits and forwards them to the provider.
    receive() external payable {
        _sendFundsToProvider(msg.value);
    }

    /// @notice Reverts on any call data to keep the interface surface narrow.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Initializes ownership and baseline deposit requirement.
    /// @param _owner EOA managed by the RPC provider.
    /// @param _minDeposit Initial deposit requirement in wei.
    function initialize(address _owner, uint256 _minDeposit) external initializer {
        minDeposit = _minDeposit;
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    /// @inheritdoc IGasTankManager
    function setMinimumDeposit(uint256 _minDeposit) external onlyOwner {
        minDeposit = _minDeposit;
        emit MinimumDepositSet(_minDeposit);
    }

    /// @inheritdoc IGasTankManager
    function topUpGasTank(uint256 _gasTankBalance) external onlyOwner {
        uint256 minDeposit_ = minDeposit;
        uint256 available = address(this).balance;
        uint256 needed = minDeposit_ - _gasTankBalance;

        _validateTopUp(_gasTankBalance, available, needed, minDeposit_);

        _sendFundsToProvider(needed);
        emit GasTankToppedUp(needed);
    }

    /// @inheritdoc IGasTankManager
    function fundGasTank() external payable {
        _sendFundsToProvider(msg.value);
    }

    /// @inheritdoc IGasTankManager
    function sendMinimumDeposit() external payable onlyThisEOA {
        require(msg.value >= minDeposit, InsufficientFundsSent(msg.value, minDeposit));
        _sendFundsToProvider(minDeposit);
    }

    /// @notice Restricts upgradeability to the owner.
    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/foundry/api/upgrades
    /// solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @notice Forwards ETH to the configured provider address.
    /// @dev Reverts when the provider rejects the transfer.
    /// @param _amount Amount of wei to transfer.
    function _sendFundsToProvider(uint256 _amount) internal {
        address provider = owner();
        require(provider != address(0), ProviderNotSet(provider));

        (bool success,) = provider.call{value: _amount}("");
        require(success, GasTankTopUpFailed(provider, _amount));
    }

    /// @notice Validates the proposed refill before forwarding funds.
    /// @dev Applies guardrails to protect configured thresholds.
    /// @param _gasTankBalance Current gas tank balance.
    /// @param _available Current contract (EOA) balance.
    /// @param _needed Amount required to reach the `minDeposit`.
    /// @param _minDeposit Minimum deposit requirement in wei.
    function _validateTopUp(uint256 _gasTankBalance, uint256 _available, uint256 _needed, uint256 _minDeposit)
        internal
        view
    {
        require(_gasTankBalance < _minDeposit, GasTankBalanceIsSufficient(_gasTankBalance, _minDeposit));
        require(_available > _needed, InsufficientEOABalance(address(this), _available, _needed));
    }
}
