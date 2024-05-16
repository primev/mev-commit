package handshake_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p/core"
	"github.com/primev/mev-commit/p2p/pkg/keykeeper"
	mockkeysigner "github.com/primev/mev-commit/p2p/pkg/keykeeper/keysigner/mock"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/p2p/libp2p/internal/handshake"
	p2ptest "github.com/primev/mev-commit/p2p/pkg/p2p/testing"
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
		privKey1, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		address1 := common.HexToAddress("0x1")
		ks1 := mockkeysigner.NewMockKeySigner(privKey1, address1)
		kk1, err := keykeeper.NewProviderKeyKeeper(ks1)
		if err != nil {
			t.Fatal(err)
		}
		privKey2, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}

		address2 := common.HexToAddress("0x2")
		ks2 := mockkeysigner.NewMockKeySigner(privKey2, address2)
		kk2, err := keykeeper.NewProviderKeyKeeper(ks2)
		if err != nil {
			t.Fatal(err)
		}
		hs1, err := handshake.New(
			kk1,
			p2p.PeerTypeProvider,
			"test",
			&testSigner{address: address2},
			&testRegister{},
			func(p core.PeerID) (common.Address, error) {
				return address2, nil
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		hs2, err := handshake.New(
			kk2,
			p2p.PeerTypeProvider,
			"test",
			&testSigner{address: address1},
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
			if !p.Keys.NIKEPublicKey.Equal(kk2.GetNIKEPublicKey()) {
				t.Errorf("expected nike pk %s, got %s", p.Keys.NIKEPublicKey.Bytes(), kk2.GetNIKEPublicKey().Bytes())
				return
			}
			if !p.Keys.PKEPublicKey.ExportECDSA().Equal(kk2.GetECIESPublicKey().ExportECDSA()) {
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
		if !p.Keys.NIKEPublicKey.Equal(kk1.GetNIKEPublicKey()) {
			t.Fatalf("expected nike pk %s, got %s", p.Keys.NIKEPublicKey.Bytes(), kk1.GetNIKEPublicKey().Bytes())
		}
		if !p.Keys.PKEPublicKey.ExportECDSA().Equal(kk1.GetECIESPublicKey().ExportECDSA()) {
			t.Fatalf("expected pke pk is not equal to present")
		}
		<-done
	})
}
