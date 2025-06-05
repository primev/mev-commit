#!/bin/sh

EMULATOR_LOG_LEVEL=${EMULATOR_LOG_LEVEL:-"debug"}
EMULATOR_LOG_FMT=${EMULATOR_LOG_FMT:-"json"}

flags="-log-level $EMULATOR_LOG_LEVEL -log-fmt $EMULATOR_LOG_FMT"

if [ -n "$EMULATOR_LOG_TAGS" ]; then
  flags="$flags -log-tags $EMULATOR_LOG_TAGS"
fi

if [ -n "$EMULATOR_IP_RPC_PORT" ]; then
  flags="$flags -server-addr $EMULATOR_IP_RPC_PORT"
fi

if [ -n "$EMULATOR_METRICS_PORT" ]; then
  flags="$flags -http-port $EMULATOR_METRICS_PORT"
fi

if [ -n "$EMULATOR_RELAY_URL" ]; then
  flags="$flags -relay $EMULATOR_RELAY_URL"
fi

if [ -n "$EMULATOR_ERROR_PROBABILITY" ]; then
  flags="$flags -error-probability $EMULATOR_ERROR_PROBABILITY"
fi

if [ -n "$EMULATOR_OTEL_COLLECTOR_ENDPOINT_URL" ]; then
  flags="$flags -otel-collector-endpoint-url $EMULATOR_OTEL_COLLECTOR_ENDPOINT_URL"
fi

provider-emulator $flags
