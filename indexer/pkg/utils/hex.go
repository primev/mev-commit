package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func Strip0x(s string) string {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return s[2:]
	}
	return s
}

func HexToBytes(s string) []byte {
	s = Strip0x(strings.TrimSpace(s))
	if s == "" {
		return nil
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}
func ParseBigString(v any) (string, bool) {
	switch t := v.(type) {
	case nil:
		return "", false
	case string:
		z := strings.ReplaceAll(strings.TrimSpace(t), ",", "")
		if z == "" {
			return "", false
		}
		if strings.HasPrefix(z, "0x") || strings.HasPrefix(z, "0X") {
			bi := new(big.Int)
			if _, ok := bi.SetString(z[2:], 16); ok {
				return bi.String(), true
			}
			return "", false
		}
		if _, ok := new(big.Int).SetString(z, 10); ok {
			return z, true
		}
		return "", false
	case float64:
		return strconv.FormatFloat(t, 'f', 0, 64), true
	case json.Number:
		return t.String(), true
	default:
		return fmt.Sprintf("%v", t), true
	}
}
func GetString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func HexToUint64(h string) (uint64, error) {
	h = strings.TrimSpace(h)
	if strings.HasPrefix(h, "0x") {
		h = h[2:]
	}
	if h == "" {
		return 0, fmt.Errorf("empty hex")
	}
	v, ok := new(big.Int).SetString(h, 16)
	if !ok {
		return 0, fmt.Errorf("bad hex %q", h)
	}
	return v.Uint64(), nil
}
