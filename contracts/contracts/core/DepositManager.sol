// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {Errors} from "../utils/Errors.sol";

contract DepositManager {

    address public immutable BIDDER_REGISTRY;
    uint256 public immutable MIN_BALANCE;

    mapping(address => uint256) public targetDeposits;

    event TargetDepositSet(address indexed provider, uint256 amount);
    event WithdrawalRequestExists(address indexed provider);
    event TargetDepositDoesNotExist(address indexed provider);
    event CurrentDepositIsSufficient(address indexed provider, uint256 currentDeposit, uint256 targetDeposit);
    event CurrentBalanceAtOrBelowMin(address indexed provider, uint256 currentBalance, uint256 minBalance);
    event TopUpReduced(address indexed provider, uint256 available, uint256 needed);
    event DepositToppedUp(address indexed provider, uint256 amount);

    error NotThisEOA(address msgSender, address thisAddress);

    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    constructor(address _registry, uint256 _minBalance) {
        BIDDER_REGISTRY = _registry;
        MIN_BALANCE = _minBalance;
    }

    receive() external payable { 
        // Eth transfers allowed.
    }

    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    function setTargetDeposits(
        address[] calldata providers,
        uint256[] calldata amounts
    ) external onlyThisEOA {
        uint256 length = providers.length;
        for (uint256 i = 0; i < length; ++i) {
            targetDeposits[providers[i]] = amounts[i];
            emit TargetDepositSet(providers[i], amounts[i]);
        }
    }

    function setTargetDeposit(address provider, uint256 amount) external onlyThisEOA {
        targetDeposits[provider] = amount;
        emit TargetDepositSet(provider, amount);
    }

    /// @notice Top-up deposits if needed, as configured by this EOA.
    /// @param provider to top-up the deposit for.
    /// @dev This function will be called automatically by external addresses.
    function topUpDeposit(address provider) external {
        if (IBidderRegistry(BIDDER_REGISTRY).withdrawalRequestExists(address(this), provider)) {
            emit WithdrawalRequestExists(provider);
            return;
        }

        uint256 target = targetDeposits[provider];
        if (target == 0) {
            emit TargetDepositDoesNotExist(provider);
            return;
        }

        uint256 currentDeposit = IBidderRegistry(BIDDER_REGISTRY).getDeposit(address(this), provider);
        if (currentDeposit >= target) {
            emit CurrentDepositIsSufficient(provider, currentDeposit, target);
            return;
        }
        uint256 needed = target - currentDeposit; // No underflow/overflow, target must be greater than current deposit

        uint256 currentBalance = address(this).balance;
        if (currentBalance <= MIN_BALANCE) {
            emit CurrentBalanceAtOrBelowMin(provider, currentBalance, MIN_BALANCE);
            return;
        }

        uint256 available = currentBalance - MIN_BALANCE; // No underflow/overflow, currentBalance must be greater than minBalance
        if (available < needed) {
            emit TopUpReduced(provider, available, needed);
            needed = available;
        }
        IBidderRegistry(BIDDER_REGISTRY).depositAsBidder{value: needed}(provider);
        emit DepositToppedUp(provider, needed);
    }
}
