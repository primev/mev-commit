// SPDX-License-Identifier: MIT

pragma solidity 0.8.20;

import "./ReputationalValReg.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract DelegationValReg is OwnableUpgradeable, UUPSUpgradeable {
    ReputationalValReg public reputationalValReg;
    IERC20 public stETHToken;

    address public constant DEFAULT_STETH_ADDRESS = 0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84;

    struct Delegation {
        address validatorEOA;
        uint256 amount;
    }

    mapping(address => Delegation) public delegations;

    event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount);
    event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount);
    event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // TODO: Determine upgrade logic and test process
    }

    function initialize(
        address _owner,
        address _reputationalValReg,
        address _stETHToken // Optional arg for testing and Holesky. Default address used for mainnet.
    ) external initializer {
        __Ownable_init(_owner);
        reputationalValReg = ReputationalValReg(payable(_reputationalValReg));
        stETHToken = IERC20(_stETHToken == address(0) ? DEFAULT_STETH_ADDRESS : _stETHToken);
    }

    function delegate(address validatorEOA, uint256 amount) external {
        require(reputationalValReg.isEOAWhitelisted(validatorEOA), "Validator EOA must be whitelisted");
        require(delegations[msg.sender].amount == 0, "Already delegated");
        require(stETHToken.transferFrom(msg.sender, address(this), amount), "Token transfer failed");

        delegations[msg.sender] = Delegation({
            validatorEOA: validatorEOA,
            amount: amount
        });

        emit Delegated(msg.sender, validatorEOA, amount);
    }

    function changeDelegation(address newValidatorEOA) external {
        require(reputationalValReg.isEOAWhitelisted(newValidatorEOA), "New validator EOA must be whitelisted");
        Delegation storage delegation = delegations[msg.sender];
        require(delegation.amount > 0, "No active delegation");

        address oldValidatorEOA = delegation.validatorEOA;
        delegation.validatorEOA = newValidatorEOA;

        emit DelegationChanged(msg.sender, oldValidatorEOA, newValidatorEOA, delegation.amount);
    }

    function withdraw() external {
        Delegation storage delegation = delegations[msg.sender];
        require(delegation.amount > 0, "No active delegation");

        uint256 amount = delegation.amount;
        address validatorEOA = delegation.validatorEOA;

        delete delegations[msg.sender];
        require(stETHToken.transfer(msg.sender, amount), "Token transfer failed");

        emit Withdrawn(msg.sender, validatorEOA, amount);
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
    // TODO: withdraw if the EOA becomes frozen or unwhitelisted, withdraw time, etc. 
}
