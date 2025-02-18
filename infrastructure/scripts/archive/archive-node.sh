#!/usr/bin/env bash

set -euo pipefail

download_artifact() {
    local url=$1
    local dest=$2
    echo "Downloading ${url} to ${dest}..."
    curl -fsSL "${url}" -o "${dest}"
}

ENV_FILE=".env"
if [[ -f "${ENV_FILE}" ]]; then
    echo "Loading environment variables from $ENV_FILE..."
    set -a
    source "${ENV_FILE}"
    set +a
else
    echo "Error: $ENV_FILE not found." >&2
    exit 1
fi

mkdir -p "${DATA_DIR}" "${ARTIFACTS_DIR}" "${BIN_DIR}"

download_artifact "${ARTIFACTS_URL}/config_v1.0.0.toml" "${CONFIG_FILE}"
download_artifact "${ARTIFACTS_URL}/genesis_v1.0.0.json" "${GENESIS_FILE}"
download_artifact "${ARTIFACTS_URL}/mev-commit-geth_v1.0.0_Linux_x86_64.tar.gz" "${GETH_ARCHIVE_TAR}"

if [[ ! -x "${GETH_BIN}" ]]; then
    tar -xzf "${GETH_ARCHIVE_TAR}" -C "${BIN_DIR}"
    chmod +x "${GETH_BIN}"
fi

GETH_ZERO_FEE_ADDRESSES="0x509b6a48fc573f0e987cb075cabee75d40e7db85"
GETH_ZERO_FEE_ADDRESSES+=",0xb8a1dac79f4674146a1e49f159107ee2817fdb81"
GETH_ZERO_FEE_ADDRESSES+=",0xf0e4285d437be60975149d5cac2dea49756a238b"
export GETH_ZERO_FEE_ADDRESSES
CHAIN_ID=$(cat "${GENESIS_FILE}" | jq -r .config.chainId)
export CHAIN_ID

echo "Starting node..."
chmod +x ${BIN_DIR}/entrypoint.sh
exec ${BIN_DIR}/entrypoint.sh
