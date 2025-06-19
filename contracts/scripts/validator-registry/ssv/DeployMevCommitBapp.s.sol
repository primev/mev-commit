// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

// solhint-disable no-console


contract DeployMevCommitBapp is Script {

    address constant public SSV_NETWORK = 0x58410Bef803ECd7E63B23664C586A6DB72DAf59c;

    function run() external {
        require(block.chainid == 560048, "must deploy on hoodi");
        vm.startBroadcast();

        MevCommitBapp mevCommitBapp = new MevCommitBapp(SSV_NETWORK, OWNER);

        vm.stopBroadcast();
    }
}
