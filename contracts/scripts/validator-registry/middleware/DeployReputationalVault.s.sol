// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable no-console

import {Script} from "forge-std/Script.sol";
import {IVault} from "symbiotic-core/interfaces/vault/IVault.sol";
import {IBurnerRouterFactory} from "symbiotic-burners/interfaces/router/IBurnerRouterFactory.sol";
import {IBurnerRouter} from "symbiotic-burners/interfaces/router/IBurnerRouter.sol";
import {MevCommitMiddlewareStorage} from "../../../contracts/validator-registry/middleware/MevCommitMiddlewareStorage.sol";

contract VaultScript is Script {

    address public constant NETWORK = 0x9101eda106A443A0fA82375936D0D1680D5a64F5;
    address public constant BURNER_ROUTER_NETWORK_RECEIVER = 0xD5881f91270550B8850127f05BD6C8C203B3D33f;
    address public constant BURNER_ROUTER_FACTORY = 0x99F2B89fB3C363fBafD8d826E5AA77b28bAB70a0;
    address public constant OWNER = 0x9101eda106A443A0fA82375936D0D1680D5a64F5;
    address public constant COLLATERAL = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2; // WETH
    uint48 public constant BURNER_ROUTER_DELAY = 2.1 days;

    uint48 public constant EPOCH_DURATION = 2 days;
    address public constant VAULT_CONFIGURATOR = 0x29300b1d3150B4E2b12fE80BE72f365E200441EC;
    uint256 public constant DEPOSIT_LIMIT = 0;
    uint64 public constant DELEGATOR_INDEX = MevCommitMiddlewareStorage._NETWORK_RESTAKE_DELEGATOR_TYPE;
    address public constant HOOK = 0x0000000000000000000000000000000000000000;
    uint64 public constant SLASHER_INDEX = MevCommitMiddlewareStorage._INSTANT_SLASHER_TYPE;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        IBurnerRouter.NetworkReceiver[] memory networkReceivers = new IBurnerRouter.NetworkReceiver[](1);
        networkReceivers[0] = IBurnerRouter.NetworkReceiver({
            network: NETWORK,
            receiver: BURNER_ROUTER_NETWORK_RECEIVER
        });
        IBurnerRouter.OperatorNetworkReceiver[] memory operatorNetworkReceivers = new IBurnerRouter.OperatorNetworkReceiver[](0);

        address burnerRouter = IBurnerRouterFactory(BURNER_ROUTER_FACTORY).create(
            IBurnerRouter.InitParams({
                owner: OWNER,
                collateral: COLLATERAL,
                delay: BURNER_ROUTER_DELAY,
                globalReceiver: OWNER,
                networkReceivers: networkReceivers,
                operatorNetworkReceivers: operatorNetworkReceivers
            })
        );
        console.log("Burner router address: ", burnerRouter);

        bytes memory vaultParams = abi.encode(
            IVault.InitParams({
                collateral: COLLATERAL,
                burner: burnerRouter,
                epochDuration: EPOCH_DURATION,
                depositWhitelist: false,
                isDepositLimit: false,
                depositLimit: 0,
                defaultAdminRoleHolder: OWNER,
                depositWhitelistSetRoleHolder: OWNER,
                depositorWhitelistRoleHolder: OWNER,
                isDepositLimitSetRoleHolder: OWNER,
                depositLimitSetRoleHolder: OWNER
            })
        );

        uint256 roleHolders = 1;
        address[] memory networkLimitSetRoleHolders = new address[](roleHolders);
        address[] memory operatorNetworkSharesSetRoleHolders = new address[](roleHolders);
        networkLimitSetRoleHolders[0] = OWNER;
        operatorNetworkSharesSetRoleHolders[0] = OWNER;

        bytes memory delegatorParams = abi.encode(
            INetworkRestakeDelegator.InitParams({
                baseParams: IBaseDelegator.BaseParams({
                    defaultAdminRoleHolder: OWNER,
                    hook: address(0),
                    hookSetRoleHolder: OWNER
                }),
                networkLimitSetRoleHolders: networkLimitSetRoleHolders,
                operatorNetworkSharesSetRoleHolders: operatorNetworkSharesSetRoleHolders
                })
            );

        bytes memory slasherParams = abi.encode(
            ISlasher.InitParams({baseParams: IBaseSlasher.BaseParams({isBurnerHook: true})})
        );

        (address vault_, address delegator_, address slasher_) = IVaultConfigurator(vaultConfigurator).create(
            IVaultConfigurator.InitParams({
                version: 1,
                owner: OWNER,
                vaultParams: vaultParams,
                delegatorIndex: DELEGATOR_INDEX,
                delegatorParams: delegatorParams,
                withSlasher: true,
                slasherIndex: SLASHER_INDEX,
                slasherParams: slasherParams
            })
        );

        console.log("Vault address: ", vault_);
        console.log("Delegator address: ", delegator_);
        console.log("Slasher address: ", slasher_);

        vm.stopBroadcast();
    }
}