#!/bin/sh
set -e

LOG_PREFIX="[init-geth]"
DATA_DIR="/data/geth"
GENESIS_FILE="/shared/genesis.json"
CHAINDATA_DIR="$DATA_DIR/geth/chaindata"

echo "$LOG_PREFIX Starting Geth initialization"

# Check if genesis file exists and display its path
echo "$LOG_PREFIX Looking for genesis file at: $GENESIS_FILE"
ls -la /shared/

# Check if genesis file exists
if [ ! -f "$GENESIS_FILE" ]; then
  echo "$LOG_PREFIX ERROR: Genesis file not found at $GENESIS_FILE, exiting"
  exit 1
fi

# Check if the genesis file is valid
if ! grep -q "\"config\"" "$GENESIS_FILE"; then
  echo "$LOG_PREFIX ERROR: Genesis file appears to be invalid or corrupted, exiting"
  exit 1
fi

# Check if the data directory already contains chaindata (restart tolerance)
if [ -d "$CHAINDATA_DIR" ]; then
  echo "$LOG_PREFIX Geth chaindata directory already exists at $CHAINDATA_DIR"
  echo "$LOG_PREFIX Checking if chaindata is valid..."
  
  # Check if the chaindata directory is not empty
  if [ "$(ls -A "$CHAINDATA_DIR" 2>/dev/null)" ]; then
    echo "$LOG_PREFIX Valid chaindata found, skipping initialization to prevent data corruption"
    exit 0
  else
    echo "$LOG_PREFIX Chaindata directory exists but is empty, will reinitialize"
  fi
fi

# Ensure data directory exists
if [ ! -d "$DATA_DIR" ]; then
  echo "$LOG_PREFIX Creating data directory at $DATA_DIR"
  mkdir -p "$DATA_DIR"
  if [ $? -ne 0 ]; then
    echo "$LOG_PREFIX ERROR: Failed to create data directory, exiting"
    exit 1
  fi
fi

# Initialize Geth with the genesis file
echo "$LOG_PREFIX Initializing Geth with genesis file"
geth --datadir "$DATA_DIR" init "$GENESIS_FILE"

# Check the exit code
if [ $? -ne 0 ]; then
  echo "$LOG_PREFIX ERROR: Geth initialization failed, exiting"
  exit 1
fi

# Verify the initialization was successful
if [ -d "$CHAINDATA_DIR" ] && [ "$(ls -A "$CHAINDATA_DIR" 2>/dev/null)" ]; then
  echo "$LOG_PREFIX Geth initialized successfully, chaindata directory populated"
else
  echo "$LOG_PREFIX WARNING: Geth initialization completed but chaindata directory appears empty"
fi

echo "$LOG_PREFIX Geth initialization completed successfully"
exit 0
