# Debugging

The mev-commit node is still in active development. As a result there is some debugging instrumentation available on the node. In the normal use-cases, the bidders might not need this, but they can always use this to get more information on the inner-workings of the node. These tools can be accessed on the HTTP API endpoint.

## Network topology

The topology of the network as seen by the mev-commit node can be obtained on the `/topology` endpoint of the node. This should show the no. of nodes connected and their roles. This would be mainly used to check if the node has sufficient connectivity with the network.

```json
{
  "self": {
    "Addresses": [
      "/ip4/127.0.0.1/tcp/13522",
      "/ip4/172.28.0.4/tcp/13522"
    ],
    "Ethereum Address": "0xB61545548948E9299Ce6eb4C01F2C31FcE6c9E83",
    "Peer Type": "bidder",
    "Underlay": "16Uiu2HAmNMdhYW6KdECZSNoCseM4cwKbV72BgCYZMz9SbMmxWScM"
  },
  "connected_peers": {
    "providers": [
      "0x6c27a32189016bde0d4a506805aa2b6c46295e8a"
    ]
  }
}
```

## Prometheus metrics

The node also emits a bunch of prometheus metrics. These can be potentially used to write dashboards that will help show different stats of the node. The node uses `libp2p` networking library. The `libp2p` default metrics are also available on the `/metrics` endpoint.

## Pprof endpoints

The pprof endpoints are also accessible on the node on the `/debug/pprof` endpoint. These are mainly used to observe how the node is performing for eg. the memory, CPU usage on the nodes. These are useful only for bidders who already know how to use them. Explaining them is beyond the scope of this documentation.
