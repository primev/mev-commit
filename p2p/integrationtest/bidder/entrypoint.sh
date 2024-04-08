#!/bin/sh

sleep 30

echo "starting bidder-emulator with : ${BIDDER_IP}"
/app/bidder-emulator --server-addr ${BIDDER_IP} --rpc-addr ${RPC_URL}