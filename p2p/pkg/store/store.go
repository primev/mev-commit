package store

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	lru "github.com/hashicorp/golang-lru/v2"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
)

type Store struct {
	*BlockStore
	*CommitmentsStore
	*BidderBalancesStore
}

type BlockStore struct {
	data map[string]uint64
	mu   sync.RWMutex
}

type CommitmentsStore struct {
	commitmentsMu               sync.RWMutex
	commitmentsByBlockNumber    map[int64][]*EncryptedPreConfirmationWithDecrypted
	commitmentsByCommitmentHash map[string]*EncryptedPreConfirmationWithDecrypted
}

type EncryptedPreConfirmationWithDecrypted struct {
	*preconfpb.EncryptedPreConfirmation
	*preconfpb.PreConfirmation
}

func NewStore() (*Store, error) {
	balancesByBlockCache, err := lru.New[string, *big.Int](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create balancesByBlockCache: %w", err)
	}
	balancesCache, err := lru.New[string, *big.Int](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create balancesCache: %w", err)
	}
	return &Store{
		BlockStore: &BlockStore{
			data: make(map[string]uint64),
		},
		CommitmentsStore: &CommitmentsStore{
			commitmentsByBlockNumber:    make(map[int64][]*EncryptedPreConfirmationWithDecrypted),
			commitmentsByCommitmentHash: make(map[string]*EncryptedPreConfirmationWithDecrypted),
		},
		BidderBalancesStore: &BidderBalancesStore{
			balances:        balancesCache,
			balancesByBlock: balancesByBlockCache,
		},
	}, nil
}

func (bs *BlockStore) LastBlock() (uint64, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	if value, exists := bs.data["last_block"]; exists {
		return value, nil
	}
	return 0, nil
}

func (bs *BlockStore) SetLastBlock(blockNum uint64) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	bs.data["last_block"] = blockNum
	return nil
}

func (cs *CommitmentsStore) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) {
	cs.commitmentsMu.Lock()
	defer cs.commitmentsMu.Unlock()

	cs.commitmentsByBlockNumber[commitment.Bid.BlockNumber] = append(cs.commitmentsByBlockNumber[commitment.Bid.BlockNumber], commitment)
	cs.commitmentsByCommitmentHash[common.Bytes2Hex(commitment.EncryptedPreConfirmation.Commitment)] = commitment
}

func (cs *CommitmentsStore) GetCommitmentsByBlockNumber(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	cs.commitmentsMu.RLock()
	defer cs.commitmentsMu.RUnlock()

	if commitments, exists := cs.commitmentsByBlockNumber[blockNum]; exists {
		return commitments, nil
	}
	return nil, nil
}

func (cs *CommitmentsStore) DeleteCommitmentByBlockNumber(blockNum int64) error {
	cs.commitmentsMu.Lock()
	defer cs.commitmentsMu.Unlock()

	for _, v := range cs.commitmentsByBlockNumber[blockNum] {
		delete(cs.commitmentsByCommitmentHash, common.Bytes2Hex(v.EncryptedPreConfirmation.Commitment))
	}
	delete(cs.commitmentsByBlockNumber, blockNum)
	return nil
}

func (cs *CommitmentsStore) DeleteCommitmentByIndex(blockNum int64, index [32]byte) error {
	cs.commitmentsMu.Lock()
	defer cs.commitmentsMu.Unlock()

	for idx, v := range cs.commitmentsByBlockNumber[blockNum] {
		if common.Bytes2Hex(v.EncryptedPreConfirmation.CommitmentIndex) == common.Bytes2Hex(index[:]) {
			cs.commitmentsByBlockNumber[blockNum] = append(cs.commitmentsByBlockNumber[blockNum][:idx], cs.commitmentsByBlockNumber[blockNum][idx+1:]...)
			break
		}
	}

	for _, v := range cs.commitmentsByCommitmentHash {
		if common.Bytes2Hex(v.EncryptedPreConfirmation.CommitmentIndex) == common.Bytes2Hex(index[:]) {
			delete(cs.commitmentsByCommitmentHash, common.Bytes2Hex(v.EncryptedPreConfirmation.Commitment))
			break
		}
	}

	return nil
}

func (cs *CommitmentsStore) SetCommitmentIndexByCommitmentDigest(cDigest, cIndex [32]byte) error {
	cs.commitmentsMu.Lock()
	defer cs.commitmentsMu.Unlock()

	// when we will have db, this will be UPDATE query, instead of inmemory update
	if commitment, exists := cs.commitmentsByCommitmentHash[common.Bytes2Hex(cDigest[:])]; exists {
		commitment.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
		return nil
	}

	return nil
}

type BidderBalancesStore struct {
	balances        *lru.Cache[string, *big.Int]
	balancesByBlock *lru.Cache[string, *big.Int]
}

func (bbs *BidderBalancesStore) SetBalance(bidder common.Address, windowNumber, depositedAmount *big.Int) error {
	bssKey := getBBSKey(bidder, windowNumber)
	bbs.balances.Add(bssKey, depositedAmount)
	return nil
}

func (bbs *BidderBalancesStore) GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error) {
	bssKey := getBBSKey(bidder, windowNumber)
	if balance, exists := bbs.balances.Get(bssKey); exists {
		return balance, nil
	}
	return nil, nil
}

func getBBSKey(bidder common.Address, windowNumber *big.Int) string {
	return bidder.String() + windowNumber.String()
}

func (bbs *BidderBalancesStore) GetBalanceForBlock(bidder common.Address, blockNumber int64) (*big.Int, error) {
	key := getBBSforBlockKey(bidder, blockNumber)
	if balance, ok := bbs.balancesByBlock.Get(key); ok {
		return balance, nil
	}
	return nil, nil
}

func (bbs *BidderBalancesStore) SetBalanceForBlock(bidder common.Address, amount *big.Int, blockNumber int64) error {
	key := getBBSforBlockKey(bidder, blockNumber)
	bbs.balancesByBlock.Add(key, amount)
	return nil
}

func (bbs *BidderBalancesStore) RefundBalanceForBlock(bidder common.Address, amount *big.Int, blockNumber int64) error {
	key := getBBSforBlockKey(bidder, blockNumber)
	if currentBalance, ok := bbs.balancesByBlock.Get(key); ok {
		// If a balance exists, simply add the amount back
		updatedBalance := new(big.Int).Add(currentBalance, amount)
		bbs.balancesByBlock.Add(key, updatedBalance)
		return nil
	}

	// If no balance found (which should be unusual for a refund), initialize to the refund amount
	bbs.balancesByBlock.Add(key, amount)
	return nil
}

func getBBSforBlockKey(bidder common.Address, blockNumber int64) string {
	return bidder.String() + fmt.Sprint(blockNumber)
}
