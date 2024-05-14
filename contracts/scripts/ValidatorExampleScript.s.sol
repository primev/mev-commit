// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "forge-std/console.sol";
import "../contracts/ValidatorRegistry.sol";

abstract contract ExampleScript is Script {

    ValidatorRegistry internal _validatorRegistry = ValidatorRegistry(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);
    address public defaultEOA = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;

    function checkStaking(bytes[] memory blsKeys) public view {

        console.log("--------------------");
        console.log("Checking Staking related state...");
        console.log("--------------------");
        
        for (uint i = 0; i < blsKeys.length; i++) {
            bool isStaked = _validatorRegistry.isStaked(blsKeys[i]);
            console.log("--------------------");
            console.log("BLS Key: ");
            console.logBytes(blsKeys[i]);
            console.log("is Staked:", isStaked);
            uint256 stakedAmount = _validatorRegistry.getStakedAmount(blsKeys[i]);
            console.log("Staked Amount:", stakedAmount);
            uint256 unstakingAmount = _validatorRegistry.getUnstakingAmount(blsKeys[i]);
            console.log("Unstaking Amount:", unstakingAmount);
            console.log("--------------------");
        }

        (uint256 numStakedValidators, uint256 stakedValsetVersion) = _validatorRegistry.getNumberOfStakedValidators();
        console.log("Num Staked Validators:", numStakedValidators);
        console.log("Staked Valset Version from len query:", stakedValsetVersion);
        if (numStakedValidators == 0) {
            return;
        }
        bytes[] memory vals;
        (vals, stakedValsetVersion) = _validatorRegistry.getStakedValidators(0, numStakedValidators);
        for (uint i = 0; i < vals.length; i++) {
            console.log("Staked validator from batch query: ");
            console.logBytes(vals[i]);
        }
        console.log("Staked Valset Version from getStakedValidators query:", stakedValsetVersion);
    }

    function checkWithdrawal(bytes[] memory blsKeys) public view {
        console.log("--------------------");
        console.log("Checking Withdrawal related state...");
        console.log("--------------------");
        for (uint i = 0; i < blsKeys.length; i++) {
            uint256 blocksTillWithdrawAllowed = _validatorRegistry.getBlocksTillWithdrawAllowed(blsKeys[i]);
            console.log("--------------------");
            console.log("BLS Key: ");
            console.logBytes(blsKeys[i]);
            console.log("Blocks till Withdraw Allowed:", blocksTillWithdrawAllowed);
            console.log("--------------------");
        }
    }
}

contract StakeExample is ExampleScript {

    function run() external {
        vm.startBroadcast();

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        checkStaking(validators);

        uint256 totalAmount = 10 ether;
        _validatorRegistry.stake{value: totalAmount}(validators);
        console.log("Stake completed ", totalAmount, "ETH");

        checkStaking(validators);

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);
        vm.stopBroadcast();
    }
}

contract UnstakeExample is ExampleScript {

    function run() external {
        vm.startBroadcast();

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        checkStaking(validators);

        _validatorRegistry.unstake(validators);
        console.log("Unstake initiated");

        checkStaking(validators);

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);

        checkWithdrawal(validators);

        vm.stopBroadcast();
    }
}

contract WithdrawExample is ExampleScript {

    function run() external {
        vm.startBroadcast();

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        checkStaking(validators);

        checkWithdrawal(validators);

        _validatorRegistry.withdraw(validators);
        console.log("Withdraw initiated");

        checkStaking(validators);

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);
        vm.stopBroadcast();
    }
}
