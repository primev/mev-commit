package preconfencryptor

import (
	"bytes"
	"crypto/ecdh"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/keykeeper"
)

var (
	ErrAlreadySignedBid             = errors.New("already contains hash or signature")
	ErrMissingHashSignature         = errors.New("missing hash or signature")
	ErrInvalidSignature             = errors.New("signature is not valid")
	ErrInvalidHash                  = errors.New("bidhash doesn't match bid payload")
	ErrAlreadySignedPreConfirmation = errors.New("preConfirmation is already hashed or signed")
	ErrInvalidCommitment            = errors.New("commitment is incorrect")
)

type Encryptor interface {
	ConstructEncryptedBid(string, string, int64, int64, int64) (*preconfpb.Bid, *preconfpb.EncryptedBid, error)
	ConstructEncryptedPreConfirmation(*preconfpb.Bid) (*preconfpb.PreConfirmation, *preconfpb.EncryptedPreConfirmation, error)
	VerifyBid(*preconfpb.Bid) (*common.Address, error)
	VerifyEncryptedPreConfirmation(providerNikePK *ecdh.PublicKey, bidHash []byte, c *preconfpb.EncryptedPreConfirmation) ([]byte, *common.Address, error)
	DecryptBidData(common.Address, *preconfpb.EncryptedBid) (*preconfpb.Bid, error)
}

type encryptor struct {
	keyKeeper      keykeeper.KeyKeeper
	bidHashesToBid map[string]*preconfpb.Bid
}

func NewEncryptor(keyKeeper keykeeper.KeyKeeper) *encryptor {
	bidHashesToBid := make(map[string]*preconfpb.Bid)
	return &encryptor{
		keyKeeper:      keyKeeper,
		bidHashesToBid: bidHashesToBid,
	}
}

func (e *encryptor) ConstructEncryptedBid(
	txHash string,
	bidAmt string,
	blockNumber int64,
	decayStartTimeStamp int64,
	decayEndTimeStamp int64,
) (*preconfpb.Bid, *preconfpb.EncryptedBid, error) {
	if txHash == "" || bidAmt == "" || blockNumber == 0 {
		return nil, nil, errors.New("missing required fields")
	}

	bid := &preconfpb.Bid{
		BidAmount:           bidAmt,
		TxHash:              txHash,
		BlockNumber:         blockNumber,
		DecayStartTimestamp: decayStartTimeStamp,
		DecayEndTimestamp:   decayEndTimeStamp,
	}

	bidHash, err := GetBidHash(bid)
	if err != nil {
		return nil, nil, err
	}

	// todo: probably sign all data including nike public key
	sig, err := e.keyKeeper.SignHash(bidHash)
	if err != nil {
		return nil, nil, err
	}

	if sig[64] == 0 || sig[64] == 1 {
		sig[64] += 27 // Transform V from 0/1 to 27/28
	}

	bidderKK := e.keyKeeper.(*keykeeper.BidderKeyKeeper)
	nikePublicKey, err := bidderKK.GenerateNIKEKeys(bidHash)
	if err != nil {
		return nil, nil, err
	}

	bid.NikePublicKey = nikePublicKey.Bytes()
	bid.Digest = bidHash
	bid.Signature = sig

	bidDataBytes, err := json.Marshal(bid)
	if err != nil {
		return nil, nil, err
	}

	e.bidHashesToBid[hex.EncodeToString(bidHash)] = bid

	encryptedBidData, err := keykeeper.EncryptWithAESGCM(bidderKK.AESKey, bidDataBytes)
	if err != nil {
		return nil, nil, err
	}

	return bid, &preconfpb.EncryptedBid{Ciphertext: encryptedBidData}, nil
}

func (e *encryptor) ConstructEncryptedPreConfirmation(bid *preconfpb.Bid) (*preconfpb.PreConfirmation, *preconfpb.EncryptedPreConfirmation, error) {
	_, err := e.VerifyBid(bid)
	if err != nil {
		return nil, nil, err
	}

	bidDataPublicKey, err := ecdh.Curve.NewPublicKey(ecdh.P256(), bid.NikePublicKey)
	if err != nil {
		return nil, nil, err
	}

	providerKK := e.keyKeeper.(*keykeeper.ProviderKeyKeeper)
	sharedSecredProviderSk, err := providerKK.GetNIKEPrivateKey().ECDH(bidDataPublicKey)
	if err != nil {
		return nil, nil, err
	}

	preConfirmation := &preconfpb.PreConfirmation{
		Bid:             bid,
		SharedSecret:    sharedSecredProviderSk,
		ProviderAddress: providerKK.GetAddress().Bytes(),
	}

	preConfirmationHash, err := GetPreConfirmationHash(preConfirmation)
	if err != nil {
		return nil, nil, err
	}

	sig, err := e.keyKeeper.SignHash(preConfirmationHash)
	if err != nil {
		return nil, nil, err
	}

	if sig[64] == 0 || sig[64] == 1 {
		sig[64] += 27 // Transform V from 0/1 to 27/28
	}

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
	pkk := e.keyKeeper.(*keykeeper.ProviderKeyKeeper)
	aesKey, exists := pkk.GetAESKey(bidderAddress)
	if !exists {
		return nil, errors.New("no AES key found for bidder")
	}
	decryptedBytes, err := keykeeper.DecryptWithAESGCM(aesKey, bid.Ciphertext)
	if err != nil {
		return nil, err
	}

	var bidData preconfpb.Bid
	if err := json.Unmarshal(decryptedBytes, &bidData); err != nil {
		return nil, err
	}

	return &bidData, nil
}

// VerifyPreConfirmation verifies the preconfirmation message, and returns the address of the provider
// that signed the preconfirmation.
func (e *encryptor) VerifyEncryptedPreConfirmation(providerNikePK *ecdh.PublicKey, bidHash []byte, c *preconfpb.EncryptedPreConfirmation) ([]byte, *common.Address, error) {
	if c.Signature == nil {
		return nil, nil, ErrMissingHashSignature
	}

	bidHashStr := hex.EncodeToString(bidHash)
	bid := e.bidHashesToBid[bidHashStr]

	bidderKK := e.keyKeeper.(*keykeeper.BidderKeyKeeper)
	sharedSecredBidderSk, err := bidderKK.BidHashesToNIKE[bidHashStr].ECDH(providerNikePK)
	if err != nil {
		return nil, nil, err
	}

	preConfirmation := &preconfpb.PreConfirmation{
		Bid:          bid,
		Digest:       bidHash,
		Signature:    c.Signature,
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
		return nil, errors.New("invalid bid amount")
	}

	// EIP712_MESSAGE_TYPEHASH
	eip712MessageTypeHash := crypto.Keccak256Hash(
		[]byte("PreConfBid(string txnHash,uint64 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp)"),
	)

	// Convert the txnHash to a byte array and hash it
	txnHashHash := crypto.Keccak256Hash([]byte(bid.TxHash))

	// Encode values similar to Solidity's abi.encode
	// The reason we use math.U256Bytes is because we want to encode the uint64 as a 32 byte array
	// The EVM does this for values due via padding to 32 bytes, as that's the base size of a word in the EVM
	data := append(eip712MessageTypeHash.Bytes(), txnHashHash.Bytes()...)
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
		nameHash           = crypto.Keccak256Hash([]byte("PreConfCommitment"))
		versionHash        = crypto.Keccak256Hash([]byte("1"))
		domainSeparatorBid = crypto.Keccak256Hash(
			append(append(domainTypeHash.Bytes(), nameHash.Bytes()...), versionHash.Bytes()...),
		)
	)

	bidAmt, ok := big.NewInt(0).SetString(c.Bid.BidAmount, 10)
	if !ok {
		return nil, errors.New("invalid bid amount")
	}

	// EIP712_MESSAGE_TYPEHASH
	eip712MessageTypeHash := crypto.Keccak256Hash(
		[]byte("PreConfCommitment(string txnHash,uint64 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,bytes32 bidHash,string signature,string sharedSecretKey)"),
	)

	// Convert the txnHash to a byte array and hash it
	txnHashHash := crypto.Keccak256Hash([]byte(c.Bid.TxHash))
	bidDigestHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.Bid.Digest)))
	bidSigHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.Bid.Signature)))
	sharedSecretHash := crypto.Keccak256Hash([]byte(hex.EncodeToString(c.SharedSecret)))

	// Encode values similar to Solidity's abi.encode
	data := append(eip712MessageTypeHash.Bytes(), txnHashHash.Bytes()...)
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
