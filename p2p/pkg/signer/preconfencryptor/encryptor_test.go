package preconfencryptor_test

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/signer/preconfencryptor"
	"github.com/primev/mev-commit/p2p/pkg/store"
	mockkeysigner "github.com/primev/mev-commit/x/keysigner/mock"
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
		_, encryptedBid, _, err := encryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end, "")
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

		bid, encryptedBid, nikePrivateKey, err := bidderEncryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end, "")
		if err != nil {
			t.Fatal(err)
		}

		decryptedBid, err := providerEncryptor.DecryptBidData(bidderAddress, encryptedBid)
		if err != nil {
			t.Fatal(err)
		}
		_, encryptedPreConfirmation, err := providerEncryptor.ConstructEncryptedPreConfirmation(decryptedBid)
		if err != nil {
			t.Fatal(err)
		}

		_, address, err := bidderEncryptor.VerifyEncryptedPreConfirmation(providerNikePrivateKey.PublicKey(), nikePrivateKey, bid.Digest, encryptedPreConfirmation)
		if err != nil {
			t.Fatal(err)
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
		expHash := "56c06a13be335eba981b780ea45dff258a7c429d0e9d993235ef2d3a7e435df8"
		if hashStr != expHash {
			t.Fatalf("hash mismatch: %s != %s", hashStr, expHash)
		}
	})

	t.Run("preConfirmation", func(t *testing.T) {
		bidHash := "56c06a13be335eba981b780ea45dff258a7c429d0e9d993235ef2d3a7e435df8"
		bidSignature := "2e7df27808c72d7d5b2543bb63b06c0ae2144e021593b8d2a7cca6a3fb2d9c4b1a82dd2a07266de9364d255bdb709476ad96b826ec855efb528eaff66682997e1c"

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
		expHash := "9d954942ad3f6cb41ccd029869be7b28036270b4754665a3783c2d6bf0ef7d08"
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

type testBid struct {
	hash        string
	amount      string
	blocknumber int64
	start       int64
	end         int64
}

func generateRandomValues() (*testBid, error) {
	start := mrand.Int63()
	end := start + mrand.Int63n(100000)
	bidHashBytes := make([]byte, 32)
	_, err := rand.Read(bidHashBytes)
	if err != nil {
		return nil, err
	}
	bidHash := hex.EncodeToString(bidHashBytes)

	bidAmount := mrand.Int63n(1000)
	blocknumber := mrand.Int63n(100000)

	return &testBid{
		hash:        bidHash,
		amount:      fmt.Sprintf("%d", bidAmount),
		blocknumber: blocknumber,
		start:       start,
		end:         end,
	}, nil
}

func BenchmarkConstructEncryptedBid(b *testing.B) {
	key, err := crypto.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	keySigner := mockkeysigner.NewMockKeySigner(key, address)
	aesKey, err := p2pcrypto.GenerateAESKey()
	if err != nil {
		b.Fatal(err)
	}
	bidderStore := store.NewStore()
	err = bidderStore.SetAESKey(address, aesKey)
	if err != nil {
		b.Fatal(err)
	}
	encryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore)
	if err != nil {
		b.Fatal(err)
	}

	bids := make([]*testBid, b.N)
	for i := 0; i < len(bids); i++ {
		bids[i], err = generateRandomValues()
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	// Benchmark loop
	for i := 0; i < b.N; i++ {
		_, _, _, err := encryptor.ConstructEncryptedBid(bids[i].hash, bids[i].amount, bids[i].blocknumber, bids[i].start, bids[i].end, "")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConstructEncryptedPreConfirmation(b *testing.B) {
	// Setup code (initialize encryptor, bid, etc.)
	bidderKey, err := crypto.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	aesKey, err := p2pcrypto.GenerateAESKey()
	if err != nil {
		b.Fatal(err)
	}

	keySigner := mockkeysigner.NewMockKeySigner(bidderKey, crypto.PubkeyToAddress(bidderKey.PublicKey))
	bidderStore := store.NewStore()
	err = bidderStore.SetAESKey(crypto.PubkeyToAddress(bidderKey.PublicKey), aesKey)
	if err != nil {
		b.Fatal(err)
	}
	bidderEncryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore)
	if err != nil {
		b.Fatal(err)
	}
	providerKey, err := crypto.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}

	keySigner = mockkeysigner.NewMockKeySigner(providerKey, crypto.PubkeyToAddress(providerKey.PublicKey))
	providerStore := store.NewStore()
	err = providerStore.SetAESKey(crypto.PubkeyToAddress(bidderKey.PublicKey), aesKey)
	if err != nil {
		b.Fatal(err)
	}
	providerNikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	err = providerStore.SetNikePrivateKey(providerNikePrivateKey)
	if err != nil {
		b.Fatal(err)
	}
	providerEncryptor, err := preconfencryptor.NewEncryptor(keySigner, providerStore)
	if err != nil {
		b.Fatal(err)
	}

	var bid *testBid
	bids := make([]*preconfpb.Bid, b.N)
	for i := 0; i < len(bids); i++ {
		bid, err = generateRandomValues()
		if err != nil {
			b.Fatal(err)
		}
		bids[i], _, _, err = bidderEncryptor.ConstructEncryptedBid(bid.hash, bid.amount, bid.blocknumber, bid.start, bid.end, "")
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := providerEncryptor.ConstructEncryptedPreConfirmation(bids[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}
