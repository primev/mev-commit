
**# MEV Commit Geth Helm Chart - Insecure Installation Guide**

⚠️ ****WARNING****: This guide demonstrates insecure installation methods intended ****ONLY**** for local development and testing. ****Never**** use these methods in production environments.

For quickly setting up a new devnet, you can use our helper script:
### From the repository root:
`./utils/setup-devnet.sh`

---

**## Quick Start**

Deploy a complete PoA network with bootnodes, signers, and members using inline secrets.

---

**### Step 1: Deploy Bootnodes**

Create bootnode values file:

```
# values-bootnode.yaml
role: "bootnode"
replicas: 2

persistence:
  size: 25Gi

resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "5Gi"
    cpu: "2000m"
```

Deploy bootnodes:

```
# Deploy bootnodes
helm install poa-bootnode . -f values-bootnode.yaml

# Wait for bootnodes to be ready. Until you see the pod 2/2 like this
kubectl get pods -l=app.kubernetes.io/instance=poa-bootnode
NAME                                   READY STATUS RESTARTS AGE
poa-bootnode-ethereum-poa-bootnode-0   1/1 Running 0 169m
```

---

**### Step 2: Get Bootnode Information**

```
# Get enode addresses
NODE_ID_1=$(kubectl logs -f poa-bootnode-ethereum-poa-bootnode-0 -c init-nodekey | grep "Node ID:" | cut -d':' -f2 | xargs)
NODE_ID_2=$(kubectl logs -f poa-bootnode-ethereum-poa-bootnode-1 -c init-nodekey | grep "Node ID:" | cut -d':' -f2 | xargs)

# Replace the public IP with the service name
BOOTNODE_1="enode://${NODE_ID_1}@poa-bootnode-ethereum-poa-bootnode:30303"
BOOTNODE_2="enode://${NODE_ID_2}@poa-bootnode-ethereum-poa-bootnode:30303"

echo "Bootnode 1: $BOOTNODE_1"
echo "Bootnode 2: $BOOTNODE_2"
```

---

**### Step 3: Deploy Signer (Insecure)**

Create `values-signer-insecure.yaml` with inline secrets:

```
role: "signer"
replicas: 1
genesisUrl: "<PLACE GENESIS URL HERE>"
chainId: ""

signer:
  address: "0xYOUR_SIGNER_ADDRESS_HERE"
  nodeBootstrapMethod: "signerKeystore"

  password:
    secretName: devnet-password
    value: "your-test-password" # INSECURE
    externalSecret:
      enabled: false

  # Singer private key is not needed as `nodeBootstrapMethod` is set to `signerKeystore`
  signerPrivateKey:
    secretName: devnet-signer-key
    value: "0xYOUR_PRIVATE_KEY_HERE" # INSECURE
    externalSecret:
      enabled: false

  signerKeystore:
    url: "<PLACE KEYSTORE URL HERE>"

nodeConfig:
  staticPeers:
    - "BOOTNODE_1"
    - "BOOTNODE_2"

persistence:
  size: 25Gi

resources:
  requests:
    memory: "256Mi"
    cpu: "100m"
  limits:
    memory: "5Gi"
    cpu: "2000m"
```

Deploy the signer:

```
helm install poa-signer . -f values-signer-insecure.yaml
```

---

**### Step 4: Deploy Member**

Create `values-member.yaml`:

```
role: "member"
replicas: 1

nodeConfig:
  staticPeers:
    - "BOOTNODE_1_ENODE_HERE"
    - "BOOTNODE_2_ENODE_HERE"

persistence:
  size: 10Gi

resources:
  requests:
    memory: "512Mi"
    cpu: "250m"
  limits:
    memory: "1Gi"
    cpu: "500m"
```

Deploy the member:

```
helm install poa-member . -f values-member.yaml
```

---

**## Signer Configuration Options**

**### Option A: Private Key Method**

```
signer:
  nodeBootstrapMethod: "signerPrivateKey"
  password:
    value: "test-password-123"
    externalSecret:
      enabled: false
  signerPrivateKey:
    value: "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    externalSecret:
      enabled: false
```

**### Option B: Keystore Method**

```
signer:
  nodeBootstrapMethod: "signerKeystore"
  password:
    value: "test-password-123"
    externalSecret:
      enabled: false
  signerKeystore:
    url: "https://raw.githubusercontent.com/yourrepo/keystore/UTC--timestamp--address"
  signerPrivateKey:
    value: ""
    externalSecret:
      enabled: false
```

---

**## Complete Example of** `signerKeystore` **method**

```
# values-signer-example.yaml
role: signer

genesisUrl: ""
signer:
  # Signer address (required for both methods)
  address: "0xc22f7d9446995d7a6c932c69d7af118bd65a6cea"

  # Bootstrap method selection
  nodeBootstrapMethod: "signerKeystore" # Options: "signerPrivateKey" or "signerKeystore"

  # Password configuration (required for both methods)
  password:
    secretName: devnet-password
    # Insecure: provide password directly in values
    value: "helloworld123"
    externalSecret:
      enabled: false
      awsSecretName: devnet-password
      refreshInterval: 1h
      secretStoreRef:
        name: aws-cluster-secret-store
        kind: ClusterSecretStore

  # Keystore method configuration (used when nodeBootstrapMethod: "signerKeystore")
  signerKeystore:
    url: "https://path/to/keystore/UTC--2025-05-12T17-38-30.308195314Z--c22f7d9446995d7a6c932c69d7af118bd65a6cea"

  # Private key method configuration (used when nodeBootstrapMethod: "signerPrivateKey")
  signerPrivateKey:
    secretName: devnet-signer-key
    # Insecure: provide password directly in values
    value: ""
    externalSecret:
      enabled: true
      awsSecretName: devnet-signer-key
      refreshInterval: 1h
      secretStoreRef:
        name: aws-cluster-secret-store
        kind: ClusterSecretStore

# Rest of the configuration remains the same
chainId:

ports:
  p2p: 30303
  http: 8545
  ws: 8546

terminationGracePeriodSeconds: 60

probes:
  enabled: false


nodeConfig:
  httpApi: "eth,net,web3,clique"
  wsApi: "eth,net,web3,clique"
  bootstrapNodes:
    - "enode://056...8311@poa-bootnode-ethereum-poa-bootnode:30303"
  staticPeers:
    - "enode://056...58311@poa-bootnode-ethereum-poa-bootnode:30303"
```

---

**## Verification**

```
# Check signer is sealing new blocks

~ kubectl logs -f poa-signer-ethereum-poa-signer-0 -c geth

INFO [05-12|21:14:45.986] Commit new sealing work                  number=2,060,963 sealhash=336bf7..0b76ff txs=0 gas=0 fees=0 elapsed="261.517µs"
INFO [05-12|21:14:45.991] Successfully sealed new block            number=2,060,963 sealhash=336bf7..0b76ff hash=779a4d..a87984 elapsed=5.216ms
INFO [05-12|21:14:45.992] Commit new sealing work                  number=2,060,964 sealhash=55ea1f..712683 txs=0 gas=0 fees=0 elapsed="266.382µs"
INFO [05-12|21:14:45.996] Successfully sealed new block            number=2,060,964 sealhash=55ea1f..712683 hash=6d89f5..4e656a elapsed=4.086ms
INFO [05-12|21:14:45.996] Commit new sealing work                  number=2,060,965 sealhash=0c2be8..ff34f5 txs=0 gas=0 fees=0 elapsed="255.769µs"
INFO [05-12|21:14:46.001] Successfully sealed new block            number=2,060,965 sealhash=0c2be8..ff34f5 hash=dd12aa..098c3f elapsed=5.133ms
INFO [05-12|21:14:46.002] Commit new sealing work                  number=2,060,966 sealhash=614953..249141 txs=0 gas=0 fees=0 elapsed="288.131µs"
INFO [05-12|21:14:46.006] Successfully sealed new block            number=2,060,966 sealhash=614953..249141 hash=39151b..20da37 elapsed=4.041ms
INFO [05-12|21:14:46.006] Commit new sealing work                  number=2,060,967 sealhash=45c3cb..b2bfe8 txs=0 gas=0 fees=0 elapsed="233.053µs"

# Check if bootnodes are importing the new blocks

~ kubectl logs -f poa-bootnode-ethereum-poa-bootnode-0 -c geth

INFO [05-12|21:15:34.860] Imported new chain segment               number=2,070,234 hash=c6583b..d0e567 blocks=190 txs=0 mgas=0.000 elapsed=174.198ms mgasps=0.000 triedirty=0.00B
INFO [05-12|21:15:35.211] Imported new chain segment               number=2,070,661 hash=429baa..50b4d1 blocks=427 txs=0 mgas=0.000 elapsed=350.093ms mgasps=0.000 triedirty=0.00B
INFO [05-12|21:15:37.936] Imported new chain segment               number=2,070,851 hash=297dc9..cb47c9 blocks=190 txs=0 mgas=0.000 elapsed=163.110ms mgasps=0.000 triedirty=0.00B

```

---

**## Important Notes**

- ALWAYS set `externalSecret.enabled: false` for insecure mode
- ALWAYS provide inline values when ESO is disabled
- Both password and private key must have inline values (unless using keystore method)
- For keystore method, only password needs inline value

---

**## Common Errors**

**### Error: Cannot use both externalSecret and inline value**

```
# WRONG - Don't do this
password:
  value: "test123"
  externalSecret:
    enabled: true # This causes the error
```

**### Error: ESO is disabled but no inline value provided**

```
# WRONG - Missing value
password:
  externalSecret:
    enabled: false
  # value: "test123" # This is missing!
```

---

**## Security Reminder**

This configuration ****exposes sensitive data**** in:

- Helm values files
- Kubernetes secrets (base64 encoded, not encrypted)
- Helm release history

_> 🔐_ ****Only use for development and testing!****
