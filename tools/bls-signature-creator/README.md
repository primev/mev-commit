# BLS Signature Creator

A tool for creating BLS signatures for provider registration.

## Installation

Download the latest release binary for your platform from the [releases page](../../releases).

The binary is named `bls-signature-creator` and is available for Linux x86_64 and arm64 architectures.

## Usage
The tool requires:
- A BLS private key (provided as a hex encoded string with optional 0x prefix) via the `--private-key` flag
- An Ethereum address (provided as a hex string with optional 0x prefix) via the `--eth-address` flag

The tool will log:
- The Ethereum address that was provided
- The hex-encoded BLS public key derived from your private key
- The hex-encoded BLS signature

The logged public key and signature should be used when making a staking request to register as a provider.
