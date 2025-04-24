// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {PreconfManager} from "../../contracts/core/PreconfManager.sol";
import {console} from "forge-std/console.sol";

/**
 * @notice This script updates the commitment dispatch window of an existing PreconfManager contract.
 *
 * Expected environment variables:
 *  - PRECONF_MANAGER_PROXY: Address of the PreconfManager proxy.
 */
contract UpdatePreconfManagerDispatchWindow is Script {
    function run() external {
        vm.startBroadcast();

        // Retrieve PreconfManager proxy address from environment variables
        address preconfManagerAddress = vm.envAddress("PRECONF_MANAGER_PROXY");
        PreconfManager preconfManager = PreconfManager(
            payable(preconfManagerAddress)
        );

        // Get current commitment dispatch window
        uint64 currentDispatchWindow = preconfManager
            .commitmentDispatchWindow();
        console.log(
            "Current commitment dispatch window:",
            currentDispatchWindow
        );

        // Set new commitment dispatch window value
        uint64 newDispatchWindow = 500;

        // Update the commitment dispatch window
        preconfManager.updateCommitmentDispatchWindow(newDispatchWindow);

        // Verify the update
        uint64 updatedDispatchWindow = preconfManager
            .commitmentDispatchWindow();
        console.log(
            "Updated commitment dispatch window:",
            updatedDispatchWindow
        );

        // Verify update was successful
        require(
            updatedDispatchWindow == newDispatchWindow,
            "Commitment dispatch window update failed"
        );

        console.log(
            "Commitment dispatch window successfully updated to",
            newDispatchWindow
        );

        vm.stopBroadcast();
    }
}
