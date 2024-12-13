package preconfencryptor

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/x/keysigner"
	"google.golang.org/protobuf/proto"
)

var (
	ErrAlreadySignedBid             = errors.New("already contains hash or signature")
	ErrMissingHashSignature         = errors.New("missing hash or signature")
	ErrInvalidSignature             = errors.New("signature is not valid")
	ErrInvalidHash                  = errors.New("bidhash doesn't match bid payload")
	ErrAlreadySignedPreConfirmation = errors.New("preConfirmation is already hashed or signed")
	ErrInvalidCommitment            = errors.New("commitment is incorrect")
	ErrMissingRequiredFields        = errors.New("missing required fields")
	ErrNoAesKeyFound                = errors.New("no AES key found for bidder")
	ErrInvalidBidAmt                = errors.New("invalid bid amount")
	ErrBidNotFound                  = errors.New("bid not found")
)

var (
	EIP712BidTypeHash        = crypto.Keccak256Hash([]byte("PreConfBid(string txnHash,string revertingTxHashes,uint256 bidAmt,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,uint256 slashAmount)"))
	EIP712CommitmentTypeHash = crypto.Keccak256Hash([]byte("OpenedCommitment(string txnHash,string revertingTxHashes,uint256 bidAmt,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,uint256 slashAmount,bytes32 bidHash,string signature,string sharedSecretKey)"))
)

type Store interface {
	GetAESKey(common.Address) ([]byte, error)
	GetNikePrivateKey() (*ecdh.PrivateKey, error)
}

type encryptor struct {
	keySigner                  keysigner.KeySigner
	address                    []byte           // set for the provider
	nikePrvKey                 *ecdh.PrivateKey // set for the provider
	aesKey                     []byte           // set for the bidder
	store                      Store
	domainSeparatorBidHash     common.Hash // Precomputed domain separator for bids
	domainSeparatorPreConfHash common.Hash // Precomputed domain separator for pre-confirmations
}

func NewEncryptor(ks keysigner.KeySigner, store Store, chainID *big.Int, preconfContract string) (*encryptor, error) {
	address := ks.GetAddress()
	// those keys are set up during the libp2p.New initialization.
	aesKey, err := store.GetAESKey(address)
	if err != nil {
		return nil, err
	}
	nikePrvKey, err := store.GetNikePrivateKey()
	if err != nil {
		return nil, err
	}

	preconfContractAddr := common.HexToAddress(preconfContract)

	domainSeparatorBidHash, err := ComputeDomainSeparator("PreConfBid", chainID, preconfContractAddr)
	if err != nil {
		return nil, err
	}

	domainSeparatorPreConfHash, err := ComputeDomainSeparator("OpenedCommitment", chainID, preconfContractAddr)
	if err != nil {
		return nil, err
	}

	return &encryptor{
		keySigner:                  ks,
		address:                    address.Bytes(), // set for the provider
		nikePrvKey:                 nikePrvKey,
		aesKey:                     aesKey, // set for the bidder
		store:                      store,
		domainSeparatorBidHash:     domainSeparatorBidHash,
		domainSeparatorPreConfHash: domainSeparatorPreConfHash,
	}, nil
}

func (e *encryptor) ConstructEncryptedBid(
	bid *preconfpb.Bid,
) (*preconfpb.EncryptedBid, *ecdh.PrivateKey, error) {
	if bid.TxHash == "" || bid.BidAmount == "" || bid.BlockNumber == 0 {
		return nil, nil, ErrMissingRequiredFields
	}

	bidHash, err := GetBidHash(bid, e.domainSeparatorBidHash)
	if err != nil {
		return nil, nil, err
	}

	sig, err := e.keySigner.SignHash(bidHash)
	if err != nil {
		return nil, nil, err
	}

	nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	transformSignatureVValue(sig)

	bid.NikePublicKey = nikePrivateKey.PublicKey().Bytes()
	bid.Digest = bidHash
	bid.Signature = sig

	bidDataBytes, err := proto.Marshal(bid)
	if err != nil {
		return nil, nil, err
	}

	encryptedBidData, err := p2pcrypto.EncryptWithAESGCM(e.aesKey, bidDataBytes)
	if err != nil {
		return nil, nil, err
	}

	return &preconfpb.EncryptedBid{Ciphertext: encryptedBidData}, nikePrivateKey, nil
}

func (e *encryptor) ConstructEncryptedPreConfirmation(
	bid *preconfpb.Bid,
) (*preconfpb.PreConfirmation, *preconfpb.EncryptedPreConfirmation, error) {
	bidDataPublicKey, err := ecdh.Curve.NewPublicKey(ecdh.P256(), bid.NikePublicKey)
	if err != nil {
		return nil, nil, err
	}

	sharedSecretProviderSk, err := e.nikePrvKey.ECDH(bidDataPublicKey)
	if err != nil {
		return nil, nil, err
	}

	preConfirmation := &preconfpb.PreConfirmation{
		Bid:             bid,
		SharedSecret:    sharedSecretProviderSk,
		ProviderAddress: e.address,
	}

	preConfirmationHash, err := GetPreConfirmationHash(preConfirmation, e.domainSeparatorPreConfHash)
	if err != nil {
		return nil, nil, err
	}

	sig, err := e.keySigner.SignHash(preConfirmationHash)
	if err != nil {
		return nil, nil, err
	}

	transformSignatureVValue(sig)

	preConfirmation.Digest = preConfirmationHash
	preConfirmation.Signature = sig

	return preConfirmation, &preconfpb.EncryptedPreConfirmation{
		Commitment: preConfirmationHash,
		Signature:  sig,
	}, nil
}

func (e *encryptor) VerifyBid(bid *preconfpb.Bid) (*common.Address, error) {
	if bid.Digest == nil || bid.Signature == nil {
		return nil, ErrMissingHashSignature
	}

	bidHash, err := GetBidHash(bid, e.domainSeparatorBidHash)
	if err != nil {
		return nil, err
	}

	return eipVerify(
		bidHash,
		bid.Digest,
		bid.Signature,
	)
}

func (e *encryptor) DecryptBidData(bidderAddress common.Address, bid *preconfpb.EncryptedBid) (*preconfpb.Bid, error) {
	aesKey, err := e.store.GetAESKey(bidderAddress)
	if err != nil {
		return nil, err
	}
	if aesKey == nil {
		return nil, ErrNoAesKeyFound
	}
	decryptedBytes, err := p2pcrypto.DecryptWithAESGCM(aesKey, bid.Ciphertext)
	if err != nil {
		return nil, err
	}

	var bidData preconfpb.Bid
	if err := proto.Unmarshal(decryptedBytes, &bidData); err != nil {
		return nil, err
	}

	return &bidData, nil
}

// VerifyPreConfirmation verifies the preconfirmation message, and returns the address of the provider
// that signed the preconfirmation.
func (e *encryptor) VerifyEncryptedPreConfirmation(
	bid *preconfpb.Bid,
	providerNikePK *ecdh.PublicKey,
	bidderNikeSC *ecdh.PrivateKey,
	c *preconfpb.EncryptedPreConfirmation,
) ([]byte, *common.Address, error) {
	if c.Signature == nil {
		return nil, nil, ErrMissingHashSignature
	}

	sharedSecredBidderSk, err := bidderNikeSC.ECDH(providerNikePK)
	if err != nil {
		return nil, nil, err
	}

	preConfirmation := &preconfpb.PreConfirmation{
		Bid:          bid,
		SharedSecret: sharedSecredBidderSk,
	}

	preConfirmationHash, err := GetPreConfirmationHash(preConfirmation, e.domainSeparatorPreConfHash)
	if err != nil {
		return nil, nil, err
	}

	address, err := eipVerify(preConfirmationHash, c.Commitment, c.Signature)
	if err != nil {
		return nil, nil, err
	}

	return sharedSecredBidderSk, address, nil
}

func eipVerify(
	payloadHash []byte,
	expectedhash []byte,
	signature []byte,
) (*common.Address, error) {
	if !bytes.Equal(payloadHash, expectedhash) {
		return nil, ErrInvalidHash
	}

	sig := make([]byte, len(signature))
	copy(sig, signature)
	if sig[64] >= 27 && sig[64] <= 28 {
		sig[64] -= 27
	}

	pubkey, err := crypto.SigToPub(payloadHash, sig)
	if err != nil {
		return nil, err
	}

	if !crypto.VerifySignature(
		crypto.FromECDSAPub(pubkey),
		payloadHash,
		sig[:len(sig)-1],
	) {
		return nil, ErrInvalidSignature
	}

	c := crypto.PubkeyToAddress(*pubkey)

	return &c, err
}

// GetBidHash returns the hash of the bid message. This is done manually to match the
// Solidity implementation. If the types change, this will need to be updated.
func GetBidHash(bid *preconfpb.Bid, domainSeparatorHash common.Hash) ([]byte, error) {
	// Compute the struct hash
	structHash, err := computeBidStructHash(bid)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid struct hash %w", err)
	}

	// Final EIP-712 hash
	eip712Hash := computeEIP712Hash(domainSeparatorHash, structHash)
	return eip712Hash, nil
}

func computeBidStructHash(bid *preconfpb.Bid) (common.Hash, error) {
	bidAmt, ok := big.NewInt(0).SetString(bid.BidAmount, 10)
	if !ok {
		return common.Hash{}, ErrInvalidBidAmt
	}

	slashAmt, ok := big.NewInt(0).SetString(bid.SlashAmount, 10)
	if !ok {
		return common.Hash{}, ErrInvalidBidAmt
	}

	txnHashHash := crypto.Keccak256Hash([]byte(bid.TxHash))
	revertingTxHashesHash := crypto.Keccak256Hash([]byte(bid.RevertingTxHashes))

	bidStructType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "EIP712BidTypeHash", Type: "bytes32"},
		{Name: "txnHash", Type: "bytes32"},
		{Name: "revertingTxHashes", Type: "bytes32"},
		{Name: "bidAmt", Type: "uint256"},
		{Name: "blockNumber", Type: "uint64"},
		{Name: "decayStartTimestamp", Type: "uint64"},
		{Name: "decayEndTimestamp", Type: "uint64"},
		{Name: "slashAmount", Type: "uint256"},
	})
	if err != nil {
		return common.Hash{}, err
	}

	// Create the arguments array
	bidStructArguments := abi.Arguments{
		{Name: "EIP712BidTypeHash", Type: *bidStructType.TupleElems[0]},
		{Name: "txnHash", Type: *bidStructType.TupleElems[1]},
		{Name: "revertingTxHashes", Type: *bidStructType.TupleElems[2]},
		{Name: "bidAmt", Type: *bidStructType.TupleElems[3]},
		{Name: "blockNumber", Type: *bidStructType.TupleElems[4]},
		{Name: "decayStartTimestamp", Type: *bidStructType.TupleElems[5]},
		{Name: "decayEndTimestamp", Type: *bidStructType.TupleElems[6]},
		{Name: "slashAmount", Type: *bidStructType.TupleElems[7]},
	}

	// Encode the bid struct using ABI encoding
	encodedBid, err := bidStructArguments.Pack(
		EIP712BidTypeHash,
		txnHashHash,
		revertingTxHashesHash,
		bidAmt,
		uint64(bid.BlockNumber),
		uint64(bid.DecayStartTimestamp),
		uint64(bid.DecayEndTimestamp),
		slashAmt,
	)
	if err != nil {
		return common.Hash{}, err
	}

	// Hash the encoded bid struct
	return crypto.Keccak256Hash(encodedBid), nil
}

// GetPreConfirmationHash returns the hash of the preconfirmation message. This is done manually to match the
// Solidity implementation. If the types change, this will need to be updated.
func GetPreConfirmationHash(c *preconfpb.PreConfirmation, domainSeparatorHash common.Hash) ([]byte, error) {
	// Compute the struct hash
	structHash, err := computePreConfStructHash(c)
	if err != nil {
		return nil, err
	}

	// Final EIP-712 hash
	eip712Hash := computeEIP712Hash(domainSeparatorHash, structHash)
	return eip712Hash, nil
}

func computePreConfStructHash(c *preconfpb.PreConfirmation) (common.Hash, error) {
	bidAmt, ok := big.NewInt(0).SetString(c.Bid.BidAmount, 10)
	if !ok {
		return common.Hash{}, ErrInvalidBidAmt
	}

	slashAmt, ok := big.NewInt(0).SetString(c.Bid.SlashAmount, 10)
	if !ok {
		return common.Hash{}, ErrInvalidBidAmt
	}

	txnHashHash := crypto.Keccak256Hash([]byte(c.Bid.TxHash))
	revertingTxHashesHash := crypto.Keccak256Hash([]byte(c.Bid.RevertingTxHashes))
	bidDigestHash, err := toBytes32(c.Bid.Digest)
	if err != nil {
		return common.Hash{}, err
	}
	signatureHash := crypto.Keccak256Hash(c.Bid.Signature)
	sharedSecretHash := crypto.Keccak256Hash(c.SharedSecret)

	preConfStructType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "EIP712CommitmentTypeHash", Type: "bytes32"},
		{Name: "txnHash", Type: "bytes32"},
		{Name: "revertingTxHashes", Type: "bytes32"},
		{Name: "bidAmt", Type: "uint256"},
		{Name: "blockNumber", Type: "uint64"},
		{Name: "decayStartTimestamp", Type: "uint64"},
		{Name: "decayEndTimestamp", Type: "uint64"},
		{Name: "slashAmount", Type: "uint256"},
		{Name: "bidDigest", Type: "bytes32"},
		{Name: "signature", Type: "bytes32"},
		{Name: "sharedSecret", Type: "bytes32"},
	})
	if err != nil {
		return common.Hash{}, err
	}

	// Create the arguments array
	preConfStructArguments := abi.Arguments{
		{Name: "EIP712CommitmentTypeHash", Type: *preConfStructType.TupleElems[0]},
		{Name: "txnHash", Type: *preConfStructType.TupleElems[1]},
		{Name: "revertingTxHashes", Type: *preConfStructType.TupleElems[2]},
		{Name: "bidAmt", Type: *preConfStructType.TupleElems[3]},
		{Name: "blockNumber", Type: *preConfStructType.TupleElems[4]},
		{Name: "decayStartTimestamp", Type: *preConfStructType.TupleElems[5]},
		{Name: "decayEndTimestamp", Type: *preConfStructType.TupleElems[6]},
		{Name: "slashAmount", Type: *preConfStructType.TupleElems[7]},
		{Name: "bidDigest", Type: *preConfStructType.TupleElems[8]},
		{Name: "signature", Type: *preConfStructType.TupleElems[9]},
		{Name: "sharedSecret", Type: *preConfStructType.TupleElems[10]},
	}

	// Encode the pre-confirmation struct using ABI encoding
	encodedPreConf, err := preConfStructArguments.Pack(
		EIP712CommitmentTypeHash,
		txnHashHash,
		revertingTxHashesHash,
		bidAmt,
		uint64(c.Bid.BlockNumber),
		uint64(c.Bid.DecayStartTimestamp),
		uint64(c.Bid.DecayEndTimestamp),
		slashAmt,
		bidDigestHash,
		signatureHash,
		sharedSecretHash,
	)
	if err != nil {
		return common.Hash{}, err
	}

	// Hash the encoded pre-confirmation struct
	return crypto.Keccak256Hash(encodedPreConf), nil
}

func ComputeDomainSeparator(name string, chainId *big.Int, verifyingContract common.Address) (common.Hash, error) {
	domainTypeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))
	nameHash := name
	versionHash := "1"

	domainSeparatorType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "domainTypeHash", Type: "bytes32"},
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	})
	if err != nil {
		return common.Hash{}, err
	}

	// Create the arguments array, dereferencing the pointers
	domainSeparatorArguments := abi.Arguments{
		{Name: "domainTypeHash", Type: *domainSeparatorType.TupleElems[0]},
		{Name: "name", Type: *domainSeparatorType.TupleElems[1]},
		{Name: "version", Type: *domainSeparatorType.TupleElems[2]},
		{Name: "chainId", Type: *domainSeparatorType.TupleElems[3]},
		{Name: "verifyingContract", Type: *domainSeparatorType.TupleElems[4]},
	}

	// Encode the domain separator using ABI encoding
	encodedDomain, err := domainSeparatorArguments.Pack(
		domainTypeHash,
		nameHash,
		versionHash,
		chainId,
		verifyingContract,
	)
	if err != nil {
		return common.Hash{}, err
	}

	// Hash the encoded domain separator
	return crypto.Keccak256Hash(encodedDomain), nil
}

func computeEIP712Hash(domainSeparatorHash, structHash common.Hash) []byte {
	// EIP-712 hash format: "\x19\x01" || domainSeparator || structHash
	eip712Data := append([]byte("\x19\x01"), append(domainSeparatorHash.Bytes(), structHash.Bytes()...)...)
	return crypto.Keccak256Hash(eip712Data).Bytes()
}

func transformSignatureVValue(sig []byte) {
	if sig[64] == 0 || sig[64] == 1 {
		sig[64] += 27 // Transform V from 0/1 to 27/28
	}
}

func toBytes32(slice []byte) ([32]byte, error) {
	var array [32]byte
	if len(slice) != 32 {
		return array, fmt.Errorf("invalid length: expected 32 bytes, got %d", len(slice))
	}
	copy(array[:], slice)
	return array, nil
}
