Summary
 - [arbitrary-send-eth](#arbitrary-send-eth) (4 results) (High)
 - [incorrect-exp](#incorrect-exp) (1 results) (High)
 - [divide-before-multiply](#divide-before-multiply) (8 results) (Medium)
 - [locked-ether](#locked-ether) (1 results) (Medium)
 - [events-access](#events-access) (3 results) (Low)
 - [missing-zero-check](#missing-zero-check) (10 results) (Low)
 - [assembly](#assembly) (5 results) (Informational)
 - [solc-version](#solc-version) (14 results) (Informational)
 - [low-level-calls](#low-level-calls) (7 results) (Informational)
 - [naming-convention](#naming-convention) (12 results) (Informational)
 - [immutable-states](#immutable-states) (4 results) (Optimization)
## arbitrary-send-eth
Impact: High
Confidence: Medium
 - [ ] ID-0
[BidderRegistry.withdrawProviderAmount(address)](contracts/BidderRegistry.sol#L186-L195) sends eth to arbitrary bidder
	Dangerous calls:
	- [(success) = provider.call{value: amount}()](contracts/BidderRegistry.sol#L193)

contracts/BidderRegistry.sol#L186-L195


 - [ ] ID-1
[ProviderRegistry.withdrawBidderAmount(address)](contracts/ProviderRegistry.sol#L191-L198) sends eth to arbitrary bidder
	Dangerous calls:
	- [(success) = bidder.call{value: bidderAmount[bidder]}()](contracts/ProviderRegistry.sol#L196)

contracts/ProviderRegistry.sol#L191-L198


 - [ ] ID-2
[ProviderRegistry.withdrawFeeRecipientAmount()](contracts/ProviderRegistry.sol#L185-L189) sends eth to arbitrary bidder
	Dangerous calls:
	- [(successFee) = feeRecipient.call{value: feeRecipientAmount}()](contracts/ProviderRegistry.sol#L187)

contracts/ProviderRegistry.sol#L185-L189


 - [ ] ID-3
[BidderRegistry.withdrawFeeRecipientAmount()](contracts/BidderRegistry.sol#L178-L184) sends eth to arbitrary bidder
	Dangerous calls:
	- [(successFee) = feeRecipient.call{value: amount}()](contracts/BidderRegistry.sol#L182)

contracts/BidderRegistry.sol#L178-L184


## incorrect-exp
Impact: High
Confidence: Medium
 - [ ] ID-4
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) has bitwise-xor operator ^ instead of the exponentiation operator **: 
	 - [inverse = (3 * denominator) ^ 2](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L184)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


## divide-before-multiply
Impact: Medium
Confidence: Medium
 - [ ] ID-5
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L190)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-6
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L193)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-7
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L188)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-8
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse = (3 * denominator) ^ 2](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L184)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-9
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [prod0 = prod0 / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L172)
	- [result = prod0 * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L199)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-10
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L192)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-11
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L191)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-12
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) performs a multiplication on the result of a division:
	- [denominator = denominator / twos](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L169)
	- [inverse *= 2 - denominator * inverse](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L189)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


## locked-ether
Impact: Medium
Confidence: High
 - [ ] ID-13
Contract locking ether found:
	Contract [PreConfCommitmentStore](contracts/PreConfirmations.sol#L16-L448) has payable functions:
	 - [PreConfCommitmentStore.fallback()](contracts/PreConfirmations.sol#L99-L101)
	 - [PreConfCommitmentStore.receive()](contracts/PreConfirmations.sol#L106-L108)
	But does not have a function to withdraw the ether

contracts/PreConfirmations.sol#L16-L448


## events-access
Impact: Low
Confidence: Medium
 - [ ] ID-14
[PreConfCommitmentStore.updateOracle(address)](contracts/PreConfirmations.sol#L393-L395) should emit an event for: 
	- [oracle = newOracle](contracts/PreConfirmations.sol#L394) 

contracts/PreConfirmations.sol#L393-L395


 - [ ] ID-15
[BidderRegistry.setPreconfirmationsContract(address)](contracts/BidderRegistry.sol#L95-L103) should emit an event for: 
	- [preConfirmationsContract = contractAddress](contracts/BidderRegistry.sol#L102) 

contracts/BidderRegistry.sol#L95-L103


 - [ ] ID-16
[ProviderRegistry.setPreconfirmationsContract(address)](contracts/ProviderRegistry.sol#L99-L107) should emit an event for: 
	- [preConfirmationsContract = contractAddress](contracts/ProviderRegistry.sol#L106) 

contracts/ProviderRegistry.sol#L99-L107


## missing-zero-check
Impact: Low
Confidence: Medium
 - [ ] ID-17
[BidderRegistry.withdrawProviderAmount(address).provider](contracts/BidderRegistry.sol#L187) lacks a zero-check on :
		- [(success) = provider.call{value: amount}()](contracts/BidderRegistry.sol#L193)

contracts/BidderRegistry.sol#L187


 - [ ] ID-18
[ProviderRegistry.constructor(uint256,address,uint16)._feeRecipient](contracts/ProviderRegistry.sol#L76) lacks a zero-check on :
		- [feeRecipient = _feeRecipient](contracts/ProviderRegistry.sol#L80)

contracts/ProviderRegistry.sol#L76


 - [ ] ID-19
[BidderRegistry.constructor(uint256,address,uint16)._feeRecipient](contracts/BidderRegistry.sol#L72) lacks a zero-check on :
		- [feeRecipient = _feeRecipient](contracts/BidderRegistry.sol#L76)

contracts/BidderRegistry.sol#L72


 - [ ] ID-20
[BidderRegistry.setNewFeeRecipient(address).newFeeRecipient](contracts/BidderRegistry.sol#L165) lacks a zero-check on :
		- [feeRecipient = newFeeRecipient](contracts/BidderRegistry.sol#L166)

contracts/BidderRegistry.sol#L165


 - [ ] ID-21
[PreConfCommitmentStore.updateOracle(address).newOracle](contracts/PreConfirmations.sol#L393) lacks a zero-check on :
		- [oracle = newOracle](contracts/PreConfirmations.sol#L394)

contracts/PreConfirmations.sol#L393


 - [ ] ID-22
[BidderRegistry.setPreconfirmationsContract(address).contractAddress](contracts/BidderRegistry.sol#L96) lacks a zero-check on :
		- [preConfirmationsContract = contractAddress](contracts/BidderRegistry.sol#L102)

contracts/BidderRegistry.sol#L96


 - [ ] ID-23
[ProviderRegistry.setPreconfirmationsContract(address).contractAddress](contracts/ProviderRegistry.sol#L100) lacks a zero-check on :
		- [preConfirmationsContract = contractAddress](contracts/ProviderRegistry.sol#L106)

contracts/ProviderRegistry.sol#L100


 - [ ] ID-24
[ProviderRegistry.setNewFeeRecipient(address).newFeeRecipient](contracts/ProviderRegistry.sol#L172) lacks a zero-check on :
		- [feeRecipient = newFeeRecipient](contracts/ProviderRegistry.sol#L173)

contracts/ProviderRegistry.sol#L172


 - [ ] ID-25
[BidderRegistry.withdrawProtocolFee(address).bidder](contracts/BidderRegistry.sol#L208) lacks a zero-check on :
		- [(success) = bidder.call{value: _protocolFeeAmount}()](contracts/BidderRegistry.sol#L214)

contracts/BidderRegistry.sol#L208


 - [ ] ID-26
[PreConfCommitmentStore.constructor(address,address,address)._oracle](contracts/PreConfirmations.sol#L127) lacks a zero-check on :
		- [oracle = _oracle](contracts/PreConfirmations.sol#L129)

contracts/PreConfirmations.sol#L127


## assembly
Impact: Informational
Confidence: High
 - [ ] ID-27
[ECDSA.tryRecover(bytes32,bytes)](lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol#L56-L73) uses assembly
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol#L64-L68)

lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol#L56-L73


 - [ ] ID-28
[Strings.toString(uint256)](lib/openzeppelin-contracts/contracts/utils/Strings.sol#L24-L44) uses assembly
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/Strings.sol#L30-L32)
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/Strings.sol#L36-L38)

lib/openzeppelin-contracts/contracts/utils/Strings.sol#L24-L44


 - [ ] ID-29
[Math.mulDiv(uint256,uint256,uint256)](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202) uses assembly
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L130-L133)
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L154-L161)
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L167-L176)

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L123-L202


 - [ ] ID-30
[MessageHashUtils.toEthSignedMessageHash(bytes32)](lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L30-L37) uses assembly
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L32-L36)

lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L30-L37


 - [ ] ID-31
[MessageHashUtils.toTypedDataHash(bytes32,bytes32)](lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L76-L85) uses assembly
	- [INLINE ASM](lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L78-L84)

lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L76-L85


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-32
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/access/Ownable.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/access/Ownable.sol#L4


 - [ ] ID-33
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/math/Math.sol#L4


 - [ ] ID-34
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/cryptography/MessageHashUtils.sol#L4


 - [ ] ID-35
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/math/SignedMath.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/math/SignedMath.sol#L4


 - [ ] ID-36
solc-0.8.20 is not recommended for deployment

 - [ ] ID-37
Pragma version[^0.8.20](contracts/ProviderRegistry.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

contracts/ProviderRegistry.sol#L2


 - [ ] ID-38
Pragma version[^0.8.20](contracts/BidderRegistry.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

contracts/BidderRegistry.sol#L2


 - [ ] ID-39
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol#L4


 - [ ] ID-40
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/Context.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/Context.sol#L4


 - [ ] ID-41
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/ReentrancyGuard.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/ReentrancyGuard.sol#L4


 - [ ] ID-42
Pragma version[^0.8.20](contracts/interfaces/IProviderRegistry.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

contracts/interfaces/IProviderRegistry.sol#L2


 - [ ] ID-43
Pragma version[^0.8.20](contracts/PreConfirmations.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

contracts/PreConfirmations.sol#L2


 - [ ] ID-44
Pragma version[^0.8.20](contracts/interfaces/IBidderRegistry.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

contracts/interfaces/IBidderRegistry.sol#L2


 - [ ] ID-45
Pragma version[^0.8.20](lib/openzeppelin-contracts/contracts/utils/Strings.sol#L4) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.

lib/openzeppelin-contracts/contracts/utils/Strings.sol#L4


## low-level-calls
Impact: Informational
Confidence: High
 - [ ] ID-46
Low level call in [ProviderRegistry.withdrawFeeRecipientAmount()](contracts/ProviderRegistry.sol#L185-L189):
	- [(successFee) = feeRecipient.call{value: feeRecipientAmount}()](contracts/ProviderRegistry.sol#L187)

contracts/ProviderRegistry.sol#L185-L189


 - [ ] ID-47
Low level call in [BidderRegistry.withdrawStakedAmount(address)](contracts/BidderRegistry.sol#L197-L205):
	- [(success) = bidder.call{value: stake}()](contracts/BidderRegistry.sol#L203)

contracts/BidderRegistry.sol#L197-L205


 - [ ] ID-48
Low level call in [BidderRegistry.withdrawProviderAmount(address)](contracts/BidderRegistry.sol#L186-L195):
	- [(success) = provider.call{value: amount}()](contracts/BidderRegistry.sol#L193)

contracts/BidderRegistry.sol#L186-L195


 - [ ] ID-49
Low level call in [BidderRegistry.withdrawProtocolFee(address)](contracts/BidderRegistry.sol#L207-L216):
	- [(success) = bidder.call{value: _protocolFeeAmount}()](contracts/BidderRegistry.sol#L214)

contracts/BidderRegistry.sol#L207-L216


 - [ ] ID-50
Low level call in [ProviderRegistry.withdrawStakedAmount(address)](contracts/ProviderRegistry.sol#L200-L223):
	- [(success) = provider.call{value: stake}()](contracts/ProviderRegistry.sol#L221)

contracts/ProviderRegistry.sol#L200-L223


 - [ ] ID-51
Low level call in [BidderRegistry.withdrawFeeRecipientAmount()](contracts/BidderRegistry.sol#L178-L184):
	- [(successFee) = feeRecipient.call{value: amount}()](contracts/BidderRegistry.sol#L182)

contracts/BidderRegistry.sol#L178-L184


 - [ ] ID-52
Low level call in [ProviderRegistry.withdrawBidderAmount(address)](contracts/ProviderRegistry.sol#L191-L198):
	- [(success) = bidder.call{value: bidderAmount[bidder]}()](contracts/ProviderRegistry.sol#L196)

contracts/ProviderRegistry.sol#L191-L198


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-53
Parameter [PreConfCommitmentStore.getPreConfHash(string,uint64,uint64,bytes32,string)._txnHash](contracts/PreConfirmations.sol#L196) is not in mixedCase

contracts/PreConfirmations.sol#L196


 - [ ] ID-54
Variable [PreConfCommitmentStore.DOMAIN_SEPARATOR_BID](contracts/PreConfirmations.sol#L39) is not in mixedCase

contracts/PreConfirmations.sol#L39


 - [ ] ID-55
Function [PreConfCommitmentStore._bytesToHexString(bytes)](contracts/PreConfirmations.sol#L437-L447) is not in mixedCase

contracts/PreConfirmations.sol#L437-L447


 - [ ] ID-56
Parameter [PreConfCommitmentStore.getPreConfHash(string,uint64,uint64,bytes32,string)._blockNumber](contracts/PreConfirmations.sol#L198) is not in mixedCase

contracts/PreConfirmations.sol#L198


 - [ ] ID-57
Variable [PreConfCommitmentStore.DOMAIN_SEPARATOR_PRECONF](contracts/PreConfirmations.sol#L36) is not in mixedCase

contracts/PreConfirmations.sol#L36


 - [ ] ID-58
Parameter [PreConfCommitmentStore.getPreConfHash(string,uint64,uint64,bytes32,string)._bid](contracts/PreConfirmations.sol#L197) is not in mixedCase

contracts/PreConfirmations.sol#L197


 - [ ] ID-59
Parameter [PreConfCommitmentStore.getBidHash(string,uint64,uint64)._txnHash](contracts/PreConfirmations.sol#L169) is not in mixedCase

contracts/PreConfirmations.sol#L169


 - [ ] ID-60
Parameter [PreConfCommitmentStore._bytesToHexString(bytes)._bytes](contracts/PreConfirmations.sol#L438) is not in mixedCase

contracts/PreConfirmations.sol#L438


 - [ ] ID-61
Parameter [PreConfCommitmentStore.getBidHash(string,uint64,uint64)._blockNumber](contracts/PreConfirmations.sol#L171) is not in mixedCase

contracts/PreConfirmations.sol#L171


 - [ ] ID-62
Parameter [PreConfCommitmentStore.getPreConfHash(string,uint64,uint64,bytes32,string)._bidHash](contracts/PreConfirmations.sol#L199) is not in mixedCase

contracts/PreConfirmations.sol#L199


 - [ ] ID-63
Parameter [PreConfCommitmentStore.getBidHash(string,uint64,uint64)._bid](contracts/PreConfirmations.sol#L170) is not in mixedCase

contracts/PreConfirmations.sol#L170


 - [ ] ID-64
Parameter [PreConfCommitmentStore.getPreConfHash(string,uint64,uint64,bytes32,string)._bidSignature](contracts/PreConfirmations.sol#L200) is not in mixedCase

contracts/PreConfirmations.sol#L200


## immutable-states
Impact: Optimization
Confidence: High
 - [ ] ID-65
[PreConfCommitmentStore.DOMAIN_SEPARATOR_BID](contracts/PreConfirmations.sol#L39) should be immutable 

contracts/PreConfirmations.sol#L39


 - [ ] ID-66
[PreConfCommitmentStore.DOMAIN_SEPARATOR_PRECONF](contracts/PreConfirmations.sol#L36) should be immutable 

contracts/PreConfirmations.sol#L36


 - [ ] ID-67
[BidderRegistry.minStake](contracts/BidderRegistry.sol#L20) should be immutable 

contracts/BidderRegistry.sol#L20


 - [ ] ID-68
[ProviderRegistry.minStake](contracts/ProviderRegistry.sol#L18) should be immutable 

contracts/ProviderRegistry.sol#L18


