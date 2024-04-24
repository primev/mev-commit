package keykeeper

import (
	"crypto/elliptic"
	"errors"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func SerializePublicKey(pub *ecies.PublicKey) []byte {
	return elliptic.MarshalCompressed(elliptic.P256(), pub.X, pub.Y)
}

func DeserializePublicKey(data []byte) (*ecies.PublicKey, error) {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecies.PublicKey{
		X:      x,
		Y:      y,
		Curve:  elliptic.P256(),
		Params: ecies.ECIES_AES128_SHA256,
	}, nil
}
