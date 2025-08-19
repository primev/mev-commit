// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {console} from "forge-std/console.sol";

// Update the import path to wherever you placed the contract
import {LidoV3Middleware} from "../../../contracts/validator-registry/lido/LidoV3Middleware.sol";

contract BaseDeploy is Script {
    function deployLidoMiddleware(
        address owner,
        address vaultHub,
        address unfreezeReceiver,
        uint256 deregistrationPeriodBlocks,
        uint256 unfreezePeriodBlocks,
        uint256 slashAmountWei,
        uint256 unfreezeFeeWei
    ) public returns (address) {
        console.log("Deploying LidoV3Middleware (UUPS) on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "LidoV3Middleware.sol:LidoV3Middleware",
            abi.encodeCall(
                LidoV3Middleware.initialize,
                (
                    owner,
                    vaultHub,
                    unfreezeReceiver,
                    deregistrationPeriodBlocks,
                    unfreezePeriodBlocks,
                    slashAmountWei,
                    unfreezeFeeWei
                )
            )
        );

        console.log("LidoV3Middleware UUPS proxy deployed to:", proxy);
        LidoV3Middleware mw = LidoV3Middleware(payable(proxy));
        console.log("LidoV3Middleware owner:", mw.owner());
        console.log("LidoV3Middleware vaultHub:", mw.vaultHub());
        console.log("LidoV3Middleware slashAmount:", mw.slashAmount());
        console.log("LidoV3Middleware unfreezeFee:", mw.unfreezeFee());
        console.log("LidoV3Middleware deregistrationPeriod:", mw.deregistrationPeriod());
        console.log("LidoV3Middleware unfreezePeriod:", mw.unfreezePeriod());

        return proxy;
    }
}

contract DeployHoodi is BaseDeploy {
    // ---- fixed params for this env ----
    address constant public OWNER             = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public UNFREEZE_RECEIVER = 0x1623fE21185c92BB43bD83741E226288B516134a;
    address constant public VAULT_HUB         = 0x26b92f0fdfeBAf43E5Ea5b5974EeBee95F17Fe08; // Hoodi VaultHub proxy

    // periods in seconds
    uint256 constant public DEREGISTRATION_PERIOD = 600;
    uint256 constant public UNFREEZE_PERIOD       = 600;

    // low on hoodi for testing
    uint256 constant public SLASH_AMOUNT_WEI = 10; // 10 wei
    uint256 constant public UNFREEZE_FEE_WEI = 0;

    function run() external {
        require(block.chainid == 560048, "must deploy on Hoodi");

        vm.startBroadcast();
        deployLidoMiddleware(
            OWNER,
            VAULT_HUB,
            UNFREEZE_RECEIVER,
            DEREGISTRATION_PERIOD,
            UNFREEZE_PERIOD,
            SLASH_AMOUNT_WEI,
            UNFREEZE_FEE_WEI
        );
        vm.stopBroadcast();
    }
}