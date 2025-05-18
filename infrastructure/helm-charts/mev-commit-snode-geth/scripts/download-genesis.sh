#!/bin/sh
set -e

LOG_PREFIX="[download-genesis]"
GENESIS_URL="{{ .Values.genesisUrl }}"
GENESIS_FILE="/shared/genesis.json"
MAX_RETRIES=3

echo "$LOG_PREFIX Starting genesis file initialization"

# Ensure shared directory exists and is writable
if [ ! -d "/shared" ]; then
  echo "$LOG_PREFIX Creating /shared directory"
  mkdir -p /shared
  if [ $? -ne 0 ]; then
    echo "$LOG_PREFIX ERROR: Failed to create /shared directory, exiting"
    exit 1
  fi
fi

echo "$LOG_PREFIX Setting permissions on /shared"
chmod 777 /shared
if [ $? -ne 0 ]; then
  echo "$LOG_PREFIX WARNING: Failed to set permissions on /shared, continuing anyway"
fi

# List directory contents before download
echo "$LOG_PREFIX Current directory content before download:"
ls -la /shared/

# Check if genesis file already exists to handle restart scenarios
if [ -f "$GENESIS_FILE" ]; then
  echo "$LOG_PREFIX Genesis file already exists, checking content"
  
  # Check if file is not empty
  if [ -s "$GENESIS_FILE" ]; then
    # Simple validation: check if file contains expected JSON format
    if grep -q "\"config\"" "$GENESIS_FILE" && grep -q "\"alloc\"" "$GENESIS_FILE"; then
      echo "$LOG_PREFIX Existing genesis file appears valid, skipping download"
      exit 0
    else
      echo "$LOG_PREFIX Existing genesis file does not appear to be valid, will redownload"
    fi
  else
    echo "$LOG_PREFIX Existing genesis file is empty, will redownload"
  fi
fi

# Download the genesis file with retries
echo "$LOG_PREFIX Downloading genesis file from $GENESIS_URL"
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if wget -O "$GENESIS_FILE" "$GENESIS_URL" 2>/dev/null; then
    echo "$LOG_PREFIX Genesis file downloaded successfully"
    break
  else
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -lt $MAX_RETRIES ]; then
      echo "$LOG_PREFIX Download failed, retrying ($RETRY_COUNT/$MAX_RETRIES)..."
      sleep 2
    else
      echo "$LOG_PREFIX ERROR: Failed to download genesis file after $MAX_RETRIES attempts, exiting"
      exit 1
    fi
  fi
done

# Validate the downloaded file
if [ -s "$GENESIS_FILE" ]; then
  echo "$LOG_PREFIX Genesis file content (first 5 lines):"
  head -5 "$GENESIS_FILE"
  echo "$LOG_PREFIX ..."
  # List directory contents to confirm file exists
  echo "$LOG_PREFIX Directory content after download:"
  ls -la /shared/
  echo "$LOG_PREFIX Genesis file download and validation completed successfully"
else
  echo "$LOG_PREFIX ERROR: Downloaded genesis file is empty, exiting"
  exit 1
fi

exit 0
