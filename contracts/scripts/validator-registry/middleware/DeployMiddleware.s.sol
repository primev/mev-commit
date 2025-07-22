// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.29;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MevCommitMiddleware} from "../../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {INetworkRegistry} from "symbiotic-core/interfaces/INetworkRegistry.sol";
import {SymbioticHoleskyDevnetConsts, SymbioticMainnetConsts} from "./ReleaseAddrConsts.s.sol";
import {IBaseDelegator} from "symbiotic-core/interfaces/delegator/IBaseDelegator.sol";
import {INetworkMiddlewareService} from "symbiotic-core/interfaces/service/INetworkMiddlewareService.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

contract BaseDeploy is Script {
    function deployMevCommitMiddleware(
        IRegistry networkRegistry,
        IRegistry operatorRegistry,
        IRegistry vaultFactory,
        IRegistry delegatorFactory,
        IRegistry slasherFactory,
        IRegistry burnerRouterFactory,
        address network,
        uint256 slashPeriodSeconds,
        address slashOracle,
        address slashReceiver,
        uint256 minBurnerRouterDelay,
        address owner
    ) public returns (address) {
        console.log("Deploying MevCommitMiddleware on chain:", block.chainid);
        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitMiddleware.sol",
            abi.encodeCall(MevCommitMiddleware.initialize, (
                networkRegistry, 
                operatorRegistry, 
                vaultFactory, 
                delegatorFactory,
                slasherFactory,
                burnerRouterFactory,
                network, 
                slashPeriodSeconds,
                slashOracle,
                slashReceiver,
                minBurnerRouterDelay,
                owner
            ))
        );
        console.log("MevCommitMiddleware UUPS proxy deployed to:", address(proxy));
        MevCommitMiddleware mevCommitMiddleware = MevCommitMiddleware(payable(proxy));
        console.log("MevCommitMiddleware owner:", mevCommitMiddleware.owner());
        return proxy;
    }
}

contract DeployMainnet is BaseDeploy {
    IRegistry constant public NETWORK_REGISTRY = IRegistry(SymbioticMainnetConsts.NETWORK_REGISTRY);
    IRegistry constant public OPERATOR_REGISTRY = IRegistry(SymbioticMainnetConsts.OPERATOR_REGISTRY);
    IRegistry constant public VAULT_FACTORY = IRegistry(SymbioticMainnetConsts.VAULT_FACTORY);
    IRegistry constant public DELEGATOR_FACTORY = IRegistry(SymbioticMainnetConsts.DELEGATOR_FACTORY);
    IRegistry constant public SLASHER_FACTORY = IRegistry(SymbioticMainnetConsts.SLASHER_FACTORY);
    IRegistry constant public BURNER_ROUTER_FACTORY = IRegistry(SymbioticMainnetConsts.BURNER_ROUTER_FACTORY);

    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public NETWORK = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    uint256 constant public SLASH_PERIOD_SECONDS = 1 days;
    address constant public SLASH_ORACLE = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    address constant public SLASH_RECEIVER = MainnetConstants.COMMITMENT_HOLDINGS_MULTISIG;
    uint256 constant public MIN_BURNER_ROUTER_DELAY = 2 days;

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        deployMevCommitMiddleware(
            NETWORK_REGISTRY, 
            OPERATOR_REGISTRY, 
            VAULT_FACTORY, 
            DELEGATOR_FACTORY,
            SLASHER_FACTORY,
            BURNER_ROUTER_FACTORY,
            NETWORK, 
            SLASH_PERIOD_SECONDS, 
            SLASH_ORACLE, 
            SLASH_RECEIVER,
            MIN_BURNER_ROUTER_DELAY,
            OWNER
        );

        vm.stopBroadcast();
    }
}

contract DeployHolesky is BaseDeploy {

    IRegistry constant public NETWORK_REGISTRY = IRegistry(SymbioticHoleskyDevnetConsts.NETWORK_REGISTRY);
    IRegistry constant public OPERATOR_REGISTRY = IRegistry(SymbioticHoleskyDevnetConsts.OPERATOR_REGISTRY);
    IRegistry constant public VAULT_FACTORY = IRegistry(SymbioticHoleskyDevnetConsts.VAULT_FACTORY);
    IRegistry constant public DELEGATOR_FACTORY = IRegistry(SymbioticHoleskyDevnetConsts.DELEGATOR_FACTORY);
    IRegistry constant public SLASHER_FACTORY = IRegistry(SymbioticHoleskyDevnetConsts.SLASHER_FACTORY);
    IRegistry constant public BURNER_ROUTER_FACTORY = IRegistry(SymbioticHoleskyDevnetConsts.BURNER_ROUTER_FACTORY);

    // On Holesky, use dev keystore account. On mainnet these will be the primev multisig.
    address constant public EXPECTED_MSG_SENDER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;
    address constant public OWNER = EXPECTED_MSG_SENDER;
    address constant public NETWORK = EXPECTED_MSG_SENDER;
    address constant public SLASH_ORACLE = EXPECTED_MSG_SENDER; // Temporary placeholder until oracle implements slashing.
    address constant public SLASH_RECEIVER = EXPECTED_MSG_SENDER; 
    uint256 constant public MIN_BURNER_ROUTER_DELAY = 2 days;

    uint96 constant public SUBNETWORK_ID = 1;
    uint256 constant public VAULT1_MAX_NETWORK_LIMIT = 100000 ether;
    uint256 constant public SLASH_PERIOD_SECONDS = 1 days; // compiles to seconds

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");
        require(msg.sender == EXPECTED_MSG_SENDER, "incorrect msg.sender");

        vm.startBroadcast();

        INetworkRegistry networkRegistry = INetworkRegistry(address(NETWORK_REGISTRY));
        if (!networkRegistry.isEntity(NETWORK)) {
            networkRegistry.registerNetwork();
        }

        address mevCommitMiddlewareProxy = deployMevCommitMiddleware(
            NETWORK_REGISTRY, 
            OPERATOR_REGISTRY, 
            VAULT_FACTORY, 
            DELEGATOR_FACTORY,
            SLASHER_FACTORY,
            BURNER_ROUTER_FACTORY,
            NETWORK, 
            SLASH_PERIOD_SECONDS, 
            SLASH_ORACLE, 
            SLASH_RECEIVER,
            MIN_BURNER_ROUTER_DELAY,
            OWNER
        );

        INetworkMiddlewareService networkMiddlewareService = INetworkMiddlewareService(address(SymbioticHoleskyDevnetConsts.NETWORK_MIDDLEWARE_SERVICE));
        if (networkMiddlewareService.middleware(msg.sender) != address(0)) {
            console.log("WARNING: overwriting existing middleware registration for network:", msg.sender);
        }
        networkMiddlewareService.setMiddleware(mevCommitMiddlewareProxy);

        IBaseDelegator vault1Delegator = IBaseDelegator(address(SymbioticHoleskyDevnetConsts.VAULT_1_DELEGATOR));
        vault1Delegator.setMaxNetworkLimit(SUBNETWORK_ID, VAULT1_MAX_NETWORK_LIMIT);
        console.log("Vault1 delegator max network limit set to:", VAULT1_MAX_NETWORK_LIMIT);

        MevCommitMiddleware mevCommitMiddleware = MevCommitMiddleware(payable(mevCommitMiddlewareProxy));
        address[] memory vaults = new address[](1);
        vaults[0] = SymbioticHoleskyDevnetConsts.VAULT_1;
        uint160[] memory slashAmounts = new uint160[](1);
        slashAmounts[0] = 0.0001 ether; 
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);
        console.log("Vault1 (representing wstETH) registered with MevCommitMiddleware with vault addr:",
            address(SymbioticHoleskyDevnetConsts.VAULT_1), "and collateral slash amount:", slashAmounts[0]);
        
        vm.stopBroadcast();
    }
}
