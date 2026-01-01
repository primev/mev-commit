# RPC_tx_insert

## Overview
This tool:
- Discovers L1 transaction hashes from Primev sources (event logs and mc transaction table).
- Fetches transaction data from Covalent (`transaction_v2`) for each hash.
- Computes:
  - `l1_timestamp`, `from_address`, `to_address`
  - `total_vol_eth`, `eth_vol`, `weth_vol`, `token_vol_eth`, `swap_vol_eth`
  - `is_swap`, `is_lending`, `is_transfer`, `is_approval`
  - `primary_class`, `protocol`
- Writes results into `mevcommit_57173.processed_l1_txns_v2`.

Write behavior:
- Inserts missing rows (event-backed + rpc-only backfill).
- Updates existing rows that are incomplete (fill-only by default).
- Optional full overwrite with `-recompute-all`.

## Data sources
Candidate discovery:
1) Event-backed candidates:
   - `OpenedCommitmentStored` joined to `CommitmentProcessed` from `mevcommit_57173.tx_view`.
2) RPC-only backfill candidates:
   - `pg_mev_commit_fastrpc.public.mctransactions_sr` where `status in ('confirmed','pre-confirmed')`,
     excluding hashes already present in `OpenedCommitmentStored` and excluding hashes already in v2.

On-chain enrichment:
- Covalent `transaction_v2` endpoint (ETH mainnet).

## Tombstone behavior for missing transactions (Covalent 404)
If Covalent returns HTTP 404 with a "Transaction hash ... not found" message:
- The tool can insert a placeholder row into `processed_l1_txns_v2` with:
  - `l1_tx_hash = <hash_norm>` (no 0x)
  - `primary_class = 'not_found'`
  - other computed fields left NULL
- This prevents repeated retries of transactions that never land on-chain.

Age gate (anti-false-tombstone):
- Tombstone insertion is gated on block age using `mctransactions_sr.block_number`.
- The tool queries StarRocks to get:
  - `head_block = MAX(block_number)` over `mctransactions_sr` for `confirmed/pre-confirmed`
  - `tx_block = block_number` for the specific hash (same table)
- Tombstone insertion happens only if:
  - `head_block - tx_block > 75`
- If the block age is <= 75 (recent), the tool skips tombstoning and will retry on later runs.

## Required environment variables
- `DB_USER`, `DB_PW`, `DB_HOST`, `DB_PORT`, `DB_NAME`  (StarRocks via MySQL protocol)
- `COVALENT_KEY`                                       (Covalent API key)

## Commands

Recommended (preview without DB writes):
- `go run . -dry-run -limit 200`

Recommended (write to DB):
- `go run . -limit 200`

Single-tx debug (no DB writes):
- `go run . -tx 0x<hash>`

Force recompute/overwrite existing non-null values:
- `go run . -recompute-all -limit 500`

Only inserts (skip updating existing incomplete rows):
- `go run . -only-inserts -limit 500`

Only updates (skip discovering/inserting missing txs):
- `go run . -only-updates -limit 500`

## Flags (summary)
- `-limit N`: cap number of txs processed (0 = no limit)
- `-dry-run`: compute and print discrepancy summary; no inserts/updates
- `-tx 0x...`: single transaction compute; prints JSON; no DB writes
- `-recompute-all`: overwrite computed columns on updates
- `-only-inserts`: only insert missing hashes; do not update existing rows
- `-only-updates`: only update existing incomplete rows; do not discover/insert missing hashes
- `-print-sample N`: log N sample hashes for insert/update sets
- `-only-old-lending`: restrict updates to rows where existing `is_lending=1`
- `-compare-only-old-swapvol-gt0`: dry-run comparison filter for swap volume discrepancies

## Notes
- `processed_l1_txns_v2.l1_tx_hash` is stored as `hash_norm` (no `0x` prefix).
- `loadExistingNeedingUpdate` excludes rows where `primary_class='not_found'` to avoid reprocessing tombstones.
