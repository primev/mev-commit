#!/bin/sh

THRESHOLD=50000000000000000000  # Set your threshold value here in wei
# THRESHOLD is set to 50 Ether in wei for comparison
echo "Starting funder at RPC endpoint $RPC_URL with threshold $THRESHOLD"
while true; do
    while IFS= read -r address; do
        balance=$(cast balance $address --rpc-url $RPC_URL)
        balance=${balance%.*}  # Remove decimals, assuming the balance is returned in the smallest unit (wei)

        below_threshold=$(echo "$balance < $THRESHOLD" | bc -l)

        if [ "$below_threshold" -eq 1 ]; then
            echo "Address $address has a balance $balance below threshold. Funding..."
            # Use Foundry's cast command to fund the address
            cast send $address --value=50ether --private-key $PRIVATE_KEY --rpc-url $RPC_URL
        fi
    done < /app/addresses.txt

    # Sleep for a specified interval before checking again
    sleep 180  # Adjust the sleep duration (in seconds) as needed
done
