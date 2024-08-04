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
	"github.com/primev/mev-commit/p2p/pkg/keysstore"
	"github.com/primev/mev-commit/p2p/pkg/signer/preconfencryptor"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
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
		bidderStore := keysstore.New(inmemstorage.New())
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

		providerStore := keysstore.New(inmemstorage.New())
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
		bidderStore := keysstore.New(inmemstorage.New())
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
		providerStore := keysstore.New(inmemstorage.New())
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
			RevertingTxHashes:   "0xkartik",
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
		expHash := "9890bcda118cfabed02ff3b9d05a54dca5310e9ace3b05f259f4731f58ad0900"
		if hashStr != expHash {
			t.Fatalf("hash mismatch: %s != %s", hashStr, expHash)
		}

		alicePrivateKey, err := crypto.HexToECDSA("9C0257114EB9399A2985F8E75DAD7600C5D89FE3824FFA99EC1C3EB8BF3B0501")
		if err != nil {
			t.Fatal(err)
		}

		expHashBytes, err := hex.DecodeString(expHash)
		if err != nil {
			t.Fatal(err)
		}

		signature, err := crypto.Sign(expHashBytes, alicePrivateKey)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Signature: %s", hex.EncodeToString(signature))

		// Log the keccak of the signature
		signatureHash := crypto.Keccak256Hash(signature)
		t.Logf("Keccak256 of Signature: %s", signatureHash.Hex())
	})

	t.Run("preConfirmation", func(t *testing.T) {
		bidHash := "9890bcda118cfabed02ff3b9d05a54dca5310e9ace3b05f259f4731f58ad0900"
		bidSignature := "f9b66c6d57dac947a3aa2b37010df745592cf57f907d437767bc0af6d44b3dc1112168e4cab311d6dfddf7f58c0d07bb95403fca2cc48d4450e088cf9ee894c81b"

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
			RevertingTxHashes:   "0xkartik",
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
		expHash := "8257770d4be5c4b622e6bd6b45ff8deb6602235f3aa844b774eb21800eb4923a"
		if hashStr != expHash {
			t.Fatalf("hash mismatch: %s != %s", hashStr, expHash)
		}

		// Sign the hash with Bob's private key
		bobPrivateKey, err := crypto.HexToECDSA("38E47A7B719DCE63662AEAF43440326F551B8A7EE198CEE35CB5D517F2D296A2")
		if err != nil {
			t.Fatal(err)
		}

		signature, err := crypto.Sign(hash, bobPrivateKey)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Bob's Signature: %s", hex.EncodeToString(signature))
	})
}

func TestVerify(t *testing.T) {
	t.Parallel()

	bidSig := "f9b66c6d57dac947a3aa2b37010df745592cf57f907d437767bc0af6d44b3dc1112168e4cab311d6dfddf7f58c0d07bb95403fca2cc48d4450e088cf9ee894c800"
	bidHash := "9890bcda118cfabed02ff3b9d05a54dca5310e9ace3b05f259f4731f58ad0900"

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

	expOwner := "0x328809Bc894f92807417D2dAD6b7C998c1aFdac6"
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
	bidderStore := keysstore.New(inmemstorage.New())
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
	bidderStore := keysstore.New(inmemstorage.New())
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
	providerStore := keysstore.New(inmemstorage.New())
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
