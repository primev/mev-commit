package preconfencryptor_test

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/keysstore"
	preconfencryptor "github.com/primev/mev-commit/p2p/pkg/preconfirmation/encryptor"
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
		encryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now().UnixMilli()
		end := start + 100000
		reqBid := &preconfpb.Bid{
			TxHash:              "0xkartik",
			BidAmount:           "10",
			BlockNumber:         2,
			DecayStartTimestamp: start,
			DecayEndTimestamp:   end,
		}
		encryptedBid, _, err := encryptor.ConstructEncryptedBid(reqBid)
		if err != nil {
			t.Fatal(err)
		}

		providerStore := keysstore.New(inmemstorage.New())
		err = providerStore.SetAESKey(address, aesKey)
		if err != nil {
			t.Fatal(err)
		}
		encryptorProvider, err := preconfencryptor.NewEncryptor(keySigner, providerStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
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
		bidderEncryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
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
		providerEncryptor, err := preconfencryptor.NewEncryptor(keySigner, providerStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
		if err != nil {
			t.Fatal(err)
		}
		start := time.Now().UnixMilli()
		end := start + 100000

		bid := &preconfpb.Bid{
			TxHash:              "0xkartik",
			BidAmount:           "10",
			BlockNumber:         2,
			DecayStartTimestamp: start,
			DecayEndTimestamp:   end,
		}

		encryptedBid, nikePrivateKey, err := bidderEncryptor.ConstructEncryptedBid(bid)
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

		_, address, err := bidderEncryptor.VerifyEncryptedPreConfirmation(
			bid,
			providerNikePrivateKey.PublicKey(),
			nikePrivateKey,
			encryptedPreConfirmation,
		)
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

		preconfAddr := common.HexToAddress("0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
		chainID := big.NewInt(31337)
		domainSeparatorBidHash, err := preconfencryptor.ComputeDomainSeparator("PreConfBid", chainID, preconfAddr)
		if err != nil {
			t.Fatal(err)
		}
		hash, err := preconfencryptor.GetBidHash(bid, domainSeparatorBidHash)
		if err != nil {
			t.Fatal(fmt.Errorf("failed to get bid hash %w", err))
		}

		hashStr := hex.EncodeToString(hash)
		// This hash is sourced from the solidity contract to ensure interoperability
		expHash := "447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf"
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
		bidHash := "447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf"
		bidSignature := "5cd1f790192a0ab79661c48f39e77937a6de537ccf6b428682583d13d30294cb113cea12822f821c064c9db918667bf74490535b35b4ef4f28f5d67b133ec22e1b"

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

		chainID := big.NewInt(31337)
		preconfContractAddr := common.HexToAddress("0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
		domainSeparatorPreConfHash, err := preconfencryptor.ComputeDomainSeparator("OpenedCommitment", chainID, preconfContractAddr)
		if err != nil {
			t.Fatal(err)
		}
	
		hash, err := preconfencryptor.GetPreConfirmationHash(preConfirmation, domainSeparatorPreConfHash)
		if err != nil {
			t.Fatal(err)
		}
		hashStr := hex.EncodeToString(hash)
		expHash := "a7f6241be0c5055f054fcbe03d98a1920f0ab874039474401323d8d95930a076"
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

	bidSig := "5cd1f790192a0ab79661c48f39e77937a6de537ccf6b428682583d13d30294cb113cea12822f821c064c9db918667bf74490535b35b4ef4f28f5d67b133ec22e1b"
	bidHash := "447b1a7d708774aa54989ab576b576242ae7fd8a37d4e8f33f0eee751bc72edf"

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
	encryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
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
		_, _, err := encryptor.ConstructEncryptedBid(&preconfpb.Bid{
			TxHash:              bids[i].hash,
			BidAmount:           bids[i].amount,
			BlockNumber:         bids[i].blocknumber,
			DecayStartTimestamp: bids[i].start,
			DecayEndTimestamp:   bids[i].end,
		})
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
	bidderEncryptor, err := preconfencryptor.NewEncryptor(keySigner, bidderStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
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
	providerEncryptor, err := preconfencryptor.NewEncryptor(keySigner, providerStore, big.NewInt(31337), "0xA4AD4f68d0b91CFD19687c881e50f3A00242828c")
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
		bids[i] = &preconfpb.Bid{
			TxHash:              bid.hash,
			BidAmount:           bid.amount,
			BlockNumber:         bid.blocknumber,
			DecayStartTimestamp: bid.start,
			DecayEndTimestamp:   bid.end,
		}
		_, _, err = bidderEncryptor.ConstructEncryptedBid(bids[i])
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
