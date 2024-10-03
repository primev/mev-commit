// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @notice Constants from https://docs.symbiotic.fi/deployments/
/// @notice Last updated 10-02-2024
library SymbioticHoleskyDevnetConsts {
    address internal constant VAULT_FACTORY = 0x18C659a269a7172eF78BBC19Fe47ad2237Be0590;
    address internal constant DELEGATOR_FACTORY = 0xdE2Ad96117b48bd614A9ed8Ff6bcf5D7eB815596;
    address internal constant SLASHER_FACTORY = 0xCeE813788eFD2edD87B2ABE96EAF4789Dbdb3d7D;
    address internal constant NETWORK_REGISTRY = 0xac5acD8A105C8305fb980734a5AD920b5920106A;
    address internal constant NETWORK_METADATA_SERVICE = 0x56286D7B00fD61AD49f423a07B4944c8E20E2397;
    address internal constant NETWORK_MIDDLEWARE_SERVICE = 0x683F470440964E353b389391CdDDf8df381C282f;
    address internal constant OPERATOR_REGISTRY = 0xAdFC41729fF447974cE27DdFa358A0f2096c3F39;
    address internal constant OPERATOR_METADATA_SERVICE = 0xb1594749596e8b1F5c4cB15f0d0583762A793482;
    address internal constant VAULT_OPT_IN_SERVICE = 0xc105215C23Ed7E45eB6Bf539e52a12c09cD504A5;
    address internal constant NETWORK_OPT_IN_SERVICE = 0xF5AFc9FA3Ca63a07E529DDbB6eae55C665cCa83E;
    address internal constant VAULT_CONFIGURATOR = 0x382e9c6fF81F07A566a8B0A3622dc85c47a891Df;
    address internal constant DEFAULT_STAKER_REWARDS_FACTORY = 0x0798d5931Fc1a807899DE5F46429c149a06e486F;
    address internal constant DEFAULT_OPERATOR_REWARDS_FACTORY = 0xeA9A0522fbC3417fA6d2bC101C49CDf2540DDB64;

    address internal constant VAULT_1 = 0x1df2fbfcD600ADd561013f44B2D055E2e974f605;
    address internal constant VAULT_1_DELEGATOR = 0x70F1450829E114A70409959ba0bF1cc9B6d8Bb67;
    address internal constant VAULT_1_SLASHER = 0x7B4A771DeB69F34dF12750efD345114B26A4271e;
}
