package handshake_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/libp2p/go-libp2p/core"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/keysstore"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/p2p/libp2p/internal/handshake"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	mockkeysigner "github.com/primev/mev-commit/x/keysigner/mock"
)

type testRegister struct{}

func (t *testRegister) CheckProviderRegistered(
	_ context.Context,
	_ common.Address,
) bool {
	return true
}

type testSigner struct {
	address common.Address
}

func (t *testSigner) Sign(_ *ecdsa.PrivateKey, _ []byte) ([]byte, error) {
	return []byte("signature"), nil
}

func (t *testSigner) Verify(_ []byte, _ []byte) (bool, common.Address, error) {
	return true, t.address, nil
}

func TestHandshake(t *testing.T) {
	t.Parallel()

	t.Run("ok", func(t *testing.T) {
		privKey1, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		address1 := common.HexToAddress("0x1")
		ks1 := mockkeysigner.NewMockKeySigner(privKey1, address1)
		privKey2, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		address2 := common.HexToAddress("0x2")
		ks2 := mockkeysigner.NewMockKeySigner(privKey2, address2)
		store1 := keysstore.New(inmemstorage.New())
		sk1, pk1, err := p2pcrypto.GenerateKeyPairBN254()
		if err != nil {
			t.Fatal(err)
		}
		err = store1.SetBN254PrivateKey(sk1)
		if err != nil {
			t.Fatal(err)
		}
		err = store1.SetBN254PublicKey(pk1)
		if err != nil {
			t.Fatal(err)
		}
		prvKey1, err := ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
		if err != nil {
			t.Fatal(err)
		}
		err = store1.SetECIESPrivateKey(prvKey1)
		if err != nil {
			t.Fatal(err)
		}
		providerKeys1 := p2p.Keys{
			PKEPublicKey:  &prvKey1.PublicKey,
			NIKEPublicKey: pk1,
		}
		hs1, err := handshake.New(
			ks1,
			p2p.PeerTypeProvider,
			"test",
			&testSigner{address: address2},
			&providerKeys1,
			&testRegister{},
			func(p core.PeerID) (common.Address, error) {
				return address2, nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}
		store2 := keysstore.New(inmemstorage.New())
		sk2, pk2, err := p2pcrypto.GenerateKeyPairBN254()
		if err != nil {
			t.Fatal(err)
		}
		err = store2.SetBN254PrivateKey(sk2)
		if err != nil {
			t.Fatal(err)
		}
		err = store2.SetBN254PublicKey(pk2)
		if err != nil {
			t.Fatal(err)
		}
		prvKey2, err := ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
		if err != nil {
			t.Fatal(err)
		}
		err = store2.SetECIESPrivateKey(prvKey2)
		if err != nil {
			t.Fatal(err)
		}

		providerKeys2 := p2p.Keys{
			PKEPublicKey:  &prvKey2.PublicKey,
			NIKEPublicKey: pk2,
		}

		hs2, err := handshake.New(
			ks2,
			p2p.PeerTypeProvider,
			"test",
			&testSigner{address: address1},
			&providerKeys2,
			&testRegister{},
			func(p core.PeerID) (common.Address, error) {
				return address1, nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		out, in := p2ptest.NewDuplexStream()

		done := make(chan struct{})
		go func() {
			defer close(done)

			p, err := hs1.Handle(context.Background(), in, core.PeerID("test2"))
			if err != nil {
				t.Error(err)
				return
			}
			if p.EthAddress != address2 {
				t.Errorf(
					"expected eth address %s, got %s",
					address2, p.EthAddress,
				)
				return
			}
			if p.Type != p2p.PeerTypeProvider {
				t.Errorf("expected peer type %s, got %s", p2p.PeerTypeProvider, p.Type)
				return
			}
			bn254pk, err := store2.BN254PublicKey()
			if err != nil {
				t.Error(err)
				return
			}
			if !p.Keys.NIKEPublicKey.Equal(bn254pk) {
				t.Errorf("expected nike pk %s, got %s", p.Keys.NIKEPublicKey.Bytes(), p2pcrypto.BN254PublicKeyToBytes(bn254pk))
				return
			}
			prvKey2, err := store2.ECIESPrivateKey()
			if err != nil {
				t.Error(err)
				return
			}

			if !p.Keys.PKEPublicKey.ExportECDSA().Equal(prvKey2.PublicKey.ExportECDSA()) {
				t.Error("expected pke pk is not equal to present")
				return
			}
		}()

		p, err := hs2.Handshake(context.Background(), core.PeerID("test1"), out)
		if err != nil {
			t.Fatal(err)
		}
		if p.EthAddress != address1 {
			t.Fatalf("expected eth address %s, got %s", address1, p.EthAddress)
		}
		if p.Type != p2p.PeerTypeProvider {
			t.Fatalf("expected peer type %s, got %s", p2p.PeerTypeProvider, p.Type)
		}
		bn254pk, err := store1.BN254PublicKey()
		if err != nil {
			t.Fatal(err)
		}
		if !p.Keys.NIKEPublicKey.Equal(bn254pk) {
			t.Fatalf("expected nike pk %s, got %s", p.Keys.NIKEPublicKey.Bytes(), p2pcrypto.BN254PublicKeyToBytes(bn254pk))
		}
		prvKey1, err = store1.ECIESPrivateKey()
		if err != nil {
			t.Fatal(err)
		}

		if !p.Keys.PKEPublicKey.ExportECDSA().Equal(prvKey1.PublicKey.ExportECDSA()) {
			t.Fatalf("expected pke pk is not equal to present")
		}
		<-done
	})
}
