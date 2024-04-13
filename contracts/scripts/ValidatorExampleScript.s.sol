// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "forge-std/Script.sol";
import "../contracts/ValidatorRegistry.sol";

contract StakeAndUnstakeExample is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0x9938f7EfB83dd3150cF3B784CEa473D29fe7cBF0);

    function run() external {
        vm.startBroadcast();
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        selfStake(3.1 ether);
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));

        checkIsStaked(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        checkIsStaked(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        checkIsStaked(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        splitStake(10 ether);
        checkIsStaked(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        checkIsStaked(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        checkIsStaked(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        
        address[] memory addrs = new address[](3);
        addrs[0] = address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317);
        addrs[1] = address(0x7802BF57d9f5a449A879E3cF89280157846651c3);
        addrs[2] = address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D);
        unstake(addrs);
        checkIsStaked(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        checkIsStaked(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        checkIsStaked(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));

        vm.stopBroadcast();
    }

    function checkIsStaked(address addr) public view {
        bool isStaked = _validatorRegistry.isStaked(addr);
        console.log("Address", addr, "isStaked:", isStaked);
    }

    function selfStake(uint256 amount) public {
        _validatorRegistry.selfStake{value: amount}();
        console.log("Performed selfStake with amount:", amount);
    }

    function splitStake(uint256 totalAmount) public {
        
        address[] memory recipients = new address[](3);
        recipients[0] = address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317);
        recipients[1] = address(0x7802BF57d9f5a449A879E3cF89280157846651c3);
        recipients[2] = address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D);
        
        _validatorRegistry.splitStake{value: totalAmount}(recipients);
        console.log("Split stake occurred");
    }

    function unstake(address[] memory unstakeAddresses) public {
        _validatorRegistry.unstake(unstakeAddresses);
        console.log("Initiated unstake process for addresses.");
    }
}

contract WithdrawExample is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0x9938f7EfB83dd3150cF3B784CEa473D29fe7cBF0);

    function run() external {
        vm.startBroadcast();

        checkIsStaked(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        checkIsStaked(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        checkIsStaked(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        uint256 initialStakedAmount = _getStakedAmount(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        console.log("Initial staked amount for 0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317:", initialStakedAmount);
        initialStakedAmount = _getStakedAmount(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        console.log("Initial staked amount for 0x7802BF57d9f5a449A879E3cF89280157846651c3:", initialStakedAmount);
        initialStakedAmount = _getStakedAmount(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        console.log("Initial staked amount for 0x8662a945619e31182894C641fC5bf74E7Cd75A7D:", initialStakedAmount);

        address[] memory addrs = new address[](3);
        addrs[0] = address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317);
        addrs[1] = address(0x7802BF57d9f5a449A879E3cF89280157846651c3);
        addrs[2] = address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D);
        withdraw(addrs);

        checkIsStaked(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        checkIsStaked(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        checkIsStaked(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        uint256 finalStakedAmount = _getStakedAmount(address(0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317));
        console.log("Final staked amount for 0x3d77fE0CeB523FAa006Cd2408F6Cb34A234C4317:", finalStakedAmount);
        finalStakedAmount = _getStakedAmount(address(0x7802BF57d9f5a449A879E3cF89280157846651c3));
        console.log("Final staked amount for 0x7802BF57d9f5a449A879E3cF89280157846651c3:", finalStakedAmount);
        finalStakedAmount = _getStakedAmount(address(0x8662a945619e31182894C641fC5bf74E7Cd75A7D));
        console.log("Final staked amount for 0x8662a945619e31182894C641fC5bf74E7Cd75A7D:", finalStakedAmount);

        vm.stopBroadcast();
    }

    function checkIsStaked(address addr) public view {
        bool isStaked = _validatorRegistry.isStaked(addr);
        console.log("Address", addr, "isStaked:", isStaked);
    }

    function withdraw(address[] memory addrs) public {
        _validatorRegistry.withdraw(addrs);
        console.log("Initiated withdraw process for addresses.");
    }

    function _getStakedAmount(address staker) internal view returns (uint256) {
        return _validatorRegistry.getStakedAmount(staker);
    }
}

contract SelfStake is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0x9938f7EfB83dd3150cF3B784CEa473D29fe7cBF0);

    function run() external {
        vm.startBroadcast();
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        stakedBalance(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        selfStake(3.1 ether);
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        vm.stopBroadcast();
    }

    function checkIsStaked(address addr) public view {
        bool isStaked = _validatorRegistry.isStaked(addr);
        console.log("Address", addr, "isStaked:", isStaked);
    }

    function stakedBalance(address addr) public view {
        uint256 balance = _validatorRegistry.getStakedAmount(addr);
        console.log("Staked balance for address", addr, "is:", balance);
    }

    function selfStake(uint256 amount) public {
        _validatorRegistry.selfStake{value: amount}();
        console.log("Performed selfStake with amount:", amount);
    }
}

contract Unstake is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0x9938f7EfB83dd3150cF3B784CEa473D29fe7cBF0);

    function run() external {
        vm.startBroadcast();
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        unstake(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        vm.stopBroadcast();
    }

    function checkIsStaked(address addr) public view {
        bool isStaked = _validatorRegistry.isStaked(addr);
        console.log("Address", addr, "isStaked:", isStaked);
    }

    function unstake(address addr) public {
        address[] memory addrs = new address[](1);
        addrs[0] = addr;
        _validatorRegistry.unstake(addrs);
        console.log("Initiated unstake process for address.");
    }
}

contract Withdraw is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0x9938f7EfB83dd3150cF3B784CEa473D29fe7cBF0);

    function run() external {
        
        console.log("Balance of msg.sender:", msg.sender.balance);

        uint256 beforeBalance = msg.sender.balance;

        vm.startBroadcast();
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        withdraw(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        checkIsStaked(address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266));
        vm.stopBroadcast();

        console.log("Balance of msg.sender:", msg.sender.balance);

        uint256 afterBalance = msg.sender.balance;
        console.log("Difference in balance:", afterBalance - beforeBalance);
    }

    function checkIsStaked(address addr) public view {
        bool isStaked = _validatorRegistry.isStaked(addr);
        console.log("Address", addr, "isStaked:", isStaked);
    }

    function withdraw(address addr) public {
        address[] memory addrs = new address[](1);
        addrs[0] = addr;
        _validatorRegistry.withdraw(addrs);
        console.log("Initiated withdraw process for address.");
    }
}
