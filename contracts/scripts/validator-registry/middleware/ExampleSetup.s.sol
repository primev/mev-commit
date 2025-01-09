// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {SymbioticHoleskyDevnetConsts} from "./ReleaseAddrConsts.s.sol";
import {IBaseDelegator} from "symbiotic-core/interfaces/delegator/IBaseDelegator.sol";
import {IBurnerRouter} from "symbiotic-burners/interfaces/router/IBurnerRouter.sol";
import {IBurnerRouterFactory} from "symbiotic-burners/interfaces/router/IBurnerRouterFactory.sol";
import {IVaultConfigurator} from "symbiotic-core/interfaces/IVaultConfigurator.sol";
import {IVault} from "symbiotic-core/interfaces/vault/IVault.sol";
import {IBaseSlasher} from "symbiotic-core/interfaces/slasher/IBaseSlasher.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IOperatorRegistry} from "symbiotic-core/interfaces/IOperatorRegistry.sol";
import {IOptInService} from "symbiotic-core/interfaces/service/IOptInService.sol";
import {ISlasher} from "symbiotic-core/interfaces/slasher/ISlasher.sol";
import {IOperatorSpecificDelegator} from "symbiotic-core/interfaces/delegator/IOperatorSpecificDelegator.sol";
import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";

contract SetupVault is Script {
    function run() external {
        vm.startBroadcast();

        // Deploy burner router
        address burnerRouterFactory = 0x32e2AfbdAffB1e675898ABA75868d92eE1E68f3b; // From https://docs.symbiotic.fi/deployments/current
        IBurnerRouter.NetworkReceiver[] memory networkReceivers = new IBurnerRouter.NetworkReceiver[](1);
        networkReceivers[0] = IBurnerRouter.NetworkReceiver({
            network: 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1, // MevCommitMiddleware.network()
            receiver: 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1 // MevCommitMiddleware.slashReceiver()
        });
        address burnerRouter = IBurnerRouterFactory(burnerRouterFactory).create(
            IBurnerRouter.InitParams({
               owner: msg.sender,                       
               collateral: 0x3F1c547b21f65e10480dE3ad8E19fAAC46C95034, // stETH                  
               delay: 3 days, // > 2 days
               globalReceiver: msg.sender,             
               networkReceivers: networkReceivers,
               operatorNetworkReceivers: new IBurnerRouter.OperatorNetworkReceiver[](0) // Empty or same as networkReceivers
         }));

        console.log("Burner router deployed to:", address(burnerRouter));

        // Deploy vault with delegator and slasher
        IVaultConfigurator vaultConfigurator = IVaultConfigurator(0xD2191FE92987171691d552C219b8caEf186eb9cA); // From https://docs.symbiotic.fi/deployments/current

        address[] memory networkLimitSetRoleHolders = new address[](1);
        networkLimitSetRoleHolders[0] = msg.sender;
        address[] memory operatorNetworkSharesSetRoleHolders = new address[](1);
        operatorNetworkSharesSetRoleHolders[0] = msg.sender;
        IVaultConfigurator.InitParams memory initParams = IVaultConfigurator.InitParams({
            version: 1,                                                                   
            owner: msg.sender,                            
            vaultParams: abi.encode(IVault.InitParams({
                collateral: 0x3F1c547b21f65e10480dE3ad8E19fAAC46C95034, // stETH
                burner: address(burnerRouter),                                                   
                epochDuration: 1 weeks,
                depositWhitelist: false, 
                isDepositLimit: false, 
                depositLimit: 0, 
                defaultAdminRoleHolder: msg.sender, 
                depositWhitelistSetRoleHolder: msg.sender, 
                depositorWhitelistRoleHolder: msg.sender, 
                isDepositLimitSetRoleHolder: msg.sender, 
                depositLimitSetRoleHolder: msg.sender
            })),
            delegatorIndex: 2, // OperatorSpecificDelegator
            delegatorParams: abi.encode(IOperatorSpecificDelegator.InitParams({
                baseParams: IBaseDelegator.BaseParams({
                    defaultAdminRoleHolder: msg.sender,
                    hook: 0x0000000000000000000000000000000000000000,
                    hookSetRoleHolder: msg.sender
                }),
                networkLimitSetRoleHolders: networkLimitSetRoleHolders,
                operator: 0xb4F13624966E874967d7C9231F2F740F03F1A832
            })),
            withSlasher: true,
            slasherIndex: 0, // Instant slasher
            slasherParams: abi.encode(ISlasher.InitParams({
                baseParams: IBaseSlasher.BaseParams({
                    isBurnerHook: true
                })
            }))
        });

        (address vault, address networkRestakeDelegator, address vetoSlasher) = vaultConfigurator.create(initParams);

        console.log("Vault deployed to:", address(vault));
        console.log("Delegator deployed to:", address(networkRestakeDelegator));
        console.log("Slasher deployed to:", address(vetoSlasher));

        vm.stopBroadcast();
    }
}

contract DepositToVault is Script {
    function run() external {
        vm.startBroadcast();

        address stEthVault = 0x5DF518571733d5F4f496D76C9087110FAe98a946;
        address stEthToken = 0x3F1c547b21f65e10480dE3ad8E19fAAC46C95034;
        
        IERC20(stEthToken).approve(stEthVault, 0.1 ether);
        
        IVault vault = IVault(stEthVault);
        vault.deposit(msg.sender, 0.1 ether);

        vm.stopBroadcast();
    }
}

contract OperatorActions is Script {
    function run() external {
        vm.startBroadcast();
        IOperatorRegistry operatorRegistry = IOperatorRegistry(SymbioticHoleskyDevnetConsts.OPERATOR_REGISTRY);
        operatorRegistry.registerOperator();

        IOptInService vaultOptInService = IOptInService(SymbioticHoleskyDevnetConsts.VAULT_OPT_IN_SERVICE);
        address stEthVault = 0x5DF518571733d5F4f496D76C9087110FAe98a946;
        vaultOptInService.optIn(stEthVault);

        IOptInService networkOptInService = IOptInService(SymbioticHoleskyDevnetConsts.NETWORK_OPT_IN_SERVICE);
        networkOptInService.optIn(0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1); // mev-commit network addr
        vm.stopBroadcast();
    }
}

contract PrimevTeamActions is Script {
    function run() external {
        vm.startBroadcast();

        IMevCommitMiddleware existingMevCommitMiddleware = IMevCommitMiddleware(payable(0x79FeCD427e5A3e5f1a40895A0AC20A6a50C95393));

        address[] memory vaults = new address[](1);
        vaults[0] = 0x5DF518571733d5F4f496D76C9087110FAe98a946;

        uint160[] memory slashAmounts = new uint160[](1);
        slashAmounts[0] = 0.01 ether;

        existingMevCommitMiddleware.registerVaults(vaults, slashAmounts);

        IBaseDelegator delegator = IBaseDelegator(IVault(vaults[0]).delegator());
        delegator.setMaxNetworkLimit(1, 1000000 ether);

        address[] memory operators = new address[](1);
        operators[0] = 0xb4F13624966E874967d7C9231F2F740F03F1A832;
        existingMevCommitMiddleware.registerOperators(operators);

        vm.stopBroadcast();
    }
}

contract VaultActions is Script { 
    function run() external {
        vm.startBroadcast();

        IOperatorSpecificDelegator delegator = IOperatorSpecificDelegator(0x75b131De299A5D343b9408081DD6A8D6a9891b8c);
        delegator.setNetworkLimit(0x4535bd6ff24860b5fd2889857651a85fb3d3c6b1000000000000000000000001, 1000000 ether);

        uint256 stake = delegator.stake(0x4535bd6ff24860b5fd2889857651a85fb3d3c6b1000000000000000000000001, msg.sender);
        console.log("Stake toward operator:", stake);

        vm.stopBroadcast();
    }
}
