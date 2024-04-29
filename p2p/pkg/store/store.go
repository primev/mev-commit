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
	commitmentsByBlockNumber      map[int64][]*EncryptedPreConfirmationWithDecrypted
	commitmentsByCommitmentHash   map[string]*EncryptedPreConfirmationWithDecrypted
	commitmentByBlockNumberMu     sync.RWMutex
	commitmentsByCommitmentHashMu sync.RWMutex
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
	return &Store{
		BlockStore: &BlockStore{
			data: make(map[string]uint64),
		},
		CommitmentsStore: &CommitmentsStore{
			commitmentsByBlockNumber:    make(map[int64][]*EncryptedPreConfirmationWithDecrypted),
			commitmentsByCommitmentHash: make(map[string]*EncryptedPreConfirmationWithDecrypted),
		},
		BidderBalancesStore: &BidderBalancesStore{
			balances:        make(map[string]*big.Int),
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

func (cs *CommitmentsStore) addCommitmentByBlockNumber(blockNum int64, commitment *EncryptedPreConfirmationWithDecrypted) {
	cs.commitmentByBlockNumberMu.Lock()
	defer cs.commitmentByBlockNumberMu.Unlock()

	cs.commitmentsByBlockNumber[blockNum] = append(cs.commitmentsByBlockNumber[blockNum], commitment)
}

func (cs *CommitmentsStore) addCommitmentByHash(hash string, commitment *EncryptedPreConfirmationWithDecrypted) {
	cs.commitmentsByCommitmentHashMu.Lock()
	defer cs.commitmentsByCommitmentHashMu.Unlock()

	cs.commitmentsByCommitmentHash[hash] = commitment
}

func (cs *CommitmentsStore) AddCommitment(commitment *EncryptedPreConfirmationWithDecrypted) {
	cs.addCommitmentByBlockNumber(commitment.Bid.BlockNumber, commitment)
	cs.addCommitmentByHash(common.Bytes2Hex(commitment.Commitment), commitment)
}

func (cs *CommitmentsStore) GetCommitmentsByBlockNumber(blockNum int64) ([]*EncryptedPreConfirmationWithDecrypted, error) {
	cs.commitmentByBlockNumberMu.RLock()
	defer cs.commitmentByBlockNumberMu.RUnlock()

	if commitments, exists := cs.commitmentsByBlockNumber[blockNum]; exists {
		return commitments, nil
	}
	return nil, nil
}

func (cs *CommitmentsStore) GetCommitmentByHash(hash string) (*EncryptedPreConfirmationWithDecrypted, error) {
	cs.commitmentsByCommitmentHashMu.RLock()
	defer cs.commitmentsByCommitmentHashMu.RUnlock()

	if commitment, exists := cs.commitmentsByCommitmentHash[hash]; exists {
		return commitment, nil
	}
	return nil, nil
}

func (cs *CommitmentsStore) DeleteCommitmentByBlockNumber(blockNum int64) error {
	cs.commitmentByBlockNumberMu.Lock()
	defer cs.commitmentByBlockNumberMu.Unlock()

	for _, v := range cs.commitmentsByBlockNumber[blockNum] {
		err := cs.deleteCommitmentByHash(common.Bytes2Hex(v.Commitment))
		if err != nil {
			return err
		}
	}
	delete(cs.commitmentsByBlockNumber, blockNum)
	return nil
}

func (cs *CommitmentsStore) deleteCommitmentByHash(hash string) error {
	cs.commitmentsByCommitmentHashMu.Lock()
	defer cs.commitmentsByCommitmentHashMu.Unlock()

	delete(cs.commitmentsByCommitmentHash, hash)
	return nil
}

func (cs *CommitmentsStore) SetCommitmentIndexByCommitmentDigest(cDigest, cIndex [32]byte) error {
	// when we will have db, this will be UPDATE query, instead of inmemory update
	commitment, err := cs.GetCommitmentByHash(common.Bytes2Hex(cDigest[:]))
	if err != nil {
		return fmt.Errorf("failed to get commitment by hash: %w", err)
	}
	if commitment == nil {
		// commitment could be not found in case this commitment is from another bidder/provider
		// so no need to return error in this case
		return nil
	}
	commitment.EncryptedPreConfirmation.CommitmentIndex = cIndex[:]
	return nil
}

type BidderBalancesStore struct {
	balances        map[string]*big.Int
	balancesByBlock *lru.Cache[string, *big.Int]
	mu              sync.RWMutex
}

func (bbs *BidderBalancesStore) SetBalance(bidder common.Address, windowNumber, depositedAmount *big.Int) error {
	bbs.mu.Lock()
	defer bbs.mu.Unlock()
	bssKey := getBBSKey(bidder, windowNumber)
	bbs.balances[bssKey] = depositedAmount
	return nil
}

func (bbs *BidderBalancesStore) GetBalance(bidder common.Address, windowNumber *big.Int) (*big.Int, error) {
	bbs.mu.RLock()
	defer bbs.mu.RUnlock()
	bssKey := getBBSKey(bidder, windowNumber)
	if balance, exists := bbs.balances[bssKey]; exists {
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
