package preconfencryptor_test

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	mockkeysigner "github.com/primev/mev-commit/x/keysigner/mock"
	"github.com/primev/mev-commit/p2p/pkg/signer/preconfencryptor"
	"github.com/primev/mev-commit/p2p/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestBids(t *testing.T) {
	t.Parallel()

	t.Run("bid", func(t *testing.T) {
		key, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		address := crypto.PubkeyToAddress(key.PublicKey)
		keySigner := mockkeysigner.NewMockKeySigner(key, address)
		aesKey, err := p2pcrypto.GenerateAESKey()
		if err != nil {
			t.Fatal(err)
		}
		bidderStore := store.NewStore()
		err = bidderStore.SetAESKey(address, aesKey)
		if err != nil {
			t.Fatal(err)
		}
		encryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore)
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now().UnixMilli()
		end := start + 100000
		_, encryptedBid, _, err := encryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end)
		if err != nil {
			t.Fatal(err)
		}

		providerStore := store.NewStore()
		err = providerStore.SetAESKey(address, aesKey)
		if err != nil {
			t.Fatal(err)
		}
		encryptorProvider, err := preconfencryptor.NewEncryptor(keySigner, providerStore)
		if err != nil {
			t.Fatal(err)
		}
		bid, err := encryptorProvider.DecryptBidData(address, encryptedBid)
		if err != nil {
			t.Fatal(err)
		}

		bidAddress, err := encryptor.VerifyBid(bid)
		if err != nil {
			t.Fatal(err)
		}

		originatorAddress, pubkey, err := encryptor.BidOriginator(bid)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, address, *originatorAddress)
		assert.Equal(t, address, *bidAddress)
		assert.Equal(t, key.PublicKey, *pubkey)
	})
	t.Run("preConfirmation", func(t *testing.T) {
		bidderKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		aesKey, err := p2pcrypto.GenerateAESKey()
		if err != nil {
			t.Fatal(err)
		}

		keySigner := mockkeysigner.NewMockKeySigner(bidderKey, crypto.PubkeyToAddress(bidderKey.PublicKey))
		bidderStore := store.NewStore()
		err = bidderStore.SetAESKey(crypto.PubkeyToAddress(bidderKey.PublicKey), aesKey)
		if err != nil {
			t.Fatal(err)
		}
		bidderEncryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore)
		if err != nil {
			t.Fatal(err)
		}
		providerKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		bidderAddress := crypto.PubkeyToAddress(bidderKey.PublicKey)
		keySigner = mockkeysigner.NewMockKeySigner(providerKey, crypto.PubkeyToAddress(providerKey.PublicKey))
		providerStore := store.NewStore()
		err = providerStore.SetAESKey(crypto.PubkeyToAddress(bidderKey.PublicKey), aesKey)
		if err != nil {
			t.Fatal(err)
		}
		providerNikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		err = providerStore.SetNikePrivateKey(providerNikePrivateKey)
		if err != nil {
			t.Fatal(err)
		}

		providerEncryptor, err := preconfencryptor.NewEncryptor(keySigner, providerStore)
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now().UnixMilli()
		end := start + 100000

		bid, encryptedBid, nikePrivateKey, err := bidderEncryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end)
		if err != nil {
			t.Fatal(err)
		}

		decryptedBid, err := providerEncryptor.DecryptBidData(bidderAddress, encryptedBid)
		if err != nil {
			t.Fatal(err)
		}
		_, encryptedPreConfirmation, err := providerEncryptor.ConstructEncryptedPreConfirmation(decryptedBid)
		if err != nil {
			t.Fail()
		}

		_, address, err := bidderEncryptor.VerifyEncryptedPreConfirmation(providerNikePrivateKey.PublicKey(), nikePrivateKey, bid.Digest, encryptedPreConfirmation)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, crypto.PubkeyToAddress(providerKey.PublicKey), *address)
	})
}

func TestHashing(t *testing.T) {
	t.Parallel()

	t.Run("bid", func(t *testing.T) {
		bid := &preconfpb.Bid{
			TxHash:              "0xkartik",
			BidAmount:           "2",
			BlockNumber:         2,
			DecayStartTimestamp: 10,
			DecayEndTimestamp:   20,
		}

		hash, err := preconfencryptor.GetBidHash(bid)
		if err != nil {
			t.Fatal(err)
		}

		hashStr := hex.EncodeToString(hash)
		// This hash is sourced from the solidity contract to ensure interoperability
		expHash := "a0327970258c49b922969af74d60299a648c50f69a2d98d6ab43f32f64ac2100"
		if hashStr != expHash {
			t.Fatalf("hash mismatch: %s != %s", hashStr, expHash)
		}
	})

	t.Run("preConfirmation", func(t *testing.T) {
		bidHash := "a0327970258c49b922969af74d60299a648c50f69a2d98d6ab43f32f64ac2100"
		bidSignature := "876c1216c232828be9fabb14981c8788cebdf6ed66e563c4a2ccc82a577d052543207aeeb158a32d8977736797ae250c63ef69a82cd85b727da21e20d030fb311b"

		bidHashBytes, err := hex.DecodeString(bidHash)
		if err != nil {
			t.Fatal(err)
		}
		bidSigBytes, err := hex.DecodeString(bidSignature)
		if err != nil {
			t.Fatal(err)
		}

		bid := &preconfpb.Bid{
			TxHash:              "0xkartik",
			BidAmount:           "2",
			BlockNumber:         2,
			DecayStartTimestamp: 10,
			DecayEndTimestamp:   20,
			Digest:              bidHashBytes,
			Signature:           bidSigBytes,
		}

		sharedSecretBytes := []byte("0xsecret")

		preConfirmation := &preconfpb.PreConfirmation{
			Bid:          bid,
			SharedSecret: sharedSecretBytes,
		}

		hash, err := preconfencryptor.GetPreConfirmationHash(preConfirmation)
		if err != nil {
			t.Fatal(err)
		}

		hashStr := hex.EncodeToString(hash)
		expHash := "65618f8f9e46b8f0790c621ca2989cfe4c949594a4a3a81261baa682e8883840"
		if hashStr != expHash {
			t.Fatalf("hash mismatch: %s != %s", hashStr, expHash)
		}
	})
}

func TestVerify(t *testing.T) {
	t.Parallel()

	bidSig := "8af22e36247e14ba05d3a5a3cc62eee708cfd9ce293c0aebcbe7f89229f6db56638af8427806247d9abb295f681c1a2f2bb127f3bf80799f80d62b252cce04d91c"
	bidHash := "2574b1ab8a90e173528ddee748be8e8e696b1f0cf687f75966550f5e9ef408b0"

	bidHashBytes, err := hex.DecodeString(bidHash)
	if err != nil {
		t.Fatal(err)
	}

	bidSigBytes, err := hex.DecodeString(bidSig)
	if err != nil {
		t.Fatal(err)
	}

	// Adjust the last byte if it's 27 or 28
	if bidSigBytes[64] >= 27 && bidSigBytes[64] <= 28 {
		bidSigBytes[64] -= 27
	}

	owner, err := preconfencryptor.EIPVerify(bidHashBytes, bidHashBytes, bidSigBytes)
	if err != nil {
		t.Fatal(err)
	}

	expOwner := "0x8339F9E3d7B2693aD8955Aa5EC59D56669A84d60"
	if owner.Hex() != expOwner {
		t.Fatalf("owner mismatch: %s != %s", owner.Hex(), expOwner)
	}
}
