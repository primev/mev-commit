package store

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/armon/go-radix"
	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
)

type Store struct {
	*radix.Tree
}

type EncryptedPreConfirmationWithDecrypted struct {
	*preconfpb.EncryptedPreConfirmation
	*preconfpb.PreConfirmation
	TxnHash common.Hash
}

func NewStore() *Store {
	return &Store{
		Tree: radix.New(),
	}
}

func (s *Store) LastBlock() (uint64, error) {
	val, ok := s.Tree.Get("last_block")
	if !ok {
		return 0, nil
	}
	return val.(uint64), nil
}

func (s *Store) SetLastBlock(blockNum uint64) error {
	_, _ = s.Tree.Insert("last_block", blockNum)
	return nil
}

func (s *Store) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) {
	key := fmt.Sprintf("cm/%d/%s",
		commitment.Bid.BlockNumber,
		string(commitment.EncryptedPreConfirmation.Commitment),
	)
	_, _ = s.Tree.Insert(key, commitment)
}

func (s *Store) GetCommitmentsByBlockNumber(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	blockCommitmentsKey := fmt.Sprintf("cm/%d", blockNum)
	commitments := make([]*EncryptedPreConfirmationWithDecrypted, 0)
	s.Tree.WalkPrefix(blockCommitmentsKey, func(key string, value interface{}) bool {
		commitments = append(commitments, value.(*EncryptedPreConfirmationWithDecrypted))
		return false
	})
	return commitments, nil
}

func (s *Store) DeleteCommitmentByBlockNumber(blockNum int64) error {
	blockCommitmentsKey := fmt.Sprintf("cm/%d", blockNum)
	_ = s.Tree.DeletePrefix(blockCommitmentsKey)
	return nil
}

func (s *Store) DeleteCommitmentByIndex(blockNum int64, index [32]byte) error {
	key := fmt.Sprintf("cm/%d/%s", blockNum, string(index[:]))
	_, _ = s.Tree.Delete(key)
	return nil
}

func (s *Store) SetCommitmentIndexByCommitmentDigest(cDigest, cIndex [32]byte) error {
	s.Tree.WalkPrefix("cm/", func(key string, value interface{}) bool {
		commitment := value.(*EncryptedPreConfirmationWithDecrypted)
		if bytes.Equal(commitment.EncryptedPreConfirmation.Commitment, cDigest[:]) {
			commitment.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
			return true
		}
		return false
	})

	return nil
}

func (s *Store) SetBalance(bidder common.Address, windowNumber, depositedAmount *big.Int) error {
	key := fmt.Sprintf("bbs/%s/%s", windowNumber.String(), bidder.String())
	_, _ = s.Tree.Insert(key, depositedAmount)
	return nil
}

func (s *Store) GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error) {
	key := fmt.Sprintf("bbs/%s/%s", windowNumber.String(), bidder.String())
	val, ok := s.Tree.Get(key)
	if !ok {
		return nil, nil
	}
	return val.(*big.Int), nil
}

func (s *Store) ClearBalances(windowNumber *big.Int) error {
	if windowNumber == nil || windowNumber.Cmp(big.NewInt(0)) == -1 {
		return nil
	}
	windows := make([]*big.Int, 0)
	s.Tree.WalkPrefix("bbs/", func(key string, value interface{}) bool {
		parts := strings.Split(key, "/")
		if len(parts) != 3 {
			return false
		}
		w, ok := new(big.Int).SetString(parts[1], 10)
		if !ok {
			return false
		}
		switch w.Cmp(windowNumber) {
		case -1:
			windows = append(windows, w)
		case 0:
			windows = append(windows, w)
			return true
		}
		return false
	})

	for _, w := range windows {
		key := fmt.Sprintf("bbs/%s", w.String())
		_ = s.Tree.DeletePrefix(key)
	}
	return nil
}

func (s *Store) GetBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	blockNumber int64,
) (*big.Int, error) {
	key := fmt.Sprintf("bbs/%s/%s/%d", window.String(), bidder.String(), blockNumber)
	val, ok := s.Tree.Get(key)
	if !ok {
		return nil, nil
	}
	return val.(*big.Int), nil
}

func (s *Store) SetBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	key := fmt.Sprintf("bbs/%s/%s/%d", window.String(), bidder.String(), blockNumber)
	_, _ = s.Tree.Insert(key, amount)
	return nil
}

func (s *Store) RefundBalanceForBlock(
	bidder common.Address,
	window *big.Int,
	amount *big.Int,
	blockNumber int64,
) error {
	key := fmt.Sprintf("bbs/%s/%s/%d", window.String(), bidder.String(), blockNumber)
	val, ok := s.Tree.Get(key)
	if !ok {
		_, _ = s.Tree.Insert(key, amount)
		return nil
	}
	amount.Add(amount, val.(*big.Int))
	_, _ = s.Tree.Insert(key, amount)
	return nil
}
