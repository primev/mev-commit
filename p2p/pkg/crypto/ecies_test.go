package crypto

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func TestSerializeEciesPublicKey(t *testing.T) {
	privKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		t.Fatalf("Failed to generate ECIES key: %v", err)
	}

	pubKey := &privKey.PublicKey
	serialized := SerializeEciesPublicKey(pubKey)

	if len(serialized) == 0 {
		t.Fatal("Serialized public key should not be empty")
	}
}

func TestDeserializeEciesPublicKey(t *testing.T) {
	privKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		t.Fatalf("Failed to generate ECIES key: %v", err)
	}

	pubKey := &privKey.PublicKey
	serialized := SerializeEciesPublicKey(pubKey)

	deserializedPubKey, err := DeserializeEciesPublicKey(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize ECIES public key: %v", err)
	}

	if pubKey.X.Cmp(deserializedPubKey.X) != 0 || pubKey.Y.Cmp(deserializedPubKey.Y) != 0 {
		t.Fatal("Deserialized public key does not match original")
	}
}

func TestDeserializeEciesPublicKeyInvalidData(t *testing.T) {
	invalidData := []byte("invalid data")

	_, err := DeserializeEciesPublicKey(invalidData)
	if err == nil {
		t.Fatal("Expected error when deserializing invalid public key data")
	}
}
