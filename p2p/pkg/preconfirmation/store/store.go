package store

import (
	"encoding/binary"
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

	cmtIndexNS = "ci/"
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
	cmtIndexKey = func(cIndex []byte) string {
		return fmt.Sprintf("%s%s", cmtIndexNS, string(cIndex))
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

func (s *Store) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var writer storage.Writer
	if w, ok := s.st.(storage.Batcher); ok {
		batch := w.Batch()
		writer = batch
		defer func() {
			switch {
			case err != nil:
				batch.Reset()
			case err == nil:
				err = batch.Write()
			}
		}()

	} else {
		writer = s.st
	}

	key := commitmentKey(commitment.Bid.BlockNumber, commitment.EncryptedPreConfirmation.Commitment)

	buf, err := msgpack.Marshal(commitment)
	if err != nil {
		return err
	}

	if err := writer.Put(key, buf); err != nil {
		return err
	}

	cIndexKey := cmtIndexKey(commitment.EncryptedPreConfirmation.Commitment)
	blkNumBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(blkNumBuf, uint64(commitment.Bid.BlockNumber))

	return writer.Put(cIndexKey, blkNumBuf)
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
	s.mu.RLock()
	blkNumBuf, err := s.st.Get(cmtIndexKey(cDigest[:]))
	s.mu.RUnlock()
	switch {
	case err == storage.ErrKeyNotFound:
		// this would happen for most of the commitments as the node only
		// stores the commitments it is involved in.
		return nil
	case err != nil:
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	blkNum := binary.LittleEndian.Uint64(blkNumBuf)
	commitmentKey := commitmentKey(int64(blkNum), cDigest[:])
	cmtBuf, err := s.st.Get(commitmentKey)
	if err != nil {
		return err
	}

	cmt := new(EncryptedPreConfirmationWithDecrypted)
	err = msgpack.Unmarshal(cmtBuf, cmt)
	if err != nil {
		return err
	}

	cmt.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
	buf, err := msgpack.Marshal(cmt)
	if err != nil {
		return err
	}

	return s.st.Put(commitmentKey, buf)
}

func (s *Store) ClearCommitmentIndexes(uptoBlock int64) error {
	keys := make([]string, 0)
	s.mu.RLock()
	err := s.st.WalkPrefix(cmtIndexNS, func(key string, val []byte) bool {
		blkNum := binary.LittleEndian.Uint64([]byte(val))
		if blkNum < uint64(uptoBlock) {
			keys = append(keys, key)
		}
		return false
	})
	s.mu.RUnlock()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		err := s.st.Delete(key)
		if err != nil {
			return err
		}
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
