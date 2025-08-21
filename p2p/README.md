# New Bidder UX Demo

The following demonstrates the new preconfirmation flow from the perspective of a bidder. The demo works locally with anvil and simulates parts of the system that are not bidder-specific.

## Actors

**Contract owner** address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

**Bidder** address: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Private Key: 59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

**Provider** address: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
Private Key: 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a

**Oracle** address: 0xa0Ee7A142d267C1f36714E4a8F75612F20a79720
Private Key: 0x2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6

## Demo

1. Start local anvil node
2. Deploy core contracts by running `make core` from `contracts` directory
3. Register provider by running `make provider-reg` from `contracts` directory

4. Start local bidder node by running `make bidder` from `p2p` directory

5. Run `make getcode` from `contracts` directory to get the code of the bidder node

6. Enable deposit manager for bidder. This command sets the code of the bidder EOA to the DepositManager implementation.
  ```
  curl -s -X POST http://localhost:13523/v1/bidder/enable_deposit_manager \
    -H 'Content-Type: application/json' -d '{}' | jq
  ```

7. Run `make getcode` from `contracts` directory again to get the code of the bidder node

8. Get valid providers. The returned list included providers which are both connected via p2p, and fully registered/staked with the provider registry.

```
curl -s http://localhost:13523/v1/bidder/get_valid_providers | jq
```

9. Set target deposit of 3 ETH for only provider returned by `get_valid_providers`

```
curl -s -X POST http://localhost:13523/v1/bidder/set_target_deposits \
  -H 'Content-Type: application/json' \
  -d '{"target_deposits":[{"provider":"0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC","target_deposit":"3000000000000000000"}]}' | jq
```

10. Get all deposits. Note bidderBalance before we bid.

```
curl -s http://localhost:13523/v1/bidder/get_all_deposits | jq
```

11. Bid 1 ETH. Commitment process is simulated. Commitment is opened onchain by the provider, and deposit relevant to that provider should be topped-up atomically.

```
curl -X POST http://localhost:13523/v1/bidder/bid \
-d '{
    "txHashes": ["0549fc7c57fffbdbfb2cf9d5e0de165fc68dadb5c27c42fdad0bdf506f4eacae"],
    "amount": "1000000000000000000",
    "blockNumber": 9999,
    "decayStartTimestamp": 1111,
    "decayEndTimestamp": 2222,
    "revertingTxHashes": []
}'
```

12. __Try calling `get_all_deposits` again. Note that the bidder balance has decreased by 1 ETH due to top-up__

```
curl -s http://localhost:13523/v1/bidder/get_all_deposits | jq
```
