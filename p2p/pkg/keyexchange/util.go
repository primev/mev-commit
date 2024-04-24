package keyexchange

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func hashData(data []byte) common.Hash {
	hash := crypto.Keccak256Hash(data)
	return hash
}

func parseTimestampMessage(msg string) (string, int64, error) {
	parts := strings.Fields(msg)
	if len(parts) != 5 || parts[0] != "mev-commit" || parts[1] != "bidder" || parts[3] != "setup" {
		return "", 0, fmt.Errorf("message format is incorrect")
	}
	address := parts[2]
	timestamp, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid timestamp")
	}

	return address, timestamp, nil
}

func isTimestampRecent(timestamp int64) bool {
	currentTime := time.Now().Unix()
	difference := currentTime - timestamp

	return difference <= 60
}
