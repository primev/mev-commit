# API clients

The mev-commit node provides two key APIs. The execution providers need to use the **Provider API** whereas the bidders need to use the **Bidder API** based on their role in the network.

## Providers

Execution providers will use the `provider` role of the mev-commit software to run their nodes. This will allow bidders to send bids to them to include in the blocks that they build. They will use the **Provider RPC API** to receive signed bids that are being propogated in the network. Once they get a bid, they need to communicate with the mev-commit node whether the bid has been **Accepted** or **Rejected**. If accepted, the mev-commit node will send a signed pre-confirmation to the bidders.

The API is implemented using gRPC framework. This allows two types of operations:

### RPC API

Bidders can find the protobuf file in the [repository](https://github.com/primevprotocol/mev-commit/blob/main/rpc/providerapi/v1/providerapi.proto). This can be used to generate the client for the RPC in the language of your choice. The go client is already generated in the repository. For other languagues, please follow the instructions in the [grpc documentation](https://grpc.io/docs/languages/) to generate them.

There are two main APIs
```proto
  // ReceiveBids is called by the execution provider to receive bids from the mev-commit node.
  // The mev-commit node will stream bids to the execution provider.
  rpc ReceiveBids(EmptyMessage) returns (stream Bid) {}
  // SendProcessedBids is called by the provider to send processed bids to the mev-commit node.
  // The execution provider will stream processed bids to the mev-commit node.
  rpc SendProcessedBids(stream BidResponse) returns (EmptyMessage) {}
```

The message definitions are as follows:
```proto
message Bid {
  string txn_hash = 1;
  int64 bid_amt = 2;
  int64 block_number = 3;
  bytes bid_hash = 4;
};

message BidResponse {
  bytes bid_hash = 1;
  enum Status {
    STATUS_UNSPECIFIED = 0;
    STATUS_ACCEPTED = 1;
    STATUS_REJECTED = 2;
  }
  Status status = 2;
};

```

### HTTP API

The same API is also available on the HTTP port configured on the node. Please go through the [API docs](https://mev-commit-docs.s3.amazonaws.com/provider.html) to understand the usage.

An [example client](https://github.com/primevprotocol/mev-commit/tree/main/examples/provideremulator) is implemented in the repository. This is mainly to demostrate how to write the client integrated in the provider's environment. The client blindly accepts each bid that it receives, however the provider needs to implement custom logic here to make the decision.

## Bidders

Bidders will use the `bidder` role of the mev-commit software to run their nodes. With this role, the node provides the **Bidder API** to submit bids to the network. The mev-commit node will sign the bid before it sends it to different providers that are accepting bids. On the response, bidders will get pre-confirmations from the providers if the bid is accepted. This is a streaming response and bidders are expected to keep the connection alive till all the preconfirmations are received by the node.

The API is implemented using gRPC framework. This allows two types of operations:

### RPC API

Bidders can find the protobuf file in the [repository](https://github.com/primevprotocol/mev-commit/blob/main/rpc/bidderapi/v1/bidderapi.proto). This can be used to generate the client for the RPC in the language of your choice. The go client is already generated in the repository. For other languagues, please follow the instructions in the [grpc documentation](https://grpc.io/docs/languages/) to generate them.

The API available is:
```proto
  rpc SendBid(Bid) returns (stream PreConfirmation)
```

The message definitions are as follows:
```proto
message Bid {
  string tx_hash = 1;
  int64 amount = 2;
  int64 block_number = 3;
};

message PreConfirmation {
  string tx_hash = 1;
  int64 amount = 2;
  int64 block_number = 3;
  string bid_digest = 4;
  string bid_signature = 5;
  string pre_confirmation_digest = 6;
  string pre_confirmation_signature = 7;
};
```

### HTTP API

The same API is also available on the HTTP port configured on the node. Please go through the [API docs](https://mev-commit-docs.s3.amazonaws.com/bidder.html) to understand the usage.

An [example CLI application](https://github.com/primevprotocol/mev-commit/tree/main/examples/biddercli) is implemented in the repository. This is mainly to demostrate how to integrate with the RPC API.


