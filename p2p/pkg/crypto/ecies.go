package crypto

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func SerializeEciesPublicKey(pub *ecies.PublicKey) []byte {
	ecdsaPub := &ecdsa.PublicKey{
		Curve: pub.Curve,
		X:     pub.X,
		Y:     pub.Y,
	}
	return crypto.CompressPubkey(ecdsaPub)
}

func DeserializeEciesPublicKey(data []byte) (*ecies.PublicKey, error) {
	ecdsaPub, err := crypto.DecompressPubkey(data)
	if err != nil {
		return nil, err
	}
	return ecies.ImportECDSAPublic(ecdsaPub), nil
}
