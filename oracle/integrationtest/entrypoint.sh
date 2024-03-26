#!/bin/sh

sed -i "s|<L1_URL>|${L1_URL}|" /config.yaml

CONFIG=$(cat /config.yaml)

echo "starting mev-commit-oracle with config: ${CONFIG}"
mev-commit-oracle start --config /config.yaml
