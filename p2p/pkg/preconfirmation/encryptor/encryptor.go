package preconfencryptor

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

// EIP712 message type hashes
var bidMessageTypeHash = crypto.Keccak256Hash(
	[]byte("PreConfBid(" +
		"string txnHash," +
		"string revertingTxHashes," +
		"uint256 bid," +
		"uint64 blockNumber," +
		"uint64 decayStartTimeStamp," +
		"uint64 decayEndTimeStamp)"),
)

var commitmentTypeHash = crypto.Keccak256Hash(
	[]byte("OpenedCommitment(" +
		"string txnHash," +
		"string revertingTxHashes," +
		"uint256 bid," +
		"uint64 blockNumber," +
		"uint64 decayStartTimeStamp," +
		"uint64 decayEndTimeStamp," +
		"bytes32 bidHash," +
		"string signature," +
		"string sharedSecretKey)"),
)

type Store interface {
	GetAESKey(common.Address) ([]byte, error)
	GetNikePrivateKey() (*ecdh.PrivateKey, error)
}

type encryptor struct {
	keySigner  keysigner.KeySigner
	address    []byte           // set for the provider
	nikePrvKey *ecdh.PrivateKey // set for the provider
	aesKey     []byte           // set for the bidder
	store      Store
}

func NewEncryptor(ks keysigner.KeySigner, store Store) (*encryptor, error) {
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

	return &encryptor{
		keySigner:  ks,
		address:    address.Bytes(), // set for the provider
		nikePrvKey: nikePrvKey,
		aesKey:     aesKey, // set for the bidder
		store:      store,
	}, nil
}

func (e *encryptor) ConstructEncryptedBid(
	bid *preconfpb.Bid,
) (*preconfpb.EncryptedBid, *ecdh.PrivateKey, error) {
	if bid.TxHash == "" || bid.BidAmount == "" || bid.BlockNumber == 0 {
		return nil, nil, ErrMissingRequiredFields
	}

	bidHash, err := GetBidHash(bid)
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

	preConfirmationHash, err := GetPreConfirmationHash(preConfirmation)
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

	bidHash, err := GetBidHash(bid)
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

	preConfirmationHash, err := GetPreConfirmationHash(preConfirmation)
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
func GetBidHash(bid *preconfpb.Bid) ([]byte, error) {
	// DOMAIN_SEPARATOR_BID
	var (
		domainTypeHash = crypto.Keccak256Hash(
			[]byte("EIP712Domain(string name,string version)"),
		)
		nameHash           = crypto.Keccak256Hash([]byte("PreConfBid"))
		versionHash        = crypto.Keccak256Hash([]byte("1"))
		domainSeparatorBid = crypto.Keccak256Hash(
			append(append(domainTypeHash.Bytes(), nameHash.Bytes()...), versionHash.Bytes()...),
		)
	)

	bidAmt, ok := big.NewInt(0).SetString(bid.BidAmount, 10)
	if !ok {
		return nil, ErrInvalidBidAmt
	}

	// Convert the txnHash to a byte array and hash it
	txnHashHash := crypto.Keccak256Hash([]byte(bid.TxHash))
	revertingTxHashesHash := crypto.Keccak256Hash([]byte(bid.RevertingTxHashes))

	// Encode values similar to Solidity's abi.encode
	// The reason we use math.U256Bytes is because we want to encode the uint64 as a 32 byte array
	// The EVM does this for values due via padding to 32 bytes, as that's the base size of a word in the EVM
	data := append(bidMessageTypeHash.Bytes(), txnHashHash.Bytes()...)
	data = append(data, revertingTxHashesHash.Bytes()...)
	data = append(data, math.U256Bytes(bidAmt)...)
	data = append(data, math.U256Bytes(big.NewInt(bid.BlockNumber))...)
	data = append(data, math.U256Bytes(big.NewInt(bid.DecayStartTimestamp))...)
	data = append(data, math.U256Bytes(big.NewInt(bid.DecayEndTimestamp))...)
	dataHash := crypto.Keccak256Hash(data)

	rawData := append([]byte("\x19\x01"), append(domainSeparatorBid.Bytes(), dataHash.Bytes()...)...)
	// Create the final hash
	return crypto.Keccak256Hash(rawData).Bytes(), nil
}

// GetPreConfirmationHash returns the hash of the preconfirmation message. This is done manually to match the
// Solidity implementation. If the types change, this will need to be updated.
func GetPreConfirmationHash(c *preconfpb.PreConfirmation) ([]byte, error) {
	// DOMAIN_SEPARATOR_BID
	var (
		domainTypeHash = crypto.Keccak256Hash(
			[]byte("EIP712Domain(string name,string version)"),
		)
		nameHash           = crypto.Keccak256Hash([]byte("OpenedCommitment"))
		versionHash        = crypto.Keccak256Hash([]byte("1"))
		domainSeparatorBid = crypto.Keccak256Hash(
			append(append(domainTypeHash.Bytes(), nameHash.Bytes()...), versionHash.Bytes()...),
		)
	)

	bidAmt, ok := big.NewInt(0).SetString(c.Bid.BidAmount, 10)
	if !ok {
		return nil, ErrInvalidBidAmt
	}

	// Convert the txnHash to a byte array and hash it
	txnHashHash := crypto.Keccak256Hash([]byte(c.Bid.TxHash))
	revertingTxHashesHash := crypto.Keccak256Hash([]byte(c.Bid.RevertingTxHashes))
	bidDigestHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.Bid.Digest)))
	bidSigHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.Bid.Signature)))
	sharedSecretHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.SharedSecret)))

	// Encode values similar to Solidity's abi.encode
	data := append(commitmentTypeHash.Bytes(), txnHashHash.Bytes()...)
	data = append(data, revertingTxHashesHash.Bytes()...)
	data = append(data, math.U256Bytes(bidAmt)...)
	data = append(data, math.U256Bytes(big.NewInt(c.Bid.BlockNumber))...)
	data = append(data, math.U256Bytes(big.NewInt(c.Bid.DecayStartTimestamp))...)
	data = append(data, math.U256Bytes(big.NewInt(c.Bid.DecayEndTimestamp))...)
	data = append(data, bidDigestHash.Bytes()...)
	data = append(data, bidSigHash.Bytes()...)
	data = append(data, sharedSecretHash.Bytes()...)
	dataHash := crypto.Keccak256Hash(data)

	rawData := append([]byte("\x19\x01"), append(domainSeparatorBid.Bytes(), dataHash.Bytes()...)...)
	// Create the final hash
	return crypto.Keccak256Hash(rawData).Bytes(), nil
}

func transformSignatureVValue(sig []byte) {
	if sig[64] == 0 || sig[64] == 1 {
		sig[64] += 27 // Transform V from 0/1 to 27/28
	}
}
