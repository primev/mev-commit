#!/bin/sh

EMULATOR_LOG_LEVEL=${EMULATOR_LOG_LEVEL:-"debug"}
EMULATOR_LOG_FMT=${EMULATOR_LOG_FMT:-"json"}
EMULATOR_BID_WORKERS=${EMULATOR_BID_WORKERS:-"1"}

flags="-log-level $EMULATOR_LOG_LEVEL -log-fmt $EMULATOR_LOG_FMT -bid-workers $EMULATOR_BID_WORKERS"

if [ -n "$EMULATOR_LOG_TAGS" ]; then
  flags="$flags -log-tags $EMULATOR_LOG_TAGS"
fi

if [ -n "$EMULATOR_IP_RPC_PORT" ]; then
  flags="$flags -server-addr $EMULATOR_IP_RPC_PORT"
fi

if [ -n "$EMULATOR_METRICS_PORT" ]; then
  flags="$flags -http-port $EMULATOR_METRICS_PORT"
fi

if [ -n "$EMULATOR_OTEL_COLLECTOR_ENDPOINT_URL" ]; then
  flags="$flags -otel-collector-endpoint-url $EMULATOR_OTEL_COLLECTOR_ENDPOINT_URL"
fi

if [ -n "$EMULATOR_L1_RPC_URL" ]; then
  flags="$flags -rpc-addr $EMULATOR_L1_RPC_URL"
fi

bidder-emulator $flags
