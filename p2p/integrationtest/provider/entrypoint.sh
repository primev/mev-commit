#!/bin/sh

sleep 25

echo "starting provider emulator for: ${PROVIDER_IP}"
/app/provider-emulator --server-addr ${PROVIDER_IP}

