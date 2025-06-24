package transfer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	bridgetransfer "github.com/primev/mev-commit/bridge/standard/pkg/transfer"
	"github.com/primev/mev-commit/x/keysigner"
)

func SetTransferFunc(f func(
	amount *big.Int,
	destAddress common.Address,
	signer keysigner.KeySigner,
	settlementRPCUrl string,
	l1RPCUrl string,
	l1ContractAddr common.Address,
	settlementContractAddr common.Address,
) (bridgetransfer.Transfer, error)) func() {
	prev := transferFunc
	transferFunc = f
	return func() {
		transferFunc = prev
	}
}
