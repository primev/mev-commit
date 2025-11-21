# Multi-ABI Support Implementation

## Summary

The indexer now supports multiple ABI versions per contract address. This allows decoding of transactions/logs from upgraded contracts where different blocks used different contract implementations.

## Changes Made

### 1. Database Schema Change

**Old schema:**
```sql
PRIMARY KEY(address)
```

**New schema:**
```sql
CREATE TABLE contract_abis (
    abi_hash VARCHAR(64) PRIMARY KEY,
    address VARCHAR(255),
    name VARCHAR(255),
    version VARCHAR(50),
    abi JSON,
    INDEX idx_address (address)
)
```

### 2. Code Changes

**Files Modified:**
- `indexer.go` - Updated table schema, changed `abiCache` from `map[string]abi.ABI` to `map[string][]abi.ABI`
- `decode_existing.go` - Updated to use multi-ABI functions
- `sql_processor.go` - Updated to use multi-ABI functions
- `kafka_processor.go` - Updated to use multi-ABI functions

**Key Functions Updated:**
- `getParsedABI()` - Now returns `[]abi.ABI` instead of single `abi.ABI`
- `decodeTxInput()` - Now accepts `[]abi.ABI` and tries each until one succeeds
- `decodeLog()` - Now accepts `[]abi.ABI` and tries each until one succeeds

### 3. How It Works

1. **Loading ABIs**: Multiple ABIs for the same address are loaded into `abiCache[address]` as a slice
2. **Decoding**: When decoding a transaction or log:
   - Get all ABIs for the contract address
   - Try each ABI in sequence
   - Return the first successful decode
   - Only log error if ALL ABIs fail

## Next Steps

### Step 1: Update the Database Schema

```sql
-- Connect to your database
mysql -h 127.0.0.1 -P 9030 -u root -D mev_commit_8855

-- Backup existing data
CREATE TABLE contract_abis_backup AS SELECT * FROM contract_abis;

-- Drop old table
DROP TABLE contract_abis;

-- Create new table
CREATE TABLE contract_abis (
    abi_hash VARCHAR(64) PRIMARY KEY,
    address VARCHAR(255),
    name VARCHAR(255),
    version VARCHAR(50),
    abi JSON,
    INDEX idx_address (address)
)
ENGINE=olap
DISTRIBUTED BY HASH(abi_hash) BUCKETS 1
PROPERTIES("replication_num"="1");
```

### Step 2: Insert Multiple ABI Versions

For each contract that has multiple versions, insert all versions:

```sql
-- Example for BlockTracker
-- Version 0.8
INSERT INTO contract_abis (abi_hash, address, name, version, abi)
VALUES (
  MD5('<full_abi_json_here>'),
  '0x0da2a367c51f2a34465acd6ae5d8a48385e9cb03',
  'BlockTracker',
  'v0.8',
  parse_json('<full_abi_json_here>')
);

-- Version 1.1.0
INSERT INTO contract_abis (abi_hash, address, name, version, abi)
VALUES (
  MD5('<full_abi_json_here>'),
  '0x0da2a367c51f2a34465acd6ae5d8a48385e9cb03',
  'BlockTracker',
  'v1.1.0',
  parse_json('<full_abi_json_here>')
);
```

**Important Notes:**
- Use `MD5()` function to calculate the hash from the ABI JSON string
- The same ABI can be inserted multiple times for different addresses (e.g., if deployed at multiple addresses)
- The `abi_hash` ensures we don't duplicate the same ABI content

### Step 3: Insert ABIs for All Contracts

You mentioned you have:
- BlockTracker: v0.8 and v1.1.0
- BidderRegistry: v0.8 and v1.1.0
- PreconfManager: v0.8 and v1.1.0 (need to verify)
- Oracle: v0.8 and v1.1.0 (need to verify)
- ProviderRegistry: v0.8 and v1.1.0 (need to verify)
- SettlementGateway: v0.8 and v1.1.0 (need to verify)

For each contract, insert both v0.8 and v1.1.0 ABIs.

### Step 4: Rebuild and Redeploy

The binary has already been rebuilt at `/Users/kant/mev-commit/tools/indexer/cmd/chain/indexer`.

You can now redeploy the decoder pod:

```bash
# Build docker image (if needed)
# Deploy with your existing helm chart

# Or test locally first
./indexer --mode=decode --decode-input --decode-logs --dry-run \
  --dsn="root:password@tcp(127.0.0.1:9030)/mev_commit_8855?parseTime=true&interpolateParams=true" \
  --abi-config="./manifest.json" \
  --batch-size=100
```

## Expected Results

After inserting multiple ABI versions:

**Before:**
- Success rate: 0.5% (12,663 out of 2,478,000 logs)
- Errors: "method signature not found" / "event signature not found"

**After:**
- Success rate: Expected 90%+
- Only logs errors when ALL ABI versions fail to decode
- Automatically handles contract upgrades

## Verification

Check logs for successful decoding:

```bash
kubectl logs <pod-name> | grep "batch_decoded"
```

You should see much higher `batch_decoded` numbers and lower `batch_failed` counts.

## Troubleshooting

### If decode rate is still low:

1. **Verify ABIs are loaded:**
   ```sql
   SELECT address, name, version, COUNT(*) as abi_count
   FROM contract_abis
   GROUP BY address, name, version;
   ```

2. **Check for missing ABI versions:**
   - Look at error logs for specific method/event signatures that fail
   - Find which ABI version contains those signatures
   - Insert missing ABI versions

3. **Verify abi_hash uniqueness:**
   ```sql
   SELECT abi_hash, COUNT(*) as count
   FROM contract_abis
   GROUP BY abi_hash
   HAVING count > 1;
   ```
   Should return empty (each hash should be unique)

## Notes

- The manifest.json file is NOT updated - it only loads ABIs once at startup
- The database is the source of truth for ABIs in decode mode
- You can keep using the same manifest.json in your helm deployment
- The code will automatically try all ABIs for each address
