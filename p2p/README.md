
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

__Run `make getcode` from `contracts` directory to get the code of the bidder node__

4. Start local bidder node by running `make bidder` from `p2p` directory
5. Enable deposit manager for bidder. This command sets the code of the bidder EOA to the DepositManager implementation.
  ```
  curl -s -X POST http://localhost:13523/v1/bidder/enable_deposit_manager \
    -H 'Content-Type: application/json' -d '{}' | jq
  ```

__Run `make getcode` from `contracts` directory again to get the code of the bidder node__

6. Get all deposits for bidder. There will be none.

```
curl -s http://localhost:13523/v1/bidder/get_all_deposits | jq
```

7. Get valid providers. The returned list included providers which are both connected via p2p, and fully registered/staked with the provider registry.

```
curl -s http://localhost:13523/v1/bidder/get_valid_providers | jq
```

8. Set target deposit of 3 ETH for only provider returned by `get_valid_providers`

```
curl -s -X POST http://localhost:13523/v1/bidder/set_target_deposits \
    -H 'Content-Type: application/json' -d '{"target_deposits": [{"provider": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "amount": "3000000000000000000"}]}' | jq
```



## TODO: Document with various commands for demo. Node API, checking balances, simulating bid, etc.


Mock out the whole preconf flow, a provider registering, being seen by bidder. balances being decremented when preconf happens etc. Will need to have some make targets for oracle doing things etc....

Make sure to demonstrate auto top-up behavior (not having to call deposit again, but eth balance being moved from eoa balance etc. ) .also demonstrate changing the target deposit amounts, then seeing how top-up behaves etc. 

Have two bidders, and some tools to simulate multiple bids identified by vanity hashes 0x001, 0x002, etc.
