# Use the latest foundry image
FROM ghcr.io/foundry-rs/foundry

# Set working directory
WORKDIR /app

# Copy our source code into the container
COPY . .

# Compile contracts using solidity compiler version 0.8.23
RUN forge build --use 0.8.23 --via-ir

# Set environment variables for RPC URL and private key
# These should be passed during the Docker build process
ARG RPC_URL
ARG PRIVATE_KEY
ARG CHAIN_ID
ARG DEPLOY_TYPE
ARG HYP_ERC20_ADDR 
ARG RELAYER_ADDR

RUN chmod +x entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]

