// SPDX-License-Identifier: MIT

pragma solidity 0.8.20;

import "./ReputationValReg.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// TODO: Update notion doc accordingly to include this contract
// TODO: hash out plan on how this and/or other contract could eventually be upgraded to incorporate slashing
// TODO: determine should use reentrancy guard
contract DelegationValReg is OwnableUpgradeable, UUPSUpgradeable {

    ReputationValReg public reputationValReg;
    IERC20 public stETHToken;
    uint256 public withdrawPeriod;

    address public constant DEFAULT_STETH_ADDRESS = 0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84;

    enum State { nonExistant, active, withdrawRequested }

    struct DelegationInfo {
        // State of delegation
        State state;
        // Account (who's being delegated to) representing group of validators 
        address validatorEOA; 
        // Amount of stETH delegated
        uint256 amount;
        // Height of withdrawal being requested
        uint256 withdrawHeight;
    }

    mapping(address => DelegationInfo) public delegations;

    event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount);
    event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount);
    event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // TODO: Determine upgrade logic and test process
    }

    function initialize(
        address _owner,
        address _reputationValReg,
        uint256 _withdrawPeriod,
        address _stETHToken // Optional arg for testing and Holesky. Default address used for mainnet.
    ) external initializer {
        __Ownable_init(_owner);
        reputationValReg = ReputationValReg(payable(_reputationValReg));
        withdrawPeriod = _withdrawPeriod;
        stETHToken = IERC20(_stETHToken == address(0) ? DEFAULT_STETH_ADDRESS : _stETHToken);
    }

    function delegate(address validatorEOA, uint256 amount) external {
        require(getSenderDelegationState() == State.nonExistant, "Delegation must not exist for sender");
        require(reputationValReg.isEOAWhitelisted(validatorEOA), "Validator EOA must be whitelisted");
        require(amount > 0, "Amount must be greater than 0");
        require(stETHToken.transferFrom(msg.sender, address(this), amount), "stETH transfer failed");

        delegations[msg.sender] = DelegationInfo({
            state: State.active,
            validatorEOA: validatorEOA,
            amount: amount,
            withdrawHeight: 0
        });

        emit Delegated(msg.sender, validatorEOA, amount);
    }

    function changeDelegation(address newValidatorEOA) external {
        require(getSenderDelegationState() == State.active, "Active delegation must exist for sender");
        require(reputationValReg.isEOAWhitelisted(newValidatorEOA), "New validator EOA must be whitelisted");
        DelegationInfo storage delegation = delegations[msg.sender];

        address oldValidatorEOA = delegation.validatorEOA;
        delegation.validatorEOA = newValidatorEOA;

        emit DelegationChanged(msg.sender, oldValidatorEOA, newValidatorEOA, delegation.amount);
    }

    function requestWithdraw() external {
        require(getSenderDelegationState() == State.active, "Active delegation must exist for sender");
        delegations[msg.sender].withdrawHeight = block.number + withdrawPeriod;
        delegations[msg.sender].state = State.withdrawRequested;
    }

    function withdraw() external {
        require(getSenderDelegationState() == State.withdrawRequested, "Withdrawal must be requested by sender");
        require(block.number >= delegations[msg.sender].withdrawHeight, "Withdraw period must be elapsed");

        uint256 amount = delegations[msg.sender].amount;
        address validatorEOA = delegations[msg.sender].validatorEOA;
        delete delegations[msg.sender];

        require(stETHToken.transfer(msg.sender, amount), "stETH transfer failed");
        emit Withdrawn(msg.sender, validatorEOA, amount);
    }

    function getSenderDelegationState() public view returns (State) {
        return delegations[msg.sender].state;
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}
