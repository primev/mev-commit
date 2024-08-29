package txmonitor

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestUnwrapABICustomError(t *testing.T) {
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"uint256","name":"available","type":"uint256"},{"internalType":"uint256","name":"required","type":"uint256"}],"name":"InsufficientBalance","type":"error"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		t.Fatal(err)
	}

	rpcErr := &rpcAPIError{
		code: 3,
		msg:  "execution reverted",
		err:  "0xcf4791810000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006f",
	}
	expected := "InsufficientBalance(available=0, required=111)"
	result := unwrapABIError(rpcErr, &contractABI)
	if result != expected {
		t.Errorf("Expected %q, but got %q", expected, result)
	}
}

type rpcAPIError struct {
	code int
	msg  string
	err  string
}

func (e *rpcAPIError) ErrorCode() int         { return e.code }
func (e *rpcAPIError) Error() string          { return e.msg }
func (e *rpcAPIError) ErrorData() interface{} { return e.err }
