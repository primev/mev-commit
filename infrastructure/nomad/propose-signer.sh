#!/usr/bin/env sh

# Check if SIGNER_ADDRESS_TO_PROPOSE environment variable is set
if [ -z "$SIGNER_ADDRESS_TO_PROPOSE" ]; then
    echo "Error: ADDRESS_TO_PROPOSE environment variable is not set."
    exit 1
fi

# Initialize an index variable
index=1

# Iterate over each argument passed to the script (each argument is an allocation ID)
for ALLOC_ID in "$@"
do
    # Define the command to be executed, using the ADDRESS_TO_PROPOSE environment variable
    COMMAND="local/geth attach --exec \"clique.propose(\\\"$SIGNER_ADDRESS_TO_PROPOSE\\\", true)\" /local/geth-signer-node${index}-data-0/geth.ipc"

    # Execute the command within the specified allocation
    echo "Executing command in allocation $ALLOC_ID..."
    nomad alloc exec $ALLOC_ID /bin/sh -c "$COMMAND"

    # Increment the index for the next iteration
    ((index++))

done
