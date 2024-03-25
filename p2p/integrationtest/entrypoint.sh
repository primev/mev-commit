#!/bin/sh

echo "Node Type: ${NODE_TYPE}"

# If this is not the bootnode, update the bootnodes entry with P2P ID
if [ "${NODE_TYPE}" != "bootnode" ]; then
    # Wait for a few seconds to ensure the bootnode is up and its API is accessible
    sleep 10
fi

sed -i "s|<BIDDER_REGISTRY>|${BIDDER_REGISTRY}|" /config.yaml
sed -i "s|<PROVIDER_REGISTRY>|${PROVIDER_REGISTRY}|" /config.yaml
sed -i "s|<RPC_URL>|${RPC_URL}|" /config.yaml

if [ "${NODE_TYPE}" == "provider" ]; then
    sed -i "s|<PRECONF_CONTRACT>|${PRECONF_CONTRACT}|" /config.yaml
fi

CONFIG=$(cat /config.yaml)

echo "starting mev-commit with config: ${CONFIG}"
/app/mev-commit --config /config.yaml
