# Bidder CLI
Bidder CLI is a command-line tool used to interact with a gRPC bidder server.

## Usage
Bidder CLI can be used with the following commands:

### 1. Send Bid
Used to send a bid.

```bash
bidder-cli send-bid --txhash <transaction_hash> --amount <bid_amount> --block <block_number>
```
* `--txhash`: Transaction hash.
* `--amount`: Bid amount.
* `--block`: Block number.

### 2. Check Status
Used to check the status of the gRPC bidder server.

```bash
bidder-cli status
```

### 3. Send Random Bid
Used to send a random bid.

```bash
bidder-cli send-rand-bid
```

## Configuration
Configuration options can be set using the config.yml file. An example configuration file is as follows:

```yaml
server_address: "localhost:13524"
log_fmt: "text"
log_level: "info"
```

* `server_address`: The address and port of the gRPC bidder server.
* `log_fmt`: Log format (text or json).
* `log_level`: Log level (debug, info, warn, error).

You can modify these settings to suit your specific environment and preferences.

