# MEV-Commit L1 Volume Pipeline

This tool builds and maintains the `processed_l1_txns` table in StarRocks. It:

- Finds MEV-Commit commitments that were actually processed on L1
- Fetches their L1 Ethereum transactions via Covalent
- Computes preconfirmation volume in ETH (native, WETH, and ERC-20s marked-to-ETH)
- Writes the results into `processed_l1_txns` for downstream analytics

It is safe to run repeatedly.

## What the Script Does

### 1. Query StarRocks (`tx_view`)
- Selects all `OpenedCommitmentStored` events.
- Finds those that also have a `CommitmentProcessed` event.
- Extracts `commitment_index`, `bidder`, `committer`, and `l1_tx_hash`.

### 2. Determine which commitments need processing
- Selects commitments from `processed_l1_txns` with non-null `total_vol_eth`.
- Skips those already processed.
- Processes all others.

### 3. Fetch & parse L1 transaction data (Covalent)
- Retrieves full tx including logs.
- Computes:
  - ETH transferred
  - WETH movements
  - ERC-20 token transfers + value in ETH using historical pricing
- Extracts L1 timestamp.

### 4. Insert enriched rows into `processed_l1_txns`
- Inserts/upserts:
  - commitment_index
  - l1_timestamp
  - bidder
  - committer
  - l1_tx_hash
  - total_vol_eth
  - eth_vol
  - weth_vol
  - token_vol_eth
- Only modifies this table; never touches others.


## Required Environment Variables
- COVALENT_KEY
- DB_USER
- DB_PW
- DB_HOST
- DB_PORT
- DB_NAME

## Usage

### Fill database
```
go run . -fill-db
```

### Single transaction mode
```
go run . 0xTX_HASH
```

### File mode
```
go run . -file txs.txt
```
