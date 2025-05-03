#!/bin/sh

DATADIR="${DATADIR:-/geth/data}"
NODEKEY_FILE="$DATADIR/geth/nodekey"
PUBKEY_FILE="$DATADIR/geth/pubkey"
GETH_NODE_TYPE="{{ .Values.role }}"

# Ensure the geth directory exists
mkdir -p "$DATADIR/geth"

# Only proceed if this is a bootnode
if [ "$GETH_NODE_TYPE" = "bootnode" ]; then
  # Check if nodekey exists (should already be there from Geth init)
  if [ -f "$NODEKEY_FILE" ]; then
    echo "✅ Existing nodekey found at $NODEKEY_FILE"
    
    # Generate pubkey file regardless of whether it exists or not
    if command -v bootnode >/dev/null 2>&1; then
      echo "🔑 Generating pubkey from existing nodekey..."
      bootnode -nodekey "$NODEKEY_FILE" -writeaddress > "$PUBKEY_FILE"
      echo "✅ Created/updated pubkey file at $PUBKEY_FILE"
      
      # Show the pubkey/node-id for reference
      NODE_ID=$(cat "$PUBKEY_FILE")
      echo "📋 Node ID: $NODE_ID"
    else
      echo "⚠️ bootnode command not available, cannot generate pubkey"
    fi
  else
    echo "❌ No nodekey found at $NODEKEY_FILE. This is unexpected."
    # Optional: Generate a nodekey if it's missing (unlikely)
    if command -v bootnode >/dev/null 2>&1; then
      echo "🔑 Generating new nodekey..."
      bootnode -genkey "$NODEKEY_FILE"
      echo "✅ Created new nodekey at $NODEKEY_FILE"
      
      echo "🔑 Generating pubkey from new nodekey..."
      bootnode -nodekey "$NODEKEY_FILE" -writeaddress > "$PUBKEY_FILE"
      echo "✅ Created pubkey file at $PUBKEY_FILE"
      
      # Show the pubkey/node-id for reference
      NODE_ID=$(cat "$PUBKEY_FILE")
      echo "📋 Node ID: $NODE_ID"
    else
      echo "⚠️ bootnode command not available, cannot generate nodekey or pubkey"
    fi
  fi
else
  echo "⏭️ Role is not bootnode (current role: $GETH_NODE_TYPE), skipping nodekey/pubkey operations"
fi
