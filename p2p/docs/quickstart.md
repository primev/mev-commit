# Quickstart

## Using the CLI

mev-commit software is a CLI program that needs to be run by the bidders in their own environments. The CLI has two commands mainly. You can run the main command with `-h`/`--help` option to see the available commands.

```
❯ mev-commit -h
NAME:
   mev-commit - Entry point for mev-commit

USAGE:
   mev-commit [global options] command [command options] [arguments...]

VERSION:
   "v1.0.0-alpha-a834960"

COMMANDS:
   start       Start the mev-commit node
   create-key  Create a new ECDSA private key and save it to a file
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Bidders can use any existing ECDSA private key or create one using the `create-key` command. This key is used to derive the ethereum address of the node taking part in the network.

```
❯ mev-commit create-key -h
NAME:
   mev-commit create-key - Create a new ECDSA private key and save it to a file

USAGE:
   mev-commit create-key [command options] <output_file>

OPTIONS:
   --help, -h  show help
```

In order to run a node, bidders need to create a configuration file in the YAML format. Example config files can be found in the [config](https://github.com/primevprotocol/mev-commit/tree/main/config) folder. The important options are defined below:

```yaml
# Path to private key file.
priv_key_file: ~/.mev-commit/keys/nodekey

# Type of peer. Options are provider and bidder.
peer_type: provider

# Port used for P2P traffic. If not configured, 13522 is the default.
p2p_port: 13522

# Port used for HTTP traffic. If not configured, 13523 is the default.
http_port: 13523

# Port used for RPC API. If not configured, 13524 is the default.
rpc_port: 13524

# Secret for the node. This is used to authorize the nodes. The value doesnt matter as long as it is sufficiently unique. It is signed using the private key.
secret: hello

# Format used for the logs. Options are "text" or "json".
log_fmt: text

# Log level. Options are "debug", "info", "warn" or "error".
log_level: debug

# Bootnodes used for bootstrapping the network.
bootnodes:
  - /ip4/35.91.118.20/tcp/13522/p2p/16Uiu2HAmAG5z3E8p7o19tEcLdGvYrJYdD1NabRDc6jmizDva5BL3

# The default is set to false for development reasons. Change it to true if you wish to accept bids on your provider instance
expose_provider_api: false
```

Place this config file in some folder. It is advised to create a new `.mev-commit` folder in the home directory and place it there. Once this is done, bidders can start the node using the `start` command.

```
❯ mev-commit start -h
NAME:
   mev-commit start - Start the mev-commit node

USAGE:
   mev-commit start [command options] [arguments...]

OPTIONS:
   --config value  path to config file [$MEV_COMMIT_CONFIG]
   --help, -h      show help
```

## Docker

Optionally bidders can use the docker images to run the client in their docker environment. You need to make sure the node is able to communicate with provider/bidder nodes.
