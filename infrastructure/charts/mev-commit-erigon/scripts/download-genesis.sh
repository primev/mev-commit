#!/bin/sh

set -e

echo "Checking for existing genesis file..."

# Check if genesis file already exists
if [ -f "/shared/genesis.json" ]; then
    echo "Genesis file already exists at /shared/genesis.json"
    
    # Validate existing file
    apk add --no-cache jq
    if jq . /shared/genesis.json > /dev/null 2>&1; then
        echo "Existing genesis file is valid JSON. Skipping download."
        exit 0
    else
        echo "Warning: Existing genesis file is invalid JSON. Re-downloading..."
        rm -f /shared/genesis.json
    fi
fi

echo "Downloading genesis from URL..."

# Check if genesis URL is provided
if [ -z "$GENESIS_URL" ]; then
    echo "Error: GENESIS_URL environment variable not set"
    exit 1
fi

echo "Downloading genesis from: $GENESIS_URL"

# Install curl and jq
apk add --no-cache curl jq

# Download genesis file
curl -L -o /shared/genesis.json "$GENESIS_URL"

if [ ! -f "/shared/genesis.json" ]; then
    echo "Error: Failed to download genesis file"
    exit 1
fi

echo "Genesis file downloaded successfully!"
echo "Genesis file location: /shared/genesis.json"

# Validate JSON format
if ! jq . /shared/genesis.json > /dev/null 2>&1; then
    echo "Error: Downloaded file is not valid JSON"
    exit 1
fi

echo "Genesis file validation passed!"
