package preconfencryptor_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primevprotocol/mev-commit/p2p/pkg/keykeeper"
	mockkeysigner "github.com/primevprotocol/mev-commit/p2p/pkg/keykeeper/keysigner/mock"
	"github.com/primevprotocol/mev-commit/p2p/pkg/signer/preconfencryptor"
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
		keyKeeper, err := keykeeper.NewBidderKeyKeeper(keySigner)
		if err != nil {
			t.Fatal(err)
		}
		encryptor := preconfencryptor.NewEncryptor(keyKeeper)

		start := time.Now().UnixMilli()
		end := start + 100000
		_, encryptedBid, err := encryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end)
		if err != nil {
			t.Fatal(err)
		}

		providerKeyKeeper, err := keykeeper.NewProviderKeyKeeper(keySigner)
		if err != nil {
			t.Fatal(err)
		}
		providerKeyKeeper.SetAESKey(address, keyKeeper.AESKey)
		encryptorProvider := preconfencryptor.NewEncryptor(providerKeyKeeper)
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

		keySigner := mockkeysigner.NewMockKeySigner(bidderKey, crypto.PubkeyToAddress(bidderKey.PublicKey))
		bidderKeyKeeper, err := keykeeper.NewBidderKeyKeeper(keySigner)
		if err != nil {
			t.Fatal(err)
		}
		bidderEncryptor := preconfencryptor.NewEncryptor(bidderKeyKeeper)
		providerKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}

		bidderAddress := crypto.PubkeyToAddress(bidderKey.PublicKey)
		keySigner = mockkeysigner.NewMockKeySigner(providerKey, crypto.PubkeyToAddress(providerKey.PublicKey))
		providerKeyKeeper, err := keykeeper.NewProviderKeyKeeper(keySigner)
		if err != nil {
			t.Fatal(err)
		}

		providerKeyKeeper.SetAESKey(bidderAddress, bidderKeyKeeper.AESKey)
		providerEncryptor := preconfencryptor.NewEncryptor(providerKeyKeeper)
		start := time.Now().UnixMilli()
		end := start + 100000

		bid, encryptedBid, err := bidderEncryptor.ConstructEncryptedBid("0xkartik", "10", 2, start, end)
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

		_, address, err := bidderEncryptor.VerifyEncryptedPreConfirmation(providerKeyKeeper.GetNIKEPublicKey(), bid.Digest, encryptedPreConfirmation)
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
