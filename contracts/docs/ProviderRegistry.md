# Solidity API

## ProviderRegistry

This contract is for provider registry and staking.

### PRECISION

```solidity
uint256 PRECISION
```

_For improved precision_

### PERCENT

```solidity
uint256 PERCENT
```

### minStake

```solidity
uint256 minStake
```

_Minimum stake required for registration_

### feePercent

```solidity
uint16 feePercent
```

_Fee percent that would be taken by protocol when provider is slashed_

### feeRecipientAmount

```solidity
uint256 feeRecipientAmount
```

_Amount assigned to feeRecipient_

### preConfirmationsContract

```solidity
address preConfirmationsContract
```

_Address of the pre-confirmations contract_

### feeRecipient

```solidity
address feeRecipient
```

_Fee recipient_

### providerRegistered

```solidity
mapping(address => bool) providerRegistered
```

_Mapping from provider address to whether they are registered or not_

### providerStakes

```solidity
mapping(address => uint256) providerStakes
```

_Mapping from provider addresses to their staked amount_

### bidderAmount

```solidity
mapping(address => uint256) bidderAmount
```

_Amount assigned to bidders_

### ProviderRegistered

```solidity
event ProviderRegistered(address provider, uint256 stakedAmount)
```

_Event for provider registration_

### FundsDeposited

```solidity
event FundsDeposited(address provider, uint256 amount)
```

_Event for depositing funds_

### FundsSlashed

```solidity
event FundsSlashed(address provider, uint256 amount)
```

_Event for slashing funds_

### FundsRewarded

```solidity
event FundsRewarded(address provider, uint256 amount)
```

_Event for rewarding funds_

### fallback

```solidity
fallback() external payable
```

_Fallback function to revert all calls, ensuring no unintended interactions._

### receive

```solidity
receive() external payable
```

_Receive function is disabled for this contract to prevent unintended interactions.
Should be removed from here in case the registerAndStake function becomes more complex_

### constructor

```solidity
constructor(uint256 _minStake, address _feeRecipient, uint16 _feePercent) public
```

_Constructor to initialize the contract with a minimum stake requirement._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _minStake | uint256 | The minimum stake required for provider registration. |
| _feeRecipient | address | The address that receives fee |
| _feePercent | uint16 | The fee percentage for protocol |

### onlyPreConfirmationEngine

```solidity
modifier onlyPreConfirmationEngine()
```

_Modifier to restrict a function to only be callable by the pre-confirmations contract._

### setPreconfirmationsContract

```solidity
function setPreconfirmationsContract(address contractAddress) external
```

_Sets the pre-confirmations contract address. Can only be called by the owner._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| contractAddress | address | The address of the pre-confirmations contract. |

### registerAndStake

```solidity
function registerAndStake() public payable
```

_Register and stake function for providers._

### checkStake

```solidity
function checkStake(address provider) external view returns (uint256)
```

_Check the stake of a provider._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| provider | address | The address of the provider. |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| [0] | uint256 | The staked amount for the provider. |

### depositFunds

```solidity
function depositFunds() external payable
```

_Deposit more funds into the provider's stake._

### slash

```solidity
function slash(uint256 amt, address provider, address payable bidder) external
```

_Slash funds from the provider and send the slashed amount to the bidder.
reenterancy not necessary but still putting here for precaution_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| amt | uint256 | The amount to slash from the provider's stake. |
| provider | address | The address of the provider. |
| bidder | address payable | The address to transfer the slashed funds to. |

### setNewFeeRecipient

```solidity
function setNewFeeRecipient(address newFeeRecipient) external
```

Sets the new fee recipient

_onlyOwner restriction_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| newFeeRecipient | address | The address to transfer the slashed funds to. |

### setNewFeePercent

```solidity
function setNewFeePercent(uint16 newFeePercent) external
```

Sets the new fee recipient

_onlyOwner restriction_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| newFeePercent | uint16 | this is the new fee percent |

### withdrawFeeRecipientAmount

```solidity
function withdrawFeeRecipientAmount() external
```

### withdrawBidderAmount

```solidity
function withdrawBidderAmount(address bidder) external
```

### withdrawStakedAmount

```solidity
function withdrawStakedAmount(address payable provider) external
```

