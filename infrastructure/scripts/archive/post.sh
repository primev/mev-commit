#!/bin/bash

set -e

SERVICE_NAME="mev-commit-geth"
MOUNT_POINT="/mnt/geth-data"

echo "=== Post-Snapshot Script: Unfreezing Filesystem and Starting ${SERVICE_NAME} ==="

echo "Unfreezing filesystem at ${MOUNT_POINT}..."
sudo fsfreeze -u "${MOUNT_POINT}"
echo "Filesystem unfrozen successfully."

echo "Starting ${SERVICE_NAME} service..."
sudo systemctl start "${SERVICE_NAME}"

if systemctl is-active --quiet "${SERVICE_NAME}"; then
    echo "${SERVICE_NAME} service started successfully."
else
    echo "Error: Failed to start ${SERVICE_NAME} service."
    exit 1
fi

echo "=== Post-Snapshot Steps Completed Successfully ==="

exit 0
