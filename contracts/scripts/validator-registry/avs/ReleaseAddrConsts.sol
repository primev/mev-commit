// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

/// @notice Constants from https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#deployments,
/// @notice Last updated 11-07-2024
library EigenMainnetReleaseConsts {

    address internal constant DELEGATION_MANAGER = 0x39053D51B77DC0d36036Fc1fCc8Cb819df8Ef37A;
    address internal constant STRATEGY_MANAGER = 0x858646372CC42E1A627fcE94aa7A7033e7CF075A;
    address internal constant EIGENPOD_MANAGER = 0x91E677b07F7AF907ec9a428aafA9fc14a0d3A338;
    address internal constant AVS_DIRECTORY = 0x135DDa560e946695d6f155dACaFC6f1F25C1F5AF;
    address internal constant SLASHER = 0xD92145c07f8Ed1D392c1B88017934E301CC1c3Cd;
    address internal constant REWARDS_COORDINATOR = 0x7750d328b314EfFa365A0402CcfD489B80B0adda;

    address internal constant STRATEGY_BASE_CBETH = 0x54945180dB7943c0ed0FEE7EdaB2Bd24620256bc;
    address internal constant STRATEGY_BASE_STETH = 0x93c4b944D05dfe6df7645A86cd2206016c51564D;
    address internal constant STRATEGY_BASE_RETH = 0x1BeE69b7dFFfA4E2d53C2a2Df135C388AD25dCD2;
    address internal constant STRATEGY_BASE_ETHX = 0x9d7eD45EE2E8FC5482fa2428f15C971e6369011d;
    address internal constant STRATEGY_BASE_ANKRETH = 0x13760F50a9d7377e4F20CB8CF9e4c26586c658ff;
    address internal constant STRATEGY_BASE_OETH = 0xa4C637e0F704745D182e4D38cAb7E7485321d059;
    address internal constant STRATEGY_BASE_OSETH = 0x57ba429517c3473B6d34CA9aCd56c0e735b94c02;
    address internal constant STRATEGY_BASE_SWETH = 0x0Fe4F44beE93503346A3Ac9EE5A26b130a5796d6;
    address internal constant STRATEGY_BASE_WBETH = 0x7CA911E83dabf90C90dD3De5411a10F1A6112184;
    address internal constant STRATEGY_BASE_SFRXETH = 0x8CA7A5d6f3acd3A7A8bC468a8CD0FB14B6BD28b6;
    address internal constant STRATEGY_BASE_LSETH = 0xAe60d8180437b5C34bB956822ac2710972584473;
    address internal constant STRATEGY_BASE_METH = 0x298aFB19A105D59E74658C4C334Ff360BadE6dd2;
    address internal constant BEACON_CHAIN_ETH = 0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0;
}

/// @notice Constants from https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-testnet-deployment
/// @notice Last updated 07-26-2024
library EigenHoleskyReleaseConsts {

    address internal constant DELEGATION_MANAGER = 0xA44151489861Fe9e3055d95adC98FbD462B948e7;
    address internal constant STRATEGY_MANAGER = 0xdfB5f6CE42aAA7830E94ECFCcAd411beF4d4D5b6;
    address internal constant EIGENPOD_MANAGER = 0x30770d7E3e71112d7A6b7259542D1f680a70e315;
    address internal constant AVS_DIRECTORY = 0x055733000064333CaDDbC92763c58BF0192fFeBf;
    address internal constant SLASHER = 0xcAe751b75833ef09627549868A04E32679386e7C;
    address internal constant REWARDS_COORDINATOR = 0xAcc1fb458a1317E886dB376Fc8141540537E68fE;

    address internal constant STRATEGY_BASE_STETH = 0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3;
    address internal constant STRATEGY_BASE_RETH = 0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0;
    address internal constant STRATEGY_BASE_WETH = 0x80528D6e9A2BAbFc766965E0E26d5aB08D9CFaF9;
    address internal constant STRATEGY_BASE_LSETH = 0x05037A81BD7B4C9E0F7B430f1F2A22c31a2FD943;
    address internal constant STRATEGY_BASE_SFRXETH = 0x9281ff96637710Cd9A5CAcce9c6FAD8C9F54631c;
    address internal constant STRATEGY_BASE_ETHX = 0x31B6F59e1627cEfC9fA174aD03859fC337666af7;
    address internal constant STRATEGY_BASE_OSETH = 0x46281E3B7fDcACdBa44CADf069a94a588Fd4C6Ef;
    address internal constant STRATEGY_BASE_CBETH = 0x70EB4D3c164a6B4A5f908D4FBb5a9cAfFb66bAB6;
    address internal constant STRATEGY_BASE_METH = 0xaccc5A86732BE85b5012e8614AF237801636F8e5;
    address internal constant STRATEGY_BASE_ANKRETH = 0x7673a47463F80c6a3553Db9E54c8cDcd5313d0ac;
    address internal constant BEACON_CHAIN_ETH = 0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0;
}

/// @notice Constants from https://github.com/Layr-Labs/eigenlayer-contracts?tab=readme-ov-file#current-testnet-deployment
/// @notice Last updated 07-25-2025 — for HOODI testnet
library EigenHoodiReleaseConsts {
    // Core
    address internal constant DELEGATION_MANAGER     = 0x867837a9722C512e0862d8c2E15b8bE220E8b87d;
    address internal constant STRATEGY_MANAGER       = 0xeE45e76ddbEDdA2918b8C7E3035cd37Eab3b5D41;
    address internal constant EIGENPOD_MANAGER       = 0xcd1442415Fc5C29Aa848A49d2e232720BE07976c;
    address internal constant AVS_DIRECTORY          = 0xD58f6844f79eB1fbd9f7091d05f7cb30d3363926;
    address internal constant REWARDS_COORDINATOR    = 0x29e8572678e0c272350aa0b4B8f304E47EBcd5e7;
    address internal constant ALLOCATION_MANAGER     = 0x95a7431400F362F3647a69535C5666cA0133CAA0;
    address internal constant PERMISSION_CONTROLLER  = 0xdcCF401fD121d8C542E96BC1d0078884422aFAD2;

    // Strategies - Deployed via StrategyFactory
    address internal constant STRATEGY_FACTORY       = 0xfB7d94501E4d4ACC264833Ef4ede70a11517422B;
    address internal constant STRATEGY_BASE          = 0x6d28cEC1659BC3a9BC814c3EFc1412878B406579;
    address internal constant STRATEGY_BASE_STETH    = 0xF8a1a66130D614c7360e868576D5E59203475FE0;
    address internal constant STRATEGY_BASE_WETH     = 0x24579aD4fe83aC53546E5c2D3dF5F85D6383420d;
    // Special strategies
    address internal constant STRATEGY_BASE_EIGEN    = 0xB27b10291DBFE6576d17afF3e251c954Ae14f1D3;
    // Beacon Chain ETH placeholder (not a real contract)
    address internal constant BEACON_CHAIN_ETH       = 0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0;
    // Tokens
    address internal constant EIGEN_TOKEN            = 0x8ae2520954db7D80D66835cB71E692835bbA45bf;
    address internal constant BACKING_EIGEN          = 0x6e60888132Cc7e637488379B4B40c42b3751f63a;
}
