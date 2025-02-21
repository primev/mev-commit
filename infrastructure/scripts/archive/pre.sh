#!/bin/bash

set -e

SERVICE_NAME="mev-commit-geth"
MOUNT_POINT="/mnt/geth-data"

echo "=== Pre-Snapshot Script: Stopping ${SERVICE_NAME} and Freezing Filesystem ==="

echo "Stopping ${SERVICE_NAME} service..."
sudo systemctl stop "${SERVICE_NAME}"

if systemctl is-active --quiet "${SERVICE_NAME}"; then
    echo "Error: ${SERVICE_NAME} service is still running."
    exit 1
else
    echo "${SERVICE_NAME} service stopped successfully."
fi

echo "Freezing filesystem at ${MOUNT_POINT}..."
sudo fsfreeze -f "${MOUNT_POINT}"
echo "Filesystem frozen successfully."

echo "=== Pre-Snapshot Steps Completed Successfully ==="

exit 0
