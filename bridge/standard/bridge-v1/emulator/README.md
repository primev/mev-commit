# Bridge user emulator

This directory contains an emulator program for stress testing the standard bridge.

## Accounts

Since a [local l1 network](https://github.com/primevprotocol/mev-commit-geth/blob/552a9cf940652156441b2608907839e635208af9/geth-poa/docker-compose.yml#L168) is used for testing. Five accounts will be allocated L1 balances on genesis. This can be confirmed with the local l1's [genesis.json](https://github.com/primevprotocol/mev-commit-geth/blob/de3ff446517b87f3ac0ad322af3a2edc382bf13a/geth-poa/local-l1/genesis.json).

To best simulate new users of the mev-commit chain, these emulator accounts will have no corresponding genesis allocation on the mev-commit chain. The emulators will continuously bridge a random value in [0.01, 10] ETH from L1 -> mev-commit chain, then bridge that random value (minus 0.009 ETH for fees) back to L1.

Address and private keys will be listed below. DO NOT use these key pairs in production. 

### Emulator 1

- Address: `0x0b1f1268f138aEEb12F54142B2359944904aaf6e`
- Private key: `0xd821c54fd6dfd4c864202c125eaedb0a1fad7f40c81863fa3038338b475ff44a`

### Emulator 2

- Address: `0x9D58dB6c050E0E708b06c8e40aE803b5c0a793B0`
- Private key: `0x6e7470ba919624df632ebe77cccda95f38e9376cf2e7e3cb42726cb23457abae`

### Emulator 3

- Address: `0x911FA3b5D45c1A5E6316830dd5B3fCcce1b421FF`
- Private key: `0x25710d3869ef44b2f026615c0734ff3f44c17d319f7fc6db318f0cacece3d575`

### Emulator 4

- Address: `0xAA7E52aF3c86Aa8b617670347024e86C3b26bcd8`
- Private key: `0xa5743a07de65be93c11e1814bfcf4b00d3e49d0e84cd488e6cb98ba4410d585f`

### Emulator 5

- Address: `0x72F17F6a137645B774c3a115530410732019fE84`
- Private key: `0x865a133ee5b85c4f4970be2f6ac590a7c4be51b936df48028f0285bed12f1f10`
