// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "forge-std/console.sol";
import "../../contracts/validator-registry/ValidatorRegistryV1.sol";
import {IValidatorRegistryV1} from "../../contracts/interfaces/IValidatorRegistryV1.sol";

// Script to e2e test the ValidatorRegistryV1 contract with anvil, also see makefile.
abstract contract ExampleScript is Script {

    ValidatorRegistryV1 internal _validatorRegistry = ValidatorRegistryV1(payable(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512));

    // When starting anvil: 
    // 
    // Available Accounts
    // ==================
    // (0) 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 (10000.000000000000000000 ETH)
    // (1) 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10000.000000000000000000 ETH)
    // (2) 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC (10000.000000000000000000 ETH)
    //
    // Private Keys
    // =============
    // (0) 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
    // (1) 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
    // (2) 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
    address public defaultEOA = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;
    address public defaultEOA2 = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address public defaultEOA3 = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;

    function checkStaking(bytes[] memory blsKeys) public view {

        console.log("--------------------");
        console.log("Checking Staking related state...");
        console.log("--------------------");
        
        for (uint i = 0; i < blsKeys.length; i++) {
            bool isStaked = _validatorRegistry.isValidatorOptedIn(blsKeys[i]);
            console.log("--------------------");
            console.log("BLS Key: ");
            console.logBytes(blsKeys[i]);
            console.log("is Validator Opted In:", isStaked);
            uint256 stakedAmount = _validatorRegistry.getStakedAmount(blsKeys[i]);
            console.log("Staked Amount:", stakedAmount);
            bool isUnstaking = _validatorRegistry.isUnstaking(blsKeys[i]);
            console.log("is Unstaking: %s", isUnstaking);
            IValidatorRegistryV1.StakedValidator memory stakedValidator = _validatorRegistry.getStakedValidator(blsKeys[i]);
            console.log("Staked Validator balance: %s", stakedValidator.balance);
            console.log("Staked Validator withdrawalAddress: %s", stakedValidator.withdrawalAddress);
            console.log("Staked Validator unstakeBlockNum: %s", stakedValidator.unstakeHeight.blockHeight);
        }
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

        require(msg.sender == defaultEOA, "must be 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb922266:", defaultEOA.balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        checkStaking(validators);

        console.log("--------------------");

        uint256 totalAmount = 10 ether;
        _validatorRegistry.stake{value: totalAmount}(validators);
        console.log("Stake completed ", totalAmount, "ETH");

        console.log("--------------------");

        checkStaking(validators);

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", defaultEOA.balance);
        vm.stopBroadcast();
    }
}

contract UnstakeExample is ExampleScript {

    function run() external {
        vm.startBroadcast();

        require(msg.sender == defaultEOA, "must be 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
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

        require(msg.sender == defaultEOA, "must be 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266");
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

contract SlashExample is ExampleScript {

    function run() external {
        vm.startBroadcast();

        require(msg.sender == defaultEOA2, "slash oracle must be 0x70997970C51812dc3A010C7d01b50e0d17dc79C8");
        console.log("Balance of slash oracle @ 0x70997970C51812dc3A010C7d01b50e0d17dc79C8:", defaultEOA2.balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        checkStaking(validators);

        _validatorRegistry.slash(validators);
        console.log("Slash initiated");

        checkStaking(validators);

        console.log("Balance of 0x70997970C51812dc3A010C7d01b50e0d17dc79C8:", defaultEOA2.balance);
        console.log("Balance of 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC:", defaultEOA3.balance);
        vm.stopBroadcast();
    }
}
