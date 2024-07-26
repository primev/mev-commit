package store

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	commitmentNS = "cm/"

	// block winners
	blockWinnerNS = "bw/"
)

var (
	commitmentKey = func(blockNum int64, index []byte) string {
		return fmt.Sprintf("%s%d/%s", commitmentNS, blockNum, string(index))
	}
	blockCommitmentPrefix = func(blockNum int64) string {
		return fmt.Sprintf("%s%d", commitmentNS, blockNum)
	}
	blockWinnerKey = func(blockNumber int64) string {
		return fmt.Sprintf("%s%d", blockWinnerNS, blockNumber)
	}
)

type Store struct {
	mu sync.RWMutex
	st storage.Storage
}

type EncryptedPreConfirmationWithDecrypted struct {
	*preconfpb.EncryptedPreConfirmation
	*preconfpb.PreConfirmation
	TxnHash common.Hash
}

type BlockWinner struct {
	BlockNumber int64
	Winner      common.Address
}

func New(st storage.Storage) *Store {
	return &Store{
		st: st,
	}
}

func (s *Store) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := commitmentKey(commitment.Bid.BlockNumber, commitment.EncryptedPreConfirmation.Commitment)

	buf, err := msgpack.Marshal(commitment)
	if err != nil {
		return err
	}
	return s.st.Put(key, buf)
}

func (s *Store) GetCommitments(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockCommitmentsKey := blockCommitmentPrefix(blockNum)
	commitments := make([]*EncryptedPreConfirmationWithDecrypted, 0)

	err := s.st.WalkPrefix(blockCommitmentsKey, func(key string, value []byte) bool {
		commitment := new(EncryptedPreConfirmationWithDecrypted)
		err := msgpack.Unmarshal(value, commitment)
		if err != nil {
			return false
		}
		commitments = append(commitments, commitment)
		return false
	})
	if err != nil {
		return nil, err
	}
	return commitments, nil
}

func (s *Store) ClearBlockNumber(blockNum int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.st.DeletePrefix(blockCommitmentPrefix(blockNum))
	if err != nil {
		return err
	}

	return s.st.Delete(blockWinnerKey(blockNum))
}

func (s *Store) DeleteCommitmentByDigest(blockNum int64, digest [32]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Delete(commitmentKey(blockNum, digest[:]))
}

func (s *Store) SetCommitmentIndexByDigest(cDigest, cIndex [32]byte) error {
	var cmt *EncryptedPreConfirmationWithDecrypted

	s.mu.RLock()
	err := s.st.WalkPrefix(commitmentNS, func(key string, value []byte) bool {
		c := new(EncryptedPreConfirmationWithDecrypted)
		err := msgpack.Unmarshal(value, c)
		if err != nil {
			return false
		}
		if bytes.Equal(c.EncryptedPreConfirmation.Commitment, cDigest[:]) {
			cmt = c
			return true
		}
		return false
	})
	s.mu.RUnlock()
	if err != nil {
		return err
	}

	if cmt != nil {
		cmt.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
		return s.AddCommitment(cmt)
	}
	return nil
}

func (s *Store) AddWinner(blockWinner *BlockWinner) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	buf, err := msgpack.Marshal(blockWinner)
	if err != nil {
		return err
	}

	return s.st.Put(blockWinnerKey(blockWinner.BlockNumber), buf)
}

func (s *Store) BlockWinners() ([]*BlockWinner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	winners := make([]*BlockWinner, 0)
	err := s.st.WalkPrefix(blockWinnerNS, func(key string, value []byte) bool {
		w := new(BlockWinner)
		err := msgpack.Unmarshal(value, w)
		if err != nil {
			return false
		}
		winners = append(winners, w)
		return false
	})
	if err != nil {
		return nil, err
	}
	return winners, nil
}
