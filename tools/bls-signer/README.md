# BLS Signer

Tool for signing payloads using BLS keys. It takes a BLS private key and payload as input and generates a BLS signature that can be used for provider registration and other signing operations.

## Installation

Download the latest release binary for your platform from the [releases page](../../releases).

The binary is named `bls-signer` and is available for Linux x86_64 and arm64 architectures.

## Usage
The tool requires:
- A BLS private key (provided as a hex encoded string with optional 0x prefix) via the `--private-key` flag
- A payload (provided as a hex string with optional 0x prefix) via the `--payload` flag

The tool will log:
- The payload that was provided
- The hex-encoded BLS public key derived from your private key
- The hex-encoded BLS signature

The logged public key and signature should be used when making a staking request to register as a provider.
