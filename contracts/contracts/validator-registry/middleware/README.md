# Mev-commit Network Middleware Implementation

## Overview

The `MevCommitMiddleware` contracts serve as the representation of network middleware for integration with Symbiotic's L1 restaking protocol. This contract enables restaked assets on L1 to be used as collateral for correct validator behavior within the mev-commit protocol. 

## DRAFT NOTES

* See https://github.com/symbioticfi/cosmos-sdk/blob/main/middleware/src/SimpleMiddleware.sol
* See https://docs.symbiotic.fi/handbooks/Handbook%20for%20Networks
* See https://github.com/symbioticfi/cosmos-sdk
* See Stride's impl?
* Note operator registration will be permissioned and will require the operator to honestly provide its val bls pub keys
* Use existing draft PR as example (impl from few months back including val whitelist)
* Look into whether operators would be disincentivized to provide bunk keys anyways cause they cant control whether that val actually does stuff correctly
* Also note more keys = more points
