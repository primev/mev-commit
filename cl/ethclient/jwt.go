package ethclient

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtRoundTripper struct {
	underlyingTransport http.RoundTripper
	jwtSecret           []byte
}

func newJWTRoundTripper(transport http.RoundTripper, jwtSecret []byte) *jwtRoundTripper {
	return &jwtRoundTripper{
		underlyingTransport: transport,
		jwtSecret:           jwtSecret,
	}
}

func (t *jwtRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(t.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("jwt token string, err: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)
	resp, err := t.underlyingTransport.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("round trip, err: %w", err)
	}

	return resp, nil
}

// LoadJWTHexFile loads a hex encoded JWT secret from the provided file.
func LoadJWTHexFile(file string) ([]byte, error) {
	jwtHex, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read jwt file, err: %w", err)
	}

	jwtHex = bytes.TrimSpace(jwtHex)
	jwtHex = bytes.TrimPrefix(jwtHex, []byte("0x"))

	jwtBytes, err := hex.DecodeString(string(jwtHex))
	if err != nil {
		return nil, fmt.Errorf("decode jwt file, err: %w", err)
	}

	return jwtBytes, nil
}
