// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

// solhint-disable no-console
import {Script} from "forge-std/Script.sol";
import {MevCommitBapp} from "../../../contracts/validator-registry/ssv/MevCommitBapp.sol";
import {console} from "forge-std/console.sol";

contract DeployMevCommitBapp is Script {

    address constant public SSV_NETWORK = 0xc7fCFeEc5FB9962bDC2234A7a25dCec739e27f9f;

    function run() external {
        require(block.chainid == 560048, "must deploy on hoodi");
        vm.startBroadcast();
        MevCommitBapp mevCommitBapp = new MevCommitBapp(SSV_NETWORK, msg.sender);
        console.log("MevCommitBapp deployed at", address(mevCommitBapp));
        vm.stopBroadcast();
    }
}

contract RegisterBapp is Script {
    function run() external {
        require(block.chainid == 560048, "must deploy on hoodi");
        vm.startBroadcast();
        MevCommitBapp mevCommitBapp = MevCommitBapp(payable(0xAebd263517047715f539D1674f86D2ee622753F4));
        address[] memory tokens = new address[](1);
        tokens[0] = 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE;
        mevCommitBapp.registerBApp(
            tokens,
            new uint32[](1),
            "https://raw.githubusercontent.com/primev/mev-commit/refs/heads/main/static/avs-metadata.json"
        );
        vm.stopBroadcast();
    }
}
