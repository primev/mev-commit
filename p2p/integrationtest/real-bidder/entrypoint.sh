#!/bin/sh

sleep 30

echo "starting bidder-emulator with : ${BIDDER_IP} and rpc url ${L1_RPC_URL}"
/app/bidder-emulator --server-addr ${BIDDER_IP} --rpc-addr ${L1_RPC_URL}
