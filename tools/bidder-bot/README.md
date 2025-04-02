# MEV-Commit Bidder Bot

A Go service that automatically places bids for block space on the MEV-Commit network. This bot monitors for upcoming Ethereum block proposers, creates bids, and sends them to the MEV-Commit bidder node.

## Architecture Overview

The bidder bot consists of several key components that work together to monitor the Ethereum network and place bids at the appropriate time.

```mermaid
flowchart TB
    User[User/Operator] --> Main
    subgraph BidderBot["Bidder Bot Service"]
        Main[Main CLI] --> Service
        Service --> Notifier
        Service --> Bidder
        Service --> BalanceChecker
        Notifier -- "Proposer Info" --> Bidder
        Bidder --> BeaconClient
    end
    
    subgraph External["External Services"]
        Bidder -- "Bids" --> BidderNode["MEV-Commit Bidder Node"]
        Notifier -- "Subscribe" --> BidderNode
        BeaconClient -- "Slot/Block Queries" --> BeaconAPI["Ethereum Beacon API"]
        Bidder -- "Transaction Creation" --> L1Client["Ethereum L1 Client"]
        BalanceChecker -- "Check L1 Balance" --> L1Client
        BalanceChecker -- "Check Settlement Balance" --> SettlementRPC["MEV-Commit Settlement"]
    end

    classDef primary fill:#c6e6ff,stroke:#1a75ff,stroke-width:2px
    classDef secondary fill:#e1f0ff,stroke:#1a75ff,stroke-width:1px
    classDef external fill:#f0f0f0,stroke:#555555,stroke-width:1px
    
    class BidderBot primary
    class Bidder,Notifier,BeaconClient,Service,Main secondary
    class External,BidderNode,BeaconAPI,L1Client,SettlementRPC external
```

## Bidding Process Flow

The following diagram illustrates the complete flow of how the bidder bot monitors for upcoming block proposers and places bids.

```mermaid
sequenceDiagram
    participant User as User/Operator
    participant Service as Service
    participant Notifier as Notifier
    participant Bidder as Bidder
    participant BeaconClient as BeaconClient
    participant L1Client as L1 Client
    participant BidderNode as MEV-Commit Bidder Node
    participant Topology as Topology Client
    
    User->>Service: Start service
    Service->>Service: Check L1 & Settlement balances
    Service->>BidderNode: Check/Enable auto-deposit
    Service->>Notifier: Start notifier service
    Service->>Bidder: Start bidder service
    
    Notifier->>BidderNode: Subscribe to upcoming proposers
    
    Note over Notifier,BidderNode: Main notification loop
    BidderNode->>Notifier: Upcoming proposer notification
    Notifier->>Notifier: Validate if newer than last seen
    Notifier->>Bidder: Pass upcoming proposer info
    
    Note over Bidder,BeaconClient: Bidding preparation
    Bidder->>BeaconClient: Get block number for slot-2
    BeaconClient->>Bidder: Return block number
    Bidder->>Bidder: Calculate target block number
    
    Bidder->>L1Client: Get nonce & chainID
    L1Client->>Bidder: Return nonce & chainID
    Bidder->>Bidder: Create self-transfer transaction
    Bidder->>Bidder: Sign transaction
    
    Note over Bidder,BidderNode: Bid submission & monitoring
    Bidder->>BidderNode: Send bid with transaction
    Bidder->>Topology: Get topology (connected providers)
    Topology->>Bidder: Return provider list
    
    loop Wait for commitments
        BidderNode->>Bidder: Stream commitment from provider
        Bidder->>Bidder: Track received commitments
    end
    
    alt All commitments received
        Bidder->>Bidder: Log success
    else Timeout or error
        Bidder->>Bidder: Log warning/error
    end
    
```

## Component Details

The bidder bot is composed of several key components, each with distinct responsibilities.

```mermaid
flowchart LR
    subgraph Service["Service Orchestration"]
        direction TB
        Config[Service Config]
        GRPCConn[gRPC Connection]
        HealthChecker[Health Checker]
        L1Client[L1 Client]
        SettlementClient[Settlement Client]
    end
    
    subgraph Core["Core Components"]
        direction TB
        Notifier[Notifier]
        Bidder[Bidder]
        BeaconClient[Beacon Client]
        
        ProposerChan[Proposer Channel]
        Notifier --> ProposerChan
        ProposerChan --> Bidder
        Bidder --> BeaconClient
    end
    
    subgraph API["API Clients"]
        direction TB
        BidderAPI[Bidder API]
        TopologyAPI[Topology API]
        NotificationsAPI[Notifications API]
    end
    
    subgraph External["External Services"]
        direction TB
        BeaconRPC[Ethereum Beacon API]
        EthRPC[Ethereum L1 RPC]
        SettlementRPC[MEV-Commit Settlement RPC]
        BidderRPC[MEV-Commit Bidder Node RPC]
    end
    
    %% Connections
    Config --> Service
    Service --> Core
    Service --> API
    
    BeaconClient --> BeaconRPC
    L1Client --> EthRPC
    SettlementClient --> SettlementRPC
    
    Notifier --> NotificationsAPI
    Bidder --> BidderAPI
    Bidder --> TopologyAPI
    Bidder --> L1Client
    
    NotificationsAPI --> BidderRPC
    BidderAPI --> BidderRPC
    TopologyAPI --> BidderRPC
    
    classDef config fill:#ffe0b2,stroke:#e67e22,stroke-width:1px
    classDef service fill:#bbdefb,stroke:#1976d2,stroke-width:1px
    classDef core fill:#c8e6c9,stroke:#2e7d32,stroke-width:1px
    classDef api fill:#e1bee7,stroke:#7b1fa2,stroke-width:1px
    classDef external fill:#f5f5f5,stroke:#616161,stroke-width:1px
    classDef channel fill:#fff9c4,stroke:#f57f17,stroke-width:1px
    
    class Config config
    class Service,GRPCConn,HealthChecker,L1Client,SettlementClient service
    class Notifier,Bidder,BeaconClient core
    class BidderAPI,TopologyAPI,NotificationsAPI api
    class BeaconRPC,EthRPC,SettlementRPC,BidderRPC external
    class ProposerChan channel
```

```mermaid
flowchart TD
    subgraph Inputs["Input Data"]
        BeaconSlots["Beacon Chain Slots"]
        ProposerNotifications["Proposer Notifications"]
        ConfigParams["Configuration Parameters"]
    end
    
    subgraph Processing["Processing"]
        SlotProcessing["Slot to Block Mapping"]
        BidCreation["Bid Creation"]
        TxCreation["Transaction Creation"]
        CommitmentTracking["Commitment Tracking"]
    end
    
    subgraph Components["System Components"]
        BeaconClient["Beacon Client"]
        Notifier["Notifier"]
        Bidder["Bidder"]
        KeySigner["Key Signer"]
    end
    
    subgraph Outputs["Output Actions"]
        BidSubmission["Bid Submission"]
        BalanceCheck["Balance Checking"]
        AutoDeposit["Auto Deposit"]
        Logging["Logging"]
    end
    
    %% Input flows
    BeaconSlots -->|"Slot data"| BeaconClient
    ProposerNotifications -->|"Upcoming proposer info"| Notifier
    ConfigParams -->|"URLs, amounts, gas prices"| Components
    
    %% Component processing
    BeaconClient -->|"Slot/block mapping"| SlotProcessing
    Notifier -->|"Filtered proposer info"| Bidder
    SlotProcessing -->|"Target block number"| BidCreation
    KeySigner -->|"Sign transaction"| TxCreation
    
    %% Output actions
    BidCreation -->|"Bid parameters"| BidSubmission
    TxCreation -->|"Signed transaction"| BidSubmission
    BidSubmission -->|"Stream"| CommitmentTracking
    CommitmentTracking -->|"Status"| Logging
    Components -->|"Periodic checks"| BalanceCheck
    BalanceCheck -->|"If low balance"| AutoDeposit
    
    classDef input fill:#ffecb3,stroke:#ff8f00,stroke-width:1px
    classDef process fill:#b3e0ff,stroke:#0277bd,stroke-width:1px
    classDef component fill:#c8e6c9,stroke:#2e7d32,stroke-width:1px
    classDef output fill:#e1bee7,stroke:#6a1b9a,stroke-width:1px
    
    class Inputs,BeaconSlots,ProposerNotifications,ConfigParams input
    class Processing,SlotProcessing,BidCreation,TxCreation,CommitmentTracking process
    class Components,BeaconClient,Notifier,Bidder,KeySigner component
    class Outputs,BidSubmission,BalanceCheck,AutoDeposit,Logging output
```

### Key Components

1. **Service** - The main orchestrator that initializes and coordinates all other components
2. **Notifier** - Subscribes to and processes notifications about upcoming block proposers
3. **Bidder** - Handles the creation and submission of bids to the MEV-Commit network
4. **BeaconClient** - Interfaces with the Ethereum Beacon chain to retrieve slot and block information
5. **KeySigner** - Manages cryptographic signing of transactions using a keystore file

## Configuration

The bidder bot accepts the following configuration parameters:

```
--keystore-dir         Directory where keystore file is stored [required]
--keystore-password    Password to access keystore [required]
--l1-rpc-urls          URLs for L1 RPC endpoints [required]
--beacon-api-urls      URLs for Beacon API endpoints [required]
--settlement-rpc-url   URL for settlement RPC [required]
--bidder-node-rpc-url  URL for mev-commit bidder node RPC [required]
--auto-deposit-amount  Amount to auto-deposit (default: 0.1 ETH)
--bid-amount           Amount to use for each bid (default: 0.005 ETH)
--gas-tip-cap          Gas tip cap (default: 0.015 gwei)
--gas-fee-cap          Gas fee cap (default: 1 gwei)
--log-fmt              Log format: 'text' or 'json' (default: text)
--log-level            Log level: 'debug', 'info', 'warn', 'error' (default: info)
```

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Ethereum keystore file with sufficient ETH on both L1 and MEV-Commit settlement layer
- Access to Ethereum L1 and Beacon chain endpoints
- Access to MEV-Commit bidder node and settlement layer endpoints

### Building

```bash
go build -o bidder-bot ./tools/bidder-bot
```

### Running

```bash
./bidder-bot \
  --keystore-dir=/path/to/keystore \
  --keystore-password=your_password \
  --l1-rpc-urls=https://ethereum-rpc.example.com \
  --beacon-api-urls=https://ethereum-beacon.example.com \
  --settlement-rpc-url=https://settlement-rpc.example.com \
  --bidder-node-rpc-url=https://bidder-node.example.com \
  --log-level=info
```

## How It Works

```mermaid
graph TD
    Start([Start Bidder Bot]) --> CheckBalances[Check L1 & Settlement Balances]
    CheckBalances --> EnableAutoDeposit[Enable Auto-Deposit if Needed]
    EnableAutoDeposit --> StartServices[Start Notifier & Bidder Services]
    
    StartServices --> SubscribeNotifications[Subscribe to Upcoming Proposer Notifications]
    
    SubscribeNotifications --> WaitForNotification{Wait for Notification}
    
    WaitForNotification -->|Notification Received| ValidateSlot[Validate Slot is Newer]
    ValidateSlot -->|Not Valid| WaitForNotification
    ValidateSlot -->|Valid| GetPreviousSlotBlock[Get Block Number for Slot-2]
    
    GetPreviousSlotBlock --> CalculateTargetBlock[Calculate Target Block Number]
    CalculateTargetBlock --> PrepareTransaction[Prepare Self-Transfer Transaction]
    PrepareTransaction --> SignTransaction[Sign Transaction]
    
    SignTransaction --> CreateBid[Create Bid with Transaction]
    CreateBid --> SendBid[Send Bid to Bidder Node]
    
    SendBid --> GetTopology[Get Network Topology]
    GetTopology --> WaitForCommitments{Wait for Commitments}
    
    WaitForCommitments -->|Commitment Received| TrackCommitment[Track Commitment]
    TrackCommitment --> CheckAllReceived{All Received?}
    CheckAllReceived -->|No| WaitForCommitments
    CheckAllReceived -->|Yes| LogSuccess[Log Success]
    
    WaitForCommitments -->|Timeout| LogError[Log Error/Warning]
    
    LogSuccess --> WaitForNotification
    LogError --> WaitForNotification
    
    classDef start fill:#d4efdf,stroke:#27ae60,stroke-width:2px
    classDef process fill:#aed6f1,stroke:#2874a6,stroke-width:1px
    classDef decision fill:#f9d79b,stroke:#d35400,stroke-width:1px
    classDef terminal fill:#f5b7b1,stroke:#c0392b,stroke-width:1px
    
    class Start start
    class CheckBalances,EnableAutoDeposit,StartServices,SubscribeNotifications,ValidateSlot,GetPreviousSlotBlock,CalculateTargetBlock,PrepareTransaction,SignTransaction,CreateBid,SendBid,GetTopology,TrackCommitment process
    class WaitForNotification,CheckAllReceived,WaitForCommitments decision
    class LogSuccess,LogError terminal
```

1. **Service Initialization**:
   - Check balances on L1 and settlement layer
   - Set up auto-deposit if needed
   - Initialize and start the Notifier and Bidder components

2. **Notification Subscription**:
   - The Notifier subscribes to upcoming proposer notifications from the MEV-Commit network
   - When a notification arrives, it's validated and passed to the Bidder

3. **Bid Preparation**:
   - The Bidder uses the BeaconClient to get the corresponding block number for the target slot
   - A self-transfer Ethereum transaction is created and signed with the configured keystore
   - This transaction is included in the bid to demonstrate user intent

4. **Bid Submission**:
   - The bid is submitted to the MEV-Commit Bidder Node with the transaction data
   - The Bidder tracks commitments from providers in response to the bid
   - Success or failures are logged appropriately