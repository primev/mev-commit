# GasTankDepositor Contract Documentation

## Overview

The `GasTankDepositor` contract coordinates on-demand ETH transfers from user EOAs to an RPC service-managed EOA for custodial gas tank management. This contract enables automatic gas tank top-ups using the ERC-7702 standard, allowing users to delegate smart contract functionality to their EOA addresses without requiring a contract wallet.

## Purpose

The contract facilitates a custodial gas tank system where:
- Users deposit ETH that is transferred to an RPC service-managed EOA
- The RPC service maintains an off-chain ledger tracking each user's custodial balance
- When users transact with the FAST RPC service, their off-chain ledger balance is debited to cover transaction costs
- The RPC service can automatically top up user accounts when balances run low

## Architecture

### Key Components

1. **RPC Service EOA**: A single EOA address managed by the RPC service that receives all gas tank deposits
2. **Off-Chain Ledger**: The RPC service maintains a ledger tracking each user's custodial balance
3. **ERC-7702 Delegation**: Users delegate their EOA to the `GasTankDepositor` contract, enabling smart contract functionality on their EOA address
4. **Minimum Deposit**: Immutable minimum amount that must be transferred in each top-up operation

### How It Works

1. **User Authorization** (One-time setup):
   - User authorizes the `GasTankDepositor` contract using ERC-7702
   - User sends a network transaction to attach the delegation
   - After delegation, the user's EOA can execute contract functions as if it were a smart contract

2. **Initial Funding**:
   - User calls `fundGasTank(uint256 _amount)` with their desired initial deposit
   - Amount must be >= `MAXIMUM_DEPOSIT`
   - ETH is transferred from user's EOA to the RPC service EOA
   - RPC service updates off-chain ledger to reflect the deposit

3. **Automatic Top-Ups**:
   - When a user's off-chain ledger balance drops below threshold, RPC service calls `fundGasTank()`
   - This always transfers exactly `MAXIMUM_DEPOSIT` amount
   - Transfer occurs directly from user's EOA balance (if sufficient funds available)
   - No user interaction required - fully automated
   - No need for `maxTransferAllowance` as RPC is restricted to minimum amount only

4. **Off-Chain Ledger Operations**:
   - RPC service tracks user balances in off-chain ledger
   - When users transact with FAST RPC service, ledger debits their account
   - When balance is low, RPC service triggers automatic top-up
   - All transfers go to the single RPC service EOA


