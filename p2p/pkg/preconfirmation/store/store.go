package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/storage"
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

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(commitment)
	if err != nil {
		return err
	}
	return s.st.Put(key, buf.Bytes())
}

func (s *Store) GetCommitments(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockCommitmentsKey := blockCommitmentPrefix(blockNum)
	commitments := make([]*EncryptedPreConfirmationWithDecrypted, 0)

	s.st.WalkPrefix(blockCommitmentsKey, func(key string, value []byte) bool {
		commitment := new(EncryptedPreConfirmationWithDecrypted)
		err := gob.NewDecoder(bytes.NewReader(value)).Decode(commitment)
		if err != nil {
			return false
		}
		commitments = append(commitments, commitment)
		return false
	})
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
	s.mu.Lock()
	defer s.mu.Unlock()

	s.st.WalkPrefix(commitmentNS, func(key string, value []byte) bool {
		c := new(EncryptedPreConfirmationWithDecrypted)
		err := gob.NewDecoder(bytes.NewReader(value)).Decode(c)
		if err != nil {
			return false
		}
		if bytes.Equal(c.EncryptedPreConfirmation.Commitment, cDigest[:]) {
			c.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
			var buf bytes.Buffer
			err = gob.NewEncoder(&buf).Encode(c)
			if err != nil {
				return false
			}
			err = s.st.Put(key, buf.Bytes())
			if err != nil {
				return false
			}
			return true
		}
		return false
	})

	return nil
}

func (s *Store) AddWinner(blockWinner *BlockWinner) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(blockWinner)
	if err != nil {
		return err
	}

	return s.st.Put(blockWinnerKey(blockWinner.BlockNumber), buf.Bytes())
}

func (s *Store) BlockWinners() ([]*BlockWinner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	winners := make([]*BlockWinner, 0)
	s.st.WalkPrefix(blockWinnerNS, func(key string, value []byte) bool {
		w := new(BlockWinner)
		err := gob.NewDecoder(bytes.NewReader(value)).Decode(w)
		if err != nil {
			return false
		}
		winners = append(winners, w)
		return false
	})
	return winners, nil
}
