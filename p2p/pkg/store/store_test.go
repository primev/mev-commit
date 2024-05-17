package store_test

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/store"
)

func TestStore(t *testing.T) {
	t.Parallel()

	st := store.NewStore()

	t.Run("last block", func(t *testing.T) {
		lastBlock, err := st.LastBlock()
		if err != nil {
			t.Fatal(err)
		}

		if lastBlock != 0 {
			t.Fatalf("expected 0, got %d", lastBlock)
		}

		if err := st.SetLastBlock(10); err != nil {
			t.Fatal(err)
		}

		lastBlock, err = st.LastBlock()
		if err != nil {
			t.Fatal(err)
		}

		if lastBlock != 10 {
			t.Fatalf("expected 10, got %d", lastBlock)
		}
	})

	t.Run("commitments", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			var blkNum int64 = 1
			if i > 5 {
				blkNum = 2
			}
			commitment := &store.EncryptedPreConfirmationWithDecrypted{
				EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
					Commitment: common.BigToHash(big.NewInt(int64(i))).Bytes(),
				},
				PreConfirmation: &preconfpb.PreConfirmation{
					Bid: &preconfpb.Bid{
						BlockNumber: blkNum,
					},
				},
			}

			st.AddCommitment(commitment)
		}

		for i := 1; i <= 10; i++ {
			err := st.SetCommitmentIndexByCommitmentDigest(
				common.BigToHash(big.NewInt(int64(i))),
				common.BigToHash(big.NewInt(int64(i))),
			)
			if err != nil {
				t.Fatal(err)
			}
		}

		commitments, err := st.GetCommitmentsByBlockNumber(1)
		if err != nil {
			t.Fatal(err)
		}

		if len(commitments) != 5 {
			t.Fatalf("expected 5, got %d", len(commitments))
		}

		for i := 1; i <= 5; i++ {
			if !bytes.Equal(commitments[i-1].Commitment, common.BigToHash(big.NewInt(int64(i))).Bytes()) {
				t.Fatalf("expected %d, got %s", i, commitments[i-1].Digest)
			}
		}

		err = st.DeleteCommitmentByBlockNumber(1)
		if err != nil {
			t.Fatal(err)
		}

		commitments, err = st.GetCommitmentsByBlockNumber(1)
		if err != nil {
			t.Fatal(err)
		}

		if len(commitments) != 0 {
			t.Fatalf("expected 0, got %d", len(commitments))
		}

		for i := 6; i <= 10; i++ {
			err := st.DeleteCommitmentByIndex(2, common.BigToHash(big.NewInt(int64(i))))
			if err != nil {
				t.Fatal(err)
			}
		}

		commitments, err = st.GetCommitmentsByBlockNumber(2)
		if err != nil {
			t.Fatal(err)
		}

		if len(commitments) != 0 {
			t.Fatalf("expected 0, got %d", len(commitments))
		}
	})

	t.Run("balances", func(t *testing.T) {
		if err := st.SetBalance(common.HexToAddress("0x123"), big.NewInt(1), big.NewInt(10)); err != nil {
			t.Fatal(err)
		}
		val, err := st.GetBalance(common.HexToAddress("0x123"), big.NewInt(1))
		if err != nil {
			t.Fatal(err)
		}
		if val.Cmp(big.NewInt(10)) != 0 {
			t.Fatalf("expected 10, got %s", val.String())
		}

		for i := 1; i <= 10; i++ {
			err := st.SetBalanceForBlock(common.HexToAddress("0x123"), big.NewInt(1), big.NewInt(10), int64(i))
			if err != nil {
				t.Fatal(err)
			}

			val, err := st.GetBalanceForBlock(common.HexToAddress("0x123"), big.NewInt(1), int64(i))
			if err != nil {
				t.Fatal(err)
			}
			if val.Cmp(big.NewInt(10)) != 0 {
				t.Fatalf("expected 10, got %s", val.String())
			}

			err = st.RefundBalanceForBlock(common.HexToAddress("0x123"), big.NewInt(1), big.NewInt(10), int64(i))
			if err != nil {
				t.Fatal(err)
			}

			val, err = st.GetBalanceForBlock(common.HexToAddress("0x123"), big.NewInt(1), int64(i))
			if err != nil {
				t.Fatal(err)
			}
			if val.Cmp(big.NewInt(20)) != 0 {
				t.Fatalf("expected 20, got %s", val.String())
			}
		}

		windows, err := st.ClearBalances(big.NewInt(12))
		if err != nil {
			t.Fatal(err)
		}
		if len(windows) != 1 {
			t.Fatalf("expected 1, got %d", len(windows))
		}

		for i := 1; i <= 10; i++ {
			val, err := st.GetBalanceForBlock(common.HexToAddress("0x123"), big.NewInt(1), int64(i))
			if err != nil {
				t.Fatal(err)
			}
			if val != nil {
				t.Fatalf("expected nil, got %s", val.String())
			}
		}
	})

	t.Run("AES keys", func(t *testing.T) {
		bidder := common.HexToAddress("0x456")
		key := []byte("my_secret_aes_key")

		// Set AES key
		err := st.SetAESKey(bidder, key)
		if err != nil {
			t.Fatal(err)
		}

		// Get AES key
		retrievedKey, err := st.GetAESKey(bidder)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(retrievedKey, key) {
			t.Fatalf("expected %s, got %s", key, retrievedKey)
		}

		// Test non-existing key retrieval
		nonExistentBidder := common.HexToAddress("0x789")
		retrievedKey, err = st.GetAESKey(nonExistentBidder)
		if err != nil {
			t.Fatal(err)
		}
		if retrievedKey != nil {
			t.Fatalf("expected nil, got %s", retrievedKey)
		}
	})

	t.Run("ECIES keys", func(t *testing.T) {
		// Generate a new ECIES private key
		ecdsaKey, err := crypto.GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		eciesKey := ecies.ImportECDSA(ecdsaKey)

		// Set ECIES private key
		err = st.SetECIESPrivateKey(eciesKey)
		if err != nil {
			t.Fatal(err)
		}

		// Get ECIES private key
		retrievedKey, err := st.GetECIESPrivateKey()
		if err != nil {
			t.Fatal(err)
		}
		if !eciesKey.ExportECDSA().Equal(retrievedKey.ExportECDSA()) {
			t.Fatalf("expected %v, got %v", eciesKey, retrievedKey)
		}

		// Ensure key retrieval when no key is set returns nil
		st2 := store.NewStore()
		retrievedKey, err = st2.GetECIESPrivateKey()
		if err != nil {
			t.Fatal(err)
		}
		if retrievedKey != nil {
			t.Fatalf("expected nil, got %v", retrievedKey)
		}
	})

	t.Run("Nike keys", func(t *testing.T) {
		// Generate a new Nike private key
		nikeCurve := ecdh.X25519()
		nikePrivateKey, err := nikeCurve.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		// Set Nike private key
		err = st.SetNikePrivateKey(nikePrivateKey)
		if err != nil {
			t.Fatal(err)
		}

		// Get Nike private key
		retrievedKey, err := st.GetNikePrivateKey()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(nikePrivateKey.Bytes(), retrievedKey.Bytes()) {
			t.Fatalf("expected %x, got %x", nikePrivateKey.Bytes(), retrievedKey.Bytes())
		}

		// Ensure key retrieval when no key is set returns nil
		st2 := store.NewStore()
		retrievedKey, err = st2.GetNikePrivateKey()
		if err != nil {
			t.Fatal(err)
		}
		if retrievedKey != nil {
			t.Fatalf("expected nil, got %x", retrievedKey.Bytes())
		}
	})
}
