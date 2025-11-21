# Timestamp Fix Instructions

## Overview
Two Kubernetes Jobs to fix incorrect timestamps in the blocks and logs tables.

## Prerequisites
- kubectl configured with access to the cluster
- kubeconfig: `/Users/kant/kubeconfig-prod`

## Step 1: Fix Blocks Table

Apply the first job to fetch correct timestamps from RPC and update the blocks table:

```bash
export KUBECONFIG=/Users/kant/kubeconfig-prod
kubectl apply -f fix-timestamps-job.yaml
```

**Monitor progress:**
```bash
kubectl logs -f job/fix-timestamps
```

**Check status:**
```bash
kubectl get job fix-timestamps
```

The job will:
- Fetch timestamps from https://chainrpc-v1.mev-commit.xyz
- Use 20 concurrent workers
- Process all blocks in the blocks table
- Log progress every 1000 blocks

**Expected duration:** Depends on block count, typically runs at 50-100 blocks/sec with 20 workers.

## Step 2: Fix Logs Table

**IMPORTANT:** Only run this AFTER the blocks fix completes successfully!

Apply the second job to update logs from the corrected blocks table:

```bash
kubectl apply -f fix-log-timestamps-job.yaml
```

**Monitor progress:**
```bash
kubectl logs -f job/fix-log-timestamps
```

The job will:
- Use a single SQL JOIN UPDATE statement
- Update all logs.block_timestamp from blocks.timestamp
- Complete very quickly (single SQL operation)

## Cleanup

After both jobs complete successfully:

```bash
kubectl delete job fix-timestamps fix-log-timestamps
```

## Troubleshooting

**If blocks fix fails:**
```bash
# Check logs
kubectl logs job/fix-timestamps

# Delete and retry
kubectl delete job fix-timestamps
kubectl apply -f fix-timestamps-job.yaml
```

**If logs fix fails:**
```bash
# Check if blocks were fixed first
kubectl logs job/fix-timestamps | tail -20

# Delete and retry logs fix
kubectl delete job fix-log-timestamps
kubectl apply -f fix-log-timestamps-job.yaml
```

## Test Mode (Optional)

To test without making changes, edit the YAML files and add `--dry-run` flag:

For blocks:
```yaml
./fix-timestamps \
  --dsn="..." \
  --workers=20 \
  --batch-size=100 \
  --dry-run  # Add this line
```

For logs:
```yaml
./fix-log-timestamps \
  --dsn="..." \
  --dry-run  # Add this line
```

## Configuration

Both jobs are pre-configured with the correct DSN. No additional configuration needed.

**Blocks fix settings:**
- Workers: 20
- Batch size: 100
- RPC: https://chainrpc-v1.mev-commit.xyz

**Logs fix settings:**
- Single UPDATE statement with JOIN
- No chunking needed (fast operation)
