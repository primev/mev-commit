// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IProviderRegistry} from "../contracts/interfaces/IProviderRegistry.sol";
import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";

contract RegisterProvider is Script {
    function run() public {
        vm.startBroadcast();
        if (msg.sender != address(0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC)) {
            revert("incorrect sender");
        }
        address providerRegistryAddr = 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853;
        IProviderRegistry(providerRegistryAddr).registerAndStake{value: 10 ether}();
        console.log("Provider registered and staked 10 ETH");
        vm.stopBroadcast();
    }
}

contract ManuallyOverrideBLSKey is Script {
    function run() public {
        vm.startBroadcast();
        if (msg.sender != address(0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266)) {
            revert("incorrect sender");
        }
        address providerRegistryAddr = 0xa513E6E4b8f2a923D98304ec87F64353C4D5C853;
        address providerAddr = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
        bytes memory blsKey = hex"aca4b5c5daf5b39514b8aa6e5f303d29f6f1bd891e5f6b6b2ae6e2ae5d95dee31cd78630c1115b6e90f3da1a66cf8edb";
        IProviderRegistry(providerRegistryAddr).overrideAddBLSKey(providerAddr, blsKey);
        console.log("BLS key set");
        vm.stopBroadcast();
    }
}

contract GetCode is Script {
    function run() public {
        vm.startBroadcast();
        address bidderAddr = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
        uint256 codeLength = bidderAddr.code.length;
        console.log("Bidder EOA's code length:", codeLength);
        console.log("Bidder EOA's code hash: ");
        console.logBytes32(bidderAddr.codehash);
        if (codeLength != 0) {
            console.log("Expected code hash: Keccak256Hash(0xef0100, depositManagerImpl)");
        }
        vm.stopBroadcast();
    }
}
