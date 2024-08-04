package keysstore_test

import (
	"crypto/ecdh"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
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
	privateKeyECIES, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
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
	retrievedKey, err := store.GetNikePrivateKey()
	assert.NoError(t, err)
	assert.Nil(t, retrievedKey)

	// Generate Nike private key
	privateKeyNike, err := ecdh.P256().GenerateKey(rand.Reader)
	assert.NoError(t, err)

	// Set and get Nike private key
	err = store.SetNikePrivateKey(privateKeyNike)
	assert.NoError(t, err)

	retrievedKey, err = store.GetNikePrivateKey()
	assert.NoError(t, err)
	assert.Equal(t, privateKeyNike.Bytes(), retrievedKey.Bytes())
}
