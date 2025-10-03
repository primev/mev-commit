// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {ProviderRegistry} from "../../contracts/core/ProviderRegistry.sol";
import {ProviderRegistryV2} from "../../contracts/core/ProviderRegistryV2.sol";
import {console} from "forge-std/console.sol";

contract UpgradeProviderRegistry is Script {
    
    function run() external {
        vm.startBroadcast();
        address proxyAddress = vm.envAddress("PROVIDER_REGISTRY_PROXY");
        console.log("Upgrading ProviderRegistry proxy at:", proxyAddress);
        Upgrades.upgradeProxy(proxyAddress, "ProviderRegistryV2.sol", "");
        console.log("ProviderRegistry successfully upgraded to:", "ProviderRegistryV2.sol");
        vm.stopBroadcast();
    }
}

contract ResetPubkeys is Script {
    
    function run() external {
        vm.startBroadcast();
        address proxyAddress = vm.envAddress("PROVIDER_REGISTRY_PROXY");
        ProviderRegistryV2 providerRegistry = ProviderRegistryV2(payable(proxyAddress));
        
        bytes[] memory pubkeys = new bytes[](33);
        pubkeys[0] = hex"910f20173e75cc036eb232cf2081043d62c028c3ddae53d5cccdbdd1213a9830c5143c431ee9e401a4d3432f3c23d18c";
        pubkeys[1] = hex"a9d0a0f9059972d775a45d8377768ff20234de91a6fbacba5737ff8803807c38021b5863e5084869e695ef1d6c2fdaef";
        pubkeys[2] = hex"b26f96664274e15fb6fcda862302e47de7e0e2a6687f8349327a9846043e42596ec44af676126e2cacbdd181f548e681";
        pubkeys[3] = hex"a94a5107948363d29e6a7c476f7e2665eaa27d3d92dcba4c68a66de71c07d1286e2755af80150a3f01e39c1fe69c4ac4";
        pubkeys[4] = hex"8216e00e1dc8e15c362ce8083ad01feeb04688dd3a18998a37db1c3c8b641372398c504e1aca2cbddd87c4075482b42f";
        pubkeys[5] = hex"b4a435cf816291596fe2e405651ec8b6c80b9cc34dace3c83202ca489a833756c9a0672ebdc17f23d9d43163db1caa5d";
        pubkeys[6] = hex"95c8cc31f8d4e54eddb0603b8f12d59d466f656f374bde2073e321bdd16082d420e3eef4d62467a7ea6b83818381f742";
        pubkeys[7] = hex"b67eaa5efcfa1d17319c344e1e5167811afbfe7922a2cf01c9a361f465597a5dc3a5472bd98843bac88d2541a78eab08";
        pubkeys[8] = hex"946fcf348bbf1044a3eaa3d27f1d01397cfc4d27495a949447c95ea7ec41db9f5c4fe3b867a937623685a0462996df09";
        pubkeys[9] = hex"b0b0de6ba411630193eed5450f82a76d66501e9838fe5fbe98d4873b5de677b3a20a360302fbe094adae63bc63e4ded4";
        pubkeys[10] = hex"b9ce2753254979122cf137288efe471791893a9cf6abcd192f33515f6ca778507a51bcab84542562efc244001a7b4c55";
        pubkeys[11] = hex"8a85062118e29045e8b199723d9af519b0c071d9e7586fa89733e798536d4637f8545a169012cc1b76005d9453603273";
        pubkeys[12] = hex"8527d16cf01edcea2cbb05e27f5f61b184578ea27c4015e4533e4cebd6d53297bd004c14fe3f247a468361bf781e0069";
        pubkeys[13] = hex"b47963246adef02cd3e61cbb648c04fd99b05e28a616aef3aa7fb688c17b10d1ce9662b61a600efbdd110e93d62d5144";
        pubkeys[14] = hex"8898a80ec199dbe15a57e3ceb51389ded413d6a2ebaaac0330af7effcdd33bba70318339ec9cbc1253b30d463c4999c4";
        pubkeys[15] = hex"807a81a9873359323966feb80fcc52b6049888f885d58134bf52bf9825f89aabc723bb5bacd1d8f4dea6322a6d166535";
        pubkeys[16] = hex"88857150299287cedfbadea1ee3fb7ac121f1e4e16bef44b4a7bad35432973c4009efb90394facca3fdc0759ba70f93f";
        pubkeys[17] = hex"A32AADB23E45595FE4981114A8230128443FD5407D557DC0C158AB93BC2B88939B5A87A84B6863B0D04A4B5A2447F847";
        pubkeys[18] = hex"AE2FFC6986C9A368C5AD2D51F86DB2031D780F6AC9B2348044DEA4E3A75808B566C935099DE8B1A1609DB322F2110E7A";
        pubkeys[19] = hex"8509ECB595DA0EDA2C6FCED4E287F0510A2C2DBA5F80EE930503EF86E268D808A6DF25E397177DA06CD479771CE66840";
        pubkeys[20] = hex"94829E6F7A598A2F2DFDD9E1246D7CFDC30A626666D9419F3C147CC954507E97184C598DC109F4D05C2139C48AF6746C";
        pubkeys[21] = hex"A3523967A7955C0244910F23B7B1FC59636F03BEC437286B622815408D51389F7F6CD54617733B93926B7860E1F6AFB0";
        pubkeys[22] = hex"8226FB149BFE7B4967FFE82ECB9084FFD5BBF0303DE0B88F68FDD8297CDFFE80F611FA27BC05506B4FBA12E2EB5BC5A5";
        pubkeys[23] = hex"AF10542267816E91ADBC8F4A6754765D492534F8325F34A2E89CAA2BA45C7158F6DEAA6E7FB454EBB6F6A1495FE63DBA";
        pubkeys[24] = hex"94A076B27F294DC44B9FD44D8E2B063FB129BC85ED047DA1CEFB82D16E1A13E6B50DE31A86F5B233D1E6BBACA3C69173";
        pubkeys[25] = hex"8B39E8F6AD0A7D2C9E893459D76AC1BC7884D5343324F7639CC883590E8914C2EDB59C3751E4B4C31D466BAEC718D440";
        pubkeys[26] = hex"B255D270445AC3A52A1A97D0D8547EEEE526D649E172663438B621FF9BE4212EEAA425002AF64E4685083E871B0BD7C6";
        pubkeys[27] = hex"A0383AAD83FA40C02CEBBC89EA396AA8545A152E60B558780173CE2B81CED85D8A3858F83AE99CBCA50DDA43CCBEACA9";
        pubkeys[28] = hex"A20E2D356EEA696F4F5F8A02D8865CB3FF287A04765978D399A677198BF2A806B42DC8DBB5B5703D8238A76A6EEB6F6D";
        pubkeys[29] = hex"88A53EC4422F50238DEF5696446E06ACD076F94D73C76E51449CF82EB5312DCB8A845A6C3199CE35DE3F2B0441DEB76B";
        pubkeys[30] = hex"B435DD63A14675EA11E4EACA6BD640C70E68843CF4C8BFE3BBAF3A7EBEDC3FC53D80050E58748EC3D088A15113C6D4C6";
        pubkeys[31] = hex"947DE8EEC641942BEE15F2754F825E24A00DDD135837C5EAAD4921BD1CF18F9EBC374E1194880676DC19CB463DCA842E";
        pubkeys[32] = hex"A7A713E0275E888A06B4A95B5C0B16557CF25FE431269F54BC898287081DC4A25110EDBAFA643EA1B86C4319C8D9977E";

        for (uint256 i = 0; i < pubkeys.length; i++) {
            providerRegistry.overrideRemoveBLSKey(address(this), pubkeys[i]);
        }
        for (uint256 i = 0; i < pubkeys.length; i++) {
            providerRegistry.overrideAddBLSKey(address(this), pubkeys[i]);
        }
        vm.stopBroadcast();
    }
}
