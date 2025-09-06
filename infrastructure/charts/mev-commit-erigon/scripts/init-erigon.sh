#!/bin/sh

set -e

echo "Initializing Erigon..."

# Check if genesis file exists
if [ ! -f "/shared/genesis.json" ]; then
    echo "Error: Genesis file not found at /shared/genesis.json"
    exit 1
fi

echo "Found genesis file at /shared/genesis.json"

# Check if Erigon data directory already exists and is initialized
if [ -d "${ERIGON_DATADIR}/chaindata" ]; then
    echo "Erigon chaindata already exists at ${ERIGON_DATADIR}/chaindata"
    echo "Skipping initialization to preserve existing data"
    exit 0
fi

echo "Initializing Erigon with genesis file..."

# Create data directory if it doesn't exist
mkdir -p "${ERIGON_DATADIR}"

# Initialize Erigon with genesis
erigon init \
    --datadir="${ERIGON_DATADIR}" \
    /shared/genesis.json

if [ $? -eq 0 ]; then
    echo "Erigon initialization completed successfully!"
else
    echo "Error: Erigon initialization failed"
    exit 1
fi

# List contents to verify
echo "Erigon data directory contents:"
ls -la "${ERIGON_DATADIR}/"
