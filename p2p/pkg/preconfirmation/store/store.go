package store

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

	indexToDigestNS = "id/"
)

var (
	MaxBidAmount, _ = new(big.Int).SetString("1000000000000000000000000000", 10) // 1e24
)

var (
	commitmentKey = func(blockNum int64, bidAmt string, index []byte) string {
		bidAmtInt, ok := new(big.Int).SetString(bidAmt, 10)
		if !ok {
			return ""
		}
		invertedBidAmount := new(big.Int).Sub(MaxBidAmount, bidAmtInt)
		paddedBidAmountHex := fmt.Sprintf("%064x", invertedBidAmount)

		return fmt.Sprintf("%s%d/%s/%s", commitmentNS, blockNum, paddedBidAmountHex, string(index))
	}
	parseBlockNumFromCommitmentKey = func(key string) (int64, error) {
		splits := strings.Split(key, "/")
		if len(splits) != 4 || !strings.HasPrefix(commitmentNS, splits[0]) {
			return 0, fmt.Errorf("invalid commitment key format: %s", key)
		}
		blockNum, err := strconv.Atoi(splits[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse block number from key %s: %w", key, err)
		}
		return int64(blockNum), nil
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
	indexToDigestKey = func(cIndex []byte) string {
		return fmt.Sprintf("%s%s", indexToDigestNS, string(cIndex))
	}
)

type Store struct {
	mu sync.RWMutex
	st storage.Storage
}

type CommitmentStatus string

const (
	CommitmentStatusPending CommitmentStatus = "pending"
	CommitmentStatusStored  CommitmentStatus = "stored"
	CommitmentStatusOpened  CommitmentStatus = "opened"
	CommitmentStatusSettled CommitmentStatus = "settled"
	CommitmentStatusSlashed CommitmentStatus = "slashed"
	CommitmentStatusFailed  CommitmentStatus = "failed"
)

type Commitment struct {
	*preconfpb.EncryptedPreConfirmation
	*preconfpb.PreConfirmation
	Status  CommitmentStatus
	Details string
	Payment string
	Refund  string
}

type BlockWinner struct {
	BlockNumber int64
	Winner      common.Address
}

type CommitmentIndexValue struct {
	BlockNumber int64
	BidAmount   string
}

func New(st storage.Storage) *Store {
	return &Store{
		st: st,
	}
}

func (s *Store) AddCommitment(commitment *Commitment) (err error) {
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

	key := commitmentKey(commitment.Bid.BlockNumber, commitment.Bid.BidAmount, commitment.Commitment)

	buf, err := msgpack.Marshal(commitment)
	if err != nil {
		return err
	}

	if err := writer.Put(key, buf); err != nil {
		return err
	}

	cIndexKey := cmtIndexKey(commitment.Commitment)
	cIndexValue := CommitmentIndexValue{
		BlockNumber: commitment.Bid.BlockNumber,
		BidAmount:   commitment.Bid.BidAmount,
	}

	cIndexValueBuf, err := msgpack.Marshal(cIndexValue)
	if err != nil {
		return err
	}

	return writer.Put(cIndexKey, cIndexValueBuf)
}

func (s *Store) SetStatus(
	blockNumber int64,
	bidAmt string,
	cDigest []byte,
	status CommitmentStatus,
	details string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := commitmentKey(blockNumber, bidAmt, cDigest)

	cmtBuf, err := s.st.Get(key)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}
	}

	commitment := new(Commitment)
	err = msgpack.Unmarshal(cmtBuf, commitment)
	if err != nil {
		return err
	}

	commitment.Status = status
	commitment.Details = details

	buf, err := msgpack.Marshal(commitment)
	if err != nil {
		return err
	}

	return s.st.Put(key, buf)
}

func (s *Store) GetCommitments(blockNum int64) ([]*Commitment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockCommitmentsKey := blockCommitmentPrefix(blockNum)
	commitments := make([]*Commitment, 0)

	err := s.st.WalkPrefix(blockCommitmentsKey, func(key string, value []byte) bool {
		commitment := new(Commitment)
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

func (s *Store) GetAllCommitments() ([]*Commitment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commitments := make([]*Commitment, 0)
	err := s.st.WalkPrefix(commitmentNS, func(key string, value []byte) bool {
		commitment := new(Commitment)
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

func (s *Store) SetCommitmentIndexByDigest(cDigest, cIndex [32]byte) (retErr error) {
	cmt, err := s.GetCommitmentByDigest(cDigest[:])
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}
		return err
	}

	cmt.CommitmentIndex = cIndex[:]
	if cmt.Status == CommitmentStatusPending {
		cmt.Status = CommitmentStatusStored
	}
	buf, err := msgpack.Marshal(cmt)
	if err != nil {
		return err
	}

	commitmentKey := commitmentKey(cmt.Bid.BlockNumber, cmt.Bid.BidAmount, cmt.Commitment)

	s.mu.Lock()
	defer s.mu.Unlock()

	var writer storage.Writer
	if w, ok := s.st.(storage.Batcher); ok {
		batch := w.Batch()
		writer = batch
		defer func() {
			switch {
			case retErr != nil:
				batch.Reset()
			case retErr == nil:
				err = batch.Write()
			}
		}()
	} else {
		writer = s.st
	}
	if err := writer.Put(commitmentKey, buf); err != nil {
		return err
	}

	indexToDigest := indexToDigestKey(cmt.CommitmentIndex)
	return writer.Put(indexToDigest, []byte(commitmentKey))
}

func (s *Store) UpdateSettlement(index []byte, isSlash bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := indexToDigestKey(index)
	commitmentKey, err := s.st.Get(key)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}
		return err
	}

	cmtBuf, err := s.st.Get(string(commitmentKey))
	if err != nil {
		return err
	}

	cmt := new(Commitment)
	err = msgpack.Unmarshal(cmtBuf, cmt)
	if err != nil {
		return err
	}

	if isSlash {
		cmt.Status = CommitmentStatusSlashed
	} else {
		cmt.Status = CommitmentStatusSettled
	}

	buf, err := msgpack.Marshal(cmt)
	if err != nil {
		return err
	}

	return s.st.Put(string(commitmentKey), buf)
}

func (s *Store) UpdatePayment(digest []byte, payment, refund string) error {
	cmt, err := s.GetCommitmentByDigest(digest)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cmt.Payment = payment
	cmt.Refund = refund

	buf, err := msgpack.Marshal(cmt)
	if err != nil {
		return err
	}

	commitmentKey := commitmentKey(cmt.Bid.BlockNumber, cmt.Bid.BidAmount, cmt.Commitment)

	return s.st.Put(string(commitmentKey), buf)
}

func (s *Store) GetCommitmentByDigest(digest []byte) (*Commitment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cIndexValueBuf, err := s.st.Get(cmtIndexKey(digest))
	if err != nil {
		return nil, err
	}

	var cIndexValue CommitmentIndexValue
	err = msgpack.Unmarshal(cIndexValueBuf, &cIndexValue)
	if err != nil {
		return nil, err
	}

	commitmentKey := commitmentKey(cIndexValue.BlockNumber, cIndexValue.BidAmount, digest)
	cmtBuf, err := s.st.Get(commitmentKey)
	if err != nil {
		return nil, err
	}

	cmt := new(Commitment)
	err = msgpack.Unmarshal(cmtBuf, cmt)
	if err != nil {
		return nil, err
	}

	return cmt, nil
}

func (s *Store) ClearCommitmentIndexes(uptoBlock int64) error {
	keys := make([]string, 0)
	s.mu.RLock()
	err := s.st.WalkPrefix(commitmentNS, func(key string, val []byte) bool {
		blockNum, err := parseBlockNumFromCommitmentKey(key)
		if err != nil {
			// If parsing fails, we might have corrupted data; discard this key
			keys = append(keys, key)
			return false
		}
		if blockNum < uptoBlock {
			var commitment Commitment
			err := msgpack.Unmarshal(val, &commitment)
			if err != nil {
				// If unmarshaling fails, we might have corrupted data; skip this key
				return false
			}
			keys = append(keys, key)
			keys = append(keys, cmtIndexKey(commitment.Commitment))
			if commitment.CommitmentIndex != nil {
				keys = append(keys, indexToDigestKey(commitment.CommitmentIndex))
			}
			return false
		}
		// DB is expected to be sorted by block number, so we can stop here
		// since all subsequent keys will also be greater than or equal to `uptoBlock`.
		return true
	})
	s.mu.RUnlock()
	if err != nil {
		return err
	}

	var (
		writer   storage.Writer
		writeErr error
	)
	if w, ok := s.st.(storage.Batcher); ok {
		batch := w.Batch()
		writer = batch
		defer func() {
			switch {
			case writeErr != nil:
				batch.Reset()
			case err == nil:
				writeErr = batch.Write()
			}
		}()
	} else {
		writer = s.st
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, key := range keys {
		writeErr = writer.Delete(key)
		if writeErr != nil {
			return writeErr
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

func (s *Store) ClearBlockNumber(blockNum int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Delete(blockWinnerKey(blockNum))
}
