#!/bin/sh

DATADIR="${DATADIR:-/geth/data}"
GENESIS_URL="{{ .Values.genesisUrl }}"
GENESIS_FILE="${GENESIS_FILE:-/tmp/genesis.json}"

# Ensure the geth directory exists
mkdir -p "$DATADIR/geth"

# Initialize Geth if not already initialized
if [ -d "$DATADIR/geth/chaindata" ]; then
  echo "ℹ️ Geth already initialized at $DATADIR. Skipping init."
else
  echo "⬇️ Downloading genesis.json from $GENESIS_URL ..."
  wget -O "$GENESIS_FILE" "$GENESIS_URL" || curl -sSL -o "$GENESIS_FILE" "$GENESIS_URL"

  if [ ! -f "$GENESIS_FILE" ]; then
    echo "❌ Failed to download genesis.json! Exiting."
    exit 1
  fi

  echo "🚀 Initializing Geth with downloaded genesis.json..."
  geth --datadir "$DATADIR" init "$GENESIS_FILE"
  
  INIT_RESULT=$?
  if [ $INIT_RESULT -eq 0 ]; then
    echo "✅ Geth init complete."
  else
    echo "❌ Geth init failed with code $INIT_RESULT"
    exit $INIT_RESULT
  fi
fi
