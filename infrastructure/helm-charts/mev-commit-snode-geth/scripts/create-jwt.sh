#!/bin/sh
set -e

LOG_PREFIX="[create-jwt]"
JWT_FILE="/shared/jwt"
JWT_SECRET="{{ .Values.secrets.jwtSecret }}"

echo "$LOG_PREFIX Starting JWT initialization"

# Check for existing valid file
if [ -f "$JWT_FILE" ]; then
  EXISTING_JWT=$(cat "$JWT_FILE" | tr -d '\n')
  if echo "$EXISTING_JWT" | grep -qE '^[a-fA-F0-9]{64}$'; then
    echo "$LOG_PREFIX Valid existing JWT found, reusing"
    exit 0
  else
    echo "$LOG_PREFIX Invalid JWT found, regenerating"
  fi
fi

# Validate input from values.yaml
if ! echo "$JWT_SECRET" | grep -qE '^[a-fA-F0-9]{64}$'; then
  echo "$LOG_PREFIX ERROR: Provided JWT secret is invalid!"
  exit 1
fi

echo "$JWT_SECRET" > "$JWT_FILE"
sync

echo "$LOG_PREFIX JWT written to $JWT_FILE"
