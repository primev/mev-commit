package keysstore_test

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/keysstore"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	"github.com/stretchr/testify/assert"
)

func TestAESKey(t *testing.T) {
	st := inmem.New()
	store := keysstore.New(st)
	bidder := common.HexToAddress("0x1")

	// Set and get AES key
	expectedKey := []byte("aes-key")
	err := store.SetAESKey(bidder, expectedKey)
	assert.NoError(t, err)

	retrievedKey, err := store.GetAESKey(bidder)
	assert.NoError(t, err)
	assert.Equal(t, expectedKey, retrievedKey)

	// Get non-existent AES key
	nonExistentBidder := common.HexToAddress("0x2")
	retrievedKey, err = store.GetAESKey(nonExistentBidder)
	assert.NoError(t, err)
	assert.Nil(t, retrievedKey)
}

func TestECIESPrivateKey(t *testing.T) {
	st := inmem.New()
	store := keysstore.New(st)

	// Get non-existent ECIES private key
	retrievedKey, err := store.GetECIESPrivateKey()
	assert.NoError(t, err)
	assert.Nil(t, retrievedKey)

	// Generate ECIES private key
	privateKeyECIES, err := ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
	assert.NoError(t, err)

	// Set and get ECIES private key
	err = store.SetECIESPrivateKey(privateKeyECIES)
	assert.NoError(t, err)

	retrievedKey, err = store.GetECIESPrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, privateKeyECIES.D, retrievedKey.D)
}

func TestNikePrivateKey(t *testing.T) {
	st := inmem.New()
	store := keysstore.New(st)

	// Get non-existent Nike private key
	retrievedKey, err := store.GetBN254PrivateKey()
	assert.NoError(t, err)
	assert.Nil(t, retrievedKey)

	// Generate Nike private key
	sk, _, err := p2pcrypto.GenerateKeyPairBN254()
	assert.NoError(t, err)

	// Set and get Nike private key
	err = store.SetBN254PrivateKey(sk)
	assert.NoError(t, err)

	retrievedKey, err = store.GetBN254PrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, sk.Bytes(), retrievedKey.Bytes())
}

func TestNikePublicKey(t *testing.T) {
	st := inmem.New()
	store := keysstore.New(st)

	// Get non-existent Nike public key
	retrievedKey, err := store.GetBN254PublicKey()
	assert.NoError(t, err)
	assert.Nil(t, retrievedKey)

	// Generate Nike key pair
	_, pk, err := p2pcrypto.GenerateKeyPairBN254()
	assert.NoError(t, err)
	
	// Set and get Nike public key
	err = store.SetBN254PublicKey(pk)
	assert.NoError(t, err)

	retrievedKey, err = store.GetBN254PublicKey()
	assert.NoError(t, err)
	assert.Equal(t, pk.Bytes(), retrievedKey.Bytes())
}
