
#!/bin/bash

# ================================
# 🚀 Geth Signer + Genesis Script
# ================================

set -e

# === Config ===
DATADIR="./signer-data"
PASSWORD_FILE="./password.txt"
KEYSTORE_DIR="$DATADIR/keystore"
GENESIS_FILE="./genesis.json"
AWS_SECRET_PREFIX="devnet-poa"  # Adjust for AWS Secrets Manager

# === 1️⃣ Create signer-data dir ===
mkdir -p "$DATADIR"

# === 2️⃣ Ask user for password ===
echo ""
read -p "🔐 Please enter a password (leave blank to auto-generate a strong one): " USER_PASSWORD_INPUT

if [ -z "$USER_PASSWORD_INPUT" ]; then
    echo "🔑 No password entered. Generating a random password..."
    PASSWORD=$(openssl rand -hex 16)
else
    PASSWORD="$USER_PASSWORD_INPUT"
    echo "✅ Using provided password."
fi

echo "$PASSWORD" > "$PASSWORD_FILE"
echo "✅ Password saved to $PASSWORD_FILE"

# === 3️⃣ Run geth account new ===
echo ""
echo "🚀 Creating new geth account..."
ACCOUNT_OUTPUT=$(geth --datadir "$DATADIR" account new --password "$PASSWORD_FILE")
echo "$ACCOUNT_OUTPUT"

echo ""
echo "📢 🔑 Public Address: Look for the 'Public address of the key' line above."
echo "   👀 Example: 0x4E47d916C5De84722972B64555D21bd914d5616E"

# === 4️⃣ Prompt user to paste address ===
while true; do
    read -p "👉 Please paste the public address here (starts with 0x): " USER_ADDRESS_INPUT
    SIGNER_ADDRESS=$(echo "$USER_ADDRESS_INPUT" | tr -d '[:space:]')
    if [[ $SIGNER_ADDRESS =~ ^0x[0-9a-fA-F]{40}$ ]]; then
        break
    else
        echo "❌ ERROR: Invalid address format. Must be 0x + 40 hex chars. Please try again."
    fi
done

# Lowercase the address for both alloc + extraData
SIGNER_ADDRESS_LOWER=$(echo "$SIGNER_ADDRESS" | tr '[:upper:]' '[:lower:]')

echo "✅ Captured Signer Address (lowercased): $SIGNER_ADDRESS_LOWER"

# === 5️⃣ Prompt for Chain ID ===
while true; do
    read -p "🔢 Please enter the Chain ID (e.g., 1121): " CHAIN_ID
    if [[ "$CHAIN_ID" =~ ^[0-9]+$ ]]; then
        break
    else
        echo "❌ ERROR: Chain ID must be a numeric value. Please try again."
    fi
done

echo "✅ Using Chain ID: $CHAIN_ID"

# === 6️⃣ Find keystore file ===
KEYSTORE_FILE=$(ls "$KEYSTORE_DIR")
echo ""
echo "📂 Keystore directory: $KEYSTORE_DIR"
echo "🔑 Keystore file: $KEYSTORE_DIR/$KEYSTORE_FILE"

# === 7️⃣ Build genesis.json ===
echo ""
echo "📝 Creating genesis file with signer address: $SIGNER_ADDRESS_LOWER and Chain ID: $CHAIN_ID ..."

# Prepare extraData (Clique: 64 zeroes + signer + 130 zeroes)
SIGNER_NO_PREFIX=$(echo "$SIGNER_ADDRESS_LOWER" | sed 's/^0x//')
EXTRA_DATA="0x$(printf '0%.0s' {1..64})${SIGNER_NO_PREFIX}$(printf '0%.0s' {1..130})"

cat > "$GENESIS_FILE" <<EOF
{
  "config": {
    "chainId": $CHAIN_ID,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "clique": {
      "period": 5,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "8000000",
  "alloc": {
    "$SIGNER_ADDRESS_LOWER": {
      "balance": "1000000000000000000000"
    }
  },
  "extraData": "$EXTRA_DATA"
}
EOF

echo "✅ Genesis file created at $GENESIS_FILE"

# === 8️⃣ (Optional) Upload to AWS Secrets Manager ===
echo ""
echo "🚀 To upload to AWS Secrets Manager, you can run the following commands:"

echo ""
echo "aws secretsmanager create-secret --name ${AWS_SECRET_PREFIX}-password --secret-string file://$PASSWORD_FILE"
echo ""
echo "aws secretsmanager create-secret --name ${AWS_SECRET_PREFIX}-signer-key --secret-string file://$KEYSTORE_DIR/$KEYSTORE_FILE"

echo ""
echo "ℹ️ NOTE: If the secrets already exist, use 'put-secret-value' instead of 'create-secret'."

# === 9️⃣ Reminder for genesis upload ===
echo ""
echo "🌐 To share the genesis with other nodes, upload $GENESIS_FILE to your preferred location (GitHub, S3, etc.)."
