# provider emulator

This example demonstrates the provider-side workflows to interact with their mev-commit
nodes. The mev-commit node will get bids from different bidders/providers through the
p2p protocols. The provider has to use the RPC API to start looking at these bids and then
make decisions on whether to confirm these bids that it receives.

The confirmed bids need to be sent back to the mev-commit node to send them back to the
bidders. This is also done using a separate RPC request. The response is mapped to
the bid hash that it receives.

The typical workflow:

1. Create a new client. This starts a global sender routine that opens a new RPC streaming
   request to send the confirmed bids back.

2. The client calls ReceiveBids to start receiving the bids over RPC.

3. The example process looks at the bid and accepts it by default after printing it. Here
   clients should have their own logic to accept/reject bids.
