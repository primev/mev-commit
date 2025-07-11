// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

// solhint-disable no-console
import {Script} from "forge-std/Script.sol";
import {MevCommitBapp} from "../../../contracts/validator-registry/ssv/MevCommitBapp.sol";
import {console} from "forge-std/console.sol";

contract DeployMevCommitBapp is Script {

    address constant public SSV_NETWORK = 0x58410Bef803ECd7E63B23664C586A6DB72DAf59c;

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
        MevCommitBapp mevCommitBapp = MevCommitBapp(payable(0x43895A07EF22560e9e1319871870b788e6458797));
        address[] memory tokens = new address[](1);
        tokens[0] = 0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE;
        mevCommitBapp.registerBApp(
            tokens,
            new uint32[](1),
            "https://github.com/primev/mev-commit/blob/main/static/avs-metadata.json"
        );
        vm.stopBroadcast();
    }
}
