// SPDX-License-Identifier: BSL 1.1
// solhint-disable one-contract-per-file
pragma solidity 0.8.29;

/// @notice Relevant constants communicated by Symbiotic team as canonical mainnet addresses.
/// @notice Last updated 1-14-2025
library SymbioticMainnetConsts {
    address internal constant NETWORK_REGISTRY = 0xC773b1011461e7314CF05f97d95aa8e92C1Fd8aA;
    address internal constant OPERATOR_REGISTRY = 0xAd817a6Bc954F678451A71363f04150FDD81Af9F;
    address internal constant VAULT_FACTORY = 0xAEb6bdd95c502390db8f52c8909F703E9Af6a346;
    address internal constant DELEGATOR_FACTORY = 0x985Ed57AF9D475f1d83c1c1c8826A0E5A34E8C7B;
    address internal constant SLASHER_FACTORY = 0x685c2eD7D59814d2a597409058Ee7a92F21e48Fd;
    address internal constant BURNER_ROUTER_FACTORY = 0x99F2B89fB3C363fBafD8d826E5AA77b28bAB70a0;
}

/// @notice Constants from https://docs.symbiotic.fi/deployments/current
/// @notice Last updated 11-30-2024
library SymbioticHoleskyDevnetConsts {
    address internal constant VAULT_FACTORY = 0x407A039D94948484D356eFB765b3c74382A050B4;
    address internal constant DELEGATOR_FACTORY = 0x890CA3f95E0f40a79885B7400926544B2214B03f;
    address internal constant SLASHER_FACTORY = 0xbf34bf75bb779c383267736c53a4ae86ac7bB299;
    address internal constant NETWORK_REGISTRY = 0x7d03b7343BF8d5cEC7C0C27ecE084a20113D15C9;
    address internal constant NETWORK_METADATA_SERVICE = 0x0F7E58Cc4eA615E8B8BEB080dF8B8FDB63C21496;
    address internal constant NETWORK_MIDDLEWARE_SERVICE = 0x62a1ddfD86b4c1636759d9286D3A0EC722D086e3;
    address internal constant OPERATOR_REGISTRY = 0x6F75a4ffF97326A00e52662d82EA4FdE86a2C548;
    address internal constant OPERATOR_METADATA_SERVICE = 0x0999048aB8eeAfa053bF8581D4Aa451ab45755c9;
    address internal constant VAULT_OPT_IN_SERVICE = 0x95CC0a052ae33941877c9619835A233D21D57351;
    address internal constant NETWORK_OPT_IN_SERVICE = 0x58973d16FFA900D11fC22e5e2B6840d9f7e13401;
    address internal constant VAULT_CONFIGURATOR = 0xD2191FE92987171691d552C219b8caEf186eb9cA;
    address internal constant DEFAULT_STAKER_REWARDS_FACTORY = 0x698C36DE44D73AEfa3F0Ce3c0255A8667bdE7cFD;
    address internal constant DEFAULT_OPERATOR_REWARDS_FACTORY = 0x00055dee9933F578340db42AA978b9c8B25640f6;

    address internal constant VAULT_1 = 0xd88dDf98fE4d161a66FB836bee4Ca469eb0E4a75;
    address internal constant VAULT_1_DELEGATOR = 0x85CF967A8DDFAf8C0DFB9c75d9E92a3C785A6532;
    address internal constant VAULT_1_SLASHER = 0x57e5Fb61981fa1b43a074B2aeb47CCF157b19223;

    address internal constant BURNER_ROUTER_FACTORY = 0x32e2AfbdAffB1e675898ABA75868d92eE1E68f3b;
}
