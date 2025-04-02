# MEV-Commit Validator Monitoring Tool

A comprehensive tool for monitoring Ethereum validators that have opted into MEV-Commit protocols and analyzing their block proposals across multiple relays.

## Description

This tool monitors the Ethereum beacon chain to identify validators that have opted into the IValidatorOptInRouter contract, tracks their proposer duties, and analyzes relay data to verify proper participation in the MEV-Commit ecosystem.

## Features

- **Validator Discovery**: Identifies validators with proposer duties over a specified time period
- **Opt-In Status Checking**: Verifies validator opt-in status to the MEV-Commit contract
- **Relay Data Analysis**: Queries multiple MEV relays to check if blocks from opted-in validators are visible
- **Data Export**: Exports detailed validator data to Excel for further analysis
- **Parallel Processing**: Uses worker pools for efficient operation
- **Resilient API Handling**: Implements retry logic with exponential backoff
- **Continuous Monitoring**: Supports running as a service with configurable intervals
- **Flexible Time Range**: Configure lookback period in days or epochs

## Installation

### Prerequisites

- Go 1.18 or higher
- Access to Ethereum beacon chain and execution layer endpoints

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/primev/mev-commit.git
   cd mev-commit/tools/monitoring
   ```

2. Build the application:
   ```bash
   go build -o mev-monitor
   ```
