#!/bin/sh

# Replace L1_URL with the value of L1_URL environment variable
sed -i "s|<L1_URL>|${L1_URL}|" /config.yaml

# Replace placeholder text with the values of oracle_user and oracle_pass environment variables
sed -i "s|oracle_user|${ORACLE_USER}|" /config.yaml
sed -i "s|oracle_pass|${ORACLE_PASS}|" /config.yaml

# Read the updated configuration
CONFIG=$(cat /config.yaml)

echo "starting mev-commit-oracle with config: ${CONFIG}"
mev-commit-oracle start --config /config.yaml
