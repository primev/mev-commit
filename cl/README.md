# README.md

## Introduction

This project sets up a local Ethereum environment with two Geth nodes, a Redis instance, and a consensus client. The setup is useful for testing and development purposes, allowing you to simulate a blockchain network locally.

## Prerequisites

- **Go Programming Language**: Ensure you have Go installed on your system.
- **Geth (Go Ethereum)**: Install Geth from the [official website](https://geth.ethereum.org/downloads/).
- **Docker and Docker Compose**: Install Docker and Docker Compose from the [official website](https://docs.docker.com/get-docker/).
- **Redis**: We will use Redis via Docker Compose.
- **Git**: To clone the repository if needed.
- **Foundry (Optional)**: For interacting with the blockchain using `cast`.

## Setup Directory Structure

Ensure you have a directory named `geth-setup` in `cl` containing `genesis.json` and `jwt.hex` files:

```bash
geth-setup/
├── genesis.json
└── jwt.hex
```

## Installation

### Set Up JWT Secret

If you don't already have a `jwt.hex` file, create one:

```bash
echo "your_jwt_secret" > geth-setup/jwt.hex
```

**Note**: Replace `"your_jwt_secret"` with your actual JWT secret. In a production environment, ensure this secret is kept secure.

## Running Geth Nodes

We will initialize and run two Geth nodes.

### Initialize Geth Nodes

First, remove any existing data directories for clean setups:

```bash
rm -rf client1 && rm -rf client2
```

Initialize the nodes with the `genesis.json` file located in `geth-setup`:

```bash
geth init --datadir ./client1 geth-setup/genesis.json
geth init --datadir ./client2 geth-setup/genesis.json
```

### Run Geth Nodes

#### Node 1

```bash
geth --verbosity 5 \
  --datadir ./client1 \
  --nodiscover \
  --http \
  --http.port 8545 \
  --port 30303 \
  --authrpc.jwtsecret ./geth-setup/jwt.hex \
  --authrpc.port 8551 \
  --networkid 141414 \
  --http.api "admin,eth,net,web3,engine" \
  --syncmode full \
  --miner.recommit 900ms
```

#### Node 2

```bash
geth --verbosity 5 \
  --datadir ./client2 \
  --nodiscover \
  --http \
  --http.port 8546 \
  --port 30304 \
  --authrpc.jwtsecret ./geth-setup/jwt.hex \
  --authrpc.port 8552 \
  --networkid 141414 \
  --http.api "admin,eth,net,web3,engine" \
  --syncmode full \
  --miner.recommit 900ms
```

### Explanation of Geth Flags

- `--verbosity 5`: Sets the logging verbosity.
- `--datadir`: Specifies the data directory for the node.
- `--nodiscover`: Disables the peer discovery mechanism.
- `--http`: Enables the HTTP-RPC server.
- `--http.port`: Port for the HTTP-RPC server.
- `--port`: Network listening port.
- `--authrpc.jwtsecret`: Path to the JWT secret file.
- `--authrpc.port`: Port for authenticated RPC endpoints.
- `--networkid`: Network identifier for the Ethereum network.
- `--http.api`: APIs exposed over the HTTP-RPC interface.
- `--syncmode full`: Synchronization mode.
- `--miner.recommit`: Frequency of miner recommit.

## Running Redis

We will use Docker Compose to run Redis.

### Docker Compose Configuration

Redis is configured in `redis-cluster` folder withing `docker-compose.yml`

### Start Redis

Stop any existing containers and remove volumes:

```bash
docker compose down -v
```

Start Redis in detached mode:

```bash
docker compose up -d
```

### Verify Redis is Running

You can verify that Redis is running by connecting to it:

```bash
redis-cli -p 7001
```

## Running the Consensus Client

The consensus client is a Go application that interacts with the Geth nodes and Redis.

### Build the Consensus Client

Ensure all dependencies are installed and build the application:

```bash
go mod tidy
go build -o consensus-client cmd/redisapp/main.go
```

### Consensus Client Configuration

The consensus client can be configured via command-line flags, environment variables, or a YAML configuration file.

#### Command-Line Flags for Streamer

- `--instance-id`: **(Required)** Unique instance ID for this node.
- `--eth-client-url`: Ethereum client URL (default: `http://localhost:8551`).
- `--jwt-secret`: JWT secret for Ethereum client.
- `--redis-addr`: Redis address (default: `127.0.0.1:7001`).
- `--evm-build-delay`: EVM build delay (default: `1s`).
- `--config`: Path to a YAML configuration file.

#### Environment Variables for Consensus Client

- `RAPP_INSTANCE_ID`
- `RAPP_ETH_CLIENT_URL`
- `RAPP_JWT_SECRET`
- `RAPP_REDIS_ADDR`
- `RAPP_EVM_BUILD_DELAY`
- `RAPP_CONFIG`

### Run the Consensus Client

Run the client using command-line flags:

```bash
./consensus-client start \
  --instance-id "node1" \
  --eth-client-url "http://localhost:8551" \
  --jwt-secret "your_jwt_secret" \
  --redis-addr "127.0.0.1:7001" \
  --evm-build-delay "1s"
```

**Note**:

- Replace `"your_jwt_secret"` with the actual JWT secret you used earlier.

### Using a Configuration File for Consensus Client

Create a `config.yaml` file:

```yaml
instance-id: "node1"
eth-client-url: "http://localhost:8551"
jwt-secret: "your_jwt_secret"
redis-addr: "127.0.0.1:7001"
evm-build-delay: "1s"
```

Run the client with the configuration file:

```bash
./consensus-client start --config config.yaml
```

## Additional Notes

- **Multiple Instances**: You can run multiple instances of the consensus client by changing the `--instance-id` and `--eth-client-url` parameters.

## Running the Streamer

The Streamer is responsible for streaming payloads to member nodes, allowing them to apply these payloads to their respective Geth instances.

### Build the Streamer

Ensure all dependencies are installed and build the Streamer application:

```bash
go mod tidy
go build -o streamer cmd/streamer/main.go
```

### Streamer Configuration

The Streamer can be configured via command-line flags, environment variables, or a YAML configuration file.

#### Command-Line Flags

- `--config`: Path to config file.
- `--redis-addr`: Redis address (default: 127.0.0.1:7001).
- `--listen-addr`: Streamer listen address (default: :50051).
- `--log-fmt`: Log format to use, options are text or json (default: text).
- `--log-level`: Log level to use, options are debug, info, warn, error (default: info).

#### Environment Variables

- `STREAMER_CONFIG`
- `STREAMER_REDIS_ADDR`
- `STREAMER_LISTEN_ADDR`
- `STREAMER_LOG_FMT`
- `STREAMER_LOG_LEVEL`

#### Run the Streamer

Run the Streamer using command-line flags:

```bash
./streamer start \
  --config "config.yaml" \
  --redis-addr "127.0.0.1:7001" \
  --listen-addr ":50051" \
  --log-fmt "json" \
  --log-level "info"
```

#### Using a Configuration File for Streamer

Create a `streamer_config.yaml` file:

```yaml
redis-addr: "127.0.0.1:7001"
listen-addr: ":50051"
log-fmt: "json"
log-level: "info"
```

Run the Streamer with the configuration file:

```bash
./streamer start --config streamer_config.yaml
```

## Running member nodes

Member nodes connect to the Streamer to receive payloads from the stream and apply them to their Geth instances.

### Build the Member Client

Ensure all dependencies are installed and build the Member Client application:

```bash
go mod tidy
go build -o memberclient cmd/member/main.go
```

### Configuration

The Member Client can be configured via command-line flags, environment variables, or a YAML configuration file.

### Command-Line Flags for Member Client

- `--config`: Path to config file.
- `--client-id`: (Required) Unique client ID for this member.
- `--streamer-addr`: (Required) Streamer address.
- `--eth-client-url`: Ethereum client URL (default: <http://localhost:8551>).
- `--jwt-secret`: JWT secret for Ethereum client.
- `--log-fmt`: Log format to use, options are text or json (default: text).
- `--log-level`: Log level to use, options are debug, info, warn, error (default: info).

### Environment Variables for Member Client

- `MEMBER_CONFIG`
- `MEMBER_CLIENT_ID`
- `MEMBER_STREAMER_ADDR`
- `MEMBER_ETH_CLIENT_URL`
- `MEMBER_JWT_SECRET`
- `MEMBER_LOG_FMT`
- `MEMBER_LOG_LEVEL`

### Run the Member Client

Run the Member Client using command-line flags:

```bash
./memberclient start \
  --client-id "member1" \
  --streamer-addr "http://localhost:50051" \
  --eth-client-url "http://localhost:8551" \
  --jwt-secret "your_jwt_secret" \
  --log-fmt "json" \
  --log-level "info"
```

Note:

Replace "your_jwt_secret" with the actual JWT secret you used earlier.

### Using a Configuration File

Create a member_config.yaml file:

```yaml
client-id: "member1"
streamer-addr: "http://localhost:50051"
eth-client-url: "http://localhost:8551"
jwt-secret: "your_jwt_secret"
log-fmt: "json"
log-level: "info"
```

Run the Member Client with the configuration file:

```bash
./memberclient start --config member_config.yaml
```

## Conclusion

You now have a local Ethereum environment with Geth nodes, Redis, and a consensus client.
