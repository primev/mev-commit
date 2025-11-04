# Contract Upgrade Structure

## Complete Example: ProviderRegistry Upgrade Journey

### Timeline of Changes

#### T0: Initial Deployment
```
contracts/
├── core/
│   ├── ProviderRegistry.sol
│   └── ProviderRegistryStorage.sol
└── upgrades/
    └── (empty)
```

#### T1: First Upgrade (V1 → V2)
```
contracts/
├── core/
│   ├── ProviderRegistryV2.sol       # Current
│   └── ProviderRegistryV2Storage.sol
└── upgrades/
    └── core/
        ├── ProviderRegistry.sol      # V1 archived
        └── ProviderRegistryStorage.sol
```

#### T2: Second Upgrade (V2 → V3)
```
contracts/
├── core/
│   ├── ProviderRegistryV3.sol       # Current
│   └── ProviderRegistryV3Storage.sol
└── upgrades/
    └── core/
        ├── ProviderRegistry.sol      # V1 archived
        ├── ProviderRegistryStorage.sol
        ├── ProviderRegistryV2.sol   # V2 archived
        └── ProviderRegistryV2Storage.sol
```

---

## Rules & Guidelines

### 1. Version Naming
- **V1 (Initial):** No suffix - `Contract.sol`
- **V2+:** Explicit suffix - `ContractV2.sol`, `ContractV3.sol`, etc.

### 2. Current Version Location
- Always in the **feature folder root**
- Always has the highest version number suffix

### 3. Previous Versions Location
- Always in the **centralized `contracts/upgrades/` folder**
- Maintains the **same subdirectory structure** as the original location
- Example: `core/ProviderRegistry.sol` → `upgrades/core/ProviderRegistry.sol`

### 4. Storage Contracts
- Follow the same versioning pattern as their implementation
- Move to `contracts/upgrades/` alongside their implementation
- Preserve the same subdirectory structure

### 5. Upgrade Process
1. Create new versioned contract (e.g., `ContractV3.sol`) in feature folder root
2. Move previous version to `contracts/upgrades/[feature-folder]/`
3. Ensure subdirectory structure matches original location
4. Update imports in dependent contracts
5. Update tests to reference new version
6. Update upgrade scripts

### 6. OpenZeppelin Annotation
Always include the reference contract annotation:
```solidity
/// @custom:oz-upgrades-from ProviderRegistryV2
contract ProviderRegistryV3 is ...
```

---

## Multi-Contract Example: validator-registry/

Shows multiple contracts at different version stages:

```
contracts/
├── validator-registry/
│   ├── ValidatorOptInHub.sol        # V1 (no suffix, never upgraded)
│   ├── ValidatorOptInHubStorage.sol
│   ├── MevCommitAVSV2.sol           # V2 (current, was upgraded once)
│   ├── MevCommitAVSV2Storage.sol
│   ├── rewards/
│   │   ├── RewardManagerV3.sol          # V3 (current, upgraded twice)
│   │   └── RewardManagerV3Storage.sol
│   └── avs/
│       └── MevCommitAVSV2.sol
└── upgrades/
    └── validator-registry/
        ├── avs/
        │   ├── MevCommitAVS.sol         # V1 archived
        │   └── MevCommitAVSStorage.sol
        └── rewards/
            ├── RewardManager.sol        # V1 archived
            ├── RewardManagerStorage.sol
            ├── RewardManagerV2.sol      # V2 archived
            └── RewardManagerV2Storage.sol
```

---

## Benefits of This Structure

1. **Clear Version History:** Easy to see all previous versions
2. **Current Version Obvious:** Highest version in root = current
3. **Organized Archives:** All old code in dedicated upgrades/ folder
4. **Git-Friendly:** Easier to track changes per version
5. **Team Clarity:** New developers immediately see current vs. archived

---

