# Decode Mode - Retroactive ABI Decoding

## Overview

The indexer now supports a **decode mode** that allows you to retroactively decode transaction inputs and event logs for existing records in the database. This is useful when you've imported data from external ETL sources that don't have the `decoded` field populated.

## Use Case

If you have:
- Transactions in the database with `input` data but empty/missing `decoded` field
- Logs in the database with `topics` and `data` but empty/missing `decoded` field
- Contract ABIs available for the contracts you want to decode

Then you can use decode mode to populate the decoded fields without re-fetching data from the blockchain.

## Features

- **Selective Decoding**: Choose to decode only transaction inputs, only logs, or both
- **Batch Processing**: Efficient batch updates with configurable batch size
- **Progress Tracking**: Real-time progress reporting with ETA
- **Dry Run Mode**: Preview what would be decoded without making changes
- **Production Safe**: Minimal changes to existing codebase, reuses existing ABI logic
- **Performance Optimized**: Batch database updates for high throughput

## Usage

### Basic Usage - Decode Both Inputs and Logs

```bash
./indexer \
  --mode=decode \
  --decode-input \
  --decode-logs \
  --dsn="root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true" \
  --abi-config="./contracts-abi/manifest.json" \
  --abi-dir="./contracts-abi/abi" \
  --batch-size=2000
```

### Cross-Database Decoding (ABIs in one DB, data in another)

If your ABIs are in `mevcommit_57173` but you need to decode data in `mev_commit_8855`:

```bash
./indexer \
  --mode=decode \
  --decode-input \
  --decode-logs \
  --dsn="root:@tcp(127.0.0.1:9030)/mev_commit_8855?parseTime=true&interpolateParams=true" \
  --abi-dsn="root:@tcp(127.0.0.1:9030)/mevcommit_57173?parseTime=true&interpolateParams=true" \
  --abi-config="./contracts-abi/manifest.json" \
  --batch-size=2000
```

### Decode Only Transaction Inputs

```bash
./indexer \
  --mode=decode \
  --decode-input \
  --dsn="root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true" \
  --abi-config="./contracts-abi/manifest.json" \
  --batch-size=1000
```

### Decode Only Event Logs

```bash
./indexer \
  --mode=decode \
  --decode-logs \
  --dsn="root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true" \
  --abi-config="./contracts-abi/manifest.json" \
  --batch-size=1000
```

### Dry Run - Preview Without Changes

```bash
./indexer \
  --mode=decode \
  --decode-input \
  --decode-logs \
  --dry-run \
  --dsn="root:@tcp(127.0.0.1:9030)/mevcommit?parseTime=true&interpolateParams=true" \
  --abi-config="./contracts-abi/manifest.json" \
  --batch-size=100
```

## Required Flags

When using `--mode=decode`, you must specify:

1. **`--abi-config`**: Path to ABI manifest JSON file
2. **At least one of**:
   - `--decode-input`: Decode transaction inputs
   - `--decode-logs`: Decode event logs

## Optional Flags

- **`--batch-size`**: Number of records to process in each batch (default: 50, recommended: 1000-5000)
- **`--dry-run`**: Preview mode - shows what would be decoded without updating the database
- **`--dsn`**: Database connection string (default: local StarRocks)
- **`--abi-dsn`**: Separate database for ABIs (optional, useful when ABIs are in different database)
- **`--abi-dir`**: Directory containing ABI files (default: "./contracts-abi/abi")

## Performance Optimizations

The decoder is highly optimized for production use:

1. **Smart Filtering**: Only queries transactions/logs for contracts that have ABIs loaded
   - Uses `WHERE LOWER(to_address) IN (contract_addresses)` for transactions
   - Uses `WHERE LOWER(address) IN (contract_addresses)` for logs
   - Avoids processing millions of records that can't be decoded

2. **Efficient Batching**: Uses `LIMIT/OFFSET` pagination to avoid loading entire table into memory
   - Processes `batch-size` records at a time
   - Each batch is independently committed to database

3. **Batch Updates**: Groups database updates into single transaction per batch
   - Uses prepared statements for efficiency
   - Commits all updates in batch atomically

4. **ABI Caching**: Loads all ABIs once at startup into memory
   - No repeated database queries for ABIs during decode
   - Fast lookup by contract address

## How It Works

### Transaction Input Decoding

1. Selects transactions where:
   - `to_address IS NOT NULL` (contract calls)
   - `decoded IS NULL OR decoded = '{}' OR decoded = ''` (not already decoded)

2. For each transaction:
   - Looks up the ABI for the `to_address` contract
   - Decodes the `input` using the contract's ABI
   - Updates the `decoded` field with JSON containing method name, signature, and arguments

3. Skips transactions where:
   - No ABI is available for the contract
   - Input cannot be decoded (unknown method or invalid data)

### Event Log Decoding

1. Selects logs where:
   - `decoded IS NULL OR decoded = '{}' OR decoded = ''` (not already decoded)

2. For each log:
   - Looks up the ABI for the log's `address` (contract address)
   - Decodes the log using the event signature from `topics[0]`
   - Updates the `decoded` field with JSON containing event name, signature, and arguments

3. Skips logs where:
   - No ABI is available for the contract
   - Log cannot be decoded (unknown event or invalid data)
   - Log has no topics (anonymous logs)

## Output Format

The decoder produces structured JSON in the `decoded` field:

### Transaction Input Example

```json
{
  "name": "transfer",
  "sig": "transfer(address,uint256)",
  "args": {
    "to": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "amount": "1000000000000000000"
  }
}
```

### Event Log Example

```json
{
  "name": "Transfer",
  "sig": "Transfer(address,address,uint256)",
  "args": {
    "from": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "to": "0x1234567890123456789012345678901234567890",
    "value": "1000000000000000000"
  }
}
```

## Performance

The decoder processes records in configurable batches with efficient database updates:

- **Transaction inputs**: ~1000-5000 tx/sec (depends on batch size and database)
- **Event logs**: ~2000-10000 logs/sec (depends on batch size and database)

### Performance Tuning

- **Increase `--batch-size`**: Larger batches = better throughput (try 2000-5000)
- **Database load**: Run during off-peak hours for production databases
- **Network**: Run on the same network as the database for lower latency

## Progress Monitoring

The decoder provides real-time progress updates:

```
level=INFO msg="Processed transaction batch" batch_size=1000 batch_decoded=847 batch_failed=153
  batch_duration=1.2s total_processed=15000 total_decoded=12705 total_failed=2295
  progress_pct=15.00% eta=1m30s
```

- **batch_decoded**: Successfully decoded in this batch
- **batch_failed**: Could not decode (no ABI or unknown method/event)
- **progress_pct**: Overall progress percentage
- **eta**: Estimated time to completion

## Error Handling

The decoder is designed to be safe and robust:

- **Graceful degradation**: Skips records it cannot decode, continues processing
- **Transaction safety**: Uses database transactions for batch updates
- **Signal handling**: Gracefully handles Ctrl+C (SIGINT/SIGTERM)
- **Idempotent**: Safe to run multiple times (only decodes empty fields)

## Building

```bash
cd /Users/kant/mev-commit/tools
go build -o indexer/cmd/chain/indexer ./indexer/cmd/chain
```

## Production Considerations

1. **Test first**: Always use `--dry-run` first to verify the results
2. **Start small**: Use a small `--batch-size` initially (100-500) to test
3. **Monitor resources**: Watch database CPU/memory during decoding
4. **Backup**: Take a database backup before large decode operations
5. **Off-peak hours**: Run during low-traffic periods for production databases
6. **Incremental**: You can stop and restart - it only processes records with empty `decoded` fields

## Technical Details

### Files Modified

- **`indexer.go`**: Added decode mode flag and execution logic (lines 157-398)
- **`decode_existing.go`**: New file with transaction and log decoding functions

### Shared Components

The decode mode reuses the existing production-tested components:
- `decodeTxInput()`: Decodes transaction input using ABI
- `decodeLog()`: Decodes event logs using ABI
- `getParsedABI()`: Retrieves cached ABI for contract address
- `loadABIs()`: Loads ABIs from manifest file
- `abiCache`: In-memory ABI cache

This ensures consistency between real-time indexing and retroactive decoding.

## Troubleshooting

### "No ABI found for contract"

Make sure your `--abi-config` manifest includes all contracts you want to decode.

### "Failed to decode"

Some transactions/logs may not decode if:
- The method/event signature is not in the ABI
- The input/log data is malformed
- The contract uses a non-standard ABI

These are expected and will be logged but won't stop processing.

### Slow performance

- Increase `--batch-size` to 2000-5000
- Ensure database has adequate resources
- Check network latency to database
- Run on the same host/network as the database

## Example Run

```bash
# Dry run to preview
./indexer --mode=decode --decode-input --decode-logs --dry-run \
  --abi-config=./contracts-abi/manifest.json --batch-size=100

# Actual decode with larger batches
./indexer --mode=decode --decode-input --decode-logs \
  --abi-config=./contracts-abi/manifest.json --batch-size=2000
```
