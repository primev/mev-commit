// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import "forge-std/Script.sol";
import "forge-std/console.sol";
import "../contracts/ValidatorRegistry.sol";

contract StakeAndUnstakeExample is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);

    function run() external {
        vm.startBroadcast();

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266).balance);

        bytes[] memory validators = new bytes[](3);
        validators[0] = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
        validators[1] = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
        validators[2] = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

        uint256 totalAmount = 10 ether;
        _validatorRegistry.stake{value: totalAmount}(validators);

        for (uint i = 0; i < validators.length; i++) {
            checkIsStaked(validators[i]);
        }

        _validatorRegistry.unstake(validators);

        for (uint i = 0; i < validators.length; i++) {
            checkIsStaked(validators[i]);
        }

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266).balance);

        vm.stopBroadcast();
    }

    function checkIsStaked(bytes memory blsKey) public view {
        bool isStaked = _validatorRegistry.isStaked(blsKey);
        console.log("--------------------");
        console.log("BLS Key: ");
        console.logBytes(blsKey);
        console.log("is Staked:", isStaked);
        uint256 stakedAmount = _validatorRegistry.getStakedAmount(blsKey);
        console.log("Staked Amount:", stakedAmount);
        console.log("--------------------");
    }
}

contract WithdrawExample is Script {
    ValidatorRegistry private _validatorRegistry = ValidatorRegistry(0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512);

    bytes public validator1Key = hex"a97794deb52ea4529d37d283213ca7e298ea9be0a2fec1bb3134a1464ab8cf9eb2c703d1b42dd68d97b5f1c8e74cc0df";
    bytes public validator2Key = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public validator3Key = hex"a840634f574c20e9a35ff80be19309dfd3ace623a093f114e7f44555e3035725b4d1d59b3ce0b2f169871fbe7abc448a";

    function run() external {
        vm.startBroadcast();

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266).balance);

        bytes[] memory blsKeys = new bytes[](3);
        blsKeys[0] = validator1Key;
        blsKeys[1] = validator2Key;
        blsKeys[2] = validator3Key;

        for (uint i = 0; i < blsKeys.length; i++) {
            checkIsStaked(blsKeys[i]);
        }

        _validatorRegistry.withdraw(blsKeys);

        for (uint i = 0; i < blsKeys.length; i++) {
            checkIsStaked(blsKeys[i]);
        }

        console.log("Balance of 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266:", address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266).balance);

        vm.stopBroadcast();
    }

    function checkIsStaked(bytes memory blsKey) public view {
        bool isStaked = _validatorRegistry.isStaked(blsKey);
        console.log("--------------------");
        console.log("BLS Key: ");
        console.logBytes(blsKey);
        console.log("is Staked:", isStaked);
        uint256 stakedAmount = _validatorRegistry.getStakedAmount(blsKey);
        console.log("Staked Amount:", stakedAmount);
        console.log("--------------------");
    }
}
